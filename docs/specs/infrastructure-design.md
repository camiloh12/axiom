# Axiom Infrastructure Design
**Date:** 2026-04-12
**Status:** Approved
**Builds on:** `axiom-spec-design.md` Section 7, `backend-architecture-design.md`, `domain-and-data-model-design.md`

---

## 1. Context and Scope

This document specifies the complete AWS infrastructure for the Axiom platform — every service, how it is configured, how services integrate, and how the infrastructure is provisioned and managed via Terraform. It is intended to be directly translatable into Terraform workspaces.

The backend architecture design defines two deployed services — a Go modular monolith (Axiom API) containing all domain modules (identity, audit core, control mapping, workpaper, reporting, provenance) and a stateless Python service (Document Processing) for PDF extraction — sharing a single database. This document specifies the AWS infrastructure that hosts and connects those services.

### Key Decisions Made in This Document

1. **Multi-account via AWS Organizations** — separate accounts for management, tooling, dev, staging, and production. Aligns with SOC 2 Type 2 and ISO 27001 compliance expectations from day one.
2. **GitHub Actions with OIDC federation** — no long-lived AWS credentials. OIDC trust policies scope deployment access by branch and environment.
3. **Route 53 for all DNS** — domain registration and hosted zones within AWS. Eliminates external DNS vendor dependency from the SOC 2 sub-processor list.
4. **S3 + DynamoDB for Terraform state** — stored in the tooling account. No additional vendors (Terraform Cloud) in the compliance scope.
5. **AWS-native observability** — CloudWatch Logs, Metrics, Dashboards, Alarms, and X-Ray. No third-party observability vendor at launch.
6. **Secrets Manager with automatic rotation** — 30-day RDS credential rotation via AWS-provided Lambda. ACM auto-renewal for TLS certificates.
7. **Enhanced backup strategy** — 35-day RDS backup retention, S3 cross-region replication for evidence files in production. Supports 99.9% uptime SLA and regulatory retention obligations.
8. **Standard 2-tier VPC** — public subnets (ALB, NAT) and private subnets (ECS, RDS). VPC endpoints for high-traffic AWS services.
9. **RDS PostgreSQL over Aurora** — cost-effective at launch scale, faster extension support (pgvector), simpler operational model. Aurora migration criteria documented.
10. **Layer-based Terraform workspaces** — infrastructure split by layer (network, data, compute, etc.) with a documented migration path to per-service compute workspaces when team growth requires it.

---

## 2. AWS Account Structure

### AWS Organizations

Five accounts organized under a single Organization:

| Account | OU | Purpose | Workloads |
|---|---|---|---|
| `axiom-management` | Root | Organizations root, consolidated billing, SCPs, IAM Identity Center (SSO), CloudTrail organization trail | None |
| `axiom-tooling` | Infrastructure | Terraform state (S3 + DynamoDB), ECR repositories, GitHub Actions OIDC identity provider, CI/CD IAM roles | CI/CD only |
| `axiom-dev` | Workloads | Development, experimentation, POCs. Safe to break. Single AZ, minimal scaling. | Full app stack (reduced) |
| `axiom-staging` | Workloads | Integration testing, QA, pre-production validation, sales demos. Multi-AZ like prod. | Full app stack (reduced) |
| `axiom-prod` | Production | Production workloads, customer data, customer-facing endpoints. | Full app stack |

### Service Control Policies

SCPs applied at the OU level enforce security invariants that no individual IAM policy can override:

**All accounts (Organization root SCP):**
- Deny all regions except `us-east-1` and global services (IAM, Route 53, CloudFront, STS, Organizations)
- Deny `cloudtrail:StopLogging`, `cloudtrail:DeleteTrail`
- Deny `organizations:LeaveOrganization`

**Production OU:**
- Deny `iam:CreateUser` with `iam:CreateLoginProfile` (force OIDC/SSO, no long-lived credentials)
- Deny `s3:PutPublicAccessBlock` with any `BlockPublicAccess` parameter set to `false`
- Deny `s3:DeleteObjectLockConfiguration`
- Deny `rds:ModifyDBInstance` with `PubliclyAccessible = true`

**Management account:**
- Deny launching workload resources: `ec2:RunInstances`, `ecs:CreateCluster`, `rds:CreateDBInstance`

### IAM Identity Center (SSO)

Configured in the management account for human console access to all accounts.

**Permission sets:**

| Permission set | Access level | Usage |
|---|---|---|
| `AdministratorAccess` | Full admin | Break-glass only. Requires MFA. Used for incident response and infrastructure recovery. |
| `ReadOnlyAccess` | Read all resources | Daily operations — investigating issues, reviewing configuration, checking metrics. |
| `DatabaseAccess` | RDS query access | Custom permission set. Allows `rds-db:connect` for interactive database debugging via RDS IAM authentication. |

### Demo Stage Account Provisioning

Before revenue, four accounts are provisioned: `axiom-management` (Organizations overhead only, no workloads), `axiom-tooling` (CI/CD only), `axiom-dev`, and `axiom-staging`. The `axiom-prod` account is provisioned when the first paying customer is onboarded.

Dev and staging run identical Terraform modules with identical sizing at demo stage. The distinction is operational discipline: dev runs feature branches and gets wiped freely; staging stays on the main branch with curated demo data.

---

## 3. Network Architecture

### VPC Design

One VPC per workload account (`axiom-dev`, `axiom-staging`, `axiom-prod`). Two-tier subnet layout: public and private.

**CIDR allocation (non-overlapping for future VPC peering):**

| Account | VPC CIDR |
|---|---|
| `axiom-dev` | `10.0.0.0/16` |
| `axiom-staging` | `10.1.0.0/16` |
| `axiom-prod` | `10.2.0.0/16` |

**Subnet layout (prod/staging — 2 AZs):**

| Subnet tier | AZ-a | AZ-b | Purpose |
|---|---|---|---|
| Public | `10.2.0.0/25` | `10.2.0.128/25` | ALB, NAT Gateways |
| Private | `10.2.16.0/20` | `10.2.32.0/20` | ECS Fargate tasks, RDS, VPC endpoint ENIs |

Dev uses the same two tiers in a single AZ with one NAT Gateway. Private subnets retain `/20` because ECS Fargate allocates an ENI per task — at scale with multiple services and autoscaling, IP consumption grows quickly. The private CIDR has room to expand to 3 AZs (`10.2.48.0/20` for AZ-c) without re-addressing.

### NAT Gateways

- Prod and staging: one per AZ (high availability — AZ-a failure does not affect AZ-b outbound traffic)
- Dev: single NAT Gateway

### VPC Endpoints

**Gateway endpoints (free):**

| Endpoint | Justification |
|---|---|
| S3 | Evidence file uploads/downloads — highest-volume data path. Avoids NAT data processing charges. |
| DynamoDB | Terraform state locking if accessed cross-account. |

**Interface endpoints (~$7.50/month each per AZ):**

| Endpoint | Justification |
|---|---|
| `ecr.api` + `ecr.dkr` | ECS image pulls stay on AWS backbone, avoid NAT costs |
| `secretsmanager` | Credential retrieval at task startup — high frequency, security-sensitive |
| `bedrock-runtime` | AI inference traffic never leaves AWS network — required by spec for compliance |
| `logs` + `monitoring` | CloudWatch log and metric shipping — high volume |
| `xray` | Trace data submission |
| `states` | Step Functions execution from Audit Core |

