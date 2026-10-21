# Teacher_Bot Native Walkthrough Handoff

Status: detailed handoff for the 2026-10-20 native Mogent dogfooding session.
Current live resume point: [docs/convo/head.md](../convo/head.md). This document
is evidence and context, not a competing resume instruction.

## Outcome

Teacher_Bot was used as a from-scratch native-format pilot. The owner initialized
an authored `mogent.build`, connected QMR, browsed and previewed content, chose
guidance/workflow/skills, rendered outputs, recorded state, exercised read-only
drift, and committed the project result:

```text
Teacher_Bot main: 878de7c skills added
```

The project currently has committed `mogent.build`, `AGENTS.md`, `.agents/skills/`
and `.mogent/state.json`. No uncommitted Teacher_Bot work was reported at handoff
inspection.

## Verified Consumer Flow

The installed CLI is `/Users/q/go/bin/mogent`; it was rebuilt with `tools/install`
after each implementation checkpoint.

```sh
cd ~/dev/apps/Teacher_Bot

# Initialization and source discovery
mogent init
mogent source add qmr ../../cdint/qmr-agents-library
mogent show --source --tree --chars unicode
mogent show qmr:intro/assistant-agent --content
mogent show qmr:intro --content --headings

# Selection and curation
mogent add qmr:intro
mogent add qmr:workflow
mogent except qmr:workflow/decision-first
mogent except qmr:workflow/developer-decision-involvement-level
mogent add qmr:skills --all

# Inspection, render, state and drift
mogent show --build
mogent show --issues
mogent build
mogent build --check
mogent drift --diff
```

Final verification before handoff reported:

```text
Build check passed
Output: AGENTS.md
Output baseline: clean
Current build: matches
No AGENTS.md diff: current build matches disk.
```

## Delivered Native Behavior

### Initialization And Sources

- Plain `mogent init` creates authored `mogent.build`; `--plan` creates intent.
  Legacy YAML/HCL initialization requires explicit `--manifest` / `--config`.
- Optional Gum initialization can choose remembered libraries, local paths, Git
  URLs, or skip. Noninteractive/dry-run/explicit-source runs bypass Gum.
- `source add`, `source list|ls`, `source register`, and `source list --known`
  support project declarations plus a user-level discovery registry.
- The registry is user-local YAML only. Projects retain explicit local/Git source
  declarations and never depend on the registry at build time.
- `show --source` presents an aligned tree; `--flat`, `--chars`, and source-alias
  scoping work. `show --source qmr` means the full qmr catalog.

### Selection And Preview

- `add`, `select`, and `module add` are current aliases for appending a validated
  native include. They preserve existing authored bytes ahead of the new include.
- `show source:path --content [--lines N]` previews logical rendered sections or
  physical payloads without selecting/installing. `--headings` works for both
  rendered sections and physical Markdown files.
- `except` is the preferred explicit-exception command; `exclude` is an alias.
  It locates the narrowest direct selected ancestor and writes a relative exclude.
  For a `choose` branch it writes `accept_defaults = true`, meaning defaults except
  the named descendant.
- QMR skills are now `foundation` rather than opt-in. Opt-in is demonstrated by
  `cdint-demo-lib:experimental/review-policy`. A bare opt-in include now appears
  in `show --issues` as an `OPT_IN` warning with `--all`/explicit-child guidance.
- `add source:path --all` creates or upgrades a bare direct include to its current
  direct child selection. It is required for opt-in physical trees and refuses to
  overwrite an existing explicit child selection.

### Build, State, And Drift

- Authored and generated native builds now both write `.mogent/state.json` with
  SHA-256 output hashes for AGENTS.md and managed skill files.
- Generated plan builds additionally record generated-build and resolved-input
  hashes; authored builds do **not** yet record a source-input baseline. State is
  output ownership/drift evidence, not the source of configuration authority.
- `build --check` verifies all managed output bytes, modes, and directory layout.
- Native `drift [--diff]` is read-only. It reports the AGENTS baseline state and
  whether AGENTS matches the current composition; `--diff` prints the latter.
  Skill-directory detail remains in `build --check`.
- A routing bug was fixed: automatic entry-file discovery previously dropped
  `build --check` flags and performed a build. Discovered builds now pass their
  flags through; plan checks now check the workspace and print success.

### Help And Formatting

- `-h` and `--help` are identical concise workflow help.
- `help --all` expands common flags; `help <stage>` / command `--help` provides
  focused advanced help. `help commands` remains the full registry reference.
- Terminal help bolds headings/commands/flags; piped output, `NO_COLOR`, and
  `TERM=dumb` are plain.
- `reorder [--dry-run]` is implemented but intentionally was not applied in
  Teacher_Bot. It groups bindings, then sources, then includes while preserving
  include order and attached options. It is not semantic optimization.

## Current Working Trees

At handoff inspection:

| Repository | Branch / state | Notes |
|---|---|---|
| `Agents` | `main` at `110c213`, substantial uncommitted implementation/docs work | Do not commit unless owner requests. `tools/check`, `tools/install`, and `git diff --check` passed after the final code changes. |
| `Teacher_Bot` | `main` at `878de7c`, clean | Includes generated skills/state and selected native configuration. |
| `qmr-agents-library` | `main`, ahead 3, modified `README.md` and `mogent.lib` | Existing uncommitted README/catalog additions were preserved. This session changes `skills` from `opt_in` to `foundation`. |
| `cdint-demo-lib` | `main`, modified `mogent.lib` | This session changes its skills tree to `foundation` and adds the `experimental` opt-in example. |

The Mogent worktree also retains the unrelated paused operation-inference stash
`5e826857a3a64633be5bf2625cd22e827c2a30af`; do not restore/drop it automatically.

## Deferred Decisions

These are recorded in [COMMAND-RECOVERY-AUDIT.md](../COMMAND-RECOVERY-AUDIT.md).

1. Canonicalize `add` / `select` / `module add`; they are aliases during
   dogfooding.
2. Add `rm` / `module rm` for one direct include. Do not alias it to `exclude`.
3. Add overlap warnings and an explicit `optimise` proposal. Parent plus child
   includes can deliberately duplicate content; do not rewrite automatically.
4. Consider native skill-directory diffs in `drift`; current `build --check`
   already verifies all skill bytes/modes/layout.
5. Decide authored source-input baseline/provenance storage separately from
   output state.
6. Add native metadata/search/TLDR/coverage and project-aware completion through
   shared native inspection rather than legacy manifest routing.
7. Viewer/editor integration remains deferred. Current pipes work:
   `mogent show qmr:intro --content | mdt` and `... | hx -`.

## Re-entry

1. Read [docs/convo/head.md](../convo/head.md), then its active work record.
2. Inspect all four worktrees above before editing; preserve unrelated changes.
3. Review/commit the owner-approved Mogent/QMR/demo-library changes only when
   explicitly requested.
4. Choose one deferred item from the audit. The recommended next dogfood item is
   native `status`/all-output drift detail, then `rm` and overlap optimization.
