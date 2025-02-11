import { ApiKeyService } from "../api-key.service";
import { mongoDBService } from "../mongodb.service";
import mongoose from "mongoose";

jest.mock("../mongodb.service");

describe("ApiKeyService", () => {
  const mockUser = {
    id: new mongoose.Types.ObjectId().toString(),
    tenantId: new mongoose.Types.ObjectId().toString(),
  };

  const mockApiKey = {
    key: "test_key",
    tenantId: mockUser.tenantId,
    userId: mockUser.id,
    description: "Test Key",
    permissions: ["read"],
    createdAt: new Date(),
    isActive: true,
  };

  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe("generateApiKey", () => {
    it("should generate a new API key with user association", async () => {
      (mongoDBService.createApiKey as jest.Mock).mockResolvedValue(mockApiKey);

      const result = await ApiKeyService.generateApiKey(
        mockUser.tenantId,
        mockUser.id,
        "Test Key",
        ["read"]
      );

      expect(result).toBeDefined();
      expect(result.key).toBe(mockApiKey.key);
      expect(mongoDBService.createApiKey).toHaveBeenCalledWith(
        mockUser.tenantId,
        mockUser.id,
        "Test Key",
        ["read"],
        undefined
      );
    });
  });

  describe("validateApiKey", () => {
    it("should validate and return active API key", async () => {
      (mongoDBService.getApiKeyByKey as jest.Mock).mockResolvedValue(
        mockApiKey
      );

      const result = await ApiKeyService.validateApiKey("test_key");

      expect(result).toBeDefined();
      expect(result?.key).toBe(mockApiKey.key);
      expect(mongoDBService.getApiKeyByKey).toHaveBeenCalledWith("test_key");
    });

    it("should return null for inactive API key", async () => {
      (mongoDBService.getApiKeyByKey as jest.Mock).mockResolvedValue({
        ...mockApiKey,
        isActive: false,
      });

      const result = await ApiKeyService.validateApiKey("test_key");

      expect(result).toBeNull();
    });

    it("should return null and deactivate expired API key", async () => {
      const expiredKey = {
        ...mockApiKey,
        expiresAt: new Date(Date.now() - 24 * 60 * 60 * 1000), // 1 day ago
      };
      (mongoDBService.getApiKeyByKey as jest.Mock).mockResolvedValue(
        expiredKey
      );
      (mongoDBService.deactivateApiKey as jest.Mock).mockResolvedValue(
        expiredKey
      );

      const result = await ApiKeyService.validateApiKey("test_key");

      expect(result).toBeNull();
      expect(mongoDBService.deactivateApiKey).toHaveBeenCalledWith("test_key");
    });
  });

  describe("getUserApiKeys", () => {
    it("should return list of user API keys", async () => {
      const mockApiKeys = [mockApiKey];
      (mongoDBService.getApiKeysByUser as jest.Mock).mockResolvedValue(
        mockApiKeys
      );

      const result = await ApiKeyService.getUserApiKeys(mockUser.id);

      expect(result).toHaveLength(1);
      expect(result[0].key).toBe(mockApiKey.key);
      expect(mongoDBService.getApiKeysByUser).toHaveBeenCalledWith(mockUser.id);
    });
  });

  describe("validateUserAccess", () => {
    it("should return true when user owns the API key", async () => {
      (mongoDBService.getApiKeyById as jest.Mock).mockResolvedValue(mockApiKey);

      const result = await ApiKeyService.validateUserAccess(
        mockUser.id,
        "test_key"
      );

      expect(result).toBe(true);
      expect(mongoDBService.getApiKeyById).toHaveBeenCalledWith("test_key");
    });

    it("should return false when user does not own the API key", async () => {
      const differentUserId = new mongoose.Types.ObjectId().toString();
      (mongoDBService.getApiKeyById as jest.Mock).mockResolvedValue({
        ...mockApiKey,
        userId: differentUserId,
      });

      const result = await ApiKeyService.validateUserAccess(
        mockUser.id,
        "test_key"
      );

      expect(result).toBe(false);
    });

    it("should return false when API key does not exist", async () => {
      (mongoDBService.getApiKeyById as jest.Mock).mockResolvedValue(null);

      const result = await ApiKeyService.validateUserAccess(
        mockUser.id,
        "test_key"
      );

      expect(result).toBe(false);
    });
  });

  describe("revokeApiKey", () => {
    it("should successfully revoke an API key", async () => {
      (mongoDBService.deactivateApiKey as jest.Mock).mockResolvedValue(
        mockApiKey
      );

      const result = await ApiKeyService.revokeApiKey("test_key");

      expect(result).toBe(true);
      expect(mongoDBService.deactivateApiKey).toHaveBeenCalledWith("test_key");
    });

    it("should return false when API key not found", async () => {
      (mongoDBService.deactivateApiKey as jest.Mock).mockResolvedValue(null);

      const result = await ApiKeyService.revokeApiKey("test_key");

      expect(result).toBe(false);
    });
  });
});
