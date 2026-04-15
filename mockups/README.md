# Axiom UI Mockups

High-fidelity HTML mockups for every key screen across all 10 Axiom user journeys.

## Viewing

Open any `.html` file directly in a browser. Each file is self-contained (inline CSS, Google Fonts via CDN). No build step required.

## Structure

```
mockups/
├── journey-01-firm-setup/          # 9 screens — FirmAdmin onboarding flow
├── journey-02-staff-onboarding/    # 5 screens — Staff invitation and setup
├── journey-03-engagement-scoping/  # 7 screens — Partner creates engagement
├── journey-04-trial-balance/       # 7 screens — TB import and analysis
├── journey-05-control-testing/     # 7 screens — Staff tests controls
├── journey-06-workpaper-review/    # 5 screens — Manager review workflow
├── journey-07-document-requests/   # 4 screens — Document request lifecycle
├── journey-08-client-hub/          # 6 screens — Client-facing portal (no sidebar)
├── journey-09-reporting/           # 6 screens — Report generation and archive
└── journey-10-eqr/                 # 6 screens — EQR read-only review
```

**Total: 62 screens across 10 journeys.**

## Design System

Design decisions are documented in `/.impeccable.md` at the project root. Key choices:

- **Typography:** Plus Jakarta Sans (UI) + JetBrains Mono (data)
- **Colors:** OKLCH color space, deep indigo primary, tinted neutrals
- **Spacing:** 4pt base scale
- **Layout:** App shell (sidebar + header) for auditor screens; centered minimal layout for Client Hub; amber read-only banner for EQR

## Sample Data

All screens use consistent, realistic sample data:

| Entity | Value |
|--------|-------|
| Firm | Meridian & Associates CPAs |
| Client | Cloudvault Technologies Inc. |
| Engagement | SOC 2 Type II, Jan 1 – Dec 31, 2025 |
| Team | Sarah Chen (Partner), James Rodriguez (Manager), Emily Park (Staff), David Kim (Staff), Lisa Nguyen (EQR) |

## Screen Index

### Journey 1 — Firm Setup (FirmAdmin: Sarah Chen)

| # | File | Screen |
|---|------|--------|
| 01 | [01-signup-form.html](journey-01-firm-setup/01-signup-form.html) | Sign-up form with SSO and email options |
| 02 | [02-email-verification.html](journey-01-firm-setup/02-email-verification.html) | Email verification interstitial |
| 03 | [03-intake-form.html](journey-01-firm-setup/03-intake-form.html) | Firm intake: name, staff count, audit types |
| 04 | [04-firm-profile.html](journey-01-firm-setup/04-firm-profile.html) | Firm profile in app shell with onboarding checklist |
| 05 | [05-methodology-templates.html](journey-01-firm-setup/05-methodology-templates.html) | Methodology template selection |
| 06 | [06-create-engagement.html](journey-01-firm-setup/06-create-engagement.html) | New engagement form wizard |
| 07 | [07-engagement-ready.html](journey-01-firm-setup/07-engagement-ready.html) | Engagement scaffold confirmation |
| 08 | [08-invite-staff.html](journey-01-firm-setup/08-invite-staff.html) | Bulk staff invitation table |
| 09 | [09-onboarding-complete.html](journey-01-firm-setup/09-onboarding-complete.html) | Onboarding completion and next steps |

### Journey 2 — Staff Onboarding (Staff: Emily Park)

| # | File | Screen |
|---|------|--------|
| 01 | [01-invitation-email.html](journey-02-staff-onboarding/01-invitation-email.html) | Invitation email mockup |
| 02 | [02-magic-link-landing.html](journey-02-staff-onboarding/02-magic-link-landing.html) | Account creation with SSO or password |
| 03 | [03-profile-setup.html](journey-02-staff-onboarding/03-profile-setup.html) | Profile setup: name, role, notifications |
| 04 | [04-guided-tour-step.html](journey-02-staff-onboarding/04-guided-tour-step.html) | Guided product tour tooltip overlay |
| 05 | [05-first-assignment.html](journey-02-staff-onboarding/05-first-assignment.html) | Engagement list with first-assignment toast |

### Journey 3 — Engagement Scoping (Partner: Sarah Chen)

| # | File | Screen |
|---|------|--------|
| 01 | [01-new-engagement-type.html](journey-03-engagement-scoping/01-new-engagement-type.html) | Engagement type selection grid |
| 02 | [02-engagement-details.html](journey-03-engagement-scoping/02-engagement-details.html) | Client details with rollforward option |
| 03 | [03-team-assignment.html](journey-03-engagement-scoping/03-team-assignment.html) | Team role assignment |
| 04 | [04-ai-control-mapping.html](journey-03-engagement-scoping/04-ai-control-mapping.html) | AI control mapping with cross-framework display |
| 05 | [05-client-acceptance.html](journey-03-engagement-scoping/05-client-acceptance.html) | SQMS 1 client acceptance checklist |
| 06 | [06-eqr-assignment.html](journey-03-engagement-scoping/06-eqr-assignment.html) | EQR reviewer assignment |
| 07 | [07-begin-fieldwork.html](journey-03-engagement-scoping/07-begin-fieldwork.html) | Fieldwork readiness checklist |

### Journey 4 — Trial Balance (Staff: Emily Park)

