# Axiom UI Mockups

High-fidelity HTML mockups for every key screen across the Axiom user journeys.

## Viewing

Open any `.html` file directly in a browser. Each file is self-contained (inline CSS, Google Fonts via CDN). No build step required.

## Structure

```
mockups/
├── journey-01-firm-setup/          # 9 screens — FirmAdmin onboarding flow
├── journey-02-staff-onboarding/    # 5 screens — Staff invitation and setup
├── journey-03-engagement-scoping/  # 7 screens — Partner creates engagement
├── journey-05-control-testing/     # 7 screens — Staff tests controls
├── journey-06-workpaper-review/    # 5 screens — Manager review workflow
├── journey-07-document-requests/   # 4 screens — Document request lifecycle
├── journey-08-client-hub/          # 6 screens — Client-facing portal (no sidebar)
├── journey-09-reporting/           # 6 screens — Report generation and archive
└── journey-10-eqr/                 # 6 screens — Engagement quality review
```

**Total: 55 screens across 9 journeys.** `journey-04-*` was removed during the compliance pivot (see `docs/specs/compliance-pivot-findings.md`); the replacement "Cross-Framework Evidence Mapping" flow (Journey 4 in `docs/user-journeys/all-journeys.md`) is served by `journey-03-engagement-scoping/04-ai-control-mapping.html` plus new screens still to be built — see "Pivot audit" below. Journey 11 (Multi-Framework Integrated Engagement) and Journey 12 (Continuous Assurance) have no mockups yet.

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
| 05 | [05-client-acceptance.html](journey-03-engagement-scoping/05-client-acceptance.html) | Client acceptance checklist |
| 06 | [06-eqr-assignment.html](journey-03-engagement-scoping/06-eqr-assignment.html) | EQR reviewer assignment |
| 07 | [07-begin-fieldwork.html](journey-03-engagement-scoping/07-begin-fieldwork.html) | Fieldwork readiness checklist |

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

## Pivot audit (outstanding mockup work)

Following the compliance pivot (see `docs/specs/compliance-pivot-findings.md`), the following mockup work remains:

### Screens needing content updates

| Screen | Update |
|---|---|
| `journey-01-firm-setup/05-methodology-templates.html` | Replace any GAAS / Financial Audit template options with the compliance roster: SOC 2, ISO 27001, ISO 27701, ISO 42001, HIPAA, PCI DSS, SOC 1. |
| `journey-01-firm-setup/06-create-engagement.html` | Engagement-type selector — remove financial audit options; same roster as above. |
| `journey-03-engagement-scoping/01-new-engagement-type.html` | Same engagement-type roster update. |
| `journey-03-engagement-scoping/04-ai-control-mapping.html` | Extend to show STRM relationship types (equivalent-to / subset-of / intersects-with), cross-framework coverage %, and partial-satisfaction gap lists. |
| `journey-09-reporting/*` | Report-type selector should cover SOC 2 Type 1/2, ISO 27001/27701/42001 certificate support letters, PCI ROC/AOC, HIPAA attestation — not financial-audit opinions. |
| `journey-10-eqr/*` | Generalize "EQR" language for non-SOC engagements; for ISO work the equivalent is firm-level internal QA, not PCAOB-style EQR. |

### New screens required for Journey 11 (Multi-Framework Integrated Engagement)

Create `mockups/journey-11-multi-framework/` with at minimum:
- `01-scope-selector.html` — pick multiple frameworks for one engagement (SOC 2 + ISO 27001 + ISO 27701)
- `02-shared-control-library.html` — CommonControl view with multi-framework satisfies edges
- `03-sampling-window-reconciliation.html` — align Type 2 window with ISO surveillance
- `04-integrated-coverage-dashboard.html` — per-framework coverage % with gap list
- `05-separate-but-linked-reports.html` — issue SOC 2 opinion + ISO readiness letter from the same engagement

### New screens required for Journey 12 (Continuous Assurance)

Create `mockups/journey-12-continuous-assurance/` with at minimum:
- `01-drift-alert.html` — auditee notification that a monitored config drifted
- `02-auto-retest-result.html` — re-test output with diff against prior evidence
- `03-risk-register-update.html` — automated risk register entry with AIDecision trail
- `04-management-response-draft.html` — agent-drafted management response with edit gate
- `05-auditor-material-change-notification.html` — auditor-side view of material drift

### Client Hub extension (Journey 8 → auditee GRC workspace)

Existing `journey-08-client-hub/` screens (6) cover the request/upload flow. Add:
- `07-continuous-monitoring-dashboard.html` — live control health for the auditee
- `08-evidence-freshness.html` — staleness indicator per control, broken down by framework
- `09-policy-library.html` — auditee policy authoring + version history
- `10-connector-status.html` — cloud/identity/dev tool connector health

### Cross-cutting new screens

- `shared/provenance-viewer.html` — cryptographic AIDecision chain for any evidence item or AI output (post-March-2026 trust differentiator)
- `shared/framework-migration-wizard.html` — AI-assisted remapping when framework versions update (e.g., ISO 27001:2013 → :2022)

This work is scoped for a future UI phase and not addressed in the current spec cascade.
