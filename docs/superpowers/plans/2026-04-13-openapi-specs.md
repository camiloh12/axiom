# OpenAPI 3.1 Service Specifications — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Create production-grade OpenAPI 3.1 specifications for all 6 Axiom microservices, detailed enough to generate server stubs (oapi-codegen) and typed clients (openapi-typescript), and to serve as the enforceable contract between frontend and backend.

**Architecture:** Each service gets its own spec file in `packages/openapi/`. Specs use OpenAPI 3.1 with JSON Schema 2020-12. Common schemas (pagination, errors, header params) are defined in a shared file and referenced via `$ref`. Each spec is self-contained enough for oapi-codegen to produce a working Go server interface.

**Tech Stack:** OpenAPI 3.1 YAML, validated against the OpenAPI spec. Consumed by oapi-codegen (Go server) and openapi-typescript (frontend client).

---

## Conventions (apply to all specs)

### Authentication

The API Gateway strips the JWT and injects trusted headers. All downstream service specs declare these as required header parameters (except public/token-auth endpoints):

```yaml
parameters:
  X-User-Id:
    name: X-User-Id
    in: header
    required: true
    schema:
      type: string
      format: uuid
  X-Firm-Id:
    name: X-Firm-Id
    in: header
    required: true
    schema:
      type: string
      format: uuid
  X-User-Role:
    name: X-User-Role
    in: header
    required: true
    schema:
      type: string
      enum: [FirmAdmin, Partner, Manager, Staff, EQReviewer, ClientAdmin, ClientUser, ViewOnly]
```

### Pagination

All list endpoints use cursor-based pagination:

```yaml
PaginatedResponse:
  type: object
  required: [data, pagination]
  properties:
    data:
      type: array
      items: {}
    pagination:
      type: object
      required: [has_more]
      properties:
        next_cursor:
          type: string
          nullable: true
        has_more:
          type: boolean
```

Query parameters: `?cursor={opaque_string}&limit={1-100, default 50}`

### Error Responses

Standard error body for all 4xx/5xx:

```yaml
ErrorResponse:
  type: object
  required: [error]
  properties:
    error:
      type: object
      required: [code, message]
      properties:
        code:
          type: string
          description: Machine-readable error code (e.g., "ENGAGEMENT_NOT_FOUND")
        message:
          type: string
          description: Human-readable error description
        details:
          type: array
          items:
            type: object
            properties:
              field:
                type: string
              reason:
                type: string
```

Standard HTTP status codes used:
- `400` — validation error, malformed request
- `401` — missing/invalid auth headers
- `403` — insufficient role/permission
- `404` — resource not found
- `409` — conflict (state transition violation, duplicate)
- `422` — business rule violation (guard not met, invariant broken)
- `500` — internal server error

### UUID Format

All entity IDs use `type: string, format: uuid`.

### Timestamps

All timestamps use `type: string, format: date-time` (RFC 3339).

### Enum Naming

Enums use PascalCase values matching the PostgreSQL enum definitions in the data model spec.

---

## File Structure

```
packages/openapi/
  common.yaml              — Shared schemas, parameters, responses
  identity-service.yaml    — Identity Service (auth, firms, users, clients, invitations, templates, control objectives)
  audit-core.yaml          — Audit Core (frameworks, engagements, controls, evidence, doc requests, client hub, AI, audit log, notifications)
  trial-balance.yaml       — Trial Balance Service (TB import, accounts, adjustments, analytics, sampling, column profiles)
  workpaper-service.yaml   — Workpaper Service (workpapers, versions, review notes; WebSocket documented as x-extension)
  reporting-service.yaml   — Reporting Service (reports, versions, generation, issuance)
  doc-processing.yaml      — Document Processing Service (single POST /extract endpoint)
```

---

## Task 1: Common Schemas (`common.yaml`)

**Files:**
- Create: `packages/openapi/common.yaml`

Defines reusable components referenced by all service specs via `$ref`.

- [ ] **Step 1: Create common.yaml with shared schemas**

Write the file containing:
- `info` block with title "Axiom Common Schemas", version "1.0.0"
- **Parameter components:**
  - `X-User-Id`, `X-Firm-Id`, `X-User-Role` header parameters
  - `CursorParam` (query, optional, string)
  - `LimitParam` (query, optional, integer, default 50, min 1, max 100)
- **Schema components:**
  - `Pagination` object (next_cursor nullable string, has_more boolean)
  - `ErrorResponse` object (as defined in Conventions above)
  - `ErrorDetail` object (field string, reason string)
- **Response components:**
  - `BadRequest` (400 + ErrorResponse)
  - `Unauthorized` (401 + ErrorResponse)
  - `Forbidden` (403 + ErrorResponse)
  - `NotFound` (404 + ErrorResponse)
  - `Conflict` (409 + ErrorResponse)
  - `UnprocessableEntity` (422 + ErrorResponse)
  - `InternalServerError` (500 + ErrorResponse)

- [ ] **Step 2: Validate the YAML syntax**

Run: `npx @redocly/cli lint packages/openapi/common.yaml || echo "Install redocly for validation"`

- [ ] **Step 3: Commit**

```bash
git add packages/openapi/common.yaml
git commit -m "feat: add shared OpenAPI common schemas"
```

---

## Task 2: Identity Service Spec (`identity-service.yaml`)

**Files:**
- Create: `packages/openapi/identity-service.yaml`

**Reference docs:**
- Data model: `docs/specs/domain-and-data-model-design.md` — Contexts 1 (Firm Identity) + 3 (Firm Methodology)
- Backend arch: `docs/specs/backend-architecture-design.md` — Service 2: Identity Service
- Journeys: 1 (Firm Setup), 2 (Staff Onboarding), 3 (Engagement Scoping — template selection)

**Database:** `identity_db` — Firm, User, Client, Invitation, MethodologyTemplate, TemplateControl, TemplateTestProcedure, TemplateDocumentRequest, FirmControlObjective, FirmControlObjectiveMapping

