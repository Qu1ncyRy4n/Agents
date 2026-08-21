# Library And Policy Review Handoff

Date: 2026-08-20
Status: awaiting owner review; no decisions recorded and no changes committed
Base revision: `2ef81b2` (`main` and `origin/main` when this handoff was written)

## Purpose

This handoff preserves the discussion and current diffs while the owner cannot
review them. The work is intentionally uncommitted. Do not treat a thought
experiment's recommendation as approval, implement a public metadata change,
move library paths, or commit this review set until the owner has reviewed it.

The pending work has two independently reviewable units:

1. a decision map and two open thought experiments about module shape,
   relationships, and organization policy; and
2. a cleanup of the previous library-clarity commit so questions are raised to
   the owner instead of being answered implicitly in reusable instructions.

## Review Unit A: Proposed Decision Documents

These are new, untracked project documents created for review:

- `docs/DECISION-MAP-LIBRARY-POLICY.md` orders the coupled decisions: module
  boundary, metadata location, relationship semantics, organization policy,
  reconciliation, required-update review, and proposal delivery.
- `docs/thought-experiments/TE-mavok-library-module-shape.md` works through a
  hybrid model: node facts in file frontmatter, spanning groups and presets in
  a source-level `library.yaml`, and physical directories as taxonomy.
- `docs/thought-experiments/TE-sulap-organization-policy-composition.md`
  compares visible per-repository requirements, explicitly trusted pinned
  policy profiles, and ambient organization detection.

The handles `mavok` and `sulap` were assigned manually because
`tools/mint-handle` is unavailable, using the existing owner-approved exception.

### Current Recommendations, Not Decisions

- Keep guidance in one file when it is always selected as one compatible unit.
  Split optional, independently selectable, or mutually exclusive material
  into separate one-top-level-heading files.
- Use frontmatter for intrinsic node facts. Reserve a source-level
  `library.yaml` for groups, presets, and other facts spanning nodes.
- Keep `self:` references for same-source relationships and defer cross-library
  relationships.
- Define `requires` as a rendered-content requirement, not a hyperlink or a
  first-item-wins rule.
- Define `exclusive_group` as at most one selected member. Express a required
  choice separately with `requires_one_of` or policy/preset semantics.
- Keep compatible modules metadata-free by default. Treat `see_also` as
  discovery only. Defer generalized logical expressions and rich provenance.
- Begin organization-policy dogfooding with visible `policy.required` entries
  in each repository. Consider a pinned trusted policy profile only after
  repetition demonstrates the need. Do not infer policy from a Git remote.
- Keep necessity and trust separate: `required` content is mandatory to
  select, but a new revision still requires verification. `--required-only`
  may filter an update preview; it must not silently accept updates.
- Prefer the initial public reconciliation verbs `status`, `preserve`,
  `inherit`, and `propose`; keep three-way reconcile as an internal operation
  until it has a distinct user job.

### Owner Questions, In Dependency Order

1. Is the proposed file boundary correct: compatible bundle together; optional,
   independently selectable, or mutually exclusive content split out?
2. Is the frontmatter plus optional source-level `library.yaml` split a useful
   single model, or should all relationships live on nodes?
3. Are the proposed meanings of `requires`, `conflicts_with`,
   `exclusive_group`, and `requires_one_of` correct?
4. Is visible per-repository `policy.required` sufficient for the first
   organization-policy dogfood version?
5. Does organization policy require exact placement as well as presence, and
   where must its origin remain visible?
6. Should the first reconciliation interface omit both a public `reconcile`
   command and a `drift` alias?
7. For proposals, should the first milestone stop at a patch, or also prepare a
   local Git branch and commit? Opening or pushing a pull request remains a
   separately authorized later operation.

Review these questions in order. Module identity and metadata ownership need to
be stable before public formats or command behavior can depend on them.

## Review Unit B: Library Cleanup Diff

This tracked, unstaged diff is relative to `2ef81b2`. Its purpose is to remove
unapproved answers added in `ef11e9e`, remove the confusing
`<!-- id: corpus-variants -->` compatibility IDs, and retain only clearer
visible headings where those headings are still useful.

### Library Files

- `libraries/cdint/engineering/code-quality.md`: removes the compatibility ID
  from `Language Fit` and restores the earlier comparison-oriented prose.
- `libraries/cdint/engineering/coordination-ids.md`: removes the newly added
  record-type and explanatory expansions, restoring the earlier compact rules.
- `libraries/cdint/engineering/error-handling.md`: removes the newly added
  narrative examples pending owner decisions about reusable library content.
- `libraries/cdint/engineering/reviewability.md`: keeps `Local Size Policies`,
  removes its compatibility ID, and restores the earlier prose.
- `libraries/cdint/engineering/test-strategy.md`: removes new definitions and
  examples, keeps `Repository-Specific Evidence`, removes its compatibility ID,
  and restores the earlier prose.
