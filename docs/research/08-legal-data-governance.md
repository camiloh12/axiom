# Legal & Data Governance Requirements for Axiom

**Research Date:** April 2026  
**Scope:** US and Canada primary launch markets; EU, UK, and Australia secondary/international clients  
**Platform type:** AI-native audit workpaper and engagement management SaaS for mid-market CPA firms (20–200 staff)

---

## 1. Immutable Archiving vs. GDPR/CCPA Deletion Rights

### The Core Conflict

Audit workpapers are subject to strict retention and immutability requirements under multiple regulatory frameworks:

- **AICPA AU-C Section 230** (private company audits): Audit documentation must be assembled and finalized no later than **60 days** after the report release date. Once the documentation completion date has passed, the auditor **must not delete or discard any audit documentation** before the end of its retention period. Minimum retention period: **5 years** from the report release date.
- **PCAOB AS 1215** (public company/issuer audits): Assembly deadline is **45 days** after the report release date; retention period is **7 years**. The standard explicitly states: *"Audit documentation must not be deleted or discarded after the documentation completion date."*
- **Sarbanes-Oxley Act Section 802**: Knowingly altering, destroying, or concealing audit records with intent to impede a federal investigation carries criminal penalties of up to **20 years imprisonment**. Audit firms auditing public companies must retain workpapers for **7 years** under SEC Rule 2-06 of Regulation S-X. Willful violation of the 5-year retention rule carries penalties of up to **10 years**.

These standards create hard legal obligations to retain and not modify finalized workpapers—directly in tension with privacy deletion rights.

### GDPR Resolution: Article 17(3)(b) Legal Obligation Exemption

GDPR Article 17(3) lists five circumstances where the right to erasure does not apply. The most directly applicable to audit records are:

- **Article 17(3)(b):** The right to erasure does not apply where processing is "necessary for compliance with a legal obligation which requires processing by Union or Member State law to which the controller is subject."  
- **Article 17(3)(e):** The right to erasure does not apply where processing is "necessary for the establishment, exercise or defence of legal claims."

For audit workpapers, the legal reasoning is straightforward: retaining an unchanged audit file is itself required by law (SOX, PCAOB AS 1215, AU-C 230, and equivalent national laws in the EU such as national company law transposing EU Accounting Directives). A data subject (e.g., an employee of the auditee whose name appears in workpapers) cannot compel deletion of records that the audit firm is criminally prohibited from deleting.

**Critical limitation:** The Article 17(3)(b) exemption is purpose-bound. Data retained under the legal obligation exemption must not be used for unrelated purposes (e.g., marketing, analytics). The same rule applies under CCPA (see below).

### CCPA/CPRA Resolution

Under California Civil Code Section 1798.105(d)(8), a business is not required to comply with a deletion request if retaining the information is "necessary to comply with a legal obligation." The California Attorney General and CPPA have confirmed this covers government-mandated record retention requirements. As with GDPR, data retained under this exemption may not be repurposed.

Unlike some state privacy laws, the CCPA does not create blanket exemptions for financial institutions—but it does provide the legal obligation carve-out that covers audit retention mandates.

### How Other Regulated Industries Resolve This

The general approach used across financial services, healthcare (HIPAA), and legal services is:

1. **Acknowledge the request** — Respond to the data subject that a deletion request was received.
2. **Invoke the legal retention exemption in writing** — Specify the applicable statute (SOX §802, PCAOB AS 1215, AU-C 230, GDPR Art. 17(3)(b)).
3. **Restrict the retained data** — Apply access controls so that retained data is not used for any purpose outside the legal retention obligation.
4. **Log the decision** — Maintain a record of the deletion request, the basis for refusal, and the date of expected deletion (i.e., when the retention period ends).
5. **Delete on schedule** — Once the mandatory retention window closes (5 or 7 years), delete the data as required by the privacy law.

### Practical Implications for Axiom

