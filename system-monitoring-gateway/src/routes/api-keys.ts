import express, { Request, Response } from "express";
import { validateJWT, requireAuth } from "../middleware/require-auth";
import { validateTenantConsistency } from "../middleware/validate-tenant";
import { ApiKeyService } from "../services/api-key.service";

const router = express.Router();

// All routes require JWT authentication
router.use(validateJWT);
router.use(requireAuth);
// Tenant validation is optional in design phase
router.use(validateTenantConsistency);
// Note: Tenant validation is currently optional during design phase

// Mount all routes under /api/v1/keys
const apiKeysRouter = express.Router();
apiKeysRouter.use("/gateway/api/v1/keys", router);

// Generate new API key
router.post("/", async (req: Request, res: Response) => {
  const { expiresInDays } = req.body;
  const tenantId = req.currentUser!.tenantId;

  const apiKey = ApiKeyService.generateApiKey(tenantId, expiresInDays);

  res.status(201).send(apiKey);
});

// List all API keys for tenant
router.get("/", async (req: Request, res: Response) => {
  const tenantId = req.currentUser!.tenantId;

  const apiKeys = await ApiKeyService.listApiKeys(tenantId);

  res.send(apiKeys);
});

// Revoke an API key
router.delete("/:key", async (req: Request, res: Response) => {
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
router.post("/:key/rotate", async (req: Request, res: Response) => {
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
});

export { apiKeysRouter };
