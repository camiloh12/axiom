import { useState, type FormEvent } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { useAuth } from '../hooks/use-auth'
import './auth.css'

const AUDIT_TYPES = [
  { value: 'FinancialAudit', label: 'Financial audits' },
  { value: 'SOC2', label: 'SOC 2' },
  { value: 'ISO27001', label: 'ISO 27001' },
  { value: 'InternalAudit', label: 'Internal audits' },
] as const

const STAFF_RANGES = ['1-10', '11-20', '21-40', '41-60', '60+'] as const

export default function RegisterPage() {
  const register = useAuth((s) => s.register)
  const navigate = useNavigate()
  const [step, setStep] = useState<0 | 1>(0)
  const [form, setForm] = useState({
    firm_name: '',
    admin_email: '',
    admin_name: '',
    password: '',
    country: 'US',
    staff_count_range: '1-10',
    primary_audit_types: [] as string[],
  })
  const [error, setError] = useState<string | null>(null)
  const [submitting, setSubmitting] = useState(false)

  function update<K extends keyof typeof form>(key: K, value: (typeof form)[K]) {
    setForm((f) => ({ ...f, [key]: value }))
  }

  function toggleAuditType(v: string) {
    setForm((f) => ({
      ...f,
      primary_audit_types: f.primary_audit_types.includes(v)
        ? f.primary_audit_types.filter((x) => x !== v)
        : [...f.primary_audit_types, v],
    }))
  }

  function onNext(e: FormEvent) {
    e.preventDefault()
    setError(null)
    if (!form.firm_name || !form.admin_email || !form.admin_name || form.password.length < 8) {
      setError('Fill in every field. Passwords must be at least 8 characters.')
      return
    }
    setStep(1)
  }

  async function onSubmit(e: FormEvent) {
    e.preventDefault()
    setError(null)
    setSubmitting(true)
    try {
      await register(form)
      navigate('/dashboard', { replace: true })
    } catch {
      setError('Registration failed. This email may already be in use.')
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <div className="auth-shell">
      <div className="auth-card auth-card--wide">
        <header className="auth-header">
          <h1>Start your Axiom firm</h1>
          <p>Set up your firm workspace in under a minute.</p>
        </header>

        <div className="stepper" aria-hidden="true">
          <span className={`stepper-dot ${step === 0 ? 'is-active' : 'is-done'}`} />
          <span className={`stepper-dot ${step === 1 ? 'is-active' : ''}`} />
        </div>

        {step === 0 && (
          <form onSubmit={onNext} noValidate aria-label="Account details">
            <div className="field">
              <label className="label" htmlFor="firm_name">Firm name</label>
              <input id="firm_name" className="input" value={form.firm_name}
                     onChange={(e) => update('firm_name', e.target.value)} required />
            </div>
            <div className="field-row">
              <div className="field">
                <label className="label" htmlFor="admin_name">Your name</label>
                <input id="admin_name" className="input" value={form.admin_name}
                       onChange={(e) => update('admin_name', e.target.value)} required />
              </div>
              <div className="field">
                <label className="label" htmlFor="admin_email">Work email</label>
                <input id="admin_email" type="email" className="input" value={form.admin_email}
                       onChange={(e) => update('admin_email', e.target.value)} required />
              </div>
            </div>
            <div className="field">
              <label className="label" htmlFor="password">Password</label>
              <input id="password" type="password" className="input" value={form.password}
                     onChange={(e) => update('password', e.target.value)} required minLength={8} />
            </div>

            {error && <div role="alert" className="error-text">{error}</div>}

            <div className="auth-actions">
              <button type="submit" className="btn btn-primary">Continue</button>
            </div>
          </form>
        )}

        {step === 1 && (
          <form onSubmit={onSubmit} noValidate aria-label="Firm profile">
            <div className="field-row">
              <div className="field">
                <label className="label" htmlFor="country">Country</label>
                <select id="country" className="select" value={form.country}
                        onChange={(e) => update('country', e.target.value)}>
                  <option value="US">United States</option>
                  <option value="CA">Canada</option>
                </select>
              </div>
              <div className="field">
                <label className="label" htmlFor="staff">Staff count</label>
                <select id="staff" className="select" value={form.staff_count_range}
                        onChange={(e) => update('staff_count_range', e.target.value)}>
                  {STAFF_RANGES.map((r) => <option key={r} value={r}>{r}</option>)}
                </select>
              </div>
            </div>

            <div className="field">
              <span className="label">Primary audit types</span>
              <div className="checkbox-group">
                {AUDIT_TYPES.map((a) => (
                  <label key={a.value} className="checkbox-option">
                    <input
                      type="checkbox"
                      checked={form.primary_audit_types.includes(a.value)}
                      onChange={() => toggleAuditType(a.value)}
                    />
                    <span>{a.label}</span>
                  </label>
                ))}
              </div>
            </div>

            {error && <div role="alert" className="error-text">{error}</div>}

            <div className="auth-actions">
              <button type="button" className="btn btn-secondary" onClick={() => setStep(0)}>
                Back
              </button>
              <button type="submit" className="btn btn-primary" disabled={submitting}>
                {submitting ? 'Creating firm…' : 'Create firm'}
              </button>
            </div>
          </form>
        )}

        <div className="auth-footer">
          Already have an account? <Link to="/login">Sign in</Link>
        </div>
      </div>
    </div>
  )
}
