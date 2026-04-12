# Axiom Product Specification v2
**Date:** April 12, 2026
**Status:** Implementation-Ready

---

**On the name:** "Axiom" represents a self-evident truth or starting point for reasoning — fitting for a platform anchored in the trial balance as the financial source of truth. It is a working name and can be changed before launch.

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
11. [Onboarding Flows](#11-onboarding-flows)
12. [Engagement Module Specifications](#12-engagement-module-specifications)
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

### Entity Inventory

#### Tenant and Identity Layer

**Firm** — The root tenant entity. All firm-owned data carries `firm_id`. Key fields: `id`, `name`, `slug`, `subscription_tier`, `settings (jsonb)`. Every tenant-scoped table has an indexed `firm_id` column enforced by PostgreSQL RLS.

**User** — Staff at the audit firm or client-side users invited to specific engagements. Key fields: `id`, `firm_id` (nullable for client users), `email`, `name`, `role`. Roles: `FirmAdmin | Partner | Manager | Staff | EQReviewer | ClientAdmin | ClientUser | ViewOnly`. Client users belong to a `Client` record, not a `Firm`; they access only engagements they are explicitly invited to.

**Client** — The entity being audited ("service organization" in SOC 2 context). Key fields: `id`, `firm_id`, `name`, `industry`, `primary_contact_email`. One firm has many clients.

#### Framework Reference Layer (Shared Across All Tenants, Not Tenant-Scoped)

**Framework** — A specific version of a standards framework. Key fields: `id`, `name`, `version`, `effective_date`, `deprecated_at`, `governing_body`. ISO 27001:2013 and ISO 27001:2022 are separate rows.

**FrameworkRequirement** — A single criterion or control within a framework (e.g., SOC 2 CC6.1, ISO 27001 A.8.3, HIPAA §164.312(a)(1)). Key fields: `id`, `framework_id`, `identifier`, `title`, `description`, `category`, `sort_order`.

#### Control Objective Layer (Framework-Agnostic Bridge)

**ControlObjectiveLibrary** — System-maintained semantic control objectives, independent of any framework. Example: "Access to production systems is restricted to authorized personnel." Key fields: `id`, `name`, `description`, `tags (jsonb)`.

**ControlObjectiveLibraryMapping** — Maps a library objective to one or more `FrameworkRequirement` rows across all frameworks. This table encodes the ~80% overlap between SOC 2 and ISO 27001, enabling a single evidence item to satisfy all mapped requirements simultaneously.

**FirmControlObjective** — A firm's customized version of a library objective, or a net-new objective. Key fields: `id`, `firm_id`, `source_library_id` (nullable), `name`, `description`, `custom_test_guidance (jsonb)`. Firms can override library objectives or define their own.

**FirmControlObjectiveMapping** — Maps a `FirmControlObjective` to specific `FrameworkRequirement` rows. This is the firm-specific layer of the cross-framework evidence architecture.

#### Engagement Layer

**Engagement** — One engagement = one client, one primary framework, one audit period. Key fields: `id`, `firm_id`, `client_id`, `name`, `engagement_type` (FinancialAudit_Private | FinancialAudit_Public | SOC1 | SOC2 | ISO27001 | HIPAA | AgreedUponProcedures | Advisory), `primary_framework_id`, `framework_version_id`, `period_start`, `period_end`, `status`, `prior_engagement_id` (FK for rollforward), `methodology_template_id`, `assembly_deadline` (computed: report date + 60 or 45 days), `retention_deadline` (computed: report date + 5 or 7 years), `report_issued_at`, `finalized_at`, `archived_at`.

**EngagementTeamMember** — Associates users to an engagement with their role on that specific engagement. Used for engagement-level access control checks in the application layer. Key fields: `id`, `engagement_id`, `user_id`, `role`, `assigned_at`, `removed_at`.

**EngagementFramework** — Supports multi-framework engagements (e.g., integrated SOC 2 + ISO 27001). One framework is primary; others are secondary. Key fields: `id`, `engagement_id`, `framework_id`, `framework_version_id`, `is_primary`.

#### Control Layer

**Control** — An instantiation of a `FirmControlObjective` within a specific engagement. Key fields: `id`, `engagement_id`, `firm_control_objective_id`, `description`, `control_owner_id`, `auditor_assigned_to_id`, `status` (NotStarted | InProgress | Complete | Exception | NotApplicable), `is_key_control`, `prior_control_id`.

**TestProcedure** — A specific test step within a control. Key fields: `id`, `control_id`, `procedure_type` (Inquiry | Observation | InspectionOfDocument | Reperformance | Analytics), `description`, `expected_result`, `population_size`, `sample_size`, `sampling_method`, `result`, `exceptions_noted`, `conclusion`, `performed_by_id`, `performed_at`, `reviewed_by_id`, `reviewed_at`, `status`.

#### Evidence Layer

Evidence is stored at the **firm + client** level, not at the engagement level. This is the architectural decision that enables year-over-year reuse and cross-framework reuse without re-uploading.

**EvidenceItem** — A single uploaded document or artifact. Key fields: `id`, `firm_id`, `client_id`, `filename`, `storage_path` (S3 key), `content_type`, `file_size_bytes`, `uploaded_by_id`, `uploaded_at`, `source_type` (ClientUpload | CloudIntegration | APIImport | AuditorGenerated), `source_integration` (Dropbox | Box | GoogleDrive | O365 | DratAuditHub | Vanta | null), `extracted_text`, `extraction_status`, `is_sensitive` (PII/PHI flag).

**EvidenceLink** — Connects an `EvidenceItem` to a specific `TestProcedure`. Key fields: `id`, `evidence_item_id`, `test_procedure_id`, `linked_by_id`, `linked_at`, `notes`, `ai_suggested` (bool), `ai_decision_id` (FK to AIDecision).

#### Framework-Agnostic Evidence Chain

The full relationship that makes cross-framework evidence reuse work:

```
EvidenceItem
  → (EvidenceLink)
  → TestProcedure
  → Control
  → FirmControlObjective
  → (FirmControlObjectiveMapping)
  → FrameworkRequirement (SOC 2 CC6.1)
  → FrameworkRequirement (ISO 27001 A.8.3)
  → FrameworkRequirement (HIPAA §164.312(a)(1))
```

One evidence upload, one link action, and the evidence simultaneously satisfies all framework requirements mapped to that control objective. This is the architectural realization of the Hyperproof/Drata cross-framework pattern on the auditor side — no competitor currently does this.

#### Document Request Layer

**DocumentRequest** — A PBC (Provided By Client) request sent to the client. Key fields: `id`, `engagement_id`, `control_id` (nullable), `assigned_to_id`, `title`, `description`, `instructions`, `due_date`, `status` (Pending | Submitted | InReview | Accepted | Rejected | Overdue), `reminder_count`, `last_reminder_sent_at`, `fulfilled_by_evidence_item_id`.

#### Trial Balance Layer (Financial Audit Only)

**TrialBalance** — Container for an imported trial balance. Key fields: `id`, `engagement_id`, `period_date`, `import_source`, `imported_at`, `imported_by_id`.

**TrialBalanceAccount** — Individual account row. Key fields: `id`, `trial_balance_id`, `account_number`, `account_name`, `account_type`, `balance_debit`, `balance_credit`, `net_balance`, `mapped_fs_line_item`, `mapping_status` (Unmapped | AISuggested | Confirmed | Overridden), `ai_decision_id`, `confirmed_by_id`.

**TrialBalanceAdjustment** — Proposed, passed, or waived adjustments. Key fields: `id`, `trial_balance_id`, `account_id`, `amount`, `description`, `adjustment_type` (Proposed | Passed | Waived), `proposed_by_id`, `approved_by_id`.

#### Workpaper Layer

**Workpaper** — A document in the engagement file. Key fields: `id`, `engagement_id`, `workpaper_type` (LeadSchedule | TestPaper | Memo | ConfirmationLetter | SamplingWorksheet | ManagementLetter | Other), `title`, `content (jsonb)`, `status` (Draft | PreparedPendingReview | InReview | ReviewNotesOpen | ReviewComplete | SignedOff), `prepared_by_id`, `reviewed_by_id`, `signed_off_by_id`, `is_locked` (true after assembly deadline), `prior_workpaper_id`.

**WorkpaperVersion** — Immutable version history. Every save creates a new row. Key fields: `id`, `workpaper_id`, `version_number`, `content (jsonb)`, `saved_by_id`, `saved_at`, `is_ai_draft` (bool — true until a human edits any content, per PCAOB AS 1105 requirement to distinguish AI-generated from auditor-authored content).

#### Quality Management Layer (SQMS 1/2)

**ClientAcceptance** — Per-engagement quality risk documentation required by SQMS 1. Key fields: `id`, `engagement_id`, `quality_risks_identified (jsonb)`, `firm_responses (jsonb)`, `independence_confirmed`, `independence_confirmed_by_id`, `accepted_by_id`, `accepted_at`. Required before Planning → Fieldwork transition.

**EngagementQualityReview** — Formal EQR record (SQMS 2 / PCAOB AS 1220). Key fields: `id`, `engagement_id`, `reviewer_id` (must have EQReviewer role; must not be on engagement team — system-enforced), `independence_documented_at`, `status` (Assigned | InProgress | Complete), `scope_notes`, `conclusion`, `signed_off_at`. Required before Review → Reporting transition where EQR is applicable.

#### AI Decision Layer

**AIDecision** — Every AI output that could affect audit content is recorded here. Required for PCAOB engagements; best practice for all others. Key fields: `id`, `firm_id`, `engagement_id` (nullable), `context_type` (EvidenceReview | ControlMapping | RiskAssessment | TrialBalanceMapping | DocumentCompleteness | WorkpaperDraft | SamplingRecommendation | AnomalyDetection), `context_id` (uuid), `context_table` (which table), `model_id` (e.g., "claude-sonnet-4-6"), `input_token_count`, `output_token_count`, `raw_output (jsonb)`, `suggested_value`, `confidence` (float 0–1), `review_action` (Pending | Accepted | Modified | Rejected), `accepted_value`, `reviewed_by_id`, `reviewed_at`.

The `context_type + context_id + context_table` triple provides a queryable link to exactly what was being analyzed without a polymorphic foreign key that would complicate migrations.

#### Audit Log (Immutable)

**AuditLog** — Append-only. No updates, no deletes. Key fields: `id` (sequential bigint — not UUID, for ordering guarantees), `firm_id`, `actor_id`, `actor_type` (User | System | AIAgent), `action` (e.g., "engagement.status.changed", "workpaper.signed_off", "evidence.linked"), `resource_type`, `resource_id`, `old_value (jsonb)`, `new_value (jsonb)`, `ip_address`, `user_agent`, `occurred_at` (timestamptz).

Implemented as a PostgreSQL insert-only table with a `RULE` preventing `UPDATE` and `DELETE`. This satisfies regulatory immutability requirements and the GDPR audit trail obligation simultaneously.

#### Reporting Layer

**Report** — Key fields: `id`, `engagement_id`, `report_type` (SOC2Type1 | SOC2Type2 | SOC1Type1 | SOC1Type2 | FinancialAuditOpinion | AgreedUponProcedures | ManagementLetter), `status` (Draft | ClientReview | FirmReview | Issued | Archived), `content (jsonb)`, `generated_at`, `issued_at`, `issued_by_id`.

**ReportVersion** — Immutable version history per report. Same pattern as WorkpaperVersion.

#### Methodology Template Layer

**MethodologyTemplate** — Firm-level reusable templates that pre-populate a new engagement with controls, test procedures, and workpaper shells. Key fields: `id`, `firm_id`, `name`, `applicable_engagement_type`, `applicable_framework_id`, `version`, `is_active`.

**TemplateControl** and **TemplateTestProcedure** — Control and test procedure definitions within a template.

### Engagement Lifecycle State Machine

```
Planning ──[Partner: acceptance complete]──► Fieldwork
Fieldwork ──[Manager/Partner: all controls have results]──► Review
Review ──[Partner: notes resolved + EQR signed off]──► Reporting
Reporting ──[Partner: report issued]──► Finalized
Finalized ──[System: assembly window elapsed]──► Archived (IMMUTABLE)

Reverse paths (exceptional):
Fieldwork ──[Partner: scope change]──► Planning
Review ──[Manager/Partner: additional procedures needed]──► Fieldwork
Reporting ──[Partner: significant issue found]──► Review
Any state ──[FirmAdmin: abandoned engagement]──► Archived
```

**Valid transitions and guards:**

| From | To | Who Triggers | Guard Condition |
|---|---|---|---|
| Planning | Fieldwork | Partner | `ClientAcceptance.accepted_at` is populated |
| Fieldwork | Review | Manager or Partner | All `Control` records have status `Complete` or `Exception` |
| Review | Reporting | Partner | All review notes resolved; `EngagementQualityReview.status = Complete` where applicable |
| Reporting | Finalized | Partner | `Report.status = Issued` |
| Finalized | Archived | System (Step Functions `Wait` state) | `report_issued_at + assembly_window` has elapsed |
| Fieldwork | Planning | Partner | (scope change requiring re-acceptance) |
| Review | Fieldwork | Manager or Partner | (additional procedures required) |
| Reporting | Review | Partner | (significant issue identified post-reporting) |
| Any | Archived | FirmAdmin | (abandoned engagement) |

**Finalized is a hard gate.** Once an engagement reaches Finalized, no workpaper content can be modified. Addenda are created as new `WorkpaperVersion` records with `is_addendum = true` and require re-sign-off. This implements AU-C 230 and PCAOB AS 1215 assembly deadline compliance at the data layer, not just in the UI.

### Year-Over-Year Rollforward Behavior

When a new engagement is created with `prior_engagement_id` set:

| Entity | Behavior |
|---|---|
| Engagement | New record; `prior_engagement_id` set; all status fields reset |
| Controls | Cloned from prior engagement; `prior_control_id` set on each; status reset to NotStarted |
| TestProcedures | Cloned per control; status reset; available as editable starting point |
| DocumentRequests | Not auto-cloned; AI suggests new requests based on prior engagement controls |
| EvidenceItems | Not touched; they exist at the firm+client level and are surfaced with "used in prior year" flag |
| TrialBalance | New import required; prior year TB accessible as read-only reference for comparatives |
| Workpapers | New drafts created; `prior_workpaper_id` set; prior year workpapers visible as read-only sidebar reference |
| Report | New document; prior year report accessible for reference only |
| AIDecisions | Not carried forward; AI re-analyzes fresh evidence for the new period |
| ClientAcceptance | New record required; quality risk acceptance must be refreshed annually |
| EngagementQualityReview | New record required if applicable |

### Multi-Tenancy Isolation

**Model: Shared PostgreSQL database with row-level security (RLS).**

All tenant-scoped tables carry `firm_id` with an index. The application sets `SET app.current_firm_id = '<uuid>'` at the start of every request. RLS policies enforce the firm isolation boundary:

```sql
CREATE POLICY firm_isolation ON engagements
  USING (firm_id = current_setting('app.current_firm_id')::uuid);
```

RLS is the safety net — defense in depth. Application-layer authorization (REST middleware) is the primary mechanism. The three authorization dimensions are:

1. **Firm isolation** — RLS + `withFirmIsolation` middleware. Every query scoped to `current_firm_id`.
2. **Engagement team membership** — Application layer `withEngagementAccess` middleware. An indexed point lookup on `EngagementTeamMember (engagement_id, user_id)`.
3. **Client user scoping** — Application layer `withClientScoping` middleware. Client users can only see `DocumentRequest` and `EvidenceItem` records linked to engagements they were invited to.

System-wide reference tables (`Framework`, `FrameworkRequirement`, `ControlObjectiveLibrary`, `ControlObjectiveLibraryMapping`) have no `firm_id` and no RLS — they are read-only reference data shared across all tenants.

---

## 6. AI Architecture

> Research: [05-ai-architecture.md](../research/05-ai-architecture.md)

### LLM Provider Decision

**Primary provider: AWS Bedrock with Claude models (Haiku and Sonnet), via the Anthropic/AWS partnership.**

Rationale:
- **PrivateLink** — model API calls are made over a VPC endpoint and never leave the AWS network. No outbound HTTPS to an external provider; stronger posture than ZDR mode at the API level.
- **IAM-based auth** — no separate vendor API keys to store in Secrets Manager, rotate, or audit. Access control is standard IAM policy.
- **Single AWS BAA** — Bedrock is covered under the existing AWS HIPAA BAA. Anthropic is not an additional sub-processor in the DPA or the SOC 2 audit scope.
- **Same models, same capability** — Bedrock provides Claude Haiku and Sonnet with the same 200K+ context window, prompt caching (0.1× the base input rate for cached reads), and batch inference (50% discount) as Anthropic direct. The cost model is effectively identical at this scale.
- **CloudWatch metrics natively** — token counts, latency, and errors report to CloudWatch without additional instrumentation. No separate dashboard or custom metrics pipeline needed.
- Single model family eliminates the prompt engineering and behavior consistency overhead of a multi-provider setup.

Model availability on Bedrock typically lags Anthropic's direct API by a few weeks on new releases. For a compliance platform where model stability is preferred over early access to new capabilities, this is a feature rather than a drawback — a validation period before switching models would be standard practice regardless.

### Vector Database Decision

**pgvector (PostgreSQL extension) at launch, with migration path to Qdrant.**

pgvector adds zero operational overhead (uses existing PostgreSQL instance), keeps vector data in the same security boundary as the rest of the application, and performs adequately (sub-20ms at 1M vectors with HNSW index) for the target scale. At maturity (50 firms × 100 engagements × 200 evidence items), the total vector count is approximately 1M — well within pgvector's comfortable range. Migrate to self-hosted Qdrant when approaching 5–10M vectors.

Embeddings are stored for: firm methodology documents, framework requirement descriptions, control objective library entries, prior engagement evidence summaries, and workpaper templates.

### Four AI Features

#### Feature 1: Document Completeness Review

**Purpose:** Eliminate the primary engagement bottleneck — the back-and-forth cycle when clients submit incomplete or wrong documents.

**Trigger:** Client uploads a document in response to a `DocumentRequest`.

**Inputs:**
- `DocumentRequest` (title, description, instructions, linked `Control`, linked `TestProcedure`)
- Uploaded document text (OCR pipeline output)
- `ControlObjective` description and `FrameworkRequirement` text
- RAG context: embeddings of prior accepted evidence for this request type (from same firm, anonymized)

**Model:** Claude Sonnet (reasoning quality is critical; errors propagate into the audit file).

**Process:** Extract document attributes (date range, parties, system names, amounts, format) → compare against request criteria → check period coverage for Type II engagements → score completeness against each test procedure → identify specific gaps → generate client-facing plain-language explanation → generate auditor-facing summary → produce recommendation: Accept | Request Clarification | Reject.

**Output:** `AIDecision` record with `context_type = DocumentCompleteness`, `suggested_value`, `confidence` score, and full `raw_output`. Status: `Pending`.

**Human review gate:** Required before any Accept action. The auditor sees the AI assessment and one-click actions (Accept / Modify / Reject). The `AIDecision` record captures which action was taken, by whom, and when. No document is marked as fulfilling a request without a named human decision.

**Failure modes:**

| Failure | Handling |
|---|---|
| Encrypted / password-protected document | Extraction fails → document flagged for manual review; client notified |
| Non-English document | Detected language flagged → manual review queue |
| AI confidence below 0.6 | Surfaced as "Low confidence — manual review recommended" without a definitive recommendation |
| Image-only PDF | OCR attempted; if confidence below threshold, flagged for manual |
| API error / timeout | Retry up to 3 times with exponential backoff; auditor notified if all retries fail |

---

#### Feature 2: Control Mapping (Framework-Agnostic Evidence Linkage)

**Purpose:** Propose `FirmControlObjectiveMapping` records across all frameworks in the engagement simultaneously — the architectural realization of the cross-framework differentiator.

**Trigger:** New engagement created from a methodology template, or new `FirmControlObjective` added to an engagement.

**Inputs:**
- `FirmControlObjective` name + description
- All `FrameworkRequirement` records for frameworks in the engagement
- RAG context: `ControlObjectiveLibrary` entries with existing mappings as few-shot examples

**Model:** Claude Haiku (structured classification task; high volume; speed matters at template instantiation).

**Process:** Embed the `FirmControlObjective` description → retrieve top-k similar library entries → score semantic similarity against each applicable `FrameworkRequirement` → apply 0.75 similarity threshold → generate explanation per proposed mapping → return proposed mappings with confidence scores.

**Output:** Mapping table displayed to auditor showing each proposed `FrameworkRequirement` link with confidence and explanation text. All proposed mappings remain in `Pending` state until confirmed.

**Human review gate:** Auditor reviews proposed mappings in bulk. All confirmed by default; any can be rejected. No `FirmControlObjectiveMapping` record is created until confirmed. `AIDecision` records created per mapping.

---

#### Feature 3: Trial Balance Account Mapping

**Purpose:** Classify each trial balance account into a standard financial statement line item, eliminating manual mapping work at engagement start.

**Trigger:** Trial balance imported (CSV/Excel upload or API).

**Inputs:** Account number, account name, balance (debit/credit), prior year mapping if rollforward.

**Model:** Claude Haiku (simple classification; high volume; prior year context improves accuracy significantly on rollforward engagements).

**Process:** Few-shot classification against standard FS line items (Cash, Accounts Receivable, Fixed Assets, etc.) with anomaly flagging for accounts where classification confidence is low or account name doesn't match expected category.

**Output:** `TrialBalanceAccount.mapping_status = AISuggested` with `ai_decision_id` populated for each mapped account. Low-confidence accounts highlighted in the Sheets UI for priority auditor review.

**Human review gate:** Auditor reviews and confirms mappings in the Sheets interface. Unmapped and low-confidence accounts surfaced prominently. `mapping_status` changes to `Confirmed` or `Overridden` on auditor action.

---

#### Feature 4: Workpaper Narrative Draft

**Purpose:** Generate a first-draft workpaper narrative that a staff auditor would write, matching firm style — reducing blank-page time and ensuring consistent structure across the team.

**Trigger:** Auditor marks a `TestProcedure` as Complete and explicitly requests a draft (not automatic).

**Inputs:**
- Control description and `TestProcedure` (type, description, expected result)
- Linked evidence items (extracted text and filenames)
- Exceptions noted (if any)
- Prior year workpaper narrative if rollforward (style reference)
- Firm workpaper template for this procedure type

**Model:** Claude Sonnet (quality matters; workpaper narratives become part of the immutable audit file).

**Output:** Draft text inserted into the workpaper editor. `WorkpaperVersion` record created with `is_ai_draft = true`. The draft is labeled "AI Draft — requires review" until any human edit changes `is_ai_draft` to false.

**Human review gate:** Mandatory. The AI draft is always editable. The workpaper's `status` cannot advance to `PreparedPendingReview` or beyond if `is_ai_draft = true` and no human edits have been made. The auditor must actively modify and sign off. This satisfies the PCAOB requirement to distinguish AI-generated from auditor-authored content.

---

### Human-in-the-Loop Policy

**Tier 1 — Fully Automated (no human approval required):**
AI may act without user interaction for: text extraction from uploaded documents, embedding generation, flagging documents as potentially incomplete (notification only, not decision), generating overdue reminder emails to clients, surfacing "this evidence was used in a prior engagement" suggestions, running anomaly detection on trial balance data (flagging in UI only), and computing trial balance analytics (variance, ratios).

**Tier 2 — Human Approval Required Before Taking Effect:**
AI suggests; human confirms. No audit file content changes without explicit auditor action: evidence links, control → framework requirement mappings, trial balance account category assignments, document request fulfillment decisions (accept/reject), risk assessment pre-population, workpaper narrative drafts.

**Tier 3 — Human-Initiated Only (AI may assist but never initiates):**
These actions can only be triggered by a named human: engagement status transitions, control conclusions (pass/fail), exception documentation, workpaper review and sign-off, EQR sign-off, report issuance, client acceptance documentation, and any action constituting professional judgment under AU-C or PCAOB standards.

This three-tier policy maps directly to PCAOB AS 1105: AI outputs inform professional judgments; they do not constitute them.

### Per-Engagement AI Cost Estimate

| AI Operation | Volume (per SOC 2 engagement) | Model | Approx. Cost |
|---|---|---|---|
| Document completeness review | 200 docs | Haiku (batch) | ~$0.30 |
| Control mapping | 50 controls | Haiku (batch) | ~$0.11 |
| Trial balance mapping | N/A for SOC 2 | — | — |
| Workpaper narrative drafts | 20 drafts | Sonnet | ~$1.20 |
| Evidence link suggestions | 200 links | Haiku | ~$0.15 |
| RAG retrieval overhead | ~50K tokens | Haiku | ~$0.04 |
| **Total per SOC 2 engagement** | | | **~$1.90** |
| **Total per financial audit** | | | **~$3–6** (larger trial balance, more workpapers) |

With prompt caching (system prompts reused within an engagement): effective cost **$0.75–$2.50 per engagement**. At 100 engagements/year for a mid-market firm, platform AI costs are $75–$250/year. This is absorbed into the subscription — AI consumption is not a separate line item on customer invoices at launch.

### What AI Does NOT Do at Launch

- Fine-tuned models per firm (insufficient training data volume at this stage)
- On-device / local model execution
- Multi-agent orchestration (LangGraph, CrewAI) — single-step AI calls are sufficient; frameworks add debugging overhead with no benefit for the defined use cases
- AI-generated audit opinions, professional conclusions, or materiality determinations
- Autonomous multi-step actions without human review gates at each step
- Auto-finalization of any audit content without a named human reviewer

---

## 7. Technology Stack

### Frontend

**TypeScript / React SPA.** Component library: Shadcn/ui (accessible, customizable, no licensing overhead). State management: TanStack Query for server state; Zustand for local UI state. API types are generated from the OpenAPI specs in `packages/openapi/` via `openapi-typescript` — a spec change automatically regenerates the client on the next build.

### Backend: Go Microservices + Python PDF Service

**Decision: Go as the primary backend language, decomposed into bounded-context microservices. Python retained for PDF extraction.**

Go was chosen for its compile-time type safety, lean container images (30–50MB per Fargate task), and strong fit for compliance SaaS — Workiva, an established player in this space, uses Go for their REST services. The backend is split into services along genuine bounded contexts in the data model, not along the five product modules. The core engagement/control/evidence cluster is too tightly coupled to split without distributed transactions and stays in one service with a shared database. Independent domains get their own services and databases.

Python is retained as a single stateless service for PDF extraction. `pdfplumber` handles complex, multi-column, scanned audit documents better than any Go library. The polyglot cost is contained — one endpoint, one job, no shared state.

**Full service decomposition, database topology, Go and Python tech stack choices, and inter-service communication patterns are specified in [`backend-architecture-design.md`](./backend-architecture-design.md).**

**Monorepo structure:**
```
apps/
  gateway/          — Go: API Gateway (JWT verification, routing)
  identity/         — Go: Identity Service (Firm, User, Client, templates)
  audit-core/       — Go: Audit Core (Engagement, Controls, Evidence, AI)
  trial-balance/    — Go: Trial Balance Service
  workpaper/        — Go: Workpaper Service (Yjs collaboration)
  reporting/        — Go: Reporting Service
  doc-processing/   — Python: PDF extraction only
packages/
  go-shared/        — Shared Go: JWT middleware, SQS wrappers, OTel setup
  openapi/          — OpenAPI 3.1 specs for all services (source of truth)
  ai/               — Go: Claude API client wrappers, AIDecision recording
```

Turborepo manages the monorepo with per-service build caching.

### API Layer: REST + OpenAPI

**Decision: REST with OpenAPI 3.1 for all services. Hasura is rejected.**

**Why REST with OpenAPI:** REST decouples the frontend from the backend language, supporting the Go services alongside the Node.js services without sharing a runtime boundary. OpenAPI enables future public API exposure (webhooks, partner integrations) without adding a separate layer.

Each service defines its API contract as an OpenAPI 3.1 spec in `packages/openapi/`. `oapi-codegen` generates typed Go server interfaces from the spec. `openapi-typescript` generates typed fetch clients for the React frontend. Authorization is enforced as composable Go middleware in each service:

- `WithFirmIsolation` — reads `firm_id` from gateway-injected headers, sets Postgres session variable for RLS
- `WithEngagementAccess` — verifies `EngagementTeamMember` record exists for the requested engagement
- `WithClientScoping` — for `ClientUser` roles, filters to invited engagements only

**Why Hasura is rejected:** The evidence authorization chain requires a five-level relationship traversal (`EvidenceItem → EvidenceLink → TestProcedure → Control → Engagement → EngagementTeamMember`). Hasura v2's `_exists` permission predicates materialize this as an `IN (...)` query that degrades as `EngagementTeamMember` grows. Hasura permissions are YAML metadata — not code, not testable with standard unit tests, difficult to reason about as authorization rules evolve. The auto-generated CRUD value Hasura provides is not justified for a platform where most queries are purpose-built for specific UI requirements.

**Public API (year 2+):** The OpenAPI specs already define the contract. Exposing selected endpoints publicly requires adding authentication scopes and rate limiting at the API Gateway — no rewrite of business logic.

### Database: PostgreSQL with RLS

PostgreSQL is the sole persistent data store. One RDS instance hosts five logical databases — one per service — with no cross-database foreign keys. The Audit Core database (`core_db`) uses row-level security (RLS) for multi-tenancy as described in Section 5; other service databases enforce tenant isolation at the application layer. Database access uses `sqlc` + `pgx/v5` (type-safe SQL generation from plain SQL query files) with `golang-migrate` for schema migrations. PgBouncer for connection pooling in transaction mode.

**pgvector** extension enabled on `core_db` for embedding storage (Section 6).

**Workpaper content** stored as typed jsonb in `Workpaper.content`. The jsonb structure supports rich text (ProseMirror document format), embedded tables, formula references, and metadata. A dedicated document store is not needed at launch scale.

### Workflow Engine: Two-Tier (River + Step Functions)

**Tier 1 — River (PostgreSQL-based job queue) for background jobs:**

All fire-and-forget background work uses River, a Go-native Postgres-backed job queue. Zero additional infrastructure — it uses the existing service database. Jobs are durable (WAL-backed), support retry with exponential backoff, and have dead-letter queues.

River jobs (running within Audit Core against `core_db`):
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

**Workpaper rich text (TipTap + Yjs):** TipTap (ProseMirror-based editor) with Yjs CRDT for real-time co-editing of workpaper narratives. Yjs handles keystroke-level text concurrency naturally. Each save (triggered by idle timeout or explicit save) creates a `WorkpaperVersion` record. The `is_ai_draft` flag on `WorkpaperVersion` is cleared when any human edits content, satisfying PCAOB's requirement to distinguish AI-generated from auditor-authored content.

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
- **ECS Fargate** — Container orchestration for all seven services (gateway, identity, audit-core, trial-balance, workpaper, reporting, doc-processing). ECS Service Connect provides internal DNS-based service discovery. Independent per-service scaling; Workpaper service scales on active WebSocket connection count via a custom CloudWatch metric.
- **ALB** — Application Load Balancer for TLS termination in front of the API Gateway
- **RDS PostgreSQL** — Single Multi-AZ instance hosting five logical databases (one per service); pgvector extension enabled on `core_db`
- **SQS** — Async cross-service event delivery (e.g., document uploaded → extraction triggered)
- **S3** — Evidence file storage; Object Lock enabled for finalized engagements (see Section 10)
- **CloudFront** — CDN for the React SPA
- **SES** — Transactional email (document request notifications, client invitations, review alerts)
- **Secrets Manager** — API keys, database credentials, OAuth tokens
- **CloudWatch + X-Ray** — Logging, monitoring, distributed tracing (via OpenTelemetry Go SDK)

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

**In transit:** TLS 1.3 minimum for all connections (browser → CloudFront → API, API → database, API → S3, API → Bedrock via VPC endpoint). No exceptions; HTTP upgrade headers redirect all HTTP traffic.

**At rest:** AES-256 server-side encryption for all data:
- RDS PostgreSQL: AES-256 via AWS-managed keys (KMS)
- S3 evidence files: AES-256 SSE-S3 (upgraded to SSE-KMS for HIPAA engagements)
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

## 11. Onboarding Flows

This section defines the step-by-step journeys for the three distinct user populations that enter the platform: audit firm administrators (self-serve firm setup), individual firm staff members (invitation-based access), and client-side contacts (per-engagement invitation). Because Axiom's self-serve model is a structural competitive advantage over Fieldguide, every onboarding flow is designed for zero-consultant completion.

---

### A. New Firm Onboarding (Self-Serve)

**Goal:** Get a new firm from trial signup to a first active engagement in under one week, without human intervention on Axiom's side.

**Steps:**

1. **Trial signup** — The prospective firm admin visits the Axiom marketing site and enters a business email address. A verification email is sent. On verification, the admin completes a brief intake form:
   - Firm name
   - Staff count (dropdown: 1–10, 11–20, 21–40, 41–60, 60+)
   - Primary audit types (multi-select: Financial Audit, SOC 2, ISO 27001, HIPAA)
   - Country (US / Canada)
   The system creates a `Firm` record with `subscription_tier = Trial` and `onboarding_status = ProfileSetup`. The 14-day trial clock starts.

2. **Firm profile setup** — The admin is taken directly into the app. A required setup step collects: firm display name, logo (optional), timezone, and billing contact email. `onboarding_status` advances to `MethodologySetup`.

3. **Methodology selection** — The admin activates one or more pre-built methodology templates:
   - AICPA/GAAS Financial Audit (always available)
   - SOC 2 Type I/II (TSC 2017)
   - ISO 27001:2022
   - HIPAA
   Growth tier: templates are read-only pre-built. Scale tier: custom template editor is also unlocked (configurable after initial activation). The admin does not need to configure anything to activate — selecting a template makes it immediately available for engagement creation. `onboarding_status` advances to `FirstEngagement`.

4. **First engagement wizard** — The system presents a guided engagement creation wizard. The admin (or any Partner/Manager role invited in the next step) selects:
   - Engagement type and framework
   - Client name (can use "Demo Client" for a test run)
   - Audit period
   - Methodology template (pre-selected based on signup intake)
   The system creates the full engagement scaffold: `Engagement` in `Planning` state, all template `Control` and `TestProcedure` records, draft `Workpaper` shells, and an empty `ClientAcceptance` record. This is the completion event for the "time-to-first-engagement" product metric. `onboarding_status` advances to `StaffInvite`.

5. **Staff invitation** — The admin invites firm staff by email. Bulk invite (paste multiple emails) is supported. For each invitee the admin assigns a role: `Partner | Manager | Staff | EQReviewer`. Invitees receive a magic link email valid for 7 days. At least one staff invite is encouraged (in-app prompt) but not required to proceed. `onboarding_status` advances to `Complete` once at least one invite is sent.

6. **Optional SSO setup** — Growth and Scale tiers: OAuth with Microsoft or Google (configured in Firm Settings → Authentication at any time). Enterprise tier: SAML with the firm's own IdP. SSO setup is not a prerequisite for using the platform and is not part of the onboarding checklist.

**In-app onboarding checklist:** A collapsible progress panel persists in the FirmAdmin sidebar until all five required steps are complete: profile setup, methodology activated, first engagement created, first staff invite sent. Each completed step is checked off. The checklist disappears permanently once all five are complete.

**Scale-tier onboarding call:** Scale customers are offered (not required) a 30-minute onboarding call with an Axiom support specialist. The offer appears as an in-app prompt on Day 3 of the trial and again on Day 10 if the trial has not yet converted. This is not a multi-week implementation program — it is a single orientation call. The scheduling link is in-app (Calendly or equivalent).

**Data entities touched:** `Firm`, `User` (FirmAdmin), `MethodologyTemplate`, `Engagement`, `EngagementTeamMember`, `Control`, `TestProcedure`, `Workpaper`, `ClientAcceptance`.

---

### B. Individual Firm Staff User Onboarding

**Goal:** Each staff member invited by the FirmAdmin reaches their first assigned task in a single session, with no friction from credential setup or platform confusion.

**Steps:**

| Step | User Action | System Response |
|---|---|---|
| **1. Invitation email** | Staff receives email: "You've been invited to [Firm Name]'s Axiom workspace as [Role]." | Email includes a magic link (valid 7 days) and a secondary link to set up a password directly. |
| **2. First access** | Staff clicks magic link → authenticated for this browser session. | If firm has OAuth SSO configured, staff is prompted to link their Google or Microsoft account. Otherwise, prompted to set a password. Staff accepts terms of service. |
| **3. Role and profile** | Staff views their assigned role (Partner / Manager / Staff / EQReviewer). Updates display name and notification preferences (email digest: real-time, daily, or weekly). | Role is displayed but cannot be self-changed — FirmAdmin only. Notification preferences are stored on the `User` record. |
| **4. Guided tour** | 5-step in-app tour, skippable at any point. | Highlights in sequence: (1) engagement list, (2) workpaper editor, (3) evidence pool, (4) document request queue, (5) AI review panel. Each step includes a one-sentence description and a "show me" pointer. Tour progress is stored per-user and never repeats after completion or skip. |
| **5. Engagement assignment** | Staff receives an in-app notification when a Partner or Manager adds them to an engagement via `EngagementTeamMember`. | The engagement appears in the staff member's engagement list. The notification links directly to the first assigned control or workpaper. This is the first meaningful platform action. |

**Note on EQReviewer role:** Users with the `EQReviewer` role follow the same onboarding path but land in a read-only view of any engagement they are assigned to review. They cannot be added to the same engagement as a team member — the system enforces this at assignment time (see Section 12, Module 1).

---

### C. Client Onboarding (Per Engagement)

**Goal:** A client contact receives a document request, accesses their secure upload portal without friction, fulfills their requests, and (if needed) delegates specific requests to internal SMEs — all without mandatory account creation.

Clients are onboarded per engagement, not at the firm level. The same client contact may be invited to multiple engagements over time; each invitation is independent. Optional account creation provides organized access across engagements.

**Steps:**

1. **Client contact invitation** — The auditor enters the client's primary contact email when creating the engagement (or later from the Client Hub settings). The system generates a tokenized Client Hub link scoped to this specific engagement and sends an invitation email: "Your audit team has started your [Engagement Name] and needs documents from you."

2. **Magic link first access** — The client contact clicks the tokenized link. No account or password is required for basic upload. They land on the Client Hub showing: engagement name and period, list of outstanding document requests with titles, instructions, and due dates, and a drag-and-drop upload interface per request. The tokenized link is engagement-scoped — it cannot be used to access any other engagement or any firm-level content.

3. **Optional account creation** — For clients with many requests or who expect multi-year engagements, the Client Hub prompts (but does not require): "Create a password-protected account to access all your requests in one place and receive organized notifications." Account creation creates a `ClientUser` record linked to the `Client`. Clients who decline continue to access via tokenized links for each request.

4. **Intake form (configurable per engagement)** — For new clients, the auditor can enable an intake form in the engagement setup. If enabled, the client is prompted to complete the intake form on first access before their request list is shown. The form collects: entity legal name, primary business description, key systems in scope, and primary client contacts. Responses are stored as a `Workpaper` of type `ClientIntake` in the engagement file. The intake form is not required by default — auditors enable it only when needed (common for first-year SOC 2 and ISO 27001 engagements).

5. **Request delegation** — A `ClientAdmin` user (elevated client role, assigned by the auditor) can delegate individual document requests to other contacts within their company. Delegation is by email: the `ClientAdmin` enters a colleague's email address on a specific request, and that colleague receives a tokenized link scoped to that single request only. The delegate sees only the request description, instructions, and upload interface — not the full Client Hub, the engagement name, or any other requests. Delegation creates an `AuditLog` entry recording who delegated, to whom, and when.

6. **Ongoing and post-archive access** — While the engagement is active, the client can view submitted documents and request statuses. After the engagement is archived, the Client Hub link becomes read-only: submitted documents are visible, new uploads are blocked, and an in-portal notice explains that the engagement has been completed. If the auditor has shared the issued report with the client, it is visible in read-only mode in the Client Hub indefinitely (within the engagement's retention period).

**Tokenized link security:**
- All client-side upload tokens expire after 90 days and require re-generation by the auditor
- Tokens are single-use per session (not per upload); a new token is issued on re-generation
- Token re-generation is available from the Document Request settings panel by any `Manager` or `Partner` on the engagement
- Token expiry and re-generation events are recorded in `AuditLog`

**Data entities touched:** `Client`, `ClientUser`, `DocumentRequest`, `EvidenceItem`, `Workpaper` (ClientIntake type), `AuditLog`, `User` (ClientAdmin / ClientUser roles).

---

## 12. Engagement Module Specifications

### Module 1: Engagement Management and Scoping (including SQMS 1 Client Acceptance)

**Purpose:** Initialize an engagement with correct methodology, framework version, team assignment, and quality documentation; enforce the SQMS 1 client acceptance workflow before fieldwork begins.

**Key user flows:**

1. **New engagement creation (from template):** Partner or Manager selects engagement type (FinancialAudit_Private, SOC2, ISO27001, HIPAA, etc.), selects the applicable `MethodologyTemplate`, selects the framework version, enters period dates, and assigns the engagement team. The system creates:
   - `Engagement` record with `status = Planning`
   - `EngagementTeamMember` records for assigned users
   - `EngagementFramework` record(s) — one per framework in scope
   - `Control` records cloned from `TemplateControl` entries (with `prior_control_id` if rollforward)
   - `TestProcedure` records cloned from `TemplateTestProcedure` entries
   - Draft `Workpaper` shells for each workpaper type in the template
   - An empty `ClientAcceptance` record flagged as incomplete

2. **SQMS 1 client acceptance:** Partner completes the `ClientAcceptance` form: quality risks identified (free text + structured categories), firm responses to those risks, independence confirmation, and acceptance sign-off. The Planning → Fieldwork transition is blocked at the state machine level until `ClientAcceptance.accepted_at` is populated by a Partner-role user. The acceptance record is immutable once signed — addenda require creating a new version.

3. **EQR assignment (SQMS 2):** For engagements requiring EQR (all PCAOB engagements; higher-risk nonissuer engagements per firm policy), the `EngagementQualityReview` record is created during engagement setup. The system validates that the assigned `reviewer_id` holds the `EQReviewer` role and is not an `EngagementTeamMember` on the same engagement — if both conditions are not met, the EQR assignment is rejected with a clear error.

4. **Rollforward from prior year:** When creating an engagement from a prior year (`prior_engagement_id` set), the system surfaces: all prior year controls (with rollforward status), prior year workpapers (read-only sidebar), prior trial balance (reference only), and prior ClientAcceptance (read-only — a new acceptance is required). Prior year evidence items are surfaced with "used in prior year" flags when the auditor links evidence to a test procedure.

5. **Scope changes:** If the Partner returns the engagement from Fieldwork to Planning (e.g., a significant scope change requiring re-evaluation of quality risks), the system creates a new `ClientAcceptance` record rather than modifying the prior one. The prior acceptance remains in the audit file with a note that it was superseded.

**Data entities involved:** `Engagement`, `EngagementTeamMember`, `EngagementFramework`, `Control`, `TestProcedure`, `ClientAcceptance`, `EngagementQualityReview`, `MethodologyTemplate`, `TemplateControl`, `TemplateTestProcedure`, `AuditLog`.

**AI touchpoints:**
- Control mapping (Feature 2): immediately after engagement creation, AI proposes `FirmControlObjectiveMapping` records across all frameworks. Partner/Manager reviews in bulk.
- Risk pre-population: AI suggests quality risk categories based on client industry and prior engagement findings. SQMS 1 Tier 2 — auditor reviews and certifies; not auto-populated.

**Regulatory constraints:**
- `ClientAcceptance` must be completed before Fieldwork begins — SQMS 1
- EQR reviewer independence check is system-enforced — SQMS 2 / PCAOB AS 1220
- Framework version must be set at engagement creation and cannot be changed after Fieldwork begins without Partner override and documented reason — from framework version management requirement in Section 4

---

### Module 2: Trial Balance and Financial Analysis (Sheets Experience)

**Purpose:** Provide a spreadsheet-like environment for importing, reviewing, and analyzing the client's trial balance; support account mapping, adjustment tracking, lead schedule generation, and analytical procedures — all without leaving the platform or opening Excel.

**Key user flows:**

1. **Trial balance import:** Staff auditor uploads a CSV or Excel file exported from the client's accounting system (QBO, NetSuite, Sage, Xero). The importer recognizes common column formats (account number, account name, debit balance, credit balance) with configurable column mapping. The result is a set of `TrialBalanceAccount` records linked to the engagement.

2. **AI account mapping (Feature 3):** Immediately after import, Claude Haiku classifies each account into a financial statement line item. Mappings are displayed in the Sheets UI with `mapping_status = AISuggested`. Low-confidence mappings are highlighted. The auditor reviews and confirms each mapping; bulk-confirm is available for accounts where the suggested mapping is unambiguous. Prior year confirmed mappings are pre-loaded on rollforward engagements and treated as the starting suggestion (with `mapping_status = AISuggested` to require re-confirmation for the current year).

3. **Lead schedule generation:** Once accounts are mapped to FS line items, the system automatically generates lead schedule workpapers grouped by financial statement section. The lead schedule aggregates balances from `TrialBalanceAccount` records and supports drill-through to individual accounts. Materiality calculations (ISA 320 / AU-C 320 methods) are available as formula functions.

4. **Adjustment tracking:** Staff auditors propose adjustments (`TrialBalanceAdjustment`). Each adjustment specifies account, amount, description, and type (Proposed | Passed | Waived). Manager or Partner approves or waives. The system tracks both the unadjusted and adjusted trial balance and reflects both in the lead schedule.

5. **Analytical procedures:** The Sheets UI exposes computed analytics: period-over-period variance by account, ratio calculations (current ratio, quick ratio, debt-to-equity), and anomaly flags generated by the AI pipeline (accounts with unusual activity relative to prior period or industry norms). Anomaly flags are Tier 1 AI actions — informational only, not decisions.

6. **Population export for sampling:** Auditors can export the GL transaction detail (if imported) as a population listing for sampling. The platform's sampling calculator (ISA 530 / AU-C 530 methods: systematic, random, monetary unit) accepts the population and produces a `TestProcedure` sample selection record.

**Data entities involved:** `TrialBalance`, `TrialBalanceAccount`, `TrialBalanceAdjustment`, `Workpaper` (lead schedules), `WorkpaperVersion`, `AIDecision`, `AuditLog`.

**AI touchpoints:**
- Account mapping (Feature 3): Haiku classification on import
- Anomaly detection (Tier 1): nightly background job flags unusual account movements for auditor attention; no human approval required to flag, required to act

**Regulatory constraints:**
- Lead schedules are workpapers subject to AU-C 230 assembly deadline and lock requirements — they are locked at the Finalized state
- Every proposed adjustment requires approval documentation per AU-C (all adjustments proposed, whether or not passed, must be retained)
- For PCAOB engagements, all technology-assisted analytical procedures (variance analysis, ratio calculations) must be documented as `AIDecision` records (AS 1105)
- Sampling documentation must record population size, sample size, selection method, and results per AU-C 530 / PCAOB AS 2315

---

### Module 3: Controls and Workpaper Management (including EQR workflow per SQMS 2)

**Purpose:** Execute and document all control testing and workpaper preparation; enforce the sign-off hierarchy from preparer to reviewer to engagement partner; support the EQR workflow for applicable engagements.

**Key user flows:**

1. **Control testing:** Staff auditors are assigned controls via `Control.auditor_assigned_to_id`. For each control, the auditor works through associated `TestProcedure` records: selecting the procedure type, documenting the population and sample (if applicable), linking evidence items, recording results, noting exceptions, and documenting conclusions. Test procedure status progresses: NotStarted → InProgress → Complete (or Exception if exceptions found).

2. **Evidence linking:** When the auditor has evidence to link to a test procedure, they access the engagement's evidence pool (all `EvidenceItem` records for this client, across all engagements) and link the relevant item. The AI may have already suggested a link (`EvidenceLink.ai_suggested = true`); the auditor accepts, modifies, or rejects. On acceptance, the `AIDecision` record is updated with the review action.

3. **Cross-framework evidence display:** When an auditor links evidence to a test procedure, the UI shows which other framework requirements that evidence simultaneously satisfies (via the `FirmControlObjective → FirmControlObjectiveMapping → FrameworkRequirement` chain). For an integrated SOC 2 + ISO 27001 engagement, the auditor sees "this evidence satisfies SOC 2 CC6.1, ISO 27001 A.8.3, and HIPAA §164.312(a)(1)." No additional action required — all mappings are populated automatically.

4. **Workpaper sign-off workflow:** After a test procedure is complete, the auditor marks the workpaper as `PreparedPendingReview`. The manager is notified. The manager reviews the workpaper (in-platform, with inline commenting and review notes), marks review notes as open or resolves them, and advances the workpaper to `ReviewComplete`. The engagement partner reviews the manager's cleared notes and signs off the workpaper. Sign-off is a timestamped, named action that creates an `AuditLog` entry — it cannot be backdated or bulk-applied without individual confirmation.

5. **AI workpaper draft (Feature 4):** Once a test procedure is marked Complete, the auditor can request an AI narrative draft. The draft appears in the workpaper editor labeled "AI Draft — requires review." The `WorkpaperVersion` record is created with `is_ai_draft = true`. The auditor must edit the text and save (changing `is_ai_draft` to false) before the workpaper can be advanced to `PreparedPendingReview`. The platform will not allow a workpaper with `is_ai_draft = true` to be signed off.

6. **EQR workflow (SQMS 2 / PCAOB AS 1220):** For engagements with EQR assigned, the EQR reviewer accesses the engagement in read-only mode (they are not an `EngagementTeamMember` and cannot modify any engagement content). The reviewer documents their review scope, findings, and conclusion in the `EngagementQualityReview` record. The Review → Reporting transition is blocked until `EngagementQualityReview.signed_off_at` is populated by the assigned reviewer. The EQR record and its sign-off timestamp become part of the immutable engagement archive.

7. **Review notes management:** Review notes (comments and findings from manager or EQR review) are tracked as structured records linked to the relevant workpaper or control. Open review notes block workpaper advancement. Resolved notes remain in the record — they cannot be deleted. The AuditLog captures the open and resolution events.

**Data entities involved:** `Control`, `TestProcedure`, `EvidenceItem`, `EvidenceLink`, `Workpaper`, `WorkpaperVersion`, `EngagementQualityReview`, `AIDecision`, `AuditLog`, `FirmControlObjective`, `FirmControlObjectiveMapping`, `FrameworkRequirement`.

**AI touchpoints:**
- Evidence link suggestions (Tier 2): AI suggests relevant evidence items for each test procedure based on the procedure description and evidence extracted text
- Workpaper narrative draft (Feature 4, Tier 2): on-demand, explicit request by auditor
- Control mapping (Feature 2): surfaces cross-framework satisfaction at evidence-link time

**Regulatory constraints:**
- Sign-off hierarchy must be enforced at data layer: workpaper cannot advance states out of order — SQMS 1, AU-C 220
- EQR reviewer must be independent of engagement team — SQMS 2, PCAOB AS 1220
- All AI suggestions that affect audit content must create `AIDecision` records — PCAOB AS 1105
- `is_ai_draft = true` workpapers cannot be signed off — PCAOB AS 1105 (AI-generated vs. auditor-authored content distinction)
- Review notes cannot be deleted — AU-C 230

---

### Module 4: Document Request and Client Hub (PBC Portal)

**Purpose:** Manage the complete cycle of PBC (Provided By Client) document requests from creation through fulfillment; provide clients with a clean, no-login-required upload experience; apply AI completeness review to uploaded documents before they enter the auditor's review queue.

**Key user flows:**

1. **Request creation:** Staff auditors or managers create `DocumentRequest` records linked to specific controls or test procedures. Requests include a title, detailed instructions (what to provide, the format required, the period to cover), a due date, and assignment to a client-side contact. Bulk request creation from a methodology template is available — the standard SOC 2 Type II template creates 80+ pre-drafted requests covering all trust services criteria.

2. **Client notification and upload portal:** On request creation (or on explicit send action), the system sends the client contact an email with a tokenized link to their secure upload portal. No client account or password is required for basic uploads — the tokenized link is sufficient. For multi-upload engagements, clients can optionally create a password-protected `ClientUser` account for organized access to their request list.

3. **Client upload experience:** The client portal presents the client's outstanding requests with descriptions, instructions, and due dates. For each request, the client uploads a file (drag-and-drop or file picker; single or bulk). The upload is stored as an `EvidenceItem` record immediately. The request status changes to `Submitted`.

4. **AI completeness review (Feature 1):** On upload, an asynchronous job queues the AI completeness review. Within minutes, Claude Sonnet analyzes the uploaded document against the request requirements and produces an `AIDecision` record with a recommendation: Accept | Request Clarification | Reject. The auditor receives an in-app notification: "AI has reviewed [document name] for [request title]."

5. **Auditor review queue:** The auditor sees all recently reviewed documents in a queue, sorted by AI confidence (low-confidence reviews surfaced first). For each review, the auditor sees the AI recommendation, the specific gaps identified (if any), and one-click actions: Accept (creates `EvidenceLink` to relevant `TestProcedure`), Send Back to Client (with auto-drafted client explanation from the AI's gap analysis), or Reject (remove from engagement). All three actions update the `AIDecision.review_action` field and create `AuditLog` entries.

6. **Automated reminders:** Overdue document requests trigger the `DocumentRequestReminderStateMachine` in Step Functions. The workflow sends reminder emails on a configurable schedule (default: 7 days before due, on due date, 7 days after). After three reminders, the workflow escalates to the auditor with a notification. Reminder frequency and content are configurable per engagement.

7. **Overdue management:** Requests that are overdue change status to `Overdue` automatically. The engagement dashboard surfaces overdue requests with a count badge. Partners can see overdue request rates across all active engagements in the analytics dashboard (Scale tier).

**Data entities involved:** `DocumentRequest`, `EvidenceItem`, `EvidenceLink`, `TestProcedure`, `Control`, `AIDecision`, `AuditLog`, `User` (ClientUser).

**AI touchpoints:**
- Document completeness review (Feature 1, Tier 2): triggered on every client upload; auditor must review and action the recommendation
- Evidence link suggestion (Tier 2): on acceptance, AI also suggests which `TestProcedure` to link the evidence to

**Regulatory constraints:**
- Period coverage check is a mandatory part of the AI completeness review for SOC 2 Type II — AT-C 320 requires evidence covering the full examination period
- Every document request acceptance creates an `AIDecision` record — PCAOB AS 1105 where applicable
- Client-side upload tokens expire after 90 days and require re-generation; tokens are single-use for the specific engagement (cannot be used to upload to a different engagement)
- Document request overdue reminders are audit-trailed in `AuditLog` (reminder sent, who was notified, when)

---

### Module 5: Reporting and Archiving (Immutable Archive, Assembly Deadline Enforcement)

**Purpose:** Generate the final engagement report; manage the transition from active engagement to finalized and then to immutable archived state; enforce assembly deadlines and retention schedules.

**Key user flows:**

1. **Report generation:** When the engagement reaches the Reporting state, the partner generates the report from a template. Report types: SOC 2 Type I, SOC 2 Type II, SOC 1 Type I, SOC 1 Type II, Financial Audit Opinion, Agreed-Upon Procedures, Management Letter. The report template is pre-populated with engagement data (client name, period, framework, controls summary, exception summary). The partner edits the narrative, adds the opinion, and iterates on the draft internally. Each save creates a `ReportVersion` record.

2. **Client review (optional):** The partner can share a draft report with client admin users for review. Client users see a read-only view of the draft report and can submit comments. The partner reviews comments in-platform and issues a revised draft.

3. **Report issuance:** The partner marks the report as Issued (`Report.status = Issued`) and the system records `report_issued_at` on the `Engagement`. This triggers: (a) computation of `assembly_deadline` (report date + 60 days for AICPA, + 45 days for PCAOB), (b) computation of `retention_deadline` (report date + 5 years for AICPA, + 7 years for PCAOB), (c) scheduling of the Finalized → Archived Step Functions wait state.

4. **Finalized state:** The partner transitions the engagement to Finalized. This transition is blocked unless `Report.status = Issued`. At Finalized: all workpaper content is locked (`Workpaper.is_locked = true`); no further edits to any workpaper, evidence link, control conclusion, or test procedure are permitted. Any modification attempt returns a hard error: "This engagement has been finalized. Modifications require an addendum."

5. **Addendum process (post-finalization):** If an error or omission is identified after finalization and before archival, the auditor creates a new `WorkpaperVersion` with `is_addendum = true`, documents the reason for the addendum, and obtains sign-off from the engagement partner. The original content is unchanged — the addendum is an additional record, not a modification. This implements AU-C 230 §.16 (subsequent discovery of information after the documentation completion date must be documented as an addendum).

6. **Automatic archival:** The `EngagementLifecycleStateMachine` in Step Functions triggers the Finalized → Archived transition when `report_issued_at + assembly_window` has elapsed. At archival:
   - All workpaper and evidence files are copied to the S3 Object Lock bucket with the engagement's `retention_deadline` set as the Object Lock retention period (COMPLIANCE mode)
   - `Engagement.archived_at` is set
   - The engagement becomes read-only in the application — no further state changes are possible (except the FirmAdmin → Abandoned path)
   - An archival confirmation email is sent to the engagement partner
   - The `AuditLog` records the system-triggered archival event

7. **Retention expiry and deletion:** When `retention_deadline` elapses (5 or 7 years from report date), the Object Lock period expires automatically. The system sends the firm admin a notification 90 days before expiry: "Engagement [name] reaches its retention limit in 90 days. Data will be deleted on [date] unless you export it." If no extension is requested, S3 lifecycle policies delete the objects after expiry. This implements the obligation to delete at schedule end per GDPR/CCPA retention obligations.

8. **Engagement export:** At any time (including before archival), firm admins can generate a complete engagement export: all workpapers as PDF, all evidence files in native format, trial balance as Excel, AuditLog as CSV, and engagement metadata as JSON. The export is a structured ZIP file suitable for standalone archival independent of the platform. This export is the primary mechanism for satisfying offboarding obligations and for firms that want their own long-term archival copies.

**Data entities involved:** `Engagement`, `Report`, `ReportVersion`, `Workpaper`, `WorkpaperVersion`, `EvidenceItem`, `EvidenceLink`, `AuditLog`, `EngagementQualityReview`, `AIDecision`.

**AI touchpoints:**
- Report narrative drafting: AI can be invoked on the report template to pre-draft the "Description of Tests of Controls" section from the control testing data; same Tier 2 rule applies (draft labeled, human must edit and sign off)
- No other AI touchpoints in reporting — the partner is performing final professional judgment at this stage

**Regulatory constraints:**
- Assembly deadline is computed at report issuance and enforced as a hard lock — AU-C 230, PCAOB AS 1215
- S3 Object Lock COMPLIANCE mode is used for archived engagements — PCAOB AS 1215 and SOX §802 immutability requirements
- Retention periods are per engagement type (5 years for private/SOC/ISO; 7 years for public/PCAOB; 6 years for HIPAA) and computed from `report_issued_at`
- Addenda post-finalization require documented reason and partner sign-off — AU-C 230 §.16
- The Finalized → Archived transition is system-triggered and cannot be overridden by any user role (including FirmAdmin) except to mark an engagement as Abandoned — which also triggers archival

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
| **Multi-agent AI orchestration** | LangGraph, CrewAI, and similar frameworks add debugging overhead without adding value for the four defined AI features, all of which are single-step or multi-step-with-human-gates. Add orchestration when genuinely needed by a specific feature, not as a framework choice. |
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
