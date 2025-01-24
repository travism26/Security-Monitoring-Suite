import { Request, Response, NextFunction } from "express";
import { NotAuthorizedError } from "../errors/not-authorized-error";

// Interface for API key validation response
interface ApiKeyValidation {
  valid: boolean;
  tenantId?: string;
}

// Augment the Express Request type to include tenant information
declare global {
  namespace Express {
    interface Request {
      tenantId?: string;
    }
  }
}

// This would typically interact with a database or external service
// For now, we'll use a simple in-memory store for demonstration
const validateApiKeyWithStore = async (
  apiKey: string
): Promise<ApiKeyValidation> => {
  // TODO: Implement actual API key validation against a database
  // This is a placeholder implementation
  if (!apiKey) {
    return { valid: false };
  }

  // In a real implementation, we would:
  // 1. Query the database for the API key
  // 2. Verify the key hasn't expired
  // 3. Check if the key is active
  // 4. Return the associated tenant information

  // Placeholder validation (replace with actual database lookup)
  const isValid = apiKey.startsWith("sms_") && apiKey.length >= 32;
  return {
    valid: isValid,
    tenantId: isValid ? "demo-tenant-id" : undefined,
  };
};

export const validateApiKey = async (
  req: Request,
  res: Response,
  next: NextFunction
) => {
  const apiKey = req.headers["x-api-key"] as string;

  if (!apiKey) {
    throw new NotAuthorizedError();
  }

  const validation = await validateApiKeyWithStore(apiKey);

  if (!validation.valid || !validation.tenantId) {
    throw new NotAuthorizedError();
  }

  // Add tenant information to the request
  req.tenantId = validation.tenantId;

  next();
};
