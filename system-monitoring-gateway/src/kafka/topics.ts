export enum Topics {
  SystemMetrics = "system-metrics",
  SystemMetricsErrors = "system-metrics-errors",
  SystemMetricsDLQ = "system-metrics-dlq",
}

export const tenantTopic = (baseTopic: Topics, tenantId: string) =>
  `${baseTopic}.${tenantId}`;
