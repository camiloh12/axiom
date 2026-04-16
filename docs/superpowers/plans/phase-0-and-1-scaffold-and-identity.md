# Phase 0 + Phase 1: Scaffold & Identity Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Go from nothing installed to a working identity system — firm registration, login, JWT auth, RLS, staff invitations, and a React frontend with login, dashboard, and user management.

**Architecture:** Go modular monolith with Chi router, PostgreSQL with RLS for multi-tenancy, RSA-signed JWTs for auth. React SPA with Vite, TanStack Query, and openapi-typescript generated types. Docker Compose for local Postgres. All local, no AWS.

**Tech Stack:** Go 1.24+, Chi, pgx/v5, sqlc, golang-migrate, bcrypt, golang-jwt/v5 | TypeScript, React 19, Vite, TanStack Query, React Router, Zustand | PostgreSQL 17 + pgvector, Docker Compose

---

## Development Methodology: Test-Driven Development

**This plan is executed test-first.** The canonical policy lives in `docs/superpowers/specs/implementation-plan-design.md` under "Cross-Cutting Methodology: Test-Driven Development" — read that section before starting. This section codifies how it applies to every task below.

### The loop for every task

For any file containing behavior (service methods, handlers, middleware, hooks, components with logic, SQL queries worth asserting against):

1. **Red** — write the test file and assertions first. Run the test suite and confirm it fails. A compile error is *not* a valid red — the test must reach an assertion that evaluates to false before you write the implementation.
2. **Green** — write the smallest implementation that makes the failing test pass. Do not stub out features that aren't under test.
3. **Refactor** — rename, extract, deduplicate. Re-run the test between each change.

Every task below that uses this pattern already shows it as three-step choreography: "Write test → Run — expect failure → Implement → Run — expect pass." If a task doesn't show these steps explicitly (e.g., a pure configuration task), it either has no behavior to test or its coverage comes from the next task that consumes it.

### What each layer looks like in this plan

| File pattern | Test type | Where in this plan |
|---|---|---|
| `internal/*/service.go` (DB-touching) | Integration tests with `platform.TestDB(t)` | Tasks 10, 13 |
| `internal/*/handler.go` | `httptest.NewRecorder` tests | Tasks 11, 13 |
| `internal/identity/jwt.go` | Pure unit tests (no DB) | Task 8 |
| `internal/gateway/middleware.go` | Unit tests with mocked JWT issuer | Task 9 |
| `internal/platform/*.go` | Pure unit tests where there's branching | Task 5 |
| RLS policies | Multi-tenant isolation tests | Task 17 |
| React hooks (`use-auth.ts`) | Vitest + RTL | Task 14 |
| React forms (`login.tsx`, `register.tsx`) | RTL interaction tests | Tasks 15, 16 |

### Rules for this phase

- **Do not move to the next task until the current task's tests pass.** A failing test is a blocking issue, not a TODO.
- **Do not commit an implementation without its tests in the same commit** (or an earlier commit in the same PR that introduced the failing test). "Tests to follow" is not acceptable.
- **Commit messages for test-first work** should either call out the tests ("feat: add JWT issuer with RSA signing, verification, and refresh — incl. unit tests") or be split into two commits ("test: add failing JWT issuer tests" → "feat: implement JWT issuer").
- **Trivial code does not need a test.** Config struct parsing, one-line constructors with no branching, and generated sqlc code are covered transitively by the tests on the code that uses them. When in doubt: if a plausible bug would be caught by a test one layer up, skip the dedicated test.
- **When adjusting tests to match reality** (e.g., sqlc generated types differ from what the test assumed), verify the adjustment preserves the original intent of the assertion. Do not weaken a test to make it pass.

The Phase 0 CI pipeline (Task 18) runs all tests on every PR. A failing test blocks merge.

---

## Git Workflow

Before starting any task, create the phase branch and push all work there:

- [ ] **Create phase branch**

```bash
git checkout master
git pull origin master
git checkout -b phase-0-1-scaffold-and-identity
```

All commits in Tasks 1–17 go to this branch. Push the branch after each commit:

```bash
git push -u origin phase-0-1-scaffold-and-identity
```

After all tasks are complete, the user creates a PR from `phase-0-1-scaffold-and-identity` → `master`, reviews, and merges.

**Before starting the next phase (Phase 2):** return to master, pull the merged changes, verify, then create a new branch for that phase.

---

## File Structure

```
axiom/                              (repo root — already exists)
  docker-compose.yml                 — Postgres 17 + pgvector, Mailhog
  turbo.json                         — Turborepo pipeline config
  apps/
    axiom-api/
      go.mod
      go.sum
      sqlc.yaml                      — sqlc configuration
      cmd/server/
        main.go                      — entrypoint: config, DB, wire modules, start server
      internal/
        platform/
          config.go                  — Config struct, envconfig loading
          database.go                — pgx pool creation, RLS helper, health check
          errors.go                  — AppError type, standard error constructors
          response.go                — JSON response helpers (writeJSON, writeError)
          testdb.go                  — test helper: create test DB, run migrations, cleanup
        gateway/
          middleware.go              — AuthMiddleware, WithFirmIsolation, WithRole
          middleware_test.go         — unit tests with mock JWT
        identity/
          service.go                 — IdentityService interface and implementation
          service_test.go            — integration tests against real Postgres
          handler.go                 — HTTP handlers for identity endpoints
          handler_test.go            — HTTP handler tests using httptest
          jwt.go                     — JWT issuer (sign, verify, refresh)
          jwt_test.go                — JWT unit tests
          queries/
            firms.sql                — sqlc queries for firms table
            users.sql                — sqlc queries for users table
            clients.sql              — sqlc queries for clients table
            invitations.sql          — sqlc queries for invitations table
      migrations/
        000001_identity_enums.up.sql
        000001_identity_enums.down.sql
        000002_identity_tables.up.sql
        000002_identity_tables.down.sql
    web/
      package.json
      tsconfig.json
      vite.config.ts
      index.html
      src/
        main.tsx                     — React root with providers
        App.tsx                      — Router setup
        api/
          client.ts                  — fetch wrapper with JWT injection
        hooks/
          use-auth.ts                — auth context: login, logout, token refresh
        pages/
          login.tsx
          register.tsx
          accept-invitation.tsx
          dashboard.tsx
          firm-settings.tsx
          users.tsx
          clients.tsx
        components/
          layout.tsx                 — sidebar + topbar shell
          protected-route.tsx        — redirects to /login if no token
  packages/
    openapi/                         — already exists with specs
```

---

### Task 1: Docker Compose and tool verification

**Files:**
- Create: `docker-compose.yml`

- [ ] **Step 1: Verify tools are installed**

Run each command and verify output:
```bash
go version        # expect: go1.24.x or later
node --version    # expect: v22.x
docker --version  # expect: Docker 27.x or later
npm --version     # expect: 10.x or later
```

If any are missing, install them:
- Go: https://go.dev/dl/
- Node.js LTS: https://nodejs.org/
- Docker Desktop: https://www.docker.com/products/docker-desktop/

- [ ] **Step 2: Create docker-compose.yml**

```yaml
services:
  postgres:
    image: pgvector/pgvector:pg17
    ports:
      - "5432:5432"
    environment:
      POSTGRES_DB: axiom_db
      POSTGRES_USER: axiom_svc
      POSTGRES_PASSWORD: localdev
    volumes:
      - pgdata:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U axiom_svc -d axiom_db"]
      interval: 5s
      timeout: 3s
      retries: 5

  mailhog:
    image: mailhog/mailhog:v1.0.1
    ports:
      - "1025:1025"   # SMTP
      - "8025:8025"   # Web UI

volumes:
  pgdata:
```

- [ ] **Step 3: Start services and verify**

```bash
docker compose up -d
```

Run: `docker compose ps`
Expected: both `postgres` and `mailhog` show as healthy/running.

Verify Postgres connectivity:
```bash
docker compose exec postgres psql -U axiom_svc -d axiom_db -c "SELECT version();"
```
Expected: PostgreSQL 17.x output.

Verify pgvector:
```bash
docker compose exec postgres psql -U axiom_svc -d axiom_db -c "CREATE EXTENSION IF NOT EXISTS vector; SELECT extversion FROM pg_extension WHERE extname = 'vector';"
```
Expected: version number (e.g., 0.8.0).

Verify Mailhog: open `http://localhost:8025` in browser. Expected: empty Mailhog inbox UI.

- [ ] **Step 4: Commit**

```bash
git add docker-compose.yml
git commit -m "chore: add docker-compose with Postgres 17 + pgvector and Mailhog"
```

---

### Task 2: Go project scaffold

**Files:**
- Create: `apps/axiom-api/go.mod`
- Create: `apps/axiom-api/cmd/server/main.go`

- [ ] **Step 1: Initialize Go module**

```bash
cd apps/axiom-api
go mod init github.com/axiom-platform/axiom/apps/axiom-api
```

- [ ] **Step 2: Create minimal main.go with Chi healthz endpoint**

Create `apps/axiom-api/cmd/server/main.go`:

```go
package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	slog.Info("starting server", "port", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}
```

- [ ] **Step 3: Install dependency and verify build**

```bash
cd apps/axiom-api
go mod tidy
go build ./cmd/server
```

Expected: builds with no errors; produces `server` (or `server.exe`) binary.

- [ ] **Step 4: Run and verify healthz**

```bash
go run ./cmd/server &
curl http://localhost:8080/healthz
```

Expected: `{"status":"ok"}`

Stop the server (Ctrl+C or kill the background process).

- [ ] **Step 5: Commit**

```bash
git add apps/axiom-api/
git commit -m "chore: scaffold Go API with Chi router and /healthz endpoint"
```

---

### Task 3: React project scaffold

**Files:**
- Create: `apps/web/` (entire Vite + React project)

- [ ] **Step 1: Create Vite project**

```bash
cd apps
npm create vite@latest web -- --template react-ts
cd web
npm install
```

- [ ] **Step 2: Install core dependencies**

```bash
npm install @tanstack/react-query react-router-dom zustand
npm install -D @types/react-router-dom
```

- [ ] **Step 3: Configure Vite proxy to Go API**

Replace `apps/web/vite.config.ts`:

```typescript
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
})
```

- [ ] **Step 4: Create minimal App with health check**

Replace `apps/web/src/App.tsx`:

```tsx
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
```

Replace `apps/web/src/main.tsx`:

```tsx
import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import App from './App'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <App />
  </StrictMode>,
)
```

- [ ] **Step 5: Verify frontend connects to backend**

In one terminal: `cd apps/axiom-api && go run ./cmd/server`
In another terminal: `cd apps/web && npm run dev`

Open `http://localhost:3000`. Expected: "Axiom" heading with green "API connected: ok" text.

- [ ] **Step 6: Commit**

```bash
git add apps/web/
git commit -m "chore: scaffold React app with Vite, TanStack Query, and API health check"
```

---

### Task 4: Turborepo and OpenAPI codegen pipeline

**Files:**
- Create: `turbo.json`
- Create: `package.json` (root)
- Create: `packages/openapi/package.json`

- [ ] **Step 1: Create root package.json for Turborepo**

Create `package.json` at repo root:

