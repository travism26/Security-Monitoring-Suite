'use client'

import { NetworkTrafficData } from "@/app/hooks/useNetworkTraffic"
import { Card, CardContent } from "@/components/ui/card"
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts'

interface NetworkTrafficChartProps {
  data: NetworkTrafficData[]
  timeRange: "1h" | "24h" | "7d"
}

export function NetworkTrafficChart({ data, timeRange }: NetworkTrafficChartProps) {
  const formatXAxis = (hour: number) => {
    if (timeRange === "1h") {
      return `${hour}:00`
    }
    return `${hour}:00`
  }

  const CustomTooltip = ({ active, payload, label }: any) => {
    if (active && payload && payload.length) {
      return (
        <div className="bg-background border rounded-lg p-3 shadow-lg">
          <p className="text-sm font-medium">{`${formatXAxis(label)}`}</p>
          <p className="text-sm text-blue-500">{`Inbound: ${payload[0].value} MB`}</p>
          <p className="text-sm text-green-500">{`Outbound: ${payload[1].value} MB`}</p>
        </div>
      )
    }
    return null
  }

  return (
    <div className="h-[300px] w-full">
      <ResponsiveContainer width="100%" height="100%">
        <BarChart data={data} margin={{ top: 10, right: 10, left: 10, bottom: 0 }}>
          <CartesianGrid strokeDasharray="3 3" className="stroke-muted" />
          <XAxis
            dataKey="hour"
            tickFormatter={formatXAxis}
            className="text-muted-foreground text-xs"
          />
          <YAxis
            className="text-muted-foreground text-xs"
            tickFormatter={(value) => `${value} MB`}
          />
          <Tooltip content={<CustomTooltip />} />
          <Bar
            dataKey="inbound"
            fill="hsl(var(--primary))"
            name="Inbound"
            radius={[4, 4, 0, 0]}
          />
          <Bar
            dataKey="outbound"
            fill="hsl(var(--success))"
            name="Outbound"
            radius={[4, 4, 0, 0]}
          />
        </BarChart>
      </ResponsiveContainer>
    </div>
  )
}
