# TODO-pupab - HCL v2 Rollout

Implement the revised v2 model in `docs/proposals/MOGENT-HCL-V2-SPEC.md` far
enough that the four user stories pass against the QMR library. This file is
the implementer's checklist; the spec is the contract.

## Decision Intent Log

ID: DI-pupab
Date: 2026-10-05
Status: active
Author: Quincy Ryan
Decision: Drop the v2 lock file and pin Git sources by a `commit` attribute written into `mogent.hcl` (Go modules model). Replace the sidecar `inclusion { policy }` block with a flat `offer` attribute whose values are `foundation | choose | optional | opt_in`. Replace v2 `build` and `build --dry-run` with `plan` and `apply` (Terraform model). Keep `tags_all` / `tags_any` as the tag selection surface. Defer template variables to the feature set after this one.
Intent: Make the consumer file the single place a reader looks to understand what is selected and from which revision. Make the library's recommendation strength read as plain English and match Promise Theory: the library offers, the consumer selects, nothing is imposed.
Constraints: V1 `agents.yaml`, `mogent.lock.yaml`, and `mogent build` are untouched. `apply` may edit `mogent.hcl` only to add or replace a `commit` attribute, through `hclwrite`, preserving comments and formatting. Rendered output never changes because of this work except where the QMR sidecar gains `offer` and the consumer selection responds to it.
Affects: v2/, internal/cli/, docs/proposals/MOGENT-HCL-V2-SPEC.md, testing-ground/hcl-v2-demo/, qmr-agents-library/library.mogent.hcl

## Conventions For The Implementer

- One commit per task below, imperative subject, no trailers or signatures.
- `tools/check` must pass before each commit. `deadcode` and `errcheck` are
  required; use `nix develop -c tools/check` if they are not on PATH.
- Tests live beside code in `v2/` and `internal/cli/`. Every task lists the
  test it must add. Existing test helpers are in `v2/test_helpers_test.go`.
- Do not widen into v1 code, the TUI, or the navigator.
- Diagnostic codes are fixed by the spec table. Do not invent new ones without
  adding them to the spec.
- When a task's acceptance criterion cannot be met without a decision that is
  not in the spec, stop that task, record the question under Open Questions
  below, and continue with the next independent task.

## Tasks

Priority order. `[dep: X]` means the task cannot start before X is committed.

### P0 — vocabulary and CLI shape

- [ ] **V2-1 Rename `inclusion` block to `offer` attribute.**
  Sidecar schema: `offer` string attribute on a section; `defaults` map and
  `default` bool become sibling attributes of the section, not a sub-block.
  Values `foundation`, `choose`, `optional`, `opt_in`. Validation rules per
  spec: `defaults` complete and only for `choose`; `default` only for
  `optional`; neither on `foundation` or `opt_in`; `offer` only on a branch.
  Rename `Inclusion` type and field to `Offer`. Update planner switch cases.
  Reject `inclusion` blocks with a load error naming the replacement.
  Tests: load valid sidecar for each value; reject each invalid combination;
  reject old `inclusion` block with a helpful message. Update the fixture in
  `v2/test_helpers_test.go` and `testing-ground/hcl-v2-demo/`.

- [ ] **V2-2 Add `MOGENT107` and `MOGENT108`.** `force_exclude` on a child
  whose parent is not `foundation`, or without a non-empty `reason`, is an
  error. `accept_defaults` under a non-`choose` branch is an error. Today both
  are silently accepted. Add `MOGENT205` warning when a `foundation` child is
  force-excluded with a valid reason, and `MOGENT204` when an `opt_in` child
  is selected.
  Tests: one per code.

- [ ] **V2-3 Replace v2 `build` with `apply`; delete v2 `--dry-run`.**
  `apply [--config] [--force]`. `build` with `--config` returns an error
  pointing at `apply`. `build` with `--manifest` is untouched. `plan` gains
  nothing yet. Update `commands()` usage lines, completion output, and
  `README.md` v2 section. `v2.BuildLocal` loses its `dryRun` parameter.
  Tests: `apply` writes; `build --config` errors; v1 `build` tests unchanged.

- [ ] **V2-4 `plan` prints a unified diff per output path.** [dep: V2-3]
  Render in memory, read current file if present, print `--- a/path` /
  `+++ b/path` unified diff, or "new file" / "unchanged". Use a small internal
  diff; `workspace/diff.go` may already have one to reuse. Diagnostics after
  the diffs. Non-zero exit on any error diagnostic.
  Tests: new file, changed file, unchanged file.

- [ ] **V2-5 `MOGENT208` drift warning in `plan` and `apply`.** [dep: V2-4]
  If an output path exists and its content hash differs from
  `.mogent/state.json`, warn in `plan` and refuse in `apply` without
  `--force`. `state.Inspect` already classifies this for v1; reuse it.
  Tests: edited output warns in plan, refuses in apply, applies with force.

### P1 — Git sources and pinning

- [ ] **V2-6 Parse `commit` on `git` sources.** Full 40-hex only; `MOGENT109`
  otherwise. `commit` forbidden on `local`. `ref` optional, default `HEAD`.
  Tests: valid, short SHA rejected, `commit` on local rejected.

