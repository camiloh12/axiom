# Axiom User Journeys

Task-level user journey maps for every major platform workflow, organized by persona and goal. Referenced from Sections 11–12 of the [product spec](../specs/axiom-spec-design.md).

Competitor workflows (Fieldguide, CaseWare, DataSnipper, Yak) informed the pain points and opportunities throughout. Where Axiom introduces flows with no competitor equivalent, those are called out explicitly.

---

## Persona Reference

| Role | Description | Journeys |
|------|-------------|----------|
| **FirmAdmin** | Managing partner or designated admin. Configures the firm, manages billing and staff. | 1 |
| **Partner** | Engagement partner. Creates engagements, approves quality documentation, signs off on reports. 10+ years experience. | 3, 9 |
| **Manager** | Engagement manager. Reviews staff work, manages review notes, advances engagement phases. 5–8 years experience. | 6 |
| **Staff Auditor** | Performs hands-on audit work — trial balance, controls, workpapers, evidence collection. 1–3 years experience. | 2, 4, 5, 7 |
| **EQR Reviewer** | Independent quality reviewer per SQMS 2 / PCAOB AS 1220. Not on the engagement team. | 10 |
| **Client Contact** | CFO, controller, or IT manager at the audited company. Non-auditor. | 8 |
| **ClientAdmin** | Elevated client role. Can delegate document requests to colleagues. | 8 |

---

## Journey Index

| # | Persona | Goal | Replaces |
|---|---------|------|----------|
| 1 | FirmAdmin | Set up firm and launch first engagement | Section 11A |
| 2 | Staff Auditor | Join platform and reach first task | Section 11B |
| 3 | Partner | Create and scope a new engagement | Module 1 |
| 4 | Staff Auditor | Import and analyze a trial balance | Module 2 |
| 5 | Staff Auditor | Test controls and prepare workpapers | Module 3 (execution) |
| 6 | Manager | Review workpapers and advance the engagement | Module 3 (review) |
| 7 | Staff Auditor | Manage document requests and collect evidence | Module 4 (auditor side) |
| 8 | Client Contact | Fulfill audit document requests | Section 11C + Module 4 (client side) |
| 9 | Partner | Generate report, finalize, and archive | Module 5 |
| 10 | EQR Reviewer | Conduct engagement quality review | Module 3 (EQR) + Module 1 (EQR assignment) |

---

# Journey 1: Set Up Firm and Launch First Engagement

## Overview
- **Persona:** FirmAdmin — managing partner or IT-responsible partner at a 20–60 person CPA firm
- **Goal:** Get the firm from trial signup to a first active engagement in under one week, without outside help
- **Trigger:** Post-busy-season frustration with CaseWare + DataSnipper + Excel toolstack, a Fieldguide quote that was too expensive, or a DataSnipper renewal that prompted the question "why are we paying for a productivity layer on top of Excel?"
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
- Completes a brief intake form: firm name, staff count (dropdown: 1–10, 11–20, 21–40, 41–60, 60+), primary audit types (multi-select: Financial Audit, SOC 2, ISO 27001, HIPAA), country (US / Canada)

### Touchpoints
- Marketing site → signup form
- Verification email (delivered within 60 seconds)
- Intake form (single page, under 2 minutes to complete)

### Thoughts & Emotions
- **Cautious optimism** — "Let me see if this actually works without a sales call"
- **Mild anxiety** — "Is this going to be another Fieldguide situation where I spend two weeks in workshops before I see anything?"
- **Relief** — verification and intake take under 3 minutes total

### Pain Points
- **Competitor context:** Fieldguide requires a demo request and sales process before any platform access. CaseWare requires partner-channel procurement and admin configuration. Yak is invoice-only with no self-serve signup.
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
- **Competitor context:** CaseWare onboarding requires configuring security roles, user groups, and access hierarchies before anyone can use the platform. This creates hours of admin work before any audit work happens.
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
- Reviews a list of pre-built methodology templates: AICPA/GAAS Financial Audit, SOC 2 Type I/II (TSC 2017), ISO 27001:2022, HIPAA
- Activates one or more templates with a single click each — no configuration required to activate
- Growth tier: templates are read-only pre-built. Scale tier: custom template editor is unlocked for later customization

### Touchpoints
- Template selection screen with framework descriptions and engagement type labels
- One-click activation (toggle or checkbox per template)
- Onboarding checklist advances to step 3/5

### Thoughts & Emotions
- **Pleasantly surprised** — "I don't have to build my own methodology from scratch? CaseWare made me configure everything."
- **Confident** — selecting known frameworks (AICPA, SOC 2) feels safe and familiar
- **Curious** — "What's in these templates? Can I customize them later?"

### Pain Points
- **Competitor context:** CaseWare firms spend significant time configuring methodology templates and control libraries. Fieldguide's Accelerator workshops include template configuration as a dedicated phase. Yak only offers SOC 1/2 and HIPAA — no financial audit templates at all.
- If the admin can't preview what's inside a template before activating it, they'll hesitate
- Firms with custom methodologies may feel constrained by read-only templates on Growth tier

### Opportunities
- Allow template preview (expand to see control count, test procedure count, workpaper types) without requiring activation
- Show a "What's included" breakdown: "This SOC 2 Type II template includes 89 controls across all 5 trust services criteria, 200+ test procedures, and 80+ pre-drafted document request templates"
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
- **Excited** — "I just created a full engagement file in 2 minutes. This took me half a day in CaseWare."
- **Impressed** — seeing 89 controls and 200+ test procedures already populated feels like the product is working for them
- **Slightly overwhelmed** — the engagement dashboard shows a lot of content; the first reaction might be "where do I start?"

### Pain Points
- **Competitor context:** Fieldguide engagement setup involves framework selection, work program configuration, and team assignment — functional but more steps. CaseWare requires selecting a template file and then populating engagement properties, milestone dates, and team assignments across multiple screens.
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
- **Nervous about adoption** — "Will my staff actually use this? The CaseWare transition was painful."
- **Ready to delegate** — can now assign staff to the engagement and start real work

### Pain Points
- **Competitor context:** Fieldguide onboarding includes dedicated staff training sessions as part of the Accelerator. CaseWare training costs $150–$850 per course. DataSnipper requires per-seat licensing and individual Excel add-in installations.
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
Starts with cautious optimism ("is this really self-serve?"), builds through each fast step ("this is working"), peaks at engagement creation ("a full audit file in 2 minutes"), and settles into confident readiness at staff invitation. The emotional trajectory is deliberately the opposite of CaseWare (weeks of setup dread) and Fieldguide (weeks of workshops before seeing real value).

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
- **Persona:** Staff Auditor — 1–3 years experience, comfortable with Excel, accustomed to CaseWare + DataSnipper workflows
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
- **Slightly skeptical** — "Another new tool? I just got used to CaseWare."
- **Relieved** — magic link means no account creation form, no confirming passwords, no "choose a username"
- **Wary** — "How long is this going to take before I can do actual work?"