### Endpoints to define:

**Authentication (public, no gateway headers):**
- `POST /api/v1/auth/register` — Firm signup (creates Firm + first FirmAdmin User). Request: email, password, firm_name, staff_count_range, primary_audit_types, country. Response: 201 with User + Firm + JWT tokens.
- `POST /api/v1/auth/login` — Email/password login. Request: email, password. Response: 200 with User + JWT access_token + refresh_token.
- `POST /api/v1/auth/refresh` — Refresh JWT. Request: refresh_token. Response: 200 with new access_token + refresh_token.
- `POST /api/v1/auth/verify-email` — Email verification. Request: verification_token. Response: 200.
- `POST /api/v1/auth/sso/callback` — SSO callback (OAuth Google/Microsoft). Request: provider, authorization_code. Response: 200 with User + tokens.
- `GET /api/v1/auth/jwks` — Public key for JWT verification (consumed by API Gateway). Response: 200 with JWKS.

**Invitation acceptance (public, token-auth):**
- `GET /api/v1/invitations/validate/{token}` — Validate magic link token. Response: 200 with invitation details (firm_name, assigned_role, email) or 410 if expired.
- `POST /api/v1/invitations/accept` — Accept invitation, create user. Request: token, password (or sso_provider + sso_token), display_name. Response: 201 with User + tokens.

**Firm management (authenticated):**
- `GET /api/v1/firms/current` — Get current firm profile (from X-Firm-Id). Response: Firm object.
- `PUT /api/v1/firms/current` — Update firm profile. Request: name, logo_url, timezone, billing_contact_email, settings. Response: 200 with updated Firm.

**User management (authenticated, FirmAdmin):**
- `GET /api/v1/users` — List firm users. Query: role, is_active, cursor, limit. Response: paginated User list.
- `GET /api/v1/users/{userId}` — Get user detail.
- `PUT /api/v1/users/{userId}` — Update user (admin: role, is_active; self: display_name, notification_frequency). Request: partial User fields. Response: 200 with User.
- `PUT /api/v1/users/me/profile` — Update own profile. Request: display_name, notification_frequency. Response: 200.
- `POST /api/v1/users/me/tour-complete` — Mark tour as completed. Response: 204.
- `POST /api/v1/users/{userId}/deactivate` — Soft-delete user (FirmAdmin only). Response: 200. Emits SQS event for Audit Core to revoke engagement access.

**Client management (authenticated):**
- `GET /api/v1/clients` — List clients for firm. Query: cursor, limit, search. Response: paginated Client list.
- `POST /api/v1/clients` — Create client. Request: name, industry, primary_contact_email. Response: 201.
- `GET /api/v1/clients/{clientId}` — Get client detail.
- `PUT /api/v1/clients/{clientId}` — Update client.

**Invitation management (authenticated, FirmAdmin):**
- `GET /api/v1/invitations` — List invitations for firm. Query: status, cursor, limit. Response: paginated Invitation list.
- `POST /api/v1/invitations` — Send invitations (bulk). Request: array of {email, assigned_role}. Response: 201 with created Invitation records.
- `DELETE /api/v1/invitations/{invitationId}` — Cancel invitation (only if status=Sent).
- `POST /api/v1/invitations/{invitationId}/resend` — Resend invitation email. Response: 200.

**Methodology Templates (authenticated, FirmAdmin for write):**
- `GET /api/v1/methodology-templates` — List templates. Query: applicable_engagement_type, is_active, cursor, limit. Response: paginated list with control_count, procedure_count, document_request_count summary fields.
- `POST /api/v1/methodology-templates` — Create custom template (Scale tier). Request: name, applicable_engagement_type, applicable_framework_id. Response: 201.
- `GET /api/v1/methodology-templates/{templateId}` — Get template with full structure (controls, procedures, document requests).
- `PUT /api/v1/methodology-templates/{templateId}` — Update template metadata.
- `POST /api/v1/methodology-templates/{templateId}/activate` — Activate template. Response: 200.
- `POST /api/v1/methodology-templates/{templateId}/deactivate` — Deactivate template. Response: 200.

**Template Controls (nested under template):**
- `GET /api/v1/methodology-templates/{templateId}/controls` — List template controls with nested test procedures.
- `POST /api/v1/methodology-templates/{templateId}/controls` — Add control to template.
- `PUT /api/v1/methodology-templates/{templateId}/controls/{controlId}` — Update template control.
- `DELETE /api/v1/methodology-templates/{templateId}/controls/{controlId}` — Remove control from template.

**Template Test Procedures (nested under template control):**
- `POST /api/v1/methodology-templates/{templateId}/controls/{controlId}/test-procedures` — Add procedure.
- `PUT /api/v1/methodology-templates/{templateId}/controls/{controlId}/test-procedures/{procedureId}` — Update procedure.
- `DELETE /api/v1/methodology-templates/{templateId}/controls/{controlId}/test-procedures/{procedureId}` — Remove procedure.

**Template Document Requests (nested under template):**
- `GET /api/v1/methodology-templates/{templateId}/document-requests` — List template document requests.
- `POST /api/v1/methodology-templates/{templateId}/document-requests` — Add document request template.
- `PUT /api/v1/methodology-templates/{templateId}/document-requests/{requestId}` — Update.
- `DELETE /api/v1/methodology-templates/{templateId}/document-requests/{requestId}` — Remove.

**Firm Control Objectives (authenticated):**
- `GET /api/v1/control-objectives` — List firm control objectives. Query: search, cursor, limit. Response: paginated list with mapping_count.
- `POST /api/v1/control-objectives` — Create control objective. Request: name, description, source_library_id (optional), custom_test_guidance. Response: 201.
- `GET /api/v1/control-objectives/{objectiveId}` — Get with mappings.
- `PUT /api/v1/control-objectives/{objectiveId}` — Update.
- `GET /api/v1/control-objectives/{objectiveId}/mappings` — List framework requirement mappings.
- `POST /api/v1/control-objectives/{objectiveId}/mappings` — Add mapping. Request: framework_requirement_id. Response: 201.
- `DELETE /api/v1/control-objectives/{objectiveId}/mappings/{mappingId}` — Remove mapping.

