# Mogent Handoff

Updated: 2026-08-18

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
- Reusable core packages now live at public top-level import paths. `workspace`
  is the supported orchestration facade; the CLI and navigator import the same
  public manifest, library, render, source, state, and workspace packages.
- Local source roots and descendants reject symlinks consistently with pinned
  sources. Errors name the link/target and recommend ordinary content or a
  separately declared source.
- Builds use only explicit manifest template variables. `init` discovers and
  materializes `repo_name` and `repo_url`, with repeatable `--var` overrides.
- Atomic replacement is centralized in `renderfs`; cleanup errors are joined
  rather than discarded. CLI dispatch, usage, command completion, and typo
  suggestions share one command registry. XDG config precedence now works
  consistently on macOS and other platforms.
- Local `main` contains the focused public-API and cleanup commit series on top
  of `origin/main` at `bec8996`; it has not been pushed by this workflow. The
  URL-source subdirectory and selectable-directory work was already merged.
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
- Source directories are first-class selectable subtrees. Inventory and
  coverage share one aligned tree presenter with explicit directory/heading
  markers; `coverage` is the compatibility preset for
  `source list --coverage --tree`.
- User-local XDG YAML supports `display.chars: ascii|unicode` and
  `display.align`, terminal-aware fitting, and explicit widths; CLI flags
  override it and ASCII remains the default.
- `add --preview=tree` expands directory selections, distinguishes inherited
  source structure, and reports exact or related existing selections. Additions
  support `--first`, `--last`, `--before`, and `--after` placement.
- The post-merge duplicate top-level Nix source has been removed. DI-pesun's
  `personal:lang/nix` remains the sole canonical Nix module and is independently
  selectable.
- DI-vurap is implemented: compact scalar sources remain compatible, explicit
  sources accept `location` plus a locked `subdir`, and
  hashing/resolution/update review operate from that root. Unit/offline
  integration checks and real-project HTTPS dogfooding pass.
- The apparent GCC requirement is traced to TUI dependencies reaching
  `os/user`, which activates standard-library CGO. Mogent contains no C code;
  `tools/install` now forces the supported pure-Go build and was verified with
  a temporary `GOBIN`.
- Multiple outputs remain unimplemented; `docs/MULTIPLE-OUTPUTS-PLAN.md` records
  the proposed schema, transactions, state migration, CLI work, and decisions
  required before implementation.
- `docs/AUTHORING-PLAN.md` and `docs/ORIGIN-RECONCILIATION-PLAN.md` now provide
  implementation sequences and decision gates for the next two phases.
- `docs/PUBLIC-API-AND-CLEANUP-PLAN.md` records the user direction to expose the
  same supported Go API used by the CLI, the decisions required before moving
  packages out of `internal/`, a focused cleanup list, local-source symlink and
  template-determinism concerns, status/drift direction, and future comparative
  research into agent skill ecosystems.
- `docs/ACTIVE-DOGFOOD.md` is the owner-maintained registry for real working
  repositories. Keep `testing-ground/` only as disposable smoke coverage unless
  real dogfooding shows it still earns a larger role.
- Quincy has completed a first note-taking pass across the current library,
  with some language/domain modules intentionally deferred, and is keeping the
  detailed questions in the untracked user-authored
  `docs/notes_on_lib_mods.md`. Preserve that file and do not stage it without
  explicit direction. `docs/LIBRARY-REVIEW-PLAN.md` classifies the resulting
  clarity, taxonomy, relationship-metadata, ownership, and source-audit work.
  `docs/LIBRARY-COVERAGE-AUDIT.md` records the first sanitized comparison with
  the now-authorized protected original guides.
- The first library clarity slice replaces extraction-facing `Corpus Variants`
  labels with meaningful visible headings while explicit IDs preserve their
  prior source paths. Testing, error, Git, coordination, documentation, request
  conflict, and usage-evidence guidance now includes the compact context and
  examples requested in the owner notes.
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

Immediate closure: continue dogfooding the directory/unified-tree work from
`docs/DOGFOOD-SESSION-3.md` and record remaining friction.

After that, the user-prioritized development order is:

1. Continue Quincy's bounded review of the intentionally deferred language and
   domain modules. In parallel, follow `docs/LIBRARY-REVIEW-PLAN.md`: preserve
   the current library as a baseline, compare the authorized protected source
   guides locally, then make focused clarity changes before taxonomy or public
   metadata changes.
2. Continue public-API dogfooding and compatibility review. The large CLI
   handler was split by command family; continue focused cleanup without
   duplicating behavior.
3. Implement the authoring contract: `source add`, move/reorder, guided
   module creation, and reference-aware suggestions.
4. Make `status` canonical for read-only workspace state, retain `drift` as a
   shorthand for its direct-edit view, and move mutations toward directional
   origin/reconcile/inherit/propose workflows with three-way provenance checks.
5. Resolve and implement the multiple-output plan for `CLAUDE.md`, `GEMINI.md`,
   `.codex/AGENTS.md`, and similar explicit targets.
6. Add XDG YAML presentation configuration, alignment/color/TLDR/hint controls,
   and Nix-installed shell completion.
7. Research native skill management across agent ecosystems before designing
   later `.agents/skills/` composition or projection support.
8. Later, design `mogent handoff` to generate a compact agent-readable summary
   from explicit Git, manifest, output, source, coverage, TODO, and design state.

## Privacy Boundary

Do not inspect, commit, expose, quote, summarize, or reconstruct anything under
`docs/other_repo_agents/` unless the user explicitly authorizes that exact work.
Do not infer its contents from filenames, ignored state, history, or adjacent
documents. A future `mogent handoff` command must exclude protected/private
sources by default.
