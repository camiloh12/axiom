import { beforeEach, describe, expect, it, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter, Route, Routes } from 'react-router-dom'
import RegisterPage from './register'
import { useAuth } from '../hooks/use-auth'

function renderRegister() {
  return render(
    <MemoryRouter initialEntries={['/register']}>
      <Routes>
        <Route path="/register" element={<RegisterPage />} />
        <Route path="/dashboard" element={<div>Dashboard</div>} />
      </Routes>
    </MemoryRouter>,
  )
}

async function fillStep1() {
  await userEvent.type(screen.getByLabelText('Firm name'), 'Acme CPAs')
  await userEvent.type(screen.getByLabelText('Your name'), 'Alice Admin')
  await userEvent.type(screen.getByLabelText('Work email'), 'alice@acme.com')
  await userEvent.type(screen.getByLabelText('Password'), 'password1')
  await userEvent.click(screen.getByRole('button', { name: /continue/i }))
}

describe('RegisterPage', () => {
  beforeEach(() => {
    useAuth.setState({ user: null, firm: null, isAuthenticated: false })
  })

  it('submits the full payload (incl. country, staff, audit types) on step 2', async () => {
    const register = vi.fn().mockResolvedValue(undefined)
    useAuth.setState({ register })

    renderRegister()
    await fillStep1()

    await userEvent.click(screen.getByLabelText(/SOC 2/i))
    await userEvent.click(screen.getByRole('button', { name: /create firm/i }))

    expect(register).toHaveBeenCalledWith(expect.objectContaining({
      firm_name: 'Acme CPAs',
      admin_email: 'alice@acme.com',
      admin_name: 'Alice Admin',
      password: 'password1',
      country: 'US',
      staff_count_range: '1-10',
      primary_audit_types: ['SOC2'],
    }))
    expect(await screen.findByText('Dashboard')).toBeInTheDocument()
  })

  it('blocks advancing when required fields or password length are missing', async () => {
    useAuth.setState({ register: vi.fn() })

    renderRegister()
    await userEvent.click(screen.getByRole('button', { name: /continue/i }))

    expect(await screen.findByRole('alert')).toHaveTextContent(/at least 8/i)
  })

  it('surfaces a server error when register rejects', async () => {
    useAuth.setState({ register: vi.fn().mockRejectedValue(new Error('dup')) })

    renderRegister()
    await fillStep1()
    await userEvent.click(screen.getByRole('button', { name: /create firm/i }))

    expect(await screen.findByRole('alert')).toHaveTextContent(/already.*in use/i)
  })
})
