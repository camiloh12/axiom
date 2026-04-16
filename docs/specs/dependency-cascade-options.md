# Dependency Cascade Automation Options

**Date:** April 15, 2026
**Status:** Reference — Option 1 implemented, Options 2-3 deferred

This document records three approaches for ensuring that changes to upstream specification documents automatically cascade to downstream artifacts. Option 1 (CLAUDE.md manifest) is implemented. Options 2 and 3 are documented here for future reference.

---

## Background

Axiom's specification documents have cascading dependencies. A change to the AI Architecture spec, for example, requires updates to the Domain Model, Backend Architecture, OpenAPI specs, Infrastructure Design, the infrastructure diagram, and the Product Spec summary. Missing any downstream artifact creates inconsistency.

The full dependency graph has 21 edges across 13 artifact types. Key characteristics:
- The Domain Model and AI Architecture have a **bidirectional** dependency
- The Product Spec is a terminal **summary hub** — updated last
- The infrastructure diagram image is a **generated artifact** (script → PNG)
- OpenAPI `common.yaml` cascades to all 6 service spec files
- Mockups depend on both User Journeys (content) and the Design System (visual)

---

## Option 1: CLAUDE.md Dependency Manifest (Implemented)

**Location:** `CLAUDE.md` → "Document Dependencies" section

Claude Code reads CLAUDE.md at session start. The dependency manifest gives Claude the knowledge to identify which downstream artifacts need updating when an upstream document changes.

**What it contains:**
- Artifact map (ID, name, path)
- Cascade rules per upstream document (which downstream artifacts to update, with specific section references)
- Bidirectional dependency callout (Domain Model ↔ AI Architecture)
- Terminal node identification (Product Spec, diagram PNG)
- Workflow steps (upstream → cascade → regenerate → summarize)

**Strengths:**
- Zero tooling, works immediately
- Claude reads it every session without additional configuration
- Human-readable, easy to update

**Weaknesses:**
- Advisory, not enforced — relies on Claude (or the human) checking the rules
- No automated detection of which files changed
- Could become stale if artifacts are added/removed without updating CLAUDE.md

---

## Option 2: Pre-Commit Hook + Changed-File Analysis

Register a Claude Code hook that fires on `UserPromptSubmit`, inspects `git diff` to detect which spec files have uncommitted changes, and outputs a reminder of downstream artifacts that may need updating.

### Implementation

**Hook registration** (`.claude/settings.json`):
```json
{
  "hooks": {
    "UserPromptSubmit": [
      {
        "type": "command",
        "command": "bash docs/scripts/check-spec-cascade.sh"
      }
    ]
  }
}
```

**Script** (`docs/scripts/check-spec-cascade.sh`):
```bash
#!/usr/bin/env bash
# Detect changed spec files and warn about downstream cascade obligations

changed=$(git diff --name-only HEAD 2>/dev/null)
if [ -z "$changed" ]; then
  exit 0
fi

warnings=""

if echo "$changed" | grep -q "docs/user-journeys/all-journeys.md"; then
  warnings+="WARNING: User Journeys changed → check: domain-and-data-model-design.md, ai-architecture-design.md, axiom-spec-design.md (§5, §6, §11-12), mockups/\n"
fi

if echo "$changed" | grep -q "docs/specs/domain-and-data-model-design.md"; then
  warnings+="WARNING: Domain Model changed → check: ai-architecture-design.md (entity refs), backend-architecture-design.md (§2-3), packages/openapi/*.yaml (schemas/enums), axiom-spec-design.md (§5)\n"
fi

if echo "$changed" | grep -q "docs/specs/ai-architecture-design.md"; then
  warnings+="WARNING: AI Architecture changed → check: domain-and-data-model-design.md (ai_content_metadata, enums), backend-architecture-design.md (River workers, Bedrock, pgvector), packages/openapi/*.yaml (AI endpoints), infrastructure-design.md (IAM, observability), axiom-spec-design.md (§6)\n"
fi

if echo "$changed" | grep -q "docs/specs/backend-architecture-design.md"; then
  warnings+="WARNING: Backend Architecture changed → check: packages/openapi/*.yaml (endpoints), infrastructure-design.md (ECS, DB, IAM), axiom-spec-design.md (§7)\n"
fi

if echo "$changed" | grep -q "docs/specs/infrastructure-design.md"; then
  warnings+="WARNING: Infrastructure Design changed → check: docs/diagrams/axiom_infrastructure.py, then regenerate .png\n"
fi

if echo "$changed" | grep -q "packages/openapi/common.yaml"; then
  warnings+="WARNING: OpenAPI common.yaml changed → check all 6 service .yaml files for $ref consistency\n"
fi

if echo "$changed" | grep -q ".impeccable.md"; then
  warnings+="WARNING: Design System changed → check mockups/ for visual consistency\n"
fi

if echo "$changed" | grep -q "docs/diagrams/axiom_infrastructure.py"; then
  warnings+="WARNING: Infrastructure diagram script changed → regenerate axiom_infrastructure.png\n"
fi

if [ -n "$warnings" ]; then
  echo "━━━ SPEC CASCADE CHECK ━━━"
  echo -e "$warnings"
  echo "See CLAUDE.md 'Document Dependencies' for full cascade rules."
  echo "━━━━━━━━━━━━━━━━━━━━━━━━━"
fi
```

