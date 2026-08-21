---
status: needs-choice
candidate_outputs: [agent-guide]
tags: [process/decision, workflow/stop-rule, governance/strict]
tldr: Resolve and record durable choices before implementation, while keeping analysis distinct from decision authority.
sources:
  - docs/other_repo_agents/ciwg_decomk-conf-cswg_refs_heads_main_AGENTS.md#decision-first-specification-and-compliance-protocol-required
requires_all: []
mutually_exclusive: []
see_also: []
---
# Lock Durable Decisions Before Implementation

Separate investigation from authority. An agent may gather evidence, model
scenarios, and narrow alternatives before a decision. Once implementation
would commit the project to a durable choice, resolve and record that choice
before editing behavior around an unapproved assumption.

## Preserve the Original Strict Rule

Verbatim source from
[`Decision-First Specification and Compliance Protocol`](../../../../docs/other_repo_agents/ciwg_decomk-conf-cswg_refs_heads_main_AGENTS.md#decision-first-specification-and-compliance-protocol-required):

> Decision-first means decisions must be locked before coding; it does not forbid pre-decision analysis such as required thought experiments.
> The agent must collect and lock user decisions before making any code edits for a task.
> Locked decisions must be recorded as Decision Intent Log entries in the relevant `TODO/*.md` file(s) with clear intent and rationale.
> The agent must ask decision questions up front in a single intake round whenever possible.
> Required decision categories are architecture, design/behavior, implementation approach, function naming, variable naming, and file/path decisions.
> The agent must ask these as multiple-choice questions whenever practical.

## Stop When a Required Decision Is Missing

Verbatim source from
[`Decision Lock and Stop Rule`](../../../../docs/other_repo_agents/ciwg_decomk-conf-cswg_refs_heads_main_AGENTS.md#decision-lock-and-stop-rule):

> The agent must produce a Decision Lock summary with decision IDs before code edits begin.
> The agent must not proceed if any required decision is missing, ambiguous, or conflicting.
> The agent must stop and ask immediately if a new decision need appears during implementation.
> The agent must not assume defaults for locked categories unless the user explicitly approves defaults.

## Review Concerns

- Requiring decisions for every function name, variable name, and touched path
  may overwhelm ordinary work and conflict with the risk-based workflow.
- Multiple-choice questions are useful when real alternatives exist but can
  manufacture false choices when one answer follows directly from established
  repository conventions.
- “Before any code edits” may need to distinguish implementation from safe
  exploratory tests, diagnostic instrumentation, or disposable prototypes.
- This module may be inseparable from the complete strict governance bundle.
