# Axiom Implementation Plan

**Status:** Draft
**Approach:** Journey-ordered vertical slices (Approach A)
**Developer context:** Solo engineer + AI coding agent, beginner-comfortable with Go and React
**Key constraint:** Delay AWS spending as long as possible — everything runs locally until Phase 10

**Git workflow:** Each phase is built on its own branch (`feature/phase-N-description`). All commits go to that branch and are pushed to remote. The user creates a PR to master, reviews, and merges. At the start of the next phase, return to master, pull, verify the merge, then create a new branch.

---

## Cross-Cutting Methodology

Two methodologies apply to every phase and are defined before the phase-by-phase plan:

- **Test-Driven Development** — every unit of behavior is written test-first. See "Cross-Cutting Methodology: Test-Driven Development" below.
- **Per-Phase Manual Testing Instructions** — every phase ends by producing a `docs/superpowers/testing/<plan-filename>.md` walkthrough. See "Cross-Cutting Methodology: Per-Phase Manual Testing Instructions" below.

## Table of Contents

1. [Local-First Architecture](#1-local-first-architecture)
2. [Phase Overview](#2-phase-overview)
3. [Phase 0 — Dev Environment & Project Scaffold](#3-phase-0--dev-environment--project-scaffold)
4. [Phase 1 — Platform Core + Identity](#4-phase-1--platform-core--identity)
5. [Phase 2 — Frameworks, Templates & Engagement Creation](#5-phase-2--frameworks-templates--engagement-creation)
6. [Phase 3 — Evidence & Document Requests](#6-phase-3--evidence--document-requests)
7. [Phase 4 — Workpapers & Review Workflow](#7-phase-4--workpapers--review-workflow)
8. [Phase 5 — Trial Balance](#8-phase-5--trial-balance)
9. [Phase 6 — Reporting & Engagement Lifecycle](#9-phase-6--reporting--engagement-lifecycle)
10. [Phase 7 — AI Features](#10-phase-7--ai-features)
11. [Phase 8 — Real-Time Collaboration & Notifications](#11-phase-8--real-time-collaboration--notifications)
12. [Phase 9 — EQR & Compliance Hardening](#12-phase-9--eqr--compliance-hardening)
13. [Phase 10 — AWS Infrastructure & Deployment](#13-phase-10--aws-infrastructure--deployment)

---

## 1. Local-First Architecture

The entire application runs on a developer laptop until Phase 10. No AWS account needed for Phases 0–9.

| AWS Service | Local Replacement | Swap Mechanism |
|---|---|---|
| RDS PostgreSQL | Docker container (postgres 17 + pgvector) | Same connection string, same SQL |
| S3 | Local filesystem directory (`./local-storage/`) | Go interface (`StorageBackend`) — swap implementation at deploy |
| Bedrock (AI) | Anthropic API directly (~$5–10 total dev cost) | Go interface (`AIClient`) — same Claude models, different transport |
| SES (email) | Console logging + Mailhog (local SMTP viewer) | Go interface (`EmailSender`) — swap implementation at deploy |
| Step Functions | Plain Go code (same guard logic as functions) | Extract to Step Functions ASL when deploying |
| CloudFront | Vite dev server proxies to Go API | N/A — dev only |
| River (jobs) | Same — River runs on Postgres | No change needed |

**Cost of Phases 0–9: $0/month** plus ~$5–10 one-time Anthropic API credit when AI features are built in Phase 7.

---

## Cross-Cutting Methodology: Test-Driven Development

**Every phase is built test-first.** This is not optional — it applies to every task in every phase from Phase 0 onward. The CI pipeline (set up in Phase 0) blocks any PR whose tests don't pass, which makes TDD the path of least resistance.

### The loop

For every unit of behavior (a service method, a handler, a React hook, a SQL query, a React component with logic):

1. **Red** — write the failing test first. The test names the behavior in terms of inputs and expected outputs, not implementation details. Run it and confirm it fails for the right reason (not a compile error masking the actual assertion).
2. **Green** — write the minimum implementation that makes the test pass. Resist adding "while I'm here" code.
3. **Refactor** — with the test locking behavior in place, clean up names, extract helpers, remove duplication. Re-run the test after each change.

Skipping the red step is the most common failure mode — writing a test *after* the code already works means the test has never actually demonstrated it can fail. If you catch yourself doing this, delete the code, re-run the test to see it fail, then re-implement.

### What counts as a test

| Layer | Test type | Example |
|---|---|---|
| Pure functions, domain logic | Go unit tests (`_test.go`, no DB) | JWT signing, slug generation, state machine guards, Levenshtein ratios |
| Service methods (DB-touching) | Integration tests with a real Postgres test database | `RegisterFirm`, `CreateEngagement`, `AdvanceToFieldwork` |
| HTTP handlers | `httptest.NewRecorder` tests covering status codes, body shape, auth failures | `POST /api/v1/auth/register`, `GET /api/v1/engagements/:id` |
| River workers | Worker unit tests with fake AI/storage/email clients | `ai-completeness-check`, `document-extract`, `notification-deliver` |
| RLS policies | Multi-tenant isolation tests (two firms, assert zero leakage) | One per new RLS-guarded table |
| React hooks & logic | Vitest + React Testing Library | `useAuth`, form validation, AI modification-ratio display |
| React components | RTL render + interaction assertions | Login form submits, engagement wizard advances steps |
| End-to-end journey | Playwright against running dev stack (from Phase 2 onward) | "Staff completes review → Partner signs off" |

### What does not need a test

- Auto-generated code (sqlc output, openapi-typescript output, oapi-codegen interfaces) — tested transitively through the service/handler tests that consume it.
- Trivial struct literals, constants, and config parsing with no branching.
- Third-party library behavior (Chi router, pgx, React Router) — test the integration, not the dependency.

Use judgment: if a mistake in the code would be caught by a compiler, a linter, or an integration test one layer up, a dedicated test adds noise without signal.

### Coverage expectations

- **Business logic packages** (`identity`, `auditcore`, `trialbalance`, `workpaper`, `reporting`, `ai`, state machine guards, AI content tracking): target **≥85% line coverage**. Every public method has at least one happy-path and one error-path test.
- **Gateway, platform infrastructure**: target **≥70%**. Middleware gets direct tests (auth, role checks, RLS isolation).
- **UI components**: target **≥60%**. Focus tests on components with logic; trivial layout components can be covered by journey/E2E tests.
- **No coverage threshold on generated code, main.go wiring, or migration SQL files.**

Coverage thresholds are enforced by CI after Phase 1 — the baseline is captured from the Phase 1 codebase, and no subsequent PR may regress it. The CI workflow defined in Phase 0 uploads `coverage.out` as an artifact; the threshold gate is wired in once real business logic exists.

### How this shows up in each phase

Each phase description below lists a "Testable Outcome." These are the acceptance criteria at the phase level — **they must be exercised by automated tests**, not only by manual browser walkthroughs. When a phase says "Staff submits workpaper → Manager reviews → Partner signs off," that entire flow has a corresponding integration test that drives it through the service layer, plus the relevant unit tests for each guard and sign-off rule it exercises.

The Phase 9 compliance validation explicitly requires automated integration tests walking the full SOC 2 and PCAOB lifecycles. That is the final backstop: by the end of Phase 9, every regulatory guard, immutability rule, and sign-off hierarchy is covered by an automated test.

---

## Cross-Cutting Methodology: Per-Phase Manual Testing Instructions

**Every phase plan ends with a task that writes manual testing instructions to `docs/superpowers/testing/<phase-plan-filename>.md`.** Automated tests verify behavior; this document verifies the *experience* — a human running the stack end-to-end in a browser or shell, confirming the phase's Testable Outcome is actually shippable.

### Why it exists

- **Testable Outcome rows** (see §2 Phase Overview) are acceptance criteria at the phase level. They need to be exercised, not just asserted.
- **CI proves the code compiles and tests pass.** It does not prove that a new user can sign up, that the buttons are labeled, or that an invitation link actually lands somewhere useful.
- **Reviewers and the user need a script.** A consistent format per phase means anyone can pick up the branch, run through the document, and either sign off or file a regression.
- **Gaps become visible.** Writing the walkthrough often exposes rough edges (missing copy, silent failures, unhandled empty states) that automated tests wouldn't catch.

### Location and naming

- Path: `docs/superpowers/testing/<same-filename-as-the-plan>.md`.
- Example: plan `docs/superpowers/plans/phase-0-and-1-scaffold-and-identity.md` → testing doc `docs/superpowers/testing/phase-0-and-1-scaffold-and-identity.md`.
- One file per plan. Phases with sub-phases that share a plan share a testing doc.

### Required structure

Every testing doc must have these sections, in this order:

1. **Header** — phase scope, link to the plan file, one-sentence purpose.
2. **Start the stack** — exact commands to bring up every service the phase depends on (Docker services, backend, frontend, any phase-specific tools). Include expected log lines so the reader knows what "ready" looks like. Include a smoke curl or equivalent.
3. **Browser walkthrough — the happy path** — numbered or sub-headed flows matching the Testable Outcome rows for the phase. Each flow states what to do and what to expect visually (copy, state transitions, styling landmarks tied to `.impeccable.md` tokens where relevant).
4. **Targeted edge cases** — a table of `Case | Expected`. Cover at minimum: wrong credentials / bad input, role-gated endpoints attempted by the wrong role, stale-session recovery, backend-down resilience on the frontend.
5. **Data-layer spot checks** — SQL or CLI commands that verify invariants the UI can't show directly: RLS isolation (register two firms, confirm zero leakage), audit-log append-only guarantees, immutability rules, AI decision records, etc. Link to the corresponding automated test so the manual check is a backup, not the primary guarantee.
6. **Integration placeholders** — services that are wired in this phase but whose full behavior ships later (e.g., Mailhog when email isn't wired yet). Explicitly call out what's expected to be empty/no-op.
7. **Known gaps** — bulleted list of things that are *deliberately* absent in this phase. Prevents false-positive regression reports.
8. **Reporting regressions** — one-line note that mismatches with §§2–5 are real regressions and should be fixed (with a new automated test) before merging.

Sections 2–5 are mandatory. Sections 1, 6, 7, 8 are mandatory but can be short. Additional sections may be added as the phase requires (e.g., AWS deployment smoke tests in Phase 10, AI feature walkthroughs in Phase 7).

### Plan-level requirements

Every phase plan must include, as its penultimate or final task (after the impeccable validation pass where applicable):

> ### Task N: Manual testing instructions
>
> Write `docs/superpowers/testing/<plan-filename>.md` covering the Testable Outcome for this phase. Follow the structure defined in "Cross-Cutting Methodology: Per-Phase Manual Testing Instructions" in `docs/superpowers/specs/implementation-plan-design.md`. Commit under a `docs:` prefix.

The testing doc is produced **after** the phase's features work end-to-end. Drafting it sooner risks the doc describing intent instead of reality.

### Relationship to automated tests

- **Automated tests remain the contract.** If a manual walkthrough step fails, the fix is (a) fix the code, and (b) add an automated test that would have caught it. The testing doc is a *human-readable map* of what's been built, not a replacement for CI.
- **Don't duplicate.** If a flow is already exhaustively covered by integration or RTL tests, the manual step can be one line ("Register → Login → Dashboard loads") rather than a blow-by-blow.
- **Keep it current.** When a phase changes a flow covered by an earlier phase's testing doc, amend that earlier doc in the same PR. Stale testing docs mislead reviewers.

### Example

The Phase 0/1 testing document at `docs/superpowers/testing/phase-0-and-1-scaffold-and-identity.md` is the reference implementation of this structure.

---

## 2. Phase Overview

| Phase | What You Build | Testable Outcome | Journeys |
|---|---|---|---|
| 0 | Dev environment + project scaffold | "Hello world" in Go and React, connected to Postgres | — |
| 1 | Platform core + Identity | Register firm → login → invite staff → staff logs in | 1, 2 |
| 2 | Frameworks, templates, engagements | Create engagement from template → scope → advance to Fieldwork | 3 |
| 3 | Evidence & document requests + PDF service | Send doc requests → client uploads → auditor reviews → link evidence | 7, 8 |
| 4 | Workpapers & review workflow | Write workpaper → submit → review notes → resolve → sign off | 5, 6 |
| 5 | Trial balance | Import CSV → map accounts → adjustments → analytics | 4 |
| 6 | Reporting & engagement lifecycle | Generate report → issue → finalize → archive | 9 |
| 7 | AI features (all 8) | Every AI feature works end-to-end with real Claude models | All |
| 8 | Real-time collaboration & notifications | Co-edit workpapers, receive in-app + email notifications | All |
| 9 | EQR & compliance hardening | Full regulatory workflow, immutable audit trail, all gates | 10 |
| 10 | AWS infrastructure & deployment | Application running on AWS, accessible via browser | — |

AI features are a separate phase rather than woven into each module. Rationale: getting the CRUD workflows right first gives a stable foundation. The database schema includes `ai_decision_id` columns and `ai_content_metadata` fields from the start — so the data model is AI-ready from Phase 2 onward. For a solo developer learning Go, keeping AI complexity out of the initial builds reduces cognitive load.

---

## 3. Phase 0 — Dev Environment & Project Scaffold

**Goal:** Go from nothing installed to a running Go API and React dev server, connected to a local Postgres database.

### Tools to Install

- **Go** (latest stable, 1.24+)
- **Node.js** (LTS 22.x) + npm
- **Docker Desktop** (for Postgres and later the Python PDF service)
- **Turborepo** (`npm install turbo --global`, manages monorepo builds)

### Docker Compose

A single `docker-compose.yml` at the repo root:

- **Postgres 17** with pgvector extension — port `5432`, database `axiom_db`, user `axiom_svc`
- **Mailhog** — SMTP on port `1025`, web UI on port `8025` (for viewing emails in later phases)

### Monorepo Scaffold

```
apps/
  axiom-api/
    cmd/server/main.go          — entrypoint: wire dependencies, start Chi router
    internal/
      platform/                 — DB pool, config, error types, slog setup
      gateway/                  — (placeholder)
      identity/                 — (placeholder)
      auditcore/                — (placeholder)
      trialbalance/             — (placeholder)
      workpaper/                — (placeholder)
      reporting/                — (placeholder)
      ai/                       — (placeholder)
    migrations/                 — (empty, first migration in Phase 1)
    go.mod
  doc-processing/               — (placeholder, built in Phase 3)
  web/
    src/
      main.tsx
      App.tsx
      routes/
    package.json
    vite.config.ts
    tsconfig.json
packages/
  openapi/                      — existing OpenAPI specs (already in repo)
docker-compose.yml
turbo.json
```

### What Gets Built

**Go API:** Chi router with `GET /healthz` returning `200 OK`. Connects to Postgres on startup, logs "connected to database" via slog.

**React app:** Vite + TypeScript + React. Single page that fetches `/healthz` from the Go API to prove connectivity. React Router and TanStack Query configured.

**OpenAPI codegen pipeline:**
- `oapi-codegen` generates Go server interfaces from `packages/openapi/*.yaml`
- `openapi-typescript` generates TypeScript types from the same specs
- Both wired into Turborepo — spec change regenerates on next build

### Testable Outcome

Run `docker compose up -d`, `go run ./cmd/server`, and `npm run dev` — browser shows "Axiom" with a green "API connected" indicator.

### Continuous Integration

GitHub Actions workflow (`.github/workflows/ci.yml`) runs on every pull request targeting `master`. All jobs must pass before merge. Jobs run in parallel where possible; Turborepo `--filter` scopes work to affected packages.

**Build & compile:**
- `go build ./...` across the Go module (fails on any compile error)
- `tsc --noEmit` for the React app (strict TypeScript compile check)
- `npm run build` via Turborepo (full production build of `apps/web`)
- OpenAPI codegen check: regenerate Go + TS clients from `packages/openapi/*.yaml`, fail if working tree dirty (enforces committed specs match generated code)

**Unit tests (TDD-enforced — see Cross-Cutting Methodology):**
- `go test ./... -race -count=1` with the race detector enabled
- `npm test` (Vitest) for the React app
- Coverage report uploaded as a job artifact. Hard thresholds activate in Phase 1: ≥85% on business-logic packages, ≥70% on platform/gateway, ≥60% on UI components. PRs that regress a package's coverage are blocked.
- Every task that adds a service method, handler, worker, RLS policy, or component-with-logic lands with its tests in the same commit or an earlier commit in the same PR — never in a follow-up

**Linting & formatting:**
- `golangci-lint run` (config at `.golangci.yml` — starts with the default linter set: `errcheck`, `govet`, `staticcheck`, `unused`, `gosimple`, `ineffassign`)
- `gofmt -l .` (fails if any file is not formatted)
- `eslint` on `apps/web/src`
- `prettier --check` on the web workspace
- `actionlint` on `.github/workflows/*.yml` (catches workflow syntax errors)
- `hadolint` on all `Dockerfile`s

**Security & vulnerability scanning** (all free for public/private repos on GitHub):
- **`govulncheck ./...`** — Go's official vulnerability scanner, cross-references the Go vuln DB against imported packages and actual call paths
- **`npm audit --audit-level=high`** — flags high/critical npm advisories
- **CodeQL** — GitHub's static analysis for Go and JavaScript/TypeScript; scheduled weekly plus on PR, results surface in the Security tab
- **Dependabot** — enabled via `.github/dependabot.yml` for Go modules, npm, GitHub Actions, and Docker base images; opens PRs for updates
- **`gitleaks`** — scans the diff for accidentally committed secrets (API keys, JWT secrets, AWS credentials)
- **Trivy** — filesystem scan for known CVEs in dependencies and (later) container images

**Status checks required for merge:** build, unit-tests, lint, govulncheck, codeql, gitleaks. Dependabot and Trivy scheduled scans report findings but don't block PRs.

---

## 4. Phase 1 — Platform Core + Identity

**Goal:** A user can register a firm, log in, see a dashboard, invite staff, and the invited staff can accept and log in.

**Journeys covered:** 1 (Firm Setup — stages 1–2, 5), 2 (Staff Onboarding)

### Backend

**Platform package** (`internal/platform`):
- Database connection pool (pgx/v5) with `SET app.current_firm_id` helper
- Config struct loaded from environment variables (envconfig)
- Structured logging (slog with JSON output)
- Standard error types (NotFound, Unauthorized, Validation, Conflict)
- Migration runner (golang-migrate — auto-runs pending migrations on startup in dev mode)

**Migrations:**
- Enum types: `user_role`, `auth_method`, `notification_frequency`, `invitation_status`
- `firms` table with RLS policy (`USING (id = current_setting('app.current_firm_id')::uuid)`)
- `users` table with RLS policy, check constraints (firm XOR client)
- `clients` table with RLS policy
- `invitations` table with RLS policy
- `SET app.current_firm_id` helper function

**Identity module** (`internal/identity`):
- Firm registration: create firm + first FirmAdmin user in one transaction
- Login: email/password (bcrypt), returns JWT access + refresh token pair
- JWT: RSA-signed, payload contains `user_id`, `firm_id`, `role`; access token 15-min TTL, refresh token 7-day TTL
- JWT refresh endpoint
- User CRUD: list firm users, update profile, deactivate (soft delete)
- Invitation: create → log magic link to console (email in Phase 8) → accept (creates user with assigned role)
- Client CRUD: list, create, update (basic — needed by Phase 2 for engagements)

**Gateway middleware** (`internal/gateway`):
- `AuthMiddleware` — verify JWT signature, inject user_id/firm_id/role into Go context
- `WithFirmIsolation` — set `app.current_firm_id` on the database connection for RLS
- `WithRole(roles...)` — check user's role against allowed roles, return 403 if disallowed
- Rate limiter: in-memory token bucket per firm_id (100 req/min default)

**API endpoints** (generated interfaces from OpenAPI, hand-written handlers):
- `POST /api/v1/auth/register` — firm + admin user creation
- `POST /api/v1/auth/login` — returns JWT pair
- `POST /api/v1/auth/refresh` — refresh access token
- `GET /api/v1/firms/me` — current firm profile
- `PATCH /api/v1/firms/me` — update firm profile
- `GET /api/v1/users` — list firm users
- `POST /api/v1/users/invite` — create invitation
- `POST /api/v1/auth/accept-invitation` — accept magic link, create user
- `GET /api/v1/clients` — list clients
- `POST /api/v1/clients` — create client
- `PATCH /api/v1/clients/:id` — update client

### Frontend

- **Auth pages:** Sign up (firm registration form with intake fields), Login, Accept invitation (set password)
- **Layout shell:** Sidebar navigation (Engagements, Clients, Users, Settings), top bar with user display name and logout
- **Dashboard:** Onboarding checklist (complete firm profile, invite staff, create first engagement) with links to each step. Empty states for engagement list.
- **Firm settings:** Edit name, timezone, billing email
- **User management:** Table of firm users, invite button, pending invitations list
- **Client management:** Table of clients, create client modal
- **Auth context:** React context holding JWT, auto-refresh before expiry, redirect to `/login` on 401
- **API client:** Generated from OpenAPI via openapi-typescript, all requests include `Authorization: Bearer <token>` header

### Testable Outcome

1. Register "Acme CPAs" with admin user → redirected to dashboard with onboarding checklist
2. Create client "TechCorp" → visible in client list
3. Invite staff@acme.com → magic link logged to console → open link → set password → staff lands on dashboard
4. Admin sees staff in user list with correct role
5. Staff sees empty engagement list
6. RLS: register a second firm "Beta LLP" → confirm zero data leakage between firms

---

## 5. Phase 2 — Frameworks, Templates & Engagement Creation

**Goal:** Seed regulatory framework data, let firms activate methodology templates, and create fully scoped engagements with controls and test procedures auto-scaffolded from templates.

**Journeys covered:** 1 (Firm Setup — stages 3–4), 3 (Engagement Scoping)

### Backend

**Migrations — system reference tables (no RLS, shared across all firms):**
- `frameworks` — versioned framework definitions
- `framework_requirements` — individual criteria within frameworks
- `control_objective_library` — system-maintained semantic control objectives
- `control_objective_library_mappings` — library-to-framework-requirement links

**Migrations — firm methodology (RLS via firm_id):**
- `methodology_templates`
- `template_controls`
- `template_test_procedures`
- `template_document_requests`
- `firm_control_objectives`
- `firm_control_objective_mappings`

**Migrations — engagement (RLS via firm_id):**
- `engagements` (with status enum, lifecycle fields)
- `engagement_team_members`
- `engagement_frameworks`
- `client_acceptances`
- `controls`
- `test_procedures`

**Migrations — cross-cutting (created now so all subsequent phases can reference them):**
- `ai_decisions` (with `ai_context_type`, `ai_review_action` enums) — table exists from Phase 2 but is not populated until Phase 7
- `audit_log` (with `actor_type` enum, PostgreSQL RULEs preventing UPDATE/DELETE) — used from Phase 2 onward
- `notifications` (with `notification_type`, `delivery_channel` enums) — table exists from Phase 2 but delivery logic is built in Phase 8

**Seed data scripts** (run via migrations or a seed command):
- SOC 2 TSC 2017 framework with all Trust Services Criteria (CC1–CC9, A1, PI1, C1, P1)
- ISO 27001:2022 framework with Annex A controls (A.5–A.8)
- HIPAA Security Rule framework with relevant sections
- GAAS (AU-C) framework with key assertions
- Pre-built methodology templates:
  - "SOC 2 Type II Standard" — ~50 controls, ~80 test procedures, ~80 document request templates
  - "GAAS Financial Audit Standard" — ~30 controls, ~60 test procedures
- Control objective library entries with cross-framework mappings

**Identity module additions:**
- Methodology template CRUD: list (including system-provided), activate/deactivate, view contents
- Firm control objective CRUD: create (from library or custom), edit, view framework mappings
- Firm control objective mapping management: add/remove framework requirement links

**Audit Core module** (`internal/auditcore`) — first implementation:
- Engagement CRUD:
  - Create from methodology template (auto-scaffolds controls + test procedures from template)
  - List engagements (filterable by status, client)
  - Get engagement detail (with team, controls, frameworks)
  - Update engagement (name, period dates)
- Engagement team member management: assign, remove, list
- Engagement framework management: add secondary frameworks to engagement
- Client acceptance workflow:
  - Create/update quality risks and firm responses
  - Independence confirmation (Partner only)
  - Sign-off (Partner only) — sets `accepted_at`, unblocks Planning → Fieldwork
- Control CRUD: list per engagement, update status, assign auditor
- Test procedure CRUD: list per control, update status/result
- Engagement state machine — first two transitions:
  - Planning → Fieldwork: guard checks `ClientAcceptance.accepted_at IS NOT NULL`
  - Reverse: Fieldwork → Planning (Partner only, for scope changes)

**Audit log** — first entries (written to the `audit_log` table created above):
- `engagement.created`, `engagement.status.changed`
- `client_acceptance.signed_off`
- `control.status.changed`

### Frontend

- **Methodology templates:** Browse available templates (system + firm-custom), view template contents (controls, test procedures, document requests), activate/deactivate templates for the firm
- **Firm control objectives:** List objectives with their framework mappings, create new (from library or custom), edit mappings
- **Engagement list:** Filterable table with status badges, engagement type, client name, date range
- **Create engagement wizard:**
  1. Select client (or create new)
  2. Select engagement type + primary framework
  3. Select methodology template
  4. Review auto-scaffolded controls → confirm
  5. Assign engagement partner and team members
- **Engagement detail page:**
  - Overview tab: status, client, framework, period, team
  - Controls tab: table of all controls with status, assigned auditor, key control flag
  - Team tab: manage team members with engagement roles
  - Client Acceptance tab: quality risk form, independence confirmation, partner sign-off action
- **Control detail page:** Description, test procedures list, status, assigned auditor
- **Phase transition actions:** "Advance to Fieldwork" button with guard condition checklist (shows what's blocking if client acceptance isn't signed)

### Testable Outcome

1. Admin activates "SOC 2 Type II Standard" methodology template
2. Partner creates engagement for TechCorp using the template → 50 controls and 80 test procedures auto-created
3. Partner assigns Manager and two Staff to the engagement team
4. Partner fills out client acceptance: quality risks, independence confirmation, signs off
5. "Advance to Fieldwork" button unblocks → engagement status changes to Fieldwork
6. Staff opens engagement → sees assigned controls with test procedures
7. Attempt to advance without client acceptance → blocked with explanation

---

## 6. Phase 3 — Evidence & Document Requests

**Goal:** Auditors can send document requests to clients, clients can upload documents via a no-login tokenized portal, and auditors can review and link evidence to test procedures.

**Journeys covered:** 7 (Document Requests — auditor side), 8 (Document Requests — client side)

### Backend

**Migrations:**
- `evidence_items` (with extraction_status enum)
- `evidence_links`
- `document_requests` (with status enum)
- `client_hub_tokens`
- `delegation_tokens`

**Storage interface** (`internal/platform`):
```go
type StorageBackend interface {
    Upload(ctx context.Context, key string, data io.Reader, contentType string) error
    Download(ctx context.Context, key string) (io.ReadCloser, error)
    Delete(ctx context.Context, key string) error
}
```
Local filesystem implementation writes to `./local-storage/evidence/`. S3 implementation swapped in Phase 10.

**Evidence management** (in `internal/auditcore`):
- File upload: validate file type/size → store via StorageBackend → create EvidenceItem record
- Evidence CRUD: list per engagement (via client_id), list per firm+client, get with extracted text
- Evidence link: create link between EvidenceItem and TestProcedure, with notes

**Document request management** (in `internal/auditcore`):
- CRUD: create (optionally from template), list per engagement, update, track status
- Status lifecycle: Pending → Submitted → InReview → Accepted/Rejected → (back to Pending if rejected)
- Batch create from methodology template's TemplateDocumentRequest records
- Accept action: atomically creates EvidenceLink when control/procedure are set on the request

**Client Hub** (in `internal/auditcore`):
- Token generation: create ClientHubToken with 90-day expiry, scoped to engagement
- Token validation endpoint: no JWT required, token grants read access to document requests for the engagement
- Client upload: via token, creates EvidenceItem + sets DocumentRequest.status to Submitted
- Delegation: ClientAdmin creates DelegationToken for a single request → delegate can upload for that request only

**Python doc-processing service** (`apps/doc-processing/`):
- FastAPI + uvicorn
- `POST /extract` — accepts file bytes, returns extracted text + metadata
- pdfplumber for digital PDFs, pytesseract for scanned documents
- Docker container definition (python:3.13-slim + tesseract-ocr apt package)
- Added to docker-compose.yml

**River integration** (`internal/platform`):
- River client initialization (uses existing Postgres connection)
- River worker startup (embedded in the main binary)
- First worker: `auditcore.document-extract`
  - Triggered on evidence upload
  - Calls doc-processing service at `http://localhost:8000/extract`
  - Stores extracted text in `EvidenceItem.extracted_text`
  - Updates `extraction_status` to Complete or Failed
  - Retry with exponential backoff on failure

**API endpoints:**
- `POST /api/v1/engagements/:id/document-requests` — create request
- `POST /api/v1/engagements/:id/document-requests/batch` — batch create from template
- `GET /api/v1/engagements/:id/document-requests` — list requests
- `PATCH /api/v1/document-requests/:id` — update request
- `POST /api/v1/document-requests/:id/accept` — accept with evidence link
- `POST /api/v1/document-requests/:id/reject` — reject with reason
- `POST /api/v1/engagements/:id/evidence` — upload evidence file
- `GET /api/v1/engagements/:id/evidence` — list evidence items
- `POST /api/v1/evidence/:id/link` — link evidence to test procedure
- `GET /api/v1/test-procedures/:id/evidence` — list evidence for a procedure
- `POST /api/v1/engagements/:id/client-hub-token` — generate client hub token
- **Client Hub (no JWT, token-authenticated):**
  - `GET /api/v1/client-hub/:token/requests` — list document requests
  - `POST /api/v1/client-hub/:token/requests/:id/upload` — upload document
  - `POST /api/v1/client-hub/:token/requests/:id/delegate` — create delegation

### Frontend

- **Document request management page:**
  - Table of requests with status, due date, assigned client contact
  - Create request form (title, instructions, due date, link to control/procedure)
  - Batch create from template action
  - Status workflow actions (accept, reject, with reason)
- **Evidence browser:**
  - Table of evidence items for the engagement
  - Filter by extraction status, source type, client
  - Preview extracted text
  - Upload button (drag-and-drop supported)
- **Evidence linking UI:**
  - On test procedure detail: "Link Evidence" panel showing available evidence items
  - Click to link with optional notes
  - View linked evidence per procedure
- **Client Hub portal** (separate route `/client-hub/:token`, no auth required):
  - Engagement name and firm branding
  - Table of document requests (pending, submitted, accepted, rejected)
  - Per-request: view instructions, upload files, see status
  - Delegation action for ClientAdmin (enter colleague's email → generates link)
- **Generate client hub link** button on engagement detail page (copies tokenized URL)

### Testable Outcome

1. Staff batch-creates document requests from engagement template → 80 requests created
2. Staff generates client hub token → URL copied
3. Open client hub URL in incognito (no login) → see pending requests
4. Client uploads PDF for a request → status changes to Submitted
5. Doc-processing extracts text (visible in River job logs) → extraction_status = Complete
6. Auditor opens request → reviews uploaded document → accepts → EvidenceLink created
7. Evidence visible on the test procedure's evidence panel
8. Rejected request returns to client with reason → client re-uploads

---

## 7. Phase 4 — Workpapers & Review Workflow

**Goal:** Auditors can create and edit rich-text workpapers, submit them for review, managers can add review notes, and the sign-off hierarchy is enforced.

**Journeys covered:** 5 (Control Testing — workpaper creation), 6 (Workpaper Review)

### Backend

**Migrations:**
- `workpapers` (with workpaper_type and workpaper_status enums)
- `workpaper_versions` (with ai_content_metadata jsonb — populated in Phase 7)
- `review_notes` (with severity and status enums, DELETE rule for immutability)

**Workpaper module** (`internal/workpaper`):
- Workpaper CRUD: create per engagement (optionally linked to a control), update content, list per engagement
- Content storage: jsonb (ProseMirror document format) — stored in `workpapers.content` for current state
- Version creation: on each save, insert a new `workpaper_versions` row with content snapshot and version_number
- Status lifecycle enforcement:
  - Draft → PreparedPendingReview: Staff submits. Checks all content is non-empty. (AI content gate added in Phase 7.)
  - PreparedPendingReview → InReview: Manager opens for review.
  - InReview → ReviewNotesOpen: Manager creates a review note. Auto-transitions.
  - ReviewNotesOpen → InReview: Staff addresses all notes and resubmits.
  - InReview → ReviewComplete: Manager marks review complete. Guard: all review notes resolved.
  - ReviewComplete → SignedOff: Partner signs off. Creates audit log entry.
- Sign-off hierarchy enforcement:
  - Submit for review: only `prepared_by_id` (Staff or above)
  - Review actions: only Manager or Partner (not the preparer)
  - Sign-off: only Partner
- Locking: after engagement finalization, `is_locked = true` on all workpapers

**Review notes** (in `internal/workpaper`):
- Create: reviewer creates note with severity (Question, Suggestion, RequiredChange) and content anchor
- Respond: staff writes response text
- Resolve: reviewer marks as resolved
- Immutability: PostgreSQL RULE prevents DELETE on review_notes
- Audit log entries on create, respond, resolve

**API endpoints:**
- `GET /api/v1/engagements/:id/workpapers` — list workpapers
- `POST /api/v1/engagements/:id/workpapers` — create workpaper
- `GET /api/v1/workpapers/:id` — get workpaper with current content
- `PUT /api/v1/workpapers/:id/content` — save content (creates new version)
- `POST /api/v1/workpapers/:id/submit-for-review` — advance to PreparedPendingReview
- `POST /api/v1/workpapers/:id/start-review` — advance to InReview
- `POST /api/v1/workpapers/:id/complete-review` — advance to ReviewComplete
- `POST /api/v1/workpapers/:id/sign-off` — advance to SignedOff
- `POST /api/v1/workpapers/:id/return-to-staff` — return to ReviewNotesOpen
- `GET /api/v1/workpapers/:id/versions` — version history
- `GET /api/v1/workpapers/:id/versions/:version` — specific version content
- `GET /api/v1/workpapers/:id/review-notes` — list review notes
- `POST /api/v1/workpapers/:id/review-notes` — create review note
- `POST /api/v1/review-notes/:id/respond` — staff responds
- `POST /api/v1/review-notes/:id/resolve` — reviewer resolves

### Frontend

- **Workpaper list** per engagement: table with title, type, status badge, preparer, reviewer, signer
- **Workpaper editor:**
  - TipTap (ProseMirror-based) rich text editor
  - Toolbar: headings (H1–H3), bold, italic, bullet/numbered lists, tables, horizontal rules
  - Section-based document structure (each H2 defines a section boundary — important for AI content tracking in Phase 7)
  - Auto-save on idle (debounced, creates version)
  - Manual save button
  - Status bar showing current status and available actions
- **Submit for review** action button (with confirmation)
- **Review interface** (Manager/Partner view):
  - Read workpaper content with review note anchors visible
  - Review note creation: select text → add note with severity dropdown and description
  - Review note list sidebar: grouped by status (Open, Responded, Resolved)
  - "Complete Review" and "Return to Staff" actions
- **Staff response interface:**
  - See review notes inline in the editor
  - Response text field per note
  - "Resubmit" action when all notes are responded to
- **Version history:** Timeline of saves, click to view any prior version (read-only)
- **Sign-off action** (Partner only, with confirmation dialog)

### Testable Outcome

1. Staff creates workpaper "Access Control Testing" linked to a control
2. Staff writes content in TipTap editor → auto-saves → version 1 created
3. Staff submits for review → status: PreparedPendingReview
4. Manager opens → starts review → creates 2 review notes (1 Question, 1 RequiredChange)
5. Status auto-transitions to ReviewNotesOpen → Staff sees notes inline
6. Staff responds to both notes → resubmits
7. Manager resolves both notes → marks review complete
8. Partner signs off → workpaper status: SignedOff → audit log entry created
9. Attempt to delete a review note → blocked by PostgreSQL RULE
10. Attempt for Staff to sign off → 403 Forbidden

---

## 8. Phase 5 — Trial Balance

**Goal:** Auditors can import trial balances from CSV/Excel, map accounts to financial statement line items, manage adjustments, and run population analysis queries.

**Journeys covered:** 4 (Trial Balance)

### Backend

**Migrations:**
- Enum types: `account_type`, `mapping_status`, `adjustment_type`
- `trial_balances`
- `trial_balance_accounts` (with generated `net_balance` column)
- `trial_balance_adjustments` (with check constraint: waived requires reason)
- `column_mapping_profiles`

**Trial Balance module** (`internal/trialbalance`):
- **CSV/Excel import:**
  - Accept uploaded file (CSV or XLSX)
  - Parse with configurable column mapping (which columns are account_number, account_name, debit, credit)
  - Validate: total debits = total credits (flag non-zero difference)
  - Create TrialBalance + TrialBalanceAccount records in one transaction
  - Save column mapping as a ColumnMappingProfile for reuse
- **Column mapping profile CRUD:** list per firm, create, load
- **Account mapping:**
  - Manual: set `mapped_fs_line_item` on individual accounts, status → Confirmed
  - Bulk confirm: select multiple accounts, set the same FS line item
  - Override: change a previously set mapping, status → Overridden
  - (AI mapping in Phase 7 sets status to AISuggested)
- **Adjustment management:**
  - Propose adjustment: Staff creates with amount, description, type (Proposed)
  - Approve: Manager/Partner approves → type changes to Passed
  - Waive: with mandatory reason → type changes to Waived
  - Cumulative waived tracking against materiality threshold
- **Population analysis** (SQL queries, results returned as JSON):
  - Gap testing: find missing account numbers in a sequence
  - Duplicate detection: accounts with identical names or numbers
  - Threshold filtering: accounts above/below a specified amount
  - Benford's law: distribution of leading digits vs. expected
  - Period-over-period variance: compare current vs. prior year balances (requires prior engagement with TB)
- **Lead schedule generation:** group accounts by `mapped_fs_line_item`, sum balances, include adjustments

**API endpoints:**
- `POST /api/v1/engagements/:id/trial-balance/import` — upload + import
- `GET /api/v1/engagements/:id/trial-balance` — get TB with accounts
- `PATCH /api/v1/trial-balance-accounts/:id` — update mapping
- `POST /api/v1/trial-balance-accounts/bulk-map` — bulk update mappings
- `GET /api/v1/firms/me/column-mapping-profiles` — list profiles
- `POST /api/v1/firms/me/column-mapping-profiles` — save profile
- `POST /api/v1/trial-balances/:id/adjustments` — propose adjustment
- `POST /api/v1/adjustments/:id/approve` — approve
- `POST /api/v1/adjustments/:id/waive` — waive with reason
- `GET /api/v1/trial-balances/:id/analytics/gap-test` — gap testing
- `GET /api/v1/trial-balances/:id/analytics/duplicates` — duplicate detection
- `GET /api/v1/trial-balances/:id/analytics/benfords` — Benford's law
- `GET /api/v1/trial-balances/:id/analytics/variance` — period-over-period
- `GET /api/v1/trial-balances/:id/lead-schedules` — grouped by FS line item

### Frontend

- **Trial balance import wizard:**
  1. Upload CSV/XLSX file
  2. Column mapping UI: preview first 5 rows, assign columns to fields (account_number, account_name, debit, credit) — or load a saved profile
  3. Validation summary (total debits, total credits, difference)
  4. Confirm import
- **Spreadsheet view** (AG Grid Community + HyperFormula):
  - All accounts with balances in a spreadsheet-like grid
  - Columns: account number, name, type, debit, credit, net balance, mapped FS line item, mapping status
  - Inline editing for mapped_fs_line_item (dropdown with standard FS line items)
  - Color coding by mapping status: red (Unmapped), yellow (AISuggested — Phase 7), green (Confirmed), blue (Overridden)
  - Sort and filter on any column
  - Row selection for bulk mapping
  - Cell-level comments (for notes on specific accounts)
  - Balance totals in footer row
- **Adjustment panel:**
  - List of proposed/passed/waived adjustments
  - Create adjustment form (select account, amount, description)
  - Approve/waive actions with permission checks
  - Cumulative waived amount vs. materiality indicator
- **Analytics dashboard:**
  - Tab-based layout: Gap Test, Duplicates, Benford's, Variance
  - Results displayed as tables with flags
  - Benford's law: bar chart comparing actual vs. expected distribution
  - Variance: table with accounts sorted by absolute variance percentage
- **Lead schedules:** Grouped view of accounts by FS line item with subtotals

### Testable Outcome

1. Staff uploads CSV trial balance (200 accounts) → column mapping → import succeeds
2. Validation: debits = credits confirmed (or difference flagged)
3. Staff maps accounts to FS line items in the AG Grid (inline dropdown)
4. Staff proposes adjustment: "Depreciation under-recorded, $5,000" → Manager approves
5. Run Benford's law analysis → chart shows distribution, flags anomalous digits
6. Run gap test → finds missing sequence numbers
7. View lead schedule → accounts grouped by FS line item with adjustment-adjusted balances
8. Save column mapping profile → reuse on next import

---

## 9. Phase 6 — Reporting & Engagement Lifecycle

**Goal:** Generate audit reports, implement the full engagement lifecycle state machine with all guards, and handle finalization and archival.

**Journeys covered:** 9 (Reporting & Archive)

### Backend

**Migrations:**
- Enum types: `report_type`, `report_status`
- `reports` (with ai_content_metadata jsonb — populated in Phase 7)
- `report_versions` (with ai_content_metadata jsonb — populated in Phase 7)

**Reporting module** (`internal/reporting`):
- Report creation: select report type, linked to engagement
- Report content: jsonb (same ProseMirror format as workpapers), section-based
- Report generation via River job (`reporting.report-generate`):
  - Reads engagement data from auditcore (controls, test results, exceptions, evidence stats)
  - Reads trial balance data from trialbalance (if financial audit)
  - Reads workpaper summaries from workpaper module
  - Renders report sections using Go html/template
  - Stores rendered content in report record
- Report status lifecycle:
  - Draft → ClientReview (share with client for factual review)
  - ClientReview → FirmReview (client review complete)
  - FirmReview → Issued (Partner signs off — triggers finalization cascade)
  - Issued → Archived (system, after assembly deadline)
- Report issuance triggers:
  - Compute `assembly_deadline` on engagement (report_issued_at + 60 days AICPA, + 45 days PCAOB)
  - Compute `retention_deadline` on engagement (+ 5 years AICPA/SOC/ISO, + 7 years PCAOB, + 6 years HIPAA)
  - Set `engagement.report_issued_at`
- Version history on each save

**Full engagement state machine** (in `internal/auditcore`):

All transitions implemented with guards:

| From → To | Guard | Who |
|---|---|---|
| Planning → Fieldwork | `ClientAcceptance.accepted_at` is set | Partner |
| Fieldwork → Review | All Controls have status Complete or Exception | Manager or Partner |
| Review → Reporting | All ReviewNotes resolved + EQR signed off (if applicable) | Partner |
| Reporting → Finalized | `Report.status = Issued` | Partner |
| Finalized → Archived | `assembly_deadline` elapsed | System (cron job via River) |
| Fieldwork → Planning | Scope change | Partner |
| Review → Fieldwork | Additional procedures needed | Manager or Partner |
| Reporting → Review | Significant issue found | Partner |

**Finalization cascade:**
- All workpapers: set `is_locked = true`
- All controls: conclusions become immutable
- Report: transitions to read-only
- Audit log: `engagement.finalized` entry

**Archive simulation** (full S3 Object Lock in Phase 10):
- Mark engagement as Archived
- Record `archived_at` timestamp
- All data becomes read-only (enforced at application layer)

**Addendum workflow** (post-finalization):
- Partner can create addendum workpaper versions with mandatory reason
- `WorkpaperVersion.is_addendum = true`, `addendum_reason` required
- Original content preserved unchanged — addendum is a new version
- Audit log: `workpaper.addendum.created`

**River job:** `reporting.report-generate` — async report assembly and rendering

**API endpoints:**
- `POST /api/v1/engagements/:id/reports` — create report
- `GET /api/v1/engagements/:id/reports` — list reports
- `GET /api/v1/reports/:id` — get report with content
- `PUT /api/v1/reports/:id/content` — save report content (creates version)
- `POST /api/v1/reports/:id/generate` — trigger async report generation
- `POST /api/v1/reports/:id/submit-client-review` — advance to ClientReview
- `POST /api/v1/reports/:id/submit-firm-review` — advance to FirmReview
- `POST /api/v1/reports/:id/issue` — issue report (triggers finalization)
- `POST /api/v1/engagements/:id/advance` — advance engagement phase (with guard check)
- `POST /api/v1/engagements/:id/reverse` — reverse engagement phase (Partner, with reason)
- `POST /api/v1/workpapers/:id/addendum` — create addendum (post-finalization)
- `GET /api/v1/engagements/:id/lifecycle` — get current state + guard conditions status

### Frontend

- **Report editor:**
  - Same TipTap editor as workpapers, section-based
  - Report type indicator (SOC 2 Type II, Financial Audit Opinion, etc.)
  - "Generate Report" button → triggers River job → content populated when complete
  - Status workflow actions (submit for client review, firm review, issue)
- **Engagement lifecycle dashboard:**
  - Visual state machine diagram (current phase highlighted, completed phases checked)
  - Guard condition checklist per transition (green check = satisfied, red X = blocking)
  - "Advance" button enabled only when all guards pass
  - Reverse transition option (Partner only, requires reason)
- **Finalization flow:**
  - Confirmation dialog: "This will lock all workpapers and evidence links. This cannot be undone."
  - Display computed assembly_deadline and retention_deadline
- **Archived engagement view:**
  - All content read-only with "Archived" banner
  - Addendum creation button (Partner only)
  - Version history accessible
- **Report issuance confirmation:**
  - Summary: X controls tested, Y exceptions, Z workpapers signed off
  - Display deadlines that will be computed

### Testable Outcome

1. Partner creates a SOC 2 Type II report → triggers report generation → content populated
2. Partner edits report → submits for client review → firm review → issues
3. Issuance computes assembly_deadline (60 days) and retention_deadline (5 years)
4. Engagement advances: Fieldwork → Review → Reporting → Finalized
5. Finalized: all workpapers locked, content read-only
6. Test each guard: try to advance without meeting conditions → blocked with explanation
7. Test reverse path: Partner reverses Review → Fieldwork with reason
8. Create addendum on finalized workpaper → new version with `is_addendum = true`
9. Simulate archive → engagement becomes fully read-only

---

## 10. Phase 7 — AI Features

**Goal:** Implement all eight AI features using real Claude models (via Anthropic API directly), integrate AIDecision tracking, and add AI content tracking to workpapers and reports.

**Journeys covered:** All (AI touches every journey)

### Backend

**Migrations:**
- `evidence_embeddings` (pgvector)
- `framework_requirement_embeddings` (pgvector)
- `control_objective_library_embeddings` (pgvector)
- `firm_control_objective_embeddings` (pgvector)

Note: `ai_decisions` table was created in Phase 2. This phase adds the embedding tables and populates `ai_decisions` via AI features.

**AI module** (`internal/ai`):
```go
type AIClient interface {
    Complete(ctx context.Context, req CompletionRequest) (CompletionResponse, error)
    Embed(ctx context.Context, text string) ([]float32, error)
}
```
- Anthropic API implementation (used locally) — uses the official Anthropic Go SDK
- Bedrock implementation (swapped in Phase 10) — same interface, different transport
- Prompt template system: each feature has a structured prompt template with input slots
- Token counting and cost tracking per request

**AIDecision service** (in `internal/auditcore`):
- Create AIDecision record (context_type, model_id, raw_output, suggested_value, confidence)
- List per engagement, filterable by context_type and review_action
- Review action: Accept / Modify / Reject (sets accepted_value, reviewed_by_id, reviewed_at)
- Audit log entry on each review action

**Embedding pipeline:**
- Generate embeddings for all framework requirements (one-time seed)
- Generate embeddings for control objective library entries (one-time seed)
- Generate embeddings for firm control objectives (on create/update)
- Generate embeddings for evidence items (on extraction complete)
- pgvector similarity queries for retrieval

**Feature 1 — Document Completeness Review:**
- River worker: `auditcore.ai-completeness-check`
- Trigger: after document extraction completes (chains from `document-extract` worker)
- Input: DocumentRequest details + extracted text + control/procedure context + framework requirements
- Model: Claude Sonnet
- Output: AIDecision with completeness assessment, confidence score, gap explanation
- Human review: Accept / Modify / Reject actions on the document request review screen

**Feature 2 — Control Mapping:**
- River worker: `auditcore.ai-batch-control-mapping`
- Trigger: new engagement created (enqueued during engagement scaffolding)
- Input: FirmControlObjective descriptions + all FrameworkRequirements for engaged frameworks + library entries as few-shot context
- Model: Claude Haiku
- Process: embed objective → retrieve similar library entries → score against requirements → threshold at 0.75
- Output: AIDecision per proposed mapping
- Human review: bulk confirm/reject in engagement scoping UI

**Feature 3 — Trial Balance Account Mapping:**
- River worker: `trialbalance.ai-account-mapping`
- Trigger: trial balance import complete
- Input: account number, name, balances + prior year mappings (if rollforward)
- Model: Claude Haiku (few-shot classification)
- Output: sets `mapping_status = AISuggested` + `ai_decision_id` on each account
- Human review: confirm/override in AG Grid

**Feature 4 — Workpaper Narrative Draft:**
- River worker: `workpaper.ai-workpaper-draft`
- Trigger: auditor clicks "Generate AI Draft" on a workpaper
- Input: control description, test procedure, linked evidence text, exceptions, prior year workpaper, firm template
- Model: Claude Sonnet
- Output: draft content inserted into workpaper with AI content tracking metadata
- Human review: mandatory editing before submit for review (gate logic below)

**Feature 5 — Evidence Link Suggestion:**
- River worker: `auditcore.ai-evidence-link-suggestion`
- Trigger 1: auditor opens test procedure evidence linking panel
- Trigger 2: document accepted via completeness review
- Input: test procedure description + evidence pool embeddings + prior year links
- Model: Claude Haiku
- Output: ranked evidence suggestions with confidence and explanation
- Human review: accept/reject per suggestion

**Feature 6 — Risk Category Suggestion:**
- River worker: `auditcore.ai-risk-category-suggestion`
- Trigger: partner opens client acceptance form
- Input: client industry, engagement type, prior year acceptance + exceptions
- Model: Claude Sonnet
- Output: suggested risk categories in sidebar (not pre-populated in form)
- Human review: partner selects from suggestions or writes manually

**Feature 7 — Trial Balance Anomaly Detection:**
- River worker: `trialbalance.ai-anomaly-detection`
- Trigger: nightly (River periodic job) + once after initial import
- Input: current year balances, prior year balances, account mappings, materiality thresholds
- Model: Claude Haiku
- Output: anomaly flags per account (type, magnitude, description)
- Tier: Tier 1 for nonissuer (informational), Tier 2 for PCAOB (creates AIDecision)
- Human review: PCAOB only — auditor acknowledges each flag

**Feature 8 — Report Section Draft:**
- River worker: `reporting.ai-report-section-draft`
- Trigger: partner clicks "Generate AI Draft" on a report section
- Input: report type, engagement-wide data, prior year report, firm template
- Model: Claude Sonnet
- Output: draft section with AI content tracking
- Human review: mandatory editing before report issuance

**AI Content Tracking:**
- `ai_content_metadata` jsonb on WorkpaperVersion and ReportVersion:
  ```json
  {
    "sections": [{
      "section_id": "scope-and-approach",
      "ai_generated": true,
      "ai_generated_at": "...",
      "human_edited": false,
      "human_edited_by": null,
      "human_edited_at": null,
      "ai_character_count": 1450,
      "current_character_count": 1450,
      "modification_ratio": 0.0
    }],
    "summary": {
      "total_sections": 8,
      "ai_generated_sections": 3,
      "ai_sections_edited": 0,
      "ai_sections_unedited": 3,
      "overall_modification_ratio": 0.0
    }
  }
  ```
- Modification ratio: Levenshtein distance / AI character count, computed on save
- Advancement gates:
  - PreparedPendingReview (workpapers): all AI sections must have `human_edited = true`
  - Report issuance: all AI sections must have `human_edited = true`
  - Soft gate: sections with `modification_ratio < 0.05` trigger confirmable warning

### Frontend

- **Document request review:** AI completeness assessment panel with confidence score, gap explanation, Accept/Modify/Reject buttons
- **Engagement scoping:** control mapping suggestions table — proposed framework links with confidence, bulk confirm/reject
- **Trial balance AG Grid:** account mapping status column now shows AISuggested (yellow) with one-click confirm, anomaly flag indicators on rows
- **Workpaper editor:**
  - "Generate AI Draft" button
  - Per-section AI indicators (badge showing "AI-generated" with modification ratio)
  - Warning on submit: "X sections have minimal edits to AI content"
  - Confirmation gate for low-modification sections
- **Evidence linking panel:** AI-suggested links with confidence bars and explanation text
- **Client acceptance form:** risk category suggestions in a sidebar panel
- **Report editor:** "Generate AI Draft" per section, same AI tracking indicators as workpapers
- **AI Decision queue:** engagement-level view of all pending AI decisions, filterable by type
- **AI audit trail:** table showing all AI decisions for an engagement (context, model, action taken, by whom)

### Testable Outcome

1. Upload document → AI completeness review runs → auditor sees assessment → accepts
2. Create engagement → AI proposes control-to-framework mappings → partner confirms in bulk
3. Import trial balance → AI maps accounts → staff confirms high-confidence, reviews low-confidence
4. Request AI workpaper draft → draft appears with AI section indicators → edit sections → modification ratio updates → submit succeeds
5. Open evidence linking → AI suggests top-3 evidence items → accept suggestions
6. Open client acceptance → AI risk categories appear in sidebar → partner selects relevant ones
7. Nightly anomaly detection flags 3 accounts → staff sees flags in AG Grid
8. Request AI report section → draft inserted → must edit before issuing
9. AI Decision queue shows all pending decisions → review and clear
10. For a PCAOB engagement: anomaly flags create AIDecision records requiring acknowledgment

---

## 11. Phase 8 — Real-Time Collaboration & Notifications

**Goal:** Multiple users can co-edit workpapers in real-time, and platform events trigger in-app and email notifications.

**Journeys covered:** All (notifications touch every journey; collaboration supports Journeys 5, 6)

### Backend

**Real-time collaboration** (`internal/workpaper`):
- WebSocket endpoint: `GET /api/v1/workpapers/:id/ws` — upgrades to WebSocket
- Yjs document provider:
  - Server holds authoritative Yjs document state in memory per open workpaper
  - Syncs updates between connected clients via WebSocket
  - Yjs awareness protocol: shows cursor positions and user identities
  - Persistence: saves Yjs document state to database on idle timeout (5 seconds) and on last client disconnect
  - On first connection: loads document state from database
- Connection management:
  - JWT required on WebSocket upgrade (passed as query param)
  - Firm isolation: validate user has access to the engagement
  - Workpaper locking: reject write connections for finalized engagements (read-only sync allowed)
  - Max connections per workpaper: configurable (default 10)

**Notifications** (in `internal/auditcore`):
- Notification service:
  - Create notification with type, recipient, title, body, deep_link
  - Delivery logic: check recipient's `notification_frequency` preference
    - RealTime: deliver immediately (in-app + email)
    - Daily: batch for daily digest
    - Weekly: batch for weekly digest
  - Mark read/unread
  - List unread count
- River worker: `auditcore.notification-deliver`
  - Sends email via EmailSender interface (Mailhog locally, SES in Phase 10)
  - HTML email templates per notification type
- Notification triggers wired into existing workflows:
  - `EngagementAssignment` — user added to engagement team
  - `ReviewNoteAdded` — review note created on a workpaper
  - `ReviewNoteResolved` — review note resolved
  - `DocumentRequestStatus` — document request status changes
  - `PhaseTransition` — engagement advances to next phase
  - `EQRNotification` — EQR assignment or finding
  - `ReminderEscalation` — document request overdue (3 reminders, then escalate to auditor)

**Document request reminder automation** (in `internal/auditcore`):
- River periodic job: check overdue document requests
- Reminder schedule: 7 days before due, on due date, 7 days after due
- After 3 reminders: escalation notification to the assigned auditor
- Updates `reminder_count` and `last_reminder_sent_at`

**API endpoints:**
- `GET /api/v1/workpapers/:id/ws` — WebSocket upgrade for Yjs collaboration
- `GET /api/v1/notifications` — list notifications (paginated, unread first)
- `GET /api/v1/notifications/unread-count` — count of unread
- `POST /api/v1/notifications/:id/read` — mark as read
- `POST /api/v1/notifications/read-all` — mark all as read

### Frontend

- **Workpaper editor — Yjs integration:**
  - TipTap configured with Yjs collaboration extension
  - Real-time cursors with user name labels and distinct colors
  - Presence indicators: avatars of connected users shown above the editor
  - Offline support: edits queue locally, sync when reconnected
  - Locked workpapers: cursor tracking works but editing is disabled
- **Notification center:**
  - Bell icon in the top bar with unread count badge
  - Dropdown panel: recent notifications (last 10) with timestamps
  - Click notification → navigate to deep_link (e.g., specific workpaper or document request)
  - "View all" link → full notification list page with filters (type, read/unread)
  - Mark read/unread actions
- **Email notifications:** viewable in Mailhog at `localhost:8025`

### Testable Outcome

1. Open same workpaper in two browser windows (different users) → see each other's cursors → edits sync in real-time
2. Third user joins → all three see each other, edits merge correctly
3. Close all windows → reopen → content persisted correctly
4. Manager creates review note → Staff receives in-app notification (bell badge increments) → click → navigates to workpaper
5. Document request goes overdue → reminder notification sent → visible in Mailhog
6. After 3 reminders → escalation notification to auditor
7. Engagement advances to Review → all team members notified
8. Notification preferences: change to Daily → notifications batch instead of immediate

---

## 12. Phase 9 — EQR & Compliance Hardening

**Goal:** Implement the Engagement Quality Review workflow and validate every regulatory compliance gate end-to-end.

**Journeys covered:** 10 (EQR)

### Backend

**Migrations:**
- Enum types: `eqr_status`, `eqr_conclusion`, `finding_severity`, `finding_status`
- `engagement_quality_reviews`
- `eqr_findings`

**EQR workflow** (in `internal/auditcore`):
- Assign EQR reviewer:
  - Validate: user has EQReviewer role
  - Validate: user is NOT in `engagement_team_members` for this engagement (independence check)
  - Create EngagementQualityReview record with status = Assigned
- Independence documentation: reviewer records `independence_documented_at`
- EQR findings:
  - Create finding (severity: Observation, Recommendation, RequiredAction)
  - Team response (engagement team member writes response)
  - Reviewer confirms response (changes finding status to Confirmed)
  - RequiredAction findings must be Addressed + Confirmed before sign-off
- EQR sign-off:
  - Guard: all RequiredAction findings are Confirmed
  - Sets conclusion (Satisfied, SatisfiedWithConcerns, NotSatisfied)
  - Sets `signed_off_at`
  - Unblocks Review → Reporting transition
- EQR reviewer access: read-only to entire engagement (all workpapers, evidence, controls, reports)

**Compliance hardening — audit log:**
- Verify PostgreSQL RULEs:
  - `audit_log`: no UPDATE, no DELETE
  - `review_notes`: no DELETE
- Verify all workflows create appropriate audit log entries:
  - Engagement status changes
  - Workpaper sign-offs (with timestamp — cannot be backdated)
  - Control conclusion changes
  - AI decision reviews
  - Client acceptance sign-off
  - EQR sign-off
  - Report issuance
  - Addendum creation
  - Delegation token creation

**Compliance hardening — engagement lifecycle:**
- Framework version lock: after Fieldwork begins, changing framework version requires Partner override with documented reason
- Addendum workflow: verify post-finalization modifications create proper addendum versions
- Assembly deadline enforcement: River periodic job checks for engagements past assembly_deadline → transitions to Archived

**Full compliance validation (automated integration tests):**
- Walk a SOC 2 Type II engagement through the entire lifecycle:
  1. Create firm, users, client
  2. Create engagement from template
  3. Complete client acceptance → advance to Fieldwork
  4. Upload evidence, create document requests, link evidence
  5. Create workpapers, review notes, sign off
  6. Assign EQR reviewer, create findings, resolve, sign off
  7. Generate report, issue → Finalized
  8. Assembly deadline passes → Archived
- Verify every guard: attempt to skip steps → confirm blocked
- Verify immutability: attempt to modify after finalization → confirm prevented
- Walk a PCAOB Financial Audit engagement:
  - Same flow + verify AI anomaly detection creates AIDecision records
  - Verify sampling documentation requirements on test procedures

**API endpoints:**
- `POST /api/v1/engagements/:id/eqr` — assign EQR reviewer
- `GET /api/v1/engagements/:id/eqr` — get EQR record with findings
- `POST /api/v1/eqr/:id/findings` — create finding
- `POST /api/v1/eqr-findings/:id/respond` — team response
- `POST /api/v1/eqr-findings/:id/confirm` — reviewer confirms
- `POST /api/v1/eqr/:id/sign-off` — EQR sign-off
- `GET /api/v1/engagements/:id/audit-trail` — paginated audit log for engagement

### Frontend

- **EQR assignment** (on engagement detail page):
  - Select reviewer dropdown (filtered to EQReviewer role users)
  - Independence validation feedback (shows error if user is on the team)
  - Assign button
- **EQR dashboard** (reviewer's view of the engagement):
  - Read-only access to all engagement content: controls, workpapers, evidence, reports
  - AI content substantiveness summary:
    - "42 workpapers used AI drafts. Average modification: 35%."
    - "3 workpapers have sections with <10% modification — review these first."
  - Filter to low-modification workpapers for priority review
- **EQR findings management:**
  - Create finding form (description, severity selector)
  - Findings list with status badges
  - Team response panel (engagement team writes response)
  - Confirm response action (reviewer)
  - Sign-off action with conclusion selector (Satisfied, SatisfiedWithConcerns, NotSatisfied)
- **Audit trail viewer:**
  - Filterable table: who did what, when, on what resource
  - Filter by: action type, actor, date range, resource type
  - Sign-off events highlighted
  - Export to CSV for external audit purposes
- **Addendum interface:**
  - On finalized workpapers: "Create Addendum" button (Partner only)
  - Reason field (mandatory)
  - New version created with `is_addendum = true` badge

### Testable Outcome

1. Assign EQR reviewer → independence validated (not on team)
2. Attempt to assign a team member as EQR → rejected with explanation
3. Reviewer sees entire engagement in read-only
4. Reviewer creates findings: 1 Observation, 1 RequiredAction
5. Team responds → reviewer confirms
6. Attempt EQR sign-off with unresolved RequiredAction → blocked
7. Resolve → sign off → Review → Reporting transition unblocked
8. Audit trail: all actions visible with timestamps and actors
9. Attempt to delete audit log entry via SQL → blocked by PostgreSQL RULE
10. Full SOC 2 lifecycle walkthrough (automated test): firm setup through archive, every guard verified

---

## 13. Phase 10 — AWS Infrastructure & Deployment

**Goal:** Provision real AWS infrastructure and deploy the application. This is the first phase that incurs cloud costs.

**Expected monthly cost at deployment:** ~$380–585/month (demo stage: dev + staging, lean config)

### Sub-Phase 10a: AWS Account Setup & Terraform Bootstrap

**What:**
- Create AWS Organization with initial accounts:
  - `axiom-management` — Organizations root, billing, SCPs
  - `axiom-tooling` — Terraform state, ECR repos, CI/CD roles
  - `axiom-dev` — development workloads
  - `axiom-staging` — pre-production (staging doubles as demo stage initially)
- Apply Service Control Policies (deny unused regions, protect CloudTrail, restrict production)
- Configure IAM Identity Center (SSO) for human console access
- Terraform bootstrap workspace:
  - S3 bucket for state (`axiom-terraform-state` in tooling account)
  - DynamoDB table for locks (`axiom-terraform-locks`)
  - OIDC identity provider for GitHub Actions
  - Cross-account IAM roles for deployment

**Testable outcome:** `terraform plan` runs from local machine against each account.

### Sub-Phase 10b: Network & Data Layer

**What:**
- **Network workspace:**
  - VPC per workload account (dev: `10.0.0.0/16`, staging: `10.1.0.0/16`)
  - Public + private subnets (dev: single AZ, staging: 2 AZs)
  - NAT Gateway (1 per environment at demo stage)
  - VPC endpoints: S3 (gateway, free), ECR, Secrets Manager, Bedrock runtime, CloudWatch Logs
  - Security groups: sg-alb, sg-ecs, sg-rds
- **Data workspace:**
  - RDS PostgreSQL 17 with pgvector (`db.t4g.medium`, single-AZ for dev/staging)
  - S3 buckets: evidence, archive, reports, spa
  - Secrets Manager: RDS credentials (axiom_svc + master), JWT keys, OAuth secrets
  - KMS keys: default (CloudWatch), HIPAA (S3 evidence), RDS
- Run database migrations against RDS

**Testable outcome:** Connect to RDS from local machine via SSM tunnel. Verify pgvector extension enabled.

### Sub-Phase 10c: Compute & DNS/CDN

**What:**
- **CICD workspace:**
  - ECR repositories: axiom-api, doc-processing, pgbouncer
  - GitHub Actions OIDC roles per environment
- Build and push Docker images to ECR
- **Compute workspace:**
  - ECS Fargate cluster
  - ALB with TLS (ACM certificate for `*.dev.axiom.com`)
  - ECS services:
    - `axiom-api` (Go binary + PgBouncer sidecar) — 1 task, 512 CPU / 1024 MB
    - `doc-processing` (Python + Tesseract) — 1 task, 512 CPU / 1024 MB
  - ECS Service Connect for internal DNS (`http://doc-processing:8000`)
  - Health check: `GET /healthz`
- **DNS/CDN workspace:**
  - Route 53 hosted zones
  - CloudFront distribution for React SPA (S3 origin with OAC)
  - ACM certificates (DNS validated)
  - SES: domain verification (DKIM + SPF + DMARC), sandbox mode for dev/staging

**Testable outcome:** Visit `https://app.dev.axiom.com` → React app loads. `https://api.dev.axiom.com/healthz` returns 200.

### Sub-Phase 10d: Service Integration Swaps

Swap local implementations for AWS services — each is a Go interface swap:

| Interface | Local Implementation | AWS Implementation |
|---|---|---|
| `StorageBackend` | Local filesystem | S3 (with SSE-KMS for HIPAA evidence) |
| `AIClient` | Anthropic API (direct) | Bedrock (via VPC endpoint, IAM auth) |
| `EmailSender` | Console log / Mailhog | SES (with bounce/complaint handling) |

**Additional integrations:**
- Step Functions: deploy `EngagementLifecycleStateMachine` and `DocumentRequestReminderStateMachine` ASL definitions. Wire Audit Core to invoke via AWS SDK.
- PgBouncer: configured as ECS sidecar, app connects to `localhost:6432`

**Testable outcome:** Upload document → stored in S3. AI features work via Bedrock. Email arrives via SES (in sandbox, to verified addresses).

### Sub-Phase 10e: CI/CD Pipeline

**GitHub Actions workflows:**

**On pull request:**
```
lint (golangci-lint, ruff, eslint)
  → build (Turborepo, affected services only)
  → unit test (go test, pytest)
  → terraform plan (all workspaces against dev, diff posted as PR comment)
```

**On merge to main:**
```
build → push images to ECR
  → terraform apply (dev)
  → ECS deploy (dev)
  → integration tests (dev)
  → terraform apply (staging)
  → ECS deploy (staging)
  → integration tests (staging)
```

**Database migrations:** run as one-shot ECS Fargate task after `data` workspace apply, before `compute` apply. Uses `master` credentials from Secrets Manager.

**Testable outcome:** Push to main → GitHub Actions deploys to dev → smoke tests pass → deploys to staging.

### Sub-Phase 10f: Observability

**What:**
- CloudWatch log groups for all services (7-day retention dev, 30-day staging)
- Basic CloudWatch alarms: ALB 5xx rate, RDS CPU, RDS storage
- SNS topic for alerts → email notification
- CloudWatch dashboard: `Axiom-{env}-Overview` (request rate, error rate, latency, RDS metrics)
- X-Ray tracing: 100% sampling on staging (not enabled on dev at demo stage)
- CloudTrail: Organization-level trail with S3 data events for evidence bucket

**Deferred (enabled when first customer is onboarded):**
- WAF (CloudFront + ALB WebACLs)
- GuardDuty
- AWS Config compliance rules
- Production account (`axiom-prod`)
- Full dashboard suite (AI, Data, Security)

**Testable outcome:** CloudWatch dashboard shows live metrics. Alarm fires on test 5xx. X-Ray traces visible in staging.

---

## Appendix: AWS Cost Timeline

| Milestone | What Changes | Monthly Cost |
|---|---|---|
| Phases 0–6 | Everything local | $0 |
| Phase 7 | Anthropic API key for AI testing | ~$5–10 one-time |
| Phases 8–9 | Still local | $0 |
| Phase 10 (demo stage) | 2 AWS accounts (dev + staging), lean config | $380–585 |
| First paying customer | Add prod account, Multi-AZ, security controls | $1,800–2,500 |
| Pre-SOC 2 audit | GuardDuty, AWS Config, conformance packs | $2,300–3,100 |
| Savings plans (3+ months stable) | Reserved instances for Fargate + RDS | $1,700–2,300 |

---

*End of Axiom Implementation Plan*
