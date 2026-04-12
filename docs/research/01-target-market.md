# Research Task 1: Target Market and ICP Definition

## Market Segments Overview

The audit software market has fragmented into four distinct buyer segments with very different needs and incumbents:

| Segment | Firm Size | Current Tools | Incumbent Software | Budget |
|---|---|---|---|---|
| **Big Four / Top 500** | 500+ staff | Enterprise platforms | Fieldguide, AuditBoard, Workiva | $150K+/year |
| **Mid-market CPA/advisory** | 20–200 staff | Fragmented stack | CaseWare + DataSnipper + Excel | $15–60K/year |
| **Boutique compliance specialists** | 1–30 staff | Lightweight SaaS | Yak, Drata, Vanta | $5–15K/year |
| **In-house compliance teams** | N/A (corporate) | Continuous monitoring | Drata, Vanta, Sprinto | $10–80K/year |

The fourth segment (in-house teams) is not the auditor market — it is the auditee market. Drata and Vanta serve companies trying to pass audits, not firms performing them. This is a separate competitive space and should not be a target.

---

## Segment Analysis

### Big Four / Top 500 — Do Not Target at Launch

Fieldguide dominates here. It is built for this segment: its enterprise positioning, Implementation Consultant onboarding model, and broad framework coverage (Financial Audit, SOC, ISO, PCI, HIPAA, ESG) map directly to the Big Four workflow. AuditBoard and Workiva compete here at the $40K–$150K/year price point.

**Why Axiom cannot win here at launch:** The sales cycle is 6–18 months, procurement requires extensive security reviews (BAA, SOC 2 evidence, pen test reports), and differentiation requires deep methodology customization and enterprise integrations. This segment is also fiercely loyal once switched — implementation cost is a switching barrier in both directions.

### Boutique Compliance Specialists — Addressable but Narrow

Yak serves this segment well: per-engagement pricing at ~$800/year, unlimited users, fast self-serve setup, focused on SOC 1/2 and HIPAA. Drata and Vanta compete from the auditee side, making this market increasingly crowded on both sides.

**Why this segment is limited:** Boutique compliance-only firms doing SOC 2 engagements have low revenue per engagement and are price-sensitive. Yak's positioning is defensible for this narrow use case. A platform trying to serve this segment would compete on price with Yak and would need to win on speed or AI quality alone — neither of which is a sustainable moat.

### Mid-Market CPA / Advisory Firms — Primary ICP

This is the largest underserved segment. Firms in the 20–200 staff range doing a mix of financial audits, compliance audits (SOC 2, ISO 27001, HIPAA), and advisory work are currently running a fragmented, expensive toolstack:

- **CaseWare Cloud** for audit file structure and workpaper management: ~€80–120/user/month
- **DataSnipper** for evidence extraction and matching within Excel: ~€500–1,500/user/year
- **Excel** for everything the other tools don't cover (trial balance, sampling, analysis)
- **Email and shared drives** for client document requests

The total cost of this stack for a 14-person fieldwork team is approximately €14,000–18,000/year (DataSnipper alone), plus CaseWare licensing on top. Despite this cost, critical gaps remain.

---

## Pain Points of Mid-Market Firms (Evidenced)

From survey data (CPA Practice Advisor, January 2025, 41 US-based auditors):

- **70%** of auditors spend over half their time in spreadsheets
- **49%** cite time-consuming reconciliations as their biggest pain point (5–20 hours/week)
- **54%** spend 5–10 hours/week on manual PDF data extraction
- **42%** report manual data entry errors as a significant problem
- **~40%** experience review delays exceeding two weeks

From market research on non-Big 4 audit tools (ciferi, 2026):

- "Evidence matching and cross-referencing accounts for 32% of fieldwork hours"
- Expensive platforms handle file structure but "frequently lack built-in calculators for ISA 320 materiality, ISA 530 sampling, ISA 570 going concern"
- "Feature bloat and cost misalignment" — firms paying for full platforms when they need one layer

From Fieldguide alternative-seekers (AuditCue research):

