# Axiom Product Specification
**Date:** April 15, 2026
**Status:** Implementation-Ready

---

**On the name:** "Axiom" represents a self-evident truth or foundational starting point for reasoning — fitting for a compliance and assurance platform whose core promise is auditor-defensible attestation: every evidence item, every control mapping, every AI decision traceable to a verifiable source of truth. It is a working name and can be changed before launch.

Alternatives under consideration for a future rename (not prioritized at this time):

| Name | Rationale |
|---|---|
| **Prism** | One piece of evidence refracted across many frameworks — a visual metaphor for cross-framework evidence mapping |
| **Tenet** | A principle held to be true; near-synonym for axiom, shorter |
| **Canon** | Authoritative body of rules / attested standard (note: potential trademark overlap with Canon Inc.) |
| **Veris** | Latin *verus* = true; attestation-native with modern startup phonetics |
| **Affirma** | Directly attestation-native — "to affirm as true" |
| **Throughline** | Evidence that runs through multiple frameworks; narrative framing |

---

## Table of Contents

1. [Product Overview](#1-product-overview)
2. [Competitive Positioning](#2-competitive-positioning)
3. [Pricing and Business Model](#3-pricing-and-business-model)
4. [Regulatory Compliance Requirements](#4-regulatory-compliance-requirements)
5. [Core Data Model](#5-core-data-model)
6. [AI Architecture](#6-ai-architecture)
7. [Technology Stack](#7-technology-stack)
8. [Integration Roadmap](#8-integration-roadmap)
9. [Legal and Data Governance](#9-legal-and-data-governance)
10. [Security Architecture](#10-security-architecture)
11. [User Journeys](#11-user-journeys)
12. [Flows Without Competitor Equivalent](#12-flows-without-competitor-equivalent)
13. [Out of Scope at Launch](#13-out-of-scope-at-launch)

---

## 1. Product Overview

> Research: [01-target-market.md](../research/01-target-market.md)

### What Axiom Is

Axiom is an AI-native audit engagement platform built for mid-market CPA and advisory firms doing both financial audits and compliance framework audits (SOC 2, ISO 27001, HIPAA) from a single subscription. It replaces the fragmented CaseWare + DataSnipper + Excel toolstack with a unified platform that handles workpaper management, trial balance analysis, framework-agnostic evidence collection, PBC request management, and regulatory-compliant archiving — all with an AI layer that eliminates the manual extraction, evidence-matching, and review-delay work that consumes 10–20 hours per auditor per week. Axiom is designed to be operational within one week of signup, without an implementation consultant.

### Target ICP

**Firm size:** 20–60 professional staff (partners, managers, seniors, staff auditors).

**Practice mix:** Mixed financial audit and compliance framework audit (SOC 2 Type II and/or ISO 27001 and/or HIPAA). Not compliance-only. Not internal audit.

**Current toolstack:** Running CaseWare (or CaseWare + DataSnipper) alongside Excel for trial balance work. Managing PBC requests via email or a bolt-on tool like Suralink or AuditDashboard. No unified AI layer.

**Geography:** US and Canada at launch. Primarily auditing private companies (AICPA nonissuer standards); may have a small number of PCAOB-registered engagements.

**Engagement volume:** 30–100 engagements per year.

**Decision-maker profile:** Managing partner or audit partner who has evaluated Fieldguide and was deterred by cost or implementation complexity, or who has not yet evaluated modern platforms and is actively experiencing the pain of their current toolstack.

### Primary and Secondary Differentiators

**Primary differentiator:** Axiom is the only mid-market-priced platform that supports financial audits (trial balance, GAAS workpapers, materiality, sampling) and compliance audits (SOC 2, ISO 27001, HIPAA) in the same engagement workspace, with a framework-agnostic evidence architecture that eliminates redundant evidence collection across frameworks. The underlying evidence chain (EvidenceItem → TestProcedure → FirmControlObjective → FrameworkRequirement) operates at all tiers. The explicit cross-framework mapping UI — showing which other frameworks an evidence item satisfies and enabling auditors to run integrated multi-framework engagements — is a Scale-tier feature.

**Secondary differentiator:** Self-serve onboarding with time-to-first-engagement under one week. No implementation consultant required for firms using standard AICPA, SOC 2, or ISO 27001 methodologies. This is not just a UX goal — it is a structural sales advantage that enables a product-led growth motion where Fieldguide cannot follow.

### What Axiom Does NOT Do at Launch

- Big Four or enterprise firm workflows (sales cycles too long, security review burden too high)
- Internal audit, SOX compliance, or enterprise GRC (AuditBoard's market)
- ESG or sustainability reporting
- White-labeling or reseller channels
- Custom AI model training per firm
- On-device or local model execution
- AI-autonomous engagement actions (AI suggests; humans decide — always)
- PCAOB-registered firm deep compliance automation beyond AS 1105 documentation requirements
- ERP write-back integrations (read-only at launch)
- Practice management billing integration (Karbon sync is read-status, not billing)

---

## 2. Competitive Positioning

> Research: [01-target-market.md](../research/01-target-market.md) · [02-competitive-differentiation.md](../research/02-competitive-differentiation.md)

### One-Paragraph Positioning vs. Fieldguide, Yak, and CaseWare+DataSnipper

Fieldguide is the correct choice for the top 100 US CPA firms with budget exceeding $50,000/year and appetite for a structured multi-week implementation program — it is not priced or designed for the firm that needs to run its first engagement next Monday. Yak is excellent for compliance-only boutiques doing SOC 2 and HIPAA, but it has no trial balance module and no financial audit support, making it irrelevant for any firm doing mixed-service work. CaseWare + DataSnipper is the de facto mid-market standard, but it is two products, two licenses, two logins, and no AI layer — firms pay $40,000–$67,000/year for a stack that still requires Excel as the connective tissue. Axiom occupies the gap that none of these tools fill: a platform with the breadth of Fieldguide (financial audit + compliance audit + AI) at mid-market pricing, with the self-serve simplicity that CaseWare and Fieldguide structurally cannot provide.

### The Switching Trigger

A mid-market firm evaluates Axiom when one or more of these events occurs:

1. **Busy-season failure:** A partner exits a January–April audit season having spent 15+ hours/week on manual PDF extraction, evidence chasing, and review delays. The pain is acute and the budget conversation is easy.
2. **The Fieldguide wall:** A firm requests a demo or quote from Fieldguide, receives either no response or a quote exceeding $40,000/year with a mandatory Accelerator program, and begins looking for alternatives.
3. **The DataSnipper renewal:** A firm receives its annual DataSnipper renewal ($10,000–$20,000/year) and asks why they are paying for a productivity layer on top of Excel rather than replacing Excel.
4. **Lost engagement:** A firm loses a prospective client because their document request process was too slow or their report delivery was delayed.

The in-app trial strategy exploits triggers 1 and 3: a firm can sign up, load a real trial balance, run an AI completeness review on a PBC document, and experience the time savings within 30 minutes — before committing to any payment.

### Why Fieldguide Won't Easily Respond Downmarket

Fieldguide has raised $125M in total funding at a $700M valuation (Series C, February 2026, led by Goldman Sachs). They serve 40+ of the top 100 US CPA firms. Their strategic incentive is to expand upmarket — into the Big Four, into ESG reporting, into enterprise GRC — not to launch a self-serve product at one-third of their current price point. A downmarket product would:

- Directly undercut their current customer base's willingness to pay (signaling that enterprise pricing was inflated)
- Require a separate product motion (PLG vs. enterprise sales) that conflicts with their Accelerator onboarding model
- Attract a buyer segment with shorter sales cycles, lower ACVs, and higher support burden per dollar

Enterprise SaaS companies at Fieldguide's stage rarely segment down successfully. The structural incentives run in the opposite direction.

---

## 3. Pricing and Business Model

> Research: [09-pricing.md](../research/09-pricing.md)

### Recommended Tiers

**Growth — $1,200/month ($14,400/year, or $12,000/year billed annually)**

Best for: 10–40 staff firms, 20–50 engagements/year, migrating off CaseWare + DataSnipper.

- Unlimited users
- Up to 35 active engagements/year included
- Financial audit: trial balance import, workpaper management, sampling calculators, materiality
- Compliance audit: SOC 2 Type I/II, HIPAA, ISO 27001
- Standard AI: document extraction, evidence completeness review, workpaper draft assist
- Pre-built methodology templates (AICPA/GAAS, SOC 2 TSC 2017, ISO 27001:2022)
- Client PBC request portal (tokenized upload links, no client login required for simple upload)
- Standard email support + in-app help
- Overage: $350/engagement beyond included 35

**Scale — $2,400/month ($28,800/year, or $24,000/year billed annually)**

Best for: 30–100 staff firms, 50–150 engagements/year, running multi-framework audits.

- Everything in Growth
- Up to 100 active engagements/year included
- Cross-framework evidence mapping (one evidence item satisfies multiple frameworks simultaneously)
- AI-assisted risk assessment and control gap analysis
- Advanced analytics dashboard (engagement cycle time, staff utilization, review bottleneck identification)
- Multi-entity / group audit support
- Custom methodology template editor
- Priority support + dedicated onboarding call
- Overage: $250/engagement beyond included 100

**Enterprise — Custom ($50,000–$120,000/year, negotiated)**

Best for: 80–200+ staff firms, 150+ engagements/year, requiring enterprise security and SLAs.

- Everything in Scale, unlimited engagements
- Dedicated Customer Success Manager
- 99.9% uptime SLA
- Signed HIPAA BAA between Axiom and the audit firm (required for the firm's compliance documentation; Growth/Scale tiers receive the platform HIPAA BAA on request, but dedicated BAA negotiation and compliance documentation are included in Enterprise)
- SAML/SSO integration
- Audit trail export and long-term data retention configuration
- Custom integrations (ERP, practice management, GL systems)
- Security review package (SOC 2 Type II report, penetration test summary, data residency documentation)
- Negotiated multi-year contracts; EU/APAC data residency option

**Annual billing incentive:** 2 months free (16.7% discount) for annual prepayment. Monthly billing available at Growth tier only. Scale and Enterprise require annual contracts.

**Engagement definition:** One engagement = one client entity, one primary framework, one audit period. A rollforward to the following year counts as a new engagement. Integrated audits (e.g., SOC 2 + ISO 27001 in one engagement) count as one engagement. This definition is stated in the terms of service and enforced in the product UI.

### Self-Serve vs. Sales-Assisted Thresholds

| Tier | Signup Flow | Sales Touch |
|---|---|---|
| Growth | Self-serve: credit card or invoice, 14-day trial, no mandatory call | Optional: automated email Day 3, human call offer Day 10 |
| Scale | Self-serve trial start (14 days), then inside sales to close | Inside sales rep follows up Day 3 and Day 10; 30-min call |
| Enterprise | Request demo → assigned AE → custom quote | Full sales-led, 4–8 week cycle, security review package |

### Trial Strategy

14-day full-feature trial, no credit card required. Trial requires a business email domain and brief intake form (firm name, staff count, primary audit types). One pre-loaded engagement template on signup (SOC 2 Type II or GAAS Financial Audit, selectable). Trial workspace is locked (not deleted) at expiry until a plan is selected. AI-driven in-app progress guidance toward completing the first engagement setup.

No permanent free tier. The trust and professional-perception risk of a free tier outweighs the acquisition benefit for a platform handling sensitive financial data.

### Revenue Model Summary

| Milestone | Firm Mix | Projected ARR |
|---|---|---|
| **50 firms** | 30 Growth, 15 Scale, 5 Enterprise | ~$1.16M ARR |
| **200 firms** | 90 Growth, 80 Scale, 30 Enterprise | ~$5.4M ARR |
| **500 firms** | 175 Growth, 225 Scale, 100 Enterprise | ~$15M ARR |

With overage revenue (5–15% of ACV) and NRR of 110–120%, realistic ARR at 500 firms is $16–18M. Gross churn target: 8–12% annually. The primary churn risk is not product dissatisfaction — it is firm-level business events (mergers, economic contraction, partner retirement removing the champion). Historical audit software has very low churn because engagement history and compliance retention obligations create a strong switching cost: migrating five years of engagement files to a new platform costs $40,000–60,000.

---

## 4. Regulatory Compliance Requirements

> Research: [03-regulatory-standards.md](../research/03-regulatory-standards.md)

### Engagement Type × Standard × Platform Requirements

| Requirement | Financial Audit (Private) | SOC 1 / SOC 2 | Financial Audit (Public / PCAOB) | ISO 27001:2022 | HIPAA |
|---|---|---|---|---|---|
| **Governing standard** | AICPA AU-C (SAS 122+) | AICPA AT-C 105/205/320 (SSAE 18) | PCAOB AS series | ISO/IEC 27001:2022 | HIPAA Security Rule (45 CFR §§164.302–318) |
| **Workpaper assembly deadline** | 60 days after report release | 60 days after report issuance | 45 days after report release | N/A | N/A |
| **Retention period** | 5 years from report date | 5 years from report date | 7 years from report date | Ongoing ISMS records | 6 years from date of creation or last effective date |
| **Sign-off hierarchy** | Engagement partner + reviewer | Partner + EQR (SQMS 2 where applicable) | Partner + mandatory EQR (AS 1220) | Certification body | N/A (agreed-upon procedures or internal audit) |
| **Period coverage** | Year-end or stub period | Full examination period (Type II: 6 or 12 months minimum) | Year-end | Certification cycle (initial + annual surveillance + triennial full) | Ongoing compliance point-in-time |
| **AI documentation required** | Best practice (PCAOB guidance used as reference) | Best practice | Mandatory (PCAOB AS 1105, eff. Dec 15, 2025) | N/A | N/A |
| **Immutable lock trigger** | After 60-day assembly window closes | After 60-day assembly window closes | After 45-day assembly window closes | After certification issued | N/A |
| **Quality management** | SQMS 1 (eff. Dec 15, 2025) | SQMS 1 (eff. Dec 15, 2025) | PCAOB QC 1000 | ISO internal audit requirements | N/A |

### SQMS 1 Implications for Platform Workflow

SQMS 1 (effective December 15, 2025) replaces the prior QC 10 quality control standard. It is not a policy document — it changes the workflow of engagements in ways the platform must support:

**Client/engagement acceptance:** When creating a new engagement, the platform requires completion of a `ClientAcceptance` record documenting quality risks identified and the firm's responses to those risks. The Planning → Fieldwork transition is blocked until this record is signed off by a partner. This is not optional and cannot be bypassed.

**Communicating quality policies:** During engagement setup, applicable firm quality policies are surfaced inline (not buried in a separate policy document the team never reads). Firms maintain their quality policies in the platform, and the system links relevant policies to each engagement type.

**Engagement Quality Review (SQMS 2):** For engagements where EQR is required (mandatory for all PCAOB engagements under AS 1220; applicable for higher-risk nonissuer engagements under SQMS 1), the platform enforces:
- EQR reviewer must be assigned at engagement setup
- EQR reviewer cannot be an `EngagementTeamMember` on the same engagement (system-enforced via role conflict check)
- Review → Reporting transition is blocked until `EngagementQualityReview.status = Complete` and `signed_off_at` is populated
- The EQR record captures: reviewer identity, independence confirmation date, scope notes, and conclusion

**Partner involvement evidence:** The AuditLog creates a timestamped record of engagement partner activity throughout the engagement. Partner sign-off events are first-class log entries. This provides the "substantial and meaningful involvement" audit trail required under AU-C 220.

**Annual internal inspection:** The platform supports inspection workflows — sampling prior engagements by stratified random selection, documenting inspection findings, tracking remediation. This is a Growth-tier feature at launch (partner-managed), with enhanced inspection analytics on Scale.

### PCAOB AS 1105 AI Documentation Requirement (Effective December 15, 2025)

PCAOB AS 1105 (as amended effective December 15, 2025) requires that for any audit procedure using technology-based analysis tools (including AI):

1. The auditor must evaluate the reliability of electronic information used as evidence
2. IT general controls over source systems must be tested
3. Technology-assisted procedures must be documented as meeting their intended purpose
4. All flagged transactions or balances must be investigated and documented

**Platform implementation:** Every AI action that affects audit content creates an `AIDecision` record (defined in Section 5). For PCAOB engagements, `AIDecision` records are included in the immutable engagement archive. The `AIDecision` schema maps directly to AS 1105 documentation requirements:

| AS 1105 Requirement | AIDecision Field |
|---|---|
| What technology procedure was performed | `context_type` + `context_id` |
| What model/tool was used | `model_id` |
| What the AI determined | `suggested_value` + `raw_output` |
| Who reviewed the AI output | `reviewed_by_id` |
| What the auditor decided (accepted/modified/rejected) | `review_action` + `accepted_value` |
| When review occurred | `reviewed_at` |

AI outputs alone never constitute sufficient audit evidence. Every Tier 2 AI action (defined in Section 6) requires explicit auditor review and documented sign-off before becoming part of the audit file.

**AS 1215 future-readiness (effective December 15, 2026):** The forthcoming AS 1215 amendment will standardize documentation structure to facilitate AI and data analytics in audit workflows. Axiom's structured, machine-readable workpaper format (content stored as typed jsonb, not free-form text) anticipates this requirement. No retroactive changes will be needed to existing engagements.

### Framework Version Management

The `Framework` table treats each version as a separate row. ISO 27001 2013 and ISO 27001 2022 are separate entries. SOC 2 TSC 2017 and any future revision will be separate entries.

**How versions are managed:**
- Each engagement records `framework_version_id` at creation time — the version in effect at that time
- When a new framework version is published, a new `Framework` row is added with its `effective_date`; no existing engagement references are changed
- Evidence linked to a control requirement must reference the specific version of the criterion satisfied (e.g., "SOC 2 TSC 2017, CC6.1" — not just "CC6.1")
- The platform prevents applying deprecated version control mappings to new engagements after the transition deadline (e.g., ISO 27001:2013 was deprecated October 31, 2025)
- When a firm opens a new engagement, the current active version for each framework is the default; firms can select a prior version only with explicit override and documentation reason (for re-engagements or comparatives)

---

## 5. Core Data Model

> **Full specification:** [Domain and Data Model Design](domain-and-data-model-design.md) — complete domain model (bounded contexts, aggregates, invariants), data model (table definitions, column types, constraints, indexes), and journey-to-entity traceability. What follows is a summary.

### Bounded Contexts

The domain is organized into 7 bounded contexts and 3 cross-cutting concerns, derived from the [user journeys](../user-journeys/all-journeys.md):

| # | Context | Key Entities | Module |
|---|---|---|---|
| 1 | Firm Identity | Firm, User, Client, Invitation | `internal/identity` |
| 2 | Regulatory Framework | Framework, FrameworkRequirement, ControlObjectiveLibrary | `internal/auditcore` |
| 3 | Firm Methodology | MethodologyTemplate, FirmControlObjective, template items | `internal/identity` |
| 4 | Audit Core | Engagement, Control, TestProcedure, EvidenceItem, EvidenceLink, DocumentRequest, ClientAcceptance, EngagementQualityReview | `internal/auditcore` |
| 5 | Trial Balance | TrialBalance, TrialBalanceAccount, TrialBalanceAdjustment | `internal/trialbalance` |
| 6 | Workpaper Authoring | Workpaper, WorkpaperVersion, ReviewNote | `internal/workpaper` |
| 7 | Reporting | Report, ReportVersion | `internal/reporting` |
| — | Cross-cutting | AIDecision, AuditLog, Notification | `internal/auditcore` |

**Total entities: 33** in a single PostgreSQL database (`axiom_db`) with RLS.

### Cross-Framework Evidence Chain

The core architectural differentiator — one evidence upload satisfies all mapped framework requirements simultaneously:

```
EvidenceItem → EvidenceLink → TestProcedure → Control
  → FirmControlObjective → FirmControlObjectiveMapping
    → FrameworkRequirement (SOC 2 CC6.1)
    → FrameworkRequirement (ISO 27001 A.8.3)
    → FrameworkRequirement (HIPAA §164.312(a)(1))
```

This chain requires ACID transactions, which is why Audit Core (Context 4) is a single bounded context with a shared database.

### Engagement Lifecycle State Machine

```
Planning ──[ClientAcceptance signed by Partner]──► Fieldwork
Fieldwork ──[All Controls: Complete or Exception]──► Review
Review ──[All ReviewNotes resolved + EQR signed off where applicable]──► Reporting
Reporting ──[Report.status = Issued]──► Finalized
Finalized ──[System: assembly_deadline elapsed]──► Archived (IMMUTABLE)
```

Reverse paths exist for exceptional cases (scope change, additional procedures, significant post-reporting issues). Once Finalized, no content can be modified — addenda only (AU-C 230, PCAOB AS 1215).

### Entities Added by Journey Analysis

The following entities were identified through the user journey analysis and are not present in earlier versions of this spec:

- **Invitation** — magic link onboarding for staff (Journey 1, 2)
- **TemplateDocumentRequest** — pre-drafted PBC requests within methodology templates (Journey 7)
- **ReviewNote** — structured, immutable review feedback on workpapers (Journey 6)
- **EQRFinding** — individual findings within an engagement quality review (Journey 10)
- **ClientHubToken** — tokenized no-login access for client document uploads (Journey 7, 8)
- **DelegationToken** — single-request scoped delegation for client contacts (Journey 8)
- **ColumnMappingProfile** — saved TB import configurations per accounting system (Journey 4)
- **Notification** — in-app and email delivery with deep links (Journey 2, 5, 6, 7)

### Multi-Tenancy

`axiom_db` uses PostgreSQL RLS with `firm_id` on all tenant-scoped tables. Three authorization dimensions: firm isolation (RLS), engagement team membership (point lookup), and client user scoping (engagement-level invitation).

See [Domain and Data Model Design](domain-and-data-model-design.md) for complete attribute definitions, data types, constraints, indexes, enum types, rollforward behavior, and the journey-to-entity traceability matrix.

---

## 6. AI Architecture

> **Full specification:** [AI Architecture Design](ai-architecture-design.md) — LLM provider decision, vector database, all eight AI features with model assignments and human review gates, AI content tracking, and cost estimates. What follows is a summary.

### LLM Provider and Vector Database

**AWS Bedrock with Claude Haiku and Sonnet**, accessed via VPC endpoint (PrivateLink). IAM-based auth, single AWS BAA, CloudWatch-native metrics. No external API keys, no additional sub-processor.

**pgvector** (PostgreSQL extension) for embedding storage at launch. Migration path to Qdrant at 5–10M vectors.

### Eight AI Features

| # | Feature | Model | Tier | Journeys |
|---|---------|-------|------|----------|
| 1 | Document completeness review | Sonnet | 2 | 7, 8 |
| 2 | Control mapping (cross-framework) | Haiku | 2 | 3 |
| 3 | Trial balance account mapping | Haiku | 2 | 4 |
| 4 | Workpaper narrative draft | Sonnet | 2 | 5 |
| 5 | Evidence link suggestion | Haiku | 2 | 5, 7 |
| 6 | Risk category suggestion | Sonnet | 2 | 3 |
| 7 | Trial balance anomaly detection | Haiku | 1 (Tier 2 for PCAOB) | 4 |
| 8 | Report section draft | Sonnet | 2 | 9 |

Every Tier 2 feature creates an `AIDecision` record and requires explicit human review before affecting the audit file. Feature 7 elevates to Tier 2 behavior for PCAOB engagements per AS 1105 documentation requirements. See [ai-architecture-design.md](ai-architecture-design.md) for full feature definitions including inputs, process, outputs, failure modes, and human review gates.

### Human-in-the-Loop Policy

**Tier 1 — Fully Automated:** Text extraction, embedding generation, notification-only flags, overdue reminders, analytics computation, anomaly detection (nonissuer engagements).

**Tier 2 — Human Approval Required:** All eight AI features above (except Feature 7 on nonissuer engagements). No audit file content changes without explicit auditor action.

**Tier 3 — Human-Initiated Only:** Engagement status transitions, control conclusions, exception documentation, sign-offs, report issuance, client acceptance, and any action constituting professional judgment.

### AI Content Tracking

AI-drafted content (Features 4 and 8) is tracked at the **section level**, not as a document-level boolean. Each AI-generated section records: origin timestamp, whether it was human-edited, editor identity, and a modification ratio (Levenshtein distance between AI output and current text). The advancement gate requires all AI-generated sections to have at least one human edit. Sections with <5% modification trigger a soft confirmation gate. EQR reviewers and managers see per-workpaper and engagement-wide AI edit substantiveness summaries. See [ai-architecture-design.md § 5](ai-architecture-design.md#5-ai-content-tracking) for the full data model and gate logic.

### Per-Engagement AI Cost Estimate

~$5.80 per SOC 2 engagement, ~$6.15 per financial audit engagement at on-demand Bedrock pricing. With prompt caching: **$3–5 per SOC 2**, **$4–6 per financial audit**. At 100 engagements/year, platform AI costs are $300–$600/year — absorbed into the subscription. See [ai-architecture-design.md § 6](ai-architecture-design.md#6-per-engagement-ai-cost-estimate) for the full breakdown.

### What AI Does NOT Do at Launch

Fine-tuned models per firm, on-device execution, multi-agent orchestration, AI-generated audit opinions or professional conclusions, autonomous multi-step actions without human review gates, auto-finalization of any audit content

---

## 7. Technology Stack

### Frontend

**TypeScript / React SPA.** Component library: Shadcn/ui (accessible, customizable, no licensing overhead). State management: TanStack Query for server state; Zustand for local UI state. API types are generated from the OpenAPI specs in `packages/openapi/` via `openapi-typescript` — a spec change automatically regenerates the client on the next build.

### Backend: Go Modular Monolith + Python PDF Service

**Decision: Go as the primary backend language, structured as a modular monolith. Python retained for PDF extraction.**

Go was chosen for its compile-time type safety, lean container images (30–50MB per Fargate task), and strong fit for compliance SaaS — Workiva, an established player in this space, uses Go for their REST services. The backend is a single Go binary organized into internal packages by bounded context (identity, audit core, trial balance, workpaper, reporting). Modules communicate via Go interfaces, not HTTP — this provides ACID transactions across the full evidence chain and eliminates distributed systems overhead for a solo developer + AI agent team.

Python is retained as a single stateless service for PDF extraction. `pdfplumber` handles complex, multi-column, scanned audit documents better than any Go library. The polyglot cost is contained — one endpoint, one job, no shared state.

**Full module descriptions, database design, Go and Python tech stack choices, and inter-module communication patterns are specified in [`backend-architecture-design.md`](./backend-architecture-design.md).**

**Monorepo structure:**
```
apps/
  axiom-api/        — Go: Modular monolith (single binary)
    internal/
      gateway/      — Chi middleware: JWT verification, routing, rate limiting
      identity/     — Auth, RBAC, firm/user/client, templates
      auditcore/    — Engagements, controls, evidence, AI decisions
      trialbalance/ — Trial balance, population analysis
      workpaper/    — Workpapers, Yjs collaboration (WebSocket)
      reporting/    — Report generation, S3 archival
      ai/           — Bedrock client, prompt templates
      platform/     — DB, config, OTel, River, common middleware
  doc-processing/   — Python: PDF extraction only
packages/
  openapi/          — OpenAPI 3.1 specs organized by module (source of truth)
```

Turborepo manages the monorepo with build caching.

### API Layer: REST + OpenAPI

**Decision: REST with OpenAPI 3.1 for all services. Hasura is rejected.**

**Why REST with OpenAPI:** REST decouples the frontend from the backend language. OpenAPI enables future public API exposure (webhooks, partner integrations) without adding a separate layer.

Each module defines its API contract as an OpenAPI 3.1 spec in `packages/openapi/`. `oapi-codegen` generates typed Go server interfaces from the spec. `openapi-typescript` generates typed fetch clients for the React frontend. Authorization is enforced as composable Go middleware:

- `WithFirmIsolation` — reads `firm_id` from gateway-injected headers, sets Postgres session variable for RLS
- `WithEngagementAccess` — verifies `EngagementTeamMember` record exists for the requested engagement
- `WithClientScoping` — for `ClientUser` roles, filters to invited engagements only

**Why Hasura is rejected:** The evidence authorization chain requires a five-level relationship traversal (`EvidenceItem → EvidenceLink → TestProcedure → Control → Engagement → EngagementTeamMember`). Hasura v2's `_exists` permission predicates materialize this as an `IN (...)` query that degrades as `EngagementTeamMember` grows. Hasura permissions are YAML metadata — not code, not testable with standard unit tests, difficult to reason about as authorization rules evolve. The auto-generated CRUD value Hasura provides is not justified for a platform where most queries are purpose-built for specific UI requirements.

**Public API (year 2+):** The OpenAPI specs already define the contract. Exposing selected endpoints publicly requires adding authentication scopes and rate limiting at the API Gateway — no rewrite of business logic.

### Database: PostgreSQL with RLS

PostgreSQL is the sole persistent data store. One RDS instance hosts a single database (`axiom_db`) with row-level security (RLS) on all tenant-scoped tables for multi-tenancy. Each module owns specific tables and accesses them via its own sqlc queries; cross-module data is accessed through Go service interfaces, not direct table queries. Database access uses `sqlc` + `pgx/v5` (type-safe SQL generation from plain SQL query files) with `golang-migrate` for schema migrations. PgBouncer for connection pooling in transaction mode.

**pgvector** extension enabled on `axiom_db` for embedding storage (Section 6).

**Workpaper content** stored as typed jsonb in `Workpaper.content`. The jsonb structure supports rich text (ProseMirror document format), embedded tables, formula references, and metadata. A dedicated document store is not needed at launch scale.

### Workflow Engine: Two-Tier (River + Step Functions)

**Tier 1 — River (PostgreSQL-based job queue) for background jobs:**

All fire-and-forget background work uses River, a Go-native Postgres-backed job queue. Zero additional infrastructure — it uses the existing database. One River instance serves all modules. Jobs are durable (WAL-backed), support retry with exponential backoff, and have dead-letter queues.

River jobs (running within the Axiom API against `axiom_db`):
- `document.extract` — PDF extraction via Python service
- `document.embed` — embedding generation and pgvector indexing
- `ai.completeness-check` — per document upload
- `ai.nightly-sweep` — engagement-level completeness review
- `ai.batch-control-mapping` — nightly, for new engagements
- `email.notification` — all transactional emails

**Tier 2 — AWS Step Functions Standard Workflows for the engagement lifecycle state machine:**

The engagement lifecycle — its state machine with explicit guard conditions, time-based Finalized → Archived transition, and auditability requirement — maps directly to Step Functions Standard Workflows. Standard Workflows are designed for long-running processes (days to years) with durable state. A `Wait` state keyed to the computed archival timestamp handles the Finalized → Archived transition reliably regardless of service restarts. The DocumentRequest reminder sequence is a `Task` → `Wait` → `Task` loop.

Step Functions is already within the AWS VPC, covered by the AWS BAA, and natively integrated with CloudWatch and X-Ray for execution history and distributed tracing. No additional vendor, no additional sub-processor. Cost at this scale is negligible ($0.025 per 1,000 state transitions; 100 engagements × ~20 transitions = cents per year).

State machines are defined in Amazon States Language (ASL) and invoked from the Audit Core service via the AWS SDK for Go. Execution history is persisted and queryable by Step Functions — engagement lifecycle transitions are auditable without a separate event store.

Step Functions state machines:
- `EngagementLifecycleStateMachine` — state machine, guards as ECS task invocations, `Wait` state for scheduled Archived transition
- `DocumentRequestReminderStateMachine` — `Task` (send) → `Wait` (7 days) → `Task` (send again) → `Wait` → `Task` (escalate)

### Real-Time Collaboration: Two-Tier by Data Type

**Workpaper rich text (TipTap + Yjs):** TipTap (ProseMirror-based editor) with Yjs CRDT for real-time co-editing of workpaper narratives. Yjs handles keystroke-level text concurrency naturally. Each save (triggered by idle timeout or explicit save) creates a `WorkpaperVersion` record. AI-generated content is tracked at the section level via `ai_content_metadata` on `WorkpaperVersion` (see [ai-architecture-design.md § 5](ai-architecture-design.md#5-ai-content-tracking)), satisfying PCAOB's requirement to distinguish AI-generated from auditor-authored content and providing EQR reviewers with edit substantiveness visibility.

**Structured audit data (server-authoritative with field locking):** Trial balance cells, AI decision acceptance, control status, and engagement state transitions use server-authoritative operations with optimistic UI updates. For high-stakes operations (AI decision review, control conclusion, workpaper sign-off), field-level locking prevents concurrent conflicting professional judgments. If Auditor A opens an AI decision for review, the server marks it "in review by User A." Auditor B sees a "locked" indicator. Conflicts (network partitions, duplicate submissions) are resolved by the server: first committed action wins; the AuditLog records both attempts; the second user is notified.

CRDTs (Yjs, Automerge) are explicitly not used for structured audit data. The automatic merge behavior of CRDTs is incompatible with the regulatory requirement that conflicting professional judgments produce an unambiguous, auditable outcome.

### Spreadsheet Component: AG Grid Community + HyperFormula

**Decision: AG Grid Community (MIT) + `hyperformula` (MIT) at launch. Univer evaluated at 6–12 months.**

AG Grid Community handles 200–10,000 rows with virtualization, has a large React integration ecosystem, and is free. `hyperformula` (the same formula engine used internally by Handsontable, now independently MIT-licensed) provides 400+ Excel-compatible functions. Cell-level comments are implemented as a custom cell renderer.

**Why AG Grid Enterprise is rejected:** The $999/developer cost is not justified when the required features (formula engine, basic export) are achievable with Community + open-source libraries.

**Why Handsontable is rejected:** The per-developer/per-year licensing model creates a recurring cost that grows with the team. The formula engine advantage is neutralized by adding `hyperformula` to AG Grid. Collaborative editing story is no better.

**Univer (Apache 2.0):** Provides native collaborative editing and a full formula engine, but is a newer project with limited production case studies outside of its creator. Evaluate via a 3-month spike at 6–12 months post-launch, after the collaborative editing requirement is confirmed by customer feedback. If the spike validates Univer's stability, migrate the spreadsheet component.

### Infrastructure: AWS + ECS Fargate

**Primary cloud:** AWS. ECS Fargate for container orchestration (serverless — no node management or control plane upgrades). Infrastructure-as-code via Terraform.

**Key AWS services:**
- **ECS Fargate** — Container orchestration for two services: Axiom API (Go modular monolith) and Document Processing (Python). ECS Service Connect provides DNS for API-to-doc-processing communication. Axiom API scales on the maximum of CPU utilization and active WebSocket connection count.
- **ALB** — Application Load Balancer for TLS termination in front of the API Gateway
- **RDS PostgreSQL** — Single Multi-AZ instance hosting one database (`axiom_db`) with RLS; pgvector extension enabled
- **S3** — Evidence file storage; Object Lock enabled for finalized engagements (see Section 10)
- **CloudFront** — CDN for the React SPA
- **SES** — Transactional email (document request notifications, client invitations, review alerts)
- **Secrets Manager** — API keys, database credentials, OAuth tokens (automatic 30-day rotation for RDS credentials)
- **CloudWatch + X-Ray** — Logging, monitoring, distributed tracing (via OpenTelemetry Go SDK)
- **AWS WAF** — Web Application Firewall on ALB and CloudFront (OWASP core rules, rate limiting, geo-restriction)
- **GuardDuty** — Threat detection across ECS, S3, and RDS
- **AWS Config** — Continuous infrastructure compliance monitoring

**Full AWS account structure, VPC design, Terraform workspace segmentation, CI/CD pipeline, security controls, observability configuration, and cost estimates are specified in [`infrastructure-design.md`](./infrastructure-design.md).**

---

## 8. Integration Roadmap

> Research: [07-integrations.md](../research/07-integrations.md)

### Launch (Must-Have)

These are prerequisites — without them, the target ICP cannot run an engagement.

| Integration | Type | What It Enables |
|---|---|---|
| **Transactional email (SES)** | Outbound notifications | Document request notifications, client portal invitations, review alerts. Without this, the PBC workflow does not function. |
| **Direct file upload (native UX)** | Evidence ingestion | Clients upload evidence via tokenized link, no login required. Handles CSV, Excel, PDF, ZIP with bulk upload and folder-structure recognition. Covers 80%+ of evidence delivery cases with no external dependency. |
| **CSV/Excel trial balance import** | Financial data | QBO, NetSuite, Sage, and Xero all export trial balances as CSV/Excel. A well-designed importer covering common format variations handles financial audit MVP. |
| **Microsoft / Google SSO (SAML/OAuth)** | Identity | Mid-market firms authenticate via Microsoft or Google. Separate Axiom credentials for every staff member is an onboarding barrier. SAML for Enterprise; OAuth for Growth/Scale. |

### First 6 Months Post-Launch

| Integration | Type | Priority Rationale |
|---|---|---|
| **Drata Audit Hub API** | Evidence ingestion | Highest-value compliance audit integration. Drata is the SOC 2 compliance automation market leader. A Drata → Axiom connection delivers structured, control-tagged evidence directly into engagements, replacing hours of manual organization. Fieldguide already has this; Axiom needs it to compete in the SOC 2 segment. |
| **Vanta evidence export** | Evidence ingestion | Analogous to Drata; 375+ source integrations, large installed base. Drata + Vanta covers the majority of well-prepared SOC 2 audit clients. |
| **Google Drive OAuth** | Cloud storage | Highest-coverage cloud storage integration. Many mid-market clients organize evidence in Drive. Clean OAuth implementation. |
| **Karbon API** | Practice management | Most important practice management integration for the ICP. Karbon is dominant in 20–200 staff US/Canada accounting firms. Engagement sync eliminates double-entry. A Karbon integration is a clear differentiator over Fieldguide, which does not offer one. |
| **QuickBooks Online direct connector** | Accounting data | Most common accounting system for clients audited by the ICP. Direct QBO API connection (pull trial balance + journal entries on authorization) removes the "ask client to export a file" step and ensures point-in-time completeness. |
| **Merge.dev HRIS layer** | HR data | Enables employee list pull for access review testing (CC6.2, CC6.3) in SOC 2 and HIPAA engagements. Covers Gusto, Rippling, BambooHR, Workday with one integration rather than six. |

### 6–18 Months Post-Launch

| Integration | Type | Notes |
|---|---|---|
| **SharePoint / Microsoft 365** | Cloud storage | Strategically important (M365 is dominant in firm IT stacks) but technically complex — Azure AD app registration, tenant admin consent, higher maintenance burden. Build after Google Drive is proven. |
| **Dropbox** | Cloud storage | Simpler OAuth than SharePoint. Lower priority than Google Drive by installed base size. |
| **NetSuite direct connector** | Accounting data | High-complexity integration; NetSuite APIs are heavily customized. Prioritize after QBO is proven. Clients at this size can provide file exports. |
| **Sage Intacct direct connector** | Accounting data | Well-documented XML API. Sage Intacct is dominant for SaaS-model companies and nonprofits — a significant SOC 2 audit client segment. |
| **Xero direct connector** | Accounting data | Clean REST API with native trial balance endpoint. Canadian market and smaller clients. |
| **TaxDome API** | Practice management | Relevant for mixed tax/audit firms. Lower priority than Karbon for audit-focused ICP. |
| **Public API + webhooks (Axiom-outbound)** | Platform extensibility | Enables Karbon workflows, Zapier/Make automation. Required before any firm can build custom workflows on top of Axiom. |
| **Sprinto / Hyperproof evidence export** | Evidence ingestion | Analogous to Drata/Vanta; smaller installed base. Add after Drata + Vanta are live. |
| **Box** | Cloud storage | Legitimate enterprise presence in regulated industries; lower footprint than Dropbox in the ICP's client population. |

### Key Insight: Drata Audit Hub Over Direct Infrastructure Integrations

Axiom is on the **auditor side**, not the auditee side. Axiom does not need to build direct integrations into AWS, Okta, GitHub, or any other client infrastructure system. That is the auditee platform's problem — Drata, Vanta, and Sprinto already have 200–375+ integrations into those systems. Axiom's role is to receive evidence packages from auditee platforms (via Drata Audit Hub, Vanta export) or via direct file upload.

The Drata Audit Hub integration is therefore more valuable than any single infrastructure integration (AWS IAM, Okta, GitHub combined), because it abstracts the entire auditee stack through one API connection. This is not obvious from the original spec's integration list — the original spec listed AWS, Dropbox, Box, Google Drive, and O365 as the integrations, without recognizing that the compliance automation platforms (Drata, Vanta) are the correct abstraction layer. Direct AWS/Okta integrations are deferred to Tier 3 and evaluated only if customer demand signals that a significant portion of audit clients do not use any compliance automation platform.

**Integration architecture principles (all tiers):**
1. Every integration is an abstraction behind an internal interface — trial balance data enters the same `TrialBalance` data model whether it came from a CSV upload, QBO API, or Codat. Integrations are data sources, not architectural dependencies.
2. OAuth credentials are stored per firm + client, not globally. Each firm's connections are isolated.
3. Evidence from integrations goes through the same AI extraction and review workflow as direct uploads. Integration source is recorded on `EvidenceItem.source_integration` but does not bypass review.
4. Graceful degradation: if Drata's API is down, the auditor falls back to file upload without blocking the engagement.
5. Read-only scopes at launch. No write-back into any client system.

---

## 9. Legal and Data Governance

> Research: [08-legal-data-governance.md](../research/08-legal-data-governance.md)

### Data Processing Agreement Structure

Axiom is a **data processor**; the audit firm is the **data controller**. The firm determines the purposes and means of processing; Axiom processes data only on the firm's documented instructions.

Every customer (US, Canada, EU, UK, AU) must execute a signed DPA before accessing the platform. Clicking "I accept" on the terms of service is not sufficient — a separately executed DPA is required.

**Required DPA contents (GDPR Article 28 and equivalent):**
1. Processing only on documented instructions from the controller
2. Confidentiality obligations on all authorized persons
3. Article 32-compliant technical and organizational measures (AES-256 encryption at rest, TLS 1.3 in transit, access controls, incident response procedures)
4. Sub-processor management: prior written authorization required; equivalent obligations imposed on sub-processors; 30-day advance notice of sub-processor changes via published sub-processor list
5. Data subject rights assistance (access, correction, deletion where not legally blocked)
6. 60-day read-only export window post-termination; deletion of all copies within 30 days of export window close; written deletion certificate provided
7. Audit rights: SOC 2 Type II certification satisfies most audit rights provisions; additional audits on request with reasonable notice
8. Breach notification: Axiom notifies the firm within 24 hours of becoming aware of a personal data breach

**Explicit DPA prohibition:** Axiom does not use customer data (uploaded workpapers, financial data, client lists, evidence) to train AI models without separate controller authorization. All AI model inference runs via AWS Bedrock over a VPC endpoint — data does not leave the AWS network and is not retained by Anthropic. The DPA states this explicitly.

**Sub-processor list:** Published publicly; current sub-processors include AWS (infrastructure, AI model inference via Bedrock, workflow execution via Step Functions, and transactional email via SES — all covered under a single AWS sub-processor entry). 30-day change notification via email to the firm's admin contact.

### GDPR/CCPA Deletion Rights vs. Immutable Archiving — Resolution

The conflict between data subject deletion rights and audit retention obligations is resolved by **GDPR Article 17(3)(b)** and **CCPA California Civil Code §1798.105(d)(8)** — the "legal obligation" exemptions. Retaining an unchanged audit file is required by law (SOX §802, PCAOB AS 1215, AU-C 230). A data subject cannot compel deletion of records that the audit firm is criminally prohibited from deleting.

**Platform implementation of the deletion request workflow:**
1. Acknowledge receipt of the deletion request in writing within 72 hours
2. Invoke the legal retention exemption, citing the specific statute (SOX §802, PCAOB AS 1215 / AU-C 230, GDPR Art. 17(3)(b))
3. Restrict the retained data — no use for any purpose outside the legal retention obligation (no analytics, no AI training)
4. Log the decision in the `AuditLog`: deletion request received, basis for refusal, date of expected deletion
5. Delete on schedule when the mandatory retention window closes (5 or 7 years from report date, per engagement type)

The platform provides templated deletion request response letters that the audit firm can use with data subjects. The response is generated automatically when a deletion request is logged.

Data minimization: the platform only captures personal data strictly necessary for audit documentation, reducing the surface area of deletion conflicts.

### Data Residency Approach

**US at launch (AWS us-east-1):** No federal data residency law mandates US-origin financial data remain in the US. All launch customers are US/Canada firms; AWS us-east-1 is the default and only deployment region.

**Canada:** AWS us-east-1 is acceptable under PIPEDA (no mandatory residency). DPA addresses PIPEDA cross-border transfer obligations. For Quebec-based clients, a Privacy Impact Assessment (PIA) is documented per Quebec Law 25 requirements.

**EU/APAC as Enterprise tier (year 2+):** EU customers (GDPR) require either EU-US Data Privacy Framework self-certification or Standard Contractual Clauses (2021 version, Module 2: Controller to Processor) incorporated into the DPA. Given DPF political risk (Schrems III challenge pending), maintaining SCCs as a parallel mechanism is maintained regardless of DPF status. For EU enterprise customers, offer AWS eu-central-1 (Frankfurt) data residency as a product tier — eliminates the transfer compliance question entirely.

For UK enterprise customers: execute the ICO IDTA or EU SCCs with UK Addendum. AWS eu-west-2 (London) residency option for FRC-regulated engagements. For Australian enterprise customers: APP 8 compliance via DPA contractual undertaking; AWS ap-southeast-2 (Sydney) residency option for APRA-regulated audit clients.

### AI Liability Policy

**Human review is mandatory; AI never auto-finalizes.** This is the single most important platform design decision for managing liability.

Every AI-generated workpaper narrative, finding suggestion, control conclusion, or evidence assessment requires explicit human review and named sign-off before it is treated as finalized in the platform. This creates an auditable paper trail demonstrating that the auditor exercised independent professional judgment, which breaks the causal chain for vendor liability in the event of an audit error.

AI liability protections embedded in the MSA:
- Disclaimer that Axiom does not provide accounting, auditing, tax, legal, or other professional services
- Liability cap equal to 12 months of fees paid (direct damages only)
- Consequential damages waiver (lost profits, business interruption, third-party claims)
- Data security super-cap: 2× annual fees for data breaches
- Gross negligence and fraud carve-out from the limitation of liability

Every AI output is logged with: which model generated it, model version, date, and the inputs that produced it. This is essential for post-incident investigation and demonstrates the AI was used as an assistive tool.

### Insurance Requirements

| Coverage | Minimum at Launch | Target at Series A |
|---|---|---|
| Tech E&O (Professional Liability) | $1M per occurrence / $1M aggregate | $2M per occurrence / $2M aggregate |
| Cyber Liability | $1M per occurrence | $2M–$5M per occurrence |
| General Commercial Liability | $1M/$2M | $1M/$2M |

Secure Tech E&O and Cyber coverage before signing the first enterprise customer contract. Most enterprise MSAs require proof of coverage before execution. Specialist brokers for early-stage SaaS with financial services exposure: Vouch, Corvus (Travelers), Corgi Insurance.

Axiom's standard MSA states coverage levels and commits to maintaining them throughout the contract term.

### Offboarding Obligations

**Data export:** Upon contract termination, the firm receives 60 days of read-only platform access for data export. Export formats: PDF (rendered workpapers), CSV/XLSX (structured data), JSON (engagement metadata and API exports), native format for uploaded documents. The export package includes the complete audit evidence chain: source documents, AI decision records, human review sign-offs, final workpaper content, and version history. A written deletion certificate is provided within 30 days of the export window close.

**Data deletion:** All customer data is deleted from production systems, backups, and sub-processor systems within 30 days of the export window close. Encrypted isolated backups are retained for up to 90 additional days then purged on the next backup rotation.

**Retention exception:** Where a legal obligation (PCAOB AS 1215, AU-C 230, SOX §802) requires the firm to retain records, the deletion obligation runs to the firm — Axiom deletes its copies; the firm is responsible for maintaining compliant archival copies. The DPA states this explicitly: "Upon termination, Processor shall provide Customer with a data export and thereafter delete Customer Data from all Processor systems within 60 days, unless a legal obligation requires Processor to retain such data, in which case Processor shall notify Customer and restrict the data to the minimum necessary for compliance with such obligation."

---

## 10. Security Architecture

### Encryption

**In transit:** TLS 1.2 minimum for all connections (browser → CloudFront → API, API → database, API → S3, API → Bedrock via VPC endpoint). TLS 1.3 is preferred and negotiated when both sides support it. The ALB uses the `ELBSecurityPolicy-TLS13-1-2-2021-06` policy, which supports TLS 1.2 and 1.3 (TLS 1.3-only would break some enterprise clients). No exceptions; HTTP requests are redirected to HTTPS.

**At rest:** AES-256 server-side encryption for all data:
- RDS PostgreSQL: AES-256 via customer-managed KMS key (`axiom-{env}-rds`) with annual automatic rotation
- S3 evidence files: AES-256 SSE-S3 default; HIPAA-flagged evidence uses SSE-KMS with a dedicated customer-managed key (`axiom-{env}-hipaa`) for CloudTrail decrypt auditing and IAM-level key access control (see [`infrastructure-design.md`](./infrastructure-design.md) Section 4 for rationale)
- S3 archive bucket (finalized engagements): SSE-KMS with the HIPAA key
- Backups: encrypted with the same KMS key as the source

### Multi-Tenancy Isolation

PostgreSQL RLS as described in Section 5. `firm_id` indexed on every tenant-scoped table. `SET app.current_firm_id` set at session start, never mutable mid-session. RLS policies enforce the firm boundary at the database layer — defense in depth against application-layer bugs. Engagement-level access control and client-user scoping enforced at the application layer (Go middleware: `WithEngagementAccess`, `WithClientScoping`) as described in Section 7.

Penetration testing on the multi-tenancy isolation (cross-tenant data leakage scenarios) is part of the pre-Series A security review package.

### Audit Log (Immutable, Append-Only)

The `AuditLog` table is PostgreSQL insert-only, with a `RULE` preventing `UPDATE` and `DELETE` statements. All significant application events are written here: engagement status changes, workpaper saves, sign-off events, user access grants, AI decision outcomes, deletion request responses, and system-triggered archival events. Log entries use sequential bigint IDs (not UUIDs) for unambiguous temporal ordering.

The `AuditLog` serves dual purposes: regulatory compliance audit trail (PCAOB inspection readiness) and security audit trail (access anomaly detection, incident investigation).

### S3 Object Lock for Finalized Engagements

When an engagement transitions to `Archived` (after the assembly window closes), all workpaper and evidence files associated with that engagement are copied to a dedicated S3 bucket configured with **Object Lock in COMPLIANCE mode**. In COMPLIANCE mode, no user — including AWS root — can delete or overwrite objects until the retention period (5 or 7 years) elapses. This implements PCAOB AS 1215 and AU-C 230 immutability requirements at the storage layer, not just the application layer.

The retention period is set per object at archive time using the engagement's `retention_deadline` field. When the deadline passes, objects transition to standard S3 lifecycle policies and are deleted automatically, satisfying both the legal retention obligation and the GDPR/CCPA obligation to delete once the retention period expires.

### HIPAA BAA

Axiom executes a HIPAA Business Associate Agreement with AWS. The AWS BAA covers all AWS services used by the platform: RDS, S3, Bedrock (AI model inference), Step Functions (workflow execution), and SES (transactional email). The BAA is in place before any HIPAA engagement data is processed, regardless of customer tier. HIPAA engagement files are stored in S3 buckets with SSE-KMS (not SSE-S3) and with more restrictive IAM policies.

Axiom also makes a HIPAA BAA available to audit firm customers who require one for their own compliance documentation. The standard HIPAA BAA is available to Growth and Scale tier customers on request (self-serve download from the customer portal). Enterprise customers receive dedicated BAA review and negotiation, plus a HIPAA compliance documentation package.

---

## 11. User Journeys

Full journey maps with stage-by-stage detail (user actions, touchpoints, emotional states, competitor context, pain points, and opportunities) are in [`docs/user-journeys/all-journeys.md`](../user-journeys/all-journeys.md). The table below summarizes each journey and the personas, goals, and key system gates involved.

| # | Persona | Goal | Key System Gates / Entities | AI Touchpoints |
|---|---------|------|-----------------------------|----------------|
| 1 | FirmAdmin | Set up firm and launch first engagement | `Firm`, `MethodologyTemplate`, `Engagement` scaffold, 5-step onboarding checklist, 14-day trial clock | — |
| 2 | Staff Auditor | Join platform and reach first task | Magic link auth, role assignment, 5-step guided tour, `EngagementTeamMember` assignment | — |
| 3 | Partner | Create and scope a new engagement | `ClientAcceptance` gate (Planning → Fieldwork blocked until signed), EQR independence validation, framework version lock after Fieldwork | Control mapping (Feature 2, Tier 2), risk category suggestions (Feature 6, Tier 2) |
| 4 | Staff Auditor | Import and analyze a trial balance | `TrialBalanceAccount` import, account mapping confirmation, lead schedule generation, sampling calculator (ISA 530 / AU-C 530) | Account mapping (Feature 3, Tier 2), anomaly detection (Feature 7, Tier 1 / Tier 2 for PCAOB) |
| 5 | Staff Auditor | Test controls and prepare workpapers | `TestProcedure` status progression (NotStarted → InProgress → Complete), AI content tracking gate, sign-off hierarchy (preparer → reviewer → partner) | Evidence link suggestions (Feature 5, Tier 2), workpaper narrative draft (Feature 4, Tier 2) |
| 6 | Manager | Review workpapers and advance the engagement | Review notes (immutable, block advancement until resolved), `ReviewComplete` sign-off, phase transition guards (Fieldwork → Review → Reporting) | — |
| 7 | Staff Auditor | Manage document requests and collect evidence | `DocumentRequest` lifecycle, AI review queue (sorted by confidence), automated reminder state machine, `EvidenceLink` on acceptance | Document completeness review (Feature 1, Tier 2), evidence link suggestion (Feature 5, Tier 2) |
| 8 | Client Contact | Fulfill audit document requests | Tokenized Client Hub link (no login, engagement-scoped, 90-day expiry), `ClientAdmin` delegation (single-request scoped), post-archive read-only access | — |
| 9 | Partner | Generate report, finalize, and archive | Report issuance triggers assembly deadline + retention computation, Finalized state locks all content, S3 Object Lock WORM archival, addendum workflow (AU-C 230 §.16) | Report section draft (Feature 8, Tier 2) |
| 10 | EQR Reviewer | Conduct engagement quality review | Read-only access (not `EngagementTeamMember`), `EngagementQualityReview` sign-off gate (Review → Reporting blocked until signed), immutable EQR record | — |

### Regulatory constraints by journey

| Constraint | Standard | Journeys |
|-----------|----------|----------|
| Client acceptance before fieldwork | SQMS 1 | 3 |
| EQR reviewer independence | SQMS 2 / PCAOB AS 1220 | 3, 10 |
| Framework version locked after fieldwork begins | Section 4 requirement | 3 |
| Sign-off hierarchy enforced at data layer | SQMS 1, AU-C 220 | 5, 6 |
| AI-drafted sections must be substantively edited before sign-off | PCAOB AS 1105 | 5, 9 |
| Review notes cannot be deleted | AU-C 230 | 6 |
| Period coverage check for SOC 2 Type II evidence | AT-C 320 | 7 |
| All AI decisions logged as `AIDecision` records | PCAOB AS 1105 | 3, 4, 5, 7, 9 |
| Anomaly detection flags documented for PCAOB engagements | PCAOB AS 1105 | 4 |
| Sampling documentation requirements | AU-C 530 / PCAOB AS 2315 | 4 |
| Assembly deadline enforcement | AU-C 230, PCAOB AS 1215 | 9 |
| WORM archival | PCAOB AS 1215, SOX §802 | 9 |
| Retention periods per engagement type | AU-C 230, PCAOB, HIPAA | 9 |
| Addenda require documented reason and partner sign-off | AU-C 230 §.16 | 9 |
| Client upload tokens expire and require re-generation | Security policy | 7, 8 |

---

## 12. Flows Without Competitor Equivalent

These flows represent genuine Axiom innovation — no competitor (Fieldguide, CaseWare, DataSnipper, Yak) currently offers them. Full design detail is in the journey maps linked above.

1. **Cross-framework evidence satisfaction display** (Journeys 3, 5) — one piece of evidence simultaneously satisfying SOC 2, ISO 27001, and HIPAA requirements, shown in real-time during testing
2. **AI document completeness review with client-facing gap explanations** (Journeys 7, 8) — AI analyzes uploaded documents against request requirements and auto-drafts specific gap explanations for the client
3. **Section-level AI content tracking and edit gate** (Journeys 5, 6, 9, 10) — AI-drafted content tracked per section with modification ratios; unedited AI sections block sign-off; EQR reviewers see engagement-wide AI edit substantiveness summaries, implementing PCAOB AS 1105 at the data layer
4. **EQR independence enforcement** (Journey 10) — system-level validation that the quality reviewer is not on the engagement team
5. **Post-finalization addendum workflow** (Journey 9) — proper AU-C 230 §.16 implementation with immutable original content and versioned addenda
6. **Full-population analytics as alternative to sampling** (Journey 4) — testing entire transaction datasets rather than statistical samples
7. **Automatic assembly deadline computation and WORM archival** (Journey 9) — computed at report issuance, enforced via S3 Object Lock COMPLIANCE mode

---

## 13. Out of Scope at Launch

The following capabilities are explicitly excluded from the MVP. Each exclusion is a deliberate decision, not an oversight.

| Capability | Rationale for Exclusion |
|---|---|
| **Enterprise-tier-only features** (SAML/enterprise SSO, custom SLA, security review package, dedicated CSM) | OAuth SSO with Microsoft and Google is available at all tiers at launch. SAML-based enterprise SSO (for firms with their own IdP), negotiated SLAs, dedicated CSMs, and security review packages are Enterprise-only. The full enterprise procurement motion is not cost-effective before achieving product-market fit in the mid-market. |
| **Internal audit / SOX compliance workflows** | AuditBoard's established territory. Different workflow (continuous monitoring vs. point-in-time engagement), different buyer (internal audit director vs. external CPA), different regulatory framework (IIA standards, not AICPA). Building for this segment dilutes the ICP focus without improving the core product. |
| **ESG / sustainability reporting** | Fieldguide and Workiva are moving into this space with enterprise-grade content. Not a pain point for the target ICP. Frameworks (GRI, ISSB, SASB) would require significant content investment with no immediate revenue. |
| **White-labeling / reseller channels** | Adds multi-tenancy complexity at the branding layer; requires reseller agreement infrastructure; dilutes the Axiom brand before it is established. |
| **Custom AI model training per firm** | Insufficient training data at launch to produce a meaningfully differentiated firm-specific model. Standard Claude models with firm methodology context (via RAG) achieve the same functional goal at launch scale. Revisit at year 3 when engagement history volume is sufficient. |
| **On-device / local model execution** | Not required for the target ICP's security posture. Adds significant infrastructure complexity. Revisit if enterprise customer demand signals an on-premise deployment requirement. |
| **Multi-agent AI orchestration** | LangGraph, CrewAI, and similar frameworks add debugging overhead without adding value for the eight defined AI features (see [ai-architecture-design.md](ai-architecture-design.md)), all of which are single-step or multi-step-with-human-gates. Add orchestration when genuinely needed by a specific feature, not as a framework choice. |
| **AI auto-finalization of any audit content** | Outside the human-in-the-loop policy. Creates professional liability exposure. The regulatory environment (PCAOB AS 1105) requires human review of AI outputs. |
| **Direct ERP/accounting system API integrations (NetSuite, Sage, Xero)** | CSV/Excel import covers the MVP use case without the engineering and maintenance cost of direct integrations. Direct QBO integration is the only accounting API integration in the first 6-month roadmap. NetSuite and Sage are deferred to 6–18 months. |
| **Direct cloud infrastructure integrations (AWS IAM, Okta, GitHub)** | Axiom is on the auditor side. The Drata/Vanta abstraction covers these at a higher value per integration. Direct infrastructure integrations are Tier 3 and evaluated only on firm-expressed demand. |
| **Practice management billing integration** | Karbon API integration (first 6 months) syncs engagement status only — not billing, invoicing, or time tracking. Billing integration involves financial data, firm-specific billing configurations, and significant scope; deferred to year 2. |
| **Mobile application** | The audit workflow is desktop-intensive: reviewing large documents, managing trial balances, sign-off workflows. A mobile-first experience would require a substantially different UX design investment. Web-responsive design is the priority; native mobile apps are year 2+. |
| **Public REST / GraphQL API** | The internal REST API is not versioned or documented for external consumption. A public API surface (versioning, rate limiting, API key management) is in the 6–18 month roadmap alongside webhooks. At launch, the Karbon integration is built as a first-party integration, not via a public API. |
| **PCAOB firm registration and inspection workflows** | PCAOB inspection preparation (assembling inspection packages, tracking inspection findings) is a distinct workflow from engagement management. The target ICP has few or no PCAOB-registered engagements. Revisit when firm base includes enough PCAOB-registered firms to justify the scope. |
| **Multi-office / global firm management** | A single firm with multiple geographic offices (e.g., a US firm with a Canadian subsidiary) can operate in Axiom as one tenant. Cross-office engagement assignment and jurisdiction-specific methodology management are deferred — the ICP is primarily single-office or single-country. |
| **EU/APAC data residency** | US and Canada at launch (AWS us-east-1). EU data residency (Frankfurt) and APAC (Sydney) offered as Enterprise tier options in year 2, after the product is validated in the primary market. |

---

*End of Axiom Product Specification v2*
