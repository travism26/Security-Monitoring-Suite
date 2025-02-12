import { Request, Response, NextFunction } from "express";
import { NotAuthorizedError } from "../errors";
import { ApiKeyService } from "../services/api-key.service";

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
    // Only validate that the API key exists and is active
    const validApiKey = await ApiKeyService.validateApiKey(apiKey);
    if (!validApiKey) {
      throw new NotAuthorizedError("Invalid API key");
    }

    // Attach API key to request for downstream use (e.g., Kafka payload)
    req.apiKey = apiKey;

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
    }
  }
}