SES, ACM validation, and GitHub Actions webhook callbacks route through NAT Gateway — low volume, not worth dedicated endpoints.

### Security Groups

| Security group | Inbound | Outbound |
|---|---|---|
| `sg-alb` | 443 from `0.0.0.0/0` | All to `sg-ecs` |
| `sg-ecs` | 8080 from `sg-alb`, 8080 from `sg-ecs` (Axiom API → doc-processing) | 5432 to `sg-rds`, 443 to VPC endpoints, 443 to NAT GW |
| `sg-rds` | 5432 from `sg-ecs` only | None |

With the modular monolith, there is no separate gateway security group — JWT verification is middleware within the Axiom API process. The self-referencing `sg-ecs` rule allows the Axiom API to call Document Processing via ECS Service Connect.

---

## 4. Data Layer

### RDS PostgreSQL

One RDS instance per environment hosting a single logical database as specified in `backend-architecture-design.md`:

| Database | Owner service | RLS |
|---|---|---|
| `axiom_db` | Axiom API (all modules) | Yes (`firm_id` RLS on all tenant-scoped tables) |

**Instance configuration by environment:**

| Parameter | Demo (dev/staging) | Staging (full) | Prod |
|---|---|---|---|
| Engine version | PostgreSQL 18.x | PostgreSQL 18.x | PostgreSQL 18.x |
| Instance class | `db.t4g.medium` | `db.r7g.large` | `db.r7g.xlarge` |
| Multi-AZ | No | Yes | Yes |
| Storage | 20 GB gp3 | 100 GB gp3 | 200 GB gp3, autoscaling to 1 TB |
| Backup retention | 1 day | 7 days | 35 days |
| Point-in-time recovery | No | Yes | Yes |
| Deletion protection | No | Yes | Yes |
| Performance Insights | Off | On (7-day free tier) | On (7-day, extend to paid before SOC 2 audit) |

**Extensions enabled:** `pgvector` on `axiom_db` (embeddings for evidence items, framework requirements, control objective library, and firm control objectives — see AI Architecture Design Section 2), `pg_stat_statements`.

**Parameter group customizations:**
- `rds.force_ssl = 1` — enforce TLS for all connections
- `shared_preload_libraries = pg_stat_statements,pgvector`
- `log_min_duration_statement = 1000` — log queries exceeding 1 second

**Authentication:** The Axiom API connects with a single database user (`axiom_svc`). Credentials stored in Secrets Manager with 30-day automatic rotation via the AWS-provided Lambda rotator. A `master` user credential is stored separately for migrations and break-glass access.

**Connection pooling:** PgBouncer deployed as an ECS sidecar container on the Axiom API task, running in transaction mode. The application connects to PgBouncer at `localhost:6432`, PgBouncer connects to RDS.

### Aurora Migration Criteria

Standard RDS PostgreSQL is chosen over Aurora for launch because:

- Aurora minimum cost is roughly 2x standard RDS for equivalent instance sizes
- Aurora Serverless v2 charges even at idle (0.5 ACU minimum, ~$43/month)
- pgvector new versions land on standard RDS before Aurora
- Launch-scale storage (< 500 GB) does not benefit from Aurora's distributed storage layer

**Reconsider Aurora when:**
- Read replicas are needed with near-zero replication lag (Aurora replicas share the storage layer)
- Storage exceeds 500 GB and IOPS become a bottleneck on gp3
- The cost premium is justified by operational benefits (automatic storage scaling, faster failover)

The Terraform `rds` module parameterizes the engine (`engine = var.rds_engine`) so switching from `postgres` to `aurora-postgresql` is a configuration change, not a module rewrite. A migration requires: snapshot the RDS instance, restore to an Aurora cluster, update ECS task environment variables to the new endpoint, and validate.

### S3 Buckets

| Bucket | Purpose | Configuration |
|---|---|---|
| `axiom-{env}-evidence` | Evidence file uploads (`EvidenceItem.storage_path`). Carries cryptographically signed artifacts (signed screenshots, hashed DOM snapshots, provenance manifests) per the compliance provenance model. | S3 Object Lock in COMPLIANCE mode for signed artifacts (WORM — artifact immutability is central to auditor-defensibility). SSE-S3 default encryption. HIPAA-flagged objects uploaded with per-request SSE-KMS using the `axiom-{env}-hipaa` key. Versioning enabled. Lifecycle: IA after 90 days, Glacier after 1 year (subject to Object Lock retention). Cross-region replication to `us-west-2` (prod only). |
| `axiom-{env}-archive` | WORM storage for finalized engagements (SOC 2 reports, ISO certs, ISO 42001 certs, PCI ROC/AOC, HIPAA attestations, HITRUST r2 reports) | S3 Object Lock in COMPLIANCE mode. Retention set per-object using `retention_deadline`. SSE-KMS with `axiom-{env}-hipaa` key. No lifecycle transitions — objects remain until retention expires then auto-delete. |
| `axiom-{env}-reports` | Generated report PDFs (SOC 2 Type 1/2, ISO 27001/27701/42001 audit reports and cert packages, PCI ROC/AOC, HIPAA/HITRUST r2 attestations) | S3 Object Lock in COMPLIANCE mode for finalized reports. SSE-S3. Versioning enabled. Lifecycle: IA after 30 days (subject to Object Lock retention). |
| `axiom-{env}-scf-catalog` | SCF (Secure Controls Framework) catalog seed data — source-of-truth crosswalks for platform-level `CommonControl` entities. OSCAL, AICPA, and CIS mapping imports land here before ingestion. | SSE-S3. Versioning enabled. Read-mostly. Private — accessed by the `scf-import` scheduled job via VPC endpoint. Lifecycle: retain all versions (small, auditable provenance of catalog changes). |
| `axiom-{env}-spa` | React SPA static assets | SSE-S3. Versioning enabled. Private — accessed only via CloudFront OAC. |

**SSE-KMS for HIPAA evidence explained:** SSE-S3 uses an AWS-managed key with no CloudTrail logging of decrypt calls and no IAM-level key access control. SSE-KMS provides: CloudTrail logging of every decrypt call (HIPAA §164.312(b) audit controls), key policy restricting decrypt to specific IAM roles (defense in depth beyond S3 bucket policy), and customer-managed annual key rotation (HIPAA §164.312(a)(2)(iv)). The application layer checks `Engagement.engagement_type` to determine which encryption to apply at upload time.

**Bucket policies (all buckets):**
- Deny non-TLS requests (`aws:SecureTransport = false`)
- Deny public access (enforced at both bucket and account level)
- Archive, evidence, and reports buckets: additional policy denying `s3:DeleteObject` and `s3:PutObjectRetention` with `bypass-governance-retention` — defense in depth on top of Object Lock COMPLIANCE mode
- SCF catalog bucket: read access restricted to the `scf-import` task role; write access restricted to the platform admin role used for curated catalog uploads

### SQS Queues

