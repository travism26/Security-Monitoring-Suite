'use client'

import { useState } from 'react'
import Link from 'next/link'
import { Home, Shield, AlertCircle, Settings, LogOut, Menu } from 'lucide-react'

export function Sidebar() {
  const [isOpen, setIsOpen] = useState(false)

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
          <li>
            <Link href="/" className="flex flex-row items-center h-12 transform hover:translate-x-2 transition-transform ease-in duration-200 text-gray-400 hover:text-white">
              <span className="inline-flex items-center justify-center h-12 w-12 text-lg text-gray-400"><Home /></span>
              <span className="text-sm font-medium">Dashboard</span>
            </Link>
          </li>
          <li>
            <Link href="/events" className="flex flex-row items-center h-12 transform hover:translate-x-2 transition-transform ease-in duration-200 text-gray-400 hover:text-white">
              <span className="inline-flex items-center justify-center h-12 w-12 text-lg text-gray-400"><Shield /></span>
              <span className="text-sm font-medium">Events</span>
            </Link>
          </li>
          <li>
            <Link href="/alerts" className="flex flex-row items-center h-12 transform hover:translate-x-2 transition-transform ease-in duration-200 text-gray-400 hover:text-white">
              <span className="inline-flex items-center justify-center h-12 w-12 text-lg text-gray-400"><AlertCircle /></span>
              <span className="text-sm font-medium">Alerts</span>
            </Link>
          </li>
          <li>
            <Link href="/settings" className="flex flex-row items-center h-12 transform hover:translate-x-2 transition-transform ease-in duration-200 text-gray-400 hover:text-white">
              <span className="inline-flex items-center justify-center h-12 w-12 text-lg text-gray-400"><Settings /></span>
              <span className="text-sm font-medium">Settings</span>
            </Link>
          </li>
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

