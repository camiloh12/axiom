# Axiom AI Architecture Design

**Date:** April 17, 2026
**Status:** Implementation-Ready (compliance/assurance pivot)
**Supersedes:** Section 6 of axiom-spec-design.md

> Research: [05-ai-architecture.md](../research/05-ai-architecture.md), [compliance-pivot-findings.md](compliance-pivot-findings.md)

This document is the authoritative specification for all AI decisions in Axiom. The [product spec](axiom-spec-design.md) Section 6 contains a summary; this document contains the full detail.

Axiom is a compliance/assurance platform covering SOC 2, ISO 27001, ISO 27701, ISO 42001, HIPAA (with HITRUST CSF r2 path), PCI DSS, and SOC 1. The AI surface is designed around framework-aware reasoning, cross-framework evidence reuse, and an ISO 42001-native human-in-the-loop ledger. No feature in this document addresses financial audit, trial balance, GAAS, or PCAOB-specific workflows.

---

## Table of Contents

1. [LLM Provider Decision](#1-llm-provider-decision)
2. [Vector Database Decision](#2-vector-database-decision)
3. [Human-in-the-Loop Policy](#3-human-in-the-loop-policy)
4. [AI Features](#4-ai-features)
5. [AI Content Tracking](#5-ai-content-tracking)
6. [ISO 42001-Native Positioning (Axiom Dogfoods Its Own Framework)](#6-iso-42001-native-positioning)
7. [Cryptographic Provenance for AI Outputs](#7-cryptographic-provenance-for-ai-outputs)
8. [Per-Engagement AI Cost Estimate](#8-per-engagement-ai-cost-estimate)
9. [What AI Does NOT Do at Launch](#9-what-ai-does-not-do-at-launch)

---

## 1. LLM Provider Decision

**Primary provider: AWS Bedrock with Claude models (Haiku and Sonnet), via the Anthropic/AWS partnership.**

Rationale:
- **PrivateLink** — model API calls are made over a VPC endpoint and never leave the AWS network. No outbound HTTPS to an external provider; stronger posture than ZDR mode at the API level.
- **IAM-based auth** — no separate vendor API keys to store in Secrets Manager, rotate, or audit. Access control is standard IAM policy.
- **Single AWS BAA** — Bedrock is covered under the existing AWS HIPAA BAA. Anthropic is not an additional sub-processor in the DPA or the SOC 2 audit scope.
- **Same models, same capability** — Bedrock provides Claude Haiku and Sonnet with the same 200K+ context window, prompt caching (0.1x the base input rate for cached reads), and batch inference (50% discount) as Anthropic direct. The cost model is effectively identical at this scale.
- **CloudWatch metrics natively** — token counts, latency, and errors report to CloudWatch without additional instrumentation. No separate dashboard or custom metrics pipeline needed.
- Single model family eliminates the prompt engineering and behavior consistency overhead of a multi-provider setup.

Model availability on Bedrock typically lags Anthropic's direct API by a few weeks on new releases. For a compliance platform where model stability is preferred over early access to new capabilities, this is a feature rather than a drawback — a validation period before switching models would be standard practice regardless.

---

## 2. Vector Database Decision

**pgvector (PostgreSQL extension) at launch, with migration path to Qdrant.**

pgvector adds zero operational overhead (uses existing PostgreSQL instance), keeps vector data in the same security boundary as the rest of the application, and performs adequately (sub-20ms at 1M vectors with HNSW index) for the target scale. At maturity (50 firms x 100 engagements x 200 evidence items), the total vector count is approximately 1M — well within pgvector's comfortable range. Migrate to self-hosted Qdrant when approaching 5–10M vectors.

**pgvector scope** (the set of content embedded into the shared vector store):

- **Framework requirement text** — every `FrameworkRequirement` across SOC 2 TSCs, ISO 27001/27701 Annex A, ISO 42001 controls, HIPAA Safeguards, PCI DSS 4.0.1 requirements, and SOC 1 control objectives. Powers semantic search over "what does this requirement cover?" and is the retrieval layer for cross-framework gap analysis.
- **Evidence content** — extracted text summaries and metadata from every `EvidenceItem` (client uploads, integration pulls, screenshots, configuration exports). Powers evidence→CommonControl mapping suggestions (Feature 2), evidence link suggestions (Feature 4), and cross-engagement reuse detection.
- **CommonControl library** — SCF-derived control descriptions, AICPA official mapping rationales, and firm-authored control objective narratives. Powers the semantic layer on top of the STRM-encoded cross-mapping graph.
- **Firm methodology and policy library** — firm-specific audit methodology, templated procedures, auditee-side policy documents (for Client Hub policy lifecycle features).
- **Prior engagement artifacts** — anonymized workpaper narratives, findings text, and report sections used as few-shot retrieval for drafting tasks.

At maturity (50 firms × 100 engagements × 200 evidence items, plus framework text and common-control library) the total vector count is approximately 1–2M — well within pgvector's HNSW-indexed comfortable range. Migrate to self-hosted Qdrant when approaching 5–10M vectors.

---

## 3. Human-in-the-Loop Policy

### Tier 1 — Fully Automated (no human approval required)

AI may act without user interaction for: text extraction from uploaded documents, embedding generation, flagging documents as potentially incomplete (notification only, not decision), generating overdue reminder emails to clients, surfacing "this evidence was used in a prior engagement" notifications, routing/classifying inbound evidence to candidate control areas (suggestion only), and drift detection on monitored configurations (detection only — re-testing that would alter control status elevates to Tier 2).

### Tier 2 — Human Approval Required Before Taking Effect

AI suggests; human confirms. No audit file content changes without explicit auditor action. Every Tier 2 action creates an `AIDecision` record. Features in this tier:

- Document completeness review (Feature 1)
- Evidence → CommonControl mapping suggestion (Feature 2)
- Cross-framework gap analysis recommendations (Feature 3)
- Workpaper narrative drafts (Feature 4)
- Evidence link suggestions (Feature 5)
- Risk category suggestions for client acceptance (Feature 6)
- Framework version migration remapping (Feature 7)
- Findings triage & severity reasoning (Feature 8)
- Drift-triggered re-testing results (Feature 9) — re-test runs are Tier 1 detection, but any proposed risk-register update or control status change is Tier 2
- Agentic management-response drafting (Feature 10) — the draft, the owner assignment, and the remediation-ticket content each require explicit approval before external side-effects (ticket creation, email) fire
- Report section drafts (Feature 11)

### Tier 3 — Human-Initiated Only (AI may assist but never initiates)

These actions can only be triggered by a named human: engagement status transitions, control conclusions (pass/fail/partial), exception documentation, workpaper sign-off **at every level of the four-level hierarchy (Tester / DetailedReviewer / GeneralReviewer / FinalReviewer)** — each level requires a named human to sign off, and AI cannot record a `WorkpaperSignOff` row at any level — EQR or independent reviewer sign-off, report issuance, client acceptance documentation, authoritative cross-framework equivalence declarations, and any action constituting professional judgment under AT-C 105, ISAE 3000 (Revised), ISO 17021-1, PCI DSS QSA sign-off, or HITRUST Authorized External Assessor attestation. The four-level reviewer hierarchy is a workflow attribute orthogonal to the three-tier HITL policy: the HITL tier governs *who can perform an action*, while the reviewer level governs *which workpaper sign-off step* a given user is satisfying.

This three-tier policy is the operational spine of Axiom's ISO 42001 alignment (see [Section 6](#6-iso-42001-native-positioning)): AI outputs inform professional judgments; they do not constitute them, and every decision is traceable to a named human reviewer.

---

## 4. AI Features

Eleven AI features at launch, organized by the user journey they support. Every feature follows the same structure: purpose, trigger, inputs, model, process, output, human review gate, and failure modes.

**Feature index:**

| # | Feature | Model | Tier |
|---|---|---|---|
| 1 | Document completeness review | Sonnet | 2 |
| 2 | Evidence → CommonControl mapping suggestion | Haiku | 2 |
| 3 | Cross-framework gap analysis | Sonnet | 2 |
| 4 | Workpaper narrative draft | Sonnet | 2 |
| 5 | Evidence link suggestion | Haiku | 2 |
| 6 | Risk category suggestion for client acceptance | Sonnet | 2 |
| 7 | Framework version migration assistance | Sonnet | 2 |
| 8 | Findings triage & severity reasoning | Sonnet | 2 |
| 9 | Drift-triggered re-testing | Haiku (+ Sonnet for severity narrative) | 1 detect / 2 conclude |
| 10 | Agentic management-response drafting | Sonnet | 2 |
| 11 | Report section draft | Sonnet | 2 |

Opus is **reserved** (not used at launch) for future complex multi-framework planning tasks where Sonnet's reasoning depth proves insufficient — specifically, whole-program scoping across 4+ concurrent frameworks with overlapping surveillance windows. Any Opus routing is gated behind a feature flag and a cost guardrail.

### Feature 1: Document Completeness Review

**Purpose:** Eliminate the primary engagement bottleneck — the back-and-forth cycle when clients submit incomplete or wrong documents.

**Journeys:** 7 (auditor reviews AI assessment), 8 (client receives gap explanation)

**Trigger:** Client uploads a document in response to a `DocumentRequest`.

**Inputs:**
- `DocumentRequest` (title, description, instructions, linked `Control`, linked `TestProcedure`)
- Uploaded document text (OCR pipeline output)
- `ControlObjective` description and `FrameworkRequirement` text
- RAG context: embeddings of prior accepted evidence for this request type (from same firm, anonymized)

**Model:** Claude Sonnet (reasoning quality is critical; errors propagate into the audit file and client-facing gap explanations require nuanced language).

**Process:** Extract document attributes (date range, parties, system names, configuration identifiers, format) → compare against request criteria → check period coverage against framework-specific windows (SOC 2 Type II per AT-C 320, ISO surveillance cycles, PCI 90-day ASV validity, HIPAA/HITRUST reporting periods) → score completeness against each test procedure → identify specific gaps → generate client-facing plain-language explanation → generate auditor-facing summary → produce recommendation: Accept | Request Clarification | Reject.

**Output:** `AIDecision` record with `context_type = document_completeness`, `suggested_value`, `confidence` score, and full `raw_output`. Status: `Pending`.

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

### Feature 2: Evidence → CommonControl Mapping Suggestion

**Purpose:** When a client uploads evidence, propose the three most likely `CommonControl` nodes it satisfies, with confidence scores and per-framework downstream coverage implications. This is the cross-framework differentiator made concrete at the evidence intake boundary — a single uploaded document can satisfy SOC 2 CC6.1, ISO 27001 A.5.16/A.5.17/A.5.18, and HITRUST access-control requirements simultaneously, and this feature surfaces all of them with one human confirmation.

**Journeys:** 7 (auditor triage of inbound evidence), 8 (client upload on Client Hub), auditee continuous-monitoring journey.

**Trigger:** `EvidenceItem` ingested — via Client Hub upload, integration pull (Okta/AWS/GitHub/Jira/etc.), or auditor-side manual attach.

**Inputs:**
- `EvidenceItem` content (extracted text, OCR output, configuration export, integration payload)
- `EvidenceItem` metadata (source system, file type, date, author, integration tags)
- Candidate `CommonControl` set from the SCF-derived library, filtered to controls in-scope for the engagement's active frameworks
- RAG context: prior confirmed evidence→CommonControl mappings for the same firm (and anonymized cross-firm patterns for cold-start)
- STRM edges from `CommonControl` into `FrameworkRequirement` for each active framework version

**Model:** Claude Haiku (classification/routing task; high volume; latency-sensitive in the upload path; semantic similarity heavy-lifted by pgvector retrieval).

**Process:** Embed the evidence content → retrieve top-k candidate `CommonControl` nodes via pgvector → for each candidate, score semantic fit against the control description and sample accepted evidence → emit top 3 candidates with confidence scores and per-framework downstream coverage implications computed from the STRM graph → generate one-line explanation per candidate ("Covers SOC 2 CC6.1 and ISO 27001 A.5.16 via the Access-Control CommonControl").

**Output:** Candidate `EvidenceControlMapping` records (status: `Pending`) with confidence scores and per-framework coverage deltas. `AIDecision` record with `context_type = evidence_control_mapping`.

**Human review gate:** Required. Auditor or auditee (depending on Client Hub policy) selects one or more candidates. The system then materializes `EvidenceControlMapping` edges with effective-dated STRM relationships. Axiom never shows a green checkmark when coverage is partial — confirmed mappings display percentage coverage and residual gap lists per framework.

**Failure modes:**

| Failure | Handling |
|---|---|
| No CommonControl candidate above 0.6 confidence | Evidence routed to "Needs manual triage" queue; auditor maps via search |
| Evidence content extraction incomplete (encrypted/image-only) | Fall back to filename + source metadata only; lower confidence noted |
| Engagement spans framework with no SCF mapping yet (ISO 42001 novel artifacts) | Firm-authored CommonControls used; AI flags "novel artifact — no cross-framework reuse available" |
| API error / timeout | Evidence uploaded successfully; mapping suggestions populate asynchronously |

---

### Feature 3: Cross-Framework Gap Analysis

**Purpose:** Given an engagement's scope (e.g., SOC 2 Type II + ISO 27001:2022 + ISO 27701:2019), show predicted coverage percentage per framework, missing evidence per control, and a prioritized recommendation of the next evidence collections that would maximize coverage across all active frameworks simultaneously. This is the economic argument for cross-mapping made visible.

**Journeys:** 3 (partner scopes integrated multi-framework engagement), continuous-assurance journey (auditee monitors standing readiness).

**Trigger:** (a) Engagement created/scoped with 2+ frameworks; (b) auditor or auditee opens the Coverage dashboard; (c) scheduled nightly recompute for engagements in active fieldwork.

**Inputs:**
- All in-scope `FrameworkRequirement` records per active framework and version
- All `CommonControl` nodes linked into scope via STRM edges
- All confirmed `EvidenceControlMapping` edges for the engagement
- Evidence freshness metadata (age, collection date, framework-specific staleness thresholds: ASV 90d, pen test 1y, background check 1y, etc.)
- Firm's standing evidence library (evidence collected across prior engagements that may be reusable)

**Model:** Claude Sonnet (reasoning over the control graph, partial-satisfaction arithmetic, and recommendation prioritization).

**Process:** For each framework in scope, walk STRM edges from `FrameworkRequirement` back to `CommonControl`s and forward to linked `EvidenceItem`s → compute per-requirement coverage (satisfied | partial | unsatisfied | stale) → aggregate to framework-level percentage → identify high-leverage gaps where a single new evidence collection would close multiple cross-framework requirements → rank recommendations by (cross-framework-coverage-gain / collection-effort) → emit narrative explanation of the top 5 recommendations.

**Output:** Coverage dashboard with per-framework progress, gap list per requirement, and a "Recommended next collections" panel. `AIDecision` record with `context_type = gap_analysis` summarizing the recommendation set and its rationale.

**Human review gate:** Required for any action the recommendations imply (creating new `DocumentRequest`s, enqueuing integration pulls, dispatching auditee questionnaires). The dashboard itself is informational.

**Failure modes:**

| Failure | Handling |
|---|---|
| Scope includes framework version with no STRM mapping | Framework displayed with "No cross-mapping available — maps manually" state; other frameworks still compute |
| Evidence has lapsed freshness threshold | Requirement marked `stale`, not `satisfied`; re-collection placed at top of recommendation list |
| Partial satisfaction (evidence covers 2 of 3 elements) | Displayed as percentage; Axiom never collapses partial to green |
| API error / timeout | Cached prior computation displayed with staleness indicator |

---

### Feature 4: Workpaper Narrative Draft

**Purpose:** Generate a first-draft workpaper narrative that a staff auditor would write, matching firm style — reducing blank-page time and ensuring consistent structure across the team.

**Journeys:** 5 (staff auditor requests draft after completing test procedure)

**Trigger:** Auditor marks a `TestProcedure` as Complete and explicitly requests a draft (not automatic).

**Inputs:**
- Control description and `TestProcedure` (type, description, expected result)
- Linked evidence items (extracted text and filenames)
- Exceptions noted (if any)
- Prior year workpaper narrative if rollforward (style reference)
- Firm workpaper template for this procedure type

**Model:** Claude Sonnet (quality matters; workpaper narratives become part of the immutable audit file).

**Output:** Draft text inserted into the workpaper editor. `WorkpaperVersion` record created with AI content tracking metadata (see [Section 5](#5-ai-content-tracking)). The draft is labeled "AI Draft — requires review" until substantive human edits are made.

**Human review gate:** Mandatory. The AI draft is always editable. The workpaper's `status` cannot advance to `PreparedPendingReview` or beyond if no substantive human edits have been made to AI-generated sections. The auditor must actively modify and sign off. This satisfies AT-C 105 / ISAE 3000 (Revised) / ISO 17021-1 documentation requirements to distinguish AI-generated from auditor-authored content, and is captured in the ISO 42001-aligned AIDecision ledger (see [Section 6](#6-iso-42001-native-positioning)).

**Failure modes:**

| Failure | Handling |
|---|---|
| Insufficient evidence linked | Draft generated with placeholder: "[Evidence not yet linked — add evidence before finalizing]" |
| Prior year workpaper not available (new engagement) | Draft uses firm template only; no style matching |
| Generated draft is excessively long or off-topic | Auditor deletes and re-generates or writes manually; AIDecision records the rejection |
| API error / timeout | Retry with backoff; auditor writes manually; no blocking dependency |

---

### Feature 5: Evidence Link Suggestion

**Purpose:** Propose which evidence items should be linked to a test procedure, reducing the manual search through hundreds of evidence items and surfacing cross-framework reuse opportunities.

**Journeys:** 5 (evidence linking during control testing), 7 (auto-link suggestion on document acceptance)

**Trigger:** Two triggers:
1. Auditor opens a test procedure for evidence linking — AI pre-populates suggested links.
2. Auditor accepts a document via the completeness review (Feature 1) — AI suggests which additional test procedures the evidence should link to beyond the originating `DocumentRequest`.

**Inputs:**
- `TestProcedure` description and expected evidence type
- `Control` description and `FrameworkRequirement` text
- `EvidenceItem` pool: extracted text summaries and metadata for all evidence in the engagement
- RAG context: prior engagement evidence links for the same control type (from same firm)
- For trigger 2: the `DocumentRequest` → `Control` → `TestProcedure` chain that originated the upload

**Model:** Claude Haiku (classification and matching task; high volume; speed matters during interactive evidence linking).

**Process:** Embed the test procedure description and expected evidence → retrieve top-k similar evidence items from the engagement pool → score relevance against the procedure's requirements → rank by confidence → return suggested links with explanation text. For trigger 2, extend suggestions beyond the originating chain to other test procedures that the evidence may satisfy.

**Output:** Suggested `EvidenceLink` records displayed in the evidence linking panel with confidence scores and explanation text. Each suggestion has `ai_suggested = true`. `AIDecision` record created per suggestion.

**Human review gate:** Required. The auditor sees suggested links and takes explicit action: Accept, Modify (link to a different procedure), or Reject. No `EvidenceLink` record is created until the auditor confirms. For trigger 2 (auto-link on acceptance), the primary link (from the `DocumentRequest` chain) is created automatically; additional cross-procedure suggestions require confirmation.

**Failure modes:**

| Failure | Handling |
|---|---|
| Empty evidence pool | No suggestions shown; auditor links manually after evidence is uploaded |
| All suggestions below confidence threshold | No auto-suggestions; evidence pool is browsable for manual linking |
| Evidence pool too large (>500 items) | Pre-filter by control area before AI ranking |
| API error / timeout | Manual linking available; suggestions populate asynchronously when service recovers |

---

### Feature 6: Risk Category Suggestion for Client Acceptance

**Purpose:** Accelerate the client-acceptance risk assessment by suggesting risk categories based on client industry, framework scope, prior engagement findings, and regulatory exposure — without pre-populating conclusions that would compromise the partner's independent judgment. Applies uniformly across SOC 2, ISO 27001/27701, ISO 42001, HIPAA, PCI DSS, and SOC 1 engagements.

**Journeys:** 3 (partner completes client acceptance before advancing to fieldwork)

**Trigger:** Partner opens the `ClientAcceptance` form for a new engagement.

**Inputs:**
- Client industry, entity type, data-sensitivity profile (PHI, PCI-regulated cardholder data, EU personal data, AI system inventory)
- Engagement type and framework(s) in scope
- Prior-period `ClientAcceptance` record (if rollforward engagement)
- Prior-period engagement exceptions and findings (if rollforward)
- Firm's quality policy categories (configured in firm settings)

**Model:** Claude Sonnet (professional judgment context; low volume — one call per engagement; quality of suggestions directly impacts regulatory documentation).

**Process:** Analyze client and engagement context → identify applicable risk categories from the firm's quality policy → assess heightened risk factors based on industry, framework-specific regulatory exposure (HIPAA breach history, PCI scope creep, ISO 42001 high-risk AI systems), engagement complexity, and prior-period findings → generate risk category suggestions with brief rationale per category → flag categories where prior-period findings indicate elevated risk.

**Output:** Suggested risk categories displayed in the `ClientAcceptance` form as selectable suggestions (not pre-populated fields). Each suggestion includes a one-line rationale. `AIDecision` record created with `context_type = risk_category_suggestion`.

**Human review gate:** Required. Suggestions are presented as a sidebar or suggestion panel — the actual `ClientAcceptance` risk documentation fields are empty until the partner writes or selects content. The partner certifies the acceptance independently; AI suggestions are reference material, not pre-filled conclusions. Firm quality management frameworks (including AICPA QM and equivalent ISO 17021-1 CB policies) require the partner's independent judgment; the AI assists the process, it does not perform it.

**Failure modes:**

| Failure | Handling |
|---|---|
| No prior year data (new client) | Suggestions based on industry and engagement type only; lower confidence noted |
| Firm quality policies not configured | Generic risk categories from standard quality-management guidance (AICPA QM, ISO 17021-1 equivalent); prompt firm to configure policies |
| API error / timeout | Acceptance form functions without suggestions; partner documents manually |

---

### Feature 7: Framework Version Migration Assistance

**Purpose:** When a framework version supersedes the one in use (ISO 27001:2013 → :2022, PCI DSS 3.2.1 → 4.0.1, NIST CSF 1.1 → 2.0, SOC 2 TSC revisions), propose a remapping of existing firm evidence and confirmed `EvidenceControlMapping` records to the new framework version's requirements. This is the step that makes cross-framework reuse durable across time — version churn otherwise invalidates the entire mapping graph at the control-ID level.

**Journeys:** Firm-admin framework lifecycle management; auditor prepares next-period engagement under new framework version.

**Trigger:** (a) Firm admin installs a new `FrameworkVersion`; (b) AI scheduler detects that an active engagement's framework version has been superseded and a cutover is due.

**Inputs:**
- Source `FrameworkVersion` and destination `FrameworkVersion`
- Published crosswalk data (AICPA official mappings, ISO normative annex for version-to-version mapping, PCI SSC summary of changes, SCF delta release notes)
- All `FrameworkRequirement` records for both versions, with their STRM edges to `CommonControl`
- All firm `EvidenceControlMapping` records currently tied to the source version
- Diff narrative from the standards body (where available)

**Model:** Claude Sonnet (reasoning over structural deltas; narrative explanation of why a mapping does/does not carry forward).

**Process:** For each requirement in the source version, identify destination-version equivalents using published crosswalks + STRM graph → classify each as (a) direct carry-forward, (b) renumbered-only, (c) scope-expanded (evidence still relevant but now insufficient), (d) scope-narrowed (evidence over-covers), (e) deprecated (no equivalent), (f) newly introduced (no prior evidence) → for each existing `EvidenceControlMapping`, propose a remapping action: keep | remap-to-new-ID | flag-for-gap-collection | retire → generate an auditor-facing migration report.

**Output:** Migration dashboard with a per-requirement disposition table and a summary of evidence impact (X mappings carry forward unchanged, Y need remapping, Z gaps require new collection). `AIDecision` record with `context_type = framework_migration`.

**Human review gate:** Required. The firm admin or engagement partner reviews each proposed disposition before any mapping is mutated. The migration is applied as a versioned, effective-dated transition — prior-period reports remain pinned to the superseded version and the `EvidenceControlMapping` edges they depend on are preserved via effective-dated STRM edges.

**Failure modes:**

| Failure | Handling |
|---|---|
| No published crosswalk for the version transition | AI flags this and falls back to semantic-similarity-only proposals with lower confidence |
| Destination version not yet available in SCF | Firm admin prompted to wait for SCF release; manual remapping UI still available |
| Large mapping set (1000+ edges) | Processed in background with batched AIDecision records; progress bar in UI |
| API error / timeout | Migration paused; resumable; no partial state persisted without human confirmation |

---

### Feature 8: Findings Triage & Severity Reasoning

**Purpose:** Ingest raw exceptions and test-procedure failures, cluster them by root cause, reason about severity against the firm's documented risk appetite and the framework-specific materiality standard, and propose a disposition (accept | remediate | escalate). Eliminates the hours-per-engagement spent manually grouping exceptions and writing severity rationale.

**Journeys:** 6 (manager reviews findings), 9 (partner finalizes disposition), EQR review journey.

**Trigger:** (a) A `TestProcedure` is recorded as failed; (b) batch run at fieldwork close on all open exceptions in the engagement.

**Inputs:**
- All open `Exception` / `Finding` records for the engagement with supporting evidence
- Control(s) implicated, their CommonControl cluster, and the upstream `FrameworkRequirement`s affected
- Firm's documented risk-appetite policy (configured in firm settings)
- Framework-specific severity vocabularies (SOC 2 deficiency/significant-deficiency/material-weakness-equivalent, ISO 17021-1 nonconformity major/minor, PCI compensating-control eligibility, HIPAA breach-notification thresholds)
- Prior-period findings on same client (recurrence signal)

**Model:** Claude Sonnet (multi-factor reasoning; full AIDecision chain with explicit rationale required).

**Process:** Cluster findings by root cause using embedding similarity + shared CommonControl → for each cluster, score severity across four axes (likelihood, impact scope, framework-specific classification, recurrence) → compare against firm risk appetite → propose disposition with explicit rationale referencing the firm's policy → identify escalation candidates (items that cross a materiality or regulatory-reporting threshold).

**Output:** Findings board grouped by cluster, with proposed severity, proposed disposition, and a rationale narrative per cluster. `AIDecision` record with `context_type = findings_triage` per cluster.

**Human review gate:** Required. The manager (and partner for escalation candidates) reviews each cluster. Disposition is never finalized without explicit human action. Escalation to management-response drafting (Feature 10) is a human-initiated step.

**Failure modes:**

| Failure | Handling |
|---|---|
| Firm risk-appetite policy not configured | AI uses framework-default severity thresholds and flags this explicitly; prompts firm admin to configure |
| Single-finding clusters (no similar exceptions) | Processed individually; no clustering benefit but severity reasoning still applies |
| Ambiguous root cause (evidence supports multiple interpretations) | AI emits multiple candidate clusters with confidence; auditor selects |
| API error / timeout | Findings remain in untriaged queue; manual triage UI available |

---

### Feature 9: Drift-Triggered Re-testing

**Purpose:** Continuous assurance. When a monitored configuration drifts (e.g., an IAM policy grants a new wildcard permission, an S3 bucket becomes public, a vendor security group opens port 22 to 0.0.0.0/0, an ISO 42001 AI system's model version changes outside change-control), autonomously re-run the affected control tests, re-assess residual risk, update the risk register, and notify the auditor with a diff. This is the feature that turns Axiom from point-in-time audit tooling into a continuous-assurance platform, and it is how Axiom competes with Delve/Drata/Vanta on the auditee side without inheriting their provenance problems.

**Journeys:** Continuous assurance (auditee), auditor real-time engagement dashboard.

**Trigger:** (a) Integration webhook fires a drift event (Okta policy change, AWS Config rule non-compliant, GitHub branch protection removed, model registry update); (b) scheduled polling for integrations without push webhooks.

**Inputs:**
- Drift event payload (before/after configuration state, timestamp, actor)
- Control(s) this configuration attests to (via `EvidenceItem` → `CommonControl` → `FrameworkRequirement` graph)
- Most recent passing test result for affected controls
- Firm's drift-response policy (re-test immediately | batch nightly | notify-only)

**Model:** Claude Haiku for detection classification and diff narrative generation. Claude Sonnet invoked only when severity reasoning is required (drift crosses a materiality threshold or implicates multiple frameworks).

**Process:** Classify the drift (material / immaterial / reversal) → if material, enqueue re-tests of affected controls as River jobs → on re-test completion, compare result to prior passing state → if regressed, generate a concise diff narrative ("IAM policy `DevelopersAdmin` added `s3:*` on `arn:aws:s3:::*` — this removes the scoped-access evidence that satisfied SOC 2 CC6.1 and ISO 27001 A.5.18 on 2026-03-04"), update the risk register entry, and notify the engagement team.

**Output:** Drift event timeline entry, re-test `TestResult` record, updated risk-register entry, and notification payload. `AIDecision` record with `context_type = drift_retest` for any state transition that affects control status.

**Tier behavior:**
- Detection and drift classification: Tier 1 (autonomous).
- Automatic re-test execution: Tier 1 when no control-status change results; Tier 2 the moment the re-test would change a control status (`Passing` → `Exception`) — the status change waits on human confirmation.
- Risk-register update: Tier 2 (auditor confirms).
- Notification to the auditee: Tier 1 (informational).

**Human review gate:** Required for any risk-register mutation or control-status change. Detection and re-test execution proceed autonomously; the conclusion does not.

**Failure modes:**

| Failure | Handling |
|---|---|
| Integration webhook missed / duplicated | Idempotent event processing via event ID; reconciliation job nightly |
| Drift is a reversal (config restored) | Classified as `reversal`; prior state re-validated; no re-test fires |
| Re-test requires evidence that has expired | Drift flagged as "re-test blocked on stale evidence"; collection request auto-proposed |
| API error / timeout | Event queued with exponential-backoff retry; eventual consistency acceptable |

---

### Feature 10: Agentic Management-Response Drafting

**Purpose:** When a finding fires, an agent carries the full round-trip: draft the management response, identify the control owner by querying the client's HRIS (BambooHR/Workday/Rippling) for the role mapped to that CommonControl, open a remediation ticket in Jira / Linear / ServiceNow with pre-filled context, track closure evidence as it arrives, and re-test the control. This is the agent feature — but every externally-visible side effect is human-gated, so the "agent" is really a staged pipeline with explicit approval points, not an unsupervised autonomous actor.

**Journeys:** Findings-to-remediation journey (cross-cutting); management-response workflow.

**Trigger:** A confirmed `Finding` with disposition = `Remediate` (from Feature 8).

**Inputs:**
- `Finding` with severity, implicated CommonControl(s), evidence, and disposition
- Client's HRIS integration data (control-owner role → named person)
- Client's ticketing integration config (Jira project, Linear team, or ServiceNow assignment group)
- Firm's management-response template for the framework(s) affected
- Prior successful remediations for similar findings on same or peer clients (few-shot context)

**Model:** Claude Sonnet (drafting + multi-step planning).

**Process (staged, each stage gated):**

1. **Draft management response** — generate narrative acknowledging the finding, proposing remediation, and committing to a timeline consistent with framework-specific windows (HIPAA breach windows, PCI Target Remediation Dates, ISO nonconformity correction periods). Tier 2 approval before draft is saved to the engagement file.
2. **Identify control owner** — query HRIS for the role assigned to the CommonControl; propose a named person. Tier 2 approval before the person is assigned as ticket owner or emailed.
3. **Open remediation ticket** — pre-fill ticket body with finding context, evidence links, and remediation plan; pass through human approval; on approval, Axiom makes the external API call. Tier 2 approval mandatory before any write to an external system.
4. **Track closure evidence** — when evidence referencing the ticket ID is uploaded (Client Hub or integration), link it to the ticket and the finding; Tier 1.
5. **Re-test** — when the ticket is marked done and closure evidence is present, queue a re-test via Feature 9's machinery; the re-test's conclusion is Tier 2 as usual.

**Output:** Draft management-response text, ticket link, closure-evidence trail, re-test result. `AIDecision` records with `context_type = management_response_draft` per stage.

**Human review gate:** Mandatory at stages 1, 2, 3, and at the re-test conclusion. External side effects (tickets, emails) never fire without explicit human approval. Axiom will not autonomously assign blame to a named employee.

**Failure modes:**

| Failure | Handling |
|---|---|
| HRIS integration not configured | Control-owner step falls back to engagement-team-designated reviewer; manual assignment prompt |
| Ticketing integration not configured | Stage 3 emits a downloadable ticket template (Markdown) for manual filing |
| Closure evidence never arrives within framework window | Escalation notification; finding remains open; ISO/PCI/HIPAA timing implications surfaced |
| Re-test fails after remediation | Loop: new finding created; severity reasoning reruns with recurrence signal |
| API error during external write | Write is transactional; on failure the AIDecision is marked as not-executed; no partial state |

---

### Feature 11: Report Section Draft

**Purpose:** Generate first drafts of standardized report sections (Description of Tests of Controls, Scope and Approach, System Description for SOC reports) so the partner writes opinion and judgment sections rather than boilerplate.

**Journeys:** 9 (partner requests AI draft of specific report sections)

**Trigger:** Partner explicitly requests a draft of a specific report section in the report editor (not automatic).

**Inputs:**
- Report type and template (SOC 2 Type I / Type II, SOC 1 Type II, ISO 27001 Stage 2 / surveillance, ISO 42001 management-system audit report, HIPAA assessment / HITRUST CSF r2 report, PCI DSS Report on Compliance (ROC) / SAQ, etc.)
- Engagement-wide data: all `Control` and `CommonControl` records with status, all `TestProcedure` results, all `Finding` records and their dispositions, `EngagementFramework` details
- `EvidenceItem` summary statistics (count by type, coverage by control area, cross-framework reuse counts)
- Prior-period report (if rollforward — style and structure reference)
- Firm report template for this report type

**Model:** Claude Sonnet (report content becomes part of the issued deliverable; quality, accuracy, and professional tone are critical).

**Process:** Aggregate engagement data by report section → populate template sections with structured data (control counts, exception summaries, testing methodology descriptions) → generate narrative for standardized sections → reference framework-specific language (AICPA/AT-C, SSAE 18, ISO 17021-1, PCI SSC, HITRUST) → produce draft with engagement-specific content rather than generic boilerplate.

**Output:** Draft text inserted into the specific report section in the editor. `ReportVersion` record created. The draft section is labeled "AI Draft — requires review" with the same AI content tracking metadata as workpapers (see [Section 5](#5-ai-content-tracking)). `AIDecision` record created with `context_type = report_section_draft`.

**Human review gate:** Mandatory. Same rules as Feature 4 (workpaper drafts): the report section cannot be finalized if no substantive human edits have been made to AI-generated content. The partner must actively modify and sign off. Report issuance (the `Report.status = Issued` transition) validates that all AI-drafted sections have been substantively edited.

**Sections AI may draft:**

| Section | Report Types | Notes |
|---|---|---|
| Description of Tests of Controls | SOC 1, SOC 2 | Populated from TestProcedure records and results |
| Scope and Approach | All | Populated from engagement configuration and framework details |
| System Description summary | SOC 1, SOC 2 | Populated from client-provided system description evidence |
| Statement of Applicability narrative | ISO 27001, ISO 27701 | Populated from CommonControl inclusion/exclusion and rationale |
| AI System Inventory and Impact Assessment summaries | ISO 42001 | Populated from client-provided AI system inventory and firm-authored impact assessments |
| Control testing summary by requirement | PCI DSS ROC | Populated from TestProcedure records mapped to PCI requirement tree |
| Assessment methodology narrative | HIPAA / HITRUST | Populated from engagement configuration and HITRUST r2 maturity scoring |

**Sections AI does NOT draft:**

| Section | Reason |
|---|---|
| Independent auditor's / assessor's opinion | Professional judgment; Tier 3 — human-initiated only |
| Management Assertions | Client-authored; not auditor content |
| Qualification / adverse / nonconformity language | Professional judgment with regulatory consequences |
| Emphasis of Matter / Other Matter | Professional judgment |
| Certification recommendation (ISO CB) | Restricted by ISO 17021-1 to accredited certification decision |
| QSA attestation of compliance | Restricted by PCI SSC to the named QSA |

**Failure modes:**

| Failure | Handling |
|---|---|
| Incomplete engagement data (controls not all complete) | Draft generated with placeholders for incomplete sections; partner warned |
| No prior year report (new client/engagement type) | Draft uses firm template only; no style matching |
| Report type not supported for AI drafting | AI draft button not shown; partner writes manually |
| API error / timeout | Retry with backoff; partner writes manually; no blocking dependency |

---

## 5. AI Content Tracking

### Problem

The original `is_ai_draft` boolean on `WorkpaperVersion` is insufficient. It answers "has any human edit been made?" but cannot answer "how much of the AI content was substantively reviewed?" — a question that independent reviewers (EQR, ISO CB technical reviewers), managers (Journey 6), and partners need answered to evaluate work quality, and a question that ISO 42001 auditability explicitly requires Axiom to answer about itself.

A binary flag creates a perverse incentive: an auditor can change a single comma, clear the gate, and advance a workpaper that is 99% unreviewed AI output. The EQR reviewer has no visibility into this.

### Design

AI content tracking operates at the **section level** within workpapers and reports. Every document type that supports AI drafting (workpapers via Feature 4, reports via Feature 11) tracks AI origin per section.

**Data model addition** (on `WorkpaperVersion` and `ReportVersion`):

```
ai_content_metadata  jsonb  nullable
```

Schema:
```json
{
  "sections": [
    {
      "section_id": "scope-and-approach",
      "ai_generated": true,
      "ai_generated_at": "2026-03-15T10:30:00Z",
      "human_edited": true,
      "human_edited_by": "user-uuid",
      "human_edited_at": "2026-03-15T11:15:00Z",
      "ai_character_count": 1450,
      "current_character_count": 1620,
      "modification_ratio": 0.42
    }
  ],
  "summary": {
    "total_sections": 8,
    "ai_generated_sections": 5,
    "ai_sections_edited": 4,
    "ai_sections_unedited": 1,
    "overall_modification_ratio": 0.38
  }
}
```

**`modification_ratio`** is computed as Levenshtein distance between AI-generated text and current text, divided by AI-generated character count. A ratio of 0.0 means no changes; 1.0 means completely rewritten. This is computed on save, not in real-time.

### Gate Logic

The advancement gate (`PreparedPendingReview` for workpapers, `Issued` for reports) checks:

1. **All AI-generated sections must have `human_edited = true`** — at least one edit per AI section. This replaces the binary `is_ai_draft` check.
2. If any AI-generated section has `modification_ratio < 0.05` (less than 5% change), a warning is surfaced: "Section [name] has minimal edits to AI-generated content. Confirm this section reflects your professional judgment." The auditor must explicitly confirm to proceed. This is a soft gate (confirmable), not a hard block.

### Reviewer Visibility

**Manager review (Journey 6):** The review view shows an AI-edit indicator per workpaper: "This workpaper was AI-drafted. The auditor modified 4 of 6 AI sections. Average modification: 38%."

**EQR review (Journey 10):** The EQR-focused view includes an AI substantiveness summary across the entire engagement: "42 workpapers used AI drafts. Average modification ratio: 35%. 3 workpapers have sections with <10% modification — review these first." The EQR reviewer can filter to low-modification workpapers for priority review.

### Backward Compatibility

The existing `is_ai_draft` boolean field is retained as a derived convenience field: `is_ai_draft = true` when `ai_content_metadata` contains any section with `ai_generated = true AND human_edited = false`. Application code that checks `is_ai_draft` continues to work. The section-level metadata provides the granularity that the boolean cannot.

### `AIDecision.context_type` Enum — Authoritative List

Every AI action that requires human review creates an `AIDecision` row with a `context_type` drawn from this enum. The OpenAPI schemas and Domain & Data Model spec reference this list verbatim.

| `context_type` | Feature | Notes |
|---|---|---|
| `document_completeness` | 1 | Client upload completeness assessment |
| `evidence_control_mapping` | 2 | Candidate CommonControl(s) for an EvidenceItem |
| `gap_analysis` | 3 | Cross-framework coverage recommendation set |
| `workpaper_draft` | 4 | Workpaper narrative AI draft |
| `evidence_link_suggestion` | 5 | Suggested EvidenceLink records |
| `risk_category_suggestion` | 6 | Client-acceptance risk categories |
| `framework_migration` | 7 | Requirement remapping on framework version change |
| `findings_triage` | 8 | Severity reasoning + disposition proposal per cluster |
| `drift_retest` | 9 | Autonomous re-test conclusion / risk-register update |
| `management_response_draft` | 10 | Each staged agent action (draft / owner / ticket / re-test) |
| `report_section_draft` | 11 | Report section AI draft |

Any new AI surface requires a new enum value added here, a matching migration in the Domain & Data Model spec, and a matching OpenAPI update.

---

## 6. ISO 42001-Native Positioning

Axiom is an ISO 42001-native platform. ISO 42001 is the international management-system standard for artificial intelligence, and it requires organizations that deploy AI to operate a documented lifecycle with governance, risk management, impact assessment, and human oversight. Axiom both *assesses* clients against ISO 42001 (Feature 11 drafts the ISO 42001 management-system audit report) and *operates under* ISO 42001 itself. This dogfooding is the defensible differentiator.

### How Axiom Satisfies ISO 42001 Operationally

The three-tier HITL policy (Section 3), the `AIDecision` ledger, and the `ai_content_metadata` audit trail together are the operational realization of ISO 42001 clauses 6 (planning), 7 (support), 8 (operation), and 9 (performance evaluation) for Axiom as a deployer of AI. Specifically:

- **Every AI invocation is logged.** The `AIDecision` row captures: `prompt_hash`, `model_id` and `model_version`, `context_type`, full `input_payload` (or hash + pointer to content-addressed store for large inputs), `raw_output`, `suggested_value`, `confidence`, `reviewer_user_id`, `reviewer_action` (accept | modify | reject | override), `reviewer_rationale` (when override), `reviewed_at`. No AI action that touches audit content bypasses this ledger.
- **Prompt, model, context, confidence, reviewer, override** — all six elements are first-class columns on `AIDecision`. An auditor (internal or external) can reconstruct exactly what the model saw, what it said, who reviewed it, and why they agreed or disagreed.
- **Impact-assessment-ready.** Feature-level impact assessments (the ISO 42001 deliverable) are pre-populated from feature metadata: each of the eleven features has a structured impact-assessment record covering intended use, foreseeable misuse, affected stakeholders, severity class, mitigation controls, and human-oversight mechanism. These records are the default view in the firm-admin "AI governance" screen.
- **Model-change management.** Bedrock model version pins are a configuration artifact with a change-control workflow: changing a model ID (e.g., Sonnet 4.6 → Sonnet 4.7) is a tracked config change with its own approval ledger. Prior `AIDecision` rows preserve the model version they were generated against, so retrospective review is unambiguous.
- **Continuous monitoring.** Feature 9 (drift-triggered re-testing) and the AI observability dashboard jointly satisfy the ISO 42001 performance-evaluation clause: quality, drift, confidence-distribution, and override-rate metrics are monitored per feature and alarm on anomaly.
- **Human oversight by design.** Tier 3 actions cannot be AI-initiated. Tier 2 actions cannot take effect without a named human. Tier 1 actions cannot alter audit content. The tiering is enforced in code, not policy.

### Competitive Positioning

The March 2026 trust crisis in the agentic-compliance category (see `compliance-pivot-findings.md` §1) made unverifiable AI output commercially toxic. Axiom's answer is not to reduce AI surface area but to make every AI output individually traceable, reviewer-attested, and cryptographically signed (Section 7). The AIDecision ledger is the artifact that a regulator, an ISO CB technical reviewer, or a client's procurement team can audit. This is the difference between *claiming* AI governance and *being* auditable under AI governance — and it is what a post-boilerplate-scandal buyer is now explicitly asking for.

Axiom does not name competitors in-product. The positioning is carried by the artifacts.

---

## 7. Cryptographic Provenance for AI Outputs

Every AI-generated artifact is cryptographically signed at creation and written to WORM storage. "Auditor-defensible by construction" is not a marketing tagline — it is the data path.

### Signed Artifacts

Provenance signing applies to:

- Workpaper AI drafts (Feature 4)
- Evidence → CommonControl mapping suggestions (Feature 2)
- Gap-analysis recommendation sets (Feature 3)
- Framework-migration disposition reports (Feature 7)
- Findings triage clusters and severity rationale (Feature 8)
- Re-test conclusions from drift events (Feature 9)
- Management-response drafts and each staged agent action (Feature 10)
- Report section drafts (Feature 11)

### Signing Envelope

At AI-output emission, Axiom constructs a provenance envelope containing: `artifact_id`, `ai_decision_id`, `engagement_id`, `model_id`, `model_version`, `prompt_hash` (SHA-256 of the canonicalized prompt), `input_content_hashes` (SHA-256 per input payload, including retrieved context), `output_hash` (SHA-256 of the raw output), `generated_at` (UTC, from authoritative time source), and `axiom_version`. The envelope is serialized as canonical JSON and signed with an AWS KMS asymmetric key (ECC_NIST_P256, SIGN_VERIFY) dedicated to AI-artifact signing. The signature, the envelope, and the raw output are written together to an S3 Object Lock (WORM) bucket with a retention period matching the engagement's record-retention policy (minimum 7 years).

On subsequent edit (human modification of an AI draft), a new envelope is signed for the edited artifact with a pointer to the parent envelope, forming a tamper-evident chain. The `ai_content_metadata.modification_ratio` computation (Section 5) references the signed original.

### Verification

A verification endpoint (and a downloadable CLI) lets any party — an EQR reviewer, an ISO CB technical reviewer, a client's security team, or a regulator — re-compute the output hash and verify the signature against Axiom's public key (published at a well-known URL and transparency-logged). This gives the verifier cryptographic certainty that the artifact they hold is byte-identical to what Axiom generated at the stated time using the stated model — the exact property whose absence caused the March 2026 category trust crisis.

### Cost and Performance Impact

Signing adds ~8ms per AI output on the ECS task (KMS Sign API) and ~50 bytes of signature payload. KMS Sign at 1M calls/year is roughly $30. S3 Object Lock storage cost scales with artifact volume; for a mid-market firm at 100 engagements/year, signed-artifact storage is <$10/month. The cost of provenance is rounding-error relative to the trust dividend.

---

## 8. Per-Engagement AI Cost Estimate

Cost estimates use Bedrock on-demand pricing with prompt caching (0.1× base input rate for cached reads). System prompts, framework-requirement text, CommonControl library excerpts, and firm methodology context are cached across calls within the same engagement session. All figures below are after-caching effective cost. Figures assume mid-market complexity; large-enterprise engagements scale roughly linearly with evidence volume.

### Per-Engagement Cost by Framework

| Framework | Typical Volume Drivers | Effective Cost Range |
|---|---|---|
| **SOC 2 Type II** | 200 evidence items, 120 controls, 20 workpaper drafts, 4 report sections | **~$3–5** |
| **ISO 27001 / 27701** | 93 Annex A controls, 150 evidence items, embedding cost for new requirement text and policy library | **~$3–5** |
| **ISO 42001** | AI system inventory + firm-authored impact assessments (token-heavy narrative), 38 controls, gap analysis across AI-lifecycle | **~$4–6** |
| **HIPAA (with HITRUST CSF r2 path)** | Administrative/physical/technical safeguards, ~180 evidence items, maturity-scoring narrative | **~$3–5** |
| **PCI DSS 4.0.1** | 12 requirements / 300+ sub-requirements, ASV scan parsing, penetration-test report extraction, network-diagram evidence (higher parsing cost) | **~$5–8** |
| **SOC 1 Type II** | User-entity control objective narratives, process-narrative drafting | **~$3–5** |

### Feature-Level Cost Composition (SOC 2 Type II Reference)

| AI Feature | Volume | Model | Approx. Cost |
|---|---|---|---|
| Document completeness review (Feature 1) | 200 docs | Sonnet | ~$4.00 |
| Evidence → CommonControl mapping (Feature 2) | 200 evidence items | Haiku | ~$0.20 |
| Cross-framework gap analysis (Feature 3) | ~20 recomputes | Sonnet | ~$0.40 |
| Workpaper narrative drafts (Feature 4) | 20 drafts | Sonnet | ~$1.20 |
| Evidence link suggestions (Feature 5) | 200 links | Haiku | ~$0.15 |
| Risk category suggestion (Feature 6) | 1 engagement | Sonnet | ~$0.05 |
| Framework migration (Feature 7) | 0 per typical engagement; amortized on version cutover | Sonnet | ~$0.00 |
| Findings triage (Feature 8) | ~10 clusters | Sonnet | ~$0.30 |
| Drift-triggered re-tests (Feature 9) | ~50 drift events in Type II window | Haiku + selective Sonnet | ~$0.25 |
| Agentic management-response drafting (Feature 10) | ~5 remediations | Sonnet | ~$0.40 |
| Report section drafts (Feature 11) | 4 sections | Sonnet | ~$0.25 |
| RAG / embedding overhead | ~60K tokens embedded + retrieval | Haiku | ~$0.10 |
| **Subtotal before caching discount** | | | **~$7.30** |
| **After prompt caching** | | | **~$3–5** |

Multi-framework integrated engagements (e.g., SOC 2 + ISO 27001 + ISO 27701 performed together) do not multiply linearly — the CommonControl graph means evidence is embedded once, gap analysis runs once per scope, and the marginal cost of adding a second framework is roughly 30–40% of its standalone cost.

### Caveats

These numbers are Bedrock on-demand with prompt caching. Batch inference (50% discount) is used where latency permits (Feature 2 nightly re-scans, Feature 9 batch drift classification). Figures exclude:

- pgvector embedding storage and HNSW index memory (amortized into database cost, not AI cost).
- KMS sign operations for provenance (~$0.03/engagement; negligible).
- S3 Object Lock WORM storage for signed AI artifacts (~$1–2/engagement at Type II retention).

### Margin Impact

At 100 engagements/year for a mid-market firm, platform AI costs are **$300–$800/year per firm**. This is absorbed into the subscription — AI consumption is not a separate line item on customer invoices at launch. AI cost represents <0.1% of subscription revenue per engagement across all framework types; the margin impact is negligible even at the PCI DSS top of the range.

---

## 9. What AI Does NOT Do at Launch

- Fine-tuned models per firm (insufficient training data volume at this stage)
- On-device / local model execution
- Multi-agent orchestration frameworks (LangGraph, CrewAI) — the staged-pipeline pattern used in Feature 10 is implemented as River-scheduled Tier-gated steps, not as an autonomous planner; framework agents add debugging overhead with no benefit for the defined use cases
- AI-generated certification decisions, attestation opinions, QSA sign-off, or HITRUST assessor attestations
- AI-generated authoritative cross-framework equivalence declarations (AI suggests mappings; the authoritative crosswalk is SCF/OSCAL/AICPA-sourced plus firm-confirmed)
- Autonomous multi-step actions without human review gates at each step
- External side effects (ticket creation, emails to named employees, risk-register mutations that change control status) without a named human approval on the AIDecision
- Auto-finalization of any audit content without a named human reviewer
- Real-time client-facing AI pre-check on upload (the completeness review runs asynchronously; the auditor reviews before the client receives feedback — adding a client-facing instant pre-check is a post-launch evaluation item that requires careful UX design to avoid creating false confidence in incomplete documents)
- Use of Opus at launch (reserved for future multi-framework-planning workloads; any Opus routing gated behind a feature flag and cost guardrail)

---

*End of Axiom AI Architecture Design*
