# Mogent Localization And Drift Contract

Status: active implementation contract.

This milestone follows the reusable workspace extraction. It implements
localization and drift through core operations first; the CLI is the first
client and the existing TUI may adopt the same operations later. Source:
DI-ravam.

## Localization

`mogent localize <manifest-heading-path>` selects one manifest entry. If the
entry composes several sources, `--from <source-ref>` selects the reference to
replace.

The workspace operation:

1. resolves the selected shared node;
2. creates ordinary Markdown under `.mogent/library`, preserving the source
   path;
3. records origin metadata in `.mogent/provenance.yaml`;
4. declares `local: .mogent/library` if needed;
5. replaces only the selected manifest reference with `local:<path>`; and
6. validates the resulting render before committing the transaction.

The provenance record includes the local reference, original source reference,
manifest-declared source location, source-relative file, localization time, and
SHA-256 hash of the original selected content. Local files and shared source
files are separate after localization; upstream changes do not flow into the
local copy automatically.

`--dry-run` reports paths and previews without workspace writes. `--rebuild`
also updates generated output through normal overwrite protection. An existing
different local artifact is an error, never an implicit overwrite.

## Drift

`mogent drift` compares the generated render, the on-disk output, and recorded
generated-output state. It reports clean, stale manifest output, untracked
output, or direct edits and returns agent-readable details.

Drift resolution is explicit:

- `mogent drift --reject --force` rebuilds from the manifest after review;
- `mogent drift --import <manifest-heading-path>` localizes the selected node
  using the corresponding edited output section only when the mapping is
  unambiguous;
- ambiguous structural edits remain unresolved and do not write files.

Import never guesses across several manifest nodes or silently converts a whole
document into one override.

## Transaction And Recovery

Manifest, local Markdown, provenance, generated output, and state writes use
temporary files and rollback on failure. Errors name any path that could not be
restored. Shared source files are never transaction targets.

## Acceptance Checks

- Shared source bytes remain unchanged.
- Localize dry-run writes nothing.
- Localization creates the expected local path and provenance record.
- The manifest visibly changes to `local:` and still renders identically before
  user edits.
- Existing conflicting local files fail loudly.
- Drift reports all existing output states without mutation.
- Reject requires an explicit force flag.
- Import succeeds only for an unambiguous selected section and otherwise leaves
  all files unchanged.
