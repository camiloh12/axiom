# Axiom User Journeys

Task-level user journey maps for every major platform workflow, organized by persona and goal. Referenced from Sections 11–12 of the [product spec](../specs/axiom-spec-design.md).

Competitor workflows (auditor-side: Fieldguide, Agentive, Yak; auditee-side: Drata, Vanta, Secureframe, Sprinto, Thoropass; GRC adjacencies: AuditBoard CrossComply, Hyperproof) informed the pain points and opportunities throughout. Axiom's ICP is compliance and assurance work (SOC 2, ISO 27001, ISO 27701, ISO 42001, HIPAA, PCI DSS, and optional SOC 1). Customer firm types in scope: CPA firms delivering compliance attestations, compliance-first consultancies (the tier below Schellman / A-LIGN), accredited ISO Certification Bodies (CBs), and PCI QSA firms. Axiom supports their engagement-delivery and evidence-collection workflow; legal sign-off authority (CB-issued certificates, QSA-signed ROCs/AOCs, AICPA-licensed attestation opinions) remains with the licensed firm. The platform is both-sided: auditor-side tooling plus a full auditee GRC workspace extending from the Client Hub. Where Axiom introduces flows with no competitor equivalent, those are called out explicitly.

---

## Persona Reference

| Role | Description | Journeys |
|------|-------------|----------|
| **FirmAdmin** | Managing partner or designated admin. Configures the firm, manages billing and staff. | 1 |
| **Partner** | Engagement partner. Creates engagements, approves quality documentation, signs off on reports. Typically serves as the **General Reviewer** and/or **Final Reviewer** in the four-level workpaper sign-off hierarchy. 10+ years experience. | 3, 9, 11 |
| **Manager** | Engagement manager. Typically serves as the **Detailed Reviewer** in the four-level workpaper sign-off hierarchy; may serve as **General Reviewer** depending on firm policy. Manages review notes, advances engagement phases. 5–8 years experience. | 6 |
| **Staff Auditor** | Performs hands-on compliance and assurance work — control testing, evidence review, workpaper drafting. Acts as the **Tester** in the four-level sign-off hierarchy. 1–3 years experience. | 2, 5, 7 |
| **EQR Reviewer** | Independent quality reviewer per SQMS 2 (for SOC attestations) and equivalent ISO 17021-1 review expectations for ISO engagements. Not on the engagement team. Independent of the four-level reviewer chain — EQR is a separate engagement-level review track. | 10 |
| **Client Contact / ClientUser** | Security, compliance, IT, HR, or privacy lead at the audited company. Non-auditor. | 8, 12 |
| **ClientAdmin** | Elevated client role. Owns the auditee GRC workspace, delegates document requests, manages continuous monitoring. | 4, 8, 12 |

---

## Journey Index

