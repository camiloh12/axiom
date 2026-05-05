# Product Owner Feedback — April 29, 2026

Product-owner feedback on the Axiom product specification (`docs/specs/axiom-spec-design.md`) and the user journeys (`docs/user-journeys/all-journeys.md`). Original messages combined English and Spanish; all content has been translated to English and expanded into full sentences while preserving the original intent. Each item flags the spec section the feedback applies to.

---

## 1. Scope — ISO Certification Bodies and PCI QSA Firms Should Be In Scope

> **Spec passage being commented on (§1, Target ICP — Practice mix):**
> "Compliance and assurance work across SOC 2 Type I/II, ISO 27001:2022, ISO 27701:2019, ISO 42001:2023, HIPAA Security Rule (with HITRUST CSF r2 path post-MVP), PCI DSS v4.0.1, and SOC 1 Type I/II. **Not** Certification Bodies (ISO CBs do not fit the engagement-workspace model), **not** QSA firms (QSA ROC signing is out of scope), **not** internal audit / SOX / enterprise GRC (AuditBoard's territory), **not** financial audit."

**Feedback:** ISO Certification Bodies and PCI QSA firms **should be in scope**, contrary to the current exclusion. The product owner wants Axiom to serve these firm types rather than rule them out at the ICP level. (Related to Item 6 below, which proposes generating the certificate/ROC artifacts themselves from templates.)

**Affected spec sections:** §1 Target ICP and §4 Explicit Scope Boundary.

---

## 2. Competitor Set — Reorient Around Auditor-Side Players

**Feedback:** The real competitors are **Fieldguide, Agentive, and Yak** — all auditor-side tools. The current spec over-indexes on Drata and Vanta, but those are auditee-side platforms and are not Axiom's primary competition. The strategic focus should shift toward differentiating against the auditor-side competitor set.

**Affected spec sections:** §2 Competitive Positioning. Fieldguide and Yak are not currently in the competitor table and need to be added. Agentive is already listed but should be elevated. The "Switching Trigger" and "Why Drata and Vanta Won't Respond Downmarket-Auditor-Side" subsections should be reweighted to reflect that Drata/Vanta are not the primary comparison set.

---

## 3. SOC Type 2 Period Coverage Range

**Feedback:** Period coverage for a SOC Type 2 engagement can range from **3 to 12 months**, not "6 or 12 months minimum" as the spec currently states. The lower bound should be 3 months.

**Affected spec sections:** §4 Regulatory Compliance Requirements (the Engagement Type × Standard × Platform Requirements table currently reads "Type II: 6 or 12 months minimum"). User-journey references to SOC 2 Type II period windows should be reviewed for the same correction.

---

## 4. SOC Type 1 Is a Point-in-Time Engagement

**Feedback:** A SOC Type 1 engagement is a **point-in-time** examination expressed as a single "as of" date (for example, "as of 12/31/26"), not a continuous period. The platform's engagement-creation UI and reporting flows must capture a Type 1 engagement as a single date, not a date range.

**Affected spec sections:** §4 Regulatory Compliance Requirements; Journey 4 (Create New Engagement) in the user journeys. The engagement-period UI should branch on report type so that Type 1 engagements take a single date and Type 2 engagements take a date range.

---

## 5. Multi-Level Sign-Off Hierarchy

**Feedback:** Workpaper and test-procedure sign-off should support **four reviewer levels**, not the simpler chain currently described in the user journeys. The four levels are:

1. **Tester** — performs the test procedure and prepares the workpaper.
2. **Detailed Reviewer** — reviews the work at the workpaper / detailed level.
3. **General Reviewer** — reviews at the engagement / section level.
4. **Final Reviewer** — provides the final partner-level approval.

Engagement Quality Review (EQR) remains a separate independent review track that runs alongside this hierarchy.

**Affected spec sections:** Journey 8 (Prepare Workpaper) and Journey 9 (Manager Review and Workpaper Sign-Off) currently model a preparer → manager → partner sign-off chain. The data model and the workpaper UI need to expand to the four-level hierarchy above. The Audit Core data model (`Workpaper`, `WorkpaperVersion`, `ReviewNote`, sign-off audit-log entries) will need to encode the four reviewer roles and enforce ordering between them.

---

## 6. Generate ISO Certificates and PCI ROCs from Templates

> **Spec passage being commented on (§4, Explicit Scope Boundary):**
> "Axiom does **not** issue ISO certificates (that is a CB function under ISO 17021-1), **not** sign PCI ROCs (that is a QSA function), and **not** issue attestation opinions on behalf of any firm."

**Feedback:** Even though Axiom is not the legal signing authority for ISO certificates or PCI ROCs, the platform **could generate these artifacts from templates** — that is, produce the draft deliverable the same way it already produces SOC 2 reports. The legal sign-off authority would still sit with the accredited CB or QSA, but Axiom would offer template-generated drafts as a deliverable type.

**Affected spec sections:** §4 Explicit Scope Boundary (the wording needs to distinguish between "Axiom does not legally sign these documents" and "Axiom can produce template drafts of these documents") and the Reporting & Archival flow in the user journeys, which would need additional report-type entries for ISO certificate drafts and PCI ROC drafts. Ties into Item 1 (CBs and QSAs in ICP scope).

---

## 7. Implementation Plan — Approved

**Feedback:** The implementation plan covering the first six months and the post-MVP roadmap **looks excellent** and is approved as currently written. No changes requested.

**Affected artifacts:** `docs/superpowers/specs/implementation-plan-design.md` and the per-phase plans under `docs/superpowers/plans/`.

---

## 8. Positioning — Strong Contender in Audit Management and Evidence Collection

**Feedback:** With the integration roadmap as currently scoped, Axiom would be a **strong contender in the audit management and evidence collection categories**. The product owner views the integration story as the lever that makes Axiom competitive in these two adjacent segments.

**Affected spec sections:** §8 Integration Roadmap (which currently emphasizes evidence-bridge integrations from Drata Audit Hub, Vanta export, Sprinto / Hyperproof export, and direct OAuth connectors). The positioning narrative in §1 and §2 should be updated so that "audit management + evidence collection" is framed as a primary category position rather than as an emergent capability of the integration set.
