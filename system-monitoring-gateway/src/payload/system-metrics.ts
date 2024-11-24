export interface SystemMetrics {
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
