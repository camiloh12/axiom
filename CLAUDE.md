# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Name

The working name for this application is **Axiom**. Alternative names under consideration for a future rename (Prism, Tenet, Canon, Veris, Affirma, Throughline) are listed in `docs/specs/axiom-spec-design.md` — the current name is not final.

## Project Status

Phase 0 (dev environment + project scaffold) and Phase 1 (platform core + identity service) are implemented in `apps/axiom-api/` (Go) and `apps/web/` (React). Phase 2 (Frameworks, Templates & Engagement Creation) is planned but not yet implemented — see `docs/superpowers/plans/phase-2-frameworks-templates-engagements.md`. All architectural guidance below reflects the current specifications; downstream phases remain pending.

## Pivot Note

This project pivoted from "mixed financial audit + compliance" to **compliance-and-assurance only** (SOC 1, SOC 2, ISO 27001, ISO 27701, ISO 42001, HIPAA, PCI DSS). See `docs/specs/compliance-pivot-findings.md` for the research brief and decisions behind the pivot. Financial-audit concepts (trial balance, GAAS, PCAOB, materiality) are out of scope.

## Specification Documents

- **Product spec (summary hub):** `docs/specs/axiom-spec-design.md`
- **Domain & data model:** `docs/specs/domain-and-data-model-design.md`
- **AI architecture:** `docs/specs/ai-architecture-design.md` — 11 AI features, three-tier HITL policy, AIDecision ledger, cryptographic provenance, ISO 42001-native governance
- **Backend architecture:** `docs/specs/backend-architecture-design.md` — supersedes the tech stack section of the product spec
- **Infrastructure:** `docs/specs/infrastructure-design.md` — AWS account structure, networking, Terraform workspaces, CI/CD, observability, security controls, cost estimates
- **Pivot findings:** `docs/specs/compliance-pivot-findings.md` — research, competitive analysis (Agentive, Delve), cross-mapping feasibility, recommended licensing stack
- **Implementation plan:** `docs/superpowers/specs/implementation-plan-design.md` — 11-phase execution plan
- **Phase plans:** `docs/superpowers/plans/phase-*.md`

## Tech Stack

| Layer | Technology |
|---|---|
| Frontend | TypeScript, React (SPA on CloudFront) |
| Backend | Go (primary), Python (PDF extraction only) |
| API | REST + OpenAPI 3.1 (oapi-codegen, openapi-typescript) |
| Database | PostgreSQL — single `axiom_db` with RLS on all tenant-scoped tables |
| DB access | sqlc + pgx/v5 (no ORM) |
| Background jobs | River (Postgres-backed, Go-native) |
| Workflow engine | AWS Step Functions |
| AI inference | AWS Bedrock (Claude via VPC endpoint) |
| Infrastructure | AWS ECS Fargate + ALB + RDS (Multi-AZ), Terraform IaC |
| CI/CD | GitHub Actions with OIDC federation to AWS |
| Observability | CloudWatch + X-Ray (via OpenTelemetry) |
| Email | AWS SES |

## Architecture Overview

Modular monolith (single Go binary) with internal packages per bounded context, plus a separate Python service for PDF extraction and an isolated Go service for cryptographic provenance signing.

**Axiom API** (Go modular monolith):
- `internal/gateway` — Chi middleware: JWT verification, routing, rate limiting
- `internal/identity` — auth, RBAC, firm/user/client management, JWT issuance
- `internal/auditcore` — engagements, controls, evidence, document requests, findings, management responses, AI decisions
- `internal/frameworks` — `CommonControl` graph, `FrameworkRequirement` catalog, NIST STRM-encoded satisfaction edges, `EvidenceItemSupports` period-aware coverage, SCF/OSCAL/AICPA/CIS crosswalk import, gap analysis
- `internal/workpaper` — collaborative workpapers with Yjs real-time sync (WebSocket)
- `internal/reporting` — async report generation, S3 WORM archiving
- `internal/ai` — Bedrock client, prompt templates, embedding helpers
- `internal/platform` — DB, config, OTel, River, common middleware

**Provenance Signer** (Go, separate ECS service) — KMS `Sign`-only IAM surface for signing AI outputs and evidence artifacts (`ECC_NIST_P256`, `SIGN_VERIFY`). Isolated for blast-radius containment.

