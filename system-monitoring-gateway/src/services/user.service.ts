import crypto from "crypto";
import { mongoDBService } from "./mongodb.service";
import { User, UserDoc } from "../models/user";
import { Types } from "mongoose";

export interface UserData {
  id: string;
  email: string;
  firstName: string;
  lastName: string;
  role: string;
  tenantId: string;
  status: string;
  lastLogin?: Date;
  emailVerified: boolean;
  createdAt: Date;
  updatedAt: Date;
}

export interface CreateUserData {
  email: string;
  password: string;
  firstName: string;
  lastName: string;
  role?: string;
  tenantId: string;
}

export interface UpdateUserData {
  firstName?: string;
  lastName?: string;
  role?: string;
  status?: string;
}

// Convert MongoDB document to API interface
const toUserData = (doc: UserDoc): UserData => ({
  id: doc.id,
  email: doc.email,
  firstName: doc.firstName,
  lastName: doc.lastName,
  role: doc.role,
  tenantId: doc.tenantId.toString(),
  status: doc.status,
  lastLogin: doc.lastLogin,
  emailVerified: doc.emailVerified,
  createdAt: doc.createdAt,
  updatedAt: doc.updatedAt,
});

export class UserService {
  static async createUser(userData: CreateUserData): Promise<UserData> {
    const user = await mongoDBService.createUser({
      ...userData,
      tenantId: new Types.ObjectId(userData.tenantId),
      verificationToken: crypto.randomBytes(32).toString("hex"),
    });
    return toUserData(user);
  }

  static async getUserById(userId: string): Promise<UserData | null> {
    const user = await mongoDBService.getUserById(userId);
    return user ? toUserData(user) : null;
  }

  static async getUserByEmail(email: string): Promise<UserData | null> {
    const user = await mongoDBService.getUserByEmail(email);
    return user ? toUserData(user) : null;
  }

  static async updateUser(
    userId: string,
    updateData: UpdateUserData
  ): Promise<UserData | null> {
    const user = await mongoDBService.updateUser(userId, updateData);
    return user ? toUserData(user) : null;
  }

  static async verifyEmail(token: string): Promise<boolean> {
    const user = await mongoDBService.getUserByVerificationToken(token);
    if (!user) {
      return false;
    }

    user.emailVerified = true;
    user.verificationToken = undefined;
    await user.save();
    return true;
  }

  static async initiatePasswordReset(email: string): Promise<boolean> {
    const user = await mongoDBService.getUserByEmail(email);
    if (!user) {
      return false;
    }

    const resetToken = crypto.randomBytes(32).toString("hex");
    const resetExpires = new Date(Date.now() + 24 * 60 * 60 * 1000); // 24 hours

    user.passwordResetToken = resetToken;
    user.passwordResetExpires = resetExpires;
    await user.save();

    return true;
  }

  static async resetPassword(
    token: string,
    newPassword: string
  ): Promise<boolean> {
    const user = await mongoDBService.getUserByResetToken(token);
    if (
      !user ||
      !user.passwordResetExpires ||
      user.passwordResetExpires < new Date()
    ) {
      return false;
    }

    user.password = newPassword;
    user.passwordResetToken = undefined;
    user.passwordResetExpires = undefined;
    await user.save();

    return true;
  }

  static async validateCredentials(
    email: string,
    password: string
  ): Promise<UserData | null> {
    const user = await mongoDBService.getUserByEmail(email);
    if (!user) {
      return null;
    }

    const isValid = await user.comparePassword(password);
    if (!isValid) {
      return null;
    }

    // Update last login
    user.lastLogin = new Date();
    await user.save();

    return toUserData(user);
  }

  static async listUsersByTenant(tenantId: string): Promise<UserData[]> {
    const users = await mongoDBService.getUsersByTenant(tenantId);
    return users.map(toUserData);
  }

  static async deactivateUser(userId: string): Promise<boolean> {
    const user = await mongoDBService.updateUser(userId, {
      status: "inactive",
    });
    return user !== null;
  }

  static async activateUser(userId: string): Promise<boolean> {
    const user = await mongoDBService.updateUser(userId, { status: "active" });
    return user !== null;
  }
}
