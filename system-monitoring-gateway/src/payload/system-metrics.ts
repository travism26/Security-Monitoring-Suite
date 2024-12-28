import { Event } from '../kafka/event';

export interface SystemMetrics extends Event {
  timestamp: string;
  data: {
    host_info?: {
      os: string;
      arch: string;
      hostname: string;
      cpu_cores: number;
      go_version: string;
    };
    metrics: {
      [key: string]: any;
    };
    metadata: {
      collection_duration: string;
      collector_count: number;
      errors?: string[];
    };
    processes?: {
      total_count: number;
      total_cpu_percent: number;
      total_memory_usage: number;
      process_list: ProcessInfo[];
    };
  };
}

export interface ProcessInfo {
  name: string;
  pid: number;
  cpu_percent: number;
  memory_usage: number;
  status: string;
}