```json
{
  "name": "axiom",
  "private": true,
  "workspaces": [
    "apps/web",
    "packages/openapi"
  ],
  "devDependencies": {
    "turbo": "^2"
  },
  "scripts": {
    "dev": "turbo dev",
    "build": "turbo build",
    "lint": "turbo lint",
    "codegen": "turbo codegen"
  }
}
```

- [ ] **Step 2: Create turbo.json**

```json
{
  "$schema": "https://turbo.build/schema.json",
  "tasks": {
    "build": {
      "dependsOn": ["codegen"],
      "outputs": ["dist/**"]
    },
    "dev": {
      "cache": false,
      "persistent": true
    },
    "lint": {},
    "codegen": {
      "outputs": ["src/api/generated/**"]
    }
  }
}
```

- [ ] **Step 3: Set up openapi-typescript codegen**

Create `packages/openapi/package.json`:

```json
{
  "name": "@axiom/openapi",
  "private": true,
  "scripts": {
    "codegen": "openapi-typescript identity-service.yaml -o ../../apps/web/src/api/generated/identity.ts && openapi-typescript common.yaml -o ../../apps/web/src/api/generated/common.ts"
  },
  "devDependencies": {
    "openapi-typescript": "^7"
  }
}
```

- [ ] **Step 4: Install dependencies and run codegen**

```bash
cd /path/to/axiom
npm install
npm run codegen
```

Expected: TypeScript types generated in `apps/web/src/api/generated/`.

- [ ] **Step 5: Install oapi-codegen for Go**

```bash
go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
```

Verify: `oapi-codegen --version` returns a version number.

Note: Go server codegen will be wired in Task 7 when we set up the identity module. For now, the toolchain is installed and ready.

- [ ] **Step 6: Commit**

```bash
git add turbo.json package.json packages/openapi/package.json apps/web/src/api/generated/
git commit -m "chore: add Turborepo config and OpenAPI TypeScript codegen pipeline"
```

---

### Task 5: Platform package — config, database, errors, response helpers

**Files:**
- Create: `apps/axiom-api/internal/platform/config.go`
- Create: `apps/axiom-api/internal/platform/config_test.go`
- Create: `apps/axiom-api/internal/platform/database.go`
- Create: `apps/axiom-api/internal/platform/errors.go`
- Create: `apps/axiom-api/internal/platform/errors_test.go`
- Create: `apps/axiom-api/internal/platform/response.go`
- Create: `apps/axiom-api/internal/platform/response_test.go`

**TDD note:** Config loading, error constructors, and response helpers all have branching behavior worth locking in with tests. `database.go` is tested transitively by the identity service tests in Task 10 (which uses `platform.NewDBPool` via `TestDB`), so a dedicated test is not required here.

- [ ] **Step 1: Create config.go**

```go
package platform

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Port        string `envconfig:"PORT" default:"8080"`
	DatabaseURL string `envconfig:"DATABASE_URL" default:"postgres://axiom_svc:localdev@localhost:5432/axiom_db?sslmode=disable"`
	JWTPrivKey  string `envconfig:"JWT_PRIVATE_KEY"` // PEM-encoded RSA private key
	JWTPubKey   string `envconfig:"JWT_PUBLIC_KEY"`  // PEM-encoded RSA public key
	Environment string `envconfig:"ENVIRONMENT" default:"development"`
}

func LoadConfig() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
```

- [ ] **Step 2: Create database.go**

```go
package platform

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDBPool(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse database URL: %w", err)
	}
	config.MaxConns = 20

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return pool, nil
}

// SetFirmID sets the RLS session variable for the current transaction.
func SetFirmID(ctx context.Context, pool *pgxpool.Pool, firmID string) error {
	_, err := pool.Exec(ctx, "SELECT set_config('app.current_firm_id', $1, true)", firmID)
	return err
}
```

- [ ] **Step 3: Create errors.go**

```go
package platform

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Code    int    `json:"-"`
	Message string `json:"error"`
	Detail  string `json:"detail,omitempty"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

func ErrNotFound(msg string) *AppError {
	return &AppError{Code: http.StatusNotFound, Message: msg}
}

func ErrUnauthorized(msg string) *AppError {
	return &AppError{Code: http.StatusUnauthorized, Message: msg}
}

func ErrForbidden(msg string) *AppError {
	return &AppError{Code: http.StatusForbidden, Message: msg}
}

func ErrBadRequest(msg string) *AppError {
	return &AppError{Code: http.StatusBadRequest, Message: msg}
}

func ErrConflict(msg string) *AppError {
	return &AppError{Code: http.StatusConflict, Message: msg}
}

func ErrValidation(msg string, detail string) *AppError {
	return &AppError{Code: http.StatusUnprocessableEntity, Message: msg, Detail: detail}
}

func ErrInternal(msg string) *AppError {
	return &AppError{Code: http.StatusInternalServerError, Message: msg}
}
```

- [ ] **Step 4: Create response.go**

```go
package platform

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to write JSON response", "error", err)
	}
}

func WriteError(w http.ResponseWriter, err error) {
	if appErr, ok := err.(*AppError); ok {
		WriteJSON(w, appErr.Code, appErr)
		return
	}
	slog.Error("unexpected error", "error", err)
	WriteJSON(w, http.StatusInternalServerError, &AppError{
		Code:    http.StatusInternalServerError,
		Message: "internal server error",
	})
}
```

- [ ] **Step 5: Install Go dependencies (needed for tests to compile)**

```bash
cd apps/axiom-api
go get github.com/kelseyhightower/envconfig
go get github.com/jackc/pgx/v5
go get github.com/stretchr/testify
go mod tidy
```

- [ ] **Step 6: Write failing tests for config, errors, and response**

Following TDD, write the test files first, run them, confirm they fail, *then* the implementations above will already satisfy them. (Ordering note: in a strict Red→Green cycle you'd write these tests before the implementation in Steps 1–4. If you did those steps first, treat this as a validation pass — delete each implementation file temporarily, re-run the relevant test, confirm red, restore the implementation, confirm green.)

`apps/axiom-api/internal/platform/config_test.go`:

```go
package platform_test

import (
	"testing"

	"github.com/axiom-platform/axiom/apps/axiom-api/internal/platform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig_Defaults(t *testing.T) {
	t.Setenv("PORT", "")
	t.Setenv("DATABASE_URL", "")
	t.Setenv("ENVIRONMENT", "")

	cfg, err := platform.LoadConfig()
	require.NoError(t, err)
	assert.Equal(t, "8080", cfg.Port)
	assert.Contains(t, cfg.DatabaseURL, "axiom_db")
	assert.Equal(t, "development", cfg.Environment)
}

func TestLoadConfig_EnvOverrides(t *testing.T) {
	t.Setenv("PORT", "9090")
	t.Setenv("DATABASE_URL", "postgres://custom/db")
	t.Setenv("ENVIRONMENT", "production")

	cfg, err := platform.LoadConfig()
	require.NoError(t, err)
	assert.Equal(t, "9090", cfg.Port)
	assert.Equal(t, "postgres://custom/db", cfg.DatabaseURL)
	assert.Equal(t, "production", cfg.Environment)
}
```

`apps/axiom-api/internal/platform/errors_test.go`:

```go
package platform_test

import (
	"net/http"
	"testing"

	"github.com/axiom-platform/axiom/apps/axiom-api/internal/platform"
	"github.com/stretchr/testify/assert"
)

func TestAppError_Error(t *testing.T) {
	err := platform.ErrNotFound("user")
	assert.Contains(t, err.Error(), "404")
	assert.Contains(t, err.Error(), "user")
}

func TestErrorConstructors_SetCorrectStatusCodes(t *testing.T) {
	cases := []struct {
		name     string
		err      *platform.AppError
		expected int
	}{
		{"NotFound", platform.ErrNotFound("x"), http.StatusNotFound},
		{"Unauthorized", platform.ErrUnauthorized("x"), http.StatusUnauthorized},
		{"Forbidden", platform.ErrForbidden("x"), http.StatusForbidden},
		{"BadRequest", platform.ErrBadRequest("x"), http.StatusBadRequest},
		{"Conflict", platform.ErrConflict("x"), http.StatusConflict},
		{"Validation", platform.ErrValidation("x", "d"), http.StatusUnprocessableEntity},
		{"Internal", platform.ErrInternal("x"), http.StatusInternalServerError},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.expected, c.err.Code)
		})
	}
}

func TestErrValidation_CarriesDetail(t *testing.T) {
	err := platform.ErrValidation("missing fields", "email is required")
	assert.Equal(t, "missing fields", err.Message)
	assert.Equal(t, "email is required", err.Detail)
}
```

`apps/axiom-api/internal/platform/response_test.go`:

```go
package platform_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/axiom-platform/axiom/apps/axiom-api/internal/platform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteJSON_SetsStatusAndBody(t *testing.T) {
	rr := httptest.NewRecorder()
	platform.WriteJSON(rr, http.StatusCreated, map[string]string{"ok": "yes"})

	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var body map[string]string
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &body))
	assert.Equal(t, "yes", body["ok"])
}

func TestWriteError_AppErrorRendersCode(t *testing.T) {
	rr := httptest.NewRecorder()
	platform.WriteError(rr, platform.ErrNotFound("firm"))

	assert.Equal(t, http.StatusNotFound, rr.Code)
	var body map[string]string
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &body))
	assert.Equal(t, "firm", body["error"])
}

func TestWriteError_GenericErrorRenders500(t *testing.T) {
	rr := httptest.NewRecorder()
	platform.WriteError(rr, errors.New("boom"))

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}
```

- [ ] **Step 7: Run tests — expect pass**

```bash
cd apps/axiom-api
go test ./internal/platform/ -v -count=1
```

Expected: all tests PASS.

- [ ] **Step 8: Verify compilation across the module**

```bash
go build ./...
```

Expected: no errors.

- [ ] **Step 9: Commit**

```bash
git add apps/axiom-api/internal/platform/ apps/axiom-api/go.mod apps/axiom-api/go.sum
git commit -m "feat: add platform package with config, database pool, errors, and response helpers (incl. unit tests)"
```

---

### Task 6: Database migrations — identity enums and tables

**Files:**
- Create: `apps/axiom-api/migrations/000001_identity_enums.up.sql`
- Create: `apps/axiom-api/migrations/000001_identity_enums.down.sql`
- Create: `apps/axiom-api/migrations/000002_identity_tables.up.sql`
- Create: `apps/axiom-api/migrations/000002_identity_tables.down.sql`

- [ ] **Step 1: Create enum migration**

`000001_identity_enums.up.sql`:

```sql
CREATE TYPE user_role AS ENUM (
  'FirmAdmin', 'Partner', 'Manager', 'Staff', 'EQReviewer',
  'ClientAdmin', 'ClientUser', 'ViewOnly'
);

CREATE TYPE auth_method AS ENUM ('Password', 'OAuth', 'SAML');

CREATE TYPE notification_frequency AS ENUM ('RealTime', 'Daily', 'Weekly');

CREATE TYPE invitation_status AS ENUM ('Sent', 'Accepted', 'Expired');
```

`000001_identity_enums.down.sql`:

```sql
DROP TYPE IF EXISTS invitation_status;
DROP TYPE IF EXISTS notification_frequency;
DROP TYPE IF EXISTS auth_method;
DROP TYPE IF EXISTS user_role;
```

- [ ] **Step 2: Create identity tables migration**

`000002_identity_tables.up.sql`:

```sql
-- RLS helper: function to get current firm ID from session variable
CREATE OR REPLACE FUNCTION current_firm_id() RETURNS uuid AS $$
  SELECT current_setting('app.current_firm_id', true)::uuid;
