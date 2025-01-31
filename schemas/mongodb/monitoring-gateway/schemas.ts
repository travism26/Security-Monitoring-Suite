/**
 * Centralized MongoDB Schemas for System Monitoring Gateway
 * Version: 1.0.0
 * Created: 2025-01-31
 * Last Modified: 2025-01-31
 */

import { JSONSchemaType } from "ajv";

// TypeScript Interfaces
export interface ITenant {
  organizationName: string;
  contactEmail: string;
  status: "active" | "inactive";
  createdAt: Date;
  updatedAt: Date;
}

export interface IApiKey {
  key: string;
  tenantId: string;
  permissions: Array<"read" | "write" | "admin">;
  expiresAt: Date;
  isActive: boolean;
  createdAt: Date;
  updatedAt: Date;
}

// MongoDB JSON Schemas
export const tenantSchema: JSONSchemaType<
  Omit<ITenant, "createdAt" | "updatedAt">
> = {
  type: "object",
  properties: {
    organizationName: { type: "string" },
    contactEmail: { type: "string", format: "email" },
    status: { type: "string", enum: ["active", "inactive"] },
  },
  required: ["organizationName", "contactEmail"],
  additionalProperties: false,
};

export const apiKeySchema: JSONSchemaType<
  Omit<IApiKey, "createdAt" | "updatedAt">
> = {
  type: "object",
  properties: {
    key: { type: "string" },
    tenantId: { type: "string" },
    permissions: {
      type: "array",
      items: { type: "string", enum: ["read", "write", "admin"] },
    },
    expiresAt: { type: "string", format: "date-time" },
    isActive: { type: "boolean" },
  },
  required: ["key", "tenantId", "expiresAt"],
  additionalProperties: false,
};

// Mongoose Schema Configuration
export const mongooseTenantConfig = {
  organizationName: {
    type: String,
    required: true,
  },
  contactEmail: {
    type: String,
    required: true,
    unique: true,
  },
  status: {
    type: String,
    required: true,
    enum: ["active", "inactive"],
    default: "active",
  },
};

export const mongooseApiKeyConfig = {
  key: {
    type: String,
    required: true,
    unique: true,
  },
  tenantId: {
    type: "ObjectId",
    ref: "Tenant",
    required: true,
  },
  permissions: {
    type: [String],
    default: ["read"],
    enum: ["read", "write", "admin"],
  },
  expiresAt: {
    type: Date,
    required: true,
    default: () => new Date(Date.now() + 365 * 24 * 60 * 60 * 1000), // 1 year
  },
  isActive: {
    type: Boolean,
    default: true,
  },
};

// Indexes Configuration
export const tenantIndexes = [{ key: { contactEmail: 1 }, unique: true }];

export const apiKeyIndexes = [
  { key: { key: 1 }, unique: true },
  { key: { tenantId: 1 } },
  { key: { expiresAt: 1 } },
];

/**
 * Schema Documentation:
 *
 * Tenant Schema:
 * - organizationName: Name of the organization
 * - contactEmail: Primary contact email (unique)
 * - status: Account status (active/inactive)
 * - createdAt: Record creation timestamp
 * - updatedAt: Last modification timestamp
 *
 * API Key Schema:
 * - key: Unique API key string
 * - tenantId: Reference to parent tenant
 * - permissions: Array of permission levels
 * - expiresAt: Key expiration timestamp
 * - isActive: Whether the key is currently active
 * - createdAt: Record creation timestamp
 * - updatedAt: Last modification timestamp
 */
