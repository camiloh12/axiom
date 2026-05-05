import { useEffect, useState } from 'react'

export interface AppConfig {
  clientHubEnabled: boolean
}

const DEFAULT_CONFIG: AppConfig = {
  clientHubEnabled: false,
}

let cached: AppConfig | null = null
let inflight: Promise<AppConfig> | null = null

async function fetchConfig(): Promise<AppConfig> {
  const res = await fetch('/api/v1/config')
  if (!res.ok) {
    return DEFAULT_CONFIG
  }
  const body = await res.json() as Partial<AppConfig>
  return { ...DEFAULT_CONFIG, ...body }
}

export function useAppConfig(): AppConfig {
  const [config, setConfig] = useState<AppConfig>(cached ?? DEFAULT_CONFIG)

  useEffect(() => {
    if (cached) return
    if (!inflight) {
      inflight = fetchConfig().then((c) => {
        cached = c
        return c
      })
    }
    void inflight.then(setConfig)
  }, [])

  return config
}

export function useClientHubEnabled(): boolean {
  return useAppConfig().clientHubEnabled
}

// Test hook: reset cached config and any in-flight fetch.
export function __resetAppConfigForTests() {
  cached = null
  inflight = null
}