- **Platform architecture:** Workpapers must be locked (write-protected) after the assembly deadline. The platform should enforce immutability via versioned, append-only storage (e.g., AWS S3 Object Lock in COMPLIANCE mode, which prevents deletion even by root/admin). COMPLIANCE mode makes deletion impossible until the configured retention period expires, satisfying both the PCAOB immutability requirement and the GDPR/CCPA "cannot delete" position.
- **Deletion request workflow:** Build a documented process for responding to data subject deletion requests. The platform should generate a templated response letter citing Article 17(3)(b) or CCPA 1798.105(d)(8) and specifying the applicable retention deadline.
- **Data minimization:** Only capture personal data strictly necessary for audit documentation. Minimize names, identifiers, and personal data not required by audit standards to reduce the surface area of deletion conflicts.
- **Retention schedule:** The platform must enforce per-engagement retention schedules (5 years for private engagements, 7 years for public company/PCAOB engagements), with automated deletion triggered at schedule expiry.
- **Audit firm contract:** The Master Services Agreement (MSA) should state that the firm, as data controller, has determined that the Article 17(3)(b)/CCPA legal obligation exemption applies to finalized audit workpapers and has instructed Axiom (as processor) to retain such records despite deletion requests.

---

## 2. Data Residency Requirements by Jurisdiction

### United States (Primary Market)

No federal data residency law mandates that US-origin financial or audit data remain in the US. US-based cloud infrastructure (AWS us-east-1 or similar) is the default acceptable choice for US CPA firms auditing US clients. SOX and PCAOB standards require that workpapers be available for inspection by the PCAOB and SEC, which is achievable regardless of physical data center location within the US.

### Canada

**Governing law:** PIPEDA (Personal Information Protection and Electronic Documents Act) at the federal level; Quebec's Law 25 (Bill 64, effective September 2023) at the provincial level with stricter requirements.

**Data residency stance:** PIPEDA does **not** require data to reside within Canada. Canadian organizations may transfer personal data to the US or other countries provided:
1. They provide notice to individuals about potential cross-border transfers.
2. They ensure the foreign processor provides **comparable protection** (typically via contractual DPA/data processing agreement).
3. The originating Canadian organization remains **accountable** for the data even when processed abroad.

**Quebec Law 25 nuance:** Quebec now requires that a Privacy Impact Assessment (PIA) be completed before transferring personal information outside Quebec if the destination jurisdiction does not provide adequate privacy protections comparable to Quebec's. The US does not have a general adequacy determination from Quebec, so a PIA must be documented.

**Practical result for Axiom:** AWS us-east-1 is acceptable for Canadian clients under PIPEDA, provided:
- The MSA/DPA with the Canadian audit firm addresses cross-border transfer obligations.
- A Privacy Impact Assessment is documented for Quebec-based clients.
- Axiom can demonstrate PIPEDA-comparable protections (SOC 2 Type II is the standard mechanism).

### European Union

**Governing law:** GDPR (Regulation 2016/679); national implementing laws; EU Accounting Directive (2013/34/EU) and country-specific company laws for retention mandates.

**Data residency:** GDPR does not require data to reside within the EU/EEA, but it **prohibits transfers** of personal data to third countries (including the US) unless an adequate safeguard is in place.

**Available transfer mechanisms for US → EU data flows:**
1. **EU-US Data Privacy Framework (DPF):** The European Commission adopted an adequacy decision for the DPF on July 10, 2023. US companies that self-certify to the DPF can receive EU personal data without Standard Contractual Clauses (SCCs). The DPF survived its first legal challenge and is the simplest mechanism for qualifying companies. However, it remains politically fragile (Schrems III challenge is pending as of early 2026).
2. **Standard Contractual Clauses (SCCs):** The Commission adopted updated SCCs in June 2021 covering four transfer scenarios (controller-to-controller, controller-to-processor, processor-to-controller, processor-to-processor). These are the standard fallback. The European Commission announced a further consultation on SCC updates in Q4 2024, with revised SCCs expected in 2025.
3. **Binding Corporate Rules (BCRs):** Only practical for large multinationals; not relevant for Axiom at this stage.

