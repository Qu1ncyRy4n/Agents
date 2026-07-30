# Mogent Milestone-Two Navigator Contract

Status: active implementation contract.

This document defines the first interactive navigator after the milestone-one
renderer. Source: DI-tuvim.

## Scope

`mogent tui` loads an existing `agents.yaml`, lets the user inspect and draft
manifest changes, previews the result, and performs an explicit save-and-build
operation. It uses Bubble Tea for terminal interaction and Lip Gloss for
layout and color. Bubbles components are used only where they remove ordinary
terminal work such as scrolling or text input.

The navigator does not modify a source library as part of ordinary selection,
ordering, exclusion, save, or build operations. Copy-on-write editing and
direct-edit import remain separate future features.

## Layout

At 140 columns or wider, show three panes by default:

```text
Manifest tree | Final AGENTS.md | Selected-source context
```

- The tree shows inclusion state, collapsed/expanded state, and source alias in
  text. Color is secondary only.
- The Final pane contains the complete rendered `AGENTS.md`, positioned at and
  highlighting the selected tree node's output.
- The Source pane shows the selected source section with its parent heading and
  adjacent sibling sections. It labels the source alias and path.

Below 140 columns, show two panes:

```text
Manifest tree | Detail: Final or Source
```

The detail-pane header always names its mode. `v` switches its mode between
Final and Source. `Tab` changes keyboard focus between visible panes. `1`,
`2`, and `3` select Tree, Final, and Source directly; on a narrow display the
latter two select the detail mode. `Esc` returns focus to the tree.

Later, full-screen Tree, Final, and Source views retain these same direct view
keys and navigation rules.

## Selection And Provenance

Selecting a tree node updates both preview panes. A selected node can map to
multiple output locations. In that case the Final pane says, for example,
`Final: 2 locations`, and `[` and `]` move among them.

The interface always displays source aliases and paths. It must remain usable
when color is unavailable.

## Drafts And Writes

Selection, ordering, and exclusions first change an in-memory draft. While it
differs from the saved manifest, show `~` in the status area.

`~` does not mean that a source is local. `L` means the node uses an explicit
copy-on-write local override. Saving a normal draft removes `~` and retains the
ordinary source marker.

The explicit **Save and build** action confirms the change, atomically writes
`agents.yaml`, then runs the normal build path to update `AGENTS.md`. It never
edits a source-library Markdown file. A cancelled confirmation or quit discards
the in-memory draft without writing either file.

## Errors

When a draft cannot render, retain the last valid Final preview and display the
current validation error. A failed save or build leaves the current files
unchanged and reports the failure clearly.

## Acceptance Checks

The milestone-two demo is ready when these checks pass against the small
manifest and shared library in `demo_script.md`:

- `mogent tui -manifest agents.yaml` opens on an existing manifest without
  changing `agents.yaml`, `AGENTS.md`, or any source-library Markdown file.
- At 140 columns or wider, the default screen shows the manifest tree, complete
  rendered `AGENTS.md`, and selected-source context at the same time.
- Below 140 columns, the screen shows the manifest tree plus one detail pane,
  and the detail pane clearly switches between Final and Source with `v`.
- The tree shows selected nodes, collapsed or expanded state, and source
  provenance in text. The same information remains understandable without
  color.
- Moving the tree selection updates the Final and Source panes. The Final pane
  keeps the whole rendered document visible and aligns to the selected node.
  The Source pane shows the selected source section with nearby source context.
- Toggling a module changes only the in-memory draft, shows `~`, updates the
  Final preview, and does not write files before confirmation.
- Save and build asks for confirmation, atomically writes `agents.yaml`, then
  uses the normal build path to update `AGENTS.md`. After a successful save,
  the dirty marker disappears.
- Quitting or cancelling with a dirty draft leaves `agents.yaml` and
  `AGENTS.md` unchanged.
- A draft that cannot render keeps the last valid Final preview visible and
  shows the current validation error.
- A failed save or build leaves `agents.yaml`, `AGENTS.md`, and source-library
  files unchanged.

## Deferred

- Copy-on-write editing and import of direct `AGENTS.md` edits.
- Full source-file view toggle.
- Wide-screen user layout preferences.
- URL sources, pinning, cache, and upstream-change review.
