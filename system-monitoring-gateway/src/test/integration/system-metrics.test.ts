import request from "supertest";
import { app } from "../../app";
import { kafkaWrapper } from "../../kafka/kafka-wrapper";
import { Topics } from "../../kafka/topics";

jest.mock("../../kafka/kafka-wrapper");

const createTestPayload = (tenantId: string) => ({
  topic: Topics.SystemMetrics,
  data: {
    timestamp: new Date().toISOString(),
    tenant: {
      id: tenantId,
      metadata: {
        environment: "test",
      },
    },
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
      },
      metadata: {
        collection_duration: "1.5s",
        collector_count: 3,
      },
      threat_indicators: [],
      processes: {
        total_count: 0,
        total_cpu_percent: 0,
        total_memory_usage: 0,
        process_list: [],
      },
    },
  },
  timestamp: new Date().toISOString(),
});

describe("System Metrics Integration", () => {
  beforeEach(() => {
    jest.clearAllMocks();
    (kafkaWrapper.isInitialized as jest.Mock).mockReturnValue(true);
    (kafkaWrapper.getProducer as jest.Mock).mockReturnValue({
      publish: jest.fn().mockResolvedValue(undefined),
    });
  });

  it("successfully processes valid metrics", async () => {
    const tenantId = "tenant-123";
    const metricPayload = createTestPayload(tenantId);

    const headers = {
      "X-Tenant-ID": tenantId,
      "X-API-Key": "test-api-key",
      "X-Tenant-Environment": "test",
      "Content-Type": "application/json",
    };

    const response = await request(app)
      .post("/api/v1/system/metrics/ingest")
      .set(headers)
      .send(metricPayload);

    expect(response.status).toBe(202);
    expect(response.body).toHaveProperty("status", "accepted");
    expect(response.body).toHaveProperty("timestamp");

    const kafkaProducer = kafkaWrapper.getProducer("system-metrics");
    expect(kafkaProducer.publish).toHaveBeenCalledWith(
      expect.objectContaining({
        data: expect.objectContaining({
          metrics: expect.any(Object),
        }),
        tenant: expect.objectContaining({
          id: tenantId,
        }),
      })
    );
  });

  it("validates tenant ID consistency", async () => {
    const headerTenantId = "tenant-789";
    const payloadTenantId = "different-tenant";
    const metricPayload = createTestPayload(payloadTenantId);

    const headers = {
      "X-Tenant-ID": headerTenantId,
      "X-API-Key": "test-api-key",
      "X-Tenant-Environment": "test",
      "Content-Type": "application/json",
    };

    const response = await request(app)
      .post("/api/v1/system/metrics/ingest")
      .set(headers)
      .send(metricPayload);

    expect(response.status).toBe(400);
    expect(response.body.errors[0].message).toContain("Tenant ID mismatch");
  });

  it("handles malformed data", async () => {
    const tenantId = "tenant-101";
    const malformedPayload = {
      topic: Topics.SystemMetrics,
      data: "invalid-json-structure",
      timestamp: new Date().toISOString(),
    };

    const headers = {
      "X-Tenant-ID": tenantId,
      "X-API-Key": "test-api-key",
      "Content-Type": "application/json",
    };

    const response = await request(app)
      .post("/api/v1/system/metrics/ingest")
      .set(headers)
      .send(malformedPayload);

    expect(response.status).toBe(400);
    expect(response.body.errors[0].message).toContain("Missing tenant headers");
  });

  it("handles Kafka unavailability", async () => {
    (kafkaWrapper.isInitialized as jest.Mock).mockReturnValue(false);

    const tenantId = "tenant-202";
    const metricPayload = createTestPayload(tenantId);

    const headers = {
      "X-Tenant-ID": tenantId,
      "X-API-Key": "test-api-key",
      "X-Tenant-Environment": "test",
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
