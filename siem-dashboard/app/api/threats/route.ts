import { NextResponse } from 'next/server'

export async function GET() {
  // Placeholder data
  const threats = [
    { name: 'Malware', count: 15, severity: 70 },
    { name: 'Phishing', count: 8, severity: 60 },
    { name: 'DDoS', count: 3, severity: 40 },
  ]

  return NextResponse.json(threats)
}