No SQS queues are required. In the previous microservices design, three SQS queues handled cross-service events (document uploaded, user deactivated, FirmControlObjective created). In the modular monolith, these are handled via direct River job enqueue or function calls within the same process.

### Secrets Manager

| Secret path | Rotation | Used by |
|---|---|---|
| `axiom/{env}/rds/axiom-svc` | 30-day auto | Axiom API |
| `axiom/{env}/rds/master` | 30-day auto | Migrations, break-glass |
| `axiom/{env}/jwt/rsa-private-key` | Manual (rotate on compromise) | Axiom API (identity module) |
| `axiom/{env}/jwt/rsa-public-key` | Matches private key | Axiom API (gateway middleware, loaded at startup) |
| `axiom/{env}/ses/smtp-credentials` | Manual | Axiom API (notification worker) |
| `axiom/{env}/oauth/{provider}` | Manual | Axiom API (identity module, SSO) |

RDS credential rotation uses the AWS-provided Lambda rotation function (`SecretsManagerRDSPostgreSQLRotationSingleUser`). The rotation Lambda runs in the same VPC private subnets as the ECS tasks, with a security group allowing 5432 to `sg-rds`.

### KMS Keys

| Key alias | Account | Purpose | Rotation |
|---|---|---|---|
| `axiom-{env}-default` | Each workload account | Default encryption for CloudWatch Logs, SSM parameters, SNS topics | Annual automatic |
| `axiom-{env}-hipaa` | Each workload account | S3 SSE-KMS for HIPAA evidence and archive buckets | Annual automatic |
| `axiom-{env}-rds` | Each workload account | RDS instance encryption at rest | Annual automatic |
| `axiom-{env}-provenance-signing` | Each workload account | **Asymmetric (`ECC_NIST_P256`, `SIGN_VERIFY`).** Used by the provenance signer to sign evidence artifacts, `AIDecision` ledger entries, and finalized report manifests. Key material is non-exportable — only signing and verification via KMS API. Public key exposed for downstream auditor verification. | Manual rotation (new key version on schedule; signatures remain verifiable under prior versions) |
| `axiom-cloudtrail` | `axiom-management` | CloudTrail log file encryption | Annual automatic |
| `axiom-terraform-state` | `axiom-tooling` | Terraform state bucket encryption | Annual automatic |

Key policies follow least privilege. The HIPAA key policy grants `kms:GenerateDataKey` and `kms:Decrypt` only to the Axiom API task role — no other service can encrypt or decrypt HIPAA evidence. The provenance-signing key grants `kms:Sign` only to the `provenance-signer` task role (isolated) and `kms:Verify` / `kms:GetPublicKey` to the `axiom-api` task role for verification flows. No role has `kms:Decrypt` or export permissions on the signing key — the key material cannot leave HSM boundary.

---

## 5. Compute Layer

### ECS Cluster

One ECS cluster per environment using Fargate capacity provider (no EC2 instances).

- Container Insights enabled (per-task CPU/memory/network metrics in CloudWatch)
- ECS Service Connect enabled with `axiom-{env}` namespace for internal service discovery (e.g., `http://audit-core:8080` within the VPC)
- Execute Command enabled on dev and staging for debugging (`aws ecs execute-command`). Disabled on prod by default — enable via break-glass IAM role only.

### Application Load Balancer

One ALB per environment in the public subnets.

| Listener | Action |
|---|---|
| HTTPS 443 | Forward to API Gateway target group. ACM certificate for `api.{env}.axiom.com` (prod: `api.axiom.com`). TLS policy: `ELBSecurityPolicy-TLS13-1-2-2021-06`. |
| HTTP 80 | Redirect to HTTPS 443. |

Single target group pointing to the Axiom API ECS service. All routing (REST and WebSocket) is handled internally by Chi within the Go binary.

**Health check:** `GET /healthz` on the Axiom API, 30-second interval, 3 consecutive failures before deregistration.

**WebSocket support:** ALB natively supports WebSocket connections. Workpaper WebSocket traffic is served at `/api/v1/workpapers/ws/*` by the same Axiom API binary. ALB idle timeout set to 3600 seconds (1 hour) for WebSocket connection persistence.

### ECS Services and Task Definitions

| Service | CPU | Memory | Min tasks (prod) | Max tasks (prod) | Scaling metric | Port |
|---|---|---|---|---|---|---|
| `axiom-api` | 1024 | 2048 MB | 2 | 8 | Max of: CPU > 60%, active WebSocket connections > threshold | 8080 |
| `doc-processing` | 1024 | 2048 MB | 1 | 4 | River job queue depth (`auditcore.document-extract`) | 8000 |
| `provenance-signer` | 512 | 1024 MB | 2 | 4 | River job queue depth (`provenance.sign`) | 8090 |

**Dev and demo-stage overrides:** `axiom-api` runs min 1 / max 2 tasks with 512 CPU / 1024 MB. `doc-processing` runs min 1 / max 2 tasks with 512 CPU / 1024 MB. `provenance-signer` runs min 1 / max 2 tasks with 256 CPU / 512 MB.

**Isolation rationale for `provenance-signer`:** kept as a dedicated ECS service with its own task role so that `kms:Sign` on the provenance-signing key is restricted to a narrow, minimal-dependency container. A signing-specific blast radius simplifies ISO 42001 and SOC 2 evidence-integrity claims. Internal callers reach the signer over ECS Service Connect at `http://provenance-signer:8090` within `sg-ecs`.

**Task definition configuration (all services):**
- Fargate platform version: `LATEST`
- Each task has its own IAM task role (least-privilege per service) and a shared execution role (ECR pull, CloudWatch Logs, Secrets Manager read)
- Log driver: `awslogs` to CloudWatch log group `/axiom/{env}/{service-name}`
- Log retention: dev 7 days, staging 30 days, prod 365 days
- Environment variables: non-sensitive via ECS task definition `environment` block; sensitive via Secrets Manager ARN references in `secrets` block
- PgBouncer sidecar container on the `axiom-api` task (256 CPU / 256 MB). Not present on `doc-processing`.

**Task IAM roles (per-service least privilege):**

| Service | Permissions |
|---|---|
| `axiom-api` | Secrets Manager read (DB creds, JWT keys, OAuth secrets, SES creds), S3 read/write to `evidence` and `archive` buckets, S3 write to `reports`, S3 read to `scf-catalog`, SES `SendEmail`, Bedrock `InvokeModel` (Haiku + Sonnet), Step Functions `StartExecution`, KMS `Decrypt`/`GenerateDataKey` (HIPAA key), KMS `Verify` + `GetPublicKey` (provenance-signing key — verification only) |
| `doc-processing` | None (stateless — receives files via HTTP, returns extracted text) |
| `provenance-signer` | KMS `Sign` + `GetPublicKey` on `axiom-{env}-provenance-signing` (no other KMS permissions). S3 read on `evidence` (to hash and sign in place), S3 write of sidecar `.sig` and manifest objects to `evidence` and `reports` with `object-lock-mode = COMPLIANCE`. CloudWatch Logs write. No Bedrock, no RDS access. |
| `scf-import` (scheduled job, runs on `axiom-api` task role or a dedicated short-lived Fargate task) | S3 read on `scf-catalog` bucket, RDS write access via `axiom-svc` credentials (ingests `CommonControl`, `FrameworkRequirement`, and STRM edges). Triggered by EventBridge schedule (quarterly aligned with SCF releases, configurable). |

