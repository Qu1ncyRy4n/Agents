# Mogent Public API And Cleanup Plan

Status: implementation direction approved by DI-zunap; detailed exported types
remain subject to compatibility review during each slice.

## Implementation Progress

Completed on 2026-08-18:

- Public `workspace`, `manifest`, `library`, `render`, `renderfs`,
  `sourcecache`, `sourcepath`, `starter`, and `state` packages; CLI and navigator
  adapters import these packages.
- An external-package workspace test that exercises load, status, source
  discovery, and dry-run addition through the same API used by the CLI.
- Default rejection of symlinked local source roots, files, and directories,
  with actionable errors and focused tests.
- Explicit, materialized `repo_name` and `repo_url` manifest variables; builds
  no longer inspect checkout directory names or Git remotes.
- One atomic file writer with joined cleanup errors, used by manifest, render,
  state, provenance, and source-lock paths.
- A single CLI command registry for dispatch, usage, completion candidates, and
  typo suggestions, followed by purpose-based handler files for init,
  workspace/status, source browsing, manifest addition, completion, and shared
  flag parsing.
- Explicit `XDG_CONFIG_HOME` precedence on every platform.

The complete `tools/check` suite passed after these slices.

## User Direction

- Move the reusable implementation out of Go's `internal/` boundary.
- Expose the same core API that the CLI uses. The CLI should be a thin client,
  not a second implementation or a privileged path through the product.
- Dogfood the public API as well as the executable.
- Keep skill management as later work. First study how each agent ecosystem
  defines, discovers, scopes, installs, and shares skills on its own terms.
- Treat the complete agent-module library review as an explicit owner task and
  keep reminding the user about it until it is checked off.

## Decisions Before Moving Packages

1. **Public promise.** Decide whether every moved package is supported API or
   whether a small supported facade sits over public-but-lower-level packages.
2. **Package map.** Choose durable names and responsibilities. A likely shape is
   `mogent` or `workspace` for orchestration, with focused `manifest`, `library`,
   `render`, and `source` packages only where direct use is genuinely valuable.
3. **CLI and TUI boundary.** Decide whether command parsing and terminal models
   are reusable public packages or remain executable adapters. Moving reusable
   behavior out does not require promising compatibility for presentation code.
4. **Operation shape.** Standardize load, inspect, preview, validate, and apply
   operations so CLI, TUI, tests, agents, and other Go programs call the same
   methods with the same defaults and safety checks.
5. **Data and error contracts.** Decide which result types, diagnostics, error
   categories, and path/reference types callers may rely on. Avoid making CLI
   strings the machine API.
6. **Effects and dependencies.** Decide how callers supply filesystem, Git,
   clock, terminal, and future network behavior. Add `context.Context` where
   cancellation matters, especially for Git operations and interactive clients.
7. **Transactions.** Make preview/dry-run and atomic apply first-class API
   concepts. Specify rollback guarantees for multi-file operations.
8. **Compatibility.** Choose a versioning and deprecation policy for the Go
   module before outside callers depend on package paths and exported symbols.
9. **Import-cycle direction.** Establish dependency rules before moving files;
   for example, model/parsing packages must not depend on CLI presentation.
10. **API dogfood.** Add at least one small external-style test or example that
    imports only public packages and performs the same workflow as the CLI.

The initial preference is one coherent workspace API used by the CLI, with
lower-level packages exported selectively. Merely renaming `internal/*` into
top-level directories would expose accidental details without defining a usable
contract.

## Specific Code Cleanup Backlog

1. Split `internal/cli/cli.go` into command definitions, argument parsing,
   execution adapters, and output formatting. Generate human help, completion,
   and future machine-readable command descriptions from one command schema.
2. Consolidate atomic file replacement. It currently exists in `manifest`,
   `render`, `renderfs`, and `state`, with different directory creation,
   permission, cleanup, and error-joining behavior.
3. Audit ignored errors. Distinguish mathematically infallible hash writes from
   best-effort temporary cleanup, explain the former, and surface or join the
   latter where failure matters.
4. Reject or safely contain symlinks for every local source, not only pinned URL
   sources. Add tests for a symlinked Markdown file, directory, and source root.
5. Make tool template variables reproducible. Do not let checkout directory or
   Git remote silently affect output unless the resolved values are written into
   the manifest or another explicit build input.
6. Make `status` the canonical read-only workspace report. Treat `drift` as a
   concise compatibility alias for its direct-edit view, while moving import,
   reject, inherit, and propose actions under explicitly directional mutation
   commands.