**Practical result for Axiom:** For EU clients, Axiom should either (a) self-certify to the EU-US DPF and include DPF as the transfer mechanism in DPAs, or (b) incorporate the EU Standard Contractual Clauses (2021 version, controller-to-processor module) into every DPA with EU audit firm customers. Given DPF political risk, maintaining SCCs as a parallel mechanism is advisable.

**Preferred architecture for EU enterprise clients:** Offer an EU data residency option (AWS eu-west-1 Frankfurt or eu-central-1) as a product tier. This removes the transfer compliance question entirely and is a competitive differentiator for EU firms.

### United Kingdom

**Governing law:** UK GDPR (as retained in domestic law post-Brexit); Data Protection Act 2018; Data (Use and Access) Act 2025 (enacted June 19, 2025, amending UK GDPR).

**Transfer mechanism:** The US does **not** have a UK adequacy decision. UK-to-US transfers require either:
1. **International Data Transfer Agreement (IDTA):** The UK's replacement for EU SCCs, published by the ICO.
2. **UK Addendum to EU SCCs:** The EU SCCs with the ICO's standard UK addendum attached.

**Financial services note:** The UK Financial Services and Markets Act 2000 and Payment Services Regulations 2017 impose additional requirements on financial institutions, including requirements for UK-accessible audit trails for AML compliance. CPA firms auditing UK-registered entities may be subject to FRC (Financial Reporting Council) inspection rights, which require workpapers to be accessible in the UK.

**EU-UK adequacy:** The EU renewed its adequacy decisions for the UK on December 19, 2025, valid through December 27, 2031.

**Practical result for Axiom:** Include the UK IDTA or UK Addendum to EU SCCs in DPAs with UK audit firm clients. For UK enterprise clients with FRC regulatory exposure, consider an AWS eu-west-2 (London) data residency option.

### Australia

**Governing law:** Privacy Act 1988; Australian Privacy Principles (APPs), particularly APP 8 on cross-border disclosure; APRA Prudential Standard CPS 234 (for APRA-regulated entities).

**Data residency:** Australia does **not** mandate data localization for most sectors. However, before disclosing personal information to an overseas recipient (including a US-based cloud provider), the organization must take **reasonable steps** to ensure the recipient does not breach the APPs. The organization remains liable for breaches by the overseas recipient.

**APRA CPS 234:** Applies to APRA-regulated entities (banks, insurers, superannuation funds)—not to CPA/audit firms directly. However, if an audit firm's client is an APRA-regulated entity, the client may impose CPS 234-equivalent security requirements on its auditors and their software vendors by contract.

**Practical result for Axiom:** AWS us-east-1 is acceptable for Australian CPA firm clients, provided the DPA addresses APP 8 cross-border disclosure obligations. For Australian clients whose audit clients are APRA-regulated entities, offer an AWS ap-southeast-2 (Sydney) data residency option and document CPS 234-compatible security controls (aligned with ISO 27001 or SOC 2 Type II).

### Summary Table

| Jurisdiction | Legal Basis for US Cloud | Mechanism Required | Preferred Architecture |
|---|---|---|---|
| US | No restriction | DPA with firm | AWS us-east-1 |
| Canada (federal) | PIPEDA – no residency mandate | DPA with comparable protections + PIA for Quebec | AWS us-east-1 or ca-central-1 |
| EU | GDPR – transfer restriction | DPF self-certification or EU SCCs (2021) | AWS eu-central-1 option |
| UK | UK GDPR – transfer restriction | IDTA or UK Addendum to EU SCCs | AWS eu-west-2 option |
| Australia | Privacy Act – APP 8 accountability | DPA with APP 8 undertaking | AWS ap-southeast-2 option |

---

## 3. Contractual Structure: Data Processing Agreement (DPA)

### Controller vs. Processor Analysis

Under GDPR Article 4, the **controller** is the party that determines the purposes and means of processing. The **processor** processes data on the controller's behalf and only under the controller's instructions.

