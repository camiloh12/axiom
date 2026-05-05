import { NavLink, useNavigate } from 'react-router-dom'
import { useAuth } from '../hooks/use-auth'
import './layout.css'

const nav = [
  { to: '/dashboard', label: 'Dashboard' },
  { to: '/clients', label: 'Clients' },
  { to: '/users', label: 'Users' },
  { to: '/settings', label: 'Settings' },
]

export function Layout({ children }: { children: React.ReactNode }) {
  const firm = useAuth((s) => s.firm)
  const user = useAuth((s) => s.user)
  const logout = useAuth((s) => s.logout)
  const navigate = useNavigate()

  function onLogout() {
    logout()
    navigate('/login', { replace: true })
  }

  return (
    <div className="layout">
      <aside className="sidebar" aria-label="Primary navigation">
        <div className="sidebar-brand">
          <span className="brand-mark">Axiom</span>
          {firm && <span className="brand-firm">{firm.name}</span>}
        </div>
        <nav>
          {nav.map((n) => (
            <NavLink
              key={n.to}
              to={n.to}
              className={({ isActive }: { isActive: boolean }) => `sidebar-link${isActive ? ' is-active' : ''}`}
            >
              {n.label}
            </NavLink>
          ))}
        </nav>
      </aside>
      <div className="main">
        <header className="topbar">
          <div className="topbar-title" />
          <div className="topbar-user">
            {user && <span className="user-name">{user.display_name}</span>}
            <button type="button" className="btn btn-ghost" onClick={onLogout}>
              Log out
            </button>
          </div>
        </header>
        <div className="content">{children}</div>
      </div>
    </div>
  )
}
