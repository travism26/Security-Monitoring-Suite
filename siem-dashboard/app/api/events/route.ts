import { NextResponse } from 'next/server'

export async function GET() {
  // Placeholder data
  const events = [
    { id: 1, timestamp: '2023-06-10 10:30:15', type: 'Login', source: '192.168.1.100', description: 'Successful login' },
    { id: 2, timestamp: '2023-06-10 10:35:22', type: 'File Access', source: '192.168.1.101', description: 'Unauthorized file access attempt' },
    { id: 3, timestamp: '2023-06-10 10:40:05', type: 'Network', source: '192.168.1.102', description: 'Unusual outbound traffic detected' },
  ]

  return NextResponse.json(events)
}

