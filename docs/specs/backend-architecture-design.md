# Axiom Backend Architecture Design
**Date:** 2026-04-16
**Status:** Approved
**Supersedes:** Section 7 (Technology Stack) of `axiom-spec-design.md`; previous microservices decomposition (2026-04-12)
**Related:** [`domain-and-data-model-design.md`](./domain-and-data-model-design.md) (authoritative domain model and entity definitions)

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

4. **Modular monolith over microservices.** The previous version of this design decomposed the backend into seven microservices aligned to bounded contexts. That architecture is technically sound but mismatched to the current development context:
   - **Team size.** The product is being built by a single engineer with an AI coding agent. Microservices impose operational overhead (7 ECS services, 5 databases, 5 sets of migrations, cross-service REST contracts, distributed tracing) that consumes development time without delivering value until there are multiple independent teams deploying on different cadences.
   - **AI agent effectiveness.** AI coding agents perform significantly better when they can see the full codebase — handlers, business logic, database queries, and domain types in one context window. In a microservices architecture, the agent must reason about HTTP contracts between services instead of function calls, and cannot trace a request end-to-end.
   - **Cross-module transactions.** The evidence chain (`EvidenceItem → EvidenceItemSupports → CommonControl → CommonControlSatisfies → FrameworkRequirement`) spans what were previously separate databases. The microservices design required cross-service REST calls to resolve framework mappings, internal REST endpoints for AIDecision creation, and SQS queues for cross-service events. In a monolith, these are function calls with ACID transactions.
   - **Industry direction.** A 2025 CNCF survey found 42% of organizations that adopted microservices have consolidated services back into larger units. Amazon Prime Video, Shopify, Basecamp, GitHub, and Stack Overflow have demonstrated that well-structured monoliths scale further than most products will ever need.
   - **Extraction path.** Go packages provide natural module boundaries. Each module exposes a Go interface; replacing an in-process call with an HTTP call is a mechanical change. Any module can be extracted to a separate service when team size or scaling requirements justify it.

   The Python Document Processing service remains separate — it has a different runtime, different resource profile (CPU-bound OCR), and no shared state.

5. **Cross-framework mapping runs in PostgreSQL.** The common-control catalog, framework-requirement crosswalks (SCF/OSCAL/AICPA/CIS), and NIST STRM-vocabulary satisfaction edges are modeled as junction tables with effective-dated rows. Gap analysis, coverage percentages, and multi-framework rollups run as SQL queries against the database. No graph database is introduced — PostgreSQL + pgvector is sufficient for the control-centric directed labeled graph. The application layer handles orchestration and result formatting only.

**AWS infrastructure** — account structure, VPC design, ECS/Fargate configuration, RDS instance sizing, S3 bucket policies, CI/CD pipeline, Terraform workspace segmentation, observability, and security controls are specified in [`infrastructure-design.md`](./infrastructure-design.md).

---

## 2. Architecture Overview

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
│                    Axiom API  (Go)                            │
│   Single binary — all domain modules + background workers    │
│                                                              │
│  ┌──────────┐ ┌──────────┐ ┌───────────┐ ┌──────────┐      │
│  │ Gateway  │ │ Identity │ │ Audit Core│ │Frameworks│      │
│  │(midware) │ │          │ │           │ │(+mapping)│      │
│  └──────────┘ └──────────┘ └───────────┘ └──────────┘      │
│  ┌──────────┐ ┌──────────┐ ┌───────────┐ ┌──────────┐      │
│  │Workpaper │ │Reporting │ │Provenance │ │    AI    │      │
│  │ (+WS)    │ │          │ │ (signing) │ │ (Bedrock)│      │
│  └──────────┘ └──────────┘ └───────────┘ └──────────┘      │
│                                                              │
│  River background workers (all modules, one instance)        │
└──────────┬───────────────────────────────────────────────────┘
           │ HTTP (document extraction only)
┌──────────▼───────────────────────────────────────────────────┐
│              Document Processing  (Python)                    │
│              FastAPI · pdfplumber · Tesseract                 │
└──────────────────────────────────────────────────────────────┘
           │
      PostgreSQL
     (single DB, RLS)
