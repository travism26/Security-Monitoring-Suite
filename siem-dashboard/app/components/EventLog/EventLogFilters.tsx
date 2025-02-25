import { EventFilters, EventType } from "@/app/hooks/useEventLog";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { CalendarIcon } from "@radix-ui/react-icons";
import { format } from "date-fns";
import { Calendar } from "@/components/ui/calendar";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";

interface EventLogFiltersProps {
  filters: EventFilters;
  onUpdateFilters: (filters: Partial<EventFilters>) => void;
  onClearFilters: () => void;
}

const EVENT_TYPES: EventType[] = [
  "Login",
  "File Access",
  "Network",
  "System",
  "Security",
  "Database",
];

const SEVERITY_LEVELS = ["low", "medium", "high"] as const;
const STATUS_OPTIONS = ["new", "investigating", "resolved"] as const;

export function EventLogFilters({
  filters,
  onUpdateFilters,
  onClearFilters,
}: EventLogFiltersProps) {
  return (
    <div className="space-y-4">
      <div className="flex flex-wrap gap-4">
        {/* Event Type Filter */}
        <div className="w-[200px]">
          <Select
            value={filters.type}
            onValueChange={(value) => onUpdateFilters({ type: value as EventType })}
          >
            <SelectTrigger>
              <SelectValue placeholder="Event Type" />
            </SelectTrigger>
            <SelectContent>
              {EVENT_TYPES.map((type) => (
                <SelectItem key={type} value={type}>
                  {type}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>

        {/* Severity Filter */}
        <div className="w-[200px]">
          <Select
            value={filters.severity}
            onValueChange={(value) =>
              onUpdateFilters({ severity: value as typeof SEVERITY_LEVELS[number] })
            }
          >
            <SelectTrigger>
              <SelectValue placeholder="Severity" />
            </SelectTrigger>
            <SelectContent>
              {SEVERITY_LEVELS.map((severity) => (
                <SelectItem key={severity} value={severity}>
                  {severity.charAt(0).toUpperCase() + severity.slice(1)}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>

        {/* Status Filter */}
        <div className="w-[200px]">
          <Select
            value={filters.status}
            onValueChange={(value) =>
              onUpdateFilters({ status: value as typeof STATUS_OPTIONS[number] })
            }
          >
            <SelectTrigger>
              <SelectValue placeholder="Status" />
            </SelectTrigger>
            <SelectContent>
              {STATUS_OPTIONS.map((status) => (
                <SelectItem key={status} value={status}>
                  {status.charAt(0).toUpperCase() + status.slice(1)}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>

        {/* Source Filter */}
        <div className="w-[200px]">
          <Input
            placeholder="Source IP"
            value={filters.source || ""}
            onChange={(e) => onUpdateFilters({ source: e.target.value })}
          />
        </div>

        {/* Date Range Filter */}
        <div className="flex gap-2">
          <Popover>
            <PopoverTrigger asChild>
              <Button variant="outline" className="w-[200px] justify-start text-left font-normal">
                <CalendarIcon className="mr-2 h-4 w-4" />
                {filters.startDate ? (
                  format(filters.startDate, "PPP")
                ) : (
                  <span>Start Date</span>
                )}
              </Button>
            </PopoverTrigger>
            <PopoverContent className="w-auto p-0">
              <Calendar
                mode="single"
                selected={filters.startDate}
                onSelect={(date) => onUpdateFilters({ startDate: date })}
                initialFocus
              />
            </PopoverContent>
          </Popover>

          <Popover>
            <PopoverTrigger asChild>
              <Button variant="outline" className="w-[200px] justify-start text-left font-normal">
                <CalendarIcon className="mr-2 h-4 w-4" />
                {filters.endDate ? (
                  format(filters.endDate, "PPP")
                ) : (
                  <span>End Date</span>
                )}
              </Button>
            </PopoverTrigger>
            <PopoverContent className="w-auto p-0">
              <Calendar
                mode="single"
                selected={filters.endDate}
                onSelect={(date) => onUpdateFilters({ endDate: date })}
                initialFocus
              />
            </PopoverContent>
          </Popover>
        </div>

        {/* Clear Filters Button */}
        <Button variant="outline" onClick={onClearFilters}>
          Clear Filters
        </Button>
      </div>
    </div>
  );
}
