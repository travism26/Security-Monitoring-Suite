interface SystemMetricsDLQ extends Event {
  timestamp: Date;
  data: {
    error_type: string;
    raw_message: any;
    failure_context: {
      tenant_id: string;
      timestamp: string;
      error_stack?: string;
    };
  };
}
