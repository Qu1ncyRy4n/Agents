---
status: needs-choice
candidate_outputs: [agent-guide]
tags: [process/handoff, process/decision, validation/evidence]
tldr: At handoff, map settled decisions to implementation evidence, validation, affected paths, and approved exceptions.
sources:
  - docs/other_repo_agents/ciwg_decomk-conf-cswg_refs_heads_main_AGENTS.md#required-final-handoff-artifacts
requires_all: []
mutually_exclusive: []
see_also: []
---
# Handoff Decisions with Implementation Evidence

For governed work, report what was decided and the evidence showing that the
result follows those decisions. Include affected runtime paths, checks run, and
exceptions rather than asking the reviewer to reconstruct compliance from a
large diff.

## Preserve the Strict Evidence Requirements

Verbatim source from
[`Required final handoff artifacts`](../../../../docs/other_repo_agents/ciwg_decomk-conf-cswg_refs_heads_main_AGENTS.md#required-final-handoff-artifacts):

> - `Decision Compliance: PASS/FAIL`
> - Decision Matrix mapping each locked decision ID to implementation evidence.
> - Inline diff annotations in the form `path:line -> decision_id -> rationale`.
> - Runtime Path Touch Matrix listing each approved runtime path/pattern, action used, and where it is implemented/validated.
> - `Exceptions:` listing only user-approved deviations.
> - Every non-trivial behavior change must include intent provenance per existing DI requirements.

## Review Concerns

- Decision-to-evidence mapping is useful for consequential work.
- Mandatory inline annotation for every changed line or path may produce more
  process material than review value for ordinary tasks.
- Runtime-path reporting may belong in high-consequence or data-mutation
  overlays rather than the default handoff.
- PASS/FAIL language is meaningful only when the applicable decisions and
  evidence requirements are explicit.

This should be compared with ordinary clear-handoff guidance before promotion.
