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
- The last validated handoff described the `workspace.Session` extraction as
  uncommitted. The current working tree still contains uncommitted workspace,
  CLI, documentation, library, and dogfood changes, so inspect and preserve the
  entire existing diff before editing or committing.
- `mogent status`, `mogent source show`, source listing/discovery, and
  `mogent coverage` exist. Current uncommitted work further improves status,
  discovery, coverage, completion, diagnostics, and dogfood material.
- Copy-on-write localization, drift import/merge, remote URL sources, and pinning
  remain future work.
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

The first four have implementations. Stabilize and validate them through the
shared workspace core, then proceed to localization rather than placing durable
behavior only in the TUI.

## Next Steps

1. Finish reviewing and validating the full uncommitted dogfood diff; do not
   overwrite or silently absorb unrelated user work.
2. Exercise `docs/DOGFOOD.md` against a disposable real repository and record
   discovery, TLDR, coverage, preview, completion, and save/rebuild friction.
3. Add coverage TLDR summaries and revisit heading-level metadata only after
   file-level summaries have been tested in real module selection.
4. Design and implement copy-on-write localization through reusable workspace
   operations with a useful noninteractive command shape.
5. Later, design `mogent handoff` to generate a compact agent-readable summary
   from explicit Git, manifest, output, source, coverage, TODO, and design state.

## Privacy Boundary

Do not inspect, commit, expose, quote, summarize, or reconstruct anything under
`docs/other_repo_agents/` unless the user explicitly authorizes that exact work.
Do not infer its contents from filenames, ignored state, history, or adjacent
documents. A future `mogent handoff` command must exclude protected/private
sources by default.