$$ LANGUAGE sql STABLE;

-- Firms
CREATE TABLE firms (
  id             uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  name           text NOT NULL,
  slug           text NOT NULL UNIQUE,
  logo_url       text,
  timezone       text NOT NULL DEFAULT 'America/New_York',
  billing_contact_email text NOT NULL,
  subscription_tier text NOT NULL DEFAULT 'Growth'
    CHECK (subscription_tier IN ('Growth', 'Scale', 'Enterprise')),
  country        text NOT NULL DEFAULT 'US'
    CHECK (country IN ('US', 'CA')),
  staff_count_range text,
  primary_audit_types jsonb DEFAULT '[]',
  settings       jsonb NOT NULL DEFAULT '{}',
  created_at     timestamptz NOT NULL DEFAULT now()
);

ALTER TABLE firms ENABLE ROW LEVEL SECURITY;
CREATE POLICY firms_isolation ON firms
  USING (id = current_firm_id());

-- Users
CREATE TABLE users (
  id                     uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  firm_id                uuid REFERENCES firms(id),
  client_id              uuid,
  email                  text NOT NULL UNIQUE,
  display_name           text NOT NULL,
  role                   user_role NOT NULL,
  auth_method            auth_method NOT NULL DEFAULT 'Password',
  password_hash          text,
  notification_frequency notification_frequency NOT NULL DEFAULT 'Daily',
  tour_completed         boolean NOT NULL DEFAULT false,
  is_active              boolean NOT NULL DEFAULT true,
  created_at             timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT users_firm_xor_client CHECK (
    (firm_id IS NOT NULL AND client_id IS NULL) OR
    (firm_id IS NULL AND client_id IS NOT NULL)
  ),
  CONSTRAINT users_client_role CHECK (
    (role IN ('ClientAdmin', 'ClientUser') AND client_id IS NOT NULL) OR
    (role NOT IN ('ClientAdmin', 'ClientUser') AND firm_id IS NOT NULL)
  )
);

CREATE INDEX idx_users_firm_id ON users(firm_id);
CREATE INDEX idx_users_client_id ON users(client_id);
CREATE INDEX idx_users_email ON users(email);

ALTER TABLE users ENABLE ROW LEVEL SECURITY;
CREATE POLICY users_isolation ON users
  USING (firm_id = current_firm_id() OR firm_id IS NULL);

