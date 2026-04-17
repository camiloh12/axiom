import { useEffect, useState, type FormEvent } from 'react'
import { Layout } from '../components/layout'
import { api } from '../api/client'
import './pages.css'

interface Firm {
  id: string
  name: string
  slug: string
  timezone?: string
  billing_contact_email?: string
}

export default function FirmSettingsPage() {
  const [form, setForm] = useState({ name: '', timezone: '', billing_contact_email: '' })
  const [error, setError] = useState<string | null>(null)
  const [notice, setNotice] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    api<Firm>('/firms/current')
      .then((firm) => setForm({
        name: firm.name ?? '',
        timezone: firm.timezone ?? '',
        billing_contact_email: firm.billing_contact_email ?? '',
      }))
      .catch(() => setError('Could not load firm profile.'))
      .finally(() => setLoading(false))
  }, [])

  async function onSave(e: FormEvent) {
    e.preventDefault()
    setError(null)
    setNotice(null)
    try {
      await api<Firm>('/firms/current', {
        method: 'PATCH',
        body: JSON.stringify({
          name: form.name || null,
          timezone: form.timezone || null,
          billing_contact_email: form.billing_contact_email || null,
        }),
      })
      setNotice('Firm profile saved.')
    } catch {
      setError('Save failed. Please try again.')
    }
  }

  return (
    <Layout>
      <div className="page-header">
        <h1>Firm settings</h1>
        <p>Billing contact, timezone, and firm branding.</p>
      </div>

      <section className="card">
        <h2>Profile</h2>
        <p>This is how Axiom addresses your firm in emails and reports.</p>

        {loading ? (
          <p className="empty">Loading…</p>
        ) : (
          <form onSubmit={onSave}>
            <div className="form-grid">
              <div className="field full">
                <label className="label" htmlFor="fn">Firm name</label>
                <input id="fn" className="input" value={form.name}
                       onChange={(e) => setForm((f) => ({ ...f, name: e.target.value }))} required />
              </div>
              <div className="field">
                <label className="label" htmlFor="tz">Timezone</label>
                <input id="tz" className="input" value={form.timezone}
                       onChange={(e) => setForm((f) => ({ ...f, timezone: e.target.value }))} />
              </div>
              <div className="field">
                <label className="label" htmlFor="be">Billing contact email</label>
                <input id="be" type="email" className="input" value={form.billing_contact_email}
                       onChange={(e) => setForm((f) => ({ ...f, billing_contact_email: e.target.value }))} />
              </div>
            </div>

            {error && <div role="alert" className="error-text" style={{ marginTop: 16 }}>{error}</div>}
            {notice && <div role="status" className="success-text" style={{ marginTop: 16 }}>{notice}</div>}

            <div className="form-actions">
              <button type="submit" className="btn btn-primary">Save changes</button>
            </div>
          </form>
        )}
      </section>
    </Layout>
  )
}
