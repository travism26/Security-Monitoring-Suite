import crypto from "crypto";

export interface ApiKey {
  key: string;
  tenantId: string;
  createdAt: Date;
  expiresAt?: Date;
  isActive: boolean;
}

export class ApiKeyService {
  // In a real implementation, this would be stored in a database
  private static apiKeys: Map<string, ApiKey> = new Map();

  static generateApiKey(tenantId: string, expiresInDays?: number): ApiKey {
    const randomBytes = crypto.randomBytes(24);
    const key = `sms_${randomBytes.toString("hex")}`;

    const apiKey: ApiKey = {
      key,
      tenantId,
      createdAt: new Date(),
      expiresAt: expiresInDays
        ? new Date(Date.now() + expiresInDays * 24 * 60 * 60 * 1000)
        : undefined,
      isActive: true,
    };

    this.apiKeys.set(key, apiKey);
    return apiKey;
  }

  static async validateApiKey(key: string): Promise<ApiKey | null> {
    const apiKey = this.apiKeys.get(key);

    if (!apiKey) {
      return null;
    }

    if (!apiKey.isActive) {
      return null;
    }

    if (apiKey.expiresAt && apiKey.expiresAt < new Date()) {
      return null;
    }

    return apiKey;
  }

  static async revokeApiKey(key: string): Promise<boolean> {
    const apiKey = this.apiKeys.get(key);

    if (!apiKey) {
      return false;
    }

    apiKey.isActive = false;
    this.apiKeys.set(key, apiKey);
    return true;
  }

  static async listApiKeys(tenantId: string): Promise<ApiKey[]> {
    return Array.from(this.apiKeys.values()).filter(
      (apiKey) => apiKey.tenantId === tenantId
    );
  }

  static async rotateApiKey(
    oldKey: string,
    expiresInDays?: number
  ): Promise<ApiKey | null> {
    const existingKey = await this.validateApiKey(oldKey);

    if (!existingKey) {
      return null;
    }

    // Revoke the old key
    await this.revokeApiKey(oldKey);

    // Generate a new key
    return this.generateApiKey(existingKey.tenantId, expiresInDays);
  }
}
