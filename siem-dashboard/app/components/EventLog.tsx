'use client'

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"

const mockEvents = [
  { id: 1, timestamp: '2023-06-10 10:30:15', type: 'Login', source: '192.168.1.100', description: 'Successful login' },
  { id: 2, timestamp: '2023-06-10 10:35:22', type: 'File Access', source: '192.168.1.101', description: 'Unauthorized file access attempt' },
  { id: 3, timestamp: '2023-06-10 10:40:05', type: 'Network', source: '192.168.1.102', description: 'Unusual outbound traffic detected' },
]

export default function EventLog() {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Event Log</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="overflow-x-auto">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Timestamp</TableHead>
                <TableHead>Type</TableHead>
                <TableHead>Source</TableHead>
                <TableHead>Description</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {mockEvents.map((event) => (
                <TableRow key={event.id}>
                  <TableCell>{event.timestamp}</TableCell>
                  <TableCell>{event.type}</TableCell>
                  <TableCell>{event.source}</TableCell>
                  <TableCell>{event.description}</TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </div>
      </CardContent>
    </Card>
  )
}

