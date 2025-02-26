// Event Types
export interface Event {
  id: string;
  timestamp: string;
  type: EventType;
  severity: EventSeverity;
  source: string;
  message: string;
  details: Record<string, unknown>;
  metadata?: Record<string, unknown>;
}

export type EventType =
  | "authentication"
  | "system"
  | "network"
  | "security"
  | "application"
  | "database";

export type EventSeverity = "info" | "warning" | "error" | "critical";

// Alert Types
export interface Alert {
  id: string;
  timestamp: string;
  title: string;
  description: string;
  severity: AlertSeverity;
  category: AlertCategory;
  source: string;
  status: AlertStatus;
  assignedTo?: string;
  relatedEvents?: string[];
  metadata?: Record<string, unknown>;
}

export type AlertSeverity = "critical" | "high" | "medium" | "low";

export type AlertCategory =
  | "security"
  | "system"
  | "network"
  | "application"
  | "database"
  | "compliance";

export type AlertStatus =
  | "new"
  | "acknowledged"
  | "investigating"
  | "resolved"
  | "false_positive";

// Network Traffic Types
export interface NetworkTrafficData {
  hour: number;
  inbound: number;
  outbound: number;
  timestamp: string;
}

export interface NetworkStats {
  totalInbound: number;
  totalOutbound: number;
  peakInbound: number;
  peakOutbound: number;
  averageInbound: number;
  averageOutbound: number;
}

// System Health Types
export interface SystemHealth {
  status: SystemStatus;
  cpu: ResourceMetric;
  memory: ResourceMetric;
  disk: ResourceMetric;
  network: ResourceMetric;
  services: ServiceStatus[];
  lastUpdated: string;
}

export type SystemStatus = "healthy" | "warning" | "critical" | "unknown";

export interface ResourceMetric {
  usage: number;
  limit: number;
  status: SystemStatus;
  trend: "up" | "down" | "stable";
}

export interface ServiceStatus {
  name: string;
  status: SystemStatus;
  uptime: number;
  lastChecked: string;
}

// System Metrics Types
export interface SystemMetrics {
  cpu: CPUMetrics;
  memory: MemoryMetrics;
  disk: DiskMetrics;
  network: NetworkMetrics;
  timestamp: string;
}

export interface CPUMetrics {
  usage: number;
  temperature: number;
  processes: number;
  loadAverage: number[];
}

export interface MemoryMetrics {
  total: number;
  used: number;
  free: number;
  cached: number;
  swapUsed: number;
  swapTotal: number;
}

export interface DiskMetrics {
  total: number;
  used: number;
  free: number;
  readRate: number;
  writeRate: number;
  partitions: DiskPartition[];
}

export interface DiskPartition {
  mount: string;
  total: number;
  used: number;
  free: number;
}

export interface NetworkMetrics {
  bytesIn: number;
  bytesOut: number;
  packetsIn: number;
  packetsOut: number;
  errors: number;
  dropped: number;
}

// Threat Summary Types
export interface ThreatSummary {
  totalThreats: number;
  criticalThreats: number;
  highThreats: number;
  mediumThreats: number;
  lowThreats: number;
  byCategory: Record<string, number>;
  recentIncidents: ThreatIncident[];
  timestamp: string;
}

export interface ThreatIncident {
  id: string;
  type: string;
  severity: AlertSeverity;
  timestamp: string;
  description: string;
  status: string;
}
