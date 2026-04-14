# Axiom Domain and Data Model

**Derived from:** [User Journeys](../user-journeys/all-journeys.md) (10 journeys, clean-slate domain modeling exercise)
**Supersedes:** Section 5 of the [product spec](axiom-spec-design.md) (which now contains an abbreviated summary pointing here)
**Related:** [Backend Architecture](backend-architecture-design.md) (service decomposition and database topology)

---

## Table of Contents

1. [Bounded Context Map](#1-bounded-context-map)
2. [Context 1: Firm Identity](#2-context-1-firm-identity)
3. [Context 2: Regulatory Framework](#3-context-2-regulatory-framework)
4. [Context 3: Firm Methodology](#4-context-3-firm-methodology)
5. [Context 4: Audit Core](#5-context-4-audit-core)
6. [Context 5: Trial Balance & Analytics](#6-context-5-trial-balance--analytics)
7. [Context 6: Workpaper Authoring](#7-context-6-workpaper-authoring)
8. [Context 7: Reporting & Archival](#8-context-7-reporting--archival)
9. [Cross-Cutting Concerns](#9-cross-cutting-concerns)
10. [Data Model](#10-data-model)
11. [Multi-Tenancy and Isolation](#11-multi-tenancy-and-isolation)
12. [Journey-to-Entity Traceability](#12-journey-to-entity-traceability)

---

## 1. Bounded Context Map

Seven bounded contexts and three cross-cutting concerns, derived from the coupling patterns observed across all 10 user journeys.

```
┌─────────────────────────────────────────────────────────────────────┐
│                         AXIOM DOMAIN                                │
│                                                                     │
│  ┌──────────────┐   ┌──────────────────┐   ┌───────────────────┐   │
│  │ 1. FIRM      │   │ 2. REGULATORY    │   │ 3. FIRM           │   │
│  │    IDENTITY  │   │    FRAMEWORK     │   │    METHODOLOGY    │   │
│  │              │   │    (shared ref)  │   │                   │   │
│  │ Firm         │   │ Framework        │   │ MethodologyTmpl   │   │
│  │ User         │   │ FrameworkReq     │   │ FirmControlObj    │   │
│  │ Client       │   │ CtrlObjLibrary   │   │ + mappings        │   │
│  │ Invitation   │   │ + mappings       │   │ + template items  │   │
│  └──────┬───────┘   └────────┬─────────┘   └────────┬──────────┘   │
│         │                    │                       │              │
│         ▼                    ▼                       ▼              │
│  ┌──────────────────────────────────────────────────────────────┐   │
│  │                    4. AUDIT CORE                              │   │
│  │  (Engagement lifecycle + control testing + evidence + PBC)   │   │
│  │                                                              │   │
│  │  Engagement ──► Control ──► TestProcedure                    │   │
│  │  EngagementTeamMember      EvidenceItem ──► EvidenceLink     │   │
│  │  EngagementFramework       DocumentRequest                   │   │
│  │  ClientAcceptance          ClientHubToken                    │   │
│  │  EngagementQualityReview   DelegationToken                   │   │
│  │  EQRFinding                                                  │   │
│  └──────────────────────┬───────────────────────────────────────┘   │
│                         │                                           │
│         ┌───────────────┼───────────────┐                          │
│         ▼               ▼               ▼                          │
│  ┌──────────────┐ ┌──────────────┐ ┌──────────────┐               │
│  │ 5. TRIAL     │ │ 6. WORKPAPER │ │ 7. REPORTING │               │
│  │    BALANCE   │ │    AUTHORING │ │              │               │
│  │              │ │              │ │ Report       │               │
│  │ TrialBalance │ │ Workpaper    │ │ ReportVer    │               │
│  │ TB Account   │ │ WP Version   │ │              │               │
│  │ TB Adjustment│ │ ReviewNote   │ │              │               │
│  │ ColMapProfile│ │              │ │              │               │
│  └──────────────┘ └──────────────┘ └──────────────┘               │
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
EvidenceItem → EvidenceLink → TestProcedure → Control
  → FirmControlObjective → FirmControlObjectiveMapping
    → FrameworkRequirement (SOC 2 CC6.1)
    → FrameworkRequirement (ISO 27001 A.8.3)
    → FrameworkRequirement (HIPAA §164.312(a)(1))
```

Splitting this chain across contexts would force distributed transactions or eventual consistency on operations that must be atomic (e.g., accepting a document request must atomically create an EvidenceLink and update DocumentRequest.status).

**Contexts 5, 6, 7 are separate** because they have distinct scaling profiles (TB is spreadsheet-intensive, Workpaper has real-time collaborative editing via WebSocket, Reporting is async generation) and communicate with Audit Core via plain UUID references rather than foreign keys.

### Aggregates within Audit Core

| Aggregate Root | Contains | Consistency Boundary |
|---|---|---|
| **Engagement** | EngagementTeamMember, EngagementFramework, ClientAcceptance, EngagementQualityReview, EQRFinding | Lifecycle state machine, quality gates |
| **Control** | TestProcedure | Control status depends on procedure results |
| **EvidenceItem** | EvidenceLink | Evidence exists at firm+client level, links are engagement-scoped |
| **DocumentRequest** | ClientHubToken, DelegationToken | PBC lifecycle, client interaction |

### Database Mapping

| Bounded Context | Database | Isolation |
|---|---|---|
| 1: Firm Identity + 3: Firm Methodology | `identity_db` | Application-layer (FK to firms) |
| 2: Regulatory Framework + 4: Audit Core + cross-cutting | `core_db` | RLS + application-layer |
| 5: Trial Balance | `trial_balance_db` | Application-layer |
| 6: Workpaper Authoring | `workpaper_db` | Application-layer |
| 7: Reporting | `reporting_db` | Application-layer |

See [Backend Architecture](backend-architecture-design.md) for service decomposition and database topology details.

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
| primary_audit_types | jsonb | Multi-select from intake: FinancialAudit, SOC2, ISO27001, HIPAA |
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

## 3. Context 2: Regulatory Framework

**Purpose:** External regulatory knowledge base. Not tenant-scoped — shared read-only reference data across all firms.
**Journeys:** 3 (Engagement Scoping), 5 (Control Testing).

### Framework (Aggregate Root)

A specific version of a standards framework. ISO 27001:2013 and ISO 27001:2022 are separate rows.

| Attribute | Type | Description |
|---|---|---|
| id | uuid | |
| name | text | e.g., "SOC 2", "ISO 27001", "HIPAA" |
| version | text | e.g., "TSC 2017", "2022" |
| effective_date | date | |
| deprecated_at | date | Nullable |
| governing_body | text | AICPA, ISO, HHS, PCAOB |

### FrameworkRequirement

A single criterion or control within a framework.

| Attribute | Type | Description |
|---|---|---|
| id | uuid | |
| framework_id | uuid | FK → Framework |
| identifier | text | e.g., "CC6.1", "A.8.3", "§164.312(a)(1)" |
| title | text | |
| description | text | |
| category | text | Grouping within framework |
| sort_order | integer | Display ordering |

### ControlObjectiveLibrary (Aggregate Root)

System-maintained semantic control objectives, independent of any framework. These encode the ~80% overlap between frameworks and enable cross-framework evidence reuse.

| Attribute | Type | Description |
|---|---|---|
| id | uuid | |
| name | text | e.g., "Access to production systems is restricted to authorized personnel" |
| description | text | |
| tags | jsonb | Categorization tags |

### ControlObjectiveLibraryMapping

Maps a library objective to framework requirements across all frameworks. This is the foundation of cross-framework evidence reuse.

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
| client_id | uuid | Ref to Client in identity_db |
| name | text | |
| engagement_type | enum | FinancialAudit_Private, FinancialAudit_Public, SOC1, SOC2, ISO27001, HIPAA, AgreedUponProcedures, Advisory |
| primary_framework_id | uuid | FK → Framework |
| period_start | date | |
| period_end | date | |
| status | enum | Planning, Fieldwork, Review, Reporting, Finalized, Archived |
| prior_engagement_id | uuid | FK → Engagement, nullable — for rollforward |
| methodology_template_id | uuid | Ref to MethodologyTemplate in identity_db |
| report_issued_at | timestamptz | Populated when report is issued |
| assembly_deadline | date | Computed: report_issued_at + 60 days (AICPA) or + 45 days (PCAOB) |
| retention_deadline | date | Computed: report_issued_at + 5 yrs (AICPA/SOC/ISO), + 7 yrs (PCAOB), + 6 yrs (HIPAA) |
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

#### EngagementTeamMember

Associates users to an engagement with an engagement-level role.

| Attribute | Type | Description |
|---|---|---|
| id | uuid | |
| engagement_id | uuid | FK → Engagement |
| user_id | uuid | Ref to User in identity_db |
| engagement_role | text | Partner, Manager, Staff — the role on this specific engagement |
| assigned_at | timestamptz | |
| removed_at | timestamptz | Nullable |

**Invariants:**
- A user with EQReviewer role on the EngagementQualityReview for this engagement cannot also be an EngagementTeamMember.
- Used for engagement-level access control checks.

#### EngagementFramework

Supports multi-framework engagements (e.g., integrated SOC 2 + ISO 27001).

| Attribute | Type | Description |
|---|---|---|
| id | uuid | |
| engagement_id | uuid | FK → Engagement |
| framework_id | uuid | FK → Framework |
| framework_version | text | |
| is_primary | boolean | Exactly one must be true per engagement |

#### ClientAcceptance

Per-engagement quality risk documentation required by SQMS 1. A regulatory gate.

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

Formal EQR record per SQMS 2 / PCAOB AS 1220.

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
- Mandatory for PCAOB engagements; optional per firm policy for nonissuer.
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
| firm_control_objective_id | uuid | Ref to FirmControlObjective in identity_db |
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
| sampling_method | enum | Systematic, Random, MonetaryUnit — nullable |
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
- Sampling documentation (population_size, sample_size, sampling_method) is required for PCAOB engagements per AS 2315.
- Status progression: NotStarted → InProgress → Complete or Exception.

---

### Evidence Aggregate

#### EvidenceItem (Aggregate Root)

A single uploaded document or artifact. **Stored at the firm + client level, not the engagement level** — this enables year-over-year and cross-framework reuse without re-uploading.

| Attribute | Type | Description |
|---|---|---|
| id | uuid | |
| firm_id | uuid | |
| client_id | uuid | Ref to Client in identity_db |
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

Connects an EvidenceItem to a specific TestProcedure. The critical junction in the cross-framework evidence chain.

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
- When an EvidenceLink is created, the evidence simultaneously satisfies all FrameworkRequirements mapped to the control's FirmControlObjective.
- Links are frozen after engagement finalization.
- AI-suggested links require auditor accept/modify/reject action.

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
| assigned_to_id | uuid | Ref to client user in identity_db |
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
| Workpapers | New drafts; prior_workpaper_id set; prior year visible as read-only sidebar |
| TrialBalance | New import required; prior year accessible as read-only reference |
| Reports | New document; prior year accessible for reference only |
| AIDecisions | Not carried forward; AI re-analyzes fresh evidence |
| ClientAcceptance | New record required; must be refreshed annually |
| EngagementQualityReview | New record required if applicable |

---

## 6. Context 5: Trial Balance & Analytics

**Purpose:** Financial audit data — import, account mapping, lead schedules, adjustments, analytical procedures, and population sampling. Exists only for FinancialAudit engagement types.
**Journeys:** 4 (Trial Balance).

### TrialBalance (Aggregate Root)

Container for an imported trial balance.

| Attribute | Type | Description |
|---|---|---|
| id | uuid | |
| engagement_id | uuid | Plain UUID ref — no FK to core_db |
| firm_id | uuid | For tenant isolation |
| period_date | date | |
| import_source | text | Accounting system name (QBO, NetSuite, Sage, Xero) |
| column_mapping_profile_id | uuid | FK → ColumnMappingProfile, nullable |
| imported_at | timestamptz | |
| imported_by_id | uuid | |

### TrialBalanceAccount

Individual account row within a trial balance.

| Attribute | Type | Description |
|---|---|---|
| id | uuid | |
| trial_balance_id | uuid | FK → TrialBalance |
| account_number | text | |
| account_name | text | |
| account_type | enum | Asset, Liability, Equity, Revenue, Expense |
| balance_debit | numeric(19,4) | |
| balance_credit | numeric(19,4) | |
| net_balance | numeric(19,4) | Computed: balance_debit - balance_credit |
| mapped_fs_line_item | text | Financial statement line item classification |
| mapping_status | enum | Unmapped, AISuggested, Confirmed, Overridden |
| ai_decision_id | uuid | Plain UUID ref to AIDecision in core_db |
| confirmed_by_id | uuid | |
| confirmed_at | timestamptz | |

**Invariants:**
- AI mapping runs immediately after import (Claude Haiku classifies each account).
- Prior year confirmed mappings are pre-loaded as starting suggestions on rollforward.
- Bulk-confirm is available for high-confidence mappings (>0.85).
- Total debits must equal total credits — non-zero difference is flagged at import validation.

### TrialBalanceAdjustment

Proposed, passed, or waived audit adjustments.

| Attribute | Type | Description |
|---|---|---|
| id | uuid | |
| trial_balance_id | uuid | FK → TrialBalance |
| account_id | uuid | FK → TrialBalanceAccount |
| amount | numeric(19,4) | |
| description | text | |
| adjustment_type | enum | Proposed, Passed, Waived |
| waiver_reason | text | Required when type = Waived (regulatory requirement) |
| proposed_by_id | uuid | |
| proposed_at | timestamptz | |
| approved_by_id | uuid | Manager or Partner |
| approved_at | timestamptz | |

**Invariants:**
- Waived adjustments require a documented reason (AU-C requirement).
- All proposed adjustments are retained whether passed or not — cannot be deleted.
- Lead schedules reflect adjustments in real-time.
- Cumulative waived adjustments are tracked against materiality — auto-flagged when approaching threshold.

### ColumnMappingProfile

Saved column mapping configurations per accounting system for reuse across imports.

| Attribute | Type | Description |
|---|---|---|
| id | uuid | |
| firm_id | uuid | |
| name | text | e.g., "NetSuite Standard Export" |
| accounting_system | text | |
| column_mappings | jsonb | Which CSV/Excel columns map to account_number, account_name, debit, credit |
| created_at | timestamptz | |

---

## 7. Context 6: Workpaper Authoring

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
| status | enum | Draft, PreparedPendingReview, InReview, ReviewNotesOpen, ReviewComplete, SignedOff |
| prepared_by_id | uuid | |
| reviewed_by_id | uuid | |
| signed_off_by_id | uuid | |
| is_locked | boolean | True after engagement finalization |
| prior_workpaper_id | uuid | FK → Workpaper, nullable — for rollforward |
| created_at | timestamptz | |

**Status lifecycle:**

```
Draft
  → PreparedPendingReview  [Staff submits; system validates is_ai_draft = false]
  → InReview               [Manager opens for review]
  → ReviewNotesOpen        [Manager raises notes; returns to Staff]
  → InReview               [Staff addresses notes; resubmits]
  → ReviewComplete         [Manager clears; all notes resolved]
  → SignedOff              [Partner final sign-off]
```

**Invariants:**
- Submit for review is blocked if any WorkpaperVersion has is_ai_draft = true (PCAOB AS 1105).
- Once submitted, the workpaper is locked for the preparer — only the reviewer can modify or return it.
- After engagement finalization, is_locked = true — modifications require an addendum.
- Sign-off creates a timestamped, named AuditLog entry — cannot be backdated.
- Sign-off hierarchy enforced: Staff prepares → Manager reviews → Partner signs off (SQMS 1, AU-C 220).

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
| is_ai_draft | boolean | True until a human edits any content — per PCAOB AS 1105 |
| is_addendum | boolean | True for post-finalization modifications — per AU-C 230 §.16 |
| addendum_reason | text | Required when is_addendum = true |

**Invariants:**
- is_ai_draft is set to true when AI generates a draft; cleared to false on any human edit.
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
- **Cannot be deleted** — AU-C 230 requires retention of all review notes.
- Open review notes block workpaper advancement (InReview → ReviewComplete).
- Resolution workflow: reviewer creates → staff responds → reviewer resolves.
- Each note creation, response, and resolution creates an AuditLog entry.

---

## 8. Context 7: Reporting & Archival

**Purpose:** Final deliverable generation, client review, issuance, and regulatory archival. Async report generation.
**Journeys:** 9 (Reporting & Archive).

### Report (Aggregate Root)

The engagement's final deliverable.

| Attribute | Type | Description |
|---|---|---|
| id | uuid | |
| engagement_id | uuid | Plain UUID ref |
| firm_id | uuid | For tenant isolation |
| report_type | enum | SOC2Type1, SOC2Type2, SOC1Type1, SOC1Type2, FinancialAuditOpinion, AgreedUponProcedures, ManagementLetter |
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
- AI-drafted report sections (Tier 2) follow the same is_ai_draft rules as workpapers.
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
| is_ai_draft | boolean | Same semantics as WorkpaperVersion |

---

## 9. Cross-Cutting Concerns

These operate across all bounded contexts.

### AIDecision

Every AI output that could affect audit content is recorded. Required for PCAOB engagements; best practice for all.

| Attribute | Type | Description |
|---|---|---|
| id | uuid | |
| firm_id | uuid | |
| engagement_id | uuid | FK → Engagement, nullable — some AI operations are firm-level |
| context_type | enum | ControlMapping, RiskAssessment, TrialBalanceMapping, DocumentCompleteness, EvidenceLinkSuggestion, WorkpaperDraft, ReportDraft, SamplingRecommendation, AnomalyDetection |
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
- Satisfies regulatory immutability requirements (AU-C 230) and GDPR audit trail obligation.
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

## 10. Data Model

The data model translates the domain model above into PostgreSQL tables across five databases. This section covers physical design decisions that go beyond the domain model: column types, constraints, indexes, enum definitions, and cross-database reference strategy.

### Database Topology

```
RDS PostgreSQL Instance (Multi-AZ production, Single-AZ dev/staging)
├── identity_db      → Contexts 1 + 3 (Firm Identity, Firm Methodology)
├── core_db          → Contexts 2 + 4 + cross-cutting (Framework, Audit Core, AI/Audit/Notification)
├── trial_balance_db → Context 5 (Trial Balance)
├── workpaper_db     → Context 6 (Workpaper Authoring)
└── reporting_db     → Context 7 (Reporting)
```

Each service has its own Postgres user with access only to its own database. No cross-database joins.

### Cross-Database Reference Strategy

References to entities in other databases use plain UUID columns with no foreign key constraint. Referential integrity is enforced at the application layer via REST calls. Within a database, standard foreign key constraints apply.

| Pattern | Usage |
|---|---|
| FK constraint | References within the same database |
| Plain UUID column | References to entities in other databases (e.g., `engagement_id` in workpaper_db) |
| REST validation | When a request requires confirming cross-database entity existence or access |

### Enum Types

**identity_db:**

```sql
CREATE TYPE user_role AS ENUM (
  'FirmAdmin','Partner','Manager','Staff','EQReviewer',
  'ClientAdmin','ClientUser','ViewOnly');
CREATE TYPE auth_method AS ENUM ('Password','OAuth','SAML');
CREATE TYPE notification_frequency AS ENUM ('RealTime','Daily','Weekly');
CREATE TYPE invitation_status AS ENUM ('Sent','Accepted','Expired');
```

**core_db:**

```sql
CREATE TYPE engagement_type AS ENUM (
  'FinancialAudit_Private','FinancialAudit_Public',
  'SOC1','SOC2','ISO27001','HIPAA',
  'AgreedUponProcedures','Advisory');
CREATE TYPE engagement_status AS ENUM (
  'Planning','Fieldwork','Review','Reporting','Finalized','Archived');
CREATE TYPE control_status AS ENUM (
  'NotStarted','InProgress','Complete','Exception','NotApplicable');
CREATE TYPE procedure_type AS ENUM (
  'Inquiry','Observation','InspectionOfDocument','Reperformance','Analytics');
CREATE TYPE procedure_status AS ENUM ('NotStarted','InProgress','Complete','Exception');
CREATE TYPE sampling_method AS ENUM ('Systematic','Random','MonetaryUnit');
CREATE TYPE evidence_source AS ENUM (
  'ClientUpload','CloudIntegration','APIImport','AuditorGenerated');
CREATE TYPE extraction_status AS ENUM ('Pending','Complete','Failed');
CREATE TYPE doc_request_status AS ENUM (
  'Pending','Submitted','InReview','Accepted','Rejected','Overdue');
CREATE TYPE hub_token_status AS ENUM ('Active','Expired','Revoked');
CREATE TYPE delegation_token_status AS ENUM ('Active','Used','Expired');
CREATE TYPE eqr_status AS ENUM ('Assigned','InProgress','Complete');
CREATE TYPE eqr_conclusion AS ENUM ('Satisfied','SatisfiedWithConcerns','NotSatisfied');
CREATE TYPE finding_severity AS ENUM ('Observation','Recommendation','RequiredAction');
CREATE TYPE finding_status AS ENUM ('Pending','Addressed','Confirmed');
CREATE TYPE ai_context_type AS ENUM (
  'ControlMapping','RiskAssessment','TrialBalanceMapping',
  'DocumentCompleteness','EvidenceLinkSuggestion','WorkpaperDraft',
  'ReportDraft','SamplingRecommendation','AnomalyDetection');
CREATE TYPE ai_review_action AS ENUM ('Pending','Accepted','Modified','Rejected');
CREATE TYPE actor_type AS ENUM ('User','System','AIAgent');
CREATE TYPE notification_type AS ENUM (
  'EngagementAssignment','ReviewNoteAdded','ReviewNoteResolved',
  'DocumentRequestStatus','PhaseTransition','EQRNotification',
  'ReminderEscalation','ArchivalConfirmation','RetentionWarning');
CREATE TYPE delivery_channel AS ENUM ('InApp','Email','Both');
```

**trial_balance_db:**

```sql
CREATE TYPE account_type AS ENUM ('Asset','Liability','Equity','Revenue','Expense');
CREATE TYPE mapping_status AS ENUM ('Unmapped','AISuggested','Confirmed','Overridden');
CREATE TYPE adjustment_type AS ENUM ('Proposed','Passed','Waived');
```

**workpaper_db:**

```sql
CREATE TYPE workpaper_type AS ENUM (
  'LeadSchedule','TestPaper','Memo','ConfirmationLetter',
  'SamplingWorksheet','ManagementLetter','AnalyticalProcedures','Other');
CREATE TYPE workpaper_status AS ENUM (
  'Draft','PreparedPendingReview','InReview',
  'ReviewNotesOpen','ReviewComplete','SignedOff');
CREATE TYPE review_note_severity AS ENUM ('Question','Suggestion','RequiredChange');
CREATE TYPE review_note_status AS ENUM ('Open','Responded','Resolved');
```

**reporting_db:**

```sql
CREATE TYPE report_type AS ENUM (
  'SOC2Type1','SOC2Type2','SOC1Type1','SOC1Type2',
  'FinancialAuditOpinion','AgreedUponProcedures','ManagementLetter');
CREATE TYPE report_status AS ENUM ('Draft','ClientReview','FirmReview','Issued','Archived');
```

### Table Definitions by Database

All UUID primary keys use `DEFAULT gen_random_uuid()`. All timestamps use `timestamptz` with `DEFAULT now()` where applicable.

#### identity_db

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

#### core_db — System Reference Tables (no RLS, not tenant-scoped)

**frameworks** — `id (uuid PK)`, `name (text NOT NULL)`, `version (text NOT NULL)`, `effective_date (date NOT NULL)`, `deprecated_at (date)`, `governing_body (text NOT NULL)`. **Unique:** `(name, version)`.

**framework_requirements** — `id (uuid PK)`, `framework_id (uuid FK → frameworks NOT NULL)`, `identifier (text NOT NULL)`, `title (text NOT NULL)`, `description (text)`, `category (text)`, `sort_order (integer NOT NULL)`. **Unique:** `(framework_id, identifier)`. **Indexes:** `(framework_id, sort_order)`.

**control_objective_library** — `id (uuid PK)`, `name (text NOT NULL)`, `description (text NOT NULL)`, `tags (jsonb NOT NULL DEFAULT '[]')`.

**control_objective_library_mappings** — `id (uuid PK)`, `library_objective_id (uuid FK → control_objective_library NOT NULL)`, `framework_requirement_id (uuid FK → framework_requirements NOT NULL)`. **Unique:** `(library_objective_id, framework_requirement_id)`.

#### core_db — Tenant-Scoped Tables (RLS via firm_id)

All tables below carry `firm_id uuid NOT NULL` with an index. RLS policy: `USING (firm_id = current_setting('app.current_firm_id')::uuid)`.

**engagements** — `id (uuid PK)`, `firm_id`, `client_id (uuid NOT NULL)`, `name (text NOT NULL)`, `engagement_type (engagement_type NOT NULL)`, `primary_framework_id (uuid FK → frameworks NOT NULL)`, `period_start (date NOT NULL)`, `period_end (date NOT NULL)`, `status (engagement_status NOT NULL DEFAULT 'Planning')`, `prior_engagement_id (uuid FK → engagements)`, `methodology_template_id (uuid)`, `report_issued_at (timestamptz)`, `assembly_deadline (date)`, `retention_deadline (date)`, `finalized_at (timestamptz)`, `archived_at (timestamptz)`, `created_at (timestamptz NOT NULL)`. **Indexes:** `(firm_id)`, `(client_id)`, `(status)`, `(prior_engagement_id)`.

**engagement_team_members** — `id (uuid PK)`, `firm_id`, `engagement_id (uuid FK → engagements NOT NULL)`, `user_id (uuid NOT NULL)`, `engagement_role (text NOT NULL)`, `assigned_at (timestamptz NOT NULL)`, `removed_at (timestamptz)`. **Unique:** `(engagement_id, user_id) WHERE removed_at IS NULL`. **Indexes:** `(engagement_id, user_id)`.

**engagement_frameworks** — `id (uuid PK)`, `firm_id`, `engagement_id (uuid FK → engagements NOT NULL)`, `framework_id (uuid FK → frameworks NOT NULL)`, `framework_version (text NOT NULL)`, `is_primary (boolean NOT NULL DEFAULT false)`. **Unique:** `(engagement_id, framework_id)`.

**client_acceptances** — `id (uuid PK)`, `firm_id`, `engagement_id (uuid FK → engagements NOT NULL UNIQUE)`, `quality_risks_identified (jsonb NOT NULL DEFAULT '[]')`, `firm_responses (jsonb NOT NULL DEFAULT '[]')`, `independence_confirmed (boolean NOT NULL DEFAULT false)`, `independence_confirmed_by_id (uuid)`, `accepted_by_id (uuid)`, `accepted_at (timestamptz)`, `created_at (timestamptz NOT NULL)`.

**engagement_quality_reviews** — `id (uuid PK)`, `firm_id`, `engagement_id (uuid FK → engagements NOT NULL UNIQUE)`, `reviewer_id (uuid NOT NULL)`, `independence_documented_at (timestamptz)`, `status (eqr_status NOT NULL DEFAULT 'Assigned')`, `scope_notes (text)`, `conclusion (eqr_conclusion)`, `signed_off_at (timestamptz)`, `created_at (timestamptz NOT NULL)`.

**eqr_findings** — `id (uuid PK)`, `firm_id`, `eqr_id (uuid FK → engagement_quality_reviews NOT NULL)`, `description (text NOT NULL)`, `severity (finding_severity NOT NULL)`, `status (finding_status NOT NULL DEFAULT 'Pending')`, `team_response (text)`, `responded_at (timestamptz)`, `confirmed_by_id (uuid)`, `confirmed_at (timestamptz)`, `created_at (timestamptz NOT NULL)`.

**controls** — `id (uuid PK)`, `firm_id`, `engagement_id (uuid FK → engagements NOT NULL)`, `firm_control_objective_id (uuid NOT NULL)`, `description (text NOT NULL)`, `control_owner_id (uuid)`, `auditor_assigned_to_id (uuid)`, `status (control_status NOT NULL DEFAULT 'NotStarted')`, `is_key_control (boolean NOT NULL DEFAULT false)`, `prior_control_id (uuid FK → controls)`, `created_at (timestamptz NOT NULL)`. **Indexes:** `(firm_id)`, `(engagement_id)`, `(auditor_assigned_to_id)`, `(engagement_id, status)`.

**test_procedures** — `id (uuid PK)`, `firm_id`, `control_id (uuid FK → controls NOT NULL)`, `procedure_type (procedure_type NOT NULL)`, `description (text NOT NULL)`, `expected_result (text)`, `population_size (integer)`, `sample_size (integer)`, `sampling_method (sampling_method)`, `result (text)`, `exceptions_noted (text)`, `conclusion (text)`, `performed_by_id (uuid)`, `performed_at (timestamptz)`, `reviewed_by_id (uuid)`, `reviewed_at (timestamptz)`, `status (procedure_status NOT NULL DEFAULT 'NotStarted')`, `prior_procedure_id (uuid FK → test_procedures)`. **Indexes:** `(control_id)`, `(performed_by_id)`.

**evidence_items** — `id (uuid PK)`, `firm_id`, `client_id (uuid NOT NULL)`, `filename (text NOT NULL)`, `storage_path (text NOT NULL)`, `content_type (text NOT NULL)`, `file_size_bytes (bigint NOT NULL)`, `uploaded_by_id (uuid NOT NULL)`, `uploaded_at (timestamptz NOT NULL)`, `source_type (evidence_source NOT NULL)`, `source_integration (text)`, `extracted_text (text)`, `extraction_status (extraction_status NOT NULL DEFAULT 'Pending')`, `is_sensitive (boolean NOT NULL DEFAULT false)`. **Indexes:** `(firm_id, client_id)`.

**evidence_embeddings** — `id (uuid PK)`, `evidence_item_id (uuid FK → evidence_items NOT NULL UNIQUE)`, `embedding (vector(1024) NOT NULL)`, `created_at (timestamptz NOT NULL)`. **Indexes:** `USING ivfflat (embedding vector_cosine_ops)`.

**evidence_links** — `id (uuid PK)`, `firm_id`, `evidence_item_id (uuid FK → evidence_items NOT NULL)`, `test_procedure_id (uuid FK → test_procedures NOT NULL)`, `linked_by_id (uuid NOT NULL)`, `linked_at (timestamptz NOT NULL)`, `notes (text)`, `ai_suggested (boolean NOT NULL DEFAULT false)`, `ai_decision_id (uuid FK → ai_decisions)`. **Unique:** `(evidence_item_id, test_procedure_id)`. **Indexes:** `(test_procedure_id)`, `(evidence_item_id)`.

**document_requests** — `id (uuid PK)`, `firm_id`, `engagement_id (uuid FK → engagements NOT NULL)`, `control_id (uuid FK → controls)`, `test_procedure_id (uuid FK → test_procedures)`, `assigned_to_id (uuid)`, `title (text NOT NULL)`, `instructions (text NOT NULL)`, `due_date (date)`, `status (doc_request_status NOT NULL DEFAULT 'Pending')`, `reminder_count (integer NOT NULL DEFAULT 0)`, `last_reminder_sent_at (timestamptz)`, `fulfilled_by_evidence_item_id (uuid FK → evidence_items)`, `sent_at (timestamptz)`, `created_at (timestamptz NOT NULL)`. **Indexes:** `(engagement_id)`, `(engagement_id, status)`, `(assigned_to_id)`.

**client_hub_tokens** — `id (uuid PK)`, `firm_id`, `engagement_id (uuid FK → engagements NOT NULL)`, `client_id (uuid NOT NULL)`, `token_hash (text NOT NULL UNIQUE)`, `valid_until (timestamptz NOT NULL)`, `created_by_id (uuid NOT NULL)`, `status (hub_token_status NOT NULL DEFAULT 'Active')`, `created_at (timestamptz NOT NULL)`. **Indexes:** `(token_hash)`.

**delegation_tokens** — `id (uuid PK)`, `firm_id`, `document_request_id (uuid FK → document_requests NOT NULL)`, `delegated_by_id (uuid NOT NULL)`, `delegate_email (text NOT NULL)`, `token_hash (text NOT NULL UNIQUE)`, `custom_message (text)`, `valid_until (timestamptz NOT NULL)`, `status (delegation_token_status NOT NULL DEFAULT 'Active')`, `created_at (timestamptz NOT NULL)`. **Indexes:** `(token_hash)`, `(document_request_id)`.

**ai_decisions** — `id (uuid PK)`, `firm_id`, `engagement_id (uuid FK → engagements)`, `context_type (ai_context_type NOT NULL)`, `context_id (uuid NOT NULL)`, `context_table (text NOT NULL)`, `model_id (text NOT NULL)`, `input_token_count (integer)`, `output_token_count (integer)`, `raw_output (jsonb NOT NULL)`, `suggested_value (text)`, `confidence (real CHECK (confidence >= 0 AND confidence <= 1))`, `explanation (text)`, `review_action (ai_review_action NOT NULL DEFAULT 'Pending')`, `accepted_value (text)`, `reviewed_by_id (uuid)`, `reviewed_at (timestamptz)`, `created_at (timestamptz NOT NULL)`. **Indexes:** `(firm_id)`, `(engagement_id)`, `(context_type, context_id)`.

**audit_log** — `id (bigint PK GENERATED ALWAYS AS IDENTITY)`, `firm_id (uuid NOT NULL)`, `actor_id (uuid)`, `actor_type (actor_type NOT NULL)`, `action (text NOT NULL)`, `resource_type (text NOT NULL)`, `resource_id (uuid)`, `old_value (jsonb)`, `new_value (jsonb)`, `ip_address (inet)`, `user_agent (text)`, `occurred_at (timestamptz NOT NULL DEFAULT now())`. **Immutability:** `CREATE RULE audit_log_no_update AS ON UPDATE TO audit_log DO INSTEAD NOTHING; CREATE RULE audit_log_no_delete AS ON DELETE TO audit_log DO INSTEAD NOTHING;`. **Indexes:** `(firm_id, occurred_at)`, `(resource_type, resource_id)`, `(actor_id)`.

**notifications** — `id (uuid PK)`, `firm_id (uuid NOT NULL)`, `recipient_id (uuid NOT NULL)`, `notification_type (notification_type NOT NULL)`, `title (text NOT NULL)`, `body (text)`, `deep_link (text)`, `is_read (boolean NOT NULL DEFAULT false)`, `delivery_channel (delivery_channel NOT NULL)`, `created_at (timestamptz NOT NULL)`. **Indexes:** `(recipient_id, is_read, created_at DESC)`.

#### trial_balance_db

Application-layer isolation (`WHERE firm_id = $1`), no RLS.

**trial_balances** — `id (uuid PK)`, `engagement_id (uuid NOT NULL)`, `firm_id (uuid NOT NULL)`, `period_date (date NOT NULL)`, `import_source (text)`, `column_mapping_profile_id (uuid FK → column_mapping_profiles)`, `imported_at (timestamptz NOT NULL)`, `imported_by_id (uuid NOT NULL)`. **Indexes:** `(firm_id)`, `(engagement_id)`.

**trial_balance_accounts** — `id (uuid PK)`, `trial_balance_id (uuid FK → trial_balances NOT NULL)`, `account_number (text NOT NULL)`, `account_name (text NOT NULL)`, `account_type (account_type)`, `balance_debit (numeric(19,4) NOT NULL DEFAULT 0)`, `balance_credit (numeric(19,4) NOT NULL DEFAULT 0)`, `net_balance (numeric(19,4) GENERATED ALWAYS AS (balance_debit - balance_credit) STORED)`, `mapped_fs_line_item (text)`, `mapping_status (mapping_status NOT NULL DEFAULT 'Unmapped')`, `ai_decision_id (uuid)`, `confirmed_by_id (uuid)`, `confirmed_at (timestamptz)`. **Indexes:** `(trial_balance_id)`, `(trial_balance_id, mapping_status)`.

**trial_balance_adjustments** — `id (uuid PK)`, `trial_balance_id (uuid FK → trial_balances NOT NULL)`, `account_id (uuid FK → trial_balance_accounts NOT NULL)`, `amount (numeric(19,4) NOT NULL)`, `description (text NOT NULL)`, `adjustment_type (adjustment_type NOT NULL)`, `waiver_reason (text)`, `proposed_by_id (uuid NOT NULL)`, `proposed_at (timestamptz NOT NULL)`, `approved_by_id (uuid)`, `approved_at (timestamptz)`. **Check:** `(adjustment_type != 'Waived' OR waiver_reason IS NOT NULL)`. **Indexes:** `(trial_balance_id)`.

**column_mapping_profiles** — `id (uuid PK)`, `firm_id (uuid NOT NULL)`, `name (text NOT NULL)`, `accounting_system (text)`, `column_mappings (jsonb NOT NULL)`, `created_at (timestamptz NOT NULL)`. **Indexes:** `(firm_id)`.

#### workpaper_db

Application-layer isolation.

**workpapers** — `id (uuid PK)`, `engagement_id (uuid NOT NULL)`, `firm_id (uuid NOT NULL)`, `control_id (uuid)`, `workpaper_type (workpaper_type NOT NULL)`, `title (text NOT NULL)`, `content (jsonb NOT NULL DEFAULT '{}')`, `status (workpaper_status NOT NULL DEFAULT 'Draft')`, `prepared_by_id (uuid)`, `reviewed_by_id (uuid)`, `signed_off_by_id (uuid)`, `is_locked (boolean NOT NULL DEFAULT false)`, `prior_workpaper_id (uuid FK → workpapers)`, `created_at (timestamptz NOT NULL)`. **Indexes:** `(firm_id)`, `(engagement_id)`, `(engagement_id, status)`.

**workpaper_versions** — `id (uuid PK)`, `workpaper_id (uuid FK → workpapers NOT NULL)`, `version_number (integer NOT NULL)`, `content (jsonb NOT NULL)`, `saved_by_id (uuid NOT NULL)`, `saved_at (timestamptz NOT NULL)`, `is_ai_draft (boolean NOT NULL DEFAULT false)`, `is_addendum (boolean NOT NULL DEFAULT false)`, `addendum_reason (text)`. **Check:** `(is_addendum = false OR addendum_reason IS NOT NULL)`. **Unique:** `(workpaper_id, version_number)`.

**review_notes** — `id (uuid PK)`, `firm_id (uuid NOT NULL)`, `workpaper_id (uuid FK → workpapers NOT NULL)`, `content_anchor (jsonb)`, `created_by_id (uuid NOT NULL)`, `description (text NOT NULL)`, `severity (review_note_severity NOT NULL)`, `status (review_note_status NOT NULL DEFAULT 'Open')`, `response (text)`, `responded_by_id (uuid)`, `responded_at (timestamptz)`, `resolved_by_id (uuid)`, `resolved_at (timestamptz)`, `created_at (timestamptz NOT NULL)`. **Immutability:** `CREATE RULE review_notes_no_delete AS ON DELETE TO review_notes DO INSTEAD NOTHING;`. **Indexes:** `(workpaper_id, status)`.

#### reporting_db

Application-layer isolation.

**reports** — `id (uuid PK)`, `engagement_id (uuid NOT NULL)`, `firm_id (uuid NOT NULL)`, `report_type (report_type NOT NULL)`, `status (report_status NOT NULL DEFAULT 'Draft')`, `content (jsonb NOT NULL DEFAULT '{}')`, `template_id (uuid)`, `generated_at (timestamptz NOT NULL)`, `issued_at (timestamptz)`, `issued_by_id (uuid)`. **Indexes:** `(firm_id)`, `(engagement_id)`.

**report_versions** — `id (uuid PK)`, `report_id (uuid FK → reports NOT NULL)`, `version_number (integer NOT NULL)`, `content (jsonb NOT NULL)`, `saved_by_id (uuid NOT NULL)`, `saved_at (timestamptz NOT NULL)`, `is_ai_draft (boolean NOT NULL DEFAULT false)`. **Unique:** `(report_id, version_number)`.

---

## 11. Multi-Tenancy and Isolation

| Database | Isolation Mechanism | firm_id Indexed | RLS |
|---|---|---|---|
| identity_db | Application-layer (FK to firms) | Yes | No |
| core_db | RLS + application-layer | Yes (all tenant tables) | Yes |
| trial_balance_db | Application-layer WHERE clause | Yes | No |
| workpaper_db | Application-layer WHERE clause | Yes | No |
| reporting_db | Application-layer WHERE clause | Yes | No |

`core_db` uses RLS because it contains the most sensitive data and has the most complex access patterns. The three authorization dimensions:

1. **Firm isolation** — RLS + middleware. Every query scoped to current_firm_id.
2. **Engagement team membership** — Application-layer middleware. Point lookup on `engagement_team_members (engagement_id, user_id)`.
3. **Client user scoping** — Application-layer middleware. Client users see only document requests and evidence items for engagements they are invited to.

System-wide reference tables (`frameworks`, `framework_requirements`, `control_objective_library`, `control_objective_library_mappings`) have no `firm_id` and no RLS — read-only reference data shared across all tenants.

---

## 12. Journey-to-Entity Traceability

| Journey | Persona | Primary Entities |
|---|---|---|
| 1: Firm Setup | FirmAdmin | Firm, User, Invitation, MethodologyTemplate, Engagement, Control, TestProcedure, Workpaper, ClientAcceptance |
| 2: Staff Onboarding | Staff Auditor | User, Invitation, Notification, EngagementTeamMember |
| 3: Engagement Scoping | Partner | Engagement, EngagementTeamMember, EngagementFramework, Client, Control, TestProcedure, ClientAcceptance, EngagementQualityReview, FirmControlObjectiveMapping, AIDecision |
| 4: Trial Balance | Staff Auditor | TrialBalance, TrialBalanceAccount, TrialBalanceAdjustment, ColumnMappingProfile, Workpaper, AIDecision |
| 5: Control Testing | Staff Auditor | Control, TestProcedure, EvidenceItem, EvidenceLink, Workpaper, WorkpaperVersion, AIDecision |
| 6: Workpaper Review | Manager | Workpaper, WorkpaperVersion, ReviewNote, AuditLog, Notification |
| 7: Document Requests | Staff Auditor | DocumentRequest, ClientHubToken, EvidenceItem, EvidenceLink, AIDecision, AuditLog |
| 8: Client Fulfillment | Client Contact | DocumentRequest, ClientHubToken, DelegationToken, EvidenceItem, AuditLog |
| 9: Reporting & Archive | Partner | Engagement, Report, ReportVersion, Workpaper, WorkpaperVersion, AuditLog |
| 10: EQR | EQR Reviewer | EngagementQualityReview, EQRFinding, AIDecision, AuditLog |

### Entity Count Summary

- **Total domain entities:** 33
- **identity_db tables:** 9
- **core_db tables:** 18 (4 system + 14 tenant-scoped)
- **trial_balance_db tables:** 4
- **workpaper_db tables:** 3
- **reporting_db tables:** 2