-- Clients
CREATE TABLE clients (
  id                    uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  firm_id               uuid NOT NULL REFERENCES firms(id),
  name                  text NOT NULL,
  industry              text,
  primary_contact_email text,
  created_at            timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_clients_firm_id ON clients(firm_id);

ALTER TABLE clients ENABLE ROW LEVEL SECURITY;
CREATE POLICY clients_isolation ON clients
  USING (firm_id = current_firm_id());

-- Add FK from users.client_id to clients
ALTER TABLE users ADD CONSTRAINT fk_users_client
  FOREIGN KEY (client_id) REFERENCES clients(id);

-- Invitations
CREATE TABLE invitations (
  id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  firm_id         uuid NOT NULL REFERENCES firms(id),
  email           text NOT NULL,
  assigned_role   user_role NOT NULL,
  token_hash      text NOT NULL UNIQUE,
  status          invitation_status NOT NULL DEFAULT 'Sent',
  expires_at      timestamptz NOT NULL,
  reminder_sent_at timestamptz,
  invited_by_id   uuid NOT NULL REFERENCES users(id),
  accepted_at     timestamptz,
  created_at      timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_invitations_firm_id ON invitations(firm_id);
CREATE INDEX idx_invitations_token_hash ON invitations(token_hash);

ALTER TABLE invitations ENABLE ROW LEVEL SECURITY;
CREATE POLICY invitations_isolation ON invitations
  USING (firm_id = current_firm_id());
```

`000002_identity_tables.down.sql`:

```sql
DROP TABLE IF EXISTS invitations;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS clients;
DROP TABLE IF EXISTS firms;
DROP FUNCTION IF EXISTS current_firm_id();
```

- [ ] **Step 3: Install golang-migrate and run migrations**

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

Run migrations against local Postgres:

```bash
migrate -path apps/axiom-api/migrations -database "postgres://axiom_svc:localdev@localhost:5432/axiom_db?sslmode=disable" up
```

Expected: `1/u identity_enums` and `2/u identity_tables` applied successfully.

- [ ] **Step 4: Verify tables and RLS**

```bash
docker compose exec postgres psql -U axiom_svc -d axiom_db -c "\dt"
```

Expected: `firms`, `users`, `clients`, `invitations` tables listed.

Verify RLS is enabled:

```bash
docker compose exec postgres psql -U axiom_svc -d axiom_db -c "SELECT tablename, rowsecurity FROM pg_tables WHERE schemaname = 'public' AND tablename IN ('firms', 'users', 'clients', 'invitations');"
```

Expected: all four tables show `rowsecurity = t`.

- [ ] **Step 5: Commit**

```bash
git add apps/axiom-api/migrations/
git commit -m "feat: add identity enum and table migrations with RLS policies"
```

---

### Task 7: sqlc setup and identity queries

**Files:**
- Create: `apps/axiom-api/sqlc.yaml`
- Create: `apps/axiom-api/internal/identity/queries/firms.sql`
- Create: `apps/axiom-api/internal/identity/queries/users.sql`
- Create: `apps/axiom-api/internal/identity/queries/clients.sql`
- Create: `apps/axiom-api/internal/identity/queries/invitations.sql`

- [ ] **Step 1: Create sqlc.yaml**

```yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "internal/identity/queries"
    schema: "migrations"
    gen:
      go:
        package: "queries"
        out: "internal/identity/queries"
        sql_package: "pgx/v5"
        emit_json_tags: true
        emit_empty_slices: true
        overrides:
          - db_type: "uuid"
            go_type: "github.com/google/uuid.UUID"
          - db_type: "timestamptz"
            go_type: "time.Time"
          - db_type: "jsonb"
            go_type: "json.RawMessage"
            import: "encoding/json"
```

- [ ] **Step 2: Create firms.sql queries**

```sql
-- name: CreateFirm :one
INSERT INTO firms (name, slug, billing_contact_email, country, staff_count_range, primary_audit_types)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetFirmByID :one
SELECT * FROM firms WHERE id = $1;

-- name: UpdateFirm :one
UPDATE firms
SET name = COALESCE(sqlc.narg('name'), name),
    logo_url = COALESCE(sqlc.narg('logo_url'), logo_url),
    timezone = COALESCE(sqlc.narg('timezone'), timezone),
    billing_contact_email = COALESCE(sqlc.narg('billing_contact_email'), billing_contact_email),
    settings = COALESCE(sqlc.narg('settings'), settings)
WHERE id = $1
RETURNING *;
```

- [ ] **Step 3: Create users.sql queries**

```sql
-- name: CreateUser :one
INSERT INTO users (firm_id, email, display_name, role, auth_method, password_hash, notification_frequency)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: ListUsersByFirmID :many
SELECT * FROM users
WHERE firm_id = $1 AND is_active = true
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateUser :one
UPDATE users
SET role = COALESCE(sqlc.narg('role'), role),
    display_name = COALESCE(sqlc.narg('display_name'), display_name),
    notification_frequency = COALESCE(sqlc.narg('notification_frequency'), notification_frequency),
    is_active = COALESCE(sqlc.narg('is_active'), is_active)
WHERE id = $1
RETURNING *;

-- name: DeactivateUser :exec
UPDATE users SET is_active = false WHERE id = $1;
```

- [ ] **Step 4: Create clients.sql queries**

```sql
-- name: CreateClient :one
INSERT INTO clients (firm_id, name, industry, primary_contact_email)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetClientByID :one
SELECT * FROM clients WHERE id = $1;

-- name: ListClientsByFirmID :many
SELECT * FROM clients
WHERE firm_id = $1
ORDER BY name ASC
LIMIT $2 OFFSET $3;

-- name: UpdateClient :one
UPDATE clients
SET name = COALESCE(sqlc.narg('name'), name),
    industry = COALESCE(sqlc.narg('industry'), industry),
    primary_contact_email = COALESCE(sqlc.narg('primary_contact_email'), primary_contact_email)
WHERE id = $1
RETURNING *;
```

- [ ] **Step 5: Create invitations.sql queries**

```sql
-- name: CreateInvitation :one
INSERT INTO invitations (firm_id, email, assigned_role, token_hash, expires_at, invited_by_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetInvitationByTokenHash :one
SELECT * FROM invitations WHERE token_hash = $1;

-- name: ListInvitationsByFirmID :many
SELECT * FROM invitations
WHERE firm_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: AcceptInvitation :one
UPDATE invitations
SET status = 'Accepted', accepted_at = now()
WHERE id = $1 AND status = 'Sent'
RETURNING *;

-- name: CancelInvitation :exec
UPDATE invitations SET status = 'Expired' WHERE id = $1 AND status = 'Sent';
```

- [ ] **Step 6: Install sqlc and generate Go code**

```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
cd apps/axiom-api
sqlc generate
```

Expected: Go files generated in `internal/identity/queries/` — `db.go`, `models.go`, `firms.sql.go`, `users.sql.go`, `clients.sql.go`, `invitations.sql.go`.

- [ ] **Step 7: Install uuid dependency and verify build**

```bash
cd apps/axiom-api
go get github.com/google/uuid
go mod tidy
go build ./...
```

Expected: compiles successfully.

- [ ] **Step 8: Commit**

```bash
git add apps/axiom-api/sqlc.yaml apps/axiom-api/internal/identity/queries/
git commit -m "feat: add sqlc config and identity query definitions with generated Go code"
```

---

### Task 8: JWT issuer

**Files:**
- Create: `apps/axiom-api/internal/identity/jwt.go`
- Create: `apps/axiom-api/internal/identity/jwt_test.go`

- [ ] **Step 1: Write JWT test**

```go
package identity_test

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	"github.com/axiom-platform/axiom/apps/axiom-api/internal/identity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testKeyPair(t *testing.T) (*rsa.PrivateKey, *rsa.PublicKey) {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	return key, &key.PublicKey
}

func TestJWTIssueAndVerify(t *testing.T) {
	privKey, pubKey := testKeyPair(t)
	issuer := identity.NewJWTIssuer(privKey, pubKey)

	userID := uuid.New()
	firmID := uuid.New()

	pair, err := issuer.Issue(userID, firmID, "FirmAdmin")
	require.NoError(t, err)
	assert.NotEmpty(t, pair.AccessToken)
	assert.NotEmpty(t, pair.RefreshToken)

	claims, err := issuer.Verify(pair.AccessToken)
	require.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, firmID, claims.FirmID)
	assert.Equal(t, "FirmAdmin", claims.Role)
}

func TestJWTExpiredToken(t *testing.T) {
	privKey, pubKey := testKeyPair(t)
	issuer := identity.NewJWTIssuerWithDurations(privKey, pubKey, -1*time.Second, 7*24*time.Hour)

	pair, err := issuer.Issue(uuid.New(), uuid.New(), "Staff")
	require.NoError(t, err)

	_, err = issuer.Verify(pair.AccessToken)
	assert.Error(t, err)
}

func TestJWTRefresh(t *testing.T) {
	privKey, pubKey := testKeyPair(t)
	issuer := identity.NewJWTIssuer(privKey, pubKey)

	userID := uuid.New()
	firmID := uuid.New()

	pair, err := issuer.Issue(userID, firmID, "Manager")
	require.NoError(t, err)

	newPair, err := issuer.Refresh(pair.RefreshToken)
	require.NoError(t, err)
	assert.NotEmpty(t, newPair.AccessToken)
	assert.NotEqual(t, pair.AccessToken, newPair.AccessToken)
}
```

- [ ] **Step 2: Run test — expect failure**

```bash
cd apps/axiom-api
go test ./internal/identity/ -v -run TestJWT
```

Expected: compilation error — `identity.NewJWTIssuer` not found.

- [ ] **Step 3: Implement jwt.go**

```go
package identity

import (
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	FirmID uuid.UUID `json:"firm_id"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
}

type JWTIssuer struct {
	privKey         *rsa.PrivateKey
	pubKey          *rsa.PublicKey
	accessDuration  time.Duration
	refreshDuration time.Duration
}

func NewJWTIssuer(privKey *rsa.PrivateKey, pubKey *rsa.PublicKey) *JWTIssuer {
	return &JWTIssuer{
		privKey:         privKey,
		pubKey:          pubKey,
		accessDuration:  15 * time.Minute,
		refreshDuration: 7 * 24 * time.Hour,
	}
}

func NewJWTIssuerWithDurations(privKey *rsa.PrivateKey, pubKey *rsa.PublicKey, access, refresh time.Duration) *JWTIssuer {
	return &JWTIssuer{
		privKey:         privKey,
		pubKey:          pubKey,
		accessDuration:  access,
		refreshDuration: refresh,
	}
}

func (j *JWTIssuer) Issue(userID, firmID uuid.UUID, role string) (*TokenPair, error) {
	now := time.Now()

	accessClaims := &Claims{
		UserID: userID,
		FirmID: firmID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(j.accessDuration)),
			IssuedAt:  jwt.NewNumericDate(now),
			Subject:   userID.String(),
			Issuer:    "axiom",
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodRS256, accessClaims)
	accessStr, err := accessToken.SignedString(j.privKey)
	if err != nil {
		return nil, fmt.Errorf("sign access token: %w", err)
	}

	refreshClaims := &Claims{
		UserID: userID,
		FirmID: firmID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(j.refreshDuration)),
			IssuedAt:  jwt.NewNumericDate(now),
			Subject:   userID.String(),
			Issuer:    "axiom",
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodRS256, refreshClaims)
	refreshStr, err := refreshToken.SignedString(j.privKey)
	if err != nil {
		return nil, fmt.Errorf("sign refresh token: %w", err)
	}

	return &TokenPair{AccessToken: accessStr, RefreshToken: refreshStr}, nil
}

func (j *JWTIssuer) Verify(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return j.pubKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}
	return claims, nil
}

func (j *JWTIssuer) Refresh(refreshTokenStr string) (*TokenPair, error) {
	claims, err := j.Verify(refreshTokenStr)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}
	return j.Issue(claims.UserID, claims.FirmID, claims.Role)
}
```

- [ ] **Step 4: Install dependencies and run tests**

```bash
cd apps/axiom-api
go get github.com/golang-jwt/jwt/v5
go get github.com/stretchr/testify
go mod tidy
go test ./internal/identity/ -v -run TestJWT
```

Expected: all 3 tests PASS.

- [ ] **Step 5: Commit**

```bash
git add apps/axiom-api/internal/identity/jwt.go apps/axiom-api/internal/identity/jwt_test.go apps/axiom-api/go.mod apps/axiom-api/go.sum
git commit -m "feat: add JWT issuer with RSA signing, verification, and refresh"
```

---

### Task 9: Gateway middleware

**Files:**
- Create: `apps/axiom-api/internal/gateway/middleware.go`
- Create: `apps/axiom-api/internal/gateway/middleware_test.go`

- [ ] **Step 1: Write middleware tests**

```go
package gateway_test

import (
	"crypto/rand"
	"crypto/rsa"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/axiom-platform/axiom/apps/axiom-api/internal/gateway"
	"github.com/axiom-platform/axiom/apps/axiom-api/internal/identity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T) (*identity.JWTIssuer, *gateway.Middleware) {
	t.Helper()
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	issuer := identity.NewJWTIssuer(privKey, &privKey.PublicKey)
	mw := gateway.NewMiddleware(issuer)
	return issuer, mw
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	issuer, mw := setup(t)
	userID := uuid.New()
	firmID := uuid.New()
	pair, err := issuer.Issue(userID, firmID, "Staff")
	require.NoError(t, err)

	handler := mw.Auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := gateway.GetClaims(r.Context())
		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, firmID, claims.FirmID)
		assert.Equal(t, "Staff", claims.Role)
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+pair.AccessToken)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestAuthMiddleware_MissingToken(t *testing.T) {
	_, mw := setup(t)
	handler := mw.Auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called")
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestWithRole_Allowed(t *testing.T) {
	issuer, mw := setup(t)
	pair, _ := issuer.Issue(uuid.New(), uuid.New(), "Partner")

	handler := mw.Auth(mw.WithRole("Partner", "FirmAdmin")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})))

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+pair.AccessToken)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestWithRole_Denied(t *testing.T) {
	issuer, mw := setup(t)
	pair, _ := issuer.Issue(uuid.New(), uuid.New(), "Staff")

	handler := mw.Auth(mw.WithRole("FirmAdmin")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called")
	})))

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+pair.AccessToken)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusForbidden, rr.Code)
}
```

- [ ] **Step 2: Run tests — expect failure**

```bash
cd apps/axiom-api
go test ./internal/gateway/ -v
```

Expected: compilation error — `gateway.NewMiddleware` not found.

- [ ] **Step 3: Implement middleware.go**

```go
package gateway

import (
	"context"
	"net/http"
	"strings"

	"github.com/axiom-platform/axiom/apps/axiom-api/internal/identity"
	"github.com/axiom-platform/axiom/apps/axiom-api/internal/platform"
)

type contextKey string

const claimsKey contextKey = "claims"

type Middleware struct {
	jwtIssuer *identity.JWTIssuer
}

func NewMiddleware(jwtIssuer *identity.JWTIssuer) *Middleware {
	return &Middleware{jwtIssuer: jwtIssuer}
}

func GetClaims(ctx context.Context) *identity.Claims {
	claims, _ := ctx.Value(claimsKey).(*identity.Claims)
	return claims
}

func (m *Middleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			platform.WriteError(w, platform.ErrUnauthorized("missing authorization header"))
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			platform.WriteError(w, platform.ErrUnauthorized("invalid authorization header format"))
			return
		}

		claims, err := m.jwtIssuer.Verify(parts[1])
		if err != nil {
			platform.WriteError(w, platform.ErrUnauthorized("invalid or expired token"))
			return
		}

		ctx := context.WithValue(r.Context(), claimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Middleware) WithRole(roles ...string) func(http.Handler) http.Handler {
	allowed := make(map[string]bool, len(roles))
	for _, r := range roles {
		allowed[r] = true
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := GetClaims(r.Context())
			if claims == nil {
				platform.WriteError(w, platform.ErrUnauthorized("no claims in context"))
				return
			}
			if !allowed[claims.Role] {
				platform.WriteError(w, platform.ErrForbidden("insufficient permissions"))
				return
			}
			next.ServeHTTP(w, r.WithContext(r.Context()))
		})
	}
}
```

- [ ] **Step 4: Run tests — expect pass**

```bash
cd apps/axiom-api
go test ./internal/gateway/ -v
```

Expected: all 4 tests PASS.

- [ ] **Step 5: Commit**

```bash
git add apps/axiom-api/internal/gateway/
git commit -m "feat: add gateway middleware for JWT auth and role-based access control"
```

---

### Task 10: Identity service — registration and login

**Files:**
- Create: `apps/axiom-api/internal/identity/service.go`
- Create: `apps/axiom-api/internal/identity/service_test.go`
- Create: `apps/axiom-api/internal/platform/testdb.go`

- [ ] **Step 1: Create test database helper**

`apps/axiom-api/internal/platform/testdb.go`:

```go
package platform

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TestDB creates a temporary test database, runs migrations, and returns a pool.
// The database is dropped when the test completes.
func TestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()
	ctx := context.Background()

	// Connect to default database to create test DB
	adminURL := "postgres://axiom_svc:localdev@localhost:5432/axiom_db?sslmode=disable"
	adminPool, err := pgxpool.New(ctx, adminURL)
	if err != nil {
		t.Fatalf("connect to admin db: %v", err)
	}

	dbName := fmt.Sprintf("axiom_test_%s", t.Name())
	// Clean the name for Postgres identifier rules
	cleanName := ""
	for _, c := range dbName {
		if (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '_' {
			cleanName += string(c)
		} else if c >= 'A' && c <= 'Z' {
			cleanName += string(c - 'A' + 'a')
		} else {
			cleanName += "_"
		}
	}
	dbName = cleanName

	_, _ = adminPool.Exec(ctx, fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName))
	_, err = adminPool.Exec(ctx, fmt.Sprintf("CREATE DATABASE %s", dbName))
	if err != nil {
		t.Fatalf("create test database: %v", err)
	}
	adminPool.Close()

	testURL := fmt.Sprintf("postgres://axiom_svc:localdev@localhost:5432/%s?sslmode=disable", dbName)

	// Run migrations
	m, err := migrate.New("file://../../migrations", testURL)
	if err != nil {
		t.Fatalf("create migrator: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		t.Fatalf("run migrations: %v", err)
	}
	srcErr, dbErr := m.Close()
	if srcErr != nil {
		t.Fatalf("close migrator source: %v", srcErr)
	}
	if dbErr != nil {
		t.Fatalf("close migrator db: %v", dbErr)
	}

	pool, err := pgxpool.New(ctx, testURL)
	if err != nil {
		t.Fatalf("connect to test db: %v", err)
	}

	t.Cleanup(func() {
		pool.Close()
		// Reconnect to admin DB to drop test DB
		cleanupPool, err := pgxpool.New(context.Background(), "postgres://axiom_svc:localdev@localhost:5432/axiom_db?sslmode=disable")
		if err == nil {
			cleanupPool.Exec(context.Background(), fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName))
			cleanupPool.Close()
		}
	})

	return pool
}
```

- [ ] **Step 2: Write service tests for registration and login**

`apps/axiom-api/internal/identity/service_test.go`:

```go
package identity_test

import (
	"context"
	"testing"

	"github.com/axiom-platform/axiom/apps/axiom-api/internal/identity"
	"github.com/axiom-platform/axiom/apps/axiom-api/internal/platform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupService(t *testing.T) *identity.Service {
	t.Helper()
	pool := platform.TestDB(t)
	privKey, pubKey := testKeyPair(t)
	issuer := identity.NewJWTIssuer(privKey, pubKey)
	return identity.NewService(pool, issuer)
}

func TestRegisterFirm(t *testing.T) {
	svc := setupService(t)
	ctx := context.Background()

	result, err := svc.RegisterFirm(ctx, identity.RegisterFirmInput{
		FirmName:     "Acme CPAs",
		AdminEmail:   "admin@acme.com",
		AdminName:    "Alice Admin",
		Password:     "securepassword123",
		Country:      "US",
		StaffCount:   "21-40",
		AuditTypes:   []string{"FinancialAudit", "SOC2"},
	})
	require.NoError(t, err)
	assert.Equal(t, "Acme CPAs", result.Firm.Name)
	assert.Equal(t, "admin@acme.com", result.User.Email)
	assert.Equal(t, "FirmAdmin", result.User.Role)
	assert.NotEmpty(t, result.Tokens.AccessToken)
}

func TestRegisterFirm_DuplicateEmail(t *testing.T) {
	svc := setupService(t)
	ctx := context.Background()

	input := identity.RegisterFirmInput{
		FirmName:   "Firm A",
		AdminEmail: "dup@example.com",
		AdminName:  "User A",
		Password:   "password123",
		Country:    "US",
		StaffCount: "1-10",
		AuditTypes: []string{"SOC2"},
	}

	_, err := svc.RegisterFirm(ctx, input)
	require.NoError(t, err)

	_, err = svc.RegisterFirm(ctx, input)
	assert.Error(t, err)
}

func TestLogin(t *testing.T) {
	svc := setupService(t)
	ctx := context.Background()

	_, err := svc.RegisterFirm(ctx, identity.RegisterFirmInput{
		FirmName:   "Test Firm",
		AdminEmail: "login@test.com",
		AdminName:  "Login User",
		Password:   "mypassword",
		Country:    "US",
		StaffCount: "1-10",
		AuditTypes: []string{"SOC2"},
	})
	require.NoError(t, err)

	result, err := svc.Login(ctx, "login@test.com", "mypassword")
	require.NoError(t, err)
	assert.NotEmpty(t, result.AccessToken)
}

func TestLogin_WrongPassword(t *testing.T) {
	svc := setupService(t)
	ctx := context.Background()

	_, err := svc.RegisterFirm(ctx, identity.RegisterFirmInput{
		FirmName:   "Test Firm",
		AdminEmail: "wrong@test.com",
		AdminName:  "User",
		Password:   "correctpassword",
		Country:    "US",
		StaffCount: "1-10",
		AuditTypes: []string{"SOC2"},
	})
	require.NoError(t, err)

	_, err = svc.Login(ctx, "wrong@test.com", "wrongpassword")
	assert.Error(t, err)
}
```

- [ ] **Step 3: Run tests — expect failure**

```bash
cd apps/axiom-api
go test ./internal/identity/ -v -run "TestRegister|TestLogin" -count=1
```

Expected: compilation error — `identity.NewService` not found.

- [ ] **Step 4: Implement service.go**

```go
package identity

import (
	"context"
	"fmt"
	"strings"

	"github.com/axiom-platform/axiom/apps/axiom-api/internal/identity/queries"
	"github.com/axiom-platform/axiom/apps/axiom-api/internal/platform"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	pool    *pgxpool.Pool
	queries *queries.Queries
	jwt     *JWTIssuer
}

func NewService(pool *pgxpool.Pool, jwt *JWTIssuer) *Service {
	return &Service{
		pool:    pool,
		queries: queries.New(pool),
		jwt:     jwt,
	}
}

type RegisterFirmInput struct {
	FirmName   string
	AdminEmail string
	AdminName  string
	Password   string
	Country    string
	StaffCount string
	AuditTypes []string
}

type RegisterFirmResult struct {
	Firm   FirmDTO
	User   UserDTO
	Tokens *TokenPair
}

type FirmDTO struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Slug string    `json:"slug"`
}

type UserDTO struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Name  string    `json:"display_name"`
	Role  string    `json:"role"`
}

func (s *Service) RegisterFirm(ctx context.Context, input RegisterFirmInput) (*RegisterFirmResult, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	slug := generateSlug(input.FirmName)

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := s.queries.WithTx(tx)

	// Create firm (bypass RLS — no firm_id set yet)
	firm, err := qtx.CreateFirm(ctx, queries.CreateFirmParams{
		Name:                input.FirmName,
		Slug:                slug,
		BillingContactEmail: input.AdminEmail,
		Country:             input.Country,
		StaffCountRange:     ptrStr(input.StaffCount),
	})
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return nil, platform.ErrConflict("firm with this name already exists")
		}
		return nil, fmt.Errorf("create firm: %w", err)
	}

	// Set RLS context for subsequent operations
	if _, err := tx.Exec(ctx, "SELECT set_config('app.current_firm_id', $1, true)", firm.ID.String()); err != nil {
		return nil, fmt.Errorf("set firm context: %w", err)
	}

	// Create admin user
	user, err := qtx.CreateUser(ctx, queries.CreateUserParams{
		FirmID:                firm.ID,
		Email:                 input.AdminEmail,
		DisplayName:           input.AdminName,
		Role:                  "FirmAdmin",
		AuthMethod:            "Password",
		PasswordHash:          ptrStr(string(hash)),
		NotificationFrequency: "RealTime",
	})
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return nil, platform.ErrConflict("email already registered")
		}
		return nil, fmt.Errorf("create user: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	tokens, err := s.jwt.Issue(user.ID, firm.ID, string(user.Role))
	if err != nil {
		return nil, fmt.Errorf("issue tokens: %w", err)
	}

	return &RegisterFirmResult{
		Firm:   FirmDTO{ID: firm.ID, Name: firm.Name, Slug: firm.Slug},
		User:   UserDTO{ID: user.ID, Email: user.Email, Name: user.DisplayName, Role: string(user.Role)},
		Tokens: tokens,
	}, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (*TokenPair, error) {
	user, err := s.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, platform.ErrUnauthorized("invalid email or password")
		}
		return nil, fmt.Errorf("get user: %w", err)
	}

	if !user.IsActive {
		return nil, platform.ErrUnauthorized("account is deactivated")
	}

	if user.PasswordHash == nil {
		return nil, platform.ErrUnauthorized("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(password)); err != nil {
		return nil, platform.ErrUnauthorized("invalid email or password")
	}

	firmID := uuid.Nil
	if user.FirmID != nil {
		firmID = *user.FirmID
	}

	return s.jwt.Issue(user.ID, firmID, string(user.Role))
}

func generateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			return r
		}
		return '-'
	}, slug)
	// Remove consecutive dashes and trim
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}
	slug = strings.Trim(slug, "-")
	// Append short random suffix for uniqueness
	slug = slug + "-" + uuid.New().String()[:8]
	return slug
}

func ptrStr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
```

Note: the sqlc generated code in `queries/` uses types that may need adjustment. The implementing agent should check the generated `CreateFirmParams` and `CreateUserParams` structs and adjust the service code to match. The key pattern is: sqlc generates Go structs from the SQL, and the service layer translates between domain types and sqlc types.

- [ ] **Step 5: Install bcrypt dependency and run tests**

```bash
cd apps/axiom-api
go get golang.org/x/crypto/bcrypt
go get github.com/golang-migrate/migrate/v4
go mod tidy
go test ./internal/identity/ -v -run "TestRegister|TestLogin" -count=1
```

Expected: all 4 tests PASS. If sqlc generated types differ from what the service expects, adjust the service code to match the generated types.

- [ ] **Step 6: Commit**

```bash
git add apps/axiom-api/internal/identity/service.go apps/axiom-api/internal/identity/service_test.go apps/axiom-api/internal/platform/testdb.go apps/axiom-api/go.mod apps/axiom-api/go.sum
git commit -m "feat: add identity service with firm registration and login against real Postgres"
```

---

### Task 11: Identity HTTP handlers

**Files:**
- Create: `apps/axiom-api/internal/identity/handler.go`
- Create: `apps/axiom-api/internal/identity/handler_test.go`

- [ ] **Step 1: Write handler tests**

```go
package identity_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/axiom-platform/axiom/apps/axiom-api/internal/identity"
	"github.com/axiom-platform/axiom/apps/axiom-api/internal/platform"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupHandler(t *testing.T) http.Handler {
	t.Helper()
	pool := platform.TestDB(t)
	privKey, pubKey := testKeyPair(t)
	issuer := identity.NewJWTIssuer(privKey, pubKey)
	svc := identity.NewService(pool, issuer)
	handler := identity.NewHandler(svc, issuer)

	r := chi.NewRouter()
	handler.RegisterRoutes(r)
	return r
}

func TestHandler_RegisterAndLogin(t *testing.T) {
	router := setupHandler(t)

	// Register
	regBody, _ := json.Marshal(map[string]any{
		"firm_name":        "Handler Test Firm",
		"admin_email":      "handler@test.com",
		"admin_name":       "Handler Admin",
		"password":         "testpass123",
		"country":          "US",
		"staff_count_range": "1-10",
		"primary_audit_types": []string{"SOC2"},
	})
	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewReader(regBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var regResp map[string]any
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &regResp))
	assert.NotEmpty(t, regResp["access_token"])

	// Login
	loginBody, _ := json.Marshal(map[string]string{
		"email":    "handler@test.com",
		"password": "testpass123",
	})
	req = httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(loginBody))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var loginResp map[string]any
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &loginResp))
	assert.NotEmpty(t, loginResp["access_token"])
}
```

- [ ] **Step 2: Run test — expect failure**

```bash
cd apps/axiom-api
go test ./internal/identity/ -v -run TestHandler -count=1
```

Expected: compilation error — `identity.NewHandler` not found.

- [ ] **Step 3: Implement handler.go**

```go
package identity

import (
	"encoding/json"
	"net/http"

	"github.com/axiom-platform/axiom/apps/axiom-api/internal/platform"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	svc *Service
	jwt *JWTIssuer
}

func NewHandler(svc *Service, jwt *JWTIssuer) *Handler {
	return &Handler{svc: svc, jwt: jwt}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Post("/api/v1/auth/register", h.register)
	r.Post("/api/v1/auth/login", h.login)
	r.Post("/api/v1/auth/refresh", h.refresh)
}

type registerRequest struct {
	FirmName        string   `json:"firm_name"`
	AdminEmail      string   `json:"admin_email"`
	AdminName       string   `json:"admin_name"`
	Password        string   `json:"password"`
	Country         string   `json:"country"`
	StaffCountRange string   `json:"staff_count_range"`
	PrimaryAuditTypes []string `json:"primary_audit_types"`
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		platform.WriteError(w, platform.ErrBadRequest("invalid request body"))
		return
	}

	if req.FirmName == "" || req.AdminEmail == "" || req.Password == "" {
		platform.WriteError(w, platform.ErrValidation("missing required fields", "firm_name, admin_email, and password are required"))
		return
	}

	result, err := h.svc.RegisterFirm(r.Context(), RegisterFirmInput{
		FirmName:   req.FirmName,
		AdminEmail: req.AdminEmail,
		AdminName:  req.AdminName,
		Password:   req.Password,
		Country:    req.Country,
		StaffCount: req.StaffCountRange,
		AuditTypes: req.PrimaryAuditTypes,
	})
	if err != nil {
		platform.WriteError(w, err)
		return
	}

	platform.WriteJSON(w, http.StatusCreated, map[string]any{
		"firm":          result.Firm,
		"user":          result.User,
		"access_token":  result.Tokens.AccessToken,
		"refresh_token": result.Tokens.RefreshToken,
	})
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		platform.WriteError(w, platform.ErrBadRequest("invalid request body"))
		return
	}

	tokens, err := h.svc.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		platform.WriteError(w, err)
		return
	}

	platform.WriteJSON(w, http.StatusOK, tokens)
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *Handler) refresh(w http.ResponseWriter, r *http.Request) {
	var req refreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		platform.WriteError(w, platform.ErrBadRequest("invalid request body"))
		return
	}

	tokens, err := h.jwt.Refresh(req.RefreshToken)
	if err != nil {
		platform.WriteError(w, platform.ErrUnauthorized("invalid refresh token"))
		return
	}

	platform.WriteJSON(w, http.StatusOK, tokens)
}
```

- [ ] **Step 4: Run tests — expect pass**

```bash
cd apps/axiom-api
go test ./internal/identity/ -v -run TestHandler -count=1
```

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add apps/axiom-api/internal/identity/handler.go apps/axiom-api/internal/identity/handler_test.go
git commit -m "feat: add identity HTTP handlers for registration, login, and token refresh"
```

---

### Task 12: Wire everything in main.go and verify end-to-end

**Files:**
- Modify: `apps/axiom-api/cmd/server/main.go`

- [ ] **Step 1: Update main.go to wire all modules**

Replace `apps/axiom-api/cmd/server/main.go` with the full wiring:

```go
package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"github.com/axiom-platform/axiom/apps/axiom-api/internal/gateway"
	"github.com/axiom-platform/axiom/apps/axiom-api/internal/identity"
	"github.com/axiom-platform/axiom/apps/axiom-api/internal/platform"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg, err := platform.LoadConfig()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	ctx := context.Background()
	pool, err := platform.NewDBPool(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()
	slog.Info("connected to database")

	// JWT key pair — in production, load from config/secrets.
	// For local dev, generate ephemeral keys if not configured.
	var privKey *rsa.PrivateKey
	var pubKey *rsa.PublicKey
	if cfg.JWTPrivKey != "" {
		// TODO: parse PEM from config
		slog.Info("using configured JWT keys")
	}
	// Fallback: generate ephemeral keys for development
	if privKey == nil {
		slog.Warn("generating ephemeral JWT keys — tokens will not survive restarts")
		privKey, err = rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			slog.Error("failed to generate RSA key", "error", err)
			os.Exit(1)
		}
		pubKey = &privKey.PublicKey
	}

	jwtIssuer := identity.NewJWTIssuer(privKey, pubKey)
	identitySvc := identity.NewService(pool, jwtIssuer)
	identityHandler := identity.NewHandler(identitySvc, jwtIssuer)
	gw := gateway.NewMiddleware(jwtIssuer)

	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Location"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health check
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	// Public routes (no auth required)
	identityHandler.RegisterRoutes(r)

	// Authenticated routes (JWT required)
	r.Group(func(r chi.Router) {
		r.Use(gw.Auth)
		// Authenticated identity endpoints will be added here as they are built
		// e.g., GET /api/v1/firms/current, GET /api/v1/users, etc.
	})

	slog.Info("starting server", "port", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}
```

- [ ] **Step 2: Install cors dependency and verify build**

```bash
cd apps/axiom-api
go get github.com/go-chi/cors
go mod tidy
go build ./cmd/server
```

Expected: builds successfully.

- [ ] **Step 3: Run migrations and start server**

```bash
migrate -path apps/axiom-api/migrations -database "postgres://axiom_svc:localdev@localhost:5432/axiom_db?sslmode=disable" up
cd apps/axiom-api && go run ./cmd/server
```

- [ ] **Step 4: Test end-to-end with curl**

Register a firm:
```bash
curl -s -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"firm_name":"Curl Test Firm","admin_email":"curl@test.com","admin_name":"Curl Admin","password":"testpass123","country":"US","staff_count_range":"1-10","primary_audit_types":["SOC2"]}' | jq .
```

Expected: 201 response with `firm`, `user`, `access_token`, `refresh_token`.

Login:
```bash
curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"curl@test.com","password":"testpass123"}' | jq .
```

Expected: 200 with `access_token` and `refresh_token`.

- [ ] **Step 5: Commit**

```bash
git add apps/axiom-api/cmd/server/main.go apps/axiom-api/go.mod apps/axiom-api/go.sum
git commit -m "feat: wire identity module into main server with auth, CORS, and healthz"
```

---

### Task 13: Authenticated identity endpoints (firms, users, clients)

**Files:**
- Modify: `apps/axiom-api/internal/identity/service.go` (add methods)
- Modify: `apps/axiom-api/internal/identity/handler.go` (add routes)
- Modify: `apps/axiom-api/internal/identity/service_test.go` (add tests)
- Modify: `apps/axiom-api/cmd/server/main.go` (wire authenticated routes)

This task adds the remaining authenticated endpoints: get/update firm, list/get/update users, create/list/get clients, and invitation CRUD. These all follow the same pattern as registration/login — service method, handler, test. The implementing agent should:

- [ ] **Step 1: Add service methods for firm, user, client, and invitation CRUD**

Add these methods to `Service` in `service.go`:
- `GetFirm(ctx, firmID) → FirmDTO`
- `UpdateFirm(ctx, firmID, input) → FirmDTO`
- `ListUsers(ctx, firmID, limit, offset) → []UserDTO`
- `GetUser(ctx, userID) → UserDTO`
- `UpdateUser(ctx, userID, input) → UserDTO`
- `DeactivateUser(ctx, userID)`
- `CreateClient(ctx, firmID, input) → ClientDTO`
- `ListClients(ctx, firmID, limit, offset) → []ClientDTO`
- `GetClient(ctx, clientID) → ClientDTO`
- `UpdateClient(ctx, clientID, input) → ClientDTO`
- `CreateInvitation(ctx, firmID, inviterID, input) → InvitationDTO`
- `ListInvitations(ctx, firmID, limit, offset) → []InvitationDTO`
- `ValidateInvitationToken(ctx, token) → InvitationDTO`
- `AcceptInvitation(ctx, token, password) → (UserDTO, TokenPair)`
- `CancelInvitation(ctx, invitationID)`

Each method sets the RLS `app.current_firm_id` before querying, uses the sqlc-generated queries, and returns domain DTOs.

- [ ] **Step 2: Write integration tests for each service method**

Test each method in `service_test.go` against a real test database. Key tests:
- `TestGetFirm` — register, then get the firm
- `TestListUsers` — register, create invitations, accept them, list users
- `TestCreateClient` — register, create client, verify RLS isolation
- `TestInvitationFlow` — create invitation → validate token → accept → user created
- `TestRLSIsolation` — register two firms, verify each firm only sees its own data

- [ ] **Step 3: Run tests — expect pass**

```bash
cd apps/axiom-api
go test ./internal/identity/ -v -count=1
```

Expected: all tests PASS.

- [ ] **Step 4: Add authenticated routes to handler.go**

Add a `RegisterAuthenticatedRoutes(r chi.Router)` method that registers:
- `GET /api/v1/firms/current`
- `PATCH /api/v1/firms/current`
- `GET /api/v1/users`
- `GET /api/v1/users/me`
- `PATCH /api/v1/users/me`
- `GET /api/v1/users/{userId}`
- `PATCH /api/v1/users/{userId}` (FirmAdmin only)
- `POST /api/v1/users/{userId}/deactivate` (FirmAdmin only)
- `GET /api/v1/clients`
- `POST /api/v1/clients`
- `GET /api/v1/clients/{clientId}`
- `PATCH /api/v1/clients/{clientId}`
- `GET /api/v1/invitations` (FirmAdmin only)
- `POST /api/v1/invitations` (FirmAdmin only)
- `DELETE /api/v1/invitations/{invitationId}` (FirmAdmin only)
- `POST /api/v1/invitations/{invitationId}/resend` (FirmAdmin only)

Each handler extracts claims from context via `gateway.GetClaims(r.Context())`, sets the RLS firm ID, delegates to the service, and returns JSON.

Use `gw.WithRole(...)` middleware for role-restricted endpoints.

- [ ] **Step 5: Write handler tests for authenticated endpoints**

Test key authenticated flows in `handler_test.go`:
- Register → use token to GET /firms/current
- Register → create client → list clients
- Register → create invitation → validate token → accept

- [ ] **Step 6: Wire authenticated routes in main.go**

Update the authenticated route group in `main.go`:
```go
r.Group(func(r chi.Router) {
    r.Use(gw.Auth)
    identityHandler.RegisterAuthenticatedRoutes(r)
})
```

- [ ] **Step 7: Run all tests**

```bash
cd apps/axiom-api
go test ./... -v -count=1
```

Expected: all tests PASS.

- [ ] **Step 8: Commit**

```bash
git add apps/axiom-api/
git commit -m "feat: add authenticated identity endpoints for firms, users, clients, and invitations"
```

---

### Task 14: React auth context and API client

**Files:**
- Create: `apps/web/src/api/client.ts`
- Create: `apps/web/src/hooks/use-auth.ts`
- Create: `apps/web/src/components/protected-route.tsx`
- Modify: `apps/web/src/App.tsx`
- Modify: `apps/web/src/main.tsx`

- [ ] **Step 1: Create API client with JWT injection**

`apps/web/src/api/client.ts`:

```typescript
const API_BASE = '/api/v1'

