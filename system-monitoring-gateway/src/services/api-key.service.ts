import crypto from "crypto";
import { mongoDBService } from "./mongodb.service";
import { ApiKeyDoc } from "../models/api-key";

export interface ApiKey {
  key: string;
  tenantId: string;
  createdAt: Date;
  expiresAt?: Date;
  isActive: boolean;
}

// Convert MongoDB document to API interface
const toApiKey = (doc: ApiKeyDoc): ApiKey => ({
  key: doc.key,
  tenantId: doc.tenantId.toString(),
  createdAt: doc.createdAt,
  expiresAt: doc.expiresAt,
  isActive: doc.isActive,
});

export class ApiKeyService {
  static async generateApiKey(
    tenantId: string,
    expiresInDays?: number
  ): Promise<ApiKey> {
    const randomBytes = crypto.randomBytes(24);
    const key = `sms_${randomBytes.toString("hex")}`;

    const expiresAt = expiresInDays
      ? new Date(Date.now() + expiresInDays * 24 * 60 * 60 * 1000)
      : undefined;

    const apiKey = await mongoDBService.createApiKey(
      tenantId,
      ["read"],
      expiresAt
    );
    return toApiKey(apiKey);
  }

  static async validateApiKey(key: string): Promise<ApiKey | null> {
    const apiKey = await mongoDBService.getApiKeyByKey(key);
    if (!apiKey) {
      return null;
    }

    // Additional validation checks
    if (!apiKey.isActive) {
      return null;
    }

    if (apiKey.expiresAt && apiKey.expiresAt < new Date()) {
      await mongoDBService.deactivateApiKey(key); // Auto-deactivate expired keys
      return null;
    }

    return toApiKey(apiKey);
  }

  static async revokeApiKey(key: string): Promise<boolean> {
    const result = await mongoDBService.deactivateApiKey(key);
    return result !== null;
  }

  static async listApiKeys(tenantId: string): Promise<ApiKey[]> {
    const apiKeys = await mongoDBService.getApiKeysByTenant(tenantId);
    return apiKeys.map(toApiKey);
  }

  static async rotateApiKey(
    oldKey: string,
    expiresInDays?: number
  ): Promise<ApiKey | null> {
    const existingKey = await mongoDBService.getApiKeyByKey(oldKey);

    if (!existingKey || !existingKey.isActive) {
      return null;
    }

    // Revoke the old key
    await mongoDBService.deactivateApiKey(oldKey);

    // Generate a new key
    return await this.generateApiKey(
      existingKey.tenantId.toString(),
      expiresInDays
    );
  }
}
