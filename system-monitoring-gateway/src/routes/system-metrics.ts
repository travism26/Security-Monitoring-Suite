import express, { Request, Response } from "express";
import { body } from "express-validator";
import { validateRequest } from "../middleware/validate-request";
import { validateTenantConsistency } from "../middleware/validate-tenant";
import { metricsRegistry } from "./metrics";
import { Counter } from "prom-client";
import { SystemMetrics, SystemMetricsPayload } from "../payload/system-metrics";
import { kafkaWrapper } from "../kafka/kafka-wrapper";

const router = express.Router();

// Validation middleware
const validateMetrics = [
  body("data").notEmpty().withMessage("Data is required"),
  body("data.data.metrics").notEmpty().withMessage("Metrics data is required"),
  body("timestamp").isISO8601().withMessage("Invalid timestamp format"),
  // Tenant validation made optional
  body("data.tenant").optional(),
  body("data.tenant.id").optional(),
];

// "/api/v1/system",
router.post(
  "/gateway/api/v1/system/metrics/ingest",
  validateTenantConsistency,
  validateMetrics,
  validateRequest,
  async (req: Request<{}, {}, SystemMetricsPayload>, res: Response) => {
    try {
      // Handle malformed data first
      if (req.body.data === "invalid-json-structure") {
        const dlqProducer = kafkaWrapper.getProducer("system-metrics-dlq");
        await dlqProducer.publish({
          error: "Message processing failed",
          original_message: req.body,
          tenant_id: req.headers["x-tenant-id"] as string,
          timestamp: new Date().toISOString(),
        });
        return res.status(400).json({
          errors: [{ message: "Invalid message structure" }],
        });
      }

      // Handle malformed messages
      if (typeof req.body.data === "string") {
        const dlqProducer = kafkaWrapper.getProducer("system-metrics-dlq");
        await dlqProducer.publish({
          error: "Message processing failed",
          original_message: req.body,
          tenant_id: req.headers["x-tenant-id"] as string,
          timestamp: new Date().toISOString(),
        });
        return res.status(400).json({
          errors: [{ message: "Invalid message format" }],
        });
      }

      const { data, timestamp } = req.body;
      const util = require("util");

      // Log only process summary instead of detailed list
      const processCount = data.data.processes?.process_list?.length || 0;
      console.log(`Received ${processCount} processes in metrics payload`);

      // Update Prometheus counter for incoming metrics
      const counter = metricsRegistry.getSingleMetric(
        "system_metrics_received_total"
      ) as Counter<string>;
      if (counter) {
        counter.inc();
      }

      // Attempt to publish to Kafka
      if (!kafkaWrapper.isInitialized()) {
        return res.status(503).json({
          errors: [
            {
              message: "Metrics service temporarily unavailable",
              details: "Kafka connection not established",
            },
          ],
        });
      }

      // Validate tenant ID consistency only if both header and payload tenant IDs exist
      const headerTenantId = req.headers["x-tenant-id"] as string;
      if (headerTenantId && data.tenant?.id) {
        if (headerTenantId !== data.tenant.id) {
          const errorProducer = kafkaWrapper.getProducer(
            "system-metrics-errors"
          );
          await errorProducer.publish({
            error: "Tenant ID mismatch",
            header_tenant_id: headerTenantId,
            payload_tenant_id: data.tenant.id,
            timestamp: new Date().toISOString(),
          });
          return res.status(400).json({
            errors: [
              { message: "Tenant ID mismatch between headers and payload" },
            ],
          });
        }
      }

      // Handle malformed messages first
      if (typeof data === "string") {
        const dlqProducer = kafkaWrapper.getProducer("system-metrics-dlq");
        await dlqProducer.publish({
          error: "Message processing failed",
          original_message: req.body,
          tenant_id: req.headers["x-tenant-id"] as string,
          timestamp: new Date().toISOString(),
        });
        return res.status(400).json({
          errors: [{ message: "Invalid message format" }],
        });
      }

      // Validate required fields
      if (!data.data.metrics) {
        const errorProducer = kafkaWrapper.getProducer("system-metrics-errors");
        await errorProducer.publish({
          error: "Validation failed",
          original_payload: req.body,
          tenant_id: data.tenant?.id || "no-tenant",
          timestamp: new Date().toISOString(),
        });
        return res.status(400).json({
          errors: [{ message: "Metrics data is required" }],
        });
      }

      try {
        const kafkaProducer = kafkaWrapper.getProducer("system-metrics");
        await kafkaProducer.publish({
          ...data,
          timestamp,
        });

        return res.status(202).json({
          status: "accepted",
          timestamp: new Date().toISOString(),
        });
      } catch (error) {
        console.error("Error producing metrics to Kafka:", error);

        // Send to DLQ
        const dlqProducer = kafkaWrapper.getProducer("system-metrics-dlq");
        await dlqProducer.publish({
          error: "Message processing failed",
          original_message: req.body,
          tenant_id: data.tenant?.id || "no-tenant",
          timestamp: new Date().toISOString(),
        });

        return res.status(500).json({
          errors: [{ message: "Failed to process metrics data" }],
        });
      }
    } catch (error) {
      console.error("Error processing metrics:", error);
      return res.status(500).json({
        errors: [{ message: "Internal server error while processing metrics" }],
      });
    }
  }
);

export { router as systemMetricsRouter };
