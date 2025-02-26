import React from "react";
import { HealthStatus } from "../../hooks/useSystemHealth";
import { HealthIndicator } from "./HealthIndicator";

interface StatusCardProps {
  title: string;
  status: HealthStatus;
  value: string | number;
  subtitle?: string;
  className?: string;
}

export function StatusCard({
  title,
  status,
  value,
  subtitle,
  className = "",
}: StatusCardProps) {
  return (
    <div
      className={`bg-white dark:bg-gray-800 rounded-lg shadow-sm p-4 ${className}`}
    >
      <div className="flex items-center justify-between mb-2">
        <h3 className="text-sm font-medium text-gray-600 dark:text-gray-300">
          {title}
        </h3>
        <HealthIndicator status={status} />
      </div>
      <div className="flex items-baseline">
        <span className="text-2xl font-semibold text-gray-900 dark:text-white">
          {value}
        </span>
        {typeof value === "number" && (
          <span className="ml-1 text-gray-500 dark:text-gray-400">%</span>
        )}
      </div>
      {subtitle && (
        <p className="mt-1 text-xs text-gray-500 dark:text-gray-400">
          {subtitle}
        </p>
      )}
    </div>
  );
}