- `libraries/cdint/shared-baseline/instructions/documentation.md`: removes the
  unapproved audience, consistency, and usage-evidence expansion.
- `libraries/cdint/shared-baseline/instructions/focused-change-loop.md`: removes
  the newly added conflict and usage-evidence instructions.
- `libraries/cdint/shared-baseline/instructions/git.md`: removes the imperative
  subject example.
- `libraries/personal/lang/nix.md`: keeps `Selection Boundary`, removes its
  compatibility ID, and restores the earlier comparison-oriented prose.
- `libraries/personal/lang/python.md`: keeps `Environment Fit`, removes its
  compatibility ID, and restores the earlier comparison-oriented prose.
- `libraries/personal/lang/rust.md`: keeps `Workspace Fit`, removes its
  compatibility ID, and restores the earlier comparison-oriented prose.

The restored prose still contains extraction/audit language such as "the
corpus." That is not presented as the final cleanup. Review whether each useful
distinction should become direct instruction, move to audit documentation, or
be removed during the later bounded reorganization.

### Progress Records

- `docs/LIBRARY-COVERAGE-AUDIT.md`: removes the progress entry that said the
  reverted clarity slice was complete.
- `docs/HANDOFF.md`: removes the matching completion claim and points to this
  pending review bundle.

### Questions To Apply To Each Library Diff

1. Is the removed explanation useful reusable guidance, or was it merely an
   answer to an owner question?
2. If useful, does it belong in the rendered library, authoring documentation,
   or a selection-time explanation from `source show`?
3. Is the retained heading a user-facing concept worth selecting, or only
   extraction/audit residue?
4. Does keeping or removing the section change a source path that consuming
   manifests may reference?

## Exact Review Commands

Review the existing tracked cleanup without showing unrelated user files:

```sh
git diff -- \
  docs/HANDOFF.md \
  docs/LIBRARY-COVERAGE-AUDIT.md \
  libraries/cdint/engineering/code-quality.md \
  libraries/cdint/engineering/coordination-ids.md \
  libraries/cdint/engineering/error-handling.md \
  libraries/cdint/engineering/reviewability.md \
  libraries/cdint/engineering/test-strategy.md \
  libraries/cdint/shared-baseline/instructions/documentation.md \
  libraries/cdint/shared-baseline/instructions/focused-change-loop.md \
  libraries/cdint/shared-baseline/instructions/git.md \
  libraries/personal/lang/nix.md \
  libraries/personal/lang/python.md \
  libraries/personal/lang/rust.md
```

New files are not included by ordinary `git diff`. Review them directly:

```sh
sed -n '1,260p' docs/DECISION-MAP-LIBRARY-POLICY.md
sed -n '1,320p' docs/thought-experiments/TE-mavok-library-module-shape.md
sed -n '1,320p' docs/thought-experiments/TE-sulap-organization-policy-composition.md
sed -n '1,360p' docs/REVIEW-HANDOFF-2026-08-20.md
```

Useful historical comparison:

```sh
git show --stat --oneline ef11e9e
git show ef11e9e -- libraries docs/LIBRARY-COVERAGE-AUDIT.md docs/HANDOFF.md
git show --stat --oneline 2ef81b2
```

## Validation Already Run

Before this handoff was written:

- `git diff --check` passed.
- `go test ./library ./render ./workspace` passed with Go caches under `/tmp`.
- `go run ./cmd/mogent source list --manifest testing-ground/basic-org/agents.yaml --tree --coverage`
  completed successfully.
- A search found no `corpus-variants` references in libraries or manifests.
  One historical mention remains in `TE-lusim` and should not be rewritten as
  though the old analysis never happened.
- No checked-in manifest references the removed IDs.

Rerun `git diff --check` after reviewing or editing this bundle. The work is
documentation/library-only; no public format or behavior implementation has
been made.

## Working-Tree Ownership Boundary

The following current paths are user-owned and are not part of this review
bundle:

- staged deletion: `AGENTS.toml`;
- unstaged deletion: `AGENTS.md`;
- untracked: `2026-08-20_quickpres.md`;
- untracked: `docs/notes_on_lib_mods.md`;
- untracked: `docs/steve-feedback.txt`; and
- untracked: `responsescratch.md`.

Do not restore, edit, stage, commit, summarize, or otherwise combine those paths
with this work without explicit owner direction. In particular, any later
commit must stage the approved review paths explicitly and leave the existing
staged `AGENTS.toml` deletion alone unless the owner says to include it.

## Resume Point

When the owner is ready, start with owner question 1 in Review Unit A. Record
the approved direction as Decision Intent before editing the public metadata
format or reorganizing source paths. Review Unit B can be accepted or revised
independently; do not use its presence to infer approval of either thought
experiment.
