'use client'

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { CheckCircle, XCircle } from 'lucide-react'

const mockSystems = [
  { name: 'Firewall', status: 'Operational' },
  { name: 'IDS', status: 'Operational' },
  { name: 'Log Server', status: 'Down' },
  { name: 'Email Filter', status: 'Operational' },
]

export default function SystemHealth() {
  return (
    <Card>
      <CardHeader>
        <CardTitle>System Health</CardTitle>
      </CardHeader>
      <CardContent>
        <ul>
          {mockSystems.map((system) => (
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

