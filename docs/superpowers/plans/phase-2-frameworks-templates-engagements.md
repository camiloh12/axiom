# Phase 2: Frameworks, Templates & Engagement Creation Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Seed the compliance/assurance framework catalog (SOC 2 TSC, ISO 27001:2022, ISO 27701:2019, ISO 42001:2023, HIPAA Security Rule, PCI DSS v4.0.1, SOC 1 SSAE 18) + platform-level `CommonControl` crosswalk + legacy control-objective library, let firms activate methodology templates, and create fully scoped compliance engagements where selecting a template auto-scaffolds controls and test procedures. Ship the Planning → Fieldwork transition gated by a signed client acceptance. Financial-audit scope (trial balance, GAAS, PCAOB) is explicitly out of scope — see `docs/specs/compliance-pivot-findings.md`.

**Architecture:** Extend the Go modular monolith with two new internal packages (`internal/frameworks` for read-only reference data — frameworks, framework versions, requirements, `common_controls`, `common_control_satisfies`, and the legacy `control_objective_library`; `internal/auditcore` for engagements/controls/test procedures) plus additions to `internal/identity` (methodology templates, firm control objectives). Cross-cutting tables (`ai_decisions`, `audit_log`, `notifications`) are created now but only the audit log is written to in this phase. The `ai_context_type` enum is compliance-pivot-aligned (no `TrialBalanceMapping` / `AnomalyDetection`; adds `EvidenceControlMapping`, `GapAnalysis`, `FrameworkMigration`, `FindingsTriage`, `DriftRetest`, `ManagementResponseDraft`). Frontend adds template browsing, firm control objectives, engagement list, a five-step create-engagement wizard, and tabbed engagement detail.

**Tech Stack:** Go 1.24+, Chi, pgx/v5, sqlc, golang-migrate, PostgreSQL 17 + pgvector | TypeScript, React 19, Vite, TanStack Query, React Router, Zustand | unchanged from Phase 0/1.

---

## UI Methodology: Mockups + Impeccable

**All React work in this plan must reference the existing mockups and the impeccable design system.** Hard rule.

- **Mockups** live in `mockups/journey-*/`. They are the visual + structural source of truth.
- **`.impeccable.md`** at the repo root is the design system. It applies to every screen.
- **The `frontend-design` skill** is used whenever a screen has no corresponding mockup, or when extending a component.

### Per-task rules

For every UI task:

1. **Before writing React code,** read the corresponding mockup. Match its structure, hierarchy, spacing, typography, and design tokens.
2. **If no mockup exists**, invoke the `frontend-design` skill to design from `.impeccable.md` first.
3. **If a new component is needed** that isn't in the mockups, invoke `frontend-design` for that component.

### Mockup → Phase 2 page mapping

| React page / component | Mockup file |
|---|---|
| `pages/templates.tsx` (list) | `mockups/journey-01-firm-setup/05-methodology-templates.html` |
| `pages/templates-detail.tsx` | `mockups/journey-01-firm-setup/05-methodology-templates.html` (detail panel) |
| `pages/engagements.tsx` (list) | **No mockup** — use `frontend-design` skill with `.impeccable.md`, align with the list patterns in `mockups/journey-01-firm-setup/09-onboarding-complete.html` |
| `pages/create-engagement.tsx` (wizard) | `mockups/journey-01-firm-setup/06-create-engagement.html` + `mockups/journey-03-engagement-scoping/01-new-engagement-type.html` + `02-engagement-details.html` + `03-team-assignment.html` |
| `pages/engagement-detail.tsx` (overview tab) | `mockups/journey-01-firm-setup/07-engagement-ready.html` |
| `pages/engagement-detail.tsx` (controls tab) | `mockups/journey-03-engagement-scoping/04-ai-control-mapping.html` (static parts only — AI comes in Phase 7) |
| `pages/engagement-detail.tsx` (client acceptance tab) | `mockups/journey-03-engagement-scoping/05-client-acceptance.html` |
| `pages/engagement-detail.tsx` (advance to fieldwork action) | `mockups/journey-03-engagement-scoping/07-begin-fieldwork.html` |
| `pages/firm-control-objectives.tsx` | **No mockup** — use `frontend-design` skill |

### End-of-phase validation (Task 26)

Every page built in Phase 2 must pass an impeccable validation pass. See Task 26 for the procedure.

---

## Development Methodology: Test-Driven Development

**This plan is executed test-first.** The canonical policy lives in `docs/superpowers/specs/implementation-plan-design.md` under "Cross-Cutting Methodology: Test-Driven Development" — read that section before starting. Phase 0/1 already codified the loop and file-pattern test expectations; the same rules apply here verbatim.

### Reminders specific to Phase 2

- **Service tests use `platform.TestDB(t)`** — it creates a temporary Postgres database, runs every migration under `apps/axiom-api/migrations/`, and drops the database on cleanup. All migrations added in this phase must be compatible with this harness.
- **Seed SQL must be idempotent** — Phase 2 seed scripts ship inside migrations, so they run under `TestDB(t)` too. Use `ON CONFLICT DO NOTHING` or equivalent.
- **RLS tests are non-negotiable** — every new tenant-scoped table ships with a multi-tenant isolation test (see Phase 0/1 `internal/identity/rls_test.go` for the pattern).
- **Trivial code does not need a test** (pure DTO conversions, sqlc-generated code). The test on the consumer covers it.
- **Do not weaken assertions** to make tests pass. If a sqlc-generated type differs from what the test assumed, fix the test to reflect the real behavior — don't remove the assertion.

---

## Git Workflow

Before starting any task, create the phase branch:

- [ ] **Create phase branch**

```bash
git checkout master
git pull origin master
git checkout -b phase-2-frameworks-templates-engagements
```

All commits in Tasks 1–26 go to this branch. Push after each commit:

```bash
git push -u origin phase-2-frameworks-templates-engagements
```

After all tasks are complete, open a PR from `phase-2-frameworks-templates-engagements` → `master`, review, and merge.

---

## File Structure

Files created or modified in this phase. Everything else from Phase 0/1 is unchanged.

```
axiom/
  apps/
    axiom-api/
      migrations/
        000003_audit_core_enums.up.sql             — engagement/control/procedure/ai/audit_log/notification enums
        000003_audit_core_enums.down.sql
        000004_system_reference_tables.up.sql      — frameworks, framework_versions, framework_requirements, common_controls, common_control_satisfies, control_objective_library(+mappings, legacy)
        000004_system_reference_tables.down.sql
        000005_firm_methodology.up.sql             — methodology_templates, template_controls, template_test_procedures, template_document_requests, firm_control_objectives, firm_control_objective_mappings
        000005_firm_methodology.down.sql
        000006_engagements.up.sql                  — engagements, engagement_team_members, engagement_frameworks, client_acceptances, controls, test_procedures
        000006_engagements.down.sql
        000007_cross_cutting_tables.up.sql         — ai_decisions, audit_log (+immutability RULEs), notifications
        000007_cross_cutting_tables.down.sql
        000008_seed_frameworks.up.sql              — SOC 2 TSC 2017, ISO 27001:2022, ISO 27701:2019, ISO 42001:2023, HIPAA Security Rule 2013, PCI DSS v4.0.1, SOC 1 SSAE 18 frameworks + requirements
        000008_seed_frameworks.down.sql
        000009_seed_control_objective_library.up.sql — legacy ControlObjectiveLibrary + cross-framework mappings + platform-seeded CommonControls + CommonControlSatisfies (SCF + AICPA crosswalks)
        000009_seed_control_objective_library.down.sql
        000010_seed_system_templates.up.sql        — six system templates: SOC 2 Type II, ISO 27001:2022, ISO 27701:2019, ISO 42001:2023, HIPAA Security Rule, PCI DSS v4.0.1
        000010_seed_system_templates.down.sql
      sqlc.yaml                                    — ADD: new query directories for frameworks + auditcore
      internal/
        frameworks/
          queries/
            frameworks.sql
            requirements.sql
            library.sql
            (generated .go files land here via sqlc)
          service.go
          service_test.go
          handler.go
          handler_test.go
        identity/
          queries/
            methodology_templates.sql              — NEW
            firm_control_objectives.sql            — NEW
          methodology.go                           — NEW: template + firm-objective service methods
          methodology_test.go                      — NEW
          methodology_handler.go                   — NEW
          methodology_handler_test.go              — NEW
          handler_extras.go                        — MODIFY: register new routes
        auditcore/
          queries/
            engagements.sql
            engagement_team.sql
            engagement_frameworks.sql
            client_acceptances.sql
            controls.sql
            test_procedures.sql
          service.go                               — engagement + control + test procedure CRUD
          service_test.go
          scaffolding.go                           — auto-build controls/test procedures from a template
          scaffolding_test.go
          state_machine.go                         — Planning ↔ Fieldwork transitions with guards
          state_machine_test.go
          client_acceptance.go
          client_acceptance_test.go
          auditlog.go                              — typed helper that writes to audit_log
          auditlog_test.go
          handler.go
          handler_test.go
          rls_test.go
      cmd/server/main.go                           — MODIFY: wire frameworks + auditcore modules
  apps/
    web/
      src/
        api/
          generated/
            audit-core.ts                          — NEW (from openapi-typescript)
        pages/
          templates.tsx                            — NEW
          templates.test.tsx
          templates-detail.tsx                     — NEW
          templates-detail.test.tsx
          firm-control-objectives.tsx              — NEW
          firm-control-objectives.test.tsx
          engagements.tsx                          — NEW
          engagements.test.tsx
          create-engagement.tsx                    — NEW (wizard)
          create-engagement.test.tsx
          engagement-detail.tsx                    — NEW (tabbed)
          engagement-detail.test.tsx
        components/
          tabs.tsx                                 — NEW shared tab component (used by engagement detail)
          tabs.test.tsx
          wizard.tsx                               — NEW shared wizard shell (used by create-engagement)
          wizard.test.tsx
        App.tsx                                    — MODIFY: new routes
  docs/
    superpowers/
      testing/
        phase-2-frameworks-templates-engagements.md  — NEW manual test doc
```

---

## Task 1: Migration — audit-core enums

**Files:**
- Create: `apps/axiom-api/migrations/000003_audit_core_enums.up.sql`
- Create: `apps/axiom-api/migrations/000003_audit_core_enums.down.sql`

These enums are referenced by migrations 4–7. Install them first.

- [ ] **Step 1: Write the up migration**

Create `apps/axiom-api/migrations/000003_audit_core_enums.up.sql`:

```sql
CREATE TYPE engagement_type AS ENUM (
  'SOC1','SOC2','ISO27001','ISO27701','ISO42001',
  'HIPAA','PCI_DSS',
  'AgreedUponProcedures','Advisory');

CREATE TYPE engagement_status AS ENUM (
  'Planning','Fieldwork','Review','Reporting','Finalized','Archived');

CREATE TYPE control_status AS ENUM (
  'NotStarted','InProgress','Complete','Exception','NotApplicable');

CREATE TYPE procedure_type AS ENUM (
  'Inquiry','Observation','InspectionOfDocument','Reperformance','Analytics');

CREATE TYPE procedure_status AS ENUM (
  'NotStarted','InProgress','Complete','Exception');

CREATE TYPE ai_context_type AS ENUM (
  'ControlMapping','EvidenceLinkSuggestion',
  'EvidenceControlMapping','GapAnalysis','FrameworkMigration',
  'FindingsTriage','DriftRetest','ManagementResponseDraft',
  'RiskCategorySuggestion','WorkpaperDraft',
  'DocumentExtraction','DocumentCompleteness','ReportSectionDraft');

CREATE TYPE strm_relationship_type AS ENUM (
  'equivalent-to','subset-of','superset-of','intersects-with','no-relationship');

CREATE TYPE common_control_source AS ENUM (
  'platform_seed','scf_import','oscal_import','aicpa_mapping','cis_mapping','firm_custom');

CREATE TYPE mapping_source AS ENUM (
  'scf','ucf','oscal','aicpa','cis','axiom_custom');

CREATE TYPE requirement_type AS ENUM ('criterion','control','specification');

CREATE TYPE ai_review_action AS ENUM ('Pending','Accepted','Modified','Rejected');

CREATE TYPE actor_type AS ENUM ('User','System','AIAgent');

CREATE TYPE notification_type AS ENUM (
  'EngagementAssignment','ReviewNoteAdded','ReviewNoteResolved',
  'DocumentRequestStatus','PhaseTransition','EQRNotification',
  'ReminderEscalation','ArchivalConfirmation','RetentionWarning');

CREATE TYPE delivery_channel AS ENUM ('InApp','Email','Both');
```

- [ ] **Step 2: Write the down migration**

Create `apps/axiom-api/migrations/000003_audit_core_enums.down.sql`:

```sql
DROP TYPE IF EXISTS delivery_channel;
DROP TYPE IF EXISTS notification_type;
DROP TYPE IF EXISTS actor_type;
DROP TYPE IF EXISTS ai_review_action;
DROP TYPE IF EXISTS requirement_type;
DROP TYPE IF EXISTS mapping_source;
DROP TYPE IF EXISTS common_control_source;
DROP TYPE IF EXISTS strm_relationship_type;
DROP TYPE IF EXISTS ai_context_type;
DROP TYPE IF EXISTS procedure_status;
DROP TYPE IF EXISTS procedure_type;
DROP TYPE IF EXISTS control_status;
DROP TYPE IF EXISTS engagement_status;
DROP TYPE IF EXISTS engagement_type;
```

- [ ] **Step 3: Apply migrations against the dev database**

```bash
cd apps/axiom-api
migrate -database "postgres://axiom_svc:localdev@localhost:5432/axiom_db?sslmode=disable" -path migrations up
```

Expected: `3/u audit_core_enums (…ms)` in output.

Verify:

```bash
docker compose exec postgres psql -U axiom_svc -d axiom_db -c "\dT+"
```

Expected: the new enum types listed.

- [ ] **Step 4: Commit**

```bash
git add apps/axiom-api/migrations/000003_*
git commit -m "feat(db): add audit-core enums (compliance-pivot engagement types, STRM + mapping source, ai_context_type) for engagement lifecycle, controls, AI, and cross-cutting tables"
git push
```

---

## Task 2: Migration — system reference tables (frameworks, common controls, legacy library)

**Files:**
- Create: `apps/axiom-api/migrations/000004_system_reference_tables.up.sql`
- Create: `apps/axiom-api/migrations/000004_system_reference_tables.down.sql`

These tables are **global reference data**: no `firm_id`, no RLS (with the exception of `common_controls`, which supports firm-custom rows via a partial RLS policy — see below). They are read by every firm.

**Design decision (Option A):** Per `domain-and-data-model-design.md` §3, Phase 2 introduces the new `common_controls` + `common_control_satisfies` cross-framework graph **alongside** the legacy `control_objective_library` + `control_objective_library_mappings` tables. The legacy tables remain in service for firm-level methodology templates (`template_controls.firm_control_objective_id` still resolves via the legacy library). The new tables become the platform-level crosswalk seeded from SCF / OSCAL / AICPA / CIS. A future migration will fold the legacy tables into `common_controls`.

- [ ] **Step 1: Write the up migration**

Create `apps/axiom-api/migrations/000004_system_reference_tables.up.sql`:

```sql
CREATE TABLE frameworks (
  id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  name            text NOT NULL,
  version         text NOT NULL,
  effective_date  date NOT NULL,
  deprecated_at   date,
  governing_body  text NOT NULL,
  description     text,
  created_at      timestamptz NOT NULL DEFAULT now(),
  UNIQUE (name, version)
);

CREATE TABLE framework_requirements (
  id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  framework_id  uuid NOT NULL REFERENCES frameworks(id) ON DELETE CASCADE,
  identifier    text NOT NULL,
  title         text NOT NULL,
  description   text,
  category      text,
  sort_order    integer NOT NULL,
  UNIQUE (framework_id, identifier)
);

CREATE INDEX idx_framework_requirements_fw_sort
  ON framework_requirements(framework_id, sort_order);

-- Legacy cross-framework library. Retained for template compatibility; new mapping
-- work lives in common_controls / common_control_satisfies below.
CREATE TABLE control_objective_library (
  id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  name         text NOT NULL,
  description  text NOT NULL,
  tags         jsonb NOT NULL DEFAULT '[]'
);

CREATE TABLE control_objective_library_mappings (
  id                        uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  library_objective_id      uuid NOT NULL REFERENCES control_objective_library(id) ON DELETE CASCADE,
  framework_requirement_id  uuid NOT NULL REFERENCES framework_requirements(id) ON DELETE CASCADE,
  UNIQUE (library_objective_id, framework_requirement_id)
);

CREATE INDEX idx_col_mappings_library_id
  ON control_objective_library_mappings(library_objective_id);
CREATE INDEX idx_col_mappings_fw_req_id
  ON control_objective_library_mappings(framework_requirement_id);

-- Platform-level common control catalog (new cross-framework graph).
-- NULL firm_id => platform-seeded row, visible to every firm.
-- Non-NULL firm_id => firm-custom adaptation, RLS-isolated to that firm.
CREATE TABLE common_controls (
  id             uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  firm_id        uuid REFERENCES firms(id) ON DELETE CASCADE,
  code           text NOT NULL,
  name           text NOT NULL,
  description    text NOT NULL,
  category       text NOT NULL,
  source         common_control_source NOT NULL,
  created_at     timestamptz NOT NULL DEFAULT now(),
  deprecated_at  timestamptz
);

-- Uniqueness: platform rows share the NULL bucket; firm-custom rows are unique per firm.
CREATE UNIQUE INDEX idx_common_controls_platform_code
  ON common_controls(code) WHERE firm_id IS NULL;
CREATE UNIQUE INDEX idx_common_controls_firm_code
  ON common_controls(firm_id, code) WHERE firm_id IS NOT NULL;
CREATE INDEX idx_common_controls_firm     ON common_controls(firm_id);
CREATE INDEX idx_common_controls_category ON common_controls(category);
CREATE INDEX idx_common_controls_source   ON common_controls(source);

ALTER TABLE common_controls ENABLE ROW LEVEL SECURITY;
-- Partial RLS: platform rows visible to all; firm rows visible only to that firm.
CREATE POLICY common_controls_read ON common_controls
  FOR SELECT
  USING (firm_id IS NULL OR firm_id = current_firm_id());
CREATE POLICY common_controls_write ON common_controls
  FOR ALL
  USING (firm_id = current_firm_id())
  WITH CHECK (firm_id = current_firm_id());

-- STRM-encoded directed edge: CommonControl → FrameworkRequirement.
CREATE TABLE common_control_satisfies (
  id                       uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  common_control_id        uuid NOT NULL REFERENCES common_controls(id) ON DELETE CASCADE,
  framework_requirement_id uuid NOT NULL REFERENCES framework_requirements(id) ON DELETE CASCADE,
  relationship_type        strm_relationship_type NOT NULL,
  strength_score           integer NOT NULL CHECK (strength_score BETWEEN 0 AND 100),
  coverage_notes           text,
  source                   mapping_source NOT NULL,
  valid_from               date NOT NULL,
  valid_to                 date,
  created_at               timestamptz NOT NULL DEFAULT now(),
  -- Partial satisfactions must carry coverage notes so the UI can surface gaps.
  CONSTRAINT ccs_partial_requires_notes CHECK (
    relationship_type NOT IN ('subset-of','intersects-with')
    OR coverage_notes IS NOT NULL
  ),
  UNIQUE (common_control_id, framework_requirement_id, valid_from)
);

CREATE INDEX idx_ccs_common_control ON common_control_satisfies(common_control_id);
CREATE INDEX idx_ccs_fw_req         ON common_control_satisfies(framework_requirement_id);
CREATE INDEX idx_ccs_validity       ON common_control_satisfies(valid_from, valid_to);
```

