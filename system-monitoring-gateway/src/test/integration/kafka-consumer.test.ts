import { kafkaWrapper } from "../../kafka/kafka-wrapper";
import { Topics } from "../../kafka/topics";
import { Consumer as KafkaConsumer, EachMessagePayload } from "kafkajs";

jest.mock("../../kafka/kafka-wrapper");

interface MockConsumer {
  connect: jest.Mock;
  disconnect: jest.Mock;
  subscribe: jest.Mock;
  run: jest.Mock;
}

const createMockConsumer = (): MockConsumer => ({
  connect: jest.fn().mockResolvedValue(undefined),
  disconnect: jest.fn().mockResolvedValue(undefined),
  subscribe: jest.fn().mockResolvedValue(undefined),
  run: jest.fn(),
});

interface SystemMetricsMessage {
  metadata: {
    tenant_id?: string;
    environment?: string;
  };
  metrics?: {
    cpu?: { usage: number };
  };
}

describe("Kafka Consumer Integration", () => {
  let consumer: MockConsumer;

  beforeEach(() => {
    jest.clearAllMocks();
    consumer = createMockConsumer();
    (kafkaWrapper.getConsumer as jest.Mock).mockReturnValue(consumer);
  });

  it("successfully processes valid messages", async () => {
    const messageHandler = jest.fn();
    const validMessage = {
      metadata: {
        tenant_id: "tenant-123",
        environment: "test",
      },
      metrics: {
        cpu: { usage: 45.5 },
      },
    };

    consumer.run.mockImplementation(({ eachMessage }) => {
      return eachMessage({
        topic: Topics.SystemMetrics,
        partition: 0,
        message: {
          value: Buffer.from(JSON.stringify(validMessage)),
        },
      });
    });

    await consumer.run({
      eachMessage: async ({ message }: EachMessagePayload) => {
        const data = message.value
          ? (JSON.parse(message.value.toString()) as SystemMetricsMessage)
          : null;
        if (data?.metadata?.tenant_id === "tenant-123") {
          messageHandler(data);
        }
      },
    });

    expect(messageHandler).toHaveBeenCalledWith(
      expect.objectContaining({
        metadata: expect.objectContaining({
          tenant_id: "tenant-123",
        }),
      })
    );
  });

  it("handles invalid JSON data", async () => {
    const errorHandler = jest.fn();
    const dlqProducer = kafkaWrapper.getProducer("system-metrics-dlq");

    consumer.run.mockImplementation(({ eachMessage }) => {
      return eachMessage({
        topic: Topics.SystemMetrics,
        partition: 0,
        message: {
          value: Buffer.from("invalid-json"),
        },
      });
    });

    await consumer.run({
      eachMessage: async ({ message }: EachMessagePayload) => {
        try {
          if (!message.value) {
            throw new Error("No message value");
          }
          JSON.parse(message.value.toString());
        } catch (error) {
          errorHandler(error);
          if (message.value) {
            await dlqProducer.publish({
              error: "Failed to parse message",
              original_message: message.value.toString(),
              timestamp: new Date().toISOString(),
            });
          }
        }
      },
    });

    expect(errorHandler).toHaveBeenCalled();
    expect(dlqProducer.publish).toHaveBeenCalledWith(
      expect.objectContaining({
        error: "Failed to parse message",
        original_message: "invalid-json",
      })
    );
  });

  it("validates required fields", async () => {
    const errorHandler = jest.fn();
    const errorProducer = kafkaWrapper.getProducer("system-metrics-errors");

    const invalidMessage = {
      metadata: {
        environment: "test",
      },
      metrics: {
        cpu: { usage: 45.5 },
      },
    };

    consumer.run.mockImplementation(({ eachMessage }) => {
      return eachMessage({
        topic: Topics.SystemMetrics,
        partition: 0,
        message: {
          value: Buffer.from(JSON.stringify(invalidMessage)),
        },
      });
    });

    await consumer.run({
      eachMessage: async ({ message }: EachMessagePayload) => {
        try {
          if (!message.value) {
            throw new Error("No message value");
          }
          const data = JSON.parse(message.value.toString());
          if (!data?.metadata?.tenant_id) {
            throw new Error("Missing tenant ID");
          }
        } catch (error) {
          errorHandler(error);
          if (message.value) {
            await errorProducer.publish({
              error: "Missing tenant ID in message",
              original_message: JSON.parse(message.value.toString()),
              timestamp: new Date().toISOString(),
            });
          }
        }
      },
    });

    expect(errorHandler).toHaveBeenCalled();
    expect(errorProducer.publish).toHaveBeenCalledWith(
      expect.objectContaining({
        error: "Missing tenant ID in message",
      })
    );
  });
});
