'use client'

import React, { createContext, useState, useContext, useEffect } from 'react'

interface Team {
  id: string
  name: string
}

interface TeamContextType {
  currentTeam: Team | null
  teams: Team[]
  switchTeam: (teamId: string) => void
}

const TeamContext = createContext<TeamContextType | undefined>(undefined)

export function TeamProvider({ children }: { children: React.ReactNode }) {
  const [currentTeam, setCurrentTeam] = useState<Team | null>(null)
  const [teams, setTeams] = useState<Team[]>([])

  useEffect(() => {
    // Mock fetching teams. In a real app, you'd call your API here.
    const mockTeams = [
      { id: '1', name: 'Team A' },
      { id: '2', name: 'Team B' },
      { id: '3', name: 'Team C' },
    ]
    setTeams(mockTeams)
    setCurrentTeam(mockTeams[0])
  }, [])

  const switchTeam = (teamId: string) => {
    const team = teams.find(t => t.id === teamId)
    if (team) {
      setCurrentTeam(team)
    }
  }

  return (
    <TeamContext.Provider value={{ currentTeam, teams, switchTeam }}>
      {children}
    </TeamContext.Provider>
  )
}

export function useTeam() {
  const context = useContext(TeamContext)
  if (context === undefined) {
    throw new Error('useTeam must be used within a TeamProvider')
  }
  return context
}

