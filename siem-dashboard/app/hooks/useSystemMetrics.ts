import { useState, useEffect } from "react";
import { SystemMetrics, ProcessInfo } from "../services/metricsService";

// Mock data that matches the expected API response structure
const mockSystemMetrics: SystemMetrics = {
  host: {
    os: "Linux",
    arch: "x64",
    hostname: "monitoring-server-01",
    cpuCores: 8,
  },
  metrics: {
    cpuUsage: 45.2,
    memoryTotal: 16384, // MB
    memoryUsed: 8192, // MB
    diskTotal: 512000, // MB
    diskUsed: 256000, // MB
    networkIn: 150, // MB/s
    networkOut: 75, // MB/s
  },
  processes: {
    totalCount: 128,
    totalCPUPercent: 65.5,
    totalMemoryUsage: 12288, // MB
    list: [
      {
        name: "nginx",
        pid: 1234,
        cpuPercent: 2.5,
        memoryUsage: 256,
        status: "running",
      },
      {
        name: "mongodb",
        pid: 1235,
        cpuPercent: 15.8,
        memoryUsage: 1024,
        status: "running",
      },
      {
        name: "node",
        pid: 1236,
        cpuPercent: 8.3,
        memoryUsage: 512,
        status: "running",
      },
    ] as ProcessInfo[],
  },
};

interface UseSystemMetricsOptions {
  useMockData?: boolean;
  refreshInterval?: number;
}

export function useSystemMetrics(options: UseSystemMetricsOptions = {}) {
  const {
    useMockData = process.env.NODE_ENV === "development",
    refreshInterval = 5000,
  } = options;

  const [metrics, setMetrics] = useState<SystemMetrics | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  useEffect(() => {
    let intervalId: NodeJS.Timeout;

    const fetchMetrics = async () => {
      try {
        if (useMockData) {
          // Simulate API delay
          await new Promise((resolve) => setTimeout(resolve, 500));
          setMetrics(mockSystemMetrics);
        } else {
          // TODO: Implement real API integration
          // const response = await metricsService.getSystemMetrics();
          // setMetrics(response);
          throw new Error("Real API integration not implemented yet");
        }
        setError(null);
      } catch (err) {
        setError(
          err instanceof Error ? err : new Error("Failed to fetch metrics")
        );
      } finally {
        setLoading(false);
      }
    };

    // Initial fetch
    fetchMetrics();

    // Set up polling interval
    if (refreshInterval > 0) {
      intervalId = setInterval(fetchMetrics, refreshInterval);
    }

    // Cleanup
    return () => {
      if (intervalId) {
        clearInterval(intervalId);
      }
    };
  }, [useMockData, refreshInterval]);

  // Utility functions for data transformation
  const getMemoryUsagePercentage = () => {
    if (!metrics) return 0;
    const { memoryTotal, memoryUsed } = metrics.metrics;
    return (memoryUsed / memoryTotal) * 100;
  };

  const getDiskUsagePercentage = () => {
    if (!metrics) return 0;
    const { diskTotal, diskUsed } = metrics.metrics;
    return (diskUsed / diskTotal) * 100;
  };

  return {
    metrics,
    loading,
    error,
    getMemoryUsagePercentage,
    getDiskUsagePercentage,
    // Expose whether we're using mock data for UI indication
    isMockData: useMockData,
  };
}
