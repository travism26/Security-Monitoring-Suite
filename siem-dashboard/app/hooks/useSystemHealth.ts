'use client'

import { useState, useEffect } from 'react'

interface System {
  name: string
  status: 'Operational' | 'Down'
}

export function useSystemHealth() {
  const [systems, setSystems] = useState<System[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<Error | null>(null)

  useEffect(() => {
    async function fetchSystemHealth() {
      try {
        const response = await fetch('/api/system-health')
        if (!response.ok) {
          throw new Error('Failed to fetch system health')
        }
        const data = await response.json()
        setSystems(data)
      } catch (err) {
        setError(err instanceof Error ? err : new Error('An error occurred'))
      } finally {
        setLoading(false)
      }
    }

    fetchSystemHealth()
  }, [])

  return { systems, loading, error }
}

