'use client'

import { useApi } from '@/hooks/useApi'
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert"
import { AlertCircle } from 'lucide-react'

interface ThreatAlert {
  id: number;
  title: string;
  description: string;
}

export function ThreatAlerts() {
  //TODO: MAKE THIS const WHEN API IS SETUP!
  let { data: alerts, loading, error } = useApi<ThreatAlert[]>('/threats');

  if (loading) return <div>Loading...</div>;
  // FOR NOW THE ENDPOINTS ARE NOT SETUP I DONT WANNA SEE ERRORS YET.
  // if (error) return <div>Error: {error.message}</div>;

  if(!alerts) {
    // Endpoints not setup yet ill do that later once core features are complete this is not super high priority yet
    alerts = [
      { id: 1, title: 'Potential Data Exfiltration', description: 'Unusual outbound traffic detected from server 192.168.1.100' },
      { id: 2, title: 'Brute Force Attack', description: 'Multiple failed login attempts on admin panel' },
      { id: 3, title: 'Malware Detected', description: 'Trojan identified on workstation WS-005' },
    ]
  }

  return (
    <div className="bg-white shadow-md rounded-lg p-4">
      <h2 className="text-xl font-semibold mb-4">Threat Alerts (XDR Data)</h2>
      <div className="space-y-4">
        {alerts && alerts.length > 0 ? (
          alerts.map((alert) => (
            <Alert key={alert.id} variant="destructive">
              <AlertCircle className="h-4 w-4" />
              <AlertTitle>{alert.title}</AlertTitle>
              <AlertDescription>
                {alert.description}
              </AlertDescription>
            </Alert>
          ))
        ) : (
          <div>No threat alerts to display</div>
        )}
      </div>
    </div>
  )
}