'use client'

import { useState } from 'react'
import Link from 'next/link'
import { Home, Shield, AlertCircle, Settings, LogOut, Menu } from 'lucide-react'
import { usePathname } from 'next/navigation'

export function Sidebar() {
  const [isOpen, setIsOpen] = useState(false)
  const pathname = usePathname()

  return (
    <>
      <button
        className="lg:hidden fixed top-4 left-4 z-20 p-2 bg-gray-800 text-white rounded-md"
        onClick={() => setIsOpen(!isOpen)}
      >
        <Menu />
      </button>
      <div className={`${isOpen ? 'translate-x-0' : '-translate-x-full'} lg:translate-x-0 transition-transform duration-300 ease-in-out fixed inset-y-0 left-0 z-10 w-64 bg-gray-800 overflow-y-auto lg:static lg:block`}>
        <div className="flex items-center justify-center h-20 shadow-md">
          <h1 className="text-3xl uppercase text-white">SIEM XDR</h1>
        </div>
        <ul className="flex flex-col py-4">
          {[
            { href: "/", icon: Home, label: "Dashboard" },
            { href: "/events", icon: Shield, label: "Events" },
            { href: "/alerts", icon: AlertCircle, label: "Alerts" },
            { href: "/settings", icon: Settings, label: "Settings" },
          ].map(({ href, icon: Icon, label }) => (
            <li key={href}>
              <Link 
                href={href} 
                className={`flex flex-row items-center h-12 transform hover:translate-x-2 transition-transform ease-in duration-200 text-gray-400 hover:text-white ${
                  pathname === href ? 'text-white' : ''
                }`}
              >
                <span className="inline-flex items-center justify-center h-12 w-12 text-lg text-gray-400">
                  <Icon />
                </span>
                <span className="text-sm font-medium">{label}</span>
              </Link>
            </li>
          ))}
        </ul>
        <div className="mt-auto pb-4">
          <Link href="/logout" className="flex flex-row items-center h-12 transform hover:translate-x-2 transition-transform ease-in duration-200 text-gray-400 hover:text-white">
            <span className="inline-flex items-center justify-center h-12 w-12 text-lg text-gray-400"><LogOut /></span>
            <span className="text-sm font-medium">Logout</span>
          </Link>
        </div>
      </div>
    </>
  )
}