In the Axiom context, there are two distinct relationships:

**Relationship 1: Audit firm ↔ Auditee (client company)**  
The audit firm is the **controller** of personal data it collects about the auditee's employees, customers, and counterparties in the course of the audit. The auditee is not the controller of data the auditor collects independently pursuant to professional standards—though the auditee may assert separate control over data it provides to the auditor.

**Relationship 2: Audit firm ↔ Axiom**  
The audit firm (as controller) engages Axiom (as **processor**) to process audit data on its behalf within the platform. GDPR Article 28 mandates a written DPA between them.

**Axiom's role:** Axiom is a **data processor** (not a controller) with respect to customer data uploaded to the platform. Axiom determines the technical means of processing but not the purpose—the audit firm determines what data to upload and why. However, if Axiom processes any data for its own purposes (e.g., training AI models on customer data, analytics), it becomes a **joint controller or independent controller** for those purposes, which has significant compliance implications. This distinction must be addressed explicitly in the DPA and terms.

### Required DPA Elements Under GDPR Article 28

A DPA between Axiom (processor) and an audit firm (controller) must contain:

1. **Processing only on documented instructions** from the controller, including instructions on international transfers.
2. **Confidentiality obligations** — all authorized persons processing data are bound by confidentiality.
3. **Security measures** — Article 32-compliant technical and organizational measures (TOMs), including encryption at rest and in transit, access controls, and incident response procedures.
4. **Sub-processor management** — Prior written authorization required before engaging sub-processors; Axiom must impose equivalent obligations on sub-processors; Axiom remains fully liable for sub-processor performance. Standard practice: maintain a public sub-processor list with 30-day advance notice of changes.
5. **Data subject rights assistance** — Axiom must assist the audit firm in responding to data subject requests (access, correction, deletion where permitted, portability).
6. **Deletion or return at termination** — Axiom must delete or return all personal data at the controller's choice upon termination, and delete existing copies unless EU/member state law requires storage.
7. **Audit rights** — Axiom must make available information necessary to demonstrate compliance; allow and contribute to audits by the controller or its appointee. In practice, providing SOC 2 Type II or ISO 27001 certification satisfies most audit rights provisions.
8. **Breach notification** — Axiom must notify the audit firm without undue delay (in practice, within 24-72 hours) upon becoming aware of a personal data breach.

### Standard Contractual Clauses for International Transfers

Where Axiom processes EU or UK personal data on US infrastructure:

- **EU to US:** The 2021 EU SCCs (Module 2: Controller to Processor) must be incorporated into the DPA. If Axiom self-certifies to the EU-US Data Privacy Framework, the DPF adequacy decision may substitute for SCCs—but maintaining SCCs as a fallback is recommended given ongoing legal challenges.
- **UK to US:** The ICO's IDTA must be executed, or the EU SCCs with the UK Addendum attached.
- **Canada to US:** No SCCs required under PIPEDA. Contractual DPA with comparable protections and Quebec PIA documentation is sufficient.
- **Australia to US:** No formal SCC analog required. APP 8 compliance is satisfied by contractual undertakings in the DPA.

### Three-Party Complexity: Auditee as Data Subject

The auditee's employees and counterparties are **data subjects**—not parties to the Axiom-firm DPA. However:
- The audit firm (controller) is responsible for providing privacy notices to data subjects explaining that their data is processed in the context of an audit.
- The audit firm's engagement letter with the auditee should reference the use of technology platforms (naming Axiom or describing the category) so the auditee can discharge its own GDPR obligations to its employees.
- Axiom should provide the audit firm with a **data processing addendum template** suitable for inclusion in engagement letters between the firm and its auditee clients, covering the audit firm's role as controller and reference to Axiom as sub-processor.

### Practical Implications for Axiom