### Schemas to define:

- Firm (all attributes from data model)
- FirmCreate (register request)
- FirmUpdate (partial)
- User (all attributes, password_hash excluded)
- UserUpdate
- UserProfile (self-update subset)
- Client, ClientCreate, ClientUpdate
- Invitation, InvitationCreate (bulk: array of email+role)
- MethodologyTemplate, MethodologyTemplateSummary (with counts), MethodologyTemplateDetail (with nested controls)
- TemplateControl, TemplateControlCreate, TemplateControlUpdate
- TemplateTestProcedure, TemplateTestProcedureCreate
- TemplateDocumentRequest, TemplateDocumentRequestCreate
- FirmControlObjective, FirmControlObjectiveCreate, FirmControlObjectiveUpdate
- FirmControlObjectiveMapping
- AuthTokens (access_token, refresh_token, expires_in)
- LoginRequest, RegisterRequest, RefreshRequest
- JWKSResponse

### Steps:

- [ ] **Step 1: Create identity-service.yaml with info, servers, and security blocks**

OpenAPI 3.1 header, server URL `http://identity:8080` (internal service connect DNS), security schemes for gateway headers and public endpoints.

- [ ] **Step 2: Define all schema components**

Write every schema listed above with full property definitions, required fields, enums matching PostgreSQL types from data model, and nullable annotations. Include `example` values for key fields.

- [ ] **Step 3: Define authentication endpoints**

`/auth/register`, `/auth/login`, `/auth/refresh`, `/auth/verify-email`, `/auth/sso/callback`, `/auth/jwks`. No gateway headers required on these — they are public.

- [ ] **Step 4: Define invitation public endpoints**

`/invitations/validate/{token}`, `/invitations/accept`. Token-based auth, no gateway headers.

- [ ] **Step 5: Define firm management endpoints**

`/firms/current` GET and PUT. Requires gateway headers.

- [ ] **Step 6: Define user management endpoints**

All `/users` and `/users/me` endpoints. Document role-based access (FirmAdmin for admin operations, any role for self-profile).

- [ ] **Step 7: Define client management endpoints**

All `/clients` CRUD endpoints.

- [ ] **Step 8: Define invitation management endpoints**

All `/invitations` CRUD + resend endpoints. FirmAdmin only.

- [ ] **Step 9: Define methodology template endpoints**

All `/methodology-templates` endpoints including nested controls, test procedures, and document requests. Document read access for all firm users, write access for FirmAdmin (Scale tier for custom templates).

- [ ] **Step 10: Define firm control objective endpoints**

All `/control-objectives` CRUD + mapping endpoints.

- [ ] **Step 11: Validate the spec**

Run: `npx @redocly/cli lint packages/openapi/identity-service.yaml`

- [ ] **Step 12: Commit**

```bash
git add packages/openapi/identity-service.yaml
git commit -m "feat: add Identity Service OpenAPI spec"
```

---

## Task 3: Audit Core Spec (`audit-core.yaml`)

**Files:**
- Create: `packages/openapi/audit-core.yaml`

**Reference docs:**
- Data model: `docs/specs/domain-and-data-model-design.md` — Context 2 (Regulatory Framework) + Context 4 (Audit Core) + Cross-Cutting
- Backend arch: `docs/specs/backend-architecture-design.md` — Service 3: Audit Core
- Journeys: 3 (Engagement Scoping), 5 (Control Testing), 6 (Workpaper Review), 7 (Document Requests — auditor), 8 (Client Fulfillment), 10 (EQR)

**Database:** `core_db` — Frameworks (system), Engagements, Controls, TestProcedures, EvidenceItems, EvidenceLinks, DocumentRequests, ClientHubTokens, DelegationTokens, ClientAcceptance, EngagementQualityReview, EQRFinding, AIDecision, AuditLog, Notification

This is the largest spec. The endpoints below are grouped by aggregate.

### Endpoints to define:

**Frameworks (read-only system reference, no firm_id filter):**
- `GET /api/v1/frameworks` — List frameworks. Query: governing_body, include_deprecated. Response: array (not paginated — small static dataset).
- `GET /api/v1/frameworks/{frameworkId}` — Get framework with requirements.
- `GET /api/v1/frameworks/{frameworkId}/requirements` — List requirements. Query: category. Response: ordered array.

**Control Objective Library (read-only system reference):**
- `GET /api/v1/control-objective-library` — List library objectives. Query: search, tags. Response: paginated.
- `GET /api/v1/control-objective-library/{objectiveId}` — Get with mappings to framework requirements.

**Engagements:**
- `GET /api/v1/engagements` — List engagements for firm. Query: status, client_id, engagement_type, cursor, limit. Response: paginated EngagementSummary list (includes progress counts: controls total/complete, open doc requests).
- `POST /api/v1/engagements` — Create engagement (scaffolds controls, test procedures, workpapers from template). Request: name, engagement_type, client_id, primary_framework_id, period_start, period_end, methodology_template_id, prior_engagement_id (optional for rollforward), team_members (array of {user_id, engagement_role}), additional_framework_ids (optional). Response: 201 with Engagement + scaffold summary (control_count, procedure_count).
- `GET /api/v1/engagements/{engagementId}` — Get engagement detail with progress summary.
- `PUT /api/v1/engagements/{engagementId}` — Update engagement metadata (name, period_start, period_end — only in Planning status).
- `POST /api/v1/engagements/{engagementId}/transitions` — Execute state transition. Request: target_status (Fieldwork, Review, Reporting, Finalized, Archived), reason (optional, required for reverse transitions). Response: 200 with updated Engagement or 422 with guard failure details.