### ECS Deployment Configuration

- **Deployment type:** rolling update
- **Minimum healthy percent:** 100% (never drop below current task count during deploy)
- **Maximum percent:** 200% (new tasks start before old tasks drain)
- **Circuit breaker:** enabled with automatic rollback — if new tasks fail health checks, ECS rolls back to the previous task definition
- **Deployment order:** `doc-processing` and `provenance-signer` first (no intra-cluster dependencies), then `axiom-api`

### Step Functions

Two state machines per environment:

| State machine | Trigger | Execution type |
|---|---|---|
| `EngagementLifecycleStateMachine` | Audit Core via AWS SDK on engagement status change | Standard (long-running, days to years) |
| `DocumentRequestReminderStateMachine` | Audit Core via AWS SDK on document request creation | Standard |

IAM role per state machine allows: invoking specific ECS tasks (guard condition checks), SQS publish (notifications), CloudWatch Logs (execution logging). State machine definitions authored as ASL JSON files in the monorepo, deployed via Terraform.

Not provisioned at demo stage — engagement lifecycle transitions managed manually until the flow is implemented.

### Bedrock

Model access enabled in `us-east-1`:
- `anthropic.claude-haiku-4-5` — evidence-to-control mapping (Feature 2), framework version migration (Feature 3), evidence link suggestion (Feature 5), drift detection / continuous assurance (Feature 7)
- `anthropic.claude-sonnet-4-6` — document completeness review (Feature 1), workpaper / audit-file narrative drafts (Feature 4), gap analysis and risk reasoning (Feature 6), findings triage and management-response drafting (Feature 8)

Traffic routes through the Bedrock VPC endpoint. The `axiom-api` task role grants access to both models — all AI features run within the single binary. Model selection per feature is enforced at the application layer (the `internal/ai` package routes each feature to its designated model).

| Model | Features |
|---|---|
| Haiku | Evidence-to-control mapping (2), framework version migration (3), evidence link suggestion (5), drift detection (7) |
| Sonnet | Completeness review (1), workpaper drafts (4), gap analysis (6), findings triage + management-response drafting (8) |

On-demand pricing at launch. Monitor `InvocationLatency` and `ThrottlingCount` CloudWatch metrics. If throttling exceeds 5% of requests, request a quota increase or evaluate provisioned throughput.

---

## 6. DNS, CDN, and TLS

### Route 53

**Hosted zones:**

| Hosted zone | Account | Purpose |
|---|---|---|
| `axiom.com` | `axiom-tooling` | Root zone, NS delegation records to environment subdomains |
| `dev.axiom.com` | `axiom-dev` | Dev environment DNS records |
| `staging.axiom.com` | `axiom-staging` | Staging environment DNS records |
| `axiom.com` (prod records) | `axiom-prod` | Production DNS records via delegation from root zone |

**DNS records (prod — other environments follow the same pattern with `{env}.axiom.com` subdomains):**

| Record | Type | Target |
|---|---|---|
| `api.axiom.com` | A (alias) | ALB |
| `app.axiom.com` | A (alias) | CloudFront distribution |
| `axiom.com` | A (alias) | CloudFront distribution |

### ACM Certificates

| Certificate | Account | Domain | Used by |
|---|---|---|---|
| `*.axiom.com` + `axiom.com` | `axiom-prod` | Wildcard + apex | ALB, CloudFront |
| `*.dev.axiom.com` | `axiom-dev` | Dev wildcard | ALB |
| `*.staging.axiom.com` | `axiom-staging` | Staging wildcard | ALB |

All certificates use DNS validation via Route 53 (automatic). ACM auto-renews 60 days before expiry as long as the validation CNAME remains in place. CloudFront requires certificates in `us-east-1` — since all infrastructure is in `us-east-1`, no cross-region certificate is needed.

### CloudFront

One distribution per environment serving the React SPA from S3.

| Setting | Value |
|---|---|
| Origin | S3 bucket `axiom-{env}-spa` (private, Origin Access Control) |
| Alternate domain | `app.axiom.com` (prod), `app.staging.axiom.com`, `app.dev.axiom.com` |
| Viewer protocol | Redirect HTTP to HTTPS |
| TLS minimum | TLSv1.2_2021 |
| Default root object | `index.html` |
| Error pages | 403 and 404 both return `/index.html` with 200 (SPA client-side routing) |
| Cache policy | `CachingOptimized` for `/assets/*`. `CachingDisabled` for `index.html` (ensures deploys are picked up immediately). |
| Response headers | `Strict-Transport-Security`, `X-Content-Type-Options: nosniff`, `X-Frame-Options: DENY`, `Content-Security-Policy` |
| WAF | WebACL attached (prod and staging; not at demo stage) |

### AWS WAF

**CloudFront WebACL:**

| Rule | Purpose |
|---|---|
| AWS Managed Core Rule Set | OWASP top-10 protections (XSS, SQLi, SSRF) |
| AWS Managed Known Bad Inputs | Log4Shell, Java deserialization, known attack patterns |
| Rate limit: 2000 requests/5 min per IP | DDoS mitigation for SPA |

**ALB WebACL:**

| Rule | Purpose |
|---|---|
| AWS Managed Core Rule Set | OWASP protections at API layer |
| AWS Managed SQL Injection Rule Set | Additional SQLi protection for API endpoints |
| Rate limit: 1000 requests/5 min per IP | API-level DDoS mitigation |
| Geo-restriction | US and Canada only at launch |

WAF logs to CloudWatch for analysis. Blocked requests logged with triggering rule for incident investigation.

Not provisioned at demo stage — acceptable risk with zero customers. Enabled when the first paying customer is onboarded.

### SES (Transactional Email)

| Setting | Value |
|---|---|
| Sending domain | `axiom.com` (DKIM + SPF + DMARC via Route 53) |
| Configuration set | `axiom-{env}-transactional` with CloudWatch event publishing (sends, bounces, complaints) |
| Sandbox | Dev and staging: sandbox mode (verified addresses only). Prod: production access. |

**DMARC policy:** `v=DMARC1; p=reject; rua=mailto:dmarc@axiom.com` — strict enforcement from day one.

**Bounce/complaint handling:** SES publishes bounce and complaint events to an SNS topic. A Lambda function processes these to suppress sends to hard-bounced addresses and flag complaint addresses in the Identity Service.

---

## 7. Observability

### CloudWatch Log Groups

| Log group | Source | Retention (dev / staging / prod) |
|---|---|---|
| `/axiom/{env}/axiom-api` | ECS awslogs | 7d / 30d / 365d |
| `/axiom/{env}/doc-processing` | ECS awslogs | 7d / 30d / 365d |
| `/axiom/{env}/provenance-signer` | ECS awslogs | 7d / 30d / 365d |
| `/axiom/{env}/pgbouncer` | ECS sidecar awslogs | 7d / 30d / 365d |
| `/axiom/{env}/waf/cloudfront` | WAF logging | 7d / 30d / 365d |
| `/axiom/{env}/waf/alb` | WAF logging | 7d / 30d / 365d |
| `/axiom/{env}/rds` | RDS PostgreSQL (error, slow query) | 7d / 30d / 365d |
| `/axiom/{env}/stepfunctions` | Step Functions execution history | 7d / 30d / 365d |
| `/axiom/{env}/vpc-flow-logs` | VPC flow logs (rejected only) | 7d / 30d / 365d |

