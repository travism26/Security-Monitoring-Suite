// routes/metrics.js
import express from "express";
import client from "prom-client";
import { kafkaWrapper } from "../kafka/kafka-wrapper";

const router = express.Router();

// Create a Registry which registers the metrics
const register = new client.Registry();
register.setDefaultLabels({
  app: "system-monitoring-gateway",
});

// Enable the collection of default metrics
client.collectDefaultMetrics({ register });

// Custom HTTP request duration histogram
const httpRequestDurationMicroseconds = new client.Histogram({
  name: "http_request_duration_seconds",
  help: "Duration of HTTP requests in seconds",
  labelNames: ["method", "route", "code"],
});
register.registerMetric(httpRequestDurationMicroseconds);

// Middleware to track the duration of HTTP requests
router.use((req, res, next) => {
  const end = httpRequestDurationMicroseconds.startTimer();
  res.on("finish", () => {
    end({ method: req.method, route: req.url, code: res.statusCode });
  });
  next();
});

// Expose metrics at the /metrics endpoint
router.get("/metrics", async (req, res) => {
  res.set("Content-Type", register.contentType);
  res.end(await register.metrics());
});

router.get("/health", (req, res) => {
  if (!kafkaWrapper.isInitialized()) {
    return res.status(503).json({
      errors: [{ message: "Kafka service temporarily unavailable" }],
    });
  }
  res.status(200).send("OK");
});

export { router as metricsRouter, register as metricsRegistry };
