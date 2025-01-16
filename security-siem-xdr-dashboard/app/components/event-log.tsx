import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
  } from "@/components/ui/table"
  
  export function EventLog() {
    const events = [
      { id: 1, type: 'Authentication', source: '192.168.1.100', timestamp: '2023-06-10 14:30:45' },
      { id: 2, type: 'File Access', source: '10.0.0.5', timestamp: '2023-06-10 14:31:12' },
      { id: 3, type: 'Network Connection', source: '172.16.0.20', timestamp: '2023-06-10 14:32:01' },
      { id: 4, type: 'System Update', source: 'Server01', timestamp: '2023-06-10 14:33:57' },
      { id: 5, type: 'User Created', source: 'AdminPanel', timestamp: '2023-06-10 14:35:23' },
    ]
  
    return (
      <div className="bg-white shadow-md rounded-lg p-4">
        <h2 className="text-xl font-semibold mb-4">Recent Events</h2>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Type</TableHead>
              <TableHead>Source</TableHead>
              <TableHead>Timestamp</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {events.map((event) => (
              <TableRow key={event.id}>
                <TableCell>{event.type}</TableCell>
                <TableCell>{event.source}</TableCell>
                <TableCell>{event.timestamp}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </div>
    )
  }
  
  