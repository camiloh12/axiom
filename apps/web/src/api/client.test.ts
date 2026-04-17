import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

async function importFresh() {
  vi.resetModules()
  return await import('./client')
}

function mockFetch(responses: Array<{ status: number; body?: unknown }>) {
  const calls: Array<{ url: string; init: RequestInit }> = []
  let i = 0
  const fetchMock = vi.fn(async (url: string, init: RequestInit = {}) => {
    calls.push({ url, init })
    const r = responses[i++] ?? { status: 500, body: 'unexpected call' }
    return {
      ok: r.status >= 200 && r.status < 300,
      status: r.status,
      json: async () => r.body,
      text: async () => typeof r.body === 'string' ? r.body : JSON.stringify(r.body),
    } as Response
  })
  vi.stubGlobal('fetch', fetchMock)
  return { calls }
}

describe('api client', () => {
  beforeEach(() => {
    localStorage.clear()
  })
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('setTokens and clearTokens round-trip via localStorage', async () => {
    const mod = await importFresh()
    mod.setTokens('a', 'r')
    expect(localStorage.getItem('access_token')).toBe('a')
    expect(localStorage.getItem('refresh_token')).toBe('r')
    expect(mod.hasTokens()).toBe(true)

    mod.clearTokens()
    expect(localStorage.getItem('access_token')).toBeNull()
    expect(mod.hasTokens()).toBe(false)
  })

  it('injects Bearer token when one is set', async () => {
    localStorage.setItem('access_token', 'my-token')
    const mod = await importFresh()
    const { calls } = mockFetch([{ status: 200, body: { ok: true } }])

    await mod.api('/x')

    expect(calls[0].init.headers).toMatchObject({ Authorization: 'Bearer my-token' })
  })

  it('refreshes token and retries on 401 when refresh token is present', async () => {
    localStorage.setItem('access_token', 'expired')
    localStorage.setItem('refresh_token', 'refresh-me')
    const mod = await importFresh()
    const { calls } = mockFetch([
      { status: 401, body: 'expired' },
      { status: 200, body: { access_token: 'new-access', refresh_token: 'new-refresh' } },
      { status: 200, body: { ok: true } },
    ])

    const result = await mod.api<{ ok: boolean }>('/x')

    expect(result.ok).toBe(true)
    expect(calls).toHaveLength(3)
    expect(calls[1].url).toContain('/auth/refresh')
    expect(localStorage.getItem('access_token')).toBe('new-access')
    expect(calls[2].init.headers).toMatchObject({ Authorization: 'Bearer new-access' })
  })

  it('clears tokens and throws when refresh fails', async () => {
    localStorage.setItem('access_token', 'expired')
    localStorage.setItem('refresh_token', 'bad-refresh')
    const mod = await importFresh()
    mockFetch([
      { status: 401, body: 'expired' },
      { status: 401, body: 'refresh rejected' },
    ])

    await expect(mod.api('/x')).rejects.toThrow(/401/)
    expect(localStorage.getItem('access_token')).toBeNull()
  })

  it('throws ApiError on non-2xx responses', async () => {
    const mod = await importFresh()
    mockFetch([{ status: 404, body: 'not found' }])

    await expect(mod.api('/x')).rejects.toBeInstanceOf(mod.ApiError)
  })
})
