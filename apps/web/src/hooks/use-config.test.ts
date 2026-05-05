import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'
import { renderHook, waitFor } from '@testing-library/react'
import { useAppConfig, useClientHubEnabled, __resetAppConfigForTests } from './use-config'

describe('useAppConfig', () => {
  beforeEach(() => {
    __resetAppConfigForTests()
    vi.spyOn(globalThis, 'fetch').mockReset()
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('returns default config until fetch resolves, then merges the response', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValueOnce(
      new Response(JSON.stringify({ clientHubEnabled: true }), {
        status: 200,
        headers: { 'content-type': 'application/json' },
      }),
    )

    const { result } = renderHook(() => useAppConfig())
    expect(result.current.clientHubEnabled).toBe(false)
    await waitFor(() => expect(result.current.clientHubEnabled).toBe(true))
  })

  it('falls back to default config (clientHubEnabled=false) when the request fails', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValueOnce(
      new Response('boom', { status: 500 }),
    )

    const { result } = renderHook(() => useAppConfig())
    await waitFor(() => expect(result.current).toEqual({ clientHubEnabled: false }))
  })

  it('useClientHubEnabled() proxies the underlying flag', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValueOnce(
      new Response(JSON.stringify({ clientHubEnabled: true }), {
        status: 200,
        headers: { 'content-type': 'application/json' },
      }),
    )

    const { result } = renderHook(() => useClientHubEnabled())
    await waitFor(() => expect(result.current).toBe(true))
  })
})
