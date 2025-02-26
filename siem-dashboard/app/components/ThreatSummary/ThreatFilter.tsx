"use client";

import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

interface ThreatFilterProps {
  filter: string;
  sortBy: "severity" | "count" | undefined;
  onFilterChange: (value: string) => void;
  onSortChange: (value: "severity" | "count") => void;
}

export function ThreatFilter({
  filter,
  sortBy,
  onFilterChange,
  onSortChange,
}: ThreatFilterProps) {
  return (
    <div className="flex gap-4 mb-6">
      <div className="flex-1">
        <Input
          placeholder="Filter threats..."
          value={filter}
          onChange={(e) => onFilterChange(e.target.value)}
          className="max-w-sm"
        />
      </div>
      <Select
        value={sortBy}
        onValueChange={(value) =>
          onSortChange(value as "severity" | "count")
        }
      >
        <SelectTrigger className="w-[180px]">
          <SelectValue placeholder="Sort by..." />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="severity">Sort by Severity</SelectItem>
          <SelectItem value="count">Sort by Count</SelectItem>
        </SelectContent>
      </Select>
    </div>
  );
}
