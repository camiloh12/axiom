# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Name

The working name for this application is **Axiom**.

## Project Status

This repository currently contains no implemented code. All guidance below reflects the intended architecture from the current specification documents.

## Specification Documents

- **Product spec:** `docs/specs/axiom-spec-design.md` — full product design, data model, and feature requirements
- **AI architecture:** `docs/specs/ai-architecture-design.md` — LLM provider, eight AI features, human-in-the-loop policy, AI content tracking, cost estimates
- **Backend architecture:** `docs/specs/backend-architecture-design.md` — supersedes the tech stack section of the product spec
- **Infrastructure:** `docs/specs/infrastructure-design.md` — AWS account structure, networking, Terraform workspaces, CI/CD, observability, security controls, cost estimates

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

Modular monolith (single Go binary) with internal packages per bounded context, plus a separate Python service for PDF extraction:

**Axiom API** (Go modular monolith):
- `internal/gateway` — Chi middleware: JWT verification, routing, rate limiting
- `internal/identity` — auth, RBAC, firm/user/client management, JWT issuance
- `internal/auditcore` — engagements, controls, evidence, document requests, AI decisions
- `internal/trialbalance` — trial balance data, population analysis (SQL-based)
- `internal/workpaper` — collaborative workpapers with Yjs real-time sync (WebSocket)
- `internal/reporting` — async report generation, S3 WORM archiving
- `internal/ai` — Bedrock client, prompt templates, embedding helpers
- `internal/platform` — DB, config, OTel, River, common middleware

**Document Processing** (Python) — stateless PDF extraction via pdfplumber + Tesseract

Modules communicate via Go interfaces (function calls), not HTTP. The only inter-service HTTP call is Axiom API → Document Processing for PDF extraction.

See `docs/specs/backend-architecture-design.md` for full module descriptions and monorepo structure.

## Key Domain Concepts

- **Engagement** — a single audit project with a defined framework, client, and lifecycle
- **Controls** — regulatory requirements being tested, mapped to framework criteria
- **Trial Balance** — financial source-of-truth for financial audit engagements
- **Workpapers / Test Papers** — documentation of test procedures and evidence
- **Client Hub** — client-facing portal for document uploads and request fulfillment
- **Population Analysis** — testing entire transaction datasets rather than samples

## User Roles

`FirmAdmin | Partner | Manager | Staff | EQReviewer | ClientAdmin | ClientUser | ViewOnly`

## Security & Compliance

Target from day one: **SOC 2 Type 2**, **ISO 27001**, **ISO 42001**.

- Encryption in transit and at rest (AWS KMS)
- RBAC + RLS; all access events in append-only audit log
- AI decisions logged via `AIDecision` table for human-in-the-loop compliance
- Finalized reports use S3 Object Lock (WORM) for regulatory immutability
- GDPR/CCPA compliant; client firm owns all engagement data

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
| MK | HTML Mockups (62 screens) | `mockups/journey-*/` |
| DS | Design System | `.impeccable.md` |

### Cascade Rules

When a document changes, update all listed downstream artifacts:

- **Research** (`docs/research/*.md`) →
  `axiom-spec-design.md` (§1-4, §8-9 reference research)

- **User Journeys** (`docs/user-journeys/all-journeys.md`) →
  `domain-and-data-model-design.md` (entity definitions tagged by journey, §12 traceability matrix),
  `ai-architecture-design.md` (§4 feature journey references),
  `axiom-spec-design.md` (§11 journey summaries, §12 innovative flows),
  `mockups/` (directories and screens map 1:1 to journey stages)

- **Domain & Data Model** (`docs/specs/domain-and-data-model-design.md`) →
  `ai-architecture-design.md` (entity name references: AIDecision, WorkpaperVersion, EvidenceItem, etc.),
  `backend-architecture-design.md` (§2-3 service entity ownership, §3 DB topology),
  `packages/openapi/*.yaml` (schemas, enums, constraints),
  `axiom-spec-design.md` (§5 data model summary)

- **AI Architecture** (`docs/specs/ai-architecture-design.md`) →
  `domain-and-data-model-design.md` (ai_content_metadata jsonb, ai_context_type enum values),
  `backend-architecture-design.md` (River workers, Bedrock model assignments, pgvector scope, cross-service AIDecision pattern),
  `packages/openapi/*.yaml` (AI endpoints, ai_content_metadata schema, enum values),
  `infrastructure-design.md` (Bedrock IAM per-service, River DLQ alarms, AI observability dashboard),
  `axiom-spec-design.md` (§6 AI summary)

- **Backend Architecture** (`docs/specs/backend-architecture-design.md`) →
  `packages/openapi/*.yaml` (endpoint distribution across service specs),
  `infrastructure-design.md` (ECS task definitions, DB topology, IAM roles),
  `axiom-spec-design.md` (§7 tech stack summary)

- **Infrastructure Design** (`docs/specs/infrastructure-design.md`) →
  `docs/diagrams/axiom_infrastructure.py` (visual rendering of infrastructure),
  then regenerate `axiom_infrastructure.png` by running the script

- **OpenAPI common.yaml** (`packages/openapi/common.yaml`) →
  all 6 service `.yaml` files (they `$ref` shared schemas)

- **Design System** (`.impeccable.md`) →
  `mockups/` (all 62 screens implement the design system)

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

## Research

Background research in `docs/research/` — market, competitive, regulatory, data model, AI architecture, tech stack evaluation, integrations, legal/governance, pricing.
