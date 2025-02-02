'use client'

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Progress } from "@/components/ui/progress"

const mockThreats = [
  { name: 'Malware', count: 15, severity: 70 },
  { name: 'Phishing', count: 8, severity: 60 },
  { name: 'DDoS', count: 3, severity: 40 },
]

export default function ThreatSummary() {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Threat Summary</CardTitle>
      </CardHeader>
      <CardContent>
        {mockThreats.map((threat) => (
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

