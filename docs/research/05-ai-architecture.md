# Research Task 5: AI Architecture and LLM Strategy

## LLM Provider Analysis

### Data Privacy Comparison

The table below covers the four viable providers for handling sensitive financial audit data (PII, financial statements, control documentation, PHI for HIPAA engagements).

| Criteria | Anthropic (Direct API) | OpenAI (Direct API) | AWS Bedrock | Azure OpenAI |
|---|---|---|---|---|
| **Training on customer data** | No — commercial API explicitly excluded | No — API excluded by default | No | No |
| **Default log retention** | 7 days (reduced Sept 2025), then deleted | 30 days | Not retained by provider — stays in your AWS account | Not retained by provider — stays in your Azure tenant |
| **Zero Data Retention option** | Yes (ZDR mode) | Yes (ZDR endpoints) | N/A — data never leaves your VPC | N/A — data stays in tenant |
| **HIPAA BAA** | Yes (qualifying enterprise customers) | Yes (email baa@openai.com) | Yes (HIPAA eligible) | Yes (HIPAA BAA) |
| **SOC 2** | Yes (Type II) | Yes (Type II) | Yes | Yes |
| **Data residency** | US-based; no region pinning at API tier | US, EU, UK, JP, CA, AU and more | Any AWS region you deploy to | Any Azure region; Resource Location Pinning |
| **Network isolation** | No (shared infrastructure) | No (shared infrastructure) | Yes — runs in your VPC, data never leaves your AWS account | Yes — VNET integration, private endpoints |
| **Models available** | Claude (Haiku, Sonnet, Opus) | GPT-4o, o3, o1, etc. | Claude + Llama + Titan + others | GPT-4o, o1, o3 (no Claude) |
| **Operational complexity** | Low — single API | Low — single API | High — AWS infra required | Medium — Azure tenant setup |

### Recommendation: Anthropic Claude API (Direct) as Primary Provider

**Rationale:**

1. **Data privacy is sufficient for the target ICP.** Mid-market US CPA firms doing private company financial audits and SOC 2 work do not require in-VPC model execution. The Anthropic API's commercial data policy (no training, 7-day deletion, ZDR available) satisfies the data handling obligations of the AICPA and PCAOB engagement types in scope. The HIPAA BAA covers HIPAA audit engagements.

2. **Claude's long-context capability is the right fit for audit work.** Document review, workpaper analysis, and evidence completeness checks involve reading lengthy documents with multiple pages. Claude's 200K+ token context window and strong performance on document comprehension tasks (legal/financial reading) makes it the best-fit model family for the primary use cases.

3. **Prompt caching eliminates redundant cost.** System prompts (firm methodology, framework criteria, engagement context) are reused across many calls within an engagement. Cached reads cost 0.1× the base input rate. For a typical engagement, this reduces AI cost by 60–80%.

4. **Haiku + Sonnet tiering allows cost optimization.** High-volume, lower-complexity tasks (document classification, account categorization, request reminders) use Haiku ($1/$5 per million tokens). Quality-critical tasks (evidence completeness reasoning, workpaper drafting, control mapping) use Sonnet ($3/$15 per million tokens). Opus is reserved for the most complex multi-step reasoning tasks if needed.

5. **Batch API for non-real-time operations.** Evidence processing, nightly completeness reviews, and background control mapping runs use the Batch API at 50% discount (Sonnet: $1.50/$7.50 per million tokens).

**AWS Bedrock as a future enterprise tier:** For customers requiring maximum data isolation (data never leaves their AWS region or VPC), offer an "Enterprise" deployment option backed by AWS Bedrock with Claude models via the Anthropic/AWS partnership. This is a future option, not an MVP requirement. Bedrock can run the same Claude models; the application code changes minimally.

**Do not use OpenAI as primary.** Not for data privacy reasons (they are comparable to Anthropic), but because splitting across two model families creates maintenance overhead for prompt engineering and behavior consistency. A single model family is simpler to reason about and optimize. OpenAI can be added later if customers require it.

---

## Vector Database Recommendation: pgvector at Launch

**Decision: pgvector (PostgreSQL extension) at launch, with migration path to Qdrant.**