### Pain Points
- **Competitor context:** CaseWare requires firm admin to configure user accounts, security roles, and group assignments. DataSnipper requires individual Excel add-in installation and license activation.
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
- **Competitor context:** CaseWare has a steep learning curve with no guided tour — new users report spending days figuring out the interface. Fieldguide's onboarding includes dedicated training sessions, which means new staff wait for the next session.
- Generic product tours that explain UI elements without connecting to audit concepts are useless to auditors
- If the tour blocks interaction ("you must complete step 3 before you can click anything else"), it creates frustration

### Opportunities
- Frame each tour step in auditor language, not product language: "This is your evidence pool — every document your clients upload and every artifact you create lives here, linked to the controls they support"
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
- **Motivated** — "OK, I have real work to do. Let me see how this compares to CaseWare."
- **Confident** — the deep-link to a specific task means they don't have to navigate a complex engagement file to find their assignment
- **Focused** — transition from "learning the tool" to "doing audit work"

### Pain Points
- **Competitor context:** In CaseWare, new staff must navigate the Document Manager hierarchy to find their assigned workpapers. There is no notification system — assignment is communicated verbally or via email outside the platform.
- If the notification arrives before the user has completed the tour, they may skip orientation entirely and get lost later
- If no engagement assignment comes quickly, the user has nothing to do and may not return

### Opportunities
- Time the first assignment to arrive during or immediately after the tour — coordinate with the FirmAdmin/Partner
- The notification deep-link should open the control or workpaper with contextual help: "This is a SOC 2 CC6.1 control test. Here's what you need to do."

---

## Journey Summary

### Emotional Arc
Starts with skepticism ("another tool"), moves to relief (magic link, fast setup), through mild engagement (tour), and peaks at the first real assignment. The critical moment is the transition from "learning" to "working" — if the user reaches real audit work within 15 minutes of clicking the magic link, adoption is likely. If they're still in setup mode after 30 minutes, they'll mentally file Axiom as "another CaseWare learning curve."

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

Users with the EQReviewer role follow the same onboarding path but land in a read-only view of any engagement they are assigned to review. They cannot be added to the same engagement as a team member — the system enforces this at assignment time (Journey 10 covers the EQR workflow in detail).

---

# Journey 3: Create and Scope a New Engagement

## Overview
- **Persona:** Partner — audit partner, 10+ years experience, responsible for engagement quality and SQMS 1 compliance
- **Goal:** Set up a new engagement with the correct framework, team, quality documentation, and (if applicable) EQR assignment so fieldwork can begin
- **Trigger:** New client signed, or recurring engagement period starting (annual SOC 2, financial audit renewal)
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
- Selects engagement type: FinancialAudit_Private, FinancialAudit_Public, SOC1, SOC2, ISO27001, HIPAA, AgreedUponProcedures, or Advisory
- Selects the applicable MethodologyTemplate (filtered by engagement type)
- For multi-framework engagements (e.g., integrated SOC 2 + ISO 27001), selects a primary framework and adds secondary frameworks via EngagementFramework
- Selects framework version (locked after Fieldwork begins unless Partner override with documented reason)

### Touchpoints
- Engagement creation wizard — step 1: type and framework
- Framework version selector with effective dates shown
- Multi-framework toggle (Scale tier: shows cross-framework evidence mapping benefit)

### Thoughts & Emotions
- **Deliberate** — framework selection is a consequential decision; the partner is careful here
- **Confident** — seeing familiar framework names (AICPA/GAAS, SOC 2 TSC 2017) with version dates
- **Curious** (for multi-framework) — "How does the cross-framework mapping actually work?"

### Pain Points
- **Competitor context:** Fieldguide supports many frameworks but multi-framework setup requires separate configuration. CaseWare uses template files that must be selected correctly — wrong template means re-doing engagement setup. Yak only supports SOC 1/2 and HIPAA.
- Framework version confusion — partners need to see clearly which version applies to the current period
- Multi-framework engagement setup must not feel like creating two separate engagements

### Opportunities
- Auto-suggest framework version based on engagement period dates
- For multi-framework, show a preview: "This integrated SOC 2 + ISO 27001 engagement will include X controls that satisfy both frameworks simultaneously"
- Rollforward detection: if the same client had a prior engagement, prompt "Roll forward from [prior engagement name]?"

---

## Stage 2: Configure Engagement Details and Assign Team

### Sub-goal
Define the engagement scope and get the right people assigned.

### User Actions
- Enters client name (autocomplete from existing Client records, or creates a new one)
- Sets audit period dates (period_start, period_end)
- For rollforward engagements: prior_engagement_id is set, and the system surfaces prior year data
- Assigns engagement team: selects users by name, assigns engagement-level roles
- System creates: Engagement (Planning status), EngagementTeamMember records, EngagementFramework records, Control records cloned from template, TestProcedure records, draft Workpaper shells, empty ClientAcceptance record

### Touchpoints
- Engagement wizard — step 2: client and period
- Team assignment panel with user search and role assignment
- Scaffold generation progress indicator
- For rollforward: prior year sidebar showing controls, workpapers, TB reference, and prior ClientAcceptance (read-only)

### Thoughts & Emotions
- **Efficient** — client autocomplete and pre-populated templates save significant setup time
- **Pleased** (rollforward) — seeing prior year data carried forward means less re-work
- **Responsible** — team assignment is a delegation decision; the partner is deciding who owns what

### Pain Points
- **Competitor context:** CaseWare rollforward is strong (mappings carry forward, prior year balances reference automatically) but is a desktop-era workflow. Fieldguide's rollforward surfaces prior year work but users report the setup can still take 30+ minutes for complex engagements.
- If scaffold generation is slow for large templates (200+ controls), the partner waits
- Rollforward from prior year needs clear visual distinction between "carried forward" and "new this year"

### Opportunities
- Show rollforward diff: "142 controls carried forward from prior year. 3 new controls added in the updated SOC 2 TSC 2017 template. 1 prior control marked as superseded."
- Allow team assignment by engagement role pattern: "Same team as last year" button for rollforward engagements
- Background scaffold generation with notification when complete, if it takes >3 seconds

---

## Stage 3: Review AI-Proposed Control Mappings

### Sub-goal
Validate the AI's proposed mappings between firm control objectives and framework requirements. This is Axiom's cross-framework differentiator in action.

### User Actions
- Immediately after scaffold creation, AI proposes FirmControlObjectiveMapping records across all in-scope frameworks
- Partner or Manager reviews a mapping table: each firm control objective is shown with proposed FrameworkRequirement links, confidence scores, and explanation text
- Reviews in bulk — all proposed mappings are confirmed by default; the partner rejects or modifies individual mappings as needed
- Each confirmed mapping creates an AIDecision record

