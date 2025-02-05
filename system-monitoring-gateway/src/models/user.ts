import mongoose from "mongoose";
import bcrypt from "bcryptjs";

// Interface representing a User document
interface UserDoc extends mongoose.Document {
  email: string;
  password: string;
  firstName: string;
  lastName: string;
  role: string;
  tenantId?: mongoose.Types.ObjectId;
  status: string;
  lastLogin?: Date;
  emailVerified: boolean;
  verificationToken?: string;
  passwordResetToken?: string;
  passwordResetExpires?: Date;
  createdAt: Date;
  updatedAt: Date;
  comparePassword(candidatePassword: string): Promise<boolean>;
}

// Interface for User Model
interface UserModel extends mongoose.Model<UserDoc> {
  build(attrs: {
    email: string;
    password: string;
    firstName: string;
    lastName: string;
    role?: string;
    tenantId?: mongoose.Types.ObjectId;
    status?: string;
    emailVerified?: boolean;
  }): UserDoc;
}

// User Schema
const userSchema = new mongoose.Schema(
  {
    email: {
      type: String,
      required: true,
      unique: true,
      lowercase: true,
      trim: true,
    },
    password: {
      type: String,
      required: true,
      minlength: 8,
    },
    firstName: {
      type: String,
      required: true,
      trim: true,
    },
    lastName: {
      type: String,
      required: true,
      trim: true,
    },
    role: {
      type: String,
      required: false,
      enum: ["admin", "team_lead", "member"],
      default: "member",
    },
    tenantId: {
      type: mongoose.Schema.Types.ObjectId,
      ref: "Tenant",
      required: false,
    },
    status: {
      type: String,
      required: true,
      enum: ["active", "inactive", "suspended"],
      default: "active",
    },
    lastLogin: {
      type: Date,
    },
    emailVerified: {
      type: Boolean,
      default: false,
    },
    verificationToken: String,
    passwordResetToken: String,
    passwordResetExpires: Date,
  },
  {
    timestamps: true,
    toJSON: {
      transform(doc, ret: any) {
        ret.id = ret._id;
        delete ret._id;
        delete ret.password;
        delete ret.__v;
        delete ret.verificationToken;
        delete ret.passwordResetToken;
        delete ret.passwordResetExpires;
      },
    },
  }
);

// Hash password before saving
userSchema.pre("save", async function (done) {
  if (this.isModified("password")) {
    const salt = await bcrypt.genSalt(12);
    const hashed = await bcrypt.hash(this.get("password"), salt);
    this.set("password", hashed);
  }
  done();
});

// Compare password method
userSchema.methods.comparePassword = async function (
  candidatePassword: string
): Promise<boolean> {
  return bcrypt.compare(candidatePassword, this.password);
};

// Add build method to schema
userSchema.statics.build = (attrs: {
  email: string;
  password: string;
  firstName: string;
  lastName: string;
  role?: string;
  tenantId?: mongoose.Types.ObjectId;
  status?: string;
  emailVerified?: boolean;
}) => {
  return new User(attrs);
};

// Create and export User model
const User = mongoose.model<UserDoc, UserModel>("User", userSchema);

export { User, UserDoc };
