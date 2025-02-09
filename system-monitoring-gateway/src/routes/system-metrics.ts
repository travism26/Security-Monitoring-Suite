import express, { Request, Response } from "express";
import { body } from "express-validator";
import { validateRequest } from "../middleware/validate-request";
import { validateTenantConsistency } from "../middleware/validate-tenant";
import { metricsRegistry } from "./metrics";
import { Counter } from "prom-client";
import {
  SystemMetrics,
  SystemMetricsPayload,
  SystemMetricsData,
} from "../payload/system-metrics";
import { kafkaWrapper } from "../kafka/kafka-wrapper";

const router = express.Router();

// Validation middleware
const validateMetrics = [
  body().isObject().withMessage("Request body must be an object"),
  body("data").isObject().withMessage("Data field is required"),
  body("data.metrics").notEmpty().withMessage("Metrics data is required"),
  body("data.timestamp").isISO8601().withMessage("Invalid timestamp format"),
  body("data.tenant_id").optional(),
  body("data.tenant_metadata").optional(),
  body("data.api_key")
    .optional()
    .isString()
    .withMessage("API key must be a string"),
  body("data.host").isObject().withMessage("Host information is required"),
  body("data.host.os").isString().withMessage("Host OS is required"),
  body("data.host.hostname")
    .isString()
    .withMessage("Host hostname is required"),
  body("data.host.cpu_cores")
    .isInt()
    .withMessage("Host CPU cores must be a number"),
  body("data.processes")
    .isObject()
    .withMessage("Process information is required"),
  body("data.processes.total_count")
    .isInt()
    .withMessage("Process total count must be a number"),
  body("data.processes.list")
    .isArray()
    .withMessage("Process list must be an array"),
  body("data.metadata").isObject().withMessage("Metadata is required"),
  body("data.metadata.collection_duration")
    .isString()
    .withMessage("Collection duration is required"),
  body("data.metadata.collector_count")
    .isInt()
    .withMessage("Collector count must be a number"),
];

// "/api/v1/system",
router.post(
  "/system-metrics/ingest",
  validateTenantConsistency,
  validateMetrics,
  validateRequest,
  async (req: Request<{}, {}, SystemMetricsData>, res: Response) => {
    console.log("[DEBUG] Received metrics request", {
      headers: req.headers,
      contentLength: req.get("content-length"),
      tenantId: req.get("x-tenant-id"),
      timestamp: new Date().toISOString(),
      apiKey: req.get("x-api-key"),
    });

    try {
      console.log("[DEBUG] Starting metrics processing");

      // Extract data from wrapper and add API key
      const { data } = req.body;
      if (!data) {
        return res.status(400).json({
          errors: [{ message: "Missing data field in payload" }],
        });
      }

      // Add API key to the data payload
      // We can safely assert this is a string since validateApiKey middleware ensures it exists
      const apiKey = req.get("x-api-key") as string;
      data.api_key = apiKey;

      // Log only process summary instead of detailed list
      const processCount = data.processes?.list?.length || 0;
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
        console.error("[ERROR] Kafka connection not established");
        return res.status(503).json({
          errors: [
            {
              message: "Metrics service temporarily unavailable",
              details: "Kafka connection not established",
            },
          ],
        });
      }
      console.log("[DEBUG] Kafka connection verified");

      // Validate tenant ID consistency if header is present
      const headerTenantId = req.get("x-tenant-id");
      if (headerTenantId && data.tenant_id) {
        if (headerTenantId !== data.tenant_id) {
          const errorProducer = kafkaWrapper.getProducer(
            "system-metrics-errors"
          );
          await errorProducer.publish({
            error: "Tenant ID mismatch",
            header_tenant_id: headerTenantId,
            payload_tenant_id: data.tenant_id,
            timestamp: new Date().toISOString(),
          });
          return res.status(400).json({
            errors: [
              { message: "Tenant ID mismatch between headers and payload" },
            ],
          });
        }
      }

      // Validate required fields
      if (!data.metrics) {
        const errorProducer = kafkaWrapper.getProducer("system-metrics-errors");
        await errorProducer.publish({
          error: "Validation failed",
          original_payload: req.body,
          tenant_id: data.tenant_id || "no-tenant",
          timestamp: new Date().toISOString(),
        });
        return res.status(400).json({
          errors: [{ message: "Metrics data is required" }],
        });
      }

      try {
        console.log("[DEBUG] Attempting to publish metrics to Kafka");
        const kafkaProducer = kafkaWrapper.getProducer("system-metrics");

        await kafkaProducer.publish(data);
        console.log("[DEBUG] Successfully published metrics to Kafka");

        return res.status(202).json({
          status: "accepted",
          timestamp: new Date().toISOString(),
        });
      } catch (error) {
        console.error("[ERROR] Failed to produce metrics to Kafka:", error);

        // Send to DLQ
        const dlqProducer = kafkaWrapper.getProducer("system-metrics-dlq");
        await dlqProducer.publish({
          error: "Message processing failed",
          original_message: req.body,
          tenant_id: data.tenant_id || "no-tenant",
          timestamp: new Date().toISOString(),
        });

        return res.status(500).json({
          errors: [{ message: "Failed to process metrics data" }],
        });
      }
    } catch (error) {
      console.error("[ERROR] Error processing metrics:", error);
      console.error(
        "[ERROR] Stack trace:",
        error instanceof Error ? error.stack : "No stack trace available"
      );
      return res.status(500).json({
        errors: [{ message: "Internal server error while processing metrics" }],
      });
    }
  }
);

export { router as systemMetricsRouter };
