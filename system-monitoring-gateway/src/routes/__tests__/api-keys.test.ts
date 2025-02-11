import request from "supertest";
import { app } from "../../app";
import { ApiKeyService } from "../../services/api-key.service";
import { mongoDBService } from "../../services/mongodb.service";
import mongoose from "mongoose";

jest.mock("../../services/api-key.service");
jest.mock("../../services/mongodb.service");

describe("API Key Routes", () => {
  const mockUser = {
    id: new mongoose.Types.ObjectId().toString(),
    tenantId: new mongoose.Types.ObjectId().toString(),
    role: "user",
  };

  const mockAdmin = {
    id: new mongoose.Types.ObjectId().toString(),
    tenantId: mockUser.tenantId,
    role: "admin",
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

  describe("POST /gateway/api/v1/users/:userId/api-keys", () => {
    it("should create API key for authenticated user", async () => {
      const mockAuthUser = { ...mockUser };
      (mongoDBService.getUserById as jest.Mock).mockResolvedValue({
        ...mockUser,
        tenantId: new mongoose.Types.ObjectId(mockUser.tenantId),
      });
      (ApiKeyService.generateApiKey as jest.Mock).mockResolvedValue(mockApiKey);

      const response = await request(app)
        .post(`/gateway/api/v1/users/${mockUser.id}/api-keys`)
        .set("Authorization", "Bearer mock-token")
        .send({
          description: "Test Key",
        });

      expect(response.status).toBe(201);
      expect(response.body.key).toBe(mockApiKey.key);
      expect(ApiKeyService.generateApiKey).toHaveBeenCalledWith(
        mockUser.tenantId,
        mockUser.id,
        "Test Key",
        ["read"],
        undefined
      );
    });

    it("should not allow creating API key for different user", async () => {
      const differentUserId = new mongoose.Types.ObjectId().toString();

      const response = await request(app)
        .post(`/gateway/api/v1/users/${differentUserId}/api-keys`)
        .set("Authorization", "Bearer mock-token")
        .send({
          description: "Test Key",
        });

      expect(response.status).toBe(403);
      expect(ApiKeyService.generateApiKey).not.toHaveBeenCalled();
    });

    it("should allow admin to create API key for any user in same tenant", async () => {
      const mockAuthAdmin = { ...mockAdmin };
      const targetUser = {
        id: new mongoose.Types.ObjectId().toString(),
        tenantId: mockAdmin.tenantId,
      };

      (mongoDBService.getUserById as jest.Mock).mockResolvedValue({
        ...targetUser,
        tenantId: new mongoose.Types.ObjectId(targetUser.tenantId),
      });
      (ApiKeyService.generateApiKey as jest.Mock).mockResolvedValue({
        ...mockApiKey,
        userId: targetUser.id,
      });

      const response = await request(app)
        .post(`/gateway/api/v1/users/${targetUser.id}/api-keys`)
        .set("Authorization", "Bearer mock-token")
        .send({
          description: "Test Key",
        });

      expect(response.status).toBe(201);
      expect(response.body.userId).toBe(targetUser.id);
    });
  });

  describe("GET /gateway/api/v1/users/:userId/api-keys", () => {
    it("should return user's API keys", async () => {
      (ApiKeyService.getUserApiKeys as jest.Mock).mockResolvedValue([
        mockApiKey,
      ]);

      const response = await request(app)
        .get(`/gateway/api/v1/users/${mockUser.id}/api-keys`)
        .set("Authorization", "Bearer mock-token");

      expect(response.status).toBe(200);
      expect(response.body).toHaveLength(1);
      expect(response.body[0].key).toBe(mockApiKey.key);
    });

    it("should not allow accessing different user's API keys", async () => {
      const differentUserId = new mongoose.Types.ObjectId().toString();

      const response = await request(app)
        .get(`/gateway/api/v1/users/${differentUserId}/api-keys`)
        .set("Authorization", "Bearer mock-token");

      expect(response.status).toBe(403);
      expect(ApiKeyService.getUserApiKeys).not.toHaveBeenCalled();
    });
  });

  describe("PUT /gateway/api/v1/users/:userId/api-keys/:keyId/revoke", () => {
    it("should revoke user's API key", async () => {
      (ApiKeyService.validateUserAccess as jest.Mock).mockResolvedValue(true);
      (ApiKeyService.revokeApiKey as jest.Mock).mockResolvedValue(true);

      const response = await request(app)
        .put(
          `/gateway/api/v1/users/${mockUser.id}/api-keys/${mockApiKey.key}/revoke`
        )
        .set("Authorization", "Bearer mock-token");

      expect(response.status).toBe(200);
      expect(ApiKeyService.revokeApiKey).toHaveBeenCalledWith(mockApiKey.key);
    });

    it("should not allow revoking API key without ownership", async () => {
      (ApiKeyService.validateUserAccess as jest.Mock).mockResolvedValue(false);

      const response = await request(app)
        .put(
          `/gateway/api/v1/users/${mockUser.id}/api-keys/${mockApiKey.key}/revoke`
        )
        .set("Authorization", "Bearer mock-token");

      expect(response.status).toBe(403);
      expect(ApiKeyService.revokeApiKey).not.toHaveBeenCalled();
    });
  });

  describe("DELETE /gateway/api/v1/users/:userId/api-keys/:keyId", () => {
    it("should delete user's API key", async () => {
      (ApiKeyService.validateUserAccess as jest.Mock).mockResolvedValue(true);
      (ApiKeyService.revokeApiKey as jest.Mock).mockResolvedValue(true);

      const response = await request(app)
        .delete(
          `/gateway/api/v1/users/${mockUser.id}/api-keys/${mockApiKey.key}`
        )
        .set("Authorization", "Bearer mock-token");

      expect(response.status).toBe(204);
      expect(ApiKeyService.revokeApiKey).toHaveBeenCalledWith(mockApiKey.key);
    });

    it("should not allow deleting API key without ownership", async () => {
      (ApiKeyService.validateUserAccess as jest.Mock).mockResolvedValue(false);

      const response = await request(app)
        .delete(
          `/gateway/api/v1/users/${mockUser.id}/api-keys/${mockApiKey.key}`
        )
        .set("Authorization", "Bearer mock-token");

      expect(response.status).toBe(403);
      expect(ApiKeyService.revokeApiKey).not.toHaveBeenCalled();
    });
  });
});
