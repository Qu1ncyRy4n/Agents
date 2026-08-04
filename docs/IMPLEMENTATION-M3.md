# Mogent Milestone-Three Localization And Drift Contract

Status: planned implementation contract.

Milestone three follows the committed milestone-two navigator. It adds the
parts needed to safely adapt shared instructions locally and handle direct
edits to generated `AGENTS.md` files. Source: DI-sovar.

## Scope

M3 is three ordered slices:

1. Copy-on-write local editing.
2. Direct `AGENTS.md` drift import.
3. Richer save/history support.

The navigator remains the main interaction surface. `mogent build` keeps the
milestone-one behavior unless a slice explicitly changes it.

## Slice 1: Copy-On-Write Local Editing

Editing a selected shared source node does not modify that shared library.
Instead, Mogent:

- copies the selected source node into a local library under `.mogent/library`;
- preserves the source heading path when practical;
- applies the user's edit to that local copy;
- changes the manifest reference from the shared source alias to `local:`;
- rebuilds the preview from the draft manifest; and
- shows `L` in the tree for the local override.

If the manifest does not already define a `local` source, Mogent adds one that
points to `.mogent/library`.

## Slice 2: Direct-Edit Drift Import

When `AGENTS.md` differs from the last generated-output hash, Mogent does not
guess. It shows the user that direct edits exist and offers explicit choices:

- keep the direct edit and stop;
- reject the direct edit by rebuilding from the manifest; or
- localize the edit by choosing a manifest node and creating a copy-on-write
  local override.

The first import implementation may be conservative. If Mogent cannot clearly
map a direct edit to one manifest node, it should ask the user to choose the
node or leave the edit unimported.

## Slice 3: Richer Save And History

Before writing, Mogent should make the pending change visible enough to review:

- manifest diff for `agents.yaml`;
- generated-output diff for `AGENTS.md`;
- clear explanation of which local override files will be created or changed.

Failed writes must leave the manifest, generated output, and local override
files in their previous state whenever practical. If a full rollback is not
possible, Mogent must report the exact paths and next action.

## Deferred

- URL source fetching, pinning, cache, and upstream-change review.
- Tags, generated source indexes, and swap-alternative groups.
- Promote-local-to-shared workflows.
- Promise Grid / CID-addressed sources.

## Acceptance Checks

- Editing a shared node creates or updates a local Markdown file under
  `.mogent/library` and changes the manifest reference to `local:`.
- The original shared library file remains byte-for-byte unchanged.
- The tree shows `L` for local overrides and `~` only for unsaved in-memory
  drafts.
- Cancelling a copy-on-write edit leaves `agents.yaml`, `AGENTS.md`, and
  `.mogent/library` unchanged.
- Direct edits to `AGENTS.md` are detected before overwrite.
- Drift import can create a local override from a user-selected manifest node.
- Rejecting drift rebuilds from the manifest only after explicit confirmation.
- Save/history output shows enough diff context to review the manifest,
  generated output, and local override file changes.