**Engagement Team:**
- `GET /api/v1/engagements/{engagementId}/team` — List team members with user details.
- `POST /api/v1/engagements/{engagementId}/team` — Add team member. Request: user_id, engagement_role. Response: 201.
- `DELETE /api/v1/engagements/{engagementId}/team/{memberId}` — Remove team member. Response: 204.

**Engagement Frameworks:**
- `GET /api/v1/engagements/{engagementId}/frameworks` — List engagement frameworks.
- `POST /api/v1/engagements/{engagementId}/frameworks` — Add secondary framework. Request: framework_id, framework_version. Response: 201.
- `DELETE /api/v1/engagements/{engagementId}/frameworks/{engagementFrameworkId}` — Remove secondary framework (not primary). Response: 204.

**Client Acceptance:**
- `GET /api/v1/engagements/{engagementId}/client-acceptance` — Get client acceptance record.
- `PUT /api/v1/engagements/{engagementId}/client-acceptance` — Update acceptance (quality_risks_identified, firm_responses, independence_confirmed). Only before signing.
- `POST /api/v1/engagements/{engagementId}/client-acceptance/sign` — Sign acceptance (Partner only). Response: 200 with signed ClientAcceptance (immutable after this).

**Engagement Quality Review (EQR):**
- `GET /api/v1/engagements/{engagementId}/quality-review` — Get EQR record.
- `POST /api/v1/engagements/{engagementId}/quality-review` — Assign EQR reviewer. Request: reviewer_id. Response: 201. Validates: reviewer has EQReviewer role AND is not an EngagementTeamMember.
- `PUT /api/v1/engagements/{engagementId}/quality-review` — Update EQR (scope_notes, conclusion — reviewer only).
- `POST /api/v1/engagements/{engagementId}/quality-review/sign-off` — Sign off EQR (reviewer only). Validates: all RequiredAction findings are Confirmed. Response: 200 with signed EQR (immutable after this).

**EQR Findings:**
- `GET /api/v1/engagements/{engagementId}/quality-review/findings` — List findings.
- `POST /api/v1/engagements/{engagementId}/quality-review/findings` — Create finding (reviewer only). Request: description, severity. Response: 201.
- `PUT /api/v1/engagements/{engagementId}/quality-review/findings/{findingId}` — Update finding description/severity (reviewer only, before confirmed).
- `POST /api/v1/engagements/{engagementId}/quality-review/findings/{findingId}/respond` — Engagement team responds. Request: team_response. Response: 200.
- `POST /api/v1/engagements/{engagementId}/quality-review/findings/{findingId}/confirm` — Reviewer confirms response. Response: 200.

**Controls:**
- `GET /api/v1/engagements/{engagementId}/controls` — List controls. Query: status, auditor_assigned_to_id, is_key_control, cursor, limit. Response: paginated Control list with linked framework requirements (resolved via control objective mappings).
- `GET /api/v1/engagements/{engagementId}/controls/{controlId}` — Get control with test procedures, evidence links, and framework requirement chain.
- `PUT /api/v1/engagements/{engagementId}/controls/{controlId}` — Update control (auditor_assigned_to_id, control_owner_id, status, description). Status transitions validated.

**Test Procedures:**
- `GET /api/v1/engagements/{engagementId}/controls/{controlId}/test-procedures` — List procedures for control.
- `GET /api/v1/engagements/{engagementId}/controls/{controlId}/test-procedures/{procedureId}` — Get procedure with evidence links.
- `PUT /api/v1/engagements/{engagementId}/controls/{controlId}/test-procedures/{procedureId}` — Update procedure fields (status, result, exceptions_noted, conclusion, population_size, sample_size, sampling_method, performed_by_id, reviewed_by_id). Status transitions validated.

**Evidence:**
- `GET /api/v1/evidence` — List evidence items for firm+client. Query: client_id (required), engagement_id (optional — filters to items linked in that engagement), search (full-text on filename + extracted_text), source_type, cursor, limit. Response: paginated EvidenceItem list.
- `POST /api/v1/evidence` — Upload evidence. Request: multipart/form-data with file, client_id, source_type. Response: 201 with EvidenceItem (extraction_status=Pending). Triggers async document.extract River job.
- `GET /api/v1/evidence/{evidenceId}` — Get evidence item detail.
- `GET /api/v1/evidence/{evidenceId}/download` — Get presigned S3 download URL. Response: 200 with {url, expires_in}.

**Evidence Links:**
- `GET /api/v1/engagements/{engagementId}/evidence-links` — List all evidence links for engagement. Query: control_id, test_procedure_id, cursor, limit.
- `POST /api/v1/engagements/{engagementId}/evidence-links` — Create link. Request: evidence_item_id, test_procedure_id, notes, ai_suggested, ai_decision_id. Response: 201 with EvidenceLink + cross-framework satisfaction summary.
- `DELETE /api/v1/engagements/{engagementId}/evidence-links/{linkId}` — Remove link (not after finalization). Response: 204.

**Document Requests:**
- `GET /api/v1/engagements/{engagementId}/document-requests` — List requests. Query: status, assigned_to_id, control_id, cursor, limit. Response: paginated list.
- `POST /api/v1/engagements/{engagementId}/document-requests` — Create request. Request: title, instructions, due_date, control_id, test_procedure_id, assigned_to_id. Response: 201.
- `POST /api/v1/engagements/{engagementId}/document-requests/bulk` — Bulk create from template. Request: methodology_template_id, customizations (optional overrides per request). Response: 201 with array of created DocumentRequests.
- `GET /api/v1/engagements/{engagementId}/document-requests/{requestId}` — Get request detail.
- `PUT /api/v1/engagements/{engagementId}/document-requests/{requestId}` — Update request (title, instructions, due_date, assigned_to_id — only before sent).
- `POST /api/v1/engagements/{engagementId}/document-requests/{requestId}/send` — Send request to client. Generates ClientHubToken if none exists. Response: 200.
- `POST /api/v1/engagements/{engagementId}/document-requests/{requestId}/accept` — Accept submitted document. Creates EvidenceLink if control/procedure linked. Response: 200.
- `POST /api/v1/engagements/{engagementId}/document-requests/{requestId}/reject` — Reject/send back with reason. Request: reason. Response: 200.