let accessToken: string | null = localStorage.getItem('access_token')
let refreshToken: string | null = localStorage.getItem('refresh_token')

export function setTokens(access: string, refresh: string) {
  accessToken = access
  refreshToken = refresh
  localStorage.setItem('access_token', access)
  localStorage.setItem('refresh_token', refresh)
}

export function clearTokens() {
  accessToken = null
  refreshToken = null
  localStorage.removeItem('access_token')
  localStorage.removeItem('refresh_token')
}

export async function api<T>(path: string, options: RequestInit = {}): Promise<T> {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(options.headers as Record<string, string>),
  }

  if (accessToken) {
    headers['Authorization'] = `Bearer ${accessToken}`
  }

  const res = await fetch(`${API_BASE}${path}`, { ...options, headers })

  if (res.status === 401 && refreshToken) {
    // Try refresh
    const refreshRes = await fetch(`${API_BASE}/auth/refresh`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ refresh_token: refreshToken }),
    })
    if (refreshRes.ok) {
      const data = await refreshRes.json()
      setTokens(data.access_token, data.refresh_token)
      headers['Authorization'] = `Bearer ${data.access_token}`
      const retryRes = await fetch(`${API_BASE}${path}`, { ...options, headers })
      if (!retryRes.ok) throw new ApiError(retryRes.status, await retryRes.text())
      return retryRes.json()
    }
    clearTokens()
    window.location.href = '/login'
    throw new ApiError(401, 'session expired')
  }

  if (!res.ok) {
    throw new ApiError(res.status, await res.text())
  }
  
  if (res.status === 204) return undefined as T
  return res.json()
}

