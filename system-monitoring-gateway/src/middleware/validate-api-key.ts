import { Request, Response, NextFunction } from "express";
import { NotAuthorizedError } from "../errors";
import { ApiKeyService } from "../services/api-key.service";

export const validateApiKey = async (
  req: Request,
  res: Response,
  next: NextFunction
) => {
  // Case-insensitive header check
  const apiKey = req.header("X-API-Key") || req.header("X-Api-Key");

  if (!apiKey) {
    throw new NotAuthorizedError("API key is required");
  }

  try {
    const validApiKey = await ApiKeyService.validateApiKey(apiKey);

    if (!validApiKey) {
      throw new NotAuthorizedError("Invalid API key");
    }

    // Add API key details to request for downstream use
    req.currentUser = {
      id: validApiKey.tenantId.toString(), // Use tenantId as the identifier for API key auth
      tenantId: validApiKey.tenantId.toString(),
      role: "api",
      email: "", // API keys don't have associated emails
    };

    next();
  } catch (error) {
    throw new NotAuthorizedError("Invalid API key");
  }
};
