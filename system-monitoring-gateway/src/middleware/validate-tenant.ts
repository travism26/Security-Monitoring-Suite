import { Request, Response, NextFunction } from "express";
import { ForbiddenError } from "../errors";

export const validateTenantConsistency = async (
  req: Request,
  res: Response,
  next: NextFunction
) => {
  const tenantId = req.currentUser?.tenantId;
  const requestTenantId = req.header("X-Tenant-ID");

  if (!tenantId) {
    throw new ForbiddenError("Tenant ID not found in user context");
  }

  if (requestTenantId && requestTenantId !== tenantId) {
    throw new ForbiddenError("Tenant ID mismatch");
  }

  next();
};