export class ApiError extends Error {
  constructor(public status: number, public body: string) {
    super(`API error ${status}: ${body}`)
  }
}
```

- [ ] **Step 2: Create auth hook with Zustand**

`apps/web/src/hooks/use-auth.ts`:

```typescript
import { create } from 'zustand'
import { api, setTokens, clearTokens } from '../api/client'

interface User {
  id: string
  email: string
  display_name: string
  role: string
}

interface Firm {
  id: string
  name: string
  slug: string
}

interface AuthState {
  user: User | null
  firm: Firm | null
  isAuthenticated: boolean
  register: (data: {
    firm_name: string
    admin_email: string
    admin_name: string
    password: string
    country: string
    staff_count_range: string
    primary_audit_types: string[]
  }) => Promise<void>
  login: (email: string, password: string) => Promise<void>
  logout: () => void
  loadProfile: () => Promise<void>
}

export const useAuth = create<AuthState>((set) => ({
  user: null,
  firm: null,
  isAuthenticated: !!localStorage.getItem('access_token'),

  register: async (data) => {
    const res = await api<any>('/auth/register', {
      method: 'POST',
      body: JSON.stringify(data),
    })
    setTokens(res.access_token, res.refresh_token)
    set({ user: res.user, firm: res.firm, isAuthenticated: true })
  },

  login: async (email, password) => {
    const res = await api<any>('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    })
    setTokens(res.access_token, res.refresh_token)
    set({ isAuthenticated: true })
  },

  logout: () => {
    clearTokens()
    set({ user: null, firm: null, isAuthenticated: false })
  },

  loadProfile: async () => {
    const [user, firm] = await Promise.all([
      api<User>('/users/me'),
      api<Firm>('/firms/current'),
    ])
    set({ user, firm })
  },
}))
```

- [ ] **Step 3: Create protected route component**

`apps/web/src/components/protected-route.tsx`:

```tsx
import { Navigate } from 'react-router-dom'
import { useAuth } from '../hooks/use-auth'

export function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const isAuthenticated = useAuth((s) => s.isAuthenticated)
  if (!isAuthenticated) return <Navigate to="/login" replace />
  return <>{children}</>
}
```

- [ ] **Step 4: Write tests (Vitest + RTL) — TDD applies**

Install test dependencies if not already present:

```bash
cd apps/web
npm install -D vitest @vitest/ui jsdom @testing-library/react @testing-library/jest-dom @testing-library/user-event
```

Create `apps/web/vitest.config.ts` (if it doesn't exist) with `environment: 'jsdom'` and a setup file that imports `@testing-library/jest-dom`.

Write the following failing tests **before** considering the implementations above final. If you already wrote the implementation, delete it, confirm the test fails, then restore — this validates the test is actually exercising the behavior.

- `src/api/client.test.ts`:
  - `api()` injects `Authorization: Bearer <token>` when a token is set (mock `fetch`)
  - On 401 with a refresh token present, calls `/auth/refresh`, saves new tokens, and retries the original request
  - On 401 with no refresh token (or refresh fails), clears tokens and throws `ApiError(401)`
  - `setTokens`/`clearTokens` round-trip via `localStorage`

- `src/hooks/use-auth.test.tsx`:
  - `login()` calls the API, stores tokens, flips `isAuthenticated` to `true`
  - `logout()` clears tokens and resets state
  - `register()` stores `user` and `firm` from the response
  - Failed login leaves state untouched and rethrows

- `src/components/protected-route.test.tsx`:
  - When authenticated, renders `children`
  - When not authenticated, renders a `<Navigate>` to `/login` (mock router assertion or assert the returned element)

Use `msw` or plain `vi.stubGlobal('fetch', ...)` to mock API calls. Reset `localStorage` between tests.

- [ ] **Step 5: Run tests — expect pass**

```bash
cd apps/web
npm test -- --run
```

Expected: all tests PASS.

- [ ] **Step 6: Update App.tsx with routing**

```tsx
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { ProtectedRoute } from './components/protected-route'
import LoginPage from './pages/login'
import RegisterPage from './pages/register'
import DashboardPage from './pages/dashboard'

