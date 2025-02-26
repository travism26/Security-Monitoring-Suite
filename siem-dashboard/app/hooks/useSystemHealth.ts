import { useState, useEffect } from "react";
import { SystemMetrics } from "../services/metricsService";

export type HealthStatus = "healthy" | "warning" | "critical" | "unknown";

export interface SystemHealth {
  cpu: {
    status: HealthStatus;
    usage: number;
    threshold: {
      warning: number;
      critical: number;
    };
  };
  memory: {
    status: HealthStatus;
    usage: number;
    threshold: {
      warning: number;
      critical: number;
    };
  };
  disk: {
    status: HealthStatus;
    usage: number;
    threshold: {
      warning: number;
      critical: number;
    };
  };
  network: {
    status: HealthStatus;
    inboundUtilization: number;
    outboundUtilization: number;
    threshold: {
      warning: number;
      critical: number;
    };
  };
  overall: HealthStatus;
  lastUpdated: string;
}

const mockSystemHealth: SystemHealth = {
  cpu: {
    status: "healthy",
    usage: 45.2,
    threshold: {
      warning: 70,
      critical: 90,
    },
  },
  memory: {
    status: "warning",
    usage: 75.5,
    threshold: {
      warning: 70,
      critical: 85,
    },
  },
  disk: {
    status: "healthy",
    usage: 50.0,
    threshold: {
      warning: 80,
      critical: 90,
    },
  },
  network: {
    status: "healthy",
    inboundUtilization: 45,
    outboundUtilization: 30,
    threshold: {
      warning: 80,
      critical: 90,
    },
  },
  overall: "warning",
  lastUpdated: new Date().toISOString(),
};

interface UseSystemHealthOptions {
  useMockData?: boolean;
  refreshInterval?: number;
}

export function useSystemHealth(options: UseSystemHealthOptions = {}) {
  const {
    useMockData = process.env.NODE_ENV === "development",
    refreshInterval = 5000,
  } = options;

  const [health, setHealth] = useState<SystemHealth | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  // Calculate health status based on value and thresholds
  const calculateStatus = (
    value: number,
    warningThreshold: number,
    criticalThreshold: number
  ): HealthStatus => {
    if (value >= criticalThreshold) return "critical";
    if (value >= warningThreshold) return "warning";
    return "healthy";
  };

  // Calculate overall system health based on individual components
  const calculateOverallHealth = (metrics: SystemMetrics): HealthStatus => {
    const statuses: HealthStatus[] = [];

    // CPU Health
    const cpuStatus = calculateStatus(
      metrics.metrics.cpuUsage,
      mockSystemHealth.cpu.threshold.warning,
      mockSystemHealth.cpu.threshold.critical
    );
    statuses.push(cpuStatus);

    // Memory Health
    const memoryUsage =
      (metrics.metrics.memoryUsed / metrics.metrics.memoryTotal) * 100;
    const memoryStatus = calculateStatus(
      memoryUsage,
      mockSystemHealth.memory.threshold.warning,
      mockSystemHealth.memory.threshold.critical
    );
    statuses.push(memoryStatus);

    // Disk Health
    const diskUsage =
      (metrics.metrics.diskUsed / metrics.metrics.diskTotal) * 100;
    const diskStatus = calculateStatus(
      diskUsage,
      mockSystemHealth.disk.threshold.warning,
      mockSystemHealth.disk.threshold.critical
    );
    statuses.push(diskStatus);

    // Determine overall status
    if (statuses.includes("critical")) return "critical";
    if (statuses.includes("warning")) return "warning";
    return "healthy";
  };

  useEffect(() => {
    let intervalId: NodeJS.Timeout;

    const updateHealth = async () => {
      try {
        if (useMockData) {
          // Simulate API delay
          await new Promise((resolve) => setTimeout(resolve, 500));
          setHealth({
            ...mockSystemHealth,
            lastUpdated: new Date().toISOString(),
          });
        } else {
          // TODO: Implement real API integration
          // const metrics = await metricsService.getSystemMetrics();
          // const health = calculateHealthFromMetrics(metrics);
          // setHealth(health);
          throw new Error("Real API integration not implemented yet");
        }
        setError(null);
      } catch (err) {
        setError(
          err instanceof Error
            ? err
            : new Error("Failed to fetch health status")
        );
      } finally {
        setLoading(false);
      }
    };

    // Initial update
    updateHealth();

    // Set up polling interval
    if (refreshInterval > 0) {
      intervalId = setInterval(updateHealth, refreshInterval);
    }

    // Cleanup
    return () => {
      if (intervalId) {
        clearInterval(intervalId);
      }
    };
  }, [useMockData, refreshInterval]);

  // Utility functions for status checks
  const isHealthy = (status: HealthStatus) => status === "healthy";
  const isWarning = (status: HealthStatus) => status === "warning";
  const isCritical = (status: HealthStatus) => status === "critical";

  return {
    health,
    loading,
    error,
    isHealthy,
    isWarning,
    isCritical,
    isMockData: useMockData,
  };
}