All log groups encrypted with the environment's `axiom-{env}-default` KMS key. All Go services emit structured JSON logs via `slog`.

**CloudWatch Logs Insights saved queries (provisioned via Terraform):**
- Slow database queries (parsed from PgBouncer and RDS slow query logs)
- 5xx error rate by service
- Authentication failures by IP
- AI inference latency by model and operation
- Step Functions execution failures

### CloudWatch Metrics

**Custom metrics published by application code (via OpenTelemetry SDK → CloudWatch):**

| Metric | Namespace | Dimensions |
|---|---|---|
| `http_request_duration_seconds` | `Axiom/{env}` | service, method, path, status_code |
| `http_requests_total` | `Axiom/{env}` | service, method, path, status_code |
| `active_websocket_connections` | `Axiom/{env}` | service |
| `river_job_duration_seconds` | `Axiom/{env}` | job_type, status |
| `river_job_queue_depth` | `Axiom/{env}` | job_type |
| `bedrock_invocation_duration_seconds` | `Axiom/{env}` | model, operation |
| `bedrock_invocation_tokens` | `Axiom/{env}` | model, operation, token_type |
| `ai_mapping_suggestion_latency_seconds` | `Axiom/{env}` | framework, control_family |
| `ai_gap_analysis_duration_seconds` | `Axiom/{env}` | framework |
| `ai_drift_detection_rate` | `Axiom/{env}` | framework, control_family |
| `ai_findings_triage_throughput` | `Axiom/{env}` | framework, severity |
| `ai_decision_created_total` | `Axiom/{env}` | context_type, framework, hitl_state |
| `provenance_sign_duration_seconds` | `Axiom/{env}` | artifact_type |
| `provenance_sign_failures_total` | `Axiom/{env}` | artifact_type, error_class |
| `scf_import_records_processed` | `Axiom/{env}` | source (`scf` / `oscal` / `aicpa` / `cis`), operation |

### CloudWatch Dashboards

Provisioned via Terraform as JSON definitions:

| Dashboard | Contents |
|---|---|
| `Axiom-{env}-Overview` | Request rate, error rate, p50/p95/p99 latency by route group (module). ECS task count. RDS connections and CPU. ALB healthy host count. Active WebSocket connections. |
| `Axiom-{env}-AI` | Bedrock invocation latency, token consumption, throttle rate by model and feature. River AI job queue depth and processing time by job type. Mapping suggestion latency (by framework). Gap analysis duration (by framework). Drift detection rate (by framework/control family). Findings triage throughput (by severity). |
| `Axiom-{env}-AIDecision` | **ISO 42001 observability surface.** AIDecision creation rate by `context_type`, `framework`, and HITL state (`suggested`, `accepted`, `overridden`, `rejected`). Override/acceptance ratio over time. Backlog of `suggested` decisions awaiting human review. Per-user reviewer throughput. Per-model decision counts (audit trail of which Claude version made each suggestion). |
| `Axiom-{env}-Provenance` | Provenance signing rate (artifacts/min), signing latency p50/p95/p99, signing failure count by `error_class`. Object Lock writes to evidence and reports buckets. KMS `Sign` call rate and error rate on the provenance-signing key. Verification failures (signed artifacts that fail re-verification). |
| `Axiom-{env}-Crosswalk` | SCF import job status and record count by source (SCF/OSCAL/AICPA/CIS). Count of `CommonControl` and STRM edges by framework and version. Partial-satisfaction coverage distribution. Framework-version migration job throughput. |
| `Axiom-{env}-Data` | RDS CPU, free storage, read/write IOPS, active connections. PgBouncer pool utilization. S3 bucket size and request count (including `scf-catalog` bucket). |
| `Axiom-{env}-Security` | WAF blocked requests by rule. Authentication failures. SES bounce/complaint rate. |

Dashboards are not provisioned at demo stage. Enabled when the full observability workspace is activated.

### CloudWatch Alarms

| Alarm | Condition | Action |
|---|---|---|
| `{env}-alb-5xx-high` | ALB 5xx > 10 in 5 min | SNS → ops email |
| `{env}-alb-unhealthy-hosts` | Unhealthy hosts > 0 for 2 consecutive periods | SNS → ops email |
| `{env}-rds-cpu-high` | RDS CPU > 80% for 10 min | SNS → ops email |
| `{env}-rds-storage-low` | Free storage < 20% of allocated | SNS → ops email |
| `{env}-rds-connections-high` | Connections > 80% of max | SNS → ops email |
| `{env}-ecs-cpu-high-{service}` | Service avg CPU > 80% for 10 min (axiom-api, doc-processing, provenance-signer) | SNS → ops email |
| `{env}-river-dlq-not-empty` | River DLQ depth > 0 (any worker: `evidence_control_mapping`, `gap_analysis`, `drift_detection`, `framework_migration`, `scf_import`, `management_response_drafter`, `findings_triage`, `provenance.sign`, `auditcore.document-extract`, `reporting.render`) | SNS → ops email |
| `{env}-bedrock-throttle-high` | Bedrock throttle > 5% of invocations in 5 min | SNS → ops email |
| `{env}-provenance-sign-failures` | `provenance_sign_failures_total` > 0 in 5 min (any failure is investigated — signing is on the audit-integrity path) | SNS → ops email + PagerDuty (prod) |
| `{env}-scf-import-failure` | Scheduled `scf_import` job has no successful completion in 48 hours | SNS → ops email |
| `{env}-aidecision-backlog-high` | Backlog of `suggested` AIDecision rows older than SLA threshold (ISO 42001 HITL compliance) | SNS → ops email |
| `{env}-ses-bounce-high` | Bounce rate > 5% in 1 hour | SNS → ops email |
| `{env}-waf-blocked-spike` | WAF blocked > 100 in 5 min | SNS → ops email |
| `{env}-stepfunctions-failed` | `ExecutionsFailed` > 0 | SNS → ops email |
| `{env}-certificate-expiry` | ACM days to expiry < 30 | SNS → ops email |

SNS topics: `axiom-prod-ops-alerts` (prod alarms), `axiom-nonprod-ops-alerts` (dev/staging). Subscriptions: email at launch, Slack webhook or PagerDuty when the team grows.

Demo stage provisions only basic alarms: ALB 5xx, RDS CPU, and RDS storage.

### X-Ray Distributed Tracing

OpenTelemetry Go SDK in the `platform` package with AWS X-Ray exporter. Within the monolith, traces propagate via Go context. The `X-Amzn-Trace-Id` header is used for the Axiom API → Document Processing HTTP call.

**Sampling rules:**
- Prod: 5% of requests
- Staging: 100%
- Dev: 100%

