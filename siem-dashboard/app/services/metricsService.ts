import axios from "axios";

export interface SystemMetrics {
  host: {
    os: string;
    arch: string;
    hostname: string;
    cpuCores: number;
  };
  metrics: {
    [key: string]: any;
  };
  processes: {
    totalCount: number;
    totalCPUPercent: number;
    totalMemoryUsage: number;
    list: ProcessInfo[];
  };
}

export interface ProcessInfo {
  name: string;
  pid: number;
  cpuPercent: number;
  memoryUsage: number;
  status: string;
}

class MetricsService {
  private baseUrl: string;

  constructor() {
    this.baseUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:3001";
  }

  async getSystemMetrics(): Promise<SystemMetrics> {
    try {
      const response = await axios.get(`${this.baseUrl}/api/v1/metrics/system`);
      return response.data;
    } catch (error) {
      console.error("Error fetching system metrics:", error);
      throw error;
    }
  }

  async getProcessMetrics(): Promise<ProcessInfo[]> {
    try {
      const response = await axios.get(
        `${this.baseUrl}/api/v1/metrics/processes`
      );
      return response.data;
    } catch (error) {
      console.error("Error fetching process metrics:", error);
      throw error;
    }
  }

  // Utility functions for formatting metrics
  static formatBytes(bytes: number): string {
    if (bytes === 0) return "0 B";
    const k = 1024;
    const sizes = ["B", "KB", "MB", "GB", "TB"];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return `${parseFloat((bytes / Math.pow(k, i)).toFixed(2))} ${sizes[i]}`;
  }

  static formatPercentage(value: number): string {
    return `${value.toFixed(1)}%`;
  }
}

export const metricsService = new MetricsService();
