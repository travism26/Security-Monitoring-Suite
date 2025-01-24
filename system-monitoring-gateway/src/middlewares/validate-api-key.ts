import { Request, Response, NextFunction } from "express";
import { NotAuthorizedError } from "../errors/not-authorized-error";

import { ApiKeyService } from "../services/api-key.service";

// Augment the Express Request type to include tenant information
declare global {
  namespace Express {
    interface Request {
      tenantId?: string;
    }
  }
}

export const validateApiKey = async (
  req: Request,
  res: Response,
  next: NextFunction
) => {
  const apiKey = req.headers["x-api-key"] as string;

  if (!apiKey) {
    throw new NotAuthorizedError();
  }

  const validatedKey = await ApiKeyService.validateApiKey(apiKey);

  if (!validatedKey) {
    throw new NotAuthorizedError();
  }

  // Add tenant information to the request
  req.tenantId = validatedKey.tenantId;

  // Add API key metadata to response headers
  res.setHeader("X-API-Key-Created", validatedKey.createdAt.toISOString());
  if (validatedKey.expiresAt) {
    res.setHeader("X-API-Key-Expires", validatedKey.expiresAt.toISOString());
  }

  next();
};