const queryClient = new QueryClient()

export default function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <Routes>
          <Route path="/login" element={<LoginPage />} />
          <Route path="/register" element={<RegisterPage />} />
          <Route path="/dashboard" element={
            <ProtectedRoute><DashboardPage /></ProtectedRoute>
          } />
          <Route path="/" element={<Navigate to="/dashboard" replace />} />
        </Routes>
      </BrowserRouter>
    </QueryClientProvider>
  )
}
```

- [ ] **Step 7: Commit**

```bash
git add apps/web/
git commit -m "feat: add React auth context, API client with token refresh, and protected routing (incl. unit tests)"
```

---

### Task 15: React pages — login, register, dashboard

**Files:**
- Create: `apps/web/src/pages/login.tsx`
- Create: `apps/web/src/pages/register.tsx`
- Create: `apps/web/src/pages/dashboard.tsx`
- Create: `apps/web/src/components/layout.tsx`

- [ ] **Step 1: Create login page**

`apps/web/src/pages/login.tsx` — form with email and password fields, calls `useAuth.login()`, navigates to `/dashboard` on success, shows error on failure, link to `/register`.

- [ ] **Step 2: Create register page**

`apps/web/src/pages/register.tsx` — form with firm name, admin email, admin name, password, country dropdown (US/CA), staff count dropdown, audit type checkboxes. Calls `useAuth.register()`, navigates to `/dashboard` on success.

- [ ] **Step 3: Create layout shell**

`apps/web/src/components/layout.tsx` — sidebar with navigation links (Dashboard, Clients, Users, Settings), top bar with firm name and user display name, logout button. Wrap protected routes in this layout.

- [ ] **Step 4: Create dashboard page**

`apps/web/src/pages/dashboard.tsx` — wrapped in `Layout`. Shows onboarding checklist:
1. Complete firm profile → link to `/settings`
2. Invite your team → link to `/users`
3. Add a client → link to `/clients`
4. Create first engagement → link to `/engagements` (disabled, Phase 2)

Fetch firm and user profile on mount via `useAuth.loadProfile()`.

- [ ] **Step 5: Write RTL tests for pages with logic — TDD applies**

For forms and components with branching behavior, write tests first. Use `@testing-library/user-event` for form interactions. Mock `useAuth` with `vi.mock(...)` so tests don't hit the network.

- `src/pages/login.test.tsx`:
  - Submitting valid credentials calls `useAuth.login` with the entered email and password
  - A rejected `login()` renders an error message
  - The "Register" link points to `/register`

- `src/pages/register.test.tsx`:
  - Submitting the full form calls `useAuth.register` with the form payload (including country, staff count, audit types)
  - Missing required fields prevents submission (or surfaces validation errors)
  - A rejected `register()` renders an error message

- `src/pages/dashboard.test.tsx`:
  - Renders the onboarding checklist with four items
  - "Create first engagement" item is disabled/marked "coming soon"
  - Calls `useAuth.loadProfile()` on mount

Layout shell is mostly static markup — cover it transitively via the dashboard test rendering `<Layout>` as a wrapper; no dedicated test required.

- [ ] **Step 6: Run tests — expect pass**

```bash
cd apps/web
npm test -- --run
```

Expected: all tests PASS.

- [ ] **Step 7: Manual browser walkthrough**

Start Go API and React dev server. Navigate to `http://localhost:3000`:
1. Redirected to `/login` (no token)
2. Click "Register" → fill form → submit → redirected to dashboard
3. See onboarding checklist
4. Click logout → redirected to login
5. Login with same credentials → back to dashboard

- [ ] **Step 8: Commit**

```bash
git add apps/web/
git commit -m "feat: add login, register, and dashboard pages with layout shell (incl. unit tests)"
```

---

### Task 16: React pages — users, clients, firm settings

**Files:**
- Create: `apps/web/src/pages/users.tsx`
- Create: `apps/web/src/pages/clients.tsx`
- Create: `apps/web/src/pages/firm-settings.tsx`
- Create: `apps/web/src/pages/accept-invitation.tsx`
- Modify: `apps/web/src/App.tsx` (add routes)

- [ ] **Step 1: Create users page**

`apps/web/src/pages/users.tsx` — table of firm users (fetched via `GET /users`). "Invite Staff" button opens a form (email, role dropdown). Pending invitations section showing sent invitations with resend/cancel actions.

- [ ] **Step 2: Create clients page**

`apps/web/src/pages/clients.tsx` — table of clients. "Add Client" button opens a form (name, industry, contact email). Click a client row to edit inline or in a modal.

- [ ] **Step 3: Create firm settings page**

`apps/web/src/pages/firm-settings.tsx` — form pre-populated with current firm data (name, timezone, billing email). Save button calls `PATCH /firms/current`.

- [ ] **Step 4: Create accept invitation page**

`apps/web/src/pages/accept-invitation.tsx` — reads token from URL query param. Calls `GET /invitations/validate/{token}` to show invitation details (email, role). Form to set display name and password. Submit calls `POST /invitations/accept`. On success, stores tokens and redirects to dashboard.

- [ ] **Step 5: Add routes to App.tsx**

Add routes:
- `/settings` → `FirmSettingsPage`
- `/users` → `UsersPage`
- `/clients` → `ClientsPage`
- `/accept-invitation` → `AcceptInvitationPage` (public, no ProtectedRoute)

All protected routes wrapped in `Layout`.

- [ ] **Step 6: Write RTL tests for pages with logic — TDD applies**

- `src/pages/users.test.tsx`:
  - Renders the user table with fetched users (mock `api` response)
  - Opening the invite form and submitting calls `POST /users/invite` with email + role
  - Pending invitation resend/cancel buttons call the right endpoints

- `src/pages/clients.test.tsx`:
  - Renders the client table
  - Creating a client calls `POST /clients` with the form payload
  - Editing a client calls `PATCH /clients/:id`

- `src/pages/firm-settings.test.tsx`:
  - Pre-populates the form from `GET /firms/current`
  - Save calls `PATCH /firms/current` with the updated fields
  - Shows a success indicator on save

- `src/pages/accept-invitation.test.tsx`:
  - Reads token from query param, calls validate endpoint, displays email + role
  - Submitting password calls accept endpoint, stores tokens, navigates to dashboard
  - Expired/invalid token surfaces an error state

Use `vi.mock` to stub the `api` module; assert the mock was called with the expected URL and body.

- [ ] **Step 7: Run tests — expect pass**

```bash
cd apps/web
npm test -- --run
```

Expected: all tests PASS.

- [ ] **Step 8: Manual full-flow browser walkthrough**

1. Register firm → dashboard
2. Navigate to Settings → update firm name → save → verify changed
3. Navigate to Users → invite staff@test.com as Staff
4. Check console/Mailhog for magic link URL
5. Open magic link URL → accept invitation page → set password → dashboard
6. Navigate to Clients → create "TechCorp" → see in list
7. Log out and log back in as the invited staff user → verify limited access

- [ ] **Step 9: Commit**

```bash
git add apps/web/
git commit -m "feat: add users, clients, firm settings, and invitation acceptance pages (incl. unit tests)"
```

---

### Task 17: RLS isolation verification test

**Files:**
- Create: `apps/axiom-api/internal/identity/rls_test.go`

- [ ] **Step 1: Write RLS isolation test**

