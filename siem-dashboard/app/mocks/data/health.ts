import { SystemHealth, SystemStatus, ResourceMetric } from "./types";
import {
  generateResourceMetric,
  getRandomSystemStatus,
  calculateUptime,
  generateRandomHostname,
} from "./utils";

const criticalServices = [
  "web-server",
  "api-gateway",
  "auth-service",
  "database",
  "cache",
  "message-queue",
  "monitoring",
  "logging",
];

function determineOverallStatus(metrics: ResourceMetric[]): SystemStatus {
  const statusPriority: Record<SystemStatus, number> = {
    critical: 3,
    warning: 2,
    unknown: 1,
    healthy: 0,
  };

  const highestPriority = metrics.reduce((highest, metric) => {
    const priority = statusPriority[metric.status];
    return priority > statusPriority[highest] ? metric.status : highest;
  }, "healthy" as SystemStatus);

  return highestPriority;
}

function generateServiceStartTimes(): Record<string, Date> {
  const now = new Date();
  return criticalServices.reduce((acc, service) => {
    // Random start time between 1 and 30 days ago
    const daysAgo = Math.floor(Math.random() * 29) + 1;
    acc[service] = new Date(now.getTime() - daysAgo * 24 * 60 * 60 * 1000);
    return acc;
  }, {} as Record<string, Date>);
}

const serviceStartTimes = generateServiceStartTimes();

function generateServiceStatuses() {
  return criticalServices.map((name) => {
    const status = Math.random() > 0.9 ? getRandomSystemStatus() : "healthy";
    return {
      name,
      status,
      uptime: calculateUptime(serviceStartTimes[name]),
      lastChecked: new Date().toISOString(),
    };
  });
}

export function generateMockSystemHealth(): SystemHealth {
  // Generate resource metrics
  const cpu = generateResourceMetric(30, 95); // CPU tends to be more volatile
  const memory = generateResourceMetric(40, 90); // Memory usage typically higher
  const disk = generateResourceMetric(50, 95); // Disk usage tends to grow
  const network = generateResourceMetric(20, 80); // Network usually more stable

  const services = generateServiceStatuses();

  // Determine overall status based on all metrics
  const status = determineOverallStatus([cpu, memory, disk, network]);

  return {
    status,
    cpu,
    memory,
    disk,
    network,
    services,
    lastUpdated: new Date().toISOString(),
  };
}

// Generate health data with specific issues
export function generateMockSystemHealthWithIssues(issues: {
  cpu?: boolean;
  memory?: boolean;
  disk?: boolean;
  network?: boolean;
  services?: string[];
}): SystemHealth {
  const health = generateMockSystemHealth();

  if (issues.cpu) {
    health.cpu = {
      ...health.cpu,
      usage: 95 + Math.random() * 5,
      status: "critical",
      trend: "up",
    };
  }

  if (issues.memory) {
    health.memory = {
      ...health.memory,
      usage: 92 + Math.random() * 8,
      status: "critical",
      trend: "up",
    };
  }

  if (issues.disk) {
    health.disk = {
      ...health.disk,
      usage: 97 + Math.random() * 3,
      status: "critical",
      trend: "up",
    };
  }

  if (issues.network) {
    health.network = {
      ...health.network,
      usage: 90 + Math.random() * 10,
      status: "warning",
      trend: "up",
    };
  }

  if (issues.services) {
    health.services = health.services.map((service) => {
      if (issues.services?.includes(service.name)) {
        return {
          ...service,
          status: Math.random() > 0.5 ? "critical" : "warning",
          uptime: Math.floor(Math.random() * 300), // Recent restart
        };
      }
      return service;
    });
  }

  // Recalculate overall status
  health.status = determineOverallStatus([
    health.cpu,
    health.memory,
    health.disk,
    health.network,
  ]);

  return health;
}

// Generate a series of health states over time
export function generateHealthTimeSeries(
  minutes: number = 60,
  intervalSeconds: number = 60,
  withIssues: boolean = false
): SystemHealth[] {
  const samples = (minutes * 60) / intervalSeconds;
  const now = new Date();

  return Array.from({ length: samples }, (_, i) => {
    const health =
      withIssues && i > samples / 2
        ? generateMockSystemHealthWithIssues({
            cpu: Math.random() > 0.7,
            memory: Math.random() > 0.8,
            disk: Math.random() > 0.9,
            services:
              Math.random() > 0.7
                ? [
                    criticalServices[
                      Math.floor(Math.random() * criticalServices.length)
                    ],
                  ]
                : undefined,
          })
        : generateMockSystemHealth();

    health.lastUpdated = new Date(
      now.getTime() - (samples - 1 - i) * intervalSeconds * 1000
    ).toISOString();

    return health;
  });
}
