'use client'

import { useAlerts, AlertFilters, Alert as AlertType, AlertCategory } from "@/app/hooks/useAlerts"
import { AlertItem } from "./AlertItem"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Button } from "@/components/ui/button"
import { Skeleton } from "@/components/ui/skeleton"
import { AlertCircle, RotateCcw } from "lucide-react"

const severityOptions: AlertType['severity'][] = ['critical', 'high', 'medium', 'low'];
const categoryOptions: AlertCategory[] = ['security', 'system', 'network', 'application', 'database', 'compliance'];
const statusOptions: AlertType['status'][] = ['new', 'acknowledged', 'investigating', 'resolved', 'false_positive'];

function AlertListSkeleton() {
  return (
    <div className="space-y-4">
      {[1, 2, 3].map((i) => (
        <Card key={i} className="mb-4">
          <CardContent className="pt-6">
            <div className="flex items-start gap-4">
              <Skeleton className="h-5 w-5 rounded-full" />
              <div className="flex-grow space-y-2">
                <Skeleton className="h-6 w-3/4" />
                <Skeleton className="h-4 w-full" />
                <div className="flex gap-2 mt-4">
                  <Skeleton className="h-8 w-24" />
                  <Skeleton className="h-8 w-24" />
                </div>
              </div>
            </div>
          </CardContent>
        </Card>
      ))}
    </div>
  )
}

export function AlertList() {
  const {
    alerts,
    loading,
    error,
    filters,
    updateFilters,
    clearFilters,
    updateAlertStatus,
    isMockData
  } = useAlerts({
    useMockData: true,
    refreshInterval: 30000
  });

  if (error) {
    return (
      <Card className="mb-4">
        <CardContent className="pt-6">
          <div className="flex items-center gap-2 text-destructive">
            <AlertCircle className="h-5 w-5" />
            <p>Error loading alerts: {error.message}</p>
          </div>
        </CardContent>
      </Card>
    );
  }

  return (
    <div className="space-y-4">
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center justify-between">
            <span>Alert Filters</span>
            <Button
              variant="outline"
              size="sm"
              className="flex items-center gap-2"
              onClick={clearFilters}
            >
              <RotateCcw className="h-4 w-4" />
              Reset
            </Button>
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <Select
              value={filters.severity}
              onValueChange={(value) => updateFilters({ severity: value as AlertType['severity'] })}
            >
              <SelectTrigger>
                <SelectValue placeholder="Filter by severity" />
              </SelectTrigger>
              <SelectContent>
                {severityOptions.map((severity) => (
                  <SelectItem key={severity} value={severity}>
                    {severity.charAt(0).toUpperCase() + severity.slice(1)}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>

            <Select
              value={filters.category}
              onValueChange={(value) => updateFilters({ category: value as AlertCategory })}
            >
              <SelectTrigger>
                <SelectValue placeholder="Filter by category" />
              </SelectTrigger>
              <SelectContent>
                {categoryOptions.map((category) => (
                  <SelectItem key={category} value={category}>
                    {category.charAt(0).toUpperCase() + category.slice(1)}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>

            <Select
              value={filters.status}
              onValueChange={(value) => updateFilters({ status: value as AlertType['status'] })}
            >
              <SelectTrigger>
                <SelectValue placeholder="Filter by status" />
              </SelectTrigger>
              <SelectContent>
                {statusOptions.map((status) => (
                  <SelectItem key={status} value={status}>
                    {status.split('_').map(word => word.charAt(0).toUpperCase() + word.slice(1)).join(' ')}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        </CardContent>
      </Card>

      {loading ? (
        <AlertListSkeleton />
      ) : (
        <div className="space-y-4">
          {alerts.length === 0 ? (
            <Card>
              <CardContent className="pt-6">
                <p className="text-center text-muted-foreground">No alerts found</p>
              </CardContent>
            </Card>
          ) : (
            alerts.map((alert) => (
              <AlertItem
                key={alert.id}
                alert={alert}
                onStatusChange={updateAlertStatus}
              />
            ))
          )}
        </div>
      )}
    </div>
  )
}