- **DPA template:** Axiom must publish and offer a GDPR-compliant DPA for all EU, UK, and Canadian firm customers. Clicking "I accept" on terms of service is not sufficient—a separately signed DPA is required.
- **No AI training on customer data:** The DPA must explicitly state that Axiom does not use customer data (uploaded workpapers, financial data, client lists) to train AI models without separate controller authorization. Any AI training on customer data would change Axiom's legal status from processor to controller for that purpose.
- **Sub-processor list:** Publish a current list of sub-processors (AWS, any AI model API providers, support tools) with a mechanism for 30-day advance change notifications.
- **Data Processing Addendum for firm-auditee engagements:** Provide template language that audit firm customers can incorporate into their engagement letters.

---

## 4. Liability for AI-Generated Workpaper Errors

### The Legal Landscape

**Auditor responsibility is non-delegable.** Under PCAOB AS 1220 (Engagement Partner Responsibility) and AICPA AU-C Section 220 (Quality Management), the engagement partner is responsible for the audit and its conclusions regardless of which tools were used. The auditor's signature on a report is a professional representation that the audit was conducted in accordance with applicable standards. An AI tool that assists the auditor does not transfer that professional responsibility to the software vendor.

**Vendor liability exposure is narrow but not zero.** The distinction that matters legally is:

- **Tools that assist professional judgment** (presenting evidence, flagging anomalies, organizing workpapers, drafting narrative for auditor review): The auditor reviews, approves, and signs off. Vendor liability for an error in AI output that the auditor reviewed and approved is extremely limited—the auditor's independent review breaks the causal chain.
- **Tools that substitute for professional judgment** (autonomously generating conclusions, automatically populating final workpaper sign-offs, making materiality determinations without human review): Vendor exposure increases significantly if the tool's output goes directly into final deliverables without meaningful human review. Courts and regulators may view this as the vendor participating in the delivery of professional services.

**The "professional services" carve-out.** Most audit software vendor contracts disclaim liability for professional services errors entirely. Standard language states that the software provides information only, does not constitute professional advice, and that the customer (the CPA firm) remains solely responsible for professional judgments and conclusions. These disclaimers are generally enforceable in commercial B2B contracts where both parties are sophisticated.

**Emerging regulatory attention (2025).** The FRC (UK Financial Reporting Council) issued AI-in-audit guidance on June 26, 2025. ICAS code of ethics revisions (effective July 1, 2025) include specific provisions on AI use. The PCAOB and AICPA CAQ published guidance on generative AI in auditing (April 2024). These frameworks universally maintain that auditors retain professional responsibility regardless of AI tool use, but they also require auditors to understand and evaluate AI-generated outputs—meaning auditors cannot blindly rely on AI-generated workpapers.

**A-LIGN, EvidencePrime, and other audit software vendors** universally position their platforms as tools that "support" rather than "replace" auditor judgment. This is both a marketing position and a legal defense strategy.

### Typical Vendor Contract Provisions

Standard enterprise SaaS audit software contracts contain:

1. **Disclaimer of professional services:** "The platform does not provide accounting, auditing, tax, legal, or other professional advice. All professional judgments are the sole responsibility of the licensed professional using the platform."
2. **Limitation of liability to fees paid:** Vendor liability capped at fees paid in the 12 months preceding the claim (or sometimes 3–6 months). This cap typically covers direct damages only.
3. **Consequential damages waiver:** Vendor disclaims liability for lost profits, loss of data, business interruption, or third-party claims resulting from errors in software output.
4. **Gross negligence/fraud carve-out:** Most enterprise contracts carve out gross negligence and willful misconduct from the limitation of liability.
5. **Data security super-cap:** An increasing number of enterprise contracts include a separate, higher liability cap (e.g., 2x–3x annual fees) specifically for data breaches or failure to maintain security.

### Practical Implications for Axiom

