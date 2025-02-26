"use client";

import { Card, CardContent } from "@/components/ui/card";

interface ThreatStatsProps {
  totalIncidents: number;
  averageSeverity: number;
  highSeverityCount: number;
}

export function ThreatStats({
  totalIncidents,
  averageSeverity,
  highSeverityCount,
}: ThreatStatsProps) {
  return (
    <div className="grid grid-cols-3 gap-4 mb-6">
      <Card>
        <CardContent className="pt-6">
          <div className="text-2xl font-bold">{totalIncidents}</div>
          <p className="text-sm text-muted-foreground">Total Incidents</p>
        </CardContent>
      </Card>
      <Card>
        <CardContent className="pt-6">
          <div className="text-2xl font-bold">
            {Math.round(averageSeverity)}%
          </div>
          <p className="text-sm text-muted-foreground">Average Severity</p>
        </CardContent>
      </Card>
      <Card>
        <CardContent className="pt-6">
          <div className="text-2xl font-bold">{highSeverityCount}</div>
          <p className="text-sm text-muted-foreground">High Severity</p>
        </CardContent>
      </Card>
    </div>
  );
}
