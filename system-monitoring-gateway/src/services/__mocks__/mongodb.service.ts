import { TenantDoc } from "../../models/tenant";
import { ApiKeyDoc } from "../../models/api-key";

class MockMongoDBService {
  private static instance: MockMongoDBService;
  private mockTenants: Map<string, any>;
  private mockApiKeys: Map<string, any>;

  private constructor() {
    this.mockTenants = new Map();
    this.mockApiKeys = new Map();
  }

  static getInstance(): MockMongoDBService {
    if (!MockMongoDBService.instance) {
      MockMongoDBService.instance = new MockMongoDBService();
    }
    return MockMongoDBService.instance;
  }

  async connect(mongoUri: string): Promise<void> {
    return Promise.resolve();
  }

  async disconnect(): Promise<void> {
    return Promise.resolve();
  }

  async createTenant(
    organizationName: string,
    contactEmail: string
  ): Promise<TenantDoc> {
    const tenant = {
      id: Math.random().toString(36).substring(7),
      organizationName,
      contactEmail,
      status: "active",
      createdAt: new Date(),
      updatedAt: new Date(),
    };
    this.mockTenants.set(tenant.id, tenant);
    return tenant as any;
  }

  async getTenantById(id: string): Promise<TenantDoc | null> {
    return (this.mockTenants.get(id) as TenantDoc) || null;
  }

  async updateTenant(
    id: string,
    updates: Partial<TenantDoc>
  ): Promise<TenantDoc | null> {
    const tenant = this.mockTenants.get(id);
    if (!tenant) return null;
    const updatedTenant = { ...tenant, ...updates, updatedAt: new Date() };
    this.mockTenants.set(id, updatedTenant);
    return updatedTenant as any;
  }

  async deleteTenant(id: string): Promise<TenantDoc | null> {
    const tenant = this.mockTenants.get(id);
    this.mockTenants.delete(id);
    return (tenant as TenantDoc) || null;
  }

  async createApiKey(
    tenantId: string,
    permissions?: string[],
    expiresAt?: Date
  ): Promise<ApiKeyDoc> {
    const apiKey = {
      key: Math.random().toString(36).substring(7),
      tenantId,
      permissions: permissions || ["read"],
      expiresAt: expiresAt || new Date(Date.now() + 365 * 24 * 60 * 60 * 1000),
      isActive: true,
      createdAt: new Date(),
      updatedAt: new Date(),
    };
    this.mockApiKeys.set(apiKey.key, apiKey);
    return apiKey as any;
  }

  async getApiKeyByKey(key: string): Promise<ApiKeyDoc | null> {
    return (this.mockApiKeys.get(key) as ApiKeyDoc) || null;
  }

  async getApiKeysByTenant(tenantId: string): Promise<ApiKeyDoc[]> {
    return Array.from(this.mockApiKeys.values()).filter(
      (key) => key.tenantId === tenantId
    ) as ApiKeyDoc[];
  }

  async deactivateApiKey(key: string): Promise<ApiKeyDoc | null> {
    const apiKey = this.mockApiKeys.get(key);
    if (!apiKey) return null;
    apiKey.isActive = false;
    apiKey.updatedAt = new Date();
    this.mockApiKeys.set(key, apiKey);
    return apiKey as any;
  }

  async deleteApiKey(key: string): Promise<ApiKeyDoc | null> {
    const apiKey = this.mockApiKeys.get(key);
    this.mockApiKeys.delete(key);
    return (apiKey as ApiKeyDoc) || null;
  }

  async healthCheck(): Promise<boolean> {
    return true;
  }
}

export const mongoDBService = MockMongoDBService.getInstance();
