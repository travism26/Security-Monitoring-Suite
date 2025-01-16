import { Sidebar } from './components/sidebar'
import { Header } from './components/header'
import { OverviewPanel } from './components/overview-panel'
import { EventLog } from './components/event-log'
import { ThreatAlerts } from './components/threat-alerts'

export default function Home() {
  return (
    <div className="flex flex-col h-screen bg-gray-100 lg:flex-row">
      <Sidebar />
      <div className="flex-1 flex flex-col overflow-hidden">
        <Header />
        <main className="flex-1 overflow-x-hidden overflow-y-auto bg-gray-100">
          <div className="container mx-auto px-4 py-8">
            <h1 className="text-3xl font-semibold text-gray-800 mb-6">Security Dashboard</h1>
            <OverviewPanel />
            <div className="mt-8 grid grid-cols-1 gap-8 lg:grid-cols-2">
              <EventLog />
              <ThreatAlerts />
            </div>
          </div>
        </main>
      </div>
    </div>
  )
}

