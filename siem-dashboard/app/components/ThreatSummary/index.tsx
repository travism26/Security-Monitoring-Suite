"use client";

import { useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useThreatSummary } from "@/app/hooks/useThreatSummary";
import { ThreatStats } from "./ThreatStats";
import { ThreatFilter } from "./ThreatFilter";
import { ThreatChart } from "./ThreatChart";

export function ThreatSummary() {
  const [filter, setFilter] = useState("");
  const [sortBy, setSortBy] = useState<"severity" | "count">("severity");

  const { threats, isLoading, error, stats } = useThreatSummary({
    filter,
    sortBy,
  });

  return (
    <Card>
      <CardHeader>
        <CardTitle>Threat Summary</CardTitle>
      </CardHeader>
      <CardContent>
        <ThreatStats
          totalIncidents={stats.totalIncidents}
          averageSeverity={stats.averageSeverity}
          highSeverityCount={stats.highSeverityCount}
        />
        <ThreatFilter
          filter={filter}
          sortBy={sortBy}
          onFilterChange={setFilter}
          onSortChange={setSortBy}
        />
        {error ? (
          <div className="text-sm text-destructive">
            Error loading threat data: {error.message}
          </div>
        ) : (
          <ThreatChart threats={threats} isLoading={isLoading} />
        )}
      </CardContent>
    </Card>
  );
}
