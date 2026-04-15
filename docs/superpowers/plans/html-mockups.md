# Axiom UI Mockups — HTML Generation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Generate high-fidelity HTML mockups for every key screen across all 10 Axiom user journeys, organized by journey and stage, so stakeholders can visualize the full UX before any React code is written.

**Architecture:** Each screen is a standalone `.html` file. Files are organized by journey number under `mockups/`. Design decisions (typography, colors, spacing, layout patterns, components) are driven by **[impeccable.style](https://impeccable.style)** — an AI design skill that establishes a project-specific design language via `/impeccable teach` and generates production-grade UI via `/impeccable craft`. The mockups depict realistic audit data (real field names, domain terminology, plausible sample data) drawn from the product spec and user journeys.

**Tech Stack:** HTML + CSS (impeccable decides the CSS approach — Tailwind, plain CSS, or hybrid). No JavaScript except optional navigation links between screens in the same journey. Impeccable's design references (OKLCH colors, modular type scales, 4pt spacing) are the design system — not a hand-rolled one.

---

## Context for Fresh Sessions

This plan generates UI mockups for **Axiom**, an AI-native audit engagement platform for mid-market CPA firms. The full product spec, data model, and user journeys are in this repository:

| Document | Path | What it contains |
|----------|------|------------------|
| Product spec | `docs/specs/axiom-spec-design.md` | Full product design, data model, pricing, regulatory requirements |
| User journeys | `docs/user-journeys/all-journeys.md` | 10 detailed user journey maps with stages, touchpoints, and screen descriptions |
| Backend architecture | `docs/specs/backend-architecture-design.md` | Service decomposition, API design |
| Domain model | `docs/specs/domain-and-data-model-design.md` | All entities, relationships, enums |

**Before generating any screen, read the corresponding journey and stage in `docs/user-journeys/all-journeys.md`.** The journey document describes exactly what each screen shows, what data is visible, what actions are available, and what emotional state the user is in. Use that as the design brief.

**User roles:** FirmAdmin | Partner | Manager | Staff | EQReviewer | ClientAdmin | ClientUser | ViewOnly

---

## Design System — Powered by Impeccable

**Do NOT hand-roll a design system.** All visual design decisions (typography, colors, spacing, layout, component styling) are handled by the [impeccable.style](https://impeccable.style) AI design skill. Impeccable must be installed and configured before any screens are generated.

### How Impeccable Works

Impeccable is an AI design skill (open source, Apache 2.0) that:
1. Runs a design discovery interview via `/impeccable teach` and saves the result as `.impeccable.md`
2. Generates production-grade UI via `/impeccable craft` following a shape-then-build workflow
3. Provides refinement commands: `/audit`, `/polish`, `/typeset`, `/layout`, `/colorize`, etc.
4. Enforces anti-patterns (no gradient text, no AI color palettes, no glassmorphism, etc.)
5. Uses modern CSS principles: OKLCH color spaces, modular type scales, 4pt spacing, container queries

### Installation (Task 1 handles this)

```bash
npx skills add pbakaus/impeccable
```

This installs the skill into `.claude/skills/`. After installation, run `/impeccable teach` to establish Axiom's design context. The teach interview will ask about brand, audience, and aesthetic direction — answer with the context below.

### Context for `/impeccable teach`

When the teach interview asks, provide these answers:

- **Product:** Axiom — an AI-native audit engagement platform for mid-market CPA firms (20–60 staff)
- **Users:** Audit professionals (partners, managers, staff auditors with 1–15 years experience) and their clients (CFOs, controllers, IT managers). Daily users during busy season. Desktop-primary.
- **Brand personality:** Trustworthy, precise, efficient. Not playful. Not enterprise-heavy. Think "modern professional tool" — closer to Linear or Notion than to Salesforce or Oracle.
- **Aesthetic direction:** Clean, information-dense but not cluttered. Data tables, status indicators, and structured forms are the primary UI patterns. The platform handles sensitive financial and compliance data — the design should convey reliability and professionalism.
- **Anti-references:** Generic SaaS dashboards with hero gradients. Anything that looks like a template. Enterprise software with dense gray chrome.
- **Key UI patterns:** App shell with sidebar navigation, data tables, form wizards, split-view panels, status badges, progress indicators, AI-generated content markers.
- **Special layout:** The Client Hub (Journey 8) is a separate, minimal layout — no sidebar, centered narrow width — designed for non-auditor clients who arrive via tokenized email links.

After `/impeccable teach` runs, the resulting `.impeccable.md` file becomes the design system. All subsequent `/impeccable craft` calls will read it automatically.

### Sidebar Navigation Structure

Impeccable handles the visual styling, but the **information architecture** of the sidebar is domain-specific. Auditor-facing screens should include these nav items:

```
Dashboard
Engagements
─ [Engagement Name] (when inside an engagement context)
  ├ Overview
  ├ Controls
  ├ Trial Balance
  ├ Workpapers
  ├ Evidence
  ├ Document Requests
  ├ Reports
  └ Settings
Staff Management
Firm Settings
```

### Sample Data Conventions

Use consistent, realistic sample data across all screens:

- **Firm:** Meridian & Associates CPAs
- **Users:** Sarah Chen (Partner), James Rodriguez (Manager), Emily Park (Staff), David Kim (Staff), Lisa Nguyen (EQR Reviewer)
- **Client:** Cloudvault Technologies Inc.
- **Engagement:** SOC 2 Type II, period Jan 1 – Dec 31, 2025
- **Secondary client:** Greenfield Manufacturing LLC (Financial Audit, FY ending Mar 31, 2026)
- **Framework controls:** Use real SOC 2 TSC references (CC6.1, CC6.2, CC7.1, etc.) and ISO 27001 Annex A references (A.5.1, A.8.1, etc.)
- **Trial balance accounts:** Use realistic account names and numbers (1000 Cash, 1200 Accounts Receivable, 4000 Revenue, etc.)
- **Dollar amounts:** Use plausible mid-market figures ($2M–$50M revenue range)

---

## Directory Structure

```
mockups/
├── README.md                          # Overview + screenshot instructions
├── journey-01-firm-setup/
│   ├── 01-signup-form.html
│   ├── 02-email-verification.html
│   ├── 03-intake-form.html
│   ├── 04-firm-profile.html
│   ├── 05-methodology-templates.html
│   ├── 06-create-engagement.html
│   ├── 07-engagement-ready.html
│   ├── 08-invite-staff.html
│   └── 09-onboarding-complete.html
├── journey-02-staff-onboarding/
│   ├── 01-invitation-email.html
│   ├── 02-magic-link-landing.html
│   ├── 03-profile-setup.html
│   ├── 04-guided-tour-step.html
│   └── 05-first-assignment.html
├── journey-03-engagement-scoping/
│   ├── 01-new-engagement-type.html
│   ├── 02-engagement-details.html
│   ├── 03-team-assignment.html
│   ├── 04-ai-control-mapping.html
│   ├── 05-client-acceptance.html
│   ├── 06-eqr-assignment.html
│   └── 07-begin-fieldwork.html
├── journey-04-trial-balance/
│   ├── 01-tb-import.html
│   ├── 02-column-mapping.html
│   ├── 03-ai-account-mapping.html
│   ├── 04-lead-schedules.html
│   ├── 05-adjustments.html
│   ├── 06-analytics.html
│   └── 07-sampling.html
├── journey-05-control-testing/
│   ├── 01-my-assignments.html
│   ├── 02-control-detail.html
│   ├── 03-test-procedure.html
│   ├── 04-evidence-linking.html
│   ├── 05-ai-workpaper-draft.html
│   ├── 06-workpaper-editor.html
│   └── 07-submit-for-review.html
├── journey-06-workpaper-review/
│   ├── 01-review-queue.html
│   ├── 02-workpaper-review.html
│   ├── 03-review-notes.html
│   ├── 04-sign-off.html
│   └── 05-engagement-progress.html
├── journey-07-document-requests/
│   ├── 01-create-requests.html
│   ├── 02-request-dashboard.html
│   ├── 03-ai-review-queue.html
│   └── 04-evidence-acceptance.html
├── journey-08-client-hub/
│   ├── 01-client-landing.html
│   ├── 02-request-list.html
│   ├── 03-upload-interface.html
│   ├── 04-delegate-request.html
│   ├── 05-sent-back-feedback.html
│   └── 06-completion-view.html
├── journey-09-reporting/
│   ├── 01-report-generation.html
│   ├── 02-report-editor.html
│   ├── 03-client-draft-review.html
│   ├── 04-issue-report.html
│   ├── 05-finalize-engagement.html
│   └── 06-archive-confirmation.html
└── journey-10-eqr/
    ├── 01-eqr-notification.html
    ├── 02-read-only-engagement.html
    ├── 03-planning-review.html
    ├── 04-testing-review.html
    ├── 05-findings.html
    └── 06-eqr-signoff.html
```

**Total: 61 screens across 10 journeys.**

---

## Screen Inventory

Each screen below maps to a specific stage in the user journeys. The journey document (`docs/user-journeys/all-journeys.md`) contains the full design brief for each — read the corresponding section before generating the HTML.

### Journey 1: Firm Setup (FirmAdmin)

| # | File | Journey Stage | What the Screen Shows |
|---|------|--------------|----------------------|
| 1 | `01-signup-form.html` | Stage 1 | Marketing-style signup page. Business email input, "Start Free Trial" button. Clean, minimal — no sidebar, centered layout. Show "No credit card required" and "14-day full-feature trial." |
| 2 | `02-email-verification.html` | Stage 1 | Interstitial: "Check your inbox" with email address shown, spam folder guidance, countdown "Your trial starts now — 14 days." No sidebar. |
| 3 | `03-intake-form.html` | Stage 1 | Single-page form: firm name, staff count dropdown (1–10, 11–20, 21–40, 41–60, 60+), primary audit types multi-select (Financial Audit, SOC 2, ISO 27001, HIPAA), country (US/Canada). No sidebar. |
| 4 | `04-firm-profile.html` | Stage 2 | First screen with the app shell (sidebar + header). Onboarding wizard step 2/5. Four fields: firm display name, logo upload (drag-and-drop with firm initials default), timezone (auto-detected), billing email. Collapsible progress panel in sidebar. |
| 5 | `05-methodology-templates.html` | Stage 3 | Onboarding step 3/5. Card grid showing 4 templates: AICPA/GAAS Financial Audit, SOC 2 Type I/II, ISO 27001:2022, HIPAA. Each card has framework name, description, control count, test procedure count, and a one-click Activate toggle. One template expanded to show preview. |
| 6 | `06-create-engagement.html` | Stage 4 | Onboarding step 4/5. Engagement creation wizard: engagement type dropdown, framework selector (pre-filtered), client name input, audit period date picker, methodology template (pre-selected). Show scaffold generation progress. |
| 7 | `07-engagement-ready.html` | Stage 4 | "Your engagement is ready" confirmation. Visual summary: 89 controls, 214 test procedures, 67 workpaper shells. "Start here" prompt pointing to client acceptance. Engagement dashboard preview behind the confirmation overlay. |
| 8 | `08-invite-staff.html` | Stage 5 | Onboarding step 5/5. Bulk email input area (paste-friendly), role assignment dropdown per invitee with role descriptions. Show 5 invitees with varied roles. "Send Invitations" button. |
| 9 | `09-onboarding-complete.html` | Stage 5 | Completion state. Onboarding checklist with all 5 steps checked green. "Your firm is set up. Here's what's next:" with action cards: Assign team to engagement, Begin fieldwork, Explore the demo engagement. |

### Journey 2: Staff Onboarding (Staff Auditor)

| # | File | Journey Stage | What the Screen Shows |
|---|------|--------------|----------------------|
| 1 | `01-invitation-email.html` | Stage 1 | Email template mockup: firm logo, "You've been invited to Meridian & Associates' Axiom workspace as Staff Auditor", prominent "Accept Invitation" button, magic link. Centered, no app shell. |
| 2 | `02-magic-link-landing.html` | Stage 1 | SSO linking screen: "Welcome to Axiom" header, SSO options (Google, Microsoft), or "Set a password" alternative. Terms of service checkbox. No sidebar. |
| 3 | `03-profile-setup.html` | Stage 1–2 | Profile setup: display name input, role badge "Staff Auditor" (non-editable) with explanation text, notification preferences (real-time/daily/weekly radio, daily pre-selected). Single page, 3 fields. |
| 4 | `04-guided-tour-step.html` | Stage 3 | App shell with a tour overlay. Show step 3/5 highlighting the evidence pool section. Spotlight effect on the evidence area, tooltip with description in auditor language: "This is your evidence pool — every document your clients upload and every artifact you create lives here, linked to the controls they support." Step indicators (1/5 through 5/5). |
| 5 | `05-first-assignment.html` | Stage 4 | App shell with engagement list. One engagement highlighted with "New" badge. Toast notification: "You've been assigned to Cloudvault Technologies SOC 2 Type II." Deep-link button to first assigned control. |

### Journey 3: Engagement Scoping (Partner)

| # | File | Journey Stage | What the Screen Shows |
|---|------|--------------|----------------------|
| 1 | `01-new-engagement-type.html` | Stage 1 | Engagement creation wizard step 1. Engagement type grid: FinancialAudit_Private, SOC2, ISO27001, HIPAA, etc. Framework version selector with effective dates. Multi-framework toggle with explanation (Scale tier). SOC 2 selected, ISO 27001 added as secondary. |
| 2 | `02-engagement-details.html` | Stage 2 | Wizard step 2. Client name autocomplete (showing "Cloudvault Technologies" match), audit period dates, rollforward detection banner: "Roll forward from Cloudvault SOC 2 Type II 2024?" |
| 3 | `03-team-assignment.html` | Stage 2 | Wizard step 3. Team assignment panel: user search, role assignment per member. Show 4 team members assigned. "Same team as last year" button for rollforward. |
| 4 | `04-ai-control-mapping.html` | Stage 3 | Full-width table: firm control objectives mapped to framework requirements. Columns: Control Name, SOC 2 Requirement, ISO 27001 Requirement, Confidence, AI Explanation. 142 rows visible (paginated). Bulk confirm button: "Confirm all 142 mappings." High-confidence in green, medium in amber, low in red. This is the "aha moment" screen — show one row satisfying CC6.1 + A.8.3 simultaneously. |
| 5 | `05-client-acceptance.html` | Stage 4 | ClientAcceptance form. Structured quality risk categories (dropdown + free text), firm responses per risk, independence confirmation checkbox with auditor name. Sign button with timestamp preview. Status indicator: "Planning — Acceptance Required." |
| 6 | `06-eqr-assignment.html` | Stage 5 | EQR panel. Reviewer dropdown (filtered to EQReviewers not on team). Independence validation indicator (green checkmark). Confirmation message. Show "EQR is required for this engagement (PCAOB)" indicator. |
| 7 | `07-begin-fieldwork.html` | Stage 6 | Engagement readiness checklist: all items green (acceptance, team, controls, EQR). "Begin Fieldwork" button enabled. Pre-transition summary: "Advancing to Fieldwork will notify 4 team members and unlock 89 controls for testing." |

### Journey 4: Trial Balance (Staff Auditor)

| # | File | Journey Stage | What the Screen Shows |
|---|------|--------------|----------------------|
| 1 | `01-tb-import.html` | Stage 1 | Trial Balance section with drag-and-drop upload zone. Supported formats listed. Import history table (empty for new engagement). |
| 2 | `02-column-mapping.html` | Stage 1 | Column mapping interface: file preview table showing raw CSV data. Auto-detected column mappings (account number, account name, debit, credit) with override dropdowns. Saved profile detection: "This looks like a QuickBooks export. Use QuickBooks mapping?" Import preview: 142 accounts, total debits/credits, difference $0.00. |
| 3 | `03-ai-account-mapping.html` | Stage 2 | Sheets-style UI: spreadsheet grid showing accounts with AI-mapped FS line items. Color-coded mapping status: blue (AI Suggested), green (Confirmed), orange (Overridden). Bulk confirm button: "Confirm all high-confidence mappings (128 of 142)." Confidence column with percentages. Prior year comparison sidebar for rollforward. |
| 4 | `04-lead-schedules.html` | Stage 3 | Lead schedule workpaper view. Sections: Current Assets, Fixed Assets, Liabilities, Equity, Revenue, Expenses. Each section shows aggregated balances from mapped accounts. Drill-through links. Materiality calculator panel (overall materiality, performance materiality, clearly trivial). |
| 5 | `05-adjustments.html` | Stage 4 | Dual-column TB view: Unadjusted and Adjusted side by side. Adjustment panel: list of proposed adjustments with status (Proposed, Approved, Waived). New adjustment form. Aggregate indicator: "Total waived adjustments: $45,000 (below materiality threshold of $50,000)." |
| 6 | `06-analytics.html` | Stage 5 | Analytics dashboard: period-over-period variance table, ratio calculations (current ratio, quick ratio, debt-to-equity), AI anomaly flags on specific accounts. Small data visualizations (bar charts for variance, trend lines). |
| 7 | `07-sampling.html` | Stage 6 | Sampling calculator: method selection (systematic, random, MUS), parameter inputs (population size, expected misstatement, confidence level), computed sample size. Sample selection results table. Population analytics alternative toggle: "Analyze full population of 4,200 transactions instead of sample of 60." |

### Journey 5: Control Testing (Staff Auditor)

| # | File | Journey Stage | What the Screen Shows |
|---|------|--------------|----------------------|
| 1 | `01-my-assignments.html` | Stage 1 | "My Assignments" personalized dashboard. Table of assigned controls: control name, framework requirements, procedure type, status, prior year status. Filter by engagement, status. Cross-framework badge on multi-framework controls: "SOC 2 CC6.1 + ISO A.8.3". |
| 2 | `02-control-detail.html` | Stage 1 | Control detail view. Description, control objective, linked framework requirements with cross-framework display. Associated test procedures list. Prior year reference panel (collapsible). Evidence requirements summary. |
| 3 | `03-test-procedure.html` | Stage 2 | Test procedure execution view. Structured fields: procedure type dropdown, population/sample reference, results text area, exceptions panel. Status: InProgress. Mix of structured and narrative input. Exception escalation: "Notify manager?" prompt. |
| 4 | `04-evidence-linking.html` | Stage 3 | Split-view: test procedure on left, evidence pool browser on right. AI-suggested links with confidence scores and Accept/Modify/Reject buttons. Cross-framework satisfaction display: "This evidence satisfies SOC 2 CC6.1, ISO 27001 A.8.3, and HIPAA 164.312(a)(1)." Search bar for evidence content. |
| 5 | `05-ai-workpaper-draft.html` | Stage 4 | Workpaper editor with "AI Draft — requires review" banner prominently displayed. AI-generated narrative text with slight visual distinction (light blue background or left border). "Generate AI Draft" button. WorkpaperVersion indicator. |
| 6 | `06-workpaper-editor.html` | Stage 5 | Full workpaper editor. Rich text with sections: objective, scope, procedures performed, results, conclusion. Evidence reference sidebar showing linked items. Version history panel. The AI draft banner is gone (human-edited). |
| 7 | `07-submit-for-review.html` | Stage 6 | Submission confirmation dialog over the workpaper: "Submit for Review?" Validation checklist: all test procedures complete, evidence linked, AI draft edited. Reviewer assignment showing manager name. Notes-for-reviewer text area. |

### Journey 6: Workpaper Review (Manager)

| # | File | Journey Stage | What the Screen Shows |
|---|------|--------------|----------------------|
| 1 | `01-review-queue.html` | Stage 1 | Cross-engagement review dashboard. Table: workpaper name, engagement, staff auditor, submission date, control area, priority. Engagement progress indicators ("34 of 89 controls reviewed"). Filter/sort controls. Staff notes visible in preview. |
| 2 | `02-workpaper-review.html` | Stage 2 | Split-view review mode: workpaper content on left, evidence viewer on right. Evidence sidebar with linked items (click to preview). Version history showing AI draft -> human edit transition with indicator: "AI-drafted. Auditor modified 5 of 8 sections." |
| 3 | `03-review-notes.html` | Stage 3 | Workpaper with inline review notes. Highlighted text with comment bubbles. Review notes panel: severity (question/suggestion/required change), open/resolved status. Resolution workflow shown. Note count: "3 open review notes." |
| 4 | `04-sign-off.html` | Stage 4 | Sign-off confirmation: "Clear workpaper for partner sign-off?" All review notes resolved (green checkmarks). Timestamped sign-off preview. AuditLog entry preview. |
| 5 | `05-engagement-progress.html` | Stage 5 | Engagement dashboard: control completion progress (pie or bar chart), phase transition readiness. Blocking items: "2 controls blocking Fieldwork -> Review. CC7.2 assigned to David Kim. CC8.1 assigned to Emily Park." Phase transition button. |

### Journey 7: Document Requests (Staff Auditor)

| # | File | Journey Stage | What the Screen Shows |
|---|------|--------------|----------------------|
| 1 | `01-create-requests.html` | Stage 1 | Document request creation. Bulk creation from template: "Generate SOC 2 Type II requests? 83 requests across 5 trust services categories." Individual request form: title, instructions, due date, client contact assignment. Request grouping by department. |
| 2 | `02-request-dashboard.html` | Stages 2–3 | Request tracking dashboard. Status per request: Pending, Submitted, InReview, Accepted, Rejected, Overdue. Due date indicators. Automated reminder configuration. Overdue count. Engagement readiness: "72% fulfilled. Estimated fieldwork start: May 15." |
| 3 | `03-ai-review-queue.html` | Stage 4 | AI-assessed upload review queue. Each item: document name, request title, AI recommendation (Accept/Request Clarification/Reject), gaps identified, confidence score. Sorted by low confidence first. One-click actions. Batch accept button for high-confidence items. |
| 4 | `04-evidence-acceptance.html` | Stage 5 | Evidence acceptance view. Accepted document with auto-created EvidenceLink. Cross-framework display. Request status: Accepted. "Send Back" option with AI-drafted gap explanation preview: "The uploaded access control policy covers January–September 2025 but the audit period extends to December 2025." |

### Journey 8: Client Hub (Client Contact) — NO SIDEBAR

| # | File | Journey Stage | What the Screen Shows |
|---|------|--------------|----------------------|
| 1 | `01-client-landing.html` | Stage 1 | Clean centered layout. Firm logo at top. "Meridian & Associates has requested documents for your SOC 2 Type II engagement." Engagement period, request count, progress bar. No login form visible — user arrived via tokenized link. |
| 2 | `02-request-list.html` | Stage 2 | Request list grouped by department: "IT Security (12 requests)", "HR Policies (8 requests)", "Financial (15 requests)." Each request: title, plain-language instructions, due date, status. Priority indicators. Progress: "5 of 23 complete." |
| 3 | `03-upload-interface.html` | Stage 3 | Single request expanded with drag-and-drop upload area. Request title and detailed instructions visible. File format guidance. Upload progress indicator. Confirmation message: "Document uploaded successfully. Your audit team will review it shortly." |
| 4 | `04-delegate-request.html` | Stage 4 | Delegation interface for ClientAdmin. Request detail with "Delegate" button. Email input for delegate. Custom message area: "Add a note for your colleague." Confirmation: "Request delegated to maria@cloudvault.com. They can only see this specific request." |
| 5 | `05-sent-back-feedback.html` | Stage 5 | Sent-back request with AI-drafted gap explanation. Original upload shown alongside the gap: "You uploaded access-control-policy-v3.pdf. Gaps identified: Policy covers January–September 2025 but the audit period extends to December 2025." Re-upload area. |
| 6 | `06-completion-view.html` | Stage 5 | All requests complete. Celebration state: "All 23 documents received. Your audit team is reviewing them." Progress 23/23 complete. Post-engagement read-only view with list of submitted documents. |

### Journey 9: Reporting & Archive (Partner)

| # | File | Journey Stage | What the Screen Shows |
|---|------|--------------|----------------------|
| 1 | `01-report-generation.html` | Stage 1 | Report type selector cards: SOC 2 Type II, Financial Audit Opinion, Management Letter, etc. Pre-populated template preview with engagement data. "Generate AI Draft" option for Description of Tests section. Report completeness checklist. |
| 2 | `02-report-editor.html` | Stage 2 | Report editor: rich text with section navigation sidebar. Version history panel. Internal sharing controls. The report shows realistic SOC 2 Type II content: scope, management assertion, description of tests, results, opinion. |
| 3 | `03-client-draft-review.html` | Stage 3 | Report in "Shared with Client" mode. Client comments inline. Comment panel with status: resolved, acknowledged, will-not-change. "Selective sharing" indicator. |
| 4 | `04-issue-report.html` | Stage 4 | Two-step issuance dialog: "You are about to issue this report." Computed deadlines: "Assembly deadline: March 1, 2026. Retention deadline: December 31, 2030." Warnings. "This action cannot be undone. Confirm?" |
| 5 | `05-finalize-engagement.html` | Stage 5 | Finalization checklist: all workpapers signed off, all review notes resolved, report issued, EQR complete. "Finalize Engagement" button. Warning: "All engagement content will be permanently locked. Modifications require an addendum." |
| 6 | `06-archive-confirmation.html` | Stage 6 | Archived engagement view. Read-only indicator. Archival details: "Archived on [date]. S3 WORM storage. Retention until December 31, 2030." Engagement export button: "Download complete engagement file (ZIP)." |

### Journey 10: EQR (EQR Reviewer)

| # | File | Journey Stage | What the Screen Shows |
|---|------|--------------|----------------------|
| 1 | `01-eqr-notification.html` | Stage 1 | Notification view: "Cloudvault Technologies SOC 2 Type II has reached Review status and is ready for your quality review." "Begin Review" button. |
| 2 | `02-read-only-engagement.html` | Stage 1 | Engagement view with prominent "EQR — Read Only" banner. Full navigation available but no edit controls visible. Yellow/amber tint on the header to visually distinguish read-only mode. |
| 3 | `03-planning-review.html` | Stage 2 | EQR-focused planning view. ClientAcceptance summary, team composition, AI control mapping decision trail: "42 AI decisions. 38 accepted, 3 modified, 1 rejected." Planning quality flags: no issues detected. |
| 4 | `04-testing-review.html` | Stage 3 | Testing sufficiency view. Workpaper sample with evidence links. AI edit indicator: "This workpaper's AI draft was edited. 5 of 8 sections modified (62.5%)." Exception summary: "3 controls with exceptions." |
| 5 | `05-findings.html` | Stage 4 | EQR findings form. Structured categories: observation/recommendation/required action. Findings list. Required action items with assignment to engagement team. Conclusion selector: satisfied / satisfied with concerns / not satisfied. |
| 6 | `06-eqr-signoff.html` | Stage 5 | Sign-off confirmation: "Sign off on Engagement Quality Review?" Summary of findings. All required actions addressed (if any). Timestamp and name. "This sign-off is immutable and will be archived with the engagement." |

---

## Execution Instructions

### Task Organization

Work through the tasks in order. Task 1 installs impeccable and establishes the design language. Tasks 2–11 generate the journey screens. Task 12 is the final review.

### Per-Screen Process

For each screen:

1. **Read the source:** Open `docs/user-journeys/all-journeys.md` and read the stage that corresponds to this screen. Pay attention to: User Actions, Touchpoints, and Opportunities sections — these describe what should appear on screen.
2. **Use `/impeccable craft`:** For each screen (or group of related screens), invoke `/impeccable craft` with a prompt describing the screen content. Impeccable will run its shape-then-build workflow. Provide:
   - The screen's purpose and user persona
   - What UI elements appear (tables, forms, wizards, etc.)
   - The sample data to populate it with
   - Whether it uses the app shell or a standalone layout (Client Hub, signup, email)
3. **Generate realistic content:** Use the sample data conventions from this plan. Fill tables with plausible rows (5–10 rows, not placeholders). Use real SOC 2 control references, real account names, realistic dollar amounts.
4. **Refine with impeccable commands:** After generation, run `/audit` to check for design issues and `/polish` for final refinement. Use `/typeset` if typography feels off, `/layout` for spatial issues.
5. **Write the HTML file:** Self-contained, opens in a browser with no build step.
6. **Verify in browser:** Open the file and visually inspect. Confirm the design feels cohesive with previously generated screens.

### Design Consistency Across Screens

After generating the first screen of the first journey (the app shell screen — `journey-01/04-firm-profile.html`), use `/impeccable extract` to identify reusable components and tokens. This establishes the shared primitives (nav shell, buttons, badges, tables, cards) that all subsequent screens should reuse. Save extracted components somewhere accessible so later screens reference them.

### Screenshot Process

After all mockups are generated, screenshots can be taken using Playwright:

```bash
# Install Playwright (one-time)
npm init -y
npm install playwright

# Screenshot script (create as screenshot.js)
node screenshot.js
```

The README in `mockups/` should include a screenshot script. Alternatively, screenshots can be taken manually in a browser.

---

## Tasks

### Task 1: Install Impeccable and Establish Design Language

**Files:**
- Create: `mockups/README.md`
- Created by impeccable: `.impeccable.md` (design config)
- Created by impeccable: `.claude/skills/` (skill files)

- [ ] **Step 1:** Install the impeccable skill.

```bash
npx skills add pbakaus/impeccable
```

- [ ] **Step 2:** Run `/impeccable teach` to establish Axiom's design context. When the teach interview asks questions, use the context provided in the "Context for `/impeccable teach`" section of this plan. The result is saved as `.impeccable.md`.

- [ ] **Step 3:** Review the generated `.impeccable.md` — confirm it captures Axiom's professional, information-dense, trustworthy aesthetic. If it doesn't, run `/impeccable teach` again with more specific input.

- [ ] **Step 4:** Create the `mockups/` directory and `mockups/README.md` explaining the structure, how to view mockups (open in browser), and how to take screenshots.

- [ ] **Step 5:** Commit.

```bash
git add .impeccable.md mockups/README.md
git commit -m "docs: add impeccable design config and mockups directory"
```

---

### Task 2: Journey 1 — Firm Setup (9 screens)

**Files:**
- Create: `mockups/journey-01-firm-setup/01-signup-form.html` through `09-onboarding-complete.html`

Read: Journey 1 in `docs/user-journeys/all-journeys.md` (lines 40–241)

Use `/impeccable craft` for each screen. Provide the screen description from the Screen Inventory table as the craft prompt, along with the user journey stage details.

- [ ] **Step 1:** `/impeccable craft` `01-signup-form.html` — marketing-style signup page, no app shell, centered layout
- [ ] **Step 2:** `/impeccable craft` `02-email-verification.html` — "Check your inbox" interstitial
- [ ] **Step 3:** `/impeccable craft` `03-intake-form.html` — firm intake form
- [ ] **Step 4:** `/impeccable craft` `04-firm-profile.html` — **first screen with app shell** (sidebar + header + content), onboarding step 2/5. This is the first app shell screen — after generating it, run `/impeccable extract` to identify reusable shell components (nav, header, cards, buttons, badges, tables) that all subsequent app shell screens should reuse.
- [ ] **Step 5:** `/impeccable craft` `05-methodology-templates.html` — template selection with previews. Reuse the app shell from step 4.
- [ ] **Step 6:** `/impeccable craft` `06-create-engagement.html` — engagement creation wizard
- [ ] **Step 7:** `/impeccable craft` `07-engagement-ready.html` — confirmation with scaffold summary
- [ ] **Step 8:** `/impeccable craft` `08-invite-staff.html` — bulk invite with role assignment
- [ ] **Step 9:** `/impeccable craft` `09-onboarding-complete.html` — completion state with next steps
- [ ] **Step 10:** Run `/audit` across all 9 files. Fix any issues found.
- [ ] **Step 11:** Run `/polish` for final pass on all 9 files.
- [ ] **Step 12:** Open each file in browser and verify layout, data, and visual consistency.
- [ ] **Step 13:** Commit.

```bash
git add mockups/journey-01-firm-setup/
git commit -m "docs: add Journey 1 (Firm Setup) HTML mockups — 9 screens"
```

---

### Task 3: Journey 2 — Staff Onboarding (5 screens)

**Files:**
- Create: `mockups/journey-02-staff-onboarding/01-invitation-email.html` through `05-first-assignment.html`

Read: Journey 2 in `docs/user-journeys/all-journeys.md` (lines 243–401)

Use `/impeccable craft` for each screen. Reuse extracted app shell components from Task 2.

- [ ] **Step 1:** `/impeccable craft` `01-invitation-email.html` — email template mockup (standalone, no app shell)
- [ ] **Step 2:** `/impeccable craft` `02-magic-link-landing.html` — SSO/password setup (standalone)
- [ ] **Step 3:** `/impeccable craft` `03-profile-setup.html` — profile and notification preferences
- [ ] **Step 4:** `/impeccable craft` `04-guided-tour-step.html` — tour overlay on app shell
- [ ] **Step 5:** `/impeccable craft` `05-first-assignment.html` — engagement list with first assignment notification
- [ ] **Step 6:** Run `/audit` and `/polish` across all 5 files.
- [ ] **Step 7:** Open each file in browser and verify.
- [ ] **Step 8:** Commit.

```bash
git add mockups/journey-02-staff-onboarding/
git commit -m "docs: add Journey 2 (Staff Onboarding) HTML mockups — 5 screens"
```

---

### Task 4: Journey 3 — Engagement Scoping (7 screens)

**Files:**
- Create: `mockups/journey-03-engagement-scoping/01-new-engagement-type.html` through `07-begin-fieldwork.html`

Read: Journey 3 in `docs/user-journeys/all-journeys.md` (lines 405–649)

Use `/impeccable craft` for each screen. Reuse extracted app shell components.

- [ ] **Step 1:** `/impeccable craft` `01-new-engagement-type.html` — type and framework selection
- [ ] **Step 2:** `/impeccable craft` `02-engagement-details.html` — client, period, rollforward detection
- [ ] **Step 3:** `/impeccable craft` `03-team-assignment.html` — team assignment panel
- [ ] **Step 4:** `/impeccable craft` `04-ai-control-mapping.html` — cross-framework mapping table (the "aha moment" — make this screen impressive, show real SOC 2 + ISO 27001 cross-mappings with confidence scores)
- [ ] **Step 5:** `/impeccable craft` `05-client-acceptance.html` — SQMS 1 acceptance form
- [ ] **Step 6:** `/impeccable craft` `06-eqr-assignment.html` — EQR reviewer assignment
- [ ] **Step 7:** `/impeccable craft` `07-begin-fieldwork.html` — readiness checklist and transition
- [ ] **Step 8:** Run `/audit` and `/polish` across all 7 files.
- [ ] **Step 9:** Open each file in browser and verify.
- [ ] **Step 10:** Commit.

```bash
git add mockups/journey-03-engagement-scoping/
git commit -m "docs: add Journey 3 (Engagement Scoping) HTML mockups — 7 screens"
```

---

### Task 5: Journey 4 — Trial Balance (7 screens)

**Files:**
- Create: `mockups/journey-04-trial-balance/01-tb-import.html` through `07-sampling.html`

Read: Journey 4 in `docs/user-journeys/all-journeys.md` (lines 652–900)

Use `/impeccable craft` for each screen. Reuse extracted app shell components.

- [ ] **Step 1:** `/impeccable craft` `01-tb-import.html` — upload zone
- [ ] **Step 2:** `/impeccable craft` `02-column-mapping.html` — column detection and mapping
- [ ] **Step 3:** `/impeccable craft` `03-ai-account-mapping.html` — sheets-style UI with AI mappings (this should look like a spreadsheet — use a grid layout with many rows/columns)
- [ ] **Step 4:** `/impeccable craft` `04-lead-schedules.html` — lead schedule workpaper with drill-through
- [ ] **Step 5:** `/impeccable craft` `05-adjustments.html` — dual-column TB with adjustment tracking
- [ ] **Step 6:** `/impeccable craft` `06-analytics.html` — analytics dashboard with variance and ratios
- [ ] **Step 7:** `/impeccable craft` `07-sampling.html` — sampling calculator and population analytics toggle
- [ ] **Step 8:** Run `/audit` and `/polish` across all 7 files.
- [ ] **Step 9:** Open each file in browser and verify.
- [ ] **Step 10:** Commit.

```bash
git add mockups/journey-04-trial-balance/
git commit -m "docs: add Journey 4 (Trial Balance) HTML mockups — 7 screens"
```

---

### Task 6: Journey 5 — Control Testing (7 screens)

**Files:**
- Create: `mockups/journey-05-control-testing/01-my-assignments.html` through `07-submit-for-review.html`

Read: Journey 5 in `docs/user-journeys/all-journeys.md` (lines 903–1151)

Use `/impeccable craft` for each screen. Reuse extracted app shell components.

- [ ] **Step 1:** `/impeccable craft` `01-my-assignments.html` — personalized control assignments
- [ ] **Step 2:** `/impeccable craft` `02-control-detail.html` — control detail with cross-framework display
- [ ] **Step 3:** `/impeccable craft` `03-test-procedure.html` — test execution with structured fields
- [ ] **Step 4:** `/impeccable craft` `04-evidence-linking.html` — split-view with AI suggestions (this is a key differentiator screen — show the split view clearly with evidence on one side, test procedure on the other)
- [ ] **Step 5:** `/impeccable craft` `05-ai-workpaper-draft.html` — AI draft with prominent banner
- [ ] **Step 6:** `/impeccable craft` `06-workpaper-editor.html` — full editor with evidence sidebar
- [ ] **Step 7:** `/impeccable craft` `07-submit-for-review.html` — submission dialog with validation
- [ ] **Step 8:** Run `/audit` and `/polish` across all 7 files.
- [ ] **Step 9:** Open each file in browser and verify.
- [ ] **Step 10:** Commit.

```bash
git add mockups/journey-05-control-testing/
git commit -m "docs: add Journey 5 (Control Testing) HTML mockups — 7 screens"
```

---

### Task 7: Journey 6 — Workpaper Review (5 screens)

**Files:**
- Create: `mockups/journey-06-workpaper-review/01-review-queue.html` through `05-engagement-progress.html`

Read: Journey 6 in `docs/user-journeys/all-journeys.md` (lines 1154–1364)

Use `/impeccable craft` for each screen. Reuse extracted app shell components.

- [ ] **Step 1:** `/impeccable craft` `01-review-queue.html` — cross-engagement review dashboard
- [ ] **Step 2:** `/impeccable craft` `02-workpaper-review.html` — split-view review mode
- [ ] **Step 3:** `/impeccable craft` `03-review-notes.html` — inline review notes with comment bubbles
- [ ] **Step 4:** `/impeccable craft` `04-sign-off.html` — sign-off confirmation
- [ ] **Step 5:** `/impeccable craft` `05-engagement-progress.html` — progress dashboard with blocking items
- [ ] **Step 6:** Run `/audit` and `/polish` across all 5 files.
- [ ] **Step 7:** Open each file in browser and verify.
- [ ] **Step 8:** Commit.

```bash
git add mockups/journey-06-workpaper-review/
git commit -m "docs: add Journey 6 (Workpaper Review) HTML mockups — 5 screens"
```

---

### Task 8: Journey 7 — Document Requests (4 screens)

**Files:**
- Create: `mockups/journey-07-document-requests/01-create-requests.html` through `04-evidence-acceptance.html`

Read: Journey 7 in `docs/user-journeys/all-journeys.md` (lines 1366–1577)

Use `/impeccable craft` for each screen. Reuse extracted app shell components.

- [ ] **Step 1:** `/impeccable craft` `01-create-requests.html` — bulk creation from template
- [ ] **Step 2:** `/impeccable craft` `02-request-dashboard.html` — status tracking
- [ ] **Step 3:** `/impeccable craft` `03-ai-review-queue.html` — AI-assessed upload queue
- [ ] **Step 4:** `/impeccable craft` `04-evidence-acceptance.html` — acceptance with send-back option
- [ ] **Step 5:** Run `/audit` and `/polish` across all 4 files.
- [ ] **Step 6:** Open each file in browser and verify.
- [ ] **Step 7:** Commit.

```bash
git add mockups/journey-07-document-requests/
git commit -m "docs: add Journey 7 (Document Requests) HTML mockups — 4 screens"
```

---

### Task 9: Journey 8 — Client Hub (6 screens)

**Files:**
- Create: `mockups/journey-08-client-hub/01-client-landing.html` through `06-completion-view.html`

Read: Journey 8 in `docs/user-journeys/all-journeys.md` (lines 1580–1793)

**IMPORTANT:** Client Hub screens use a DIFFERENT layout — no sidebar, centered narrow width, minimal design. The client is a non-auditor (CFO, controller, IT manager). Use plain language, not audit jargon. Tell `/impeccable craft` explicitly that this is a separate, minimal client-facing portal — not the main app shell.

- [ ] **Step 1:** `/impeccable craft` `01-client-landing.html` — firm-branded landing page (standalone, centered, minimal)
- [ ] **Step 2:** `/impeccable craft` `02-request-list.html` — grouped request list with progress
- [ ] **Step 3:** `/impeccable craft` `03-upload-interface.html` — drag-and-drop upload per request
- [ ] **Step 4:** `/impeccable craft` `04-delegate-request.html` — delegation to colleague
- [ ] **Step 5:** `/impeccable craft` `05-sent-back-feedback.html` — AI gap explanation with re-upload
- [ ] **Step 6:** `/impeccable craft` `06-completion-view.html` — all-done celebration state
- [ ] **Step 7:** Run `/audit` and `/polish` across all 6 files.
- [ ] **Step 8:** Open each file in browser and verify.
- [ ] **Step 9:** Commit.

```bash
git add mockups/journey-08-client-hub/
git commit -m "docs: add Journey 8 (Client Hub) HTML mockups — 6 screens"
```

---

### Task 10: Journey 9 — Reporting & Archive (6 screens)

**Files:**
- Create: `mockups/journey-09-reporting/01-report-generation.html` through `06-archive-confirmation.html`

Read: Journey 9 in `docs/user-journeys/all-journeys.md` (lines 1795–2036)

Use `/impeccable craft` for each screen. Reuse extracted app shell components.

- [ ] **Step 1:** `/impeccable craft` `01-report-generation.html` — report type selector
- [ ] **Step 2:** `/impeccable craft` `02-report-editor.html` — rich report editor with sections
- [ ] **Step 3:** `/impeccable craft` `03-client-draft-review.html` — shared draft with client comments
- [ ] **Step 4:** `/impeccable craft` `04-issue-report.html` — two-step issuance dialog with deadlines
- [ ] **Step 5:** `/impeccable craft` `05-finalize-engagement.html` — finalization checklist
- [ ] **Step 6:** `/impeccable craft` `06-archive-confirmation.html` — archived state with WORM details
- [ ] **Step 7:** Run `/audit` and `/polish` across all 6 files.
- [ ] **Step 8:** Open each file in browser and verify.
- [ ] **Step 9:** Commit.

```bash
git add mockups/journey-09-reporting/
git commit -m "docs: add Journey 9 (Reporting & Archive) HTML mockups — 6 screens"
```

---

### Task 11: Journey 10 — EQR (6 screens)

**Files:**
- Create: `mockups/journey-10-eqr/01-eqr-notification.html` through `06-eqr-signoff.html`

Read: Journey 10 in `docs/user-journeys/all-journeys.md` (lines 2039–2244)

Use `/impeccable craft` for each screen. Reuse extracted app shell components.

**IMPORTANT:** EQR screens use the standard app shell BUT with a prominent read-only indicator ("EQR — Read Only" banner across the top of the content area). No edit buttons visible anywhere. Tell `/impeccable craft` this is a read-only review mode.

- [ ] **Step 1:** `/impeccable craft` `01-eqr-notification.html` — review notification
- [ ] **Step 2:** `/impeccable craft` `02-read-only-engagement.html` — read-only engagement overview with prominent read-only banner
- [ ] **Step 3:** `/impeccable craft` `03-planning-review.html` — planning quality review
- [ ] **Step 4:** `/impeccable craft` `04-testing-review.html` — testing sufficiency with AI edit indicators
- [ ] **Step 5:** `/impeccable craft` `05-findings.html` — structured findings form
- [ ] **Step 6:** `/impeccable craft` `06-eqr-signoff.html` — immutable sign-off confirmation
- [ ] **Step 7:** Run `/audit` and `/polish` across all 6 files.
- [ ] **Step 8:** Open each file in browser and verify.
- [ ] **Step 9:** Commit.

```bash
git add mockups/journey-10-eqr/
git commit -m "docs: add Journey 10 (EQR) HTML mockups — 6 screens"
```

---

### Task 12: Final Review and Index

- [ ] **Step 1:** Verify all 61 HTML files exist and open correctly in a browser.
- [ ] **Step 2:** Run `/audit` across the full `mockups/` directory for a final design consistency check.
- [ ] **Step 3:** Update `mockups/README.md` with a complete screen index linking to each file.
- [ ] **Step 4:** Commit.

```bash
git add mockups/
git commit -m "docs: complete mockup index — 61 screens across 10 journeys"
```

---

## Quality Checklist (Run After Each Journey)

- [ ] All screens use the design language established by `/impeccable teach` (consistent colors, typography, spacing)
- [ ] Sample data is consistent across screens (same firm name, same users, same engagement)
- [ ] Screens open in a browser without unstyled content or broken assets
- [ ] Screens within a journey feel like a connected flow (navigation context is consistent)
- [ ] Client Hub screens (Journey 8) use the centered, no-sidebar layout
- [ ] EQR screens (Journey 10) show the read-only banner
- [ ] AI-related content has a clear visual indicator distinguishing it from human content
- [ ] Cross-framework displays show real framework references (CC6.1, A.8.3, etc.)
- [ ] Data tables have 5–10 realistic rows, not placeholder text
- [ ] Role badges in headers match the journey persona
- [ ] `/audit` passes with no critical issues
- [ ] No impeccable anti-patterns present (gradient text, AI color palette, glassmorphism, etc.)
