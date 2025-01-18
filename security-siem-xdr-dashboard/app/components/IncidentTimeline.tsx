'use client'

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"

const timelineEvents = [
  { id: 1, title: 'Malware Detected', description: 'Endpoint ABC-123 reported malware detection', timestamp: '2023-06-10 16:45:00' },
  { id: 2, title: 'Quarantine Action', description: 'Malicious file quarantined on Endpoint ABC-123', timestamp: '2023-06-10 16:45:30' },
  { id: 3, title: 'Alert Generated', description: 'High-priority alert created for SOC team', timestamp: '2023-06-10 16:46:00' },
  { id: 4, title: 'Investigation Started', description: 'SOC analyst began investigation of the incident', timestamp: '2023-06-10 16:50:00' },
]

export default function IncidentTimeline() {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Incident Timeline</CardTitle>
      </CardHeader>
      <CardContent>
        <ol className="relative border-l border-gray-200 dark:border-gray-700">
          {timelineEvents.map((event, index) => (
            <li key={event.id} className="mb-10 ml-4">
              <div className="absolute w-3 h-3 bg-gray-200 rounded-full mt-1.5 -left-1.5 border border-white dark:border-gray-900 dark:bg-gray-700"></div>
              <time className="mb-1 text-sm font-normal leading-none text-gray-400 dark:text-gray-500">{event.timestamp}</time>
              <h3 className="text-lg font-semibold text-gray-900 dark:text-white">{event.title}</h3>
              <p className="mb-4 text-base font-normal text-gray-500 dark:text-gray-400">{event.description}</p>
            </li>
          ))}
        </ol>
      </CardContent>
    </Card>
  )
}