| Factor | pgvector | Qdrant | Pinecone |
|---|---|---|---|
| **Operational complexity** | Zero — already have Postgres | Medium — separate service | Low — managed SaaS |
| **Data residency** | In your own Postgres instance | Self-hosted option available | Third-party SaaS — data leaves your infra |
| **Vendor lock-in** | None | None (open source) | High |
| **Performance (< 5M vectors)** | Sufficient — sub-20ms at 1M vectors with HNSW index | Excellent | Excellent |
| **Cost** | Free (Postgres extension) | Free (self-hosted) | $70+/month minimum |
| **Suitable for audit data** | Yes — same security boundary as main DB | Yes — self-hosted | No — data residency concern |
| **Migration** | Embeddings are portable numeric vectors | — | — |

**What goes in the vector store:**
- Firm methodology documents (uploaded during onboarding)
- Framework requirement descriptions (FrameworkRequirement.description)
- ControlObjectiveLibrary descriptions
- Prior engagement evidence summaries (for "similar evidence used before" suggestions)
- Workpaper templates and examples

**Scale estimate for target ICP:**
- 50 firms × 20 methodology documents × ~50 chunks/doc = 50,000 vectors
- Framework requirements: ~2,000 vectors (all frameworks combined)
- Evidence summaries: 50 firms × 100 engagements × 200 evidence items = 1,000,000 vectors (at maturity)

This is well within pgvector's comfortable range. Migrate to Qdrant if/when approaching 5–10M vectors.

---

## AI Feature Architecture

### Feature 1: Document Completeness Review

**The highest-value AI feature.** Eliminates the primary bottleneck in the engagement: the back-and-forth cycle when clients submit incomplete or wrong documents.

**Trigger:** Client uploads a document in response to a DocumentRequest.

**Inputs:**
- DocumentRequest (title, description, instructions, linked Control, linked TestProcedure)
- Uploaded document (extracted text from OCR pipeline)
- ControlObjective description and FrameworkRequirement text
- RAG context: embeddings of prior accepted evidence for this type of request (from same firm or similar firms, anonymized)

**Process:**
```
1. Extract document attributes (date range, parties, system names, amounts, format)
2. Compare extracted attributes against DocumentRequest criteria
3. Check period coverage (critical for Type II: evidence must cover full audit period)
4. Score completeness against each test procedure requirement
5. Identify specific gaps: what is missing, what is wrong period, what is wrong format
6. Generate plain-language client explanation (for display in client portal)
7. Generate auditor-facing summary (for review workflow)
8. Produce recommendation: Accept | Request Clarification | Reject
```

**Model:** Sonnet (reasoning quality matters; errors here propagate into the audit file)

**Outputs → AIDecision record:**
```json
{
  "context_type": "DocumentCompleteness",
  "suggested_value": "ACCEPT — document covers full audit period, matches requested system access log format",
  "confidence": 0.87,
  "raw_output": { ... full model response ... },
  "review_action": "Pending"
}
```

**Human review gate:** Always required before Accept. The auditor sees the AI assessment and recommendation in the review queue. They can Accept (with one click if they agree), Modify (edit the assessment), or Reject (return to client with different explanation). The AIDecision record captures their choice.

**Failure modes and fallbacks:**
| Failure | Handling |
|---|---|
| Document encrypted / password-protected | Extraction fails → document flagged for manual review, client notified |
| Non-English document | Flag with detected language → manual review |
| AI confidence below 0.6 | Surfaced as "Low confidence — manual review recommended" without a definitive recommendation |
| Extraction succeeds but document is an image-only PDF | OCR confidence score returned; if below threshold, flag for manual |
| Model returns error / timeout | Request queued for retry; auditor notified if retry fails after 3 attempts |

---

### Feature 2: Control Mapping (Framework-Agnostic Evidence Linkage)

**The architectural differentiator.** Maps FirmControlObjectives to FrameworkRequirements across all frameworks in the engagement simultaneously.

**Trigger:** New engagement created from methodology template, or new FirmControlObjective added.

**Inputs:**
- FirmControlObjective name + description
- All FrameworkRequirements for frameworks included in the engagement (from DB, not from vector search)
- RAG context: ControlObjectiveLibrary entries with existing mappings (as few-shot examples of good mappings)

**Process:**
```
1. Embed the FirmControlObjective description
2. Retrieve top-k similar ControlObjectiveLibrary entries (pgvector similarity search)
3. For each framework in the engagement, retrieve candidate FrameworkRequirements
4. Score semantic similarity between FirmControlObjective and each FrameworkRequirement
5. Apply threshold filtering (only propose mappings above 0.75 similarity)
6. Generate explanation for each proposed mapping
7. Return proposed FirmControlObjectiveMapping records with confidence scores
```

**Model:** Haiku is sufficient for this structured classification task. Sonnet if confidence scores are consistently low.