```

Two deployed services:

- **Axiom API** — a single Go binary containing all domain modules (identity, audit core, frameworks/control mapping, workpaper, reporting, provenance) plus the gateway middleware, AI integration layer, and River background workers. Runs as an ECS Fargate service.
- **Document Processing** — a stateless Python service for PDF text extraction and OCR. Called via HTTP from the Axiom API's `auditcore.document-extract` River worker. Not exposed through the ALB.

Modules communicate via Go function calls, not HTTP. Each module defines a service interface that other modules depend on. Dependencies are wired at application startup.

---

## 3. Internal Modules

### Module: Gateway (Middleware)

**Package:** `internal/gateway`

Not a separate service — implemented as Chi middleware functions composed into the router.

Responsibilities:
- Verify JWT signature using the Identity module's RSA public key (loaded at startup, refreshed on rotation via a background goroutine).
- Inject `X-User-Id`, `X-Firm-Id`, `X-User-Role` into the request context for downstream handlers.
- Enforce rate limits per `firm_id` (token bucket, in-memory — sufficient at launch scale).
- Return 401 for missing/invalid JWT before any handler is invoked.

Target: under 300 lines of Go.

---

### Module: Identity

**Package:** `internal/identity`

Owns:
- `Firm` — root tenant entity
- `User` — firm staff and client-side users; roles
- `Client` — entity being audited
- `Invitation` — magic link staff onboarding tokens
- `MethodologyTemplate`, `TemplateControl`, `TemplateTestProcedure`, `TemplateDocumentRequest` — firm-level reusable templates

Responsibilities:
- User registration, login (email/password + MFA), SSO via SAML (Microsoft/Google).
- JWT issuance and refresh. JWTs are signed with an RSA private key held only by this module's configuration.
- RBAC: role assignment, permission checks. Roles: `FirmAdmin | Partner | Manager | Staff | EQReviewer | ClientAdmin | ClientUser | ViewOnly`.
- Firm settings and subscription tier management.
- Staff invitation management (magic link issuance, expiry, day-5 reminders).
- Methodology template CRUD (including pre-drafted document request templates).

Client users (`ClientAdmin`, `ClientUser`) belong to a `Client` record, not a `Firm`. Their JWT encodes the `client_id` and the specific engagement IDs they are invited to.

---

### Module: Audit Core

**Package:** `internal/auditcore`

The largest module by entity count. Owns the central domain model.

Owns:
- `Engagement`, `EngagementTeamMember`, `EngagementFramework`
- `ClientAcceptance`, `EngagementQualityReview`, `EQRFinding`
- `Control`, `TestProcedure` — engagement-scoped instances tied to a `CommonControl`
- `Finding`, `ManagementResponse`, `CorrectiveActionPlan` — findings and remediation lifecycle
- `EvidenceItem`, `evidence_embeddings` (pgvector)
- `DocumentRequest`, `ClientHubToken`, `DelegationToken`
- `AIDecision` — shared cross-service HITL ledger; every AI output written here
- `AuditLog` (append-only, insert-only PostgreSQL rule)
- `Notification` (in-app + email delivery)

**Step Functions state machines** (invoked via AWS SDK):
- `EngagementLifecycleStateMachine` — engagement status transitions with guard conditions
- `DocumentRequestReminderStateMachine` — reminder escalation sequence

River workers:
- `auditcore.document-extract` — calls Document Processing service via HTTP; stores extracted text in `EvidenceItem.extracted_text`
- `auditcore.document-embed` — generates embeddings via Bedrock; stores vectors in pgvector
- `auditcore.ai-completeness-check` — per document upload, checks completeness against request. Claude Haiku.
- `auditcore.ai-nightly-sweep` — engagement-level completeness review. Claude Haiku.
- `auditcore.ai-risk-category-suggestion` — triggered when a partner opens the ClientAcceptance form. Suggests risk categories based on client/engagement context. Creates AIDecision record. Claude Sonnet. Low volume (one per engagement).
- `auditcore.ai-findings-triage` — on new or updated finding, proposes severity classification and cross-framework impact surface. Creates AIDecision record. Claude Sonnet.
- `auditcore.ai-management-response-drafter` — agentic remediation loop: drafts management response, opens/ties to ticketing system (Jira/Linear/GitHub), round-trips closure evidence back to the finding. Creates AIDecision record per drafting step. Claude Sonnet.
- `auditcore.notification-deliver` — creates Notification records and delivers transactional email via SES based on recipient notification preferences.

The evidence chain (`EvidenceItem → EvidenceItemSupports → CommonControl → CommonControlSatisfies → FrameworkRequirement`) is fully within one database with ACID transactions. `Control.common_control_id` references `CommonControl` in the frameworks module's tables — since they share the same database, this is a real foreign key with referential integrity.

---

### Module: Frameworks (Control Mapping)

**Package:** `internal/frameworks`

The cross-framework intelligence layer. Owns the common-control catalog, framework-requirement crosswalks, and NIST STRM-encoded satisfaction edges. Named `frameworks` to align with the Phase 2 implementation plan (`internal/frameworks`), which already positions it as the reference-data and mapping home.

Owns:
- `Framework`, `FrameworkVersion`, `FrameworkRequirement` — platform-seeded framework catalog (SOC 2, ISO 27001/27701, ISO 42001, HIPAA/HITRUST r2, PCI DSS 4.0.1, SOC 1). Versioned; `FrameworkVersion` carries effective-dated windows.
- `CommonControl`, `CommonControlVersion` — the unified control catalog. Platform seed from SCF (primary) + OSCAL (NIST-family) + AICPA official mappings + CIS Controls v8.1 (secondary cross-check). Firms may extend with firm-scoped controls (RLS-enforced).
- `CommonControlSatisfies` — directed, labeled edges from `CommonControl → FrameworkRequirement`. Carries NIST STRM relationship vocabulary (`equivalent-to | subset-of | superset-of | intersects-with | no-relationship`), a strength score, a coverage percentage, and an effective-dated window. Multi-edge per pair permitted across framework versions.
- `EvidenceItemSupports` — directed, labeled edges from `EvidenceItem → CommonControl`. Carries coverage percentage, period window (respecting framework-specific age tolerances — e.g., ASV scan 90d, pen test 1y, background check 1y), and an STRM relationship tag.
- `PolicyLibrary`, `PolicyLibraryEmbedding` — firm-scoped policy corpus used for semantic retrieval during mapping and response drafting.
- `framework_requirement_embeddings`, `common_control_embeddings`, `policy_library_embeddings` (pgvector).

Responsibilities:
- CommonControl CRUD; SCF/OSCAL/AICPA/CIS crosswalk import and reconciliation.
- FrameworkRequirement catalog maintenance, including version migration helpers (PCI 3.2→4.0, ISO 27001:2013→2022, NIST CSF 1.1→2.0).
- Gap analysis queries across multi-framework engagements — percent coverage per `FrameworkRequirement`, list of unsatisfied requirements, partial-coverage ranking.
- Semantic search over requirement text and policy library (pgvector cosine similarity).
- Exposes `FrameworkCatalog` and `ControlMapper` Go interfaces consumed by `auditcore`, `workpaper`, and `reporting`.

River workers:
- `frameworks.evidence-control-mapping` — on new evidence upload or edit, suggests which `CommonControl`s the evidence supports and the coverage percentage. Creates `AIDecision` record via `auditcore.AIDecisionService.Create` (function call). Claude Haiku.
- `frameworks.gap-analysis` — scheduled (nightly per active engagement) and on-demand. Produces a ranked gap list per framework version in the engagement. Creates `AIDecision` when remediation suggestions are generated. Claude Sonnet for reasoning, Haiku for routing.
- `frameworks.drift-detection` — monitors connector-sourced configs (cloud providers, identity providers, dev tools, HRIS) for drift against previously accepted evidence; when drift exceeds threshold, enqueues a retest job and creates an `AIDecision`. Claude Haiku for classification.
- `frameworks.framework-migration` — bulk remapping when a framework version is superseded. Proposes `CommonControlSatisfies` edge updates; each proposal is an `AIDecision`. Claude Sonnet.
- `frameworks.scf-import` — quarterly refresh of the SCF catalog and NIST STRM mappings. Diff-based import; generates `AIDecision` records only for edges that conflict with firm-accepted mappings. Claude Haiku for classification; no LLM call if the diff is structurally clean.

---

### Module: Provenance

**Package:** `internal/provenance`

Cryptographic provenance for evidence and AI outputs. Post-Delve differentiator (see compliance-pivot-findings.md §3.3): evidence artifacts and AI decisions must be auditor-defensible by construction.

Owns:
- `ProvenanceRecord` — signed envelope for each evidence artifact or AI output; stores content hash (SHA-256), signature, signing key identifier (KMS key ARN), timestamp, and WORM S3 object lock reference.
- `SignedScreenshot`, `HashedDOMSnapshot` — typed provenance records for connector-captured browser evidence.
- `AIDecisionProvenance` — links `AIDecision` records to the prompt hash, model identifier, Bedrock request ID, and output hash; together these form the AIDecision audit chain.

Responsibilities:
- Sign evidence artifacts on ingestion using AWS KMS (asymmetric signing key, per-firm or platform-level depending on artifact class).
- Hash and sign DOM snapshots and screenshots from browser-based evidence capture.
- Persist finalized artifacts to S3 with Object Lock (WORM) in compliance mode; mirror the lock metadata into `ProvenanceRecord`.
- Verify provenance on retrieval (hash match, signature valid, lock status intact) and expose verification status to `reporting` and `workpaper`.
- Propagate provenance tags through the AIDecision audit chain so every downstream artifact that cites AI output carries its signed lineage.
- Exposes the `ProvenanceSigner` Go interface consumed by `auditcore`, `frameworks`, `workpaper`, `reporting`, and `ai`.

River workers:
- `provenance.sign-evidence` — triggered on evidence upload or connector capture. Signs artifact, writes WORM copy to S3, creates `ProvenanceRecord`.
- `provenance.verify-chain` — scheduled verification sweep over recent `AIDecision` records and their cited evidence; flags any broken chain for auditor review.

---

### Module: Workpaper

**Package:** `internal/workpaper`

Owns:
- `Workpaper`, `WorkpaperVersion`
- `ReviewNote` — structured reviewer feedback anchored to workpaper content (cannot be deleted — AU-C 230)
- Yjs document awareness state (in-memory per document, persisted to database on save)

Handles WebSocket connections for Yjs real-time collaboration on a dedicated route group (`/api/v1/workpapers/ws/*`). The ALB natively supports WebSocket upgrade requests.

River workers:
- `workpaper.ai-workpaper-draft` — triggered when an auditor explicitly requests a narrative draft after marking a TestProcedure as Complete. Calls Bedrock (Claude Sonnet) with control description, test procedure, linked evidence text (with provenance references), exceptions, prior year workpaper (if rollforward), and firm template. Inserts draft text into the workpaper editor with `ai_content_metadata` tracking per section. Creates `AIDecision` record with `context_type = WorkpaperDraft` and attaches provenance tags from cited evidence.

**AI content tracking:** `WorkpaperVersion` carries `ai_content_metadata` (jsonb) that tracks AI origin per section — `ai_generated`, `human_edited`, and `modification_ratio` (Levenshtein distance between AI-generated text and current text, divided by AI character count). The `modification_ratio` is computed on each save, not in real-time. The `is_ai_draft` boolean is retained as a derived convenience field: true when any section has `ai_generated = true AND human_edited = false`. The advancement gate (`PreparedPendingReview`) checks: (1) all AI-generated sections must have `human_edited = true`, (2) sections with `modification_ratio < 0.05` trigger a confirmable warning. This supports ISO 42001 AI management controls and SOC 2 / ISO 27001 evidence-of-work expectations (human-in-the-loop).

---

### Module: Reporting

**Package:** `internal/reporting`

Owns:
- `Report`, `ReportVersion`

Report generation is an async operation. When a report is requested, a River job is enqueued. The worker:
1. Reads engagement data, controls, evidence, findings, and management responses from the auditcore module (function calls).
2. Reads framework/common-control mappings and gap-analysis results from the frameworks module (function call).
3. Reads workpaper content from the workpaper module (function call).
4. Verifies provenance for every cited artifact via the provenance module (function call); reports cannot be rendered if any cited artifact fails chain verification.
5. Renders the report using a Go template.
6. Stores the rendered report in S3 and the metadata in the database.

River workers:
- `reporting.report-generate` — assembles and renders reports as described above.
- `reporting.ai-report-section-draft` — triggered when a partner explicitly requests an AI draft of a specific report section. Calls Bedrock (Claude Sonnet) with report type/template, engagement-wide data (controls, test results, findings, cross-framework coverage statistics), prior year report (if rollforward), and firm template. Inserts draft text into the specific report section with `ai_content_metadata` tracking. Creates `AIDecision` record with `context_type = ReportSectionDraft` and provenance references to every cited artifact. AI may draft: Description of Tests of Controls (SOC 1/2), Scope and Approach (all), System Description summary (SOC 1/2), Statement of Applicability narrative (ISO 27001), AI system description (ISO 42001). AI does NOT draft: opinions, management assertions, scope limitations, or qualification language.

**AI content tracking:** `ReportVersion` carries `ai_content_metadata` (jsonb) with the same section-level tracking schema as `WorkpaperVersion` — `ai_generated`, `human_edited`, `modification_ratio`. Report issuance (`Report.status = Issued`) validates that all AI-drafted sections have been substantively edited: all sections must have `human_edited = true`, and sections with `modification_ratio < 0.05` trigger a confirmable warning.

Finalized and archived reports use S3 Object Lock (WORM) to satisfy regulatory immutability requirements; lock metadata is mirrored into the provenance module's `ProvenanceRecord` so the entire chain (evidence → AI output → workpaper → report) is verifiable post-issuance. `Report` records transition to read-only in the database at the same time.

---

### External Service: Document Processing

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

Called exclusively by the Axiom API's `auditcore.document-extract` River worker. Not exposed through the ALB. Accessed via ECS Service Connect DNS (`http://doc-processing`).

---

## 4. Database

### Single Database with RLS

One RDS PostgreSQL instance (Multi-AZ for production, Single-AZ for staging/dev), one logical database:

```
RDS PostgreSQL (Multi-AZ)
└── axiom_db    → Axiom API (RLS on all tenant-scoped tables)
```

The Axiom API connects with a single database user (`axiom_svc`). A separate `master` user is used for migrations and break-glass access.

PgBouncer (transaction-mode connection pooling) is deployed as an ECS sidecar container in the Axiom API's task definition. The application connects to PgBouncer at `localhost:6432`, not directly to RDS.

### Row-Level Security

All tenant-scoped tables use PostgreSQL Row-Level Security. `firm_id` is indexed on every tenant-scoped table. The application sets `SET app.current_firm_id = $1` at transaction start. RLS policies enforce the firm boundary at the database layer.

The three authorization middleware functions are implemented as Go middleware in the gateway module:
- `WithFirmIsolation` — reads `firm_id` from JWT context, sets Postgres session variable
- `WithEngagementAccess` — verifies `EngagementTeamMember` record exists for the requested engagement
- `WithClientScoping` — for `ClientUser` roles, filters to invited engagements only

### Module Table Ownership

Each module owns specific tables and accesses them via its own sqlc queries. Cross-module data is accessed through the module's service interface, not by querying another module's tables directly. This preserves clean boundaries and makes future service extraction straightforward.

| Module | Tables |
|---|---|
| Identity | `firms`, `users`, `clients`, `invitations`, `methodology_templates`, `template_controls`, `template_test_procedures`, `template_document_requests` |
| Audit Core | `engagements`, `engagement_team_members`, `engagement_frameworks`, `client_acceptances`, `engagement_quality_reviews`, `eqr_findings`, `controls`, `test_procedures`, `findings`, `management_responses`, `corrective_action_plans`, `evidence_items`, `evidence_embeddings`, `document_requests`, `client_hub_tokens`, `delegation_tokens`, `ai_decisions`, `audit_logs`, `notifications` |
| Frameworks | `frameworks`, `framework_versions`, `framework_requirements`, `framework_requirement_embeddings`, `common_controls`, `common_control_versions`, `common_control_embeddings`, `common_control_satisfies`, `evidence_item_supports`, `policy_library`, `policy_library_embeddings` |
| Provenance | `provenance_records`, `signed_screenshots`, `hashed_dom_snapshots`, `ai_decision_provenance` |
| Workpaper | `workpapers`, `workpaper_versions`, `review_notes` |
| Reporting | `reports`, `report_versions` |

Foreign keys that cross module boundaries (e.g., `controls.common_control_id → common_controls.id`, `evidence_item_supports.evidence_item_id → evidence_items.id`, `ai_decision_provenance.ai_decision_id → ai_decisions.id`) are enforced at the database level — a significant advantage over the previous multi-database design where these were plain UUID references without referential integrity.

**Platform-seeded vs. firm-scoped controls.** `frameworks`, `framework_versions`, `framework_requirements`, and the platform-provided rows of `common_controls` / `common_control_satisfies` are seeded platform-wide and not tenant-scoped — no RLS. Firm-authored extensions in `common_controls` and edges in `common_control_satisfies` and `evidence_item_supports` carry a `firm_id` and are RLS-enforced like any other tenant-scoped table. The SCF quarterly refresh updates only platform-seeded rows and flags conflicts with firm-accepted mappings for review.

### pgvector

The pgvector extension is enabled on `axiom_db`. Embedding vectors are stored for: `EvidenceItem` records (`evidence_embeddings`), `FrameworkRequirement` records (`framework_requirement_embeddings`), `CommonControl` records (`common_control_embeddings`), and `PolicyLibrary` records (`policy_library_embeddings`). These embeddings support evidence→control mapping suggestions, semantic gap analysis, framework version migration, and management-response drafting retrieval.

---

## 5. Inter-Module Communication

### Internal (Go Function Calls)

Modules communicate via Go interfaces. Each module exports a service interface; other modules depend on the interface, not the implementation. Dependencies are wired at application startup via constructor injection.

Principal inter-module interfaces:

| Interface | Exported by | Consumers |
|---|---|---|
| `AIDecisionService` | `auditcore` | All modules that produce AI output — `frameworks`, `workpaper`, `reporting`, `provenance`, `ai` |
| `FrameworkCatalog` | `frameworks` | `auditcore` (engagement scoping), `reporting` (coverage rollups), `workpaper` (procedure drafting) |
| `ControlMapper` | `frameworks` | `auditcore` (evidence linking), `reporting` (gap analysis), `ai` (mapping suggestions) |
| `ProvenanceSigner` | `provenance` | `auditcore` (evidence ingestion), `frameworks` (connector captures), `workpaper` (draft citation), `reporting` (finalization), `ai` (AIDecision signing) |
| `EvidenceStore` | `auditcore` | `frameworks` (mapping workers), `workpaper`, `reporting` |

Example: when the Frameworks module's `evidence-control-mapping` worker produces a suggestion, it calls `auditcore.AIDecisionService.Create()` and `provenance.ProvenanceSigner.SignAIOutput()` — both function calls within the same process, same database transaction if needed. No HTTP, no serialization, no retry logic.

**Cross-service AIDecision pattern.** Every AI output produced anywhere in the system — mapping suggestion, gap analysis, findings triage, management response draft, workpaper draft, report section draft — is recorded via the single shared `AIDecisionService` and signed through `ProvenanceSigner`. This is ISO 42001-native: there is exactly one ledger, one signature chain, one human-in-the-loop gate.

This replaces:
- All internal REST API calls between services (the previous design had 6 cross-service REST patterns)
- All SQS queues for cross-service events (3 queues eliminated)
- The internal `/internal/ai-decisions` endpoint pattern

### External (HTTP to Document Processing)

The only HTTP inter-service call is from the Axiom API to the Document Processing service. The `auditcore.document-extract` River worker sends the file bytes via `POST /extract` to `http://doc-processing:8000` (ECS Service Connect DNS). Standard retry with exponential backoff on failure.

### Background Jobs (River)

One River instance backed by `axiom_db`. All modules register their workers with the shared River client. Job types are prefixed by module to avoid collisions:

| Module | Job types |
|---|---|
| Audit Core | `auditcore.document-extract`, `auditcore.document-embed`, `auditcore.ai-completeness-check`, `auditcore.ai-nightly-sweep`, `auditcore.ai-risk-category-suggestion`, `auditcore.ai-findings-triage`, `auditcore.ai-management-response-drafter`, `auditcore.notification-deliver` |
| Frameworks | `frameworks.evidence-control-mapping`, `frameworks.gap-analysis`, `frameworks.drift-detection`, `frameworks.framework-migration`, `frameworks.scf-import` |
| Provenance | `provenance.sign-evidence`, `provenance.verify-chain` |
| Workpaper | `workpaper.ai-workpaper-draft` |
| Reporting | `reporting.report-generate`, `reporting.ai-report-section-draft` |

**Bedrock model assignments.** Haiku for high-volume classification and routing work (`frameworks.evidence-control-mapping`, `frameworks.drift-detection`, `frameworks.scf-import`, `auditcore.ai-completeness-check`, `auditcore.ai-nightly-sweep`). Sonnet for reasoning-heavy workloads (`frameworks.gap-analysis`, `frameworks.framework-migration`, `auditcore.ai-risk-category-suggestion`, `auditcore.ai-findings-triage`, `auditcore.ai-management-response-drafter`, `workpaper.ai-workpaper-draft`, `reporting.ai-report-section-draft`). Cost assumptions remain aligned with `ai-architecture-design.md` — Haiku dominates request volume; Sonnet dominates token cost.

River provides durable job execution with retry and dead-letter handling, using the existing PostgreSQL connection — no additional infrastructure.

### Workflow Engine (Step Functions)

Two Step Functions state machines, unchanged from the previous design:
- `EngagementLifecycleStateMachine` — long-running engagement status transitions with guard conditions and scheduled archival
- `DocumentRequestReminderStateMachine` — reminder escalation sequence

Invoked from the Audit Core module via the AWS SDK for Go.

---

## 6. Go Tech Stack

| Concern | Choice | Rationale |
|---|---|---|
| HTTP framework | [Chi](https://github.com/go-chi/chi) | Lightweight, idiomatic, uses standard `net/http` — no lock-in to framework-specific types |
| API contract | OpenAPI 3.1 spec (written first) + [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) | API-first: spec is the contract; codegen produces server interfaces and typed request/response structs. Specs organized per module for clean code generation. |
| Database access | [sqlc](https://sqlc.dev/) + [pgx/v5](https://github.com/jackc/pgx) | Type-safe SQL — queries are plain SQL files, sqlc generates Go functions; no ORM magic hiding query behavior (important for audit-grade explainability) |
| Migrations | [golang-migrate](https://github.com/golang-migrate/migrate) | SQL migration files, supports Postgres, integrates with CI. Single migration directory for all modules. |
| Background jobs | [River](https://riverqueue.com/) | Postgres-backed job queue, Go-native. One shared instance across all modules. |
| Config | [envconfig](https://github.com/kelseyhightower/envconfig) | 12-factor config from environment variables with struct tags |
| Testing | [testify](https://github.com/stretchr/testify) + [httptest](https://pkg.go.dev/net/testing/httptest) | Standard Go HTTP testing; integration tests use a real Postgres instance (not mocks) |
| Logging | [slog](https://pkg.go.dev/log/slog) (stdlib) | Structured logging; no external dependency |
| Tracing | OpenTelemetry Go SDK → AWS X-Ray | AWS-native tracing, no additional vendor |
| WebSocket | [nhooyr/websocket](https://github.com/nhooyr/websocket) or Gorilla WebSocket | Yjs collaboration in the workpaper module, served on the same binary |
| AWS SDK | [aws-sdk-go-v2](https://github.com/aws/aws-sdk-go-v2) | Bedrock (AI model inference via VPC endpoint), Step Functions (workflow state machines), S3 (evidence and report storage), SES (transactional email) |

---

## 7. Python Tech Stack (Document Processing)

| Concern | Choice |
|---|---|
| Framework | FastAPI + uvicorn |
| PDF extraction | pdfplumber |
| OCR | pytesseract (wraps system Tesseract install) |
| Dependency management | uv + pyproject.toml |
| Container | python:3.13-slim base; Tesseract installed via apt |

---

## 8. Monorepo Structure

```
apps/
  axiom-api/                — Go: Modular monolith (single binary)
    cmd/server/             — main.go entrypoint, dependency wiring
    internal/
      gateway/              — Chi middleware: JWT verification, rate limiting,
                              header injection, request context
      identity/             — Auth, RBAC, firm/user/client management,
                              methodology templates
      auditcore/            — Engagements, controls, evidence, findings,
                              management responses, document requests,
                              AI decisions, audit log, notifications
      frameworks/           — Framework catalog, common-control catalog,
                              SCF/OSCAL/AICPA/CIS crosswalks, STRM
                              satisfaction edges, evidence→control edges,
                              gap analysis, version migration, policy library
      provenance/           — Signed screenshots, hashed DOM snapshots,
                              KMS signing, WORM S3 Object Lock integration,
                              AIDecision audit chain verification
      workpaper/            — Workpapers, Yjs collaboration (WebSocket),
                              review notes
      reporting/            — Report generation, S3 archival
      ai/                   — Bedrock client, prompt templates, embedding
                              helpers, shared AI utilities
      integrations/         — Connector framework for cloud providers
                              (AWS, Azure, GCP), identity (Okta, Google
                              Workspace, Microsoft Entra), dev tools
                              (GitHub, GitLab, Jira, Linear), HRIS
                              (Rippling, Gusto, BambooHR) — evidence ingress
      platform/             — DB connection pool, config, error types,
                              OpenTelemetry setup, River client, common middleware
    migrations/             — All SQL migrations (single ordered sequence)
    go.mod
  doc-processing/           — Python: PDF extraction only
    pyproject.toml

packages/
  openapi/                  — OpenAPI 3.1 specs (organized per module for
                              code generation; single API surface for frontend)

infra/
  modules/                  — Reusable Terraform modules (vpc, ecs-service, rds, etc.)
  workspaces/               — Layer-based Terraform workspaces (network, data, compute, etc.)
  envs/                     — Per-environment tfvars (dev, staging, prod)
```

The Axiom API has a single `go.mod`. Each internal package is a Go module boundary — packages import each other's exported interfaces, not internal types. The `platform` package provides shared infrastructure (database pool, configuration, telemetry) that all modules depend on.

Turborepo manages the monorepo with build caching. The Go service builds to a single binary; the Python service builds a Docker image. Infrastructure is provisioned via Terraform with layer-based workspace segmentation — see [`infrastructure-design.md`](./infrastructure-design.md) for full details.

---

## 9. Infrastructure Changes from v2 Spec

| Component | v2 Spec | This Design |
|---|---|---|
| API layer | tRPC | REST + OpenAPI |
| Primary language | TypeScript/Node.js | Go |
| Architecture | Microservices (7 services) | Modular monolith (1 Go binary) + Python extraction service |
| PDF service | Python FastAPI (unchanged) | Python FastAPI (unchanged) |
| ORM/DB access | Prisma | sqlc + pgx |
| Background jobs | pg-boss (Node.js) | River (Go) — single instance, all modules |
| Database structure | Single shared Postgres | Single database with RLS on all tenant-scoped tables |
| Inter-module communication | N/A (was tRPC) | Go function calls (interfaces); HTTP only for doc processing |
| Cross-service messaging | SQS queues | Eliminated — direct function calls + River job enqueue |
| Hasura | Rejected (unchanged) | N/A |
| Temporal | TypeScript SDK | Step Functions Standard Workflows |
| AI model API | Anthropic direct | AWS Bedrock (PrivateLink, IAM auth, single AWS sub-processor) |
| Transactional email | SES or SendGrid (undecided) | SES |
| Frontend type sharing | tRPC inferred types | openapi-typescript (generated from OpenAPI spec) |
| Container orchestration | EKS (Kubernetes) | ECS Fargate — 2 services (Axiom API + doc-processing) |
| Infrastructure-as-code | Undefined | Terraform with layer-based workspace segmentation |
| CI/CD | Undefined | GitHub Actions with OIDC federation to AWS |

The frontend (`apps/web`) remains TypeScript + React. `openapi-typescript` generates typed API clients from the OpenAPI specs in `packages/openapi/`. The code-generation step runs as part of the Turborepo build pipeline — a spec change automatically regenerates the client on the next build.

---

## 10. What This Design Defers

- **Service extraction** — any module can be extracted to a separate service by replacing its Go interface with an HTTP client. The most likely extraction candidate is the Workpaper module (WebSocket connections have different scaling characteristics than REST). Trigger: WebSocket connection count drives autoscaling beyond what's efficient for the combined service, or a second team needs independent deployment.
- **Database-per-module** — all modules share one database. If a module's query patterns diverge significantly (e.g., frameworks module gap-analysis queries competing with audit core transactional writes at scale), a read replica or database split can be introduced. The module table ownership boundaries make this a data migration, not a code rewrite.
- **mTLS for doc processing** — the Axiom API calls Document Processing over plaintext HTTP within the VPC. Amazon VPC Lattice provides mTLS natively if the security posture requires it.
- **gRPC for doc processing** — REST is simpler to debug. gRPC can replace it if profiling shows serialization overhead is significant for large documents.
