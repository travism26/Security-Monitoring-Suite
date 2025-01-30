import mongoose from "mongoose";
import { Tenant, TenantDoc } from "../models/tenant";
import { ApiKey, ApiKeyDoc } from "../models/api-key";

export class MongoDBService {
  private static instance: MongoDBService;

  private constructor() {}

  static getInstance(): MongoDBService {
    if (!MongoDBService.instance) {
      MongoDBService.instance = new MongoDBService();
    }
    return MongoDBService.instance;
  }

  async connect(mongoUri: string): Promise<void> {
    try {
      await mongoose.connect(mongoUri, {
        maxPoolSize: 10,
        serverSelectionTimeoutMS: 5000,
      });
      console.log("Connected to MongoDB");
    } catch (error) {
      console.error("MongoDB connection error:", error);
      throw error;
    }
  }

  async disconnect(): Promise<void> {
    try {
      await mongoose.disconnect();
      console.log("Disconnected from MongoDB");
    } catch (error) {
      console.error("MongoDB disconnection error:", error);
      throw error;
    }
  }

  // Tenant Operations
  async createTenant(
    organizationName: string,
    contactEmail: string
  ): Promise<TenantDoc> {
    const tenant = Tenant.build({ organizationName, contactEmail });
    return await tenant.save();
  }

  async getTenantById(id: string): Promise<TenantDoc | null> {
    return await Tenant.findById(id);
  }

  async updateTenant(
    id: string,
    updates: Partial<TenantDoc>
  ): Promise<TenantDoc | null> {
    return await Tenant.findByIdAndUpdate(id, updates, { new: true });
  }

  async deleteTenant(id: string): Promise<TenantDoc | null> {
    // Delete associated API keys first
    await ApiKey.deleteMany({ tenantId: id });
    return await Tenant.findByIdAndDelete(id);
  }

  // API Key Operations
  async createApiKey(
    tenantId: string,
    permissions?: string[],
    expiresAt?: Date
  ): Promise<ApiKeyDoc> {
    const key = this.generateApiKey();
    const apiKey = ApiKey.build({
      key,
      tenantId,
      permissions,
      expiresAt,
    });
    return await apiKey.save();
  }

  async getApiKeyByKey(key: string): Promise<ApiKeyDoc | null> {
    return await ApiKey.findOne({ key, isActive: true });
  }

  async getApiKeysByTenant(tenantId: string): Promise<ApiKeyDoc[]> {
    return await ApiKey.find({ tenantId });
  }

  async deactivateApiKey(key: string): Promise<ApiKeyDoc | null> {
    const apiKey = await ApiKey.findOne({ key });
    if (!apiKey) {
      return null;
    }
    apiKey.isActive = false;
    return await apiKey.save();
  }

  async deleteApiKey(key: string): Promise<ApiKeyDoc | null> {
    const apiKey = await ApiKey.findOne({ key });
    if (!apiKey) {
      return null;
    }
    await apiKey.deleteOne();
    return apiKey;
  }

  // Helper Methods
  private generateApiKey(): string {
    // Generate a random 32-character hexadecimal string
    return Array.from({ length: 32 }, () =>
      Math.floor(Math.random() * 16).toString(16)
    ).join("");
  }

  // Health Check
  async healthCheck(): Promise<boolean> {
    try {
      await mongoose.connection.db.admin().ping();
      return true;
    } catch (error) {
      console.error("MongoDB health check failed:", error);
      return false;
    }
  }
}

export const mongoDBService = MongoDBService.getInstance();
