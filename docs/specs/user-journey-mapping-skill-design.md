# User Journey Mapping Skill — Design Spec

## Purpose

A general-purpose skill that guides Claude through creating task-level user journey maps for any product or feature. Produces a structured markdown document describing stages, touchpoints, emotions, pain points, and opportunities.

## Skill Identity

- **Name:** `user-journey-mapping`
- **Type:** Technique (concrete method with steps to follow)
- **Location:** `~/.claude/skills/user-journey-mapping/SKILL.md`

## Trigger Conditions

### Standalone

Invoked explicitly when someone wants to map a user journey:

- "Map the user journey for X"
- "What's the user experience for Y"
- "Walk me through the flow for Z"

### From Brainstorming

The brainstorming skill suggests invoking it when the feature being designed involves user-facing flows. Signals that should trigger the suggestion:

- The feature spans multiple screens or views
- Multiple user roles interact with the feature
- There's a sequential workflow (create → configure → submit → review)
- The user mentions "flow", "experience", "steps", or "process"

Suggested prompt from brainstorming:

> "This feature involves a multi-step user flow. Want me to map out the user journey before we finalize the design? I can invoke the user-journey-mapping skill."

The resulting journey map becomes an input to the design spec — it informs what screens/components are needed and what the interaction model looks like.

## Input Modes

The skill supports two input modes, determined at invocation:

### Interactive Mode (default)

Used when the user invokes the skill directly without pre-gathered context. The skill asks questions one at a time to build the journey map collaboratively.

### Research-Fed Mode

Used when pre-gathered research (personas, flows, domain context) is provided as input. The skill skips the interview and goes straight to stage decomposition and mapping, generating a complete draft journey map from the research.

**How it works:**

