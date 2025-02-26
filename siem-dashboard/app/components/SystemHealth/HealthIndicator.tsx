import React from "react";
import { HealthStatus } from "../../hooks/useSystemHealth";

interface HealthIndicatorProps {
  status: HealthStatus;
  size?: "sm" | "md" | "lg";
}

const statusColors = {
  healthy: "bg-green-500",
  warning: "bg-yellow-500",
  critical: "bg-red-500",
  unknown: "bg-gray-400",
};

const statusPulse = {
  healthy: "animate-pulse-slow",
  warning: "animate-pulse-medium",
  critical: "animate-pulse-fast",
  unknown: "",
};

const sizes = {
  sm: "w-2 h-2",
  md: "w-3 h-3",
  lg: "w-4 h-4",
};

export function HealthIndicator({ status, size = "md" }: HealthIndicatorProps) {
  return (
    <div
      className={`rounded-full ${sizes[size]} ${statusColors[status]} ${
        statusPulse[status]
      }`}
      title={`Status: ${status.charAt(0).toUpperCase() + status.slice(1)}`}
    />
  );
}
