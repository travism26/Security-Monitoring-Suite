'use client'

import React, { createContext, useState, useContext, useEffect } from 'react'
import { useRouter } from 'next/navigation'

interface User {
  id: string
  name: string
  email: string
}

interface AuthContextType {
  user: User | null
  login: (email: string, password: string) => Promise<void>
  logout: () => void
  signup: (email: string, password: string) => Promise<void>
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const router = useRouter()

  useEffect(() => {
    // Check for existing session
    const storedUser = localStorage.getItem('user')
    if (storedUser) {
      setUser(JSON.parse(storedUser))
    }
  }, [])

  const login = async (email: string, password: string) => {
    // This is a mock login. In a real app, you'd call your API here.
    if (email === 'test@example.com' && password === 'password123') {
      const mockUser = { id: '1', name: 'Test User', email: email }
      setUser(mockUser)
      localStorage.setItem('user', JSON.stringify(mockUser))
      router.push('/dashboard')
    } else {
      throw new Error('Invalid credentials')
    }
  }

  const logout = () => {
    setUser(null)
    localStorage.removeItem('user')
    router.push('/login')
  }

  const signup = async (email: string, password: string) => {
    // This is a mock signup. In a real app, you'd call your API here.
    const mockUser = { id: '2', name: 'New User', email: email }
    setUser(mockUser)
    localStorage.setItem('user', JSON.stringify(mockUser))
    router.push('/dashboard')
  }

  return (
    <AuthContext.Provider value={{ user, login, logout, signup }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}

