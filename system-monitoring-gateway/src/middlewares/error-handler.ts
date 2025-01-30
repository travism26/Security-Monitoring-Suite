import { Request, Response, NextFunction } from "express";
import { CustomError } from "../errors/custom-error";
import { Topics } from "../kafka/topics";
import { kafkaWrapper } from "../kafka/kafka-wrapper";

interface ErrorResponse {
  errors: { message: string; field?: string }[];
  tenantId?: string;
  timestamp: string;
  requestId: string;
}

export const errorHandler = async (
  err: Error,
  req: Request,
  res: Response,
  next: NextFunction
) => {
  const publishError = async (topic: Topics, errorData: any): Promise<void> => {
    try {
      const producer = kafkaWrapper.getProducer(topic);
      await producer.publish(errorData);
    } catch (kafkaError) {
      console.error("Failed to publish error to Kafka:", kafkaError);
    }
  };

  // Generate unique request ID for tracking
  const requestId = Math.random().toString(36).substring(2, 15);

  // Get tenant context if available
  const tenantId = req.tenantId || req.currentUser?.tenantId;

  if (err instanceof CustomError) {
    const errorData = {
      error: err.message,
      original_payload: req.body,
      tenant_id: tenantId || "unknown",
      timestamp: new Date().toISOString(),
      status_code: err.statusCode,
    };

    if (err.statusCode >= 400 && err.statusCode < 500) {
      await publishError(Topics.SystemMetricsErrors, errorData);
    } else if (err.statusCode >= 500) {
      await publishError(Topics.SystemMetricsDLQ, errorData);
    }

    const response: ErrorResponse = {
      errors: err.serializeErrors(),
      timestamp: new Date().toISOString(),
      requestId,
    };

    if (tenantId) {
      response.tenantId = tenantId;
    }

    // Add tenant-specific headers
    if (tenantId) {
      res.setHeader("X-Tenant-ID", tenantId);
    }
    res.setHeader("X-Request-ID", requestId);

    return res.status(err.statusCode).send(response);
  }

  // Log unexpected errors with tenant context
  console.error("Unexpected error:", {
    error: err,
    tenantId,
    requestId,
    path: req.path,
    method: req.method,
  });

  const response: ErrorResponse = {
    errors: [{ message: "Something went wrong" }],
    timestamp: new Date().toISOString(),
    requestId,
  };

  if (tenantId) {
    response.tenantId = tenantId;
    res.setHeader("X-Tenant-ID", tenantId);
  }
  res.setHeader("X-Request-ID", requestId);

  res.status(500).send(response);
};
