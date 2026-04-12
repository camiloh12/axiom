# Research Task 7: Integration Landscape and Roadmap

## Scope and Constraints

This analysis covers integrations for Axiom targeting mid-market CPA and advisory firms (20–200 staff) doing financial audits and compliance framework audits (SOC 2, ISO 27001, HIPAA) in the US and Canada. The existing spec identifies Jira, Gusto, AWS, Dropbox, Box, Google Drive, and O365 as integrations. This research evaluates those, surfaces gaps (particularly ERP/accounting systems), and produces a prioritized roadmap.

---

## 1. Financial Audit Integrations: ERP and Accounting Systems

### Accounting Software Market Structure (US Mid-Market)

The accounting software landscape for private mid-market companies (the clients of Axiom's target ICP) follows a clear segmentation by company size and complexity:

**Under $5M revenue / under 25 employees:** QuickBooks Online dominates, with over 62% of the overall US accounting software market and an even higher share at the small end of mid-market. These clients run single-entity, single-currency books. QuickBooks Desktop (Enterprise edition) extends the lifecycle for slightly larger single-entity manufacturers and distributors.

**$5M–$50M revenue:** The transition zone. Over half of SaaS companies below $5M ARR use QuickBooks; by $50M ARR, half have transitioned to NetSuite. The actual trigger for transition is not revenue size alone — it is operational complexity: multi-entity consolidation, international operations, complex revenue recognition, or audit trail requirements. A $30M single-entity services business often still runs on QuickBooks Online. A $15M company with three subsidiaries likely needs NetSuite or Sage Intacct. This is the range where the ICP's clients live.

**$50M–$200M revenue (upper mid-market):** NetSuite dominates for technology and services companies. By $100M ARR, two-thirds of companies are on NetSuite. Sage Intacct leads among SaaS-model businesses in the $11M–$50M range (Sage's own survey data). Microsoft Dynamics 365 Business Central has significant presence in manufacturing, distribution, and professional services.

**Important: Xero is primarily a Canadian and international small-business platform** in the North American market. It has meaningful Canadian market share in smaller firms. For Axiom's launch focus (US and Canada), Xero clients are mostly at the smaller end of the client spectrum and tend to be very simple books.

**Practical implication for Axiom's ICP:** The financial audit clients being audited by a 20–60 person CPA firm will primarily run:
- QuickBooks Online (most common for clients under ~$20M revenue)
- QuickBooks Desktop / Enterprise (common in construction, manufacturing, regional businesses)
- Sage Intacct (SaaS companies, nonprofits, services companies in the $5M–$100M range)
- NetSuite (growing tech companies, multi-entity, higher-complexity)
- Xero (smaller clients, especially Canadian)
- Microsoft Dynamics 365 Business Central (less common but present in manufacturing)

SAP and Oracle are enterprise systems. Mid-market CPA firms auditing 20–60 person firms will rarely if ever audit a client on full SAP or Oracle EBS — those clients are large enough to have Big Four auditors.

### What Financial Audit Data Is Actually Needed

The primary data a financial auditor needs from a client's accounting system is:

**Trial Balance** — A list of all accounts with their debit/credit balances at period end. This is the foundation of every financial audit. It is used to populate the lead schedules, tie out to the financial statements, and calculate materiality. Most audit work starts here.

**General Ledger (GL) Transactions** — The full detail of every posting to each account. Used for:
- Journal entry testing (PCAOB AS 2410 / AU-C 240) — auditors examine unusual, manual, or year-end journal entries for fraud indicators
- Substantive testing of specific accounts (e.g., all transactions in accounts payable over a threshold)
- Population definition for sampling (auditors pull the GL to define their population before sampling)
- Variance analysis (comparing period-over-period GL activity)

**Subledger Reports** — Accounts receivable aging, accounts payable aging, fixed asset schedules, inventory listings, payroll registers. These are used to test the completeness and accuracy of balance sheet accounts. These are typically available as reports (PDF or Excel export) rather than structured data feeds.

**Journal Entry Listings** — All manually posted journal entries for the period, with user, date, amounts, and accounts. This is a specific PCAOB/AICPA requirement: auditors must test all manual JEs and all unusual JEs. The listing needs to be complete and show who posted each entry.

**The fundamental API challenge:** Most accounting platforms do not expose a clean, structured "trial balance" endpoint. A trial balance is computed by summing all GL activity — it is a report, not a record. Auditors most commonly receive trial balances as:
1. Excel/CSV export from the client's accounting system (most common in practice)
2. PDF report export (still common, especially for QuickBooks Desktop)
3. In larger engagements, a structured data extract from the ERP

For the GL transaction detail and journal entry listing, the data is more structured and more available via API, but the schema varies significantly across platforms.

### API Availability by Platform

**QuickBooks Online (Intuit):** Has a well-documented REST API. Exposes accounts, journal entries, transactions (purchase, payment, invoice), and can reconstruct a trial balance by querying account balances. Does not have a single "trial balance" endpoint — requires querying `Account` objects and summing. GL transaction detail requires querying individual transaction types. The Intuit API is rate-limited and the data model is specific to QBO — not a clean financial data model.

**NetSuite (Oracle):** Exposes a REST API and SuiteScript API. Trial balance is accessible as a saved search or via SuiteAnalytics. GL transactions are queryable. The API is powerful but requires NetSuite-specific knowledge. NetSuite accounts are often heavily customized, making normalized extraction complex.

**Sage Intacct:** XML-based API (legacy but comprehensive). Exposes journal entries, GL detail, and financial reports. The `GLJOURNAL` and `GLENTRY` objects give full GL transaction detail. Well-documented for developers. Less clean than REST but more complete than some REST APIs for financial data specifically.

**Xero:** REST API with good documentation. Exposes accounts, journal lines, and reports. Has a specific `Reports/TrialBalance` endpoint — one of the few platforms with a direct trial balance endpoint. Good for the simpler-book clients.

**Microsoft Dynamics 365 Business Central:** REST API available. Exposes ledger entries, accounts, and financial statements. Less commonly used for direct API integration in audit contexts — most firms get Excel exports.

### The Aggregator Case: Codat

Codat is the primary candidate for normalizing access to accounting systems. Key findings:

**Coverage:** Integrates with 35+ accounting platforms including QuickBooks Online, QuickBooks Desktop, Sage Intacct, NetSuite, Xero, Wave, MYOB, FreshBooks, Microsoft Dynamics, and others. Claims 90% coverage of small businesses with a single API integration.

**Data types available through Codat's Accounting API:**
- Accounts (chart of accounts)
- Account transactions (GL-level transactions)
- Journal entries (with transactional amount and currency support added in 2024)
- Balance sheet (computed)
- Profit & loss (computed)
- Credit notes, invoices, purchase orders, payments, suppliers, customers
- Bank accounts and bank transactions (separate Banking API)

**Critical limitation for audit use:** Codat is primarily designed for lending and fintech (underwriting, expense management, corporate cards). Its "Assess" product provides derived financial metrics for lenders — liquidity ratios, creditworthiness indicators. Codat does not have an audit-specific data product. This means:
- The trial balance would need to be reconstructed from account balances (Codat's balance sheet endpoint), not pulled as a discrete report
- Journal entry completeness depends on the underlying platform's API — Codat normalizes what the API exposes, but if the platform doesn't fully expose all JEs via API, Codat cannot fill the gap
- Data freshness varies: Codat caches data with platform-dependent sync frequencies, which is acceptable for underwriting but may be insufficient for audit evidence (auditors need a point-in-time snapshot that is demonstrably complete)

**Codat pricing:** Custom/quote-based, no public pricing. Positioned for fintech and financial services companies. Given their enterprise orientation and the fact that Axiom would be using them to connect audit clients' accounting systems (not direct Axiom customers), the commercial model may not fit cleanly. Codat's model charges per linked company — in Axiom's case, that means charging per audit client accounting system connection, which could become expensive as the client base grows.

**Verdict on Codat for Axiom:** Codat solves the integration engineering problem (building 10 separate accounting connectors) but is optimized for a different use case (lending). The data model is close to what audit needs but not purpose-built for it. The commercial model (per linked company) could become expensive at scale. Codat is worth evaluating as an accelerator for the accounting integration layer, but Axiom should not architect itself as dependent on Codat — the integration should be encapsulated such that direct connectors can replace Codat later if economics warrant it.

**Apideck accounting coverage:** 22+ accounting platforms, including QuickBooks, NetSuite, Sage, and Xero. Consumer-based pricing (not per-connection), which is more predictable for Axiom's use case. Real-time proxy (no caching), which is better for audit point-in-time data. The accounting coverage is narrower than Codat but the pricing model is more favorable.

**Merge.dev for accounting:** Merge covers accounting as one of eight verticals. Merge caches data, which raises the same point-in-time completeness concern as Codat. Merge's strength is in HRIS/ATS — its accounting coverage is not as deep as Codat's. However, Merge is more relevant for the HR/payroll data needed for SOC 2 access review evidence (employee lists, termination dates) — see section 2.

---

## 2. Compliance Audit Integrations: SOC 2, ISO 27001, HIPAA

### The Evidence Ecosystem for SOC 2

SOC 2 evidence collection is highly system-specific. The auditee (the company being audited) uses specific cloud infrastructure and SaaS tools, and auditors need to collect screenshots, exports, logs, and configurations from those systems to test controls. The key systems and what auditors actually request:

**AWS (Amazon Web Services):**
- IAM users list and permissions — evidences that access is appropriately restricted (CC6.1)
- CloudTrail logs — authentication events, configuration changes, API calls — evidences monitoring and detection (CC7.1, CC7.2)
- S3 bucket configuration (public access block settings, encryption settings) — evidences data protection controls (CC6.7)
- VPC configuration, security groups — evidences logical access controls (CC6.6)
- CloudWatch alarms and Config rules — evidences monitoring (CC7.1)
- Multi-factor authentication enforcement configuration
- Typically delivered as: screenshots of console + CSV export of IAM reports + CloudTrail export (S3 download or SIEM query)

**GCP (Google Cloud Platform) and Azure:** Analogous to AWS. GCP IAM roles, Cloud Audit Logs, Azure AD conditional access policies, Azure Security Center configurations. The specific evidence items mirror AWS with platform-specific naming.

**Okta (Identity Provider):**
- User list and group memberships — evidences access provisioning (CC6.2)
- MFA enrollment status — evidences authentication controls (CC6.1)
- Application access policies — evidences access restriction (CC6.1)
- Audit logs (login events, user lifecycle events) — evidences monitoring (CC7.2)
- Administrator list — evidences privileged access management
- Typically delivered as: Okta admin console screenshots + CSV exports of user reports

**GitHub / GitLab:**
- Repository access list — who has access to production code (CC6.1)
- Branch protection rules — evidences change management (CC8.1)
- Pull request review requirements — evidences four-eyes principle for deployments (CC8.1)
- Audit log of permission changes (last 90 days — must be exported to external system because GitHub only retains 90 days)
- Dependency review / security scanning configuration
- Typically delivered as: screenshots + exported audit logs + CSV of collaborator list

**Jira / Linear / project management:**
- Change tickets — evidences that changes went through an approval process (CC8.1)
- Ticket lifecycle (created → approved → deployed) — evidences change authorization
- Typically delivered as: exported ticket data or screenshots of specific change tickets

**HR / Payroll (Gusto, Rippling, BambooHR, Workday, etc.):**
- Active employee list at period end — evidences that access reviews are based on accurate headcount
- New hire and termination dates — evidences that access is provisioned/deprovisioned timely (CC6.2, CC6.3)
- Offer letter templates / onboarding acknowledgment — evidences security training and background check controls (CC1.3, CC1.4)
- Typically delivered as: HR system exports or screenshots

**MDM / Endpoint Management (Jamf, Kandji, Intune):**
- Device enrollment list and compliance status — evidences endpoint security controls (CC6.8)
- OS version and encryption status per device
- Typically delivered as: MDM console screenshot or CSV export

**SIEM / Log Management (Splunk, Datadog, SumoLogic):**
- Alerting rules configuration — evidences monitoring program (CC7.1)
- Alert response evidence (tickets, incident records) for sampled events
- Typically delivered as: configuration screenshots + incident record exports

### Automated vs. Manual Evidence Collection

**The key distinction for Axiom:** Drata and Vanta (auditee-side platforms) do automated, continuous evidence collection directly from the client's systems — they hold OAuth connections into AWS, Okta, GitHub, Jamf, HR systems, and pull evidence automatically. This is 375+ integrations in Vanta's case.

**Axiom's role is different.** Axiom is on the auditor side. The evidence flow is:
1. Client connects their systems to Drata/Vanta/Sprinto (or manages evidence manually)
2. Evidence is available to the auditor via the auditee platform (e.g., Drata Audit Hub) or exported as files
3. The auditor receives and organizes evidence in Axiom

This means Axiom does not need to build 375 integrations into client infrastructure systems. That is the auditee platform's problem. Axiom's integration scope is:
- Receiving evidence packages from auditee platforms (Drata, Vanta, Sprinto exports)
- Receiving file uploads directly from clients
- Pulling files from client cloud storage (Dropbox, Box, Google Drive, SharePoint) where evidence was placed
- For firms that do not have clients on Drata/Vanta, managing the PBC request and document upload workflow directly

**The Drata/Vanta/Sprinto connection is high-value but narrow.** Drata already has a Fieldguide integration: auditors work in Fieldguide, evidence syncs from Drata without auditors needing a Drata login. This is the right architecture pattern for Axiom. An API connection to Drata's Audit Hub (and Vanta's equivalent) would allow Axiom to pull structured evidence directly into engagements, with evidence already tagged to SOC 2 controls. This is more valuable than any single infrastructure integration (AWS, Okta, etc.) because it abstracts the entire auditee stack.

**Prioritized compliance audit integrations for Axiom:**
1. Drata Audit Hub API (or bulk export) — single integration covering all Drata-connected client systems
2. Vanta audit evidence export — similar benefit
3. Direct file upload — covers clients not using any compliance automation platform
4. Gusto / BambooHR / Rippling via Merge.dev — for employee list data when clients do not use Drata/Vanta

### What Yak Integrates With

Yak is built specifically by auditors for auditors doing SOC 1/2/HIPAA. Based on available information, Yak manages the control mapping, document request collection, review cycles, and client communication workflow — but it functions as a document-request and workpaper management tool, not as an infrastructure integration hub. Yak does not appear to offer direct API connections into AWS, Okta, or GitHub from the auditor side. Its evidence collection approach is closer to a structured PBC request portal than an automated pull from client infrastructure. This is consistent with its positioning as a tool that eliminates "yak shaving" (manual coordination tasks) rather than automating evidence extraction from source systems.

---

## 3. Document Storage Integrations

### How Mid-Market Clients Actually Deliver Evidence

The practical reality for a 20–60 person CPA firm:

**The most common evidence delivery methods, in order of frequency:**
1. **Email with attachments** — Still used, especially by smaller clients. Universally problematic: hard to track, no audit trail, difficult to version-control.
2. **Dedicated PBC portal (Suralink, AuditDashboard, DataSnipper UpLink)** — Increasingly common. Clients upload directly to a structured request list. No need for cloud storage integration.
3. **SharePoint / Teams** — Common for clients using Microsoft 365. A shared folder is created for the engagement.
4. **Google Drive** — Common for clients using Google Workspace. A shared folder is set up.
5. **Dropbox / Box** — Less common than Microsoft/Google but still present, particularly in technology-adjacent client organizations.

**The key insight:** The shift away from email is driven by the auditor, not the client. Firms that adopt PBC portals (Suralink, AuditDashboard) have largely eliminated email for document collection. But the majority of mid-market firms have not yet done this — they are still managing PBC requests via email and shared drives.

**Implication for Axiom:** A well-built, native PBC request and direct-upload experience eliminates the need for cloud storage integrations in most cases. A client who gets a link to upload documents directly into Axiom will often prefer this to being asked to put files in a Dropbox folder. The friction of setting up a new Dropbox shared folder for each engagement is actually higher than using a dedicated upload portal.

### Are Cloud Storage OAuth Integrations Worth Building?

**Dropbox, Box, Google Drive, and SharePoint OAuth integrations have legitimate value but are not blockers:**

**When they matter:** Some audit clients will already have organized their evidence in a Google Drive or SharePoint folder as part of their SOC 2 preparation or year-end close process. An auditor who can connect directly and pull files from that folder saves a round-trip of "please export these 40 files to our portal."

**The setup friction for OAuth integrations:**
- Dropbox, Google Drive, and Box have well-documented OAuth 2.0 flows and are relatively straightforward to implement for a developer
- SharePoint (Microsoft 365) OAuth is significantly more complex — it involves Azure AD app registration, tenant-specific permissions, and often requires IT admin consent, not just user consent. This is the hardest of the four to implement and maintain
- Token refresh, revocation handling, and per-client credential storage add ongoing maintenance burden
- For a mid-market audit firm, each audit client connection is a separate OAuth grant — if a firm has 100 clients, that is potentially 100 separate cloud storage connections to manage

**Verdict on cloud storage integrations:** Google Drive and Dropbox are worth building at the mid-term (not MVP) — they have large coverage and relatively clean OAuth implementations. SharePoint requires significantly more engineering and maintenance per the complexity of Microsoft's OAuth model, but is strategically important because Microsoft 365 is dominant in mid-market firms' own IT stacks. Box is lower priority — it has legitimate enterprise presence but smaller mid-market footprint than Google Drive or SharePoint.

**The key point:** Axiom's direct upload experience should be good enough to make cloud storage integrations optional, not required. A polished drag-and-drop bulk upload with folder-structure recognition and automatic request matching substantially reduces demand for cloud storage integrations.

---

## 4. Minimum Integration Set for MVP

### What Must Be True at Launch

The MVP integration set is the smallest set of integrations that allows the target ICP (20–60 person firm doing SOC 2 and financial audits) to run a real engagement without being blocked by an integration gap.

**Financial audit MVP:** The primary data input is a trial balance, delivered as a CSV or Excel file. The vast majority of mid-market financial audits begin with the client exporting a trial balance from QuickBooks, NetSuite, or Sage and emailing it or uploading it to the audit firm. This is current practice universally — even firms using Fieldguide receive trial balances as file uploads. Axiom needs:
- CSV and Excel file import for trial balance (structured data model from Task 4)
- PDF import for financial statements and support documents
- Bulk file upload for GL transaction detail (Excel/CSV format)

No accounting system API integration is needed at MVP for financial audit. A well-designed import experience that handles common formats (QBO report export, NetSuite CSV export, Sage Intacct export) covers the first 6–12 months.

**Compliance audit MVP:** SOC 2 and ISO 27001 engagements are fundamentally document-collection and control-testing exercises. The auditor sends document requests; the client provides evidence; the auditor reviews and links evidence to controls. This flow works entirely with:
- Native document request management (PBC portal, built into Axiom)
- Direct file upload by clients (no login required)
- File attachment by auditors

No infrastructure integration (AWS, Okta, GitHub) is needed at MVP because Axiom is on the auditor side, not the auditee side. Clients either upload directly or the auditor uploads files received from the client.

**What can be solved with good file upload instead of native integration:**
- Accounting system data (trial balance, GL export) → structured CSV/Excel import
- Cloud infrastructure evidence (AWS IAM export, Okta user report) → file upload
- HR system data (employee list) → file upload or CSV import
- Client cloud storage (Dropbox, Google Drive, SharePoint) → direct upload portal with drag-and-drop

**The one integration that genuinely matters at launch:** Email notifications. Clients receive document request notifications by email and upload via a tokenized link. This is table stakes — without it, the PBC workflow is broken. SMTP / transactional email (SendGrid, Postmark) is the one external integration that is truly required at MVP.

---

## 5. Public API and Webhooks

### What Audit Firms Would Actually Use a Public API For

The audit firm's primary workflow is inside the platform — work is done in Axiom directly. A public API would primarily serve:

**Practice management system (Karbon, TaxDome, Financial Cents) synchronization:** Mid-market CPA firms run Karbon or TaxDome as their central work management and billing system. An engagement exists as a project/job in Karbon and as an engagement in Axiom. Firms want bidirectional sync:
- Create an Axiom engagement when a new job is created in Karbon
- Sync engagement status changes from Axiom back to Karbon (so partners see completion status in their practice management tool)
- Sync client data (contact details, entity names) so it does not need to be entered twice

Karbon has a rebuilt API with webhooks and supports integrations with other tools. TaxDome has an integration ecosystem. This is a real, firm-expressed pain point — managing client data in multiple systems is a documented friction in mid-market firms.

**Client onboarding automation:** Some larger firms have intake processes where a new engagement triggers multiple downstream tasks (create Axiom engagement, send welcome email, set up shared folder, assign team). A public API enables this automation via Zapier/Make workflows that not every customer can build but sophisticated ones will.

**BI and reporting:** Partner-level users want engagement status, workpaper sign-off rates, and turnaround time data in their own dashboards. A read-only API that exposes engagement status and milestone timestamps enables this without a native reporting dashboard.

**Audit evidence retrieval:** Some clients or their IT teams may want to programmatically deliver evidence rather than manually uploading. An API endpoint for evidence submission (POST /evidence-items) enables this.

### Webhook Events That Would Be Most Useful

Based on the workflows above, the highest-value webhook events:

| Event | Receiver Use Case |
|---|---|
| `engagement.status.changed` | Update Karbon/TaxDome job status; trigger billing workflows |
| `document_request.overdue` | Trigger escalation in practice management tool; send partner alert |
| `document_request.fulfilled` | Confirm client delivery in external tracker; trigger next-step workflow |
| `workpaper.signed_off` | Update engagement completion tracking; trigger EQR assignment |
| `engagement.finalized` | Trigger assembly deadline countdown in practice management; archive in DMS |
| `evidence_item.uploaded` | Trigger AI analysis queue; notify assigned auditor |
| `review_note.opened` | Notify preparer in Slack/Teams (via Zapier bridge) |
| `review_note.resolved` | Update review completion dashboard |

### Practice Management Tool Priority

**Karbon** is the most important practice management integration for the target ICP. Karbon is dominant in the 20–200 staff accounting firm segment in the US and Canada. It has a mature API, existing integrations with Xero, Ignition, and others, and firms actively express desire to connect their audit tooling to Karbon. A Karbon integration (either native or via their API) is a meaningful competitive advantage — Fieldguide does not currently offer one.

**TaxDome** is more common in tax-focused practices. Firms doing primarily tax + audit would find TaxDome integration valuable, but for firms where audit is the primary service, Karbon is more relevant.

**Financial Cents** is smaller and oriented toward bookkeeping and tax firms. Lower priority for Axiom's ICP.

---

## 6. Integration Aggregator Evaluation

### Codat (Accounting Data)

| Dimension | Assessment |
|---|---|
| **Coverage** | 35+ accounting platforms; QuickBooks Online and Desktop, Sage Intacct, NetSuite, Xero, Wave, MYOB, FreshBooks, Microsoft Dynamics |
| **Data types** | Accounts, journal entries, transactions, balance sheet (computed), P&L (computed), bank data |
| **Trial balance support** | Not a native endpoint — requires reconstruction from account/transaction data |
| **GL transaction depth** | Available for most platforms; journal entry data added transactional currency support in 2024 |
| **Audit-specific fit** | Low-medium. Built for lending/fintech. No audit-specific data product. Data completeness at point-in-time (vs. continuous) is uncertain. |
| **Pricing** | Per linked company, custom/quote. Enterprise-oriented pricing model may be expensive at scale for audit use case. |
| **Build-time savings** | High — replaces 8–10 direct accounting platform integrations |
| **Recommendation** | Evaluate for Phase 2 (first 6 months post-launch) as an accelerator. Do not depend on it at MVP. Encapsulate behind an abstraction layer. |

### Merge.dev (HRIS/HR)

| Dimension | Assessment |
|---|---|
| **Coverage** | 200+ HRIS, ATS, payroll, accounting, CRM, file storage integrations. HRIS coverage: BambooHR, Gusto, Rippling, Workday, ADP, Paylocity, Paycor, UKG, Namely, TriNet |
| **Relevant data for audit** | Employee list (active/inactive status), hire dates, termination dates, departments, managers |
| **SOC 2 use case** | Pull employee list for access review testing — verify that departed employees' access was revoked (CC6.2, CC6.3) |
| **Data model** | Caches data (syncs periodically). For access review purposes (quarterly user list reconciliation), this is acceptable. |
| **Gap** | No QuickBooks Payroll or ADP Run support — but Gusto, Rippling, BambooHR, and Workday cover the majority of the ICP's client base |
| **Pricing** | Consumer-based (per customer/linked account). Free tier for development. Enterprise plan for production. Not publicly disclosed. |
| **Compliance** | SOC 2 Type II, ISO 27001, HIPAA, GDPR certified — appropriate for handling HR data |
| **Recommendation** | High value for SOC 2/HIPAA compliance audits where employee list reconciliation is a test procedure. Worth evaluating for Phase 2. Replaces 6–8 direct HRIS integrations with one. |

### Apideck (Multi-Category)

| Dimension | Assessment |
|---|---|
| **Coverage** | 190+ connectors across accounting, CRM, HRIS, file storage, ecommerce |
| **Architecture** | Real-time proxy (no data caching) — better for audit point-in-time snapshots than Merge's caching model |
| **Accounting depth** | 22+ accounting platforms. Slightly narrower than Codat but covers the key players |
| **Pricing** | Consumer-based pricing (as of Jan 2026). Launch/Growth/Enterprise tiers. 30-day free trial. Roughly $500–$2,000/month for Growth. More predictable than Codat's per-connection model. |
| **File storage unification** | Covers Box, Dropbox, Google Drive, OneDrive — one integration covers all four cloud storage platforms |
| **Recommendation** | High value specifically for **file storage unification** — a single Apideck file storage integration surfaces files from Dropbox, Box, Google Drive, and OneDrive without four separate OAuth implementations. This is where Apideck's multi-category approach is most useful for Axiom. |

### Aggregator Trade-off Summary

| Use Case | Best Option | Rationale |
|---|---|---|
| Accounting system trial balance / GL data | Codat (or direct QBO/Xero connectors) | Codat has deepest accounting coverage; direct connectors for the two most common platforms (QBO, Xero) may be more cost-effective |
| Employee list / HR data for access reviews | Merge.dev | Deepest HRIS coverage; built specifically for this use case |
| Cloud storage file access | Apideck | Single integration covering all four major cloud storage platforms |
| Infrastructure evidence (AWS, Okta, GitHub) | Not needed — receive from Drata/Vanta exports or direct file upload | Axiom is auditor-side; auditee-side connectors are not Axiom's responsibility |

---

## 7. Prioritized Integration Roadmap

### Tier 1: Must-Have at Launch

These are not features — they are prerequisites. Without them, the target ICP cannot run an engagement in Axiom.

| Integration | Type | Rationale |
|---|---|---|
| **Transactional email (SendGrid / Postmark)** | Outbound notifications | Document request notifications, client portal invitations, review alerts. Without this, the PBC workflow does not function. |
| **Direct file upload (no integration required)** | Native UX | Clients upload evidence directly via tokenized link. Covers the majority of evidence delivery cases without any external integration. Must handle CSV, Excel, PDF, ZIP with bulk upload. |
| **CSV/Excel trial balance import** | Native UX | QBO, NetSuite, Sage, and Xero all export trial balances as CSV/Excel. A well-designed importer that handles common format variations covers financial audit for MVP. |
| **O365/Google SSO (SAML/OAuth)** | Identity | Mid-market firms authenticate via Microsoft or Google. Requiring separate username/password credentials for every staff member is a significant onboarding barrier. |

### Tier 2: High-Value in First 6 Months Post-Launch

These integrations unlock meaningful workflow improvements and competitive differentiation. They are not blockers at launch but should be built within the first 6 months based on early customer feedback confirming priority.

| Integration | Type | Priority Rationale |
|---|---|---|
| **Drata Audit Hub API** | Evidence ingestion | The single highest-value compliance audit integration. Drata is the market leader in SOC 2 compliance automation. A Drata → Axiom connection means structured, control-tagged evidence delivered directly into engagements, replacing hours of manual evidence organization. Fieldguide already has this integration — Axiom needs it to be competitive in the SOC 2 audit segment. |
| **Vanta evidence export** | Evidence ingestion | Analogous to Drata. Vanta has 375+ source integrations and a large installed base. Covering Drata + Vanta covers the majority of well-prepared SOC 2 audit clients. |
| **Google Drive** | Cloud storage | Highest-coverage cloud storage integration after SharePoint. Many mid-market clients use Google Workspace and will have evidence organized in Drive folders. OAuth implementation is relatively clean. |
| **Karbon API** | Practice management | The most important practice management integration for the ICP. Engagement sync eliminates double-entry in Karbon + Axiom, which is the daily friction point for any firm running both tools. A Karbon integration is a clear differentiator over Fieldguide, which does not offer one. |
| **QuickBooks Online direct connector** | Accounting data | QBO is the most common accounting system for the clients being audited by the ICP. A direct QBO connector (pull trial balance + journal entries via API on authorization) reduces the "ask client to export a file" step. Lower priority than file import because file import works, but the API connection improves data quality (complete, point-in-time, no manual steps). |
| **Merge.dev HRIS layer** | HR data | Enables employee list pull for access review testing in SOC 2/HIPAA engagements. Covers Gusto, Rippling, BambooHR, Workday via one integration rather than six. Directly supports the access provisioning/deprovisioning testing that accounts for a significant portion of SOC 2 fieldwork. |

### Tier 3: Longer-Term (6–18 Months)

These integrations are valuable and have a clear use case but are not urgent at launch or in the first 6 months. They represent the second expansion wave, informed by which engagements the firm is winning.

| Integration | Type | Rationale / Dependencies |
|---|---|---|
| **SharePoint / Microsoft 365** | Cloud storage | Strategically important (Microsoft is dominant in mid-market firm IT stacks) but technically complex — Azure AD app registration, tenant admin consent, significantly more maintenance burden than Google Drive. Build after Google Drive is proven. |
| **Dropbox** | Cloud storage | Meaningful coverage in technology-adjacent clients. Simpler OAuth than SharePoint. Lower priority than Google Drive because Google Workspace has larger mid-market penetration. |
| **Box** | Cloud storage | Enterprise storage with real mid-market presence in regulated industries. Lower priority than Dropbox and Google Drive by installed base size in the ICP's client population. |
| **NetSuite direct connector** | Accounting data | High-complexity integration. NetSuite API requires NetSuite-specific knowledge and accounts are often heavily customized. By the time a CPA firm is auditing clients on NetSuite, those clients likely have the sophistication to provide file exports. However, at the upper end of mid-market (clients with $50M+ revenue), this becomes a meaningful differentiator. Build after QBO is proven. |
| **Sage Intacct direct connector** | Accounting data | Sage Intacct's XML API is well-documented and Sage Intacct is the dominant system for SaaS-model companies and nonprofits — a significant portion of SOC 2 audit clients. Medium-complexity build. |
| **Xero direct connector** | Accounting data | Xero has a clean REST API including a native trial balance endpoint. Canadian market and smaller clients. Lower volume than QBO but clean implementation. |
| **TaxDome API** | Practice management | Relevant for firms that use TaxDome as their practice management tool. Lower priority than Karbon for the audit-focused ICP but relevant for mixed tax/audit practices. |
| **Public API + webhooks (Axiom-outbound)** | Platform extensibility | Enables Karbon workflows, Zapier/Make automation, and custom firm integrations. The webhook event list in section 5 above covers the most valuable events. Required before any firm can build custom workflows on top of Axiom — important once the firm base grows beyond early adopters. |
| **Sprinto / Hyperproof evidence export** | Evidence ingestion | Analogous to Drata/Vanta but smaller installed base. Once Drata and Vanta are live, adding Sprinto and Hyperproof is incremental effort for expanded coverage. |
| **AWS CloudTrail / IAM export** | Direct infrastructure | Only relevant if a significant portion of audit clients do not use Drata/Vanta but their auditors want to pull evidence directly. Given the Drata/Vanta abstraction covers most well-prepared clients, direct AWS integration is a niche need. Evaluate based on firm demand signals. |

---

## 8. Integration Architecture Principles

These principles should govern integration design from the start to avoid costly rewrites later:

**1. Abstract every integration behind an internal interface.** Whether trial balance data comes from a CSV upload, a QBO API call, or Codat, it should enter the same `TrialBalance` data model. The integration is a data source, not an architectural dependency. This allows replacing Codat with a direct connector (or vice versa) without touching audit logic.

**2. Store OAuth credentials per firm + client, not globally.** Each audit firm has its own client connections. A firm connecting a client's Google Drive is different from a different firm connecting the same client. Multi-tenant credential storage with firm-scoped OAuth grants.

**3. Evidence from integrations is not automatically trusted.** Files pulled from a client's Dropbox or Google Drive must go through the same AI extraction, OCR, and review workflow as direct uploads. Integration source should be recorded on the `EvidenceItem.source_integration` field (already in the data model) but not treated as pre-validated.

**4. Graceful degradation beats hard dependency.** If Drata's API is down, the auditor should be able to fall back to file upload without blocking the engagement. No integration should create a hard dependency that stops audit work.

**5. Scope integrations to read-only at launch.** The risk profile of writing data back into a client's QuickBooks or NetSuite is high. No proposed integration requires write-back at launch. Read-only OAuth scopes reduce risk, simplify consent flows, and make security reviews easier.

---

## 9. Gaps in the Original Spec

The spec listed Jira, Gusto, AWS, Dropbox, Box, Google Drive, and O365 as integrations. This research identifies the following adjustments:

| Original Spec Item | Assessment |
|---|---|
| **Jira** | Relevant for SOC 2 change management evidence, but lower priority than Drata/Vanta abstractions. Move to Tier 3. |
| **Gusto** | Relevant for HRIS data (employee list for access reviews). Better addressed via Merge.dev HRIS layer than direct Gusto integration — Merge covers Gusto plus 200+ others. Keep but address via aggregator. |
| **AWS** | As discussed: Axiom is auditor-side. Direct AWS integration is not needed when Drata/Vanta abstractions cover AWS evidence. Move to Tier 3 / evaluate-on-demand. |
| **Dropbox** | Genuine value but not Tier 1. Move to Tier 3. |
| **Box** | Genuine value but not Tier 1. Move to Tier 3. |
| **Google Drive** | Genuine value — move to Tier 2 (ahead of Box/Dropbox). |
| **O365 (SharePoint)** | SSO/SAML (O365 login) is Tier 1. SharePoint file storage integration is Tier 3 due to complexity. Split into two separate items. |

**Critical gaps in the original spec not addressed:**
- **QuickBooks Online / accounting system connectors** — The most important financial audit integration. Not in original spec. Tier 2 (QBO direct connector) / Tier 3 (NetSuite, Sage).
- **Drata / Vanta audit evidence integration** — The single highest-ROI compliance audit integration. Not in original spec. Tier 2.
- **Karbon / practice management** — The most important practice management integration for the ICP. Not in original spec. Tier 2.
- **Merge.dev HRIS layer** — Aggregator for employee list data. Not in original spec. Tier 2.
- **Transactional email** — Table stakes, not a "feature integration." Tier 1.

---

## Sources Consulted

- [QuickBooks Market Share: Global & Industry Insights](https://www.acecloudhosting.com/blog/quickbooks-market-share/)
- [Survey: Sage Intacct Tops NetSuite in Market Share Among SaaS Customers](https://www.sage.com/en-us/blog/sage-intacct-survey-saas-customers/)
- [QuickBooks vs. NetSuite: When SaaS Companies Should Upgrade](https://www.m3ter.com/blog/saas-accounting-systems-quickbooks-vs-netsuite-upgrade)
- [NetSuite vs QuickBooks: ERP Comparison and Migration Guide](https://www.houseblend.io/articles/netsuite-quickbooks-erp-evaluation)
- [When QuickBooks Hits Its Limits: Signs You've Outgrown QBO](https://www.eaglerockcfo.com/blog/quickbooks-guide/limitations-at-scale)
- [Codat Accounting Integrations Overview](https://docs.codat.io/integrations/accounting/overview)
- [Codat Supported Data Types](https://docs.codat.io/lending/data-types)
- [Guide to Accounting Integration in Financial Services — Codat](https://codat.io/blog/accounting-integration/)
- [15 Accounting APIs to Integrate With in 2026 — Unified.to](https://unified.to/blog/15_accounting_apis_to_integrate_with_in_2026_quickbooks_xero_freshbooks_and_unified_accounting_apis)
- [Apideck Accounting Connector Coverage](https://developers.apideck.com/apis/accounting/coverage)
- [Apideck Pricing](https://www.apideck.com/pricing)
- [Breaking Down Unified API Pricing — Apideck](https://www.apideck.com/blog/breaking-down-unified-api-pricing-why-api-call-pricing-stands-out)
- [Merge.dev HRIS Supported Features](https://docs.merge.dev/integrations/hris/supported-features/)
- [Best Unified APIs for HRIS and Payroll — Finch](https://www.tryfinch.com/blog/best-unified-apis-hris-payroll)
- [Merge vs Apideck — Truto](https://truto.one/blog/merge-vs-apideck-which-unified-api-is-better-in-2026)
- [Top Unified Accounting API Platforms of 2025 — Satva Solutions](https://satvasolutions.com/blog/top-unified-accounting-api-platforms)
- [AWS SOC 2 Compliance: What Auditors Actually Look For — DEV Community](https://dev.to/dannysteenman/aws-soc-2-compliance-what-auditors-actually-look-for-ghh)
- [SOC 2 Logging and Monitoring Requirements — SecurityDocs](https://security-docs.com/blog/soc2-logging-monitoring)
- [GitHub Configuration Checklist for SOC 2 — Delve](https://delve.co/blog/github-configuration-checklist-for-soc-2-compliance)
- [Drata x Fieldguide: Streamlining Audit Readiness](https://drata.com/blog/fieldguide-streamline-audit-readiness)
- [Understanding the Drata + Fieldguide Integration](https://help.drata.com/en/articles/11504009-understanding-the-drata-fieldguide-integration)
- [Drata Audit Hub](https://drata.com/product/audit-hub)
- [Vanta SOC 2 Compliance Automation](https://www.vanta.com/products/soc-2)
- [SOC 2 Compliance Platforms Compared 2026 — Cavanex](https://cavanex.com/blog/soc-2-compliance-platforms-compared-2026)
- [Fieldguide Integrations](https://www.fieldguide.io/product/integrations)
- [Fieldguide Integrations List](https://www.fieldguide.io/product/integrations-list)
- [Karbon Integrations](https://karbonhq.com/integrations/)
- [New Ignition and Karbon Integration — CPA Practice Advisor](https://www.cpapracticeadvisor.com/2025/07/08/new-ignition-and-karbon-integration-aimed-at-unifying-billing-and-workflow-automation/164520/)
- [5 Best PBC Software Tools for Audit Teams — DataSnipper](https://www.datasnipper.com/resources/best-pbc-software-tools-audit)
- [Guide to Secure File Sharing for Accountants — DataSnipper](https://www.datasnipper.com/resources/secure-file-sharing-for-accountants)
- [Top 5 File Storage APIs to Integrate With — Apideck](https://www.apideck.com/blog/top-5-file-storage-apis-to-integrate-with)
- [Sage Intacct API Integration — Merge.dev](https://www.merge.dev/blog/sage-intacct-api)
- [Financial Statement APIs — Apideck](https://www.apideck.com/blog/financial-statement-api)
- [Gusto Payroll API Capabilities](https://embedded.gusto.com/blog/payroll-api-capabilities-payroll-data-management/)
- [Yak Technologies About](https://yaktech.io/about/)
