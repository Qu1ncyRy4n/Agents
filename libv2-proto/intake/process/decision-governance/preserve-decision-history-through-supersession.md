---
status: reliable
candidate_outputs: [agent-guide, specification]
tags: [process/decision, records/history, provenance/intent]
tldr: Preserve prior decisions as history and add an explicit superseding record when intent changes.
sources:
  - docs/other_repo_agents/ciwg_decomk-conf-cswg_refs_heads_main_AGENTS.md#comment-preservation-protocol-required
requires_all: []
mutually_exclusive: []
see_also: []
---
# Preserve Decision History Through Supersession

Treat durable decision records as historical evidence. When intent changes,
add a new decision that identifies the earlier record it supersedes instead of
rewriting the earlier decision to look as though the new intent was always in
force.

## Preserve the Source Record Shape

Verbatim source from
[`Comment Preservation Protocol`](../../../../docs/other_repo_agents/ciwg_decomk-conf-cswg_refs_heads_main_AGENTS.md#comment-preservation-protocol-required):

> Maintain a `## Decision Intent Log` at the top of relevant `TODO/*.md` files.
> Treat DI logs as append-only history. Do not rewrite or delete prior entries.
> When intent evolves, add a new DI entry and set `Supersedes: <old-di-id>`.
> DI entries must include:
>
> - `ID: DI-<handle>`
> - `Date: YYYY-MM-DD HH:MM:SS`
> - `Status: active|superseded`
> - `Decision:`
> - `Intent:`
> - `Constraints:`
> - `Affects:`
> - `Supersedes:` (optional)

## Review Concerns

- Historical preservation is broadly useful; storing the log inside TODO files
  is an organization-specific layout choice.
- The exact ID format and shared proquint namespace belong to a coordination-ID
  module or specification.
- “Append-only” should mean historical decisions are not silently rewritten;
  it may still permit transparent typo or metadata corrections under an
  explicit record-maintenance policy.
- This content may ultimately belong partly in an agent guide and partly in a
  generated decision-record specification.