- [ ] **Step 2: Write the down migration**

Create `apps/axiom-api/migrations/000004_system_reference_tables.down.sql`:

```sql
DROP TABLE IF EXISTS common_control_satisfies;
DROP TABLE IF EXISTS common_controls;
DROP TABLE IF EXISTS control_objective_library_mappings;
DROP TABLE IF EXISTS control_objective_library;
DROP TABLE IF EXISTS framework_requirements;
DROP TABLE IF EXISTS frameworks;
```

- [ ] **Step 3: Apply and verify**

```bash
cd apps/axiom-api
migrate -database "postgres://axiom_svc:localdev@localhost:5432/axiom_db?sslmode=disable" -path migrations up
docker compose exec postgres psql -U axiom_svc -d axiom_db -c "\d frameworks"
```

Expected: table listing matching the up migration.

- [ ] **Step 4: Commit**

```bash
git add apps/axiom-api/migrations/000004_*
git commit -m "feat(db): add frameworks, framework_requirements, common_controls (+ STRM satisfies edges), and legacy control objective library tables"
git push
```

---

## Task 3: Migration — firm methodology tables

**Files:**
- Create: `apps/axiom-api/migrations/000005_firm_methodology.up.sql`
- Create: `apps/axiom-api/migrations/000005_firm_methodology.down.sql`

These are firm-scoped and get RLS. `methodology_templates` allows `is_system_provided=true` with `firm_id` = a sentinel system firm — to keep the data shape uniform, we still require a `firm_id`, and system templates use NULL with a `firm_id IS NULL` carve-out in the RLS policy so every firm sees them. (Domain model permits this — system templates are readable by all firms.)

- [ ] **Step 1: Write the up migration**

Create `apps/axiom-api/migrations/000005_firm_methodology.up.sql`:

```sql
CREATE TABLE methodology_templates (
  id                         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  firm_id                    uuid REFERENCES firms(id) ON DELETE CASCADE,
  name                       text NOT NULL,
  description                text,
  applicable_engagement_type text NOT NULL,
  applicable_framework_id    uuid REFERENCES frameworks(id),
  version                    integer NOT NULL DEFAULT 1,
  is_active                  boolean NOT NULL DEFAULT true,
  is_system_provided         boolean NOT NULL DEFAULT false,
  created_at                 timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT templates_firm_xor_system CHECK (
    (is_system_provided = true  AND firm_id IS NULL) OR
    (is_system_provided = false AND firm_id IS NOT NULL)
  )
);

CREATE INDEX idx_templates_firm_type
  ON methodology_templates(firm_id, applicable_engagement_type);

ALTER TABLE methodology_templates ENABLE ROW LEVEL SECURITY;
CREATE POLICY methodology_templates_isolation ON methodology_templates
  USING (firm_id = current_firm_id() OR is_system_provided = true);

-- Firm activation of system-provided templates (so firms can deactivate them).
CREATE TABLE firm_template_activations (
  id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  firm_id     uuid NOT NULL REFERENCES firms(id) ON DELETE CASCADE,
  template_id uuid NOT NULL REFERENCES methodology_templates(id) ON DELETE CASCADE,
  is_active   boolean NOT NULL DEFAULT true,
  activated_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (firm_id, template_id)
);

CREATE INDEX idx_firm_template_activations_firm
  ON firm_template_activations(firm_id);

ALTER TABLE firm_template_activations ENABLE ROW LEVEL SECURITY;
CREATE POLICY firm_template_activations_isolation ON firm_template_activations
  USING (firm_id = current_firm_id());

CREATE TABLE template_controls (
  id                          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  template_id                 uuid NOT NULL REFERENCES methodology_templates(id) ON DELETE CASCADE,
  firm_control_objective_id   uuid,
  description                 text NOT NULL,
  is_key_control              boolean NOT NULL DEFAULT false,
  sort_order                  integer NOT NULL
);

CREATE INDEX idx_template_controls_template_sort
  ON template_controls(template_id, sort_order);

CREATE TABLE template_test_procedures (
  id                  uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  template_control_id uuid NOT NULL REFERENCES template_controls(id) ON DELETE CASCADE,
  procedure_type      procedure_type NOT NULL,
  description         text NOT NULL,
  expected_result     text,
  sort_order          integer NOT NULL DEFAULT 0
);

CREATE INDEX idx_template_test_procedures_control
  ON template_test_procedures(template_control_id);

CREATE TABLE template_document_requests (
  id                    uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  template_id           uuid NOT NULL REFERENCES methodology_templates(id) ON DELETE CASCADE,
  template_control_id   uuid REFERENCES template_controls(id) ON DELETE SET NULL,
  title                 text NOT NULL,
  instructions_template text NOT NULL,
  sort_order            integer NOT NULL
);

CREATE TABLE firm_control_objectives (
  id                   uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  firm_id              uuid NOT NULL REFERENCES firms(id) ON DELETE CASCADE,
  source_library_id    uuid REFERENCES control_objective_library(id),
  name                 text NOT NULL,
  description          text NOT NULL,
  custom_test_guidance jsonb,
  created_at           timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_firm_control_objectives_firm
  ON firm_control_objectives(firm_id);

ALTER TABLE firm_control_objectives ENABLE ROW LEVEL SECURITY;
CREATE POLICY firm_control_objectives_isolation ON firm_control_objectives
  USING (firm_id = current_firm_id());

CREATE TABLE firm_control_objective_mappings (
  id                         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  firm_control_objective_id  uuid NOT NULL REFERENCES firm_control_objectives(id) ON DELETE CASCADE,
  framework_requirement_id   uuid NOT NULL REFERENCES framework_requirements(id),
  UNIQUE (firm_control_objective_id, framework_requirement_id)
);

CREATE INDEX idx_firm_co_mappings_objective
  ON firm_control_objective_mappings(firm_control_objective_id);
CREATE INDEX idx_firm_co_mappings_fw_req
  ON firm_control_objective_mappings(framework_requirement_id);

-- Add FK now that firm_control_objectives exists.
ALTER TABLE template_controls
  ADD CONSTRAINT template_controls_firm_co_fk
  FOREIGN KEY (firm_control_objective_id)
  REFERENCES firm_control_objectives(id) ON DELETE SET NULL;
```

- [ ] **Step 2: Write the down migration**

Create `apps/axiom-api/migrations/000005_firm_methodology.down.sql`:

```sql
ALTER TABLE template_controls DROP CONSTRAINT IF EXISTS template_controls_firm_co_fk;
DROP TABLE IF EXISTS firm_control_objective_mappings;
DROP TABLE IF EXISTS firm_control_objectives;
DROP TABLE IF EXISTS template_document_requests;
DROP TABLE IF EXISTS template_test_procedures;
DROP TABLE IF EXISTS template_controls;
DROP TABLE IF EXISTS firm_template_activations;
DROP TABLE IF EXISTS methodology_templates;
```

- [ ] **Step 3: Apply and verify**

```bash
migrate -database "postgres://axiom_svc:localdev@localhost:5432/axiom_db?sslmode=disable" -path migrations up
docker compose exec postgres psql -U axiom_svc -d axiom_db -c "\d methodology_templates"
```

- [ ] **Step 4: Commit**

```bash
git add apps/axiom-api/migrations/000005_*
git commit -m "feat(db): add methodology templates and firm control objective tables (firm-scoped, RLS)"
git push
```

---

## Task 4: Migration — engagements, controls, test procedures

**Files:**
- Create: `apps/axiom-api/migrations/000006_engagements.up.sql`
- Create: `apps/axiom-api/migrations/000006_engagements.down.sql`

- [ ] **Step 1: Write the up migration**

Create `apps/axiom-api/migrations/000006_engagements.up.sql`:

```sql
-- Period semantics:
--   - Type 1 (point-in-time) engagements (SOC 1 Type I, SOC 2 Type I): period_start = period_end (single as-of date).
--   - Type 2 (continuous period) engagements (SOC 1 Type II, SOC 2 Type II): period_end > period_start;
--     SOC Type 2 must cover 3 to 12 months; outside this range requires a partner-override marker
--     enforced at the application layer (engagement creation service).
--   - ISO and PCI engagements use the framework's native window.
CREATE TABLE engagements (
  id                        uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  firm_id                   uuid NOT NULL REFERENCES firms(id),
  client_id                 uuid NOT NULL REFERENCES clients(id),
  name                      text NOT NULL,
  engagement_type           engagement_type NOT NULL,
  primary_framework_id      uuid NOT NULL REFERENCES frameworks(id),
  period_start              date NOT NULL,
  period_end                date NOT NULL,
  status                    engagement_status NOT NULL DEFAULT 'Planning',
  prior_engagement_id       uuid REFERENCES engagements(id),
  methodology_template_id   uuid REFERENCES methodology_templates(id),
  report_issued_at          timestamptz,
  assembly_deadline         date,
  retention_deadline        date,
  finalized_at              timestamptz,
  archived_at               timestamptz,
  created_at                timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT engagements_period_valid CHECK (period_end >= period_start)
);

CREATE INDEX idx_engagements_firm ON engagements(firm_id);
CREATE INDEX idx_engagements_client ON engagements(client_id);
CREATE INDEX idx_engagements_status ON engagements(status);
CREATE INDEX idx_engagements_prior ON engagements(prior_engagement_id);

ALTER TABLE engagements ENABLE ROW LEVEL SECURITY;
CREATE POLICY engagements_isolation ON engagements
  USING (firm_id = current_firm_id());

CREATE TABLE engagement_team_members (
  id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  firm_id         uuid NOT NULL REFERENCES firms(id),
  engagement_id   uuid NOT NULL REFERENCES engagements(id) ON DELETE CASCADE,
  user_id         uuid NOT NULL REFERENCES users(id),
  engagement_role text NOT NULL,
  assigned_at     timestamptz NOT NULL DEFAULT now(),
  removed_at      timestamptz
);

CREATE UNIQUE INDEX idx_engagement_team_active
  ON engagement_team_members(engagement_id, user_id)
  WHERE removed_at IS NULL;
CREATE INDEX idx_engagement_team_engagement_user
  ON engagement_team_members(engagement_id, user_id);

ALTER TABLE engagement_team_members ENABLE ROW LEVEL SECURITY;
CREATE POLICY engagement_team_members_isolation ON engagement_team_members
  USING (firm_id = current_firm_id());

CREATE TABLE engagement_frameworks (
  id                uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  firm_id           uuid NOT NULL REFERENCES firms(id),
  engagement_id     uuid NOT NULL REFERENCES engagements(id) ON DELETE CASCADE,
  framework_id      uuid NOT NULL REFERENCES frameworks(id),
  framework_version text NOT NULL,
  is_primary        boolean NOT NULL DEFAULT false,
  UNIQUE (engagement_id, framework_id)
);

ALTER TABLE engagement_frameworks ENABLE ROW LEVEL SECURITY;
CREATE POLICY engagement_frameworks_isolation ON engagement_frameworks
  USING (firm_id = current_firm_id());

CREATE TABLE client_acceptances (
  id                             uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  firm_id                        uuid NOT NULL REFERENCES firms(id),
  engagement_id                  uuid NOT NULL UNIQUE REFERENCES engagements(id) ON DELETE CASCADE,
  quality_risks_identified       jsonb NOT NULL DEFAULT '[]',
  firm_responses                 jsonb NOT NULL DEFAULT '[]',
  independence_confirmed         boolean NOT NULL DEFAULT false,
  independence_confirmed_by_id   uuid REFERENCES users(id),
  accepted_by_id                 uuid REFERENCES users(id),
  accepted_at                    timestamptz,
  created_at                     timestamptz NOT NULL DEFAULT now()
);

ALTER TABLE client_acceptances ENABLE ROW LEVEL SECURITY;
CREATE POLICY client_acceptances_isolation ON client_acceptances
  USING (firm_id = current_firm_id());

CREATE TABLE controls (
  id                          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  firm_id                     uuid NOT NULL REFERENCES firms(id),
  engagement_id               uuid NOT NULL REFERENCES engagements(id) ON DELETE CASCADE,
  firm_control_objective_id   uuid REFERENCES firm_control_objectives(id),
  description                 text NOT NULL,
  control_owner_id            uuid REFERENCES users(id),
  auditor_assigned_to_id      uuid REFERENCES users(id),
  status                      control_status NOT NULL DEFAULT 'NotStarted',
  is_key_control              boolean NOT NULL DEFAULT false,
  prior_control_id            uuid REFERENCES controls(id),
  created_at                  timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_controls_firm ON controls(firm_id);
CREATE INDEX idx_controls_engagement ON controls(engagement_id);
CREATE INDEX idx_controls_auditor ON controls(auditor_assigned_to_id);
CREATE INDEX idx_controls_engagement_status ON controls(engagement_id, status);

ALTER TABLE controls ENABLE ROW LEVEL SECURITY;
CREATE POLICY controls_isolation ON controls
  USING (firm_id = current_firm_id());

CREATE TABLE test_procedures (
  id                   uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  firm_id              uuid NOT NULL REFERENCES firms(id),
  control_id           uuid NOT NULL REFERENCES controls(id) ON DELETE CASCADE,
  procedure_type       procedure_type NOT NULL,
  description          text NOT NULL,
  expected_result      text,
  population_size      integer,
  sample_size          integer,
  sampling_method      text,
  result               text,
  exceptions_noted     text,
  conclusion           text,
  performed_by_id      uuid REFERENCES users(id),
  performed_at         timestamptz,
  reviewed_by_id       uuid REFERENCES users(id),
  reviewed_at          timestamptz,
  status               procedure_status NOT NULL DEFAULT 'NotStarted',
  prior_procedure_id   uuid REFERENCES test_procedures(id),
  sort_order           integer NOT NULL DEFAULT 0
);

CREATE INDEX idx_test_procedures_control ON test_procedures(control_id);
CREATE INDEX idx_test_procedures_performed_by ON test_procedures(performed_by_id);

ALTER TABLE test_procedures ENABLE ROW LEVEL SECURITY;
CREATE POLICY test_procedures_isolation ON test_procedures
  USING (firm_id = current_firm_id());
```

- [ ] **Step 2: Write the down migration**

Create `apps/axiom-api/migrations/000006_engagements.down.sql`:

```sql
DROP TABLE IF EXISTS test_procedures;
DROP TABLE IF EXISTS controls;
DROP TABLE IF EXISTS client_acceptances;
DROP TABLE IF EXISTS engagement_frameworks;
DROP TABLE IF EXISTS engagement_team_members;
DROP TABLE IF EXISTS engagements;
```

- [ ] **Step 3: Apply and verify**

```bash
migrate -database "postgres://axiom_svc:localdev@localhost:5432/axiom_db?sslmode=disable" -path migrations up
docker compose exec postgres psql -U axiom_svc -d axiom_db -c "\d engagements"
```

- [ ] **Step 4: Commit**

```bash
git add apps/axiom-api/migrations/000006_*
git commit -m "feat(db): add engagement aggregate, controls, and test procedures (firm-scoped, RLS)"
git push
```

---

## Task 5: Migration — cross-cutting tables (ai_decisions, audit_log, notifications)

**Files:**
- Create: `apps/axiom-api/migrations/000007_cross_cutting_tables.up.sql`
- Create: `apps/axiom-api/migrations/000007_cross_cutting_tables.down.sql`

Only `audit_log` is populated this phase. The other two exist so downstream phases don't block on schema churn.

- [ ] **Step 1: Write the up migration**

Create `apps/axiom-api/migrations/000007_cross_cutting_tables.up.sql`:

```sql
CREATE TABLE ai_decisions (
  id                  uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  firm_id             uuid NOT NULL REFERENCES firms(id),
  engagement_id       uuid REFERENCES engagements(id),
  context_type        ai_context_type NOT NULL,
  context_id          uuid NOT NULL,
  context_table       text NOT NULL,
  model_id            text NOT NULL,
  input_token_count   integer,
  output_token_count  integer,
  raw_output          jsonb NOT NULL,
  suggested_value     text,
  confidence          real CHECK (confidence >= 0 AND confidence <= 1),
  explanation         text,
  review_action       ai_review_action NOT NULL DEFAULT 'Pending',
  accepted_value      text,
  reviewed_by_id      uuid REFERENCES users(id),
  reviewed_at         timestamptz,
  created_at          timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_ai_decisions_firm ON ai_decisions(firm_id);
CREATE INDEX idx_ai_decisions_engagement ON ai_decisions(engagement_id);
CREATE INDEX idx_ai_decisions_context ON ai_decisions(context_type, context_id);

ALTER TABLE ai_decisions ENABLE ROW LEVEL SECURITY;
CREATE POLICY ai_decisions_isolation ON ai_decisions
  USING (firm_id = current_firm_id());

CREATE TABLE audit_log (
  id             bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  firm_id        uuid NOT NULL,
  actor_id       uuid,
  actor_type     actor_type NOT NULL,
  action         text NOT NULL,
  resource_type  text NOT NULL,
  resource_id    uuid,
  old_value      jsonb,
  new_value      jsonb,
  ip_address     inet,
  user_agent     text,
  occurred_at    timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_audit_log_firm_time
  ON audit_log(firm_id, occurred_at DESC);
CREATE INDEX idx_audit_log_resource
  ON audit_log(resource_type, resource_id);
CREATE INDEX idx_audit_log_actor
  ON audit_log(actor_id);

-- Immutability: writes only.
CREATE RULE audit_log_no_update AS ON UPDATE TO audit_log DO INSTEAD NOTHING;
CREATE RULE audit_log_no_delete AS ON DELETE TO audit_log DO INSTEAD NOTHING;

ALTER TABLE audit_log ENABLE ROW LEVEL SECURITY;
CREATE POLICY audit_log_isolation ON audit_log
  USING (firm_id = current_firm_id());

CREATE TABLE notifications (
  id                 uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  firm_id            uuid NOT NULL REFERENCES firms(id),
  recipient_id       uuid NOT NULL REFERENCES users(id),
  notification_type  notification_type NOT NULL,
  title              text NOT NULL,
  body               text,
  deep_link          text,
  is_read            boolean NOT NULL DEFAULT false,
  delivery_channel   delivery_channel NOT NULL DEFAULT 'InApp',
  created_at         timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_notifications_recipient
  ON notifications(recipient_id, is_read, created_at DESC);

ALTER TABLE notifications ENABLE ROW LEVEL SECURITY;
CREATE POLICY notifications_isolation ON notifications
  USING (firm_id = current_firm_id());
```

