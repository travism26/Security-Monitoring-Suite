import { Event } from "./types";
import {
  generateId,
  getRandomEventType,
  getRandomEventSeverity,
  getRandomTimestamp,
  generateRandomIP,
  generateRandomHostname,
} from "./utils";

const eventMessages = {
  authentication: [
    "Failed login attempt",
    "Successful login",
    "Password reset requested",
    "Account locked",
    "New device authentication",
    "Session expired",
    "Multi-factor authentication failure",
  ],
  system: [
    "System startup",
    "Service restart",
    "High CPU usage detected",
    "Low disk space warning",
    "Memory threshold exceeded",
    "Process terminated unexpectedly",
    "System update available",
  ],
  network: [
    "High network latency detected",
    "Port scan detected",
    "Unusual outbound traffic",
    "Network interface down",
    "DNS resolution failure",
    "Bandwidth threshold exceeded",
    "New device connected",
  ],
  security: [
    "Malware detected",
    "Suspicious file activity",
    "Firewall rule violation",
    "SSL certificate expired",
    "Unauthorized access attempt",
    "Security patch available",
    "Intrusion attempt blocked",
  ],
  application: [
    "Application error",
    "Database connection failure",
    "API rate limit exceeded",
    "Cache invalidation",
    "Job queue overflow",
    "Configuration change",
    "Service dependency failure",
  ],
  database: [
    "Slow query detected",
    "Deadlock detected",
    "Connection pool exhausted",
    "Backup failure",
    "Replication lag",
    "Index rebuild required",
    "Storage threshold warning",
  ],
};

function generateEventDetails(type: Event["type"]): Record<string, unknown> {
  const baseDetails = {
    hostname: generateRandomHostname(),
    timestamp: new Date().toISOString(),
  };

  switch (type) {
    case "authentication":
      return {
        ...baseDetails,
        username: "user" + Math.floor(Math.random() * 1000),
        ipAddress: generateRandomIP(),
        userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
        method: Math.random() > 0.5 ? "password" : "sso",
      };

    case "system":
      return {
        ...baseDetails,
        cpuUsage: Math.floor(Math.random() * 100),
        memoryUsage: Math.floor(Math.random() * 100),
        diskUsage: Math.floor(Math.random() * 100),
        processId: Math.floor(Math.random() * 10000),
      };

    case "network":
      return {
        ...baseDetails,
        sourceIp: generateRandomIP(),
        destinationIp: generateRandomIP(),
        protocol: ["TCP", "UDP", "HTTP", "HTTPS"][
          Math.floor(Math.random() * 4)
        ],
        bytesTransferred: Math.floor(Math.random() * 1000000),
      };

    case "security":
      return {
        ...baseDetails,
        threatType: [
          "malware",
          "intrusion",
          "vulnerability",
          "policy_violation",
        ][Math.floor(Math.random() * 4)],
        sourceIp: generateRandomIP(),
        targetAsset: generateRandomHostname(),
        action: ["blocked", "detected", "quarantined"][
          Math.floor(Math.random() * 3)
        ],
      };

    case "application":
      return {
        ...baseDetails,
        applicationName: [
          "web-server",
          "api-gateway",
          "auth-service",
          "db-service",
        ][Math.floor(Math.random() * 4)],
        errorCode: Math.floor(Math.random() * 1000),
        stackTrace: "Error: Something went wrong\n  at Function.Module._load",
      };

    case "database":
      return {
        ...baseDetails,
        databaseName: ["users", "transactions", "analytics", "logs"][
          Math.floor(Math.random() * 4)
        ],
        queryId: generateId("query-"),
        executionTime: Math.floor(Math.random() * 10000),
        rowsAffected: Math.floor(Math.random() * 1000),
      };

    default:
      return baseDetails;
  }
}

export function generateMockEvent(): Event {
  const type = getRandomEventType();
  const severity = getRandomEventSeverity();
  const messages = eventMessages[type];
  const message = messages[Math.floor(Math.random() * messages.length)];

  return {
    id: generateId("evt-"),
    timestamp: getRandomTimestamp(),
    type,
    severity,
    source: `${type}-service`,
    message,
    details: generateEventDetails(type),
    metadata: {
      version: "1.0",
      environment: ["production", "staging", "development"][
        Math.floor(Math.random() * 3)
      ],
      datacenter: ["us-east", "us-west", "eu-central"][
        Math.floor(Math.random() * 3)
      ],
    },
  };
}

export function generateMockEvents(count: number = 50): Event[] {
  return Array.from({ length: count }, () => generateMockEvent()).sort(
    (a, b) => new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime()
  );
}
