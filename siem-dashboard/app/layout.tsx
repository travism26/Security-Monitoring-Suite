import './globals.css'
import { Inter } from 'next/font/google'
import { AuthProvider } from './contexts/AuthContext'
import { TeamProvider } from './contexts/TeamContext'

const inter = Inter({ subsets: ['latin'] })

export const metadata = {
  title: 'SIEM Dashboard',
  description: 'Security Information and Event Management Dashboard',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en">
      <body className={inter.className}>
        <AuthProvider>
          <TeamProvider>
            {children}
          </TeamProvider>
        </AuthProvider>
      </body>
    </html>
  )
}

