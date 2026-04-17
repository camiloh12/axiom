import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter } from 'react-router-dom'
import FirmSettingsPage from './firm-settings'
import { useAuth } from '../hooks/use-auth'
import * as client from '../api/client'

function renderPage() {
  return render(<MemoryRouter><FirmSettingsPage /></MemoryRouter>)
}

describe('FirmSettingsPage', () => {
  beforeEach(() => {
    useAuth.setState({
      user: { id: 'u', email: 'a', display_name: 'A', role: 'FirmAdmin' },
      firm: { id: 'f', name: 'Firm', slug: 's' },
      isAuthenticated: true,
    })
  })
  afterEach(() => { vi.restoreAllMocks() })

  it('pre-populates the form from GET /firms/current', async () => {
    vi.spyOn(client, 'api').mockResolvedValue({
      id: 'f', name: 'Acme CPAs', slug: 'a', timezone: 'America/New_York', billing_contact_email: 'b@a.com',
    })

    renderPage()
    await waitFor(() => expect((screen.getByLabelText('Firm name') as HTMLInputElement).value).toBe('Acme CPAs'))
    expect((screen.getByLabelText('Timezone') as HTMLInputElement).value).toBe('America/New_York')
    expect((screen.getByLabelText('Billing contact email') as HTMLInputElement).value).toBe('b@a.com')
  })

  it('PATCHes /firms/current on save and shows a success notice', async () => {
    const apiMock = vi.spyOn(client, 'api').mockImplementation(async (path: string, init?: RequestInit) => {
      if (path === '/firms/current' && !init?.method) {
        return { id: 'f', name: 'Old', slug: 's', timezone: 'UTC', billing_contact_email: 'old@x.com' }
      }
      if (path === '/firms/current' && init?.method === 'PATCH') {
        return { id: 'f', name: 'New', slug: 's', timezone: 'UTC', billing_contact_email: 'old@x.com' }
      }
      throw new Error(`unexpected ${path}`)
    })

    renderPage()
    await waitFor(() => expect((screen.getByLabelText('Firm name') as HTMLInputElement).value).toBe('Old'))

    await userEvent.clear(screen.getByLabelText('Firm name'))
    await userEvent.type(screen.getByLabelText('Firm name'), 'New')
    await userEvent.click(screen.getByRole('button', { name: /save changes/i }))

    await waitFor(() => expect(apiMock).toHaveBeenCalledWith('/firms/current', expect.objectContaining({ method: 'PATCH' })))
    expect(await screen.findByRole('status')).toHaveTextContent(/saved/i)
  })
})