| # | File | Screen |
|---|------|--------|
| 01 | [01-tb-import.html](journey-04-trial-balance/01-tb-import.html) | Trial balance import with drag-and-drop |
| 02 | [02-column-mapping.html](journey-04-trial-balance/02-column-mapping.html) | Column mapping with AI profile detection |
| 03 | [03-ai-account-mapping.html](journey-04-trial-balance/03-ai-account-mapping.html) | AI account-to-FS-line mapping grid |
| 04 | [04-lead-schedules.html](journey-04-trial-balance/04-lead-schedules.html) | Lead schedules with materiality calculator |
| 05 | [05-adjustments.html](journey-04-trial-balance/05-adjustments.html) | Adjustments with unadjusted/adjusted comparison |
| 06 | [06-analytics.html](journey-04-trial-balance/06-analytics.html) | Ratio analytics and AI anomaly flags |
| 07 | [07-sampling.html](journey-04-trial-balance/07-sampling.html) | Statistical sampling configuration |

### Journey 5 — Control Testing (Staff: Emily Park)

| # | File | Screen |
|---|------|--------|
| 01 | [01-my-assignments.html](journey-05-control-testing/01-my-assignments.html) | My assignments dashboard |
| 02 | [02-control-detail.html](journey-05-control-testing/02-control-detail.html) | Control detail with prior year reference |
| 03 | [03-test-procedure.html](journey-05-control-testing/03-test-procedure.html) | Test procedure execution |
| 04 | [04-evidence-linking.html](journey-05-control-testing/04-evidence-linking.html) | Evidence pool browser with AI suggestions |
| 05 | [05-ai-workpaper-draft.html](journey-05-control-testing/05-ai-workpaper-draft.html) | AI-generated workpaper draft |
| 06 | [06-workpaper-editor.html](journey-05-control-testing/06-workpaper-editor.html) | Workpaper editor with version history |
| 07 | [07-submit-for-review.html](journey-05-control-testing/07-submit-for-review.html) | Submit for review modal with validation |

### Journey 6 — Workpaper Review (Manager: James Rodriguez)

| # | File | Screen |
|---|------|--------|
| 01 | [01-review-queue.html](journey-06-workpaper-review/01-review-queue.html) | Review queue with prioritization |
| 02 | [02-workpaper-review.html](journey-06-workpaper-review/02-workpaper-review.html) | Split-view workpaper review with AI indicators |
| 03 | [03-review-notes.html](journey-06-workpaper-review/03-review-notes.html) | Inline review notes and comment threads |
| 04 | [04-sign-off.html](journey-06-workpaper-review/04-sign-off.html) | Sign-off confirmation |
| 05 | [05-engagement-progress.html](journey-06-workpaper-review/05-engagement-progress.html) | Engagement progress dashboard |

### Journey 7 — Document Requests (Manager: James Rodriguez)

| # | File | Screen |
|---|------|--------|
| 01 | [01-create-requests.html](journey-07-document-requests/01-create-requests.html) | Bulk request creation |
| 02 | [02-request-dashboard.html](journey-07-document-requests/02-request-dashboard.html) | Request status dashboard |
| 03 | [03-ai-review-queue.html](journey-07-document-requests/03-ai-review-queue.html) | AI evidence review queue with confidence scores |
| 04 | [04-evidence-acceptance.html](journey-07-document-requests/04-evidence-acceptance.html) | Evidence acceptance with cross-framework display |

### Journey 8 — Client Hub (Client: Michael Reeves)

| # | File | Screen |
|---|------|--------|
| 01 | [01-client-landing.html](journey-08-client-hub/01-client-landing.html) | Client portal landing page |
| 02 | [02-request-list.html](journey-08-client-hub/02-request-list.html) | Document requests grouped by department |
| 03 | [03-upload-interface.html](journey-08-client-hub/03-upload-interface.html) | File upload interface |
| 04 | [04-delegate-request.html](journey-08-client-hub/04-delegate-request.html) | Delegate request to colleague |
| 05 | [05-sent-back-feedback.html](journey-08-client-hub/05-sent-back-feedback.html) | Sent-back notification with AI gap explanation |
| 06 | [06-completion-view.html](journey-08-client-hub/06-completion-view.html) | Portal completion celebration |

### Journey 9 — Reporting & Archive (Partner: Sarah Chen)

| # | File | Screen |
|---|------|--------|
| 01 | [01-report-generation.html](journey-09-reporting/01-report-generation.html) | Report type selector with template preview |
| 02 | [02-report-editor.html](journey-09-reporting/02-report-editor.html) | Rich text report editor with version history |
| 03 | [03-client-draft-review.html](journey-09-reporting/03-client-draft-review.html) | Client draft review with comment threads |
| 04 | [04-issue-report.html](journey-09-reporting/04-issue-report.html) | Report issuance with regulatory deadlines |
| 05 | [05-finalize-engagement.html](journey-09-reporting/05-finalize-engagement.html) | Engagement finalization checklist and lock |
| 06 | [06-archive-confirmation.html](journey-09-reporting/06-archive-confirmation.html) | Archive confirmation with WORM storage details |

### Journey 10 — EQR (EQR Reviewer: Lisa Nguyen)

| # | File | Screen |
|---|------|--------|
| 01 | [01-eqr-notification.html](journey-10-eqr/01-eqr-notification.html) | EQR assignment notification |
| 02 | [02-read-only-engagement.html](journey-10-eqr/02-read-only-engagement.html) | Read-only engagement overview with amber banner |
| 03 | [03-planning-review.html](journey-10-eqr/03-planning-review.html) | Planning review with AI decision trail |
| 04 | [04-testing-review.html](journey-10-eqr/04-testing-review.html) | Testing sufficiency review |
| 05 | [05-findings.html](journey-10-eqr/05-findings.html) | EQR findings form |
| 06 | [06-eqr-signoff.html](journey-10-eqr/06-eqr-signoff.html) | Immutable EQR sign-off |
