# Phase 0/1 — Manual Testing Instructions

This document describes how to run and verify the identity system built in Phase 0/1. Automated tests (Go + Vitest) cover the underlying logic and are re-run on every PR via CI; this walkthrough is for confirming the real browser/HTTP experience end-to-end.

**Scope:** Tasks 1–19 of `docs/superpowers/plans/phase-0-and-1-scaffold-and-identity.md`.

---

## 1. Start the stack

Open three terminals in the repo root (`C:\Users\camil\projects\axiom` on Windows).

**Terminal 1 — database & Mailhog:**

```bash
docker compose up -d
docker compose ps
```

Expected: `axiom-postgres-1` (healthy), `axiom-mailhog-1` (running).

**Terminal 2 — Go API:**

```bash
migrate -path apps/axiom-api/migrations \
  -database "postgres://axiom_svc:localdev@localhost:5432/axiom_db?sslmode=disable" up
cd apps/axiom-api && go run ./cmd/server
```

Expected log lines:
- `"connected to database"`
- `"generating ephemeral JWT keys — tokens will not survive restarts"` (expected in dev)
- `"starting server" port=8080`

**Terminal 3 — React dev server:**

```bash
cd apps/web && npm run dev
```

Vite prints `Local: http://localhost:3000/`.

**Smoke test from a fourth terminal:**

```bash
curl -s http://localhost:3000/api/healthz    # expect: {"status":"ok"}
```

---

## 2. Browser walkthrough — the happy path

Open `http://localhost:3000`. Unauthenticated, you should be redirected to `/login`.

### Register flow (`/register`)

- **Step 1 card:** firm name, your name, work email, password (min 8 chars).
  - Submitting empty → red inline alert.
  - Filling all fields → "Continue" advances to step 2.
- **Step 2 card:** country (US/CA), staff count range, audit-type checkboxes.
  - Toggling a checkbox tints the option card indigo (`primary-50` bg, `primary-500` border). This is the `:has(input:checked)` pattern from `auth.css`.
  - "Create firm" → redirects to `/dashboard`.

### Dashboard (`/dashboard`)

- Left sidebar (240px, `neutral-100` bg) with Dashboard active (indigo pill), Clients / Users / Settings below.
- Topbar shows your display name + "Log out".
- Main area: `Welcome to Axiom, {firm name}` then four numbered onboarding items.
- The fourth item ("Create your first engagement") renders grayed out with a `COMING SOON` pill and no Start button.

### Firm settings (`/settings`)

- Form pre-populates with firm name; timezone and billing email start empty unless provided during registration.
- Change the firm name → "Save changes" → green `Firm profile saved.` banner.
- Next page load: the sidebar brand name updates (because `loadProfile` re-fetches on dashboard mount).

### Users (`/users`)

- Invite form: enter `staff@test.com`, pick `Staff`, click "Send invitation".
- Green banner: **"Invitation sent. Magic link token: {long hex string}"**. Copy that token.
- In the "Pending invitations" section, the row appears with a `Sent` status pill and a Cancel button.

> Mailhog isn't wired to email delivery yet, so the token is surfaced in-app as a stand-in. When email delivery lands in Phase 2, `http://localhost:8025` will catch outgoing invitation mail.

### Invitation acceptance (public route)

- Open `http://localhost:3000/accept-invitation?token={paste-the-token}` in a different browser or incognito window.
- You should see "Join your firm on Axiom" with the invitee email and role populated from `GET /invitations/validate/{token}`.
- Fill display name + password (min 8 chars) → "Accept and sign in" → logged in as Staff, lands on `/dashboard`.
- Error paths to verify:
  - `?token=bad` → "Invitation unavailable" card.
  - Missing `?token=` param → same error card.
  - Reusing a already-accepted token → same error card.

### Clients (`/clients`)

- Add Client form: name is required (empty → red alert). Create "TechCorp" → appears in the table.
- Industry and primary contact email are optional; empty values render as `—` in the table.

### Logout

- Top-right "Log out" → redirects to `/login`.
- Refresh and hit `/dashboard` directly → redirected back to `/login` (protected route).

---

## 3. Targeted edge cases

| Case | Expected |
|---|---|
| Wrong password on login | Red `Invalid email or password.` alert, state unchanged |
| Register with an email that already exists | Red `…may already be in use.` alert on step 2 |
| Staff user (from invitation) visits `/users` | Users table renders; invitations list is silently empty — `GET /invitations` returns 403 and is swallowed |
| Staff tries to submit the invite form | `Could not send invitation.` — backend returns 403 `insufficient permissions` |
| Kill the Go server, then click around the React app | No crash — pages surface "Failed to load…" messages and React keeps working |
| Let the access token expire (15 min) and trigger any authenticated action | The API client auto-refreshes via `/auth/refresh` and retries the original request transparently |
| Manually delete `access_token` from localStorage and refresh | Redirects to `/login` |

---

## 4. RLS sanity check (via psql)

Confirm row-level security is filtering across firms:

```bash
docker compose exec postgres psql -U axiom_svc -d axiom_db \
  -c "SELECT id, name, slug FROM firms;"
```

All firms you've registered appear (no RLS context set at the SQL prompt).

Now set a bogus firm context and re-query clients:

```bash
docker compose exec postgres psql -U axiom_svc -d axiom_db \
  -c "SET app.current_firm_id = '00000000-0000-0000-0000-000000000000'; SELECT * FROM clients;"
```

Expected: zero rows. The RLS policy `firm_id = current_firm_id()` filters out everything the current firm can't see.

For a programmatic isolation check, run the automated test:

```bash
cd apps/axiom-api && go test ./internal/identity/ -run TestRLSIsolation -v
```

---

## 5. Mailhog (placeholder for Phase 2)

`http://localhost:8025` — currently empty because email isn't wired. Once email delivery lands, invitation emails will arrive here.

---

## 6. Known gaps (out of scope for Phase 0/1)

These are expected to be missing; they are not regressions:

- No email delivery (magic-link token is surfaced in the UI instead).
- No keyboard shortcuts on the dashboard.
- No "resend invitation" button (plan mentions it; Phase 2 adds it alongside email).
- `Create first engagement` on the dashboard is a deliberate placeholder — engagements ship in Phase 2.
- JWT keys are ephemeral — restarting the Go server invalidates all existing tokens. Configure `JWT_PRIVATE_KEY` + `JWT_PUBLIC_KEY` env vars for durable tokens (intended once AWS infra lands in Phase 10).

---

## 7. Reporting regressions

If anything in §2, §3, or §4 doesn't behave as described, that's a regression worth fixing before opening the PR. The corresponding automated test should be strengthened to catch the drift.