- [ ] **Step 2: Write the down migration**

Create `apps/axiom-api/migrations/000007_cross_cutting_tables.down.sql`:

```sql
DROP TABLE IF EXISTS notifications;
DROP RULE IF EXISTS audit_log_no_delete ON audit_log;
DROP RULE IF EXISTS audit_log_no_update ON audit_log;
DROP TABLE IF EXISTS audit_log;
DROP TABLE IF EXISTS ai_decisions;
```

- [ ] **Step 3: Apply and verify immutability**

```bash
migrate -database "postgres://axiom_svc:localdev@localhost:5432/axiom_db?sslmode=disable" -path migrations up

docker compose exec postgres psql -U axiom_svc -d axiom_db -c "
INSERT INTO audit_log (firm_id, actor_type, action, resource_type)
VALUES (gen_random_uuid(), 'System', 'test', 'test') RETURNING id;
"
```

Record the returned `id`, then try to update and delete it — both must be no-ops:

```bash
docker compose exec postgres psql -U axiom_svc -d axiom_db -c "
UPDATE audit_log SET action = 'hacked' WHERE id = 1;
DELETE FROM audit_log WHERE id = 1;
SELECT action FROM audit_log WHERE id = 1;
"
```

Expected: `action` still equals `test` (RULEs silently no-op). Clean up: `TRUNCATE audit_log RESTART IDENTITY;`

- [ ] **Step 4: Commit**

```bash
git add apps/axiom-api/migrations/000007_*
git commit -m "feat(db): add ai_decisions, immutable audit_log, and notifications tables"
git push
```

---

## Task 6: Seed data — frameworks

**Files:**
- Create: `apps/axiom-api/migrations/000008_seed_frameworks.up.sql`
- Create: `apps/axiom-api/migrations/000008_seed_frameworks.down.sql`

Seed is a migration so it runs under `TestDB(t)`. Writes must be idempotent via `ON CONFLICT DO NOTHING`.

The frameworks seeded here reflect the compliance/assurance pivot (`docs/specs/compliance-pivot-findings.md`). Financial-audit frameworks (GAAS / PCAOB AS) are NOT seeded; financial audits are out of scope.

Frameworks to seed (stable UUID prefix `11111111-0000-0000-0000-...`):
- `001` — **SOC 2 TSC 2017** (AICPA Trust Services Criteria, 2022 revision)
- `002` — **ISO 27001:2022** (Annex A — 93 controls across A.5/A.6/A.7/A.8)
- `003` — **HIPAA Security Rule 2013** (45 CFR §164 Subpart C; note 2024 NPRM tracking)
- `004` — **ISO 27701:2019** (PIMS — Annex A Controllers 31 + Annex B Processors 18)
- `005` — **ISO 42001:2023** (AI management system — Annex A, ~38 controls / 9 objectives)
- `006` — **PCI DSS v4.0.1** (12 high-level requirements; ~300 sub-requirements loaded in a later migration)
- `007` — **SOC 1 SSAE 18** (AICPA SSAE No. 18; control objectives are client-defined — seeded as a generic Common Objectives Library entry rather than a fixed catalog)

- [ ] **Step 1: Write the up migration**

The file is long (seven frameworks × many requirements); show the structure and a slice of each framework. The implementing engineer fills in the full requirement lists from the governing-body source material.

Create `apps/axiom-api/migrations/000008_seed_frameworks.up.sql`:

```sql
-- SOC 2 TSC 2017 (AICPA Trust Services Criteria, 2022 revision)
INSERT INTO frameworks (id, name, version, effective_date, governing_body, description)
VALUES (
  '11111111-0000-0000-0000-000000000001',
  'SOC 2 TSC', '2017',
  '2017-05-01',
  'AICPA',
  'Trust Services Criteria for Security, Availability, Processing Integrity, Confidentiality, and Privacy (2017, revised 2022).'
)
ON CONFLICT (name, version) DO NOTHING;

INSERT INTO framework_requirements (framework_id, identifier, title, description, category, sort_order)
VALUES
  ('11111111-0000-0000-0000-000000000001', 'CC1.1',
   'COSO Principle 1: Commitment to integrity and ethical values',
   'The entity demonstrates a commitment to integrity and ethical values.',
   'CC', 1),
  ('11111111-0000-0000-0000-000000000001', 'CC1.2',
   'COSO Principle 2: Board independence and oversight',
   'The board demonstrates independence from management and exercises oversight.',
   'CC', 2)
  -- ... continue CC1.3–CC1.5, CC2.1–CC2.3, CC3.1–CC3.4, CC4.1–CC4.2, CC5.1–CC5.3,
  --     CC6.1–CC6.8, CC7.1–CC7.5, CC8.1, CC9.1–CC9.2, A1.1–A1.3, PI1.1–PI1.5,
  --     C1.1–C1.2, P1.1–P8.1 (see AICPA 2017 TSC publication for full list).
ON CONFLICT (framework_id, identifier) DO NOTHING;

-- ISO/IEC 27001:2022 (Annex A controls)
INSERT INTO frameworks (id, name, version, effective_date, governing_body, description)
VALUES (
  '11111111-0000-0000-0000-000000000002',
  'ISO 27001', '2022',
  '2022-10-25',
  'ISO/IEC',
  'Information security management — Annex A controls (A.5 Organisational, A.6 People, A.7 Physical, A.8 Technological; 93 controls total).'
)
ON CONFLICT (name, version) DO NOTHING;

INSERT INTO framework_requirements (framework_id, identifier, title, description, category, sort_order)
VALUES
  ('11111111-0000-0000-0000-000000000002', 'A.5.1',
   'Policies for information security',
   'Information security policy and topic-specific policies shall be defined, approved by management.',
   'A.5 Organisational', 1),
  ('11111111-0000-0000-0000-000000000002', 'A.5.2',
   'Information security roles and responsibilities',
   'Information security roles and responsibilities shall be defined and allocated.',
   'A.5 Organisational', 2)
  -- ... continue A.5.3–A.5.37, A.6.1–A.6.8, A.7.1–A.7.14, A.8.1–A.8.34
  --     (see ISO/IEC 27001:2022 Annex A for full list).
ON CONFLICT (framework_id, identifier) DO NOTHING;

-- HIPAA Security Rule (45 CFR §164, Subpart C)
INSERT INTO frameworks (id, name, version, effective_date, governing_body, description)
VALUES (
  '11111111-0000-0000-0000-000000000003',
  'HIPAA Security Rule', '2013',
  '2013-09-23',
  'HHS OCR',
  'HIPAA Security Rule administrative, physical, and technical safeguards (45 CFR §164 Subpart C; 2024 NPRM updates tracked in a follow-up seed).'
)
ON CONFLICT (name, version) DO NOTHING;

INSERT INTO framework_requirements (framework_id, identifier, title, description, category, sort_order)
VALUES
  ('11111111-0000-0000-0000-000000000003', '164.308(a)(1)(i)',
   'Security Management Process',
   'Implement policies and procedures to prevent, detect, contain, and correct security violations.',
   'Administrative Safeguards', 1)
  -- ... continue §164.308(a)(1)–(a)(8), §164.310(a)–(d), §164.312(a)–(e), §164.314, §164.316
ON CONFLICT (framework_id, identifier) DO NOTHING;

-- ISO/IEC 27701:2019 (PIMS — Annex A Controllers + Annex B Processors)
INSERT INTO frameworks (id, name, version, effective_date, governing_body, description)
VALUES (
  '11111111-0000-0000-0000-000000000004',
  'ISO 27701', '2019',
  '2019-08-06',
  'ISO/IEC',
  'Privacy Information Management System extension to ISO 27001 — Annex A Controllers (31 controls) + Annex B Processors (18 controls).'
)
ON CONFLICT (name, version) DO NOTHING;

INSERT INTO framework_requirements (framework_id, identifier, title, description, category, sort_order)
VALUES
  ('11111111-0000-0000-0000-000000000004', 'A.7.2.1',
   'Identify and document purpose',
   'The organization shall identify and document the specific purposes for which the PII will be processed.',
   'Annex A — Controllers', 1),
  ('11111111-0000-0000-0000-000000000004', 'B.8.2.1',
   'Customer agreement',
   'The organization shall ensure that the customer agreement addresses the processing of PII.',
   'Annex B — Processors', 100)
  -- ... continue Annex A (A.7.2–A.7.5 with sub-items: 31 controller controls)
  --     and Annex B (B.8.2–B.8.5: 18 processor controls).
ON CONFLICT (framework_id, identifier) DO NOTHING;

-- ISO/IEC 42001:2023 (AI Management System — Annex A controls / 9 objectives)
INSERT INTO frameworks (id, name, version, effective_date, governing_body, description)
VALUES (
  '11111111-0000-0000-0000-000000000005',
  'ISO 42001', '2023',
  '2023-12-18',
  'ISO/IEC',
  'Artificial Intelligence Management System — Annex A controls spanning AI policy, accountability, resources, impact assessment, AI system lifecycle, data management, information for interested parties, use of AI systems, and third-party relationships (~38 controls / 9 objectives).'
)
ON CONFLICT (name, version) DO NOTHING;

INSERT INTO framework_requirements (framework_id, identifier, title, description, category, sort_order)
VALUES
  ('11111111-0000-0000-0000-000000000005', 'A.2.2',
   'AI policy',
   'The organization shall document a policy for the development or use of AI systems.',
   'A.2 Policies related to AI', 1),
  ('11111111-0000-0000-0000-000000000005', 'A.6.2.1',
   'AI system impact assessment',
   'The organization shall assess the potential consequences for individuals, groups, and society of the AI system.',
   'A.6 AI system lifecycle', 10)
  -- ... continue A.2–A.10: ~38 controls across 9 objectives
  --     (see ISO/IEC 42001:2023 Annex A for the full catalog).
ON CONFLICT (framework_id, identifier) DO NOTHING;

-- PCI DSS v4.0.1 (12 high-level requirements; sub-requirements loaded in a later migration)
INSERT INTO frameworks (id, name, version, effective_date, governing_body, description)
VALUES (
  '11111111-0000-0000-0000-000000000006',
  'PCI DSS', 'v4.0.1',
  '2024-06-01',
  'PCI SSC',
  'Payment Card Industry Data Security Standard. Seeded here with the 12 top-level requirements; the ~300 sub-requirements (e.g., 1.1.1, 8.3.6, 11.4.2) are loaded in a follow-up migration to keep this seed reviewable.'
)
ON CONFLICT (name, version) DO NOTHING;

INSERT INTO framework_requirements (framework_id, identifier, title, description, category, sort_order)
VALUES
  ('11111111-0000-0000-0000-000000000006', '1',
   'Install and maintain network security controls',
   'Network security controls such as firewalls and other network security technologies protect the cardholder data environment.',
   'Build and Maintain a Secure Network and Systems', 1),
  ('11111111-0000-0000-0000-000000000006', '2',
   'Apply secure configurations to all system components',
   'Secure configurations are applied to all system components.',
   'Build and Maintain a Secure Network and Systems', 2)
  -- ... continue requirements 3–12:
  -- 3 Protect stored account data; 4 Protect cardholder data with strong cryptography during transmission;
  -- 5 Protect all systems and networks from malicious software; 6 Develop and maintain secure systems and software;
  -- 7 Restrict access to system components and cardholder data by business need to know;
  -- 8 Identify users and authenticate access; 9 Restrict physical access;
  -- 10 Log and monitor all access; 11 Test security of systems and networks regularly;
  -- 12 Support information security with organizational policies and programs.
ON CONFLICT (framework_id, identifier) DO NOTHING;

-- SOC 1 SSAE 18 (AICPA) — control objectives are client-defined per engagement,
-- so we seed a framework row only. Specific objectives live in the per-engagement
-- Common Objectives Library (see migration 000009).
INSERT INTO frameworks (id, name, version, effective_date, governing_body, description)
VALUES (
  '11111111-0000-0000-0000-000000000007',
  'SOC 1', 'SSAE 18',
  '2017-05-01',
  'AICPA',
  'SOC 1 examination under SSAE No. 18 — reports on controls at a service organization relevant to user entities'' internal control over financial reporting. Control objectives are defined per engagement rather than catalogued.'
)
ON CONFLICT (name, version) DO NOTHING;
-- No framework_requirements rows for SOC 1 — seed a generic Common Objectives
-- Library entry in migration 000009 instead.
```

- [ ] **Step 2: Fill in the requirement lists**

Replace the `-- ... continue ...` comments with actual `INSERT` rows. Use the governing-body documents as the source of truth:
- **SOC 2 TSC**: AICPA 2017 Trust Services Criteria (2022 revision) publication. Expect ~60 rows (CC1–CC9 + A1 + PI1 + C1 + P1–P8).
- **ISO 27001:2022**: ISO/IEC 27001:2022 Annex A — 93 controls across A.5 Organisational (37), A.6 People (8), A.7 Physical (14), A.8 Technological (34).
- **HIPAA Security Rule 2013**: 45 CFR §164 Subpart C — ~55 rows including addressable/required sub-items across Administrative, Physical, Technical, Organizational, and Policies/Procedures safeguards.
- **ISO 27701:2019**: Annex A Controllers (31 controls, A.7.2–A.7.5) + Annex B Processors (18 controls, B.8.2–B.8.5). Expect ~49 rows total.
- **ISO 42001:2023**: Annex A — ~38 controls across 9 objectives (A.2 Policies, A.3 Internal organization, A.4 Resources, A.5 Impact assessment, A.6 Lifecycle, A.7 Data, A.8 Information for interested parties, A.9 Use, A.10 Third-party).
- **PCI DSS v4.0.1**: 12 high-level requirements only (rows 1–12). Sub-requirements (~300 rows — e.g., 1.1.1, 3.4.1, 8.3.6, 11.4.2) are loaded in a follow-up migration so this seed stays reviewable.
- **SOC 1 SSAE 18**: no requirement rows — control objectives are engagement-specific and are seeded as a Common Objectives Library entry in migration 000009.

**Do not leave the `-- ... continue ...` comments in the final migration.**

- [ ] **Step 3: Write the down migration**

Create `apps/axiom-api/migrations/000008_seed_frameworks.down.sql`:

```sql
DELETE FROM framework_requirements
  WHERE framework_id IN (
    '11111111-0000-0000-0000-000000000001', -- SOC 2 TSC 2017
    '11111111-0000-0000-0000-000000000002', -- ISO 27001:2022
    '11111111-0000-0000-0000-000000000003', -- HIPAA Security Rule 2013
    '11111111-0000-0000-0000-000000000004', -- ISO 27701:2019
    '11111111-0000-0000-0000-000000000005', -- ISO 42001:2023
    '11111111-0000-0000-0000-000000000006', -- PCI DSS v4.0.1
    '11111111-0000-0000-0000-000000000007'  -- SOC 1 SSAE 18
  );
DELETE FROM frameworks
  WHERE id IN (
    '11111111-0000-0000-0000-000000000001',
    '11111111-0000-0000-0000-000000000002',
    '11111111-0000-0000-0000-000000000003',
    '11111111-0000-0000-0000-000000000004',
    '11111111-0000-0000-0000-000000000005',
    '11111111-0000-0000-0000-000000000006',
    '11111111-0000-0000-0000-000000000007'
  );
```

- [ ] **Step 4: Apply and verify counts**

```bash
migrate -database "postgres://axiom_svc:localdev@localhost:5432/axiom_db?sslmode=disable" -path migrations up

docker compose exec postgres psql -U axiom_svc -d axiom_db -c "
SELECT f.name, f.version, COUNT(r.id) AS reqs
FROM frameworks f LEFT JOIN framework_requirements r ON r.framework_id = f.id
GROUP BY f.id ORDER BY f.name;
"
```

Expected: seven frameworks rows. All except `SOC 1 SSAE 18` have non-zero requirement counts; SOC 1 has zero (objectives live in the Common Objectives Library entry from migration 000009).

- [ ] **Step 5: Commit**

```bash
git add apps/axiom-api/migrations/000008_*
git commit -m "feat(db): seed SOC 2 TSC, ISO 27001/27701/42001, HIPAA, PCI DSS, and SOC 1 frameworks"
git push
```

---

## Task 7: Seed data — control objective library, CommonControls, and cross-framework mappings

**Files:**
- Create: `apps/axiom-api/migrations/000009_seed_control_objective_library.up.sql`
- Create: `apps/axiom-api/migrations/000009_seed_control_objective_library.down.sql`

This migration seeds two parallel catalogs per the Option A coexistence strategy (see Task 2 note and `domain-and-data-model-design.md` §3):

1. The **legacy `control_objective_library`** — ~25 semantic objectives used by `template_controls.firm_control_objective_id` and firm methodology workflows, with simple `control_objective_library_mappings` to `framework_requirements`.
2. The **platform-level `common_controls` catalog** + STRM-encoded `common_control_satisfies` edges. This is the new cross-framework graph. Seeds are sourced from the recommended licensing stack (`docs/specs/compliance-pivot-findings.md` §1): **SCF** (primary crosswalk, free + CC-licensed + NIST STRM-encoded), **AICPA official mappings** for SOC 2 ↔ ISO 27001 (auditor-defensible), **OSCAL** for NIST-family catalogs (future-proofing), **CIS Controls v8.1** mappings as secondary cross-check. Every edge row records `source` (`scf`, `aicpa`, `oscal`, `cis`, `axiom_custom`), `relationship_type` (NIST STRM: `equivalent-to | subset-of | superset-of | intersects-with | no-relationship`), `strength_score`, and `valid_from` so mapping churn across framework versions is preserved in history.

