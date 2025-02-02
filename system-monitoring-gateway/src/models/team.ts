import mongoose from "mongoose";

// Interface representing a Team document
interface TeamDoc extends mongoose.Document {
  name: string;
  description?: string;
  tenantId: mongoose.Types.ObjectId;
  parentTeamId?: mongoose.Types.ObjectId;
  owner: mongoose.Types.ObjectId;
  members: Array<{
    userId: mongoose.Types.ObjectId;
    role: string;
    joinedAt: Date;
  }>;
  settings: {
    resourceQuota?: {
      maxMembers: number;
      maxStorage: number;
      maxProjects: number;
    };
    visibility: string;
    allowResourceSharing: boolean;
  };
  createdAt: Date;
  updatedAt: Date;
}

// Interface for Team Model
interface TeamModel extends mongoose.Model<TeamDoc> {
  build(attrs: {
    name: string;
    description?: string;
    tenantId: mongoose.Types.ObjectId;
    parentTeamId?: mongoose.Types.ObjectId;
    owner: mongoose.Types.ObjectId;
    settings?: {
      resourceQuota?: {
        maxMembers?: number;
        maxStorage?: number;
        maxProjects?: number;
      };
      visibility?: string;
      allowResourceSharing?: boolean;
    };
  }): TeamDoc;
}

// Team Schema
const teamSchema = new mongoose.Schema(
  {
    name: {
      type: String,
      required: true,
      trim: true,
    },
    description: {
      type: String,
      trim: true,
    },
    tenantId: {
      type: mongoose.Schema.Types.ObjectId,
      ref: "Tenant",
      required: true,
    },
    parentTeamId: {
      type: mongoose.Schema.Types.ObjectId,
      ref: "Team",
    },
    owner: {
      type: mongoose.Schema.Types.ObjectId,
      ref: "User",
      required: true,
    },
    members: [
      {
        userId: {
          type: mongoose.Schema.Types.ObjectId,
          ref: "User",
          required: true,
        },
        role: {
          type: String,
          enum: ["admin", "member"],
          required: true,
        },
        joinedAt: {
          type: Date,
          default: Date.now,
        },
      },
    ],
    settings: {
      resourceQuota: {
        maxMembers: {
          type: Number,
          default: 10,
        },
        maxStorage: {
          type: Number, // in GB
          default: 100,
        },
        maxProjects: {
          type: Number,
          default: 5,
        },
      },
      visibility: {
        type: String,
        enum: ["private", "public"],
        default: "private",
      },
      allowResourceSharing: {
        type: Boolean,
        default: true,
      },
    },
  },
  {
    timestamps: true,
    toJSON: {
      transform(doc, ret: any) {
        ret.id = ret._id;
        delete ret._id;
        delete ret.__v;
      },
    },
  }
);

// Add indexes for performance
teamSchema.index({ tenantId: 1, name: 1 }, { unique: true });
teamSchema.index({ parentTeamId: 1 });

// Add build method to schema
teamSchema.statics.build = (attrs: {
  name: string;
  description?: string;
  tenantId: mongoose.Types.ObjectId;
  parentTeamId?: mongoose.Types.ObjectId;
  owner: mongoose.Types.ObjectId;
  settings?: {
    resourceQuota?: {
      maxMembers?: number;
      maxStorage?: number;
      maxProjects?: number;
    };
    visibility?: string;
    allowResourceSharing?: boolean;
  };
}) => {
  return new Team(attrs);
};

// Create and export Team model
const Team = mongoose.model<TeamDoc, TeamModel>("Team", teamSchema);

export { Team, TeamDoc };
