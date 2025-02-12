import mongoose from "mongoose";
import { ApiKey } from "../models/api-key";

export const up = async () => {
  try {
    // Get all existing API keys
    const existingKeys = await ApiKey.find({});

    // For each existing key, add default values for new required fields
    for (const key of existingKeys) {
      // Get the first admin user from the tenant as default owner
      const adminUser = await mongoose.model("User").findOne({
        tenantId: key.tenantId,
        role: "admin",
      });

      if (!adminUser) {
        console.warn(
          `No admin user found for tenant ${key.tenantId}, skipping key ${key.key}`
        );
        continue;
      }

      // Update the key with new required fields
      key.userId = adminUser._id;
      key.description = `Legacy API key (auto-migrated)`;
      await key.save();
    }

    console.log(`Successfully migrated ${existingKeys.length} API keys`);
  } catch (error) {
    console.error("Migration failed:", error);
    throw error;
  }
};

export const down = async () => {
  try {
    // Remove the new fields from all documents
    await ApiKey.updateMany(
      {},
      {
        $unset: {
          userId: "",
          description: "",
        },
      }
    );

    // Remove the new indexes
    const collection = mongoose.connection.collection("apikeys");
    await collection.dropIndex("userId_1");
    await collection.dropIndex("tenantId_1_userId_1");

    console.log("Successfully reverted API key user association migration");
  } catch (error) {
    console.error("Migration reversion failed:", error);
    throw error;
  }
};
