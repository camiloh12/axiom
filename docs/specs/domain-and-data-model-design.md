# Axiom Domain and Data Model

**Derived from:** [User Journeys](../user-journeys/all-journeys.md) (10 journeys, clean-slate domain modeling exercise)
**Supersedes:** Section 5 of the [product spec](axiom-spec-design.md) (which now contains an abbreviated summary pointing here)
**Related:** [Backend Architecture](backend-architecture-design.md) (service decomposition and database topology)

---

## Table of Contents

1. [Bounded Context Map](#1-bounded-context-map)
2. [Context 1: Firm Identity](#2-context-1-firm-identity)
3. [Context 2: Regulatory Framework & Common Controls](#3-context-2-regulatory-framework--common-controls)
4. [Context 3: Firm Methodology](#4-context-3-firm-methodology)
5. [Context 4: Audit Core](#5-context-4-audit-core)
6. [Context 5: Workpaper Authoring](#6-context-5-workpaper-authoring)
7. [Context 6: Reporting & Archival](#7-context-6-reporting--archival)
8. [Cross-Cutting Concerns](#8-cross-cutting-concerns)
9. [Data Model](#9-data-model)
10. [Multi-Tenancy and Isolation](#10-multi-tenancy-and-isolation)
11. [Journey-to-Entity Traceability](#11-journey-to-entity-traceability)

---

## 1. Bounded Context Map

Six bounded contexts and three cross-cutting concerns, derived from the coupling patterns observed across the compliance and assurance user journeys.

```
┌─────────────────────────────────────────────────────────────────────┐
│                         AXIOM DOMAIN                                │
│                                                                     │
│  ┌──────────────┐   ┌──────────────────┐   ┌───────────────────┐   │
│  │ 1. FIRM      │   │ 2. REGULATORY    │   │ 3. FIRM           │   │
│  │    IDENTITY  │   │    FRAMEWORK &   │   │    METHODOLOGY    │   │
│  │              │   │    COMMON CTRLS  │   │                   │   │
│  │ Firm         │   │ Framework        │   │ MethodologyTmpl   │   │
│  │ User         │   │ FrameworkVersion │   │ FirmControlObj    │   │
│  │ Client       │   │ FrameworkReq     │   │ + mappings        │   │
│  │ Invitation   │   │ CommonControl    │   │ + template items  │   │
│  │              │   │ CtrlObjLibrary   │   │                   │   │
│  │              │   │ CommonCtrlSats   │   │                   │   │
│  │              │   │ (STRM edges)     │   │                   │   │
│  └──────┬───────┘   └────────┬─────────┘   └────────┬──────────┘   │
│         │                    │                       │              │
│         ▼                    ▼                       ▼              │
│  ┌──────────────────────────────────────────────────────────────┐   │
│  │                    4. AUDIT CORE                              │   │
│  │  (Engagement lifecycle + control testing + evidence + PBC)   │   │
│  │                                                              │   │
│  │  Engagement ──► Control ──► TestProcedure                    │   │
│  │  EngagementTeamMember      EvidenceItem ──► EvidenceLink     │   │
│  │  EngagementFramework         + EvidenceItemSupports edges    │   │
│  │  ClientAcceptance          DocumentRequest                   │   │
│  │  EngagementQualityReview   ClientHubToken                    │   │
│  │  EQRFinding                DelegationToken                   │   │
│  └──────────────────────┬───────────────────────────────────────┘   │
│                         │                                           │
│                 ┌───────┴────────┐                                  │
│                 ▼                ▼                                  │
│         ┌──────────────┐ ┌──────────────┐                           │
│         │ 5. WORKPAPER │ │ 6. REPORTING │                           │
│         │    AUTHORING │ │              │                           │
│         │              │ │ Report       │                           │
│         │ Workpaper    │ │ ReportVer    │                           │
│         │ WP Version   │ │              │                           │
│         │ ReviewNote   │ │              │                           │
│         └──────────────┘ └──────────────┘                           │
│                                                                     │
│  ═══════════════════ CROSS-CUTTING ═══════════════════════════     │
│  AIDecision (all AI outputs)                                       │
│  AuditLog (append-only event trail)                                │
│  Notification (in-app + email delivery)                            │
└─────────────────────────────────────────────────────────────────────┘
```

### Why this decomposition

**Context 4 (Audit Core) is deliberately large.** The evidence chain requires ACID transactions across its entities:

```
EvidenceItem ──► EvidenceItemSupports ──► CommonControl
                                             │
                                             ├──► CommonControlSatisfies ──► FrameworkRequirement (SOC 2 CC6.1)
                                             ├──► CommonControlSatisfies ──► FrameworkRequirement (ISO 27001 A.5.15)
                                             └──► CommonControlSatisfies ──► FrameworkRequirement (HIPAA §164.312(a)(1))

EvidenceItem ──► EvidenceLink ──► TestProcedure ──► Control
                                                      │
                                                      └──► FirmControlObjective
                                                              │
                                                              └──► FirmControlObjectiveMapping ──► FrameworkRequirement
```

The `CommonControl` node is the pivot: it decouples "what a control does" (the semantic intent, platform-seeded from SCF / OSCAL / AICPA / CIS catalogs) from "how a firm tests it" (FirmControlObjective) and "how a framework expresses it" (FrameworkRequirement). STRM-encoded edges (`CommonControlSatisfies`) capture partial-satisfaction, strength, and effective dating. Evidence is supported against CommonControls directly (enabling cross-framework reuse) and linked through TestProcedures during engagement execution (preserving audit trail).

Splitting this chain across contexts would force distributed transactions or eventual consistency on operations that must be atomic (e.g., accepting a document request must atomically create an EvidenceLink and update DocumentRequest.status).

**Contexts 5 and 6 are separate** because they have distinct scaling profiles (Workpaper has real-time collaborative editing via WebSocket, Reporting is async generation) and communicate with Audit Core via plain UUID references rather than foreign keys.

### Aggregates within Audit Core

| Aggregate Root | Contains | Consistency Boundary |
|---|---|---|
| **Engagement** | EngagementTeamMember, EngagementFramework, ClientAcceptance, EngagementQualityReview, EQRFinding | Lifecycle state machine, quality gates |
| **Control** | TestProcedure | Control status depends on procedure results |
| **EvidenceItem** | EvidenceLink | Evidence exists at firm+client level, links are engagement-scoped |
| **DocumentRequest** | ClientHubToken, DelegationToken | PBC lifecycle, client interaction |

### Database Mapping

| Bounded Context | Database | Module | Isolation |
|---|---|---|---|
| 1: Firm Identity + 3: Firm Methodology | `axiom_db` | `internal/identity` | RLS (`firm_id`) |
| 2: Regulatory Framework & Common Controls + 4: Audit Core + cross-cutting | `axiom_db` | `internal/auditcore` (+ `internal/controlmapping` for CommonControl graph) | RLS (`firm_id`) on tenant tables; system reference tables are read-only |
| 5: Workpaper Authoring | `axiom_db` | `internal/workpaper` | RLS (`firm_id`) |
| 6: Reporting | `axiom_db` | `internal/reporting` | RLS (`firm_id`) |

All bounded contexts share a single database (`axiom_db`) with RLS on all tenant-scoped tables. Each module owns specific tables and accesses them via its own sqlc queries. Cross-module data is accessed through Go service interfaces, preserving clean boundaries for future service extraction.

See [Backend Architecture](backend-architecture-design.md) for module descriptions and database topology details.

---

## 2. Context 1: Firm Identity

**Purpose:** Organizational identity and people on the platform.
**Journeys:** 1 (Firm Setup), 2 (Staff Onboarding).

### Firm (Aggregate Root)

The root tenant entity. All firm-owned data across the system carries a firm_id reference.

| Attribute | Type | Description |
|---|---|---|
| id | uuid | Primary identifier |
| name | text | Display name |
| slug | text | URL-safe unique identifier |
| logo_url | text | Optional — defaults to generated initials |
| timezone | text | Auto-detected at signup, editable |
| billing_contact_email | text | For subscription and invoicing |
| subscription_tier | enum | Growth, Scale, Enterprise |
| country | text | US or Canada at launch |
| staff_count_range | text | From intake form: 1–10, 11–20, 21–40, 41–60, 60+ |
| primary_audit_types | jsonb | Multi-select from intake: SOC2, SOC1, ISO27001, ISO27701, ISO42001, HIPAA, PCI_DSS |
| settings | jsonb | General firm configuration |
| created_at | timestamptz | |

**Invariants:**
- Every tenant-scoped entity in the system references a Firm.
- A Firm always has at least one User with FirmAdmin role.

### User

People who use the platform — both firm staff and client-side users.

| Attribute | Type | Description |
|---|---|---|
| id | uuid | |
| firm_id | uuid | Nullable — null for client users |
| client_id | uuid | Nullable — set for ClientAdmin/ClientUser |
| email | text | Unique across system |
| display_name | text | |
| role | enum | FirmAdmin, Partner, Manager, Staff, EQReviewer, ClientAdmin, ClientUser, ViewOnly |
| auth_method | enum | Password, OAuth (Google/Microsoft), SAML |
| notification_frequency | enum | RealTime, Daily, Weekly — defaults vary by role |
| tour_completed | boolean | Per-user, never repeats after completion or skip |
| is_active | boolean | Soft-delete flag |
| created_at | timestamptz | |

**Invariants:**
- Firm staff (FirmAdmin, Partner, Manager, Staff, EQReviewer) must have a firm_id.
- Client users (ClientAdmin, ClientUser) must have a client_id and no firm_id.
- A User cannot hold both a firm role and a client role.
- Notification frequency defaults: Daily for Staff, RealTime for Partner and Manager.

### Client

The entity being audited.

| Attribute | Type | Description |
|---|---|---|
| id | uuid | |
| firm_id | uuid | The firm that owns this client relationship |
| name | text | |
| industry | text | |
| primary_contact_email | text | |
| created_at | timestamptz | |

**Invariants:**
- A Client belongs to exactly one Firm.
- Client users access only engagements they are explicitly invited to.

### Invitation

Magic link mechanism for onboarding staff.

| Attribute | Type | Description |
|---|---|---|
| id | uuid | |
| firm_id | uuid | |
| email | text | Invitee's email |
| assigned_role | enum | Role the invitee will receive |
| token_hash | text | Hashed magic link token |
| status | enum | Sent, Accepted, Expired |
| expires_at | timestamptz | 7 days from creation |
| reminder_sent_at | timestamptz | Day-5 reminder timestamp |
| invited_by_id | uuid | User who sent the invitation |
| accepted_at | timestamptz | |
| created_at | timestamptz | |

**Invariants:**
- Token is valid for 7 days only.
- Accepting an invitation creates a User with the assigned role.
- A day-5 reminder is sent if the link hasn't been used.

---

## 3. Context 2: Regulatory Framework & Common Controls

**Purpose:** External regulatory knowledge base and Axiom's internal common-control catalog — the cross-framework mapping core. Framework/requirement reference data is not tenant-scoped (shared read-only across all firms). CommonControls are platform-seeded but firms can adapt them via firm-specific rows.

The data model is a **control-centric directed labeled graph modeled in PostgreSQL with junction tables** (node tables + edge tables carrying NIST STRM relationship vocabulary and effective dating). No graph database. Semantic search over requirement text is powered by pgvector embeddings.

**Journeys:** 3 (Engagement Scoping), 5 (Control Testing), new cross-framework mapping and multi-framework integrated engagement journeys.

### Framework (Aggregate Root)

A named standards framework. One row per framework family: SOC 2, ISO 27001, ISO 27701, ISO 42001, HIPAA, PCI DSS, SOC 1.

| Attribute | Type | Description |
|---|---|---|
| id | uuid | |
| name | text | e.g., "SOC 2", "ISO 27001", "ISO 27701", "ISO 42001", "HIPAA", "PCI DSS", "SOC 1" |
| governing_body | text | AICPA, ISO, HHS, PCI SSC |
| description | text | |

### FrameworkVersion

A specific published version of a framework. ISO 27001:2013 and ISO 27001:2022 are separate rows; PCI DSS 4.0 and 4.0.1 are separate rows.

| Attribute | Type | Description |
|---|---|---|
| id | uuid | |
| framework_id | uuid | FK → Framework |
| version | text | e.g., "2022", "TSC 2017 (rev. 2022)", "4.0.1" |
| valid_from | date | When this version took effect |
| valid_to | date | Nullable — when superseded |
| published_at | date | |

**Invariants:**
- `(framework_id, version)` is unique.
- Version churn is expected (PCI 3.2→4.0, ISO 27001:2013→2022, NIST CSF 1.1→2.0). Mappings reference `framework_requirement_id`, so requirement-level effective dating cascades into cross-framework edges.

### FrameworkRequirement

A single criterion, control, or specification within a framework version.

| Attribute | Type | Description |
|---|---|---|
| id | uuid | |
| framework_version_id | uuid | FK → FrameworkVersion |
| requirement_code | text | e.g., "CC6.1", "A.5.15", "8.3.1", "§164.312(a)(1)" |
| title | text | |
| description | text | |
| requirement_type | enum | criterion, control, specification |
| category | text | Grouping within framework (e.g., "Access Control", "Trust Services Criterion: Security") |
| sort_order | integer | Display ordering |

**Indexes:** pgvector embedding on `description` for semantic search, maintained in `framework_requirement_embeddings`.

### CommonControl (Aggregate Root) — the pivot node

Axiom's internal common-control catalog. Every `CommonControl` represents a semantic control intent ("Access to production systems is restricted to authorized personnel") that can map to many framework requirements via STRM-encoded edges. Platform-seeded from SCF (primary), OSCAL (NIST family), AICPA official mappings, and CIS Controls v8.1. Firms may add or adapt their own.

| Attribute | Type | Description |
|---|---|---|
| id | uuid | |
| firm_id | uuid | Nullable — NULL means platform-seeded (visible to all firms); populated for firm-specific adaptations |
| code | text | Axiom-assigned identifier, e.g., "CC-AC-01" |
| name | text | |
| description | text | |
| category | text | e.g., Access Control, Encryption, Logging, Governance, Risk Management, AI Governance |
| source | enum | platform_seed, scf_import, oscal_import, aicpa_mapping, cis_mapping, firm_custom |
| created_at | timestamptz | |
| deprecated_at | timestamptz | Nullable — effective dating for catalog evolution |

**Invariants:**
- Platform-seeded rows (`firm_id IS NULL`) are read-only for firms; Axiom maintains them via quarterly catalog updates.
- Firm-custom rows are tenant-scoped and RLS-isolated.
- Never deleted — deprecated_at preserves historical mappings for archived engagements.
- pgvector embedding on `description` maintained in `common_control_embeddings` for AI-assisted mapping suggestions.

### CommonControlSatisfies (Edge) — the cross-framework mapping core

Directed, labeled, effective-dated edge from `CommonControl` to `FrameworkRequirement` carrying NIST STRM relationship vocabulary. This is the table that powers cross-framework evidence reuse.

| Attribute | Type | Description |
|---|---|---|
| id | uuid | |
| common_control_id | uuid | FK → CommonControl |
| framework_requirement_id | uuid | FK → FrameworkRequirement |
| relationship_type | enum | NIST STRM vocabulary: `equivalent-to`, `subset-of`, `superset-of`, `intersects-with`, `no-relationship` |
| strength_score | integer | 0–100, how strong the relationship is (for AI ranking and partial-coverage surfacing) |
| coverage_notes | text | Free text explaining partial satisfaction ("Covers network segmentation but not quarterly ASV scans") |
| source | enum | scf, ucf, oscal, aicpa, cis, axiom_custom |
| valid_from | date | |
| valid_to | date | Nullable — when this mapping is superseded (framework version churn) |
| created_at | timestamptz | |

**Invariants:**
- `(common_control_id, framework_requirement_id, valid_from)` is unique — a new edge row is created when a mapping changes rather than mutating the existing row.
- Partial-satisfaction (`subset-of`, `intersects-with`) must always carry `coverage_notes`.
- UI never renders a green check for `subset-of` or `intersects-with` — always surfaces percentage coverage and gap list.
- AI does NOT author authoritative crosswalks; edges ingested from SCF/OSCAL/AICPA/CIS are authoritative. AI-suggested mappings go through `AIDecision` and are confirmed by a human before becoming a `CommonControlSatisfies` row.

### ControlObjectiveLibrary (Aggregate Root)

Legacy system-maintained library of semantic control objectives. Retained for backward compatibility with the firm methodology templates; new mapping work goes through `CommonControl`. The `ControlObjectiveLibrary` is effectively a view over platform-seeded `CommonControl` rows — in the current implementation they coexist, but a future migration will fold `ControlObjectiveLibrary` into `CommonControl`.

| Attribute | Type | Description |
|---|---|---|
| id | uuid | |
| name | text | e.g., "Access to production systems is restricted to authorized personnel" |
| description | text | |
| tags | jsonb | Categorization tags |

### ControlObjectiveLibraryMapping

Maps a library objective to framework requirements across all frameworks. Being superseded by `CommonControlSatisfies`; both are maintained during the transition.

| Attribute | Type | Description |
|---|---|---|
| id | uuid | |
| library_objective_id | uuid | FK → ControlObjectiveLibrary |
| framework_requirement_id | uuid | FK → FrameworkRequirement |

**Relationship:** Many-to-many between ControlObjectiveLibrary and FrameworkRequirement.

---

## 4. Context 3: Firm Methodology

**Purpose:** A firm's reusable audit approach — templates and customized control objectives.
**Journeys:** 1 (Firm Setup — template activation), 3 (Engagement Scoping — template selection, control mapping).

### MethodologyTemplate (Aggregate Root)

Firm-level reusable templates that scaffold new engagements with controls, test procedures, and document requests.

| Attribute | Type | Description |
|---|---|---|
| id | uuid | |
| firm_id | uuid | |
| name | text | e.g., "SOC 2 Type II Standard" |
| applicable_engagement_type | text | |
| applicable_framework_id | uuid | Ref to Framework |
| version | integer | |
| is_active | boolean | |
| is_system_provided | boolean | True for pre-built; false for firm-customized (Scale tier) |
| created_at | timestamptz | |

### TemplateControl

A control definition within a methodology template.

| Attribute | Type | Description |
|---|---|---|
| id | uuid | |
| template_id | uuid | FK → MethodologyTemplate |
| firm_control_objective_id | uuid | Ref to FirmControlObjective |
| description | text | |
| is_key_control | boolean | |
| sort_order | integer | |

### TemplateTestProcedure

A test procedure definition within a template control.

| Attribute | Type | Description |
|---|---|---|
| id | uuid | |
| template_control_id | uuid | FK → TemplateControl |
| procedure_type | enum | Inquiry, Observation, InspectionOfDocument, Reperformance, Analytics |
| description | text | |
| expected_result | text | |

### TemplateDocumentRequest

Pre-drafted PBC request definitions within a methodology template. A standard SOC 2 Type II template includes 80+ pre-drafted requests.

| Attribute | Type | Description |
|---|---|---|
| id | uuid | |
| template_id | uuid | FK → MethodologyTemplate |
| template_control_id | uuid | FK → TemplateControl, nullable — linked to specific control if applicable |
| title | text | |
| instructions_template | text | Customizable instructions text |
| sort_order | integer | |

### FirmControlObjective (Aggregate Root)

A firm's customized control objectives — derived from the system library or created from scratch.

| Attribute | Type | Description |
|---|---|---|
| id | uuid | |
| firm_id | uuid | |
| source_library_id | uuid | Nullable — null means firm-created, not derived from library |
| name | text | |
| description | text | |
| custom_test_guidance | jsonb | Firm-specific testing guidance |
| created_at | timestamptz | |

### FirmControlObjectiveMapping

Maps a firm's control objectives to specific framework requirements. This is the firm-specific layer of cross-framework evidence architecture.

| Attribute | Type | Description |
|---|---|---|
| id | uuid | |
| firm_control_objective_id | uuid | FK → FirmControlObjective |
| framework_requirement_id | uuid | Ref to FrameworkRequirement |

**Cross-framework evidence chain (the full relationship):**

```
EvidenceItem
  → EvidenceLink → TestProcedure → Control
    → FirmControlObjective
      → FirmControlObjectiveMapping
        → FrameworkRequirement (SOC 2 CC6.1)
        → FrameworkRequirement (ISO 27001 A.8.3)
        → FrameworkRequirement (HIPAA §164.312(a)(1))
```

One evidence upload, one link action, and the evidence simultaneously satisfies all framework requirements mapped to that control objective. This is the architectural realization of cross-framework evidence reuse.

---

## 5. Context 4: Audit Core

**Purpose:** The heart of the platform — engagement lifecycle, control testing, evidence management, and document request fulfillment.
**Journeys:** 3 (Engagement Scoping), 5 (Control Testing), 6 (Workpaper Review), 7 (Document Requests — auditor side), 8 (Document Requests — client side), 10 (EQR).

### Engagement Aggregate

#### Engagement (Aggregate Root)

One engagement = one client, one primary framework, one audit period. The central organizing entity.

| Attribute | Type | Description |
|---|---|---|
| id | uuid | |
| firm_id | uuid | |
| client_id | uuid | FK → clients (identity module) |
| name | text | |
| engagement_type | enum | SOC1, SOC2, ISO27001, ISO27701, ISO42001, HIPAA, PCI_DSS, AgreedUponProcedures, Advisory |
| primary_framework_id | uuid | FK → Framework |
| primary_framework_version_id | uuid | FK → FrameworkVersion — locked at Planning → Fieldwork transition |
| period_start | date | For Type 1 engagements (SOC 1 Type I, SOC 2 Type I), the as-of date — period_start equals period_end. For Type 2 / continuous engagements, the start of the examination period. |
| period_end | date | For Type 1 engagements, equals period_start (single as-of date). For SOC Type 2, an examination period of 3–12 months from period_start. For ISO surveillance and PCI-style point-in-time, set per the framework's native window. |
| status | enum | Planning, Fieldwork, Review, Reporting, Finalized, Archived |
| prior_engagement_id | uuid | FK → Engagement, nullable — for rollforward |
| methodology_template_id | uuid | FK → methodology_templates (identity module) |
| report_issued_at | timestamptz | Populated when report is issued |
| assembly_deadline | date | Computed: report_issued_at + framework-specific window (SOC: 60 days; ISO CB assessment: per CB policy; PCI ROC: per QSA policy) |
| retention_deadline | date | Computed: report_issued_at + framework-specific retention (SOC/ISO: 5 yrs; HIPAA: 6 yrs; PCI: 3 yrs per PCI DSS 12.10.1) |
| finalized_at | timestamptz | |
| archived_at | timestamptz | |
| created_at | timestamptz | |

**State machine:**

```
Planning ──[ClientAcceptance signed by Partner]──► Fieldwork
Fieldwork ──[All Controls: Complete or Exception]──► Review
Review ──[All ReviewNotes resolved + EQR signed off where applicable]──► Reporting
Reporting ──[Report.status = Issued]──► Finalized
Finalized ──[System: assembly_deadline elapsed]──► Archived (IMMUTABLE)

Reverse paths (exceptional, Partner only):
  Fieldwork → Planning  (scope change requiring re-acceptance)
  Review → Fieldwork    (additional procedures needed)
  Reporting → Review    (significant issue found)

Emergency path:
  Any → Archived        (FirmAdmin: abandoned engagement)
```

**Transition guards:**

| From | To | Who | Guard |
|---|---|---|---|
| Planning | Fieldwork | Partner | ClientAcceptance.accepted_at is populated |
| Fieldwork | Review | Manager or Partner | All Controls have status Complete or Exception |
| Review | Reporting | Partner | All ReviewNotes resolved; EQR.status = Complete where applicable |
| Reporting | Finalized | Partner | Report.status = Issued |
| Finalized | Archived | System (Step Functions) | report_issued_at + assembly_window elapsed |
| Fieldwork | Planning | Partner | Scope change requiring re-acceptance |
| Review | Fieldwork | Manager or Partner | Additional procedures required |
| Reporting | Review | Partner | Significant issue identified |
| Any | Archived | FirmAdmin | Abandoned engagement |

**Invariants:**
- Framework version is locked once Fieldwork begins (Partner override with documented reason only).
- Assembly deadline and retention deadline are computed at report issuance, never manually set.
- Once Finalized, no workpaper content can be modified — addenda only.
- Once Archived, no state changes are possible.
- **Period semantics by engagement subtype:** Type 1 engagements (SOC 1 Type I, SOC 2 Type I) are point-in-time — `period_start = period_end` (a single as-of date). SOC Type 2 engagements (SOC 1 Type II, SOC 2 Type II) cover a continuous examination period of **3–12 months** (`period_end > period_start`; range validated at engagement creation). ISO certification cycles use surveillance windows; PCI assessments use the framework's annual cycle with 90-day ASV scan validity. Engagement creation UI and validation must branch on the report-type to enforce these.

#### EngagementTeamMember

Associates users to an engagement with an engagement-level role.

| Attribute | Type | Description |
|---|---|---|
| id | uuid | |
| engagement_id | uuid | FK → Engagement |
| user_id | uuid | FK → users (identity module) |
| engagement_role | text | Partner, Manager, Staff — the role on this specific engagement |
| assigned_at | timestamptz | |
| removed_at | timestamptz | Nullable |

**Invariants:**
- A user with EQReviewer role on the EngagementQualityReview for this engagement cannot also be an EngagementTeamMember.
- Used for engagement-level access control checks.

#### EngagementFramework

Supports multi-framework integrated engagements — a single engagement can simultaneously test SOC 2 + ISO 27001 + ISO 27701, reusing evidence across frameworks via the `CommonControl` graph. This is a first-class Axiom differentiator.

| Attribute | Type | Description |
|---|---|---|
| id | uuid | |
| engagement_id | uuid | FK → Engagement |
| framework_id | uuid | FK → Framework |
| framework_version_id | uuid | FK → FrameworkVersion |
| is_primary | boolean | Exactly one must be true per engagement |

**Invariants:**
- An engagement may include any combination of activated frameworks; AI-powered cross-mapping identifies where a single `CommonControl` satisfies requirements in multiple frameworks, so each control is tested once and evidence flows through `EvidenceItemSupports → CommonControl → CommonControlSatisfies` to all in-scope requirements.
- Period windows may differ across frameworks (SOC 2 Type 2 window vs. ISO surveillance vs. PCI 90-day scan validity); `EvidenceItemSupports` rows carry `period_start`/`period_end` to express per-framework coverage.

#### ClientAcceptance

Per-engagement quality risk documentation required by AICPA SQMS 1 (for SOC engagements by CPA firms) and ISO 17021-1 (for ISO certification bodies). A regulatory gate.

| Attribute | Type | Description |
|---|---|---|
| id | uuid | |
| engagement_id | uuid | FK → Engagement, unique |
| quality_risks_identified | jsonb | Structured categories + free-text narrative |
| firm_responses | jsonb | Structured responses per identified risk |
| independence_confirmed | boolean | |
| independence_confirmed_by_id | uuid | Must be Partner role |
| accepted_by_id | uuid | Must be Partner role |
| accepted_at | timestamptz | Once populated, Planning → Fieldwork is unblocked |
| created_at | timestamptz | |

**Invariants:**
- Required before Planning → Fieldwork transition.
- Only a Partner-role user can sign acceptance.
- Immutable once signed — modifications require a new version with documented reason.
- AI risk category suggestions are Tier 2 (auditor reviews and certifies, not auto-populated).

#### EngagementQualityReview

Formal engagement quality review per AICPA SQMS 2 (SOC engagements) / ISO 17021-1 impartiality review (ISO certification audits) / ISAE 3000 (Revised) equivalents.

| Attribute | Type | Description |
|---|---|---|
| id | uuid | |
| engagement_id | uuid | FK → Engagement, unique |
| reviewer_id | uuid | Ref to User — must have EQReviewer role; must NOT be an EngagementTeamMember |
| independence_documented_at | timestamptz | |
| status | enum | Assigned, InProgress, Complete |
| scope_notes | text | What was reviewed, sampling approach, time spent |
| conclusion | enum | Satisfied, SatisfiedWithConcerns, NotSatisfied |
| signed_off_at | timestamptz | |
| created_at | timestamptz | |

**Invariants:**
- Mandatory for framework assessments where the governing body or firm methodology requires independent quality review (e.g., SOC 2 per AICPA SQMS 2, ISO certifications per ISO 17021-1 impartiality rules); optional per firm policy otherwise.
- System validates: reviewer has EQReviewer role AND is not on the engagement team — assignment rejected if validation fails.
- Required before Review → Reporting transition (where applicable).
- Sign-off is immutable once recorded.
- Reviewer has read-only access to the entire engagement file.

#### EQRFinding

Individual findings within an engagement quality review.

| Attribute | Type | Description |
|---|---|---|
| id | uuid | |
| eqr_id | uuid | FK → EngagementQualityReview |
| description | text | |
| severity | enum | Observation, Recommendation, RequiredAction |
| status | enum | Pending, Addressed, Confirmed |
| team_response | text | Engagement team's response |
| responded_at | timestamptz | |
| confirmed_by_id | uuid | Reviewer confirms the response |
| confirmed_at | timestamptz | |
| created_at | timestamptz | |

**Invariants:**
- Findings with severity RequiredAction must be Addressed and Confirmed before EQR can be signed off.
- Cannot be deleted once created.

---

### Control Aggregate

#### Control (Aggregate Root)

An instantiation of a FirmControlObjective within a specific engagement.

| Attribute | Type | Description |
|---|---|---|
| id | uuid | |
| engagement_id | uuid | FK → Engagement |
| firm_control_objective_id | uuid | FK → firm_control_objectives (identity module) |
| description | text | |
| control_owner_id | uuid | Ref to User — who owns the control at the client |
| auditor_assigned_to_id | uuid | Ref to User — who is testing this control |
| status | enum | NotStarted, InProgress, Complete, Exception, NotApplicable |
| is_key_control | boolean | |
| prior_control_id | uuid | FK → Control, nullable — for rollforward |
| created_at | timestamptz | |

**Invariants:**
- Status can only advance to Complete or Exception when all child TestProcedures have a result.
- All Controls must be Complete or Exception before Fieldwork → Review transition.
- After engagement is Finalized, control conclusions are immutable.

#### TestProcedure

A specific test step within a control.

| Attribute | Type | Description |
|---|---|---|
| id | uuid | |
| control_id | uuid | FK → Control |
| procedure_type | enum | Inquiry, Observation, InspectionOfDocument, Reperformance, Analytics |
| description | text | |
| expected_result | text | |
| population_size | integer | |
| sample_size | integer | |
| sampling_method | enum | Systematic, Random, Judgmental — nullable |
| result | text | What was observed/measured/confirmed |
| exceptions_noted | text | |
| conclusion | text | |
| performed_by_id | uuid | |
| performed_at | timestamptz | |
| reviewed_by_id | uuid | |
| reviewed_at | timestamptz | |
| status | enum | NotStarted, InProgress, Complete, Exception |
| prior_procedure_id | uuid | FK → TestProcedure, nullable — for rollforward |

**Invariants:**
- Sampling documentation (population_size, sample_size, sampling_method) is required where the framework methodology calls for sample-based testing (SOC 2 Type 2, ISO surveillance); for prescriptive population-level checks (e.g., PCI DSS ASV scans, log-stream checks) populate `population_size` with total population and leave sampling fields null.
- Status progression: NotStarted → InProgress → Complete or Exception.

---

### Evidence Aggregate

#### EvidenceItem (Aggregate Root)

A single uploaded document or artifact. **Stored at the firm + client level, not the engagement level** — this enables year-over-year and cross-framework reuse without re-uploading.

| Attribute | Type | Description |
|---|---|---|
| id | uuid | |
| firm_id | uuid | |
| client_id | uuid | FK → clients (identity module) |
| filename | text | |
| storage_path | text | S3 key |
| content_type | text | MIME type |
| file_size_bytes | bigint | |
| uploaded_by_id | uuid | |
| uploaded_at | timestamptz | |
| source_type | enum | ClientUpload, CloudIntegration, APIImport, AuditorGenerated |
| source_integration | text | Nullable: Dropbox, Box, GoogleDrive, O365, Vanta |
| extracted_text | text | From Document Processing Service |
| extraction_status | enum | Pending, Complete, Failed |
| is_sensitive | boolean | PII/PHI flag |

**Invariants:**
- Evidence items are never deleted — they persist at the firm+client level across engagements.
- The same evidence item can be linked to test procedures in multiple engagements.
- For rollforward: prior year evidence is surfaced with a "used in prior year" flag.

#### EvidenceLink

Connects an EvidenceItem to a specific TestProcedure within an engagement. Preserves the engagement-scoped audit trail of which evidence supported which test step.

| Attribute | Type | Description |
|---|---|---|
| id | uuid | |
| evidence_item_id | uuid | FK → EvidenceItem |
| test_procedure_id | uuid | FK → TestProcedure |
| linked_by_id | uuid | User who created the link |
| linked_at | timestamptz | |
| notes | text | |
| ai_suggested | boolean | |
| ai_decision_id | uuid | FK → AIDecision, nullable |

**Invariants:**
- Links are frozen after engagement finalization.
- AI-suggested links require auditor accept/modify/reject action.

#### EvidenceItemSupports (Edge) — the cross-framework coverage edge

Directed, period-aware edge from `EvidenceItem` to `CommonControl`. A single evidence item (e.g., an access-review screenshot) can support many common controls, and each common control can be satisfied by many evidence items. The edge carries the period window (for SOC 2 Type 2 / ISO surveillance / PCI 90-day scan alignment), coverage percentage, and AI provenance. An `EvidenceItem` can attach to multiple `CommonControl` rows through this edge without re-uploading.

| Attribute | Type | Description |
|---|---|---|
| id | uuid | |
| firm_id | uuid | For RLS |
| evidence_item_id | uuid | FK → EvidenceItem |
| common_control_id | uuid | FK → CommonControl |
| coverage_pct | integer | 0–100, how completely this evidence covers the common control |
| period_start | date | Start of the window this evidence covers (for Type 2 / surveillance alignment) |
| period_end | date | End of the window; null means point-in-time |
| gap_notes | text | Free text describing what is NOT covered (when coverage_pct < 100) |
| ai_suggested | boolean | |
| ai_decision_id | uuid | FK → AIDecision, nullable — populated when AI proposed this mapping |
| confirmed_by_user_id | uuid | FK → users — populated once a human accepts the mapping |
| confirmed_at | timestamptz | |
| created_at | timestamptz | |

**Invariants:**
- AI-suggested rows (`ai_suggested = true AND confirmed_by_user_id IS NULL`) are proposals only. They do NOT count toward coverage until a human confirms.
- Once confirmed, the evidence simultaneously satisfies all `FrameworkRequirement` rows reachable from the `CommonControl` via `CommonControlSatisfies` edges (respecting the edge's `relationship_type` and `strength_score`).
- Period windows on this edge align evidence to framework-specific validity: PCI ASV scans are valid 90 days, pen tests 1 year, background checks 1 year, SOC 2 Type 2 requires continuous period coverage.
- Rows are not deleted — a removed mapping is expressed by setting `coverage_pct = 0` and documenting `gap_notes`, preserving the audit trail.
- UI surfaces partial coverage explicitly (never a green check for coverage_pct < 100).

---

### Document Request Aggregate

#### DocumentRequest (Aggregate Root)

A PBC (Provided By Client) request sent to the client.

| Attribute | Type | Description |
|---|---|---|
| id | uuid | |
| engagement_id | uuid | FK → Engagement |
| control_id | uuid | FK → Control, nullable |
| test_procedure_id | uuid | FK → TestProcedure, nullable |
| assigned_to_id | uuid | FK → users (identity module, client user) |
| title | text | |
| instructions | text | Detailed description of what to provide |
| due_date | date | |
| status | enum | Pending, Submitted, InReview, Accepted, Rejected, Overdue |
| reminder_count | integer | |
| last_reminder_sent_at | timestamptz | |
| fulfilled_by_evidence_item_id | uuid | FK → EvidenceItem, nullable |
| sent_at | timestamptz | When request was sent to client |
| created_at | timestamptz | |

**Invariants:**
- Automated reminders fire at: 7 days before due, on due date, 7 days after due.
- After 3 reminders, an escalation notification goes to the auditor.
- Accepting a request atomically creates an EvidenceLink to the relevant TestProcedure (when control/procedure links exist).
- For SOC 2 Type II: period coverage check (AT-C 320) is mandatory on every AI completeness review.

#### ClientHubToken

Tokenized access link for client-side document upload — no account or password required.

| Attribute | Type | Description |
|---|---|---|
| id | uuid | |
| engagement_id | uuid | FK → Engagement — scoped to this engagement only |
| client_id | uuid | |
| token_hash | text | Hashed token value |
| valid_until | timestamptz | 90 days from creation |
| created_by_id | uuid | The auditor who generated the token |
| status | enum | Active, Expired, Revoked |
| created_at | timestamptz | |

**Invariants:**
- Token provides access to document requests for this engagement only.
- After engagement is archived, the Client Hub becomes read-only.

#### DelegationToken

Single-request scoped token for ClientAdmin to delegate a specific request to a colleague.

| Attribute | Type | Description |
|---|---|---|
| id | uuid | |
| document_request_id | uuid | FK → DocumentRequest |
| delegated_by_id | uuid | Must be a ClientAdmin user |
| delegate_email | text | |
| token_hash | text | |
| custom_message | text | Note from ClientAdmin to delegate |
| valid_until | timestamptz | |
| status | enum | Active, Used, Expired |
| created_at | timestamptz | |

**Invariants:**
- Only ClientAdmin role can create delegation tokens.
- Delegate sees only: request description, instructions, and upload interface — not the full Client Hub.
- Delegation creates an AuditLog entry.

---

### Rollforward Behavior (Cross-Aggregate)

When a new engagement is created with prior_engagement_id set:

| Entity | Behavior |
|---|---|
| Engagement | New record; prior_engagement_id set; all status fields reset |
| Controls | Cloned from prior; prior_control_id set; status reset to NotStarted |
| TestProcedures | Cloned per control; status reset; editable starting point |
| DocumentRequests | Not auto-cloned; AI suggests new requests based on prior controls |
| EvidenceItems | Not touched; exist at firm+client level; surfaced with "used in prior year" flag |
| EvidenceItemSupports edges | Prior-year rows remain. On rollforward, AI proposes new `EvidenceItemSupports` rows for the new period window; prior-period rows are retained as history. |
| Workpapers | New drafts; prior_workpaper_id set; prior year visible as read-only sidebar |
| Reports | New document; prior year accessible for reference only |
| AIDecisions | Not carried forward; AI re-analyzes fresh evidence |
| ClientAcceptance | New record required; must be refreshed annually |
| EngagementQualityReview | New record required if applicable |

---

## 6. Context 5: Workpaper Authoring

**Purpose:** Collaborative document creation, version history, and review workflow. Has distinct scaling characteristics — real-time collaborative editing via WebSocket (Yjs).
**Journeys:** 5 (Control Testing — workpaper creation), 6 (Workpaper Review).

### Workpaper (Aggregate Root)

A document in the engagement file.

| Attribute | Type | Description |
|---|---|---|
| id | uuid | |
| engagement_id | uuid | Plain UUID ref |
| firm_id | uuid | For tenant isolation |
| control_id | uuid | Plain UUID ref, nullable — not all workpapers are tied to a specific control |
| workpaper_type | enum | LeadSchedule, TestPaper, Memo, ConfirmationLetter, SamplingWorksheet, ManagementLetter, AnalyticalProcedures, Other |
| title | text | |
| content | jsonb | Structured rich text |
| status | enum | Draft, PreparedPendingReview, DetailedReviewInProgress, GeneralReviewInProgress, FinalReviewInProgress, ReviewNotesOpen, ReviewComplete, SignedOff |
| prepared_by_id | uuid | The Tester (Staff or Manager performing the test procedure) |
| current_reviewer_level | enum | nullable — `DetailedReviewer`, `GeneralReviewer`, `FinalReviewer`. Tracks where the workpaper currently sits in the four-level review chain. |
| signed_off_by_id | uuid | Populated when the Final Reviewer signs off. |
| is_locked | boolean | True after engagement finalization |
| prior_workpaper_id | uuid | FK → Workpaper, nullable — for rollforward |
| created_at | timestamptz | |

**Status lifecycle (four-level sign-off):**

```
Draft
  → PreparedPendingReview        [Tester submits; system validates all AI sections human_edited = true; soft gate on modification_ratio < 0.05]
  → DetailedReviewInProgress     [Detailed Reviewer opens for review]
  → ReviewNotesOpen              [Detailed Reviewer raises notes; returns to Tester]
  → DetailedReviewInProgress     [Tester addresses notes; resubmits]
  → GeneralReviewInProgress      [Detailed Reviewer signs off; General Reviewer opens]
  → ReviewNotesOpen              [General Reviewer raises notes; returns to Tester or Detailed Reviewer]
  → GeneralReviewInProgress      [notes resolved; resubmitted]
  → FinalReviewInProgress        [General Reviewer signs off; Final Reviewer opens]
  → ReviewNotesOpen              [Final Reviewer raises notes; returns]
  → FinalReviewInProgress        [notes resolved; resubmitted]
  → ReviewComplete               [Final Reviewer clears; all notes resolved across all levels]
  → SignedOff                    [Final Reviewer records sign-off]
```

The four reviewer levels (Tester, Detailed Reviewer, General Reviewer, Final Reviewer) are workflow positions on the workpaper, not user roles — see `WorkpaperSignOff` below. Engagement Quality Review (EQR) under SQMS 2 / ISO 17021-1 §9.6 remains a separate independent track (`EngagementQualityReview`), gating Review → Reporting at the engagement level.

**Invariants:**
- Submit for review is blocked if ai_content_metadata contains any section with ai_generated = true AND human_edited = false (ISO 42001 human-in-the-loop, AICPA SSAE 21 / ISAE 3000 documentation standards, and Axiom provenance policy). Sections with modification_ratio < 0.05 trigger a confirmable warning ("Section [name] has minimal edits to AI-generated content") — soft gate, not a hard block.
- Once submitted, the workpaper is locked for the preparer — only the assigned reviewer at the current level can modify or return it.
- After engagement finalization, is_locked = true — modifications require an addendum.
- Sign-off at every reviewer level creates a timestamped, named AuditLog entry via `WorkpaperSignOff` — cannot be backdated.
- **Four-level sign-off hierarchy** enforced at the data layer (Tester → Detailed Reviewer → General Reviewer → Final Reviewer). The `current_reviewer_level` advances strictly forward; level skipping is rejected. Firm policy maps each reviewer level to eligible firm roles (e.g., Final Reviewer typically Partner; General Reviewer typically Partner or senior Manager; Detailed Reviewer typically Manager). Grounded in AICPA SQMS 1, ISO 17021-1 competence requirements, and ISAE 3000 (Revised).

### WorkpaperSignOff

Append-only sign-off ledger. One row per reviewer level signed off on a workpaper. Encodes the four-level review hierarchy as workflow attributes orthogonal to user roles.

| Attribute | Type | Description |
|---|---|---|
| id | uuid | |
| workpaper_id | uuid | FK → Workpaper |
| firm_id | uuid | For tenant isolation |
| reviewer_level | enum | `Tester`, `DetailedReviewer`, `GeneralReviewer`, `FinalReviewer` — the workflow level being signed off |
| signed_by_id | uuid | The user signing off at this level (engagement team member; their UserRole and engagement_role are recorded for audit) |
| signed_by_role_at_signoff | enum | Snapshot of signer's UserRole at sign-off time (FirmAdmin, Partner, Manager, Staff, EQReviewer) — for historical auditability when roles change |
| signed_off_at | timestamptz | System-set; cannot be backdated |
| comments | text | Optional reviewer commentary attached to the sign-off |
| superseded_at | timestamptz | Nullable — set when the workpaper is returned to a prior level (e.g., review notes raised at a higher level invalidate lower sign-offs); the prior sign-off must be re-recorded after rework |
| superseded_by_event | text | Reason the sign-off was superseded (e.g., "review_notes_raised_by_general_reviewer") |
| created_at | timestamptz | |

**Invariants:**
- One active (non-superseded) `WorkpaperSignOff` per (`workpaper_id`, `reviewer_level`) at any time.
- Sign-off ordering is enforced — `DetailedReviewer` cannot sign off before `Tester`; `GeneralReviewer` cannot sign off before `DetailedReviewer`; `FinalReviewer` cannot sign off before `GeneralReviewer`.
- The signer cannot equal the preparer (`signed_by_id ≠ workpaper.prepared_by_id`) for any non-Tester level — a tester cannot review their own work. Independence at higher levels is firm-policy driven and configurable per engagement.
- Eligibility constraints (which `UserRole` is allowed at which `reviewer_level`) are firm-policy driven; the platform records the firm's chosen mapping in firm settings and validates at sign-off time.
- When higher-level review surfaces notes and the workpaper is returned, the affected lower-level sign-offs are marked `superseded_at` (with `superseded_by_event` populated). Sign-offs must be re-recorded after rework — preserving an immutable timeline of every sign-off action and revocation.
- Each row creates an `AuditLog` entry (`workpaper.signed_off`).
- EQR is a separate track recorded in `EngagementQualityReview`, not in `WorkpaperSignOff`.

### WorkpaperVersion

Immutable version history. Every save creates a new row.

| Attribute | Type | Description |
|---|---|---|
| id | uuid | |
| workpaper_id | uuid | FK → Workpaper |
| version_number | integer | Sequential |
| content | jsonb | Structured rich text |
| saved_by_id | uuid | |
| saved_at | timestamptz | |
| is_ai_draft | boolean | Derived convenience field: true when ai_content_metadata contains any section with ai_generated = true AND human_edited = false. Retained for Axiom provenance/ISO 42001 human-in-the-loop audit trail. |
| ai_content_metadata | jsonb | Section-level AI content tracking: per-section ai_generated flag, human_edited flag, modification_ratio (Levenshtein distance / AI character count), character counts, editor identity and timestamps. See AI Architecture Design Section 5. |
| is_addendum | boolean | True for post-finalization modifications — per ISAE 3000 (Revised) / AICPA AT-C documentation rules on subsequent discoveries |
| addendum_reason | text | Required when is_addendum = true |

**Invariants:**
- ai_content_metadata tracks AI origin per section. modification_ratio is computed on save (not real-time).
- Addenda preserve the original content unchanged — they are appended records, not edits.
- Addenda require partner sign-off and a documented reason.

### ReviewNote

Structured feedback from a reviewer linked to specific workpaper content.

| Attribute | Type | Description |
|---|---|---|
| id | uuid | |
| workpaper_id | uuid | FK → Workpaper |
| content_anchor | jsonb | Reference to specific section/position in the workpaper |
| created_by_id | uuid | The reviewer (Manager or Partner) |
| description | text | |
| severity | enum | Question, Suggestion, RequiredChange |
| status | enum | Open, Responded, Resolved |
| response | text | Staff auditor's response text |
| responded_by_id | uuid | |
| responded_at | timestamptz | |
| resolved_by_id | uuid | The reviewer who resolves |
| resolved_at | timestamptz | |
| created_at | timestamptz | |

**Invariants:**
- **Cannot be deleted** — AICPA AT-C 105 / ISAE 3000 (Revised) / ISO 17021-1 documentation rules require retention of all review notes.
- Open review notes block workpaper advancement (InReview → ReviewComplete).
- Resolution workflow: reviewer creates → staff responds → reviewer resolves.
- Each note creation, response, and resolution creates an AuditLog entry.

---

## 7. Context 6: Reporting & Archival

**Purpose:** Final deliverable generation, client review, issuance, and regulatory archival. Async report generation.
**Journeys:** 9 (Reporting & Archive).

### Report (Aggregate Root)

The engagement's final deliverable.

| Attribute | Type | Description |
|---|---|---|
| id | uuid | |
| engagement_id | uuid | Plain UUID ref |
| firm_id | uuid | For tenant isolation |
| report_type | enum | SOC2Type1, SOC2Type2, SOC1Type1, SOC1Type2, ISO27001Certification, ISO27701Certification, ISO42001Certification, ISOCertificateDraft, HIPAAReport, PCIDSS_ROC, PCIDSS_AOC, PCIROCDraft, AgreedUponProcedures, ManagementLetter. **Note:** `ISOCertificateDraft` is a template-generated draft of the certification certificate produced for accredited Certification Body customers — the legal certification decision and signature remain with the CB under ISO 17021-1, and Axiom does not issue accredited certificates. `PCIROCDraft` and `PCIDSS_AOC` similarly produce QSA-deliverable templates; sign-off remains with the accredited QSA under PCI SSC. |
| status | enum | Draft, ClientReview, FirmReview, Issued, Archived |
| content | jsonb | Structured rich text |
| template_id | uuid | Nullable — firm-customizable report template used |
| generated_at | timestamptz | |
| issued_at | timestamptz | |
| issued_by_id | uuid | Must be Partner role |
| created_at | timestamptz | |

**Invariants:**
- Issuance triggers automatic computation of assembly_deadline and retention_deadline on the Engagement.
- Only a Partner can issue a report.
- AI-drafted report sections (Tier 2) follow the same ai_content_metadata gate rules as workpapers: all AI-generated sections must have human_edited = true; sections with modification_ratio < 0.05 trigger a confirmable warning. Report issuance (Report.status = Issued) validates that all AI-drafted sections have been substantively edited.
- After issuance, report transitions to read-only.
- Archived reports use S3 Object Lock (COMPLIANCE mode) with retention_deadline for WORM storage.

### ReportVersion

Immutable version history per report.

| Attribute | Type | Description |
|---|---|---|
| id | uuid | |
| report_id | uuid | FK → Report |
| version_number | integer | Sequential |
| content | jsonb | |
| saved_by_id | uuid | |
| saved_at | timestamptz | |
| is_ai_draft | boolean | Derived convenience field: true when ai_content_metadata contains any section with ai_generated = true AND human_edited = false |
| ai_content_metadata | jsonb | Section-level AI content tracking — same schema as WorkpaperVersion. See AI Architecture Design Section 5. |

---

## 8. Cross-Cutting Concerns

These operate across all bounded contexts.

### AIDecision

Every AI output that could affect audit content is recorded. Core to ISO 42001 human-in-the-loop compliance and Axiom's provenance ledger (the "auditor-defensible by construction" positioning). Required for all engagement types.

| Attribute | Type | Description |
|---|---|---|
| id | uuid | |
| firm_id | uuid | |
| engagement_id | uuid | FK → Engagement, nullable — some AI operations are firm-level |
| context_type | enum | ControlMapping, RiskCategorySuggestion, DocumentCompleteness, EvidenceLinkSuggestion, EvidenceControlMapping, GapAnalysis, FrameworkMigration, FindingsTriage, DriftRetest, ManagementResponseDraft, WorkpaperDraft, ReportSectionDraft, AnomalyDetection |
| context_id | uuid | UUID of the entity being analyzed |
| context_table | text | Which table the context_id refers to |
| model_id | text | e.g., "claude-sonnet-4-6", "claude-haiku-4-5" |
| input_token_count | integer | |
| output_token_count | integer | |
| raw_output | jsonb | Full AI response |
| suggested_value | text | The AI's recommendation |
| confidence | real | Float 0–1 |
| explanation | text | AI's reasoning for the suggestion |
| review_action | enum | Pending, Accepted, Modified, Rejected |
| accepted_value | text | What the human decided (may differ from suggested) |
| reviewed_by_id | uuid | |
| reviewed_at | timestamptz | |
| created_at | timestamptz | |

**Invariants:**
- The context_type + context_id + context_table triple provides a queryable link to the analyzed entity without polymorphic foreign keys.
- AI tiers determine human interaction requirements:
  - Tier 1 (informational): anomaly detection — no action required.
  - Tier 2 (auditor reviews): all other types — human must act.

### AuditLog

Append-only immutable event trail.

| Attribute | Type | Description |
|---|---|---|
| id | bigint | Sequential (not UUID) — for ordering guarantees |
| firm_id | uuid | |
| actor_id | uuid | |
| actor_type | enum | User, System, AIAgent |
| action | text | Namespaced string, e.g., "engagement.status.changed", "workpaper.signed_off", "evidence.linked", "review_note.created", "delegation.created", "eqr.signed_off" |
| resource_type | text | |
| resource_id | uuid | |
| old_value | jsonb | |
| new_value | jsonb | |
| ip_address | inet | |
| user_agent | text | |
| occurred_at | timestamptz | |

**Invariants:**
- **No updates, no deletes** — enforced by PostgreSQL RULE.
- Satisfies regulatory immutability requirements (AICPA AT-C 105, ISAE 3000 (Revised), ISO 17021-1 documentation rules, ISO 42001 AI system audit trail) and GDPR audit trail obligation.
- Sign-off actions cannot be backdated — occurred_at is system-set.

### Notification

In-app and email delivery system for platform events.

| Attribute | Type | Description |
|---|---|---|
| id | uuid | |
| recipient_id | uuid | |
| notification_type | enum | EngagementAssignment, ReviewNoteAdded, ReviewNoteResolved, DocumentRequestStatus, PhaseTransition, EQRNotification, ReminderEscalation, ArchivalConfirmation, RetentionWarning |
| title | text | |
| body | text | |
| deep_link | text | URL path to specific content |
| is_read | boolean | |
| delivery_channel | enum | InApp, Email, Both — governed by recipient's notification_frequency |
| created_at | timestamptz | |

---

## 9. Data Model

The data model translates the domain model above into PostgreSQL tables in a single database (`axiom_db`). This section covers physical design decisions that go beyond the domain model: column types, constraints, indexes, enum definitions, and cross-module reference strategy.

### Database Topology

```
RDS PostgreSQL Instance (Multi-AZ production, Single-AZ dev/staging)
└── axiom_db
    ├── Identity module tables     → Contexts 1 + 3 (Firm Identity, Firm Methodology)
    ├── Audit Core module tables   → Contexts 2 + 4 + cross-cutting (Framework & Common Controls, Audit Core, AI/Audit/Notification)
    ├── Workpaper module tables    → Context 5 (Workpaper Authoring)
    └── Reporting module tables    → Context 6 (Reporting)
```

All tables share one database. The Axiom API connects with a single database user (`axiom_svc`). RLS is enabled on all tenant-scoped tables.

### Cross-Module Reference Strategy

All cross-module references use standard foreign key constraints with full referential integrity — a significant advantage of the single-database modular monolith over the previous multi-database design. Each module accesses other modules' data through Go service interfaces, not direct SQL queries against other modules' tables. This preserves clean module boundaries for future service extraction.

| Pattern | Usage |
|---|---|
| FK constraint | All entity references, including cross-module (e.g., `engagement_id` in workpaper tables → `engagements` in auditcore tables) |
| Go service interface | Cross-module data access at the application layer (e.g., workpaper module calls auditcore module to validate engagement access) |

### Enum Types

**Identity module enums:**

```sql
CREATE TYPE user_role AS ENUM (
  'FirmAdmin','Partner','Manager','Staff','EQReviewer',
  'ClientAdmin','ClientUser','ViewOnly');
CREATE TYPE auth_method AS ENUM ('Password','OAuth','SAML');
CREATE TYPE notification_frequency AS ENUM ('RealTime','Daily','Weekly');
CREATE TYPE invitation_status AS ENUM ('Sent','Accepted','Expired');
```

**Audit Core module enums:**

```sql
CREATE TYPE engagement_type AS ENUM (
  'SOC1','SOC2','ISO27001','ISO27701','ISO42001',
  'HIPAA','PCI_DSS','AgreedUponProcedures','Advisory');
CREATE TYPE engagement_status AS ENUM (
  'Planning','Fieldwork','Review','Reporting','Finalized','Archived');
CREATE TYPE control_status AS ENUM (
  'NotStarted','InProgress','Complete','Exception','NotApplicable');
CREATE TYPE procedure_type AS ENUM (
  'Inquiry','Observation','InspectionOfDocument','Reperformance','Analytics');
CREATE TYPE procedure_status AS ENUM ('NotStarted','InProgress','Complete','Exception');
CREATE TYPE sampling_method AS ENUM ('Systematic','Random','Judgmental');
CREATE TYPE evidence_source AS ENUM (
  'ClientUpload','CloudIntegration','APIImport','AuditorGenerated','AgenticCollected');
CREATE TYPE extraction_status AS ENUM ('Pending','Complete','Failed');
CREATE TYPE doc_request_status AS ENUM (
  'Pending','Submitted','InReview','Accepted','Rejected','Overdue');
CREATE TYPE hub_token_status AS ENUM ('Active','Expired','Revoked');
CREATE TYPE delegation_token_status AS ENUM ('Active','Used','Expired');
CREATE TYPE eqr_status AS ENUM ('Assigned','InProgress','Complete');
CREATE TYPE eqr_conclusion AS ENUM ('Satisfied','SatisfiedWithConcerns','NotSatisfied');
CREATE TYPE finding_severity AS ENUM ('Observation','Recommendation','RequiredAction');
CREATE TYPE finding_status AS ENUM ('Pending','Addressed','Confirmed');
CREATE TYPE requirement_type AS ENUM ('criterion','control','specification');
CREATE TYPE strm_relationship_type AS ENUM (
  'equivalent-to','subset-of','superset-of','intersects-with','no-relationship');
CREATE TYPE common_control_source AS ENUM (
  'platform_seed','scf_import','oscal_import','aicpa_mapping','cis_mapping','firm_custom');
CREATE TYPE mapping_source AS ENUM (
  'scf','ucf','oscal','aicpa','cis','axiom_custom');
CREATE TYPE ai_context_type AS ENUM (
  'ControlMapping','RiskCategorySuggestion',
  'DocumentCompleteness','EvidenceLinkSuggestion',
  'EvidenceControlMapping','GapAnalysis','FrameworkMigration',
  'FindingsTriage','DriftRetest','ManagementResponseDraft',
  'WorkpaperDraft','ReportSectionDraft','AnomalyDetection');
CREATE TYPE ai_review_action AS ENUM ('Pending','Accepted','Modified','Rejected');
CREATE TYPE actor_type AS ENUM ('User','System','AIAgent');
CREATE TYPE notification_type AS ENUM (
  'EngagementAssignment','ReviewNoteAdded','ReviewNoteResolved',
  'DocumentRequestStatus','PhaseTransition','EQRNotification',
  'ReminderEscalation','ArchivalConfirmation','RetentionWarning',
  'DriftDetected','EvidenceExpiring');
CREATE TYPE delivery_channel AS ENUM ('InApp','Email','Both');
```

**Workpaper module enums:**

```sql
CREATE TYPE workpaper_type AS ENUM (
  'LeadSchedule','TestPaper','Memo','ConfirmationLetter',
  'SamplingWorksheet','ManagementLetter','AnalyticalProcedures','Other');
CREATE TYPE workpaper_status AS ENUM (
  'Draft','PreparedPendingReview',
  'DetailedReviewInProgress','GeneralReviewInProgress','FinalReviewInProgress',
  'ReviewNotesOpen','ReviewComplete','SignedOff');
CREATE TYPE reviewer_level AS ENUM (
  'Tester','DetailedReviewer','GeneralReviewer','FinalReviewer');
CREATE TYPE review_note_severity AS ENUM ('Question','Suggestion','RequiredChange');
CREATE TYPE review_note_status AS ENUM ('Open','Responded','Resolved');
```

**Reporting module enums:**

```sql
CREATE TYPE report_type AS ENUM (
  'SOC2Type1','SOC2Type2','SOC1Type1','SOC1Type2',
  'ISO27001Certification','ISO27701Certification','ISO42001Certification',
  'ISOCertificateDraft',
  'HIPAAReport','PCIDSS_ROC','PCIDSS_AOC','PCIROCDraft',
  'AgreedUponProcedures','ManagementLetter');
-- ISOCertificateDraft: template-generated certificate for accredited Certification Body
-- customers; legal certification decision and signature remain with the CB under
-- ISO 17021-1.
-- PCIROCDraft / PCIDSS_AOC: template-generated ROC and Attestation of Compliance for
-- QSA firm customers; legal sign-off remains with the accredited QSA under PCI SSC.
CREATE TYPE report_status AS ENUM ('Draft','ClientReview','FirmReview','Issued','Archived');
```

### Table Definitions by Module

All tables reside in `axiom_db`. All UUID primary keys use `DEFAULT gen_random_uuid()`. All timestamps use `timestamptz` with `DEFAULT now()` where applicable.

#### Identity Module

**firms** — `id (uuid PK)`, `name (text NOT NULL)`, `slug (text NOT NULL UNIQUE)`, `logo_url (text)`, `timezone (text NOT NULL DEFAULT 'America/New_York')`, `billing_contact_email (text NOT NULL)`, `subscription_tier (text NOT NULL CHECK IN ('Growth','Scale','Enterprise'))`, `country (text NOT NULL CHECK IN ('US','CA'))`, `staff_count_range (text)`, `primary_audit_types (jsonb)`, `settings (jsonb NOT NULL DEFAULT '{}')`, `created_at (timestamptz NOT NULL)`.

**users** — `id (uuid PK)`, `firm_id (uuid FK → firms)`, `client_id (uuid)`, `email (text NOT NULL UNIQUE)`, `display_name (text NOT NULL)`, `role (user_role NOT NULL)`, `auth_method (auth_method NOT NULL DEFAULT 'Password')`, `password_hash (text)`, `notification_frequency (notification_frequency NOT NULL)`, `tour_completed (boolean NOT NULL DEFAULT false)`, `is_active (boolean NOT NULL DEFAULT true)`, `created_at (timestamptz NOT NULL)`. **Check constraints:** user belongs to firm XOR client; client roles require client_id. **Indexes:** `(firm_id)`, `(client_id)`, `(email)`.

**clients** — `id (uuid PK)`, `firm_id (uuid FK → firms NOT NULL)`, `name (text NOT NULL)`, `industry (text)`, `primary_contact_email (text)`, `created_at (timestamptz NOT NULL)`. **Indexes:** `(firm_id)`.

**invitations** — `id (uuid PK)`, `firm_id (uuid FK → firms NOT NULL)`, `email (text NOT NULL)`, `assigned_role (user_role NOT NULL)`, `token_hash (text NOT NULL UNIQUE)`, `status (invitation_status NOT NULL DEFAULT 'Sent')`, `expires_at (timestamptz NOT NULL)`, `reminder_sent_at (timestamptz)`, `invited_by_id (uuid FK → users NOT NULL)`, `accepted_at (timestamptz)`, `created_at (timestamptz NOT NULL)`. **Indexes:** `(firm_id)`, `(token_hash)`.

**methodology_templates** — `id (uuid PK)`, `firm_id (uuid FK → firms NOT NULL)`, `name (text NOT NULL)`, `applicable_engagement_type (text NOT NULL)`, `applicable_framework_id (uuid)`, `version (integer NOT NULL DEFAULT 1)`, `is_active (boolean NOT NULL DEFAULT true)`, `is_system_provided (boolean NOT NULL DEFAULT false)`, `created_at (timestamptz NOT NULL)`. **Indexes:** `(firm_id, applicable_engagement_type)`.

**template_controls** — `id (uuid PK)`, `template_id (uuid FK → methodology_templates NOT NULL)`, `firm_control_objective_id (uuid)`, `description (text NOT NULL)`, `is_key_control (boolean NOT NULL DEFAULT false)`, `sort_order (integer NOT NULL)`. **Indexes:** `(template_id, sort_order)`.

**template_test_procedures** — `id (uuid PK)`, `template_control_id (uuid FK → template_controls NOT NULL)`, `procedure_type (text NOT NULL)`, `description (text NOT NULL)`, `expected_result (text)`.

**template_document_requests** — `id (uuid PK)`, `template_id (uuid FK → methodology_templates NOT NULL)`, `template_control_id (uuid FK → template_controls)`, `title (text NOT NULL)`, `instructions_template (text NOT NULL)`, `sort_order (integer NOT NULL)`.

**firm_control_objectives** — `id (uuid PK)`, `firm_id (uuid FK → firms NOT NULL)`, `source_library_id (uuid)`, `name (text NOT NULL)`, `description (text NOT NULL)`, `custom_test_guidance (jsonb)`, `created_at (timestamptz NOT NULL)`. **Indexes:** `(firm_id)`.

**firm_control_objective_mappings** — `id (uuid PK)`, `firm_control_objective_id (uuid FK → firm_control_objectives NOT NULL)`, `framework_requirement_id (uuid NOT NULL)`. **Unique:** `(firm_control_objective_id, framework_requirement_id)`. **Indexes:** `(firm_control_objective_id)`, `(framework_requirement_id)`.

**firm_control_objective_embeddings** — `id (uuid PK)`, `firm_control_objective_id (uuid FK → firm_control_objectives NOT NULL UNIQUE)`, `embedding (vector(1024) NOT NULL)`, `created_at (timestamptz NOT NULL)`. **Indexes:** `USING ivfflat (embedding vector_cosine_ops)`. Requires pgvector extension on `axiom_db`. Used by Feature 2 (Control Mapping) for semantic similarity retrieval.

#### Audit Core Module — System Reference Tables (no RLS, not tenant-scoped)

**frameworks** — `id (uuid PK)`, `name (text NOT NULL UNIQUE)`, `governing_body (text NOT NULL)`, `description (text)`. Rows: SOC 2, SOC 1, ISO 27001, ISO 27701, ISO 42001, HIPAA, PCI DSS.

**framework_versions** — `id (uuid PK)`, `framework_id (uuid FK → frameworks NOT NULL)`, `version (text NOT NULL)`, `valid_from (date NOT NULL)`, `valid_to (date)`, `published_at (date)`. **Unique:** `(framework_id, version)`. **Indexes:** `(framework_id, valid_from)`.

**framework_requirements** — `id (uuid PK)`, `framework_version_id (uuid FK → framework_versions NOT NULL)`, `requirement_code (text NOT NULL)`, `title (text NOT NULL)`, `description (text)`, `requirement_type (requirement_type NOT NULL)`, `category (text)`, `sort_order (integer NOT NULL)`. **Unique:** `(framework_version_id, requirement_code)`. **Indexes:** `(framework_version_id, sort_order)`.

**common_controls** — `id (uuid PK)`, `firm_id (uuid)` — nullable; NULL = platform-seeded, populated = firm-custom, `code (text NOT NULL)`, `name (text NOT NULL)`, `description (text NOT NULL)`, `category (text NOT NULL)`, `source (common_control_source NOT NULL)`, `created_at (timestamptz NOT NULL)`, `deprecated_at (timestamptz)`. **Unique:** `(firm_id, code)` (with `firm_id IS NULL` treated as a distinct tenant for uniqueness). **Indexes:** `(firm_id)`, `(category)`, `(source)`. **RLS:** partial — rows with `firm_id IS NULL` are visible to all; `firm_id IS NOT NULL` rows are RLS-isolated to that firm.

**common_control_satisfies** — `id (uuid PK)`, `common_control_id (uuid FK → common_controls NOT NULL)`, `framework_requirement_id (uuid FK → framework_requirements NOT NULL)`, `relationship_type (strm_relationship_type NOT NULL)`, `strength_score (integer NOT NULL CHECK (strength_score BETWEEN 0 AND 100))`, `coverage_notes (text)`, `source (mapping_source NOT NULL)`, `valid_from (date NOT NULL)`, `valid_to (date)`, `created_at (timestamptz NOT NULL)`. **Check:** `(relationship_type IN ('subset-of','intersects-with') AND coverage_notes IS NOT NULL) OR relationship_type NOT IN ('subset-of','intersects-with')`. **Unique:** `(common_control_id, framework_requirement_id, valid_from)`. **Indexes:** `(common_control_id)`, `(framework_requirement_id)`, `(valid_from, valid_to)`.

**common_control_embeddings** — `id (uuid PK)`, `common_control_id (uuid FK → common_controls NOT NULL UNIQUE)`, `embedding (vector(1024) NOT NULL)`, `created_at (timestamptz NOT NULL)`. **Indexes:** `USING ivfflat (embedding vector_cosine_ops)`. Used by Feature: evidence→control mapping and cross-framework semantic search.

**control_objective_library** — `id (uuid PK)`, `name (text NOT NULL)`, `description (text NOT NULL)`, `tags (jsonb NOT NULL DEFAULT '[]')`. Legacy table — superseded by `common_controls`; retained during migration for template compatibility.

**control_objective_library_mappings** — `id (uuid PK)`, `library_objective_id (uuid FK → control_objective_library NOT NULL)`, `framework_requirement_id (uuid FK → framework_requirements NOT NULL)`. **Unique:** `(library_objective_id, framework_requirement_id)`.

#### Audit Core Module — Tenant-Scoped Tables (RLS via firm_id)

All tables below carry `firm_id uuid NOT NULL` with an index. RLS policy: `USING (firm_id = current_setting('app.current_firm_id')::uuid)`.

**engagements** — `id (uuid PK)`, `firm_id`, `client_id (uuid NOT NULL)`, `name (text NOT NULL)`, `engagement_type (engagement_type NOT NULL)`, `primary_framework_id (uuid FK → frameworks NOT NULL)`, `primary_framework_version_id (uuid FK → framework_versions NOT NULL)`, `period_start (date NOT NULL)`, `period_end (date NOT NULL)`, `status (engagement_status NOT NULL DEFAULT 'Planning')`, `prior_engagement_id (uuid FK → engagements)`, `methodology_template_id (uuid)`, `report_issued_at (timestamptz)`, `assembly_deadline (date)`, `retention_deadline (date)`, `finalized_at (timestamptz)`, `archived_at (timestamptz)`, `created_at (timestamptz NOT NULL)`. **Indexes:** `(firm_id)`, `(client_id)`, `(status)`, `(prior_engagement_id)`.

**engagement_team_members** — `id (uuid PK)`, `firm_id`, `engagement_id (uuid FK → engagements NOT NULL)`, `user_id (uuid NOT NULL)`, `engagement_role (text NOT NULL)`, `assigned_at (timestamptz NOT NULL)`, `removed_at (timestamptz)`. **Unique:** `(engagement_id, user_id) WHERE removed_at IS NULL`. **Indexes:** `(engagement_id, user_id)`.

**engagement_frameworks** — `id (uuid PK)`, `firm_id`, `engagement_id (uuid FK → engagements NOT NULL)`, `framework_id (uuid FK → frameworks NOT NULL)`, `framework_version_id (uuid FK → framework_versions NOT NULL)`, `is_primary (boolean NOT NULL DEFAULT false)`. **Unique:** `(engagement_id, framework_id)`. Enables multi-framework integrated engagements (SOC 2 + ISO 27001 + ISO 27701 in a single engagement).

**client_acceptances** — `id (uuid PK)`, `firm_id`, `engagement_id (uuid FK → engagements NOT NULL UNIQUE)`, `quality_risks_identified (jsonb NOT NULL DEFAULT '[]')`, `firm_responses (jsonb NOT NULL DEFAULT '[]')`, `independence_confirmed (boolean NOT NULL DEFAULT false)`, `independence_confirmed_by_id (uuid)`, `accepted_by_id (uuid)`, `accepted_at (timestamptz)`, `created_at (timestamptz NOT NULL)`.

**engagement_quality_reviews** — `id (uuid PK)`, `firm_id`, `engagement_id (uuid FK → engagements NOT NULL UNIQUE)`, `reviewer_id (uuid NOT NULL)`, `independence_documented_at (timestamptz)`, `status (eqr_status NOT NULL DEFAULT 'Assigned')`, `scope_notes (text)`, `conclusion (eqr_conclusion)`, `signed_off_at (timestamptz)`, `created_at (timestamptz NOT NULL)`.

**eqr_findings** — `id (uuid PK)`, `firm_id`, `eqr_id (uuid FK → engagement_quality_reviews NOT NULL)`, `description (text NOT NULL)`, `severity (finding_severity NOT NULL)`, `status (finding_status NOT NULL DEFAULT 'Pending')`, `team_response (text)`, `responded_at (timestamptz)`, `confirmed_by_id (uuid)`, `confirmed_at (timestamptz)`, `created_at (timestamptz NOT NULL)`.

**controls** — `id (uuid PK)`, `firm_id`, `engagement_id (uuid FK → engagements NOT NULL)`, `firm_control_objective_id (uuid NOT NULL)`, `description (text NOT NULL)`, `control_owner_id (uuid)`, `auditor_assigned_to_id (uuid)`, `status (control_status NOT NULL DEFAULT 'NotStarted')`, `is_key_control (boolean NOT NULL DEFAULT false)`, `prior_control_id (uuid FK → controls)`, `created_at (timestamptz NOT NULL)`. **Indexes:** `(firm_id)`, `(engagement_id)`, `(auditor_assigned_to_id)`, `(engagement_id, status)`.

**test_procedures** — `id (uuid PK)`, `firm_id`, `control_id (uuid FK → controls NOT NULL)`, `procedure_type (procedure_type NOT NULL)`, `description (text NOT NULL)`, `expected_result (text)`, `population_size (integer)`, `sample_size (integer)`, `sampling_method (sampling_method)`, `result (text)`, `exceptions_noted (text)`, `conclusion (text)`, `performed_by_id (uuid)`, `performed_at (timestamptz)`, `reviewed_by_id (uuid)`, `reviewed_at (timestamptz)`, `status (procedure_status NOT NULL DEFAULT 'NotStarted')`, `prior_procedure_id (uuid FK → test_procedures)`. **Indexes:** `(control_id)`, `(performed_by_id)`.

**evidence_items** — `id (uuid PK)`, `firm_id`, `client_id (uuid NOT NULL)`, `filename (text NOT NULL)`, `storage_path (text NOT NULL)`, `content_type (text NOT NULL)`, `file_size_bytes (bigint NOT NULL)`, `uploaded_by_id (uuid NOT NULL)`, `uploaded_at (timestamptz NOT NULL)`, `source_type (evidence_source NOT NULL)`, `source_integration (text)`, `extracted_text (text)`, `extraction_status (extraction_status NOT NULL DEFAULT 'Pending')`, `is_sensitive (boolean NOT NULL DEFAULT false)`. **Indexes:** `(firm_id, client_id)`.

**evidence_embeddings** — `id (uuid PK)`, `evidence_item_id (uuid FK → evidence_items NOT NULL UNIQUE)`, `embedding (vector(1024) NOT NULL)`, `created_at (timestamptz NOT NULL)`. **Indexes:** `USING ivfflat (embedding vector_cosine_ops)`.

**framework_requirement_embeddings** — `id (uuid PK)`, `framework_requirement_id (uuid FK → framework_requirements NOT NULL UNIQUE)`, `embedding (vector(1024) NOT NULL)`, `created_at (timestamptz NOT NULL)`. **Indexes:** `USING ivfflat (embedding vector_cosine_ops)`. Used by Feature 2 (Control Mapping) for semantic similarity matching.

**control_objective_library_embeddings** — `id (uuid PK)`, `library_objective_id (uuid FK → control_objective_library NOT NULL UNIQUE)`, `embedding (vector(1024) NOT NULL)`, `created_at (timestamptz NOT NULL)`. **Indexes:** `USING ivfflat (embedding vector_cosine_ops)`. Used by Feature 2 (Control Mapping) as few-shot examples for cross-framework mapping.

**evidence_links** — `id (uuid PK)`, `firm_id`, `evidence_item_id (uuid FK → evidence_items NOT NULL)`, `test_procedure_id (uuid FK → test_procedures NOT NULL)`, `linked_by_id (uuid NOT NULL)`, `linked_at (timestamptz NOT NULL)`, `notes (text)`, `ai_suggested (boolean NOT NULL DEFAULT false)`, `ai_decision_id (uuid FK → ai_decisions)`. **Unique:** `(evidence_item_id, test_procedure_id)`. **Indexes:** `(test_procedure_id)`, `(evidence_item_id)`.

**evidence_item_supports** — `id (uuid PK)`, `firm_id`, `evidence_item_id (uuid FK → evidence_items NOT NULL)`, `common_control_id (uuid FK → common_controls NOT NULL)`, `coverage_pct (integer NOT NULL CHECK (coverage_pct BETWEEN 0 AND 100))`, `period_start (date NOT NULL)`, `period_end (date)`, `gap_notes (text)`, `ai_suggested (boolean NOT NULL DEFAULT false)`, `ai_decision_id (uuid FK → ai_decisions)`, `confirmed_by_user_id (uuid FK → users)`, `confirmed_at (timestamptz)`, `created_at (timestamptz NOT NULL)`. **Check:** `(coverage_pct = 100 OR gap_notes IS NOT NULL)`. **Indexes:** `(evidence_item_id)`, `(common_control_id)`, `(firm_id, common_control_id)`, `(period_start, period_end)`. RLS policy: `USING (firm_id = current_firm_id())`. The cross-framework reuse edge — coverage flows from here through `common_control_satisfies` to every mapped `framework_requirement`.

**document_requests** — `id (uuid PK)`, `firm_id`, `engagement_id (uuid FK → engagements NOT NULL)`, `control_id (uuid FK → controls)`, `test_procedure_id (uuid FK → test_procedures)`, `assigned_to_id (uuid)`, `title (text NOT NULL)`, `instructions (text NOT NULL)`, `due_date (date)`, `status (doc_request_status NOT NULL DEFAULT 'Pending')`, `reminder_count (integer NOT NULL DEFAULT 0)`, `last_reminder_sent_at (timestamptz)`, `fulfilled_by_evidence_item_id (uuid FK → evidence_items)`, `sent_at (timestamptz)`, `created_at (timestamptz NOT NULL)`. **Indexes:** `(engagement_id)`, `(engagement_id, status)`, `(assigned_to_id)`.

**client_hub_tokens** — `id (uuid PK)`, `firm_id`, `engagement_id (uuid FK → engagements NOT NULL)`, `client_id (uuid NOT NULL)`, `token_hash (text NOT NULL UNIQUE)`, `valid_until (timestamptz NOT NULL)`, `created_by_id (uuid NOT NULL)`, `status (hub_token_status NOT NULL DEFAULT 'Active')`, `created_at (timestamptz NOT NULL)`. **Indexes:** `(token_hash)`.

**delegation_tokens** — `id (uuid PK)`, `firm_id`, `document_request_id (uuid FK → document_requests NOT NULL)`, `delegated_by_id (uuid NOT NULL)`, `delegate_email (text NOT NULL)`, `token_hash (text NOT NULL UNIQUE)`, `custom_message (text)`, `valid_until (timestamptz NOT NULL)`, `status (delegation_token_status NOT NULL DEFAULT 'Active')`, `created_at (timestamptz NOT NULL)`. **Indexes:** `(token_hash)`, `(document_request_id)`.

**ai_decisions** — `id (uuid PK)`, `firm_id`, `engagement_id (uuid FK → engagements)`, `context_type (ai_context_type NOT NULL)`, `context_id (uuid NOT NULL)`, `context_table (text NOT NULL)`, `model_id (text NOT NULL)`, `input_token_count (integer)`, `output_token_count (integer)`, `raw_output (jsonb NOT NULL)`, `suggested_value (text)`, `confidence (real CHECK (confidence >= 0 AND confidence <= 1))`, `explanation (text)`, `review_action (ai_review_action NOT NULL DEFAULT 'Pending')`, `accepted_value (text)`, `reviewed_by_id (uuid)`, `reviewed_at (timestamptz)`, `created_at (timestamptz NOT NULL)`. **Indexes:** `(firm_id)`, `(engagement_id)`, `(context_type, context_id)`.

**audit_log** — `id (bigint PK GENERATED ALWAYS AS IDENTITY)`, `firm_id (uuid NOT NULL)`, `actor_id (uuid)`, `actor_type (actor_type NOT NULL)`, `action (text NOT NULL)`, `resource_type (text NOT NULL)`, `resource_id (uuid)`, `old_value (jsonb)`, `new_value (jsonb)`, `ip_address (inet)`, `user_agent (text)`, `occurred_at (timestamptz NOT NULL DEFAULT now())`. **Immutability:** `CREATE RULE audit_log_no_update AS ON UPDATE TO audit_log DO INSTEAD NOTHING; CREATE RULE audit_log_no_delete AS ON DELETE TO audit_log DO INSTEAD NOTHING;`. **Indexes:** `(firm_id, occurred_at)`, `(resource_type, resource_id)`, `(actor_id)`.

**notifications** — `id (uuid PK)`, `firm_id (uuid NOT NULL)`, `recipient_id (uuid NOT NULL)`, `notification_type (notification_type NOT NULL)`, `title (text NOT NULL)`, `body (text)`, `deep_link (text)`, `is_read (boolean NOT NULL DEFAULT false)`, `delivery_channel (delivery_channel NOT NULL)`, `created_at (timestamptz NOT NULL)`. **Indexes:** `(recipient_id, is_read, created_at DESC)`.

#### Workpaper Module

Application-layer isolation.

**workpapers** — `id (uuid PK)`, `engagement_id (uuid NOT NULL)`, `firm_id (uuid NOT NULL)`, `control_id (uuid)`, `workpaper_type (workpaper_type NOT NULL)`, `title (text NOT NULL)`, `content (jsonb NOT NULL DEFAULT '{}')`, `status (workpaper_status NOT NULL DEFAULT 'Draft')`, `prepared_by_id (uuid)`, `current_reviewer_level (reviewer_level)`, `signed_off_by_id (uuid)`, `is_locked (boolean NOT NULL DEFAULT false)`, `prior_workpaper_id (uuid FK → workpapers)`, `created_at (timestamptz NOT NULL)`. **Indexes:** `(firm_id)`, `(engagement_id)`, `(engagement_id, status)`. **Note:** `signed_off_by_id` is populated when the `FinalReviewer` records sign-off; per-level sign-offs are recorded in `workpaper_sign_offs`.

**workpaper_sign_offs** — `id (uuid PK)`, `firm_id (uuid NOT NULL)`, `workpaper_id (uuid FK → workpapers NOT NULL)`, `reviewer_level (reviewer_level NOT NULL)`, `signed_by_id (uuid NOT NULL)`, `signed_by_role_at_signoff (user_role NOT NULL)`, `signed_off_at (timestamptz NOT NULL DEFAULT now())`, `comments (text)`, `superseded_at (timestamptz)`, `superseded_by_event (text)`, `created_at (timestamptz NOT NULL DEFAULT now())`. **Unique:** `(workpaper_id, reviewer_level) WHERE superseded_at IS NULL` — at most one active sign-off per workpaper per reviewer level. **Indexes:** `(workpaper_id)`, `(firm_id, signed_by_id)`. **Immutability:** `CREATE RULE workpaper_sign_offs_no_update AS ON UPDATE TO workpaper_sign_offs DO INSTEAD NOTHING;` — sign-offs cannot be edited; supersession is achieved by updating `superseded_at` only via a designated function. `CREATE RULE workpaper_sign_offs_no_delete AS ON DELETE TO workpaper_sign_offs DO INSTEAD NOTHING;` — append-only ledger. **Application-layer constraints:** sign-off ordering enforced (Tester → DetailedReviewer → GeneralReviewer → FinalReviewer); signer cannot equal preparer at non-Tester levels (independence); firm-policy maps each `reviewer_level` to eligible `user_role` values, validated at sign-off time.

**workpaper_versions** — `id (uuid PK)`, `workpaper_id (uuid FK → workpapers NOT NULL)`, `version_number (integer NOT NULL)`, `content (jsonb NOT NULL)`, `saved_by_id (uuid NOT NULL)`, `saved_at (timestamptz NOT NULL)`, `is_ai_draft (boolean NOT NULL DEFAULT false)`, `ai_content_metadata (jsonb)`, `is_addendum (boolean NOT NULL DEFAULT false)`, `addendum_reason (text)`. **Check:** `(is_addendum = false OR addendum_reason IS NOT NULL)`. **Unique:** `(workpaper_id, version_number)`. **Note:** `is_ai_draft` is a derived convenience field — true when `ai_content_metadata` contains any section with `ai_generated = true AND human_edited = false`.

**review_notes** — `id (uuid PK)`, `firm_id (uuid NOT NULL)`, `workpaper_id (uuid FK → workpapers NOT NULL)`, `raised_at_level (reviewer_level NOT NULL)`, `content_anchor (jsonb)`, `created_by_id (uuid NOT NULL)`, `description (text NOT NULL)`, `severity (review_note_severity NOT NULL)`, `status (review_note_status NOT NULL DEFAULT 'Open')`, `response (text)`, `responded_by_id (uuid)`, `responded_at (timestamptz)`, `resolved_by_id (uuid)`, `resolved_at (timestamptz)`, `created_at (timestamptz NOT NULL)`. **Immutability:** `CREATE RULE review_notes_no_delete AS ON DELETE TO review_notes DO INSTEAD NOTHING;`. **Indexes:** `(workpaper_id, status)`. **Note:** `raised_at_level` records which review level raised the note — informs which prior `WorkpaperSignOff` rows must be superseded when notes return the workpaper to the preparer.

#### Reporting Module

Application-layer isolation.

**reports** — `id (uuid PK)`, `engagement_id (uuid NOT NULL)`, `firm_id (uuid NOT NULL)`, `report_type (report_type NOT NULL)`, `status (report_status NOT NULL DEFAULT 'Draft')`, `content (jsonb NOT NULL DEFAULT '{}')`, `template_id (uuid)`, `generated_at (timestamptz NOT NULL)`, `issued_at (timestamptz)`, `issued_by_id (uuid)`. **Indexes:** `(firm_id)`, `(engagement_id)`.

**report_versions** — `id (uuid PK)`, `report_id (uuid FK → reports NOT NULL)`, `version_number (integer NOT NULL)`, `content (jsonb NOT NULL)`, `saved_by_id (uuid NOT NULL)`, `saved_at (timestamptz NOT NULL)`, `is_ai_draft (boolean NOT NULL DEFAULT false)`, `ai_content_metadata (jsonb)`. **Unique:** `(report_id, version_number)`. **Note:** `is_ai_draft` is a derived convenience field — true when `ai_content_metadata` contains any section with `ai_generated = true AND human_edited = false`.

---

## 10. Multi-Tenancy and Isolation

All tables reside in `axiom_db`. RLS is enabled on all tenant-scoped tables.

| Module | Tables | firm_id Indexed | RLS |
|---|---|---|---|
| Identity | 11 (includes firm_control_objective_embeddings) | Yes | Yes |
| Audit Core | 27 (9 system reference + 18 tenant-scoped) | Yes (tenant tables) | Yes (tenant tables); `common_controls` has partial RLS — NULL-firm rows shared across tenants |
| Workpaper | 4 (workpapers, workpaper_sign_offs, workpaper_versions, review_notes) | Yes | Yes |
| Reporting | 2 | Yes | Yes |

The three authorization dimensions:

1. **Firm isolation** — RLS + middleware. Every query scoped to current_firm_id.
2. **Engagement team membership** — Application-layer middleware. Point lookup on `engagement_team_members (engagement_id, user_id)`.
3. **Client user scoping** — Application-layer middleware. Client users see only document requests and evidence items for engagements they are invited to.

System-wide reference tables (`frameworks`, `framework_versions`, `framework_requirements`, `common_controls` WHERE `firm_id IS NULL`, `common_control_satisfies`, `common_control_embeddings`, `framework_requirement_embeddings`, `control_objective_library`, `control_objective_library_mappings`) have no `firm_id` (or allow NULL) and no conventional RLS — read-only reference data shared across all tenants, maintained by Axiom from SCF / OSCAL / AICPA / CIS catalogs on a quarterly cadence.

---

## 11. Journey-to-Entity Traceability

| Journey | Persona | Primary Entities |
|---|---|---|
| 1: Firm Setup | FirmAdmin | Firm, User, Invitation, MethodologyTemplate, Engagement, Control, TestProcedure, Workpaper, ClientAcceptance |
| 2: Staff Onboarding | Staff Auditor | User, Invitation, Notification, EngagementTeamMember |
| 3: Engagement Scoping | Partner | Engagement, EngagementTeamMember, EngagementFramework, Client, Control, TestProcedure, ClientAcceptance, EngagementQualityReview, FirmControlObjectiveMapping, AIDecision |
| 4: Cross-Framework Mapping | Manager / FirmAdmin | CommonControl, CommonControlSatisfies, FrameworkRequirement, FrameworkVersion, FirmControlObjective, FirmControlObjectiveMapping, AIDecision (EvidenceControlMapping, GapAnalysis, FrameworkMigration) |
| 5: Control Testing | Staff Auditor | Control, TestProcedure, EvidenceItem, EvidenceLink, EvidenceItemSupports, CommonControl, Workpaper, WorkpaperVersion, WorkpaperSignOff, AIDecision |
| 6: Workpaper Review | Detailed / General / Final Reviewer | Workpaper, WorkpaperVersion, WorkpaperSignOff, ReviewNote, AuditLog, Notification |
| 7: Document Requests | Staff Auditor | DocumentRequest, ClientHubToken, EvidenceItem, EvidenceLink, EvidenceItemSupports, AIDecision, AuditLog |
| 8: Client Fulfillment | Client Contact | DocumentRequest, ClientHubToken, DelegationToken, EvidenceItem, AuditLog |
| 9: Reporting & Archive | Partner | Engagement, Report, ReportVersion, Workpaper, WorkpaperVersion, AuditLog |
| 10: EQR | EQR Reviewer | EngagementQualityReview, EQRFinding, AIDecision, AuditLog |
| 11: Multi-Framework Integrated Engagement | Partner / Manager | Engagement, EngagementFramework (multiple per engagement), FrameworkVersion, CommonControl, CommonControlSatisfies, EvidenceItemSupports, Control, TestProcedure, AIDecision (EvidenceControlMapping) — one control tested once, evidence satisfies requirements across all in-scope frameworks |
| 12: Continuous Assurance / Drift-Triggered Re-Testing | Client / Auditor | EvidenceItem, EvidenceItemSupports (period windows), CommonControl, Notification (DriftDetected, EvidenceExpiring), AIDecision (DriftRetest, FindingsTriage), AuditLog |

### Entity Count Summary

- **Total domain entities:** ~47 (added: FrameworkVersion, CommonControl, CommonControlSatisfies, CommonControlEmbedding, EvidenceItemSupports, WorkpaperSignOff; removed: TrialBalance, TrialBalanceAccount, TrialBalanceAdjustment, ColumnMappingProfile)
- **Total tables in `axiom_db`:** ~47
- **Identity module:** 11 (includes firm_control_objective_embeddings; pgvector required)
- **Audit Core module:** 27 (9 system reference + 18 tenant-scoped; system tables include framework_requirement_embeddings, common_control_embeddings, control_objective_library_embeddings)
- **Workpaper module:** 4 (workpapers, workpaper_sign_offs, workpaper_versions, review_notes)
- **Reporting module:** 2
