import { Request, Response, NextFunction } from "express";
import { NotAuthorizedError } from "../errors";
import { ApiKeyService, ApiKey } from "../services/api-key.service";

export const validateApiKey = async (
  req: Request,
  res: Response,
  next: NextFunction
) => {
  // Extract API key using case-insensitive header check
  const apiKey = req.get("x-api-key");

  if (!apiKey) {
    throw new NotAuthorizedError("API key is required");
  }

  try {
    // Validate API key and get associated user data
    const validApiKey = await ApiKeyService.validateApiKey(apiKey);
    if (!validApiKey) {
      throw new NotAuthorizedError("Invalid API key");
    }

    // Attach API key and user data to request for downstream use
    req.apiKey = apiKey;
    req.userId = validApiKey.userId;
    req.tenantId = validApiKey.tenantId;

    next();
  } catch (error) {
    throw new NotAuthorizedError("Invalid API key");
  }
};

// Add type definition for Request
declare global {
  namespace Express {
    interface Request {
      apiKey?: string;
      userId?: string;
      tenantId?: string;
    }
  }
}