**X-Ray groups:**
- `high-latency` — response time > 2 seconds
- `errors` — fault or error status
- `ai-operations` — requests hitting Bedrock endpoints (annotated with `ai_feature` and `framework`)
- `provenance-signing` — spans that call the provenance-signer service or KMS `Sign` (annotated with `artifact_type`)
- `crosswalk` — requests touching `CommonControl` / STRM edge resolution (annotated with `framework` and `framework_version`)

Not enabled at demo stage.

### CloudTrail

Enabled at the Organization level from the management account:

| Setting | Value |
|---|---|
| Management events | All accounts, all regions |
| S3 data events | `axiom-prod-evidence`, `axiom-prod-archive` (object-level read/write tracking) |
| KMS data events | Enabled (tracks decrypt calls against HIPAA keys) |
| Log file integrity validation | Enabled |
| Log encryption | `axiom-cloudtrail` KMS key |
| Retention | 365 days in CloudTrail. S3 lifecycle: Glacier after 90 days, retained 7 years (SOC 2 Type 2 historical trail requirements, ISO 27001 evidence retention, ISO 42001 decision-audit requirements, HIPAA §164.316(b)(2) 6-year minimum). |

---

## 8. CI/CD Pipeline

### ECR Repositories

One repository per service in the tooling account:

| Repository | Image source |
|---|---|
| `axiom/axiom-api` | Go binary (modular monolith) |
| `axiom/doc-processing` | Python + Tesseract |
| `axiom/provenance-signer` | Go binary (isolated signing service) |
| `axiom/pgbouncer` | PgBouncer with config |

**Configuration:**
- Image scanning on push (CVE detection)
- Lifecycle policy: retain last 20 tagged images, expire untagged images after 7 days
- Immutable image tags — a pushed tag cannot be overwritten
- Cross-account pull access: IAM policies allow dev, staging, and prod execution roles to pull

Demo stage provisions all 4 repositories: `axiom-api`, `doc-processing`, `provenance-signer`, `pgbouncer`.

### GitHub Actions OIDC Federation

No long-lived AWS credentials. Each environment has a dedicated IAM role assumed via OIDC:

| Role | Account | Trust condition |
|---|---|---|
| `axiom-ci-tooling` | `axiom-tooling` | `repo:your-org/axiom:*` (ECR push, Terraform state) |
| `axiom-ci-dev` | `axiom-dev` | `repo:your-org/axiom:ref:refs/heads/*` (any branch) |
| `axiom-ci-staging` | `axiom-staging` | `repo:your-org/axiom:ref:refs/heads/main` (main only) |
| `axiom-ci-prod` | `axiom-prod` | `repo:your-org/axiom:ref:refs/heads/main` + GitHub Environment approval |

Trust policies use the OIDC `sub` claim to restrict which branches can assume which roles. A feature branch can deploy to dev but cannot assume the staging or prod role.

### Pipeline Workflows

**On pull request (any branch):**

```
lint → build → unit test → terraform plan (dev)
```

- Lint: `golangci-lint` (Go), `ruff` (Python), `eslint` (frontend)
- Build: Turborepo builds affected services only (change detection by directory)
- Unit test: `go test ./...` per affected service, `pytest` for doc-processing
- Terraform plan: runs plan for all workspaces against dev, posts diff as PR comment. No apply.

**On merge to main:**

```
build → push images → deploy staging → integration test → manual gate → deploy prod
```

1. **Build + push:** Turborepo builds affected apps. Docker images pushed to ECR with `main-{short-sha}` and `latest` tags.
2. **Deploy staging:** Terraform apply for affected workspaces in `axiom-staging`. ECS services updated to new task definitions. Wait for deployment stability (circuit breaker confirms healthy).
3. **Integration tests:** Run against staging — API smoke tests, critical path end-to-end tests.
4. **Manual approval gate:** GitHub Environment protection rule on `prod`. Requires one approval from a designated deployer.
5. **Deploy prod:** Same Terraform apply + ECS update flow as staging.

**Terraform apply order (follows dependency graph):**

```
bootstrap (manual only, never auto-applied)
    ↓
network
    ↓
data
    ↓
compute ← dns-cdn
    ↓
observability
```

`cicd` workspace applies independently — no downstream dependencies. Only workspaces with changes in the current commit are applied (change detection: diff `infra/{workspace}/` against main).

**Database migrations (separate from Terraform):**

Migrations run as a dedicated one-shot ECS Fargate task within the deploy pipeline:

1. After `data` workspace is applied (RDS exists)
2. Before `compute` workspace is applied (services expect the new schema)
3. Fargate task runs `golang-migrate` against `axiom_db` with pending migrations
4. Migration task uses `axiom/{env}/rds/master` credentials from Secrets Manager
5. Pipeline proceeds to `compute` only if all migrations succeed

### Rollback Strategy

- **ECS:** Circuit breaker handles automatic rollback during deployment. Manual rollback: update ECS service to previous task definition revision.
- **Terraform:** State is versioned in S3. Apply a prior commit's configuration to roll back infrastructure.
- **Database migrations:** `golang-migrate` supports down migrations. Destructive down migrations (drop column, drop table) require explicit review. Non-destructive migrations (add column, add table) don't need rollback — previous application version ignores new columns.

---

## 9. Terraform Workspace Design

### Directory Structure

```
infra/
  modules/                    — Reusable Terraform modules
    vpc/
    ecs-cluster/
    ecs-service/              — One module, parameterized per service
    rds/
    s3-bucket/
    cloudfront/
    waf/
    alarms/
  workspaces/
    bootstrap/                — Organizations, accounts, OIDC, state backend
    network/                  — VPC, subnets, NAT, VPC endpoints, security groups
    data/                     — RDS, S3, SQS, Secrets Manager, KMS keys
    compute/                  — ECS cluster, ALB, services, task defs, Step Functions, Bedrock
    dns-cdn/                  — Route 53, CloudFront, ACM, WAF, SES
    observability/            — CloudWatch log groups, dashboards, alarms, X-Ray, CloudTrail
    cicd/                     — ECR repositories, GitHub Actions OIDC roles
  envs/
    demo.tfvars               — Demo stage (dev + staging, lean config)
    dev.tfvars                — Full dev environment
    staging.tfvars            — Full staging environment
    prod.tfvars               — Production
```

### State Backend

All state in a single S3 bucket in the tooling account, keyed by workspace and environment:

```
s3://axiom-terraform-state/
  bootstrap/terraform.tfstate
  network/dev/terraform.tfstate
  network/staging/terraform.tfstate
  network/prod/terraform.tfstate
  data/dev/terraform.tfstate
  ...
```

DynamoDB lock table: `axiom-terraform-locks`, partition key `LockID`. One table handles all workspaces — lock keys include the full state path.

The `bootstrap` workspace is not environment-scoped — it manages Organization-level resources that exist once.

S3 bucket configuration: versioning enabled, `axiom-terraform-state` KMS key encryption, bucket policy denying non-TLS and public access.

### Cross-Workspace Data Flow

Workspaces share outputs via **SSM Parameter Store** in each workload account. This avoids `terraform_remote_state` data sources, which tightly couple workspaces and require shared state bucket permissions.

**Pattern:** each workspace writes outputs to SSM parameters on apply. Downstream workspaces read parameters as data sources.

