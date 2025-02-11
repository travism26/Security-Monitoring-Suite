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
  key: string;
  description: string;
  permissions: string[];
  createdAt: Date;
  expiresAt?: Date;
  isActive: boolean;
}

export default function ApiKeysPage() {
  const { user } = useAuth() as { user: { id: string } | null }
  const router = useRouter()
  const [apiKeys, setApiKeys] = useState<ApiKey[]>([])
  const [newKeyDescription, setNewKeyDescription] = useState('')
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    if (!user) {
      router.push('/login')
    } else {
      // Fetch API keys
      fetchApiKeys()
    }
  }, [user, router])

  const fetchApiKeys = async () => {
    try {
      const response = await fetch(`/gateway/api/v1/users/${user?.id}/api-keys`);
      if (!response.ok) throw new Error('Failed to fetch API keys');
      const data = await response.json();
      setApiKeys(data);
    } catch (error) {
      toast({
        title: "Error",
        description: "Failed to load API keys",
        variant: "destructive",
      });
    }
  }

  const createNewApiKey = async () => {
    if (!newKeyDescription.trim()) {
      toast({
        title: "Error",
        description: "Please enter a name for the new API key.",
        variant: "destructive",
      })
      return
    }

    setLoading(true);
    try {
      const response = await fetch(`/gateway/api/v1/users/${user?.id}/api-keys`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          description: newKeyDescription.trim(),
        }),
      });

      if (!response.ok) throw new Error('Failed to create API key');

      const newKey = await response.json();
      setApiKeys([...apiKeys, newKey]);
      setNewKeyDescription('');
      toast({
        title: "Success",
        description: "New API key created successfully.",
      });
    } catch (error) {
      toast({
        title: "Error",
        description: "Failed to create API key",
        variant: "destructive",
      });
    } finally {
      setLoading(false);
    }
  }

  const revokeApiKey = async (keyId: string) => {
    try {
      const response = await fetch(
        `/gateway/api/v1/users/${user?.id}/api-keys/${keyId}/revoke`,
        {
          method: 'PUT',
        }
      );

      if (!response.ok) throw new Error('Failed to revoke API key');

      // Update the local state to reflect the change
      setApiKeys(apiKeys.map(key => 
        key.key === keyId ? { ...key, isActive: false } : key
      ));

      toast({
        title: "Success",
        description: "API key revoked successfully.",
      });
    } catch (error) {
      toast({
        title: "Error",
        description: "Failed to revoke API key",
        variant: "destructive",
      });
    }
  }

  const copyToClipboard = async (text: string) => {
    try {
      await navigator.clipboard.writeText(text);
      toast({
        title: "Success",
        description: "API key copied to clipboard",
      });
    } catch (error) {
      toast({
        title: "Error",
        description: "Failed to copy API key",
        variant: "destructive",
      });
    }
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
                    placeholder="Enter API key description"
                    value={newKeyDescription}
                    onChange={(e) => setNewKeyDescription(e.target.value)}
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
                      <TableHead>Description</TableHead>
                      <TableHead>Key</TableHead>
                      <TableHead>Created</TableHead>
                      <TableHead>Expires</TableHead>
                      <TableHead>Status</TableHead>
                      <TableHead>Actions</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {apiKeys.map((apiKey) => (
                      <TableRow key={apiKey.key}>
                        <TableCell>{apiKey.description}</TableCell>
                        <TableCell>
                          <code className="relative rounded bg-muted px-[0.3rem] py-[0.2rem] font-mono text-sm">
                            {apiKey.key.substring(0, 10)}...
                          </code>
                        </TableCell>
                        <TableCell>{new Date(apiKey.createdAt).toLocaleDateString()}</TableCell>
                        <TableCell>
                          {apiKey.expiresAt ? new Date(apiKey.expiresAt).toLocaleDateString() : "Never"}
                        </TableCell>
                        <TableCell>
                          <span
                            className={`inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium ${
                              apiKey.isActive
                                ? "bg-green-100 text-green-800"
                                : "bg-red-100 text-red-800"
                            }`}
                          >
                            {apiKey.isActive ? "Active" : "Revoked"}
                          </span>
                        </TableCell>
                        <TableCell>
                          <div className="flex gap-2">
                            <Button
                              variant="outline"
                              size="sm"
                              onClick={() => copyToClipboard(apiKey.key)}
                            >
                              Copy
                            </Button>
                            {apiKey.isActive && (
                              <Button
                                variant="destructive"
                                size="sm"
                                onClick={() => revokeApiKey(apiKey.key)}
                              >
                                Revoke
                              </Button>
                            )}
                          </div>
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
