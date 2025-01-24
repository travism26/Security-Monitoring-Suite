import { Request, Response, NextFunction } from "express";
import { NotAuthorizedError } from "../errors/not-authorized-error";

interface TenantContext {
  tenantId: string;
  rateLimit: {
    requestsPerMinute: number;
    currentRequests: number;
    lastReset: Date;
  };
}

// In-memory store for tenant contexts (should be replaced with Redis/database in production)
const tenantContexts: Map<string, TenantContext> = new Map();

// Default rate limit per minute per tenant
const DEFAULT_RATE_LIMIT = 1000;

export const validateTenant = (
  req: Request,
  res: Response,
  next: NextFunction
) => {
  const tenantId = req.tenantId || req.currentUser?.tenantId;

  if (!tenantId) {
    throw new NotAuthorizedError();
  }

  // Get or create tenant context
  let tenantContext = tenantContexts.get(tenantId);
  if (!tenantContext) {
    tenantContext = {
      tenantId,
      rateLimit: {
        requestsPerMinute: DEFAULT_RATE_LIMIT,
        currentRequests: 0,
        lastReset: new Date(),
      },
    };
    tenantContexts.set(tenantId, tenantContext);
  }

  // Check and update rate limiting
  const now = new Date();
  const timeDiff = now.getTime() - tenantContext.rateLimit.lastReset.getTime();

  // Reset counter if a minute has passed
  if (timeDiff >= 60000) {
    tenantContext.rateLimit.currentRequests = 0;
    tenantContext.rateLimit.lastReset = now;
  }

  // Check rate limit
  if (
    tenantContext.rateLimit.currentRequests >=
    tenantContext.rateLimit.requestsPerMinute
  ) {
    res.status(429).json({
      errors: [
        {
          message: "Rate limit exceeded. Please try again later.",
        },
      ],
    });
    return;
  }

  // Increment request counter
  tenantContext.rateLimit.currentRequests++;
  tenantContexts.set(tenantId, tenantContext);

  // Add rate limit headers
  res.setHeader(
    "X-RateLimit-Limit",
    tenantContext.rateLimit.requestsPerMinute.toString()
  );
  res.setHeader(
    "X-RateLimit-Remaining",
    (
      tenantContext.rateLimit.requestsPerMinute -
      tenantContext.rateLimit.currentRequests
    ).toString()
  );
  res.setHeader(
    "X-RateLimit-Reset",
    new Date(tenantContext.rateLimit.lastReset.getTime() + 60000).toISOString()
  );

  next();
};