| Producer | SSM parameter examples | Consumer |
|---|---|---|
| `network` | `/axiom/{env}/vpc-id` | `data`, `compute` |
| `network` | `/axiom/{env}/private-subnet-ids` | `data`, `compute` |
| `network` | `/axiom/{env}/public-subnet-ids` | `compute`, `dns-cdn` |
| `network` | `/axiom/{env}/sg-alb` | `compute` |
| `network` | `/axiom/{env}/sg-ecs-gateway` | `compute` |
| `network` | `/axiom/{env}/sg-ecs-services` | `compute`, `data` |
| `network` | `/axiom/{env}/sg-rds` | `data` |
| `data` | `/axiom/{env}/rds-endpoint` | `compute` |
| `data` | `/axiom/{env}/rds-port` | `compute` |
| `data` | `/axiom/{env}/s3-evidence-bucket-arn` | `compute` |
| `data` | `/axiom/{env}/s3-archive-bucket-arn` | `compute` |
| `data` | `/axiom/{env}/s3-reports-bucket-arn` | `compute` |
| `data` | `/axiom/{env}/s3-scf-catalog-bucket-arn` | `compute` |
| `data` | `/axiom/{env}/kms-hipaa-key-arn` | `compute` |
| `data` | `/axiom/{env}/kms-provenance-signing-key-arn` | `compute` |
| `compute` | `/axiom/{env}/alb-dns-name` | `dns-cdn` |
| `compute` | `/axiom/{env}/alb-hosted-zone-id` | `dns-cdn` |
| `cicd` | `/axiom/ecr-repo-axiom-api-url` | `compute` (all envs) |
| `cicd` | `/axiom/ecr-repo-doc-processing-url` | `compute` (all envs) |
| `cicd` | `/axiom/ecr-repo-provenance-signer-url` | `compute` (all envs) |
| `cicd` | `/axiom/ecr-repo-pgbouncer-url` | `compute` (all envs) |

SSM parameters are typed as `String`, encrypted with the account's default KMS key.

### Environment Differentiation

Each workspace uses shared `.tf` files with environment-specific values via tfvars:

```hcl
# prod.tfvars
environment           = "prod"
rds_instance_class    = "db.r7g.xlarge"
rds_multi_az          = true
rds_backup_retention  = 35
rds_storage_gb        = 200
rds_storage_autoscale = 1000
ecs_axiom_api_min     = 2
ecs_axiom_api_max     = 8
ecs_doc_processing_min = 1
ecs_doc_processing_max = 4
ecs_provenance_signer_min = 2
ecs_provenance_signer_max = 4
log_retention_days    = 365
single_nat            = false
enable_s3_replication = true
enable_s3_object_lock = true
enable_exec_command   = false
enable_waf            = true
enable_guardduty      = true
enable_aws_config     = true
enable_xray           = true
enable_step_functions = true
xray_sampling_rate    = 0.05
```

```hcl
# demo.tfvars
environment              = "demo"
rds_instance_class       = "db.t4g.medium"
rds_multi_az             = false
rds_backup_retention     = 1
rds_storage_gb           = 20
rds_storage_autoscale    = 0
ecs_axiom_api_min        = 1
ecs_axiom_api_max        = 2
ecs_doc_processing_min   = 1
ecs_doc_processing_max   = 2
ecs_provenance_signer_min = 1
ecs_provenance_signer_max = 2
enable_step_functions    = false
enable_waf               = false
enable_guardduty         = false
enable_aws_config        = false
enable_xray              = false
enable_s3_replication    = false
enable_s3_object_lock    = false
log_retention_days       = 7
single_nat               = true
xray_sampling_rate       = 0.0
```

CI/CD passes the correct tfvars: `terraform apply -var-file=../../envs/{env}.tfvars`

### Migration Path: Monolith → Extracted Services

With the modular monolith architecture, compute Terraform is already simple (2 ECS services). The per-service workspace split described below applies only if a module is extracted from the monolith into its own service.

**When to extract a module:**
- A second team needs independent deployment of a specific module
- WebSocket scaling (workpaper) forces over-provisioning of the combined service
- A module's resource requirements diverge significantly from the rest

**Extraction steps:**

1. Create a new `apps/{module}/` directory with its own `go.mod`, pulling the module's code out of `internal/`.
2. Replace in-process calls with HTTP client calls (the Go interface boundary makes this mechanical).
3. Add a new ECS service in the `compute` workspace using the existing `ecs-service` module.
4. Optionally move the new service's resources to a per-service workspace (`infra/workspaces/svc-{module}/`) if independent Terraform applies are needed.

The `ecs-service` Terraform module is parameterized to support any number of services without restructuring.

---

## 10. Security Controls (Infrastructure Layer)

### S3 Account-Level Public Access Block

Enabled on all five accounts. No S3 bucket in any account can be made public regardless of bucket-level settings. Single highest-impact control for preventing accidental data exposure.

### VPC Flow Logs

Enabled on all VPCs. Capture rejected traffic only (accepted traffic is high volume, low signal at launch). Destination: CloudWatch log group `/axiom/{env}/vpc-flow-logs`.

When security monitoring matures, switch to all traffic and ship to S3 for batch analysis.

### GuardDuty

Enabled at the Organization level from the management account.

**Detectors:**
- EC2/ECS runtime monitoring (compromised containers)
- S3 protection (anomalous access to evidence buckets)
- RDS login monitoring (brute-force or anomalous database authentication)

High-severity findings (`CRITICAL`, `HIGH`) trigger SNS → ops email. Expected cost: $20-50/month across all accounts.

Not enabled at demo stage.

### AWS Config

Enabled in all workload accounts with compliance rules:

| Rule | Purpose |
|---|---|
| `rds-instance-public-access-check` | RDS must not be publicly accessible |
| `rds-storage-encrypted` | RDS must be encrypted at rest |
| `s3-bucket-ssl-requests-only` | Buckets must deny non-TLS requests |
| `s3-bucket-public-read-prohibited` | No public read |
| `s3-bucket-public-write-prohibited` | No public write |
| `ecs-task-definition-log-configuration` | All ECS tasks must ship logs |
| `secretsmanager-rotation-enabled-check` | All secrets must have rotation configured |
| `iam-no-inline-policy-check` | No inline IAM policies |
| `vpc-flow-logs-enabled` | Flow logs must be on |
| `cloudtrail-enabled` | CloudTrail must be on |

Non-compliant resources trigger CloudWatch Event → SNS notification. No auto-remediation at launch — alert and investigate. SOC 2 and ISO 27001 conformance packs available from AWS for layering on before certification audits.

Not enabled at demo stage.

### IAM Boundaries and Conventions

**No IAM users with long-lived credentials in any workload account.** SCPs enforce this in production. All human access via IAM Identity Center. All CI/CD access via OIDC federation.

**Role naming convention:**
- `axiom-{env}-axiom-api-task-role` — Axiom API ECS task role
- `axiom-{env}-doc-processing-task-role` — Document Processing ECS task role
- `axiom-{env}-provenance-signer-task-role` — Provenance signer ECS task role (isolated KMS `Sign` permission on the provenance-signing key)
- `axiom-{env}-scf-import-role` — Scheduled SCF/OSCAL/AICPA/CIS catalog import (S3 read on `scf-catalog`, RDS write via Secrets Manager credentials)
- `axiom-{env}-ecs-execution-role` — shared execution role
- `axiom-ci-{env}` — GitHub Actions deployment role
- `axiom-{env}-migration-role` — database migration task role

