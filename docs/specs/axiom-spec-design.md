# Axiom Product Specification
**Date:** April 17, 2026
**Status:** Implementation-Ready (compliance/assurance pivot)

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

> Research: [01-target-market.md](../research/01-target-market.md) <!-- STALE: research doc predates compliance pivot --> · [compliance-pivot-findings.md](compliance-pivot-findings.md)

### What Axiom Is

Axiom is an AI-native **audit management and evidence collection** platform built for mid-market CPA firms, compliance-first consultancies, accredited ISO Certification Bodies, and PCI QSA firms that deliver SOC 2, ISO 27001, ISO 27701, ISO 42001, HIPAA, PCI DSS, and SOC 1 engagement work. It is deliberately **both-sided**: the auditor-side workspace (engagement scoping, fieldwork, workpapers, reporting) and the auditee-side workspace (Client Hub with continuous monitoring, policy library, drift-triggered re-testing) operate on the same `CommonControl` graph, so one evidence artifact clears requirements across every in-scope framework simultaneously. The AI layer is ISO 42001-native by construction and every AI output is cryptographically signed and WORM-stored — auditor-defensible by the data path, not by marketing claim. Axiom is designed to be operational within one week of signup, without an implementation consultant.

### Target ICP

**Firm type:** Mid-market CPA firms delivering AICPA SOC attestations + compliance-first consultancies (the tier below Schellman / A-LIGN) delivering ISO readiness, PCI gap assessment, HIPAA/HITRUST, and ISO 42001 engagements.

**Firm size:** 20–60 professional staff.

