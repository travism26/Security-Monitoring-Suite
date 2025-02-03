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
    console.log("Headers", req.headers);
    throw new BadRequestError("Missing required tenant headers");
  }

  if (tenantId && headerTenantId !== tenantId) {
    console.log("Tenant ID mismatch", tenantId, headerTenantId);
    throw new ForbiddenError(
      "Tenant ID mismatch between user context and headers"
    );
  }

  // For POST requests, validate tenant ID and environment in body matches headers
  if (req.method === "POST" && req.body?.data?.metadata) {
    console.log("Validating tenant consistency in POST request");
    const payloadTenantId = req.body.data.metadata.tenant_id;
    const payloadEnvironment = req.body.data.metadata.environment;

    if (headerTenantId !== payloadTenantId) {
      console.log("Tenant ID mismatch", headerTenantId, payloadTenantId);
      throw new BadRequestError(
        "Tenant ID mismatch between headers and payload"
      );
    }

    if (headerEnvironment !== payloadEnvironment) {
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