- Fieldguide's pricing is "unsustainable" for mid-market firms who don't see matching productivity gains
- "Bloated risk assessment and unintuitive interface"
- Steep learning curve requiring Implementation Consultant involvement
- AuditBoard notably **lacks a built-in trial balance module**, forcing Excel workarounds for financial audits

---

## Competitive White Space

The mid-market is genuinely underserved because:

1. **Fieldguide is priced and positioned for enterprise.** Its opaque, quote-based pricing, Implementation Consultant-led onboarding, and deep feature set are correctly matched to the Big Four but create friction for a 40-person firm.

2. **CaseWare + DataSnipper is the de facto mid-market standard but is not integrated.** These are two separate tools with separate licenses, separate logins, and no unified AI layer. CaseWare is strong at file/workpaper management but is legacy architecture. DataSnipper is strong at evidence extraction but lives inside Excel, meaning the "spreadsheet problem" is still the foundation.

3. **Yak is too narrow.** It serves compliance-only firms doing SOC 1/2/HIPAA. It does not support financial audit (no trial balance), and its scope is explicitly limited to a small subset of frameworks.

4. **No competitor has unified financial audit and compliance audit in a mid-market-priced, AI-native platform.** This is the gap.

---

## Recommended ICP

**Mid-market accounting and advisory firms, 20–200 staff, performing a mix of financial audits and compliance framework audits (SOC 2, ISO 27001, HIPAA), operating in the US and Canada, who are currently running CaseWare and/or DataSnipper alongside Excel and are experiencing visible bottlenecks in evidence collection and review turnaround.**

More specifically, the most reachable early adopter within this segment is:

- A **partner-led firm of 20–60 people** doing 30–100 engagements per year
- Already aware of Fieldguide or Yak but priced out of Fieldguide and under-served by Yak's compliance-only scope
- Frustrated specifically by: PDF evidence extraction time, review delays, client document request back-and-forth, and the Excel reconciliation burden
- Willing to pay $10,000–$40,000/year for a platform that meaningfully replaces 2–3 tools in their current stack

---

## Segments Explicitly Out of Scope at Launch

| Segment | Reason |
|---|---|
| Big Four / Top 500 | Sales cycle too long, security requirements too heavy, Fieldguide entrenched |
| In-house compliance teams | Different market (auditee, not auditor) — Drata/Vanta territory |
| Solo practitioners / micro-firms (<10 staff) | Price sensitivity too high, ROI unclear below critical mass of engagements |
| Internal audit departments (corporate) | Different workflows, TeamMate/AuditBoard territory, different buying center |
| International-only firms (EU/APAC, no US presence) | Regulatory complexity (GDPR, PCAOB jurisdictions) best deferred until US is proven |

---

## Sales Cycle Implications

For the recommended ICP:

- **Sales cycle**: 4–12 weeks. Typically a partner-level decision, sometimes with IT and compliance sign-off. Not a 6-month enterprise procurement process.
- **Key trigger event**: Busy season pain (January–April for US calendar-year clients), a failed engagement review, or losing a client because the document request process was too slow.
- **Champion profile**: A managing partner or audit manager who has already evaluated Fieldguide and was deterred by cost or implementation burden.
- **Buying criteria**: Time-to-first-engagement (can we run a real engagement within 2 weeks of signing?), replacement of at least one existing tool (usually DataSnipper or CaseWare), demonstrable AI reduction in manual work, transparent pricing.
- **Self-serve feasibility**: High. A mid-market firm doing SOC 2 or financial audit should be able to onboard without an Implementation Consultant, using pre-built methodology templates. This is a structural advantage over Fieldguide.

---

## Sources Consulted

- Fieldguide customer reviews: Capterra, G2, Gartner Peer Insights
- CPA Practice Advisor: "Many Auditors Are Still Married to Their Excel Spreadsheets" (Jan 2025)
- ciferi: "Best Audit Tools for Non-Big 4 Firms in 2026"
- AuditCue: "Alternatives to Fieldguide.io"
- Market research: GMInsights, Fortune Business Insights (audit software market 2025–2032)
- Competitor analysis: AuditBoard (Sprinto), Vanta (Drata), Workiva, TeamMate, Suralink, CaseWare
