import express from "express";
import "express-async-errors"; // to handle async errors in express
import { json } from "body-parser";
import cookieSession from "cookie-session";
import { NotFoundError } from "./errors/not-found-error";
import { errorHandler } from "./middlewares/error-handlers";
// Import routers and middleware
import { metricsRouter } from "./routes/metrics";
import { systemMetricsRouter } from "./routes/system-metrics";
import { apiKeysRouter } from "./routes/api-keys";
import { validateApiKey } from "./middlewares/validate-api-key";

const app = express();
app.set("trust proxy", true);
app.use(json());
app.use(
  cookieSession({ signed: false, secure: false }) //process.env.NODE_ENV !== 'test'
);

// API key management routes (JWT protected)
app.use(apiKeysRouter);

// Routes that require API key authentication
app.use(validateApiKey);

// Conditionally use the metrics router
if (process.env.NODE_ENV !== "test") {
  app.use(metricsRouter);
}
app.use(systemMetricsRouter);
app.all("*", async (req, res) => {
  throw new NotFoundError();
});
app.use(errorHandler);

export { app };
