import { Request, Response, NextFunction } from "express";
import { NotAuthorizedError } from "../errors";
import { ApiKeyService } from "../services/api-key.service";

export const validateApiKey = async (
  req: Request,
  res: Response,
  next: NextFunction
) => {
  const apiKey = req.header("X-API-Key");

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
      id: validApiKey.userId.toString(),
      tenantId: validApiKey.tenantId.toString(),
      role: "api",
      email: "",
    };

    next();
  } catch (error) {
    throw new NotAuthorizedError("Invalid API key");
  }
};