**Permission boundaries** applied to all roles created by Terraform. The boundary denies:
- `iam:CreateUser` (no backdoor user creation)
- `organizations:LeaveOrganization`
- `cloudtrail:StopLogging`, `cloudtrail:DeleteTrail`
- `s3:PutBucketPublicAccessBlock` with `BlockPublicAccess = false`

Defense-in-depth: even if a role's policy is overly broad, the boundary prevents highest-impact destructive actions.

---

## 11. Cost Estimates

### Demo Stage (Pre-Revenue)

Two workload accounts (dev + staging), 2 services each, single AZ, no security controls beyond CloudTrail and S3 public access blocks.

| Category | Service | Dev | Staging | Tooling/Mgmt | Total |
|---|---|---|---|---|---|
| Compute | ECS Fargate (3 svc × 1-2 tasks — axiom-api, doc-processing, provenance-signer) | $45-75 | $45-75 | — | $90-150 |
| Database | RDS db.t4g.medium, Single-AZ | $55-65 | $55-65 | — | $110-130 |
| Networking | NAT Gateway (×1) | $35-45 | $35-45 | — | $70-90 |
| Networking | VPC Endpoints (×4, 1 AZ) | $30 | $30 | — | $60 |
| Networking | ALB | $20-25 | $20-25 | — | $40-50 |
| Storage | S3 (evidence only) | $1 | $1 | — | $2 |
| CDN | CloudFront | — | $1-2 | — | $1-2 |
| DNS | Route 53 | $1 | $1 | $2 | $4 |
| AI | Bedrock | $5-10 | $5-15 | — | $10-25 |
| Email | SES (sandbox) | $0 | $0 | — | $0 |
| Secrets | Secrets Manager | $1 | $1 | — | $2 |
| KMS | 1 key/account | $1 | $1 | $1 | $3 |
| Observability | CloudWatch (logs + basic alarms) | $5-10 | $5-10 | — | $10-20 |
| Registry | ECR | — | — | $1-2 | $1-2 |
| Org/Trail | CloudTrail | — | — | $5 | $5 |
| | **Subtotals** | **$190-280** | **$190-290** | **$10-15** | |
| | **Demo stage total** | | | | **$380-585/mo** |

### Full Architecture (Post First Customer)

| Category | Service | Est. monthly (prod) |
|---|---|---|
| Compute | ECS Fargate (3 svc — axiom-api, doc-processing, provenance-signer — 2-4 tasks avg per service) | $200-320 |
| Database | RDS db.r7g.xlarge, Multi-AZ, 200 GB | $550-600 |
| Networking | NAT Gateways (×2) | $70-90 |
| Networking | VPC Endpoints (×8 interface, 2 AZs) | $120 |
| Networking | ALB | $25-40 |
| Storage | S3 (evidence + archive + reports + SPA) | $5-10 |
| Storage | S3 cross-region replication | $5-10 |
| CDN | CloudFront | $5-10 |
| DNS | Route 53 | $5 |
| Security | WAF (2 WebACLs) | $20-30 |
| Security | GuardDuty | $20-50 |
| Security | KMS (5 keys) | $5-10 |
| Observability | CloudWatch | $30-50 |
| Observability | X-Ray (5% sampling) | $5-10 |
| Observability | CloudTrail | $10-20 |
| Config | AWS Config | $10-20 |
| Workflow | Step Functions | $1 |
| AI | Bedrock (on-demand) — compliance workload: ~$3-8 per engagement across frameworks (evidence-to-control mapping, gap analysis, drift detection, findings triage, management-response drafting). Population not financial-audit-sized. | $100-250 |
| Email | SES | $5-10 |
| Secrets | Secrets Manager | $3-5 |
| Registry | ECR | $2-5 |
| | **Production total** | **$1,150-1,550/mo** |

**Staging:** same topology, smaller instances — **$550-750/mo**.
**Dev:** single AZ, minimum tasks — **$280-400/mo**.
**Total across all environments:** **$1,980-2,700/mo ($24k-32k/year)**.

### Growth Path

| Milestone | Action | Monthly cost |
|---|---|---|
| Demo stage (pre-revenue) | 2 workload accounts, 2 services, single AZ | $380-585 |
| First paying customer | Add prod (Multi-AZ, security controls). Upgrade staging. | $1,800-2,500 |
| Stable operations | Full dashboards, X-Ray, alarms. | $2,100-2,900 |
| Pre-SOC 2 audit | GuardDuty, AWS Config, conformance packs in all accounts. | $2,300-3,100 |
| Savings plans | Fargate + RDS reserved instances (3+ months stable). 25-35% savings. | $1,700-2,300 |

### Cost Scaling Notes

- **Biggest cost driver:** RDS instance size. `db.r7g.2xlarge` doubles the RDS line. Aurora becomes competitive when read replicas are needed.
- **ECS scales linearly:** ~$30-60/month per additional Fargate task. With 2 services instead of 7, the Fargate baseline is significantly lower.
- **Bedrock scales with engagements:** compliance engagements cost ~$3-8 each in AI inference (per `ai-architecture-design.md`). 500 engagements/month → $1,500-4,000/month, still trivial relative to subscription revenue.
- **NAT Gateway data charges:** evidence file volumes could surprise. S3 gateway endpoint (free) already mitigates the largest NAT traffic source.
- **Savings plans:** after 3+ months stable, Fargate Savings Plans (1-year no-upfront, 20-30% savings) and RDS Reserved Instances (1-year partial-upfront, 30-40% savings).

---

## 12. What This Design Defers

| Item | Current posture | Trigger to revisit |
|---|---|---|
| Aurora PostgreSQL | Standard RDS. Simpler, cheaper at launch scale. | Storage > 500 GB, need read replicas with low replication lag, IOPS bottleneck on gp3 |
| Service extraction | Modular monolith — all Go modules in one binary. | Second team needs independent deploy cadence, WebSocket scaling diverges from REST, module resource needs diverge significantly |
| mTLS for doc-processing | Plaintext HTTP within VPC. | SOC 2 auditor requires service-to-service authentication. VPC Lattice provides mTLS natively. |
| Per-service Terraform workspaces | Single `compute` workspace for 2 services. | A module is extracted to its own service and needs independent Terraform lifecycle |
| Grafana Cloud / Datadog | CloudWatch-native only. | CloudWatch dashboarding UX becomes a bottleneck, team needs advanced alerting or APM |
| Multi-region | `us-east-1` only. | EU/APAC enterprise customers require data residency. Architecture supports adding `eu-central-1` deployment. |
| WAF Bot Control | Managed core rules + rate limiting only. | Bot traffic becomes a meaningful fraction of requests |
| Auto-remediation (AWS Config) | Alert only, no auto-remediation. | Confidence in remediation actions, team has runbooks for common violations |
| VPC flow logs (all traffic) | Rejected traffic only. | Security monitoring maturity, SIEM integration, forensic requirements |
| PagerDuty / Opsgenie | SNS → email for all alarms. | Team grows, on-call rotation is formalized |
