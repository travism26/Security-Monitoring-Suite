import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert"
import { AlertCircle } from 'lucide-react'

export function ThreatAlerts() {
  const alerts = [
    { id: 1, title: 'Potential Data Exfiltration', description: 'Unusual outbound traffic detected from server 192.168.1.100' },
    { id: 2, title: 'Brute Force Attack', description: 'Multiple failed login attempts on admin panel' },
    { id: 3, title: 'Malware Detected', description: 'Trojan identified on workstation WS-005' },
  ]

  return (
    <div className="bg-white shadow-md rounded-lg p-4">
      <h2 className="text-xl font-semibold mb-4">Threat Alerts</h2>
      <div className="space-y-4">
        {alerts.map((alert) => (
          <Alert key={alert.id} variant="destructive">
            <AlertCircle className="h-4 w-4" />
            <AlertTitle>{alert.title}</AlertTitle>
            <AlertDescription>
              {alert.description}
            </AlertDescription>
          </Alert>
        ))}
      </div>
    </div>
  )
}