**Practice mix:** Compliance and assurance work across SOC 2 Type I/II, ISO 27001:2022, ISO 27701:2019, ISO 42001:2023, HIPAA Security Rule (with HITRUST CSF r2 path post-MVP), PCI DSS v4.0.1, and SOC 1 Type I/II. **In scope as customer firms:** AICPA-licensed CPA firms delivering SOC attestations, compliance-first consultancies delivering ISO readiness / ISO 42001 engagements, **accredited ISO Certification Bodies (CBs) under ISO 17021-1**, **PCI Qualified Security Assessor (QSA) firms**, and firms delivering HIPAA / HITRUST-r2 (post-MVP) engagements. Axiom serves the engagement-delivery and evidence-collection workflow for these firms; the legal sign-off authority (CB-issued certificates, QSA-signed ROCs, AICPA-licensed attestation opinions) remains with the licensed firm. **Out of scope:** internal audit / SOX / enterprise GRC (AuditBoard's territory) and financial audit.

**Geography:** US and Canada at launch. AWS us-east-1.

**Engagement volume:** 30–100 engagements per year.

**Decision-maker profile:** Managing partner, practice lead, or head of attestation services who is either (a) paying $30K+/year for Drata/Vanta/Hyperproof and realizing these tools are auditee-facing with bolted-on auditor features, (b) struggling through multi-framework engagements where evidence is re-collected per framework, (c) selling ISO 42001 engagements with no adequate tooling, or (d) scrutinizing AI-assisted work after the March 2026 category trust crisis and asking "how do we prove our AI outputs are defensible?"

### Primary and Secondary Differentiators

**Primary differentiator:** Axiom is a **both-sided (auditor + auditee) AI-native compliance platform with STRM-grade cross-framework evidence mapping, cryptographic AIDecision provenance, and ISO 42001-native human-in-the-loop AI governance.** Each of those is concrete:

- **Both-sided** — auditor-side tooling (engagement lifecycle, workpapers, multi-level sign-off, reporting, archival) and a full auditee GRC workspace (Client Hub continuous monitoring, policy library, drift-triggered re-testing) operate on the same data model. Fieldguide, Agentive, and Yak are auditor-only. Drata / Vanta / Secureframe / Delve are auditee-only. None cover both.
- **STRM-grade cross-framework evidence mapping** — the `CommonControl` graph uses NIST Secure Controls Framework (SCF) crosswalks with NIST STRM relationship vocabulary (`equivalent-to | subset-of | superset-of | intersects-with | no-relationship`), strength scores, and effective-dated edges. One evidence artifact flows coverage through the `EvidenceItemSupports → CommonControl → CommonControlSatisfies → FrameworkRequirement` chain to every mapped framework, with period-aware staleness checks (ASV scan 90d, pen test 1y, background check 1y, SOC 2 Type II windows, ISO surveillance cycles) applied per-framework. Axiom never shows a green checkmark on partial coverage — percentages and gap lists only.
- **Cryptographic AIDecision provenance** — every AI output is signed with an AWS KMS asymmetric key (`ECC_NIST_P256`, `SIGN_VERIFY`) at emission and written to S3 Object Lock (WORM) with the provenance envelope (`artifact_id`, `model_id`, `model_version`, `prompt_hash`, `input_content_hashes`, `output_hash`, `generated_at`, reviewer identity, reviewer action). Any party — EQR reviewer, ISO CB technical reviewer, client procurement, regulator — can recompute and verify the signature. This is the direct answer to the March 2026 agentic-compliance trust crisis.
- **ISO 42001-native HITL AI governance** — Axiom dogfoods ISO 42001. The three-tier HITL policy, `AIDecision` ledger, and `ai_content_metadata` tracking are the operational realization of ISO 42001 clauses 6–9 for Axiom as a deployer of AI. Impact-assessment records, model-change-management logs, prompt/model/context/confidence/reviewer/override auditability are first-class platform features, not a policy PDF.

**Secondary differentiator:** Self-serve onboarding with time-to-first-engagement under one week. No implementation consultant required for standard SOC 2, ISO 27001/27701, HIPAA, ISO 42001, PCI DSS, or SOC 1 methodologies. This is a structural sales advantage against auditor-side incumbents (Fieldguide enterprise sales motion) and against GRC adjacencies (AuditBoard, Hyperproof) whose ACVs require sales-led procurement cycles. Axiom's integration roadmap (cloud, identity, dev tools, evidence-bridge ingestion from Drata / Vanta / Sprinto / Hyperproof) makes it a strong contender in the **audit management and evidence collection** category positioning, not just in the auditor-workspace niche.

### What Axiom Does NOT Do at Launch

- **Issue ISO certificates on behalf of a Certification Body** — accredited CBs are valid customer firms and Axiom supports their engagement-delivery workflow, including generating the certificate document from a template. The legal certification decision and signature remain with the CB under ISO 17021-1; Axiom does not become an accreditation body and does not replace CB technical-review authority.
- **Sign PCI ROCs on behalf of a QSA firm** — accredited QSA firms are valid customer firms and Axiom supports their gap-assessment, evidence-collection, and ROC-assembly workflow, including generating the ROC and AOC from templates. Legal ROC/AOC sign-off remains with the QSA under PCI SSC accreditation; Axiom does not become a QSA.
- **PCAOB public-company audits** — Axiom does not support issuer audits under PCAOB standards.
- **Financial audit / GAAS** — no trial balance, sampling, materiality, or financial-statement workpapers. This is a deliberate exit from the audit vertical.
- **Internal audit / SOX compliance / enterprise GRC** — AuditBoard's market; different buyer, different workflow, different regulatory framework (IIA vs. AICPA / ISO).
- **ESG or sustainability reporting**.
- **White-labeling or reseller channels**.
- **Custom AI model training per firm** (insufficient data, managed via RAG instead).
- **On-device or local model execution** (Bedrock over VPC endpoint is sufficient).
- **AI-autonomous engagement actions** — AI suggests, humans decide; ISO 42001 Tier 3 actions are human-initiated only.

---

## 2. Competitive Positioning

> Research: [02-competitive-differentiation.md](../research/02-competitive-differentiation.md) <!-- STALE: research doc predates compliance pivot --> · [compliance-pivot-findings.md](compliance-pivot-findings.md)

### One-Paragraph Positioning

Axiom's primary competition is **auditor-side**: **Fieldguide** (Series-funded AI-native audit-engagement platform; the most capitalized auditor-AI player in the market), **Agentive** (YC S23, ~$500K seed, ~6 people; narrow workpaper-drafting surface with no framework intelligence), and **Yak** (auditor-side AI tooling; smaller and less established). **AuditBoard CrossComply** and **Hyperproof** are GRC-side adjacencies that license UCF / pivot to SCF for crosswalks but carry internal-audit baggage that mis-fits mid-market attestation firms. **Drata, Vanta, and Secureframe** sit in the *auditee* category and are not the primary comparison set — they automate control implementation and evidence collection for the company being audited and hand off a report package to whichever external auditor the client brings; Axiom appears alongside them on the auditee-side workspace (Client Hub) but the buying decision happens on the auditor-side workspace, where they have no product. **Sprinto** and **Thoropass** are mid-market auditee SaaS, with Thoropass owning an audit practice (an independence concern for external-auditor buyers). The March 2026 trust crisis in the agentic-compliance category made unverifiable AI output commercially toxic, and no competitor today produces cryptographically signed AI outputs with a reviewer-attested ledger. Axiom is the only platform that (a) serves both auditor and auditee sides on one data model, (b) rides SCF/OSCAL/AICPA/CIS crosswalks with STRM-grade period-aware edges, (c) signs every AI output at emission and WORM-stores it, and (d) is ISO 42001-native by construction — and does so at mid-market self-serve pricing.

### Competitor Set

The competitor set is grouped by where the buying decision sits. Axiom's primary competition is on the **auditor side**; auditee-side and GRC-adjacent platforms are secondary.

**Primary — auditor-side platforms (the actual head-to-head):**

| Competitor | Scope | Strengths | Gap vs. Axiom |
|---|---|---|---|
| **Fieldguide** | Auditor-side AI audit-engagement platform | Best-funded auditor-AI player; strong workpaper automation and PBC workflow; established mid-market traction | No both-sided workspace (no auditee-facing Client Hub / continuous monitoring); no STRM-grade cross-framework graph; no cryptographic AI provenance / WORM at emission; not ISO 42001-native by construction |
| **Agentive** (YC S23) | Auditor-only workpaper drafting | Narrow PBC/workpaper UX | Tiny (~$500K funding, ~6 people); no framework intelligence, no cross-mapping, no both-sided surface, no signed AI provenance |
| **Yak Tech** (yaktech.io) | Auditor-side AI audit-lifecycle platform; SOC 1 and SOC 2 only | Founded 2024, ~7 employees, Portland OR. Founder Shannon Smith (CISA); CTO Saxon D'Aubin. "Built by auditors for auditors." Features: AI Reviewer (document evaluation), one-click report generation, control mapping, evidence request workflow, client communication. Won 2025 Accounting Today Top New Product and 2025 Columbia River Pitch Best Overall | **SOC 1/2 only** — no ISO 27001, ISO 27701, ISO 42001, HIPAA, or PCI DSS support; no STRM-grade cross-framework mapping (single-framework architecture); no auditee-side workspace / Client Hub for continuous monitoring; no cryptographic AI provenance; not ISO 42001-native; capital-constrained relative to Fieldguide |

**Secondary — GRC adjacencies (mis-fit on workflow, not direct competitors but show up in evaluations):**

| Competitor | Scope | Strengths | Gap vs. Axiom |
|---|---|---|---|
| **AuditBoard CrossComply** | Enterprise GRC + UCF license | Strongest internal-audit pedigree, UCF crosswalk | Internal-audit workflow mis-fit for mid-market attestation; enterprise sales cycle |
| **Hyperproof** | GRC with 142 frameworks | SCF-pivot crosswalk, breadth | Internal-audit tilt; no both-sided product; no signed provenance |

**Tertiary — auditee-side platforms (not the buying decision; appear alongside Axiom's Client Hub):**

| Competitor | Scope | Strengths | Gap vs. Axiom |
|---|---|---|---|
| **Drata** | Auditee-only | 30+ frameworks, DCF proprietary control framework, strong SOC 2 auditor hand-off | No auditor workspace, no cross-mapping UX for external auditors, no signed AI provenance |
| **Vanta** | Auditee-only | 35+ frameworks, strongest enterprise adoption, Vanta Control Framework (VCF) launched March 2026 | Auditee-only; enterprise pricing; AI features unsigned |
| **Secureframe** | Auditee-only | 40+ frameworks, similar depth to Vanta | Auditee-only; weaker enterprise story |
| **Sprinto** | Mid-market auditee | "Magic Mapping", fast-to-SOC-2 | Auditee-only |
| **Thoropass** | Mid-market auditee + audit practice | Bundled audit + tooling | Owns audit practice — independence concerns for external-auditor buyers |

**Category context — March 2026 agentic-compliance trust crisis:** Across all three tiers, no competitor today produces cryptographically signed AI outputs with a reviewer-attested ledger. The crisis makes unverifiable AI output commercially untenable; Axiom's response is signed, reviewer-attested `AIDecision` records. Axiom does not name the crisis-implicated competitor in-product.

### The Switching Trigger

A firm evaluates Axiom when one or more of the following occurs:

1. **Auditor-side platform limits.** A firm is using Fieldguide, Agentive, or Yak and hits a ceiling — no cross-framework evidence graph, no auditee-facing Client Hub for continuous monitoring, no cryptographic AI provenance. A multi-framework engagement or a post–March-2026 procurement review surfaces these gaps.
2. **Failed audit on evidence gaps.** A firm's engagement goes sideways because evidence staleness, period coverage, or partial-satisfaction gaps were not caught. A platform that surfaces these automatically is an easy budget conversation.
3. **Continuous-monitoring drift not caught.** An incumbent tool flagged a drift but didn't tie it to re-testing; the auditor only discovered during fieldwork that a control had silently broken mid-period.
4. **Multi-framework expansion.** A firm picking up a SOC 2 + ISO 27001 + ISO 27701 integrated engagement, or adding ISO 42001 to the offering, needs cross-mapping and cannot get it cleanly from an auditor-side tool that lacks a STRM graph.
5. **Post-crisis AI scrutiny.** After the March 2026 category trust crisis, a partner asks "how do we prove our AI-assisted work is defensible?" Signed provenance + AIDecision ledger is the answer.
6. **Auditee-side renewal review (secondary).** A firm receives a Drata / Vanta / Hyperproof renewal at $50K+/year and re-examines whether they are paying for features they actually use. This trigger now matters mainly for firms that bought an auditee tool and tried to retrofit it for auditor workflow; Axiom's both-sided model means they don't have to.

The 14-day trial exploits triggers 1 and 4: a firm can sign up, load a real SOC 2 Type II template, upload one piece of evidence, and watch the cross-framework coverage dashboard update live — before committing to any payment.

### Why Auditor-Side Incumbents Won't Close the Gap Quickly

**Fieldguide** is well-capitalized but has built around the auditor workflow as the singular surface. Adding an auditee-facing Client Hub with continuous monitoring and drift-triggered re-testing is a separate product motion with separate buyer behavior; the cross-framework STRM graph and cryptographic provenance are foundational data-model decisions, not features that bolt on. A retrofit takes a year-plus of platform work that the team hasn't signaled.

**Agentive** is capital-constrained (~$500K seed, ~6 people), with no framework intelligence and a narrow auditor-side surface. They can grow within PBC/workpaper drafting but cannot cross into cross-mapping + auditee workspace without a ground-up rebuild.

**Yak Tech** is small (~7 employees, founded 2024) and currently SOC-only. To compete on Axiom's multi-framework footprint they would need to build out ISO 27001/27701/42001, HIPAA, and PCI DSS framework intelligence, plus a STRM-grade cross-mapping graph, plus an auditee-side workspace — multi-year platform work for a 7-person team. The auditor-AI category as a whole has not yet absorbed the post–March-2026 provenance requirement.

**Drata and Vanta** sit in the auditee category and won't respond downmarket-auditor-side. Drata is a unicorn ($2B+ valuation) and Vanta is the category leader; both are raised to sustain enterprise-auditee growth. Their strategic incentive is to expand **upmarket within the auditee segment** — larger customers, more frameworks, more integrations, deeper verticals — not to launch an auditor-side workspace that would (a) require a separate product motion, (b) create a conflict of interest with their installed base (they serve the auditees their auditor customers audit), and (c) dilute their core message. A both-sided pivot by an incumbent would require organizational duplication that enterprise SaaS companies at this stage rarely execute successfully.

---

## 3. Pricing and Business Model

> Research: [09-pricing.md](../research/09-pricing.md)

### Recommended Tiers

**Growth — $1,200/month ($14,400/year, or $12,000/year billed annually)**

Best for: 10–40 staff firms, 20–50 engagements/year, one to two frameworks.

- Unlimited users
- Up to 35 active engagements/year included
- **Framework templates at launch:** SOC 2 Type I/II, ISO 27001:2022, ISO 27701:2019, HIPAA, PCI DSS SAQ, SOC 1 Type I/II
- Client Hub with continuous monitoring (auditee workspace)
- AI evidence → `CommonControl` mapping suggestions
- Standard AI: document extraction, evidence completeness review, workpaper draft assist, report section draft
- Pre-built methodology templates (AICPA SOC 2 TSC, ISO 27001 Annex A, HIPAA Safeguards, PCI DSS requirements)
- Tokenized Client Hub upload links (no client login required for simple upload) + full ClientAdmin workspace
- Standard email support + in-app help
- Overage: $350/engagement beyond included 35

**Scale — $2,400/month ($28,800/year, or $24,000/year billed annually)**

Best for: 30–100 staff firms, 50–150 engagements/year, running multi-framework integrated engagements.

- Everything in Growth
- Up to 100 active engagements/year included
- **Cross-framework evidence mapping** (one evidence item flows coverage through the STRM graph to every mapped framework — note: this is a tier feature; the primary product differentiator is the both-sided workspace with provenance, of which cross-mapping is one component)
- AI-assisted gap analysis and cross-framework risk assessment
- **Drift-triggered re-testing** (configuration-drift detection + autonomous re-test with AIDecision ledger)
- **ISO 42001 engagement support** (AI system inventory, impact-assessment templates, model-card management)
- **PCI DSS ROC preparation** (gap assessment, evidence binder assembly for QSA sign-off)
- Advanced analytics dashboard (engagement cycle time, staff utilization, review bottleneck identification, coverage freshness)
- Multi-entity / group audit support
- Custom methodology template editor
- Priority support + dedicated onboarding call
- Overage: $250/engagement beyond included 100

**Enterprise — Custom ($50,000–$120,000/year, negotiated)**

Best for: 80–200+ staff firms, 150+ engagements/year, requiring enterprise security, SLAs, and custom integrations.

- Everything in Scale, unlimited engagements
- Dedicated Customer Success Manager
- 99.9% uptime SLA
- SAML/SSO integration
- Signed HIPAA BAA (dedicated BAA review and negotiation)
- Audit trail export and long-term data retention configuration
- Custom integrations (firm-specific connectors, enterprise identity providers, bespoke reporting pipelines)
- Security review package (Axiom's own SOC 2 Type II report, ISO 27001 certificate, ISO 42001 certificate, penetration test summary, data residency documentation)
- **Optional SKU:** UCF (Unified Compliance Framework) license passthrough for firms that require UCF content in addition to the platform's default SCF/OSCAL/AICPA/CIS stack
- Negotiated multi-year contracts; EU/APAC data residency option

**Annual billing incentive:** 2 months free (16.7% discount) for annual prepayment. Monthly billing available at Growth tier only. Scale and Enterprise require annual contracts.

**Engagement definition:** One engagement = one client entity, one primary framework, one period. A rollforward to the following period counts as a new engagement. **An integrated engagement (e.g., SOC 2 + ISO 27001 + ISO 27701 scoped together per Journey 11) counts as one engagement**, not three — the whole value proposition of cross-mapping is defeated if the pricing model re-counts it. This definition is stated in the terms of service and enforced in the product UI.

### Self-Serve vs. Sales-Assisted Thresholds

| Tier | Signup Flow | Sales Touch |
|---|---|---|
| Growth | Self-serve: credit card or invoice, 14-day trial, no mandatory call | Optional: automated email Day 3, human call offer Day 10 |
| Scale | Self-serve trial start (14 days), then inside sales to close | Inside sales rep follows up Day 3 and Day 10; 30-min call |
| Enterprise | Request demo → assigned AE → custom quote | Full sales-led, 4–8 week cycle, security review package |

### Trial Strategy

14-day full-feature trial, no credit card required. Trial requires a business email domain and brief intake form (firm name, staff count, primary engagement types: multi-select from SOC 2 / ISO 27001 / ISO 27701 / ISO 42001 / HIPAA / PCI DSS / SOC 1). One pre-loaded engagement template on signup — **SOC 2 Type II or ISO 27001:2022, selectable**. Trial workspace is locked (not deleted) at expiry until a plan is selected. AI-driven in-app progress guidance toward completing the first engagement setup.

No permanent free tier. The trust and professional-perception risk of a free tier outweighs the acquisition benefit for a platform handling sensitive compliance data and cryptographically signed evidence.

### Revenue Model Summary

| Milestone | Firm Mix | Projected ARR |
|---|---|---|
| **50 firms** | 30 Growth, 15 Scale, 5 Enterprise | ~$1.16M ARR |
| **200 firms** | 90 Growth, 80 Scale, 30 Enterprise | ~$5.4M ARR |
| **500 firms** | 175 Growth, 225 Scale, 100 Enterprise | ~$15M ARR |

With overage revenue (5–15% of ACV) and NRR of 110–120%, realistic ARR at 500 firms is $16–18M. Gross churn target: 8–12% annually. Primary churn risk is not product dissatisfaction but firm-level business events (mergers, practice consolidation, champion turnover). Historical compliance/assurance software has very low churn because engagement history and regulatory retention obligations create a strong switching cost.

---

## 4. Regulatory Compliance Requirements

> Research: [03-regulatory-standards.md](../research/03-regulatory-standards.md) <!-- STALE: research doc predates compliance pivot --> · [compliance-pivot-findings.md](compliance-pivot-findings.md)

### Engagement Type × Standard × Platform Requirements

| Requirement | SOC 1 | SOC 2 | ISO 27001 | ISO 27701 | ISO 42001 | HIPAA | PCI DSS |
|---|---|---|---|---|---|---|---|
| **Governing standard** | AICPA AT-C 105/205/320 (SSAE 18) | AICPA AT-C 105/205/320 (SSAE 18) | ISO/IEC 27001:2022 | ISO/IEC 27701:2019 | ISO/IEC 42001:2023 | HIPAA Security Rule (45 CFR §§164.302–318) | PCI DSS v4.0.1 |
| **Sign-off authority** | CPA firm partner + EQR (SQMS 2 where applicable) | CPA firm partner + EQR (SQMS 2 where applicable) | Accredited CB (ISO 17021-1) — Axiom supports readiness; CB issues certificate | Accredited CB — Axiom supports readiness | Accredited CB — Axiom supports readiness; dogfooded internally | Firm attestation or internal audit; HITRUST r2 via Authorized External Assessor (post-MVP) | QSA (accredited by PCI SSC) signs ROC; Axiom supports gap assessment and evidence assembly |
| **Period coverage** | Type I: point-in-time as-of date (single date). Type II: continuous examination period of **3 to 12 months** | Type I: point-in-time as-of date (single date). Type II: continuous examination period of **3 to 12 months** | Certification cycle (initial + annual surveillance + triennial full) | Certification cycle + privacy-program surveillance | Certification cycle; AI-lifecycle management-system scope | Ongoing compliance point-in-time | Annual assessment; ASV scans quarterly; pen test annually |
| **Retention period** | 5 years from report date | 5 years from report date | Ongoing ISMS records | Ongoing PIMS records | Ongoing AIMS records | 6 years from creation / last effective date | 3 years minimum (some evidence longer) |
| **Workpaper assembly deadline** | 60 days after report issuance | 60 days after report issuance | N/A (ongoing ISMS) | N/A | N/A | N/A | N/A |
| **Immutable lock trigger** | After 60-day assembly window closes | After 60-day assembly window closes | On certificate issuance | On certificate issuance | On certificate issuance | Policy-driven point-in-time | On ROC issuance |
| **AI documentation required** | Best practice | Best practice | Best practice (ISO 42001 if the org is deploying AI) | Best practice | **Mandatory** — ISO 42001 itself requires documented AI lifecycle, impact assessments, HITL, model change management | Best practice | Best practice; required if AI is used in CDE |

### Explicit Scope Boundary

Axiom is the engagement-delivery and evidence-collection platform for licensed firms; it is **not** an accreditation authority. Specifically:

- **Axiom does not issue ISO certificates.** ISO certification is a Certification Body (CB) function under ISO 17021-1. Accredited CBs are valid customer firms — Axiom supports their engagement workflow and **generates the certificate document from a template** as a deliverable artifact, but the certification decision and signature remain with the CB.
- **Axiom does not sign PCI ROCs or AOCs.** ROC/AOC sign-off is a PCI SSC–accredited QSA function. QSA firms are valid customer firms — Axiom supports their gap-assessment, evidence-collection, and ROC-assembly workflow and **generates the ROC and AOC documents from templates** as deliverable artifacts, but the legal sign-off remains with the QSA.
- **Axiom does not issue attestation opinions on behalf of any firm.** Attestation opinions remain the licensed firm's responsibility (AICPA-licensed CPA firms for SOC attestations).

Customer firm types in scope: AICPA-licensed CPA firms for SOC attestations, compliance-first consultancies for ISO readiness and ISO 42001 engagements, accredited ISO Certification Bodies for ISO certification engagements, QSA firms for PCI gap assessment / ROC assembly, and firms delivering HIPAA/HITRUST-r2 (post-MVP) engagements. This boundary is also stated in the MSA and in product in-app.

### Firm Quality Management + Internal Inspection

Firms running SOC engagements fall under AICPA SQMS 1 (effective December 15, 2025) and SQMS 2 (for EQR). Firms running ISO engagements fall under ISO 17021-1 §9.6 internal-audit expectations for CBs (Axiom does not become a CB but its consultancy customers produce readiness artifacts a CB will later review). Firms running PCI engagements fall under PCI SSC QSA QA requirements.

The platform supports a **generic firm quality management + internal inspection workflow** that applies across all of these:

- **Client/engagement acceptance:** On engagement creation, the platform requires completion of a `ClientAcceptance` record documenting quality risks and firm responses. Planning → Fieldwork is blocked until this record is signed off by a partner (or partner-equivalent for consultancies).
- **Firm quality policies:** Firms maintain policies in the platform; applicable policies surface inline during engagement setup.
- **Engagement Quality Review (EQR):** For SOC attestations requiring EQR, the platform enforces reviewer assignment, reviewer independence (the EQR reviewer cannot be an `EngagementTeamMember`), and blocks the Review → Reporting transition until EQR sign-off. For ISO engagements, an analogous internal-review workflow applies.
- **Partner involvement evidence:** The `AuditLog` records timestamped partner activity throughout the engagement.
- **Annual internal inspection:** The platform supports inspection sampling by stratified random selection, findings documentation, and remediation tracking. Growth-tier firms get partner-managed inspection; Scale adds inspection analytics.

### AIDecision Ledger Maps to Regulatory Documentation Expectations

Every AI action that affects engagement content creates an `AIDecision` row. The schema is regulator-ready:

| Documentation Expectation | AIDecision Field |
|---|---|
| What technology procedure was performed | `context_type` + `context_id` + `context_table` |
| What model/tool was used | `model_id` (pinned Bedrock model version) |
| What the AI determined | `suggested_value` + `raw_output` |
| Who reviewed the AI output | `reviewed_by_id` |
| What the reviewer decided (accepted/modified/rejected) | `review_action` + `accepted_value` |
| When review occurred | `reviewed_at` |
| Reviewer rationale when overriding | `explanation` |
| Cryptographic proof of artifact | `ProvenanceRecord` (signed envelope — see §10) |

AI outputs alone never constitute sufficient audit evidence. Every Tier 2 AI action (see §6) requires explicit reviewer action and a signed `AIDecision` before becoming part of the engagement record.

### Framework Version Management

`Framework` and `FrameworkVersion` are separate rows. ISO 27001:2013 and ISO 27001:2022 are distinct. SOC 2 TSC 2017, ISO 42001:2023, PCI DSS 3.2.1 and 4.0.1 are distinct. Each engagement records `framework_version_id` at creation time; new framework versions do not change existing engagement references. Evidence → `CommonControl` → `FrameworkRequirement` edges carry effective-dated windows so version churn does not silently invalidate historical coverage. AI Feature 7 (framework version migration) proposes `CommonControlSatisfies` edge updates on version cutover; each proposal is a signed `AIDecision`.

---

## 5. Core Data Model

> **Full specification:** [Domain and Data Model Design](domain-and-data-model-design.md) — complete domain model (bounded contexts, aggregates, invariants), table definitions, column types, constraints, indexes, enum types, and journey-to-entity traceability. What follows is a summary.

### Bounded Contexts

Six bounded contexts and three cross-cutting concerns, derived from the [user journeys](../user-journeys/all-journeys.md) (updated for the compliance/assurance scope):

| # | Context | Key Entities | Module |
|---|---|---|---|
| 1 | Firm Identity | Firm, User, Client, Invitation | `internal/identity` |
| 2 | Regulatory Framework & Common Controls | Framework, FrameworkVersion, FrameworkRequirement, CommonControl, CommonControlSatisfies (STRM), EvidenceItemSupports | `internal/frameworks` |
| 3 | Firm Methodology | MethodologyTemplate, FirmControlObjective, template items | `internal/identity` |
| 4 | Audit Core | Engagement, EngagementFramework, Control, TestProcedure, EvidenceItem, EvidenceLink, DocumentRequest, ClientAcceptance, EngagementQualityReview, Finding, ManagementResponse, CorrectiveActionPlan | `internal/auditcore` |
| 5 | Workpaper Authoring | Workpaper, WorkpaperSignOff, WorkpaperVersion, ReviewNote | `internal/workpaper` |
| 6 | Reporting & Archival | Report, ReportVersion | `internal/reporting` |
| — | Cross-cutting | AIDecision, AuditLog, Notification, ProvenanceRecord | `internal/auditcore` + `internal/provenance` |

**Total domain entities:** ~46 in a single PostgreSQL database (`axiom_db`) with RLS on all tenant-scoped tables. See [Domain and Data Model Design](domain-and-data-model-design.md) for the authoritative list.

### Cross-Framework Evidence Chain

The core architectural differentiator — one evidence upload satisfies all mapped framework requirements simultaneously, through the STRM graph:

```
EvidenceItem ──► EvidenceItemSupports ──► CommonControl
                                             │
                                             ├──► CommonControlSatisfies ──► FrameworkRequirement (SOC 2 CC6.1)
                                             ├──► CommonControlSatisfies ──► FrameworkRequirement (ISO 27001 A.5.15)
                                             ├──► CommonControlSatisfies ──► FrameworkRequirement (ISO 27701 §6.x)
                                             └──► CommonControlSatisfies ──► FrameworkRequirement (HIPAA §164.312(a)(1))

EvidenceItem ──► EvidenceLink ──► TestProcedure ──► Control
                                                      │
                                                      └──► FirmControlObjective (engagement-specific testing narrative)
```

`CommonControl` is the pivot node: it decouples **what a control does** (semantic intent, platform-seeded from SCF / OSCAL / AICPA / CIS) from **how a firm tests it** (`FirmControlObjective`) and **how a framework expresses it** (`FrameworkRequirement`). STRM-encoded `CommonControlSatisfies` edges carry NIST relationship vocabulary, strength score, coverage percentage, and effective-dated windows — so partial satisfaction and version churn are modeled explicitly, not papered over.

This chain requires ACID transactions, which is why Audit Core and Frameworks are in the same `axiom_db` with function-call (not HTTP) interfaces. See [Backend Architecture §3](backend-architecture-design.md).

### Engagement Lifecycle State Machine

```
Planning ──[ClientAcceptance signed by Partner]──► Fieldwork
Fieldwork ──[All Controls: Complete or Exception]──► Review
Review ──[All ReviewNotes resolved + EQR signed off where applicable]──► Reporting
Reporting ──[Report.status = Issued]──► Finalized
Finalized ──[System: assembly_deadline elapsed]──► Archived (IMMUTABLE)
```

Reverse paths exist for exceptional cases (scope change, additional procedures, significant post-reporting issues). Once Finalized, no content can be modified — addenda only.

### Integrated Multi-Framework Engagement Model

`engagement_frameworks` is a 1:N join between `Engagement` and `FrameworkVersion` with `is_primary` flag. Per Journey 11, a single engagement can carry SOC 2 + ISO 27001 + ISO 27701 simultaneously; one `Control` tested once produces coverage edges across every in-scope framework via the `CommonControl` graph, with an auditor-defensible separation of opinion per framework at the report layer.

### Multi-Tenancy

`axiom_db` uses PostgreSQL RLS with `firm_id` on all tenant-scoped tables. Platform-seeded reference data (`frameworks`, `framework_versions`, `framework_requirements`, platform rows of `common_controls` / `common_control_satisfies`) is shared across tenants. Firm-authored extensions carry `firm_id` and are RLS-enforced. Three authorization dimensions: firm isolation (RLS), engagement team membership (point lookup), and client user scoping (engagement-level invitation).

See [Domain and Data Model Design](domain-and-data-model-design.md) for complete attribute definitions, data types, constraints, indexes, enum types, and journey-to-entity traceability.

---

## 6. AI Architecture

> **Full specification:** [AI Architecture Design](ai-architecture-design.md) — LLM provider decision, vector database, all eleven AI features with model assignments and human review gates, AI content tracking, ISO 42001-native positioning, cryptographic provenance, and cost estimates. What follows is a summary.

### LLM Provider and Vector Database

**AWS Bedrock with Claude Haiku and Sonnet**, accessed via VPC endpoint (PrivateLink). IAM-based auth, single AWS BAA, CloudWatch-native metrics. No external API keys, no additional sub-processor. Opus is reserved behind a feature flag for future multi-framework-planning workloads.

**pgvector** (PostgreSQL extension) for embedding storage at launch — framework requirement text, evidence content, `CommonControl` library, firm methodology/policy library, and prior engagement artifacts. Migration path to Qdrant at 5–10M vectors.

### Eleven AI Features

| # | Feature | Model | Tier | Journeys |
|---|---------|-------|------|----------|
| 1 | Document completeness review | Sonnet | 2 | 7, 8 |
| 2 | Evidence → CommonControl mapping suggestion | Haiku | 2 | 4, 7, 8 |
| 3 | Cross-framework gap analysis | Sonnet | 2 | 3, 4, 11, 12 |
| 4 | Workpaper narrative draft | Sonnet | 2 | 5 |
| 5 | Evidence link suggestion | Haiku | 2 | 5, 7 |
| 6 | Risk category suggestion for client acceptance | Sonnet | 2 | 3 |
| 7 | Framework version migration assistance | Sonnet | 2 | 3, 4 |
| 8 | Findings triage & severity reasoning | Sonnet | 2 | 6, 9 |
| 9 | Drift-triggered re-testing | Haiku (+ Sonnet for severity narrative) | 1 detect / 2 conclude | 12 |
| 10 | Agentic management-response drafting | Sonnet | 2 | 9, 12 |
| 11 | Report section draft | Sonnet | 2 | 9 |

Every Tier 2 feature creates an `AIDecision` record and requires explicit human review before affecting engagement content. See [ai-architecture-design.md](ai-architecture-design.md) for full feature definitions.

### Human-in-the-Loop Policy (Three Tiers)

**Tier 1 — Fully Automated:** Text extraction, embedding generation, notification-only flags, overdue reminders, analytics computation, drift *detection* (not conclusion), evidence routing suggestions.

**Tier 2 — Human Approval Required:** All eleven AI features above (Feature 9 detection is Tier 1; any control-status-changing conclusion elevates to Tier 2). No engagement content changes without explicit named human action.

**Tier 3 — Human-Initiated Only:** Engagement status transitions, control conclusions, exception documentation, sign-offs, report issuance, client acceptance, EQR / internal-review sign-off, authoritative cross-framework equivalence declarations, and any action constituting professional judgment under AT-C 105, ISAE 3000 (Revised), ISO 17021-1, PCI DSS QSA sign-off, or HITRUST AEA attestation.

This tiering is the operational spine of Axiom's ISO 42001 alignment — it is enforced in code, not policy.

### ISO 42001-Native Positioning

Axiom is an ISO 42001-native platform. ISO 42001 is the international management-system standard for AI, and Axiom both *assesses* clients against it (Feature 11 drafts the ISO 42001 management-system audit report) and *operates under* it. The three-tier HITL policy, the `AIDecision` ledger with `prompt_hash`/`model_id`/`model_version`/`confidence`/`reviewer`/`override` as first-class columns, the per-feature impact-assessment records, and model-change-management workflows together satisfy ISO 42001 clauses 6–9 for Axiom as a deployer of AI. See [ai-architecture-design.md §6](ai-architecture-design.md#6-iso-42001-native-positioning) for detail.

### Cryptographic Provenance for AI Outputs

Every AI output is signed at emission. At creation, Axiom constructs a provenance envelope (`artifact_id`, `ai_decision_id`, `model_id`, `model_version`, `prompt_hash`, `input_content_hashes`, `output_hash`, `generated_at`, `axiom_version`), signs it with an AWS KMS asymmetric key (`ECC_NIST_P256`, `SIGN_VERIFY`), and writes the signature + envelope + raw output to an S3 Object Lock (WORM) bucket with retention matching the engagement's record-retention policy (minimum 7 years). Human edits to AI drafts sign a new envelope with a pointer to the parent — a tamper-evident chain. A public verification endpoint and CLI let any party recompute the output hash and verify the signature. This is the direct answer to the March 2026 category trust crisis. See [ai-architecture-design.md §7](ai-architecture-design.md#7-cryptographic-provenance-for-ai-outputs) for detail.

### AI Content Tracking

AI-drafted content (Features 4 and 11) is tracked at the **section level** via `ai_content_metadata` on `WorkpaperVersion` and `ReportVersion`. Each AI-generated section records origin timestamp, human-edited flag, editor identity, and modification ratio (Levenshtein distance between AI output and current text, divided by AI character count). The advancement gate requires all AI sections to have at least one human edit; sections with <5% modification trigger a soft confirmation gate. EQR reviewers and managers see per-workpaper and engagement-wide AI edit substantiveness summaries. See [ai-architecture-design.md §5](ai-architecture-design.md#5-ai-content-tracking).

### Per-Engagement AI Cost Estimate

After prompt caching, per-engagement effective cost by framework:

| Framework | Cost Range |
|---|---|
| SOC 2 Type II | $3–5 |
| ISO 27001 / 27701 | $3–5 |
| ISO 42001 | $4–6 |
| HIPAA (with HITRUST CSF r2 path post-MVP) | $3–5 |
| PCI DSS 4.0.1 | $5–8 |
| SOC 1 Type II | $3–5 |

Multi-framework integrated engagements (Journey 11) do not multiply linearly — the `CommonControl` graph means evidence is embedded once, gap analysis runs once per scope, and the marginal cost of adding a second framework is roughly 30–40% of its standalone cost. At 100 engagements/year, platform AI costs are **$300–$800/year per firm** — absorbed into the subscription. KMS signing (~$0.03/engagement) and S3 Object Lock WORM storage (~$1–2/engagement) are rounding error. See [ai-architecture-design.md §8](ai-architecture-design.md#8-per-engagement-ai-cost-estimate).

### What AI Does NOT Do at Launch

Fine-tuned models per firm, on-device execution, multi-agent orchestration frameworks (LangGraph / CrewAI — staged pipelines in Feature 10 use River-scheduled Tier-gated steps), AI-generated certification decisions / attestation opinions / QSA sign-off / HITRUST AEA attestations, AI-generated authoritative cross-framework equivalence declarations (AI suggests mappings; authoritative crosswalk is SCF/OSCAL/AICPA-sourced plus firm-confirmed), external side effects without named human approval on the AIDecision, auto-finalization of any audit content, use of Opus (reserved, feature-flagged).

---

## 7. Technology Stack

### Frontend

**TypeScript / React SPA.** Component library: Shadcn/ui (accessible, customizable, no licensing overhead). State management: TanStack Query for server state; Zustand for local UI state. API types are generated from the OpenAPI specs in `packages/openapi/` via `openapi-typescript`.

### Backend: Go Modular Monolith + Python PDF Service + Provenance Signer

**Go as the primary backend language, structured as a modular monolith. Python retained for PDF extraction. A separate `provenance-signer` ECS service isolates KMS signing.**

Go was chosen for compile-time type safety, lean container images (30–50MB per Fargate task), and strong fit for compliance SaaS. The backend is a single Go binary organized into internal packages by bounded context. Modules communicate via Go interfaces, not HTTP — this provides ACID transactions across the evidence chain and eliminates distributed-systems overhead.

Python is retained as a stateless service for PDF extraction (`pdfplumber` + Tesseract). The `provenance-signer` is a minimal Go service with narrow IAM (only `kms:Sign` on the provenance-signing key) to reduce signing blast radius.

**Full module descriptions, database design, and inter-module communication patterns are specified in [`backend-architecture-design.md`](./backend-architecture-design.md).**

**Monorepo structure:**
```
apps/
  axiom-api/        — Go: Modular monolith (single binary)
    internal/
      gateway/      — Chi middleware: JWT verification, routing, rate limiting
      identity/     — Auth, RBAC, firm/user/client, templates
      auditcore/    — Engagements, controls, evidence, AI decisions, findings
      frameworks/   — Framework catalog, CommonControl graph, STRM edges, cross-mapping, drift detection
      provenance/   — AI output signing, evidence provenance, WORM manifest management
      workpaper/    — Workpapers, Yjs collaboration (WebSocket)
      reporting/    — Report generation, S3 archival
      ai/           — Bedrock client, prompt templates, embedding helpers
      platform/     — DB, config, OTel, River, common middleware
  doc-processing/   — Python: PDF extraction only
  provenance-signer/— Go: Isolated ECS service with narrow KMS:Sign IAM
packages/
  openapi/          — OpenAPI 3.1 specs organized by module
```

Turborepo manages the monorepo with build caching.

### API Layer: REST + OpenAPI

REST with OpenAPI 3.1 for all services. Each module defines its contract as an OpenAPI 3.1 spec in `packages/openapi/`. `oapi-codegen` generates typed Go server interfaces; `openapi-typescript` generates typed fetch clients. Authorization is composable Go middleware (`WithFirmIsolation`, `WithEngagementAccess`, `WithClientScoping`).

### Database: PostgreSQL with RLS

PostgreSQL is the sole persistent data store. One RDS instance hosts a single database (`axiom_db`) with row-level security on all tenant-scoped tables. Each module owns specific tables and accesses them via its own `sqlc` queries; cross-module data goes through Go service interfaces. `golang-migrate` for schema migrations. PgBouncer (transaction mode) as ECS sidecar.

**pgvector** extension enabled for embeddings (framework requirements, evidence, `CommonControl`, firm methodology/policy library).

### Workflow Engine: Two-Tier (River + Step Functions)

**River** (Go-native Postgres-backed queue) for fire-and-forget background jobs — document extraction, embedding generation, AI mapping/gap analysis, drift detection, notifications. **AWS Step Functions** Standard Workflows for the engagement lifecycle state machine and long-running reminder sequences.

### Real-Time Collaboration

Workpaper rich text uses TipTap + Yjs CRDT. Structured audit data (evidence mapping acceptance, control status, engagement transitions) uses server-authoritative operations with field-level locking for high-stakes decisions. CRDTs are explicitly not used for structured audit data — the automatic-merge behavior is incompatible with the regulatory requirement that conflicting professional judgments produce an unambiguous, auditable outcome.

### Spreadsheet Component: AG Grid Community + HyperFormula

AG Grid Community (MIT) + `hyperformula` (MIT) at launch — used for coverage dashboards, evidence ledgers, and tabular evidence capture. Univer evaluated at 6–12 months.

### Infrastructure: AWS + ECS Fargate

ECS Fargate, ALB, RDS PostgreSQL (Multi-AZ), S3 (with Object Lock for evidence, archive, reports, and signed AI artifacts), CloudFront, SES, Secrets Manager, CloudWatch + X-Ray, WAF, GuardDuty, AWS Config. Three ECS services: `axiom-api`, `doc-processing`, and `provenance-signer`. Dedicated KMS asymmetric key (`axiom-{env}-provenance-signing`, `ECC_NIST_P256`, `SIGN_VERIFY`) for AI and evidence provenance. Dedicated S3 bucket (`axiom-{env}-scf-catalog`) for SCF/OSCAL/AICPA/CIS crosswalk import.

**Full AWS account structure, VPC design, Terraform workspaces, CI/CD, security controls, observability, and cost estimates are specified in [`infrastructure-design.md`](./infrastructure-design.md).**

---

## 8. Integration Roadmap

> Research: [07-integrations.md](../research/07-integrations.md) <!-- STALE: research doc predates compliance pivot -->

### Integration Philosophy

Axiom prefers **standards-based ingestion (OAuth 2.0, SAML, SCIM, webhooks)** over per-vendor custom connectors. Cloud-provider, identity-provider, and dev-tool integrations exist to pull evidence (SOC 2 CC6.x, ISO 27001 A.5.x, A.8.x, PCI DSS 2.x/7.x/8.x, HIPAA §164.312) directly into `EvidenceItem` records, where they enter the same AI mapping / human-review pipeline as direct uploads.

### Launch (Must-Have)

| Integration | Type | What It Enables |
|---|---|---|
| **Transactional email (SES)** | Outbound | Document request notifications, client portal invitations, review alerts, drift alerts |
| **Direct file upload** | Evidence ingestion | Clients upload evidence via tokenized link, no login required. Handles CSV, Excel, PDF, ZIP |
| **Microsoft / Google SSO (SAML/OAuth)** | Identity | Firm staff authentication |
| **SCF quarterly import** | Platform seed data | Updates the `CommonControl` and `CommonControlSatisfies` catalog with the quarterly SCF release |
| **OSCAL catalog import** | Platform seed data | NIST-family catalogs (future-proofs for FedRAMP) |

### First 6 Months Post-Launch

| Integration | Type | Purpose |
|---|---|---|
| **AWS, Azure, GCP** | Cloud evidence | Evidence collection for SOC 2 CC6, ISO 27001 A.5/A.8, PCI DSS 2.x, HIPAA technical safeguards |
| **Okta, Google Workspace, Microsoft Entra** | Identity | Access reviews, SSO, MFA evidence (SOC 2 CC6.2/6.3, ISO A.5.15-18, HIPAA §164.308(a)(4)) |
| **GitHub, GitLab, Bitbucket, Jira, Linear** | Dev tools | Change management, ticket evidence, SDLC audit trail (SOC 2 CC8.1, ISO A.8.25-32) |
| **Slack, Zendesk** | Productivity / support | Incident response evidence (SOC 2 CC7.x, ISO A.5.24-29) |
| **Drata Audit Hub + Vanta evidence export** | Evidence bridge | Ingest structured, control-tagged evidence from auditee-side tools when the client is already on Drata / Vanta — abstracts the client's infrastructure stack |

### 6–18 Months Post-Launch

| Integration | Type | Notes |
|---|---|---|
| **Rippling, Gusto, BambooHR** | HRIS | Onboarding/offboarding/training evidence (SOC 2 CC1.4, ISO A.6.x, HIPAA §164.308(a)(5)) |
| **Sprinto / Hyperproof evidence export** | Evidence bridge | Analogous to Drata/Vanta; smaller installed base |
| **Box / SharePoint / Google Drive / Dropbox** | Cloud storage | Client-side evidence repositories (Box retained — useful for regulated-industry clients) |
| **Public API + webhooks (Axiom-outbound)** | Platform extensibility | Required before firms can build custom workflows on top of Axiom |
| **UCF license passthrough** | Optional Enterprise SKU | Commercial UCF crosswalk content in addition to SCF/OSCAL/AICPA/CIS defaults |

### Integration Architecture Principles (all tiers)

1. Every integration is an abstraction behind an internal interface — evidence enters the same `EvidenceItem` data model whether it came from a CSV upload, an OAuth connector, or a Drata Audit Hub pull.
2. OAuth credentials are stored per firm + client, not globally. Each firm's connections are isolated.
3. Evidence from integrations goes through the same AI mapping + human-review workflow as direct uploads. `EvidenceItem.source_type` and `source_integration` record provenance but do not bypass review.
4. Connector-captured browser evidence (screenshots, DOM snapshots) is signed at capture via the Provenance module (`SignedScreenshot`, `HashedDOMSnapshot`).
5. Graceful degradation: if a connector is down, auditor or auditee falls back to file upload without blocking the engagement.
6. Read-only scopes at launch. No write-back into client systems.

---

## 9. Legal and Data Governance

> Research: [08-legal-data-governance.md](../research/08-legal-data-governance.md)

### Data Processing Agreement Structure

Axiom is a **data processor**; the customer firm is the **data controller**. Every customer (US, Canada, EU, UK, AU) must execute a signed DPA before accessing the platform.

**Required DPA contents (GDPR Article 28 and equivalent):**
1. Processing only on documented instructions from the controller
2. Confidentiality obligations on all authorized persons
3. Article 32-compliant technical and organizational measures (AES-256 at rest, TLS 1.3 in transit, access controls, incident response)
4. Sub-processor management: prior written authorization, equivalent obligations, 30-day advance notice
5. Data subject rights assistance (access, correction, deletion where not legally blocked)
6. 60-day read-only export window post-termination; deletion within 30 days thereafter; written deletion certificate
7. Audit rights: Axiom's own SOC 2 Type II + ISO 27001 + ISO 42001 certificates satisfy most audit provisions
8. Breach notification: Axiom notifies within 24 hours of awareness

**Explicit DPA prohibition:** Axiom does not use customer data (control evidence, policies, risk assessments, AI decision artifacts, client lists, uploaded documents) to train AI models without separate controller authorization. All AI model inference runs via AWS Bedrock over a VPC endpoint — data does not leave the AWS network and is not retained by Anthropic. The DPA states this explicitly.

**Sub-processor list:** Published publicly; current sub-processors include AWS (infrastructure, AI model inference via Bedrock, workflow execution via Step Functions, transactional email via SES, KMS signing — all covered under a single AWS sub-processor entry).

### GDPR/CCPA Deletion Rights vs. Immutable Archiving

The conflict between data subject deletion rights and record-retention obligations is resolved by **GDPR Article 17(3)(b)** and **CCPA §1798.105(d)(8)** — the legal-obligation exemptions. Retention of ISMS/PIMS/AIMS records, SOC engagement files, and PCI ROC evidence is required by the applicable regulatory regime; a data subject cannot compel deletion of records the firm is obligated to retain.

**Platform implementation of the deletion-request workflow:**
1. Acknowledge receipt within 72 hours
2. Invoke legal retention exemption, citing the specific regime (AU-C 230 / AT-C 105 / ISO 27001:2022 §7.5 / ISO 42001:2023 §7.5 / HIPAA §164.530(j) / PCI DSS 12.10.1 / GDPR Art. 17(3)(b))
3. Restrict retained data — no use outside the legal retention obligation
4. Log the decision in the `AuditLog`
5. Delete on schedule when the mandatory retention window closes

The platform provides templated deletion request response letters the firm can use with data subjects. Data minimization: the platform only captures personal data strictly necessary for attestation documentation.

### Data Residency Approach

**US at launch (AWS us-east-1).** Canada customers covered under PIPEDA (no mandatory residency; DPA addresses cross-border obligations; Quebec Law 25 PIA documented for Quebec-based clients). **EU/APAC as Enterprise tier (year 2+)** — AWS eu-central-1 (Frankfurt), eu-west-2 (London), ap-southeast-2 (Sydney) data residency options under SCCs (2021 Module 2) incorporated into the DPA.

### AI Liability Policy

**Human review is mandatory; AI never auto-finalizes.** Every AI-generated workpaper narrative, finding suggestion, control conclusion, evidence assessment, or mapping requires explicit human review and named sign-off before it is treated as finalized. This creates an auditable paper trail (the `AIDecision` ledger + `ProvenanceRecord` signature chain) demonstrating that the reviewer exercised independent professional judgment.

AI liability protections embedded in the MSA:
- Disclaimer that Axiom does not provide accounting, auditing, certification, QSA, HITRUST, or other professional services
- Liability cap equal to 12 months of fees paid (direct damages only)
- Consequential damages waiver
- Data security super-cap: 2× annual fees for data breaches
- Gross negligence and fraud carve-out from the limitation of liability

### Insurance Requirements

| Coverage | Minimum at Launch | Target at Series A |
|---|---|---|
| Tech E&O (Professional Liability) | $1M per occurrence / $1M aggregate | $2M / $2M |
| Cyber Liability | $1M per occurrence | $2M–$5M |
| General Commercial Liability | $1M / $2M | $1M / $2M |

### Offboarding Obligations

60-day read-only export window; written deletion certificate within 30 days of export-window close. Export formats include PDF (rendered workpapers/reports), CSV/XLSX (structured data), JSON (engagement metadata and API exports), native format for uploaded documents, and **signed provenance manifests** so the firm retains verifiable proof that exported AI artifacts are byte-identical to what Axiom generated. All customer data is deleted from production systems, backups, and sub-processor systems within 30 days of export-window close; encrypted isolated backups are retained up to 90 additional days then purged on next rotation.

---

## 10. Security Architecture

### Encryption

**In transit:** TLS 1.2 minimum (1.3 preferred and negotiated) on all connections. HTTP redirects to HTTPS; ALB uses `ELBSecurityPolicy-TLS13-1-2-2021-06`.

**At rest:** AES-256 server-side encryption:
- RDS PostgreSQL via customer-managed KMS key (`axiom-{env}-rds`), annual automatic rotation
- S3 evidence files: SSE-S3 default; HIPAA-flagged evidence uses SSE-KMS with a dedicated key (`axiom-{env}-hipaa`) for CloudTrail decrypt auditing
- S3 archive bucket (finalized engagements): SSE-KMS with the HIPAA key
- Backups encrypted with the same KMS key as source

### Multi-Tenancy Isolation

PostgreSQL RLS with `firm_id` indexed on every tenant-scoped table. `SET app.current_firm_id` set at session start. Application-layer middleware enforces engagement team membership and client-user scoping. Penetration testing on multi-tenant isolation is part of the pre-Series A security review package.

### Audit Log

`AuditLog` is PostgreSQL insert-only (`RULE` preventing UPDATE / DELETE). All significant events: engagement status changes, workpaper saves, sign-offs, access grants, AIDecision outcomes, deletion-request responses, archival events. Sequential bigint IDs for unambiguous temporal ordering.

### Cryptographic Evidence and AI Provenance (Category Differentiator)

Axiom signs every AI output and every connector-captured evidence artifact at creation using **AWS KMS `ECC_NIST_P256`** (`SIGN_VERIFY`, non-exportable, verifiable via public key). The signature + envelope + raw output/artifact are written to an S3 bucket with **Object Lock in COMPLIANCE mode** — no user, including AWS root, can delete or overwrite until retention elapses (minimum 7 years; per-object retention set to the engagement's `retention_deadline`). This implements evidence and AI-output immutability at the storage layer, not just the application layer.

Every finalized AI output is signed and WORM-stored. On retrieval, the `ProvenanceRecord` is re-verified (hash match + signature valid + lock status intact). Reports cannot be rendered if any cited artifact fails chain verification. A public verification endpoint and CLI let any party — EQR reviewer, ISO CB technical reviewer, client procurement, regulator — verify the signature independently. The provenance-signing key material is non-exportable; `kms:Sign` is restricted to the isolated `provenance-signer` ECS task role.

This is the data-path answer to the March 2026 agentic-compliance trust crisis. See [`ai-architecture-design.md §7`](ai-architecture-design.md#7-cryptographic-provenance-for-ai-outputs) and [`infrastructure-design.md §4`](infrastructure-design.md).

### Axiom's Own Certifications (Day-One Posture)

Axiom targets from day one: **SOC 2 Type 2**, **ISO 27001**, **ISO 42001**. The AIDecision ledger, three-tier HITL policy, and cryptographic provenance are Axiom's own ISO 42001 controls, not just features it sells. Customers can inspect the dogfooding artifacts as part of procurement.

### HIPAA BAA

Axiom executes a HIPAA Business Associate Agreement with AWS. The AWS BAA covers RDS, S3, Bedrock, Step Functions, SES, and KMS. HIPAA-flagged engagement files use SSE-KMS with more restrictive IAM. Axiom makes a HIPAA BAA available to customer firms (self-serve for Growth/Scale; dedicated negotiation for Enterprise).

---

## 11. User Journeys

Full journey maps with stage-by-stage detail (user actions, touchpoints, emotional states, competitor context, pain points, and opportunities) are in [`docs/user-journeys/all-journeys.md`](../user-journeys/all-journeys.md). The table below summarizes each journey.

| # | Persona | Goal | Key System Gates / Entities | AI Touchpoints |
|---|---------|------|-----------------------------|----------------|
| 1 | FirmAdmin | Set up firm and launch first engagement | `Firm`, `MethodologyTemplate`, `Engagement` scaffold, 5-step onboarding, 14-day trial clock | — |
| 2 | Staff Auditor | Join platform and reach first task | Magic link auth, role assignment, guided tour, `EngagementTeamMember` | — |
| 3 | Partner | Create and scope a new engagement | `ClientAcceptance` gate, EQR independence validation (where applicable), framework version lock after Fieldwork | Risk category (Feature 6), Framework migration (Feature 7), Gap analysis (Feature 3) |
| 4 | ClientAdmin | **Cross-framework evidence mapping** — upload one artifact, see which `FrameworkRequirement` nodes it satisfies across every in-scope framework | `EvidenceItem` → `EvidenceItemSupports` → `CommonControl` → `CommonControlSatisfies` → `FrameworkRequirement`, period-aware staleness, partial-satisfaction percentages with gap lists | Evidence → CommonControl mapping (Feature 2), Gap analysis (Feature 3), Framework migration (Feature 7) |
| 5 | Staff Auditor | Test controls and prepare workpapers | `TestProcedure` status progression, AI content tracking gate, four-level sign-off (Tester → Detailed Reviewer → General Reviewer → Final Reviewer) via `WorkpaperSignOff` | Evidence link (Feature 5), Workpaper draft (Feature 4) |
| 6 | Manager | Review workpapers and advance the engagement | `ReviewNote` (immutable), `WorkpaperSignOff` reviewer-level enforcement, phase transition guards | Findings triage (Feature 8) |
| 7 | Staff Auditor | Manage document requests and collect evidence | `DocumentRequest` lifecycle, AI review queue sorted by confidence, reminder state machine | Document completeness (Feature 1), Evidence link (Feature 5), Evidence → CommonControl mapping (Feature 2) |
| 8 | Client Contact | Fulfill audit document requests | Tokenized Client Hub link (no login, engagement-scoped), `ClientAdmin` delegation | Document completeness (Feature 1), Evidence → CommonControl mapping (Feature 2) |
| 9 | Partner | Generate report, finalize, and archive | Report issuance triggers retention computation; Finalized locks content; S3 Object Lock WORM archival; addendum workflow | Report section draft (Feature 11), Findings triage (Feature 8), Management-response drafting (Feature 10) |
| 10 | EQR Reviewer | Conduct engagement quality review | Read-only (not `EngagementTeamMember`), `EngagementQualityReview` sign-off gate, immutable EQR record, AI edit substantiveness summaries | — |
| 11 | Partner | **Multi-framework integrated engagement** — scope SOC 2 + ISO 27001 + ISO 27701 in one engagement; one control tested once, evidence satisfies all in-scope frameworks via `CommonControl` graph; separate opinion per framework at reporting | `EngagementFramework` (multiple per `Engagement`), shared `CommonControl` library, reconciled sampling windows | Gap analysis (Feature 3), Evidence → CommonControl mapping (Feature 2) |
| 12 | ClientAdmin | **Continuous assurance** — respond to drift alerts, auto-retest affected controls, update risk register, notify auditor when drift is material | Drift detection on connector-sourced configs, period-window thresholds, `Notification (DriftDetected / EvidenceExpiring)`, auto-enqueued re-test jobs | Drift-triggered re-testing (Feature 9), Findings triage (Feature 8), Management-response drafting (Feature 10) |

### Regulatory / methodology constraints by journey

| Constraint | Standard | Journeys |
|-----------|----------|----------|
| Client acceptance before fieldwork | AICPA SQMS 1 / internal equivalent for ISO | 3 |
| EQR / internal-reviewer independence | SQMS 2 (SOC) / ISO 17021-1 §9.6 analog | 3, 10 |
| Framework version locked after fieldwork begins | §4 | 3, 4 |
| Four-level sign-off hierarchy enforced at data layer (Tester → Detailed Reviewer → General Reviewer → Final Reviewer) | SQMS 1 / equivalent firm quality frameworks | 5, 6 |
| AI-drafted sections must be substantively edited before sign-off | ISO 42001 §7–9 + firm quality policies | 5, 9 |
| Review notes cannot be deleted | Evidence preservation standards | 6 |
| Period coverage check for SOC 2 Type II evidence | AT-C 320 | 7, 12 |
| Framework-specific evidence staleness (ASV 90d, pen test 1y, etc.) | PCI DSS v4.0.1, vendor frameworks | 4, 7, 12 |
| All AI decisions logged as `AIDecision` records + signed `ProvenanceRecord` | ISO 42001 + firm AI governance | 3, 4, 5, 6, 7, 9, 12 |
| Assembly deadline enforcement (SOC attestations) | AU-C 230 | 9 |
| WORM archival | AICPA retention + ISO record-keeping + PCI 12.10 | 9 |
| Retention periods per engagement type | Per §4 | 9 |
| Client upload tokens expire and require re-generation | Security policy | 7, 8 |
| Integrated multi-framework engagement produces separate opinion per framework | AT-C 105 / ISO 17021-1 / PCI QSA | 11 |
| Drift-triggered re-test results feed AIDecision ledger before any control status change | ISO 42001 + evidence-integrity | 12 |

---

## 12. Flows Without Competitor Equivalent

These flows represent genuine Axiom differentiation — no competitor (Fieldguide, Agentive, Yak, AuditBoard CrossComply, Hyperproof, Drata, Vanta, Secureframe, Sprinto, Thoropass) currently offers them.

1. **Cross-framework evidence mapping (STRM-grade, period-aware)** (Journeys 4, 11) — top-level product differentiator. One `EvidenceItem` flows coverage through `EvidenceItemSupports → CommonControl → CommonControlSatisfies → FrameworkRequirement` with NIST STRM relationship vocabulary, strength scores, coverage percentages, and effective-dated windows that honor framework-specific staleness (ASV 90d, pen test 1y, background check 1y, SOC 2 Type II periods, ISO surveillance cycles). Never a green checkmark on partial coverage — percentages and gap lists only.
2. **Cryptographic AIDecision provenance** (all AI-producing journeys) — every AI output signed with AWS KMS at emission and WORM-stored with a public verification endpoint. Post–March-2026 category differentiator; a direct answer to the trust crisis.
3. **Both-sided auditor + auditee workspace** (Journeys 4, 7, 8, 11, 12) — the auditor-side engagement workflow and the auditee-side Client Hub / continuous-monitoring workspace operate on the same data model. Competitors are one-sided.
4. **ISO 42001-native HITL ledger** (all AI journeys) — three-tier HITL policy, impact-assessment records per feature, model-change-management workflow, prompt/model/context/confidence/reviewer/override as first-class columns on `AIDecision`. Axiom dogfoods its own framework.
5. **Agentic management-response drafting with round-trip remediation** (Journeys 9, 12) — drafts management response to a finding, assigns owner, opens a ticket in Jira / Linear / GitHub, tracks closure evidence back through the `Finding → ManagementResponse → CorrectiveActionPlan` chain, and re-tests when remediation is declared complete. Each staged action is a signed `AIDecision` with explicit human approval.
6. **Drift-triggered continuous re-testing** (Journey 12) — connector-sourced configuration drift detection enqueues autonomous re-test jobs; proposed control-status changes require Tier 2 human approval. Converts a point-in-time audit into a continuous-attestation posture without giving up auditor-defensible human-in-the-loop.
7. **Section-level AI content tracking and edit gate** (Journeys 5, 6, 9, 10) — AI-drafted content tracked per section via `ai_content_metadata` with modification ratios; unedited AI sections block sign-off; soft-warning gate at <5% modification; EQR reviewers see engagement-wide AI edit substantiveness summaries.
8. **EQR / internal-reviewer independence enforcement** (Journey 10) — system-level validation that the quality reviewer is not on the engagement team.
9. **Post-finalization addendum workflow** (Journey 9) — proper implementation with immutable original content and versioned addenda signed into the provenance chain.
10. **Automatic retention computation and WORM archival** (Journey 9) — retention deadline computed per engagement type at report issuance, enforced via S3 Object Lock COMPLIANCE mode with signed manifests.

---

## 13. Out of Scope at Launch

The following capabilities are explicitly excluded from the MVP. Each exclusion is deliberate.

| Capability | Rationale for Exclusion |
|---|---|
| **Financial audit / GAAS** | Exited vertical. No trial balance, sampling, materiality, or financial-statement workpapers. |
| **PCAOB public-company audits** | Out of scope. Different audit regime, different liability profile, different buyer. |
| **SOX internal controls / internal audit / enterprise GRC** | AuditBoard's market. Different buyer (internal audit director vs. external attestation partner), different workflow (continuous monitoring vs. point-in-time engagement with sign-off hierarchy). |
| **Replacing CB accreditation authority** | CBs are valid customer firms; Axiom supports their engagement-delivery workflow and generates draft certificate documents from templates, but does not become an accreditation authority and does not replace the CB's certification decision under ISO 17021-1. |
| **Replacing QSA accreditation authority** | QSA firms are valid customer firms; Axiom supports their gap assessment, evidence collection, and ROC/AOC assembly (including template-generated ROC and AOC drafts), but does not become a QSA and does not replace QSA legal sign-off under PCI SSC accreditation. |
| **HITRUST CSF r2 (full assessment mode)** | De facto market standard for HIPAA-heavy environments but requires HITRUST Authorized External Assessor partnership. Post-MVP per compliance-pivot-findings decision 8. |
| **ESG / sustainability reporting** | Adjacent category, different buyer, content investment not justified by immediate revenue. |
| **White-labeling / reseller channels** | Adds multi-tenancy complexity at the branding layer; dilutes brand before established. |
| **Custom AI model training per firm** | Insufficient training data at launch. RAG over firm methodology achieves the same functional goal. Revisit at year 3. |
| **On-device / local model execution** | Not required for the target ICP. Adds significant infrastructure complexity. |
| **Multi-agent AI orchestration (LangGraph, CrewAI)** | Staged pipelines in Feature 10 are implemented as River-scheduled Tier-gated steps. Orchestration frameworks add debugging overhead with no benefit for the 11 defined features. |
| **AI auto-finalization of any audit content** | Outside the three-tier HITL policy. Professional liability exposure. ISO 42001 Tier 3 actions are human-initiated only. |
| **AI-issued certification / attestation / QSA / HITRUST AEA sign-offs** | Professional licensure responsibility; not delegable to software. |
| **Mobile application** | Desktop-intensive workflow (document review, coverage dashboards, sign-off). Web-responsive design is the priority; native mobile is year 2+. |
| **Public REST / GraphQL API at launch** | The internal REST API is not versioned or documented for external consumption. Public API + webhooks are in the 6–18 month roadmap. |
| **Multi-office / global firm management** | A single firm with multiple offices can operate in Axiom as one tenant. Cross-office assignment and jurisdiction-specific methodology are deferred. |
| **EU/APAC data residency** | US and Canada at launch (AWS us-east-1). EU (Frankfurt), UK (London), and APAC (Sydney) residency are Enterprise-tier options in year 2. |

---

*End of Axiom Product Specification (compliance/assurance pivot)*
