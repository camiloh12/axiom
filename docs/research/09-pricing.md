# Research Task 9: Pricing Strategy and Business Model

## Executive Summary

Axiom should adopt a **hybrid flat-base + per-engagement-overage model** organized into three tiers: a self-serve Growth tier (~$12,000–$18,000/year), a Scale tier (~$24,000–$36,000/year), and an Enterprise tier (custom). This model aligns price with value delivery, supports self-serve adoption at the entry tier, creates natural expansion revenue, and avoids the two failure modes of per-user pricing (penalizes breadth) and pure per-engagement pricing (unpredictable revenue). At 50 firms, Axiom can generate ~$900K ARR; at 200 firms, ~$4.2M; at 500 firms, ~$12M+.

---

## 1. Competitor Pricing Benchmarks

### Fieldguide

Fieldguide does not publish pricing. All available evidence points to a quote-based, enterprise sales motion with no self-serve option. Key facts established from research:

- Used by more than 40 of the Top 100 US CPA firms and claims Big Four presence; $700M valuation after $125M total funding (Series C closed February 2026, led by Goldman Sachs Growth Equity at $75M).
- Pricing is **per-engagement in philosophy** — their own blog states that per-seat pricing "disincentivizes efficiency" and that per-engagement aligns the vendor's incentive (firms grow without hiring more staff) with the customer's goal. This is a notable signal about how they position internally, even if the billing mechanics are not public.
- Estimated market pricing for mid-market firms: **$30,000–$100,000+/year** based on comparable enterprise audit platforms and patterns from analyst review sites (SoftwareAdvice, Capterra, GetApp). No public confirmation of actual contract values.
- The Accelerator onboarding program (multiple workshop phases, dedicated implementation consultant) is a structural cost that is likely baked into first-year pricing, creating a high minimum viable deal size.
- **The key competitive fact:** Fieldguide is opaque, requires a sales call, and has no self-serve path. Mid-market firms who have gone through their sales process and emerged without a quote are the most accessible Axiom prospects.

### AuditBoard

- Publicly reported via Vendr marketplace data: **$40,000–$150,000/year**; median contract ~$42,800/year.
- CrossComply Essentials starts at ~$32,831/year; SOXHUB Professional at ~$48,204/year.
- Positioned for internal audit, SOX, and enterprise GRC. Not a direct competitor for external financial audit. Included here for market context only.

### CaseWare

- **~€80–120/user/month** for CaseWare Cloud, consistent with prior research.
- Reference data point from ciferi: a firm called "Bakker Audit" (representative mid-market firm) paid approximately **€28,000/year** for CaseWare before DataSnipper; migrating to CaseWare Cloud added ~€6,000/year.
- One-time migration cost for a 22-person firm switching platforms: **€40,000–60,000** (including template rebuilds and file migration). This switching cost is a significant retention anchor for existing CaseWare customers — but also a reason firms resist switching to CaseWare once they have already left.
- **No AI layer.** Deep methodology library but desktop-first architecture. DataSnipper is a mandatory add-on for modern evidence workflows.

### DataSnipper

- Three tiers as of 2025–2026: **Start, Accelerate, Elevate** — all require custom quotes, but third-party estimates:
  - Start: ~$64/user/month (~$768/user/year)
  - Accelerate: ~$175/user/month (~$2,100/user/year)
  - Cloud Collaboration Suite add-on: ~$62/user/month additional
- For a team of 14 fieldwork staff on the Accelerate tier: **€14,000–18,000/year** (ciferi data, consistent with €500–1,500/user/year range in prior research).
- For a 20-person fieldwork team: estimated **€10,000–30,000/year** depending on tier.
- DataSnipper is Excel-native and does not replace an audit platform — it is a productivity layer. Firms pay for both CaseWare and DataSnipper simultaneously.

### Yak

- Per-engagement pricing at approximately **$800/year** (unlimited users per engagement).
- Compliance-only (SOC 1, SOC 2, HIPAA). No financial audit support.
- Attractive for boutique compliance-specialist boutiques but irrelevant to mixed-service CPA firms.
- Yak's pricing serves as a **market anchor**: it establishes that compliance-only audit tooling is perceived as a commodity worth ~$800/engagement-year, not $50,000/year.