**Strengths:**
- Automatic detection on every prompt — hard to miss
- Shows warnings contextually when relevant files are dirty
- Can be combined with Option 1 (CLAUDE.md provides the detail, hook provides the trigger)

**Weaknesses:**
- Fires every prompt (could be noisy if files stay dirty across a long session)
- Shell script maintenance — must be kept in sync with the dependency graph
- Only detects uncommitted changes (doesn't help with changes made in the current prompt before the next fires)

---

## Option 3: Structured Dependency Manifest + Validation Skill

Create a machine-readable YAML dependency manifest and a custom Claude Code skill (`/cascade`) that reads it, analyzes changed files, and produces a structured checklist.

### Implementation

**Manifest** (`docs/dependencies.yaml`):
```yaml
# Document dependency graph — machine-readable
# Used by the /cascade skill to detect incomplete propagation

artifacts:
  research:
    paths: ["docs/research/*.md"]
    downstream: [product-spec]
    sections_affected:
      product-spec: "§1-4, §8-9"

  user-journeys:
    paths: ["docs/user-journeys/all-journeys.md"]
    downstream: [domain-model, ai-architecture, product-spec, mockups]
    sections_affected:
      domain-model: "entity journey tags, §12 traceability matrix"
      ai-architecture: "§4 feature journey references"
      product-spec: "§11 journey summaries, §12 innovative flows"
      mockups: "directories and screens map 1:1 to journey stages"

  domain-model:
    paths: ["docs/specs/domain-and-data-model-design.md"]
    downstream: [ai-architecture, backend-architecture, openapi-specs, product-spec]
    bidirectional: [ai-architecture]
    sections_affected:
      ai-architecture: "entity name references (AIDecision, WorkpaperVersion, EvidenceItem, etc.)"
      backend-architecture: "§2-3 service entity ownership, §3 DB topology"
      openapi-specs: "schemas, enums, constraints"
      product-spec: "§5 data model summary"

  ai-architecture:
    paths: ["docs/specs/ai-architecture-design.md"]
    downstream: [domain-model, backend-architecture, openapi-specs, infrastructure, product-spec]
    bidirectional: [domain-model]
    sections_affected:
      domain-model: "ai_content_metadata jsonb structure, ai_context_type enum values"
      backend-architecture: "River workers, Bedrock model assignments, pgvector scope, cross-service AIDecision pattern"
      openapi-specs: "AI endpoints, ai_content_metadata schema, enum values"
      infrastructure: "Bedrock IAM per-service model access, River DLQ alarms per service, AI observability dashboard"
      product-spec: "§6 AI summary"

  backend-architecture:
    paths: ["docs/specs/backend-architecture-design.md"]
    downstream: [openapi-specs, infrastructure, product-spec]
    sections_affected:
      openapi-specs: "endpoint distribution across service specs"
      infrastructure: "ECS task definitions, DB topology, IAM roles"
      product-spec: "§7 tech stack summary"

  infrastructure:
    paths: ["docs/specs/infrastructure-design.md"]
    downstream: [diagram-script]
    sections_affected:
      diagram-script: "visual rendering of infrastructure components and connections"

  diagram-script:
    paths: ["docs/diagrams/axiom_infrastructure.py"]
    downstream: [diagram-image]
    action: "cd docs/diagrams && python axiom_infrastructure.py"

  diagram-image:
    paths: ["docs/diagrams/axiom_infrastructure.png"]
    downstream: []
    generated: true

  openapi-common:
    paths: ["packages/openapi/common.yaml"]
    downstream: [openapi-specs]
    sections_affected:
      openapi-specs: "all 6 service specs $ref shared schemas, parameters, and responses"

  openapi-specs:
    paths:
      - "packages/openapi/identity-service.yaml"
      - "packages/openapi/audit-core.yaml"
      - "packages/openapi/trial-balance.yaml"
      - "packages/openapi/workpaper-service.yaml"
      - "packages/openapi/reporting-service.yaml"
      - "packages/openapi/doc-processing.yaml"
    downstream: []

  design-system:
    paths: [".impeccable.md"]
    downstream: [mockups]
    sections_affected:
      mockups: "typography, colors, spacing, layout patterns"

  mockups:
    paths: ["mockups/**/*.html"]
    downstream: []

  product-spec:
    paths: ["docs/specs/axiom-spec-design.md"]
    downstream: []
    terminal: true
    note: "Summary hub — aggregates from all other specs. Update last."

  claude-md:
    paths: ["CLAUDE.md"]
    downstream: []
    note: "Meta-file with dependency rules. Update when artifacts are added/removed."
```

**Skill** (`.claude/skills/cascade.md`):
```markdown
---
name: cascade
description: Analyze changed files and produce a cascade checklist of downstream artifacts that need updating
user_invocable: true
---

# /cascade — Dependency Cascade Checker

When invoked, perform these steps:

1. Read `docs/dependencies.yaml` to load the dependency graph.
2. Run `git diff --name-only HEAD` to identify changed files.
3. If the user passed a file path argument, use that as the changed file instead of git diff.
4. For each changed file, match it against `artifacts[*].paths` globs.
5. Walk the dependency graph recursively (follow `downstream` edges, including transitive dependencies).
6. For bidirectional edges, flag them but do not recurse infinitely — check each artifact at most once.
7. Output a structured checklist:

   ## Cascade Checklist
   
   **Changed:** [list of changed upstream artifacts]
   
   | Priority | Artifact | What to check | Status |
   |----------|----------|---------------|--------|
   | 1 | ... | specific sections/areas | [ ] |
   | 2 | ... | ... | [ ] |
   
   Priority ordering: direct dependents first, then transitive, then terminal nodes (product-spec) last.

8. If any artifact has an `action` field, include it as a required step.
9. If the user asks to execute the cascade, proceed through each artifact in priority order.
```

**Strengths:**
- Machine-readable manifest can power multiple consumers (skill, hook, CI check)
- The `/cascade` skill gives an on-demand audit of cascade completeness
- `sections_affected` metadata provides specific guidance, not just file names
- The YAML could later drive a CI job that fails the PR if downstream artifacts weren't touched
- Extensible — add new artifacts without changing code

**Weaknesses:**
- Most setup effort of the three options
- YAML must be kept in sync with reality (though the skill itself could validate this)
- Skill requires Claude Code skill authoring

### Future CI Integration

The YAML manifest could power a GitHub Actions check:
```yaml
# .github/workflows/spec-cascade.yml
- name: Check spec cascade completeness
  run: |
    python docs/scripts/validate-cascade.py \
      --manifest docs/dependencies.yaml \
      --changed $(git diff --name-only origin/master...HEAD)
```

This would flag PRs that change an upstream spec without touching its downstream artifacts — not as a hard block (some changes genuinely don't cascade), but as a reviewer hint.

---

## Recommended Path

1. **Now:** Option 1 (CLAUDE.md manifest) — implemented
2. **When starting implementation:** Option 3 (YAML manifest + `/cascade` skill) — the structured format supports automated validation as the codebase grows
3. **Optional:** Option 2 (hook) can be added alongside Option 3 as an always-on reminder layer

---

*End of Dependency Cascade Options*
