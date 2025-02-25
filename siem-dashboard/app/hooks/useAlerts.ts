import { useState, useEffect, useMemo } from "react";

export interface Alert {
  id: string;
  timestamp: string;
  title: string;
  description: string;
  severity: "critical" | "high" | "medium" | "low";
  category: AlertCategory;
  source: string;
  status:
    | "new"
    | "acknowledged"
    | "investigating"
    | "resolved"
    | "false_positive";
  assignedTo?: string;
  relatedEvents?: string[];
  metadata?: Record<string, unknown>;
}

export type AlertCategory =
  | "security"
  | "system"
  | "network"
  | "application"
  | "database"
  | "compliance";

export interface AlertFilters {
  severity?: Alert["severity"];
  category?: AlertCategory;
  status?: Alert["status"];
  startDate?: Date;
  endDate?: Date;
  source?: string;
  assignedTo?: string;
}

interface UseAlertsOptions {
  useMockData?: boolean;
  refreshInterval?: number;
  enableWebSocket?: boolean;
}

// Mock data with realistic security alerts
const mockAlerts: Alert[] = [
  {
    id: "alert-001",
    timestamp: new Date().toISOString(),
    title: "Suspicious Login Activity",
    description:
      "Multiple failed login attempts detected from IP 192.168.1.100",
    severity: "high",
    category: "security",
    source: "auth-service",
    status: "new",
    relatedEvents: ["event-123", "event-124"],
    metadata: {
      ip: "192.168.1.100",
      attempts: 5,
      timeWindow: "5 minutes",
    },
  },
  {
    id: "alert-002",
    timestamp: new Date(Date.now() - 300000).toISOString(), // 5 minutes ago
    title: "High CPU Usage",
    description: "System CPU usage exceeded 90% for 5 minutes",
    severity: "medium",
    category: "system",
    source: "system-monitor",
    status: "investigating",
    assignedTo: "admin",
    metadata: {
      cpuUsage: 92,
      duration: "5 minutes",
      affectedServices: ["web-server", "database"],
    },
  },
  {
    id: "alert-003",
    timestamp: new Date(Date.now() - 600000).toISOString(), // 10 minutes ago
    title: "Potential Data Exfiltration",
    description: "Unusual outbound data transfer detected",
    severity: "critical",
    category: "network",
    source: "network-monitor",
    status: "investigating",
    relatedEvents: ["event-125"],
    metadata: {
      dataVolume: "500MB",
      destination: "unknown-endpoint",
      protocol: "HTTPS",
    },
  },
];

export function useAlerts(options: UseAlertsOptions = {}) {
  const {
    useMockData = process.env.NODE_ENV === "development",
    refreshInterval = 30000, // 30 seconds
    enableWebSocket = false,
  } = options;

  const [alerts, setAlerts] = useState<Alert[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);
  const [filters, setFilters] = useState<AlertFilters>({});
  const [sortConfig, setSortConfig] = useState<{
    field: keyof Alert;
    direction: "asc" | "desc";
  }>({ field: "timestamp", direction: "desc" });

  // Fetch alerts with optional filtering
  useEffect(() => {
    let intervalId: NodeJS.Timeout;

    const fetchAlerts = async () => {
      try {
        if (useMockData) {
          // Simulate API delay
          await new Promise((resolve) => setTimeout(resolve, 500));
          setAlerts(mockAlerts);
        } else {
          // TODO: Implement real API integration
          // const response = await alertService.getAlerts(filters);
          // setAlerts(response.alerts);
          throw new Error("Real API integration not implemented yet");
        }
        setError(null);
      } catch (err) {
        setError(
          err instanceof Error ? err : new Error("Failed to fetch alerts")
        );
      } finally {
        setLoading(false);
      }
    };

    fetchAlerts();

    if (refreshInterval > 0 && !enableWebSocket) {
      intervalId = setInterval(fetchAlerts, refreshInterval);
    }

    return () => {
      if (intervalId) {
        clearInterval(intervalId);
      }
    };
  }, [useMockData, refreshInterval, enableWebSocket, filters]);

  // WebSocket connection for real-time updates
  useEffect(() => {
    if (!enableWebSocket || useMockData) return;

    // TODO: Implement WebSocket connection
    // const ws = new WebSocket(WEBSOCKET_URL);
    //
    // ws.onmessage = (event) => {
    //   const newAlert = JSON.parse(event.data);
    //   setAlerts(prev => [newAlert, ...prev]);
    // };
    //
    // return () => ws.close();
  }, [enableWebSocket, useMockData]);

  // Filter and sort alerts
  const processedAlerts = useMemo(() => {
    let filtered = alerts.filter((alert) => {
      if (filters.severity && alert.severity !== filters.severity) return false;
      if (filters.category && alert.category !== filters.category) return false;
      if (filters.status && alert.status !== filters.status) return false;
      if (filters.source && !alert.source.includes(filters.source))
        return false;
      if (filters.assignedTo && alert.assignedTo !== filters.assignedTo)
        return false;
      if (filters.startDate && new Date(alert.timestamp) < filters.startDate)
        return false;
      if (filters.endDate && new Date(alert.timestamp) > filters.endDate)
        return false;
      return true;
    });

    // Sort alerts
    filtered.sort((a, b) => {
      const aValue = String(a[sortConfig.field] ?? "");
      const bValue = String(b[sortConfig.field] ?? "");
      return sortConfig.direction === "asc"
        ? aValue.localeCompare(bValue)
        : bValue.localeCompare(aValue);
    });

    return filtered;
  }, [alerts, filters, sortConfig]);

  // Alert management functions
  const updateAlertStatus = async (
    alertId: string,
    status: Alert["status"],
    assignedTo?: string
  ) => {
    try {
      if (useMockData) {
        setAlerts((prev) =>
          prev.map((alert) =>
            alert.id === alertId
              ? {
                  ...alert,
                  status,
                  assignedTo,
                  timestamp: new Date().toISOString(),
                }
              : alert
          )
        );
      } else {
        // TODO: Implement real API integration
        // await alertService.updateAlertStatus(alertId, status, assignedTo);
        throw new Error("Real API integration not implemented yet");
      }
    } catch (err) {
      setError(
        err instanceof Error ? err : new Error("Failed to update alert status")
      );
      throw err;
    }
  };

  const updateFilters = (newFilters: Partial<AlertFilters>) => {
    setFilters((prev) => ({ ...prev, ...newFilters }));
  };

  const clearFilters = () => {
    setFilters({});
  };

  const updateSort = (field: keyof Alert) => {
    setSortConfig((prev) => ({
      field,
      direction:
        prev.field === field && prev.direction === "asc" ? "desc" : "asc",
    }));
  };

  return {
    alerts: processedAlerts,
    loading,
    error,
    filters,
    sortConfig,
    updateAlertStatus,
    updateFilters,
    clearFilters,
    updateSort,
    isMockData: useMockData,
  };
}