### Total Cost of Ownership: 30-Person CPA Firm Using CaseWare + DataSnipper

| Cost Component | Annual Estimate |
|---|---|
| CaseWare Cloud (30 users @ ~€100/user/month) | ~€36,000/year |
| DataSnipper Accelerate (20 fieldwork staff @ ~€1,000/user/year) | ~€20,000/year |
| Excel + Microsoft 365 (already paid) | ~€3,600/year |
| Training, onboarding, implementation amortized | ~€3,000/year |
| **Total annual toolstack spend** | **~€62,000/year (~$67,000 at parity)** |

Note: the €36,000 CaseWare figure assumes all 30 staff are licensed. Many firms license only active audit staff, so real-world CaseWare costs may be lower (18–22 users), bringing the combined total to **$40,000–$55,000/year**. The $15,000–$35,000 estimate used in prior research likely reflects smaller firms (12–18 staff) or firms that subset-license CaseWare aggressively.

**The Axiom value proposition in numbers:** A 30-person firm spending $40,000–$55,000/year on fragmented tools should be willing to pay $15,000–$30,000/year for a unified replacement that eliminates DataSnipper entirely and reduces CaseWare dependency. At 50% of the current spend, there is no savings argument against it.

---

## 2. Pricing Model Options: Analysis and Trade-offs

### Option A: Per-User/Month Subscription

**Mechanics:** $X/user/month × number of licensed seats.

**What it incentivizes:** Firms license only the seats they need, minimizing cost. Axiom's revenue grows when firms add headcount. Firms are incentivized to minimize the number of licensed users, which reduces platform stickiness (fewer people in the product = less organizational dependency).

**Pros:**
- Simple to understand and budget for.
- Easy to forecast revenue (ARR = seats × monthly rate × 12).
- Familiar to buyers; no explanation required.

**Cons:**
- Penalizes firm growth. A firm that grows from 20 to 40 staff faces a cost doubling. This is the opposite of the value proposition.
- Encourages "seat rationing" — firms give access only to senior staff, preventing junior auditor adoption, reducing stickiness.
- CaseWare already uses this model and it is a known friction point. Positioning Axiom on per-user pricing directly against CaseWare puts Axiom on the same unfavorable terms.
- Fieldguide explicitly criticized per-seat pricing as misaligned with firm efficiency goals. A per-user model signals the same misalignment.
- Does not capture value from high-engagement firms (a 30-person firm doing 200 engagements/year pays the same as one doing 20 engagements/year).

**Verdict:** Weak fit for Axiom's ICP. The model is familiar but structurally misaligned with the efficiency value proposition.

### Option B: Per-Engagement Pricing

**Mechanics:** Flat fee per engagement started/completed on the platform.

**What it incentivizes:** Fast engagement setup (each engagement is a billable event worth optimizing). Firms that run more engagements pay more. Axiom's revenue correlates directly with platform usage.

**Pros:**
- Aligns cost with output — firms pay when they get value.
- Self-limiting risk: a firm doing fewer engagements (slower period) pays less, reducing churn pressure during downturns.
- Transparent and auditable (a firm knows exactly what they paid per engagement).
- Yak's model validates that CPA firms accept per-engagement billing.

**Cons:**
- Revenue is unpredictable month-to-month; engagement volume is seasonal (US calendar-year audits peak January–April, SOC 2 renewals cluster by client anniversary).
- Firms facing a high-volume season may experience sticker shock from an unexpectedly large bill.
- Incentivizes firms to avoid starting engagements on the platform ("we'll do this one in Excel to save the fee"), undermining adoption.
- Makes it hard to price for compliance engagements vs. financial audits (complexity varies enormously — a 3-control HIPAA review vs. a 200-control SOC 2 Type 2 should not cost the same).
- CAC is high relative to a small per-engagement fee if ACVs are low.

