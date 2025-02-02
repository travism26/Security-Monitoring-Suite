'use client'

import { useState } from 'react'
import Link from 'next/link'
import { useAuth } from '../contexts/AuthContext'
import { Card, CardContent, CardHeader, CardTitle, CardDescription, CardFooter } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import { Alert, AlertDescription } from "@/components/ui/alert"
import { Separator } from "@/components/ui/separator"
import { FaGoogle, FaMicrosoft } from 'react-icons/fa'

export default function LoginPage() {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const { login } = useAuth()

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    await login(email, password)
  }

  return (
    <div className="flex items-center justify-center min-h-screen bg-gradient-to-r from-blue-400 to-purple-500">
      <div className="absolute inset-0 bg-black opacity-50"></div>
      <Card className="w-[400px] z-10">
        <CardHeader className="space-y-1">
          <CardTitle className="text-2xl text-center">Login to SIEM Dashboard</CardTitle>
          <CardDescription className="text-center">Enter your credentials to access the dashboard</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="space-y-2">
              <Input
                id="email"
                type="email"
                placeholder="Email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
              />
            </div>
            <div className="space-y-2">
              <Input
                id="password"
                type="password"
                placeholder="Password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
              />
            </div>
            <Button type="submit" className="w-full">Login</Button>
          </form>
          <Separator />
          <div className="space-y-2">
            <Button variant="outline" className="w-full" onClick={() => console.log('Google login')}>
              <FaGoogle className="mr-2" /> Login with Google
            </Button>
            <Button variant="outline" className="w-full" onClick={() => console.log('Microsoft login')}>
              <FaMicrosoft className="mr-2" /> Login with Microsoft
            </Button>
          </div>
        </CardContent>
        <CardFooter>
          <p className="text-sm text-center w-full">
            Don't have an account? <Link href="/signup" className="text-blue-500 hover:underline">Sign up</Link>
          </p>
        </CardFooter>
      </Card>
      <Alert className="fixed bottom-4 right-4 w-auto">
        <AlertDescription>
          Test Credentials:<br />
          Email: test@example.com<br />
          Password: password123
        </AlertDescription>
      </Alert>
    </div>
  )
}

