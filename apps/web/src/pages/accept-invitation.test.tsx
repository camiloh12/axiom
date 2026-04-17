import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter, Route, Routes } from 'react-router-dom'
import AcceptInvitationPage from './accept-invitation'

function renderAt(search: string) {
  return render(
    <MemoryRouter initialEntries={[`/accept-invitation${search}`]}>
      <Routes>
        <Route path="/accept-invitation" element={<AcceptInvitationPage />} />
        <Route path="/dashboard" element={<div>Dashboard</div>} />
      </Routes>
    </MemoryRouter>,
  )
}

function mockFetch(responses: Array<{ status: number; body?: unknown }>) {
  let i = 0
  const fetchMock = vi.fn(async () => {
    const r = responses[i++] ?? { status: 500, body: 'unexpected call' }
    return {
      ok: r.status >= 200 && r.status < 300,
      status: r.status,
      json: async () => r.body,
      text: async () => typeof r.body === 'string' ? r.body : JSON.stringify(r.body),
    } as Response
  })
  vi.stubGlobal('fetch', fetchMock)
}

describe('AcceptInvitationPage', () => {
  beforeEach(() => { localStorage.clear() })
  afterEach(() => { vi.unstubAllGlobals() })

  it('validates the token and displays email + role', async () => {
    mockFetch([{ status: 200, body: { email: 'new@firm.com', assigned_role: 'Staff', expires_at: '' } }])

    renderAt('?token=abc')

    expect(await screen.findByText(/new@firm.com/)).toBeInTheDocument()
    expect(screen.getByText(/Staff/)).toBeInTheDocument()
  })

  it('accepts the invitation, stores tokens, and navigates to dashboard', async () => {
    mockFetch([
      { status: 200, body: { email: 'new@firm.com', assigned_role: 'Staff', expires_at: '' } },
      { status: 201, body: {
        access_token: 'a', refresh_token: 'r',
        user: { id: 'u', email: 'new@firm.com', display_name: 'Staffer', role: 'Staff' },
      }},
    ])

    renderAt('?token=abc')
    await screen.findByText(/new@firm.com/)

    await userEvent.type(screen.getByLabelText('Your name'), 'Staffer')
    await userEvent.type(screen.getByLabelText('Create a password'), 'pwpwpwpw')
    await userEvent.click(screen.getByRole('button', { name: /accept and sign in/i }))

    await waitFor(() => expect(screen.getByText('Dashboard')).toBeInTheDocument())
    expect(localStorage.getItem('access_token')).toBe('a')
  })

  it('shows an error state for an invalid token', async () => {
    mockFetch([{ status: 404, body: 'not found' }])

    renderAt('?token=bad')
    expect(await screen.findByText(/invalid or has expired/i)).toBeInTheDocument()
  })
})
