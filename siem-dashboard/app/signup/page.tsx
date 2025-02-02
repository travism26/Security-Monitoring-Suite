'use client'

import { useState } from "react"
import Link from "next/link"
import { useAuth } from "../contexts/AuthContext"
import { Card, CardContent, CardHeader, CardTitle, CardDescription, CardFooter } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import { Separator } from "@/components/ui/separator"
import { FaGoogle, FaMicrosoft } from "react-icons/fa"
import { Alert, AlertDescription } from "@/components/ui/alert"

export default function SignupPage() {
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [confirmPassword, setConfirmPassword] = useState("")
  const [firstName, setFirstName] = useState("")
  const [lastName, setLastName] = useState("")
  const [error, setError] = useState("")
  const { signup } = useAuth()

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError("")

    if (password !== confirmPassword) {
      setError("Passwords don't match")
      return
    }

    try {
      await signup(email, password, firstName, lastName)
    } catch (error) {
      if (error instanceof Error) {
        setError(error.message)
      } else {
        setError("Failed to create account")
      }
    }
  }

  return (
    <div className="flex items-center justify-center min-h-screen bg-gradient-to-r from-green-400 to-blue-500">
      <div className="absolute inset-0 bg-black opacity-50"></div>
      <Card className="w-[400px] z-10">
        <CardHeader className="space-y-1">
          <CardTitle className="text-2xl text-center">Create an Account</CardTitle>
          <CardDescription className="text-center">Sign up for SIEM Dashboard</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="space-y-2">
              <Input
                id="firstName"
                type="text"
                placeholder="First Name"
                value={firstName}
                onChange={(e) => setFirstName(e.target.value)}
                required
              />
            </div>
            <div className="space-y-2">
              <Input
                id="lastName"
                type="text"
                placeholder="Last Name"
                value={lastName}
                onChange={(e) => setLastName(e.target.value)}
                required
              />
            </div>
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
            <div className="space-y-2">
              <Input
                id="confirmPassword"
                type="password"
                placeholder="Confirm Password"
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
                required
              />
            </div>
            {error && (
              <Alert variant="destructive">
                <AlertDescription>{error}</AlertDescription>
              </Alert>
            )}
            <Button type="submit" className="w-full">Sign Up</Button>
          </form>
          <Separator />
          <div className="space-y-2">
            <Button variant="outline" className="w-full" onClick={() => console.log("Google signup")}>
              <FaGoogle className="mr-2" /> Sign up with Google
            </Button>
            <Button variant="outline" className="w-full" onClick={() => console.log("Microsoft signup")}>
              <FaMicrosoft className="mr-2" /> Sign up with Microsoft
            </Button>
          </div>
        </CardContent>
        <CardFooter>
          <p className="text-sm text-center w-full">
            Already have an account? <Link href="/login" className="text-blue-500 hover:underline">Log in</Link>
          </p>
        </CardFooter>
      </Card>
    </div>
  )
}