The crosswalk covers: **SOC 2 TSC ↔ ISO 27001 Annex A ↔ ISO 27701 PIMS ↔ HIPAA safeguards ↔ PCI DSS requirements ↔ ISO 42001 AI controls**. Use AICPA official mappings for SOC 2↔ISO 27001 edges (mark `source = 'aicpa'`); use SCF for everything else (`source = 'scf'`).

Also seed a generic SOC 1 "Common Objectives Library" entry here (SOC 1 objectives are client-defined per engagement, so the framework row from migration 000008 has no `framework_requirements`; an engagement cloning a SOC 1 template instead pulls from this library.)

- [ ] **Step 1: Write the up migration**

Create `apps/axiom-api/migrations/000009_seed_control_objective_library.up.sql` using deterministic UUIDs so mapping FKs are stable:

```sql
INSERT INTO control_objective_library (id, name, description, tags) VALUES
  ('22222222-0000-0000-0000-000000000001',
   'Access provisioning and deprovisioning',
   'Users are granted, modified, and revoked access based on documented authorisation and business need.',
   '["access-management","identity"]'),
  ('22222222-0000-0000-0000-000000000002',
   'Privileged access monitoring',
   'Privileged account usage is logged, reviewed, and constrained by least privilege.',
   '["access-management","monitoring"]'),
  ('22222222-0000-0000-0000-000000000003',
   'Change management approval',
   'Changes to production systems are tested, approved, and documented before deployment.',
   '["change-management"]'),
  -- 22 more: vendor management, incident response, data classification, backups,
  -- MFA enforcement, encryption-at-rest, encryption-in-transit, vulnerability mgmt,
  -- SDLC, logging, log review, anti-malware, network segmentation, physical access,
  -- environmental controls, BCP/DR, risk assessment, employee training,
  -- background checks, policy review, customer data segregation, sub-service
  -- monitoring, monitoring tools, capacity planning.
ON CONFLICT (id) DO NOTHING;

-- Example mappings for "Access provisioning and deprovisioning":
-- Replace the FW-req UUIDs with the actual IDs from Task 6 (use the identifier
-- columns: SOC 2 TSC CC6.1/CC6.2/CC6.3; ISO 27001 A.5.15/A.5.16/A.5.18;
-- ISO 27701 A.7.2.1; HIPAA §164.308(a)(4); PCI DSS 7, 8.)
INSERT INTO control_objective_library_mappings (library_objective_id, framework_requirement_id)
SELECT '22222222-0000-0000-0000-000000000001', fr.id
FROM framework_requirements fr
JOIN frameworks f ON f.id = fr.framework_id
WHERE (f.name = 'SOC 2 TSC'  AND fr.identifier IN ('CC6.1','CC6.2','CC6.3'))
   OR (f.name = 'ISO 27001'  AND fr.identifier IN ('A.5.15','A.5.16','A.5.18'))
   OR (f.name = 'ISO 27701'  AND fr.identifier IN ('A.7.2.1'))
   OR (f.name = 'HIPAA Security Rule' AND fr.identifier IN ('164.308(a)(4)(i)'))
   OR (f.name = 'PCI DSS'    AND fr.identifier IN ('7','8'))
ON CONFLICT (library_objective_id, framework_requirement_id) DO NOTHING;

-- Repeat the INSERT…SELECT pattern for each of the 25 library objectives.

-- ---------------------------------------------------------------------------
-- Platform-level CommonControl catalog + STRM-encoded satisfies edges.
-- (New cross-framework graph, per domain-and-data-model-design.md §3.)
-- Stable UUID prefix: 44444444-0000-0000-0000-... for common_controls.
-- ---------------------------------------------------------------------------

INSERT INTO common_controls (id, firm_id, code, name, description, category, source)
VALUES
  ('44444444-0000-0000-0000-000000000001', NULL,
   'CC-AC-01',
   'Access to production systems is restricted to authorized personnel',
   'Access to production systems is granted only upon documented approval, is periodically reviewed, and is revoked upon role change or departure.',
   'Access Control',
   'platform_seed'),
  ('44444444-0000-0000-0000-000000000002', NULL,
   'CC-CH-01',
   'Changes to production systems are approved, tested, and documented',
   'Change-management process requires documented approval, test evidence, and post-deployment verification before promotion to production.',
   'Change Management',
   'platform_seed'),
  ('44444444-0000-0000-0000-000000000003', NULL,
   'CC-AI-01',
   'AI system impact assessment is performed before deployment',
   'Before any production AI system is deployed, the organization performs and documents an impact assessment covering individuals, groups, and society, per ISO 42001 A.6.',
   'AI Governance',
   'platform_seed')
  -- Additional platform-seeded common controls for: privileged access monitoring,
  -- vendor management, incident response, data classification, backups, MFA,
  -- encryption-at-rest, encryption-in-transit, vulnerability management, SDLC,
  -- logging, log review, anti-malware, network segmentation, physical access,
  -- BCP/DR, risk assessment, employee training, PII purpose limitation (27701),
  -- PII retention (27701), cardholder-data cryptography (PCI DSS 3/4),
  -- quarterly ASV scans (PCI DSS 11.3), AI data governance (42001 A.7), etc.
ON CONFLICT DO NOTHING;

-- STRM-encoded satisfies edges. Example for CC-AC-01 (access control):
-- - AICPA official mapping SOC 2 CC6.1 ↔ ISO 27001 A.5.15 (equivalent-to)
-- - SCF-derived HIPAA §164.308(a)(4) (subset-of — HIPAA only covers workforce access)
-- - SCF-derived PCI DSS 7 + 8 (intersects-with — PCI adds prescriptive password/MFA specificity)
INSERT INTO common_control_satisfies
  (common_control_id, framework_requirement_id,
   relationship_type, strength_score, coverage_notes, source, valid_from)
SELECT '44444444-0000-0000-0000-000000000001', fr.id,
       'equivalent-to', 95, NULL, 'aicpa', DATE '2024-01-01'
FROM framework_requirements fr JOIN frameworks f ON f.id = fr.framework_id
WHERE (f.name = 'SOC 2 TSC' AND fr.identifier = 'CC6.1')
   OR (f.name = 'ISO 27001' AND fr.identifier = 'A.5.15')
ON CONFLICT DO NOTHING;

INSERT INTO common_control_satisfies
  (common_control_id, framework_requirement_id,
   relationship_type, strength_score, coverage_notes, source, valid_from)
SELECT '44444444-0000-0000-0000-000000000001', fr.id,
       'subset-of', 70,
       'HIPAA §164.308(a)(4) covers workforce access authorization but not broader non-workforce production-system access.',
       'scf', DATE '2024-01-01'
FROM framework_requirements fr JOIN frameworks f ON f.id = fr.framework_id
WHERE f.name = 'HIPAA Security Rule' AND fr.identifier = '164.308(a)(4)(i)'
ON CONFLICT DO NOTHING;

INSERT INTO common_control_satisfies
  (common_control_id, framework_requirement_id,
   relationship_type, strength_score, coverage_notes, source, valid_from)
SELECT '44444444-0000-0000-0000-000000000001', fr.id,
       'intersects-with', 60,
       'PCI DSS 7/8 overlap access restriction but add prescriptive authentication and session-management requirements.',
       'scf', DATE '2024-01-01'
FROM framework_requirements fr JOIN frameworks f ON f.id = fr.framework_id
WHERE f.name = 'PCI DSS' AND fr.identifier IN ('7','8')
ON CONFLICT DO NOTHING;

-- Repeat this pattern for every platform-seeded common_control. Keep edges
-- typed correctly: ONLY 'equivalent-to' and 'superset-of' may omit coverage_notes;
-- 'subset-of' and 'intersects-with' MUST include coverage_notes (enforced by CHECK).

-- SOC 1 "Common Objectives Library" entry — SOC 1 objectives are client-defined,
-- so we seed a single generic library objective rather than a fixed catalog.
INSERT INTO control_objective_library (id, name, description, tags)
VALUES
  ('22222222-0000-0000-0000-000000000099',
   'SOC 1 common objectives starter library',
   'Starter library of typical service-organization control objectives for SOC 1 engagements (IT general controls, transaction processing, change management, computer operations). Firms and engagement teams customize objectives per service organization; this row serves as the cloning template.',
   '["soc1","starter-library","it-general-controls"]')
ON CONFLICT (id) DO NOTHING;
-- No control_objective_library_mappings rows for this objective — SOC 1 control
-- objectives do not map to a fixed framework catalog.
```

- [ ] **Step 2: Fill in the remaining 24 library objectives, their mappings, and the CommonControl satisfies edges**

For each of the 25 library objectives:
- Use the same `INSERT…SELECT` pattern to map it to its cross-framework requirements. Every security-related objective must have at least one mapping in SOC 2 TSC, ISO 27001, and at least one of HIPAA / PCI DSS / ISO 27701.
- AI-governance objectives (model cards, AI system inventory, AI impact assessments) map primarily to **ISO 42001** Annex A.
- Privacy-oriented objectives (PII purpose limitation, PII retention, DSR handling) map to **ISO 27701** Annex A/B plus ISO 27001 where applicable.

Also populate the `common_controls` table and `common_control_satisfies` edges:
- Seed **~40–60 platform-level common controls** covering access, change management, vendor management, incident response, backups, cryptography, logging, physical access, BCP/DR, risk management, training, PII controls (27701), AI controls (42001), and cardholder-data controls (PCI DSS).
- For each common control, create STRM-encoded edges to the relevant `framework_requirements` across SOC 2, ISO 27001, ISO 27701, ISO 42001, HIPAA, and PCI DSS. Use `source = 'aicpa'` for SOC 2 ↔ ISO 27001 edges (their official mapping is auditor-defensible); `source = 'scf'` for everything else.
- Use `relationship_type = 'equivalent-to'` only where the mapping is bidirectionally faithful; otherwise pick `subset-of`, `superset-of`, or `intersects-with` and provide `coverage_notes`.
- Every edge carries `valid_from = '2024-01-01'` (the effective date of the current framework versions); leave `valid_to` NULL.

- [ ] **Step 3: Write the down migration**

Create `apps/axiom-api/migrations/000009_seed_control_objective_library.down.sql`:

```sql
DELETE FROM common_control_satisfies;
DELETE FROM common_controls WHERE firm_id IS NULL;
DELETE FROM control_objective_library_mappings;
DELETE FROM control_objective_library;
```

- [ ] **Step 4: Apply and verify**

```bash
migrate -database "postgres://axiom_svc:localdev@localhost:5432/axiom_db?sslmode=disable" -path migrations up

docker compose exec postgres psql -U axiom_svc -d axiom_db -c "
SELECT col.name, COUNT(m.id) AS mappings
FROM control_objective_library col
LEFT JOIN control_objective_library_mappings m ON m.library_objective_id = col.id
GROUP BY col.id ORDER BY col.name;

SELECT cc.code, cc.name, COUNT(ccs.id) AS satisfies_edges
FROM common_controls cc
LEFT JOIN common_control_satisfies ccs ON ccs.common_control_id = cc.id
WHERE cc.firm_id IS NULL
GROUP BY cc.id ORDER BY cc.code;
"
```

Expected: 25 legacy library objectives (each security-related objective ≥3 mappings; the SOC 1 starter has 0); 40–60 platform-seeded common controls, each with ≥3 satisfies edges spanning at least two frameworks.

- [ ] **Step 5: Commit**

```bash
git add apps/axiom-api/migrations/000009_*
git commit -m "feat(db): seed legacy control objective library + platform CommonControls with STRM cross-framework mappings"
git push
```

---

## Task 8: Seed data — system methodology templates

**Files:**
- Create: `apps/axiom-api/migrations/000010_seed_system_templates.up.sql`
- Create: `apps/axiom-api/migrations/000010_seed_system_templates.down.sql`

**Six** system-provided templates ship in MVP (SOC 1 is deferred post-MVP — service-organization-specific and harder to standardize):

| Template name | `applicable_engagement_type` | Framework | Controls | Procedures | Doc requests |
|---|---|---|---|---|---|
| SOC 2 Type II Standard | `SOC2` | SOC 2 TSC 2017 | ~50 | ~80 | ~80 |
| ISO 27001:2022 Standard | `ISO27001` | ISO 27001:2022 | ~93 (Annex A) | ~100 | ~80 |
| ISO 27701:2019 Standard | `ISO27701` | ISO 27701:2019 | ~49 PIMS controls on top of 27001 (seed as additions) | ~60 | ~50 |
| ISO 42001:2023 Standard | `ISO42001` | ISO 42001:2023 | ~38 (AI management) | ~50 | ~40 |
| HIPAA Security Rule Standard | `HIPAA` | HIPAA Security Rule 2013 | ~50 (mapped to Administrative/Physical/Technical safeguards) | ~60 | ~50 |
| PCI DSS v4.0.1 Standard | `PCI_DSS` | PCI DSS v4.0.1 | ~50 starter (the 12 high-level requirements + their key sub-requirements) | ~80 | ~80 |

Stable UUID prefixes (one per template):
- `33333333-0000-0000-0000-000000000001` — SOC 2 Type II Standard
- `33333333-0000-0000-0000-000000000002` — ISO 27001:2022 Standard
- `33333333-0000-0000-0000-000000000003` — ISO 27701:2019 Standard
- `33333333-0000-0000-0000-000000000004` — ISO 42001:2023 Standard
- `33333333-0000-0000-0000-000000000005` — HIPAA Security Rule Standard
- `33333333-0000-0000-0000-000000000006` — PCI DSS v4.0.1 Standard

- [ ] **Step 1: Write the up migration (structure)**

Create `apps/axiom-api/migrations/000010_seed_system_templates.up.sql`:

```sql
-- SOC 2 Type II Standard
INSERT INTO methodology_templates
  (id, firm_id, name, description, applicable_engagement_type, applicable_framework_id, version, is_active, is_system_provided)
VALUES
  ('33333333-0000-0000-0000-000000000001', NULL,
   'SOC 2 Type II Standard',
   'Baseline control, test procedure, and document request set for a SOC 2 Type II engagement under the 2017 TSC.',
   'SOC2',
   '11111111-0000-0000-0000-000000000001',
   1, true, true)
ON CONFLICT (id) DO NOTHING;

-- Controls (~50) — example slice:
INSERT INTO template_controls (id, template_id, description, is_key_control, sort_order) VALUES
  ('33333333-0001-0000-0000-000000000001',
   '33333333-0000-0000-0000-000000000001',
   'The entity maintains a formally documented and board-approved information-security policy, reviewed annually.',
   true, 10),
  ('33333333-0001-0000-0000-000000000002',
   '33333333-0000-0000-0000-000000000001',
   'User access to production systems is granted only upon documented approval from the system owner.',
   true, 20)
  -- ... ~48 more
ON CONFLICT (id) DO NOTHING;

-- Test procedures per control (~80 total across the template):
INSERT INTO template_test_procedures
  (id, template_control_id, procedure_type, description, expected_result, sort_order)
VALUES
  ('33333333-0002-0000-0000-000000000001',
   '33333333-0001-0000-0000-000000000001',
   'InspectionOfDocument',
   'Inspect the signed information-security policy and evidence of annual board approval.',
   'Policy document and meeting minutes showing annual approval within the audit period.',
   10),
  ('33333333-0002-0000-0000-000000000002',
   '33333333-0001-0000-0000-000000000001',
   'Inquiry',
   'Inquire of the security officer about the policy review cadence and distribution.',
   'Corroborates inspected evidence.',
   20)
  -- ... ~78 more
ON CONFLICT (id) DO NOTHING;

-- Document request templates (~80):
INSERT INTO template_document_requests
  (id, template_id, template_control_id, title, instructions_template, sort_order)
VALUES
  ('33333333-0003-0000-0000-000000000001',
   '33333333-0000-0000-0000-000000000001',
   '33333333-0001-0000-0000-000000000001',
   'Information security policy (current version)',
   'Please upload the most recent board-approved information security policy covering the audit period.',
   10)
  -- ... ~79 more
ON CONFLICT (id) DO NOTHING;

-- ISO 27001:2022 Standard
INSERT INTO methodology_templates
  (id, firm_id, name, description, applicable_engagement_type, applicable_framework_id, version, is_active, is_system_provided)
VALUES
  ('33333333-0000-0000-0000-000000000002', NULL,
   'ISO 27001:2022 Standard',
   'Baseline control, procedure, and document request set for an ISO 27001:2022 certification audit. Controls align 1:1 with Annex A (93 controls across A.5/A.6/A.7/A.8).',
   'ISO27001',
   '11111111-0000-0000-0000-000000000002',
   1, true, true)
ON CONFLICT (id) DO NOTHING;
-- ~93 controls, ~100 procedures, ~80 doc requests under id prefixes
-- 33333333-0001-0002-…, 33333333-0002-0002-…, 33333333-0003-0002-…

-- ISO 27701:2019 Standard (PIMS add-on to ISO 27001)
INSERT INTO methodology_templates
  (id, firm_id, name, description, applicable_engagement_type, applicable_framework_id, version, is_active, is_system_provided)
VALUES
  ('33333333-0000-0000-0000-000000000003', NULL,
   'ISO 27701:2019 Standard',
   'Baseline privacy-information-management control set. Adds ~49 PIMS controls (Annex A Controllers + Annex B Processors) on top of ISO 27001:2022. Typically used jointly with the ISO 27001:2022 template in a multi-framework integrated engagement.',
   'ISO27701',
   '11111111-0000-0000-0000-000000000004',
   1, true, true)
ON CONFLICT (id) DO NOTHING;
-- ~49 controls, ~60 procedures, ~50 doc requests under id prefixes
-- 33333333-0001-0003-…, 33333333-0002-0003-…, 33333333-0003-0003-…

-- ISO 42001:2023 Standard
INSERT INTO methodology_templates
  (id, firm_id, name, description, applicable_engagement_type, applicable_framework_id, version, is_active, is_system_provided)
VALUES
  ('33333333-0000-0000-0000-000000000004', NULL,
   'ISO 42001:2023 Standard',
   'Baseline AI-management-system control, procedure, and document request set. Controls align with ISO 42001 Annex A (~38 controls / 9 objectives) covering AI policy, accountability, resources, impact assessment, lifecycle, data management, information for interested parties, use of AI systems, and third-party relationships.',
   'ISO42001',
   '11111111-0000-0000-0000-000000000005',
   1, true, true)
ON CONFLICT (id) DO NOTHING;
-- ~38 controls, ~50 procedures, ~40 doc requests under id prefixes
-- 33333333-0001-0004-…, 33333333-0002-0004-…, 33333333-0003-0004-…

-- HIPAA Security Rule Standard
INSERT INTO methodology_templates
  (id, firm_id, name, description, applicable_engagement_type, applicable_framework_id, version, is_active, is_system_provided)
VALUES
  ('33333333-0000-0000-0000-000000000005', NULL,
   'HIPAA Security Rule Standard',
   'Baseline assessment set for HIPAA Security Rule (45 CFR §164 Subpart C) — controls mapped to Administrative, Physical, and Technical safeguards. A HITRUST r2 assessment mode is out of scope for MVP (see compliance-pivot-findings.md §4.8).',
   'HIPAA',
   '11111111-0000-0000-0000-000000000003',
   1, true, true)
ON CONFLICT (id) DO NOTHING;
-- ~50 controls, ~60 procedures, ~50 doc requests under id prefixes
-- 33333333-0001-0005-…, 33333333-0002-0005-…, 33333333-0003-0005-…

-- PCI DSS v4.0.1 Standard
INSERT INTO methodology_templates
  (id, firm_id, name, description, applicable_engagement_type, applicable_framework_id, version, is_active, is_system_provided)
VALUES
  ('33333333-0000-0000-0000-000000000006', NULL,
   'PCI DSS v4.0.1 Standard',
   'Baseline PCI DSS v4.0.1 control set covering the 12 high-level requirements and their key sub-requirements. MVP seeds ~50 starter controls; the full ~300 sub-requirement catalog is added in a follow-up migration. Quarterly ASV scans and annual pen tests are tracked as population-level procedures.',
   'PCI_DSS',
   '11111111-0000-0000-0000-000000000006',
   1, true, true)
ON CONFLICT (id) DO NOTHING;
-- ~50 controls, ~80 procedures, ~80 doc requests under id prefixes
-- 33333333-0001-0006-…, 33333333-0002-0006-…, 33333333-0003-0006-…

-- (SOC 1 SSAE 18 is NOT seeded in MVP — service-org-specific, post-MVP.)
```

