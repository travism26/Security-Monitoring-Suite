'use client'

import { useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { useAuth } from '../contexts/AuthContext'
import { useTeam } from '../contexts/TeamContext'
import { SidebarNav } from '../components/Sidebar'
import { SidebarProvider, SidebarInset } from "@/components/ui/sidebar"

export default function AlertsPage() {
  const { user } = useAuth()
  const { currentTeam } = useTeam()
  const router = useRouter()

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
              {currentTeam.name} - Alerts
            </h1>
            <p>Alerts content goes here.</p>
          </main>
        </SidebarInset>
      </div>
    </SidebarProvider>
  )
}

