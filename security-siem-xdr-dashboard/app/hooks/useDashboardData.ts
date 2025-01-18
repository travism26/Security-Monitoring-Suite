import { useState, useEffect } from "react";
import { apiService } from "../services/api";

interface Event {
  id: number;
  timestamp: string;
  type: string;
  source: string;
  description: string;
}

interface Threat {
  name: string;
  count: number;
  severity: number;
}

interface SystemHealth {
  name: string;
  status: string;
}

interface Alert {
  id: number;
  message: string;
  timestamp: string;
}

export function useDashboardData() {
  const [events, setEvents] = useState<Event[]>([]);
  const [threats, setThreats] = useState<Threat[]>([]);
  const [systemHealth, setSystemHealth] = useState<SystemHealth[]>([]);
  const [alerts, setAlerts] = useState<Alert[]>([]);

  useEffect(() => {
    fetchDashboardData();
  }, []);

  const fetchDashboardData = async () => {
    try {
      const [
        fetchedEvents,
        fetchedThreats,
        fetchedSystemHealth,
        fetchedAlerts,
      ] = await Promise.all([
        apiService.getEvents(),
        apiService.getThreats(),
        apiService.getSystemHealth(),
        apiService.getAlerts(),
      ]);

      setEvents(fetchedEvents);
      setThreats(fetchedThreats);
      setSystemHealth(fetchedSystemHealth);
      setAlerts(fetchedAlerts);
    } catch (error) {
      console.error("Failed to fetch dashboard data:", error);
    }
  };

  return {
    events,
    threats,
    systemHealth,
    alerts,
    refreshData: fetchDashboardData,
  };
}
