# Research Task 3: Regulatory and Compliance Standards Landscape

## Engagement Type → Governing Standard → Platform Requirements

### 1. Financial Audit — Nonissuers (Private Companies)

**Governing standard:** AICPA AU-C sections (codified in SAS No. 122 and subsequent amendments, effective through 2025/2026)

**What it governs:** The complete financial statement audit lifecycle — planning, risk assessment, evidence, substantive procedures, review, and reporting — for companies that are not publicly traded (i.e., not subject to PCAOB oversight).

**Mandatory platform requirements:**

| Requirement | Standard | Detail |
|---|---|---|
| Workpaper documentation | AU-C 230 | Each workpaper must document its purpose, source, and conclusions. Must demonstrate compliance with GAAS and support every assertion in the auditor's report. |
| Assembly deadline | AU-C 230 | Workpapers must be finalized (locked) within **60 days** of the report release date. No modifications after lockdown except via documented addendum. |
| Retention | AU-C 230 | **Minimum 5 years** from report release date for nonissuers. |
| Sign-off and review hierarchy | AU-C 220 | Engagement partner must provide evidence of "substantial and meaningful" involvement throughout the engagement. Supervisor review of each section must be documented. |
| Risk linkage | AU-C 315/330 | Each audit procedure must be linked to a specific identified risk and financial statement assertion. |
| Sampling documentation | AU-C 530 | Statistical or non-statistical sampling approach must be documented; population, sample size, selection method, and results must be recorded. |
| Quality management | SQMS 1 (eff. Dec 15, 2025) | Firm's quality risk responses must be documented at client acceptance. Platform must support risk-based quality documentation per the 8-component SQMS 1 framework (see section 5 below). |

---

### 2. SOC 1 and SOC 2 — Attestation Engagements

**Governing standard:** AICPA AT-C sections (SSAE No. 18, Clarification and Recodification); specifically:
- AT-C 105: Concepts common to all attestation engagements
- AT-C 205: Examination engagements (Type II SOC reports)
- AT-C 320: Reporting on an examination of controls at a service organization (SOC 2 specifically)

**Key distinction:** SOC 1 and SOC 2 are *attestation* engagements, not *audit* engagements. The practitioner is an "examiner," not an "auditor" in the AU-C sense. This affects what standards apply and what the engagement deliverable looks like, but the workpaper and documentation obligations are functionally similar.

**Mandatory platform requirements:**

| Requirement | Standard | Detail |
|---|---|---|
| Workpaper documentation | AT-C 105/205 | Same concept as AU-C 230: document purpose, source, conclusions for each procedure. |
| Assembly and retention | AT-C 105 | **60-day assembly** deadline after report issuance; **minimum 5-year retention**. |
| Period coverage (Type II) | AT-C 320 | Evidence must cover the **entire audit period** (commonly 6 or 12 months). Requires population listings, sample selections, timestamps, and proof of control operation across the period — not just point-in-time. |
| Trust Services Criteria version | AICPA TSC 2017 (rev. 2022) | The current applicable criteria are the 2017 Trust Services Criteria with 2022 revised Points of Focus. The 2022 update is non-breaking — it does not alter the underlying criteria, only guidance on points of focus. Engagements using 2017 criteria remain fully compliant. Platform must track which version of criteria was in effect for each engagement. |
| Control testing documentation | AT-C 205 | Each control test must document: control description, test procedure performed, sample selected, exceptions noted, conclusion. |
| Management assertion | AT-C 320 | The service organization's management must provide a written assertion. Platform should support the assertion workflow and store the signed document linked to the engagement. |

---

### 3. Financial Audit — Issuers (Public Companies)

**Governing standard:** PCAOB auditing standards (AS series). Separate from AICPA standards — PCAOB is the regulator for auditors of SEC-registered companies.

**Key difference from nonissuers:** PCAOB standards are stricter, longer retention, heavier documentation burden, and subject to PCAOB inspection. The target ICP (mid-market firms, 20–200 staff) primarily performs nonissuer audits. Public company audits are edge cases — most mid-market CPA firms have limited or no PCAOB-registered engagements.

**Mandatory platform requirements (where they differ from AICPA standards):**

