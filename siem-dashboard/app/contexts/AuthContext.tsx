'use client'

import React, { createContext, useState, useContext, useEffect } from "react"
import { useRouter } from "next/navigation"
import { AuthService, User } from "../services/auth.service"

interface AuthContextType {
  user: User | null
  loading: boolean
  login: (email: string, password: string) => Promise<void>
  logout: () => void
  signup: (email: string, password: string, firstName: string, lastName: string) => Promise<void>
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const router = useRouter()

  const [loading, setLoading] = useState(true)

  useEffect(() => {
    console.log("[Auth Context] Initializing auth context");
    const checkSession = async () => {
      console.log("[Auth Context] Checking for existing session");
      try {
        const currentUser = await AuthService.getCurrentUser();
        if (currentUser) {
          console.log("[Auth Context] Session validated successfully");
          setUser(currentUser);
        } else {
          console.log("[Auth Context] No valid session found");
        }
      } catch (error) {
        console.error("[Auth Context] Error validating session:", error);
      } finally {
        setLoading(false);
      }
    };

    checkSession();
  }, [])

  const login = async (email: string, password: string) => {
    console.log("[Auth Context] Attempting login");
    try {
      const { user } = await AuthService.login(email, password);
      console.log("[Auth Context] Login successful");
      setUser(user);
      console.log("[Auth Context] Navigating to dashboard");
      router.push("/dashboard");
    } catch (error) {
      console.error("[Auth Context] Login failed:", error);
      if (error instanceof Error) {
        throw new Error(error.message);
      }
      throw new Error("Failed to login");
    }
  }

  const logout = () => {
    console.log("[Auth Context] Logging out user");
    setUser(null);
    console.log("[Auth Context] Session cleared, redirecting to login");
    router.push("/login");
  }

  const signup = async (email: string, password: string, firstName: string, lastName: string) => {
    console.log("[Auth Context] Attempting signup");
    try {
      const { user } = await AuthService.signup(
        email,
        password,
        firstName,
        lastName
      );
      console.log("[Auth Context] Signup successful");
      setUser(user);
      console.log("[Auth Context] Navigating to dashboard");
      router.push("/dashboard");
    } catch (error) {
      console.error("[Auth Context] Signup failed:", error);
      if (error instanceof Error) {
        throw new Error(error.message);
      }
      throw new Error("Failed to signup");
    }
  }

  return (
    <AuthContext.Provider value={{ user, loading, login, logout, signup }}>
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
