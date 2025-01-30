import { Event } from "../kafka/event";
import { Topics } from "../kafka/topics";

interface TenantContext {
  id: string;
  metadata?: { [key: string]: string };
}

interface ThreatIndicator {
  type: string;
  description: string;
  severity: string;
  score: number;
  timestamp: string;
  metadata: {
    tags: string[];
  };
  details?: { [key: string]: any };
}

interface HostInfo {
  os: string;
  arch: string;
  hostname: string;
  cpu_cores: number;
  go_version: string;
}

interface ProcessInfo {
  name: string;
  pid: number;
  cpu_percent: number;
  memory_usage: number;
  status: string;
}

interface SystemProcessStats {
  total_count: number;
  total_cpu_percent: number;
  total_memory_usage: number;
  process_list: ProcessInfo[];
}

interface MetadataInfo {
  collection_duration: string;
  collector_count: number;
  errors?: string[];
  tenant_metadata?: { [key: string]: string };
}

interface MetricData {
  host_info: HostInfo;
  metrics: { [key: string]: any };
  threat_indicators: ThreatIndicator[];
  metadata: MetadataInfo;
  processes: SystemProcessStats;
}

export interface SystemMetricsData {
  timestamp: string;
  tenant: TenantContext;
  data: MetricData;
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
