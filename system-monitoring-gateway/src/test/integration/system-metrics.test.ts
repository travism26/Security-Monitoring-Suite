import request from "supertest";
import { app } from "../../app";
import { kafkaWrapper } from "../../kafka/kafka-wrapper";
import { Topics } from "../../kafka/topics";

jest.mock("../../kafka/kafka-wrapper");

describe("System Metrics Integration", () => {
  beforeEach(() => {
    jest.clearAllMocks();
    // Mock Kafka initialization check
    (kafkaWrapper.isInitialized as jest.Mock).mockReturnValue(true);
    // Mock Kafka producer
    (kafkaWrapper.getProducer as jest.Mock).mockReturnValue({
      publish: jest.fn().mockResolvedValue(undefined),
    });
  });

  it("successfully processes valid metrics from agent format", async () => {
    // This payload matches the agent's MetricPayload structure
    const metricPayload = {
      topic: Topics.SystemMetrics,
      data: {
        host_info: {
          os: "darwin",
          arch: "amd64",
          hostname: "test-host",
          cpu_cores: 8,
          go_version: "1.21",
        },
        metrics: {
          cpu: {
            usage: 45.5,
            cores: 8,
          },
          memory: {
            total: 16000000000,
            used: 8000000000,
            free: 8000000000,
          },
          disk: {
            total: 500000000000,
            used: 250000000000,
            free: 250000000000,
          },
        },
        metadata: {
          collection_duration: "1.5s",
          collector_count: 3,
        },
        processes: {
          total_count: 1,
          total_cpu_percent: 2.5,
          total_memory_usage: 1500000,
          process_list: [
            {
              pid: 1234,
              name: "test-process",
              cpu_percent: 2.5,
              memory_usage: 1500000,
              status: "running",
            },
          ],
        },
      },
      timestamp: new Date().toISOString(),
    };

    // Simulate agent headers
    const headers = {
      "X-Tenant-ID": "test-tenant",
      "X-API-Key": "test-api-key",
      "X-Tenant-Environment": "test",
      "X-Tenant-Type": "agent",
      "Content-Type": "application/json",
    };

    const response = await request(app)
      .post("/api/v1/system/metrics/ingest")
      .set(headers)
      .send(metricPayload);

    expect(response.status).toBe(202);
    expect(response.body).toHaveProperty("status", "accepted");
    expect(response.body).toHaveProperty("timestamp");

    // Verify Kafka producer was called with correct data
    const kafkaProducer = kafkaWrapper.getProducer("system-metrics");
    expect(kafkaProducer.publish).toHaveBeenCalledWith({
      ...metricPayload.data,
      timestamp: metricPayload.timestamp,
    });
  });

  it("rejects metrics with missing required fields", async () => {
    const invalidPayload = {
      topic: Topics.SystemMetrics,
      data: {
        host_info: {
          os: "darwin",
          arch: "amd64",
          hostname: "test-host",
          cpu_cores: 8,
          go_version: "1.21",
        },
        metadata: {
          collection_duration: "1.5s",
          collector_count: 3,
        },
        // Missing metrics field
      },
      timestamp: new Date().toISOString(),
    };

    const headers = {
      "X-Tenant-ID": "test-tenant",
      "X-API-Key": "test-api-key",
      "Content-Type": "application/json",
    };

    const response = await request(app)
      .post("/api/v1/system/metrics/ingest")
      .set(headers)
      .send(invalidPayload);

    expect(response.status).toBe(400);
    expect(response.body.errors).toBeDefined();
    expect(response.body.errors[0].message).toContain(
      "Metrics data is required"
    );
  });

  it("handles Kafka unavailability gracefully", async () => {
    (kafkaWrapper.isInitialized as jest.Mock).mockReturnValue(false);

    const metricPayload = {
      topic: Topics.SystemMetrics,
      data: {
        host_info: {
          os: "darwin",
          arch: "amd64",
          hostname: "test-host",
          cpu_cores: 8,
          go_version: "1.21",
        },
        metrics: {
          cpu: { usage: 45.5 },
        },
        metadata: {
          collection_duration: "1.5s",
          collector_count: 3,
        },
      },
      timestamp: new Date().toISOString(),
    };

    const headers = {
      "X-Tenant-ID": "test-tenant",
      "X-API-Key": "test-api-key",
      "X-Tenant-Environment": "test",
      "X-Tenant-Type": "agent",
      "Content-Type": "application/json",
    };

    const response = await request(app)
      .post("/api/v1/system/metrics/ingest")
      .set(headers)
      .send(metricPayload);

    expect(response.status).toBe(503);
    expect(response.body.errors[0].message).toBe(
      "Metrics service temporarily unavailable"
    );
  });
});
