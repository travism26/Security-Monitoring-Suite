import { Request, Response, NextFunction } from "express";
import { CustomError } from "../errors/custom-error";

interface ErrorResponse {
  errors: { message: string; field?: string }[];
  tenantId?: string;
  timestamp: string;
  requestId: string;
}

export const errorHandler = (
  err: Error,
  req: Request,
  res: Response,
  next: NextFunction
) => {
  // Generate unique request ID for tracking
  const requestId = Math.random().toString(36).substring(2, 15);

  // Get tenant context if available
  const tenantId = req.tenantId || req.currentUser?.tenantId;

  if (err instanceof CustomError) {
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