### Touchpoints
- Control mapping review screen — sortable table with control name, proposed framework links, confidence percentage, and AI explanation
- Bulk confirm button with count: "Confirm all 142 mappings"
- Individual reject/modify actions per mapping
- Low-confidence mappings highlighted for priority review

### Thoughts & Emotions
- **Impressed** — "It mapped 142 controls to SOC 2 and ISO 27001 simultaneously. This would have taken me a full day."
- **Professionally cautious** — "I need to check the low-confidence ones. My name goes on this engagement."
- **Delighted** (for multi-framework) — seeing one control satisfy CC6.1 + A.8.3 + §164.312(a)(1) in a single row is the "aha moment" for the cross-framework value proposition

### Pain Points
- **Competitor context:** No competitor does this. Fieldguide has framework-specific control mapping but not cross-framework. CaseWare has no AI mapping. Yak has "automated control mapping" but only for SOC 1/2 and HIPAA. DataSnipper's Document Matching is evidence-level, not control-level.
- If the AI confidence is consistently low, bulk review becomes tedious rather than efficient
- Partners unfamiliar with AI-assisted workflows may not trust the suggestions and review every single mapping individually, defeating the time savings

### Opportunities
- Show accuracy metrics: "Based on prior engagements at your firm, AI control mapping accuracy is 94%"
- Group review by confidence tier: high confidence (>0.85) in a bulk-confirm block, medium (0.65–0.85) for quick scan, low (<0.65) for individual attention
- For rollforward engagements, pre-load prior year confirmed mappings as the starting suggestion (re-confirmation still required)

---

## Stage 4: Complete SQMS 1 Client Acceptance

### Sub-goal
Document quality risks and acceptance decision before fieldwork can begin — a regulatory gate, not a platform feature.

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
- **Confident** — the system enforces the gate, which means the firm can demonstrate SQMS 1 compliance to regulators

### Pain Points
- **Competitor context:** Most competitors don't enforce client acceptance as a system gate. CaseWare leaves it to firm policy. Fieldguide includes engagement acceptance but enforcement varies by configuration. The system-enforced gate is an Axiom differentiator for quality management.
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
- For PCAOB engagements (mandatory) or higher-risk nonissuer engagements (per firm policy), the partner opens the EQR assignment panel
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
- Show a clear indicator: "EQR is required for this engagement (PCAOB)" or "EQR is optional for this engagement (nonissuer)"
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
- Staff auditors can now begin control testing (Journey 5), trial balance work (Journey 4), and document requests (Journey 7)

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
- If any guard fails (acceptance incomplete, EQR not assigned for PCAOB engagement), the partner needs clear guidance on what's missing — not just a disabled button
- The transition notification should be actionable for staff — "Fieldwork has begun. Your assigned controls are ready for testing."

### Opportunities
- Show a pre-transition summary: "Advancing to Fieldwork will notify 4 team members and unlock 89 controls for testing. Proceed?"
- After transition, redirect the partner to the engagement dashboard with a Fieldwork-phase overview: progress by control, team workload distribution

---

## Handoffs

| From | To | Information Transferred | Trigger |
|------|-----|------------------------|---------|
| Partner | Staff Auditor (Journey 4, 5, 7) | Engagement scaffold with assigned controls and procedures | Fieldwork transition |
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

# Journey 4: Import and Analyze a Trial Balance

## Overview
- **Persona:** Staff Auditor — assigned to a financial audit engagement, responsible for trial balance work
- **Goal:** Import the client's trial balance, map accounts to financial statement line items, generate lead schedules, and prepare the analytical foundation for substantive testing
- **Trigger:** Engagement is in Fieldwork status and the client has provided a TB export from their accounting system
- **Stages:**
  1. Import trial balance data
  2. Review AI account mappings
  3. Generate lead schedules
  4. Propose and track adjustments
  5. Run analytical procedures
  6. Prepare population for sampling

## Stage 1: Import Trial Balance Data

### Sub-goal
Get the client's trial balance into the platform in a clean, structured format.

### User Actions
- Navigates to the Trial Balance section of the engagement
- Uploads a CSV or Excel file exported from the client's accounting system (QBO, NetSuite, Sage, Xero)
- The importer recognizes common column formats (account number, account name, debit balance, credit balance) with configurable column mapping
- Reviews the import preview — row count, account count, total debits/credits
- Confirms the import — TrialBalanceAccount records are created

### Touchpoints
- Trial Balance section with "Import" button
- File upload (drag-and-drop or file picker)
- Column mapping screen with auto-detected mappings and manual override
- Import preview with summary statistics
- Sheets UI populated with imported accounts

### Thoughts & Emotions
- **Focused** — TB import is foundational; errors here cascade through the entire engagement
- **Hopeful** — "Please recognize the column format so I don't have to map every field manually"
- **Relieved** — if auto-detection works cleanly on the first try

### Pain Points
- **Competitor context:** CaseWare handles TB import from 60+ accounting systems with robust mapping — this is their strongest feature. Fieldguide recently added TB import but users report buggy imports with misaligned columns. DataSnipper uses Table Snip to extract TB from PDFs, then pivot tables in Excel for analysis.
- Non-standard column formats (especially from smaller accounting systems) require manual mapping every time
- If the import fails silently or partially (some rows imported, others dropped), the auditor may not notice until later

### Opportunities
- Save column mapping profiles per accounting system: "This looks like a NetSuite export. Use NetSuite column mapping?"
- Show a pre-confirmation validation: "142 accounts imported. Total debits: $X. Total credits: $Y. Difference: $0.00." Any non-zero difference is a red flag.
- For rollforward engagements, show prior year TB as a reference panel: "Prior year had 138 accounts. 4 new accounts detected."

---

## Stage 2: Review AI Account Mappings

### Sub-goal
Classify each trial balance account into a standard financial statement line item using AI assistance.

### User Actions
- Immediately after import, Claude Haiku classifies each account into a FS line item (Cash, Accounts Receivable, Fixed Assets, Revenue, etc.)
- Mappings appear in the Sheets UI with mapping_status = AISuggested
- Low-confidence mappings are highlighted in yellow; unmapped accounts in red
- The auditor reviews and confirms each mapping — bulk-confirm available for high-confidence mappings
- Prior year confirmed mappings are pre-loaded on rollforward engagements as starting suggestions

### Touchpoints
- Sheets UI with AI mapping column (suggested FS line item + confidence indicator)
- Bulk-confirm button: "Confirm all high-confidence mappings (128 of 142)"
- Individual override dropdown per account
- Mapping status indicators: AISuggested (blue), Confirmed (green), Overridden (orange)

