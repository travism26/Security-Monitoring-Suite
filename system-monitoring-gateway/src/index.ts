import { app } from "./app";
import { kafkaWrapper } from "./kafka/kafka-wrapper";
import {
  SystemMetricsPublisher,
  SystemMetricsErrorProducer,
  SystemMetricsDLQProducer,
} from "./events/system-metrics-publisher";
import { Topics } from "./kafka/topics";
import { mongoDBService } from "./services/mongodb.service";

const start = async () => {
  try {
    console.log("Starting server...");

    // Initialize Kafka if broker is configured
    if (process.env.KAFKA_BROKER) {
      const clientId =
        process.env.KAFKA_CLIENT_ID || "system-monitoring-gateway";
      await kafkaWrapper.initialize([process.env.KAFKA_BROKER], clientId);
      await kafkaWrapper.addProducer(
        Topics.SystemMetrics,
        new SystemMetricsPublisher(kafkaWrapper.getClient())
      );
      await kafkaWrapper.addProducer(
        Topics.SystemMetricsErrors,
        new SystemMetricsErrorProducer(kafkaWrapper.getClient())
      );
      await kafkaWrapper.addProducer(
        Topics.SystemMetricsDLQ,
        new SystemMetricsDLQProducer(kafkaWrapper.getClient())
      );
      console.log("Kafka initialized successfully");
    }

    const server = app.listen(3000, () => {
      console.log("Listening on port 3000");
    });

    const shutdown = async () => {
      console.log("Shutting down gracefully...");
      server.close(async () => {
        try {
          if (process.env.KAFKA_BROKER) {
            await kafkaWrapper.disconnect();
            console.log("Kafka disconnected");
          }
          await mongoDBService.disconnect();
          console.log("MongoDB disconnected");
          process.exit(0);
        } catch (err) {
          console.error("Error during shutdown:", err);
          process.exit(1);
        }
      });
    };

    process.on("SIGTERM", shutdown);
    process.on("SIGINT", shutdown);
  } catch (err) {
    console.error("Startup error:", err);
    process.exit(1);
  }
};

start().catch((err) => {
  console.error("Fatal error during startup:", err);
  process.exit(1);
});
