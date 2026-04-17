import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter } from 'react-router-dom'
import ClientsPage from './clients'
import { useAuth } from '../hooks/use-auth'
import * as client from '../api/client'

function renderPage() {
  return render(<MemoryRouter><ClientsPage /></MemoryRouter>)
}

describe('ClientsPage', () => {
  beforeEach(() => {
    useAuth.setState({
      user: { id: 'u', email: 'a', display_name: 'A', role: 'FirmAdmin' },
      firm: { id: 'f', name: 'F', slug: 's' },
      isAuthenticated: true,
    })
  })
  afterEach(() => { vi.restoreAllMocks() })

  it('renders fetched clients', async () => {
    vi.spyOn(client, 'api').mockResolvedValue({
      items: [{ id: 'c1', name: 'TechCorp', industry: 'Software', primary_contact_email: 'ops@t.com' }],
    })

    renderPage()
    expect(await screen.findByText('TechCorp')).toBeInTheDocument()
    expect(screen.getByText('Software')).toBeInTheDocument()
  })

  it('POSTs a new client with form payload', async () => {
    const apiMock = vi.spyOn(client, 'api').mockImplementation(async (path: string, init?: RequestInit) => {
      if (path === '/clients' && !init?.method) return { items: [] }
      if (path === '/clients' && init?.method === 'POST') {
        return { id: 'c2', name: 'NewCo', industry: 'Retail', primary_contact_email: '' }
      }
      throw new Error(`unexpected ${path}`)
    })

    renderPage()
    await waitFor(() => expect(screen.getByRole('button', { name: /add client/i })).toBeEnabled())

    await userEvent.type(screen.getByLabelText('Client name'), 'NewCo')
    await userEvent.type(screen.getByLabelText('Industry'), 'Retail')
    await userEvent.click(screen.getByRole('button', { name: /add client/i }))

    expect(apiMock).toHaveBeenCalledWith('/clients', expect.objectContaining({
      method: 'POST',
      body: JSON.stringify({ name: 'NewCo', industry: 'Retail', primary_contact_email: '' }),
    }))
  })
})