- **Product design:** Every AI-generated workpaper narrative, finding, or conclusion must require explicit human review and sign-off (a named reviewer, date, and approval action) before it is treated as finalized in the platform. This creates an auditable paper trail showing the auditor exercised judgment.
- **Never auto-finalize:** The platform must not automatically promote AI-generated content to "final" status without a required human review step. This is the single most important design decision for managing liability.
- **Contract language:** The MSA must include:
  - A clear disclaimer that Axiom does not provide professional audit, accounting, or legal services.
  - A cap on liability equal to 12 months of fees paid.
  - A consequential damages waiver.
  - A data security carve-out with a higher cap (recommend 2x annual fees).
  - A carve-out for gross negligence and fraud.
- **Audit trail for AI outputs:** Log which AI model generated each output, the version of the model, the date, and the prompts used. This is essential for post-incident investigation and demonstrates the AI was used as an assistive tool.
- **Terms of AI model use:** If using third-party AI models (e.g., OpenAI, Anthropic), ensure the API terms prohibit the model provider from using customer data for training, and reflect this commitment to customers in Axiom's DPA.

---

## 5. Insurance Requirements

### Technology Errors & Omissions (Tech E&O)

Tech E&O (also called Professional Liability for Technology Companies) covers claims arising from errors, omissions, or failures in the software that cause financial loss to customers. For an audit platform where errors could contribute to materially misstated financial statements, this is the most critical coverage.

**Recommended coverage for Axiom:**
- **At seed / pre-revenue:** $1M per occurrence / $1M aggregate is the minimum that satisfies most enterprise MSAs.
- **At Series A with enterprise customers:** $2M per occurrence / $2M aggregate is the standard requirement in enterprise financial services contracts. Some enterprise contracts (particularly those with Big 4 affiliates or large regional firms) require $5M.
- **For fintech/financial data handling:** WTW's 2025 FinTech E&O analysis identifies E&O as the second most common source of fintech insurance claims (behind cyber). Coverage amounts of $2M–$5M are standard for SaaS companies handling regulated financial data.
- **Annual premium estimate (2025/2026):** $3,000–$12,000/year for a $1M–$2M policy for an early-stage SaaS with no revenue history; $8,000–$30,000/year for $5M coverage with fintech risk profile.

### Cyber Liability Insurance

Cyber insurance covers first-party costs (breach response, forensics, notification, credit monitoring, business interruption) and third-party costs (customer lawsuits, regulatory fines and defense, PCI-DSS penalties).

**Recommended coverage for Axiom:**
- **Minimum at Series A:** $2M per occurrence / $2M aggregate.
- **Target for enterprise financial services:** $5M per occurrence, with sub-limits covering:
  - Breach response / forensics: $500K–$1M
  - Regulatory defense: $1M
  - Business interruption: 6–12 months revenue
  - Ransomware / extortion: $500K
- **Annual premium estimate:** $8,000–$25,000/year for $2M coverage; $20,000–$60,000/year for $5M coverage with financial data risk profile.
- **Key underwriting factors:** SOC 2 Type II certification, MFA enforcement, encryption at rest and in transit, incident response plan, employee security training.

### General Commercial Liability (CGL)

Standard CGL ($1M/$2M) is required for office leases and general contracting purposes. Less critical than Tech E&O and Cyber for a software-only business.

### Directors & Officers (D&O)

Not directly relevant to data governance but required by institutional investors at Series A. Typically $2M–$5M for an early-stage company.

### Insurance Procurement Timing

Secure Tech E&O and Cyber coverage before signing the first enterprise contract. Most enterprise MSA redlines require proof of coverage. Specialist brokers for early-stage SaaS with financial services exposure include Vouch, Corvus (a Travelers company), Corgi Insurance, and WHINS.

### Key Contract Requirement

Enterprise customer MSAs will require Axiom to maintain specified insurance and provide certificates of insurance on request. Axiom's standard MSA should state its coverage levels and commit to maintaining them throughout the contract term.

---

## 6. Data Ownership and Offboarding

### Legal Baseline: Customer Owns Its Data

The industry-standard position—codified in nearly every enterprise SaaS contract including Fieldguide's terms—is that **the customer owns all data it uploads to the platform**. The platform has a limited license to process the data for the purpose of providing the service. This principle is not negotiable with enterprise audit firm customers and should be the starting position in Axiom's MSA.

