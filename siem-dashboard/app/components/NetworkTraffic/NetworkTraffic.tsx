'use client'

import { useNetworkTraffic } from "@/app/hooks/useNetworkTraffic"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { NetworkTrafficChart } from "./NetworkTrafficChart"
import { NetworkTrafficStats } from "./NetworkTrafficStats"
import { NetworkTrafficSkeleton } from "./NetworkTrafficSkeleton"
import { AlertCircle } from "lucide-react"

export function NetworkTraffic() {
  const {
    data,
    stats,
    formattedStats,
    loading,
    error,
    isMockData
  } = useNetworkTraffic({
    useMockData: true,
    refreshInterval: 60000, // 1 minute
    timeRange: "24h"
  })

  if (error) {
    return (
      <Card>
        <CardContent className="pt-6">
          <div className="flex items-center gap-2 text-destructive">
            <AlertCircle className="h-5 w-5" />
            <p>Error loading network traffic data: {error.message}</p>
          </div>
        </CardContent>
      </Card>
    )
  }

  if (loading) {
    return <NetworkTrafficSkeleton />
  }

  return (
    <div className="space-y-6">
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle>Network Traffic</CardTitle>
          <Select defaultValue="24h">
            <SelectTrigger className="w-[180px]">
              <SelectValue placeholder="Select time range" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="1h">Last Hour</SelectItem>
              <SelectItem value="24h">Last 24 Hours</SelectItem>
              <SelectItem value="7d">Last 7 Days</SelectItem>
            </SelectContent>
          </Select>
        </CardHeader>
        <CardContent>
          <NetworkTrafficStats stats={stats} formattedStats={formattedStats} />
          <NetworkTrafficChart data={data} timeRange="24h" />
        </CardContent>
      </Card>
    </div>
  )
}
