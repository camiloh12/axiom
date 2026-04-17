import { useEffect, useState, type FormEvent } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { setTokens } from '../api/client'
import { useAuth } from '../hooks/use-auth'
import './auth.css'

interface Invitation {
  email: string
  assigned_role: string
  expires_at: string
}

// Direct fetch (not via the auth'd api() wrapper) — these endpoints are public.
async function publicFetch<T>(path: string, init: RequestInit = {}): Promise<T> {
  const res = await fetch(`/api/v1${path}`, {
    ...init,
    headers: { 'Content-Type': 'application/json', ...(init.headers || {}) },
  })
  if (!res.ok) throw new Error(await res.text())
  return res.json()
}

export default function AcceptInvitationPage() {
  const [params] = useSearchParams()
  const navigate = useNavigate()
  const token = params.get('token') ?? ''
  const [invitation, setInvitation] = useState<Invitation | null>(null)
  const [loadError, setLoadError] = useState<string | null>(null)
  const [displayName, setDisplayName] = useState('')
  const [password, setPassword] = useState('')
  const [submitError, setSubmitError] = useState<string | null>(null)
  const [submitting, setSubmitting] = useState(false)

  useEffect(() => {
    if (!token) return
    publicFetch<Invitation>(`/invitations/validate/${token}`)
      .then(setInvitation)
      .catch(() => setLoadError('This invitation is invalid or has expired.'))
  }, [token])

  // Render the "missing token" state directly rather than setting state in an effect.
  const missingToken = !token

  async function onSubmit(e: FormEvent) {
    e.preventDefault()
    setSubmitError(null)
    if (password.length < 8) {
      setSubmitError('Password must be at least 8 characters.')
      return
    }
    setSubmitting(true)
    try {
      const res = await publicFetch<{
        access_token: string
        refresh_token: string
        user: { id: string; email: string; display_name: string; role: string }
      }>('/invitations/accept', {
        method: 'POST',
        body: JSON.stringify({ token, display_name: displayName, password }),
      })
      setTokens(res.access_token, res.refresh_token)
      useAuth.setState({ user: res.user, isAuthenticated: true })
      navigate('/dashboard', { replace: true })
    } catch {
      setSubmitError('Could not accept invitation. Please request a new link.')
    } finally {
      setSubmitting(false)
    }
  }

  if (missingToken || loadError) {
    return (
      <div className="auth-shell">
        <div className="auth-card">
          <header className="auth-header">
            <h1>Invitation unavailable</h1>
            <p>{loadError ?? 'Missing invitation token.'}</p>
          </header>
        </div>
      </div>
    )
  }

  if (!invitation) {
    return (
      <div className="auth-shell">
        <div className="auth-card">
          <p>Checking invitation…</p>
        </div>
      </div>
    )
  }

  return (
    <div className="auth-shell">
      <div className="auth-card">
        <header className="auth-header">
          <h1>Join your firm on Axiom</h1>
          <p>
            You've been invited as <strong>{invitation.assigned_role}</strong>.
            Sign in details will be tied to <strong>{invitation.email}</strong>.
          </p>
        </header>

        <form onSubmit={onSubmit} noValidate>
          <div className="field">
            <label className="label" htmlFor="display_name">Your name</label>
            <input id="display_name" className="input" value={displayName}
                   onChange={(e) => setDisplayName(e.target.value)} required />
          </div>
          <div className="field">
            <label className="label" htmlFor="password">Create a password</label>
            <input id="password" type="password" className="input" value={password}
                   onChange={(e) => setPassword(e.target.value)} required minLength={8} />
          </div>

          {submitError && <div role="alert" className="error-text">{submitError}</div>}

          <div className="auth-actions">
            <button type="submit" className="btn btn-primary" disabled={submitting}>
              {submitting ? 'Setting up…' : 'Accept and sign in'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}
