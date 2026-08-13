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
mogent origin status [manifest-heading-path]
mogent reconcile <manifest-heading-path> [--dry-run]
mogent inherit <manifest-heading-path> [--accept]
mogent propose <manifest-heading-path> --to <writable-source> [--dry-run]
```

Names remain candidates, but every command must state direction. The current
`drift` readout can become a compatibility alias or migrate into this family
after dogfooding.

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

- Final command family and compatibility treatment for `drift`.
- Content-addressed base files versus embedded base content in provenance YAML.
- Conflict artifact format and lifecycle.
- Which local sources are considered writable and how that authority is
  declared.