**Client Hub Tokens:**
- `POST /api/v1/engagements/{engagementId}/client-hub-tokens` — Generate new Client Hub token. Request: client_id. Response: 201 with token URL.
- `POST /api/v1/engagements/{engagementId}/client-hub-tokens/{tokenId}/revoke` — Revoke token. Response: 200.

**Client Hub (token-authenticated, NO gateway headers — public endpoints):**
- `GET /api/v1/client-hub/validate` — Validate hub token. Query: token. Response: 200 with engagement summary (name, period, firm_name) + outstanding request count.
- `GET /api/v1/client-hub/requests` — List document requests for this token's engagement. Query: token, status, cursor, limit. Response: paginated DocumentRequest list (client-facing subset: title, instructions, due_date, status).
- `POST /api/v1/client-hub/requests/{requestId}/upload` — Upload document for request. Request: multipart/form-data with file + token. Response: 201 with EvidenceItem. Changes request status to Submitted.
- `POST /api/v1/client-hub/requests/{requestId}/delegate` — Delegate request (ClientAdmin only via token). Request: token, delegate_email, custom_message. Response: 201 with DelegationToken.

**Delegation (token-authenticated, NO gateway headers — public endpoints):**
- `GET /api/v1/delegation/validate` — Validate delegation token. Query: token. Response: 200 with request title, instructions.
- `POST /api/v1/delegation/upload` — Upload document via delegation token. Request: multipart/form-data with file + token. Response: 201.

**AI Decisions:**
- `GET /api/v1/engagements/{engagementId}/ai-decisions` — List AI decisions for engagement. Query: context_type, review_action, cursor, limit. Response: paginated list.
- `GET /api/v1/ai-decisions/{decisionId}` — Get AI decision detail.
- `POST /api/v1/ai-decisions/{decisionId}/review` — Review AI decision. Request: review_action (Accepted, Modified, Rejected), accepted_value (if Modified). Response: 200.

**Audit Log:**
- `GET /api/v1/engagements/{engagementId}/audit-log` — List audit log entries for engagement. Query: action, resource_type, actor_id, cursor, limit. Response: paginated AuditLogEntry list. Read-only.
- `GET /api/v1/audit-log` — Firm-level audit log. Query: action, resource_type, actor_id, from_date, to_date, cursor, limit. Response: paginated list.

**Notifications:**
- `GET /api/v1/notifications` — List notifications for current user. Query: is_read, notification_type, cursor, limit. Response: paginated list.
- `POST /api/v1/notifications/{notificationId}/read` — Mark notification as read. Response: 204.
- `POST /api/v1/notifications/read-all` — Mark all as read. Response: 204.

### Schemas to define:

- Framework, FrameworkRequirement
- ControlObjectiveLibraryItem, ControlObjectiveLibraryMapping
- Engagement, EngagementSummary, EngagementCreate, EngagementUpdate
- EngagementTransition, TransitionGuardFailure
- EngagementTeamMember, EngagementTeamMemberCreate
- EngagementFramework, EngagementFrameworkCreate
- ClientAcceptance, ClientAcceptanceUpdate
- EngagementQualityReview, EQRCreate, EQRUpdate
- EQRFinding, EQRFindingCreate, EQRFindingResponse
- Control, ControlUpdate, ControlWithFrameworkChain
- TestProcedure, TestProcedureUpdate
- EvidenceItem, EvidenceUpload, EvidenceDownloadUrl
- EvidenceLink, EvidenceLinkCreate, CrossFrameworkSatisfaction
- DocumentRequest, DocumentRequestCreate, DocumentRequestBulkCreate, DocumentRequestClientView
- ClientHubTokenCreate, ClientHubValidation
- DelegationTokenCreate, DelegationValidation
- AIDecision, AIDecisionReview
- AuditLogEntry
- Notification

### Steps:

- [ ] **Step 1: Create audit-core.yaml with info, servers, and security blocks**
- [ ] **Step 2: Define all schema components (Framework, Engagement, Control aggregates)**
- [ ] **Step 3: Define all schema components (Evidence, Document Request, AI Decision, Audit Log, Notification)**
- [ ] **Step 4: Define framework and control objective library endpoints (read-only)**
- [ ] **Step 5: Define engagement CRUD and transition endpoints**
- [ ] **Step 6: Define engagement team and framework endpoints**
- [ ] **Step 7: Define client acceptance and EQR endpoints**
- [ ] **Step 8: Define control and test procedure endpoints**
- [ ] **Step 9: Define evidence and evidence link endpoints**
- [ ] **Step 10: Define document request endpoints (auditor-side)**
- [ ] **Step 11: Define client hub and delegation endpoints (token-auth)**
- [ ] **Step 12: Define AI decision, audit log, and notification endpoints**
- [ ] **Step 13: Validate the spec**
- [ ] **Step 14: Commit**

```bash
git add packages/openapi/audit-core.yaml
git commit -m "feat: add Audit Core OpenAPI spec"
```

---

## Task 4: Trial Balance Service Spec (`trial-balance.yaml`)

**Files:**
- Create: `packages/openapi/trial-balance.yaml`

**Reference docs:**
- Data model: `docs/specs/domain-and-data-model-design.md` — Context 5 (Trial Balance & Analytics)
- Backend arch: `docs/specs/backend-architecture-design.md` — Service 4: Trial Balance Service
- Journeys: 4 (Trial Balance Import and Analysis)

**Database:** `trial_balance_db` — TrialBalance, TrialBalanceAccount, TrialBalanceAdjustment, ColumnMappingProfile

### Endpoints to define:

**Trial Balances:**
- `GET /api/v1/engagements/{engagementId}/trial-balances` — List trial balances for engagement. Response: array (typically one per engagement, could be multiple for interim/final).
- `POST /api/v1/engagements/{engagementId}/trial-balances` — Import trial balance. Request: multipart/form-data with file (CSV/XLSX), period_date, import_source, column_mapping_profile_id (optional). Response: 201 with TrialBalance + import summary (account_count, total_debit, total_credit, difference).
- `GET /api/v1/trial-balances/{trialBalanceId}` — Get trial balance detail with summary stats.
- `DELETE /api/v1/trial-balances/{trialBalanceId}` — Delete trial balance (only if no confirmed mappings or adjustments exist). Response: 204.

**Trial Balance Accounts:**
- `GET /api/v1/trial-balances/{trialBalanceId}/accounts` — List accounts. Query: account_type, mapping_status, search, sort_by (account_number, net_balance), cursor, limit. Response: paginated list with computed net_balance.
- `GET /api/v1/trial-balances/{trialBalanceId}/accounts/{accountId}` — Get account detail with adjustments and adjusted balance.
- `PUT /api/v1/trial-balances/{trialBalanceId}/accounts/{accountId}` — Update account mapping. Request: mapped_fs_line_item, mapping_status (Confirmed or Overridden). Response: 200.
- `POST /api/v1/trial-balances/{trialBalanceId}/accounts/bulk-confirm` — Bulk confirm AI-suggested mappings. Request: account_ids (array) or min_confidence (confirm all above threshold). Response: 200 with confirmed_count.

**Trial Balance Adjustments:**
- `GET /api/v1/trial-balances/{trialBalanceId}/adjustments` — List adjustments. Query: adjustment_type, account_id, cursor, limit. Response: paginated list.
- `POST /api/v1/trial-balances/{trialBalanceId}/adjustments` — Propose adjustment. Request: account_id, amount, description, adjustment_type. Response: 201.
- `GET /api/v1/trial-balances/{trialBalanceId}/adjustments/{adjustmentId}` — Get adjustment detail.
- `PUT /api/v1/trial-balances/{trialBalanceId}/adjustments/{adjustmentId}` — Update adjustment (only before approved). Request: amount, description, adjustment_type, waiver_reason. Response: 200.
- `POST /api/v1/trial-balances/{trialBalanceId}/adjustments/{adjustmentId}/approve` — Approve adjustment (Manager/Partner only). Response: 200.

**Analytics (SQL-based, all computed server-side):**
- `GET /api/v1/trial-balances/{trialBalanceId}/analytics/summary` — Aggregated trial balance summary: total by account_type, total debits/credits, net difference, adjustment impact. Response: 200 with TrialBalanceSummary.
- `GET /api/v1/trial-balances/{trialBalanceId}/analytics/lead-schedules` — Lead schedule data grouped by mapped_fs_line_item. Response: 200 with array of LeadScheduleSection (fs_line_item, accounts, unadjusted_total, adjusted_total, prior_year_total if rollforward).
- `GET /api/v1/trial-balances/{trialBalanceId}/analytics/variance` — Period-over-period variance analysis. Query: prior_trial_balance_id (required). Response: 200 with array of VarianceItem (account, current_balance, prior_balance, variance_amount, variance_percent).
- `GET /api/v1/trial-balances/{trialBalanceId}/analytics/ratios` — Financial ratios. Response: 200 with FinancialRatios (current_ratio, quick_ratio, debt_to_equity, etc.).
- `GET /api/v1/trial-balances/{trialBalanceId}/analytics/anomalies` — AI-flagged anomalies. Response: 200 with array of AnomalyFlag (account_id, anomaly_type, description, severity).

**Population Analysis (SQL-based):**
- `GET /api/v1/trial-balances/{trialBalanceId}/population` — Get population data for sampling. Query: account_type, min_amount, max_amount. Response: paginated transaction list.
- `POST /api/v1/trial-balances/{trialBalanceId}/population/sample` — Compute sample selection. Request: sampling_method (Systematic, Random, MonetaryUnit), population_size, confidence_level, expected_misstatement_rate. Response: 200 with SampleSelection (sample_size, selected_items).
- `GET /api/v1/trial-balances/{trialBalanceId}/population/benfords-law` — Benford's law first-digit distribution. Response: 200 with array of {digit, expected_pct, actual_pct, deviation}.
- `GET /api/v1/trial-balances/{trialBalanceId}/population/gap-test` — Gap testing for sequence analysis. Query: account_type. Response: 200 with gaps found.
- `GET /api/v1/trial-balances/{trialBalanceId}/population/duplicates` — Duplicate detection. Query: threshold_amount, date_range_days. Response: 200 with potential duplicates.

**Column Mapping Profiles:**
- `GET /api/v1/column-mapping-profiles` — List profiles for firm. Response: array.
- `POST /api/v1/column-mapping-profiles` — Create profile. Request: name, accounting_system, column_mappings. Response: 201.
- `GET /api/v1/column-mapping-profiles/{profileId}` — Get profile detail.
- `PUT /api/v1/column-mapping-profiles/{profileId}` — Update profile.
- `DELETE /api/v1/column-mapping-profiles/{profileId}` — Delete profile. Response: 204.

### Schemas to define:

- TrialBalance, TrialBalanceCreate (multipart metadata), TrialBalanceImportSummary
- TrialBalanceAccount, TrialBalanceAccountUpdate, TrialBalanceAccountBulkConfirm
- TrialBalanceAdjustment, TrialBalanceAdjustmentCreate, TrialBalanceAdjustmentUpdate
- TrialBalanceSummary, LeadScheduleSection, VarianceItem, FinancialRatios, AnomalyFlag
- PopulationItem, SampleSelectionRequest, SampleSelection
- BenfordsLawResult, GapTestResult, DuplicateDetectionResult
- ColumnMappingProfile, ColumnMappingProfileCreate, ColumnMappingProfileUpdate

### Steps:

