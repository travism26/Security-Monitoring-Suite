import { Kafka } from "kafkajs";
import { Publisher } from "../kafka/base-kafka-producer";
import { Topics } from "../kafka/topics";
import { SystemMetrics } from "../payload/system-metrics";

export class SystemMetricsPublisher extends Publisher<SystemMetrics> {
  readonly topic = Topics.SystemMetrics;

  constructor(client: Kafka) {
    super(client);
  }
}

interface SystemMetricsErrors extends Event {
  topic: Topics.SystemMetricsErrors;
  data: {
    error: string;
    original_payload: any;
    metadata: {
      tenant_id: string;
      timestamp: string;
    };
  };
}

export class SystemMetricsErrorProducer extends Publisher<SystemMetricsErrors> {
  readonly topic = Topics.SystemMetricsErrors;

  constructor(client: Kafka) {
    super(client);
  }
}

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
