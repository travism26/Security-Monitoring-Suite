'use client'

import { Card, CardHeader, CardTitle } from "@/components/ui/card"
import { AlertList } from "./AlertList"

export function Alert() {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Security Alerts</CardTitle>
      </CardHeader>
      <AlertList />
    </Card>
  )
}