**Fieldguide's terms** explicitly state: "Customer Data is owned by the Customer. The Customer Agreement provides Customer with many choices and control over that Customer Data." Specific export and deletion procedures are handled in the MSA rather than public-facing terms.

### Industry Standard Offboarding Commitments

Based on analysis of enterprise SaaS contracts in regulated industries:

**Data export:**
- Customer must be provided a mechanism to export all customer data in **machine-readable, standard formats** (CSV, JSON, XML, PDF) within a reasonable window after termination.
- Export should include: workpapers, engagement files, sign-off logs, audit trails, configuration data, and metadata—not just raw records.
- **Timeline standard:** 30–90 days of continued platform access (or read-only access) following contract termination. A minimum 30-day export window is the acceptable floor; 60–90 days is customer-friendly and competitive.
- API access during the transition window is increasingly expected by enterprise customers.

**Data deletion:**
- After the export window, the vendor must **delete all customer data** from production systems, backups, and sub-processor systems.
- Written confirmation of deletion (a "deletion certificate") is increasingly required by enterprise and regulated customers.
- **Timeline for deletion:** Industry standard is 30–60 days after the export window closes (i.e., 60–150 days after termination). Some enterprise contracts require deletion within 30 days of termination.
- **Backup exception:** Many contracts allow retention in encrypted, isolated backup systems for up to 90 days after deletion from production, with a commitment to overwrite/purge on the next backup rotation cycle.

**Retention obligation exception:**
- If applicable law (e.g., SOX, PCAOB AS 1215, AU-C 230) requires the **firm** to retain records, the deletion obligation runs to the firm—not to Axiom. However, Axiom's DPA should clarify that upon termination, Axiom will delete its copies and the firm is responsible for maintaining its own copies in compliant archival storage.
- The DPA should address this explicitly: "Upon termination, Processor shall provide Customer with a data export and thereafter delete Customer Data from all Processor systems within 60 days, unless a legal obligation requires Processor to retain such data, in which case Processor shall notify Customer and restrict the data to the minimum necessary for compliance with such obligation."

### Practical Implications for Axiom

- **MSA provisions to include:**
  1. Customer owns all customer data.
  2. Axiom has a limited, purpose-restricted license to process data solely to provide the service.
  3. Axiom will not use customer data for AI training, benchmarking, or product improvement without explicit opt-in consent from the customer.
  4. Upon termination: 60-day read-only export window; deletion of all copies within 30 days of export window close; written deletion certificate provided to customer.
  5. Export formats: PDF (for rendered workpapers), CSV/XLSX (for structured data), JSON (for engagement metadata and API exports), and native format for all uploaded documents.
  6. Audit trail export: Complete sign-off history, user access logs, and version history must be included in the export package.

- **Product design:**
  - Build a self-service data export feature available at all times (not only at offboarding).
  - Design exports to include the full audit evidence chain (source documents, AI outputs, human review sign-offs, final workpaper) so the exported package can stand alone as an archival record.
  - Consider offering an "archival export" format specifically designed for long-term retention: a structured ZIP or PDF/A package with all evidence, metadata, and sign-off records that can be stored independently of any software platform.

- **Competitive positioning:** Fieldguide and similar platforms defer specific offboarding terms to the MSA, which creates negotiation friction with enterprise customers. Axiom can differentiate by publishing explicit, customer-friendly data portability and deletion commitments in public-facing documentation.

---

## Cross-Cutting Practical Recommendations for Platform Design

