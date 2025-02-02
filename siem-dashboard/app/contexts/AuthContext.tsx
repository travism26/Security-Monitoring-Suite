'use client'

import React, { createContext, useState, useContext, useEffect } from "react"
import { useRouter } from "next/navigation"
import { AuthService, User } from "../services/auth.service"

interface AuthContextType {
  user: User | null
  login: (email: string, password: string) => Promise<void>
  logout: () => void
  signup: (email: string, password: string, firstName: string, lastName: string) => Promise<void>
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const router = useRouter()

  useEffect(() => {
    // Check for existing session
    const storedUser = localStorage.getItem("user")
    if (storedUser) {
      setUser(JSON.parse(storedUser))
    }
  }, [])

  const login = async (email: string, password: string) => {
    try {
      const { user, token } = await AuthService.login(email, password)
      setUser(user)
      localStorage.setItem("user", JSON.stringify(user))
      localStorage.setItem("token", token)
      router.push("/dashboard")
    } catch (error) {
      if (error instanceof Error) {
        throw new Error(error.message)
      }
      throw new Error("Failed to login")
    }
  }

  const logout = () => {
    setUser(null)
    localStorage.removeItem("user")
    localStorage.removeItem("token")
    router.push("/login")
  }

  const signup = async (email: string, password: string, firstName: string, lastName: string) => {
    try {
      const { user, token } = await AuthService.signup({
        email,
        password,
        firstName,
        lastName,
      })
      setUser(user)
      localStorage.setItem("user", JSON.stringify(user))
      localStorage.setItem("token", token)
      router.push("/dashboard")
    } catch (error) {
      if (error instanceof Error) {
        throw new Error(error.message)
      }
      throw new Error("Failed to signup")
    }
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
    throw new Error("useAuth must be used within an AuthProvider")
  }
  return context
}
