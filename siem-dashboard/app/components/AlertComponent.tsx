'use client'

import { useEffect, useState } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert"
import { AlertCircle } from 'lucide-react'

export default function AlertComponent() {
  const [alerts, setAlerts] = useState<string[]>([])

  useEffect(() => {
    // Simulating real-time alerts
    const interval = setInterval(() => {
      setAlerts(prev => [...prev, `New alert at ${new Date().toLocaleTimeString()}`].slice(-5))
    }, 5000)

    return () => clearInterval(interval)
  }, [])

  return (
    <Card>
      <CardHeader>
        <CardTitle>Recent Alerts</CardTitle>
      </CardHeader>
      <CardContent>
        {alerts.map((alert, index) => (
          <Alert key={index} variant="destructive" className="mb-2">
            <AlertCircle className="h-4 w-4" />
            <AlertTitle>Alert</AlertTitle>
            <AlertDescription>{alert}</AlertDescription>
          </Alert>
        ))}
      </CardContent>
    </Card>
  )
}

