import { useEffect } from 'react'
import { Link } from 'react-router-dom'
import { Layout } from '../components/layout'
import { useAuth } from '../hooks/use-auth'
import './dashboard.css'

interface ChecklistItem {
  id: string
  title: string
  description: string
  to?: string
  comingSoon?: boolean
}

const CHECKLIST: ChecklistItem[] = [
  {
    id: 'firm',
    title: 'Complete firm profile',
    description: 'Timezone, branding, and billing contact.',
    to: '/settings',
  },
  {
    id: 'team',
    title: 'Invite your team',
    description: 'Add partners, managers, and staff.',
    to: '/users',
  },
  {
    id: 'client',
    title: 'Add your first client',
    description: 'Set up a client record so you can start engagements.',
    to: '/clients',
  },
  {
    id: 'engagement',
    title: 'Create your first engagement',
    description: 'Engagement scoping ships next phase.',
    comingSoon: true,
  },
]

export default function DashboardPage() {
  const loadProfile = useAuth((s) => s.loadProfile)
  const firm = useAuth((s) => s.firm)

  useEffect(() => {
    loadProfile().catch(() => { /* swallow — topbar/sidebar just hide firm name */ })
  }, [loadProfile])

  return (
    <Layout>
      <div className="page-header">
        <h1>Welcome to Axiom{firm ? `, ${firm.name}` : ''}</h1>
        <p>Set up your firm in four steps. Work through them in any order.</p>
      </div>

      <ol className="checklist" aria-label="Onboarding checklist">
        {CHECKLIST.map((item, i) => (
          <li key={item.id} className={`checklist-item${item.comingSoon ? ' is-disabled' : ''}`}>
            <span className="checklist-index" aria-hidden="true">{i + 1}</span>
            <div className="checklist-body">
              <div className="checklist-title-row">
                <h3>{item.title}</h3>
                {item.comingSoon && <span className="badge">Coming soon</span>}
              </div>
              <p>{item.description}</p>
            </div>
            {item.to && !item.comingSoon && (
              <Link to={item.to} className="btn btn-secondary">Start</Link>
            )}
          </li>
        ))}
      </ol>
    </Layout>
  )
}
