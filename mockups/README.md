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

**Total: 61 screens across 10 journeys.**

## Design System

Design decisions are documented in `/.impeccable.md` at the project root. Key choices:

- **Typography:** Plus Jakarta Sans (UI) + JetBrains Mono (data)
- **Colors:** OKLCH color space, deep indigo primary, tinted neutrals
- **Spacing:** 4pt base scale
- **Layout:** App shell (sidebar + header) for auditor screens; centered minimal layout for Client Hub

## Sample Data

All screens use consistent, realistic sample data:

| Entity | Value |
|--------|-------|
| Firm | Meridian & Associates CPAs |
| Client | Cloudvault Technologies Inc. |
| Engagement | SOC 2 Type II, Jan 1 – Dec 31, 2025 |
| Team | Sarah Chen (Partner), James Rodriguez (Manager), Emily Park (Staff), David Kim (Staff), Lisa Nguyen (EQR) |

## Screenshots

To take screenshots programmatically, use Playwright:

```bash
npm install playwright
node screenshot.js
```

Or open each file in a browser and screenshot manually.

## Screen Index

*(See individual journey directories for screen listings.)*
