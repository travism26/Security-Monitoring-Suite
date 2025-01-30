import { Kafka } from "kafkajs";
import { Publisher } from "../kafka/base-kafka-producer";
import { Topics } from "../kafka/topics";
import { Event } from "../kafka/event";

interface SystemMetricsDLQ extends Event {
  topic: Topics.SystemMetricsDLQ;
  data: {
    error: string;
    original_message: any;
    tenant_id: string;
    timestamp: string;
  };
}

/**
 * Publisher for handling unprocessable system metric messages
 * Publishes failed messages to a Dead Letter Queue for later analysis
 */
export class SystemMetricsDLQProducer extends Publisher<SystemMetricsDLQ> {
  readonly topic = Topics.SystemMetricsDLQ;

  constructor(client: Kafka) {
    super(client);
  }
}