**Verdict:** Best for commodity compliance-only tools (Yak's domain). Not ideal as the primary model for a platform aspiring to $15,000–$30,000/year ACV because it either under-prices or creates billing volatility.

### Option C: Per-Firm Flat Fee (Unlimited Users and Engagements)

**Mechanics:** One annual price per firm, regardless of headcount or engagement volume.

**What it incentivizes:** Maximum breadth of usage (every team member should be in the platform), maximum engagement volume (no marginal cost disincentive to running another engagement), and deep organizational dependency. Firms are incentivized to get full value out of the platform.

**Pros:**
- Simplest possible billing. Zero resistance at renewal ("same as last year").
- Maximum platform stickiness — everyone in the firm uses it, so switching cost is enormous.
- Aligns with the "two tools become one" value proposition. Firms see one line item instead of two.
- Strong NPS driver: firms never feel "charged for more success."
- Basecamp ($299/month flat for unlimited users) proved this model works for collaboration software. The audit context has a higher justifiable price point.

**Cons:**
- Revenue does not scale with firm size or usage. A 200-person firm and a 20-person firm on the same tier generate the same revenue despite vastly different value capture.
- Mispricing risk: set too low and Axiom leaves money on the table with large firms; set too high and small firms self-select out.
- Requires accurate segmentation by firm size at sign-up (tiering by headcount or revenue) to prevent systematic mispricing.
- No expansion revenue mechanism unless firm upgrades tiers.

**Verdict:** Strong model for the SMB and lower mid-market segment (20–60 staff). Becomes problematic without firm-size tiers above that. Works best as the structure within a tier rather than as the sole pricing axis.

### Option D: Hybrid Flat Base + Per-Engagement Overage

**Mechanics:** An annual flat fee covers a defined number of engagements (e.g., 30 per year). Additional engagements above the included allotment are billed at a flat per-engagement rate. Users are unlimited within the flat fee.

**What it incentivizes:** Firms start engagements freely up to their included volume, building strong habits and platform dependency. The overage mechanism creates natural expansion revenue without penalizing moderate users. Firms self-select into higher tiers as they grow.

**Pros:**
- Predictable base revenue (flat annual fee per firm) with upside from high-volume users.
- The "unlimited users" component is a differentiator vs. CaseWare and creates stickiness.
- Overage billing is familiar to SaaS buyers (similar to Twilio, AWS, and usage-based SaaS across many categories).
- Expansion ARR occurs organically as firm engagement volume grows — no sales call required to upgrade.
- Allows pricing segmentation by firm size (tier 1: up to 30 engagements; tier 2: up to 75 engagements; tier 3: unlimited) without per-user mechanics.

**Cons:**
- Slightly more complex than pure flat fee. Firms need to track included engagements vs. used.
- Overage billing creates potential sticker shock if not communicated clearly.
- Requires robust engagement-counting instrumentation in the product.
- Overage rates must be set carefully — too high and firms feel penalized for success; too low and the model barely outperforms a flat fee.

**Verdict:** Best fit for Axiom's target ICP. Captures base revenue predictably, scales with high-volume firms, aligns with platform value, and supports self-serve adoption via a well-defined entry tier.

---

## 3. Self-Serve vs. Sales-Assisted: Price and Complexity Thresholds

### General SaaS Benchmarks

Research consistently shows:
- **Under $5,000 ACV**: fully self-serve; credit card signup; no sales touch.
- **$5,000–$25,000 ACV**: inside sales or sales-assisted self-serve; demo may be offered but not required; credit card or invoice.
- **$25,000–$50,000 ACV**: inside sales required; demo expected; contract review expected.
- **$50,000+ ACV**: field sales required; formal procurement; security review; multi-stakeholder approval.

Axiom's target ACV range of **$12,000–$30,000/year** places it squarely in the inside-sales-assisted zone, but with a strong product and low complexity, self-serve can extend further up the range than typical.

### Audit Software Specifically

Audit firms have specific purchasing behaviors that shift these thresholds:

- **Firms expect a demo** before signing any software that will touch client engagement files. This is not negotiable for the partner making the decision — they will not click "Start Free Trial" and give Axiom access to client data without understanding the security posture first.
- **Security and data residency questions** arise at any price point. A single partner asking "where is client data stored?" before entering a credit card is enough to break a self-serve flow.
- **The $10,000–$15,000/year threshold** is where CPA firms typically involve a second decision-maker (managing partner or IT director). Below this, a single partner can approve. Above it, a firm meeting or vote is likely.
- **Contract language matters.** Audit firms have professional liability insurance requirements. They may need to verify that their E&O policy covers data held in third-party platforms. This creates demand for a BAA or data processing agreement regardless of price.

### Recommended Self-Serve Threshold for Axiom

- **Under $15,000/year (Growth tier):** Support self-serve signup with credit card OR invoice. Provide in-app demo flow (interactive walkthrough, not a human call). Offer an optional "onboarding call" rather than requiring one.
- **$15,000–$30,000/year (Scale tier):** Offer a "book a demo" flow prominently, but allow firms to start a 14-day trial without a call. Convert the trial with a short inside sales touch (email + 30-minute call). Contract sent via DocuSign.
- **$30,000+/year (Enterprise tier):** Sales-led. No self-serve. Custom SOW, dedicated CSM, BAA negotiation, security review package provided.

The key insight: **self-serve does not mean no human touch.** It means the firm can start using the product, explore the value, and make a low-risk commitment without a mandatory sales call. A sales-assisted touch should be triggered by behavior (trial started, first engagement created, invite sent to a colleague) not by price alone.

---

## 4. Freemium vs. Free Trial Strategy

### What Comparable B2B SaaS Tools Do

- **Free trial (time-limited, full features):** Most common in professional services SaaS. 14-day or 30-day trials. ChartMogul data shows 57% of B2B SaaS products use a free trial as the primary acquisition path. Conversion rate benchmarks: 18–25% (opt-in) to 48% (opt-out/credit-card-required).
- **Freemium (permanent free tier, limited features):** More common in collaboration and productivity tools. Conversion rate: 5–8% of free users convert to paid. Works when the free tier has broad consumer reach (Slack, Notion, Figma). Rarely effective in vertical B2B where the buyer is a professional making a deliberate tool decision.
- **Interactive demo / sandbox (no account required):** Increasingly popular in B2B. Allows a prospect to experience core workflows without provisioning real data. Low commitment; good for early-funnel education.

### Should Axiom Offer a Free Tier?

**Recommendation: No permanent free tier. Yes to a 14-day full-feature trial.**

Reasons against a free tier:

1. **Trust and perception risk is real in audit.** Audit software handles sensitive client financial data. A free tier signals one or more of: the product is not yet enterprise-grade, the business model is uncertain, or client data may be monetized. CPA partners are trained risk managers — they will notice these signals.
2. **Freemium in professional services dilutes brand positioning.** Axiom is positioned to replace a $40,000/year toolstack. A free tier undercuts the "premium, unified platform" narrative before a firm even starts a trial.
3. **Free tier users rarely convert in vertical SaaS.** The target buyer (a managing partner at a 30-person CPA firm) is not a product-led self-discovery buyer. They need a business case, not a freemium sandbox.
4. **Free tier creates support burden.** Audit software is complex. Free users who cannot figure out the trial will generate support tickets, churn without conversion, and leave negative reviews. These users are not the ICP.

Reasons to offer a 14-day trial:

1. **Removes adoption risk.** The ICP has been burned by CaseWare's implementation complexity and Fieldguide's sales process. A trial that lets a firm run one real engagement proves value without a commitment.
2. **Conversion benchmarks justify it.** B2B accounting/finance software converts free trials at 20%+ (above average). A 14-day trial with a real engagement is a high-intent activity — the partner who completes it has demonstrated willingness to use the product.
3. **Data point from auditee-side compliance tools:** Vanta, Drata, and Secureframe all use trial/demo flows, not freemium, for their CPA and enterprise buyers. The pattern is established.
4. **It can be gated.** The trial can require a business email domain and a brief intake form (firm name, number of staff, primary audit types). This qualifies the lead without requiring a sales call, while filtering out non-ICP trial signups.

**Trial design recommendations:**
- 14 days, full features, no credit card required to start.
- One included engagement template pre-loaded on signup (e.g., SOC 2 Type 2 or financial audit — selectable).
- In-app progress nudges ("You've completed 3 of 5 setup steps — you're 2 steps from running your first engagement").
- Outreach from a customer success rep on Day 3 (automated email) and Day 10 (human email or call offer) if engagement has been created.
- Trial does not delete data on expiry — it locks the workspace until a plan is selected (no data loss anxiety).

---

## 5. Pricing Model Implications for Product Design

The choice of pricing model has direct consequences for product architecture, feature prioritization, and the incentive structures embedded in the UX. These are first-order and second-order effects of the recommended hybrid model.

### Engagement as the Core Unit of Value

If the pricing unit is the engagement (included volume + overage), then the product must make engagement creation frictionless and engagement completion visible. The product implications:

- **Engagement setup must be fast.** A firm that can create and staff a new SOC 2 engagement in under 10 minutes will use the platform for more engagements. This reinforces both the value prop and the pricing model.
- **Engagement count dashboard is a first-class feature.** Firms need to know how many included engagements they have used vs. remaining. This is not a billing footnote — it is a feature. A visible "12 of 30 engagements used" counter in the dashboard creates gentle urgency to upgrade before a busy season.
- **Engagement archiving vs. deletion matters.** If firms can archive completed engagements and re-open them (for a re-engagement the following year), the pricing model needs to define whether a rollforward counts as a new engagement. Clear rules prevent disputes; ambiguity creates churn.

### Unlimited Users as a Differentiation and Stickiness Driver

Because the model is per-firm with unlimited users (not per-seat), the product must be designed to benefit from multi-user participation:

- **Role-based access is a must-have from day one.** With unlimited users, firms will add partners, managers, seniors, and staff. The product must support different permission levels or senior partners will not add junior staff.
- **Real-time collaboration is a competitive feature, not a nice-to-have.** If the pricing model rewards adding more users, the product must make collaboration (concurrent review, comment threads, review sign-offs) visibly better than the single-user CaseWare experience.
- **Invite flows should be prominent and frictionless.** Every additional person added to a firm's Axiom account increases switching cost. The product should actively nudge team expansion ("Add your audit staff to collaborate on this engagement →").

### Overage Mechanism Creates a Natural Upgrade Path

The per-engagement overage creates a self-service upgrade path:

- When a firm hits their included engagement count, the UX should prompt tier upgrade before blocking new engagement creation. "You've used all 30 included engagements. Upgrade to Scale for unlimited engagements + AI analytics — or add a 5-pack for $X."
- **Usage analytics** (engagements per month, AI actions per engagement, evidence items processed) should be available to firm admins and to Axiom's CS team. These are expansion signals.
- The product should make it easy to upgrade mid-year (prorated billing) so a firm hitting their limit in October does not defect to a competitor rather than upgrade.

### Per-Engagement Pricing Rewards Speed; Product Must Deliver

If engagements are the value unit, the product must make each engagement demonstrably faster than the alternative:

- **Time-to-first-value instrumentation.** Axiom should track time from engagement creation to first AI-extracted evidence item, first review comment, and first signed-off workpaper. These metrics are both product KPIs and selling points.
- **Template quality is a direct revenue lever.** A poor SOC 2 template forces firms to spend hours customizing before starting fieldwork. Every hour of setup time is a reason to question the platform's cost. Pre-built, methodology-accurate templates for the top 10 engagement types are a business imperative, not a nice-to-have.

### Second-Order Effects to Monitor

- **Gaming risk:** Firms may try to combine multiple engagements into one to stay within their included count. The product should define "engagement" clearly (one engagement = one client, one framework, one period). This needs to be in the terms of service and the product UI should enforce it.
- **Seasonal revenue cliff:** If most firms' busy season is January–April (calendar-year financial audits), a material number of overage events will cluster in Q1. This is predictable and positive for ARR but should be anticipated in cash flow planning.
- **Data export pressure:** Because unlimited users are included, Axiom will accumulate large amounts of engagement data per firm. At renewal, firms will ask "can we export everything?" The product needs a robust export capability — and the export experience should be good enough that firms feel safe, but not so easy that switching to a competitor is trivial.

---

## 6. Revenue Model at Scale

### Assumptions for Modeling

| Parameter | Assumption |
|---|---|
| Growth tier ACV | $14,400/year ($1,200/month) |
| Scale tier ACV | $28,800/year ($2,400/month) |
| Enterprise tier ACV | $60,000/year (blended average) |
| Tier distribution (early stage) | 60% Growth, 30% Scale, 10% Enterprise |
| Tier distribution (mature) | 40% Growth, 40% Scale, 20% Enterprise |
| Annual gross churn | 8–12% (see below) |
| Net revenue retention | 110–120% (expansion from overage + tier upgrades) |
| Average CAC (self-serve) | $2,000–$4,000 |
| Average CAC (sales-assisted) | $8,000–$15,000 |

### ARR Model by Milestone

**50 Firms (Early Stage)**

- Mix: 30 Growth, 15 Scale, 5 Enterprise
- ARR: (30 × $14,400) + (15 × $28,800) + (5 × $60,000) = $432,000 + $432,000 + $300,000 = **~$1.16M ARR**
- Monthly burn context: at $1.16M ARR with ~70% gross margins, this supports a 5–8 person team.

**200 Firms (Growth Stage)**

- Mix: 90 Growth, 80 Scale, 30 Enterprise
- ARR: (90 × $14,400) + (80 × $28,800) + (30 × $60,000) = $1,296,000 + $2,304,000 + $1,800,000 = **~$5.4M ARR**
- This is a realistic 24–36 month milestone for a well-executed PLG + inside sales motion.

**500 Firms (Scale Stage)**

- Mix: 175 Growth, 225 Scale, 100 Enterprise
- ARR: (175 × $14,400) + (225 × $28,800) + (100 × $60,000) = $2,520,000 + $6,480,000 + $6,000,000 = **~$15M ARR**
- With NRR of 115%, organic expansion adds ~$1.5–2M/year on top of new logo revenue.

Note: these models do not include overage revenue (5–15% of ACV is a reasonable estimate for overage in a volume-based model) or add-on modules (e.g., client portal, advanced analytics, multi-framework mapping). Including those, the realistic ARR at 500 firms is **$16–18M**.

### LTV/CAC Targets

| Metric | Target | Notes |
|---|---|---|
| LTV/CAC ratio | 4:1 to 6:1 | Best-in-class B2B SaaS is 3:1 minimum; 5:1 is healthy growth |
| Payback period | 12–18 months | Achievable with low CAC self-serve motion |
| Blended CAC | $5,000–$8,000 | Weighted average of self-serve ($2,500) and sales-assisted ($12,000) |
| LTV at 8% annual churn | $14,400 / 0.08 = $180,000 (Growth tier); $28,800 / 0.08 = $360,000 (Scale tier) — blended ~$200,000 | |
| LTV/CAC (Growth tier, $6,000 blended CAC) | $180,000 / $6,000 = **30:1** (theoretical max; gross margin adjusted ~10:1 at 70% GM) | |

At 70% gross margin, adjusted LTV/CAC is approximately 8:1 to 12:1 for the Growth tier — well above the 3:1 minimum. This is achievable because audit software is very sticky once adopted, and CAC can be kept low with a PLG motion.

### Churn Rate Analysis for Audit Software

**Why audit software churns less than the average B2B SaaS:**

1. **Engagement history is irreplaceable.** Prior-year audit files, workpaper rollforwards, and evidence trails are stored in the platform. Moving five years of engagement history to a new platform costs €40,000–60,000 in migration effort (ciferi data). This switching cost is not theoretical — it was the dominant reason mid-market firms do not leave CaseWare even when unhappy with it.
2. **Compliance record retention requirements.** AICPA standards require firms to retain audit documentation for a minimum of five years (seven for public company audits). A firm that switches platforms must either migrate historical records or maintain two systems. Neither is attractive.
3. **Staff training investment.** Audit software touches every level of the firm. Once senior managers and partners build muscle memory on the platform, retraining is a significant hidden cost.
4. **The annual engagement cycle creates a natural renewal moment, not a churn moment.** CPA firms renew software in August–September (before US fiscal year) or January (after busy season). This is a high-consideration renewal, not an automatic churn risk.

**Churn rate benchmarks:**
- Average B2B SaaS annual churn: 3.5% monthly = ~35% annually (high; dominated by SMB horizontal tools)
- B2B SaaS with high switching costs (enterprise infrastructure): 1–5% annually
- Professional services vertical SaaS (legal, accounting): industry benchmarks show 8–15% annual churn for mid-market, lower for enterprise
- **Axiom target:** 8–12% gross annual churn; 90–95% net revenue retention baseline; 110–120% NRR with expansion

**The key churn risk for Axiom is not retention of happy customers — it is firm-level business volatility:** mergers and acquisitions (two firms merging onto one platform), economic downturns causing firms to cut software spend, or a partner retirement that removes the internal champion. These are manageable with proactive CS touchpoints but not fully eliminable.

---

## 7. Recommended Pricing Model and Tiers

### Model: Hybrid Flat Base (Unlimited Users) + Per-Engagement Tier Structure

### Tier Design

**Growth — $1,200/month ($14,400/year)**

Best for: Partner-led firms of 10–40 staff, 20–50 engagements/year, evaluating or migrating off CaseWare + DataSnipper.

- Unlimited users
- Up to 35 active engagements/year included
- Financial audit support (trial balance, workpaper management, sampling, materiality calculators)
- Compliance audit support (SOC 2 Type 1/2, HIPAA, ISO 27001)
- Standard AI evidence extraction (PDF, Excel, image)
- Pre-built methodology templates (AICPA/GAAS, SOC 2 TSP, ISO 27001)
- Client document request portal (PBC list management)
- Standard email support + in-app help
- Overage: $350/additional engagement

**Scale — $2,400/month ($28,800/year)**

Best for: Firms of 30–100 staff, 50–150 engagements/year, running multi-framework audits and needing advanced collaboration.

- Everything in Growth
- Up to 100 active engagements/year included
- Cross-framework evidence mapping (one evidence item maps to multiple frameworks)
- AI-assisted risk assessment and control gap analysis
- Advanced analytics dashboard (engagement profitability, staff utilization, review cycle time)
- Multi-entity / group audit support
- Custom methodology template editor
- Priority support + dedicated onboarding call
- Overage: $250/additional engagement

**Enterprise — Custom Pricing (~$50,000–$120,000/year)**

Best for: Firms of 80–200+ staff, 150+ engagements/year, requiring deep customization, white-glove onboarding, and enterprise security controls.

- Everything in Scale
- Unlimited engagements
- Dedicated Customer Success Manager
- Custom SLA (99.9% uptime guarantee)
- BAA / data processing agreement (required for HIPAA-covered engagements at this scale)
- SSO / SAML integration
- Audit trail export and long-term data retention configuration
- Custom integrations (ERP, practice management, GL systems)
- Security review package (SOC 2 Type 2 report, penetration test summary, data residency documentation)
- Negotiated contract terms (multi-year discounting, custom billing cadence)

### Price Point Rationale

**Why $1,200/month ($14,400/year) for Growth:**
- This is below the single-decision-maker approval threshold for most CPA firm partners (~$15,000/year). A managing partner at a 30-person firm can approve this without a firm vote.
- It represents approximately 25–35% of the current toolstack spend (CaseWare + DataSnipper), making the ROI case trivially easy to make: consolidate two tools into one for less than half the price.
- Yak's $800/year compliance-only engagement establishes the floor. Axiom serves both financial and compliance audit — $14,400/year (18× Yak's price) is justified by the breadth of coverage.
- At 70% gross margin, Axiom earns ~$10,000/year per Growth customer. With a $5,000 blended CAC, payback is 6 months.

