'use client'

import { useState } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"

const mockIncidents = [
  { id: 1, title: 'Suspicious Login Attempt', severity: 'High', status: 'Open', timestamp: '2023-06-10 15:30:00' },
  { id: 2, title: 'Malware Detected', severity: 'Critical', status: 'In Progress', timestamp: '2023-06-10 16:45:00' },
  { id: 3, title: 'Unusual File Access', severity: 'Medium', status: 'Closed', timestamp: '2023-06-09 09:15:00' },
]

export default function IncidentList() {
  const [incidents, setIncidents] = useState(mockIncidents)

  const getSeverityColor = (severity: string) => {
    switch (severity.toLowerCase()) {
      case 'critical': return 'bg-red-500 text-white'
      case 'high': return 'bg-orange-500 text-white'
      case 'medium': return 'bg-yellow-500 text-white'
      default: return 'bg-blue-500 text-white'
    }
  }

  const getStatusColor = (status: string) => {
    switch (status.toLowerCase()) {
      case 'open': return 'bg-green-500 text-white'
      case 'in progress': return 'bg-yellow-500 text-white'
      case 'closed': return 'bg-gray-500 text-white'
      default: return 'bg-blue-500 text-white'
    }
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>Active Incidents</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="overflow-x-auto">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Title</TableHead>
                <TableHead>Severity</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>Timestamp</TableHead>
                <TableHead>Action</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {incidents.map((incident) => (
                <TableRow key={incident.id}>
                  <TableCell>{incident.title}</TableCell>
                  <TableCell>
                    <Badge className={getSeverityColor(incident.severity)}>{incident.severity}</Badge>
                  </TableCell>
                  <TableCell>
                    <Badge className={getStatusColor(incident.status)}>{incident.status}</Badge>
                  </TableCell>
                  <TableCell>{incident.timestamp}</TableCell>
                  <TableCell>
                    <Button size="sm">View Details</Button>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </div>
      </CardContent>
    </Card>
  )
}

