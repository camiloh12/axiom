# Axiom Backend Architecture Design
**Date:** 2026-04-12
**Status:** Approved
**Supersedes:** Section 7 (Technology Stack) of `axiom-spec-v2-design.md`

---

## 1. Context and Decisions

### What Changed and Why

The v2 spec chose tRPC as the API layer because it provides end-to-end TypeScript type safety with zero code generation. This is correct when the frontend and backend are both TypeScript. It becomes a constraint when the backend language changes, because tRPC requires TypeScript on both sides.

The following decisions were made in sequence:

1. **REST over tRPC.** REST with OpenAPI decouples the frontend from the backend language. It also enables future public API exposure and partner integrations without adding a separate GraphQL or REST layer on top of tRPC.

2. **Go as the primary backend language.** Go was chosen for the following reasons:
   - Workiva, an established compliance SaaS platform, uses Go for their REST services and open-sourced a Go REST framework for exactly this use case.
   - Static typing with compile-time safety — well-suited for long-lived, regulated financial software.
   - Single compiled binary per service, lean container images (30–50MB vs 200MB+ for Python/Node.js), low memory footprint per Fargate task.
   - AWS SDK for Go (v2) is well-supported; Bedrock and Step Functions are first-party services, no additional vendor dependencies.
   - River (Go-native Postgres-backed job queue) provides a direct equivalent to pg-boss with no additional infrastructure.
   - Strong AI coding agent support — Go has a large and well-indexed training corpus.

3. **Python retained for PDF extraction.** Python's `pdfplumber` handles complex, multi-column, scanned audit documents (financial statements, SOC 2 system descriptions) better than any Go library. UniPDF (the closest Go alternative) is commercial and still weaker on layout analysis for the hardest documents. The Python service is stateless — one endpoint, one job — so the polyglot cost is contained.

4. **Microservices with bounded-context decomposition.** Services are split along genuine bounded contexts in the data model, not along the five product modules. The core engagement/control/evidence cluster is too tightly coupled to split without distributed transactions, so it stays in one service with a shared database. Independent domains (identity, trial balance, workpapers, reporting) get their own services and databases.

