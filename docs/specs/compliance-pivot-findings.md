# Findings Brief: Axiom Pivot to Compliance/Assurance

**Status:** Research complete — awaiting product owner decisions before spec cascade

---

## TL;DR

The pivot is **strategically sound and timely**. SOC 2 is a ~$850M market with strong tailwinds; ISO 42001 has a severe auditor supply shortage making it a prime wedge; Delve — the most aggressive "agentic AI" competitor — was hit with a major fraud scandal three weeks ago (March 2026), creating a provenance/trust opening. Cross-mapping is feasible, legally supported, and market-validated (Drata, Vanta, Secureframe, Hyperproof, AuditBoard all do it). The spec needs a substantive rewrite but most structural bones carry over. Roughly **30% delete, 20% rewrite, 50% repurpose**.

---

## 1. Research highlights

### Frameworks

| Framework | Market signal | Overlap w/ SOC 2 | Notes |
|---|---|---|---|
| **SOC 2** | Revenue engine (~$850M, 83% of enterprise buyers require it) | — | Auditor = CPA firm. Sample-based. 5 TSCs, 100–150+ controls per report. |
| **ISO 27001** | Steady growth, paired with 27701 for GDPR | 60–95% | Accredited CB (not CPA). 93 Annex A controls. 3-year cert cycle + annual surveillance. |
| **ISO 27701** | Market pair to 27001 (privacy) | High via 27001 | **"ISO 2717" is almost certainly this.** Confirm below. GDPR/state-privacy-law driven. |
| **ISO 27017** | Cloud-specific; hyperscaler-focused | Moderate | Less common mid-market. Likely NOT the typo. |
| **ISO 42001** | **Steepest demand curve** — 76% plan audit in 24mo; severe auditor shortage | Partial | Strongest differentiator. Immature tooling. EU AI Act alignment. |
| **HIPAA** | No single certifier; HITRUST CSF r2 is de facto market standard | High (technical) | Multiple assessment modes required. |
| **PCI DSS 4.0.1** | Large remediation wave 2025–2027 | Moderate | Prescriptive; QSA-only; quarterly ASV scans + pen tests. |
| **SOC 1** | Slower growth; service-org focused | Low | Financial-process oriented. Keep optional. |

**Universal workflow backbone** all seven frameworks share: scoping → control universe → RFI/document requests → evidence collection → control testing → findings/CAPA → report generation. Axiom's existing engagement model generalizes cleanly.

**Key divergences:** PCI is uniquely prescriptive (population-level scans); HIPAA has no single certifier; ISO 42001 has novel artifacts (AI system inventory, model cards) and immature auditor ecosystem.

### Competitors

- **Agentive (YC S23, ~$500K funding, ~6 people):** auditor-side workspace; narrow agentic surface (per-request workpaper drafting); financial/SOX tilt; NO framework-aware intelligence. Small, capital-constrained, but genuinely good at PBC request UX. Named customers: Larson & Company, Stambaugh Ness, Cropper Accountancy.
- **Delve (Series A $32M, Insight-led, auditee-side):** most aggressive "agentic" marketing (browser agents, continuous monitoring, 25+ frameworks claimed). **March 2026: TechCrunch/Inc. reported 493 of 494 SOC 2 reports used identical boilerplate including a shared typo; YC reportedly asked them to leave.** Trust crisis is live. Integration depth (~100) thinner than Vanta (400+). No auditor-side product.

**Differentiation opportunities identified (8):**
1. Cross-framework evidence mapping (the wedge)
2. Auditor-side copilot as first-class surface
3. Autonomous population-level control testing (not just sampling)
4. Agentic management-response drafting with round-tripping to Jira/Linear
5. Continuous assurance with drift-triggered re-testing
6. AI findings triage with severity reasoning
7. Deep tool-use for non-API evidence with cryptographic provenance
8. Human-in-the-loop compliance ledger (ISO 42001-native AIDecision table)

### Cross-mapping feasibility

- **Feasible and market-standard.** Every serious player has a common-control catalog with 1:N framework mapping (Drata DCF, Vanta VCF, Secureframe, AuditBoard via UCF, Hyperproof via SCF, Sprinto Magic Mapping, Thoropass, OneTrust/Tugboat).
- **License, don't build.** Recommended stack:
  - **SCF** (free, CC-licensed, NIST STRM-encoded, quarterly updates, 1,400+ controls, 200+ frameworks) as primary crosswalk
  - **OSCAL** for NIST-family catalogs (future-proofs for FedRAMP)
  - **AICPA official mappings** for SOC-2-to-anything (auditor-defensible)
  - **CIS Controls v8.1 mappings** as secondary cross-check
  - **UCF** as optional enterprise SKU (commercial license)
