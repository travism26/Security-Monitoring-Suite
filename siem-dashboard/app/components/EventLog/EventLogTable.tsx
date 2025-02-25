import { SecurityEvent } from "@/app/hooks/useEventLog";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";

interface EventLogTableProps {
  events: SecurityEvent[];
}

export function EventLogTable({ events }: EventLogTableProps) {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Timestamp</TableHead>
          <TableHead>Type</TableHead>
          <TableHead>Source</TableHead>
          <TableHead>Description</TableHead>
          <TableHead>Severity</TableHead>
          <TableHead>Status</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {events.map((event) => (
          <TableRow key={event.id}>
            <TableCell>{event.timestamp}</TableCell>
            <TableCell>{event.type}</TableCell>
            <TableCell>{event.source}</TableCell>
            <TableCell>{event.description}</TableCell>
            <TableCell>
              <span
                className={`inline-block px-2 py-1 rounded-full text-xs font-medium ${
                  event.severity === "high"
                    ? "bg-red-100 text-red-800"
                    : event.severity === "medium"
                    ? "bg-yellow-100 text-yellow-800"
                    : "bg-green-100 text-green-800"
                }`}
              >
                {event.severity}
              </span>
            </TableCell>
            <TableCell>
              <span
                className={`inline-block px-2 py-1 rounded-full text-xs font-medium ${
                  event.status === "new"
                    ? "bg-blue-100 text-blue-800"
                    : event.status === "investigating"
                    ? "bg-purple-100 text-purple-800"
                    : "bg-gray-100 text-gray-800"
                }`}
              >
                {event.status}
              </span>
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
