import React from "react";
import { useSystemHealth } from "../../hooks/useSystemHealth";
import { StatusCard } from "./StatusCard";
import { SystemHealthSkeleton } from "./SystemHealthSkeleton";

interface SystemHealthProps {
  className?: string;
  useMockData?: boolean;
  refreshInterval?: number;
}

export function SystemHealth({
  className = "",
  useMockData,
  refreshInterval,
}: SystemHealthProps) {
  const { health, loading, error, isMockData } = useSystemHealth({
    useMockData,
    refreshInterval,
  });

  if (loading) {
    return <SystemHealthSkeleton />;
  }

  if (error) {
    return (
      <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4">
        <p className="text-red-600 dark:text-red-400">
          Error loading system health: {error.message}
        </p>
      </div>
    );
  }

  if (!health) {
    return null;
  }

  return (
    <div className={`space-y-4 ${className}`}>
      <div className="flex items-center justify-between">
        <h2 className="text-xl font-semibold text-gray-900 dark:text-white">
          System Health
        </h2>
        {isMockData && (
          <span className="text-xs text-gray-500 dark:text-gray-400 bg-gray-100 dark:bg-gray-800 px-2 py-1 rounded">
            Mock Data
          </span>
        )}
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <StatusCard
          title="CPU Usage"
          status={health.cpu.status}
          value={health.cpu.usage}
          subtitle={`Warning at ${health.cpu.threshold.warning}%, Critical at ${health.cpu.threshold.critical}%`}
        />
        <StatusCard
          title="Memory Usage"
          status={health.memory.status}
          value={health.memory.usage}
          subtitle={`Warning at ${health.memory.threshold.warning}%, Critical at ${health.memory.threshold.critical}%`}
        />
        <StatusCard
          title="Disk Usage"
          status={health.disk.status}
          value={health.disk.usage}
          subtitle={`Warning at ${health.disk.threshold.warning}%, Critical at ${health.disk.threshold.critical}%`}
        />
        <StatusCard
          title="Network Status"
          status={health.network.status}
          value={`${health.network.inboundUtilization}/${health.network.outboundUtilization}`}
          subtitle="Inbound/Outbound Utilization %"
        />
      </div>

      <div className="mt-4 flex items-center justify-between text-sm text-gray-500 dark:text-gray-400">
        <div className="flex items-center space-x-2">
          <span
            className={`inline-block w-2 h-2 rounded-full ${
              health.overall === "healthy"
                ? "bg-green-500"
                : health.overall === "warning"
                ? "bg-yellow-500"
                : "bg-red-500"
            }`}
          />
          <span>Overall Status: {health.overall}</span>
        </div>
        <span>Last Updated: {new Date(health.lastUpdated).toLocaleTimeString()}</span>
      </div>
    </div>
  );
}
