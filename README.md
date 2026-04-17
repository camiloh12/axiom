# Axiom

An AI-native compliance and assurance platform for mid-market CPA firms and compliance consultancies.

## What It Is

Axiom is a both-sided (auditor + auditee) AI-native platform for running SOC 1, SOC 2, ISO 27001, ISO 27701, ISO 42001, HIPAA, and PCI DSS engagements. It supplants the fragmented Drata/Vanta/Secureframe (auditee-only) and Agentive (auditor-only) toolstacks with a unified workspace covering:

- **Cross-framework evidence mapping** — STRM-graded (NIST Secure Controls Framework) crosswalk that maps one evidence artifact to every applicable framework requirement, with period-aware coverage and explicit partial-satisfaction gap lists (never a green checkmark on partial coverage)
- **Engagement lifecycle** — scoping → readiness → fieldwork → reporting → archiving; supports integrated multi-framework engagements (e.g., SOC 2 + ISO 27001 + ISO 27701 in one scope per Journey 11)
- **Client Hub / auditee GRC workspace** — continuous monitoring, policy library, evidence freshness dashboard, drift-triggered re-testing (Journey 12)
- **PBC document requests** — client-facing portal for document uploads and request fulfillment
- **Cryptographic AIDecision provenance** — every AI output signed at emission with an AWS KMS asymmetric key (`ECC_NIST_P256`, `SIGN_VERIFY`) and written to S3 Object Lock (WORM) for tamper-evident, publicly verifiable artifacts
- **ISO 42001-native HITL AI governance** — three-tier human-in-the-loop policy, `AIDecision` ledger, `ai_content_metadata` tracking; Axiom dogfoods the framework it sells

Target firm size: 20–60 professional staff, running 30–100 engagements per year. Designed for time-to-first-engagement under one week with no implementation consultant.

## Tech Stack

| Layer | Technology |
|---|---|
| Frontend | TypeScript, React (SPA on CloudFront) |
| Backend | Go (primary), Python (PDF extraction only) |
| API | REST + OpenAPI 3.1 |
| Database | PostgreSQL with RLS (sqlc + pgx/v5) |
| Background jobs | River (Postgres-backed) |
| Workflow engine | AWS Step Functions |
| AI inference | AWS Bedrock (Claude via VPC endpoint) |
| Infrastructure | AWS ECS Fargate + ALB + RDS Multi-AZ, Terraform |
| CI/CD | GitHub Actions with OIDC federation |
| Observability | CloudWatch + X-Ray via OpenTelemetry |

## Architecture

Modular monolith (single Go binary) with internal packages per bounded context, plus a separate Python service for PDF extraction and an isolated signer service for provenance:

1. **Gateway** (`internal/gateway`) — JWT verification, routing, rate limiting
2. **Identity** (`internal/identity`) — auth, RBAC, firm/user/client management
3. **Audit Core** (`internal/auditcore`) — engagements, controls, evidence, document requests, findings, management responses, AI decisions
4. **Frameworks** (`internal/frameworks`) — `CommonControl` graph, `FrameworkRequirement` catalog, NIST STRM-encoded `CommonControlSatisfies` edges, `EvidenceItemSupports` period-aware coverage, SCF/OSCAL/AICPA/CIS crosswalk import
5. **Workpaper** (`internal/workpaper`) — collaborative workpapers with Yjs real-time sync
6. **Reporting** (`internal/reporting`) — async report generation and S3 archiving
7. **AI** (`internal/ai`) — Bedrock client, prompt templates, embedding helpers
8. **Provenance** (`internal/provenance`) — cryptographic signing of AI outputs and evidence artifacts; runs as a separate ECS service with KMS `Sign`-only IAM surface
9. **Document Processing** (Python, separate service) — stateless PDF extraction via pdfplumber + Tesseract

## Compliance Targets

Axiom itself targets SOC 2 Type 2, ISO 27001, ISO 27701, and ISO 42001 from day one. All AI decisions are logged via the `AIDecision` table for human-in-the-loop accountability.

## Specification Documents

- [axiom-spec-design.md](docs/specs/axiom-spec-design.md) — product summary hub
- [domain-and-data-model-design.md](docs/specs/domain-and-data-model-design.md) — domain model and data schema
- [ai-architecture-design.md](docs/specs/ai-architecture-design.md) — 11 AI features, HITL policy, provenance
- [backend-architecture-design.md](docs/specs/backend-architecture-design.md) — service decomposition and monorepo structure
- [infrastructure-design.md](docs/specs/infrastructure-design.md) — AWS account structure, Terraform workspaces, CI/CD, observability
- [compliance-pivot-findings.md](docs/specs/compliance-pivot-findings.md) — research brief behind the compliance pivot

## User Journeys

- [all-journeys.md](docs/user-journeys/all-journeys.md) — end-to-end user journeys across all roles and flows
