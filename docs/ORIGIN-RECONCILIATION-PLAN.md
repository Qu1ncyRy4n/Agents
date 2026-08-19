# Mogent Origin And Reconciliation Plan

Status: design-ready plan; command vocabulary and provenance storage require a
recorded decision before implementation.

## Goal

Let ordinary edits to generated or localized guidance be preserved, compared,
inherited from upstream, or proposed back to a writable source without silently
overwriting either side.

## Model

For each localized module, reconciliation needs three durable inputs:

1. **Base**: the exact content localized or last reconciled.
2. **Origin**: the current source content at its pinned revision or local path.
3. **Local**: the repository-owned customized module.

The existing provenance hash proves identity but is insufficient for a
three-way comparison by itself. The next schema must preserve retrievable base
content, likely as content-addressed files under `.mogent/bases/`, while keeping
workspace history out of rendered Markdown.

## Candidate Surface

```text
mogent status [--manifest agents.yaml]
mogent drift [--manifest agents.yaml] # read-only status shorthand
mogent reconcile <manifest-heading-path> [--dry-run]
mogent preserve <manifest-heading-path> [--as heading] [--dry-run]
mogent inherit <manifest-heading-path> (--ff-only|--merge|--replace-local) [--accept]
mogent propose <manifest-heading-path> --to <writable-source> [--dry-run]
```

Names remain candidates, but every command must state direction. The current
`drift` readout can become a compatibility alias or migrate into this family
after dogfooding.

## Vocabulary And Direction

`reconcile` means "compare the relationship and prepare a safe choice," not
"move content in one predetermined direction." It examines:

```text
base -> origin changes
base -> local changes
```

and reports whether the changes are identical, independent, overlapping, or
based on an origin identity that no longer matches provenance. A reconcile
preview may recommend a direction, but does not choose one implicitly.

The directional operations are:

- **preserve**: take a direct edit from generated output and make it durable as
  repository-local source content;
- **inherit**: bring a newer origin version into a localized module;
- **propose**: turn a local change into a reviewable change aimed at a writable
  origin.

`status` remains the canonical read-only workspace view. `drift` should become
a compatibility shorthand for its direct-output-edit view. Existing
`drift --import` maps naturally to `preserve`; existing
`drift --reject --force` should migrate to the already explicit `build --force`.

## Preserve Shape

The normal preserve operation replaces the selected manifest reference with a
localized copy, as today's conservative drift import does.

A second mode may keep the original selection and add the edited content as a
locally named sibling:

```text
mogent preserve Instructions/Testing --as "Testing (Local Variant)" --after Instructions/Testing --dry-run
```

This is a document branch, not necessarily a Git branch. Use `--as` (plus the
existing placement vocabulary) rather than `--branch`, which readers will
reasonably interpret as version control. The preview must show both manifest
entries and the resulting rendered sections. It must reject duplicate headings
or ambiguous section mappings.

## Inherit Modes

Inheritance needs explicit strategy rather than one vague merge operation:

- `--ff-only`: accept origin only when local content is unchanged from base;
- `--merge`: perform a three-way merge when origin and local edits do not
  overlap, and otherwise produce a reviewable conflict artifact without writes;
- `--replace-local`: replace the local override with reviewed origin content,
  requiring explicit acceptance because local changes are discarded.

The safe first implementation is `--ff-only`, followed by dry-run three-way
classification. Automatic non-overlapping merge can come after the base storage
and conflict artifact contracts are proven through dogfooding.

## Proposal Adapters

The core proposal operation does not require messaging or a forge account. Its
first job is to produce a patch or change bundle against an explicitly writable
source checkout. That artifact can be reviewed or transported manually.

Adapters can add delivery later:

- a Git adapter creates a named branch and optional commit in a user-approved
  writable checkout, but does not push or open a change request implicitly;
- a forge adapter may open a pull/change request only with explicit authority;
- a future Promise Grid adapter may advertise or route the proposal without
  making the filesystem reconciliation model depend on Promise Grid.

Private code is therefore not blocked on a shared messaging system. It can stop
at a local patch or branch. Source writability, credentials, commit messages,
remote publication, and recipient routing remain separate permissions.

## Required Behaviors

- Preserve a direct generated-file edit by extracting the selected section into
  a local override when mapping is unambiguous.
- Report base-to-local and base-to-origin changes separately.
- Auto-apply inheritance only when changes do not overlap; otherwise produce a
  reviewable conflict artifact without changing source, local, or output files.
- `propose` creates a patch against an explicitly writable source. Git branch,
  commit, or forge integration is a later adapter, never an implicit side
  effect of reconciliation.
- `status` reports pinned revision, local divergence, and the last explicit
  upstream check. It never performs network access itself.
- Localized URL provenance records immutable revision and lock identity.

## Safety And Transactions

- Ordinary source libraries and pinned caches remain read-only.
- Every mutation has dry-run output and validates the complete resulting render.
- Base/provenance/local writes commit atomically or not at all.
- Reject missing bases, changed source identity, ambiguous section mapping, and
  concurrent local changes rather than guessing.
- Personal/private origins remain excluded from generated handoff records by
  default.

## Implementation Sequence

1. Decide and migrate provenance v2 with durable base content and URL lock
   identity.
2. Add read-only `origin status` and three-way classification.
3. Add direct-edit-to-local preservation using the existing single-section
   mapper.
4. Add non-conflicting `inherit` with patch preview.
5. Add source-targeted `propose` that emits a patch only.
6. Add optional Git/forge adapters after the filesystem workflow is stable.

## Decisions Required

- Whether directional verbs remain top-level or live under an `origin` command
  family. `status` plus read-only `drift` alias is the current preference.
- Content-addressed base files versus embedded base content in provenance YAML.
- Conflict artifact format and lifecycle.
- Which local sources are considered writable and how that authority is
  declared.
- Final names for preserve-as-sibling and inheritance strategies.
- Proposal artifact format and the boundary between core, Git, forge, and
  Promise Grid adapters.
