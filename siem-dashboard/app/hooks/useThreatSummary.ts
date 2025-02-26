"use client";

import { useState } from "react";

export interface Threat {
  name: string;
  count: number;
  severity: number;
  timestamp?: string;
  details?: string;
}

interface UseThreatSummaryOptions {
  filter?: string;
  sortBy?: "severity" | "count";
}

const mockThreats: Threat[] = [
  {
    name: "Malware",
    count: 15,
    severity: 70,
    timestamp: new Date().toISOString(),
    details: "Multiple malware variants detected across network",
  },
  {
    name: "Phishing",
    count: 8,
    severity: 60,
    timestamp: new Date().toISOString(),
    details: "Targeted phishing campaigns detected",
  },
  {
    name: "DDoS",
    count: 3,
    severity: 40,
    timestamp: new Date().toISOString(),
    details: "Attempted DDoS attacks blocked",
  },
  {
    name: "Unauthorized Access",
    count: 5,
    severity: 65,
    timestamp: new Date().toISOString(),
    details: "Multiple failed login attempts detected",
  },
];

export function useThreatSummary(options: UseThreatSummaryOptions = {}) {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  // Filter threats based on name
  const filteredThreats = options.filter
    ? mockThreats.filter((threat) =>
        threat.name.toLowerCase().includes(options.filter!.toLowerCase())
      )
    : mockThreats;

  // Sort threats based on specified criteria
  const sortedThreats = [...filteredThreats].sort((a, b) => {
    if (options.sortBy === "severity") {
      return b.severity - a.severity;
    }
    if (options.sortBy === "count") {
      return b.count - a.count;
    }
    return 0;
  });

  // Calculate summary statistics
  const totalIncidents = sortedThreats.reduce(
    (sum, threat) => sum + threat.count,
    0
  );
  const averageSeverity =
    sortedThreats.reduce((sum, threat) => sum + threat.severity, 0) /
    sortedThreats.length;

  const highSeverityThreats = sortedThreats.filter(
    (threat) => threat.severity >= 70
  );

  return {
    threats: sortedThreats,
    isLoading,
    error,
    stats: {
      totalIncidents,
      averageSeverity,
      highSeverityCount: highSeverityThreats.length,
    },
  };
}
