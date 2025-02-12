import mongoose from "mongoose";
import { TenantDoc } from "./tenant";

// Interface representing an API Key document
interface ApiKeyDoc extends mongoose.Document {
  key: string;
  tenantId?: TenantDoc["_id"];
  userId: mongoose.Types.ObjectId;
  description: string;
  permissions: string[];
  expiresAt: Date;
  isActive: boolean;
  createdAt: Date;
  updatedAt: Date;
}

// Interface for API Key Model
interface ApiKeyModel extends mongoose.Model<ApiKeyDoc> {
  build(attrs: {
    key: string;
    tenantId: string;
    userId: string;
    description: string;
    permissions?: string[];
    expiresAt?: Date;
    isActive?: boolean;
  }): ApiKeyDoc;
}

// API Key Schema
const apiKeySchema = new mongoose.Schema(
  {
    key: {
      type: String,
      required: true,
      unique: true,
    },
    tenantId: {
      type: mongoose.Schema.Types.ObjectId,
      ref: "Tenant",
      required: false,
    },
    userId: {
      type: mongoose.Schema.Types.ObjectId,
      ref: "User",
      required: true,
    },
    description: {
      type: String,
      required: true,
      maxlength: 200,
    },
    permissions: {
      type: [String],
      default: ["read"],
      enum: ["read", "write", "admin"],
    },
    expiresAt: {
      type: Date,
      required: true,
      default: () => new Date(Date.now() + 365 * 24 * 60 * 60 * 1000), // 1 year from now
    },
    isActive: {
      type: Boolean,
      default: true,
    },
  },
  {
    timestamps: true,
    toJSON: {
      transform(doc, ret) {
        ret.id = ret._id;
        delete ret._id;
        delete ret.__v;
      },
    },
  }
);

// Add indexes for performance
apiKeySchema.index({ key: 1 });
apiKeySchema.index({ tenantId: 1 });
apiKeySchema.index({ userId: 1 });
apiKeySchema.index({ expiresAt: 1 });
apiKeySchema.index({ tenantId: 1, userId: 1 }); // Compound index for tenant and user queries

// Add build method to schema
apiKeySchema.statics.build = (attrs: {
  key: string;
  tenantId?: string;
  userId: string;
  description: string;
  permissions?: string[];
  expiresAt?: Date;
  isActive?: boolean;
}) => {
  return new ApiKey(attrs);
};

// Create and export API Key model
const ApiKey = mongoose.model<ApiKeyDoc, ApiKeyModel>("ApiKey", apiKeySchema);

export { ApiKey, ApiKeyDoc };