- [ ] **V2-7 Resolve and fetch `git` sources into `.mogent/sources/<alias>/<commit>/`.** [dep: V2-6]
  Reuse `sourcecache.fetchCandidate` mechanics (shallow clone at a ref,
  `rev-parse HEAD`) but do not touch `mogent.lock.yaml`. When `commit` is set
  and the checkout exists, read it with no network. When `commit` is set and
  absent, fetch exactly that commit. When `commit` is absent, resolve `ref`,
  emit `MOGENT206`, and expose the resolved commit on the plan. Honor
  `subdir`. `LoadLocalLibraries` becomes `LoadLibraries`.
  Tests: use a local bare repository created in `t.TempDir()` as the `git`
  URL (file path is acceptable for tests; validate `https://` in production
  through the existing `ValidateRemote`). Cover: pinned and cached, pinned
  and missing, unpinned.

- [ ] **V2-8 `apply` writes resolved `commit` into `mogent.hcl`.** [dep: V2-7]
  Use `hclwrite` to set the attribute on the exact `source` block. Preserve
  every other byte, including comments. Write atomically after outputs and
  state succeed; include the config file in the transaction snapshot so a
  failed apply restores it.
  Tests: commit written; comments preserved; failed state write restores
  config.

- [ ] **V2-9 `mogent update [alias] [--accept]`.** [dep: V2-8]
  Re-resolve `ref` for one or all `git` sources. Print old and new commit.
  Render the configured outputs under both and print the unified diff. With
  `--accept`, rewrite `commit` and nothing else. Without, write nothing.
  Tests: no change, change shown without accept, change accepted.

- [ ] **V2-10 `MOGENT207` local source Git info.** [dep: V2-7]
  If a `local` source directory is inside a Git work tree, print its HEAD
  short commit and dirty flag as an info diagnostic in `plan`. If `git` is
  not on PATH or the directory is not a repository, print nothing.
  Tests: one with a temp repo, one with a plain directory.

### P2 — discovery and library polish

- [ ] **V2-11 `MOGENT201` zero-match tag warning.** Each `tags_all` or
  `tags_any` query that matched no leaf is a warning naming the query.
  Tests: matching and non-matching query.

- [ ] **V2-12 `mogent source list [alias] [--tldr] [--tags]` for v2.**
  Tree view of the sidecar: path, title, `offer` with `default`/`defaults`,
  tags with `--tags`, TLDR with `--tldr`. Trees listed after sections with
  their entries. Reuse `internal/presentation` tree helpers. Dispatch: when
  `--config` is given or `mogent.hcl` exists and `agents.yaml` does not, use
  the v2 listing.
  Tests: golden output for the QMR fixture.

- [ ] **V2-13 `mogent source show alias:path` for v2.** [dep: V2-12]
  One section: path, title, offer, tags (effective, with inherited marked),
  TLDR, source file, first N body lines.
  Tests: branch and leaf.

- [ ] **V2-14 Promote QMR metadata into the live sidecar.** [dep: V2-1]
  In `qmr-agents-library`: move `tldr`, `tags`, and `offer` from
  `library-tmp.mogent.hcl` into `library.mogent.hcl` using the new
  vocabulary. Add tree `entry` blocks for every skill directory. Delete
  `library-tmp.mogent.hcl`. Update `README.md` there to drop the YAML
  template reference once V2-15 lands.
  Acceptance: `mogent plan` against the QMR library with `select = { all =
  true }` produces zero errors and the expected `MOGENT203` warnings for
  `workflow`.

- [ ] **V2-15 Replace `templates/core-and-constraints/agents.yaml` with an HCL consumer template.** [dep: V2-14]
  `mogent init --template qmr-core --source qmr=<path>` should produce a
  `mogent.hcl` that selects `intro` and `constraints` as foundation and
  `workflow` with `accept_defaults`. Keep the YAML file under
  `archive/` in the QMR repo.

### P3 — user story verification

- [ ] **V2-16 Script the four acceptance scenarios.** [dep: V2-9, V2-14]
  Add `tools/v2-stories` (sh) that creates a temp consumer directory, runs
  each scenario from the spec against `../qmr-agents-library`, and asserts
  exit codes and key output lines. The scenarios map to US-1, US-2, US-4,
  and US-5 in `docs/user-stories.md`; mark those stories' v2 status there. Not part of `tools/check`; documented in
  `docs/DOGFOOD.md`.

- [ ] **V2-17 Reconcile docs.** [dep: V2-16]
  `README.md` v2 section, `docs/DESIGN.md` §6 command table, and
  `docs/HANDOFF.md` resume point reflect `plan`/`apply`/`update`, `offer`,
  and no-lock pinning. Delete stale `mogent lock` mentions. `HCL-CONFIG-SPIKE`
  gets one line pointing at the revised spec. `docs/yaml_psudeo.md` open
  question 1 (selection in manifest only, sidecar describes only) is answered
  yes by v2; note that there. `docs/user-stories-ciwg.md` acceptance criteria
  still name `agents.yaml` and `build --dry-run`; add the v2 equivalents.

## Deferred To The Next Feature Set

- Template variables in content. Steve's ask; after V2-17.
- Required `offer` on every branch (spec open question).
- `plan --out`.
- Non-AGENTS outputs (slides). Scope undefined.
- V1-to-V2 migration command.

## Open Questions

Record blockers here as they appear, with the task number.

- (none yet)
