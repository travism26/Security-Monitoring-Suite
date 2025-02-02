'use client'

import { useState, useEffect } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts'

const generateMockData = () => {
  return Array.from({ length: 24 }, (_, i) => ({
    hour: i,
    inbound: Math.floor(Math.random() * 1000),
    outbound: Math.floor(Math.random() * 1000),
  }))
}

export default function NetworkTraffic() {
  const [data, setData] = useState(generateMockData())

  useEffect(() => {
    const interval = setInterval(() => {
      setData(generateMockData())
    }, 60000) // Update every minute

    return () => clearInterval(interval)
  }, [])

  return (
    <Card>
      <CardHeader>
        <CardTitle>Network Traffic (24h)</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="h-[300px]">
          <ResponsiveContainer width="100%" height="100%">
            <BarChart data={data}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="hour" />
              <YAxis />
              <Tooltip />
              <Bar dataKey="inbound" fill="#8884d8" name="Inbound" />
              <Bar dataKey="outbound" fill="#82ca9d" name="Outbound" />
            </BarChart>
          </ResponsiveContainer>
        </div>
      </CardContent>
    </Card>
  )
}