5. **Population analysis runs in PostgreSQL.** Trial balance population analysis (gap testing, duplicate detection, threshold filtering, Benford's law) runs as SQL queries against the trial balance database. The application layer handles orchestration and result formatting only. This removes performance as a differentiator between language options and makes Python's GIL a non-issue for this workload.

**AWS infrastructure** — account structure, VPC design, ECS/Fargate configuration, RDS instance sizing, S3 bucket policies, CI/CD pipeline, Terraform workspace segmentation, observability, and security controls are specified in [`infrastructure-design.md`](./infrastructure-design.md).

---

## 2. Service Decomposition

### Overview

```
┌──────────────────────────────────────────────────────────────┐
│                    React SPA (CloudFront)                     │
└───────────────────────────┬──────────────────────────────────┘
                            │ HTTPS
┌───────────────────────────▼──────────────────────────────────┐
│                  AWS ALB (TLS termination)                    │
└───────────────────────────┬──────────────────────────────────┘
                            │
┌───────────────────────────▼──────────────────────────────────┐
│                   API Gateway  (Go)                           │
│   JWT verification · request routing · rate limiting         │
│   Loads Identity public key at startup — no per-request      │
│   call to Identity on the hot path.                          │
└──┬──────────┬─────────────┬──────────┬──────────┬────────────┘
   │          │             │          │          │
   ▼          ▼             ▼          ▼          ▼
Identity  Audit Core   Trial Bal   Workpaper  Reporting
(Go)       (Go)          (Go)        (Go)       (Go)
                          ▲
                          │ HTTP (from River worker)
                    Doc Processing
                       (Python)
```

All services run as ECS Fargate Services. The API Gateway routes by URL prefix (e.g., `/api/v1/engagements/*` → Audit Core, `/api/v1/trial-balance/*` → Trial Balance). Internal service-to-service calls use ECS Service Connect DNS (`http://audit-core` within the VPC namespace).

---

### Service 1: API Gateway

**Language:** Go
**Database:** None
**Framework:** Standard library `net/http` + `httputil.ReverseProxy`

Responsibilities:
- Verify JWT signature using Identity Service's RSA public key (loaded at startup, refreshed on rotation via a background goroutine).
- Inject `X-User-Id`, `X-Firm-Id`, `X-User-Role` headers for downstream services. Downstream services trust these headers — they never re-validate the JWT.
- Route requests to the appropriate downstream service by URL prefix.
- Enforce rate limits per `firm_id` (token bucket, in-memory per pod — sufficient at launch scale).
- Return 401 for missing/invalid JWT before any downstream call is made.

This is infrastructure, not a domain service. It has no business logic and no database. Target: under 500 lines of Go.

---

### Service 2: Identity Service

**Language:** Go
**Database:** `identity_db` (own Postgres database on the shared RDS instance)
**Framework:** Chi + oapi-codegen

Owns:
- `Firm` — root tenant entity
- `User` — firm staff and client-side users; roles
- `Client` — entity being audited
- `MethodologyTemplate`, `TemplateControl`, `TemplateTestProcedure` — firm-level reusable templates

Responsibilities:
- User registration, login (email/password + MFA), SSO via SAML (Microsoft/Google).
- JWT issuance and refresh. JWTs are signed with an RSA private key held only by this service.
- RBAC: role assignment, permission checks. Roles: `FirmAdmin | Partner | Manager | Staff | EQReviewer | ClientAdmin | ClientUser | ViewOnly`.
- Firm settings and subscription tier management.
- Methodology template CRUD.

Client users (`ClientAdmin`, `ClientUser`) belong to a `Client` record in this service, not a `Firm`. Their JWT encodes the `client_id` and the specific engagement IDs they are invited to.

---

### Service 3: Audit Core

**Language:** Go
**Database:** `core_db` — shared Postgres with RLS (see Section 3)
**Framework:** Chi + oapi-codegen
**Background jobs:** River (Postgres-backed, uses `core_db`)
**Step Functions state machines:** `EngagementLifecycleStateMachine`, `DocumentRequestReminderStateMachine`

Owns:
- `Engagement`, `EngagementTeamMember`, `EngagementFramework`
- `ClientAcceptance`, `EngagementQualityReview`
- `FirmControlObjective`, `FirmControlObjectiveMapping`
- `Control`, `TestProcedure`
- `EvidenceItem`, `EvidenceLink`
- `DocumentRequest`
- `AIDecision`
- `AuditLog` (append-only, insert-only PostgreSQL rule)

System tables (seeded, not tenant-scoped, updated via migrations):
- `Framework`, `FrameworkRequirement`
- `ControlObjectiveLibrary`, `ControlObjectiveLibraryMapping`

River workers (background jobs within `core_db`):
- `document.extract` — calls Document Processing Service via HTTP; stores extracted text in `EvidenceItem.extracted_text`
- `document.embed` — generates embeddings via Claude API; stores vectors in pgvector
- `ai.completeness-check` — per document upload, checks completeness against request
- `ai.nightly-sweep` — engagement-level completeness review
- `ai.batch-control-mapping` — maps new FirmControlObjectives to FrameworkRequirements
- `email.notification` — all transactional email via SES

This is the largest service by entity count and intentionally so. The evidence chain (`EvidenceItem → EvidenceLink → TestProcedure → Control → FirmControlObjective → FrameworkRequirement`) is the product's core differentiator and requires ACID transactions across these entities. Splitting it would require distributed transactions or eventual consistency on operations that must be atomic (e.g., accepting a document request must atomically create an `EvidenceLink` and update `DocumentRequest.status`).

---

### Service 4: Trial Balance Service

**Language:** Go
**Database:** `trial_balance_db` (own Postgres database)
**Framework:** Chi + oapi-codegen

Owns:
- `TrialBalance`, `TrialBalanceAccount`, `TrialBalanceAdjustment`

`engagement_id` is stored as a plain UUID column — no foreign key, no join to `core_db`. If a request requires validating that the engagement exists and the requesting user has access, the service makes one REST call to Audit Core at the start of the handler and caches the result for the duration of the request.

Exists only for `FinancialAudit` engagement types. Population analysis (gap testing, duplicate detection, threshold filtering, Benford's law distribution analysis) runs as SQL queries against `trial_balance_db`. No application-layer computation for bulk analytics.

The spreadsheet UI (AG Grid + HyperFormula) has distinct scaling and collaboration requirements from the rest of the product, which justifies independent deployment. This service is the most likely candidate to be rewritten (e.g., if a dedicated spreadsheet service like Univer is adopted).

---

### Service 5: Workpaper Service

**Language:** Go
**Database:** `workpaper_db` (own Postgres database)
**Framework:** Chi + oapi-codegen + WebSocket (Gorilla WebSocket or nhooyr/websocket)

Owns:
- `Workpaper`, `WorkpaperVersion`
- Yjs document awareness state (in-memory per document, persisted to `workpaper_db` on save)

`engagement_id` and `control_id` are plain UUID references — no foreign keys to `core_db`.

The WebSocket server for Yjs real-time collaboration has different scaling characteristics from REST API services (long-lived connections vs stateless request/response). This is the primary reason to keep workpapers as a separate service. Workpaper tasks are scaled independently via a custom CloudWatch metric (active WebSocket connection count), not request throughput.

The `is_ai_draft` flag on `WorkpaperVersion` is cleared when any human edits content (satisfying PCAOB AS 1105). This logic runs entirely within this service.

---

### Service 6: Reporting Service

**Language:** Go
**Database:** `reporting_db` (own Postgres database)
**Framework:** Chi + oapi-codegen

Owns:
- `Report`, `ReportVersion`

Report generation is an async operation (not a synchronous API response). The Reporting Service runs its own River instance backed by `reporting_db`. When a report is requested, a River job is enqueued in `reporting_db`. The worker:
1. Calls Audit Core REST API for engagement data, controls, evidence, and workpapers.
2. Calls Trial Balance REST API for trial balance data (if financial audit).
3. Calls Workpaper REST API for workpaper content.
4. Renders the report using a Go template.
5. Stores the rendered report in S3 and the metadata in `reporting_db`.

Finalized and archived reports use S3 Object Lock (WORM) to satisfy regulatory immutability requirements. `Report` records transition to read-only in `reporting_db` at the same time.

---

### Service 7: Document Processing Service

**Language:** Python
**Database:** None (stateless)
**Framework:** FastAPI + uvicorn

Single endpoint: `POST /extract`

Request: multipart form with the file bytes and content type.
Response: `{ "text": "...", "pages": [...], "metadata": { "page_count": N, "has_tables": bool, "is_scanned": bool } }`

Internally:
- `pdfplumber` for digital PDFs (layout-aware text extraction, table detection).
- `pytesseract` (wrapping Tesseract) for scanned documents where `is_scanned` is detected.
- No state retained between calls.

Called exclusively by Audit Core's `document.extract` River worker. Not exposed through the API Gateway.

---

## 3. Database Topology

### Physical Layout

One RDS PostgreSQL instance (Multi-AZ for production, Single-AZ for staging/dev), five logical databases:

```
RDS PostgreSQL (Multi-AZ)
├── identity_db          → Identity Service
├── core_db              → Audit Core (RLS enabled)
├── trial_balance_db     → Trial Balance Service
├── workpaper_db         → Workpaper Service
└── reporting_db         → Reporting Service
```

Each service has its own Postgres user with access only to its own database. No cross-database queries. Cross-service data is resolved via REST API calls at runtime.

PgBouncer (transaction-mode connection pooling) is deployed as an ECS sidecar container in each service's task definition. Each service connects to PgBouncer at `localhost:6432`, not directly to RDS.

### core_db: Multi-Tenancy via RLS

`core_db` uses PostgreSQL Row-Level Security for multi-tenancy. `firm_id` is indexed on every tenant-scoped table. Application sets `SET app.current_firm_id = $1` at session/transaction start. RLS policies enforce the firm boundary at the database layer.

The three authorization middleware functions from the v2 spec are preserved as Go middleware:
- `WithFirmIsolation` — reads `firm_id` from JWT headers, sets Postgres session variable
- `WithEngagementAccess` — verifies `EngagementTeamMember` record exists for the requested engagement
- `WithClientScoping` — for `ClientUser` roles, filters to invited engagements only

### Other databases: Application-Layer Isolation

`identity_db`, `trial_balance_db`, `workpaper_db`, and `reporting_db` enforce tenant isolation at the application layer (query always includes `WHERE firm_id = $1`). RLS is not required because these services are simpler and the isolation logic is less complex.

### pgvector

pgvector extension is enabled on `core_db`. Embedding vectors for `EvidenceItem` records are stored in a `evidence_embeddings` table in `core_db`, colocated with the evidence data they reference.

---

## 4. Inter-Service Communication

### Synchronous (REST over HTTP)

Used for request/response queries where the caller needs an immediate result. Examples:
- API Gateway → any service (all client-initiated requests)
- Trial Balance Service → Audit Core (validate engagement access)
- Workpaper Service → Audit Core (validate engagement and control access)
- Reporting Service → Audit Core, Trial Balance, Workpaper (assemble report data)

All internal service calls use ECS Service Connect DNS. No service mesh at launch — direct HTTP with standard retry/timeout middleware in each Go client. mTLS can be added via Amazon VPC Lattice post-launch if the security posture requires it.

### Asynchronous (AWS SQS)

Used for cross-service events where the producer does not need an immediate response. Examples:
- Audit Core → SQS → Document Processing triggered when a new `EvidenceItem` is uploaded
- Identity Service → SQS → Audit Core notified when a user is deactivated (revoke engagement access)

SQS standard queues (at-least-once delivery). Consumers are idempotent — processing the same event twice has no side effects (idempotency key stored in the relevant table).

### Internal Async (River in core_db)

All background jobs that stay within Audit Core use River, backed by `core_db`. These never cross a service boundary and do not use SQS. River provides durable job execution with retry and dead-letter queues, using the existing PostgreSQL connection — no additional infrastructure.

---

## 5. Go Tech Stack (Per Service)

| Concern | Choice | Rationale |
|---|---|---|
| HTTP framework | [Chi](https://github.com/go-chi/chi) | Lightweight, idiomatic, uses standard `net/http` — no lock-in to framework-specific types |
| API contract | OpenAPI 3.1 spec (written first) + [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) | API-first: spec is the contract; codegen produces server interfaces and typed request/response structs |
| Database access | [sqlc](https://sqlc.dev/) + [pgx/v5](https://github.com/jackc/pgx) | Type-safe SQL — queries are plain SQL files, sqlc generates Go functions; no ORM magic hiding query behavior (important for audit-grade explainability) |
| Migrations | [golang-migrate](https://github.com/golang-migrate/migrate) | SQL migration files, supports Postgres, integrates with CI |
| Background jobs | [River](https://riverqueue.com/) | Postgres-backed job queue, Go-native equivalent of pg-boss; uses existing `core_db`, no additional infrastructure |
| Config | [envconfig](https://github.com/kelseyhightower/envconfig) | 12-factor config from environment variables with struct tags |
| Testing | [testify](https://github.com/stretchr/testify) + [httptest](https://pkg.go.dev/net/testing/httptest) | Standard Go HTTP testing; integration tests use a real Postgres instance (not mocks) |
| Logging | [slog](https://pkg.go.dev/log/slog) (stdlib) | Structured logging; no external dependency |
| Tracing | OpenTelemetry Go SDK → AWS X-Ray | AWS-native tracing, no additional vendor |
| AWS SDK | [aws-sdk-go-v2](https://github.com/aws/aws-sdk-go-v2) | Bedrock (AI model inference via VPC endpoint), Step Functions (workflow state machines), SQS (cross-service events) |

---

## 6. Python Tech Stack (Document Processing Service)

| Concern | Choice |
|---|---|
| Framework | FastAPI + uvicorn |
| PDF extraction | pdfplumber |
| OCR | pytesseract (wraps system Tesseract install) |
| Dependency management | uv + pyproject.toml |
| Container | python:3.13-slim base; Tesseract installed via apt |

---

## 7. Monorepo Structure

```
apps/
  gateway/           — Go: API Gateway
  identity/          — Go: Identity Service
  audit-core/        — Go: Audit Core (largest service)
  trial-balance/     — Go: Trial Balance Service
  workpaper/         — Go: Workpaper Service
  reporting/         — Go: Reporting Service
  doc-processing/    — Python: Document Processing Service

packages/
  go-shared/         — Shared Go: JWT middleware, SQS client wrappers,
                        OpenTelemetry setup, common error types
  openapi/           — OpenAPI specs for all services (source of truth)
  ai/                — Go: Bedrock client wrappers, AIDecision recording

infra/
  modules/           — Reusable Terraform modules (vpc, ecs-service, rds, etc.)
  workspaces/        — Layer-based Terraform workspaces (network, data, compute, etc.)
  envs/              — Per-environment tfvars (dev, staging, prod)
```

Each Go service has its own `go.mod`. `go-shared` is a local module referenced by each service. The Python service has its own `pyproject.toml`.

Turborepo manages the monorepo with per-service build caching. Go services build to single binaries; the Python service builds a Docker image. Infrastructure is provisioned via Terraform with layer-based workspace segmentation — see [`infrastructure-design.md`](./infrastructure-design.md) for full details.

---

## 8. Infrastructure Changes from v2 Spec

| Component | v2 Spec | This Design |
|---|---|---|
| API layer | tRPC | REST + OpenAPI |
| Primary language | TypeScript/Node.js | Go |
| PDF service | Python FastAPI (unchanged) | Python FastAPI (unchanged) |
| ORM/DB access | Prisma | sqlc + pgx |
| Background jobs | pg-boss (Node.js) | River (Go) |
| Database structure | Single shared Postgres | Shared `core_db` + 4 separate databases |
| Hasura | Rejected (unchanged) | N/A |
| Temporal | TypeScript SDK | Step Functions Standard Workflows (no self-hosted or cloud Temporal dependency) |
| AI model API | Anthropic direct | AWS Bedrock (PrivateLink, IAM auth, single AWS sub-processor) |
| Transactional email | SES or SendGrid (undecided) | SES |
| Frontend type sharing | tRPC inferred types | openapi-typescript (generated from OpenAPI spec) |
| Container orchestration | EKS (Kubernetes) | ECS Fargate — no node management, no control plane upgrades; ECS Service Connect for internal DNS; Amazon VPC Lattice for deferred mTLS |
| Infrastructure-as-code | Undefined | Terraform with layer-based workspace segmentation (network, data, compute, dns-cdn, observability, cicd) — see [`infrastructure-design.md`](./infrastructure-design.md) |
| CI/CD | Undefined | GitHub Actions with OIDC federation to AWS (no long-lived credentials) |
| AWS account structure | Undefined | Multi-account via AWS Organizations (management, tooling, dev, staging, prod) |
| Observability | Undefined | CloudWatch Logs + Metrics + Dashboards + Alarms, X-Ray via OpenTelemetry |

The frontend (`apps/web`) remains TypeScript + React. `openapi-typescript` generates typed API clients from the OpenAPI specs in `packages/openapi/`. The code-generation step runs as part of the Turborepo build pipeline — a spec change automatically regenerates the client on the next build.

---

## 9. What This Design Defers

- **mTLS between services** — deferred until compliance review requires it. Amazon VPC Lattice provides mTLS natively for ECS services without a sidecar service mesh. Current posture: VPC internal network trusted, JWT already validates user identity at the gateway.
- **Database-per-service for core_db** — the tightly coupled evidence chain makes this impractical without distributed transactions. If a specific entity cluster within `core_db` needs to scale or be rewritten independently, it can be extracted at that time with a defined migration path.
- **gRPC for internal service communication** — REST is simpler to debug, observe, and test. gRPC can replace REST for high-frequency internal calls if profiling shows REST overhead is significant.