### Thoughts & Emotions
- **Impressed** — "It mapped 128 of 142 accounts correctly in seconds. This usually takes me 2 hours."
- **Vigilant** — low-confidence accounts need manual attention; the auditor takes professional responsibility
- **Comfortable** — the Sheets UI feels familiar (like Excel) rather than alien

### Pain Points
- **Competitor context:** CaseWare has manual account mapping with carry-forward from prior year — accurate but time-consuming. No AI assistance. DataSnipper has no account mapping. Fieldguide's AI TB mapping is new and users report accuracy issues on non-standard charts of accounts.
- If AI accuracy is below 90%, bulk confirm becomes risky and the auditor reviews everything individually
- Non-standard account names (abbreviations, foreign languages, industry jargon) reduce AI accuracy

### Opportunities
- Show per-firm accuracy improvement over time: "AI mapping accuracy for your firm: 91% (up from 84% last quarter)"
- For rollforward, pre-load prior year confirmed mappings — auditor re-confirms but starts from a known baseline
- Allow custom mapping rules at the firm level: "Accounts starting with '4xxx' are always Revenue"

---

## Stage 3: Generate Lead Schedules

### Sub-goal
Produce summary workpapers that aggregate trial balance accounts by financial statement section.

### User Actions
- Once account mappings are confirmed, clicks "Generate Lead Schedules"
- System automatically creates Workpaper records of type LeadSchedule, grouped by FS section (Current Assets, Fixed Assets, Liabilities, Equity, Revenue, Expenses)
- Each lead schedule aggregates balances from mapped TrialBalanceAccount records with drill-through to individual accounts
- Materiality calculations (ISA 320 / AU-C 320 methods) are available as formula functions within the Sheets UI

### Touchpoints
- "Generate Lead Schedules" button (available once mappings are confirmed)
- Lead schedule workpapers in the engagement file — one per FS section
- Drill-through links from lead schedule totals to individual accounts
- Materiality calculator (overall materiality, performance materiality, clearly trivial threshold)

### Thoughts & Emotions
- **Delighted** — lead schedules that would take 2–4 hours to build in Excel are generated instantly
- **Trusting** — the drill-through to individual accounts lets the auditor verify that totals are correct
- **Professional** — materiality calculations are built in rather than maintained in a separate spreadsheet

### Pain Points
- **Competitor context:** CaseWare auto-generates lead schedules from TB groupings and maintains live links — this is mature, well-tested functionality. DataSnipper has no lead schedule capability. Fieldguide added this recently.
- Lead schedules must update dynamically when adjustments are posted — if they're static snapshots, they become stale
- Materiality calculation methods differ between ISA and AICPA; the calculator must support both

### Opportunities
- Live lead schedules that update in real-time as the TB changes (adjustments posted, accounts remapped)
- Show prior year lead schedule comparison side-by-side for rollforward engagements
- Flag line items that exceed materiality thresholds automatically

---

## Stage 4: Propose and Track Adjustments

### Sub-goal
Document proposed, passed, and waived audit adjustments and track their impact on the trial balance.

### User Actions
- Identifies an account requiring adjustment (from analytical procedures or substantive testing)
- Creates a TrialBalanceAdjustment: specifies account, amount, description, and type (Proposed, Passed, or Waived)
- Manager or Partner reviews and approves or waives each adjustment
- The Sheets UI tracks both unadjusted and adjusted trial balance — both views are available
- Lead schedules reflect adjustments in real-time

### Touchpoints
- Adjustment entry form within the Sheets UI
- Adjustment tracking panel: list of all proposed adjustments with status
- Dual-view TB: unadjusted and adjusted columns
- Approval workflow: Manager/Partner action per adjustment

### Thoughts & Emotions
- **Methodical** — adjustments are consequential; the auditor documents carefully
- **Reassured** — seeing both unadjusted and adjusted balances side-by-side prevents confusion
- **Collaborative** — the approval workflow means the auditor isn't making decisions alone

### Pain Points
- **Competitor context:** CaseWare has dedicated adjustment worksheets with auto-propagation to lead schedules and financial statements. This is mature. DataSnipper doesn't handle adjustments.
- If adjustments don't propagate to lead schedules and financial statements in real-time, the auditor must manually verify consistency
- Waived adjustments must be documented per AU-C (all proposed adjustments, whether passed or not, must be retained) — the system must prevent deletion

### Opportunities
- Aggregate waived adjustments and show total impact: "Total waived adjustments: $45,000 (below materiality threshold of $50,000)"
- Auto-flag when cumulative waived adjustments approach or exceed materiality
- Require a documented reason for each waived adjustment (regulatory requirement)

---

## Stage 5: Run Analytical Procedures

### Sub-goal
Identify unusual account activity, compute financial ratios, and flag anomalies for further investigation.

### User Actions
- Opens the Analytics panel in the Sheets UI
- Reviews computed analytics: period-over-period variance by account, ratio calculations (current ratio, quick ratio, debt-to-equity), trend analysis
- Reviews AI-generated anomaly flags: accounts with unusual activity relative to prior period or industry norms
- Investigates flagged accounts by drilling into transaction detail
- Documents analytical procedure results in a workpaper

### Touchpoints
- Analytics panel with computed metrics and visualizations
- Anomaly flag indicators on individual accounts (AI-generated, Tier 1: informational only)
- Drill-through to GL transaction detail (if imported)
- Analytical procedures workpaper template

### Thoughts & Emotions
- **Analytical** — this is the intellectually engaging part of TB work
- **Curious** — AI anomaly flags may surface things the auditor wouldn't have noticed
- **Careful** — anomaly flags are starting points for investigation, not conclusions

### Pain Points
- **Competitor context:** CaseWare provides variance analysis through linked spreadsheets. DataSnipper's Document Matching can cross-reference TB to financial statements. Neither has AI anomaly detection. Fieldguide offers AI-powered analytics but specifics on financial audit analytics are limited.
- If anomaly flags are too frequent (noisy), the auditor ignores them; if too rare, they miss real issues
- PCAOB engagements require all technology-assisted analytical procedures to be documented as AIDecision records

### Opportunities
- Configurable anomaly sensitivity per engagement: "Flag variances >10% from prior year" or "Flag variances >$X"
- Industry benchmarking: "This client's current ratio (1.2) is below the industry median (1.8)"
- Auto-generate the analytical procedures workpaper from the analytics panel output

---

## Stage 6: Prepare Population for Sampling

### Sub-goal
Export transaction populations and compute statistically valid sample selections for substantive testing.

### User Actions
- Exports GL transaction detail as a population listing
- Opens the sampling calculator — selects sampling method (systematic, random, monetary unit per ISA 530 / AU-C 530)
- Enters parameters: population size, expected misstatement rate, confidence level
- System computes sample size and generates a sample selection
- The sample selection is recorded as a TestProcedure sample selection record

