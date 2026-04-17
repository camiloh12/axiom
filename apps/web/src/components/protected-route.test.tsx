import { describe, expect, it } from 'vitest'
import { render, screen } from '@testing-library/react'
import { MemoryRouter, Route, Routes } from 'react-router-dom'
import { ProtectedRoute } from './protected-route'
import { useAuth } from '../hooks/use-auth'

function renderAt(path: string) {
  return render(
    <MemoryRouter initialEntries={[path]}>
      <Routes>
        <Route path="/login" element={<div>Login</div>} />
        <Route
          path="/secret"
          element={
            <ProtectedRoute>
              <div>Secret</div>
            </ProtectedRoute>
          }
        />
      </Routes>
    </MemoryRouter>,
  )
}

describe('ProtectedRoute', () => {
  it('renders children when authenticated', () => {
    useAuth.setState({ isAuthenticated: true })
    renderAt('/secret')
    expect(screen.getByText('Secret')).toBeInTheDocument()
  })

  it('redirects to /login when not authenticated', () => {
    useAuth.setState({ isAuthenticated: false })
    renderAt('/secret')
    expect(screen.getByText('Login')).toBeInTheDocument()
  })
})