| Requirement | Platform Feature | Legal Basis |
|---|---|---|
| Immutable workpaper archiving | AWS S3 Object Lock (COMPLIANCE mode) with per-engagement retention periods | PCAOB AS 1215, AU-C 230, SOX §802 |
| Deletion request workflow | Automated response templates citing Art. 17(3)(b) / CCPA 1798.105(d)(8) | GDPR, CCPA/CPRA |
| EU/UK data residency option | AWS eu-central-1 / eu-west-2 deployment tier | GDPR, UK GDPR |
| AI output review enforcement | Mandatory human sign-off before finalization; no auto-approve | Professional liability risk management |
| DPA with every customer | Pre-signed DPA template including EU SCCs Module 2 | GDPR Art. 28, UK GDPR |
| Sub-processor transparency | Public sub-processor list with 30-day change notice | GDPR Art. 28(2) |
| No AI training on customer data | Explicit DPA prohibition; third-party model API terms reviewed | GDPR controller/processor distinction |
| Insurance | Tech E&O $2M, Cyber $2M–$5M at Series A | Enterprise MSA requirements |
| Data export on offboarding | 60-day read-only access; deletion certificate; PDF/CSV/JSON formats | Industry standard; GDPR Art. 28(3)(g) |

---

## Key Sources and Citations

- PCAOB AS 1215: Audit Documentation — https://pcaobus.org/oversight/standards/auditing-standards/details/AS1215
- AICPA AU-C Section 230: Audit Documentation — https://us.aicpa.org/content/dam/aicpa/research/standards/auditattest/downloadabledocuments/au-c-00230.pdf
- SOX Section 802 — https://www.sarbanes-oxley-101.com/SOX-802.htm
- GDPR Article 17 (Right to Erasure) — https://gdpr-info.eu/art-17-gdpr/
- GDPR Article 28 (Processor obligations) — https://gdpr-info.eu/art-28-gdpr/
- EU Standard Contractual Clauses (2021) — https://commission.europa.eu/law/law-topic/data-protection/international-dimension-data-protection/standard-contractual-clauses-scc_en
- EU-US Data Privacy Framework adequacy decision (July 2023) — https://ec.europa.eu/commission/presscorner/detail/en/ip_23_3721
- ICO International Transfers guidance (UK GDPR) — https://ico.org.uk/for-organisations/uk-gdpr-guidance-and-resources/international-transfers/international-transfers-a-guide/
- ICO IDTA — https://ico.org.uk/for-organisations/data-protection-and-the-eu/data-protection-and-the-eu-in-detail/the-uk-gdpr/international-data-transfers/
- OAIC APP 8: Cross-border disclosure of personal information — https://www.oaic.gov.au/privacy/australian-privacy-principles/australian-privacy-principles-guidelines/chapter-8-app-8-cross-border-disclosure-of-personal-information
- APRA CPS 234 — https://www.apra.gov.au/sites/default/files/cps_234_july_2019_for_public_release.pdf
- PIPEDA overview — https://www.priv.gc.ca/en/privacy-topics/privacy-laws-in-canada/the-personal-information-protection-and-electronic-documents-act-pipeda/
- CCPA deletion exemptions — https://www.clarip.com/data-privacy/ccpa-erasure-exemptions/
- CAQ: Auditing in the Age of Generative AI (April 2024) — https://www.thecaq.org/wp-content/uploads/2024/04/caq_auditing-in-the-age-of-generative-ai__2024-04.pdf
- AI vendor liability: Jones Walker LLP (2025) — https://www.joneswalker.com/en/insights/blogs/ai-law-blog/ai-vendor-liability-squeeze-courts-expand-accountability-while-contracts-shift-r.html
- WTW FinTech E&O Risk (2025) — https://www.wtwco.com/en-us/insights/2025/08/fintech-errors-and-omissions-e-and-o-risk-what-do-the-numbers-say
- Fieldguide Terms of Service — https://www.fieldguide.io/terms
- SaaS Capital: Insurance for Growth-Stage SaaS — https://www.saas-capital.com/blog-posts/what-insurance-and-how-much-should-a-growth-stage-saas-company-carry/
- Accountancy Europe: GDPR impact on auditors — https://accountancyeurope.eu/stories/gdpr-one-year-on-its-impact-on-auditors-and-accountants/
- UK Data (Use and Access) Act 2025 — https://www.globalprivacyblog.com/2025/07/uk-adequacy-holds-firm-under-new-data-use-and-access-act-2025/
