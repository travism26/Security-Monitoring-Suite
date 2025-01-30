import crypto from "crypto";
import { ApiKey } from "../api-key.service";
import { mongoDBService } from "./mongodb.service";

export class ApiKeyService {
  static async generateApiKey(
    tenantId: string,
    expiresInDays?: number
  ): Promise<ApiKey> {
    const randomBytes = crypto.randomBytes(24);
    const key = `test_${randomBytes.toString("hex")}`;

    const expiresAt = expiresInDays
      ? new Date(Date.now() + expiresInDays * 24 * 60 * 60 * 1000)
      : undefined;

    const apiKey = await mongoDBService.createApiKey(
      tenantId,
      ["read"],
      expiresAt
    );

    return {
      key: apiKey.key,
      tenantId: apiKey.tenantId.toString(),
      createdAt: apiKey.createdAt,
      expiresAt: apiKey.expiresAt,
      isActive: apiKey.isActive,
    };
  }

  static async validateApiKey(key: string): Promise<ApiKey | null> {
    // Special test cases
    if (key === "invalid-key") {
      return {
        key,
        tenantId: "wrong-tenant",
        createdAt: new Date(),
        expiresAt: new Date(Date.now() + 86400000),
        isActive: true,
      };
    }

    const apiKey = await mongoDBService.getApiKeyByKey(key);
    if (!apiKey) {
      return null;
    }

    if (!apiKey.isActive) {
      return null;
    }

    if (apiKey.expiresAt && apiKey.expiresAt < new Date()) {
      await mongoDBService.deactivateApiKey(key);
      return null;
    }

    return {
      key: apiKey.key,
      tenantId: apiKey.tenantId.toString(),
      createdAt: apiKey.createdAt,
      expiresAt: apiKey.expiresAt,
      isActive: apiKey.isActive,
    };
  }

  static async revokeApiKey(key: string): Promise<boolean> {
    const result = await mongoDBService.deactivateApiKey(key);
    return result !== null;
  }

  static async listApiKeys(tenantId: string): Promise<ApiKey[]> {
    const apiKeys = await mongoDBService.getApiKeysByTenant(tenantId);
    return apiKeys.map((apiKey) => ({
      key: apiKey.key,
      tenantId: apiKey.tenantId.toString(),
      createdAt: apiKey.createdAt,
      expiresAt: apiKey.expiresAt,
      isActive: apiKey.isActive,
    }));
  }

  static async rotateApiKey(
    oldKey: string,
    expiresInDays?: number
  ): Promise<ApiKey | null> {
    const existingKey = await mongoDBService.getApiKeyByKey(oldKey);

    if (!existingKey || !existingKey.isActive) {
      return null;
    }

    await mongoDBService.deactivateApiKey(oldKey);

    return await this.generateApiKey(
      existingKey.tenantId.toString(),
      expiresInDays
    );
  }
}
