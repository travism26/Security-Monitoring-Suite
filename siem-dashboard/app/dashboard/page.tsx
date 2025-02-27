'use client'

import { useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { useAuth } from '../contexts/AuthContext'
import { SidebarNav } from '../components/Sidebar'
import { SidebarProvider, SidebarInset } from "@/components/ui/sidebar"
import { EventLog } from '../components/EventLog/EventLog'
import { ThreatSummary } from '../components/ThreatSummary'
import { SystemHealth } from '../components/SystemHealth/SystemHealth'
import { Alert } from '../components/Alert'
import { NetworkTraffic } from '../components/NetworkTraffic'

export default function Dashboard() {
  const { user, loading } = useAuth()
  const router = useRouter()

  useEffect(() => {
    if (!loading && !user) {
      console.log("[Dashboard] No authenticated user found, redirecting to login");
      router.replace('/login');
      return;
    }

    if (user) {
      console.log("[Dashboard] User authenticated:", {
        userId: user.id,
        email: user.email,
        role: user.role
      });
    }
  }, [user, loading, router]);

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-xl">Loading...</div>
      </div>
    );
  }

  if (!user) {
    return null;
  }

  return (
    <SidebarProvider>
      <div className="flex h-screen overflow-hidden">
        <SidebarNav />
        <SidebarInset className="flex-1 overflow-auto">
          <main className="p-4 md:p-6 bg-background">
            <h1 className="text-3xl font-bold mb-6">
              SIEM Dashboard
            </h1>
            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-2 xl:grid-cols-2">
              <EventLog />
              <ThreatSummary />
              <SystemHealth />
              <Alert />
              <NetworkTraffic />
            </div>
          </main>
        </SidebarInset>
      </div>
    </SidebarProvider>
  )
}
