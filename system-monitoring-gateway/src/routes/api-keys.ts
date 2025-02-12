import express, { Request, Response } from "express";
import { body } from "express-validator";
import { validateRequest } from "../middleware/validate-request";
import { requireAuth } from "../middleware/require-auth";
import { ApiKeyService } from "../services/api-key.service";
import { mongoDBService } from "../services/mongodb.service";
import { BadRequestError, NotAuthorizedError } from "../errors";

const router = express.Router();

// Validation middleware
const validateApiKeyCreation = [
  body("description")
    .isString()
    .trim()
    .notEmpty()
    .withMessage("Description is required")
    .isLength({ max: 200 })
    .withMessage("Description must be less than 200 characters"),
  body("permissions")
    .optional()
    .isArray()
    .withMessage("Permissions must be an array")
    .custom((value: string[]) => {
      const validPermissions = ["read", "write", "admin"];
      return value.every((permission) => validPermissions.includes(permission));
    })
    .withMessage("Invalid permissions"),
  body("expiresInDays")
    .optional()
    .isInt({ min: 1, max: 365 })
    .withMessage("Expiration days must be between 1 and 365"),
];

// Create API key for user
router.post(
  "/users/:userId/api-keys",
  requireAuth,
  validateApiKeyCreation,
  validateRequest,
  async (req: Request, res: Response) => {
    const { userId } = req.params;
    const { description, permissions, expiresInDays } = req.body;

    // Ensure user can only create keys for themselves unless they're an admin
    if (userId !== req.currentUser!.id && req.currentUser!.role !== "admin") {
      throw new NotAuthorizedError();
    }

    try {
      // Validate user ID format and existence
      const targetUser = await mongoDBService.getUserById(userId);
      if (!targetUser) {
        throw new BadRequestError("User not found or invalid user ID format");
      }

      // Only check tenant match if both users have tenants
      if (req.currentUser!.tenantId && targetUser.tenantId) {
        if (targetUser.tenantId.toString() !== req.currentUser!.tenantId) {
          throw new NotAuthorizedError("Users must belong to the same tenant");
        }
      }
    } catch (error) {
      console.error("[API Keys] Error validating user:", error);
      throw new BadRequestError("Failed to validate user");
    }

    const apiKey = await ApiKeyService.generateApiKey(
      userId,
      description,
      req.currentUser!.tenantId,
      permissions,
      expiresInDays
    );

    res.status(201).json(apiKey);
  }
);

// List user's API keys
router.get(
  "/users/:userId/api-keys",
  requireAuth,
  async (req: Request, res: Response) => {
    const { userId } = req.params;

    // Ensure user can only view their own keys unless they're an admin
    if (userId !== req.currentUser!.id && req.currentUser!.role !== "admin") {
      throw new NotAuthorizedError();
    }

    const apiKeys = await ApiKeyService.getUserApiKeys(userId);
    res.json(apiKeys);
  }
);

// Revoke API key
router.put(
  "/users/:userId/api-keys/:keyId/revoke",
  requireAuth,
  async (req: Request, res: Response) => {
    const { userId, keyId } = req.params;

    // Ensure user can only revoke their own keys unless they're an admin
    if (userId !== req.currentUser!.id && req.currentUser!.role !== "admin") {
      throw new NotAuthorizedError();
    }

    // Validate user owns this key
    const hasAccess = await ApiKeyService.validateUserAccess(userId, keyId);
    if (!hasAccess) {
      throw new NotAuthorizedError();
    }

    const success = await ApiKeyService.revokeApiKey(keyId);
    if (!success) {
      throw new BadRequestError("Failed to revoke API key");
    }

    res.status(200).json({ message: "API key revoked successfully" });
  }
);

// Delete API key
router.delete(
  "/users/:userId/api-keys/:keyId",
  requireAuth,
  async (req: Request, res: Response) => {
    const { userId, keyId } = req.params;

    // Ensure user can only delete their own keys unless they're an admin
    if (userId !== req.currentUser!.id && req.currentUser!.role !== "admin") {
      throw new NotAuthorizedError();
    }

    // Validate user owns this key
    const hasAccess = await ApiKeyService.validateUserAccess(userId, keyId);
    if (!hasAccess) {
      throw new NotAuthorizedError();
    }

    const success = await ApiKeyService.revokeApiKey(keyId);
    if (!success) {
      throw new BadRequestError("Failed to delete API key");
    }

    res.status(204).send();
  }
);

export { router as apiKeysRouter };
