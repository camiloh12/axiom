import { beforeEach, describe, expect, it, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import DashboardPage from './dashboard'
import { useAuth } from '../hooks/use-auth'

function renderDash() {
  return render(
    <MemoryRouter>
      <DashboardPage />
    </MemoryRouter>,
  )
}

describe('DashboardPage', () => {
  beforeEach(() => {
    useAuth.setState({
      user: { id: 'u1', email: 'e', display_name: 'Alice', role: 'FirmAdmin' },
      firm: { id: 'f1', name: 'Acme', slug: 's' },
      isAuthenticated: true,
    })
  })

  it('renders the onboarding checklist with four items', () => {
    useAuth.setState({ loadProfile: vi.fn().mockResolvedValue(undefined) })
    renderDash()

    const items = screen.getAllByRole('listitem')
    expect(items).toHaveLength(4)
    expect(screen.getByRole('heading', { name: /complete firm profile/i })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: /invite your team/i })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: /add your first client/i })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: /create your first engagement/i })).toBeInTheDocument()
  })

  it('marks "create first engagement" as coming soon and disables it', () => {
    useAuth.setState({ loadProfile: vi.fn().mockResolvedValue(undefined) })
    renderDash()

    const heading = screen.getByRole('heading', { name: /create your first engagement/i })
    const item = heading.closest('li')!
    expect(item.className).toMatch(/is-disabled/)
    expect(screen.getByText(/coming soon/i)).toBeInTheDocument()
  })

  it('calls loadProfile on mount', () => {
    const loadProfile = vi.fn().mockResolvedValue(undefined)
    useAuth.setState({ loadProfile })
    renderDash()
    expect(loadProfile).toHaveBeenCalled()
  })
})
