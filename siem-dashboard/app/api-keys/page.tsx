'use client'

import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { useAuth } from '../contexts/AuthContext'
import { SidebarNav } from '../components/Sidebar'
import { SidebarProvider, SidebarInset } from "@/components/ui/sidebar"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { toast } from "@/components/ui/use-toast"

interface ApiKey {
  id: string
  name: string
  key: string
  createdAt: string
}

export default function ApiKeysPage() {
  const { user } = useAuth()
  const router = useRouter()
  const [apiKeys, setApiKeys] = useState<ApiKey[]>([])
  const [newKeyName, setNewKeyName] = useState('')

  useEffect(() => {
    if (!user) {
      router.push('/login')
    } else {
      // Fetch API keys
      fetchApiKeys()
    }
  }, [user, router])

  const fetchApiKeys = async () => {
    // In a real application, you would fetch this from your backend
    const mockApiKeys: ApiKey[] = [
      { id: '1', name: 'Production Server', key: 'prod_abcdefghijklmnop', createdAt: '2023-06-15T10:00:00Z' },
      { id: '2', name: 'Development Server', key: 'dev_qrstuvwxyz123456', createdAt: '2023-06-16T11:30:00Z' },
    ]
    setApiKeys(mockApiKeys)
  }

  const createNewApiKey = async () => {
    if (!newKeyName.trim()) {
      toast({
        title: "Error",
        description: "Please enter a name for the new API key.",
        variant: "destructive",
      })
      return
    }

    // In a real application, you would call your backend to create a new API key
    const newKey: ApiKey = {
      id: String(apiKeys.length + 1),
      name: newKeyName,
      key: `new_${Math.random().toString(36).substr(2, 16)}`,
      createdAt: new Date().toISOString(),
    }

    setApiKeys([...apiKeys, newKey])
    setNewKeyName('')
    toast({
      title: "Success",
      description: "New API key created successfully.",
    })
  }

  const deleteApiKey = async (id: string) => {
    // In a real application, you would call your backend to delete the API key
    setApiKeys(apiKeys.filter(key => key.id !== id))
    toast({
      title: "Success",
      description: "API key deleted successfully.",
    })
  }

  if (!user) {
    return <div className="flex items-center justify-center min-h-screen">Loading...</div>
  }

  return (
    <SidebarProvider>
      <div className="flex h-screen overflow-hidden">
        <SidebarNav />
        <SidebarInset className="flex-1 overflow-auto">
          <main className="p-4 md:p-6 bg-background">
            <h1 className="text-3xl font-bold mb-6">
              API Keys
            </h1>
            <Card className="mb-6">
              <CardHeader>
                <CardTitle>Create New API Key</CardTitle>
                <CardDescription>Generate a new API key for your monitoring agents</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="flex space-x-2">
                  <Input
                    placeholder="Enter API key name"
                    value={newKeyName}
                    onChange={(e) => setNewKeyName(e.target.value)}
                  />
                  <Button onClick={createNewApiKey}>Create Key</Button>
                </div>
              </CardContent>
            </Card>
            <Card>
              <CardHeader>
                <CardTitle>Existing API Keys</CardTitle>
                <CardDescription>Manage your existing API keys</CardDescription>
              </CardHeader>
              <CardContent>
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>Name</TableHead>
                      <TableHead>Key</TableHead>
                      <TableHead>Created At</TableHead>
                      <TableHead>Action</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {apiKeys.map((apiKey) => (
                      <TableRow key={apiKey.id}>
                        <TableCell>{apiKey.name}</TableCell>
                        <TableCell>{apiKey.key}</TableCell>
                        <TableCell>{new Date(apiKey.createdAt).toLocaleString()}</TableCell>
                        <TableCell>
                          <Button variant="destructive" onClick={() => deleteApiKey(apiKey.id)}>Delete</Button>
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </CardContent>
            </Card>
          </main>
        </SidebarInset>
      </div>
    </SidebarProvider>
  )
}
