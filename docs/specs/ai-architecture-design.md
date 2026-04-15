# Axiom AI Architecture Design

**Date:** April 15, 2026
**Status:** Implementation-Ready
**Supersedes:** Section 6 of axiom-spec-design.md

> Research: [05-ai-architecture.md](../research/05-ai-architecture.md)

This document is the authoritative specification for all AI decisions in Axiom. The [product spec](axiom-spec-design.md) Section 6 contains a summary; this document contains the full detail.

---

## Table of Contents

1. [LLM Provider Decision](#1-llm-provider-decision)
2. [Vector Database Decision](#2-vector-database-decision)
3. [Human-in-the-Loop Policy](#3-human-in-the-loop-policy)
4. [AI Features](#4-ai-features)
5. [AI Content Tracking](#5-ai-content-tracking)
6. [Per-Engagement AI Cost Estimate](#6-per-engagement-ai-cost-estimate)
7. [What AI Does NOT Do at Launch](#7-what-ai-does-not-do-at-launch)

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

Embeddings are stored for: firm methodology documents, framework requirement descriptions, control objective library entries, prior engagement evidence summaries, and workpaper templates.

---

## 3. Human-in-the-Loop Policy

### Tier 1 — Fully Automated (no human approval required)

AI may act without user interaction for: text extraction from uploaded documents, embedding generation, flagging documents as potentially incomplete (notification only, not decision), generating overdue reminder emails to clients, surfacing "this evidence was used in a prior engagement" notifications, computing trial balance analytics (variance, ratios).

### Tier 2 — Human Approval Required Before Taking Effect

AI suggests; human confirms. No audit file content changes without explicit auditor action. Every Tier 2 action creates an `AIDecision` record. Features in this tier:

- Document completeness review (Feature 1)
- Control-to-framework-requirement mappings (Feature 2)
- Trial balance account category assignments (Feature 3)
- Workpaper narrative drafts (Feature 4)
- Evidence link suggestions (Feature 5)
- Risk category suggestions for client acceptance (Feature 6)
- Report section drafts (Feature 8)

### Tier 3 — Human-Initiated Only (AI may assist but never initiates)

These actions can only be triggered by a named human: engagement status transitions, control conclusions (pass/fail), exception documentation, workpaper review and sign-off, EQR sign-off, report issuance, client acceptance documentation, and any action constituting professional judgment under AU-C or PCAOB standards.

This three-tier policy maps directly to PCAOB AS 1105: AI outputs inform professional judgments; they do not constitute them.

### PCAOB Tier Elevation for Trial Balance Anomaly Detection

Trial balance anomaly detection (Feature 7) operates at **Tier 1** for nonissuer engagements — informational flags in the UI, no AIDecision record required.

For **PCAOB engagements**, AS 1105 requires that all technology-assisted analytical procedures be documented. Anomaly detection flags in PCAOB engagements therefore create `AIDecision` records and must be reviewed by the auditor, effectively elevating to **Tier 2 behavior** for those engagements. The feature implementation checks `engagement.type` and applies the appropriate tier.

---

## 4. AI Features

Eight AI features at launch, organized by the user journey they support. Every feature follows the same structure: purpose, trigger, inputs, model, process, output, human review gate, and failure modes.

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

**Process:** Extract document attributes (date range, parties, system names, amounts, format) → compare against request criteria → check period coverage for Type II engagements (AT-C 320) → score completeness against each test procedure → identify specific gaps → generate client-facing plain-language explanation → generate auditor-facing summary → produce recommendation: Accept | Request Clarification | Reject.

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

### Feature 2: Control Mapping (Framework-Agnostic Evidence Linkage)

**Purpose:** Propose `FirmControlObjectiveMapping` records across all frameworks in the engagement simultaneously — the architectural realization of the cross-framework differentiator.

**Journeys:** 3 (partner reviews proposed mappings at engagement creation)

**Trigger:** New engagement created from a methodology template, or new `FirmControlObjective` added to an engagement.

**Inputs:**
- `FirmControlObjective` name + description
- All `FrameworkRequirement` records for frameworks in the engagement
- RAG context: `ControlObjectiveLibrary` entries with existing mappings as few-shot examples

**Model:** Claude Haiku (structured classification task; high volume; speed matters at template instantiation).

**Process:** Embed the `FirmControlObjective` description → retrieve top-k similar library entries → score semantic similarity against each applicable `FrameworkRequirement` → apply 0.75 similarity threshold → generate explanation per proposed mapping → return proposed mappings with confidence scores.

**Output:** Mapping table displayed to auditor showing each proposed `FrameworkRequirement` link with confidence and explanation text. All proposed mappings remain in `Pending` state until confirmed.

**Human review gate:** Auditor reviews proposed mappings in bulk. All confirmed by default; any can be rejected. No `FirmControlObjectiveMapping` record is created until confirmed. `AIDecision` records created per mapping.

**Failure modes:**

| Failure | Handling |
|---|---|
| No similar library entries found | Mapping proposed with lower confidence; flagged for manual review |
| All mappings below 0.75 threshold | No auto-suggestions; auditor maps manually |
| Large template (200+ controls) | Process in background with notification on completion |
| API error / timeout | Retry with backoff; auditor notified; engagement creation proceeds without mappings |

---

### Feature 3: Trial Balance Account Mapping

**Purpose:** Classify each trial balance account into a standard financial statement line item, eliminating manual mapping work at engagement start.

**Journeys:** 4 (staff auditor reviews AI mappings after TB import)

**Trigger:** Trial balance imported (CSV/Excel upload or API).

**Inputs:** Account number, account name, balance (debit/credit), prior year mapping if rollforward.

**Model:** Claude Haiku (simple classification; high volume; prior year context improves accuracy significantly on rollforward engagements).

**Process:** Few-shot classification against standard FS line items (Cash, Accounts Receivable, Fixed Assets, etc.) with anomaly flagging for accounts where classification confidence is low or account name doesn't match expected category.

**Output:** `TrialBalanceAccount.mapping_status = AISuggested` with `ai_decision_id` populated for each mapped account. Low-confidence accounts highlighted in the Sheets UI for priority auditor review.

**Human review gate:** Auditor reviews and confirms mappings in the Sheets interface. Unmapped and low-confidence accounts surfaced prominently. `mapping_status` changes to `Confirmed` or `Overridden` on auditor action.

**Failure modes:**

| Failure | Handling |
|---|---|
| Non-standard account names (abbreviations, jargon) | Lower confidence score; flagged for manual classification |
| Foreign-language account names | Flagged for manual classification |
| Duplicate account numbers | Import validation catches before AI runs |
| API error / timeout | Retry with backoff; auditor maps manually; TB import is not blocked |

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

**Human review gate:** Mandatory. The AI draft is always editable. The workpaper's `status` cannot advance to `PreparedPendingReview` or beyond if no substantive human edits have been made to AI-generated sections. The auditor must actively modify and sign off. This satisfies the PCAOB requirement to distinguish AI-generated from auditor-authored content.

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

### Feature 6: Risk Category Suggestion

**Purpose:** Accelerate the SQMS 1 client acceptance process by suggesting quality risk categories based on client industry, engagement type, and prior engagement findings — without pre-populating conclusions that would compromise the partner's independent judgment.

**Journeys:** 3 (partner completes client acceptance before advancing to fieldwork)

**Trigger:** Partner opens the `ClientAcceptance` form for a new engagement.

**Inputs:**
- Client industry and entity type
- Engagement type and framework(s)
- Prior year `ClientAcceptance` record (if rollforward engagement)
- Prior year engagement exceptions and findings (if rollforward)
- Firm's quality policy categories (configured in firm settings)

**Model:** Claude Sonnet (professional judgment context; low volume — one call per engagement; quality of suggestions directly impacts regulatory compliance documentation).

**Process:** Analyze client and engagement context → identify applicable risk categories from the firm's quality policy → assess heightened risk factors based on industry, engagement complexity, and prior year findings → generate risk category suggestions with brief rationale per category → flag categories where prior year findings indicate elevated risk.

**Output:** Suggested risk categories displayed in the `ClientAcceptance` form as selectable suggestions (not pre-populated fields). Each suggestion includes a one-line rationale. `AIDecision` record created with `context_type = RiskCategorySuggestion`.

**Human review gate:** Required. Suggestions are presented as a sidebar or suggestion panel — the actual `ClientAcceptance` risk documentation fields are empty until the partner writes or selects content. The partner certifies the acceptance independently; AI suggestions are reference material, not pre-filled conclusions. SQMS 1 requires the partner's independent judgment; the AI assists the process, it does not perform it.

**Failure modes:**

| Failure | Handling |
|---|---|
| No prior year data (new client) | Suggestions based on industry and engagement type only; lower confidence noted |
| Firm quality policies not configured | Generic risk categories from standard SQMS 1 guidance; prompt firm to configure policies |
| API error / timeout | Acceptance form functions without suggestions; partner documents manually |

---

### Feature 7: Trial Balance Anomaly Detection

**Purpose:** Flag unusual account activity, outlier balances, and unexpected patterns in the trial balance for auditor investigation — turning raw financial data into directed inquiry starting points.

**Journeys:** 4 (staff auditor reviews anomaly flags during analytical procedures)

**Trigger:** Nightly background job on all engagements in Fieldwork status with an imported trial balance. Also runs once immediately after initial TB import and account mapping confirmation.

**Inputs:**
- All `TrialBalanceAccount` records for the engagement (current year balances)
- Prior year `TrialBalanceAccount` records (if rollforward)
- Confirmed account mappings (FS line item classifications)
- Engagement-level parameters: materiality thresholds, industry context
- Configurable sensitivity settings per engagement (default: flag variances >10% from prior year)

**Model:** Claude Haiku (pattern detection and classification; batch processing; nightly runs across all active engagements).

**Process:** Compute period-over-period variance by account → compute financial ratios (current ratio, quick ratio, debt-to-equity) → identify accounts with unusual activity relative to prior period → flag accounts where balance direction conflicts with account type (e.g., credit balance in a debit-normal account) → apply configurable sensitivity thresholds → generate anomaly descriptions.

**Output:** Anomaly flag indicators on individual accounts in the Sheets UI. Each flag includes: anomaly type, magnitude, and a brief description (e.g., "Accounts Receivable increased 47% year-over-year; prior year increase was 3%"). Flags are informational in the UI — they do not change any account status or audit content.

**Tier behavior:** See [Section 3](#pcaob-tier-elevation-for-trial-balance-anomaly-detection). For nonissuer engagements, anomaly flags are Tier 1 (informational, no AIDecision record). For PCAOB engagements, each flag creates an `AIDecision` record with `context_type = AnomalyDetection` and the auditor must document their investigation response.

**Human review gate:**
- Nonissuer: None required. Flags are informational; the auditor investigates at their professional discretion.
- PCAOB: Required. Each anomaly flag must be acknowledged (investigated, dismissed with reason, or noted for further procedures) and the response recorded in the `AIDecision` record.

**Failure modes:**

| Failure | Handling |
|---|---|
| No prior year data (new client) | Limited to within-period anomalies (balance direction, zero balances in expected accounts); year-over-year analysis unavailable |
| Materiality thresholds not set | Use default thresholds; prompt the auditor to set engagement-specific thresholds |
| Too many flags (noisy) | Apply materiality filter: suppress flags below clearly trivial threshold |
| API error / timeout | Nightly job retries on next cycle; no blocking dependency on any workflow |

---

### Feature 8: Report Section Draft

**Purpose:** Generate first drafts of standardized report sections (Description of Tests of Controls, Scope and Approach, System Description for SOC reports) so the partner writes opinion and judgment sections rather than boilerplate.

**Journeys:** 9 (partner requests AI draft of specific report sections)

**Trigger:** Partner explicitly requests a draft of a specific report section in the report editor (not automatic).

**Inputs:**
- Report type and template (SOC 2 Type II, Financial Audit Opinion, etc.)
- Engagement-wide data: all `Control` records with status, all `TestProcedure` results, all exceptions, `EngagementFramework` details
- `EvidenceItem` summary statistics (count by type, coverage by control area)
- Prior year report (if rollforward — style and structure reference)
- Firm report template for this report type

**Model:** Claude Sonnet (report content becomes part of the issued deliverable; quality, accuracy, and professional tone are critical).

**Process:** Aggregate engagement data by report section → populate template sections with structured data (control counts, exception summaries, testing methodology descriptions) → generate narrative for standardized sections → reference framework-specific language (AICPA, PCAOB, SSAE 18) → produce draft with engagement-specific content rather than generic boilerplate.

**Output:** Draft text inserted into the specific report section in the editor. `ReportVersion` record created. The draft section is labeled "AI Draft — requires review" with the same AI content tracking metadata as workpapers (see [Section 5](#5-ai-content-tracking)). `AIDecision` record created with `context_type = ReportSectionDraft`.

**Human review gate:** Mandatory. Same rules as Feature 4 (workpaper drafts): the report section cannot be finalized if no substantive human edits have been made to AI-generated content. The partner must actively modify and sign off. Report issuance (the `Report.status = Issued` transition) validates that all AI-drafted sections have been substantively edited.

**Sections AI may draft:**

| Section | Report Types | Notes |
|---|---|---|
| Description of Tests of Controls | SOC 1, SOC 2 | Populated from TestProcedure records and results |
| Scope and Approach | All | Populated from engagement configuration and framework details |
| System Description summary | SOC 1, SOC 2 | Populated from client-provided system description evidence |
| Control Environment Overview | Financial Audit | Populated from control testing results |

**Sections AI does NOT draft:**

| Section | Reason |
|---|---|
| Audit Opinion | Professional judgment; Tier 3 — human-initiated only |
| Management Assertions | Client-authored; not auditor content |
| Going Concern paragraphs | Professional judgment with legal implications |
| Emphasis of Matter / Other Matter | Professional judgment |
| Qualification language | Professional judgment with regulatory consequences |

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

The original `is_ai_draft` boolean on `WorkpaperVersion` is insufficient. It answers "has any human edit been made?" but cannot answer "how much of the AI content was substantively reviewed?" — a question that EQR reviewers (Journey 10), managers (Journey 6), and partners need answered to evaluate work quality.

A binary flag creates a perverse incentive: an auditor can change a single comma, clear the gate, and advance a workpaper that is 99% unreviewed AI output. The EQR reviewer has no visibility into this.

### Design

AI content tracking operates at the **section level** within workpapers and reports. Every document type that supports AI drafting (workpapers via Feature 4, reports via Feature 8) tracks AI origin per section.

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

---

## 6. Per-Engagement AI Cost Estimate

Cost estimates use Bedrock on-demand pricing. Prompt caching (0.1x base input rate for cached reads) applies to system prompts reused within an engagement — system prompt, framework text, and firm methodology context are cached across calls within the same engagement session.

### SOC 2 Type II Engagement

| AI Feature | Volume | Model | Approx. Cost |
|---|---|---|---|
| Document completeness review (Feature 1) | 200 docs | Sonnet | ~$4.00 |
| Control mapping (Feature 2) | 50 controls | Haiku (batch) | ~$0.11 |
| Evidence link suggestions (Feature 5) | 200 links | Haiku | ~$0.15 |
| Risk category suggestion (Feature 6) | 1 engagement | Sonnet | ~$0.05 |
| Workpaper narrative drafts (Feature 4) | 20 drafts | Sonnet | ~$1.20 |
| Report section drafts (Feature 8) | 4 sections | Sonnet | ~$0.25 |
| RAG retrieval overhead | ~50K tokens | Haiku | ~$0.04 |
| **Total per SOC 2 engagement** | | | **~$5.80** |

### Financial Audit Engagement

| AI Feature | Volume | Model | Approx. Cost |
|---|---|---|---|
| Document completeness review (Feature 1) | 150 docs | Sonnet | ~$3.00 |
| Control mapping (Feature 2) | 30 controls | Haiku (batch) | ~$0.07 |
| TB account mapping (Feature 3) | 200 accounts | Haiku (batch) | ~$0.15 |
| Anomaly detection (Feature 7) | ~5 runs | Haiku (batch) | ~$0.10 |
| Evidence link suggestions (Feature 5) | 150 links | Haiku | ~$0.12 |
| Risk category suggestion (Feature 6) | 1 engagement | Sonnet | ~$0.05 |
| Workpaper narrative drafts (Feature 4) | 40 drafts | Sonnet | ~$2.40 |
| Report section drafts (Feature 8) | 3 sections | Sonnet | ~$0.20 |
| RAG retrieval overhead | ~80K tokens | Haiku | ~$0.06 |
| **Total per financial audit** | | | **~$6.15** |

### With Prompt Caching

With prompt caching (system prompts and framework context reused within an engagement): effective cost **$3–5 per SOC 2 engagement**, **$4–6 per financial audit**. At 100 engagements/year for a mid-market firm, platform AI costs are **$300–$600/year**. This is absorbed into the subscription — AI consumption is not a separate line item on customer invoices at launch.

### Cost Change from Prior Estimate

The prior estimate ($0.75–$2.50 per engagement) was based on Haiku for document completeness review and did not include Features 5–8. Document completeness review is the highest-cost feature because it uses Sonnet at high volume — this is justified because errors in completeness assessment directly cause the PBC back-and-forth cycle that is the platform's primary value proposition to eliminate. Even at $5–6 per engagement, AI costs represent <0.05% of subscription revenue per engagement ($250–$350 overage rate). The margin impact is negligible.

---

## 7. What AI Does NOT Do at Launch

- Fine-tuned models per firm (insufficient training data volume at this stage)
- On-device / local model execution
- Multi-agent orchestration (LangGraph, CrewAI) — single-step AI calls are sufficient; frameworks add debugging overhead with no benefit for the defined use cases
- AI-generated audit opinions, professional conclusions, or materiality determinations
- Autonomous multi-step actions without human review gates at each step
- Auto-finalization of any audit content without a named human reviewer
- Real-time client-facing AI pre-check on upload (the completeness review runs asynchronously; the auditor reviews before the client receives feedback — adding a client-facing instant pre-check is a post-launch evaluation item that requires careful UX design to avoid creating false confidence in incomplete documents)

---

*End of Axiom AI Architecture Design*
