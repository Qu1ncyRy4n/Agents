---
status: needs-choice
candidate_outputs: [agent-guide]
tags: [workflow/risk, process/escalation, process/decision]
tldr: Handle routine reversible work directly, but pause for review before durable, risky, or ownership-changing decisions.
sources:
  - libraries/cdint/process/lightweight-escalation.md
  - docs/other_repo_agents/ciwg_decomk-conf-cswg_refs_heads_main_AGENTS.md#decision-first-specification-and-compliance-protocol-required
requires_all: []
mutually_exclusive: []
see_also: []
---
# Use Risk to Choose Routine Work or Decision Review

Use an ordinary focused workflow for routine, reversible changes with clear
requirements. Escalate before making a durable choice when the work changes an
architecture, public interface, persistent state, trust boundary, irreversible
operation, public specification, or another repository's ownership boundary.

<!-- review-note: This is the major policy fork in the current sources. Some
guides require locked decisions before every code edit; the v1 cdint module uses
risk-based escalation. Preserve both until the owner decides whether strict
decision-first remains a selectable alternative or a compatible deeper layer. -->

## Use the Focused Loop for Routine Work

Verbatim v1 library source from
[`Lightweight Escalation`](../../../libraries/cdint/process/lightweight-escalation.md):

> Use the focused change loop for routine fixes, small documentation edits, and
> local improvements. Pause and ask before making a durable choice when the work
> changes architecture, a public or wire format, persistent data, a security
> boundary, an irreversible operation, a public specification, or another
> repository's ownership boundary.
>
> Recommend a wider check when the scope, result, or risk looks uncertain.

This is an existing library adaptation, not verbatim source-guide wording. Its
source-guide derivation still needs heading-level traceability.

## Preserve the Strict Decision-First Option

Verbatim source from
[`Decision-First Specification and Compliance Protocol`](../../../docs/other_repo_agents/ciwg_decomk-conf-cswg_refs_heads_main_AGENTS.md#decision-first-specification-and-compliance-protocol-required):

> Decision-first means decisions must be locked before coding; it does not forbid pre-decision analysis such as required thought experiments.
> The agent must collect and lock user decisions before making any code edits for a task.
> Locked decisions must be recorded as Decision Intent Log entries in the relevant `TODO/*.md` file(s) with clear intent and rationale.

The strict source continues with mandatory decisions for architecture,
behavior, implementation approach, function names, variable names, and every
file or runtime path. Those requirements are preserved in the source guide and
must be extracted before deciding whether this option remains one complete
bundle.

## Proposed Decision Boundary

Proposed example:

```text
Routine: Correct a typo, repair a clearly failing focused test, or implement a
fully specified local change with no durable alternatives.

Decision review: Choose a persisted data format, rename a public command,
change an authorization boundary, perform an irreversible migration, or decide
which repository owns shared behavior.
```

## Review Questions

1. Is risk-based escalation the canonical default?
2. Is strict decision-first a selectable alternative, or a deeper bundle used
   only for particular projects and tasks?
3. Which strict requirements are valuable independently, and which become
   counterproductive when separated from the complete governance system?
