import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter } from 'react-router-dom'
import UsersPage from './users'
import { useAuth } from '../hooks/use-auth'
import * as client from '../api/client'

function renderPage() {
  return render(<MemoryRouter><UsersPage /></MemoryRouter>)
}

describe('UsersPage', () => {
  beforeEach(() => {
    useAuth.setState({
      user: { id: 'u', email: 'a@b.com', display_name: 'Admin', role: 'FirmAdmin' },
      firm: { id: 'f', name: 'Firm', slug: 's' },
      isAuthenticated: true,
    })
  })
  afterEach(() => { vi.restoreAllMocks() })

  it('renders fetched users and invitations', async () => {
    vi.spyOn(client, 'api').mockImplementation(async (path: string) => {
      if (path === '/users') return { items: [{ id: 'u1', email: 'a@b.com', display_name: 'Alice', role: 'Partner' }] }
      if (path === '/invitations') return { items: [] }
      throw new Error(`unexpected ${path}`)
    })

    renderPage()
    const row = (await screen.findByText('Alice')).closest('tr')!
    expect(row).toHaveTextContent('Partner')
  })

  it('calls POST /invitations when inviting, with email and role', async () => {
    const apiMock = vi.spyOn(client, 'api').mockImplementation(async (path: string, init?: RequestInit) => {
      if (path === '/users') return { items: [] }
      if (path === '/invitations' && init?.method === 'POST') {
        return { id: 'i1', email: 'new@test.com', assigned_role: 'Staff', status: 'Sent', expires_at: '', token: 'tok' }
      }
      if (path === '/invitations') return { items: [] }
      throw new Error(`unexpected ${path}`)
    })

    renderPage()
    await waitFor(() => expect(screen.getByRole('button', { name: /send invitation/i })).toBeEnabled())

    await userEvent.type(screen.getByLabelText('Email'), 'new@test.com')
    await userEvent.selectOptions(screen.getByLabelText('Role'), 'Staff')
    await userEvent.click(screen.getByRole('button', { name: /send invitation/i }))

    expect(apiMock).toHaveBeenCalledWith('/invitations', expect.objectContaining({
      method: 'POST',
      body: JSON.stringify({ email: 'new@test.com', assigned_role: 'Staff' }),
    }))
  })
})