**Why $2,400/month ($28,800/year) for Scale:**
- Positions Axiom directly against Fieldguide's estimated mid-market entry price ($30,000–$50,000/year) at a significant discount, with a self-serve path rather than a sales-required process.
- The 2× step-up from Growth is justified by the 3× increase in included engagements, cross-framework capabilities, and advanced AI features that firms running 50–150 engagements genuinely need.
- Firms on Scale are spending approximately 50–70% of their current toolstack cost, with no separate DataSnipper license required.

**Why "custom" for Enterprise:**
- Firms at this scale expect a negotiated price. Publishing a specific number creates anchoring that limits deal size flexibility.
- The Enterprise tier needs to justify a CSM, security review, and dedicated onboarding — these have real costs that need to be reflected.
- A soft anchor of "$50,000–$120,000" communicated verbally in sales calls is appropriate.

### Annual vs. Monthly Billing

- Annual billing is strongly preferred and should be incentivized: offer 2 months free (16.7% discount) for annual prepayment.
- Monthly billing available at Growth tier only; Scale and Enterprise require annual contracts.
- Rationale: annual contracts materially reduce churn risk (firms cannot cancel month-to-month during busy season); they also improve cash flow predictability. The discount aligns with industry norms.

### Self-Serve vs. Sales-Assisted by Tier

