import mongoose from "mongoose";

// Interface representing a Tenant document
interface TenantDoc extends mongoose.Document {
  organizationName: string;
  contactEmail: string;
  status: string;
  createdAt: Date;
  updatedAt: Date;
}

// Interface for Tenant Model
interface TenantModel extends mongoose.Model<TenantDoc> {
  build(attrs: {
    organizationName: string;
    contactEmail: string;
    status?: string;
  }): TenantDoc;
}

// Tenant Schema
const tenantSchema = new mongoose.Schema(
  {
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

// Add build method to schema
tenantSchema.statics.build = (attrs: {
  organizationName: string;
  contactEmail: string;
  status?: string;
}) => {
  return new Tenant(attrs);
};

// Create and export Tenant model
const Tenant = mongoose.model<TenantDoc, TenantModel>("Tenant", tenantSchema);

export { Tenant, TenantDoc };
