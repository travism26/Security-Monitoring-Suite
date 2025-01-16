import { Bell, User } from 'lucide-react'

export function Header() {
  return (
    <header className="flex justify-between items-center py-4 px-6 bg-white border-b-4 border-gray-800">
      <div className="flex items-center">
        <span className="text-gray-800 text-xl font-semibold">Security Operations Center</span>
      </div>
      <div className="flex items-center">
        <button className="flex mx-4 text-gray-600 focus:outline-none">
          <Bell className="h-6 w-6" />
        </button>
        <div className="relative">
          <button className="flex items-center text-gray-600 focus:outline-none">
            <User className="h-6 w-6" />
          </button>
        </div>
      </div>
    </header>
  )
}

