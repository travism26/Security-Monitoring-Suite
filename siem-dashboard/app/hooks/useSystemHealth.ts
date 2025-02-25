import { useState, useEffect, useMemo } from "react";

export type SystemStatus =
  | "operational"
  | "degraded"
  | "warning"
  | "critical"
  | "down";
export type SystemType =
  | "firewall"
  | "ids"
  | "log_server"
  | "email_filter"
  | "database"
  | "api_gateway"
  | "monitoring";

export interface SystemMetrics {
  cpu: number;
  memory: number;
  disk: number;
  latency: number;
  uptime: number;
  lastChecked: string;
}

export interface SystemHealthData {
  id: string;
  name: string;
  type: SystemType;
  status: SystemStatus;
  metrics: SystemMetrics;
  lastIncident?: string;
  message?: string;
}

interface HealthThresholds {
  cpu: {
    warning: number;
    critical: number;
  };
  memory: {
    warning: number;
    critical: number;
  };
  disk: {
    warning: number;
    critical: number;
  };
  latency: {
    warning: number; // ms
    critical: number; // ms
  };
}

interface UseSystemHealthOptions {
  useMockData?: boolean;
  refreshInterval?: number;
  customThresholds?: Partial<HealthThresholds>;
}

// Default thresholds for system health monitoring
const DEFAULT_THRESHOLDS: HealthThresholds = {
  cpu: {
    warning: 70,
    critical: 90,
  },
  memory: {
    warning: 80,
    critical: 95,
  },
  disk: {
    warning: 85,
    critical: 95,
  },
  latency: {
    warning: 1000, // 1 second
    critical: 5000, // 5 seconds
  },
};

// Generate realistic mock data for a system
const generateMockSystemData = (
  type: SystemType,
  forceStatus?: SystemStatus
): SystemHealthData => {
  const now = new Date();
  const isSimulatedIssue = forceStatus ? true : Math.random() < 0.1; // 10% chance of issue

  // Base metrics with some randomization
  const baseMetrics = {
    cpu: Math.floor(Math.random() * 40) + (isSimulatedIssue ? 55 : 0),
    memory: Math.floor(Math.random() * 30) + (isSimulatedIssue ? 65 : 0),
    disk: Math.floor(Math.random() * 20) + (isSimulatedIssue ? 75 : 0),
    latency: Math.floor(Math.random() * 500) + (isSimulatedIssue ? 2000 : 0),
    uptime: Math.floor(Math.random() * 30 * 24 * 60 * 60), // Random uptime up to 30 days
    lastChecked: now.toISOString(),
  };

  // System-specific adjustments
  const systemSpecifics: Record<SystemType, Partial<SystemMetrics>> = {
    firewall: { cpu: baseMetrics.cpu * 1.2 }, // Firewalls tend to be CPU-intensive
    ids: { memory: baseMetrics.memory * 1.3 }, // IDS systems are memory-intensive
    log_server: { disk: baseMetrics.disk * 1.4 }, // Log servers use more disk
    email_filter: { latency: baseMetrics.latency * 0.8 }, // Email filters should be fast
    database: {
      memory: baseMetrics.memory * 1.2,
      disk: baseMetrics.disk * 1.2,
    },
    api_gateway: { latency: baseMetrics.latency * 0.7 },
    monitoring: { cpu: baseMetrics.cpu * 0.8 }, // Monitoring should be light
  };

  const metrics = {
    ...baseMetrics,
    ...systemSpecifics[type],
  };

  return {
    id: `${type}-1`,
    name: type
      .split("_")
      .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
      .join(" "),
    type,
    status: forceStatus || "operational",
    metrics,
    lastIncident: isSimulatedIssue
      ? new Date(now.getTime() - 1000 * 60 * 15).toISOString()
      : undefined,
    message: isSimulatedIssue ? "Performance degradation detected" : undefined,
  };
};

// Generate mock data for all systems
const generateMockData = (): SystemHealthData[] => {
  return [
    generateMockSystemData("firewall"),
    generateMockSystemData("ids"),
    generateMockSystemData(
      "log_server",
      Math.random() < 0.2 ? "down" : undefined
    ), // 20% chance of being down
    generateMockSystemData("email_filter"),
    generateMockSystemData("database"),
    generateMockSystemData("api_gateway"),
    generateMockSystemData("monitoring"),
  ];
};

