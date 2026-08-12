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
- `main` is synchronized with `origin/main` at merge commit `7efe946`. The
  localization/pinning branch, real-project dogfood feedback, and the remote
  atomic-library split are all merged. The current feature work and retained
  follow-up documentation are on `codex/url-source-subdir`.
- `mogent status`, `mogent source show`, source listing/discovery, and
  `mogent coverage` exist. Real-project dogfooding is active in
  `/home/qix/dev/omnicortex/todo_app_project`.
- Core-first copy-on-write localization, conservative single-section drift
  import, immutable URL pinning, guided starter templates, file-level TLDRs, and
  coverage TLDR presentation are implemented on `main`.
- URL manifests use committed `mogent.lock.yaml` records and verified ignored
  `.mogent/sources/` checkouts. Ordinary commands never fetch.
- A live HTTPS pin against this GitHub repository succeeded and verified an
  immutable commit/hash. Consuming the repository root as one source then found
  duplicate documentation heading paths, confirming that pinned source
  subdirectories are required before `libraries/cdint` or `libraries/personal`
  can be clean remote aliases.
- Source libraries are now split into atomic metadata-bearing Markdown files.
- The remote atomic split also reintroduced `libraries/nix/nix/*`, conflicting
  with active DI-pesun's `personal:lang/nix` boundary. Both currently exist;
  reconcile the intended canonical location before publishing library URLs.
- Current work is on `codex/url-source-subdir`. DI-vurap is implemented: compact
  scalar sources remain compatible, explicit sources accept `location` plus a
  locked `subdir`, and hashing/resolution/update review operate from that root.
  Unit and offline integration checks pass; the live HTTPS combined smoke test
  was blocked before execution by the approval service.
- The apparent GCC requirement is traced to TUI dependencies reaching
  `os/user`, which activates standard-library CGO. Mogent contains no C code;
  `tools/install` now forces the supported pure-Go build and was verified with
  a temporary `GOBIN`.
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

All five have core/CLI implementations on `main`. The TUI can adopt
the same workspace operations later.

## Next Steps

1. Reconcile the Nix source-boundary conflict introduced by the remote atomic
   split; do not expose two canonical Nix sources accidentally.
2. Repeat live remote pin/build/browse dogfooding against `libraries/cdint` or
   `libraries/personal` when network-command approval works, then commit and
   merge the subdirectory branch.
3. Design the ASCII-first shared list/coverage presenter and XDG YAML personal
   presentation config, including Nix-installed shell completion.
4. Run a thought experiment replacing the directionless `drift` surface with
   origin/reconcile/inherit/propose workflows and three-way provenance checks.
5. Add the guided module-creation workflow after its destination and metadata
   contract are recorded.
6. Resolve the questions in `docs/MULTIPLE-OUTPUTS-PLAN.md` before implementation.
7. Later, design `mogent handoff` to generate a compact agent-readable summary
   from explicit Git, manifest, output, source, coverage, TODO, and design state.

## Privacy Boundary

Do not inspect, commit, expose, quote, summarize, or reconstruct anything under
`docs/other_repo_agents/` unless the user explicitly authorizes that exact work.
Do not infer its contents from filenames, ignored state, history, or adjacent
documents. A future `mogent handoff` command must exclude protected/private
sources by default.