### Touchpoints
- Population export from GL transaction detail
- Sampling calculator with method selection and parameter inputs
- Sample selection results: list of selected transactions with amounts
- TestProcedure record updated with population_size, sample_size, sampling_method

### Thoughts & Emotions
- **Technical** — sampling is a statistical exercise; the auditor wants precision
- **Confident** — built-in ISA 530 / AU-C 530 methods mean they don't need to compute manually in Excel
- **Relieved** — the sample selection is documented automatically, satisfying documentation requirements

### Pain Points
- **Competitor context:** CaseWare has sampling tools but they're often supplemented with Excel. DataSnipper has no sampling capability. Fieldguide supports statistical sampling. Yak has no sampling or financial audit tools.
- If the GL transaction detail wasn't imported (only summary TB), population-based sampling isn't available — this must be communicated clearly
- PCAOB documentation requirements for sampling are strict: population size, sample size, selection method, and results must all be recorded per AS 2315

### Opportunities
- Support both summary-level TB import and detailed GL import — make the distinction clear at import time
- For population analysis (Axiom differentiator): offer full-population analytics as an alternative to sampling where the data supports it — "Test the entire population of 4,200 transactions instead of a sample of 60"
- Save sampling parameters as templates for recurring engagement types

---

## Handoffs

| From | To | Information Transferred | Trigger |
|------|-----|------------------------|---------|
| Staff Auditor | Manager/Partner | Proposed adjustments requiring approval | Adjustment created |
| Staff Auditor | Staff Auditor (Journey 5) | Mapped TB, lead schedules, and sample selections | TB work complete |

---

## Journey Summary

### Emotional Arc
Opens with the tension of import ("will the data come in clean?"), climbs through the satisfying AI mapping moment ("128 of 142 in seconds"), peaks at lead schedule generation ("this took me 4 hours in Excel"), and sustains through the methodical adjustment and analytics work. The emotional through-line is "the platform handles the tedious parts so I can focus on judgment."

### Cross-Cutting Pain Points
- Import quality determines everything downstream — garbage in, garbage out
- Real-time propagation (TB → mappings → lead schedules → adjustments → analytics) must be seamless
- The gap between CaseWare's 20-year TB maturity and a new platform's TB is the biggest adoption risk for financial audit firms

### Prioritized Opportunities
1. **Full-population analytics as alternative to sampling** (high impact, high effort) — a genuine innovation over all competitors
2. **Saved column mapping profiles per accounting system** (high impact, low effort) — eliminates repeated import friction
3. **Live lead schedules with adjustment propagation** (high impact, medium effort) — must match CaseWare's auto-propagation to be credible
4. **AI accuracy tracking per firm** (medium impact, low effort) — builds trust in AI mapping over time

---

# Journey 5: Test Controls and Prepare Workpapers

## Overview
- **Persona:** Staff Auditor — assigned to specific controls within a financial or compliance audit engagement
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
- **Competitor context:** In CaseWare, staff must navigate the Document Manager to find assigned workpapers — there's no personalized assignment view. Fieldguide has assignment dashboards with real-time status. DataSnipper has no assignment concept — it's a tool, not a workflow.
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
- Documents the population and sample (if applicable — from Journey 4 sampling)
- Records results: what was observed, measured, or confirmed
- Notes any exceptions found
- Documents conclusions per procedure
- Changes status to Complete (or Exception if exceptions found)

### Touchpoints
- Test procedure detail view with structured fields
- Population/sample reference (linked from TB sampling in Journey 4)
- Results entry: structured fields + narrative area
- Exception documentation panel
- Status progression: NotStarted → InProgress → Complete / Exception

### Thoughts & Emotions
- **Focused** — this is the core audit work; the auditor is in professional execution mode
- **Methodical** — structured fields guide the documentation without being overly rigid
- **Concerned** (when exceptions found) — "I need to document this carefully and escalate if significant"

### Pain Points
- **Competitor context:** Fieldguide's AI Field Agents can automate up to 70% of testing — validating journal entries, performing substantive tests, summarizing results, flagging exceptions. CaseWare test documentation is manual but structured. DataSnipper automates the evidence-matching part but not the procedure documentation.
- Procedure documentation that is too structured (all dropdowns, no narrative) doesn't capture the auditor's professional judgment
- Procedure documentation that is too freeform (blank text box) leads to inconsistent quality across the team

### Opportunities
- Balanced structure: required fields (procedure type, population, sample, result, exceptions) plus narrative area for professional judgment
- Exception escalation: when exceptions are noted, prompt "Notify manager?" with a one-click notification
- AI-assisted population testing (Axiom differentiator): for procedures that support it, offer "Analyze full population" alongside traditional sampling

---

## Stage 3: Link Evidence to Test Procedures

### Sub-goal
Connect supporting documents to test procedures to demonstrate that conclusions are evidence-based.

### User Actions
- From the test procedure view, opens the evidence linking panel
- Browses the engagement's evidence pool (all EvidenceItem records for this client, across all engagements)
- AI may have already suggested evidence links (EvidenceLink.ai_suggested = true) — the auditor accepts, modifies, or rejects
- Links selected evidence items to the test procedure
- For multi-framework engagements: the UI shows which other framework requirements that evidence simultaneously satisfies

### Touchpoints
- Evidence linking panel with search and browse
- AI-suggested links with confidence scores and "Accept / Modify / Reject" actions
- Cross-framework satisfaction display: "This evidence satisfies SOC 2 CC6.1, ISO 27001 A.8.3, and HIPAA §164.312(a)(1)"
- EvidenceLink record created with linked_by_id and timestamp
- AIDecision record updated on acceptance/modification/rejection

### Thoughts & Emotions
- **Efficient** — AI suggestions mean the auditor doesn't have to search through hundreds of evidence items manually
- **Delighted** (multi-framework) — seeing one piece of evidence satisfy three frameworks is the payoff of the cross-framework architecture
- **Professionally responsible** — accepting an AI suggestion still requires the auditor to verify that the evidence is relevant and sufficient

### Pain Points
- **Competitor context:** DataSnipper's Document Matching automates evidence-to-workpaper cross-referencing within Excel — this is their core strength, used by 500,000+ auditors. Fieldguide has evidence linking but the "no split-view" pain point means auditors can't see evidence and workpapers side-by-side. CaseWare requires manual file linking in the Document Manager.
- If the evidence pool is large (hundreds of items across multiple engagements), search and filtering are critical
- AI suggestions that are wrong erode trust in the system — the auditor starts ignoring suggestions

