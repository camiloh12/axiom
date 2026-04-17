import { create } from 'zustand'
import { api, setTokens, clearTokens, hasTokens } from '../api/client'

export interface User {
  id: string
  email: string
  display_name: string
  role: string
}

export interface Firm {
  id: string
  name: string
  slug: string
}

export interface RegisterInput {
  firm_name: string
  admin_email: string
  admin_name: string
  password: string
  country: string
  staff_count_range: string
  primary_audit_types: string[]
}

interface AuthState {
  user: User | null
  firm: Firm | null
  isAuthenticated: boolean
  register: (data: RegisterInput) => Promise<void>
  login: (email: string, password: string) => Promise<void>
  logout: () => void
  loadProfile: () => Promise<void>
}

export const useAuth = create<AuthState>((set) => ({
  user: null,
  firm: null,
  isAuthenticated: hasTokens(),

  register: async (data) => {
    const res = await api<{ access_token: string; refresh_token: string; user: User; firm: Firm }>(
      '/auth/register',
      { method: 'POST', body: JSON.stringify(data) },
    )
    setTokens(res.access_token, res.refresh_token)
    set({ user: res.user, firm: res.firm, isAuthenticated: true })
  },

  login: async (email, password) => {
    const res = await api<{ access_token: string; refresh_token: string }>('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    })
    setTokens(res.access_token, res.refresh_token)
    set({ isAuthenticated: true })
  },

  logout: () => {
    clearTokens()
    set({ user: null, firm: null, isAuthenticated: false })
  },

  loadProfile: async () => {
    const [user, firm] = await Promise.all([
      api<User>('/users/me'),
      api<Firm>('/firms/current'),
    ])
    set({ user, firm })
  },
}))
