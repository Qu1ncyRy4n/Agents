# Mogent Authoring Workflow Plan

Status: implementation-ready plan; command details remain reviewable until the
first implementation slice.

## Goal

Make source declaration, manifest placement, reordering, and module creation
available through reusable workspace operations. Every mutation must support a
non-writing preview and preserve the core/CLI/TUI boundary.

## Proposed Surface

```text
mogent source add <alias> <path-or-url> [--subdir path] [--dry-run]
mogent move <manifest-heading-path> (--before path | --after path | --under path [--first|--last]) [--dry-run]
mogent module new --source <alias> --path <path> --heading <text> --tldr <text> [metadata flags] [--dry-run]
```

`add` already provides the shared placement vocabulary: `--under` with
`--first|--last`, plus `--before` and `--after`. `move` should reuse the same
core insertion primitive rather than implementing separate ordering rules.

## Source Add

- Add an explicit alias to `agents.yaml`; never create an ambient search path.
- Validate duplicate aliases, local directory readability, URL syntax, and
  `subdir` safety before previewing.
- Local sources can be added without network access.
- URL sources require a separate explicit pin step in the first slice. A later
  `--pin` convenience may compose source declaration and pinning as one
  transactional operation, but must still preview network and lock effects.
- Adding a source does not select any module.

## Move And Reorder

- Address the moved node and anchors by full manifest heading path.
- Reject ambiguous paths, moving a node beneath itself, and sibling heading
  collisions.
- Preserve the complete entry, including children, `from`, and exclusions.
- Preview both the manifest tree and rendered-output patch.

## Module Creation

- Permit writes only to an explicitly selected writable local source.
- Reject URL caches, pinned checkouts, unresolved aliases, existing paths, and
  any destination that escapes the source root.
- Collect heading, Markdown content, TLDR, destination, tags, priority, scope,
  requirements, and conflicts. Heading/content/TLDR/destination are the first
  required guided steps; other metadata remains optional.
- Preview the file path, resulting source reference, frontmatter, and Markdown
  before writing.
- Do not automatically select the new module; print the exact `mogent add`
  command as the next hint.

## Implementation Sequence

1. Extract reusable manifest location/removal/insertion operations from `add`.
2. Implement and test local `source add` with dry-run.
3. Implement `move` using the shared placement operations.
4. Implement noninteractive `module new` flags and deterministic previews.
5. Add an optional guided prompt client over the same operation.
6. Extend completion and dogfood documentation.

## Decisions Before URL Convenience

- Whether `source add URL --pin` is one recoverable transaction or remains two
  explicit commands.
- Whether module creation can target a separate writable library repository
  outside the consuming project without crossing an ownership boundary.
- The first stable metadata flag names and whether content comes from a flag,
  stdin, or an editor handoff.