- [ ] **Step 2: Fill in the content**

Flesh out each template's rows. Keep the sort_order gap at 10 so firms can insert custom rows between defaults. All rows must have stable UUIDs (the prefix convention above makes collision avoidance easy).

- [ ] **Step 3: Write the down migration**

Create `apps/axiom-api/migrations/000010_seed_system_templates.down.sql`:

```sql
DELETE FROM template_document_requests
  WHERE template_id IN (
    '33333333-0000-0000-0000-000000000001',
    '33333333-0000-0000-0000-000000000002',
    '33333333-0000-0000-0000-000000000003',
    '33333333-0000-0000-0000-000000000004',
    '33333333-0000-0000-0000-000000000005',
    '33333333-0000-0000-0000-000000000006'
  );
DELETE FROM template_test_procedures
  WHERE template_control_id IN (
    SELECT id FROM template_controls WHERE template_id IN (
      '33333333-0000-0000-0000-000000000001',
      '33333333-0000-0000-0000-000000000002',
      '33333333-0000-0000-0000-000000000003',
      '33333333-0000-0000-0000-000000000004',
      '33333333-0000-0000-0000-000000000005',
      '33333333-0000-0000-0000-000000000006'
    )
  );
DELETE FROM template_controls
  WHERE template_id IN (
    '33333333-0000-0000-0000-000000000001',
    '33333333-0000-0000-0000-000000000002',
    '33333333-0000-0000-0000-000000000003',
    '33333333-0000-0000-0000-000000000004',
    '33333333-0000-0000-0000-000000000005',
    '33333333-0000-0000-0000-000000000006'
  );
DELETE FROM methodology_templates
  WHERE id IN (
    '33333333-0000-0000-0000-000000000001',
    '33333333-0000-0000-0000-000000000002',
    '33333333-0000-0000-0000-000000000003',
    '33333333-0000-0000-0000-000000000004',
    '33333333-0000-0000-0000-000000000005',
    '33333333-0000-0000-0000-000000000006'
  );
```

- [ ] **Step 4: Apply and verify**

```bash
migrate -database "postgres://axiom_svc:localdev@localhost:5432/axiom_db?sslmode=disable" -path migrations up

docker compose exec postgres psql -U axiom_svc -d axiom_db -c "
SELECT mt.name,
  (SELECT COUNT(*) FROM template_controls tc WHERE tc.template_id = mt.id)       AS controls,
  (SELECT COUNT(*) FROM template_test_procedures tp
     JOIN template_controls tc ON tc.id = tp.template_control_id
     WHERE tc.template_id = mt.id)                                               AS procedures,
  (SELECT COUNT(*) FROM template_document_requests dr WHERE dr.template_id = mt.id) AS doc_requests
FROM methodology_templates mt WHERE mt.is_system_provided;
"
```

Expected (approximate row counts):
- `SOC 2 Type II Standard` ≈ 50/80/80
- `ISO 27001:2022 Standard` ≈ 93/100/80
- `ISO 27701:2019 Standard` ≈ 49/60/50
- `ISO 42001:2023 Standard` ≈ 38/50/40
- `HIPAA Security Rule Standard` ≈ 50/60/50
- `PCI DSS v4.0.1 Standard` ≈ 50/80/80

- [ ] **Step 5: Commit**

```bash
git add apps/axiom-api/migrations/000010_*
git commit -m "feat(db): seed six compliance-framework system methodology templates (SOC 2, ISO 27001/27701/42001, HIPAA, PCI DSS)"
git push
```

---

## Task 9: sqlc config update + frameworks queries

**Files:**
- Modify: `apps/axiom-api/sqlc.yaml`
- Create: `apps/axiom-api/internal/frameworks/queries/frameworks.sql`
- Create: `apps/axiom-api/internal/frameworks/queries/requirements.sql`
- Create: `apps/axiom-api/internal/frameworks/queries/library.sql`

- [ ] **Step 1: Update sqlc.yaml to generate three new packages**

Replace the single-entry `sql:` block in `apps/axiom-api/sqlc.yaml` with a list:

```yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "internal/identity/queries"
    schema: "migrations"
    gen:
      go:
        package: "queries"
        out: "internal/identity/queries"
        sql_package: "pgx/v5"
        emit_json_tags: true
        emit_empty_slices: true
        overrides:
          - db_type: "uuid"
            go_type: "github.com/google/uuid.UUID"
          - db_type: "timestamptz"
            go_type: "time.Time"
          - db_type: "jsonb"
            go_type: "encoding/json.RawMessage"

  - engine: "postgresql"
    queries: "internal/frameworks/queries"
    schema: "migrations"
    gen:
      go:
        package: "queries"
        out: "internal/frameworks/queries"
        sql_package: "pgx/v5"
        emit_json_tags: true
        emit_empty_slices: true
        overrides:
          - db_type: "uuid"
            go_type: "github.com/google/uuid.UUID"
          - db_type: "timestamptz"
            go_type: "time.Time"
          - db_type: "jsonb"
            go_type: "encoding/json.RawMessage"

  - engine: "postgresql"
    queries: "internal/auditcore/queries"
    schema: "migrations"
    gen:
      go:
        package: "queries"
        out: "internal/auditcore/queries"
        sql_package: "pgx/v5"
        emit_json_tags: true
        emit_empty_slices: true
        overrides:
          - db_type: "uuid"
            go_type: "github.com/google/uuid.UUID"
          - db_type: "timestamptz"
            go_type: "time.Time"
          - db_type: "jsonb"
            go_type: "encoding/json.RawMessage"
```

- [ ] **Step 2: Write frameworks queries**

Create `apps/axiom-api/internal/frameworks/queries/frameworks.sql`:

```sql
-- name: ListFrameworks :many
SELECT * FROM frameworks
WHERE deprecated_at IS NULL OR sqlc.arg(include_deprecated)::boolean
ORDER BY name, version;

-- name: GetFrameworkByID :one
SELECT * FROM frameworks WHERE id = $1;
```

Create `apps/axiom-api/internal/frameworks/queries/requirements.sql`:

```sql
-- name: ListRequirementsByFramework :many
SELECT * FROM framework_requirements
WHERE framework_id = $1
ORDER BY sort_order;

-- name: ListRequirementsByFrameworkAndCategory :many
SELECT * FROM framework_requirements
WHERE framework_id = $1 AND category = $2
ORDER BY sort_order;
```

Create `apps/axiom-api/internal/frameworks/queries/library.sql`:

```sql
-- name: ListControlObjectiveLibrary :many
SELECT * FROM control_objective_library
ORDER BY name;

-- name: GetControlObjectiveLibraryByID :one
SELECT * FROM control_objective_library WHERE id = $1;

-- name: ListLibraryMappings :many
SELECT m.*, r.identifier AS requirement_identifier, r.title AS requirement_title,
       f.name AS framework_name, f.version AS framework_version
FROM control_objective_library_mappings m
JOIN framework_requirements r ON r.id = m.framework_requirement_id
JOIN frameworks f ON f.id = r.framework_id
WHERE m.library_objective_id = $1;
```

- [ ] **Step 3: Run sqlc**

```bash
cd apps/axiom-api
sqlc generate
```

Expected: three new `queries` directories each containing `db.go`, `models.go`, and `*.sql.go` files. Build succeeds:

```bash
go build ./...
```

- [ ] **Step 4: Commit**

```bash
git add apps/axiom-api/sqlc.yaml apps/axiom-api/internal/frameworks/queries/
git commit -m "feat(api): add sqlc config and queries for frameworks module"
git push
```

---

## Task 10: Frameworks service + handler

**Files:**
- Create: `apps/axiom-api/internal/frameworks/service.go`
- Create: `apps/axiom-api/internal/frameworks/service_test.go`
- Create: `apps/axiom-api/internal/frameworks/handler.go`
- Create: `apps/axiom-api/internal/frameworks/handler_test.go`

Read-only module. No RLS concerns (global tables).

- [ ] **Step 1: Write the failing service test**

Create `apps/axiom-api/internal/frameworks/service_test.go`:

```go
package frameworks_test

import (
	"context"
	"testing"

	"github.com/axiom-platform/axiom/apps/axiom-api/internal/frameworks"
	"github.com/axiom-platform/axiom/apps/axiom-api/internal/platform"
)

func TestListFrameworksReturnsSeededData(t *testing.T) {
	pool := platform.TestDB(t)
	svc := frameworks.NewService(pool)

	got, err := svc.ListFrameworks(context.Background(), false)
	if err != nil {
		t.Fatalf("ListFrameworks: %v", err)
	}
	names := map[string]bool{}
	for _, f := range got {
		names[f.Name+" "+f.Version] = true
	}
	for _, want := range []string{
		"SOC 2 TSC 2017",
		"ISO 27001 2022",
		"ISO 27701 2019",
		"ISO 42001 2023",
		"HIPAA Security Rule 2013",
		"PCI DSS v4.0.1",
		"SOC 1 SSAE 18",
	} {
		if !names[want] {
			t.Errorf("expected framework %q in list, got %v", want, names)
		}
	}
}

func TestListRequirementsByFramework(t *testing.T) {
	pool := platform.TestDB(t)
	svc := frameworks.NewService(pool)
	all, _ := svc.ListFrameworks(context.Background(), false)

	var soc2ID string
	for _, f := range all {
		if f.Name == "SOC 2 TSC" && f.Version == "2017" {
			soc2ID = f.ID.String()
			break
		}
	}
	if soc2ID == "" {
		t.Fatal("SOC 2 TSC not found")
	}

	reqs, err := svc.ListRequirements(context.Background(), soc2ID, "")
	if err != nil {
		t.Fatalf("ListRequirements: %v", err)
	}
	if len(reqs) < 40 {
		t.Errorf("expected ≥40 SOC 2 requirements, got %d", len(reqs))
	}
}
```

- [ ] **Step 2: Run — expect failure**

```bash
cd apps/axiom-api
go test ./internal/frameworks/...
```

Expected: compile error (`frameworks.NewService undefined`).

- [ ] **Step 3: Implement the service**

Create `apps/axiom-api/internal/frameworks/service.go`:

```go
package frameworks

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/axiom-platform/axiom/apps/axiom-api/internal/frameworks/queries"
)

type Service struct {
	q *queries.Queries
}

func NewService(pool *pgxpool.Pool) *Service {
	return &Service{q: queries.New(pool)}
}

type Framework struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	Version        string    `json:"version"`
	EffectiveDate  string    `json:"effective_date"`
	GoverningBody  string    `json:"governing_body"`
	Description    string    `json:"description"`
}

type Requirement struct {
	ID          uuid.UUID `json:"id"`
	FrameworkID uuid.UUID `json:"framework_id"`
	Identifier  string    `json:"identifier"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	SortOrder   int32     `json:"sort_order"`
}

func (s *Service) ListFrameworks(ctx context.Context, includeDeprecated bool) ([]Framework, error) {
	rows, err := s.q.ListFrameworks(ctx, includeDeprecated)
	if err != nil {
		return nil, fmt.Errorf("list frameworks: %w", err)
	}
	out := make([]Framework, 0, len(rows))
	for _, r := range rows {
		desc := ""
		if r.Description.Valid {
			desc = r.Description.String
		}
		out = append(out, Framework{
			ID: r.ID, Name: r.Name, Version: r.Version,
			EffectiveDate: r.EffectiveDate.Format("2006-01-02"),
			GoverningBody: r.GoverningBody, Description: desc,
		})
	}
	return out, nil
}

func (s *Service) ListRequirements(ctx context.Context, frameworkID, category string) ([]Requirement, error) {
	fid, err := uuid.Parse(frameworkID)
	if err != nil {
		return nil, fmt.Errorf("invalid framework id: %w", err)
	}
	var rows []queries.FrameworkRequirement
	if category == "" {
		rows, err = s.q.ListRequirementsByFramework(ctx, fid)
	} else {
		rows, err = s.q.ListRequirementsByFrameworkAndCategory(ctx, queries.ListRequirementsByFrameworkAndCategoryParams{
			FrameworkID: fid, Category: pgText(category),
		})
	}
	if err != nil {
		return nil, fmt.Errorf("list requirements: %w", err)
	}
	out := make([]Requirement, 0, len(rows))
	for _, r := range rows {
		out = append(out, Requirement{
			ID: r.ID, FrameworkID: r.FrameworkID, Identifier: r.Identifier,
			Title: r.Title, Description: nullStr(r.Description),
			Category: nullStr(r.Category), SortOrder: r.SortOrder,
		})
	}
	return out, nil
}
```

Add helpers (`pgText`, `nullStr`) in a small `internal.go` file in the same package (mirroring the identity module's local pg helpers).

- [ ] **Step 4: Run — expect pass**

```bash
go test ./internal/frameworks/...
```

Expected: PASS.

- [ ] **Step 5: Write failing handler test**

Create `apps/axiom-api/internal/frameworks/handler_test.go`:

```go
package frameworks_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/axiom-platform/axiom/apps/axiom-api/internal/frameworks"
	"github.com/axiom-platform/axiom/apps/axiom-api/internal/platform"
)

