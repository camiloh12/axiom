import { useEffect, useState, type FormEvent } from 'react'
import { Layout } from '../components/layout'
import { api } from '../api/client'
import './pages.css'

interface Client {
  id: string
  name: string
  industry?: string
  primary_contact_email?: string
}

export default function ClientsPage() {
  const [clients, setClients] = useState<Client[]>([])
  const [form, setForm] = useState({ name: '', industry: '', primary_contact_email: '' })
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)

  async function reload() {
    try {
      const res = await api<{ items: Client[] }>('/clients')
      setClients(res.items)
    } catch {
      setError('Failed to load clients.')
    } finally {
      setLoading(false)
    }
  }

  // eslint-disable-next-line react-hooks/set-state-in-effect -- setState is awaited inside reload(), not synchronous
  useEffect(() => { void reload() }, [])

  async function onCreate(e: FormEvent) {
    e.preventDefault()
    setError(null)
    if (!form.name.trim()) {
      setError('Client name is required.')
      return
    }
    try {
      await api<Client>('/clients', { method: 'POST', body: JSON.stringify(form) })
      setForm({ name: '', industry: '', primary_contact_email: '' })
      await reload()
    } catch {
      setError('Could not create client.')
    }
  }

  return (
    <Layout>
      <div className="page-header">
        <h1>Clients</h1>
        <p>Client records anchor engagements, documents, and audit history.</p>
      </div>

      <section className="card" aria-label="Add client">
        <h2>Add a client</h2>
        <p>You can edit industry and contact details later.</p>
        <form onSubmit={onCreate}>
          <div className="form-grid">
            <div className="field full">
              <label className="label" htmlFor="c-name">Client name</label>
              <input id="c-name" className="input" value={form.name}
                     onChange={(e) => setForm((f) => ({ ...f, name: e.target.value }))} required />
            </div>
            <div className="field">
              <label className="label" htmlFor="c-ind">Industry</label>
              <input id="c-ind" className="input" value={form.industry}
                     onChange={(e) => setForm((f) => ({ ...f, industry: e.target.value }))} />
            </div>
            <div className="field">
              <label className="label" htmlFor="c-email">Primary contact email</label>
              <input id="c-email" type="email" className="input" value={form.primary_contact_email}
                     onChange={(e) => setForm((f) => ({ ...f, primary_contact_email: e.target.value }))} />
            </div>
          </div>
          {error && <div role="alert" className="error-text" style={{ marginTop: 16 }}>{error}</div>}
          <div className="form-actions">
            <button type="submit" className="btn btn-primary">Add client</button>
          </div>
        </form>
      </section>

      <section className="card" aria-label="Client list">
        <h2>All clients</h2>
        {loading ? (
          <p className="empty">Loading…</p>
        ) : clients.length === 0 ? (
          <p className="empty">No clients yet. Add your first above.</p>
        ) : (
          <table className="table">
            <thead>
              <tr><th>Name</th><th>Industry</th><th>Primary contact</th></tr>
            </thead>
            <tbody>
              {clients.map((c) => (
                <tr key={c.id}>
                  <td>{c.name}</td>
                  <td className="muted">{c.industry || '—'}</td>
                  <td className="muted">{c.primary_contact_email || '—'}</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </section>
    </Layout>
  )
}