| Requirement | Standard | Detail |
|---|---|---|
| Retention | AS 1215 | **Minimum 7 years** (vs. 5 years for nonissuers). |
| Assembly deadline | AS 1215 | **45 days** after report release date (vs. 60 days for nonissuers). |
| Technology-assisted analysis | AS 1105 (eff. Dec 15, 2025) | If audit procedures use technology-based analysis tools (including AI), auditors must: (1) evaluate the reliability of electronic information used as evidence; (2) test IT general controls over the source systems; (3) document that technology-assisted procedures met their intended purpose; (4) investigate and document all flagged transactions or balances. |
| AI output limitation | AS 1105 / PCAOB guidance | **AI outputs alone do not constitute sufficient appropriate audit evidence.** The auditor must independently corroborate AI-flagged findings. Platform must enforce that AI recommendations require human review and documented sign-off before inclusion in the audit file. |
| Engagement quality review | AS 1220 | EQR reviewer must be formally documented, independent of the engagement team, and their review must be documented before the report is issued. |

**Key upcoming change — AS 1215 effective December 15, 2026:**
Amended to explicitly address documentation requirements for technology-assisted audit procedures, including standardized structure requirements to facilitate AI and data analytics in audit workflows. Firms planning to use AI in PCAOB engagements should design their documentation around this forthcoming standard, not just the current one.

---

### 4. ISO 27001

**Governing standard:** ISO/IEC 27001:2022 (most recent revision, superseding 2013 edition). Certification body (e.g., BSI, Bureau Veritas) performs the audit; not a CPA firm engagement in most cases.

**Key structural difference from SOC 2:** ISO 27001 is an *information security management system* (ISMS) standard — it certifies the organization's management framework, not specific control operation over a period. The audit is against Annex A controls (93 controls in 2022 version) plus the ISMS requirements in clauses 4–10.

**Overlap with SOC 2:** Approximately 80% control overlap. A SOC 2 engagement covering Security + Availability + Confidentiality produces evidence that largely satisfies ISO 27001 Annex A requirements. This is the foundation for framework-agnostic evidence reuse.

**Platform implications:**
- ISO 27001 certification cycles differ from SOC 2: initial certification (Stage 1 + Stage 2 audit), then annual surveillance audits (partial scope), then full recertification every 3 years.
- Evidence retention for ISO 27001 is ongoing (continuous ISMS records), not point-in-time.
- Control mapping must account for the ISO 27001 clause/control numbering scheme (A.5.x through A.8.x in 2022 version) alongside SOC 2 Trust Services Criteria numbering (CC6.x, A1.x, etc.).

---

### 5. HIPAA

**Governing standard:** HIPAA Security Rule (45 CFR §§ 164.302–318) and Privacy Rule. Federal law — not a voluntary framework.

**Key structural difference from SOC 2 and ISO 27001:** HIPAA is a legal obligation, not a certification. There is no formal "HIPAA certification." A HIPAA audit assesses compliance with the Security Rule's required and addressable safeguards. Auditors are not performing an attestation under AT-C — they are conducting a compliance assessment under agreed-upon procedures or an internal audit framework.

**Platform implications:**
- Engagements should be structured as agreed-upon procedures (AT-C 215) or an internal audit, not a SOC-style examination.
- Evidence must demonstrate compliance with specific regulatory citations (§164.312(a)(1), etc.), not Trust Services Criteria.
- Control numbering and framework must accommodate HIPAA's administrative, technical, and physical safeguard structure.

---

### 6. Quality Management — SQMS 1 (Effective December 15, 2025)

**What it is:** AICPA's new risk-based quality management standard, replacing the prior quality control standard (QC 10). Applies to any CPA firm performing audit, review, or attestation engagements.

**Why it matters for platform design:** SQMS 1 is not just a firm policy document — it changes the *workflow* of engagements in ways that a platform must support.

**Platform requirements SQMS 1 creates:**

| Requirement | What the platform must do |
|---|---|
| Quality risk documentation at acceptance | Client/engagement acceptance workflow must prompt documentation of quality risks identified and the firm's responses. |
| Communicating QM policies to engagement teams | Platform must surface applicable firm policies during engagement setup, not just in a manual the team never reads. |
| Engagement Quality Review (SQMS 2) | Formal EQR workflow: assign reviewer, document reviewer independence, track review scope, require EQR sign-off before report issuance. EQR reviewer cannot be part of the engagement team. |
| Annual internal inspection | Platform should support inspection workflow: sampling prior engagements, documenting findings, tracking remediation. |
| Partner involvement evidence | Platform must create a timestamped record of engagement partner activity throughout the engagement, demonstrating "substantial and meaningful" involvement. |

---

### 7. AI Use in Audits — Regulatory Position Summary

