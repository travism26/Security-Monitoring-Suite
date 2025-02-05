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

  // If no tenant headers are provided, log and proceed
  if (!headerTenantId || !headerEnvironment) {
    console.log("No tenant headers provided - proceeding in non-tenant mode");
    return next();
  }

  // Only validate tenant consistency if both user context and header tenant IDs exist
  if (tenantId && headerTenantId && headerTenantId !== tenantId) {
    console.log("Tenant ID mismatch", tenantId, headerTenantId);
    throw new ForbiddenError(
      "Tenant ID mismatch between user context and headers"
    );
  }

  // For POST requests, validate tenant ID and environment in body matches headers
  // Only validate POST request tenant consistency if tenant information is provided
  if (req.method === "POST" && req.body?.data?.metadata) {
    console.log("Checking tenant consistency in POST request");
    const payloadTenantId = req.body.data.metadata.tenant_id;
    const payloadEnvironment = req.body.data.metadata.environment;

    // Skip validation if payload doesn't include tenant information
    if (!payloadTenantId || !payloadEnvironment) {
      console.log("No tenant information in payload - proceeding");
      return next();
    }

    // Only validate if both header and payload tenant IDs exist
    if (
      headerTenantId &&
      payloadTenantId &&
      headerTenantId !== payloadTenantId
    ) {
      console.log("Tenant ID mismatch", headerTenantId, payloadTenantId);
      throw new BadRequestError(
        "Tenant ID mismatch between headers and payload"
      );
    }

    // Only validate if both header and payload environments exist
    if (
      headerEnvironment &&
      payloadEnvironment &&
      headerEnvironment !== payloadEnvironment
    ) {
      console.log(
        `Environment mismatch`,
        headerEnvironment,
        payloadEnvironment
      );
      throw new BadRequestError(
        "Environment mismatch between headers and payload"
      );
    }
  }

  console.log("Tenant consistency validated");
  next();
};