### Opportunities
- Split-view: evidence document viewer on one side, test procedure on the other (addresses Fieldguide's top pain point)
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
- **Competitor context:** Fieldguide's Field Assist automates workpaper narrative drafting with 50–75% time savings reported. CaseWare has no AI drafting. DataSnipper has no workpaper narrative capability. Yak has no workpaper drafting.
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
- **Competitor context:** Fieldguide uses real-time collaborative editing. CaseWare's check-in/check-out means only one person can edit at a time — a major collaboration bottleneck. DataSnipper embeds evidence references in Excel but has no workpaper editor.
- If the editor is clunky (poor formatting, no undo, slow saves), it undermines the entire AI-assisted workflow
- Version history must be clear — the auditor needs to see what changed between versions, especially to demonstrate human edits on AI drafts

### Opportunities
- Real-time collaborative editing (not check-in/check-out) for team workpapers
- Auto-save with version history — no explicit "save" action needed, every change creates a recoverable version
- Evidence reference sidebar that updates contextually based on cursor position in the workpaper

---

## Stage 6: Submit Workpaper for Review

### Sub-goal
Mark the workpaper as ready for manager review, triggering the sign-off workflow.

### User Actions
- Reviews the workpaper one final time
- Clicks "Submit for Review" — status changes to PreparedPendingReview
- System validates: is_ai_draft must be false (if an AI draft was generated, the auditor must have edited it)
- The assigned manager receives a notification
- The workpaper is now locked for the preparer — only the reviewer can make changes or return it

### Touchpoints
- "Submit for Review" button with validation check
- Validation error if is_ai_draft = true: "This workpaper contains an unedited AI draft. You must review and edit the AI-generated content before submitting."
- Notification to assigned manager
- Workpaper status badge updates

### Thoughts & Emotions
- **Accomplished** — a completed workpaper represents meaningful professional output
- **Hopeful** — "I hope the manager doesn't have too many review notes"
- **Trusting** — the system validation (AI draft check) gives the auditor confidence that they're not submitting incomplete work

### Pain Points
- **Competitor context:** Fieldguide has digital sign-off workflow. CaseWare has electronic sign-off but the workflow varies by firm configuration. No competitor enforces the AI-draft-must-be-edited gate — this is an Axiom regulatory compliance differentiator (PCAOB AS 1105).
- If the manager review takes days (a common audit bottleneck), the auditor's momentum stalls
- The "locked for preparer" state must be communicated clearly — the auditor needs to know they can't make further edits until the manager returns it or clears it

### Opportunities
- Show estimated review queue time: "Manager [name] has 4 workpapers ahead of yours in the review queue"
- Allow the auditor to add notes for the reviewer: "FYI — this is a new control this year. Prior year test procedure was different."

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
1. **Split-view evidence + workpaper** (high impact, medium effort) — addresses Fieldguide's most-cited pain point
2. **Cross-framework evidence satisfaction display** (high impact, already designed) — the core architectural differentiator experienced at the task level
3. **Firm-specific AI draft style matching** (high impact, high effort) — the difference between a generic draft and a useful one
4. **Review queue visibility** (medium impact, low effort) — transparency reduces the waiting frustration

---

# Journey 6: Review Workpapers and Advance the Engagement

## Overview
- **Persona:** Manager — 5–8 years experience, responsible for reviewing staff work and managing the engagement's progression through its phases
- **Goal:** Review all submitted workpapers, provide feedback via review notes, clear completed workpapers for partner sign-off, and advance the engagement from Fieldwork through Review to Reporting
- **Trigger:** Staff auditor submits a workpaper for review (PreparedPendingReview status)
- **Stages:**
  1. Triage the review queue
  2. Review workpaper content and evidence
  3. Raise and manage review notes
  4. Clear workpaper for partner sign-off
  5. Advance engagement status

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
- **Competitor context:** Fieldguide's engagement dashboard gives managers real-time visibility but lacks bulk review capabilities. CaseWare's review workflow is per-document with no cross-engagement queue. DataSnipper has no review workflow.
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
- **Competitor context:** Fieldguide lacks split-view for evidence alongside workpapers — reviewers must switch between tabs. CaseWare's Document Manager requires navigating between the workpaper and evidence files. DataSnipper's embedded hyperlinks in Excel provide instant cross-reference — the one area where the Excel paradigm actually works well for reviewers.
- If the evidence viewer is a separate window or requires navigation, the review takes twice as long
- Reviewing AI-drafted content requires additional attention — the manager must ensure the auditor actually applied judgment, not just accepted the AI output

### Opportunities
- Split-view review mode: workpaper on the left, evidence on the right (inspired by DataSnipper's inline hyperlinks but in a modern UI)
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
- Resolved notes remain in the record — they cannot be deleted (AU-C 230)

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
- **Competitor context:** CaseWare tracks "Issues" (review notes) linked to documents. Fieldguide has inline commenting. Both support open/resolved status tracking. None prevent note deletion — Axiom's immutability is a compliance differentiator.
- If review notes are too informal (just text comments), they lack the structure needed for regulatory documentation
- The back-and-forth cycle (note → response → resolution) can be slow if notifications are delayed

### Opportunities
- Note categorization for pattern detection: "This manager raised 14 'missing evidence description' notes this quarter — consider adding to the staff training curriculum"
- One-click note response templates: "Evidence updated" / "Narrative revised" / "Disagree — see response below"
- Real-time notification when review notes are added or responded to

---

## Stage 4: Clear Workpaper for Partner Sign-Off

### Sub-goal
Advance the workpaper from InReview to ReviewComplete, signaling to the partner that the manager's review is satisfied.

### User Actions
- Confirms all review notes are resolved (open notes block advancement)
- Changes workpaper status to ReviewComplete
- The sign-off action creates a timestamped, named AuditLog entry — it cannot be backdated
- The engagement partner is notified that the workpaper is ready for final sign-off
- The partner reviews the manager's cleared notes and signs off the workpaper (SignedOff status)

### Touchpoints
- Review status advancement button (disabled if open notes exist)
- Sign-off confirmation with name and timestamp
- Notification to engagement partner
- AuditLog entry: "Workpaper [name] reviewed by [Manager] at [timestamp]"

### Thoughts & Emotions
- **Confident** — the workpaper meets professional standards
- **Relieved** — one more item cleared from the review queue
- **Professional pride** — the sign-off carries their reputation

### Pain Points
- **Competitor context:** CaseWare has configurable sign-off schemes but the enforcement varies. Fieldguide has digital sign-off. No competitor enforces the sign-off hierarchy at the data layer as strictly as Axiom specifies (cannot advance states out of order — SQMS 1, AU-C 220).
- If the partner sign-off takes additional days, the manager's work sits idle
- Bulk sign-off scenarios (many workpapers completing simultaneously) need efficiency without sacrificing individual attention

### Opportunities
- Partner review queue mirroring the manager review queue — cross-engagement visibility
- Allow the partner to batch-sign low-risk workpapers with individual confirmation: "Sign off on 8 workpapers? Each will be individually timestamped."

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
| Manager | Staff Auditor (Journey 5) | Review notes requiring response | Notes raised |
| Manager | Partner (Journey 9) | Reviewed workpaper ready for sign-off | ReviewComplete status |
| Manager | EQR Reviewer (Journey 10) | Engagement ready for quality review | Review phase reached |

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
- **Competitor context:** Fieldguide's AI Field Agents auto-generate PBC request lists from engagement context — practitioners approve before sending. CaseWare's PBC Requests is a separate add-on product. DataSnipper's UpLink provides document request tracking. Yak has basic request creation.
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
- **Competitor context:** Fieldguide sends requests through its Client Hub. DataSnipper's UpLink provides no-login upload links. CaseWare requires the PBC add-on. Yak has "one-click document upload" for clients.
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
- **Competitor context:** Fieldguide has automated reminders and real-time status tracking. CaseWare PBC has notification capability but as a separate product. DataSnipper's UpLink tracks request status. AuditDashboard (standalone product) has a "trigger fieldwork at 80% completion" AI feature.
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
- **Competitor context:** Fieldguide's Request Analysis Agent validates documents on upload and flags gaps. DataSnipper's UpLink has AI-powered document pre-validation. CaseWare has no AI document review. Yak's AI Reviewer evaluates documents for key information.
- If the AI frequently misjudges completeness (too many false Accept or false Reject recommendations), the auditor stops trusting the queue
- The "Send Back" action must communicate clearly to the client what's missing — a generic "incomplete" message generates confusion and more back-and-forth
- For SOC 2 Type II: period coverage check is mandatory (AT-C 320) and must be part of every AI review

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
- **Persona:** Client Contact — CFO, controller, or IT manager at the audited company. Not an auditor. May have limited time and multiple people who hold different pieces of information.
- **Goal:** Understand what documents the audit team needs, upload them correctly the first time, and get through the PBC process with minimal disruption to daily work
- **Trigger:** Email notification: "Your audit team has started your [Engagement Name] and needs documents from you."
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
- **Competitor context:** Fieldguide's Client Hub requires the client to see a full dashboard. DataSnipper's UpLink provides no-login upload links — the closest competitor experience. CaseWare's PBC portal is a separate add-on that may or may not be configured.
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
- **Competitor context:** Fieldguide organizes requests in a dashboard with progress tracking. DataSnipper's UpLink presents requests with descriptions. Yak has "simple, one-click document upload."
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
- **Competitor context:** All competitors support basic file upload. DataSnipper's UpLink and Fieldguide's Client Hub both provide clean upload experiences. CaseWare varies by configuration.
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
- **Competitor context:** Fieldguide's Client Hub allows clients to delegate tasks to team members within the portal. This delegation model (single-request scoped) is more restrictive but more secure.
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
- **Competitor context:** Fieldguide allows clients to comment on report sections and collaborate in-platform. DataSnipper's UpLink provides status tracking. CaseWare's PBC portal has basic notification capability.
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
- Selects report type: SOC 2 Type I, SOC 2 Type II, SOC 1 Type I, SOC 1 Type II, Financial Audit Opinion, Agreed-Upon Procedures, or Management Letter
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
- **Competitor context:** Fieldguide has one-click report generation from templates, but users report "reports take a lot of work to create and templates can be finicky." CaseWare automates financial statement drafting and report wording. Yak has "one-click report generation" for compliance reports.
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
- **Competitor context:** Fieldguide allows report editing in Microsoft Word before delivery. CaseWare generates reports within its template system. None offer real-time collaborative report editing.
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
- **Competitor context:** Fieldguide allows clients to collaborate on report drafts directly in-platform. CaseWare has no built-in client report sharing.
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
- System computes automatically: assembly_deadline (report date + 60 days for AICPA, + 45 days for PCAOB), retention_deadline (report date + 5 years for AICPA/SOC/ISO, + 7 years for PCAOB, + 6 years for HIPAA)
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
- **Competitor context:** CaseWare computes and enforces assembly deadlines via its Lockdown feature. Fieldguide has archiving but no mention of WORM or assembly deadline enforcement. Yak has no archiving capability mentioned.
- If the partner issues accidentally (premature click), reversal is complex (must withdraw and re-issue)
- The distinction between AICPA and PCAOB deadline windows must be clear

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
- **Competitor context:** CaseWare's Lockdown is irreversible and creates a permanent read-only copy — well-established. No competitor has the addendum workflow that Axiom specifies (new WorkpaperVersion with is_addendum = true and partner sign-off per AU-C 230 §.16).
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
- At any time: FirmAdmin can generate a complete engagement export (all workpapers as PDF, evidence files in native format, TB as Excel, AuditLog as CSV, metadata as JSON — structured ZIP file)

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
- **Competitor context:** CaseWare's Lockdown creates a point-in-time copy but doesn't use WORM storage. No competitor mentions S3 Object Lock COMPLIANCE mode. This is a strong regulatory differentiator.
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
- Report template quality determines the partner's experience — finicky templates are the most-cited Fieldguide complaint
- The issuance → finalization → archival chain must be flawless — errors here have regulatory consequences
- Engagement export is a critical trust feature — firms need to know they own their data independent of the platform

### Prioritized Opportunities
1. **Automatic assembly deadline computation and enforcement** (high impact, medium effort) — a genuine regulatory differentiator over all competitors
2. **S3 WORM archival** (high impact, high effort) — strongest archival guarantee in the market
3. **Addendum workflow** (medium impact, medium effort) — proper AU-C 230 §.16 implementation; no competitor has this
4. **Engagement export** (high impact, medium effort) — critical for trust, offboarding, and regulatory compliance

---

# Journey 10: Conduct Engagement Quality Review

## Overview
- **Persona:** EQR Reviewer — experienced partner or external reviewer, independent of the engagement team, assigned per SQMS 2 / PCAOB AS 1220
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
- Sees all engagement content: workpapers, evidence, controls, test procedures, trial balance, review notes, AI decisions, audit log
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
- **Competitor context:** Few competitors have a dedicated EQR workflow. Fieldguide supports EQR but specifics are limited. CaseWare leaves EQR to firm-level processes. No competitor enforces the independence check (reviewer ≠ team member) at the system level.
- If the read-only mode doesn't provide sufficient navigation tools, the reviewer spends too long finding things
- The reviewer needs to see everything without being overwhelmed — they're reviewing, not performing the audit

### Opportunities
- EQR-focused view: pre-organized sections aligned with SQMS 2 review scope (planning quality, risk assessment, testing sufficiency, documentation adequacy, significant judgments)
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
1. **EQR-focused navigation view** (high impact, medium effort) — aligns the review interface with SQMS 2 review scope rather than the engagement's natural structure
2. **AI edit substantiveness indicator** (high impact, low effort) — surfaces the most important quality concern in AI-assisted engagements
3. **Structured findings with action tracking** (medium impact, medium effort) — streamlines the reviewer-to-team communication cycle
4. **Historical EQR comparison** (low impact, low effort) — useful context but not critical

---

# Cross-Journey Insights

## Competitor Pain Points Addressed by Axiom

| Competitor Pain Point | Source | Axiom Response | Journeys |
|----------------------|--------|---------------|----------|
| Multi-week implementation with dedicated consultant | Fieldguide Accelerator | Self-serve onboarding, <1 week to first engagement | 1 |
| Steep learning curve, costly training | CaseWare ($150–$850/course) | 5-step guided tour, intuitive UI, deep-link to first task | 2 |
| No split-view for evidence alongside workpapers | Fieldguide (top G2 complaint) | Split-view evidence panel during testing and review | 5, 6 |
| Check-in/check-out file locking | CaseWare desktop | Real-time collaborative editing | 5 |
| PBC as separate add-on product | CaseWare PBC Requests | Integrated document request management | 7, 8 |
| Excel dependency for all audit work | DataSnipper (lives inside Excel) | Native web platform with sheets-like TB experience | 4 |
| Report templates are finicky | Fieldguide (G2 reviews) | Pre-populated templates with engagement data | 9 |
| No AI layer | CaseWare | AI across 4 features: completeness review, control mapping, TB mapping, workpaper draft | 3, 4, 5, 7 |
| Per-seat pricing penalizes mid-market | DataSnipper, Fieldguide | Unlimited users at all tiers | 1 |
| No financial audit support | Yak (SOC 1/2 and HIPAA only) | Full TB, lead schedules, sampling, GAAS workpapers | 4 |
| No cross-framework evidence mapping | All competitors | Framework-agnostic evidence architecture | 3, 5 |
| No system-enforced regulatory gates | Most competitors | SQMS 1 acceptance gate, EQR independence check, assembly deadline enforcement, AI draft edit requirement | 3, 9, 10 |

## Flows Without Competitor Equivalent

These Axiom flows represent genuine innovation — no competitor currently offers them:

1. **Cross-framework evidence satisfaction display** (Journeys 3, 5) — one piece of evidence simultaneously satisfying SOC 2, ISO 27001, and HIPAA requirements, shown in real-time during testing
2. **AI document completeness review with client-facing gap explanations** (Journeys 7, 8) — AI analyzes uploaded documents against request requirements and auto-drafts specific gap explanations for the client
3. **System-enforced AI draft edit gate** (Journeys 5, 6) — workpapers with unedited AI content cannot be signed off, implementing PCAOB AS 1105 at the data layer
4. **EQR independence enforcement** (Journey 10) — system-level validation that the quality reviewer is not on the engagement team
5. **Post-finalization addendum workflow** (Journey 9) — proper AU-C 230 §.16 implementation with immutable original content and versioned addenda
6. **Full-population analytics as alternative to sampling** (Journey 4) — testing entire transaction datasets rather than statistical samples
7. **Automatic assembly deadline computation and WORM archival** (Journey 9) — computed at report issuance, enforced via S3 Object Lock COMPLIANCE mode

## Data Entities by Journey

| Journey | Primary Entities |
|---------|-----------------|
| 1: Firm Setup | Firm, User, MethodologyTemplate, Engagement, Control, TestProcedure, Workpaper, ClientAcceptance |
| 2: Staff Onboarding | User, EngagementTeamMember |
| 3: Engagement Scoping | Engagement, EngagementTeamMember, EngagementFramework, Control, TestProcedure, ClientAcceptance, EngagementQualityReview, FirmControlObjectiveMapping, AIDecision |
| 4: Trial Balance | TrialBalance, TrialBalanceAccount, TrialBalanceAdjustment, Workpaper, AIDecision |
| 5: Control Testing | Control, TestProcedure, EvidenceItem, EvidenceLink, Workpaper, WorkpaperVersion, AIDecision |
| 6: Workpaper Review | Workpaper, WorkpaperVersion, AuditLog |
| 7: Document Requests | DocumentRequest, EvidenceItem, EvidenceLink, AIDecision, AuditLog |
| 8: Client Fulfillment | DocumentRequest, EvidenceItem, Client, User (ClientAdmin/ClientUser), AuditLog |
| 9: Reporting & Archive | Engagement, Report, ReportVersion, Workpaper, WorkpaperVersion, EvidenceItem, AuditLog |
| 10: EQR | EngagementQualityReview, AIDecision, AuditLog |

## Regulatory Constraints by Journey

| Constraint | Standard | Journeys |
|-----------|----------|----------|
| Client acceptance before fieldwork | SQMS 1 | 3 |
| EQR reviewer independence | SQMS 2 / PCAOB AS 1220 | 3, 10 |
| Framework version locked after fieldwork begins | Section 4 requirement | 3 |
| Sign-off hierarchy enforced at data layer | SQMS 1, AU-C 220 | 5, 6 |
| AI draft must be edited before sign-off | PCAOB AS 1105 | 5 |
| Review notes cannot be deleted | AU-C 230 | 6 |
| Period coverage check for SOC 2 Type II evidence | AT-C 320 | 7 |
| All AI decisions logged as AIDecision records | PCAOB AS 1105 | 3, 4, 5, 7 |
| Sampling documentation requirements | AU-C 530 / PCAOB AS 2315 | 4 |
| Assembly deadline enforcement | AU-C 230, PCAOB AS 1215 | 9 |
| WORM archival | PCAOB AS 1215, SOX §802 | 9 |
| Retention periods per engagement type | AU-C 230, PCAOB, HIPAA | 9 |
| Addenda require documented reason and partner sign-off | AU-C 230 §.16 | 9 |
| Client upload tokens expire and require re-generation | Security policy | 7, 8 |

## AI Touchpoints by Journey

| AI Feature | Tier | Journey | Trigger |
|-----------|------|---------|---------|
| Control mapping (cross-framework) | Tier 2 — auditor reviews | 3 | Engagement creation |
| Risk category suggestions (client acceptance) | Tier 2 — auditor certifies | 3 | Client acceptance form opened |
| Trial balance account mapping | Tier 2 — auditor confirms | 4 | TB imported |
| Anomaly detection (unusual account activity) | Tier 1 — informational only | 4 | Nightly background job |
| Evidence link suggestions | Tier 2 — auditor accepts/rejects | 5 | Test procedure opened |
| Workpaper narrative draft | Tier 2 — auditor must edit | 5 | Explicit request after procedure complete |
| Document completeness review | Tier 2 — auditor must action | 7 | Client uploads document |
| Evidence link suggestion (on acceptance) | Tier 2 — auditor accepts/rejects | 7 | Document accepted |
| Report section draft | Tier 2 — partner must edit | 9 | Explicit request in report editor |
