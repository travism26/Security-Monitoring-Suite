'use client'

import { Alert as AlertUI, AlertDescription, AlertTitle } from "@/components/ui/alert"
import { Button } from "@/components/ui/button"
import { Card, CardContent } from "@/components/ui/card"
import { AlertCircle, CheckCircle, Clock, HelpCircle, Search } from 'lucide-react'
import { Alert } from "@/app/hooks/useAlerts"

interface AlertItemProps {
  alert: Alert;
  onStatusChange: (alertId: string, status: Alert['status'], assignedTo?: string) => Promise<void>;
}

const severityIcons = {
  critical: <AlertCircle className="h-5 w-5 text-red-600" />,
  high: <AlertCircle className="h-5 w-5 text-orange-500" />,
  medium: <AlertCircle className="h-5 w-5 text-yellow-500" />,
  low: <AlertCircle className="h-5 w-5 text-blue-500" />
}

const statusIcons = {
  new: <AlertCircle className="h-4 w-4" />,
  acknowledged: <Clock className="h-4 w-4" />,
  investigating: <Search className="h-4 w-4" />,
  resolved: <CheckCircle className="h-4 w-4" />,
  false_positive: <HelpCircle className="h-4 w-4" />
}

export function AlertItem({ alert, onStatusChange }: AlertItemProps) {
  const handleStatusChange = async (newStatus: Alert['status']) => {
    try {
      await onStatusChange(alert.id, newStatus);
    } catch (error) {
      console.error('Failed to update alert status:', error);
    }
  };

  return (
    <Card className="mb-4">
      <CardContent className="pt-6">
        <AlertUI variant={alert.severity === 'critical' ? 'destructive' : 'default'} className="mb-2">
          <div className="flex items-start gap-4">
            <div className="flex-shrink-0">
              {severityIcons[alert.severity]}
            </div>
            <div className="flex-grow">
              <AlertTitle className="flex items-center gap-2 text-lg font-semibold">
                {alert.title}
                <span className="text-sm font-normal text-muted-foreground">
                  ({alert.category})
                </span>
              </AlertTitle>
              <AlertDescription className="mt-2">
                <p className="text-sm text-muted-foreground mb-2">{alert.description}</p>
                <div className="flex flex-wrap gap-2 mt-4">
                  <div className="flex items-center gap-1 text-sm">
                    {statusIcons[alert.status]}
                    <span className="capitalize">{alert.status}</span>
                  </div>
                  <span className="text-sm text-muted-foreground">
                    {new Date(alert.timestamp).toLocaleString()}
                  </span>
                  {alert.assignedTo && (
                    <span className="text-sm text-muted-foreground">
                      Assigned to: {alert.assignedTo}
                    </span>
                  )}
                </div>
                <div className="flex gap-2 mt-4">
                  {alert.status === 'new' && (
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => handleStatusChange('acknowledged')}
                    >
                      Acknowledge
                    </Button>
                  )}
                  {['new', 'acknowledged'].includes(alert.status) && (
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => handleStatusChange('investigating')}
                    >
                      Investigate
                    </Button>
                  )}
                  {alert.status === 'investigating' && (
                    <>
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => handleStatusChange('resolved')}
                      >
                        Mark Resolved
                      </Button>
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => handleStatusChange('false_positive')}
                      >
                        False Positive
                      </Button>
                    </>
                  )}
                </div>
              </AlertDescription>
            </div>
          </div>
        </AlertUI>
      </CardContent>
    </Card>
  )
}
