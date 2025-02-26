'use client'

import { NetworkStats } from "@/app/hooks/useNetworkTraffic"
import { Card, CardContent } from "@/components/ui/card"
import { ArrowDownIcon, ArrowUpIcon } from "lucide-react"

interface NetworkTrafficStatsProps {
  stats: NetworkStats
  formattedStats: {
    totalInbound: string
    totalOutbound: string
    peakInbound: string
    peakOutbound: string
    averageInbound: string
    averageOutbound: string
  }
}

export function NetworkTrafficStats({ stats, formattedStats }: NetworkTrafficStatsProps) {
  return (
    <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
      <Card>
        <CardContent className="pt-6">
          <div className="space-y-2">
            <h4 className="text-sm font-medium">Total Traffic</h4>
            <div className="grid grid-cols-2 gap-4">
              <div className="flex items-center space-x-2">
                <ArrowDownIcon className="h-4 w-4 text-primary" />
                <span className="text-sm text-muted-foreground">
                  {formattedStats.totalInbound}
                </span>
              </div>
              <div className="flex items-center space-x-2">
                <ArrowUpIcon className="h-4 w-4 text-success" />
                <span className="text-sm text-muted-foreground">
                  {formattedStats.totalOutbound}
                </span>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardContent className="pt-6">
          <div className="space-y-2">
            <h4 className="text-sm font-medium">Peak Traffic</h4>
            <div className="grid grid-cols-2 gap-4">
              <div className="flex items-center space-x-2">
                <ArrowDownIcon className="h-4 w-4 text-primary" />
                <span className="text-sm text-muted-foreground">
                  {formattedStats.peakInbound}
                </span>
              </div>
              <div className="flex items-center space-x-2">
                <ArrowUpIcon className="h-4 w-4 text-success" />
                <span className="text-sm text-muted-foreground">
                  {formattedStats.peakOutbound}
                </span>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardContent className="pt-6">
          <div className="space-y-2">
            <h4 className="text-sm font-medium">Average Traffic</h4>
            <div className="grid grid-cols-2 gap-4">
              <div className="flex items-center space-x-2">
                <ArrowDownIcon className="h-4 w-4 text-primary" />
                <span className="text-sm text-muted-foreground">
                  {formattedStats.averageInbound}
                </span>
              </div>
              <div className="flex items-center space-x-2">
                <ArrowUpIcon className="h-4 w-4 text-success" />
                <span className="text-sm text-muted-foreground">
                  {formattedStats.averageOutbound}
                </span>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
