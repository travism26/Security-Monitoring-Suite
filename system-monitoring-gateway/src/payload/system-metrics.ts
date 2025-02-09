import { Event } from "../kafka/event";
import { Topics } from "../kafka/topics";

interface ProcessInfo {
  name: string;
  pid: number;
  cpu_percent: number;
  memory_usage: number;
  status: string;
}

interface ThreatIndicator {
  type: string;
  description: string;
  severity: string;
  score: number;
  timestamp: string;
  tags: string[];
  details?: { [key: string]: any };
}

export interface MetricsPayload {
  timestamp: string;
  tenant_id: string;
  api_key: string;
  tenant_metadata?: { [key: string]: string };

  host: {
    os: string;
    arch: string;
    hostname: string;
    cpu_cores: number;
    go_version: string;
  };

  metrics: { [key: string]: any };

  processes: {
    total_count: number;
    total_cpu_percent: number;
    total_memory_usage: number;
    list: ProcessInfo[];
  };

  threat_indicators?: ThreatIndicator[];

  metadata: {
    collection_duration: string;
    collector_count: number;
    errors?: string[];
  };
}

export interface SystemMetricsData {
  data: MetricsPayload;
}

export interface SystemMetrics extends Event<SystemMetricsData> {
  timestamp: string;
}

export type SystemMetricsPayload =
  | SystemMetrics
  | {
      topic: Topics.SystemMetrics;
      data: "invalid-json-structure" | string;
      timestamp: string;
    };