| Regulator | Position | Key Platform Implication |
|---|---|---|
| **PCAOB** (public companies) | Formal standard (AS 1105, eff. Dec 15, 2025): AI outputs alone are not sufficient evidence. Must test IT controls over source data. Must document all AI-assisted procedures. | Human review gate before any AI output is accepted into the audit file. AI decision trail stored as first-class audit documentation. |
| **AICPA** (private companies, SOC) | No formal AI-specific standard yet. Encourages technology adoption. PCAOB guidance is increasingly used as reference even for nonissuer engagements. | Same practical implication: build human-in-the-loop into every AI workflow. Review states must be auditable. |
| **PCAOB AS 1215** (eff. Dec 15, 2026) | Will standardize documentation structure to "facilitate AI and data analytics" in audit workflows. | Platform documentation schema should anticipate this standard. Structured, machine-readable workpaper formats will become a compliance requirement for public company audits. |

**Critical platform design constraint:** Any AI-generated content in a workpaper must be distinguishable from auditor-authored content in the audit file. The platform must record: (1) what the AI suggested, (2) who reviewed it, (3) what action the auditor took (accepted, modified, rejected), and (4) when. This is an audit trail of the AI decision, not just the audit trail of the workpaper.

---

### 8. Framework Version Management

**How criteria updates work in practice:**

- **SOC 2 (AICPA TSC):** Last substantive criteria change was 2017. Points of focus were updated in 2022 (non-breaking — old criteria remain valid). The AICPA does not set a mandatory transition date for points-of-focus changes; the underlying criteria govern whether a report is valid. Future criteria revisions (if any) would likely follow the same backward-compatible pattern with a 1–2 year transition window.

- **ISO 27001:** Updated from 2013 to 2022 edition. IAF (International Accreditation Forum) set a **3-year transition period ending October 31, 2025**, after which all certificates must be against the 2022 standard. Platform must track which standard version applies to each engagement and prevent applying 2013-version control mappings to post-transition engagements.

- **HIPAA Security Rule:** Static since 2003 (minor updates). A regulatory update is under active HHS review as of 2024–2025 but not yet finalized. Platform should monitor for changes.

**Platform implication:** The framework version must be a first-class field on every engagement. Evidence linked to a control must reference the specific version of the criterion it satisfies (e.g., "SOC 2 TSC 2017, CC6.1" not just "CC6.1"). When criteria versions are updated, existing engagements must remain linked to the version in effect at the time of the engagement.

---

## Summary: Platform Requirements by Engagement Type

| Requirement | Financial Audit (Private) | SOC 1 / SOC 2 | Financial Audit (Public) | ISO 27001 | HIPAA |
|---|---|---|---|---|---|
| **Workpaper assembly deadline** | 60 days | 60 days | 45 days | N/A | N/A |
| **Retention period** | 5 years | 5 years | 7 years | Ongoing ISMS records | 6 years (HIPAA requirement) |
| **Sign-off hierarchy** | Partner + reviewer sign-off | Partner + EQR (if applicable) | Partner + mandatory EQR | Certification body | N/A |
| **Period coverage** | Year-end or stub period | Full examination period (Type II) | Year-end | Ongoing | Ongoing |
| **AI documentation required** | Best practice (PCAOB guidance as reference) | Best practice | Mandatory (AS 1105, Dec 2025) | N/A | N/A |
| **Immutable lock** | After 60-day window | After 60-day window | After 45-day window | After certification issued | N/A |
| **Quality management** | SQMS 1 (Dec 2025) | SQMS 1 (Dec 2025) | PCAOB QC 1000 | ISO internal audit | N/A |

---

## Sources Consulted

- AICPA AU-C Section 230, Audit Documentation
- AICPA AT-C Sections 105, 205, 320 (SSAE No. 18)
- PCAOB AS 1215: Audit Documentation (current and effective Dec 15, 2026)
- PCAOB AS 1105: Audit Evidence (amended, effective Dec 15, 2025)
- PCAOB speech: "AI and the Pursuit of Audit Quality: A Regulatory Perspective"
- AICPA SQMS 1 implementation guidance (CPA Journal, April 2025)
- AICPA Trust Services Criteria 2017 (revised points of focus 2022)
- ISO/IEC 27001:2022 transition guidance (IAF deadline Oct 31, 2025)
- ISMS.online: SOC 2 vs ISO 27001 vs HIPAA framework comparisons
- Linfordco: Trust Services Criteria SOC 2 audit guidance
