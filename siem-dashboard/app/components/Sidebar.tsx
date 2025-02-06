'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { Home, BarChart2, Shield, Bell, Settings, AlertOctagon, LogOut, Key } from 'lucide-react'
import { cn } from "@/lib/utils"
import {
  Sidebar,
  SidebarContent,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuItem,
  SidebarMenuButton,
  SidebarTrigger,
  SidebarFooter,
} from "@/components/ui/sidebar"
import { useAuth } from '../contexts/AuthContext'
import { Button } from "@/components/ui/button"

const menuItems = [
  { icon: Home, label: 'Dashboard', href: '/dashboard' },
  { icon: BarChart2, label: 'Analytics', href: '/analytics' },
  { icon: Shield, label: 'Threats', href: '/threats' },
  { icon: Bell, label: 'Alerts', href: '/alerts' },
  { icon: AlertOctagon, label: 'Incident Response', href: '/incident-response' },
  { icon: Key, label: 'API Keys', href: '/api-keys' },
  { icon: Settings, label: 'Settings', href: '/settings' },
]

export function SidebarNav() {
  const pathname = usePathname()
  const { user, logout } = useAuth()

  return (
    <Sidebar>
      <SidebarHeader>
        <SidebarTrigger className="absolute right-2 top-2 md:hidden" />
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton size="lg" asChild>
              <Link href="/dashboard">
                <div className="flex aspect-square size-8 items-center justify-center rounded-lg bg-primary text-primary-foreground">
                  <Shield className="size-4" />
                </div>
                <div className="flex flex-col gap-0.5 leading-none">
                  <span className="font-semibold">SIEM Dashboard</span>
                  <span className="text-xs text-muted-foreground">Security at a glance</span>
                </div>
              </Link>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>
      <SidebarContent>
        <SidebarMenu>
          {menuItems.map((item) => (
            <SidebarMenuItem key={item.href}>
              <SidebarMenuButton asChild isActive={pathname === item.href}>
                <Link href={item.href}>
                  <item.icon className="mr-2 h-4 w-4" />
                  {item.label}
                </Link>
              </SidebarMenuButton>
            </SidebarMenuItem>
          ))}
        </SidebarMenu>
      </SidebarContent>
      <SidebarFooter>
        <div className="p-4">
          <div className="flex items-center justify-between">
            <span className="text-sm">{user?.firstName}</span>
            <Button variant="ghost" size="sm" onClick={logout}>
              <LogOut className="mr-2 h-4 w-4" />
              Logout
            </Button>
          </div>
        </div>
      </SidebarFooter>
    </Sidebar>
  )
}