- [ ] **Step 1: Create trial-balance.yaml with info, servers, and security blocks**
- [ ] **Step 2: Define all schema components**
- [ ] **Step 3: Define trial balance CRUD and import endpoints**
- [ ] **Step 4: Define account listing, update, and bulk-confirm endpoints**
- [ ] **Step 5: Define adjustment CRUD and approval endpoints**
- [ ] **Step 6: Define analytics endpoints (summary, lead schedules, variance, ratios, anomalies)**
- [ ] **Step 7: Define population analysis endpoints (sampling, Benford's, gap test, duplicates)**
- [ ] **Step 8: Define column mapping profile CRUD endpoints**
- [ ] **Step 9: Validate the spec**
- [ ] **Step 10: Commit**

```bash
git add packages/openapi/trial-balance.yaml
git commit -m "feat: add Trial Balance Service OpenAPI spec"
```

---

## Task 5: Workpaper Service Spec (`workpaper-service.yaml`)

**Files:**
- Create: `packages/openapi/workpaper-service.yaml`

**Reference docs:**
- Data model: `docs/specs/domain-and-data-model-design.md` — Context 6 (Workpaper Authoring)
- Backend arch: `docs/specs/backend-architecture-design.md` — Service 5: Workpaper Service
- Journeys: 5 (Control Testing — workpaper creation), 6 (Workpaper Review)

**Database:** `workpaper_db` — Workpaper, WorkpaperVersion, ReviewNote

### Endpoints to define:

**Workpapers:**
- `GET /api/v1/engagements/{engagementId}/workpapers` — List workpapers. Query: status, workpaper_type, control_id, prepared_by_id, cursor, limit. Response: paginated WorkpaperSummary list.
- `POST /api/v1/engagements/{engagementId}/workpapers` — Create workpaper. Request: title, workpaper_type, control_id (optional), content (optional initial content). Response: 201.
- `GET /api/v1/workpapers/{workpaperId}` — Get workpaper with current content.
- `PUT /api/v1/workpapers/{workpaperId}` — Update workpaper content. Request: content (jsonb structured rich text). Automatically creates a new WorkpaperVersion. Clears is_ai_draft flag on latest version if AI draft existed. Response: 200. Returns 422 if is_locked=true.
- `POST /api/v1/workpapers/{workpaperId}/submit-for-review` — Submit for review. Validates: latest WorkpaperVersion.is_ai_draft must be false. Changes status: Draft → PreparedPendingReview. Locks workpaper for preparer. Response: 200 or 422 if AI draft not edited.
- `POST /api/v1/workpapers/{workpaperId}/start-review` — Manager opens review. Changes status: PreparedPendingReview → InReview. Response: 200.
- `POST /api/v1/workpapers/{workpaperId}/return-to-preparer` — Return with review notes. Changes status: InReview → ReviewNotesOpen. Response: 200.
- `POST /api/v1/workpapers/{workpaperId}/resubmit` — Staff resubmits after addressing notes. Changes status: ReviewNotesOpen → InReview. Response: 200.
- `POST /api/v1/workpapers/{workpaperId}/complete-review` — Manager clears review. Validates: all ReviewNotes resolved. Changes status: InReview → ReviewComplete. Response: 200 or 422 if open notes exist.
- `POST /api/v1/workpapers/{workpaperId}/sign-off` — Partner final sign-off. Changes status: ReviewComplete → SignedOff. Creates AuditLog entry (via Audit Core REST call). Response: 200.
- `POST /api/v1/workpapers/{workpaperId}/lock` — Lock workpaper (engagement finalization). Sets is_locked=true. Response: 200.
- `POST /api/v1/workpapers/{workpaperId}/addendum` — Create post-finalization addendum. Request: content, addendum_reason. Creates WorkpaperVersion with is_addendum=true. Requires Partner role. Response: 201.

**Workpaper Versions:**
- `GET /api/v1/workpapers/{workpaperId}/versions` — List version history. Response: array of WorkpaperVersion (version_number, saved_by_id, saved_at, is_ai_draft, is_addendum).
- `GET /api/v1/workpapers/{workpaperId}/versions/{versionNumber}` — Get specific version content.

**Review Notes:**
- `GET /api/v1/workpapers/{workpaperId}/review-notes` — List review notes. Query: status, severity, cursor, limit. Response: paginated list.
- `POST /api/v1/workpapers/{workpaperId}/review-notes` — Create review note (Manager/Partner only). Request: description, severity, content_anchor. Response: 201.
- `GET /api/v1/workpapers/{workpaperId}/review-notes/{noteId}` — Get review note detail.
- `POST /api/v1/workpapers/{workpaperId}/review-notes/{noteId}/respond` — Staff responds to note. Request: response. Changes note status: Open → Responded. Response: 200.
- `POST /api/v1/workpapers/{workpaperId}/review-notes/{noteId}/resolve` — Reviewer resolves note. Changes note status: Responded → Resolved. Response: 200.

**WebSocket (Yjs Real-Time Sync) — documented as x-extension:**
- `WS /ws/workpapers/{workpaperId}` — Yjs awareness and document sync protocol. Documented via `x-websocket` extension with message format description. Auth: token query parameter derived from JWT.

### Schemas to define:

- Workpaper, WorkpaperSummary, WorkpaperCreate, WorkpaperUpdate
- WorkpaperVersion, WorkpaperAddendum
- ReviewNote, ReviewNoteCreate, ReviewNoteResponse
- WorkpaperStatusTransition (describes valid transitions)

### Steps:

- [ ] **Step 1: Create workpaper-service.yaml with info, servers, and security blocks**
- [ ] **Step 2: Define all schema components**
- [ ] **Step 3: Define workpaper CRUD endpoints**
- [ ] **Step 4: Define workpaper status transition endpoints (submit, review, sign-off, lock, addendum)**
- [ ] **Step 5: Define workpaper version endpoints**
- [ ] **Step 6: Define review note endpoints**
- [ ] **Step 7: Document WebSocket endpoint via x-extension**
- [ ] **Step 8: Validate the spec**
- [ ] **Step 9: Commit**

```bash
git add packages/openapi/workpaper-service.yaml
git commit -m "feat: add Workpaper Service OpenAPI spec"
```

---

## Task 6: Reporting Service Spec (`reporting-service.yaml`)

**Files:**
- Create: `packages/openapi/reporting-service.yaml`

**Reference docs:**
- Data model: `docs/specs/domain-and-data-model-design.md` — Context 7 (Reporting & Archival)
- Backend arch: `docs/specs/backend-architecture-design.md` — Service 6: Reporting Service
- Journeys: 9 (Reporting, Finalize, Archive)

**Database:** `reporting_db` — Report, ReportVersion

### Endpoints to define:

**Reports:**
- `GET /api/v1/engagements/{engagementId}/reports` — List reports for engagement. Response: array of ReportSummary.
- `POST /api/v1/engagements/{engagementId}/reports` — Create report (triggers async generation via River). Request: report_type, template_id (optional). Response: 202 Accepted with Report (status=Draft). River job assembles data from Audit Core, Trial Balance, and Workpaper services.
- `GET /api/v1/reports/{reportId}` — Get report with current content and status.
- `PUT /api/v1/reports/{reportId}` — Update report content. Request: content. Creates new ReportVersion. Response: 200. Returns 422 if report is Issued or Archived.
- `POST /api/v1/reports/{reportId}/share-with-client` — Share draft with client for review. Changes status: Draft → ClientReview. Response: 200.
- `POST /api/v1/reports/{reportId}/complete-client-review` — Mark client review complete. Changes status: ClientReview → FirmReview. Response: 200.
- `POST /api/v1/reports/{reportId}/issue` — Issue final report (Partner only). Changes status: FirmReview or Draft → Issued. Records issued_at, issued_by_id. Triggers engagement assembly_deadline and retention_deadline computation (via Audit Core REST call). Response: 200. Returns confirmation with computed deadlines.
- `GET /api/v1/reports/{reportId}/download` — Get presigned S3 download URL for rendered report (PDF). Response: 200 with {url, expires_in}.

**Report Versions:**
- `GET /api/v1/reports/{reportId}/versions` — List version history. Response: array of ReportVersion.
- `GET /api/v1/reports/{reportId}/versions/{versionNumber}` — Get specific version content.

**Engagement Export (available at any engagement status):**
- `POST /api/v1/engagements/{engagementId}/export` — Request full engagement export. Response: 202 Accepted with export_id. Async River job assembles ZIP (workpapers as PDF, evidence in native format, TB as XLSX, AuditLog as CSV, metadata as JSON).
- `GET /api/v1/exports/{exportId}` — Get export status and download URL. Response: 200 with {status: Pending|Processing|Complete|Failed, download_url, expires_at}.

### Schemas to define:

- Report, ReportSummary, ReportCreate, ReportUpdate
- ReportVersion
- ReportIssuanceConfirmation (includes computed deadlines)
- ReportDownload
- EngagementExportRequest, EngagementExportStatus

### Steps:

- [ ] **Step 1: Create reporting-service.yaml with info, servers, and security blocks**
- [ ] **Step 2: Define all schema components**
- [ ] **Step 3: Define report CRUD endpoints**
- [ ] **Step 4: Define report lifecycle endpoints (share, review, issue, download)**
- [ ] **Step 5: Define report version endpoints**
- [ ] **Step 6: Define engagement export endpoints**
- [ ] **Step 7: Validate the spec**
- [ ] **Step 8: Commit**

```bash
git add packages/openapi/reporting-service.yaml
git commit -m "feat: add Reporting Service OpenAPI spec"
```

---

## Task 7: Document Processing Service Spec (`doc-processing.yaml`)

**Files:**
- Create: `packages/openapi/doc-processing.yaml`

**Reference docs:**
- Backend arch: `docs/specs/backend-architecture-design.md` — Service 7: Document Processing Service

**Database:** None (stateless)

This is a single-endpoint internal service. Not exposed through the API Gateway — called exclusively by Audit Core's `document.extract` River worker via HTTP within the VPC.

### Endpoints to define:

- `POST /extract` — Extract text from a document. Request: multipart/form-data with file (PDF, image). Response: 200 with ExtractionResult (text, pages array, metadata: page_count, has_tables, is_scanned).

### Schemas to define:

- ExtractionResult
- ExtractionPage (page_number, text, tables)
- ExtractionMetadata (page_count, has_tables, is_scanned)

### Steps:

- [ ] **Step 1: Create doc-processing.yaml with info and servers**

No security block — internal service, no JWT/gateway headers. Server URL: `http://doc-processing:8080`.

- [ ] **Step 2: Define schemas**
- [ ] **Step 3: Define POST /extract endpoint**
- [ ] **Step 4: Validate the spec**
- [ ] **Step 5: Commit**

```bash
git add packages/openapi/doc-processing.yaml
git commit -m "feat: add Document Processing Service OpenAPI spec"
```

---

## Dependency Graph

```
Task 1 (common.yaml) ──► Task 2 (identity)
                     ──► Task 3 (audit-core)
                     ──► Task 4 (trial-balance)
                     ──► Task 5 (workpaper)
                     ──► Task 6 (reporting)
                     ──► Task 7 (doc-processing, no common dependency)
```

Tasks 2–7 are independent of each other and can be executed in parallel after Task 1 completes. Task 7 does not depend on common.yaml (internal service, no shared auth headers).

---

## Validation Checklist (run after all tasks)

- [ ] All specs pass `npx @redocly/cli lint packages/openapi/*.yaml`
- [ ] Every entity from the data model has a corresponding schema in the appropriate spec
- [ ] Every user journey action maps to at least one API endpoint
- [ ] All enum values match PostgreSQL enum definitions in the data model
- [ ] Cross-service references are documented (e.g., Workpaper Service notes that engagement_id is validated via Audit Core REST call)
- [ ] Error responses are consistent across all specs
- [ ] Pagination is consistent across all list endpoints
