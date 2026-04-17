import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom'
import { ProtectedRoute } from './components/protected-route'
import LoginPage from './pages/login'
import RegisterPage from './pages/register'
import DashboardPage from './pages/dashboard'
import UsersPage from './pages/users'
import ClientsPage from './pages/clients'
import FirmSettingsPage from './pages/firm-settings'
import AcceptInvitationPage from './pages/accept-invitation'
import './styles/tokens.css'

const queryClient = new QueryClient()

export default function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <Routes>
          <Route path="/login" element={<LoginPage />} />
          <Route path="/register" element={<RegisterPage />} />
          <Route path="/accept-invitation" element={<AcceptInvitationPage />} />
          <Route
            path="/dashboard"
            element={<ProtectedRoute><DashboardPage /></ProtectedRoute>}
          />
          <Route
            path="/users"
            element={<ProtectedRoute><UsersPage /></ProtectedRoute>}
          />
          <Route
            path="/clients"
            element={<ProtectedRoute><ClientsPage /></ProtectedRoute>}
          />
          <Route
            path="/settings"
            element={<ProtectedRoute><FirmSettingsPage /></ProtectedRoute>}
          />
          <Route path="/" element={<Navigate to="/dashboard" replace />} />
          <Route path="*" element={<Navigate to="/dashboard" replace />} />
        </Routes>
      </BrowserRouter>
    </QueryClientProvider>
  )
}
