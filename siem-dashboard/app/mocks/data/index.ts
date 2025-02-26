// Type exports
export * from "./types";

// Utility exports
export {
  getRandomTimestamp,
  getRandomEnum,
  getRandomEventType,
  getRandomEventSeverity,
  getRandomAlertSeverity,
  getRandomAlertCategory,
  getRandomAlertStatus,
  getRandomSystemStatus,
  generateResourceMetric,
  generateId,
  generateTrafficValue,
  isBusinessHours,
  formatBytes,
  calculateUptime,
  generateRandomIP,
  generateRandomPort,
  generateRandomMAC,
  generateRandomHostname,
} from "./utils";

// Event data exports
export { generateMockEvent, generateMockEvents } from "./events";

// Alert data exports
export { generateMockAlert, generateMockAlerts } from "./alerts";

// System metrics exports
export {
  generateMockSystemMetrics,
  generateMetricsTimeSeries,
  generateTrendingMetrics,
} from "./metrics";

// System health exports
export {
  generateMockSystemHealth,
  generateMockSystemHealthWithIssues,
  generateHealthTimeSeries,
} from "./health";
