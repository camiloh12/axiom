import { QueryClient, QueryClientProvider, useQuery } from '@tanstack/react-query'

const queryClient = new QueryClient()

function HealthCheck() {
  const { data, isLoading, isError } = useQuery({
    queryKey: ['health'],
    queryFn: () => fetch('/api/healthz').then(res => res.json()),
  })

  if (isLoading) return <p>Connecting to API...</p>
  if (isError) return <p style={{ color: 'red' }}>API disconnected</p>

  return <p style={{ color: 'green' }}>API connected: {data?.status}</p>
}

export default function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <div style={{ padding: '2rem' }}>
        <h1>Axiom</h1>
        <HealthCheck />
      </div>
    </QueryClientProvider>
  )
}
