---
tldr: Build narrow typed Rust slices with focused checks and explicit compatibility boundaries.
---
# Rust

## Development

Run `cargo fmt`, `cargo check`, and focused `cargo test` unless the repository
wraps these with `just`, `make`, or another documented command. Use the wrapper
when it encodes repository policy.

## API Shape

Prefer typed structures and enums for domain concepts instead of strings spread
through command handling. Avoid boolean or ambiguous `Option` parameters that
make call sites hard to read; prefer enums, named methods, newtypes, or builder
style when that clarifies intent.

## Dependencies

Do not introduce a database, async runtime, GUI framework, or broad dependency
unless the change clearly needs it. If dependency files change, run the
repository's lockfile refresh command and include the generated updates.

## First Implementation Slice

For a Rust successor to an existing program, do not begin major implementation
until the repository has recorded the first slice's responsibility and its
public compatibility boundary. Resolve the crate layout, runtime choice,
process layout, configuration/state schema boundary, and the typed operation
that the slice owns. A conventional choice is not a substitute for a recorded
decision when these choices constrain later GUI, IPC, or migration work.

Keep the first slice narrow and independently testable. Put parsing and domain
operations behind typed interfaces so CLI, GUI, and local-service adapters can
share behavior without making the whole application asynchronous or
multi-crate prematurely.

## Tests

Keep tests deterministic. Prefer whole-object equality when it makes failures
clear. Avoid tests for static constants or negative tests for behavior that was
removed.

## Corpus Variants

Rust guidance ranges from early CLI prototypes to large multi-crate workspaces.
Small projects emphasize `cargo check` and plain output; large workspaces may
require `just test -p <crate>`, snapshot updates, schema generation, or Bazel
lock refreshes.
