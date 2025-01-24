import express, { Request, Response } from "express";
import { validateJWT, requireAuth } from "../middlewares/require-auth";
import { ApiKeyService } from "../services/api-key.service";

const router = express.Router();

// All routes require JWT authentication
router.use(validateJWT);
router.use(requireAuth);

// Generate new API key
router.post("/api/v1/api-keys", async (req: Request, res: Response) => {
  const { expiresInDays } = req.body;
  const tenantId = req.currentUser!.tenantId;

  const apiKey = ApiKeyService.generateApiKey(tenantId, expiresInDays);

  res.status(201).send(apiKey);
});

// List all API keys for tenant
router.get("/api/v1/api-keys", async (req: Request, res: Response) => {
  const tenantId = req.currentUser!.tenantId;

  const apiKeys = await ApiKeyService.listApiKeys(tenantId);

  res.send(apiKeys);
});

// Revoke an API key
router.delete("/api/v1/api-keys/:key", async (req: Request, res: Response) => {
  const { key } = req.params;
  const tenantId = req.currentUser!.tenantId;

  // Validate the key belongs to the tenant
  const apiKey = await ApiKeyService.validateApiKey(key);
  if (!apiKey || apiKey.tenantId !== tenantId) {
    res.status(404).send({ message: "API key not found" });
    return;
  }

  const success = await ApiKeyService.revokeApiKey(key);

  if (success) {
    res.status(204).send();
  } else {
    res.status(404).send({ message: "API key not found" });
  }
});

// Rotate API key
router.post(
  "/api/v1/api-keys/:key/rotate",
  async (req: Request, res: Response) => {
    const { key } = req.params;
    const { expiresInDays } = req.body;
    const tenantId = req.currentUser!.tenantId;

    // Validate the key belongs to the tenant
    const apiKey = await ApiKeyService.validateApiKey(key);
    if (!apiKey || apiKey.tenantId !== tenantId) {
      res.status(404).send({ message: "API key not found" });
      return;
    }

    const newApiKey = await ApiKeyService.rotateApiKey(key, expiresInDays);

    if (newApiKey) {
      res.send(newApiKey);
    } else {
      res.status(404).send({ message: "API key not found" });
    }
  }
);

export { router as apiKeysRouter };
