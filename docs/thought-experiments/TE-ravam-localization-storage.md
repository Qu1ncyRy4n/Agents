# TE: Localization Storage And Core Operation

TE ID: TE-ravam
Status: concluded by DI-ravam

## Decision Under Test

How should Mogent create a copy-on-write local override while keeping shared
sources unchanged, provenance reviewable, and the operation reusable by CLI,
TUI, and future clients?

## Constraints

- The manifest must visibly switch from the shared reference to a local one.
- Local Markdown should remain normal source-library Markdown.
- Provenance must survive beyond chat or process memory.
- Dry runs must not mutate the workspace.
- A failed operation must not leave a manifest pointing at a missing override.
- The core operation must not depend on opening an interactive editor.

## Alternatives

### Put provenance in Markdown frontmatter

This keeps provenance next to content, but mixes source-selection metadata with
workspace history and would shape every local file around Mogent's internal
needs. It also makes upstream hashes and timestamps appear to apply to every
heading in the file.

### Put provenance in the manifest

This makes history visible but turns the document outline into a verbose
operation log. Provenance is about one local artifact, not its rendered heading,
and should not complicate the stable manifest schema.

### Use a workspace provenance sidecar

Store ordinary local Markdown under `.mogent/library` and record its origin in
`.mogent/provenance.yaml`. The manifest contains only the explicit `local:`
reference and source declaration. The sidecar can evolve independently and can
record original hashes without changing rendered content.

## Scenarios

For `shared:instructions/testing`, localization creates
`.mogent/library/instructions/testing.md` with `# Testing`, copies the selected
subtree, adds `local: .mogent/library` when needed, and changes only the chosen
manifest entry to `local:instructions/testing`. The source file remains
byte-for-byte unchanged.

If two manifest entries use the same shared reference, targeting a manifest
heading changes only that entry. A composed entry requires the caller to name
which source reference is being replaced.

If the local path already exists with different provenance or content, Mogent
fails instead of silently overwriting it. Updating an established local override
is a later explicit edit operation.

For a dry run, Mogent reports the manifest change, local path, provenance record,
and rendered preview without persisting workspace files.

## Conclusion

Use `.mogent/library`, the explicit `local` source alias, and a versioned
`.mogent/provenance.yaml` sidecar. Address the operation by manifest heading
path, with an optional source reference for composed entries. Keep editor launch
outside the core operation; identical-copy localization is useful on its own and
future clients can supply edited content explicitly.