- **Legally supported:** AT-C 105 and ISAE 3000 (Revised) permit *artifact* reuse; *opinions* are not transferable. ISO 17021-1 similar.
- **Hard problems are real:**
  - **Semantic mismatch** — SOC 2 (principle-based) vs. ISO Annex A (catalog) vs. PCI (prescriptive) vs. ISO 42001 (management-system)
  - **Period coverage** — Type 2 window ≠ ISO surveillance ≠ PCI 90-day scan validity
  - **Partial satisfaction** — evidence for SOC 2 CC6.1 may only partially cover ISO A.5.16/A.5.17/A.5.18
  - **Version churn** — PCI 3.2→4.0, ISO 27001:2013→2022, NIST CSF 1.1→2.0 all invalidate mappings at control-ID level
  - **Evidence staleness** — framework-and-control-specific age tolerances (ASV scan 90d, pen test 1y, background check 1y, etc.)
- **Data model recommendation:** control-centric directed labeled graph with effective-dated edges. Node types: `Framework`, `FrameworkVersion`, `FrameworkRequirement`, `CommonControl`, `EvidenceItem`, `Engagement`, `Test`, `Sample`. Edge types carry NIST STRM relationship vocabulary (`equivalent-to | subset-of | superset-of | intersects-with | no-relationship`) and strength score. Fits PostgreSQL + pgvector — no Neo4j required.
- **AI's role:** NOT authoring the authoritative crosswalk (auditors reject LLM-generated equivalence). AI suggests evidence→control mappings, does gap analysis, assists framework version migration, powers semantic search over requirement text.

---

## 2. Spec impact analysis

### What must be DELETED

- **Name rationale** — "anchored in the trial balance" (§On the name). Name still works; only the rationale changes.
- **ICP: "mixed financial audit + compliance framework audit"** (§1) → compliance/assurance-only
- **§4 regulatory table columns:** Financial Audit (Private, PCAOB), SQMS 1, AS 1220 EQR, 60-day workpaper assembly, PCAOB AS 1105 AI documentation
- **§5 TrialBalance data model** (entity 5 in the table: TrialBalance, TrialBalanceAccount, TrialBalanceAdjustment)
- **§6 AI features:** Feature 3 (Trial balance account mapping), Feature 7 (Trial balance anomaly detection)
- **§7 tech stack:** `internal/trialbalance` module
- **§8 integrations:** QBO/Xero/NetSuite/Sage direct connectors and Codat (all trial-balance-oriented)
- **Journey 4 entirely** — trial balance + population analytics
- **§12 innovative flow:** "Full-population analytics as alternative to sampling" (as framed — see repurpose below)
- **Pricing feature lines:** trial balance import, sampling calculators, materiality, GAAS (§3 Growth tier)
- **Competitive positioning vs. Fieldguide / CaseWare / DataSnipper** (§2) — wrong competitor set entirely

### What must be REWRITTEN

- **§1 Product Overview** — positioning, ICP, differentiators (now vs. Drata/Vanta/Secureframe/Delve/Agentive, not Fieldguide)
- **§2 Competitive positioning** — new competitor set; lean into provenance/auditor-defensibility post-Delve scandal
- **§3 Pricing tier features** — replace financial-audit bullets with framework-specific ones
- **§4 Regulatory requirements** — add ISO CB accreditation, PCI QSA, HITRUST r2, ISO 42001 surveillance cycles; remove PCAOB/GAAS
- **§5 Core Data Model** — introduce `CommonControl`, STRM-encoded `satisfies` edges, period-aware `EvidenceItem → CommonControl` edges, HITRUST r2 maturity scoring
- **§6 AI features** — replace trial balance features with: evidence→control mapping, gap analysis, framework version migration assistance, agentic management-response drafting, findings triage, drift-triggered re-testing
- **§11 user journeys** — journey reframes (Journey 4 removed; new journey for multi-framework integrated engagement)
- **§12 innovative flows** — cross-framework evidence mapping moves from "Scale tier feature" to **primary differentiator**; add provenance-verified evidence capture, continuous assurance with drift-triggered re-testing

### What CARRIES OVER (structurally intact)

- Modular monolith architecture
- Engagement / Control / Evidence / Document Request / AIDecision data model core
- Client Hub (PBC portal)
- Workpapers (now "evidence binders" / "audit files")
- WORM archiving (S3 Object Lock) — still required
- Tech stack (Go + React + Postgres + Bedrock + ECS)
- Security & compliance posture (SOC 2 + ISO 27001 + ISO 42001 targets — Axiom now dogfoods these)
- RBAC + RLS + audit log + AIDecision HITL pattern
- Self-serve onboarding / trial strategy
- Pricing *tier structure* (just feature re-mix)

---

## 3. Recommended direction

