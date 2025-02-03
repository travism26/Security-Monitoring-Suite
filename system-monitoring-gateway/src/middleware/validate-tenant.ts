import { Request, Response, NextFunction } from "express";
import { ForbiddenError, BadRequestError } from "../errors";

export const validateTenantConsistency = async (
  req: Request,
  res: Response,
  next: NextFunction
) => {
  const tenantId = req.currentUser?.tenantId;
  const headerTenantId = req.headers["x-tenant-id"] as string;
  const headerEnvironment = req.headers["x-tenant-environment"] as string;

  if (!headerTenantId || !headerEnvironment) {
    throw new BadRequestError("Missing required tenant headers");
  }

  if (tenantId && headerTenantId !== tenantId) {
    throw new ForbiddenError(
      "Tenant ID mismatch between user context and headers"
    );
  }

  // For POST requests, validate tenant ID and environment in body matches headers
  if (req.method === "POST" && req.body?.data?.metadata) {
    const payloadTenantId = req.body.data.metadata.tenant_id;
    const payloadEnvironment = req.body.data.metadata.environment;

    if (headerTenantId !== payloadTenantId) {
      throw new BadRequestError(
        "Tenant ID mismatch between headers and payload"
      );
    }

    if (headerEnvironment !== payloadEnvironment) {
      throw new BadRequestError(
        "Environment mismatch between headers and payload"
      );
    }
  }

  next();
};
