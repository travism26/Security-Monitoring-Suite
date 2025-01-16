'use client'

import { useApi } from '@/hooks/useApi'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"

interface Event {
  id: number;
  type: string;
  source: string;
  timestamp: string;
}

export function EventLog() {
  //TODO: MAKE THIS const WHEN API IS SETUP!
  let { data: events, loading, error } = useApi<Event[]>('/events');

  if(!events) {
    console.log("mtravis - event is null so lets set some test data")
    events = [
      { id: 1, type: 'Authentication', source: '192.168.1.100', timestamp: '2023-06-10 14:30:45' },
      { id: 2, type: 'File Access', source: '10.0.0.5', timestamp: '2023-06-10 14:31:12' },
      { id: 3, type: 'Network Connection', source: '172.16.0.20', timestamp: '2023-06-10 14:32:01' },
      { id: 4, type: 'System Update', source: 'Server01', timestamp: '2023-06-10 14:33:57' },
      { id: 5, type: 'User Created', source: 'AdminPanel', timestamp: '2023-06-10 14:35:23' },
    ]
  }

  if (loading) return <div>Loading...</div>;
  // FOR NOW THE ENDPOINTS ARE NOT SETUP I DONT WANNA SEE ERRORS YET.
  // if (error) return <div>Error: {error.message}</div>; Recent Events (Log aggregator Data)

  return (
    <div className="bg-white shadow-md rounded-lg p-4 overflow-x-auto">
      <h2 className="text-xl font-semibold mb-4">Recent Events (Log aggregator Data)</h2>
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Type</TableHead>
            <TableHead>Source</TableHead>
            <TableHead>Timestamp</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {events && events.length > 0 ? (
            events.map((event) => (
              <TableRow key={event.id}>
                <TableCell>{event.type}</TableCell>
                <TableCell>{event.source}</TableCell>
                <TableCell>{event.timestamp}</TableCell>
              </TableRow>
            ))
          ) : (
            <TableRow>
              <TableCell colSpan={3} className="text-center">No events to display</TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>
    </div>
  )}

