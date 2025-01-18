'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { useAuth } from '../contexts/AuthContext'
import { useTeam } from '../contexts/TeamContext'
import { SidebarNav } from '../components/Sidebar'
import { SidebarProvider, SidebarInset } from "@/components/ui/sidebar"
import EventLog from '../components/EventLog'
import ThreatSummary from '../components/ThreatSummary'
import SystemHealth from '../components/SystemHealth'
import AlertComponent from '../components/AlertComponent'
import NetworkTraffic from '../components/NetworkTraffic'
import { Alert, AlertDescription } from "@/components/ui/alert"

export default function Dashboard() {
  const { user } = useAuth()
  const { currentTeam } = useTeam()
  const router = useRouter()
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (!user) {
      router.push('/login')
    }
  }, [user, router])

  if (!user || !currentTeam) {
    return <div className="flex items-center justify-center min-h-screen">Loading...</div>
  }

  return (
    <SidebarProvider>
      <div className="flex h-screen overflow-hidden">
        <SidebarNav />
        <SidebarInset className="flex-1 overflow-auto">
          <main className="p-4 md:p-6 bg-background">
            <h1 className="text-3xl font-bold mb-6">
              {currentTeam.name} - SIEM Dashboard
            </h1>
            {error && (
              <Alert variant="destructive" className="mb-4">
                <AlertDescription>{error}</AlertDescription>
              </Alert>
            )}
            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-2 xl:grid-cols-2">
              <EventLog />
              <ThreatSummary />
              <SystemHealth />
              <AlertComponent />
              <NetworkTraffic />
            </div>
          </main>
        </SidebarInset>
      </div>
    </SidebarProvider>
  )
}

