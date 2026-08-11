# Mogent Handoff

Updated: 2026-08-11

This is the compact resume point for another agent. Read `AGENTS.md` and
`docs/DESIGN.md`, then inspect the current Git status before changing anything.
Repository docs and Git state are authoritative if this file or a chat summary
is stale.

## Durable Context Map

- `AGENTS.md`: workflow, constraints, and repository-specific working rules.
- `docs/DESIGN.md` and decision records: product design and rationale.
- `TODO/`: planned, completed, and unresolved work.
- `docs/HANDOFF.md`: current state, unfinished work, and next steps only.

Keep durable project memory in these repository artifacts rather than relying on
ChatGPT, Codex, or another product to share conversation history.

## Current State

- Architecture is core/CLI-first. The TUI remains a client of reusable workspace
  operations, not the only home of product behavior.
- The prior dogfood baseline is committed as `8bf094b`. Current feature work is
  on `codex/localization-and-pinning`.
- `mogent status`, `mogent source show`, source listing/discovery, and
  `mogent coverage` exist. Current uncommitted work further improves status,
  discovery, coverage, completion, diagnostics, and dogfood material.
- Core-first copy-on-write localization, conservative single-section drift
  import, immutable URL pinning, and guided starter templates are implemented on
  the feature branch.
- URL manifests use committed `mogent.lock.yaml` records and verified ignored
  `.mogent/sources/` checkouts. Ordinary commands never fetch.
- Multiple outputs remain unimplemented; `docs/MULTIPLE-OUTPUTS-PLAN.md` records
  the proposed schema, transactions, state migration, CLI work, and decisions
  required before implementation.
- DI-pesun resolves the Nix portion of DR-garom: Nix now lives at
  `personal:lang/nix`, remains independently selectable, and retains explicit
  machine-safety guidance.
- `docs/DOGFOOD.md` stages the implemented features for repeatable exercises and
  issue reports. Seed library files now have concise file-level TLDR metadata
  for compact source selection.

## Capability Sequence

Continue to organize the agent-facing workspace surface in this order:

1. `status` - inspect manifest, output, and workspace state.
2. show/inspect - inspect one source or rendered relationship.
3. sources/discovery - explore what is available before editing.
4. coverage - compare selected and available material.
5. localization - copy-on-write editing with explicit provenance.

All five have core/CLI implementations on the feature branch. The TUI can adopt
the same workspace operations later.

## Next Steps

1. Complete the branch-wide validation and requirement audit.
2. Run the HTTPS pin smoke test when network-command approval is available; the
   approval service rejected the attempted disposable fetch before execution.
3. Exercise `docs/DOGFOOD.md` manually when the user returns, especially URL
   update review and real-repository starter selection.
4. Resolve the questions in `docs/MULTIPLE-OUTPUTS-PLAN.md` before implementation.
5. Later, design `mogent handoff` to generate a compact agent-readable summary
   from explicit Git, manifest, output, source, coverage, TODO, and design state.

## Privacy Boundary

Do not inspect, commit, expose, quote, summarize, or reconstruct anything under
`docs/other_repo_agents/` unless the user explicitly authorizes that exact work.
Do not infer its contents from filenames, ignored state, history, or adjacent
documents. A future `mogent handoff` command must exclude protected/private
sources by default.
