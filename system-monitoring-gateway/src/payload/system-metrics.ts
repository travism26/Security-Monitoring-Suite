import { Event } from '../kafka/event';

export interface SystemMetrics extends Event {
  timestamp: string;
  data: {
    host_info?: {
      os: string;
      arch: string;
    };
    metrics: {
      cpu_usage?: number;
      memory_usage?: number;
      memory_total?: number;
      memory_used_percent?: number;
      disk?: {
        total: number;
        used: number;
        free: number;
        used_percent: number;
      };
      network?: {
        bytes_sent: number;
        bytes_received: number;
      };
      processes?: {
        total_count: number;
        total_cpu_percent: number;
        total_memory_usage: number;
        process_list: Array<{
          pid: number;
          name: string;
          cpu_percent: number;
          memory_usage: number;
          status: string;
        }>;
      };
    };
    threat_indicators?: Array<{
      type: string;
      description: string;
      severity: string;
      score: number;
      metadata?: Record<string, unknown>;
    }>;
  };
}
