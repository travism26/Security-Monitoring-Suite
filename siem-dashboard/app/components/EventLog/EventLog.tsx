"use client";

import { useEventLog } from "@/app/hooks/useEventLog";
import { EventLogTable } from "./EventLogTable";
import { EventLogSkeleton } from "./EventLogSkeleton";
import { EventLogFilters } from "./EventLogFilters";
import { Button } from "@/components/ui/button";
import {
  ChevronLeftIcon,
  ChevronRightIcon,
  ReloadIcon,
} from "@radix-ui/react-icons";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { AlertCircle } from "lucide-react";

interface EventLogProps {
  refreshInterval?: number;
}

export function EventLog({ refreshInterval }: EventLogProps) {
  const {
    events,
    loading,
    error,
    filters,
    currentPage,
    totalPages,
    updateFilters,
    clearFilters,
    goToPage,
    isMockData,
  } = useEventLog({
    refreshInterval,
    useMockData: process.env.NODE_ENV === "development",
  });

  return (
    <div className="space-y-4">
      {/* Mock Data Indicator */}
      {isMockData && (
        <Alert>
          <AlertCircle className="h-4 w-4" />
          <AlertDescription>
            Using mock data. Real API integration pending.
          </AlertDescription>
        </Alert>
      )}

      {/* Filters Section */}
      <EventLogFilters
        filters={filters}
        onUpdateFilters={updateFilters}
        onClearFilters={clearFilters}
      />

      {/* Error Message */}
      {error && (
        <Alert variant="destructive">
          <AlertCircle className="h-4 w-4" />
          <AlertDescription>{error.message}</AlertDescription>
        </Alert>
      )}

      {/* Event Log Table */}
      <div className="rounded-md border">
        {loading ? <EventLogSkeleton /> : <EventLogTable events={events} />}
      </div>

      {/* Pagination Controls */}
      {!loading && totalPages > 0 && (
        <div className="flex items-center justify-between">
          <div className="text-sm text-muted-foreground">
            Page {currentPage} of {totalPages}
          </div>
          <div className="flex items-center space-x-2">
            <Button
              variant="outline"
              size="sm"
              onClick={() => goToPage(currentPage - 1)}
              disabled={currentPage === 1}
            >
              <ChevronLeftIcon className="h-4 w-4" />
              Previous
            </Button>
            <Button
              variant="outline"
              size="sm"
              onClick={() => goToPage(currentPage + 1)}
              disabled={currentPage === totalPages}
            >
              Next
              <ChevronRightIcon className="h-4 w-4" />
            </Button>
          </div>
        </div>
      )}
    </div>
  );
}