| Tier | Signup Flow | Sales Touch |
|---|---|---|
| Growth | Self-serve: credit card or invoice, 14-day trial | Optional: automated email sequence, 30-min call offered on Day 10 |
| Scale | Self-serve trial start (14 days), then inside sales to close | Inside sales rep follows up at trial Day 3 and Day 10 |
| Enterprise | Request demo → assigned AE → custom quote | Full sales-led, 4–8 week cycle |

---

## 8. What This Model Implies for Product Roadmap Priorities

The recommended pricing model directly shapes product priorities. In order of impact:

1. **Engagement count instrumentation** (Day 1 requirement). Firms must see how many engagements they have created, active, and completed. This is part of the core data model, not a reporting afterthought.

2. **Frictionless engagement creation** (Day 1 requirement). The entire pricing model is anchored on engagements as the value unit. If creating a new engagement takes more than 5 minutes, the model fails.

3. **Self-serve onboarding flow** (Day 1 requirement for Growth tier). A firm signing up without a sales call must be able to go from account creation to first engagement created within one hour. This requires: firm profile setup, template selection, default roles, client workspace creation, and PBC request list — all guided by in-app UX.

4. **Unlimited-user role-based access control** (Required before Scale launch). Without RBAC, firms will not add junior staff. Without junior staff, platform stickiness is low and expansion revenue from invite flows is unavailable.

