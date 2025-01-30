import express from "express";
import "express-async-errors"; // to handle async errors in express
import { json } from "body-parser";
import cookieSession from "cookie-session";
import { NotFoundError } from "./errors/not-found-error";
import { errorHandler } from "./middlewares/error-handler";
// Import routers and middleware
import { metricsRouter } from "./routes/metrics";
import { systemMetricsRouter } from "./routes/system-metrics";
import { apiKeysRouter } from "./routes/api-keys";
import { validateApiKey } from "./middlewares/validate-api-key";
import { validateTenantConsistency } from "./middlewares/validate-tenant";
import { validateJWT } from "./middlewares/require-auth";
import { mongoDBService } from "./services/mongodb.service";

const app = express();

// Initialize MongoDB connection
const MONGODB_URI =
  process.env.MONGODB_URI || "mongodb://localhost:27017/monitoring";
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

// API key management routes
app.use(apiKeysRouter);

// Routes that require API key authentication
app.use(validateApiKey);
app.use(validateTenantConsistency);

// Apply tenant context and rate limiting to all metric routes
if (process.env.NODE_ENV !== "test") {
  app.use("/api/v1/metrics", metricsRouter);
}
app.use(systemMetricsRouter);
app.all("*", async (req, res) => {
  throw new NotFoundError();
});
app.use(errorHandler);

export { app };