**Outputs — presented to auditor as a reviewable mapping table:**
```
FirmControlObjective: "Access to production systems is restricted to authorized personnel"
  → SOC 2 CC6.1 (confidence: 0.94) — "Reason: both describe logical access controls..."
  → ISO 27001 A.8.3 (confidence: 0.91) — "Reason: ISO 27001 A.8.3 addresses access rights..."
  → HIPAA §164.312(a)(1) (confidence: 0.88) — "Reason: HIPAA technical safeguard for access..."
```

**Human review gate:** Auditor reviews proposed mappings in bulk (checkboxes). All confirmed; any rejected. No mapping takes effect until confirmed. AIDecision records created for each mapping.

**Failure modes:**
| Failure | Handling |
|---|---|
| Novel control with no library match | Low confidence → surface with note "no similar controls found in library; manual mapping required" |
| Ambiguous description matching multiple criteria | Surface all candidates with explanations; auditor chooses |
| Framework not in library | Graceful degradation — skip that framework, flag for manual |

---

### Feature 3: Trial Balance Account Mapping

**Trigger:** Trial balance imported (CSV/Excel upload).

**Inputs:** Account number, account name, account balance (debit/credit), prior year mapping (if rollforward).

**Process:** Classify each account into a standard financial statement line item (Cash, Accounts Receivable, Fixed Assets, etc.) using few-shot classification.

**Model:** Haiku (simple classification; high volume; speed matters for large trial balances).

