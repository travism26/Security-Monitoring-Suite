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
  };
}
