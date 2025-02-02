'use client'

import { useState, useEffect } from 'react'

interface Alert {
  id: number
  message: string
  timestamp: string
}

export function useAlerts() {
  const [alerts, setAlerts] = useState<Alert[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<Error | null>(null)

  useEffect(() => {
    async function fetchAlerts() {
      try {
        const response = await fetch('/api/alerts')
        if (!response.ok) {
          throw new Error('Failed to fetch alerts')
        }
        const data = await response.json()
        setAlerts(data)
      } catch (err) {
        setError(err instanceof Error ? err : new Error('An error occurred'))
      } finally {
        setLoading(false)
      }
    }

    fetchAlerts()

    // Simulating real-time updates
    const interval = setInterval(fetchAlerts, 5000)

    return () => clearInterval(interval)
  }, [])

  return { alerts, loading, error }
}

