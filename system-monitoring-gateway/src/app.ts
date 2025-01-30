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

const app = express();
app.set("trust proxy", true);
app.use(json());
app.use(
  cookieSession({ signed: false, secure: false }) //process.env.NODE_ENV !== 'test'
);

// Health check endpoint - no auth required
app.get("/health", (req, res) => {
  res.status(200).send({ status: "healthy" });
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
