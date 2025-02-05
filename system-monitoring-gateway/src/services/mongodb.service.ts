import mongoose from "mongoose";
import { Tenant, TenantDoc } from "../models/tenant";
import { ApiKey, ApiKeyDoc } from "../models/api-key";
import { User, UserDoc } from "../models/user";

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

  // User Operations
  async createUser(userData: {
    email: string;
    password: string;
    firstName: string;
    lastName: string;
    role?: string;
    tenantId?: mongoose.Types.ObjectId;
    verificationToken?: string;
  }): Promise<UserDoc> {
    const user = User.build(userData);
    return await user.save();
  }

  async getUserById(id: string): Promise<UserDoc | null> {
    return await User.findById(id);
  }

  async getUserByEmail(email: string): Promise<UserDoc | null> {
    return await User.findOne({ email });
  }

  async getUserByVerificationToken(token: string): Promise<UserDoc | null> {
    return await User.findOne({ verificationToken: token });
  }

  async getUserByResetToken(token: string): Promise<UserDoc | null> {
    return await User.findOne({
      passwordResetToken: token,
      passwordResetExpires: { $gt: new Date() },
    });
  }

  async getUsersByTenant(tenantId: string): Promise<UserDoc[]> {
    try {
      // During design phase, handle invalid tenant IDs gracefully
      const objectId = mongoose.Types.ObjectId.isValid(tenantId)
        ? new mongoose.Types.ObjectId(tenantId)
        : null;

      if (!objectId) {
        console.log("Invalid tenant ID format - returning empty list");
        return [];
      }

      return await User.find({
        tenantId: objectId,
        $and: [{ tenantId: { $exists: true } }],
      });
    } catch (error) {
      console.log(
        "Error fetching users by tenant - returning empty list:",
        error
      );
      return [];
    }
  }

  async updateUser(
    id: string,
    updates: Partial<Omit<UserDoc, "password">>
  ): Promise<UserDoc | null> {
    return await User.findByIdAndUpdate(id, updates, { new: true });
  }

  async deleteUser(id: string): Promise<UserDoc | null> {
    return await User.findByIdAndDelete(id);
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
