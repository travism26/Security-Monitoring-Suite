import { NextResponse } from "next/server";

export async function GET() {
  // Placeholder data
  const systems = [
    { name: "Firewall", status: "Operational" },
    { name: "IDS", status: "Operational" },
    { name: "Log Server", status: "Down" },
    { name: "Email Filter", status: "Operational" },
  ];

  return NextResponse.json(systems);
}
