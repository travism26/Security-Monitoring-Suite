import { v4 as uuidv4 } from "uuid";
import {
  EventType,
  EventSeverity,
  AlertSeverity,
  AlertCategory,
  AlertStatus,
  SystemStatus,
  ResourceMetric,
} from "./types";

// Date utilities
export function getRandomTimestamp(
  startHoursAgo: number = 24,
  endHoursAgo: number = 0
): string {
  const end = new Date();
  const start = new Date(end.getTime() - startHoursAgo * 60 * 60 * 1000);
  const randomDate = new Date(
    start.getTime() + Math.random() * (end.getTime() - start.getTime())
  );
  return randomDate.toISOString();
}

// Random data generators
export function getRandomEnum<T extends { [key: string]: string }>(
  enumObj: T
): T[keyof T] {
  const enumValues = Object.values(enumObj);
  return enumValues[
    Math.floor(Math.random() * enumValues.length)
  ] as T[keyof T];
}

export function getRandomEventType(): EventType {
  const types: EventType[] = [
    "authentication",
    "system",
    "network",
    "security",
    "application",
    "database",
  ];
  return types[Math.floor(Math.random() * types.length)];
}

export function getRandomEventSeverity(): EventSeverity {
  const severities: EventSeverity[] = ["info", "warning", "error", "critical"];
  return severities[Math.floor(Math.random() * severities.length)];
}

export function getRandomAlertSeverity(): AlertSeverity {
  const severities: AlertSeverity[] = ["critical", "high", "medium", "low"];
  return severities[Math.floor(Math.random() * severities.length)];
}

export function getRandomAlertCategory(): AlertCategory {
  const categories: AlertCategory[] = [
    "security",
    "system",
    "network",
    "application",
    "database",
    "compliance",
  ];
  return categories[Math.floor(Math.random() * categories.length)];
}

export function getRandomAlertStatus(): AlertStatus {
  const statuses: AlertStatus[] = [
    "new",
    "acknowledged",
    "investigating",
    "resolved",
    "false_positive",
  ];
  return statuses[Math.floor(Math.random() * statuses.length)];
}

export function getRandomSystemStatus(): SystemStatus {
  const statuses: SystemStatus[] = [
    "healthy",
    "warning",
    "critical",
    "unknown",
  ];
  return statuses[Math.floor(Math.random() * statuses.length)];
}

// Resource metrics utilities
export function generateResourceMetric(
  minUsage: number = 0,
  maxUsage: number = 100,
  limit: number = 100
): ResourceMetric {
  const usage = minUsage + Math.random() * (maxUsage - minUsage);
  let status: SystemStatus = "healthy";
  if (usage > limit * 0.9) status = "critical";
  else if (usage > limit * 0.7) status = "warning";

  const rand = Math.random();
  const trend = rand < 0.33 ? "up" : rand < 0.66 ? "down" : "stable";

  return {
    usage,
    limit,
    status,
    trend,
  };
}

// ID generation
export function generateId(prefix: string = ""): string {
  return `${prefix}${uuidv4()}`;
}

// Network traffic utilities
export function generateTrafficValue(
  baseValue: number = 500,
  variance: number = 200
): number {
  return Math.max(0, baseValue + (Math.random() - 0.5) * variance);
}

// Business hours simulation
export function isBusinessHours(hour: number): boolean {
  return hour >= 8 && hour <= 18;
}

// Data volume formatting
export function formatBytes(bytes: number): string {
  if (bytes === 0) return "0 B";
  const k = 1024;
  const sizes = ["B", "KB", "MB", "GB", "TB"];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return `${parseFloat((bytes / Math.pow(k, i)).toFixed(2))} ${sizes[i]}`;
}

// Service uptime calculation
export function calculateUptime(startDate: Date): number {
  return Math.floor((new Date().getTime() - startDate.getTime()) / 1000);
}

// Random IP address generation
export function generateRandomIP(): string {
  return Array.from({ length: 4 }, () => Math.floor(Math.random() * 256)).join(
    "."
  );
}

// Random port generation
export function generateRandomPort(): number {
  return Math.floor(Math.random() * (65535 - 1024) + 1024);
}

// Random MAC address generation
export function generateRandomMAC(): string {
  return Array.from({ length: 6 }, () =>
    Math.floor(Math.random() * 256)
      .toString(16)
      .padStart(2, "0")
  ).join(":");
}

// Random hostname generation
export function generateRandomHostname(): string {
  const prefixes = ["srv", "web", "db", "app", "auth", "api"];
  const environments = ["prod", "dev", "stage", "test"];
  const numbers = Math.floor(Math.random() * 100)
    .toString()
    .padStart(2, "0");
  return `${prefixes[Math.floor(Math.random() * prefixes.length)]}-${
    environments[Math.floor(Math.random() * environments.length)]
  }-${numbers}`;
}