func TestListFrameworksHandler(t *testing.T) {
	pool := platform.TestDB(t)
	svc := frameworks.NewService(pool)
	h := frameworks.NewHandler(svc)

	r := chi.NewRouter()
	h.RegisterRoutes(r)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/frameworks", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
	}
	var body struct {
		Items []frameworks.Framework `json:"items"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(body.Items) < 4 {
		t.Errorf("want ≥4 frameworks, got %d", len(body.Items))
	}
}
```

- [ ] **Step 6: Run — expect failure**

```bash
go test ./internal/frameworks/...
```

- [ ] **Step 7: Implement the handler**

Create `apps/axiom-api/internal/frameworks/handler.go`:

```go
package frameworks

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/axiom-platform/axiom/apps/axiom-api/internal/platform"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/api/v1/frameworks", h.list)
	r.Get("/api/v1/frameworks/{frameworkId}/requirements", h.listRequirements)
	r.Get("/api/v1/control-objective-library", h.listLibrary)
	r.Get("/api/v1/control-objective-library/{id}/mappings", h.listLibraryMappings)
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	includeDeprecated := r.URL.Query().Get("include_deprecated") == "true"
	items, err := h.svc.ListFrameworks(r.Context(), includeDeprecated)
	if err != nil {
		platform.WriteError(w, err)
		return
	}
	platform.WriteJSON(w, http.StatusOK, map[string]any{"items": items})
}

func (h *Handler) listRequirements(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "frameworkId")
	category := r.URL.Query().Get("category")
	items, err := h.svc.ListRequirements(r.Context(), id, category)
	if err != nil {
		platform.WriteError(w, platform.ErrBadRequest(err.Error()))
		return
	}
	platform.WriteJSON(w, http.StatusOK, map[string]any{"items": items})
}

// listLibrary / listLibraryMappings: same pattern as list/listRequirements — reach
// into svc for ListControlObjectiveLibrary / ListMappings (add those methods in
// service.go alongside ListFrameworks).
```

Add the service methods `ListControlObjectiveLibrary` and `ListMappings` following the pattern in `ListFrameworks`. Add a unit test for each.

- [ ] **Step 8: Run — expect pass**

```bash
go test ./internal/frameworks/...
```

- [ ] **Step 9: Commit**

```bash
git add apps/axiom-api/internal/frameworks/
git commit -m "feat(api): frameworks service + HTTP handlers for read-only reference data"
git push
```

---

## Task 11: Identity — methodology templates queries + service

**Files:**
- Create: `apps/axiom-api/internal/identity/queries/methodology_templates.sql`
- Create: `apps/axiom-api/internal/identity/methodology.go`
- Create: `apps/axiom-api/internal/identity/methodology_test.go`

- [ ] **Step 1: Write queries**

Create `apps/axiom-api/internal/identity/queries/methodology_templates.sql`:

```sql
-- name: ListTemplatesForFirm :many
-- Returns system-provided templates (firm_id IS NULL) and firm's own custom
-- templates. Activation is separate; see firm_template_activations.
SELECT mt.*,
       COALESCE(fta.is_active, mt.is_active) AS effective_is_active
FROM methodology_templates mt
LEFT JOIN firm_template_activations fta
  ON fta.template_id = mt.id AND fta.firm_id = $1
WHERE mt.is_system_provided = true OR mt.firm_id = $1
ORDER BY mt.is_system_provided DESC, mt.name;

-- name: GetTemplateByID :one
SELECT * FROM methodology_templates WHERE id = $1;

-- name: ListTemplateControlsByTemplate :many
SELECT * FROM template_controls WHERE template_id = $1 ORDER BY sort_order;

-- name: ListTemplateProceduresByTemplate :many
SELECT tp.* FROM template_test_procedures tp
JOIN template_controls tc ON tc.id = tp.template_control_id
WHERE tc.template_id = $1
ORDER BY tc.sort_order, tp.sort_order;

-- name: ListTemplateDocRequestsByTemplate :many
SELECT * FROM template_document_requests
WHERE template_id = $1 ORDER BY sort_order;

-- name: UpsertTemplateActivation :one
INSERT INTO firm_template_activations (firm_id, template_id, is_active)
VALUES ($1, $2, $3)
ON CONFLICT (firm_id, template_id)
DO UPDATE SET is_active = EXCLUDED.is_active, activated_at = now()
RETURNING *;
```

- [ ] **Step 2: Run sqlc**

```bash
sqlc generate
go build ./...
```

- [ ] **Step 3: Write failing service test**

Create `apps/axiom-api/internal/identity/methodology_test.go`. Test three behaviors:
1. `ListTemplates` returns the two system templates for a fresh firm.
2. `ActivateTemplate` toggles the `firm_template_activations` row.
3. `GetTemplateDetail` returns controls/procedures/doc-requests in sort order.

Use `platform.TestDB(t)` and create a firm via the existing `identity.NewService`. Follow the pattern in `apps/axiom-api/internal/identity/service_test.go`.

- [ ] **Step 4: Run — expect failure**

```bash
go test ./internal/identity/... -run Methodology
```

- [ ] **Step 5: Implement**

Create `apps/axiom-api/internal/identity/methodology.go`. Add methods on the existing `*Service`:

```go
func (s *Service) ListTemplates(ctx context.Context, firmID uuid.UUID) ([]TemplateDTO, error) { ... }
func (s *Service) GetTemplateDetail(ctx context.Context, templateID uuid.UUID) (*TemplateDetailDTO, error) { ... }
func (s *Service) SetTemplateActivation(ctx context.Context, firmID, templateID uuid.UUID, active bool) error { ... }
```

DTOs:

```go
type TemplateDTO struct {
	ID                       uuid.UUID `json:"id"`
	Name                     string    `json:"name"`
	Description              string    `json:"description"`
	ApplicableEngagementType string    `json:"applicable_engagement_type"`
	ApplicableFrameworkID    *uuid.UUID `json:"applicable_framework_id,omitempty"`
	Version                  int32     `json:"version"`
	IsSystemProvided         bool      `json:"is_system_provided"`
	IsActiveForFirm          bool      `json:"is_active_for_firm"`
}

type TemplateDetailDTO struct {
	TemplateDTO
	Controls         []TemplateControlDTO        `json:"controls"`
	TestProcedures   []TemplateProcedureDTO      `json:"test_procedures"`
	DocumentRequests []TemplateDocRequestDTO     `json:"document_requests"`
}
```

Set `app.current_firm_id` on the pgx transaction/connection before RLS-scoped queries (use the same pattern as `RegisterFirm`).

- [ ] **Step 6: Run — expect pass**

```bash
go test ./internal/identity/... -run Methodology
```

- [ ] **Step 7: Commit**

```bash
git add apps/axiom-api/internal/identity/queries/methodology_templates.sql \
        apps/axiom-api/internal/identity/queries/methodology_templates.sql.go \
        apps/axiom-api/internal/identity/queries/models.go \
        apps/axiom-api/internal/identity/methodology.go \
        apps/axiom-api/internal/identity/methodology_test.go
git commit -m "feat(identity): methodology template list, detail, and firm activation"
git push
```

---

## Task 12: Identity — firm control objectives

**Files:**
- Create: `apps/axiom-api/internal/identity/queries/firm_control_objectives.sql`
- Modify: `apps/axiom-api/internal/identity/methodology.go` (add objective methods)
- Create: `apps/axiom-api/internal/identity/firm_control_objectives_test.go`

- [ ] **Step 1: Write queries**

Create `apps/axiom-api/internal/identity/queries/firm_control_objectives.sql`:

```sql
-- name: ListFirmControlObjectives :many
SELECT * FROM firm_control_objectives
WHERE firm_id = $1
ORDER BY name
LIMIT $2 OFFSET $3;

-- name: CreateFirmControlObjective :one
INSERT INTO firm_control_objectives (firm_id, source_library_id, name, description, custom_test_guidance)
VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: UpdateFirmControlObjective :one
UPDATE firm_control_objectives
SET name = COALESCE(sqlc.narg('name'), name),
    description = COALESCE(sqlc.narg('description'), description),
    custom_test_guidance = COALESCE(sqlc.narg('custom_test_guidance'), custom_test_guidance)
WHERE id = $1
RETURNING *;

-- name: GetFirmControlObjectiveByID :one
SELECT * FROM firm_control_objectives WHERE id = $1;

-- name: ListMappingsForObjective :many
SELECT m.*, r.identifier, r.title, f.name AS framework_name, f.version AS framework_version
FROM firm_control_objective_mappings m
JOIN framework_requirements r ON r.id = m.framework_requirement_id
JOIN frameworks f ON f.id = r.framework_id
WHERE m.firm_control_objective_id = $1;

-- name: AddObjectiveMapping :one
INSERT INTO firm_control_objective_mappings (firm_control_objective_id, framework_requirement_id)
VALUES ($1, $2)
ON CONFLICT (firm_control_objective_id, framework_requirement_id) DO NOTHING
RETURNING *;

-- name: RemoveObjectiveMapping :exec
DELETE FROM firm_control_objective_mappings
WHERE firm_control_objective_id = $1 AND framework_requirement_id = $2;
```

- [ ] **Step 2: sqlc generate; build**

```bash
sqlc generate
go build ./...
```

- [ ] **Step 3: Write failing tests**

Test cases in `firm_control_objectives_test.go`:
- Create objective from library (copies name/description, records `source_library_id`).
- Create custom objective (no library id).
- Update objective (partial fields).
- Add and remove mappings (idempotency — re-add returns same row).
- List objectives for a firm returns only that firm's rows (RLS-aware; call `SET LOCAL app.current_firm_id`).

- [ ] **Step 4: Run — expect failure**

```bash
go test ./internal/identity/... -run ControlObjective
```

- [ ] **Step 5: Implement**

Add methods on `*Service` (in `methodology.go` or a new `firm_control_objectives.go` — keep under 400 lines per file):

```go
func (s *Service) ListFirmControlObjectives(ctx context.Context, firmID uuid.UUID, limit, offset int32) ([]FirmControlObjectiveDTO, error)
func (s *Service) CreateFirmControlObjective(ctx context.Context, firmID uuid.UUID, in CreateFirmControlObjectiveInput) (*FirmControlObjectiveDTO, error)
func (s *Service) UpdateFirmControlObjective(ctx context.Context, id uuid.UUID, in UpdateFirmControlObjectiveInput) (*FirmControlObjectiveDTO, error)
func (s *Service) ListObjectiveMappings(ctx context.Context, objectiveID uuid.UUID) ([]ObjectiveMappingDTO, error)
func (s *Service) AddObjectiveMapping(ctx context.Context, objectiveID, frameworkReqID uuid.UUID) error
func (s *Service) RemoveObjectiveMapping(ctx context.Context, objectiveID, frameworkReqID uuid.UUID) error
```

When `source_library_id` is non-nil on create, copy `name`/`description` from the library row *and* copy each library mapping into `firm_control_objective_mappings` so the firm's objective starts mapped the same way.

- [ ] **Step 6: Run — expect pass**

```bash
go test ./internal/identity/...
```

- [ ] **Step 7: Commit**

```bash
git add apps/axiom-api/internal/identity/
git commit -m "feat(identity): firm control objectives with library seeding and mapping management"
git push
```

---

## Task 13: Identity — methodology + objective HTTP handlers

**Files:**
- Create: `apps/axiom-api/internal/identity/methodology_handler.go`
- Create: `apps/axiom-api/internal/identity/methodology_handler_test.go`
- Modify: `apps/axiom-api/internal/identity/handler_extras.go`

- [ ] **Step 1: Write failing handler tests**

Cases:
- `GET /api/v1/methodology-templates` returns system + firm templates.
- `GET /api/v1/methodology-templates/{id}` returns detail.
- `POST /api/v1/methodology-templates/{id}/activations` toggles activation (FirmAdmin only).
- `GET /api/v1/firm-control-objectives` returns firm's objectives.
- `POST /api/v1/firm-control-objectives` creates (FirmAdmin/Partner only).
- `POST /api/v1/firm-control-objectives/{id}/mappings` adds mapping.
- `DELETE /api/v1/firm-control-objectives/{id}/mappings/{reqId}` removes mapping.

Follow the pattern in `handler_test.go` from Phase 0/1: spin up `platform.TestDB`, create a firm + user, mint a JWT, hit the router through `httptest`.

- [ ] **Step 2: Run — expect failure**

- [ ] **Step 3: Implement the handler**

Create `methodology_handler.go`. Register routes under a new method on `*Handler`:

```go
func (h *Handler) RegisterMethodologyRoutes(r chi.Router, gw RoleGuard) {
	r.Get("/api/v1/methodology-templates", h.listTemplates)
	r.Get("/api/v1/methodology-templates/{id}", h.getTemplate)
	r.With(gw.WithRole("FirmAdmin")).
		Post("/api/v1/methodology-templates/{id}/activations", h.setTemplateActivation)

	r.Get("/api/v1/firm-control-objectives", h.listObjectives)
	r.With(gw.WithRole("FirmAdmin", "Partner")).
		Post("/api/v1/firm-control-objectives", h.createObjective)
	r.Get("/api/v1/firm-control-objectives/{id}", h.getObjective)
	r.With(gw.WithRole("FirmAdmin", "Partner")).
		Patch("/api/v1/firm-control-objectives/{id}", h.updateObjective)
	r.Get("/api/v1/firm-control-objectives/{id}/mappings", h.listObjectiveMappings)
	r.With(gw.WithRole("FirmAdmin", "Partner")).
		Post("/api/v1/firm-control-objectives/{id}/mappings", h.addObjectiveMapping)
	r.With(gw.WithRole("FirmAdmin", "Partner")).
		Delete("/api/v1/firm-control-objectives/{id}/mappings/{reqId}", h.removeObjectiveMapping)
}
```

Handler bodies mirror the existing pattern in `handler_extras.go` (decode → validate → call service → `platform.WriteJSON`).

Call `RegisterMethodologyRoutes(r, gw)` from `RegisterAuthenticatedRoutes` (or directly from `main.go` — same effect).

- [ ] **Step 4: Run — expect pass**

```bash
go test ./internal/identity/...
```

- [ ] **Step 5: Commit**

```bash
git add apps/axiom-api/internal/identity/methodology_handler*.go \
        apps/axiom-api/internal/identity/handler_extras.go
git commit -m "feat(identity): HTTP handlers for methodology templates and firm control objectives"
git push
```

---

## Task 14: Auditcore — queries

**Files:**
- Create: `apps/axiom-api/internal/auditcore/queries/engagements.sql`
- Create: `apps/axiom-api/internal/auditcore/queries/engagement_team.sql`
- Create: `apps/axiom-api/internal/auditcore/queries/engagement_frameworks.sql`
- Create: `apps/axiom-api/internal/auditcore/queries/client_acceptances.sql`
- Create: `apps/axiom-api/internal/auditcore/queries/controls.sql`
- Create: `apps/axiom-api/internal/auditcore/queries/test_procedures.sql`

- [ ] **Step 1: engagements.sql**

```sql
-- name: CreateEngagement :one
INSERT INTO engagements (
  firm_id, client_id, name, engagement_type, primary_framework_id,
  period_start, period_end, methodology_template_id
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;

-- name: GetEngagementByID :one
SELECT * FROM engagements WHERE id = $1;

-- name: ListEngagementsByFirm :many
SELECT * FROM engagements
WHERE firm_id = $1
  AND (sqlc.narg('status')::engagement_status IS NULL OR status = sqlc.narg('status'))
  AND (sqlc.narg('client_id')::uuid IS NULL OR client_id = sqlc.narg('client_id'))
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateEngagement :one
UPDATE engagements
SET name         = COALESCE(sqlc.narg('name'), name),
    period_start = COALESCE(sqlc.narg('period_start'), period_start),
    period_end   = COALESCE(sqlc.narg('period_end'), period_end)
WHERE id = $1
RETURNING *;

-- name: UpdateEngagementStatus :one
UPDATE engagements SET status = $2 WHERE id = $1 RETURNING *;
```

- [ ] **Step 2: engagement_team.sql**

```sql
-- name: AddTeamMember :one
INSERT INTO engagement_team_members (firm_id, engagement_id, user_id, engagement_role)
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: RemoveTeamMember :one
UPDATE engagement_team_members SET removed_at = now()
WHERE engagement_id = $1 AND user_id = $2 AND removed_at IS NULL
RETURNING *;

-- name: ListActiveTeam :many
SELECT * FROM engagement_team_members
WHERE engagement_id = $1 AND removed_at IS NULL
ORDER BY assigned_at;
```

- [ ] **Step 3: engagement_frameworks.sql**

```sql
-- name: AddEngagementFramework :one
INSERT INTO engagement_frameworks (firm_id, engagement_id, framework_id, framework_version, is_primary)
VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: ListEngagementFrameworks :many
SELECT ef.*, f.name AS framework_name
FROM engagement_frameworks ef
JOIN frameworks f ON f.id = ef.framework_id
WHERE ef.engagement_id = $1
ORDER BY ef.is_primary DESC, f.name;
```

- [ ] **Step 4: client_acceptances.sql**

```sql
-- name: UpsertClientAcceptance :one
INSERT INTO client_acceptances (firm_id, engagement_id, quality_risks_identified, firm_responses)
VALUES ($1, $2, $3, $4)
ON CONFLICT (engagement_id) DO UPDATE
  SET quality_risks_identified = EXCLUDED.quality_risks_identified,
      firm_responses           = EXCLUDED.firm_responses
RETURNING *;

-- name: GetClientAcceptanceByEngagement :one
SELECT * FROM client_acceptances WHERE engagement_id = $1;

-- name: ConfirmIndependence :one
UPDATE client_acceptances
SET independence_confirmed = true,
    independence_confirmed_by_id = $2
WHERE engagement_id = $1
RETURNING *;

-- name: SignOffClientAcceptance :one
UPDATE client_acceptances
SET accepted_by_id = $2, accepted_at = now()
WHERE engagement_id = $1 AND independence_confirmed = true
RETURNING *;
```

- [ ] **Step 5: controls.sql and test_procedures.sql**

```sql
-- controls.sql
-- name: CreateControl :one
INSERT INTO controls (firm_id, engagement_id, firm_control_objective_id, description, is_key_control)
VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: ListControlsByEngagement :many
SELECT * FROM controls WHERE engagement_id = $1 ORDER BY created_at;

-- name: GetControlByID :one
SELECT * FROM controls WHERE id = $1;

-- name: UpdateControlStatus :one
UPDATE controls SET status = $2 WHERE id = $1 RETURNING *;

-- name: UpdateControlAssignment :one
UPDATE controls SET auditor_assigned_to_id = $2 WHERE id = $1 RETURNING *;
```

```sql
-- test_procedures.sql
-- name: CreateTestProcedure :one
INSERT INTO test_procedures (firm_id, control_id, procedure_type, description, expected_result, sort_order)
VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: ListProceduresByControl :many
SELECT * FROM test_procedures WHERE control_id = $1 ORDER BY sort_order, id;

-- name: UpdateProcedureStatus :one
UPDATE test_procedures SET status = $2 WHERE id = $1 RETURNING *;

-- name: UpdateProcedureResult :one
UPDATE test_procedures
SET result = $2, exceptions_noted = $3, conclusion = $4,
    performed_by_id = $5, performed_at = now()
WHERE id = $1 RETURNING *;
```

- [ ] **Step 6: sqlc generate and build**

```bash
sqlc generate
go build ./...
```

- [ ] **Step 7: Commit**

```bash
git add apps/axiom-api/internal/auditcore/queries/
git commit -m "feat(auditcore): sqlc queries for engagements, team, controls, procedures"
git push
```

---

## Task 15: Auditcore — engagement CRUD + scaffolding service

**Files:**
- Create: `apps/axiom-api/internal/auditcore/service.go`
- Create: `apps/axiom-api/internal/auditcore/scaffolding.go`
- Create: `apps/axiom-api/internal/auditcore/service_test.go`
- Create: `apps/axiom-api/internal/auditcore/scaffolding_test.go`

The **scaffolding** behavior is the headline feature of this phase: creating an engagement from a template must auto-create its controls and test procedures in a single transaction.

- [ ] **Step 1: Write failing scaffolding test**

`scaffolding_test.go`:

```go
package auditcore_test

import (
	"context"
	"testing"
	"time"

	"github.com/axiom-platform/axiom/apps/axiom-api/internal/auditcore"
	"github.com/axiom-platform/axiom/apps/axiom-api/internal/identity"
	"github.com/axiom-platform/axiom/apps/axiom-api/internal/platform"
)

func TestCreateEngagementFromSOC2Template(t *testing.T) {
	pool := platform.TestDB(t)
	// Register a firm (reuses Phase 1 code).
	idSvc := identity.NewService(pool, identity.NewJWTIssuer(testPrivKey, testPubKey))
	reg, err := idSvc.RegisterFirm(context.Background(), identity.RegisterFirmInput{
		FirmName: "Test LLP", AdminEmail: "a@b.com", AdminName: "A", Password: "pw-123456",
		Country: "US",
	})
	if err != nil { t.Fatal(err) }

	// Create a client.
	client, err := idSvc.CreateClient(context.Background(), reg.Firm.ID, identity.CreateClientInput{Name: "TechCorp"})
	if err != nil { t.Fatal(err) }

	// Find SOC 2 Type II Standard template and SOC 2 framework.
	tmpl, fw := mustFindSOC2Template(t, pool)

	svc := auditcore.NewService(pool)
	eng, err := svc.CreateEngagement(context.Background(), auditcore.CreateEngagementInput{
		FirmID:               reg.Firm.ID,
		ClientID:             client.ID,
		Name:                 "TechCorp SOC 2 2026",
		EngagementType:       "SOC2",
		PrimaryFrameworkID:   fw.ID,
		PeriodStart:          time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		PeriodEnd:            time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC),
		MethodologyTemplateID: &tmpl.ID,
	})
	if err != nil { t.Fatal(err) }

	// Scaffolding assertions.
	controls, err := svc.ListControls(context.Background(), eng.ID)
	if err != nil { t.Fatal(err) }
	if len(controls) < 40 {
		t.Errorf("want ≥40 scaffolded controls, got %d", len(controls))
	}
	// Spot-check that procedures exist for the first control.
	procs, err := svc.ListProcedures(context.Background(), controls[0].ID)
	if err != nil { t.Fatal(err) }
	if len(procs) == 0 {
		t.Error("want ≥1 procedure per scaffolded control")
	}
}

func TestCreateEngagementWithoutTemplateCreatesNoControls(t *testing.T) {
	// Same setup, pass nil MethodologyTemplateID, assert 0 controls.
}

func TestScaffoldingIsAtomic(t *testing.T) {
	// Inject a failure mid-scaffolding (e.g. force a constraint violation on a
	// template procedure row) and assert no engagement + no controls remain.
	// Use a nested transaction trick or a deliberately bad template state.
}
```

Add a helper `mustFindSOC2Template` and test key material in a `testhelpers_test.go`.

- [ ] **Step 2: Run — expect failure**

- [ ] **Step 3: Implement the service and scaffolding**

`service.go` exposes `CreateEngagement` which delegates to `scaffolding.go::scaffoldFromTemplate(tx, engagementID, templateID)`.

Key rules:
- **Single transaction.** The engagement row, team membership for the creating Partner, `engagement_frameworks` row for the primary framework, and every scaffolded control + procedure are written inside one `pool.Begin` → `Commit` block.
- **RLS context.** Call `SET LOCAL app.current_firm_id = $1` inside the transaction before any RLS-scoped write.
- **Scaffolding** reads `template_controls` and `template_test_procedures` joined to the template, then inserts a `controls` row per template control (carrying `firm_control_objective_id`) and a `test_procedures` row per template procedure (cascading the `control_id` that was just inserted).
- Preserve `sort_order` in procedures.
- **Period validation** is enforced before the transaction begins:
  - For SOC Type 1 engagements (`engagement_type` SOC1 or SOC2 *and* `report_type` SOC*Type1` from the wizard), require `period_start == period_end`. Reject otherwise.
  - For SOC Type 2 engagements, require `period_end > period_start` AND `period_end - period_start BETWEEN 90 AND 366 days` (3 to 12 months). Outside this range, require an explicit `partner_override_reason` field on the input; record the override in the audit log.
  - For ISO certification cycles and PCI engagements, validate per the framework's native window rules.

Skeleton:

```go
func (s *Service) CreateEngagement(ctx context.Context, in CreateEngagementInput) (*EngagementDTO, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil { return nil, err }
	defer func() { _ = tx.Rollback(ctx) }()

	if _, err := tx.Exec(ctx, "SET LOCAL app.current_firm_id = $1", in.FirmID.String()); err != nil {
		return nil, err
	}

	qtx := s.q.WithTx(tx)

	eng, err := qtx.CreateEngagement(ctx, queries.CreateEngagementParams{...})
	if err != nil { return nil, err }

	// Primary framework link.
	_, _ = qtx.AddEngagementFramework(ctx, queries.AddEngagementFrameworkParams{
		FirmID: in.FirmID, EngagementID: eng.ID,
		FrameworkID: in.PrimaryFrameworkID, FrameworkVersion: "…lookup…", IsPrimary: true,
	})

	if in.MethodologyTemplateID != nil {
		if err := scaffoldFromTemplate(ctx, qtx, in.FirmID, eng.ID, *in.MethodologyTemplateID); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil { return nil, err }
	return toEngagementDTO(eng), nil
}
```

`scaffoldFromTemplate` queries `template_controls` (sorted) and inserts `controls`, tracking a `map[templateControlID]controlID`. Then queries `template_test_procedures` and inserts `test_procedures` with `control_id = map[tp.TemplateControlID]`.

- [ ] **Step 4: Run — expect pass**

```bash
go test ./internal/auditcore/...
```

- [ ] **Step 5: Add remaining CRUD tests and methods**

Add to `service.go`:
- `ListEngagements(firmID, filters)`
- `GetEngagement(id)` returning team + primary framework
- `UpdateEngagement(id, in)` (name, period)
- `AddTeamMember / RemoveTeamMember`
- `AddEngagementFramework` (add a secondary framework after creation)

Test each with `TestDB(t)`.

- [ ] **Step 6: Commit**

```bash
git add apps/axiom-api/internal/auditcore/service.go \
        apps/axiom-api/internal/auditcore/scaffolding.go \
        apps/axiom-api/internal/auditcore/service_test.go \
        apps/axiom-api/internal/auditcore/scaffolding_test.go
git commit -m "feat(auditcore): engagement CRUD with atomic template-driven scaffolding"
git push
```

---

## Task 16: Auditcore — audit log writer

**Files:**
- Create: `apps/axiom-api/internal/auditcore/auditlog.go`
- Create: `apps/axiom-api/internal/auditcore/auditlog_test.go`

A typed helper every auditcore write path calls.

- [ ] **Step 1: Write failing test**

```go
func TestAuditLogRecordsEngagementCreation(t *testing.T) {
	pool := platform.TestDB(t)
	logger := auditcore.NewAuditLogger(pool)

	firmID := uuid.New(); userID := uuid.New(); engID := uuid.New()
	if err := logger.Record(context.Background(), auditcore.AuditEvent{
		FirmID:       firmID,
		ActorID:      &userID,
		ActorType:    "User",
		Action:       "engagement.created",
		ResourceType: "engagement",
		ResourceID:   &engID,
		NewValue:     json.RawMessage(`{"name":"Foo"}`),
	}); err != nil { t.Fatal(err) }

	var count int
	_ = pool.QueryRow(ctx,
		"SELECT COUNT(*) FROM audit_log WHERE action = 'engagement.created' AND firm_id = $1",
		firmID).Scan(&count)
	if count != 1 { t.Errorf("want 1 row, got %d", count) }
}

func TestAuditLogIsAppendOnly(t *testing.T) {
	// Insert a row, try UPDATE/DELETE, assert row still exists unchanged.
}
```

- [ ] **Step 2: Run — expect failure**

- [ ] **Step 3: Implement**

`auditlog.go` defines `AuditEvent` and `AuditLogger.Record` that INSERTs into `audit_log`. Use raw SQL (simpler than sqlc for a single table with jsonb diffs).

- [ ] **Step 4: Wire into CreateEngagement and state transitions**

Call `logger.Record` inside the scaffolding transaction for `engagement.created` and, later, for `engagement.status.changed` and `control.status.changed`. The audit log writer participates in the same transaction as the business change.

- [ ] **Step 5: Run — expect pass**

- [ ] **Step 6: Commit**

```bash
git add apps/axiom-api/internal/auditcore/auditlog*.go apps/axiom-api/internal/auditcore/scaffolding.go
git commit -m "feat(auditcore): append-only audit log writer wired into engagement writes"
git push
```

---

## Task 17: Auditcore — client acceptance workflow

**Files:**
- Create: `apps/axiom-api/internal/auditcore/client_acceptance.go`
- Create: `apps/axiom-api/internal/auditcore/client_acceptance_test.go`

Partner-only sign-off that unblocks the Planning → Fieldwork transition.

- [ ] **Step 1: Write failing tests**

Cases (in order):
1. `UpsertClientAcceptance` saves quality risks + firm responses.
2. `ConfirmIndependence` is **Partner-only** — Manager calls return `ErrForbidden`.
3. `SignOffClientAcceptance` requires `independence_confirmed = true`; otherwise returns a validation error.
4. Successful sign-off sets `accepted_at` and records an `audit_log` entry `client_acceptance.signed_off`.

- [ ] **Step 2: Run — expect failure**

- [ ] **Step 3: Implement**

```go
func (s *Service) UpsertClientAcceptance(ctx context.Context, firmID, engID uuid.UUID, in ClientAcceptanceInput) (*ClientAcceptanceDTO, error)
func (s *Service) ConfirmIndependence(ctx context.Context, firmID, engID, partnerID uuid.UUID, role string) error
func (s *Service) SignOffClientAcceptance(ctx context.Context, firmID, engID, partnerID uuid.UUID, role string) (*ClientAcceptanceDTO, error)
```

Role check inside the service (defense in depth — the handler also gates it). The service does the hardening; handlers are dumb glue.

- [ ] **Step 4: Run — expect pass**

- [ ] **Step 5: Commit**

```bash
git add apps/axiom-api/internal/auditcore/client_acceptance*.go
git commit -m "feat(auditcore): client acceptance workflow with partner sign-off"
git push
```

---

## Task 18: Auditcore — state machine (Planning ↔ Fieldwork)

**Files:**
- Create: `apps/axiom-api/internal/auditcore/state_machine.go`
- Create: `apps/axiom-api/internal/auditcore/state_machine_test.go`

- [ ] **Step 1: Write failing tests**

Cases:
1. `Advance(Planning → Fieldwork)` succeeds when `client_acceptance.accepted_at IS NOT NULL`.
2. `Advance(Planning → Fieldwork)` returns a structured validation error (`independence_confirmed=false`, `accepted_at IS NULL`) when acceptance isn't signed. The error payload includes `blockers: ["client_acceptance_not_signed"]`.
3. `Revert(Fieldwork → Planning)` is Partner-only; Manager/Staff return `ErrForbidden`.
4. An audit-log entry `engagement.status.changed` is written on every successful transition, with `old_value` and `new_value` capturing the old and new status.
5. No other transitions are accepted (attempt `Planning → Review` returns `ErrBadRequest`).

- [ ] **Step 2: Run — expect failure**

- [ ] **Step 3: Implement**

```go
type Transition struct { From, To string }

func (s *Service) AdvanceEngagement(ctx context.Context, firmID, engID, actorID uuid.UUID, role string) (*EngagementDTO, error)
func (s *Service) RevertEngagement(ctx context.Context,  firmID, engID, actorID uuid.UUID, role string) (*EngagementDTO, error)
```

Structure:

```go
func (s *Service) guardPlanningToFieldwork(ctx context.Context, q *queries.Queries, engID uuid.UUID) []string {
	ca, err := q.GetClientAcceptanceByEngagement(ctx, engID)
	if err != nil { return []string{"client_acceptance_missing"} }
	var blockers []string
	if !ca.IndependenceConfirmed { blockers = append(blockers, "independence_not_confirmed") }
	if !ca.AcceptedAt.Valid       { blockers = append(blockers, "client_acceptance_not_signed") }
	return blockers
}
```

Return the blocker list in a typed error:

```go
type TransitionBlocked struct { Blockers []string }
func (e *TransitionBlocked) Error() string { return "transition blocked" }
```

Handler layer serializes this into a 422 response with the blockers array.

- [ ] **Step 4: Run — expect pass**

- [ ] **Step 5: Commit**

```bash
git add apps/axiom-api/internal/auditcore/state_machine*.go
git commit -m "feat(auditcore): engagement state machine with guarded planning↔fieldwork transitions"
git push
```

---

## Task 19: Auditcore — HTTP handlers

**Files:**
- Create: `apps/axiom-api/internal/auditcore/handler.go`
- Create: `apps/axiom-api/internal/auditcore/handler_test.go`

- [ ] **Step 1: Write failing handler tests**

Cover every endpoint below. Use the identity test helpers to mint a JWT and exercise the full router stack (auth → role guard → handler → DB).

- [ ] **Step 2: Run — expect failure**

- [ ] **Step 3: Implement routes**

```go
func (h *Handler) RegisterRoutes(r chi.Router, gw RoleGuard) {
	// Engagements
	r.Get("/api/v1/engagements", h.listEngagements)
	r.With(gw.WithRole("FirmAdmin", "Partner")).
		Post("/api/v1/engagements", h.createEngagement)
	r.Get("/api/v1/engagements/{engagementId}", h.getEngagement)
	r.With(gw.WithRole("FirmAdmin", "Partner", "Manager")).
		Patch("/api/v1/engagements/{engagementId}", h.updateEngagement)
	r.With(gw.WithRole("FirmAdmin", "Partner")).
		Post("/api/v1/engagements/{engagementId}/transitions", h.transitionEngagement)

	// Team
	r.Get("/api/v1/engagements/{engagementId}/team", h.listTeam)
	r.With(gw.WithRole("FirmAdmin", "Partner", "Manager")).
		Post("/api/v1/engagements/{engagementId}/team", h.addTeamMember)
	r.With(gw.WithRole("FirmAdmin", "Partner", "Manager")).
		Delete("/api/v1/engagements/{engagementId}/team/{userId}", h.removeTeamMember)

	// Frameworks on engagement
	r.Get("/api/v1/engagements/{engagementId}/frameworks", h.listEngFrameworks)
	r.With(gw.WithRole("FirmAdmin", "Partner")).
		Post("/api/v1/engagements/{engagementId}/frameworks", h.addEngFramework)

	// Client acceptance
	r.Get("/api/v1/engagements/{engagementId}/client-acceptance", h.getClientAcceptance)
	r.With(gw.WithRole("Partner", "Manager", "FirmAdmin")).
		Put("/api/v1/engagements/{engagementId}/client-acceptance", h.upsertClientAcceptance)
	r.With(gw.WithRole("Partner", "FirmAdmin")).
		Post("/api/v1/engagements/{engagementId}/client-acceptance/independence", h.confirmIndependence)
	r.With(gw.WithRole("Partner", "FirmAdmin")).
		Post("/api/v1/engagements/{engagementId}/client-acceptance/signoff", h.signOffClientAcceptance)

	// Controls & procedures
	r.Get("/api/v1/engagements/{engagementId}/controls", h.listControls)
	r.Get("/api/v1/controls/{controlId}", h.getControl)
	r.Patch("/api/v1/controls/{controlId}", h.updateControl)
	r.Get("/api/v1/controls/{controlId}/procedures", h.listProcedures)
	r.Patch("/api/v1/procedures/{procedureId}", h.updateProcedure)
}
```

The transition endpoint accepts `{ "to": "Fieldwork" }` or `{ "to": "Planning" }`. On `TransitionBlocked` errors, respond with 422 and `{ "error": "transition blocked", "blockers": [...] }`.

- [ ] **Step 4: Run — expect pass**

- [ ] **Step 5: Commit**

```bash
git add apps/axiom-api/internal/auditcore/handler*.go
git commit -m "feat(auditcore): HTTP handlers for engagements, team, client acceptance, transitions"
git push
```

---

## Task 20: Auditcore — multi-tenant RLS tests

**Files:**
- Create: `apps/axiom-api/internal/auditcore/rls_test.go`

Prove that RLS prevents cross-firm reads for every tenant-scoped table introduced this phase.

- [ ] **Step 1: Write failing test**

Pattern: create two firms (A, B); create an engagement + control + client acceptance under each; `SET LOCAL app.current_firm_id = firm_A`; list engagements/controls/client_acceptances/engagement_team_members; assert only firm A's rows visible.

Include a test that writing with a spoofed `firm_id` fails when RLS is set to a different firm (the INSERT raises because the row wouldn't satisfy the policy).

- [ ] **Step 2: Run — expect failure**

- [ ] **Step 3: Implement (RLS already in migrations)**

If the migrations are correct, the tests should pass immediately. If any fail, fix the migration (missing `ENABLE ROW LEVEL SECURITY` or policy) rather than weakening the test.

- [ ] **Step 4: Run — expect pass**

- [ ] **Step 5: Commit**

```bash
git add apps/axiom-api/internal/auditcore/rls_test.go
git commit -m "test(auditcore): multi-tenant RLS isolation for engagements, controls, acceptances"
git push
```

---

## Task 21: Wire modules into main.go

**Files:**
- Modify: `apps/axiom-api/cmd/server/main.go`

- [ ] **Step 1: Wire frameworks + auditcore**

Add imports and construction after the existing identity wiring:

```go
"github.com/axiom-platform/axiom/apps/axiom-api/internal/auditcore"
"github.com/axiom-platform/axiom/apps/axiom-api/internal/frameworks"
```

```go
frameworksSvc     := frameworks.NewService(pool)
frameworksHandler := frameworks.NewHandler(frameworksSvc)

auditLogger       := auditcore.NewAuditLogger(pool)
auditcoreSvc      := auditcore.NewService(pool, auditLogger)
auditcoreHandler  := auditcore.NewHandler(auditcoreSvc)
```

Register routes — frameworks read endpoints are public-to-firm (any authenticated user); auditcore uses role-guarded routes:

```go
r.Group(func(r chi.Router) {
    r.Use(gw.Auth)
    identityHandler.RegisterAuthenticatedRoutes(r, gw)
    identityHandler.RegisterMethodologyRoutes(r, gw)
    frameworksHandler.RegisterRoutes(r)
    auditcoreHandler.RegisterRoutes(r, gw)
})
```

- [ ] **Step 2: Run the server**

```bash
cd apps/axiom-api
go run ./cmd/server
```

In another terminal:

```bash
curl -s -X POST http://localhost:8080/api/v1/auth/register -H 'content-type: application/json' \
  -d '{"firm_name":"Sample LLP","admin_email":"a@b.com","admin_name":"A","password":"hunter22hunter22","country":"US"}' \
  | jq -r .access_token > /tmp/tok

curl -s http://localhost:8080/api/v1/frameworks -H "authorization: Bearer $(cat /tmp/tok)" | jq '.items | length'
# expect 4

curl -s http://localhost:8080/api/v1/methodology-templates -H "authorization: Bearer $(cat /tmp/tok)" | jq '.items | length'
# expect 2 (the two system templates)
```

- [ ] **Step 3: Commit**

```bash
git add apps/axiom-api/cmd/server/main.go
git commit -m "feat(api): wire frameworks and auditcore modules into server"
git push
```

---

## Task 22: Regenerate frontend OpenAPI types

**Files:**
- Create: `apps/web/src/api/generated/audit-core.ts`

- [ ] **Step 1: Run openapi-typescript for the audit-core spec**

```bash
cd apps/web
npx openapi-typescript ../../packages/openapi/audit-core.yaml -o src/api/generated/audit-core.ts
```

- [ ] **Step 2: Verify the import graph**

```bash
npm run build
```

Expected: build succeeds (no references to the new file are required yet; later tasks import from it).

- [ ] **Step 3: Commit**

```bash
git add apps/web/src/api/generated/audit-core.ts
git commit -m "chore(web): generate TS types from audit-core OpenAPI spec"
git push
```

---

## Task 23: Frontend — shared Tabs and Wizard components

**Files:**
- Create: `apps/web/src/components/tabs.tsx`
- Create: `apps/web/src/components/tabs.test.tsx`
- Create: `apps/web/src/components/wizard.tsx`
- Create: `apps/web/src/components/wizard.test.tsx`

Both components are used by multiple Phase 2 pages. Build them first.

- [ ] **Step 1: Read the mockups that exercise these patterns**

- Tabs: `mockups/journey-03-engagement-scoping/02-engagement-details.html`
- Wizard: `mockups/journey-01-firm-setup/06-create-engagement.html` (has a step indicator in the left rail)

Invoke the `frontend-design` skill on both components to generate impeccable token-aligned styles.

- [ ] **Step 2: Write failing component tests**

For `tabs.tsx`:
- Renders a list of tab labels; clicking a tab calls `onSelect` with that tab's key.
- The `activeKey` prop controls which tab is rendered as active.
- Initial render with no `activeKey` selects the first tab.

For `wizard.tsx`:
- Given N steps, renders the current step's children.
- "Next" moves forward; "Back" moves backward; "Back" is disabled on step 0; "Next" becomes "Finish" on the last step.
- Emits `onComplete` when "Finish" is clicked.

- [ ] **Step 3: Run — expect failure**

```bash
cd apps/web
npm test -- components/tabs.test.tsx components/wizard.test.tsx
```

- [ ] **Step 4: Implement both components**

Use semantic HTML (`<button role="tab">`, `<ol>` for wizard steps) and map styles from `.impeccable.md`. Avoid inline styles; use CSS classes in a co-located `.css` file per component (match the Phase 0/1 pattern in `components/layout.css`).

- [ ] **Step 5: Run — expect pass**

- [ ] **Step 6: Commit**

```bash
git add apps/web/src/components/tabs* apps/web/src/components/wizard*
git commit -m "feat(web): shared tabs and wizard components following impeccable design system"
git push
```

---

## Task 24: Frontend — templates list, detail, and firm control objectives

**Files:**
- Create: `apps/web/src/pages/templates.tsx`
- Create: `apps/web/src/pages/templates.test.tsx`
- Create: `apps/web/src/pages/templates-detail.tsx`
- Create: `apps/web/src/pages/templates-detail.test.tsx`
- Create: `apps/web/src/pages/firm-control-objectives.tsx`
- Create: `apps/web/src/pages/firm-control-objectives.test.tsx`
- Modify: `apps/web/src/App.tsx`

- [ ] **Step 1: Read mockups and plan impeccable**

- Templates list + detail: `mockups/journey-01-firm-setup/05-methodology-templates.html`.
- Firm control objectives: no mockup. Invoke `frontend-design` skill with `.impeccable.md` to design a list + detail pattern aligned with the templates page.

- [ ] **Step 2: Write failing page tests (templates)**

`templates.test.tsx`:
- Renders a grid of templates from a mocked fetch; system-provided templates show an "Included" badge.
- FirmAdmin sees an "Activate" / "Deactivate" toggle; non-FirmAdmin does not.
- Clicking a template card navigates to `/templates/:id`.

`templates-detail.test.tsx`:
- Fetches template detail and renders counts ("50 controls · 80 test procedures · 80 document requests") and tabs (Controls, Procedures, Doc Requests).

`firm-control-objectives.test.tsx`:
- Lists objectives with their mapping counts.
- Clicking "Add from library" opens a modal that lists library entries and adds on selection.
- Clicking "Add custom" opens a form; submitting creates the objective.
- Mapping editor: a framework-requirement picker adds and removes mappings.

Mock `fetch` via the existing `client.ts` pattern from Phase 0/1.

- [ ] **Step 3: Run — expect failure**

- [ ] **Step 4: Implement the pages**

Use TanStack Query for list/detail fetches. Reuse the `client.ts` wrapper from Phase 0/1 so JWT injection is automatic.

- [ ] **Step 5: Wire routes**

Edit `apps/web/src/App.tsx` to add:

```tsx
<Route path="/templates" element={<ProtectedRoute><TemplatesPage /></ProtectedRoute>} />
<Route path="/templates/:id" element={<ProtectedRoute><TemplateDetailPage /></ProtectedRoute>} />
<Route path="/control-objectives" element={<ProtectedRoute><FirmControlObjectivesPage /></ProtectedRoute>} />
```

Add sidebar links in `components/layout.tsx`.

- [ ] **Step 6: Run all web tests**

```bash
npm test
```

- [ ] **Step 7: Commit**

```bash
git add apps/web/src/pages/templates* apps/web/src/pages/firm-control-objectives* apps/web/src/App.tsx apps/web/src/components/layout.tsx
git commit -m "feat(web): methodology templates and firm control objectives pages"
git push
```

---

## Task 25: Frontend — engagements list, create wizard, and detail

**Files:**
- Create: `apps/web/src/pages/engagements.tsx`
- Create: `apps/web/src/pages/engagements.test.tsx`
- Create: `apps/web/src/pages/create-engagement.tsx`
- Create: `apps/web/src/pages/create-engagement.test.tsx`
- Create: `apps/web/src/pages/engagement-detail.tsx`
- Create: `apps/web/src/pages/engagement-detail.test.tsx`
- Modify: `apps/web/src/App.tsx`, `apps/web/src/components/layout.tsx`

The create-engagement wizard is the largest UI unit. Walk through the five steps carefully.

- [ ] **Step 1: Read the mockups**

- Engagement list: `mockups/journey-01-firm-setup/09-onboarding-complete.html` (for the shell/list pattern; no dedicated list mockup exists — invoke `frontend-design` for the list).
- Create wizard:
  - Step 1 (Client): `mockups/journey-01-firm-setup/06-create-engagement.html`
  - Step 2 (Type + framework): `mockups/journey-03-engagement-scoping/01-new-engagement-type.html`
  - Step 3 (Template + review scaffold): `mockups/journey-03-engagement-scoping/02-engagement-details.html`
  - Step 4 (Team): `mockups/journey-03-engagement-scoping/03-team-assignment.html`
  - Step 5 (Confirm): `mockups/journey-01-firm-setup/07-engagement-ready.html`
- Engagement detail (tabs: Overview, Controls, Team, Client Acceptance):
  - Overview: `mockups/journey-01-firm-setup/07-engagement-ready.html`
  - Controls: `mockups/journey-03-engagement-scoping/04-ai-control-mapping.html` (drop the AI-suggestion column for now; Phase 7 adds it)
  - Client Acceptance: `mockups/journey-03-engagement-scoping/05-client-acceptance.html`
  - Advance to Fieldwork: `mockups/journey-03-engagement-scoping/07-begin-fieldwork.html`

- [ ] **Step 2: Write failing engagements-list tests**

`engagements.test.tsx`:
- Renders a table of engagements with status badges, client name, period, engagement type.
- Filters: by status, by client. Filter selection triggers refetch.
- "New engagement" button navigates to `/engagements/new` (FirmAdmin/Partner only).

- [ ] **Step 3: Write failing create-engagement wizard tests**

`create-engagement.test.tsx` covers each step:

- Step 1 (Client): existing-client dropdown populates from `/api/v1/clients`. "Create new client" opens an inline form; submit uses `POST /api/v1/clients`.
- Step 2 (Type + framework): selecting an engagement type filters the framework dropdown. The mapping (engagement_type → primary framework version) is:
  - `SOC1` → SOC 1 SSAE 18
  - `SOC2` → SOC 2 TSC 2017
  - `ISO27001` → ISO 27001:2022
  - `ISO27701` → ISO 27701:2019
  - `ISO42001` → ISO 42001:2023
  - `HIPAA` → HIPAA Security Rule 2013
  - `PCI_DSS` → PCI DSS v4.0.1
  - `AgreedUponProcedures`, `Advisory` → no preselected framework (user picks any).
- Step 3 (Template): template dropdown filtered by engagement type. When selected, fetches `/api/v1/methodology-templates/:id` and renders a preview with control + procedure counts. "Skip template" is allowed.
- Step 4 (Team): multi-select for Partner/Manager/Staff assignment (required: exactly one Partner).
- Step 5 (Confirm): summary; Submit calls `POST /api/v1/engagements`, then navigates to `/engagements/:id`.

- [ ] **Step 4: Write failing engagement-detail tests**

`engagement-detail.test.tsx`:
- Tabs render correctly (Overview, Controls, Team, Client Acceptance).
- Overview shows status, client name, framework, period, Partner.
- Controls tab shows scaffolded controls from the template in their sort order.
- Client Acceptance tab renders an editable quality-risk form; Partner-only "Confirm independence" and "Sign off" buttons.
- "Advance to Fieldwork" button is disabled when `client_acceptance.accepted_at IS NULL`; the tooltip lists blockers.
- Clicking Advance when unlocked calls `POST /api/v1/engagements/:id/transitions { to: 'Fieldwork' }`; page refetches and status shows Fieldwork.

- [ ] **Step 5: Run — expect failures**

```bash
cd apps/web
npm test -- engagements create-engagement engagement-detail
```

- [ ] **Step 6: Implement the pages**

Key points:
- Use the shared `Wizard` from Task 23. Each step is a controlled component that reads/writes a top-level `formState` via props.
- Validate per step before enabling "Next".
- Use TanStack Query mutations for `POST /api/v1/engagements`; on success, invalidate `['engagements']`.
- Engagement detail uses the shared `Tabs` component.
- Blockers rendering: after a 422 with `{ "blockers": [...] }`, render each blocker as a human-readable string. Maintain a single `blockers.ts` util that maps machine-readable blockers to copy. (This is the only new util file; keep it under `apps/web/src/api/blockers.ts`.)

- [ ] **Step 7: Wire routes**

```tsx
<Route path="/engagements" element={<ProtectedRoute><EngagementsPage /></ProtectedRoute>} />
<Route path="/engagements/new" element={<ProtectedRoute><CreateEngagementPage /></ProtectedRoute>} />
<Route path="/engagements/:id" element={<ProtectedRoute><EngagementDetailPage /></ProtectedRoute>} />
```

- [ ] **Step 8: Run all tests**

```bash
npm test
```

Everything should pass.

- [ ] **Step 9: Commit**

```bash
git add apps/web/src/pages/engagements* apps/web/src/pages/create-engagement* \
        apps/web/src/pages/engagement-detail* apps/web/src/api/blockers.ts \
        apps/web/src/App.tsx apps/web/src/components/layout.tsx
git commit -m "feat(web): engagement list, create wizard, detail with client acceptance and transitions"
git push
```

---

## Task 26: Impeccable validation pass

**Files:** None created — this is a review pass.

Every page built or modified this phase must pass impeccable review before the phase is considered complete.

- [ ] **Step 1: Assemble the list**

Pages to review:
- `pages/templates.tsx`, `pages/templates-detail.tsx`
- `pages/firm-control-objectives.tsx`
- `pages/engagements.tsx`
- `pages/create-engagement.tsx` (all 5 steps)
- `pages/engagement-detail.tsx` (all 4 tabs)

Shared components:
- `components/tabs.tsx`, `components/wizard.tsx`
- `components/layout.tsx` (updated sidebar)

- [ ] **Step 2: Run `frontend-design` skill review for each**

For each page + shared component, invoke the `frontend-design` skill with the context: "Validate against `.impeccable.md` — typography scale, OKLCH color tokens, spacing tokens, voice and copy, component hierarchy." Capture findings and address them.

- [ ] **Step 3: Cross-reference with Phase 0/1 pages**

Open `pages/users.tsx`, `pages/clients.tsx`, and `pages/firm-settings.tsx` side-by-side with the new pages. Visual density, header patterns, empty-state copy, and button hierarchy should match.

- [ ] **Step 4: Run the server + web and manually walk the end-to-end flow**

```bash
# terminal 1
cd apps/axiom-api && go run ./cmd/server

# terminal 2
cd apps/web && npm run dev
```

Walk the full Phase 2 testable outcome (see Task 27) in the browser. Note any polish issues (misaligned baselines, loading states, focus outlines) and fix them.

- [ ] **Step 5: Commit any fixes**

```bash
git add apps/web/src/
git commit -m "fix(web): impeccable polish for phase 2 pages"
git push
```

---

## Task 27: Manual testing doc

**Files:**
- Create: `docs/superpowers/testing/phase-2-frameworks-templates-engagements.md`

The plan-level testing methodology (`docs/superpowers/specs/implementation-plan-design.md` §"Per-Phase Manual Testing Instructions") requires a checklist a non-engineer can follow.

- [ ] **Step 1: Write the doc**

Structure (verbatim to phase 0/1 convention):

```markdown
# Phase 2: Frameworks, Templates & Engagement Creation — Manual Testing

## Prerequisites

- Repo at tip of `phase-2-frameworks-templates-engagements` branch.
- `docker compose up -d` running (Postgres + Mailhog healthy).
- No migrations pending: `migrate -database … -path apps/axiom-api/migrations up` is a no-op.
- Backend running: `go run ./apps/axiom-api/cmd/server`.
- Frontend running: `cd apps/web && npm run dev`.

## Test data

- A fresh firm registered via `/register` — call it "Acme LLP" with FirmAdmin.
- Invite a Partner (via Phase 1 staff invitations) — call her "Pat Partner".
- Invite a Manager and two Staff.

## Scenario 1 — Activate a system template (FirmAdmin)

1. Log in as FirmAdmin.
2. Navigate to **Templates**.
3. Verify six system templates appear: "SOC 2 Type II Standard", "ISO 27001:2022 Standard", "ISO 27701:2019 Standard", "ISO 42001:2023 Standard", "HIPAA Security Rule Standard", "PCI DSS v4.0.1 Standard".
4. Click **SOC 2 Type II Standard**.
5. Verify detail panel shows ~50 controls, ~80 test procedures, ~80 document requests.
6. Click **Activate** — confirm the toggle flips and a success toast appears.

## Scenario 2 — Create an engagement with auto-scaffolded controls (Partner)

1. Log in as Pat Partner.
2. Navigate to **Engagements** → click **New engagement**.
3. Step 1: Pick client "TechCorp" or create one.
4. Step 2: Select engagement type `SOC 2`, primary framework `SOC 2 TSC 2017`.
5. Step 3: Select template `SOC 2 Type II Standard`. Verify preview shows 50/80/80 counts.
6. Step 4: Assign yourself as Partner, a Manager, and two Staff.
7. Step 5: Confirm summary → **Create engagement**.
8. You are redirected to `/engagements/:id`. Controls tab shows ~50 controls matching the template.

## Scenario 3 — Client acceptance blocks fieldwork (Partner)

1. On the engagement page, open the **Client Acceptance** tab.
2. Click **Advance to Fieldwork** — verify it's disabled and the tooltip says "Client acceptance not signed".
3. Fill in two quality risks and firm responses. Save.
4. Click **Confirm independence** — succeeds (Partner).
5. Sign in as Manager, reload the page. Independence confirm button is hidden; sign-off button is hidden.
6. Return to Partner. Click **Sign off**.
7. `Advance to Fieldwork` is now enabled. Click it.
8. Engagement status updates to **Fieldwork**.

## Scenario 4 — Cross-firm isolation

1. Register a second firm "BigAudit LLP".
2. Log in as BigAudit's FirmAdmin.
3. Navigate to **Engagements** — the Acme engagement must NOT appear.
4. Attempt to request Acme's engagement directly with `curl`:
   ```
   curl -H "authorization: Bearer <BigAudit_TOKEN>" http://localhost:8080/api/v1/engagements/<ACME_ENG_ID>
   ```
   Expected: 404 (RLS masks the row).

## Scenario 5 — Audit log is append-only

After the transition in Scenario 3, connect to Postgres and run:
```sql
SELECT action, resource_type, old_value->>'status' AS from, new_value->>'status' AS to
FROM audit_log WHERE firm_id = '<ACME_FIRM_ID>' ORDER BY occurred_at;
```
Expected rows (chronological):
- `engagement.created`
- `client_acceptance.signed_off`
- `engagement.status.changed` with `from='Planning'`, `to='Fieldwork'`

Attempt `UPDATE audit_log SET action = 'x' WHERE id = 1` — verify the row is unchanged.
```

- [ ] **Step 2: Commit**

```bash
git add docs/superpowers/testing/phase-2-frameworks-templates-engagements.md
git commit -m "docs: manual testing guide for phase 2"
git push
```

---

## Task 28: Open the PR

- [ ] **Step 1: Run full test + build from a clean state**

```bash
cd apps/axiom-api && go test ./... && go build ./...
cd apps/web && npm test && npm run lint && npm run build
```

All green before opening.

- [ ] **Step 2: Open the PR**

```bash
gh pr create --title "Phase 2: frameworks, templates, and engagement creation" --body "$(cat <<'EOF'
## Summary
- Adds system reference data (four frameworks + cross-framework control objective library) and two system methodology templates.
- Adds `frameworks` and `auditcore` modules with engagement CRUD, atomic template-driven scaffolding, client acceptance workflow, and the Planning↔Fieldwork state machine.
- Wires `audit_log` as the append-only journal for every engagement and control mutation.
- Ships templates pages, firm control objectives, engagement list, a 5-step create-engagement wizard, and a tabbed engagement detail.
- Adds RLS on every new tenant-scoped table; multi-tenant isolation tests added.

## Test plan
- [ ] Manual test doc in `docs/superpowers/testing/phase-2-frameworks-templates-engagements.md` passes end to end.
- [ ] `go test ./...` in `apps/axiom-api` passes.
- [ ] `npm test` in `apps/web` passes.
- [ ] Impeccable review pass completed for every new/changed page.

🤖 Generated with [Claude Code](https://claude.com/claude-code)
EOF
)"
```

- [ ] **Step 3: Review the PR, squash-merge to master, pull master**

```bash
git checkout master
git pull origin master
```

Phase 2 is complete.

---

## Plan Self-Review Notes

Coverage check against `implementation-plan-design.md` §Phase 2:

- Backend migrations — system reference, firm methodology, engagement, cross-cutting: **Tasks 1–5** ✓
- Seed data — frameworks, library, templates: **Tasks 6–8** ✓
- Identity module additions (methodology CRUD, firm control objectives, mapping management): **Tasks 11–13** ✓
- Audit Core module (engagement CRUD from template, team, framework mgmt, client acceptance, controls/procedures CRUD, state machine, audit log): **Tasks 14–20** ✓
- Frontend (templates browse/activate, firm control objectives, engagement list, create wizard, detail page, phase transition): **Tasks 23–25** ✓
- Impeccable + manual testing + PR: **Tasks 26–28** ✓

All spec requirements have a task. No placeholders. Method signatures referenced across tasks use consistent names (`CreateEngagement`, `AdvanceEngagement`, `SetTemplateActivation`, `UpsertClientAcceptance`).
