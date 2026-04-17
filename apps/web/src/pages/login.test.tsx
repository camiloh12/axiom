import { beforeEach, describe, expect, it, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter, Route, Routes } from 'react-router-dom'
import LoginPage from './login'
import { useAuth } from '../hooks/use-auth'

function renderLogin() {
  return render(
    <MemoryRouter initialEntries={['/login']}>
      <Routes>
        <Route path="/login" element={<LoginPage />} />
        <Route path="/dashboard" element={<div>Dashboard</div>} />
        <Route path="/register" element={<div>Register</div>} />
      </Routes>
    </MemoryRouter>,
  )
}

describe('LoginPage', () => {
  beforeEach(() => {
    useAuth.setState({ user: null, firm: null, isAuthenticated: false })
  })

  it('submits entered credentials and navigates to dashboard on success', async () => {
    const login = vi.fn().mockResolvedValue(undefined)
    useAuth.setState({ login })

    renderLogin()
    await userEvent.type(screen.getByLabelText('Email'), 'a@b.com')
    await userEvent.type(screen.getByLabelText('Password'), 'password1')
    await userEvent.click(screen.getByRole('button', { name: /sign in/i }))

    expect(login).toHaveBeenCalledWith('a@b.com', 'password1')
    expect(await screen.findByText('Dashboard')).toBeInTheDocument()
  })

  it('renders an error when login rejects', async () => {
    useAuth.setState({ login: vi.fn().mockRejectedValue(new Error('bad')) })

    renderLogin()
    await userEvent.type(screen.getByLabelText('Email'), 'a@b.com')
    await userEvent.type(screen.getByLabelText('Password'), 'wrong')
    await userEvent.click(screen.getByRole('button', { name: /sign in/i }))

    expect(await screen.findByRole('alert')).toHaveTextContent(/invalid/i)
  })

  it('links to the register page', () => {
    renderLogin()
    const link = screen.getByRole('link', { name: /start a free trial/i })
    expect(link).toHaveAttribute('href', '/register')
  })
})