1. **Sequencing:** SOC 2 → ISO 27001 + 27701 → ISO 42001 → HIPAA (with HITRUST r2 path) → PCI DSS → SOC 1 (optional, service-org only). ISO 42001 gets outsized marketing emphasis despite smaller near-term revenue — it's the market wedge and Axiom credibly dogfoods it.
2. **Both-sided product.** Delve is auditee-only; Agentive is auditor-only. Neither competitor has both. Axiom's original spec is auditor-side; extend Client Hub into a full auditee workspace (continuous monitoring, evidence freshness, policy library). This is the big bet.
3. **Provenance as a category.** Post-Delve scandal: build cryptographic evidence provenance (signed screenshots, hashed DOM snapshots, WORM artifacts, full AIDecision ledger) and market it explicitly. "Auditor-defensible by construction." Defensible wedge that's live in the news *right now*.
4. **Cross-mapping stack:** license SCF (free) + OSCAL + AICPA official mappings + CIS. Do NOT build from scratch. Do NOT pay for UCF at launch (offer as enterprise SKU only).
5. **Control-centric data model** with STRM relationship vocabulary, effective-dated edges, partial-satisfaction tracking. Never show a green checkmark when coverage is partial — always show percentage and gap list.
6. **Kill Journey 4 entirely**, but repurpose "full-population analytics" engineering muscle as **agentic autonomous control testing** (log-stream checks against SOC 2 CC6, PCI 10 logging requirements, etc.). Same capability, different packaging.

---

## 4. Open decisions (need your input)

Before cascading changes through the spec:

1. **ISO 2717 typo:** Confirm **27701 (privacy)** — market-standard pair with 27001 — or 27017 (cloud)? Recommendation: 27701 primary, 27017 as secondary cloud-vendor feature.
2. **Both-sided vs. auditor-only:** Extend Client Hub into a full auditee GRC workspace (compete with Drata/Vanta/Delve directly), or stay auditor-centric (compete with Agentive)? Recommendation: both-sided — biggest differentiator, and cross-mapping value explodes with both sides.
3. **ICP rethink:** Keep CPA firms (20–60 staff), or also target (a) ISO-accredited Certification Bodies, (b) QSA firms, (c) compliance-first consultancies (e.g., Schellman, A-LIGN tier down)? Recommendation: CPA firms + compliance consultancies; skip CBs/QSAs at launch (different sales motions).
4. **Product name:** Keep "Axiom"? The trial-balance rationale is dead, but the name still fits (axiom = self-evident truth / foundational assurance). Alternatives worth considering if you want a reset.
5. **Financial audit bones:** Delete cleanly, or shelve as a dormant future module (e.g., partner-request)? Recommendation: delete cleanly. Every line of financial-audit code is technical debt pointing at a market we're leaving.
6. **Sequencing confirmation:** SOC 2 first is obvious. Do you want ISO 42001 emphasized as the *wedge* (higher marketing prominence than revenue alone would justify), or follow pure revenue ranking?
7. **Delve scandal positioning:** How aggressive do you want to be? Options: (a) subtle — market provenance/AIDecision as value prop without naming Delve, (b) direct — explicit comparison content ("why your AI compliance tool needs a ledger"), (c) stay above the fray. Recommendation: (a). Market will connect dots.
8. **HITRUST:** Support HITRUST CSF r2 as an assessment mode within the HIPAA offering? It's the de facto market standard but requires HITRUST Authorized External Assessor partnership. Recommendation: yes, but post-MVP.

---

## 5. Cascade plan (once decisions are resolved)

Per CLAUDE.md document dependency rules, changes will propagate in this order:

1. **User Journeys** (`docs/user-journeys/all-journeys.md`) — remove financial-audit journeys, add multi-framework integrated engagement journey, add continuous-assurance journey (if auditee-side is approved)
2. **Domain & Data Model** (`docs/specs/domain-and-data-model-design.md`) — remove TrialBalance entities, add CommonControl/FrameworkRequirement/STRM-edge entities, update traceability matrix
3. **AI Architecture** (`docs/specs/ai-architecture-design.md`) — rework 8 AI features around compliance-only surface, update AIDecision coverage
4. **Backend Architecture** (`docs/specs/backend-architecture-design.md`) — remove `internal/trialbalance`, add `internal/controlmapping` (or similar), update River workers
5. **OpenAPI specs** (`packages/openapi/*.yaml`) — endpoint redistribution
6. **Infrastructure Design** (`docs/specs/infrastructure-design.md`) — minimal changes; possibly new IAM for provenance/signing
7. **Infrastructure diagram** (`docs/diagrams/axiom_infrastructure.py` → regenerate `.png`)
8. **Product Spec summary hub** (`docs/specs/axiom-spec-design.md`) — updated last as aggregation of all above
9. **Mockups** (`mockups/journey-*/`) — reviewed for which screens survive, which need replacement, which need new
10. **Design System** (`.impeccable.md`) — likely unaffected

---

## Appendix: source material

The three full research outputs (compliance frameworks, Agentive/Delve competitive analysis, cross-mapping feasibility) were produced by parallel research agents on 2026-04-17. Key sources cited in each research output include AICPA, ISO, HHS, PCI SSC, NIST OSCAL/CSRC, TechCrunch/Inc. (Delve scandal), Y Combinator, Drata/Vanta/Secureframe/AuditBoard/Hyperproof product documentation, and Secure Controls Framework.
