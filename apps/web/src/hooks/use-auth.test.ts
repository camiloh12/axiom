import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

async function freshAuth() {
  vi.resetModules()
  return await import('./use-auth')
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

describe('useAuth', () => {
  beforeEach(() => {
    localStorage.clear()
  })
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('login stores tokens and flips isAuthenticated', async () => {
    const { useAuth } = await freshAuth()
    mockFetch([{ status: 200, body: { access_token: 'a', refresh_token: 'r' } }])

    await useAuth.getState().login('x@y.com', 'pw')

    expect(useAuth.getState().isAuthenticated).toBe(true)
    expect(localStorage.getItem('access_token')).toBe('a')
  })

  it('register stores user and firm from response', async () => {
    const { useAuth } = await freshAuth()
    mockFetch([{
      status: 201,
      body: {
        access_token: 'a', refresh_token: 'r',
        user: { id: 'u', email: 'e', display_name: 'n', role: 'FirmAdmin' },
        firm: { id: 'f', name: 'F', slug: 's' },
      },
    }])

    await useAuth.getState().register({
      firm_name: 'F', admin_email: 'e', admin_name: 'n', password: 'pw',
      country: 'US', staff_count_range: '1-10', primary_audit_types: ['SOC2'],
    })

    const s = useAuth.getState()
    expect(s.user?.role).toBe('FirmAdmin')
    expect(s.firm?.name).toBe('F')
    expect(s.isAuthenticated).toBe(true)
  })

  it('logout clears tokens and resets state', async () => {
    const { useAuth } = await freshAuth()
    mockFetch([{ status: 200, body: { access_token: 'a', refresh_token: 'r' } }])
    await useAuth.getState().login('x', 'y')

    useAuth.getState().logout()

    const s = useAuth.getState()
    expect(s.isAuthenticated).toBe(false)
    expect(s.user).toBeNull()
    expect(localStorage.getItem('access_token')).toBeNull()
  })

  it('failed login leaves state untouched and rethrows', async () => {
    const { useAuth } = await freshAuth()
    mockFetch([{ status: 401, body: 'bad creds' }])

    await expect(useAuth.getState().login('x', 'wrong')).rejects.toThrow()
    expect(useAuth.getState().isAuthenticated).toBe(false)
  })
})
