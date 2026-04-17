const API_BASE = '/api/v1'

let accessToken: string | null = localStorage.getItem('access_token')
let refreshToken: string | null = localStorage.getItem('refresh_token')

export function setTokens(access: string, refresh: string) {
  accessToken = access
  refreshToken = refresh
  localStorage.setItem('access_token', access)
  localStorage.setItem('refresh_token', refresh)
}

export function clearTokens() {
  accessToken = null
  refreshToken = null
  localStorage.removeItem('access_token')
  localStorage.removeItem('refresh_token')
}

export function hasTokens() {
  return !!accessToken
}

export class ApiError extends Error {
  constructor(public status: number, public body: string) {
    super(`API error ${status}: ${body}`)
  }
}

export async function api<T>(path: string, options: RequestInit = {}): Promise<T> {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(options.headers as Record<string, string>),
  }

  if (accessToken) {
    headers['Authorization'] = `Bearer ${accessToken}`
  }

  const res = await fetch(`${API_BASE}${path}`, { ...options, headers })

  if (res.status === 401 && refreshToken) {
    const refreshRes = await fetch(`${API_BASE}/auth/refresh`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ refresh_token: refreshToken }),
    })
    if (refreshRes.ok) {
      const data = await refreshRes.json()
      setTokens(data.access_token, data.refresh_token)
      headers['Authorization'] = `Bearer ${data.access_token}`
      const retryRes = await fetch(`${API_BASE}${path}`, { ...options, headers })
      if (!retryRes.ok) throw new ApiError(retryRes.status, await retryRes.text())
      if (retryRes.status === 204) return undefined as T
      return retryRes.json()
    }
    clearTokens()
    throw new ApiError(401, 'session expired')
  }

  if (!res.ok) {
    throw new ApiError(res.status, await res.text())
  }

  if (res.status === 204) return undefined as T
  return res.json()
}