**Document Processing** (Python, separate ECS service) — stateless PDF extraction via pdfplumber + Tesseract.

Axiom API modules communicate via Go interfaces (function calls), not HTTP. Inter-service HTTP calls: Axiom API → Document Processing (PDF extraction) and Axiom API → Provenance Signer (signing).

See `docs/specs/backend-architecture-design.md` for full module descriptions and monorepo structure.

## Key Domain Concepts

- **Engagement** — a single compliance/assurance project with a defined framework (or set of frameworks, for integrated engagements), client, and lifecycle
- **Framework / FrameworkVersion / FrameworkRequirement** — the regulatory catalog (e.g., ISO 27001:2022 Annex A control A.5.15)
- **CommonControl** — Axiom's internal pivot node; one CommonControl can satisfy many FrameworkRequirements across frameworks
- **CommonControlSatisfies** — NIST STRM-encoded edge (`equivalent-to | subset-of | superset-of | intersects-with | no-relationship`) with a strength score, effective dating, and source attribution (SCF, UCF, OSCAL, AICPA, CIS, Axiom custom)
- **EvidenceItemSupports** — period-aware edge from an evidence artifact to a CommonControl, carrying `coverage_pct`, `period_start`/`period_end`, and AI provenance (so a SOC 2 Type II sampling window can be reconciled against an ISO 27001 surveillance window)
- **Finding / ManagementResponse / CorrectiveActionPlan** — exception lifecycle: deviation tracking → root-cause → owner → remediation plan → re-test
- **Workpapers / Test Papers** — documentation of test procedures and evidence
- **Client Hub** — both a PBC document-request portal and the auditee-side GRC workspace (continuous monitoring, policy library, drift alerts)
- **AIDecision** — append-only ledger of every AI-assisted decision (prompt, model, context, confidence, reviewer, action); the operational realization of ISO 42001 HITL compliance

## User Roles

`FirmAdmin | Partner | Manager | Staff | EQReviewer | ClientAdmin | ClientUser | ViewOnly`

## Security & Compliance

Axiom itself targets from day one: **SOC 2 Type 2**, **ISO 27001**, **ISO 27701**, **ISO 42001**.

- Encryption in transit and at rest (AWS KMS)
- RBAC + RLS; all access events in append-only audit log
- AI decisions logged via `AIDecision` table for human-in-the-loop compliance (ISO 42001-native)
- Cryptographic provenance — every AI output and connector-captured evidence artifact is signed at emission (KMS ECC_NIST_P256, SIGN_VERIFY), written to S3 Object Lock (WORM), and publicly verifiable
- Finalized reports use S3 Object Lock in COMPLIANCE mode for regulatory immutability
- GDPR/CCPA compliant; client firm owns all engagement data

Out of scope: PCAOB public-company audits, ISO Certification Body issuance, PCI QSA ROC sign-off, internal audit / SOX, HITRUST CSF r2 (deferred post-MVP).

## Document Dependencies

Specification documents have cascading dependencies. When an upstream document changes, propagate changes to all downstream artifacts before considering the work complete.

### Artifact Map

| ID | Artifact | Path |
|---|---|---|
| R | Research docs | `docs/research/*.md` |
| J | User Journeys | `docs/user-journeys/all-journeys.md` |
| D | Domain & Data Model | `docs/specs/domain-and-data-model-design.md` |
| AI | AI Architecture | `docs/specs/ai-architecture-design.md` |
| BE | Backend Architecture | `docs/specs/backend-architecture-design.md` |
| INF | Infrastructure Design | `docs/specs/infrastructure-design.md` |
| OA | OpenAPI specs (7 files) | `packages/openapi/*.yaml` |
| DPY | Infrastructure diagram script | `docs/diagrams/axiom_infrastructure.py` |
| DPNG | Infrastructure diagram image | `docs/diagrams/axiom_infrastructure.png` |
| PS | Product Spec (summary hub) | `docs/specs/axiom-spec-design.md` |
| MK | HTML Mockups | `mockups/journey-*/` + `mockups/README.md` (pivot-audit inventory) |
| DS | Design System | `.impeccable.md` |
| IP | Implementation Plan | `docs/superpowers/specs/implementation-plan-design.md` |
| PP | Phase Plans | `docs/superpowers/plans/phase-*.md` |

