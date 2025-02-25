import { useState, useEffect, useMemo } from "react";

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

interface UseNetworkTrafficOptions {
  useMockData?: boolean;
  refreshInterval?: number;
  timeRange?: "1h" | "24h" | "7d";
  dataPoints?: number;
}

// Utility function to format bytes
const formatBytes = (bytes: number): string => {
  if (bytes === 0) return "0 B";
  const k = 1024;
  const sizes = ["B", "KB", "MB", "GB", "TB"];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return `${parseFloat((bytes / Math.pow(k, i)).toFixed(2))} ${sizes[i]}`;
};

// Generate realistic mock data
const generateMockData = (
  dataPoints: number = 24,
  timeRange: UseNetworkTrafficOptions["timeRange"] = "24h"
): NetworkTrafficData[] => {
  const now = new Date();
  const hourInterval = timeRange === "1h" ? 1 / 60 : 1; // For 1h, generate per-minute data

  return Array.from({ length: dataPoints }, (_, i) => {
    // Create base traffic with daily patterns
    const hour = now.getHours() - (dataPoints - 1 - i) * hourInterval;
    const normalizedHour = ((hour % 24) + 24) % 24;

    // Simulate higher traffic during business hours (8-18)
    const businessHourMultiplier =
      normalizedHour >= 8 && normalizedHour <= 18 ? 2 : 1;

    // Add some randomness but keep it somewhat realistic
    const baseTraffic =
      Math.floor(Math.random() * 500) * businessHourMultiplier;

    return {
      hour: normalizedHour,
      inbound: baseTraffic + Math.floor(Math.random() * 200),
      outbound: baseTraffic * 0.7 + Math.floor(Math.random() * 150),
      timestamp: new Date(
        now.getTime() - (dataPoints - 1 - i) * hourInterval * 3600000
      ).toISOString(),
    };
  });
};

export function useNetworkTraffic(options: UseNetworkTrafficOptions = {}) {
  const {
    useMockData = process.env.NODE_ENV === "development",
    refreshInterval = 60000, // 1 minute
    timeRange = "24h",
    dataPoints = timeRange === "1h" ? 60 : 24, // 60 points for 1h, 24 for 24h
  } = options;

  const [trafficData, setTrafficData] = useState<NetworkTrafficData[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  // Fetch network traffic data
  useEffect(() => {
    let intervalId: NodeJS.Timeout;

    const fetchTrafficData = async () => {
      try {
        if (useMockData) {
          // Simulate API delay
          await new Promise((resolve) => setTimeout(resolve, 500));
          setTrafficData(generateMockData(dataPoints, timeRange));
        } else {
          // TODO: Implement real API integration
          // const response = await networkService.getTrafficData({
          //   timeRange,
          //   dataPoints,
          // });
          // setTrafficData(response.data);
          throw new Error("Real API integration not implemented yet");
        }
        setError(null);
      } catch (err) {
        setError(
          err instanceof Error
            ? err
            : new Error("Failed to fetch network traffic data")
        );
      } finally {
        setLoading(false);
      }
    };

    fetchTrafficData();

    if (refreshInterval > 0) {
      intervalId = setInterval(fetchTrafficData, refreshInterval);
    }

    return () => {
      if (intervalId) {
        clearInterval(intervalId);
      }
    };
  }, [useMockData, refreshInterval, timeRange, dataPoints]);

  // Calculate network statistics
  const stats = useMemo((): NetworkStats => {
    if (!trafficData.length) {
      return {
        totalInbound: 0,
        totalOutbound: 0,
        peakInbound: 0,
        peakOutbound: 0,
        averageInbound: 0,
        averageOutbound: 0,
      };
    }

    const totals = trafficData.reduce(
      (acc, curr) => ({
        inbound: acc.inbound + curr.inbound,
        outbound: acc.outbound + curr.outbound,
      }),
      { inbound: 0, outbound: 0 }
    );

    return {
      totalInbound: totals.inbound,
      totalOutbound: totals.outbound,
      peakInbound: Math.max(...trafficData.map((d) => d.inbound)),
      peakOutbound: Math.max(...trafficData.map((d) => d.outbound)),
      averageInbound: totals.inbound / trafficData.length,
      averageOutbound: totals.outbound / trafficData.length,
    };
  }, [trafficData]);

  // Format traffic data for display
  const formattedData = useMemo(() => {
    return trafficData.map((data) => ({
      ...data,
      inboundFormatted: formatBytes(data.inbound),
      outboundFormatted: formatBytes(data.outbound),
    }));
  }, [trafficData]);

  // Format stats for display
  const formattedStats = useMemo(() => {
    return {
      totalInbound: formatBytes(stats.totalInbound),
      totalOutbound: formatBytes(stats.totalOutbound),
      peakInbound: formatBytes(stats.peakInbound),
      peakOutbound: formatBytes(stats.peakOutbound),
      averageInbound: formatBytes(stats.averageInbound),
      averageOutbound: formatBytes(stats.averageOutbound),
    };
  }, [stats]);

  return {
    data: trafficData,
    formattedData,
    stats,
    formattedStats,
    loading,
    error,
    isMockData: useMockData,
    // Utility functions
    formatBytes,
  };
}
