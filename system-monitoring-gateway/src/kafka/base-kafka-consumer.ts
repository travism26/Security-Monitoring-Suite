import { Kafka, Consumer as KafkaConsumer, EachMessagePayload } from "kafkajs";
import { Event } from "./event";
import { kafkaWrapper } from "./kafka-wrapper";

export abstract class Consumer<T extends Event> {
  abstract topic: T["topic"];
  abstract onMessage(data: T["data"]): void;

  protected consumer: KafkaConsumer;
  protected groupId: string;

  constructor(client: Kafka, groupId: string) {
    this.groupId = groupId;
    this.consumer = client.consumer({ groupId: this.groupId });
  }

  async connect(): Promise<void> {
    await this.consumer.connect();
  }

  async listen() {
    await this.consumer.subscribe({ topic: this.topic, fromBeginning: true });
    await this.consumer.run({
      eachMessage: async (payload: EachMessagePayload) => {
        console.log(`Message received: ${this.topic} / ${this.groupId}`);
        try {
          if (!payload.message.value) {
            throw new Error("No message value");
          }

          const parsedData = await this.parseMessage(payload.message.value);
          await this.validateMessage(parsedData);
          await this.onMessage(parsedData);
        } catch (err) {
          const error = err as Error;
          console.error(`Error processing message: ${error.message}`);
          await this.handleError(error, payload);
        }
      },
    });
  }

  protected async parseMessage(data: Buffer): Promise<any> {
    try {
      return JSON.parse(data.toString());
    } catch (err) {
      const dlqProducer = kafkaWrapper.getProducer("system-metrics-dlq");
      await dlqProducer.publish({
        error: "Failed to parse message",
        original_message: data.toString(),
        topic: this.topic,
        timestamp: new Date().toISOString(),
      });
      throw new Error("Failed to parse message");
    }
  }

  protected async validateMessage(data: any): Promise<void> {
    if (!data?.metadata?.tenant_id) {
      const errorProducer = kafkaWrapper.getProducer("system-metrics-errors");
      await errorProducer.publish({
        error: "Missing tenant ID in message",
        original_payload: data,
        tenant_id: data?.metadata?.tenant_id || "unknown",
        timestamp: new Date().toISOString(),
      });
      throw new Error("Invalid message: Missing tenant ID");
    }
  }

  protected async handleError(
    error: Error,
    payload: EachMessagePayload
  ): Promise<void> {
    const dlqProducer = kafkaWrapper.getProducer("system-metrics-dlq");
    await dlqProducer.publish({
      error: error.message,
      original_message: payload.message.value?.toString() || "no message",
      topic: this.topic,
      timestamp: new Date().toISOString(),
    });
  }

  async disconnect(): Promise<void> {
    await this.consumer.disconnect();
  }
}