### Cascade Rules

When a document changes, update all listed downstream artifacts:

- **Research** (`docs/research/*.md`) →
  `axiom-spec-design.md` (§1-4, §8-9 reference research)

- **User Journeys** (`docs/user-journeys/all-journeys.md`) →
  `domain-and-data-model-design.md` (entity definitions tagged by journey, §11 traceability matrix),
  `ai-architecture-design.md` (feature journey references),
  `axiom-spec-design.md` (§11 journey summaries, §12 innovative flows),
  `mockups/` (directories and screens map 1:1 to journey stages; `mockups/README.md` tracks screen inventory and outstanding work)

- **Domain & Data Model** (`docs/specs/domain-and-data-model-design.md`) →
  `ai-architecture-design.md` (entity name references: AIDecision, EvidenceItem, CommonControl, CommonControlSatisfies, EvidenceItemSupports, etc.),
  `backend-architecture-design.md` (module entity ownership, DB topology),
  `packages/openapi/*.yaml` (schemas, enums, constraints),
  `axiom-spec-design.md` (§5 data model summary)

- **AI Architecture** (`docs/specs/ai-architecture-design.md`) →
  `domain-and-data-model-design.md` (`ai_content_metadata` jsonb, `ai_context_type` enum values),
  `backend-architecture-design.md` (River workers, Bedrock model assignments, pgvector scope, cross-service AIDecision pattern),
  `packages/openapi/*.yaml` (AI endpoints, ai_content_metadata schema, enum values),
  `infrastructure-design.md` (Bedrock IAM per-service, River DLQ alarms, AI observability dashboard, KMS provenance-signing key),
  `axiom-spec-design.md` (§6 AI summary)

- **Backend Architecture** (`docs/specs/backend-architecture-design.md`) →
  `packages/openapi/*.yaml` (endpoint distribution across service specs; new `frameworks.yaml` owns CommonControl/STRM endpoints),
  `infrastructure-design.md` (ECS task definitions, DB topology, IAM roles, provenance-signer service),
  `axiom-spec-design.md` (§7 tech stack summary)

- **Infrastructure Design** (`docs/specs/infrastructure-design.md`) →
  `docs/diagrams/axiom_infrastructure.py` (visual rendering of infrastructure),
  then regenerate `axiom_infrastructure.png` by running the script

- **OpenAPI common.yaml** (`packages/openapi/common.yaml`) →
  all 6 service `.yaml` files (`audit-core`, `doc-processing`, `frameworks`, `identity-service`, `reporting-service`, `workpaper-service`) via `$ref`

- **Design System** (`.impeccable.md`) →
  `mockups/` screens implement the design system

- **Specifications → Implementation Plan** (any spec change) →
  review `docs/superpowers/specs/implementation-plan-design.md` and the per-phase plans under `docs/superpowers/plans/` for impact

### Bidirectional Dependencies

The Domain Model (D) and AI Architecture (AI) have a **bidirectional** dependency. D defines entity structures that AI references; AI defines `ai_content_metadata` and enum values that go into D. Changes to either may require updates to the other.

### Terminal Nodes

`axiom-spec-design.md` is the **summary hub** — it aggregates summaries from all other specs but nothing depends on it. Update it last, after all upstream changes are complete.

`axiom_infrastructure.png` is a generated artifact. Never edit it directly; update the `.py` script and regenerate.

### Workflow

1. Make changes to the upstream document
2. Walk the cascade rules to identify all downstream artifacts
3. Update each downstream artifact
4. If `axiom_infrastructure.py` changed, run it to regenerate the `.png`
5. Update `axiom-spec-design.md` summary sections last
6. Review the implementation plan and any in-flight phase plan for impact

## Research

Background research in `docs/research/` — market, competitive, regulatory, data model, AI architecture, tech stack evaluation, integrations, legal/governance, pricing. Some research documents predate the compliance pivot and are flagged as stale inside `axiom-spec-design.md`; refresh is a known follow-up.