5. **Cross-framework evidence mapping** (Scale tier differentiator). This is the architectural feature that separates Scale from Growth and from all competitors. Required for Scale launch, not Growth.

6. **Overage billing and tier upgrade flows** (Required at launch). The product must handle the overage state gracefully: warning at 80% usage, prompt at 100%, no hard block that locks a firm out of active engagements.

7. **Engagement history export** (Required before any enterprise deal). An enterprise firm will not commit without knowing they can get their data out. The export must be comprehensive (PDF workpapers, evidence files, metadata) and available on request.

8. **Usage analytics dashboard** (Required for CS-led expansion). The Scale and Enterprise CSM motion depends on visibility into which firms are hitting capacity constraints, how frequently AI features are used, and where engagement completion rates drop off. This is both a product feature and a revenue enablement tool.

---

## Sources Consulted

- SoftwareAdvice, Capterra, GetApp: Fieldguide, AuditBoard, CaseWare, DataSnipper pricing profiles (2025–2026)
- Vendr Marketplace: AuditBoard contract data ($42,833 median ACV)
- ciferi: "Audit Software for Non-Big 4 Firms: CaseWare Alternatives" (2026) — Bakker Audit TCO data
- SoftwareFinder: DataSnipper pricing guide 2026 (tier estimates: Start $64/user/month, Accelerate ~$175/user/month)
- PricingNow: Fieldguide pricing analysis 2025
- CPA Practice Advisor: Fieldguide Series C announcement (Feb 2026, $75M, Goldman Sachs, $700M valuation)
- Fieldguide blog: "Modern Pricing Strategies Shaping the Firm of the Future" (per-engagement philosophy)
- Fishbowl App: DataSnipper community pricing discussion
- ChartMogul: SaaS Conversion Report (free trial vs. freemium conversion benchmarks)
- Pathmonk / Userpilot: Free-to-paid conversion rate benchmarks (18–25% opt-in; 48% opt-out)
- SaaS Capital: ACV benchmarks and self-serve vs. enterprise motion thresholds
- SaaS Hero: B2B SaaS LTV/CAC ratio benchmarks 2026 (3:1–5:1 healthy range)
- Vitally / Kalungi: B2B SaaS churn rate benchmarks (3.5% average monthly; enterprise 1–5% annually)
- CustomerGauge: Average churn rate by industry (professional services: 27% annual — notable outlier)
- Focus Digital: SaaS churn rate by industry 2025
- HubiFi: B2B SaaS benchmarks — churn, CAC, LTV, growth rates
- Tight.com: Vertical SaaS churn reduction for SMB
- Close.com: Self-serve to sales-supported SaaS transition analysis
- Chaotic Flow: Three SaaS sales models (self-service / transactional / enterprise)
- Monetizely: Freemium for enterprise open-source analysis
- Logicalcommander: B2B SaaS risk management and compliance trust signals
- Scalekit: Security and authentication as B2B SaaS trust driver
- MRRSaver: SaaS churn benchmarks 2026 by segment and ACV
