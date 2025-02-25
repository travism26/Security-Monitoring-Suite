import { useState, useEffect, useMemo } from "react";

export interface SecurityEvent {
  id: number;
  timestamp: string;
  type: EventType;
  source: string;
  description: string;
  severity?: "low" | "medium" | "high";
  status?: "new" | "investigating" | "resolved";
}

export type EventType =
  | "Login"
  | "File Access"
  | "Network"
  | "System"
  | "Security"
  | "Database";

export interface EventFilters {
  type?: EventType;
  severity?: SecurityEvent["severity"];
  status?: SecurityEvent["status"];
  startDate?: Date;
  endDate?: Date;
  source?: string;
}

interface UseEventLogOptions {
  useMockData?: boolean;
  pageSize?: number;
  refreshInterval?: number;
}

// Extended mock data with more realistic security events
const mockEvents: SecurityEvent[] = [
  {
    id: 1,
    timestamp: "2024-02-25 10:30:15",
    type: "Login",
    source: "192.168.1.100",
    description: "Successful login attempt",
    severity: "low",
    status: "resolved",
  },
  {
    id: 2,
    timestamp: "2024-02-25 10:35:22",
    type: "File Access",
    source: "192.168.1.101",
    description: "Unauthorized file access attempt",
    severity: "high",
    status: "investigating",
  },
  {
    id: 3,
    timestamp: "2024-02-25 10:40:05",
    type: "Network",
    source: "192.168.1.102",
    description: "Unusual outbound traffic detected",
    severity: "medium",
    status: "new",
  },
  {
    id: 4,
    timestamp: "2024-02-25 10:45:30",
    type: "Security",
    source: "192.168.1.103",
    description: "Failed SSH authentication attempts",
    severity: "high",
    status: "investigating",
  },
  {
    id: 5,
    timestamp: "2024-02-25 10:50:15",
    type: "Database",
    source: "192.168.1.104",
    description: "Database connection error",
    severity: "medium",
    status: "resolved",
  },
];

export function useEventLog(options: UseEventLogOptions = {}) {
  const {
    useMockData = process.env.NODE_ENV === "development",
    pageSize = 10,
    refreshInterval = 30000, // 30 seconds
  } = options;

  const [events, setEvents] = useState<SecurityEvent[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);
  const [filters, setFilters] = useState<EventFilters>({});
  const [currentPage, setCurrentPage] = useState(1);

  // Fetch events with optional filtering
  useEffect(() => {
    let intervalId: NodeJS.Timeout;

    const fetchEvents = async () => {
      try {
        if (useMockData) {
          // Simulate API delay
          await new Promise((resolve) => setTimeout(resolve, 500));
          setEvents(mockEvents);
        } else {
          // TODO: Implement real API integration
          // const response = await eventService.getEvents({
          //   page: currentPage,
          //   pageSize,
          //   ...filters
          // });
          // setEvents(response.events);
          throw new Error("Real API integration not implemented yet");
        }
        setError(null);
      } catch (err) {
        setError(
          err instanceof Error ? err : new Error("Failed to fetch events")
        );
      } finally {
        setLoading(false);
      }
    };

    fetchEvents();

    if (refreshInterval > 0) {
      intervalId = setInterval(fetchEvents, refreshInterval);
    }

    return () => {
      if (intervalId) {
        clearInterval(intervalId);
      }
    };
  }, [useMockData, currentPage, pageSize, refreshInterval, filters]);

  // Filter events based on current filters
  const filteredEvents = useMemo(() => {
    return events.filter((event) => {
      if (filters.type && event.type !== filters.type) return false;
      if (filters.severity && event.severity !== filters.severity) return false;
      if (filters.status && event.status !== filters.status) return false;
      if (filters.source && !event.source.includes(filters.source))
        return false;
      if (filters.startDate && new Date(event.timestamp) < filters.startDate)
        return false;
      if (filters.endDate && new Date(event.timestamp) > filters.endDate)
        return false;
      return true;
    });
  }, [events, filters]);

  // Calculate pagination
  const totalPages = Math.ceil(filteredEvents.length / pageSize);
  const paginatedEvents = filteredEvents.slice(
    (currentPage - 1) * pageSize,
    currentPage * pageSize
  );

  // Utility functions
  const updateFilters = (newFilters: Partial<EventFilters>) => {
    setFilters((prev) => ({ ...prev, ...newFilters }));
    setCurrentPage(1); // Reset to first page when filters change
  };

  const clearFilters = () => {
    setFilters({});
    setCurrentPage(1);
  };

  const goToPage = (page: number) => {
    const targetPage = Math.max(1, Math.min(page, totalPages));
    setCurrentPage(targetPage);
  };

  return {
    events: paginatedEvents,
    loading,
    error,
    filters,
    currentPage,
    totalPages,
    pageSize,
    updateFilters,
    clearFilters,
    goToPage,
    isMockData: useMockData,
  };
}
