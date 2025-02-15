import express from "express";
import "express-async-errors"; // to handle async errors in express
import { json } from "body-parser";
import cookieSession from "cookie-session";
import { NotFoundError } from "./errors";
import { errorHandler } from "./middleware/error-handler";
// Import routers and middleware
import { metricsRouter } from "./routes/metrics";
import { systemMetricsRouter } from "./routes/system-metrics";
import { apiKeysRouter } from "./routes/api-keys";
import { usersRouter } from "./routes/users";
import { teamsRouter } from "./routes/teams";
import { validateApiKey } from "./middleware/validate-api-key";
import { validateTenantConsistency } from "./middleware/validate-tenant";
import { requireAuth } from "./middleware/require-auth";
import { mongoDBService } from "./services/mongodb.service";

const app = express();

// Initialize MongoDB connection
const MONGODB_URI =
  process.env.MONGODB_URI ||
  "mongodb://localhost:30090/monitoring" ||
  "mongodb://localhost:27017/monitoring";
mongoDBService.connect(MONGODB_URI).catch((err) => {
  console.error("Failed to connect to MongoDB:", err);
  process.exit(1);
});
app.set("trust proxy", true);
app.use(json());
app.use(
  cookieSession({ signed: false, secure: false }) //process.env.NODE_ENV !== 'test'
);

// Health check endpoint - no auth required
app.get("/health", async (req, res) => {
  const mongoHealth = await mongoDBService.healthCheck();
  const status = mongoHealth ? "healthy" : "degraded";
  res.status(mongoHealth ? 200 : 503).send({
    status,
    services: {
      mongodb: mongoHealth ? "connected" : "disconnected",
    },
  });
});

// User, Team, and API key management routes
app.use("/gateway/api/v1", usersRouter);
app.use("/gateway/api/v1", teamsRouter);
app.use("/gateway/api/v1", apiKeysRouter);

// Apply authentication to protected routes
app.use("/api/teams", requireAuth);

// Apply API key validation only to agent metrics routes
app.use(
  "/gateway/api/v1",
  validateApiKey,
  // Note: Tenant validation is optional during design phase
  validateTenantConsistency,
  systemMetricsRouter
);
app.all("*", async (req, res) => {
  throw new NotFoundError();
});
app.use(errorHandler);

export { app };