export function useSystemHealth(options: UseSystemHealthOptions = {}) {
  const {
    useMockData = process.env.NODE_ENV === "development",
    refreshInterval = 30000, // 30 seconds
    customThresholds = {},
  } = options;

  const [systems, setSystems] = useState<SystemHealthData[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  // Merge custom thresholds with defaults
  const thresholds = useMemo(
    () => ({
      ...DEFAULT_THRESHOLDS,
      ...customThresholds,
    }),
    [customThresholds]
  );

  // Calculate system status based on metrics and thresholds
  const calculateStatus = (metrics: SystemMetrics): SystemStatus => {
    if (
      metrics.cpu >= thresholds.cpu.critical ||
      metrics.memory >= thresholds.memory.critical ||
      metrics.disk >= thresholds.disk.critical ||
      metrics.latency >= thresholds.latency.critical
    ) {
      return "critical";
    }

    if (
      metrics.cpu >= thresholds.cpu.warning ||
      metrics.memory >= thresholds.memory.warning ||
      metrics.disk >= thresholds.disk.warning ||
      metrics.latency >= thresholds.latency.warning
    ) {
      return "warning";
    }

    if (
      metrics.cpu >= thresholds.cpu.warning * 0.8 ||
      metrics.memory >= thresholds.memory.warning * 0.8 ||
      metrics.disk >= thresholds.disk.warning * 0.8 ||
      metrics.latency >= thresholds.latency.warning * 0.8
    ) {
      return "degraded";
    }

    return "operational";
  };

  // Fetch system health data
  useEffect(() => {
    let intervalId: NodeJS.Timeout;

    const fetchHealthData = async () => {
      try {
        if (useMockData) {
          // Simulate API delay
          await new Promise((resolve) => setTimeout(resolve, 500));
          const mockData = generateMockData();

          // Update status based on metrics
          const updatedData = mockData.map((system) => ({
            ...system,
            status:
              system.status === "down"
                ? "down"
                : calculateStatus(system.metrics),
          }));

          setSystems(updatedData);
        } else {
          // TODO: Implement real API integration
          // const response = await healthService.getSystemHealth();
          // setSystems(response.systems);
          throw new Error("Real API integration not implemented yet");
        }
        setError(null);
      } catch (err) {
        setError(
          err instanceof Error
            ? err
            : new Error("Failed to fetch system health data")
        );
      } finally {
        setLoading(false);
      }
    };

    fetchHealthData();

    if (refreshInterval > 0) {
      intervalId = setInterval(fetchHealthData, refreshInterval);
    }

    return () => {
      if (intervalId) {
        clearInterval(intervalId);
      }
    };
  }, [useMockData, refreshInterval, thresholds]);

  // Calculate overall system health
  const overallHealth = useMemo(() => {
    if (!systems.length) return "unknown";

    const statusCounts = systems.reduce(
      (acc, system) => ({
        ...acc,
        [system.status]: (acc[system.status] || 0) + 1,
      }),
      {} as Record<SystemStatus, number>
    );

    if (statusCounts.down > 0) return "critical";
    if (statusCounts.critical > 0) return "critical";
    if (statusCounts.warning > 0) return "warning";
    if (statusCounts.degraded > 0) return "degraded";
    return "operational";
  }, [systems]);

  // Format uptime
  const formatUptime = (seconds: number): string => {
    const days = Math.floor(seconds / (24 * 60 * 60));
    const hours = Math.floor((seconds % (24 * 60 * 60)) / (60 * 60));
    const minutes = Math.floor((seconds % (60 * 60)) / 60);

    if (days > 0) return `${days}d ${hours}h`;
    if (hours > 0) return `${hours}h ${minutes}m`;
    return `${minutes}m`;
  };

  return {
    systems,
    loading,
    error,
    overallHealth,
    thresholds,
    isMockData: useMockData,
    // Utility functions
    formatUptime,
    calculateStatus,
  };
}
