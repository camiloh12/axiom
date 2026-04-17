import { Navigate } from 'react-router-dom'
import { useAuth } from '../hooks/use-auth'

export function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const isAuthenticated = useAuth((s) => s.isAuthenticated)
  if (!isAuthenticated) return <Navigate to="/login" replace />
  return <>{children}</>
}
