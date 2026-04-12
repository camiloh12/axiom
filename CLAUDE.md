# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Name

The working name for this application is **Axiom**.

## Project Status

This repository currently contains no implemented code. All guidance below reflects the intended architecture from the current specification documents.

## Specification Documents

- **Product spec:** `docs/specs/axiom-spec-v2-design.md` — full product design, data model, and feature requirements
- **Backend architecture:** `docs/specs/backend-architecture-design.md` — supersedes the tech stack section of the product spec

## Tech Stack

| Layer | Technology |
|---|---|
| Frontend | TypeScript, React (SPA on CloudFront) |
| Backend | Go (primary), Python (PDF extraction only) |
| API | REST + OpenAPI 3.1 (oapi-codegen, openapi-typescript) |
| Database | PostgreSQL — `core_db` (RLS), plus separate DBs per service |
| DB access | sqlc + pgx/v5 (no ORM) |
| Background jobs | River (Postgres-backed, Go-native) |
| Workflow engine | AWS Step Functions |
| AI inference | AWS Bedrock (Claude via VPC endpoint) |
| Infrastructure | AWS ECS Fargate + ALB + RDS (Multi-AZ) |
| Email | AWS SES |

## Architecture Overview

Microservices decomposed by bounded context, all behind an API Gateway (Go):

1. **API Gateway** — JWT verification, routing, rate limiting
2. **Identity Service** — auth, RBAC, firm/user/client management, JWT issuance
3. **Audit Core** — engagements, controls, evidence, document requests, AI decisions
4. **Trial Balance Service** — trial balance data, population analysis (SQL-based)
5. **Workpaper Service** — collaborative workpapers with Yjs real-time sync (WebSocket)
6. **Reporting Service** — async report generation, S3 WORM archiving
7. **Document Processing Service** (Python) — PDF extraction via pdfplumber + Tesseract

See `docs/specs/backend-architecture-design.md` for full service decomposition and monorepo structure.

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

## Research

Background research in `docs/research/` — market, competitive, regulatory, data model, AI architecture, tech stack evaluation, integrations, legal/governance, pricing.
