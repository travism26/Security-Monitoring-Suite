import { NextResponse } from 'next/server'

export async function GET() {
  // Placeholder data
  const alerts = [
    { id: 1, message: 'High CPU usage detected', timestamp: new Date().toISOString() },
    { id: 2, message: 'Unusual network activity', timestamp: new Date().toISOString() },
    { id: 3, message: 'Failed login attempts', timestamp: new Date().toISOString() },
  ]

  return NextResponse.json(alerts)
}

