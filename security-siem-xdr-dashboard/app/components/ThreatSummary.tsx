'use client'

import { useThreats } from '../hooks/useThreats'
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Progress } from "@/components/ui/progress"

export default function ThreatSummary() {
  const { threats, loading, error } = useThreats()

  if (loading) return <div>Loading threats...</div>
  if (error) return <div>Error: {error.message}</div>

  return (
    <Card>
      <CardHeader>
        <CardTitle>Threat Summary</CardTitle>
      </CardHeader>
      <CardContent>
        {threats.map((threat) => (
          <div key={threat.name} className="mb-4">
            <div className="flex justify-between mb-1">
              <span>{threat.name}</span>
              <span>{threat.count} incidents</span>
            </div>
            <Progress value={threat.severity} className="w-full" />
          </div>
        ))}
      </CardContent>
    </Card>
  )
}

