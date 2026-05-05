import { useEffect, useState, type FormEvent } from 'react'
import { Layout } from '../components/layout'
import { api, ApiError } from '../api/client'
import './pages.css'

interface FirmUser {
  id: string
  email: string
  display_name: string
  role: string
}

interface Invitation {
  id: string
  email: string
  assigned_role: string
  status: string
  expires_at: string
  token?: string
}

// ClientAdmin / ClientUser are deliberately not invitable from the firm-staff
// page. Those roles belong to a specific client engagement (users.client_id is
// required by a CHECK constraint) and will be invited via the Client Hub flow,
// which is gated behind the CLIENT_HUB_ENABLED launch-posture flag. Backend
// CreateInvitation enforces the same restriction as defense-in-depth.
const ROLES = ['Partner', 'Manager', 'Staff', 'EQReviewer', 'ViewOnly'] as const

export default function UsersPage() {
  const [users, setUsers] = useState<FirmUser[]>([])
  const [invitations, setInvitations] = useState<Invitation[]>([])
  const [email, setEmail] = useState('')
  const [role, setRole] = useState<typeof ROLES[number]>('Staff')
  const [error, setError] = useState<string | null>(null)
  const [notice, setNotice] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)

  async function reload() {
    try {
      const [u, i] = await Promise.all([
        api<{ items: FirmUser[] }>('/users'),
        api<{ items: Invitation[] }>('/invitations'),
      ])
      setUsers(u.items)
      setInvitations(i.items)
    } catch (err) {
      if (err instanceof ApiError && err.status === 403) {
        // Non-admins cannot list invitations.
        const u = await api<{ items: FirmUser[] }>('/users')
        setUsers(u.items)
      } else {
        setError('Failed to load team.')
      }
    } finally {
      setLoading(false)
    }
  }

  // eslint-disable-next-line react-hooks/set-state-in-effect -- setState is awaited inside reload(), not synchronous
  useEffect(() => { void reload() }, [])

  async function onInvite(e: FormEvent) {
    e.preventDefault()
    setError(null)
    setNotice(null)
    if (!email) {
      setError('Email is required.')
      return
    }
    try {
      const inv = await api<Invitation>('/invitations', {
        method: 'POST',
        body: JSON.stringify({ email, assigned_role: role }),
      })
      setEmail('')
      setNotice(`Invitation sent. Magic link token: ${inv.token}`)
      await reload()
    } catch {
      setError('Could not send invitation.')
    }
  }

  async function onCancel(id: string) {
    try {
      await api(`/invitations/${id}`, { method: 'DELETE' })
      await reload()
    } catch {
      setError('Could not cancel invitation.')
    }
  }

  return (
    <Layout>
      <div className="page-header">
        <h1>Users</h1>
        <p>Invite staff and manage firm access.</p>
      </div>

      <section className="card" aria-label="Invite staff">
        <h2>Invite staff</h2>
        <p>They'll receive a magic link to set a password and sign in.</p>
        <form onSubmit={onInvite}>
          <div className="form-grid">
            <div className="field">
              <label className="label" htmlFor="invite-email">Email</label>
              <input id="invite-email" type="email" className="input"
                     value={email} onChange={(e) => setEmail(e.target.value)} required />
            </div>
            <div className="field">
              <label className="label" htmlFor="invite-role">Role</label>
              <select id="invite-role" className="select"
                      value={role} onChange={(e) => setRole(e.target.value as typeof ROLES[number])}>
                {ROLES.map((r) => <option key={r} value={r}>{r}</option>)}
              </select>
            </div>
          </div>
          {error && <div role="alert" className="error-text" style={{ marginTop: 16 }}>{error}</div>}
          {notice && <div role="status" className="success-text" style={{ marginTop: 16 }}>{notice}</div>}
          <div className="form-actions">
            <button type="submit" className="btn btn-primary">Send invitation</button>
          </div>
        </form>
      </section>

      <section className="card" aria-label="Team members">
        <h2>Team members</h2>
        {loading ? (
          <p className="empty">Loading…</p>
        ) : users.length === 0 ? (
          <p className="empty">No users yet.</p>
        ) : (
          <table className="table">
            <thead><tr><th>Name</th><th>Email</th><th>Role</th></tr></thead>
            <tbody>
              {users.map((u) => (
                <tr key={u.id}>
                  <td>{u.display_name}</td>
                  <td className="muted">{u.email}</td>
                  <td>{u.role}</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </section>

      {invitations.length > 0 && (
        <section className="card" aria-label="Pending invitations">
          <h2>Pending invitations</h2>
          <table className="table">
            <thead><tr><th>Email</th><th>Role</th><th>Status</th><th aria-label="Actions" /></tr></thead>
            <tbody>
              {invitations.map((inv) => (
                <tr key={inv.id}>
                  <td>{inv.email}</td>
                  <td>{inv.assigned_role}</td>
                  <td>
                    <span className={`status-pill status-${inv.status.toLowerCase()}`}>{inv.status}</span>
                  </td>
                  <td>
                    {inv.status === 'Sent' && (
                      <button className="btn btn-ghost" onClick={() => onCancel(inv.id)}>Cancel</button>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </section>
      )}
    </Layout>
  )
}