1. The user (or another skill/agent) gathers research on personas, goals, and flows
2. The research is passed to the skill as context when invoking it
3. The skill assesses whether the research provides sufficient detail for each dimension (actions, touchpoints, emotions, pain points, opportunities)
4. If sufficient, it generates the full journey map and presents the draft for review
5. If ambiguous, it applies this rule:
   - **Minor ambiguity** (single stage, doesn't change flow structure): make best judgment, flag the assumption inline with `[Assumption: ...]` for user review
   - **Significant ambiguity** (affects multiple stages, changes flow structure, or involves role/handoff decisions): stop and ask the user before proceeding

**Detecting the mode:** If the user's message or context includes research artifacts, persona descriptions, or flow documentation, default to research-fed mode. If the user just says "map the journey for X" with no supporting context, use interactive mode.

## Process Flow

### Interactive Mode

#### 1. Identify the Persona and Goal

Ask the user:

- **Who is the user?** Role, experience level, context.
- **What are they trying to accomplish?** The end goal.
- **What triggers them to start?** The entry point or initiating event.

One question at a time. Prefer multiple choice when the project context makes options obvious.

#### 2. Decompose Goal into Sub-Goals

Break the end goal into 3–7 sub-goals. Each sub-goal becomes a journey stage. Name each stage with a verb phrase (e.g., "Set up engagement", "Collect evidence", "Submit for review").

Present the proposed stages to the user and validate before proceeding.

#### 3. Map Each Stage (Task-Level Detail)

For every stage, ask the user about and document:

| Dimension | What to capture |
|-----------|----------------|
| **User Actions** | Specific tasks and interactions the user performs |
| **Touchpoints** | Screens, notifications, emails, or system elements involved |
| **Thoughts & Emotions** | What the user is thinking/feeling (confident, confused, frustrated, relieved) |
| **Pain Points** | Friction, confusion, or blockers at this stage |
| **Opportunities** | Ways to reduce friction or delight the user |

Work through stages sequentially. For each stage, ask the user about actions and touchpoints (what they concretely do), then pain points and opportunities (what frustrates them and what could be better). Infer thoughts & emotions from context and confirm with the user.

#### 4. Detect Handoffs

If a stage involves another role (e.g., "Manager reviews the submission"), flag it and ask whether to expand into multi-persona mapping.

When expanding, document the handoff:

- Who hands off and who receives
- What information transfers
- What each party sees at the handoff point
- The receiving persona's experience through their part of the flow

#### 5. Synthesize

After all stages are mapped:

- **Emotional arc** — narrative of how user confidence/satisfaction changes across stages. Where are the highs? Where are the lows?
- **Cross-cutting pain points** — pain points that appear in multiple stages or compound across the journey.
- **Prioritized opportunities** — ranked by impact and effort. Which improvements would most transform the experience?

### Research-Fed Mode

#### 1. Assess Research Completeness

Review the provided research and check for:

- Clear persona definition (role, context, experience level)
- Defined goal and trigger
- Enough detail to identify discrete stages
- Sufficient information about what the user does at each stage

If the research covers these, proceed. If it's missing fundamentals (no clear persona or goal), fall back to interactive mode for those gaps only.

#### 2. Generate Stage Decomposition

Derive 3–7 stages from the research. Present the proposed stages to the user for validation before generating the full map.

#### 3. Generate Full Journey Map

For each stage, populate all five dimensions (actions, touchpoints, thoughts & emotions, pain points, opportunities) from the research. Flag assumptions inline with `[Assumption: ...]` where the research was ambiguous but the ambiguity is minor.

Stop and ask the user when encountering significant ambiguity.

#### 4. Detect and Expand Handoffs

If the research describes role transitions, automatically expand into multi-persona mapping. If it's unclear whether a handoff exists, ask.

#### 5. Synthesize and Present Draft

Generate the full journey summary (emotional arc, cross-cutting pain points, prioritized opportunities) and present the complete draft for user review.

## Output Document Structure

```markdown
# User Journey: [Goal Description]

## Overview
- **Persona:** [role, context]
- **Goal:** [what they're trying to accomplish]
- **Trigger:** [what initiates the journey]
- **Stages:** [numbered list of stage names]

## Stage 1: [Verb Phrase]
### Sub-goal
What the user is trying to achieve in this stage.

### User Actions
- Step-by-step tasks the user performs

### Touchpoints
- Screens, components, notifications involved

### Thoughts & Emotions
- What the user is thinking/feeling at this point

### Pain Points
- Friction, confusion, or blockers

### Opportunities
- Ways to improve this stage

---
(repeat for each stage)

## Handoffs
(only present if multi-persona was triggered)

| From | To | Information Transferred | Trigger |
|------|-----|----------------------|---------|
| Staff | Manager | Completed workpaper | Submit for review |

### [Receiving Persona]'s Experience
Mapping of what the receiving persona sees and does at the handoff point.

## Journey Summary

### Emotional Arc
Narrative of how the user's confidence/satisfaction changes across stages.

### Cross-Cutting Pain Points
Pain points that appear in multiple stages or affect the overall experience.

### Prioritized Opportunities
Ranked list of improvements, considering impact and effort.
```

**Output location:** `docs/user-journeys/[persona]-[goal-slug].md`

## Scope Boundaries

**In scope:**
- Single flow per invocation (one persona pursuing one goal)
- Task-level detail within each stage
- Multi-persona expansion when handoffs are detected
- Structured markdown output

**Out of scope:**
- Visual/diagrammatic journey maps
- Competitive analysis or benchmarking
- Business metrics or KPI definition
- Technical implementation details
- Service blueprints (backend process mapping)

## Key Principles

- **One question at a time** — don't overwhelm with multiple questions
- **Goal-driven stages** — stages are derived from user goals, not imposed from a template
- **Task-level granularity** — each stage includes specific actions and interactions
- **Smart handoff detection** — single persona by default, expand only when the flow involves role transitions
- **Synthesize, don't just list** — the summary section should surface insights, not repeat what's already in the stages