| # | Persona | Goal |
|---|---------|------|
| [1](#journey-1-set-up-firm-and-launch-first-engagement) | FirmAdmin | Set up firm and launch first engagement |
| [2](#journey-2-join-platform-and-reach-first-task) | Staff Auditor | Join platform and reach first task |
| [3](#journey-3-create-and-scope-a-new-engagement) | Partner | Create and scope a new engagement |
| [4](#journey-4-cross-framework-evidence-mapping) | ClientAdmin | Upload one piece of evidence and see which framework requirements it satisfies across SOC 2, ISO 27001, HIPAA, PCI, and more |
| [5](#journey-5-test-controls-and-prepare-workpapers) | Staff Auditor | Test controls and prepare workpapers |
| [6](#journey-6-review-workpapers-and-advance-the-engagement) | Manager | Review workpapers and advance the engagement |
| [7](#journey-7-manage-document-requests-and-collect-evidence) | Staff Auditor | Manage document requests and collect evidence |
| [8](#journey-8-fulfill-audit-document-requests) | Client Contact | Fulfill audit document requests |
| [9](#journey-9-generate-report-finalize-and-archive) | Partner | Generate report, finalize, and archive |
| [10](#journey-10-conduct-engagement-quality-review) | EQR Reviewer | Conduct engagement quality review |
| [11](#journey-11-multi-framework-integrated-engagement) | Partner | Scope a single engagement covering SOC 2 + ISO 27001 + ISO 27701 simultaneously |
| [12](#journey-12-continuous-assurance-auditee-side) | ClientAdmin | Respond to drift alerts, auto-retest affected controls, and notify the auditor of material changes |

---

# Journey 1: Set Up Firm and Launch First Engagement

## Overview
- **Persona:** FirmAdmin — managing partner or IT-responsible partner at a 20–60 person CPA firm doing compliance attestations, or at a compliance-first consultancy (the tier below Schellman / A-LIGN)
- **Goal:** Get the firm from trial signup to a first active engagement in under one week, without outside help
- **Trigger:** A Drata / Vanta / Hyperproof quote that felt bolted-on for auditor workflows, a painful SOC 2 + ISO 27001 crossover where evidence had to be re-collected, a new ISO 42001 prospect the firm has no tooling for, or post–Delve-scandal scrutiny from a partner asking "how do we know our AI-assisted work is defensible?"
- **Stages:**
  1. Sign up and verify identity
  2. Configure firm profile
  3. Activate methodology templates
  4. Create first engagement
  5. Invite staff

## Stage 1: Sign Up and Verify Identity

### Sub-goal
Create a firm account on Axiom and prove business legitimacy.

### User Actions
- Visits the Axiom marketing site and clicks "Start Free Trial"
- Enters a business email address (personal email domains are accepted but flagged for follow-up)
- Receives and clicks the verification email
- Completes a brief intake form: firm name, staff count (dropdown: 1–10, 11–20, 21–40, 41–60, 60+), primary engagement types (multi-select: SOC 2, ISO 27001, ISO 27701, ISO 42001, HIPAA, PCI DSS, SOC 1), country (US / Canada)

### Touchpoints
- Marketing site → signup form
- Verification email (delivered within 60 seconds)
- Intake form (single page, under 2 minutes to complete)

### Thoughts & Emotions
- **Cautious optimism** — "Let me see if this actually works without a sales call"
- **Mild anxiety** — "Is this going to be another AuditBoard situation where I spend two weeks in implementation workshops before I see anything?"
- **Relief** — verification and intake take under 3 minutes total

### Pain Points
- **Competitor context:** Drata, Vanta, and Secureframe are all self-serve but designed around auditee-only workflows. Hyperproof and AuditBoard require sales cycles and onboarding implementations. Agentive requires a demo request before platform access.
- If the verification email is slow or lands in spam, the entire momentum of "let me try this right now" dies

### Opportunities
- Instant email delivery with clear "check your inbox" interstitial and spam-folder guidance
- Show a countdown: "Your trial starts now — 14 days to run your first engagement"
- Pre-populate engagement type selection in the next step based on intake form answers

---

## Stage 2: Configure Firm Profile

### Sub-goal
Establish the firm's identity on the platform so engagements and reports carry the correct branding.

### User Actions
- Enters firm display name, uploads logo (optional), sets timezone, enters billing contact email
- Reviews and confirms — this is a lightweight step, not a configuration labyrinth

### Touchpoints
- In-app setup wizard (step 2 of 5 on the onboarding checklist)
- Collapsible progress panel in the sidebar showing completed and remaining steps

### Thoughts & Emotions
- **Comfortable** — this feels like normal SaaS setup, nothing audit-specific yet
- **Slightly impatient** — "When do I get to see the actual product?"

### Pain Points
- **Competitor context:** Hyperproof and AuditBoard onboarding require configuring security roles, user groups, and access hierarchies before anyone can use the platform. This creates hours of admin work before any compliance work happens.
- If the firm profile step requires too many fields, it feels like a registration form rather than product setup

### Opportunities
- Keep this to 4 fields maximum — name, logo, timezone, billing email
- Auto-detect timezone from browser
- Make logo upload truly optional with a clean default (firm initials)

---

## Stage 3: Activate Methodology Templates

### Sub-goal
Select which audit frameworks the firm will use, so engagement creation has the right controls and test procedures available.

### User Actions
- Reviews a list of pre-built methodology templates: SOC 2 Type I/II (TSC 2017), ISO 27001:2022 (+ ISO 27701:2019 extension), ISO 42001:2023, HIPAA Security/Privacy Rules (with optional HITRUST CSF r2 assessment mode, post-MVP), PCI DSS 4.0.1, SOC 1 Type I/II
- Activates one or more templates with a single click each — no configuration required to activate
- Growth tier: templates are read-only pre-built. Scale tier: custom template editor is unlocked for later customization, including custom `CommonControl` nodes that layer on top of the licensed SCF / AICPA / OSCAL / CIS crosswalk stack

### Touchpoints
- Template selection screen with framework descriptions and engagement type labels
- One-click activation (toggle or checkbox per template)
- Onboarding checklist advances to step 3/5

### Thoughts & Emotions
- **Pleasantly surprised** — "I don't have to build my own methodology from scratch? Hyperproof made me configure everything."
- **Confident** — selecting known frameworks (SOC 2, ISO 27001) feels safe and familiar
- **Curious** — "What's in these templates? Can I customize them later?"

### Pain Points
- **Competitor context:** Hyperproof firms spend significant time configuring methodology templates and control libraries. AuditBoard requires Unified Controls Framework (UCF) configuration as a dedicated phase. Drata / Vanta / Secureframe are auditee-side and don't ship auditor methodology templates at all — auditors must bring their own working papers.
- If the admin can't preview what's inside a template before activating it, they'll hesitate
- Firms with custom methodologies may feel constrained by read-only templates on Growth tier

### Opportunities
- Allow template preview (expand to see control count, test procedure count, workpaper types) without requiring activation
- Show a "What's included" breakdown: "This SOC 2 Type II template includes 89 controls across all 5 trust services criteria, 200+ test procedures, 80+ pre-drafted document request templates, and pre-built STRM-encoded crosswalks to ISO 27001 Annex A, HIPAA Security Rule, and PCI DSS 4.0.1"
- Surface the custom template editor as a Scale-tier upsell naturally: "Want to customize? Upgrade to Scale for the template editor."

---

## Stage 4: Create First Engagement

### Sub-goal
Create a real (or test) engagement that scaffolds the full audit file — controls, test procedures, workpapers, and client acceptance.

### User Actions
- Clicks "New Engagement" from the engagement list
- Selects engagement type and framework (pre-filtered based on activated templates)
- Enters client name (can use "Demo Client" for a trial run), audit period dates
- Selects methodology template (pre-selected based on signup intake)
- System creates the full engagement scaffold: Engagement in Planning status, all template Control and TestProcedure records, draft Workpaper shells, and an empty ClientAcceptance record

### Touchpoints
- Engagement creation wizard (3–4 screens: type, client, period, confirmation)
- Scaffold generation happens in seconds — progress indicator while controls and workpapers are created
- Engagement dashboard appears immediately after creation
- Onboarding checklist advances to step 4/5

### Thoughts & Emotions
- **Excited** — "I just created a full engagement file in 2 minutes. This took me half a day in Hyperproof."
- **Impressed** — seeing 89 controls and 200+ test procedures already populated feels like the product is working for them
- **Slightly overwhelmed** — the engagement dashboard shows a lot of content; the first reaction might be "where do I start?"

### Pain Points
- **Competitor context:** AuditBoard and Hyperproof engagement setup involve framework selection, work program configuration, and team assignment — functional but many steps. Agentive is auditor-side but framework-agnostic — no pre-built SOC 2 / ISO control libraries, so the firm must configure from scratch.
- If scaffold generation is slow (>5 seconds), the "instant gratification" moment is lost
- The first engagement dashboard needs clear orientation — where to go first

### Opportunities
- Show a "Your engagement is ready" confirmation with a visual summary: X controls, Y test procedures, Z workpapers
- Include a "Start here" prompt pointing to the first logical action (complete client acceptance → assign team → begin fieldwork)
- For trial users, offer a pre-loaded "Demo Engagement" with sample data so they can explore a fully populated engagement immediately
- This is the completion event for the "time-to-first-engagement" product metric

---

## Stage 5: Invite Staff

### Sub-goal
Get team members onto the platform so they can be assigned to engagements.

### User Actions
- Navigates to Staff Management (prompted by onboarding checklist)
- Enters email addresses — bulk invite supported (paste multiple emails)
- Assigns a role to each invitee: Partner, Manager, Staff, or EQReviewer
- Sends invitations — each invitee receives a magic link email valid for 7 days

### Touchpoints
- Staff invitation screen with bulk-paste support
- Role assignment dropdown per invitee
- Confirmation: "5 invitations sent"
- Onboarding checklist completes (all 5 steps done) — checklist collapses permanently

### Thoughts & Emotions
- **Satisfied** — "I set up the entire firm in under an hour. No consultant, no workshops."
- **Nervous about adoption** — "Will my staff actually use this? The last tool migration was painful."
- **Ready to delegate** — can now assign staff to the engagement and start real work

### Pain Points
- **Competitor context:** AuditBoard and Hyperproof onboarding includes dedicated staff training sessions as part of implementation. Drata / Vanta are single-sided (auditee-only), so firms using them can't invite their auditors directly — evidence still gets emailed around.
- If the FirmAdmin doesn't know what role to assign (especially EQReviewer vs. Partner), the role picker needs clear descriptions
- The 7-day magic link expiry could cause problems if invitations are sent during a busy period

### Opportunities
- Include one-line role descriptions in the assignment dropdown: "Staff — performs audit work, prepares workpapers" / "EQReviewer — independent quality review, read-only engagement access"
- Send a reminder email at day 5 if the magic link hasn't been used
- Offer an optional 30-minute onboarding call for Scale-tier firms (shown as in-app prompt on Day 3 and Day 10 of trial)
- SSO setup (OAuth with Microsoft/Google on Growth/Scale, SAML on Enterprise) is available in Firm Settings but not required during onboarding

---

## Handoffs

| From | To | Information Transferred | Trigger |
|------|-----|------------------------|---------|
| FirmAdmin | Staff Auditor (Journey 2) | Magic link email with role assignment | Invitation sent |
| FirmAdmin | Partner (Journey 3) | Activated templates + first engagement scaffold | Staff invitation complete |

---

## Journey Summary

### Emotional Arc
Starts with cautious optimism ("is this really self-serve?"), builds through each fast step ("this is working"), peaks at engagement creation ("a full compliance engagement in 2 minutes"), and settles into confident readiness at staff invitation. The emotional trajectory is deliberately the opposite of AuditBoard / Hyperproof (weeks of implementation) and of Drata / Vanta from an auditor's perspective (a product built for the other side of the table).

### Cross-Cutting Pain Points
- Speed is everything — any step that takes more than 2 minutes breaks the "under one hour" promise
- Template opacity — admins need to see what's inside before committing
- Staff adoption anxiety persists even after the admin is onboarded

### Prioritized Opportunities
1. **Pre-loaded demo engagement** (high impact, low effort) — lets the admin explore a fully populated engagement before creating their own
2. **Template preview with content counts** (high impact, low effort) — builds confidence in methodology selection
3. **Guided "start here" prompt on first engagement dashboard** (high impact, medium effort) — prevents the "I created it, now what?" drop-off
4. **Day 3 / Day 10 Scale-tier onboarding call prompt** (medium impact, low effort) — safety net for firms that want guidance without requiring it

---

# Journey 2: Join Platform and Reach First Task

## Overview
- **Persona:** Staff Auditor — 1–3 years experience, comfortable with Excel, accustomed to Hyperproof / AuditBoard / spreadsheet-driven workflows
- **Goal:** Accept the invitation, get oriented on the platform, and reach their first assigned task in a single session
- **Trigger:** Invitation email from FirmAdmin
- **Stages:**
  1. Accept invitation and set up access
  2. Complete profile and notification preferences
  3. Take guided platform tour
  4. Receive first engagement assignment

## Stage 1: Accept Invitation and Set Up Access

### Sub-goal
Get into the platform with minimal friction.

### User Actions
- Opens invitation email: "You've been invited to [Firm Name]'s Axiom workspace as [Role]"
- Clicks the magic link — authenticated immediately for this browser session
- If the firm has OAuth SSO configured, prompted to link Google or Microsoft account; otherwise, prompted to set a password
- Accepts terms of service

### Touchpoints
- Invitation email with magic link (valid 7 days) and secondary password-setup link
- Browser session auto-authentication
- SSO linking prompt or password creation screen
- Terms of service acceptance (single checkbox)

### Thoughts & Emotions
- **Slightly skeptical** — "Another new tool? I just got used to the last one."
- **Relieved** — magic link means no account creation form, no confirming passwords, no "choose a username"
- **Wary** — "How long is this going to take before I can do actual work?"

### Pain Points
- **Competitor context:** Hyperproof and AuditBoard require firm admin to configure user accounts, security roles, and group assignments before staff can do any work. Agentive requires a demo-gated account provisioning flow.
- If the magic link has expired (sent >7 days ago), the user hits a dead end and has to request a new invite
- SSO linking can be confusing if the user's Google/Microsoft account doesn't match their work email

### Opportunities
- Show clear "link expired" messaging with a one-click "request new invitation" button
- Pre-select the SSO provider if the firm has already configured it — don't show both options

---

## Stage 2: Complete Profile and Notification Preferences

### Sub-goal
Set up identity and communication preferences.

### User Actions
- Views their assigned role (Partner / Manager / Staff / EQReviewer) — displayed but not self-changeable
- Updates display name
- Sets notification preferences: email digest frequency (real-time, daily, or weekly)

### Touchpoints
- Profile setup screen (single page, 3 fields)
- Role badge displayed with explanation text

### Thoughts & Emotions
- **Neutral** — standard profile setup, nothing surprising
- **Curious about role** — "What does 'Staff' mean in this system? What can I do?"

### Pain Points
- If the role explanation is too terse, the user doesn't understand their capabilities and limitations
- Notification defaults matter — real-time is noisy for most staff; daily digest is a safer default

### Opportunities
- Default to daily digest for Staff, real-time for Partner and Manager
- Show a one-line role description: "As Staff, you prepare workpapers, test controls, and manage document requests. Your work is reviewed by Managers and Partners."

---

## Stage 3: Take Guided Platform Tour

### Sub-goal
Understand the platform's layout and core features before starting real work.

### User Actions
- A 5-step in-app tour begins automatically (skippable at any point)
- Tour highlights in sequence: (1) engagement list, (2) workpaper editor, (3) evidence pool, (4) document request queue, (5) AI review panel
- Each step shows a one-sentence description and a "show me" pointer
- Tour progress is stored per-user and never repeats after completion or skip

### Touchpoints
- Guided tour overlay with step indicators (1/5, 2/5...)
- Spotlight highlighting on each feature area
- "Skip tour" button always visible

### Thoughts & Emotions
- **Engaged** — if the tour is fast and shows real value ("the AI review panel analyzes uploaded documents automatically")
- **Impatient** — if the tour feels like a slideshow rather than interactive; experienced auditors want to click around, not read tooltips
- **Oriented** — by step 5, the user has a mental map of where things live

### Pain Points
- **Competitor context:** AuditBoard and Hyperproof have steep learning curves with no strong guided tour — new users report spending days figuring out the interface. Agentive's onboarding includes dedicated training sessions, which means new staff wait for the next session.
- Generic product tours that explain UI elements without connecting to audit concepts are useless to auditors
- If the tour blocks interaction ("you must complete step 3 before you can click anything else"), it creates frustration

### Opportunities
- Frame each tour step in auditor language, not product language: "This is your evidence pool — every document your clients upload and every artifact you create lives here, linked to the `CommonControl` nodes and framework requirements they support"
- Make the tour non-blocking — let users click around freely while tour steps remain available
- Offer a "replay tour" option in the help menu for users who skipped but later want guidance

---

## Stage 4: Receive First Engagement Assignment

### Sub-goal
Land on a real piece of audit work — the first meaningful platform action.

### User Actions
- Receives an in-app notification when a Partner or Manager adds them to an engagement via EngagementTeamMember
- The engagement appears in their engagement list
- Clicks the notification link — taken directly to their first assigned control or workpaper
- Begins audit work (continues in Journeys 4, 5, or 7 depending on assignment)

### Touchpoints
- In-app notification (bell icon + toast)
- Email notification (if preferences set to real-time)
- Engagement list with the new engagement highlighted
- Direct deep-link to assigned control or workpaper

### Thoughts & Emotions
- **Motivated** — "OK, I have real work to do."
- **Confident** — the deep-link to a specific task means they don't have to navigate a complex engagement file to find their assignment
- **Focused** — transition from "learning the tool" to "doing audit work"

### Pain Points
- **Competitor context:** In AuditBoard and Hyperproof, new staff must navigate a hierarchy of workspaces and folders to find their assigned tasks. There is no notification system — assignment is often communicated verbally or via email outside the platform.
- If the notification arrives before the user has completed the tour, they may skip orientation entirely and get lost later
- If no engagement assignment comes quickly, the user has nothing to do and may not return

### Opportunities
- Time the first assignment to arrive during or immediately after the tour — coordinate with the FirmAdmin/Partner
- The notification deep-link should open the control or workpaper with contextual help: "This is a SOC 2 CC6.1 control test. Here's what you need to do."

---

## Journey Summary

### Emotional Arc
Starts with skepticism ("another tool"), moves to relief (magic link, fast setup), through mild engagement (tour), and peaks at the first real assignment. The critical moment is the transition from "learning" to "working" — if the user reaches real audit work within 15 minutes of clicking the magic link, adoption is likely. If they're still in setup mode after 30 minutes, they'll mentally file Axiom as "another GRC tool learning curve."

### Cross-Cutting Pain Points
- Time from invitation to real work must be under 15 minutes
- Role clarity — staff need to understand what they can and can't do
- The gap between "invited" and "assigned to an engagement" is dead time

### Prioritized Opportunities
1. **Deep-link to first assigned task** (high impact, low effort) — the single most important onboarding moment
2. **Audit-language tour framing** (medium impact, medium effort) — differentiates from generic SaaS onboarding
3. **Smart notification defaults by role** (medium impact, low effort) — prevents notification fatigue on day one

---

## EQR Reviewer Onboarding Note

Users with the EQReviewer role follow the same onboarding path but land in a read-only view of any engagement they are assigned to review. They cannot be added to the same engagement as a team member — the system enforces this at assignment time (Journey 10 covers the EQR workflow in detail). The same pattern is used for ISO 17021-1 internal quality review expectations on ISO 27001 / 27701 / 42001 engagements.

---

# Journey 3: Create and Scope a New Engagement

## Overview
- **Persona:** Partner — engagement partner, 10+ years experience, responsible for engagement quality and SQMS 1 compliance (or ISO 17021-1 equivalent for ISO engagements)
- **Goal:** Set up a new engagement with the correct framework(s), team, quality documentation, and (if applicable) EQR assignment so fieldwork can begin
- **Trigger:** New client signed, or recurring engagement period starting (annual SOC 2, ISO 27001 surveillance visit, PCI ROC renewal, ISO 42001 initial certification)
- **Stages:**
  1. Select engagement type and framework
  2. Configure engagement details and assign team
  3. Review AI-proposed control mappings
  4. Complete SQMS 1 client acceptance
  5. Assign EQR reviewer (if applicable)
  6. Advance to fieldwork

## Stage 1: Select Engagement Type and Framework

### Sub-goal
Choose the right engagement structure and regulatory framework.

### User Actions
- Clicks "New Engagement" from the engagement list
- Selects engagement type: SOC2_TypeI, SOC2_TypeII, SOC1_TypeI, SOC1_TypeII, ISO27001, ISO27701, ISO42001, HIPAA (optionally with HITRUST CSF r2 assessment mode, post-MVP), PCI_DSS, AgreedUponProcedures, or Advisory
- Selects the applicable MethodologyTemplate (filtered by engagement type)
- For multi-framework engagements (e.g., integrated SOC 2 + ISO 27001 + ISO 27701), selects a primary framework and adds secondary frameworks via EngagementFramework (see Journey 11 for the full multi-framework flow)
- Selects framework version — e.g., ISO 27001:2022 vs. 2013, PCI DSS 4.0.1 vs. 3.2.1 (locked after Fieldwork begins unless Partner override with documented reason; the system also warns when an inbound version migration invalidates existing mappings)

### Touchpoints
- Engagement creation wizard — step 1: type and framework
- Framework version selector with effective dates shown
- Multi-framework toggle (Scale tier: shows cross-framework evidence mapping benefit)

### Thoughts & Emotions
- **Deliberate** — framework selection is a consequential decision; the partner is careful here
- **Confident** — seeing familiar framework names (SOC 2 TSC 2017, ISO 27001:2022, PCI DSS 4.0.1, ISO 42001:2023) with version dates
- **Curious** (for multi-framework) — "How does the cross-framework mapping actually work?"

### Pain Points
- **Competitor context:** Drata and Vanta advertise 25+ frameworks but auditee-side only — they assume an external auditor is plugged in, not that the audit firm itself is running the engagement. Delve claims 25+ frameworks but was hit in March 2026 with the 493-of-494-reports boilerplate scandal, eroding trust in cross-framework output. Hyperproof and AuditBoard support multi-framework via UCF but setup is labor-intensive. Agentive is framework-agnostic — every engagement is a blank canvas.
- Framework version confusion — partners need to see clearly which version applies to the current period
- Multi-framework engagement setup must not feel like creating two separate engagements

### Opportunities
- Auto-suggest framework version based on engagement period dates
- For multi-framework, show a preview: "This integrated SOC 2 + ISO 27001 + ISO 27701 engagement will include X `CommonControl` nodes that satisfy all three frameworks simultaneously. Y requirements will have only partial coverage and will be flagged in the gap list."
- Rollforward detection: if the same client had a prior engagement, prompt "Roll forward from [prior engagement name]?"

---

## Stage 2: Configure Engagement Details and Assign Team

### Sub-goal
Define the engagement scope and get the right people assigned.

### User Actions
- Enters client name (autocomplete from existing Client records, or creates a new one)
- Sets the audit period — the wizard branches by report type:
  - **Type 1 (point-in-time)** engagements (SOC 1 Type I, SOC 2 Type I): a single "as of" date picker, e.g. "as of 2026-12-31". The system stores `period_start = period_end = as-of date`.
  - **Type 2 (continuous-period)** engagements (SOC 1 Type II, SOC 2 Type II): a date-range picker validating that the period is between **3 and 12 months** (the SOC Type 2 minimum is 3 months; longer continuous periods up to 12 months are common; periods outside this range require partner override with documented reason).
  - **ISO certification cycles**: surveillance window per framework cycle.
  - **PCI DSS**: annual cycle with 90-day ASV scan validity surfaced in scope reminders.
- For rollforward engagements: prior_engagement_id is set, and the system surfaces prior year data
- Assigns engagement team: selects users by name, assigns engagement-level roles
- System creates: Engagement (Planning status), EngagementTeamMember records, EngagementFramework records, Control records cloned from template, TestProcedure records, draft Workpaper shells, empty ClientAcceptance record

### Touchpoints
- Engagement wizard — step 2: client and period
- Team assignment panel with user search and role assignment
- Scaffold generation progress indicator
- For rollforward: prior year sidebar showing controls, workpapers, evidence reuse candidates, and prior ClientAcceptance (read-only)

### Thoughts & Emotions
- **Efficient** — client autocomplete and pre-populated templates save significant setup time
- **Pleased** (rollforward) — seeing prior year data carried forward means less re-work
- **Responsible** — team assignment is a delegation decision; the partner is deciding who owns what

### Pain Points
- **Competitor context:** AuditBoard's rollforward surfaces prior year work but users report setup can still take 30+ minutes for complex engagements. Hyperproof has strong evidence reuse but weaker engagement-scaffold rollforward. Drata and Vanta treat each assessment period as a fresh configuration.
- If scaffold generation is slow for large templates (200+ controls), the partner waits
- Rollforward from prior year needs clear visual distinction between "carried forward" and "new this year"

### Opportunities
- Show rollforward diff: "142 controls carried forward from prior year. 3 new controls added in the updated SOC 2 TSC 2017 template. 1 prior control marked as superseded. 5 mappings invalidated by ISO 27001:2022 migration — requires partner reconfirmation."
- Allow team assignment by engagement role pattern: "Same team as last year" button for rollforward engagements
- Background scaffold generation with notification when complete, if it takes >3 seconds

---

## Stage 3: Review AI-Proposed Control Mappings

### Sub-goal
Validate the AI's proposed mappings between firm control objectives and framework requirements. This is Axiom's cross-framework differentiator in action.

### User Actions
- Immediately after scaffold creation, AI proposes mappings from firm `CommonControl` nodes to each in-scope framework's requirements, using the licensed SCF + AICPA + OSCAL + CIS crosswalk as the authoritative backbone (AI never authors authoritative equivalence — it only suggests candidate links from an existing catalog)
- Partner or Manager reviews a mapping table: each `CommonControl` is shown with proposed `FrameworkRequirement` links, STRM relationship type (`equivalent-to | subset-of | superset-of | intersects-with | no-relationship`), coverage strength (full / partial / not), confidence score, and explanation text
- Reviews in bulk — all proposed mappings are confirmed by default; the partner rejects or modifies individual mappings as needed
- Partial-satisfaction rows surface a gap list — the partner sees exactly which sub-requirements of a given framework node are NOT covered by the linked evidence class
- Each confirmed mapping creates an AIDecision record

### Touchpoints
- Control mapping review screen — sortable table with control name, proposed framework links, confidence percentage, and AI explanation
- Bulk confirm button with count: "Confirm all 142 mappings"
- Individual reject/modify actions per mapping
- Low-confidence mappings highlighted for priority review

### Thoughts & Emotions
- **Impressed** — "It mapped 142 controls across SOC 2, ISO 27001, and ISO 27701 simultaneously. This would have taken me a full day."
- **Professionally cautious** — "I need to check the low-confidence ones, and I need to see every partial-coverage gap. My name goes on this engagement."
- **Delighted** (for multi-framework) — seeing one `CommonControl` satisfy CC6.1 + ISO 27001 A.5.15 + HIPAA §164.312(a)(1) + PCI 8.3 in a single row is the "aha moment" for the cross-framework value proposition
- **Reassured** — every AI suggestion links back to the licensed SCF / OSCAL / AICPA crosswalk row that grounds it; nothing is an LLM hallucination dressed up as authoritative mapping

### Pain Points
- **Competitor context:** Every serious GRC player has a common-control catalog (Drata DCF, Vanta VCF, Secureframe, AuditBoard via UCF, Hyperproof via SCF, Sprinto Magic Mapping). The differentiation is in auditor-defensibility: mappings grounded in a licensed crosswalk with STRM-encoded edges, partial-satisfaction surfaced explicitly, and every accept/modify/reject decision recorded as an AIDecision. Delve's March 2026 scandal (493 of 494 SOC 2 reports shared identical boilerplate and a single typo) is the cautionary tale — provenance and the AIDecision ledger are the direct response.
- If the AI confidence is consistently low, bulk review becomes tedious rather than efficient
- Partners unfamiliar with AI-assisted workflows may not trust the suggestions and review every single mapping individually, defeating the time savings

### Opportunities
- Show accuracy metrics: "Based on prior engagements at your firm, AI control mapping accuracy is 94%"
- Group review by confidence tier: high confidence (>0.85) in a bulk-confirm block, medium (0.65–0.85) for quick scan, low (<0.65) for individual attention
- Never show a green checkmark when coverage is partial — always show percentage and the explicit gap list
- For rollforward engagements, pre-load prior year confirmed mappings as the starting suggestion (re-confirmation still required); flag mappings invalidated by framework version migration

---

## Stage 4: Complete SQMS 1 Client Acceptance

### Sub-goal
Document quality risks and acceptance decision before fieldwork can begin — a regulatory gate, not a platform feature. (SQMS 1 for SOC engagements; ISO 17021-1 clauses 9.1–9.4 equivalents for ISO engagements; PCI DSS scoping worksheet for PCI engagements.)

### User Actions
- Opens the ClientAcceptance record from the engagement dashboard (flagged as incomplete)
- Documents quality risks identified: structured categories plus free-text narrative
- Documents firm responses to each identified risk
- Confirms independence (checkbox with named confirmation)
- Signs the acceptance — timestamped, named action
- The Planning → Fieldwork transition is now unblocked

### Touchpoints
- ClientAcceptance form (structured + narrative sections)
- Independence confirmation with auditor's name displayed
- Sign-off button with timestamp preview
- Engagement status indicator changes from "Planning — Acceptance Required" to "Planning — Ready for Fieldwork"

### Thoughts & Emotions
- **Dutiful** — this is compliance work, not exciting, but the partner understands its importance
- **Appreciative** — if the form is well-structured with clear categories, it's faster than a blank Word document
- **Confident** — the system enforces the gate, which means the firm can demonstrate SQMS 1 (or ISO 17021-1 equivalent) compliance to regulators

### Pain Points
- **Competitor context:** Most compliance GRC platforms don't enforce client acceptance as a system gate. Hyperproof and AuditBoard leave it to firm policy. Drata and Vanta don't model an acceptance decision at all — they're auditee-onboarding tools, not engagement-management systems. The system-enforced gate is an Axiom differentiator.
- If the form is too rigid (only structured categories, no narrative flexibility), it won't fit all engagement types
- AI-suggested risk categories are helpful but must be clearly marked as suggestions, not pre-populated conclusions — SQMS 1 requires the partner's independent judgment

### Opportunities
- AI suggests quality risk categories based on client industry and prior engagement findings (Tier 2 — auditor reviews and certifies, not auto-populated)
- For rollforward engagements, show prior year's acceptance as a reference sidebar (read-only)
- The acceptance record is immutable once signed — addenda require creating a new version. Surface this clearly: "This acceptance is final. Changes require a new version with documented reason."

---

## Stage 5: Assign EQR Reviewer (If Applicable)

### Sub-goal
For engagements requiring engagement quality review, assign an independent reviewer who is not on the engagement team.

### User Actions
- For higher-risk engagements (per firm policy — commonly SOC 2 Type II, ISO 27001 initial certification, ISO 42001, PCI DSS, or first-year HIPAA engagements), the partner opens the EQR assignment panel
- Selects a reviewer from users with the EQReviewer role
- System validates: reviewer has EQReviewer role AND is not an EngagementTeamMember on this engagement
- If validation fails, the assignment is rejected with a clear error explaining why
- EngagementQualityReview record is created with status = Assigned

### Touchpoints
- EQR assignment panel in engagement setup
- Reviewer selection dropdown (filtered to EQReviewer role users)
- Independence validation with pass/fail indicator
- Confirmation: "EQR assigned to [name]. They will be notified when the engagement reaches Review status."

### Thoughts & Emotions
- **Careful** — EQR assignment has regulatory implications; the partner wants to get this right
- **Confident** — system-enforced independence check means they can't accidentally violate SQMS 2

### Pain Points
- **Competitor context:** Few competitors enforce EQR independence at the system level. This is a compliance differentiator.
- Small firms may have limited EQReviewer-eligible users, making the selection constrained
- The partner needs to understand when EQR is required vs. optional

### Opportunities
- Show a clear indicator: "EQR is required for this engagement (firm policy — SOC 2 Type II)" or "EQR is optional for this engagement"
- If no eligible reviewers are available, suggest: "No EQReviewers are available who are not on this engagement's team. Consider inviting an external reviewer."

---

## Stage 6: Advance to Fieldwork

### Sub-goal
Transition the engagement from Planning to Fieldwork so the audit team can begin testing.

### User Actions
- Reviews the engagement readiness checklist: acceptance complete ✓, team assigned ✓, controls mapped ✓, EQR assigned (if applicable) ✓
- Clicks "Begin Fieldwork"
- System validates the Planning → Fieldwork guard: ClientAcceptance.accepted_at must be populated by a Partner-role user
- Engagement status changes to Fieldwork — team members are notified
- Staff auditors can now begin control testing (Journey 5) and document requests (Journey 7); for multi-framework engagements the Partner may jump to Journey 11 to reconcile sampling windows across frameworks

### Touchpoints
- Engagement readiness checklist (all green before transition is available)
- "Begin Fieldwork" button (disabled until all guards pass)
- Notification to all EngagementTeamMembers
- Engagement dashboard updates to show Fieldwork-phase views

### Thoughts & Emotions
- **Satisfied** — the engagement is properly set up with all regulatory requirements met
- **Delegating** — the partner's intensive setup work is done; staff now execute
- **Watchful** — will monitor progress from the engagement dashboard

### Pain Points
- If any guard fails (acceptance incomplete, EQR not assigned where required), the partner needs clear guidance on what's missing — not just a disabled button
- The transition notification should be actionable for staff — "Fieldwork has begun. Your assigned controls are ready for testing."

### Opportunities
- Show a pre-transition summary: "Advancing to Fieldwork will notify 4 team members and unlock 89 controls for testing. Proceed?"
- After transition, redirect the partner to the engagement dashboard with a Fieldwork-phase overview: progress by control, team workload distribution

---

## Handoffs

| From | To | Information Transferred | Trigger |
|------|-----|------------------------|---------|
| Partner | Staff Auditor (Journey 5, 7) | Engagement scaffold with assigned controls and procedures | Fieldwork transition |
| Partner | EQR Reviewer (Journey 10) | EngagementQualityReview assignment | EQR record created |
| Partner | Client Contact (Journey 8) | Engagement exists; document requests can now be sent | Fieldwork begins |

---

## Journey Summary

### Emotional Arc
Begins with deliberate decision-making (type, framework, team), peaks at the AI control mapping moment ("this would have taken me a day"), passes through the compliance-necessary but emotionally flat client acceptance, and resolves with the satisfying "Begin Fieldwork" transition. The emotional high point is Stage 3 — the cross-framework mapping is where the platform's architectural differentiator becomes tangible.

### Cross-Cutting Pain Points
- Rollforward setup must be faster than new engagement setup, not the same amount of work
- Regulatory gates (acceptance, EQR) must feel like helpful guardrails, not obstacles
- Multi-framework engagements need to feel like one engagement, not two parallel ones

### Prioritized Opportunities
1. **Cross-framework mapping preview** (high impact, medium effort) — the single strongest differentiator moment in the platform
2. **Rollforward diff summary** (high impact, medium effort) — makes year-over-year engagements feel effortless
3. **AI risk category suggestions for client acceptance** (medium impact, medium effort) — accelerates the compliance gate
4. **"Same team as last year" button** (medium impact, low effort) — small time saver with high satisfaction

---

# Journey 4: Cross-Framework Evidence Mapping

## Overview
- **Persona:** ClientAdmin — the security, compliance, or IT lead at the auditee, operating in their Client Hub / auditee GRC workspace
- **Goal:** Upload one piece of evidence (e.g., an MFA configuration screenshot, an access-review export, a vulnerability-scan report) and immediately see — with auditor-defensible grounding — which `FrameworkRequirement` nodes across every in-scope framework it satisfies in full, partial, or not at all, with a clear gap list for the partial ones and period-aware coverage for time-windowed frameworks
- **Trigger:** A new artifact becomes available (quarterly access review completed, MFA rollout finished, ASV scan posted, employee background check processed) and the ClientAdmin wants to route it once rather than upload it separately for SOC 2, ISO 27001, HIPAA, and PCI
- **Stages:**
  1. Upload the evidence artifact
  2. Classify and fingerprint the artifact
  3. Review AI-suggested mappings to `CommonControl` nodes
  4. View framework-requirement coverage (full / partial / not)
  5. Review period-aware coverage and staleness
  6. Record the AIDecision and push into engagement evidence pools

## Stage 1: Upload the Evidence Artifact

### Sub-goal
Get the artifact into Axiom once, with provenance metadata captured at the moment of upload.

### User Actions
- In the Client Hub / auditee GRC workspace, clicks "Upload Evidence" (not scoped to a specific request — this is free-form evidence the ClientAdmin wants to register)
- Drops in a file (PDF, PNG screenshot, XLSX export, JSON from an API snapshot, signed CSV from an IdP export, etc.) or pastes a link to a connected SaaS source
- Optionally provides a short title and an effective-date window ("This MFA config applies from 2025-10-01 forward")
- At upload the system fingerprints the file (SHA-256), records the uploader identity, captures signed browser / OS metadata where available, and writes the artifact to S3 with Object Lock staging — this is the provenance bedrock the product bets on post–Delve scandal

### Touchpoints
- "Upload Evidence" entry point in the auditee workspace
- Drag-and-drop upload panel with optional metadata form
- Provenance pill shown immediately after upload: "SHA-256 recorded · uploader verified · effective window 2025-10-01 → open"

### Thoughts & Emotions
- **Efficient** — "I upload this once, not four times"
- **Reassured** — the provenance pill signals auditor-defensibility without the user having to understand the detail
- **Curious** — "What will it actually match to?"

### Pain Points
- Evidence uploaded without an effective-date window can't be tested against period-aware framework requirements later
- Password-protected or encrypted files block AI mapping entirely

### Opportunities
- Auto-suggest effective window from file metadata (PDF creation date, export timestamp)
- Detect and warn on password-protected files at upload time: "We can't analyze this. Provide an unlocked copy or share a connected-source link."

---

## Stage 2: Classify and Fingerprint the Artifact

### Sub-goal
Identify what kind of evidence this is, so the mapping step has a typed anchor.

### User Actions
- AI (Claude Haiku on the fast path; Sonnet if low confidence) classifies the artifact into an evidence class (e.g., `AccessReview`, `MFAConfiguration`, `VulnerabilityScanReport`, `BackgroundCheck`, `EncryptionAtRestAttestation`, `IncidentResponsePlan`, `ChangeApprovalRecord`)
- Confidence score and top 3 class candidates are shown; ClientAdmin can override
- Classification is recorded on the `EvidenceItem` and becomes the join key to `CommonControl` nodes in the mapping step

### Touchpoints
- Classification card with AI-suggested evidence class, confidence, and 2 alternate candidates
- "This is something else" override with searchable evidence-class taxonomy

### Thoughts & Emotions
- **Oriented** — the ClientAdmin understands how Axiom is "thinking" about the artifact
- **Trusting** — classification with confidence and alternates feels more honest than a single hard guess

### Pain Points
- Evidence classes that span multiple types (e.g., a combined "access review + MFA config" screenshot) confuse single-label classification
- Unknown evidence classes (new SaaS vendor, novel artifact) require taxonomy extension

### Opportunities
- Allow multi-class classification for composite artifacts
- Firm-level extensions to the evidence-class taxonomy (Scale tier)

---

## Stage 3: Review AI-Suggested Mappings to CommonControl Nodes

### Sub-goal
See which `CommonControl` nodes this artifact could satisfy, using the licensed crosswalk as the authoritative backbone.

### User Actions
- Axiom proposes candidate `satisfies` edges from this `EvidenceItem` to one or more `CommonControl` nodes, grounded in the SCF + AICPA + OSCAL + CIS crosswalk (AI suggests, does not author authoritative equivalence)
- Each candidate edge shows: the `CommonControl` name, the evidence-class pattern that triggered the suggestion, an explanation ("MFA configuration screenshots typically satisfy authentication-strength common controls"), and confidence
- ClientAdmin accepts, modifies, or rejects each candidate — this writes STRM-encoded edges and an AIDecision record for every action (accept / modify / reject)

### Touchpoints
- Mapping review table: `CommonControl` name, explanation, confidence, action buttons
- Filter "high confidence only" for bulk accept

### Thoughts & Emotions
- **Productive** — "One upload, multiple controls covered"
- **Professional** — "Every decision I make is logged; I can explain this to my auditor"

### Pain Points
- If candidate suggestions aren't grounded in the licensed crosswalk, they feel like hallucinations — this erodes trust fast
- Too many low-confidence candidates overwhelm the ClientAdmin

### Opportunities
- Link every suggestion back to the crosswalk row that grounds it (SCF ID, OSCAL control ID, AICPA mapping reference)
- Default-hide candidates below a configurable confidence floor

---

## Stage 4: View Framework-Requirement Coverage (Full / Partial / Not)

### Sub-goal
Answer the question the ClientAdmin actually came to answer: "Which framework requirements does this evidence satisfy?"

### User Actions
- The mapping panel expands into a coverage matrix: each confirmed `CommonControl` edge fans out to every framework requirement it links to (via the crosswalk)
- Each framework-requirement cell renders one of three states: **Full** (green), **Partial** (amber with percentage), or **Not** (gray)
- Partial rows expose a gap list — the specific sub-requirements NOT covered by this artifact — and a "What would close this gap?" prompt with suggested additional evidence classes
- Example rendering for a single MFA configuration export: SOC 2 CC6.1 — Full · ISO 27001 A.5.15 — Full · ISO 27001 A.8.5 — Partial (2 of 3 authentication factors evidenced; privileged-access MFA not covered) · HIPAA §164.312(a)(1) — Full · PCI DSS 8.3 — Partial (console MFA covered; remote-access MFA requires separate evidence) · ISO 42001 A.6.2.5 — Not (this AI-impact control is unrelated)

### Touchpoints
- Coverage matrix per framework — collapsible rows per framework, explicit gap list under every Partial cell
- Never show a green check when coverage is partial; always render the percentage and the explicit missing sub-requirements

### Thoughts & Emotions
- **Empowered** — the ClientAdmin knows, for this one artifact, exactly what they've bought and what they still need
- **Honest** — "Partial with gap list" beats "green check that falls apart at audit time"

### Pain Points
- Framework version mismatch (evidence evaluated against 2013 mapping when the engagement is on 2022) — must be surfaced explicitly
- Partial-satisfaction calculations need to be transparent — not a black-box score

### Opportunities
- Click-through from any framework-requirement cell to the crosswalk evidence: "Why is this only 67% coverage?"
- Pre-built "what's next" recommendations: "To close the PCI 8.3 gap, upload a remote-access MFA configuration export"

---

## Stage 5: Review Period-Aware Coverage and Staleness

### Sub-goal
Confirm the evidence is current enough for each target framework's window.

### User Actions
- For each Full / Partial cell, the system also evaluates the evidence's effective-date window against that framework's age tolerance (ASV scan 90 days for PCI, pen test 1 year, background check 1 year, SOC 2 Type 2 continuous-coverage window, ISO surveillance interval, etc.)
- Cells render a second-line status: "Current" / "Expires in 14 days" / "Stale — re-collect before report issuance"
- ClientAdmin can set a freshness alert to fire before expiry — this feeds Journey 12 (Continuous Assurance)

### Touchpoints
- Freshness indicator on every coverage cell
- "Remind me before expiry" toggle per evidence item

### Thoughts & Emotions
- **Reassured** — a green + "Current" cell is genuinely green; staleness doesn't sneak up
- **Planning-oriented** — upcoming expiries are visible weeks in advance

### Pain Points
- Framework-and-control-specific age tolerances are genuinely complex; the UI must not hide them behind a single "freshness" score

### Opportunities
- Show the exact rule being applied: "PCI DSS 4.0.1 requires ASV scans within 90 days of report"
- Bulk freshness view at the workspace level (this also powers the evidence freshness dashboard extension of the Client Hub)

---

## Stage 6: Record the AIDecision and Push Into Engagement Evidence Pools

### Sub-goal
Make the evidence available to every in-flight auditor engagement the artifact is relevant to, without the ClientAdmin having to think about which engagement is which.

### User Actions
- On final confirmation, the system writes AIDecision records for every accept / modify / reject action taken across Stages 3–5
- The `EvidenceItem` is linked to every active engagement whose scope includes a framework requirement this artifact covers — auditors see it appear in their engagement evidence pool with the cross-framework coverage card already populated (Journey 5 stage 3 consumes this)
- A summary card is shown to the ClientAdmin: "This artifact now satisfies 6 requirements across 4 frameworks in 2 active engagements. 3 gaps remain; suggested evidence to close them: [list]."

### Touchpoints
- AIDecision writeback confirmation
- "Pushed to N engagements" summary card
- Gap-closing suggestion list (feeds future uploads and the auditee GRC workspace to-do)

### Thoughts & Emotions
- **Accomplished** — one upload cleared real work across multiple engagements
- **Forward-looking** — the gap list is a concrete to-do, not a vague alarm

### Pain Points
- If the push-to-engagement step fails silently, the auditor never sees the evidence and the ClientAdmin thinks the work is done

### Opportunities
- Post-push verification: "Confirmed visible in engagement [SOC 2 FY26] and engagement [ISO 27001 surveillance 2026]"
- Nudge on the top gap: "Your biggest open gap is PCI 8.3 (remote-access MFA). Upload a remote-access MFA export next."

---

## Handoffs

| From | To | Information Transferred | Trigger |
|------|-----|------------------------|---------|
| ClientAdmin | Staff Auditor (Journey 5, 7) | EvidenceItem with confirmed CommonControl edges and per-framework coverage | Mapping confirmed |
| ClientAdmin | ClientAdmin (Journey 12) | Freshness alerts for evidence approaching staleness | "Remind me before expiry" set |

---

## Journey Summary

### Emotional Arc
Opens with the simple act of upload ("I have an MFA screenshot and I want to use it"), climbs through the reassuring classification step ("Axiom understands what this is"), peaks at the coverage matrix ("one artifact covers 6 requirements across 4 frameworks — and here are the 3 gaps, explicitly"), and resolves with the satisfying realization that this single interaction propagated value across every active engagement the artifact touches. The emotional through-line is "I can see exactly what I bought, and I can trust what I see."

### Cross-Cutting Pain Points
- Suggestions that aren't grounded in a licensed crosswalk erode trust instantly
- Green checkmarks on partial coverage are a product-level integrity failure — always show percentages and explicit gap lists
- Period-aware staleness and framework-specific age tolerances are invisible to most competitors; Axiom must surface them without overwhelming the user

### Prioritized Opportunities
1. **Crosswalk-grounded AI suggestions with explicit link-backs to SCF / OSCAL / AICPA rows** (high impact, medium effort) — post–Delve-scandal trust differentiator
2. **Partial-satisfaction gap lists rendered inline** (high impact, medium effort) — the #1 failure mode of every competitor's cross-mapping UI
3. **Period-aware freshness surfaced per framework per requirement** (high impact, medium effort) — where Vanta/Drata/Secureframe most often hand-wave
4. **Push-to-engagement with post-push verification** (high impact, low effort) — closes the auditor ↔ auditee loop that single-sided competitors structurally cannot

---

# Journey 5: Test Controls and Prepare Workpapers

## Overview
- **Persona:** Staff Auditor — assigned to specific controls within a compliance or assurance engagement (SOC 2, ISO 27001, ISO 27701, ISO 42001, HIPAA, PCI DSS, or SOC 1)
- **Goal:** Complete all assigned control tests with linked evidence and draft workpapers ready for manager review
- **Trigger:** Engagement is in Fieldwork status and controls are assigned via Control.auditor_assigned_to_id
- **Stages:**
  1. Review assigned controls and test procedures
  2. Execute test procedures and document results
  3. Link evidence to test procedures
  4. Request AI workpaper narrative draft
  5. Edit and finalize workpaper content
  6. Submit workpaper for review

## Stage 1: Review Assigned Controls and Test Procedures

### Sub-goal
Understand what needs to be tested, the expected results, and what evidence is required.

### User Actions
- Opens the engagement dashboard and navigates to "My Assignments"
- Reviews assigned Control records: description, control objective, framework requirements satisfied, prior year status (for rollforward)
- For each control, reviews associated TestProcedure records: procedure type (Inquiry, Observation, InspectionOfDocument, Reperformance, Analytics), description, expected result
- Notes which test procedures require evidence and what type

### Touchpoints
- "My Assignments" panel — personalized view of assigned controls
- Control detail view with linked framework requirements, test procedures, and prior year reference
- For multi-framework engagements: cross-framework display showing which other requirements this control satisfies

### Thoughts & Emotions
- **Oriented** — the assignments view gives clear scope of work
- **Informed** — seeing expected results and prior year context reduces ambiguity about what "done" looks like
- **Motivated** (multi-framework) — "Testing this one control covers SOC 2 CC6.1 AND ISO 27001 A.8.3. Two birds, one stone."

### Pain Points
- **Competitor context:** In Hyperproof and AuditBoard, staff must navigate workspace hierarchies to find assigned tasks — there's no personalized assignment view. Agentive has assignment workflows but no framework-aware control context. Drata / Vanta / Secureframe don't model auditor assignments at all.
- If the control description or expected result is too vague, the auditor doesn't know what "sufficient" looks like
- Rollforward prior year context must be reference-only — it shouldn't pre-fill current year conclusions

### Opportunities
- Show estimated time per control based on procedure types and evidence requirements
- For rollforward: highlight what changed from prior year ("Test procedure updated: added population size requirement")

---

## Stage 2: Execute Test Procedures and Document Results

### Sub-goal
Perform each test procedure and record the results, including any exceptions.

### User Actions
- Opens a TestProcedure record and changes status from NotStarted to InProgress
- Selects the procedure type (if not pre-set from template)
- Documents the population and sample where applicable (all seven frameworks involve sampling at some point — SOC 2 Type II period samples, ISO 27001 Annex A control sampling, PCI DSS population-level scan evidence)
- Records results: what was observed, measured, or confirmed
- Notes any exceptions found
- Documents conclusions per procedure
- Changes status to Complete (or Exception if exceptions found)

### Touchpoints
- Test procedure detail view with structured fields
- Population/sample reference for the relevant framework's sampling expectation
- Results entry: structured fields + narrative area
- Exception documentation panel
- Status progression: NotStarted → InProgress → Complete / Exception

### Thoughts & Emotions
- **Focused** — this is the core audit work; the auditor is in professional execution mode
- **Methodical** — structured fields guide the documentation without being overly rigid
- **Concerned** (when exceptions found) — "I need to document this carefully and escalate if significant"

### Pain Points
- **Competitor context:** Agentive's agentic workpaper drafting automates parts of control-test documentation but lacks framework-aware control context. Delve claims continuous monitoring across 25+ frameworks but the March 2026 scandal showed the "testing" was largely boilerplate. Hyperproof and AuditBoard document tests manually but structurally.
- Procedure documentation that is too structured (all dropdowns, no narrative) doesn't capture the auditor's professional judgment
- Procedure documentation that is too freeform (blank text box) leads to inconsistent quality across the team

### Opportunities
- Balanced structure: required fields (procedure type, population, sample, result, exceptions) plus narrative area for professional judgment
- Exception escalation: when exceptions are noted, prompt "Notify manager?" with a one-click notification
- Agentic autonomous control testing (Axiom differentiator): for procedures where the population is log-stream or configuration data (SOC 2 CC6 access logs, PCI DSS §10 logging requirements, ISO 27001 A.8.15 logging, ISO 42001 AI system monitoring), offer "Analyze full population" alongside traditional sampling — evidence collected once, tested continuously, with every decision written to the AIDecision ledger

---

## Stage 3: Link Evidence to Test Procedures

### Sub-goal
Connect supporting documents to test procedures to demonstrate that conclusions are evidence-based.

### User Actions
- From the test procedure view, opens the evidence linking panel
- Browses the engagement's evidence pool (all EvidenceItem records for this client, across all engagements)
- AI may have already suggested evidence links (EvidenceLink.ai_suggested = true) — the auditor accepts, modifies, or rejects
- Links selected evidence items to the test procedure
- For multi-framework engagements: the UI shows which other `FrameworkRequirement` nodes that evidence simultaneously satisfies (consuming the coverage data already recorded in Journey 4)

### Touchpoints
- Evidence linking panel with search and browse
- AI-suggested links with confidence scores and "Accept / Modify / Reject" actions
- Cross-framework satisfaction display: "This evidence satisfies SOC 2 CC6.1 (Full), ISO 27001 A.5.15 (Full), HIPAA §164.312(a)(1) (Full), PCI DSS 8.3 (Partial — remote-access MFA still required)"
- EvidenceLink record created with linked_by_id and timestamp
- AIDecision record updated on acceptance/modification/rejection

### Thoughts & Emotions
- **Efficient** — AI suggestions mean the auditor doesn't have to search through hundreds of evidence items manually
- **Delighted** (multi-framework) — seeing one piece of evidence satisfy three or four frameworks is the payoff of the cross-framework architecture
- **Professionally responsible** — accepting an AI suggestion still requires the auditor to verify that the evidence is relevant and sufficient

### Pain Points
- **Competitor context:** Hyperproof and AuditBoard have evidence linking but poor split-view — auditors can't see evidence and workpapers side-by-side. Drata / Vanta / Secureframe track evidence but don't model auditor test procedures. Agentive's workpaper drafting references evidence but not with framework-aware coverage indicators.
- If the evidence pool is large (hundreds of items across multiple engagements), search and filtering are critical
- AI suggestions that are wrong erode trust in the system — the auditor starts ignoring suggestions

### Opportunities
- Split-view: evidence document viewer on one side, test procedure on the other (addresses the top Hyperproof / AuditBoard pain point)
- Evidence search by content: "Find evidence mentioning 'access control policy'" using extracted text
- Show evidence reuse: "This document was used for 3 other controls in this engagement. First linked by [name] on [date]."
- Prior year evidence flagging: "This evidence was linked to the same control last year"

---

## Stage 4: Request AI Workpaper Narrative Draft

### Sub-goal
Get a first draft of the workpaper narrative that the auditor would otherwise write from scratch.

### User Actions
- After marking a test procedure as Complete, clicks "Generate AI Draft"
- This is an explicit, on-demand action — never automatic
- AI generates a narrative draft using: control description, test procedure details, linked evidence (extracted text), exceptions noted, prior year workpaper narrative (for rollforward), and firm workpaper template
- Draft appears in the workpaper editor labeled "AI Draft — requires review"
- A WorkpaperVersion record is created with is_ai_draft = true

### Touchpoints
- "Generate AI Draft" button on completed test procedures
- AI draft label: "AI Draft — requires review" banner in the workpaper editor
- WorkpaperVersion record with is_ai_draft = true

### Thoughts & Emotions
- **Hopeful** — "Let's see if the AI can capture what I would have written"
- **Critical** — reading the draft with a professional editor's eye
- **Relieved** — even a 70% accurate draft saves 30–60 minutes of blank-page writing time

### Pain Points
- **Competitor context:** Agentive's narrow agentic surface drafts workpapers per request with reported time savings. Hyperproof and AuditBoard have no AI drafting. Drata / Vanta / Secureframe don't author auditor workpapers at all — they're on the other side of the table.
- If the AI draft is too generic (boilerplate language that doesn't reflect the specific evidence reviewed), the auditor spends as much time editing as they would writing from scratch
- The is_ai_draft flag must be clearly explained — the auditor needs to understand they MUST edit before the workpaper can advance

### Opportunities
- Use firm-specific style reference: learn from prior year workpapers at the same firm to match tone, terminology, and structure
- Show a diff view when the auditor edits: "You changed 3 of 8 paragraphs. is_ai_draft will be set to false when you save."
- Generate section-by-section rather than all-at-once, so the auditor can review and edit incrementally

---

## Stage 5: Edit and Finalize Workpaper Content

### Sub-goal
Transform the AI draft (or blank workpaper) into a complete, auditor-authored workpaper.

### User Actions
- Edits the workpaper content in the platform editor — the AI draft becomes the starting point
- Adds professional judgment, specific findings, and conclusions that the AI cannot generate
- Saves edits — WorkpaperVersion created with is_ai_draft = false (once any human edit is made)
- Reviews the complete workpaper: test procedure results, linked evidence, narrative, conclusions
- Ensures all content is accurate and reflects their professional work

### Touchpoints
- Workpaper editor (rich text with structured sections)
- Version history panel showing all WorkpaperVersion records
- Evidence reference sidebar (linked evidence items visible while editing)
- Status indicator: Draft → ready for PreparedPendingReview

### Thoughts & Emotions
- **Focused** — editing requires concentrated professional judgment
- **Ownership** — the workpaper carries their name; the auditor takes personal responsibility for accuracy
- **Satisfied** — the AI draft accelerated the work but the final product is genuinely theirs

### Pain Points
- **Competitor context:** Agentive supports real-time collaborative editing. Hyperproof and AuditBoard allow simultaneous edits but with varying conflict-resolution quality. Drata / Vanta / Secureframe don't have a workpaper editor.
- If the editor is clunky (poor formatting, no undo, slow saves), it undermines the entire AI-assisted workflow
- Version history must be clear — the auditor needs to see what changed between versions, especially to demonstrate human edits on AI drafts

### Opportunities
- Real-time collaborative editing (not check-in/check-out) for team workpapers
- Auto-save with version history — no explicit "save" action needed, every change creates a recoverable version
- Evidence reference sidebar that updates contextually based on cursor position in the workpaper

---

## Stage 6: Submit Workpaper for Review

### Sub-goal
Record the Tester sign-off on the workpaper and submit it into the four-level review chain (Detailed Reviewer → General Reviewer → Final Reviewer).

### User Actions
- Reviews the workpaper one final time
- Clicks "Submit for Review" — system records a `WorkpaperSignOff` row at `reviewer_level = Tester`, status changes to `PreparedPendingReview`, `current_reviewer_level` advances to `DetailedReviewer`
- System validates: is_ai_draft must be false (if an AI draft was generated, the auditor must have edited it)
- The assigned Detailed Reviewer (typically Manager) receives a notification
- The workpaper is now locked for the preparer — only the reviewer at the current level can make changes or return it

### Touchpoints
- "Submit for Review" button with validation check
- Validation error if is_ai_draft = true: "This workpaper contains an unedited AI draft. You must review and edit the AI-generated content before submitting."
- Notification to the assigned Detailed Reviewer
- Workpaper status badge updates with current reviewer level
- Sign-off ledger panel: "Tester sign-off recorded — [name] at [timestamp]"

### Thoughts & Emotions
- **Accomplished** — a completed workpaper represents meaningful professional output
- **Hopeful** — "I hope the reviewers don't have too many review notes"
- **Trusting** — the system validation (AI draft check) gives the auditor confidence that they're not submitting incomplete work

### Pain Points
- **Competitor context:** Agentive, Hyperproof, and AuditBoard have digital sign-off workflows but enforcement varies by firm configuration. None enforces a four-level reviewer chain (Tester → Detailed → General → Final) at the data layer. None enforces the AI-draft-must-be-edited gate at the data layer — together these are core Axiom differentiators grounded in ISO 42001 HITL discipline, AICPA SQMS 1 / ISO 17021-1 firm quality frameworks, and the post–Delve-scandal provenance bet.
- If review at any level takes days (a common audit bottleneck), the auditor's momentum stalls
- The "locked for preparer" state must be communicated clearly — the auditor needs to know they can't make further edits until a reviewer returns the workpaper for rework
- A workpaper returned from a higher level (e.g., General Reviewer raising notes after the Detailed Reviewer signed off) supersedes lower sign-offs, which can feel demoralizing — UI must communicate the workflow rationale

### Opportunities
- Show estimated review queue time per reviewer level: "Detailed Reviewer [name] has 4 workpapers ahead of yours in the review queue"
- Allow the auditor to add notes for the reviewer chain: "FYI — this is a new control this year. Prior year test procedure was different."
- Sign-off ledger panel showing each reviewer level (Tester, Detailed, General, Final) with status, signer name, and timestamp; supersession history shown inline so the workflow is fully transparent

---

## Handoffs

| From | To | Information Transferred | Trigger |
|------|-----|------------------------|---------|
| Staff Auditor | Manager (Journey 6) | Completed workpaper with linked evidence | Submit for Review |
| Staff Auditor | Staff Auditor (Journey 7) | Evidence needs identified during testing | Evidence gaps found |

---

## Journey Summary

### Emotional Arc
Begins with orientation and planning (reviewing assignments), builds through the focused execution of testing (the core professional work), peaks at the AI-assisted drafting moment ("a first draft in 30 seconds"), sustains through the careful editing phase (professional ownership), and resolves with the satisfying submission for review. The emotional through-line is "I'm doing real audit work, and the AI handles the tedium."

### Cross-Cutting Pain Points
- Evidence linking is the critical integration point — if it's clunky, the entire testing workflow suffers
- The AI draft quality determines whether the feature saves time or wastes it
- The gap between submission and review is the biggest bottleneck in audit workflows industry-wide

### Prioritized Opportunities
1. **Split-view evidence + workpaper** (high impact, medium effort) — addresses the top Hyperproof / AuditBoard pain point
2. **Cross-framework evidence satisfaction display** (high impact, already designed) — the core architectural differentiator experienced at the task level
3. **Firm-specific AI draft style matching** (high impact, high effort) — the difference between a generic draft and a useful one
4. **Review queue visibility** (medium impact, low effort) — transparency reduces the waiting frustration

---

# Journey 6: Review Workpapers and Advance the Engagement

## Overview
- **Persona:** Reviewer at one of three levels in the four-level sign-off hierarchy — **Detailed Reviewer** (typically Manager, 5–8 years experience), **General Reviewer** (Manager or Partner, depending on firm policy and engagement risk), or **Final Reviewer** (typically Partner). The same person may serve different levels on different workpapers within the same engagement; firm policy maps each `reviewer_level` to eligible `UserRole` values.
- **Goal:** Review all submitted workpapers at the assigned level, provide feedback via review notes, sign off when satisfied, and (for the engagement-managing reviewer) advance the engagement from Fieldwork through Review to Reporting
- **Trigger:** A workpaper enters the reviewer's level — `current_reviewer_level` matches the reviewer's assignment (Detailed at submit, General after Detailed sign-off, Final after General sign-off)
- **Stages:**
  1. Triage the review queue (filtered by reviewer level assignment)
  2. Review workpaper content and evidence
  3. Raise and manage review notes
  4. Sign off at the current reviewer level (advances workpaper to the next level, or to ReviewComplete after Final)
  5. Advance engagement status (engagement-managing reviewer only)

## Stage 1: Triage the Review Queue

### Sub-goal
Prioritize which workpapers to review first across multiple engagements and team members.

### User Actions
- Opens the review dashboard showing all workpapers in PreparedPendingReview status across assigned engagements
- Sorts and filters by: engagement, control area, staff auditor, submission date, priority
- Selects a workpaper to begin review

### Touchpoints
- Review dashboard — cross-engagement view of all pending reviews
- Engagement-level progress indicators (e.g., "34 of 89 controls reviewed")
- Staff auditor notes visible in the queue: "FYI — new control this year"

### Thoughts & Emotions
- **Busy** — multiple engagements, multiple staff submitting simultaneously during busy season
- **Strategic** — prioritizing based on engagement deadlines and control risk
- **Responsible** — every review carries their professional endorsement

### Pain Points
- **Competitor context:** Agentive's manager dashboard offers per-request visibility but lacks cross-engagement review queues. Hyperproof and AuditBoard review workflows are per-document with no cross-engagement queue. Drata / Vanta don't model a reviewer role at all.
- If the review queue has no prioritization tools, the manager reviews in submission order rather than by importance
- Busy season creates a bottleneck where the review queue grows faster than the manager can clear it

### Opportunities
- Priority flags on high-risk controls or approaching-deadline engagements
- Batch review mode: open multiple workpapers in tabs for rapid sequential review
- Show reviewer workload: "You have 12 workpapers to review across 3 engagements. Estimated time: 4 hours."

---

## Stage 2: Review Workpaper Content and Evidence

### Sub-goal
Verify that the workpaper is complete, accurate, properly supported by evidence, and reflects sound professional judgment.

### User Actions
- Opens the workpaper in review mode (full edit access as reviewer)
- Reads the narrative — checks for accuracy, completeness, and professional quality
- Examines linked evidence items — opens each to verify relevance and sufficiency
- Checks test procedure results and conclusions
- For AI-drafted content: reviews WorkpaperVersion history to confirm the auditor made substantive edits (is_ai_draft was set to false)
- Checks cross-framework evidence satisfaction (for multi-framework engagements)

### Touchpoints
- Workpaper review view with evidence sidebar
- Version history panel showing all edits, including the AI draft → human edit transition
- Evidence viewer (inline preview without leaving the workpaper)
- Cross-framework indicator: "This evidence satisfies 3 framework requirements"

### Thoughts & Emotions
- **Evaluative** — reading with a critical eye, checking both substance and form
- **Mentoring** — identifying areas where the staff auditor can improve
- **Efficient** — needs to review quickly without sacrificing thoroughness

### Pain Points
- **Competitor context:** Hyperproof and AuditBoard lack true split-view for evidence alongside workpapers — reviewers switch between tabs. Agentive is adding split-view but without framework-aware coverage metadata. Most GRC tools require navigating between the workpaper and evidence files in separate panes.
- If the evidence viewer is a separate window or requires navigation, the review takes twice as long
- Reviewing AI-drafted content requires additional attention — the manager must ensure the auditor actually applied judgment, not just accepted the AI output

### Opportunities
- Split-view review mode: workpaper on the left, evidence on the right, with framework-coverage badges inline
- AI-edit indicator: "This workpaper was AI-drafted. The auditor modified 5 of 8 sections." — helps the reviewer focus their attention
- Inline annotations: the reviewer can highlight text and add comments without opening a separate review notes panel

---

## Stage 3: Raise and Manage Review Notes

### Sub-goal
Document specific feedback, questions, or required changes as structured review notes that the staff auditor must address.

### User Actions
- Creates a review note linked to a specific section of the workpaper: selects text → adds comment
- Each note has a description, severity (question, suggestion, required change), and linked content
- Open review notes block workpaper advancement — the auditor must address each one
- When the auditor responds, the manager reviews and resolves or escalates
- Resolved notes remain in the record — they cannot be deleted (audit-documentation integrity applies equally across SOC 2 AICPA AT-C, ISO 17021-1 audit records, HIPAA §164.312(b), and PCI DSS §12.10)

### Touchpoints
- Inline review note creation (select text → comment)
- Review notes panel: list of all notes with open/resolved status
- Notification to staff auditor when notes are added
- Resolution workflow: auditor responds → manager resolves
- AuditLog entries for note creation, response, and resolution

### Thoughts & Emotions
- **Constructive** — the goal is improving the work, not finding fault
- **Frustrated** (if recurring issues) — the same error across multiple workpapers from the same auditor is a training signal
- **Satisfied** — when notes are addressed well and the workpaper improves

### Pain Points
- **Competitor context:** Hyperproof and AuditBoard track review notes linked to controls. Agentive has inline commenting. All support open/resolved status tracking. None prevent note deletion — Axiom's immutability is a compliance differentiator.
- If review notes are too informal (just text comments), they lack the structure needed for regulatory documentation
- The back-and-forth cycle (note → response → resolution) can be slow if notifications are delayed

### Opportunities
- Note categorization for pattern detection: "This manager raised 14 'missing evidence description' notes this quarter — consider adding to the staff training curriculum"
- One-click note response templates: "Evidence updated" / "Narrative revised" / "Disagree — see response below"
- Real-time notification when review notes are added or responded to

---

## Stage 4: Sign Off at the Current Reviewer Level

### Sub-goal
Record the reviewer's sign-off (`WorkpaperSignOff` row at the current `reviewer_level`) and advance the workpaper to the next level — or to `ReviewComplete` after the Final Reviewer signs off.

### User Actions
- Confirms all review notes raised at the current level are resolved (open notes at this level block advancement)
- Clicks "Sign off at [level]" — system creates a `WorkpaperSignOff` row at the current `reviewer_level` and advances `current_reviewer_level` to the next position:
  - Detailed Reviewer signs → workpaper moves to General Reviewer (`GeneralReviewInProgress`)
  - General Reviewer signs → workpaper moves to Final Reviewer (`FinalReviewInProgress`)
  - Final Reviewer signs → workpaper status becomes `ReviewComplete` then `SignedOff`; `signed_off_by_id` populated
- The sign-off action creates a timestamped, named AuditLog entry — it cannot be backdated
- The next-level reviewer (or, after Final, the partner managing the engagement) is notified

### Touchpoints
- "Sign off at [level]" button (disabled if open notes at this level exist)
- Sign-off confirmation with name, level, and timestamp
- Notification to next-level reviewer
- AuditLog entry: "Workpaper [name] signed off at [level] by [Reviewer] at [timestamp]"
- Sign-off ledger panel showing all four levels with status (Pending / Signed / Superseded), signer name, and timestamp

### Thoughts & Emotions
- **Confident** — the workpaper meets professional standards at this review level
- **Relieved** — one more item cleared from the review queue
- **Professional pride** — the sign-off at this level carries the reviewer's reputation

### Pain Points
- **Competitor context:** Hyperproof and AuditBoard have configurable sign-off schemes but enforcement varies. Agentive has digital sign-off. **No competitor enforces a four-level reviewer chain (Tester → Detailed → General → Final) at the data layer** with strict ordering, supersession on rework, and per-level eligibility validation — this is grounded in SQMS 1 / ISO 17021-1 / internal firm quality frameworks and is a core Axiom differentiator.
- If sign-off at any later level takes additional days, the prior reviewer's work sits idle
- Bulk sign-off scenarios (many workpapers completing simultaneously) need efficiency without sacrificing individual attention
- Supersession when a higher level returns the workpaper can be confusing — the UI must clearly indicate which prior sign-offs were invalidated and why

### Opportunities
- Per-level review queue mirroring the prior-level review queue — cross-engagement visibility for each reviewer level
- Allow Final Reviewers to batch-sign low-risk workpapers with individual confirmation: "Sign off on 8 workpapers at Final level? Each will be individually timestamped."
- Visual sign-off ladder showing the four-level chain with completion status — helps reviewers understand workflow position at a glance

---

## Stage 5: Advance Engagement Status

### Sub-goal
Move the engagement from Fieldwork to Review, and from Review to Reporting, once all controls are complete.

### User Actions
- Reviews the engagement progress dashboard: all Control records must have status Complete or Exception
- For Fieldwork → Review: confirms all controls have results
- For Review → Reporting: all review notes resolved AND EngagementQualityReview.status = Complete where applicable
- Each transition has a guard condition enforced at the data layer

### Touchpoints
- Engagement progress dashboard with control status summary
- Phase transition button with guard validation
- Transition confirmation: "Moving to Review. 89 controls complete, 2 with exceptions. All review notes resolved."
- Team notifications on phase transition

### Thoughts & Emotions
- **Accomplished** — advancing the engagement represents major project milestones
- **Careful** — the manager verifies that nothing was missed before advancing
- **Collaborative** — the transition affects the entire engagement team

### Pain Points
- If one or two controls are blocking the transition (not yet complete), it's unclear who is responsible and what's left
- The Fieldwork → Review transition may be blocked by a single outstanding control that was assigned to a staff auditor who is out

### Opportunities
- Blocking item dashboard: "2 controls blocking Fieldwork → Review transition. CC7.2 assigned to [name], due [date]. CC8.1 assigned to [name], due [date]."
- Allow Manager to reassign blocking controls to available staff
- Show a pre-transition checklist summarizing what was accomplished and what exceptions exist

---

## Handoffs

| From | To | Information Transferred | Trigger |
|------|-----|------------------------|---------|
| Detailed / General / Final Reviewer | Tester (Journey 5) | Review notes requiring response | Notes raised at any level |
| Detailed Reviewer | General Reviewer (this journey, next level) | Workpaper signed off at Detailed level | DetailedReviewer sign-off recorded |
| General Reviewer | Final Reviewer (this journey, next level) | Workpaper signed off at General level | GeneralReviewer sign-off recorded |
| Final Reviewer | Partner (Journey 9) | Workpaper SignedOff; report assembly can proceed | All four levels signed; engagement reaches ReviewComplete |
| Engagement-managing Partner | EQR Reviewer (Journey 10) | Engagement ready for quality review | Review phase reached |

---

## Journey Summary

### Emotional Arc
Starts with the pressure of a growing review queue (busy season reality), moves through focused evaluation (professional judgment), includes the constructive cycle of review notes (mentoring and quality assurance), and resolves with the satisfaction of clearing workpapers and advancing engagement phases. The emotional low point is the review backlog; the high point is the phase transition.

### Cross-Cutting Pain Points
- The review bottleneck is the single biggest workflow problem in audit — managers are always the constraint
- Evidence review alongside workpaper content must be seamless, not a multi-window juggling act
- Review note management can become tedious if the tooling doesn't streamline the back-and-forth

### Prioritized Opportunities
1. **Split-view review mode** (high impact, medium effort) — the reviewer's version of the split-view evidence panel
2. **Blocking item dashboard** (high impact, low effort) — makes phase transitions transparent and actionable
3. **Review note pattern detection** (medium impact, medium effort) — transforms review data into training insights
4. **Batch sign-off with individual confirmation** (medium impact, low effort) — addresses busy-season efficiency without sacrificing compliance

---

# Journey 7: Manage Document Requests and Collect Evidence

## Overview
- **Persona:** Staff Auditor or Manager — responsible for PBC (Provided By Client) document collection
- **Goal:** Create document requests, send them to the client, track fulfillment, review AI-assessed uploads, and accept evidence into the engagement
- **Trigger:** Engagement is in Fieldwork and control testing requires client-provided evidence
- **Stages:**
  1. Create document requests
  2. Send requests to client
  3. Monitor status and manage reminders
  4. Review AI-assessed uploads
  5. Accept evidence and link to test procedures

## Stage 1: Create Document Requests

### Sub-goal
Define exactly what documents the client needs to provide, linked to specific controls and test procedures.

### User Actions
- Creates DocumentRequest records linked to specific controls or test procedures
- Each request includes: title, detailed instructions (what to provide, format required, period to cover), due date, and assignment to a client contact
- Bulk request creation from methodology template is available — the standard SOC 2 Type II template creates 80+ pre-drafted requests covering all trust services criteria
- Reviews and customizes auto-generated request descriptions as needed

### Touchpoints
- Document request creation form (individual or bulk from template)
- Request list with status overview
- Template-based bulk creation: "Generate SOC 2 Type II requests? This will create 83 requests across 5 trust services categories."

### Thoughts & Emotions
- **Efficient** — bulk creation from templates saves hours of manual request writing
- **Thorough** — wants to make sure every required document is requested upfront to avoid back-and-forth
- **Empathetic** — writing clear instructions because confusing requests lead to wrong uploads

### Pain Points
- **Competitor context:** Agentive auto-generates PBC request lists from engagement context — practitioners approve before sending. Hyperproof and AuditBoard have request management as core features. Drata and Vanta have auditee-initiated evidence request tracking but not auditor-initiated PBC workflow.
- If requests are too vague, clients upload wrong documents and the cycle repeats
- If bulk creation generates too many requests without customization, the client is overwhelmed

### Opportunities
- AI-assisted request descriptions: generate specific, clear instructions based on the control description and framework requirement
- Request grouping by client contact: "These 12 requests are all for the IT department — assign to [IT Contact]"
- Prior year request reuse for rollforward: "Last year you requested 78 documents. Carry forward and update for the current period?"

---

## Stage 2: Send Requests to Client

### Sub-goal
Deliver requests to the client and give them access to the upload portal.

### User Actions
- Reviews the request list and clicks "Send to Client" (individual or batch)
- System sends the client contact an email with a tokenized Client Hub link scoped to this engagement
- No client account or password is required for basic uploads — the tokenized link is sufficient
- The requests appear in the Client Hub immediately (Journey 8 continues from the client's perspective)

### Touchpoints
- Send confirmation: "15 requests sent to [client email]. Token valid for 90 days."
- Email delivery to client contact with tokenized link
- Client Hub populated with outstanding requests

### Thoughts & Emotions
- **Proactive** — getting requests out early reduces delays later
- **Confident** — the tokenized link system means the client doesn't need to create an account or remember a password

### Pain Points
- **Competitor context:** Agentive sends requests through a Client Hub. Hyperproof and AuditBoard include client portals. Drata / Vanta provide auditee-side portals but they're configured by the auditee, not sent by the auditor.
- If the client email goes to spam or the tokenized link is confusing, the client never accesses the portal
- Token expiry (90 days) must be managed — expired tokens need easy re-generation

### Opportunities
- Delivery confirmation: track whether the client opened the email and accessed the portal
- Tokenized link landing page with clear branding: "Your audit team at [Firm Name] has requested documents for your [Engagement Type] engagement"

---

## Stage 3: Monitor Status and Manage Reminders

### Sub-goal
Track which requests are fulfilled, which are overdue, and escalate as needed.

### User Actions
- Monitors the request dashboard: Pending, Submitted, InReview, Accepted, Rejected, Overdue status per request
- Reviews automated reminder schedule: default reminders at 7 days before due, on due date, 7 days after
- Customizes reminder frequency and content per engagement if needed
- After three automated reminders, receives an escalation notification from the system
- Can manually send follow-up messages for critical requests

### Touchpoints
- Request dashboard with status badges and due date indicators
- Automated reminder configuration (per-engagement customizable)
- Escalation notifications after 3 reminders
- Overdue count badge on engagement dashboard
- For Partners (Scale tier): overdue request rates across all active engagements in analytics dashboard

### Thoughts & Emotions
- **Monitoring** — checking status is a daily ritual during active engagements
- **Frustrated** — overdue requests are the most common engagement delay
- **Grateful** — automated reminders mean the auditor doesn't have to send manual follow-up emails

### Pain Points
- **Competitor context:** Agentive has automated reminders and real-time status tracking. Hyperproof and AuditBoard track request status with configurable notifications. Drata / Vanta reminder loops are auditee-side (reminders to the company's own staff to maintain evidence).
- Client non-responsiveness is the #1 cause of engagement delays across the industry
- If reminders are too frequent, clients feel harassed; if too infrequent, requests slip

### Opportunities
- Smart reminder escalation: increase urgency tone progressively rather than sending the same message three times
- Engagement readiness indicator: "72% of requests fulfilled. Estimated fieldwork start: [date] based on current fulfillment rate."
- Client contact delegation prompt: "These 5 overdue requests were sent to [email]. Would you like to try a different contact?"

---

## Stage 4: Review AI-Assessed Uploads

### Sub-goal
Review documents that the client uploaded and the AI has pre-assessed for completeness, before formally accepting them as evidence.

### User Actions
- When a client uploads a document, an async AI completeness review runs (Feature 1: Claude Sonnet analyzes the document against request requirements)
- The auditor receives a notification: "AI has reviewed [document name] for [request title]"
- Opens the review queue — sorted by AI confidence (low-confidence reviews surfaced first)
- For each upload, sees: AI recommendation (Accept / Request Clarification / Reject), specific gaps identified, confidence score
- Takes action: Accept, Send Back to Client (with auto-drafted explanation from AI gap analysis), or Reject

### Touchpoints
- AI review notification (in-app + email based on preferences)
- Evidence review queue with AI recommendations and confidence scores
- Document viewer with highlighted gaps
- One-click actions: Accept / Send Back / Reject
- AIDecision record created and updated with review_action
- AuditLog entries for each action

### Thoughts & Emotions
- **Guided** — the AI pre-assessment means the auditor isn't reviewing blind; they know what to look for
- **Efficient** — accepting clearly-complete documents is a one-click action
- **Professionally cautious** — low-confidence reviews require careful manual examination
- **Appreciative** — the auto-drafted "send back" explanation saves time writing rejection emails

### Pain Points
- **Competitor context:** Agentive's Request Analysis Agent validates documents on upload and flags gaps. Hyperproof and AuditBoard have AI-assisted pre-validation in newer releases. Drata / Vanta validate at the auditee's own evidence maintenance layer but don't assess submission against a specific auditor request.
- If the AI frequently misjudges completeness (too many false Accept or false Reject recommendations), the auditor stops trusting the queue
- The "Send Back" action must communicate clearly to the client what's missing — a generic "incomplete" message generates confusion and more back-and-forth
- Period-coverage checks vary by framework: SOC 2 Type II continuous window (AT-C 320), ISO 27001 surveillance visit window, PCI DSS ASV scan validity (90 days), HIPAA risk-analysis refresh cadence — every AI review must apply the right rule for the target framework

### Opportunities
- AI confidence calibration display: "This AI has reviewed 340 documents for your firm. Acceptance recommendation accuracy: 93%."
- Client-facing gap explanation auto-drafted by AI: "The uploaded access control policy covers January–September 2025 but the audit period extends to December 2025. Please provide an updated policy covering the full period."
- Batch review mode for high-confidence accepts: "AI recommends accepting 12 documents with >95% confidence. Review and batch accept?"

---

## Stage 5: Accept Evidence and Link to Test Procedures

### Sub-goal
Formally accept uploaded documents as engagement evidence and connect them to the relevant test procedures.

### User Actions
- On accepting a document, the system creates an EvidenceLink to the relevant TestProcedure
- AI also suggests which TestProcedure to link the evidence to (if not already linked via the DocumentRequest → Control chain)
- The DocumentRequest status changes to Accepted
- The evidence item is now available in the engagement's evidence pool for use across controls

### Touchpoints
- Accept action creates EvidenceLink automatically
- AI link suggestion with accept/modify interface
- Document request status update: Submitted → InReview → Accepted
- Evidence pool updated with new item
- Cross-framework display: accepted evidence automatically satisfies all mapped framework requirements

### Thoughts & Emotions
- **Closing the loop** — request created → sent → fulfilled → reviewed → accepted → linked. The cycle is complete.
- **Productive** — each accepted document moves the engagement forward
- **Organized** — everything is connected: request → evidence → test procedure → control → framework requirement

### Pain Points
- If evidence linking is manual after acceptance, the auditor does double work (accept the document, then go find the right test procedure to link it to)
- Rejected or sent-back documents restart the cycle — the auditor needs to track what's been sent back and why

### Opportunities
- Auto-link on acceptance: when a DocumentRequest is linked to a specific Control and TestProcedure, acceptance automatically creates the EvidenceLink
- "Sent back" tracking: show a history of back-and-forth per request so the client's second attempt addresses the right gaps
- Evidence reuse detection: "This document was uploaded for a different engagement. Do you want to link it here as well?"

---

## Handoffs

| From | To | Information Transferred | Trigger |
|------|-----|------------------------|---------|
| Staff Auditor | Client Contact (Journey 8) | Document requests with instructions and due dates | Requests sent |
| Client Contact (Journey 8) | Staff Auditor | Uploaded documents with AI assessment | Client uploads |
| Staff Auditor | Staff Auditor (Journey 5) | Accepted evidence linked to test procedures | Evidence accepted |

---

## Journey Summary

### Emotional Arc
Begins with the productive energy of request creation (bulk templates are satisfying), settles into the monitoring rhythm of waiting for client responses (the emotional low point — waiting is frustrating), picks up with AI-assisted review (the AI recommendation queue is engaging and efficient), and closes with the satisfying loop of acceptance and evidence linking. The PBC cycle is the most frustrating part of audit work industry-wide; the AI completeness review is the feature that most directly addresses the pain.

### Cross-Cutting Pain Points
- Client responsiveness is outside the platform's control but inside its influence (reminders, escalation, clear instructions)
- The AI completeness review quality determines whether PBC is a one-cycle or multi-cycle process
- Evidence linking must be automatic where possible — manual linking after acceptance is wasted effort

### Prioritized Opportunities
1. **AI-drafted client-facing gap explanations** (high impact, medium effort) — reduces the back-and-forth that causes engagement delays
2. **Auto-link evidence on acceptance** (high impact, low effort) — eliminates the manual step after acceptance
3. **Prior year request rollforward** (medium impact, low effort) — saves significant setup time for recurring engagements
4. **Batch review for high-confidence AI recommendations** (medium impact, low effort) — accelerates the review queue during bulk uploads

---

# Journey 8: Fulfill Audit Document Requests

## Overview
- **Persona:** Client Contact — security, compliance, IT, HR, or privacy lead at the audited company. Not an auditor. May have limited time and multiple people who hold different pieces of information.
- **Goal:** Understand what documents the audit team needs, upload them correctly the first time, and get through the PBC process with minimal disruption to daily work
- **Trigger:** Email notification: "Your audit team has started your [Engagement Name] and needs documents from you."
- **Note:** The Client Hub is also the entry point to the broader auditee GRC workspace — Journey 4 (cross-framework evidence mapping), Journey 12 (continuous assurance), evidence freshness dashboard, and policy library all extend from this same surface. A ClientAdmin user sees those extensions; a basic Client Contact typically only sees their PBC request queue.
- **Stages:**
  1. Access the Client Hub
  2. Review outstanding requests
  3. Upload documents
  4. Delegate specialized requests
  5. Track status and respond to follow-ups

## Stage 1: Access the Client Hub

### Sub-goal
Get into the upload portal without friction — no account creation, no password, no app installation.

### User Actions
- Opens the email from the audit firm
- Clicks the tokenized Client Hub link — no login required
- Lands on the Client Hub showing: engagement name and period, list of outstanding requests, and a drag-and-drop upload interface
- The tokenized link is scoped to this specific engagement — cannot access other engagements or firm content

### Touchpoints
- Email notification with prominent "View Your Requests" button
- Tokenized link landing page (no login screen)
- Client Hub dashboard: engagement name, period, request count, progress indicator

### Thoughts & Emotions
- **Relieved** — "No account? No password? I just click the link?"
- **Oriented** — the landing page clearly shows what's needed and how many items
- **Mildly stressed** — "I have 23 document requests and my day job to do"

### Pain Points
- **Competitor context:** Agentive's Client Hub requires the client to see a full dashboard. Hyperproof and AuditBoard ship client portals but require per-user account configuration. Drata / Vanta / Secureframe are the auditee's workspace — they don't ship a no-login auditor-to-auditee link at all.
- If the tokenized link is expired (>90 days), the client hits a dead end — they must contact the auditor to re-generate
- If the email looks like spam or the branding is unfamiliar, the client may not click

### Opportunities
- Clear audit firm branding on the email and landing page: "[Firm Name] logo + engagement details"
- Token expiry warning: "This link expires on [date]. Bookmark it for easy access."
- Optional account creation prompt for returning clients: "Create a password-protected account to access all your engagements in one place" — but never required

---

## Stage 2: Review Outstanding Requests

### Sub-goal
Understand what's needed, prioritize by due date, and identify what the client can provide themselves vs. what needs to be delegated.

### User Actions
- Scans the request list: each request shows title, instructions, due date, and status
- Reads the detailed instructions for each request to understand exactly what's needed
- Mentally categorizes: "I can provide this one," "IT needs to handle this," "I need to ask the CFO"
- If an intake form is enabled (common for first-year SOC 2 and ISO engagements), completes it before the request list appears

### Touchpoints
- Request list with title, instructions, due date, status
- Detailed instruction expandable per request
- Due date indicators (overdue in red, approaching in yellow)
- Optional intake form (entity legal name, business description, key systems, contacts)

### Thoughts & Emotions
- **Organized** — the request list is clear and structured
- **Overwhelmed** (if 80+ requests) — "This is a lot. Where do I start?"
- **Uncertain** — some requests may be unclear even with detailed instructions

### Pain Points
- **Competitor context:** Agentive organizes requests in a dashboard with progress tracking. Hyperproof and AuditBoard present requests with descriptions. Drata / Vanta surface evidence gaps in the auditee's own dashboard.
- Audit terminology in request descriptions may confuse non-auditor clients
- 80+ requests for a SOC 2 Type II is standard but psychologically overwhelming — the client needs prioritization guidance

### Opportunities
- Request grouping by topic or department: "IT Security (12 requests) | HR Policies (8 requests) | Financial (15 requests)"
- Priority indicators set by the auditor: "Start with these 5 high-priority requests"
- Plain-language descriptions: avoid "Provide evidence of logical access controls for production environment" — instead "Upload your company's policy document that shows who has access to production systems and how access is approved"

---

## Stage 3: Upload Documents

### Sub-goal
Provide the right documents in the right format.

### User Actions
- Selects a request from the list
- Uploads a file: drag-and-drop or file picker, single or bulk upload
- The upload is stored as an EvidenceItem immediately
- The request status changes to Submitted
- AI completeness review begins asynchronously in the background (the client may not know this is happening)

### Touchpoints
- Upload interface per request: drag-and-drop zone + file picker button
- Upload progress indicator
- Confirmation: "Document uploaded successfully. Your audit team will review it shortly."
- Request status changes to "Submitted" in the client's view

### Thoughts & Emotions
- **Accomplished** — each upload reduces the outstanding count
- **Hopeful** — "I hope this is what they wanted"
- **Efficient** — drag-and-drop is fast; no forms to fill out, no metadata to enter

### Pain Points
- **Competitor context:** All competitors support basic file upload. Agentive's Client Hub provides a clean upload experience. Drata / Vanta / Secureframe emphasize connected-source ingestion (SaaS integrations) over file upload.
- If the client uploads the wrong document, they don't know until the auditor reviews and sends it back — this delay is the core PBC problem
- Large file uploads (>100MB) may fail or time out
- Clients sometimes upload password-protected files without realizing the auditor can't open them

### Opportunities
- Real-time AI pre-check: immediately after upload, show the client a brief assessment: "This document appears to match the request. Your audit team will confirm." — manages expectations without making the client wait
- File format guidance: "Please upload in PDF, XLSX, or DOCX format. Password-protected files cannot be processed automatically."
- Upload progress persistence: if the upload fails, resume from where it left off rather than restarting

---

## Stage 4: Delegate Specialized Requests

### Sub-goal
Route specific requests to colleagues who have the relevant information, without giving them access to the full engagement.

### User Actions
- A ClientAdmin user (elevated role, assigned by the auditor) opens a request they can't fulfill themselves
- Enters a colleague's email address on that specific request
- The colleague receives a tokenized link scoped to that single request only
- The delegate sees only: request description, instructions, and upload interface — not the full Client Hub, engagement name, or other requests
- Delegation creates an AuditLog entry

### Touchpoints
- "Delegate" button on each request (visible only to ClientAdmin role)
- Email address input for the delegate
- Confirmation: "Request delegated to [email]. They can only see this specific request."
- Email to delegate with single-request tokenized link
- AuditLog: who delegated, to whom, when

### Thoughts & Emotions
- **Relieved** — "I don't have to chase my IT director for this. They can upload it directly."
- **In control** — delegation is specific and auditable; the client isn't sharing broad access
- **Efficient** — the delegate gets just enough context to fulfill the request

### Pain Points
- **Competitor context:** Agentive's Client Hub allows clients to delegate tasks to team members within the portal. Drata / Vanta have team invitations but scoped to the full workspace, not a single request. Axiom's single-request-scoped delegation is more restrictive but more secure.
- If the delegate doesn't respond, the ClientAdmin has no visibility into whether the email was received or the link was clicked
- Delegates who receive a narrow scoped link without context may be confused: "Why am I getting a link to upload a document?"

### Opportunities
- Delegation tracking: show the ClientAdmin whether the delegate has opened the link and uploaded anything
- Custom message on delegation: "Add a note for your colleague: 'Hi Maria, can you upload the Q4 access review report?'"
- Reminder capability: the ClientAdmin can re-send the delegation email

---

## Stage 5: Track Status and Respond to Follow-Ups

### Sub-goal
Monitor which requests are accepted, which need rework, and respond to auditor feedback.

### User Actions
- Returns to the Client Hub (via bookmarked link or new email notification)
- Sees request statuses: Submitted, Accepted, or Sent Back with explanation
- For "Sent Back" requests: reads the auditor's explanation (AI-drafted from gap analysis), understands what's missing, uploads a corrected document
- After the engagement is archived: the Client Hub becomes read-only with all submitted documents visible and an explanation that the engagement is complete

### Touchpoints
- Client Hub dashboard with status per request
- "Sent Back" notification with gap explanation
- Re-upload interface for corrected documents
- Post-archive read-only view with submitted documents and completion notice
- If the auditor shared the issued report: read-only report access in the Client Hub

### Thoughts & Emotions
- **Frustrated** (when sent back) — "I thought I uploaded the right thing." The clarity of the gap explanation determines whether frustration becomes action or abandonment.
- **Satisfied** (when accepted) — each accepted document is a closed loop
- **Relieved** (post-archive) — the engagement is done; access to the final report is a nice touch

### Pain Points
- **Competitor context:** Agentive allows clients to comment on report sections and collaborate in-platform. Hyperproof and AuditBoard provide request-status tracking. Drata / Vanta keep the client in their own workspace continuously — the engagement "ending" is a softer concept there.
- Vague "send back" explanations are the #1 source of client frustration in the PBC process
- Clients who need to upload a corrected version may not remember which version they uploaded originally

### Opportunities
- Show the original upload alongside the gap explanation: "You uploaded [filename]. The audit team identified these gaps: [specific list]"
- Accepted request celebration: "18 of 23 requests complete! Only 5 remaining."
- Post-engagement client satisfaction survey: one question embedded in the read-only Client Hub

---

## Handoffs

| From | To | Information Transferred | Trigger |
|------|-----|------------------------|---------|
| Client Contact | Staff Auditor (Journey 7) | Uploaded document + AI assessment | Document uploaded |
| ClientAdmin | Delegate (internal colleague) | Single-request scope + instructions | Delegation action |
| Staff Auditor (Journey 7) | Client Contact | Sent-back explanation + gap analysis | Request rejected / sent back |

---

## Journey Summary

### Emotional Arc
Starts with relief (no-login access), dips at the overwhelming request list (especially for 80+ request SOC 2 engagements), steadies through the rhythm of uploading documents, and fluctuates with sent-back requests (frustration → understanding → re-upload). The emotional high points are the first successful upload ("this is easy") and the final accepted request ("done!"). The emotional low is receiving a "sent back" notification with a vague explanation.

### Cross-Cutting Pain Points
- The client experience determines the audit timeline — frustrated clients are slow clients
- Audit terminology in request descriptions creates confusion for non-auditor clients
- The sent-back cycle is the most critical moment — clear explanations prevent multi-cycle delays

### Prioritized Opportunities
1. **AI-drafted plain-language gap explanations** (high impact, medium effort) — the single highest-impact feature for reducing PBC cycle time
2. **Request grouping by department/topic** (medium impact, low effort) — makes 80+ requests manageable
3. **Delegation tracking** (medium impact, low effort) — gives the ClientAdmin visibility into delegated requests
4. **Real-time upload pre-check** (medium impact, medium effort) — immediate feedback before the auditor reviews

---

# Journey 9: Generate Report, Finalize, and Archive

## Overview
- **Persona:** Partner — engagement partner responsible for the final deliverable, regulatory compliance, and engagement closure
- **Goal:** Produce the final engagement report, issue it, finalize the engagement file, and ensure proper archival per assembly deadline and retention requirements
- **Trigger:** Engagement reaches Reporting status (all review notes resolved, EQR signed off where applicable)
- **Stages:**
  1. Generate report from template
  2. Edit and iterate on draft
  3. Share draft with client (optional)
  4. Issue the final report
  5. Finalize the engagement
  6. Monitor automatic archival

## Stage 1: Generate Report from Template

### Sub-goal
Create a first draft of the engagement report populated with engagement data.

### User Actions
- Opens the Reporting section of the engagement
- Selects report type: SOC 2 Type I, SOC 2 Type II, SOC 1 Type I, SOC 1 Type II, ISO 27001 certificate support letter, ISO 27701 certificate support letter, ISO 42001 certificate support letter, **ISO Certificate (template draft)** for accredited Certification Body customers, HIPAA attestation letter (with optional HITRUST CSF r2 validated assessment report, post-MVP), PCI DSS Report on Compliance (ROC) + Attestation of Compliance (AOC), **PCI ROC (template draft)** and **PCI AOC (template draft)** for QSA firm customers, Agreed-Upon Procedures, or Management Letter. **Note:** the "(template draft)" report types produce the deliverable document for the licensed firm to review and sign — Axiom does not act as the issuing CB or QSA. Legal certification decisions and ROC/AOC sign-offs remain with the accredited firm under ISO 17021-1 / PCI SSC accreditation.
- The report template is pre-populated with: client name, audit period, framework, controls summary, exception summary, testing results
- A ReportVersion record is created for the initial draft
- Optionally, requests an AI draft of the "Description of Tests of Controls" section (Tier 2 AI — same rules as workpaper draft: labeled, human must edit and sign off)

### Touchpoints
- Report type selector
- Pre-populated report template with engagement data
- AI draft option for specific report sections
- Report editor with version tracking

### Thoughts & Emotions
- **Focused** — this is the final deliverable; quality and accuracy are paramount
- **Efficient** — pre-populated data means the partner writes narrative and opinion, not boilerplate
- **Authoritative** — the partner is exercising final professional judgment on every statement in the report

### Pain Points
- **Competitor context:** Agentive has one-click report generation from templates, but users report the templates can be finicky. Hyperproof and AuditBoard have report generation but require heavy template customization. Delve claims one-click SOC 2 report generation but the March 2026 scandal (493 of 494 reports shared identical boilerplate) is the live cautionary tale for that approach.
- If the template pre-population has errors (wrong period, missing exceptions), the partner must catch them — auto-populated content is trusted but must be verified
- Report templates that are too rigid don't accommodate firm-specific language and formatting

### Opportunities
- Customizable report templates per firm (headers, fonts, boilerplate language) that persist across engagements
- AI section drafting for standardized sections (Description of Tests of Controls, Scope and Approach) while leaving Opinion and Management Assertions to the partner
- Show a report completeness checklist: "All required sections populated. Exception summary includes 2 noted exceptions."

---

## Stage 2: Edit and Iterate on Draft

### Sub-goal
Refine the report narrative, add the professional opinion, and prepare for issuance.

### User Actions
- Edits the report in the platform editor
- Writes the audit opinion (or management assertion for SOC reports)
- Reviews the controls summary and exception handling sections
- Each save creates a new ReportVersion record
- Shares internally with other partners or the manager for feedback
- Iterates until satisfied

### Touchpoints
- Report editor with rich text, section navigation, and version history
- Internal sharing with commenting capability
- Version comparison view

### Thoughts & Emotions
- **Careful** — every word in the report has professional liability implications
- **Collaborative** — may consult with other partners on the opinion
- **Iterative** — multiple drafts are normal; the version history provides safety

### Pain Points
- **Competitor context:** Agentive allows report editing in Microsoft Word before delivery. Hyperproof and AuditBoard generate reports within their template system. None offer real-time collaborative report editing with cross-framework provenance baked in.
- Report editing is time-sensitive — the partner wants to issue quickly once testing is complete
- Formatting issues (headers, tables, pagination) can consume disproportionate time

### Opportunities
- Export to Word/PDF for offline review with tracked changes that can be re-imported
- Version diff: "Changes since last version: 3 paragraphs modified, exception summary updated"

---

## Stage 3: Share Draft with Client (Optional)

### Sub-goal
Get client feedback on the draft report before final issuance.

### User Actions
- Shares a draft report view with ClientAdmin users
- Clients see a read-only view of the draft and can submit comments
- Partner reviews client comments in-platform
- Issues a revised draft if needed

### Touchpoints
- "Share with Client" button creating a read-only client view
- Client commenting interface (inline comments)
- Comment review panel for the partner

### Thoughts & Emotions
- **Transparent** — sharing the draft builds client trust
- **Guarded** — the partner controls what the client sees and when
- **Responsive** — client comments are addressed promptly to avoid delays

### Pain Points
- **Competitor context:** Agentive allows clients to collaborate on report drafts directly in-platform. Hyperproof and AuditBoard have basic draft sharing capabilities.
- Client comments must be managed carefully — not all feedback should change the report
- Sharing too early (before the partner is confident in the draft) creates unnecessary revision cycles

### Opportunities
- Comment status tracking: resolved, acknowledged, will-not-change (with explanation)
- Selective sharing: share specific sections rather than the entire draft

---

## Stage 4: Issue the Final Report

### Sub-goal
Mark the report as issued, triggering regulatory compliance computations.

### User Actions
- Reviews the final report one last time
- Clicks "Issue Report" — Report.status = Issued, report_issued_at recorded on the Engagement
- System computes automatically: assembly_deadline (report date + 60 days for AICPA SOC engagements per AT-C; ISO engagements track to the certification body's file-closure requirements; PCI ROC/AOC assembly per PCI SSC), retention_deadline (report date + 5 years for AICPA SOC / ISO, + 6 years for HIPAA, + 3 years for PCI DSS ROC archives)
- Scheduling of the Finalized → Archived Step Functions wait state begins

### Touchpoints
- "Issue Report" button with confirmation dialog showing computed deadlines
- Deadline display: "Assembly deadline: [date]. Retention deadline: [date]."
- Issuance confirmation with AuditLog entry

### Thoughts & Emotions
- **Decisive** — issuing is a one-way action; the partner confirms with full awareness
- **Relieved** — the primary engagement deliverable is done
- **Aware** — the computed deadlines are shown clearly so the partner understands the regulatory timeline

### Pain Points
- **Competitor context:** Agentive has archiving but no explicit WORM or assembly deadline enforcement. Hyperproof and AuditBoard retain reports but don't enforce regulator-specific assembly windows. Drata / Vanta don't issue reports — they support the audit but don't finalize it.
- If the partner issues accidentally (premature click), reversal is complex (must withdraw and re-issue)
- The distinction between AICPA SOC, ISO certification-body, and PCI DSS assembly windows must be clear

### Opportunities
- Two-step issuance confirmation: "You are about to issue this report. Assembly deadline will be [date]. This action cannot be undone. Confirm?"
- Auto-share issued report with client in the Client Hub (configurable per engagement)

---

## Stage 5: Finalize the Engagement

### Sub-goal
Lock the engagement file — no further modifications to any workpaper, evidence, or test procedure.

### User Actions
- Transitions the engagement to Finalized status (blocked unless Report.status = Issued)
- At Finalized: all workpapers locked (is_locked = true), all evidence links frozen, all control conclusions immutable
- Any modification attempt returns a hard error: "This engagement has been finalized. Modifications require an addendum."
- If a post-finalization error is discovered, the partner can create an addendum: new WorkpaperVersion with is_addendum = true, documented reason, and partner sign-off

### Touchpoints
- Finalize transition button with guard validation
- Engagement-wide lock confirmation
- Addendum creation interface (post-finalization)
- Hard error on any modification attempt

### Thoughts & Emotions
- **Final** — this is the point of no return for the engagement file
- **Secure** — the immutability gives regulatory confidence
- **Attentive** — if an addendum is needed, the process is formal and documented

### Pain Points
- **Competitor context:** No GRC competitor has a formal addendum workflow equivalent to Axiom's (new WorkpaperVersion with is_addendum = true, documented reason, partner sign-off). ISO 17021-1 audit-record maintenance and AICPA AT-C §A.60 assembly rules both require this pattern.
- The gap between issuance and finalization is the "cleanup window" — firms may need a few days to resolve administrative items before locking
- Partners need to understand that addenda are not edits — they are new records appended to the file

### Opportunities
- Pre-finalization checklist: "All workpapers signed off ✓. All review notes resolved ✓. Report issued ✓. Ready to finalize."
- Addendum workflow with clear explanation: "This addendum will be added to the engagement file as a new record. The original content is unchanged."

---

## Stage 6: Monitor Automatic Archival

### Sub-goal
Ensure the engagement transitions correctly from Finalized to Archived when the assembly deadline elapses, with proper WORM archiving and retention enforcement.

### User Actions
- The EngagementLifecycleStateMachine in Step Functions triggers the Finalized → Archived transition automatically when report_issued_at + assembly_window has elapsed
- At archival: all files copied to S3 Object Lock bucket (COMPLIANCE mode) with retention_deadline, engagement becomes read-only, no further state changes possible
- Partner receives archival confirmation email
- Before retention expiry (90 days warning): FirmAdmin receives notification with export option
- At any time: FirmAdmin can generate a complete engagement export (all workpapers as PDF, evidence files in native format, cross-framework coverage matrix as CSV, AIDecision ledger as CSV, AuditLog as CSV, metadata as JSON — structured ZIP file)

### Touchpoints
- Archival confirmation email to engagement partner
- Engagement read-only indicator in the application
- Retention expiry warning (90 days prior) to FirmAdmin
- Engagement export interface (available at any time)
- AuditLog: system-triggered archival event

### Thoughts & Emotions
- **Confident** — the system handles archival automatically per regulatory requirements
- **Reassured** — S3 WORM storage means the file is tamper-proof
- **Future-aware** — retention deadlines are years away, but the system tracks them

### Pain Points
- **Competitor context:** Hyperproof and AuditBoard retain reports but don't use immutable WORM storage. No competitor specifies S3 Object Lock COMPLIANCE mode for audit-file retention. This is a strong regulatory differentiator and reinforces the post–Delve-scandal provenance positioning.
- The automatic transition must never fail silently — the partner needs confirmation
- Engagement export must be comprehensive enough for standalone archival independent of the platform (critical for offboarding or firm dissolution scenarios)

### Opportunities
- Annual archival report: "12 engagements archived this year. All within assembly deadline windows. Export available."
- Retention countdown: visible on archived engagements showing years remaining until deletion
- GDPR/CCPA compliance: automatic deletion at retention expiry with notification and export opportunity

---

## Handoffs

| From | To | Information Transferred | Trigger |
|------|-----|------------------------|---------|
| Partner | Client Contact (Journey 8) | Draft report for review | Draft shared |
| Partner | Client Contact (Journey 8) | Final report in Client Hub | Report issued |
| System | FirmAdmin | Archival confirmation + retention timeline | Automatic archival |

---

## Journey Summary

### Emotional Arc
Begins with the focused energy of report creation (the final deliverable), sustains through the careful editing and client review cycle, peaks at the decisive issuance moment ("this engagement is done"), settles into the administrative satisfaction of finalization, and fades into the quiet confidence of automatic archival. The emotional through-line is "the system handles the regulatory complexity so I can focus on professional judgment."

### Cross-Cutting Pain Points
- Report template quality determines the partner's experience — finicky templates are the most-cited complaint across every GRC tool
- The issuance → finalization → archival chain must be flawless — errors here have regulatory consequences
- Engagement export is a critical trust feature — firms need to know they own their data independent of the platform

### Prioritized Opportunities
1. **Automatic assembly deadline computation and enforcement** (high impact, medium effort) — a genuine regulatory differentiator over all competitors
2. **S3 WORM archival** (high impact, high effort) — strongest archival guarantee in the market
3. **Addendum workflow** (medium impact, medium effort) — proper ISO 17021-1 / AICPA AT-C assembly-integrity implementation; no competitor has this
4. **Engagement export** (high impact, medium effort) — critical for trust, offboarding, and regulatory compliance

---

# Journey 10: Conduct Engagement Quality Review

## Overview
- **Persona:** EQR Reviewer — experienced partner or external reviewer, independent of the engagement team, assigned per firm quality policy (SQMS 2 for SOC engagements; ISO 17021-1 §9.6 internal review analogs for ISO engagements)
- **Goal:** Review the engagement file, document review scope and findings, and sign off so the engagement can advance from Review to Reporting
- **Trigger:** Engagement reaches Review status and the reviewer is notified
- **Stages:**
  1. Access engagement in read-only mode
  2. Review engagement scope and planning quality
  3. Review testing and evidence sufficiency
  4. Document findings and conclusion
  5. Sign off on EQR

## Stage 1: Access Engagement in Read-Only Mode

### Sub-goal
Enter the engagement file with full visibility but zero ability to modify content.

### User Actions
- Receives notification: "Engagement [name] has reached Review status and is ready for your quality review"
- Clicks through to the engagement — the system enforces read-only access (the reviewer is not an EngagementTeamMember)
- Sees all engagement content: workpapers, evidence, controls, test procedures, cross-framework coverage matrix, review notes, AI decisions, audit log
- Cannot edit, sign off, or modify any engagement content

### Touchpoints
- EQR notification (in-app + email)
- Engagement view with "EQR — Read Only" indicator
- Full navigation of all engagement sections
- EngagementQualityReview record with status = InProgress

### Thoughts & Emotions
- **Independent** — the entire point is fresh eyes with no involvement in the work
- **Methodical** — will review systematically, not skimming
- **Authoritative** — the reviewer's sign-off carries significant regulatory weight

### Pain Points
- **Competitor context:** Few GRC competitors have a dedicated EQR workflow. Agentive supports EQR but specifics are limited. Hyperproof and AuditBoard leave EQR to firm-level processes. No competitor enforces the independence check (reviewer ≠ team member) at the system level.
- If the read-only mode doesn't provide sufficient navigation tools, the reviewer spends too long finding things
- The reviewer needs to see everything without being overwhelmed — they're reviewing, not performing the audit

### Opportunities
- EQR-focused view: pre-organized sections aligned with SQMS 2 (and ISO 17021-1 equivalent) review scope — planning quality, risk assessment, testing sufficiency, documentation adequacy, significant judgments, AI-assisted content review
- Highlight areas that may need attention: controls with exceptions, low-confidence AI decisions, open review notes that were resolved

---

## Stage 2: Review Engagement Scope and Planning Quality

### Sub-goal
Assess whether the engagement was properly planned — correct framework, appropriate team, adequate risk assessment, and proper client acceptance.

### User Actions
- Reviews the ClientAcceptance record: quality risks, firm responses, independence confirmation
- Reviews the engagement setup: framework selection, team composition, EQR independence documentation
- Reviews the AI control mapping decisions: were the AI suggestions reviewed and confirmed by the engagement team?
- Documents scope notes in the EngagementQualityReview record

### Touchpoints
- ClientAcceptance record (read-only)
- Engagement configuration details
- AI control mapping review trail (AIDecision records)
- EQR scope notes entry field

### Thoughts & Emotions
- **Evaluative** — checking the foundation before examining the building
- **Questioning** — "Did the partner adequately consider this quality risk?"
- **Documenting** — everything observed goes into the EQR record

### Pain Points
- If the ClientAcceptance is too brief or uses only template language, the reviewer can't assess the quality of the risk evaluation
- AI decision records must be accessible and clear — the reviewer needs to understand what the AI suggested and what the team decided

### Opportunities
- AI decision summary: "42 AI decisions in this engagement. 38 accepted as-is, 3 modified, 1 rejected."
- Planning quality indicators: flags for unusual patterns (e.g., acceptance completed in <5 minutes, all AI mappings bulk-confirmed without individual review)

---

## Stage 3: Review Testing and Evidence Sufficiency

### Sub-goal
Assess whether the testing was thorough, evidence is sufficient, and conclusions are supported.

### User Actions
- Reviews a sample of completed workpapers across different control areas
- Examines evidence links: is the evidence relevant, sufficient, and properly linked?
- Reviews exception handling: were exceptions properly investigated and documented?
- Checks AI-assisted content: reviews WorkpaperVersion history to confirm AI drafts were substantively edited
- Reviews review notes history: were manager review notes properly raised and resolved?

### Touchpoints
- Workpaper list with completion status and sign-off trail
- Evidence linking summary per control
- Exception summary with investigation documentation
- WorkpaperVersion history (AI draft → human edit trail)
- Review notes archive (open → resolved history)

### Thoughts & Emotions
- **Thorough** — this is the substance of the review; the reviewer is applying deep professional judgment
- **Critical** — looking for gaps, not confirming the team's conclusions
- **Supportive** — the goal is quality assurance, not blame

### Pain Points
- If the engagement has 89 controls and 200+ workpapers, the reviewer needs efficient sampling tools — reviewing everything is impractical
- AI-drafted content that wasn't substantively edited (is_ai_draft was set to false by a trivial edit) is a quality risk the reviewer must catch
- The review notes archive must distinguish between substantive notes (required changes) and minor notes (formatting)

### Opportunities
- EQR sampling guide: "Based on engagement risk and size, we suggest reviewing at least 15% of workpapers across all control categories"
- AI edit substantiveness indicator: "This workpaper's AI draft was edited by [name]. 2 of 8 sections modified. Modification percentage: 25%."
- Exception highlight: "3 controls have exceptions noted. Review these first."

---

## Stage 4: Document Findings and Conclusion

### Sub-goal
Record the EQR findings and reach a conclusion about the engagement's quality.

### User Actions
- Documents review scope: what was reviewed, sampling approach, time spent
- Documents findings: any quality concerns, suggestions, or required actions
- Records conclusion: satisfied / satisfied with concerns noted / not satisfied
- Findings that require action must be addressed by the engagement team before sign-off
- Required actions create an obligation — the partner must respond before the EQR can be signed off

### Touchpoints
- EQR findings documentation interface (structured + narrative)
- Conclusion selection with explanation
- Required action items with assignment to engagement team
- Communication channel with engagement partner

### Thoughts & Emotions
- **Decisive** — the conclusion has regulatory weight
- **Fair** — balanced assessment considering the engagement's complexity and constraints
- **Responsible** — the EQR sign-off becomes part of the immutable engagement archive

### Pain Points
- If findings require engagement team action, the back-and-forth cycle delays the Review → Reporting transition
- The reviewer must balance thoroughness with timeliness — the engagement can't sit in Review indefinitely

### Opportunities
- Structured findings template: categorize by severity (observation, recommendation, required action)
- Action tracking: each required action shows status (pending, addressed, confirmed) with engagement team responses

---

## Stage 5: Sign Off on EQR

### Sub-goal
Formally sign off on the engagement quality review, unblocking the Review → Reporting transition.

### User Actions
- Confirms all required actions have been addressed (if any were raised)
- Signs off on the EngagementQualityReview: signed_off_at is populated
- The Review → Reporting transition is now unblocked for the Partner
- The EQR record and sign-off timestamp become part of the immutable engagement archive

### Touchpoints
- Sign-off button (disabled if required actions are unresolved)
- Sign-off confirmation with name, timestamp, and summary of findings
- Notification to engagement partner: "EQR complete. Engagement can proceed to Reporting."
- AuditLog entry for EQR sign-off

### Thoughts & Emotions
- **Confident** — the engagement meets quality standards
- **Final** — the EQR sign-off is immutable once recorded
- **Professional pride** — the reviewer's name is permanently attached to this quality assessment

### Pain Points
- **Competitor context:** No competitor implements the EQR gate at the system level with immutable sign-off records. This is a genuine regulatory compliance differentiator.
- If the reviewer discovers a significant issue at this late stage, the engagement team faces rework pressure
- The immutability of the EQR record must be clearly communicated — this is not a form that can be amended

### Opportunities
- Historical EQR comparison: "In prior year, this engagement had 2 findings. This year: 0 findings."
- EQR completion notification to the engagement team with summary: "EQR complete with [0/N] findings. Proceed to Reporting."

---

## Handoffs

| From | To | Information Transferred | Trigger |
|------|-----|------------------------|---------|
| EQR Reviewer | Partner (Journey 9) | EQR sign-off unblocking Review → Reporting | EQR complete |
| EQR Reviewer | Engagement Team | Required action items (if any) | Findings documented |

---

## Journey Summary

### Emotional Arc
Starts with the disciplined entry into read-only review mode (independence and objectivity), builds through the methodical examination of planning and testing (professional thoroughness), navigates the potentially tense findings documentation phase (where quality concerns must be raised diplomatically but clearly), and resolves with the formal sign-off (professional confidence and regulatory compliance). The emotional high is the clean sign-off; the emotional low is discovering a significant quality issue that requires engagement team rework.

### Cross-Cutting Pain Points
- EQR efficiency vs. thoroughness — the reviewer must be thorough without becoming a bottleneck
- AI-assisted content adds a new dimension to quality review — the reviewer must evaluate both the AI's work and the auditor's review of the AI's work
- The EQR gate delays the engagement timeline; anything that accelerates the review without reducing quality is valuable

### Prioritized Opportunities
1. **EQR-focused navigation view** (high impact, medium effort) — aligns the review interface with SQMS 2 / ISO 17021-1 review scope rather than the engagement's natural structure
2. **AI edit substantiveness indicator** (high impact, low effort) — surfaces the most important quality concern in AI-assisted engagements
3. **Structured findings with action tracking** (medium impact, medium effort) — streamlines the reviewer-to-team communication cycle
4. **Historical EQR comparison** (low impact, low effort) — useful context but not critical

---

# Journey 11: Multi-Framework Integrated Engagement

## Overview
- **Persona:** Partner — scoping a single engagement that covers SOC 2 + ISO 27001 + ISO 27701 simultaneously (the most common "year-one credibility stack" for an early-stage security-sensitive company; PCI DSS or ISO 42001 may be added as a fourth track depending on client)
- **Goal:** Run one engagement, not three parallel ones — shared `CommonControl` library, shared evidence, reconciled sampling windows, one cohesive fieldwork plan — with an auditor-defensible separation of opinion per framework
- **Trigger:** Client signs for an integrated SOC 2 Type II + ISO 27001 initial certification + ISO 27701 extension package (or a rollforward that combines an ISO 27001 surveillance visit with a SOC 2 Type II)
- **Stages:**
  1. Scope the integrated engagement
  2. Reconcile sampling windows across frameworks
  3. Unify the shared `CommonControl` library
  4. Coordinate fieldwork across tracks
  5. Separate-but-linked report issuance

## Stage 1: Scope the Integrated Engagement

### Sub-goal
Define the integrated engagement with explicit framework tracks, without letting the convenience collapse into a single undifferentiated opinion.

### User Actions
- From Journey 3's Stage 1, selects "Integrated multi-framework engagement"
- Chooses a primary framework track (typically SOC 2 Type II) and one or more secondary framework tracks (ISO 27001, ISO 27701 extension, optionally PCI DSS or ISO 42001)
- Each framework track gets its own `EngagementFramework` record, period window, version, and scope statement — a single engagement in the UI, multiple opinions / certificates downstream
- The system pulls the pre-built STRM-encoded crosswalks between the selected frameworks and presents a coverage forecast: "With the licensed SCF + AICPA mappings, your SOC 2 evidence is expected to cover ~82% of ISO 27001 A.5–A.8, ~71% of ISO 27701 privacy requirements. Explicit gaps will be tracked per framework."

### Touchpoints
- Integrated engagement wizard
- Coverage forecast preview with framework-by-framework expected overlap
- Per-framework scope statement form

### Thoughts & Emotions
- **Strategic** — "One engagement for our client, three opinions, no double-collection"
- **Cautious** — "I need to keep the opinion boundaries clear; mixing SOC 2 and ISO certification scope would be a regulatory error"

### Pain Points
- Drata / Vanta / Secureframe market multi-framework support but the auditor still has to run multiple parallel engagements in their own tooling
- Hyperproof and AuditBoard via UCF can model multi-framework but sampling-window reconciliation is manual

### Opportunities
- Visually distinct framework tracks with independent status, so the partner can see at a glance that SOC 2 is in fieldwork while ISO 27001 is still in planning
- Per-track scope statement prompts pulled from the framework body (AICPA AT-C 205 scope for SOC 2; ISO 17021-1 §9.2 scope worksheet for ISO)

---

## Stage 2: Reconcile Sampling Windows Across Frameworks

### Sub-goal
Handle the fact that SOC 2 Type II wants a continuous period (3–12 months; commonly 6 or 12), SOC 2 Type I is a single as-of date, ISO 27001 surveillance samples at a point in time, and PCI DSS wants 90-day ASV scan validity — a single calendar plan must honor all of these.

### User Actions
- Opens the sampling window reconciliation view
- System shows each framework's native window: "SOC 2 Type II: 2026-01-01 → 2026-12-31 continuous; ISO 27001: 2026-11-15 surveillance sample point with 3-month look-back acceptable; PCI DSS: ASV scans within 90 days of report date"
- AI proposes a unified evidence-collection schedule that satisfies all three windows with the fewest redundant asks on the client
- Partner reviews and adjusts; every adjustment is an AIDecision
- Output: a master sampling calendar — which evidence collected when, which framework each collection satisfies, and where a single collection does double or triple duty

### Touchpoints
- Sampling window reconciliation view (horizontal framework lanes against a shared timeline)
- AI-proposed schedule with explanation per item
- Partner override per scheduled collection

### Thoughts & Emotions
- **Relieved** — "I don't have to build this calendar in Excel"
- **Attentive** — sampling is where cross-framework shortcuts most often break; the partner reads this carefully

### Pain Points
- Getting this wrong means either over-collection (client pain) or under-coverage (audit failure). The UI must make both failure modes visible.
- Framework version differences (ISO 27001:2013 vs. 2022) shift window expectations; migration must be surfaced

### Opportunities
- Show each framework's sampling rule inline with the schedule: "Why is this required for PCI 11.3? Because the ROC requires an ASV scan within 90 days of report issuance."
- Flag over-collection opportunities: "Item collected 4 times across frameworks — consolidate to a single quarterly collection?"

---

## Stage 3: Unify the Shared CommonControl Library

### Sub-goal
One control library for the whole engagement, with explicit per-framework coverage.

### User Actions
- The engagement's `CommonControl` library is assembled from the activated methodology templates plus any firm-level custom `CommonControl` nodes
- Each `CommonControl` shows which `FrameworkRequirement` nodes it satisfies across every in-scope framework, with STRM relationship type and strength
- Partner reviews the library — same cross-framework mapping review as Journey 3 Stage 3 but framed around integrated coverage
- Any `FrameworkRequirement` not covered by any `CommonControl` surfaces as a "gap-at-scoping" — these must be addressed before fieldwork begins

### Touchpoints
- Integrated CommonControl library view
- Gap-at-scoping panel: "4 ISO 27701 requirements have no `CommonControl` covering them — add or map before advancing"

### Thoughts & Emotions
- **Organized** — one library, one mental model
- **Vigilant** — gaps surfaced now are easier to fix than gaps discovered in fieldwork

### Pain Points
- A gap-at-scoping list that's too long demoralizes the team
- CommonControl semantics drift between frameworks (e.g., "access review" means different things in SOC 2 CC6.3 vs ISO 27001 A.5.18 vs HIPAA §164.308(a)(4))

### Opportunities
- Suggested CommonControl additions to close gap-at-scoping items
- Semantic-drift warnings when the same CommonControl name spans heterogeneous framework nodes

---

## Stage 4: Coordinate Fieldwork Across Tracks

### Sub-goal
Run fieldwork once, with per-framework progress tracking, so the team doesn't accidentally skip a track.

### User Actions
- Staff auditors pick up assigned controls via Journey 5, but each control's "complete" state is evaluated per framework track — a control can be Complete-for-SOC-2 while still being Incomplete-for-ISO-27001 if an ISO-specific sub-requirement is unmet
- Manager review (Journey 6) and EQR (Journey 10) happen at the engagement level but with per-framework completion rollups
- The engagement progress dashboard shows three (or more) parallel track progress bars rather than a single bar

### Touchpoints
- Per-track progress dashboard
- Control detail view with per-framework completion state

### Thoughts & Emotions
- **Coordinated** — the team sees progress in the same shape the engagement will be reported in
- **Clear-headed** — no accidental "we're done" when one track is still open

### Pain Points
- If track-level progress is buried under a single aggregate, mistakes happen
- Staff auditors may not understand why a control is "Complete for SOC 2 but not ISO 27001" — needs clear explanation

### Opportunities
- Track-specific completion criteria inline on each control
- Staff-friendly "what's missing for this track" prompts

---

## Stage 5: Separate-but-Linked Report Issuance

### Sub-goal
Issue each framework's deliverable on its own regulatory cadence, from a single cohesive engagement file.

### User Actions
- Each framework track ends in its own report: SOC 2 report, ISO 27001 certificate support letter (the certification body issues the actual certificate; Axiom produces the supporting assurance file), ISO 27701 certificate support letter, PCI ROC + AOC, HIPAA attestation letter, ISO 42001 certificate support letter
- Each report follows Journey 9 independently — issuance, finalization, and assembly-deadline computation per framework
- The engagement-level Finalized state is reached only when all tracks have reached Finalized
- Cross-report references are preserved: the ISO 27001 support letter references that SOC 2 evidence was reused, with explicit provenance trail

### Touchpoints
- Per-framework report selection from Journey 9 Stage 1
- Cross-report reference block auto-populated from the cross-framework coverage matrix
- Engagement-level finalization that waits on all tracks

### Thoughts & Emotions
- **Decisive** — each track closes on its own cadence
- **Provenance-aware** — auditors of the ISO report can trace exactly which SOC 2 evidence was reused and why that reuse is defensible

### Pain Points
- Misaligned track finalization dates can cause staff to "forget" an open track
- Each framework has its own report format idiosyncrasies that must be honored

### Opportunities
- Pre-finalization checklist at engagement level showing each track's state
- Track-specific report template library

---

## Handoffs

| From | To | Information Transferred | Trigger |
|------|-----|------------------------|---------|
| Partner | Staff Auditor (Journey 5) | Per-track assigned controls with framework completion criteria | Fieldwork begins |
| Partner | EQR Reviewer (Journey 10) | Integrated engagement file with per-framework rollups | Review phase reached |
| Partner | Partner (Journey 9) | Per-framework reports, each on its own assembly/retention clock | Each track's Reporting phase |

---

## Journey Summary

### Emotional Arc
Opens with the strategic clarity of integrated scoping ("one engagement, three opinions"), steadies through the careful reconciliation of sampling windows (the intellectually demanding core), cruises through unified fieldwork (the team sees one coherent plan), and resolves with the decisive separate-but-linked issuance (each framework closes on its own terms). The emotional high is the sampling-window reconciliation moment — where the product does what every competitor pretends to do but doesn't. The emotional risk is letting integrated convenience blur opinion boundaries.

### Cross-Cutting Pain Points
- Integrated ≠ homogenized — opinion boundaries must stay crisp per framework
- Sampling-window reconciliation is the #1 technical challenge in multi-framework engagements and is largely invisible in every competitor's product today
- Framework version drift (ISO 27001:2013 → 2022, PCI 3.2 → 4.0.1) re-introduces gaps every migration cycle

### Prioritized Opportunities
1. **Visual framework-track lanes with independent progress** (high impact, medium effort) — the conceptual backbone of multi-framework UX
2. **AI-proposed sampling calendar with per-item framework rationale** (high impact, high effort) — the technical differentiator
3. **Gap-at-scoping panel surfaced before fieldwork begins** (high impact, low effort) — fixes the most expensive class of integrated-engagement error
4. **Cross-report provenance references** (medium impact, medium effort) — reinforces auditor-defensibility post–Delve-scandal

---

# Journey 12: Continuous Assurance (Auditee-Side)

## Overview
- **Persona:** ClientAdmin — at the auditee, operating the full auditee GRC workspace (the extension of Journey 8's Client Hub). Continuous monitoring is enabled on the engagement.
- **Goal:** When a configuration drift, access change, or other control-relevant event occurs between formal audit windows, get an immediate alert, auto-retest the affected controls, update the risk register, and — where the drift is material — notify the auditor so the official audit file reflects reality
- **Trigger:** A drift event — IdP MFA policy downgraded, production S3 bucket reconfigured, a privileged access role granted outside the approved change process, a new AI model deployed without an ISO 42001 impact review, an ASV scan producing a high-severity finding
- **Note:** This journey is the "both-sided" product bet. Single-sided competitors (Drata/Vanta on the auditee side; Agentive on the auditor side) structurally cannot run this loop end-to-end. Delve claims continuous monitoring but the March 2026 scandal raised the question of what "continuous" actually means without a credible AIDecision ledger.
- **Stages:**
  1. Detect the drift
  2. Auto-retest affected controls
  3. Update the risk register
  4. Draft the management response
  5. Notify the auditor on material change

## Stage 1: Detect the Drift

### Sub-goal
Catch the change within minutes, with enough context to know what it affects.

### User Actions
- Continuous monitoring jobs (River workers polling connected SaaS sources, listening on Webhook receivers, or running scheduled config snapshots) detect a change that affects a `CommonControl` node the engagement tracks
- The change is fingerprinted, timestamped, and persisted as a `DriftEvent` with before/after state and the identity of the actor (where known)
- ClientAdmin receives an in-app alert and an email: "MFA policy changed on Okta — `Require MFA for Admins` toggled off by [user] at [timestamp]. This affects SOC 2 CC6.1 and ISO 27001 A.5.15."

### Touchpoints
- Drift alert notification (in-app + email, routing configurable)
- DriftEvent detail view with before/after diff and affected controls
- Link into the evidence freshness dashboard

### Thoughts & Emotions
- **Alert** — the ClientAdmin treats this like a monitoring page rather than a compliance nag
- **Oriented** — "I know exactly what changed, who did it, and which controls are now in question"

### Pain Points
- False positives erode trust fast — config snapshots that flap generate alert fatigue
- Connected-source gaps: a drift event in a source Axiom doesn't monitor (e.g., an on-prem system) is a known blind spot that must be labeled

### Opportunities
- Severity tiers on drift alerts based on blast radius (how many controls and frameworks are affected)
- Connected-source coverage dashboard so the ClientAdmin knows what is and isn't monitored

---

## Stage 2: Auto-Retest Affected Controls

### Sub-goal
Answer, quickly and automatically, the question "is the control still satisfied?"

### User Actions
- The platform auto-triggers the control tests that the drifted evidence supports — for log-stream / configuration controls, Axiom reruns the autonomous control test (Journey 5 Stage 2's agentic testing) against the new state
- Each re-test writes an AIDecision record — the outcome (still satisfied / partially satisfied / no longer satisfied) is recorded with the evidence it was based on
- Re-test results feed back into the framework-coverage matrix — SOC 2 CC6.1 may drop from Full → Not, ISO 27001 A.5.15 may drop from Full → Partial

### Touchpoints
- Auto-retest status panel: which controls queued, running, completed
- Updated coverage matrix with the newly Not or Partial cells highlighted

### Thoughts & Emotions
- **Empowered** — the platform tells the ClientAdmin what the change actually means, in framework-native terms
- **Pragmatic** — "I have a real gap now. What's the fix?"

### Pain Points
- Auto-retest that can't run (source unavailable, manual evidence required) must be clearly handled, not left silently pending
- Every AIDecision record must capture the evidence and reasoning, not just the verdict

---

## Stage 3: Update the Risk Register

### Sub-goal
Reflect the drift-induced risk in the live risk register — don't let gaps exist only in the compliance tool.

### User Actions
- The platform proposes a `RiskRegisterEntry` update based on the auto-retest outcome: new entry for a newly-Not control, severity rescored for a newly-Partial control
- ClientAdmin reviews and accepts / modifies / rejects (AIDecision record each time)
- The risk register and the compliance coverage matrix stay synchronized — a fundamental product commitment

### Touchpoints
- Risk register update proposal card
- Accept / modify / reject buttons with AIDecision capture

### Thoughts & Emotions
- **Integrated** — "My risk register and my compliance file are telling the same story"
- **Accountable** — every accept/modify/reject is on the record

### Pain Points
- Risk register hygiene varies enormously across organizations; the auto-proposal must defer to the ClientAdmin's policy
- Severity scoring is judgment-heavy — AI should propose, not decide

---

## Stage 4: Draft the Management Response

### Sub-goal
Get from "we have a problem" to "we have a plan" within minutes, with the plan in auditor-defensible form.

### User Actions
- ClientAdmin requests an agentic management-response draft
- AI drafts a remediation plan tied to the specific drift, the affected controls, the risk register entry, and the ClientAdmin's ticketing system of choice (Jira / Linear ID is generated with appropriate labels)
- ClientAdmin edits the draft (is_ai_draft must flip to false before submission — same HITL gate as Journey 5 workpaper drafts)
- On submission, the ticket is round-tripped into Jira/Linear and the remediation target date is written onto the RiskRegisterEntry

### Touchpoints
- "Draft management response" button on the DriftEvent
- Draft editor with linked ticket preview
- Round-trip status indicator

### Thoughts & Emotions
- **Quick** — from alert to actionable ticket in under 10 minutes
- **Trust-aware** — the AI drafts, the human commits

### Pain Points
- Ticketing integrations are brittle; failed round-trip must be surfaced
- Over-reliance on AI drafts creates the exact Delve-scandal risk — boilerplate remediations. The edit-before-submit gate is non-negotiable.

---

## Stage 5: Notify the Auditor on Material Change

### Sub-goal
Keep the auditor's engagement file honest when material changes happen between formal windows.

### User Actions
- When a DriftEvent's severity or blast radius exceeds a configurable threshold (firm-level default + client override), the system notifies the assigned auditor team on the active engagement
- The auditor sees the DriftEvent, the auto-retest results, the RiskRegisterEntry update, and the management response in their engagement view (read-only from the auditor's side unless they choose to assess further)
- Auditor judgment determines whether the change requires a workpaper addendum (Journey 9 Stage 5's addendum flow) and whether the change alters the in-process opinion
- All of this flows through the AIDecision + AuditLog backbone, producing an unambiguous provenance trail: "This is what changed, this is how we auto-retested, this is what the auditee did about it, this is how the auditor responded"

### Touchpoints
- Material-change notification to auditor team
- Auditor-side drift review panel with the full chain of evidence
- Link into Journey 9 addendum workflow if the change requires one

### Thoughts & Emotions
- **Transparent** — the auditee knows exactly what the auditor sees
- **Defensible** — the auditor has a continuous record, not a point-in-time snapshot that the world may have already invalidated
- **Collaborative** — the drift-response loop is a shared workflow, not a gotcha

### Pain Points
- Threshold tuning is hard — too low and the auditor is deluged, too high and material changes slip
- The auditee-auditor relationship is delicate; the notification UX must feel like professional communication, not a surveillance tool

### Opportunities
- Firm-level thresholds tunable per client and per framework — SOC 2 Type II has different significance criteria than an ISO 27001 surveillance visit
- Historical drift rate displayed on the engagement dashboard: "3 material drift events this period; auditor acknowledged all 3"
- Drift-to-opinion trail exportable as part of the engagement export (Journey 9 Stage 6)

---

## Handoffs

| From | To | Information Transferred | Trigger |
|------|-----|------------------------|---------|
| ClientAdmin | Auditor team (Journeys 5, 6, 9) | Material DriftEvent, auto-retest outcome, RiskRegisterEntry update, management response | Threshold exceeded |
| ClientAdmin | ClientAdmin (Journey 4) | Updated coverage matrix after auto-retest | Auto-retest complete |
| Auditor | Partner (Journey 9) | Material change may require workpaper addendum | Auditor review decision |

---

## Journey Summary

### Emotional Arc
Opens with the immediate alertness of a real-time drift alert (the auditee is paying attention now, not at audit time), passes through the reassuring automation of auto-retest and risk-register sync, builds through the actionable management-response draft, and resolves with the transparent material-change notification to the auditor. The emotional through-line is the thing every single-sided GRC tool cannot promise: "my compliance story is true right now, not just at the last audit window."

### Cross-Cutting Pain Points
- Alert fatigue from false positives is the #1 adoption risk
- The auditee ↔ auditor communication must feel professional and collaborative, not adversarial
- Continuous assurance without a credible AIDecision ledger is exactly the Delve-scandal posture; the provenance discipline is not optional

### Prioritized Opportunities
1. **Configurable significance thresholds with firm defaults and client overrides** (high impact, medium effort) — the essential knob for making the auditor-notification channel usable
2. **Round-trip ticketing into Jira / Linear with labels** (high impact, medium effort) — the feature that makes drift response real work, not a compliance ritual
3. **Connected-source coverage dashboard with explicit blind-spot labeling** (high impact, low effort) — honest monitoring beats marketed monitoring
4. **Drift-to-opinion export bundled with engagement archive** (medium impact, medium effort) — completes the provenance story

---

# Cross-Journey Insights

## Competitor Pain Points Addressed by Axiom

| Competitor Pain Point | Source | Axiom Response | Journeys |
|----------------------|--------|---------------|----------|
| Multi-week implementation with dedicated consultant | Hyperproof / AuditBoard onboarding | Self-serve onboarding, <1 week to first engagement | 1 |
| Steep learning curve, long time-to-value | Hyperproof, AuditBoard | 5-step guided tour, intuitive UI, deep-link to first task | 2 |
| No split-view for evidence alongside workpapers | Hyperproof, AuditBoard | Split-view evidence panel during testing and review | 5, 6 |
| Auditor-side and auditee-side are separate products | Drata/Vanta/Secureframe (auditee-only); Agentive (auditor-only) | Unified platform — Client Hub extends into full auditee GRC workspace | 4, 7, 8, 12 |
| PBC as separate add-on or email workflow | Legacy tools and Drata/Vanta | Integrated document request management | 7, 8 |
| Boilerplate-driven report generation with trust issues | Delve (March 2026 scandal — 493 of 494 reports shared identical text and typo) | Grounded crosswalk, AIDecision ledger, WORM archival, firm-specific template matching | 3, 9 |
| Report templates are finicky | All competitors | Pre-populated templates with engagement data and cross-framework coverage | 9 |
| No cross-framework evidence mapping grounded in a licensed crosswalk | Every competitor claims it but few ground it in SCF/OSCAL/AICPA | SCF + OSCAL + AICPA + CIS backbone with STRM-encoded edges, partial-satisfaction gap lists, and period-aware freshness | 3, 4, 5, 11 |
| No system-enforced regulatory gates | Most GRC competitors | SQMS 1 (or ISO 17021-1 equivalent) acceptance gate, EQR independence check, assembly deadline enforcement, AI draft edit requirement | 3, 9, 10 |
| No continuous assurance loop between auditor and auditee | Single-sided competitors structurally can't | Drift-triggered re-testing with material-change notification to the auditor | 12 |
| No AI provenance / HITL ledger | All competitors, acutely after March 2026 Delve scandal | AIDecision table records every accept/modify/reject; cryptographic evidence fingerprinting at upload | 3, 4, 5, 7, 9 |

## Flows Without Competitor Equivalent

These Axiom flows represent genuine innovation — no competitor currently offers them:

1. **Cross-framework evidence mapping with partial-satisfaction gap lists** (Journeys 3, 4, 5, 11) — one piece of evidence simultaneously satisfying SOC 2, ISO 27001, ISO 27701, HIPAA, and PCI requirements, grounded in the licensed SCF/OSCAL/AICPA crosswalk, with explicit gap lists for partial coverage and period-aware freshness per framework
2. **AI document completeness review with client-facing gap explanations** (Journeys 7, 8) — AI analyzes uploaded documents against request requirements and auto-drafts specific gap explanations for the client
3. **System-enforced AI draft edit gate with AIDecision ledger** (Journeys 5, 6) — workpapers with unedited AI content cannot be signed off; every AI accept/modify/reject is recorded. This is the direct product response to the March 2026 Delve scandal.
4. **EQR independence enforcement** (Journey 10) — system-level validation that the quality reviewer is not on the engagement team
5. **Post-finalization addendum workflow** (Journey 9) — immutable original content with versioned addenda, sign-off required
6. **Agentic autonomous control testing for log-stream and configuration populations** (Journey 5) — continuous testing against SOC 2 CC6 logs, PCI DSS §10 logging, ISO 27001 A.8.15 logging, and ISO 42001 AI system monitoring, with every decision written to the AIDecision ledger
7. **Automatic assembly deadline computation and WORM archival** (Journey 9) — computed at report issuance, enforced via S3 Object Lock COMPLIANCE mode
8. **Unified auditor + auditee GRC surface with continuous assurance loop** (Journeys 4, 8, 11, 12) — the Client Hub extends into a full auditee GRC workspace; drift alerts propagate from the auditee side to the auditor side automatically

## Data Entities by Journey

| Journey | Primary Entities |
|---------|-----------------|
| 1: Firm Setup | Firm, User, MethodologyTemplate, Engagement, Control, TestProcedure, Workpaper, ClientAcceptance |
| 2: Staff Onboarding | User, EngagementTeamMember |
| 3: Engagement Scoping | Engagement, EngagementTeamMember, EngagementFramework, Control, TestProcedure, ClientAcceptance, EngagementQualityReview, FirmControlObjectiveMapping, AIDecision |
| 4: Cross-Framework Evidence Mapping | EvidenceItem, CommonControl, FrameworkRequirement, FrameworkVersion, satisfies edges (STRM-encoded), AIDecision |
| 5: Control Testing | Control, TestProcedure, EvidenceItem, EvidenceLink, Workpaper, WorkpaperVersion, AIDecision |
| 6: Workpaper Review | Workpaper, WorkpaperVersion, AuditLog |
| 7: Document Requests | DocumentRequest, EvidenceItem, EvidenceLink, AIDecision, AuditLog |
| 8: Client Fulfillment | DocumentRequest, EvidenceItem, Client, User (ClientAdmin/ClientUser), AuditLog |
| 9: Reporting & Archive | Engagement, Report, ReportVersion, Workpaper, WorkpaperVersion, EvidenceItem, AuditLog |
| 10: EQR | EngagementQualityReview, AIDecision, AuditLog |
| 11: Multi-Framework Integrated Engagement | Engagement, EngagementFramework, CommonControl, FrameworkRequirement, SamplingWindow, AIDecision |
| 12: Continuous Assurance | EvidenceItem, CommonControl, DriftEvent, RiskRegisterEntry, Engagement, AIDecision, AuditLog |

## Regulatory Constraints by Journey

| Constraint | Standard | Journeys |
|-----------|----------|----------|
| Client acceptance before fieldwork | SQMS 1 (SOC) / ISO 17021-1 §9.1–9.4 (ISO) / PCI DSS scoping worksheet | 3, 11 |
| EQR reviewer independence | SQMS 2 (SOC) / ISO 17021-1 §9.6 (ISO) / firm policy | 3, 10 |
| Framework version locked after fieldwork begins | Section 4 requirement | 3, 11 |
| Four-level sign-off hierarchy enforced at data layer (Tester → Detailed Reviewer → General Reviewer → Final Reviewer; supersession on rework) | SQMS 1, ISO 17021-1, ISAE 3000 (Revised), firm quality frameworks | 5, 6 |
| AI draft must be edited before sign-off | ISO 42001 HITL discipline, AICPA AT-C AI guidance | 5 |
| Review notes cannot be deleted | AICPA AT-C, ISO 17021-1, HIPAA §164.312(b), PCI §12.10 | 6 |
| Period coverage check per framework | AT-C 320 (SOC 2 Type II), ISO surveillance cycle, PCI ASV 90-day validity, HIPAA risk-analysis refresh | 4, 7, 11, 12 |
| All AI decisions logged as AIDecision records | ISO 42001 HITL, AICPA AT-C AI guidance | 3, 4, 5, 7, 11, 12 |
| Sampling windows reconciled across frameworks | AT-C 205 (SOC 2), ISO 17021-1 sampling, PCI ASV | 11 |
| Assembly deadline enforcement | AT-C §A.60 (SOC), ISO CB file-closure policy, PCI SSC ROC assembly | 9 |
| WORM archival | AT-C retention, ISO 17021-1 records retention, HIPAA §164.316, PCI §12.10 | 9 |
| Retention periods per engagement type | AICPA SOC / ISO / HIPAA / PCI | 9 |
| Addenda require documented reason and partner sign-off | AICPA AT-C §A.60, ISO 17021-1 | 9 |
| Client upload tokens expire and require re-generation | Security policy | 7, 8 |
| Drift alerts propagated to auditor for material changes | ISO 27001 Annex A.5.24 (incident communication analog), SOC 2 CC2.3 (internal comms) | 12 |

## AI Touchpoints by Journey

| AI Feature | Tier | Journey | Trigger |
|-----------|------|---------|---------|
| Control mapping (cross-framework, grounded in SCF/OSCAL/AICPA) | Tier 2 — auditor reviews | 3, 11 | Engagement creation |
| Risk category suggestions (client acceptance) | Tier 2 — auditor certifies | 3 | Client acceptance form opened |
| Evidence classification and CommonControl mapping | Tier 2 — ClientAdmin/auditor reviews | 4 | Evidence uploaded |
| Framework-requirement coverage evaluation (full / partial / not) | Tier 2 — auditor confirms | 4, 5, 11 | Evidence mapped |
| Period-aware freshness evaluation | Tier 1 — informational + alerting | 4, 12 | Scheduled + on evidence change |
| Evidence link suggestions | Tier 2 — auditor accepts/rejects | 5 | Test procedure opened |
| Workpaper narrative draft | Tier 2 — auditor must edit | 5 | Explicit request after procedure complete |
| Document completeness review | Tier 2 — auditor must action | 7 | Client uploads document |
| Evidence link suggestion (on acceptance) | Tier 2 — auditor accepts/rejects | 7 | Document accepted |
| Report section draft | Tier 2 — partner must edit | 9 | Explicit request in report editor |
| Sampling window reconciliation across frameworks | Tier 2 — partner reviews | 11 | Multi-framework engagement scoped |
| Drift detection on config / evidence changes | Tier 1 — alerting, then Tier 2 re-test | 12 | Continuous monitoring signal |
| Agentic management-response drafting | Tier 2 — ClientAdmin + auditor review | 12 | Drift event requires remediation |