7. Finish the root dogfood migration from legacy `AGENTS.toml` to canonical
   `agents.yaml`, then regenerate `AGENTS.md` through the supported workflow.
8. Refresh stale handoff/roadmap statements as features merge; avoid recording
   branch facts that immediately become false after merge.
9. Add package documentation and public examples as APIs move. Tests should use
   public entry points where that is the behavior being promised, while focused
   package tests continue to cover internals of an exported package.
10. Pass cancellation and command execution through testable seams for Git
    source operations; preserve ordinary commands' no-network guarantee.
11. Remove duplicate or superseded backlog entries during the next TODO
    grooming pass without rewriting the historical Decision Intent log.
12. Run `gofmt`, focused tests, the public-API dogfood example, `go test ./...`,
    `go vet ./...`, and `errcheck ./...` for the eventual package move.

## Local-Source Symlink Safety

`library.Load` walks the declared local source root and collects paths ending in
`.md`. `filepath.WalkDir` does not follow a directory symlink, but a symlink
whose own name ends in `.md` is collected and later opened by `os.Open`, which
follows it. `os.Stat` also follows a symlink used as the source root.

For example, a declared source could contain:

```text
library/private.md -> ../../private-notes.md
```

Mogent would treat `private.md` as part of the library even though its bytes are
outside the declared root. This can violate provenance, privacy boundaries, and
the claim that a source location defines all readable content. It can also make
the same manifest render differently depending on machine-local symlink targets.

Pinned URL source hashing already rejects every symlink. The simplest consistent
rule is to reject symlinked roots, directories, and Markdown files for local
sources too. A more permissive alternative would resolve every link and prove
that its final target stays beneath the canonical source root, but that adds
platform-specific complexity and still weakens portability. Choose and record
the rule before changing source semantics.

## `repo_name` And `repo_url` Determinism

When the manifest does not define them, rendering currently derives:

- `repo_name` from the basename of the directory containing `agents.yaml`;
- `repo_url` from that checkout's `git remote.origin.url`, or an empty string if
  the command fails.

As a result, identical libraries and manifest bytes can render different output
when cloned into differently named directories, used from a copied workspace,
or configured with different remotes. That makes the design's shorthand
determinism statement—same libraries plus manifest gives the same output—true
only when those implicit variables are unused or happen to match.

Preferred direction: repository discovery can suggest values during `init`, but
the accepted values should be written under manifest `vars`. A build then uses
only explicit inputs. If implicit discovery remains available, it should be an
explicit, visible mode and part of the declared reproducibility contract.

## Status And Drift Direction

A small compatibility alias is useful if it reduces vocabulary without hiding
effects:

```text
mogent status                 # concise aggregate workspace state
mogent status --drift         # detailed generated-output/direct-edit state
mogent drift                  # shorthand for the same read-only view
```

The alias should remain read-only. Existing `drift --import` and
`drift --reject` behavior should migrate, with deprecation guidance, to commands
whose names state direction: preserve/import a hand edit locally, reject it,
inherit origin changes, or propose a local change upstream. Making `drift` an
alias for a mixture of reporting and mutation would save typing but make safety
and help text harder to understand.

## Skills Ecosystem Research

Managing repository skills, including a possible `.agents/skills/` surface, is
later work. Do not assume all ecosystems share one storage or execution model.
Research Codex, Claude Code, Gemini CLI, GitHub Copilot, Cursor, OpenCode, Goose,
and other relevant agent tools against a common matrix:

- repository-local, user-local, organization, and built-in discovery paths;
- skill manifest and instruction formats;
- names, versions, dependencies, and conflict rules;
- trust, permissions, executable hooks, and network behavior;
- installation, update, pinning, lockfiles, and removal;
- precedence and override rules across scopes;
- portability versus ecosystem-specific extensions;
- how skills coexist with `AGENTS.md` or equivalent instruction files;
- whether Mogent should compose, install, link, mirror, or only inventory them.

The first deliverable should be a dated comparison based on each ecosystem's
current primary documentation and a small real repository fixture. Only then
should Mogent decide whether `.agents/skills/` is a canonical authored format,
an interoperability projection, or one output among several.

## Dogfood And Owner Review

- Exercise each public API workflow through a Go example and through the CLI;
  compare results, diagnostics, and file effects.
- Continue real-repository CLI dogfooding for authoring, source selection,
  status, and reconciliation vocabulary.
- **Recurring owner reminder:** Quincy must read every checked-in agent module
  under `libraries/`, recording content corrections, metadata changes, overlap,
  and source-boundary concerns. Until complete, agents should mention this open
  review in substantive project handoffs and ask whether Quincy wants to take
  the next bounded library section.
