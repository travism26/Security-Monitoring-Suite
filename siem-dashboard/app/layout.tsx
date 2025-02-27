import './globals.css'
import { Inter } from 'next/font/google'
import { AuthProvider } from './contexts/AuthContext'
// Import the EventEmitter configuration to increase max listeners limit
import './utils/eventEmitterConfig'
// Import signal handlers for proper cleanup on termination
import './utils/signalHandlers'
const inter = Inter({ subsets: ['latin'] })

export const metadata = {
  title: 'SIEM Dashboard',
  description: 'Security Information and Event Management Dashboard',
}

const RootLayout = ({
  children,
}: {
  children: React.ReactNode
}) => (
  <html lang="en">
    <body className={inter.className}>
      <AuthProvider>
        {children}
      </AuthProvider>
    </body>
  </html>
)

export default RootLayout