**Human review gate:** Auditor reviews and confirms mappings in the Sheets interface. Anomalous mappings (low confidence, account name doesn't match expected category) highlighted for priority review.

---

### Feature 4: Workpaper Narrative Draft

**Trigger:** Auditor marks a TestProcedure as Complete and requests a draft.

**Inputs:**
- Control description
- TestProcedure description, procedure type, and expected result
- Evidence items linked (their extracted text and filenames)
- Exceptions noted (if any)
- Prior year workpaper narrative (if rollforward, as style reference)
- Firm workpaper template for this procedure type

**Process:** Generate a first-draft narrative that a staff auditor would write, matching firm style.

**Model:** Sonnet (quality matters; workpaper narratives go into the audit file).

**Human review gate:** Mandatory. AI draft is presented in the workpaper editor as editable text. The auditor must actively edit and sign off. The workpaper content is never finalized based solely on AI output. The draft is explicitly labeled "AI Draft — requires review" until modified.

**Critical constraint:** AI-generated workpaper content must be distinguishable from auditor-authored content in the version history (WorkpaperVersion records include `is_ai_draft: bool`). This satisfies emerging PCAOB documentation requirements.

---

## Human-in-the-Loop Policy

Three tiers based on audit risk and regulatory requirements:

### Tier 1 — Fully Automated (No Human Approval Required)

AI may take these actions without user interaction:

- Text extraction from uploaded documents (OCR + LLM)
- Embedding generation and vector store indexing
- Flagging documents as potentially incomplete in the UI (notification, not decision)
- Generating and sending overdue request reminders to clients
- Surfacing "this evidence was used in a prior engagement" suggestions
- Running anomaly detection on trial balance data (flagging outliers in UI)
- Computing trial balance analytics (variance analysis, ratio calculations)

These actions are either purely informational (not affecting audit conclusions) or are reversible notifications.

### Tier 2 — Human Approval Required Before Taking Effect

AI suggests; human confirms. No audit file content changes without explicit auditor action:

- Evidence links (EvidenceLink records)
- Control → FrameworkRequirement mappings (FirmControlObjectiveMapping records)
- Trial balance account category assignments (TrialBalanceAccount.mapped_fs_line_item)
- Document request fulfillment decisions (accept/reject uploaded document)
- Risk assessment pre-population (auditor must review and certify)
- Workpaper narrative drafts (always editable before sign-off)

### Tier 3 — Human-Initiated Only (AI May Assist But Never Initiates)

These actions can only be triggered by a human, regardless of AI capability:

- Engagement status transitions
- Control conclusions (pass/fail)
- Exception documentation
- Workpaper review and sign-off
- Engagement quality review sign-off
- Report issuance
- Client acceptance documentation
- Any action that constitutes professional judgment under AU-C or PCAOB standards

This policy maps directly to PCAOB AS 1105: AI outputs inform but do not constitute professional judgments.

---

## AI Decision Storage (Regulatory Compliance)

Every Tier 2 AI action creates an AIDecision record (defined in Task 4 data model). The key fields that satisfy regulatory requirements:

| PCAOB AS 1105 Requirement | AIDecision Field |
|---|---|
| What technology procedure was performed | `context_type` + `context_id` |
| What model/tool was used | `model_id` |
| What the AI determined | `suggested_value` + `raw_output` |
| Who reviewed the AI output | `reviewed_by_id` |
| What the auditor decided | `review_action` + `accepted_value` |
| When review occurred | `reviewed_at` |

For PCAOB engagements, AIDecision records are included in the engagement's immutable archive. They form the "AI audit trail" that inspectors can examine.

---

## Per-Engagement AI Cost Model

Based on current Anthropic pricing and typical mid-market engagement volumes:

| AI Operation | Volume (per engagement) | Model | Tokens | Cost |
|---|---|---|---|---|
| Document completeness review | 200 docs × ~2K tokens | Haiku (batch) | 400K | ~$0.30 |
| Control mapping | 50 controls × ~3K tokens | Haiku (batch) | 150K | ~$0.11 |
| Trial balance mapping | 200 accounts × ~500 tokens | Haiku | 100K | ~$0.10 |
| Workpaper drafts | 20 drafts × ~8K tokens | Sonnet | 160K | ~$1.20 |
| Evidence link suggestions | 200 × ~1K tokens | Haiku | 200K | ~$0.15 |
| RAG retrieval overhead | ~50K tokens/engagement | Haiku | 50K | ~$0.04 |
| **Total per SOC 2 engagement** | | | | **~$1.90** |
| **Total per financial audit engagement** | | | | **~$3–6** (larger trial balance, more workpapers) |

With prompt caching (system prompts reused across calls within an engagement): reduce by 60–70%.

**Effective AI cost per engagement: $0.75–$2.50.**

This is negligible relative to engagement fees. At 100 engagements/year for a mid-market firm, platform AI costs are $75–$250/year. Absorbing this into the platform subscription is straightforward. There is no need to expose AI consumption costs to customers as a separate line item at launch.

---

## AI Infrastructure Architecture

```
Request arrives at API layer
  │
  ├── Synchronous (user is waiting):
  │     Document completeness review → Claude API (ZDR mode) → AIDecision created
  │     Control mapping suggestions → Claude API → AIDecision created
  │
  └── Asynchronous (background job):
        Document text extraction → OCR service → Claude API (extract + summarize)
        Nightly completeness sweep → Batch API (50% cost reduction)
        Evidence embedding generation → Claude API → pgvector store
        Workpaper draft on request → Claude API → WorkpaperVersion (is_ai_draft: true)

All Claude API calls:
  - Authenticated with firm-scoped API key rotation
  - System prompt includes: firm name, engagement type, applicable standards
  - Response logged to AIDecision table before returning to caller
  - Retried up to 3 times on transient failure
  - Hard timeout: 30 seconds synchronous, unlimited asynchronous
```

**Document processing pipeline (for uploaded evidence):**
```
Client uploads document
  → S3 storage (server-side encryption, AES-256)
  → Extraction job queued (async)
  → PDF: pdfplumber + Claude vision for image-heavy pages
  → Word/Excel: python-docx / openpyxl text extraction
  → Extracted text stored on EvidenceItem.extracted_text
  → Embedding generated and stored in pgvector
  → DocumentCompleteness AIDecision created (batch queue)
  → Auditor notified: "AI has reviewed [document name]"
```

---

## What This Architecture Does Not Include at Launch

| Capability | Reason Deferred |
|---|---|
| Fine-tuned models per firm | High cost, requires training data volume not available at launch |
| On-device / local model execution | Not required for target ICP; adds infrastructure complexity |
| Multi-agent orchestration frameworks (LangGraph, CrewAI) | Single-step AI calls are sufficient for the defined features; frameworks add debugging overhead |
| AI-generated audit opinions or conclusions | Outside human-in-the-loop policy; may create liability |
| Autonomous AI agents that take multi-step actions without human review gates | Incompatible with Tier 2/3 policy; not currently compliant with PCAOB guidance |

---

## Sources Consulted

- Anthropic Privacy Center: API data retention and training policies
- OpenAI Enterprise Privacy page
- AWS Bedrock vs Azure OpenAI compliance comparison (CloudOptimo, Reintech, DEV Community)
- pgvector vs Qdrant vs Pinecone comparison (CloudMagazin, TigerData, Firecrawl)
- Anthropic API pricing page (platform.claude.com/docs/en/about-claude/pricing)
- PCAOB AS 1105 (effective Dec 15, 2025): technology-assisted analysis documentation requirements
