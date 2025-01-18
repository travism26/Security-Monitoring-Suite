'use client'

import { useSystemHealth } from '../hooks/useSystemHealth'
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { CheckCircle, XCircle } from 'lucide-react'

export default function SystemHealth() {
  const { systems, loading, error } = useSystemHealth()

  if (loading) return <div>Loading system health...</div>
  if (error) return <div>Error: {error.message}</div>

  return (
    <Card>
      <CardHeader>
        <CardTitle>System Health</CardTitle>
      </CardHeader>
      <CardContent>
        <ul>
          {systems.map((system) => (
            <li key={system.name} className="flex items-center justify-between mb-2">
              <span>{system.name}</span>
              {system.status === 'Operational' ? (
                <CheckCircle className="text-green-500" />
              ) : (
                <XCircle className="text-red-500" />
              )}
            </li>
          ))}
        </ul>
      </CardContent>
    </Card>
  )
}

