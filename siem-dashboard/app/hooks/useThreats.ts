'use client'

import { useState, useEffect } from 'react'

interface Threat {
  name: string
  count: number
  severity: number
}

export function useThreats() {
  const [threats, setThreats] = useState<Threat[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<Error | null>(null)

  useEffect(() => {
    async function fetchThreats() {
      try {
        const response = await fetch('/api/threats')
        if (!response.ok) {
          throw new Error('Failed to fetch threats')
        }
        const data = await response.json()
        setThreats(data)
      } catch (err) {
        setError(err instanceof Error ? err : new Error('An error occurred'))
      } finally {
        setLoading(false)
      }
    }

    fetchThreats()
  }, [])

  return { threats, loading, error }
}

