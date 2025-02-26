import { Alert, AlertCategory, AlertSeverity } from "./types";
import {
  generateId,
  getRandomAlertCategory,
  getRandomAlertSeverity,
  getRandomAlertStatus,
  getRandomTimestamp,
  generateRandomIP,
  generateRandomHostname,
} from "./utils";

const alertTitles: Record<AlertCategory, string[]> = {
  security: [
    "Potential Data Exfiltration Detected",
    "Brute Force Attack Attempt",
    "Malware Activity Detected",
    "Suspicious Process Execution",
    "Unauthorized Access Attempt",
    "Security Policy Violation",
  ],
  system: [
    "High CPU Usage Alert",
    "Memory Usage Threshold Exceeded",
    "Disk Space Critical",
    "System Service Failure",
    "Kernel Panic Detected",
    "Hardware Error Detected",
  ],
  network: [
    "Network Anomaly Detected",
    "DDoS Attack Suspected",
    "Unusual Traffic Pattern",
    "Port Scan Detected",
    "Network Interface Down",
    "Bandwidth Threshold Exceeded",
  ],
  application: [
    "Application Error Rate Spike",
    "API Response Time Degradation",
    "Database Connection Failure",
    "Cache Hit Rate Drop",
    "Application Crash Detected",
    "Service Dependencies Issue",
  ],
  database: [
    "Database Performance Degradation",
    "Replication Lag Alert",
    "Database Space Critical",
    "Deadlock Detected",
    "Backup Failure Alert",
    "Index Corruption Warning",
  ],
  compliance: [
    "Compliance Policy Violation",
    "Audit Log Gap Detected",
    "Unauthorized Data Access",
    "Configuration Drift Detected",
    "Security Control Failure",
    "Regulatory Requirement Breach",
  ],
};

const alertDescriptions: Record<AlertCategory, string[]> = {
  security: [
    "Large volume of data transferred to external IP {ip}",
    "Multiple failed login attempts from IP {ip}",
    "Malicious file activity detected on {host}",
    "Unauthorized process execution on {host}",
    "Failed authentication attempts from unusual location",
    "Security policy violation detected on {host}",
  ],
  system: [
    "CPU usage exceeded 90% for over 5 minutes",
    "Available memory dropped below 10%",
    "Disk usage exceeded 95% on {host}",
    "Critical system service failed to start",
    "System kernel panic occurred on {host}",
    "Hardware failure detected on {component}",
  ],
  network: [
    "Unusual traffic pattern detected from {ip}",
    "High volume of incoming traffic suggesting DDoS",
    "Abnormal outbound traffic to {ip}",
    "Sequential port scan detected from {ip}",
    "Network interface {interface} is down",
    "Network bandwidth exceeded threshold",
  ],
  application: [
    "Error rate exceeded 5% in last 5 minutes",
    "API latency increased by 200%",
    "Failed to connect to database after 3 retries",
    "Cache hit rate dropped below 50%",
    "Application crashed due to unhandled exception",
    "Multiple service dependencies failing",
  ],
  database: [
    "Query response time exceeded 10 seconds",
    "Replication lag exceeded 300 seconds",
    "Database space usage above 90%",
    "Multiple deadlocks detected in last 5 minutes",
    "Database backup failed on {host}",
    "Critical index corruption detected",
  ],
  compliance: [
    "Unauthorized access to sensitive data detected",
    "Missing audit logs for time period",
    "Unauthorized data export detected",
    "System configuration differs from baseline",
    "Failed security control check on {host}",
    "Non-compliant activity detected",
  ],
};

function getRandomTitle(category: AlertCategory): string {
  const titles = alertTitles[category];
  return titles[Math.floor(Math.random() * titles.length)];
}

function getRandomDescription(category: AlertCategory): string {
  const descriptions = alertDescriptions[category];
  let description =
    descriptions[Math.floor(Math.random() * descriptions.length)];

  // Replace placeholders with random values
  description = description.replace("{ip}", generateRandomIP());
  description = description.replace("{host}", generateRandomHostname());
  description = description.replace(
    "{interface}",
    "eth" + Math.floor(Math.random() * 4)
  );
  description = description.replace(
    "{component}",
    ["CPU", "RAM", "Disk", "Network"][Math.floor(Math.random() * 4)]
  );

  return description;
}

function getRelatedEvents(): string[] {
  const count = Math.floor(Math.random() * 3) + 1; // 1-3 related events
  return Array.from({ length: count }, () => generateId("evt-"));
}

export function generateMockAlert(): Alert {
  const category = getRandomAlertCategory();
  const severity = getRandomAlertSeverity();
  const title = getRandomTitle(category);
  const description = getRandomDescription(category);

  return {
    id: generateId("alt-"),
    timestamp: getRandomTimestamp(),
    title,
    description,
    severity,
    category,
    source: `${category}-monitor`,
    status: getRandomAlertStatus(),
    assignedTo: Math.random() > 0.7 ? "admin" : undefined,
    relatedEvents: Math.random() > 0.5 ? getRelatedEvents() : undefined,
    metadata: {
      version: "1.0",
      environment: ["production", "staging", "development"][
        Math.floor(Math.random() * 3)
      ],
      datacenter: ["us-east", "us-west", "eu-central"][
        Math.floor(Math.random() * 3)
      ],
      priority:
        severity === "critical"
          ? 1
          : severity === "high"
          ? 2
          : severity === "medium"
          ? 3
          : 4,
    },
  };
}

export function generateMockAlerts(count: number = 20): Alert[] {
  return Array.from({ length: count }, () => generateMockAlert()).sort(
    (a, b) => new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime()
  );
}