```go
package identity_test

import (
	"context"
	"testing"

	"github.com/axiom-platform/axiom/apps/axiom-api/internal/identity"
	"github.com/axiom-platform/axiom/apps/axiom-api/internal/platform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRLSIsolation(t *testing.T) {
	pool := platform.TestDB(t)
	privKey, pubKey := testKeyPair(t)
	issuer := identity.NewJWTIssuer(privKey, pubKey)
	svc := identity.NewService(pool, issuer)
	ctx := context.Background()

	// Register two firms
	firm1, err := svc.RegisterFirm(ctx, identity.RegisterFirmInput{
		FirmName:   "Firm Alpha",
		AdminEmail: "alpha@test.com",
		AdminName:  "Alpha Admin",
		Password:   "pass123",
		Country:    "US",
		StaffCount: "1-10",
		AuditTypes: []string{"SOC2"},
	})
	require.NoError(t, err)

	firm2, err := svc.RegisterFirm(ctx, identity.RegisterFirmInput{
		FirmName:   "Firm Beta",
		AdminEmail: "beta@test.com",
		AdminName:  "Beta Admin",
		Password:   "pass123",
		Country:    "US",
		StaffCount: "1-10",
		AuditTypes: []string{"SOC2"},
	})
	require.NoError(t, err)

	// Create a client in Firm Alpha
	_, err = svc.CreateClient(ctx, firm1.Firm.ID, identity.CreateClientInput{
		Name:     "Alpha Client",
		Industry: "Tech",
	})
	require.NoError(t, err)

	// List clients as Firm Alpha — should see 1
	alphaClients, err := svc.ListClients(ctx, firm1.Firm.ID, 50, 0)
	require.NoError(t, err)
	assert.Len(t, alphaClients, 1)

	// List clients as Firm Beta — should see 0
	betaClients, err := svc.ListClients(ctx, firm2.Firm.ID, 50, 0)
	require.NoError(t, err)
	assert.Len(t, betaClients, 0)
}
```

- [ ] **Step 2: Run test**

```bash
cd apps/axiom-api
go test ./internal/identity/ -v -run TestRLSIsolation -count=1
```

Expected: PASS — Firm Beta cannot see Firm Alpha's client.

- [ ] **Step 3: Commit**

```bash
git add apps/axiom-api/internal/identity/rls_test.go
git commit -m "test: add RLS multi-tenant isolation verification test"
```

---

### Task 18: GitHub Actions CI pipeline

**Files:**
- Create: `.github/workflows/ci.yml`
- Create: `.github/dependabot.yml`
- Create: `.golangci.yml`
- Create: `.gitleaks.toml` (optional — defaults are fine for most repos)

This task sets up the CI pipeline that runs on every pull request from a feature branch to `master`. It validates build, tests, lint, and security on the full codebase produced by Tasks 1–17.

- [ ] **Step 1: Create golangci-lint config**

`.golangci.yml`:

```yaml
run:
  timeout: 5m
  tests: true

linters:
  enable:
    - errcheck
    - govet
    - staticcheck
    - unused
    - gosimple
    - ineffassign
    - gofmt
    - goimports
    - misspell
    - bodyclose
    - rowserrcheck
    - sqlclosecheck

linters-settings:
  goimports:
    local-prefixes: github.com/axiom-platform/axiom
```

- [ ] **Step 2: Create Dependabot config**

`.github/dependabot.yml`:

```yaml
version: 2
updates:
  - package-ecosystem: "gomod"
    directory: "/apps/axiom-api"
    schedule:
      interval: "weekly"
    open-pull-requests-limit: 5

  - package-ecosystem: "npm"
    directory: "/"
    schedule:
      interval: "weekly"
    open-pull-requests-limit: 5

  - package-ecosystem: "npm"
    directory: "/apps/web"
    schedule:
      interval: "weekly"

  - package-ecosystem: "npm"
    directory: "/packages/openapi"
    schedule:
      interval: "weekly"

  - package-ecosystem: "github-actions"
    directory: "/"
    schedule:
      interval: "weekly"

  - package-ecosystem: "docker"
    directory: "/"
    schedule:
      interval: "weekly"
```

- [ ] **Step 3: Create CI workflow**

`.github/workflows/ci.yml`:

```yaml
name: CI

on:
  pull_request:
    branches: [master]
  push:
    branches: [master]

concurrency:
  group: ci-${{ github.ref }}
  cancel-in-progress: true

permissions:
  contents: read
  security-events: write  # required for CodeQL upload
  pull-requests: read

jobs:
  go-build-test:
    name: Go build & test
    runs-on: ubuntu-latest
    services:
      postgres:
        image: pgvector/pgvector:pg17
        env:
          POSTGRES_DB: axiom_db
          POSTGRES_USER: axiom_svc
          POSTGRES_PASSWORD: localdev
        ports:
          - 5432:5432
        options: >-
          --health-cmd "pg_isready -U axiom_svc -d axiom_db"
          --health-interval 5s
          --health-timeout 3s
          --health-retries 10
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version-file: apps/axiom-api/go.mod
          cache-dependency-path: apps/axiom-api/go.sum
      - name: Install migrate
        run: go install -tags postgres github.com/golang-migrate/migrate/v4/cmd/migrate@latest
      - name: Run migrations
        run: migrate -path apps/axiom-api/migrations -database "postgres://axiom_svc:localdev@localhost:5432/axiom_db?sslmode=disable" up
      - name: Build
        working-directory: apps/axiom-api
        run: go build ./...
      - name: Test (with race detector)
        working-directory: apps/axiom-api
        run: go test ./... -race -count=1 -coverprofile=coverage.out
      - name: Upload coverage
        uses: actions/upload-artifact@v4
        with:
          name: go-coverage
          path: apps/axiom-api/coverage.out

  go-lint:
    name: Go lint
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version-file: apps/axiom-api/go.mod
          cache-dependency-path: apps/axiom-api/go.sum
      - uses: golangci/golangci-lint-action@v6
        with:
          version: latest
          working-directory: apps/axiom-api
      - name: gofmt check
        working-directory: apps/axiom-api
        run: |
          if [ -n "$(gofmt -l .)" ]; then
            echo "gofmt found unformatted files:"
            gofmt -l .
            exit 1
          fi

  go-vuln:
    name: Go vulnerability scan
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version-file: apps/axiom-api/go.mod
          cache-dependency-path: apps/axiom-api/go.sum
      - name: Install govulncheck
        run: go install golang.org/x/vuln/cmd/govulncheck@latest
      - name: Run govulncheck
        working-directory: apps/axiom-api
        run: govulncheck ./...

  web-build-test:
    name: Web build, test, lint
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: "22"
          cache: npm
      - run: npm ci
      - name: Codegen (verify specs match generated code)
        run: |
          npm run codegen
          if ! git diff --exit-code apps/web/src/api/generated/; then
            echo "::error::Generated TypeScript types are out of sync with OpenAPI specs."
            echo "Run 'npm run codegen' and commit the result."
            exit 1
          fi
      - name: Type check
        working-directory: apps/web
        run: npx tsc --noEmit
      - name: Lint
        working-directory: apps/web
        run: npm run lint
      - name: Test
        working-directory: apps/web
        run: npm test -- --run
      - name: Build
        run: npm run build
      - name: npm audit
        run: npm audit --audit-level=high
        continue-on-error: false

  codeql:
    name: CodeQL
    runs-on: ubuntu-latest
    permissions:
      security-events: write
    strategy:
      fail-fast: false
      matrix:
        language: [go, javascript-typescript]
    steps:
      - uses: actions/checkout@v4
      - uses: github/codeql-action/init@v3
        with:
          languages: ${{ matrix.language }}
      - uses: github/codeql-action/autobuild@v3
      - uses: github/codeql-action/analyze@v3
        with:
          category: "/language:${{ matrix.language }}"

  gitleaks:
    name: Secret scan
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
      - uses: gitleaks/gitleaks-action@v2
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}

  trivy:
    name: Trivy filesystem scan
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: aquasecurity/trivy-action@master
        with:
          scan-type: fs
          scan-ref: .
          severity: HIGH,CRITICAL
          exit-code: "1"
          ignore-unfixed: true

  workflow-lint:
    name: actionlint
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: raven-actions/actionlint@v2
```

- [ ] **Step 4: Add test scripts to apps/web/package.json**

Ensure `apps/web/package.json` has these scripts (add any missing):

```json
{
  "scripts": {
    "dev": "vite",
    "build": "tsc -b && vite build",
    "lint": "eslint .",
    "test": "vitest",
    "preview": "vite preview"
  }
}
```

Install Vitest if not already present:

```bash
cd apps/web
npm install -D vitest @vitest/ui jsdom @testing-library/react @testing-library/jest-dom
```

- [ ] **Step 5: Configure branch protection (manual — user action)**

Document that the user should configure branch protection on `master` in GitHub settings:
- Require a pull request before merging
- Require status checks to pass before merging, selecting: `Go build & test`, `Go lint`, `Go vulnerability scan`, `Web build, test, lint`, `CodeQL`, `Secret scan`
- Require branches to be up to date before merging

`Trivy filesystem scan` and `actionlint` are recommended but can be left non-blocking while the codebase stabilizes.

- [ ] **Step 6: Verify locally before push**

Run each check locally against the current branch to confirm the pipeline will pass:

```bash
# Go checks
cd apps/axiom-api
go build ./...
go test ./... -race -count=1
gofmt -l .  # should output nothing

# Web checks
cd ../..
npm ci
npm run codegen
npx tsc --noEmit --project apps/web
cd apps/web && npm run lint && npm test -- --run && cd ../..
npm run build
```

Expected: all commands succeed with no errors.

- [ ] **Step 7: Commit**

```bash
git add .github/ .golangci.yml apps/web/package.json apps/web/package-lock.json
git commit -m "ci: add GitHub Actions pipeline with build, test, lint, and security scans"
```

- [ ] **Step 8: Verify CI runs on PR**

When the user opens the PR from `phase-0-1-scaffold-and-identity` → `master`, confirm:
- All required jobs trigger and complete within ~10 minutes
- Any failures are surfaced inline on the PR
- CodeQL results appear in the Security tab
- Dependabot opens its first update PRs within a day or two

---

## Self-Review Checklist

**Spec coverage:**
- [x] Docker Compose with Postgres + pgvector — Task 1
- [x] Go project scaffold with Chi — Task 2
- [x] React scaffold with Vite, TanStack Query — Task 3
- [x] Turborepo + OpenAPI codegen — Task 4
- [x] Platform package (config, DB pool, errors) — Task 5
- [x] Identity enums and tables with RLS — Task 6
- [x] sqlc queries for all identity tables — Task 7
- [x] JWT issuer with RSA signing — Task 8
- [x] Gateway middleware (auth, role checks) — Task 9
- [x] Firm registration and login — Tasks 10, 11
- [x] User, client, invitation CRUD — Task 13
- [x] Server wiring with CORS — Task 12
- [x] React auth (context, API client, token refresh) — Task 14
- [x] React pages (login, register, dashboard, settings, users, clients) — Tasks 15, 16
- [x] RLS isolation verification — Task 17
- [x] Invitation acceptance flow — Tasks 13, 16
- [x] GitHub Actions CI pipeline (build, test, lint, vuln scan, CodeQL, secret scan) — Task 18

**Not in scope (deferred to Phase 2+):**
- Methodology templates (Phase 2)
- Firm control objectives (Phase 2)
- SSO/OAuth (Phase 10 — Enterprise feature)
- Email verification (deferred — magic link for invitations covers the invite flow)
- Tour completion endpoint (trivial — add when dashboard tour is built)

**Type consistency:** Verified — `TokenPair`, `Claims`, `JWTIssuer`, `Service`, `Handler`, `Middleware`, `AppError` types are used consistently across all tasks.

---

*End of Phase 0 + Phase 1 Implementation Plan*
