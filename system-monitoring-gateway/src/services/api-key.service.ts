import crypto from "crypto";
import { mongoDBService } from "./mongodb.service";
import { ApiKeyDoc } from "../models/api-key";

export interface ApiKey {
  key: string;
  tenantId?: string;
  userId: string;
  description: string;
  permissions: string[];
  createdAt: Date;
  expiresAt?: Date;
  isActive: boolean;
}

// Convert MongoDB document to API interface
const toApiKey = (doc: ApiKeyDoc): ApiKey => ({
  key: doc.key,
  tenantId: doc.tenantId.toString(),
  userId: doc.userId.toString(),
  description: doc.description,
  permissions: doc.permissions,
  createdAt: doc.createdAt,
  expiresAt: doc.expiresAt,
  isActive: doc.isActive,
});

export class ApiKeyService {
  static async generateApiKey(
    userId: string,
    description: string,
    tenantId?: string,
    permissions: string[] = ["read"],
    expiresInDays?: number
  ): Promise<ApiKey> {
    const randomBytes = crypto.randomBytes(24);
    const key = `sms_${randomBytes.toString("hex")}`;

    const expiresAt = expiresInDays
      ? new Date(Date.now() + expiresInDays * 24 * 60 * 60 * 1000)
      : undefined;

    const apiKey = await mongoDBService.createApiKey(
      userId,
      description,
      tenantId,
      permissions,
      expiresAt
    );
    return toApiKey(apiKey);
  }

  /**
   * Validates an API key for agent authentication
   * Only checks if the key exists and is active
   * Keeps user/tenant associations for tracking but doesn't use them for validation
   */
  static async validateApiKey(key: string): Promise<ApiKey | null> {
    const apiKey = await mongoDBService.getApiKeyByKey(key);
    console.log("[DEBUG] validateApiKey API key:", apiKey);
    if (!apiKey || !apiKey.isActive) {
      return null;
    }

    // Auto-deactivate expired keys
    if (apiKey.expiresAt && apiKey.expiresAt < new Date()) {
      await mongoDBService.deactivateApiKey(key);
      return null;
    }

    return toApiKey(apiKey);
  }

  static async revokeApiKey(key: string): Promise<boolean> {
    const result = await mongoDBService.deactivateApiKey(key);
    return result !== null;
  }

  static async listApiKeys(
    userId?: string,
    tenantId?: string
  ): Promise<ApiKey[]> {
    const apiKeys = userId
      ? await mongoDBService.getApiKeysByUser(userId)
      : await mongoDBService.getApiKeysByTenant(tenantId);
    return apiKeys.map(toApiKey);
  }

  static async getUserApiKeys(userId: string): Promise<ApiKey[]> {
    const apiKeys = await mongoDBService.getApiKeysByUser(userId);
    return apiKeys.map(toApiKey);
  }

  static async validateUserAccess(
    userId: string,
    keyId: string
  ): Promise<boolean> {
    const apiKey = await mongoDBService.getApiKeyById(keyId);
    if (!apiKey) {
      return false;
    }
    return apiKey.userId.toString() === userId;
  }

  static async rotateApiKey(
    oldKey: string,
    description?: string,
    expiresInDays?: number
  ): Promise<ApiKey | null> {
    const existingKey = await mongoDBService.getApiKeyByKey(oldKey);

    if (!existingKey || !existingKey.isActive) {
      return null;
    }

    // Revoke the old key
    await mongoDBService.deactivateApiKey(oldKey);

    // Generate a new key with same permissions but new description if provided
    return await this.generateApiKey(
      existingKey.userId.toString(),
      description || `Rotated key (${new Date().toISOString()})`,
      existingKey.tenantId?.toString(),
      existingKey.permissions,
      expiresInDays
    );
  }
}
