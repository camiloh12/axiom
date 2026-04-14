# Axiom

An AI-native audit engagement platform for mid-market CPA and advisory firms.

## What It Is

Axiom replaces the fragmented CaseWare + DataSnipper + Excel toolstack with a unified platform covering:

- **Workpaper management** — GAAS-compliant test procedures and documentation
- **Trial balance analysis** — population-level testing, materiality, and sampling
- **Framework-agnostic evidence collection** — SOC 2, ISO 27001, HIPAA, and financial audits in one workspace
- **PBC request management** — client-facing portal for document uploads and request fulfillment
- **Regulatory-compliant archiving** — S3 Object Lock (WORM) for finalized reports
- **AI assistance** — evidence matching, review suggestions, and audit decision logging (AI suggests; humans decide)

Target firm size: 20–60 staff, running 30–100 engagements per year. Designed for time-to-first-engagement under one week with no implementation consultant.

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

Microservices behind an API Gateway (Go), decomposed by bounded context:

1. **API Gateway** — JWT verification, routing, rate limiting
2. **Identity Service** — auth, RBAC, firm/user/client management
3. **Audit Core** — engagements, controls, evidence, document requests, AI decisions
4. **Trial Balance Service** — trial balance data and population analysis
5. **Workpaper Service** — collaborative workpapers with Yjs real-time sync
6. **Reporting Service** — async report generation and S3 archiving
7. **Document Processing Service** (Python) — PDF extraction via pdfplumber + Tesseract

## Compliance Targets

SOC 2 Type 2, ISO 27001, ISO 42001 from day one. All AI decisions are logged for human-in-the-loop compliance.

## Specification Documents

- [axiom-spec-design.md](docs/specs/axiom-spec-design.md) — full product design and feature requirements
- [backend-architecture-design.md](docs/specs/backend-architecture-design.md) — service decomposition and monorepo structure
- [infrastructure-design.md](docs/specs/infrastructure-design.md) — AWS account structure, Terraform workspaces, CI/CD, observability
- [domain-and-data-model-design.md](docs/specs/domain-and-data-model-design.md) — domain model and data schema design
- [user-journey-mapping-skill-design.md](docs/specs/user-journey-mapping-skill-design.md) — user journey mapping skill design

## User Journeys

- [all-journeys.md](docs/user-journeys/all-journeys.md) — end-to-end user journeys across all roles and flows
