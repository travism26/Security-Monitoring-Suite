import { Request, Response, NextFunction } from "express";

export const validateTenantConsistency = (
  req: Request,
  res: Response,
  next: NextFunction
) => {
  const headerTenantId = req.headers["x-tenant-id"];
  const headerEnvironment = req.headers["x-tenant-environment"];

  if (!headerTenantId || !headerEnvironment) {
    return res.status(400).json({
      errors: [{ message: "Missing tenant headers" }],
    });
  }

  // For POST requests, validate tenant ID in body matches header
  if (req.method === "POST" && req.body?.data?.metadata) {
    const payloadTenantId = req.body.data.metadata.tenant_id;
    const payloadEnvironment = req.body.data.metadata.environment;

    if (headerTenantId !== payloadTenantId) {
      return res.status(400).json({
        errors: [{ message: "Tenant ID mismatch between headers and payload" }],
      });
    }

    if (headerEnvironment !== payloadEnvironment) {
      return res.status(400).json({
        errors: [
          { message: "Environment mismatch between headers and payload" },
        ],
      });
    }
  }

  next();
};
