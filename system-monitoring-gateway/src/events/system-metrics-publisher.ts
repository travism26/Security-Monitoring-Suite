import { Publisher } from '../kafka/base-kafka-producer';
import { Topics } from '../kafka/topics';
import { SystemMetrics } from '../payload/system-metrics';

export class SystemMetricsPublisher extends Publisher<SystemMetrics> {
  topic: Topics.SystemMetrics = Topics.SystemMetrics;
}
