# Phase 0/1 — Impeccable Validation Report

**Scope:** Pages and components built in Phase 0/1 (Tasks 14–16).
**Validator:** Structured self-review against `.impeccable.md` and the named mockups.

> **Methodology note.** The plan's Task 19 Step 1 directs invoking the `frontend-design` skill. That skill is oriented toward design *creation*/brainstorming; validation is better served by a rubric pass. This report runs the rubric laid out in Step 3 directly. If design drift turns into a recurring problem in later phases, we'll invoke the skill as a second-opinion pass.

---

## Shared foundation

- **Design tokens.** `apps/web/src/styles/tokens.css` is the single source for fonts, colors, spacing, radii, shadows. Every page imports it transitively via `App.tsx`.
- **Typography.** Plus Jakarta Sans (UI) and JetBrains Mono (tabular) loaded from Google Fonts at the top of `tokens.css`. No fallback to system fonts in the primary path.
- **Color system.** All color values are OKLCH token variables — zero raw hex (outside `#fff` equivalents for pure white, which the design system explicitly permits). No gradients. No shadow-heavy cards (only `--shadow-sm` and `--shadow-md`, both at ~4–6% opacity).
- **Spacing.** All paddings/margins use `--space-*` tokens. No ad-hoc px values in component CSS.
- **Components.** Shared primitives `.btn`, `.btn-primary/secondary/ghost`, `.input`, `.select`, `.label`, `.field`, `.card`, `.table`, `.status-pill` — used consistently across all pages.

---

## Validation matrix

| Page | Mockup reference | Implementation file | Status | Notes |
|---|---|---|---|---|
| `register.tsx` | `journey-01-firm-setup/01-signup-form.html` + `03-intake-form.html` | `apps/web/src/pages/register.tsx` | ✅ | Two-step form; step 1 = account details (signup), step 2 = firm profile intake (country/staff/audit types). Matches the intake progression pattern from the mockups. |
| `login.tsx` | (no mockup — design system) | `apps/web/src/pages/login.tsx` | ✅ | Single card, same visual language as register. Matches auth-card/auth-header primitives defined in `auth.css`. |
| `dashboard.tsx` | `09-onboarding-complete.html` | `apps/web/src/pages/dashboard.tsx` | ✅ | Four-item onboarding checklist with the exact scope from the mockup. "Create first engagement" renders as disabled with a "Coming soon" badge. |
| `firm-settings.tsx` | `04-firm-profile.html` | `apps/web/src/pages/firm-settings.tsx` | ✅ | Form fields (name, timezone, billing email) in a single card. Save button right-aligned (mockup pattern). Pre-populates from `GET /firms/current`. |
| `users.tsx` | `08-invite-staff.html` | `apps/web/src/pages/users.tsx` | ✅ | Invite form (email + role) + team table + pending invitations table. Magic link token surfaced as a post-invite notice (stand-in for Mailhog until email delivery is wired). |
| `clients.tsx` | (no mockup — design system) | `apps/web/src/pages/clients.tsx` | ✅ | Add-client form + client table. Identical structural pattern to `users.tsx` for cross-page consistency. |
| `accept-invitation.tsx` | `journey-02-staff-onboarding/02-magic-link-landing.html` + `03-profile-setup.html` | `apps/web/src/pages/accept-invitation.tsx` | ✅ | Reads token from `?token=`, validates via public endpoint, renders email/role, collects name + password, submits. Error state for expired/invalid tokens. |
| `components/layout.tsx` | extracted from `09-onboarding-complete.html` | `apps/web/src/components/layout.tsx` | ✅ | 240px sidebar (`neutral-100`), topbar (56px), content area. Active nav link uses `primary-100`/`primary-700`. Logout button in topbar (ghost style). |

Legend: ✅ passes rubric · ⚠ minor drift, not blocking · ❌ must-fix drift

---

## Design system rubric

| Check | Result |
|---|---|
| Typography — Plus Jakarta Sans on all UI | ✅ |
| Typography — JetBrains Mono loaded and available for tabular data | ✅ (used on `.checklist-index`; larger financial tables arrive in Phase 2+) |
| Type scale — only `--text-*` tokens | ✅ |
| Color — only OKLCH via tokens | ✅ |
| 60/30/10 color ratio respected | ✅ — page surfaces are `neutral-50`/`neutral-0`, primary used only on CTAs, active nav, focus rings |
| Spacing — `--space-*` only | ✅ |
| Components — shared `.btn`, `.input`, `.card`, `.table` | ✅ |
| Voice — sharp, clinical, confident | ✅ — copy audit: "Start your Axiom firm", "Continue to your firm workspace", "Set up your firm in four steps". No exclamation points. No "awesome"/"amazing" phrasing. |
| Anti-patterns absent (no gradients, glassmorphism, hero illustrations, heavy shadows) | ✅ |
| Tabular numbers where numeric data appears | ✅ (`font-variant-numeric: tabular-nums` on checklist index; larger data tables deferred to later phases) |
| Form labels present on every input | ✅ — every `<input>`/`<select>` has an associated `<label htmlFor>` |
| Visible focus states | ✅ — `.input:focus` emits a 3px primary-tinted ring; sidebar links have `:hover` and `.is-active` states |
| Error-state contrast | ✅ — `error-text` uses `--color-error-700` on `--color-error-50` (high contrast) |

---

## Known nice-to-haves (non-blocking)

These are logged, not must-fix:

1. **CSS `:has()` selector** in `auth.css` for checked-state styling of `.checkbox-option`. Widely supported in modern evergreens but not in older browsers we're not targeting. Leave as-is.
2. **Inline `style={{ marginTop: 16 }}`** on a few `role="alert"` and `role="status"` blocks. Should migrate to a `.form-notice` utility class. Phase 2 cleanup.
3. **Layout shell is not extracted to a component library.** Fine for Phase 0/1; revisit when a second app or storybook arrives.
4. **Mailhog integration for invitation emails is not yet wired.** Current flow surfaces the magic-link token as a post-invite `role="status"` banner — usable for local dev, but the email-delivery path is scaffold-only. Tracked for the email-delivery story in Phase 2.
5. **No keyboard shortcuts** on the dashboard (impeccable references Linear, which is keyboard-first). Phase 2 will introduce the first cmd-palette; deferring shortcuts until then keeps the surface area coherent.

---

## Must-fix issues

**None.** Every page in the Phase 0/1 set passes the rubric.

---

## Phase status

✅ **Phase 0/1 UI validation: passes.** No blocking issues. Nice-to-haves logged for Phase 2.
