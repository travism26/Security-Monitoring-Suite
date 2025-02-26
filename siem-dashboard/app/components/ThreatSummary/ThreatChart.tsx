"use client";

import { Progress } from "@/components/ui/progress";
import { Skeleton } from "@/components/ui/skeleton";
import { type Threat } from "@/app/hooks/useThreatSummary";
import { cn } from "@/lib/utils";

interface ThreatChartProps {
  threats: Threat[];
  isLoading?: boolean;
}

export function ThreatChartSkeleton() {
  return (
    <div className="space-y-4">
      {[1, 2, 3].map((i) => (
        <div key={i} className="space-y-2">
          <div className="flex justify-between">
            <Skeleton className="h-4 w-[100px]" />
            <Skeleton className="h-4 w-[60px]" />
          </div>
          <Skeleton className="h-2 w-full" />
        </div>
      ))}
    </div>
  );
}

export function ThreatChart({ threats, isLoading }: ThreatChartProps) {
  if (isLoading) {
    return <ThreatChartSkeleton />;
  }

  return (
    <div className="space-y-4">
      {threats.map((threat) => (
        <div key={threat.name} className="space-y-2">
          <div className="flex justify-between items-center">
            <div className="flex items-center gap-2">
              <span className="font-medium">{threat.name}</span>
              {threat.severity >= 70 && (
                <span className="px-2 py-1 text-xs bg-destructive text-destructive-foreground rounded-full">
                  High
                </span>
              )}
            </div>
            <span className="text-sm text-muted-foreground">
              {threat.count} incidents
            </span>
          </div>
          <Progress
            value={threat.severity}
            className={cn(
              "h-2",
              threat.severity >= 70
                ? "[&>div]:bg-destructive"
                : threat.severity >= 50
                ? "[&>div]:bg-warning"
                : "[&>div]:bg-primary"
            )}
          />
          {threat.details && (
            <p className="text-sm text-muted-foreground mt-1">
              {threat.details}
            </p>
          )}
        </div>
      ))}
    </div>
  );
}
