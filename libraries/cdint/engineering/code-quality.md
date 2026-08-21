---
tldr: Prefer clear, focused changes that preserve repository conventions.
---
# Code Quality

## Local Design Fit

Prefer the repository's existing architecture, helper APIs, and naming style
over introducing a new pattern. Add an abstraction only when it removes real
complexity, reduces meaningful duplication, or matches an established local
pattern.

## Small Surfaces

Keep public APIs, crate/package exports, configuration keys, and command-line
surfaces as small as practical. Do not add test-only helpers to production code
when a test can use an existing public path or a local fixture.

## Dependency Restraint

Prefer standard libraries and well-known dependencies. Add or upgrade a
dependency only when the task needs it, and include the associated lockfile or
schema updates required by the repository.

## Module Size

Avoid growing large central files. Prefer adding a focused module or file when
new behavior would otherwise make a high-touch orchestration file harder to
review. Move tests and module documentation with extracted logic so invariants
stay near their owner.

## Helper Discipline

Do not create a small helper that is referenced only once unless it names a
non-obvious concept or isolates a risky boundary. Remove duplication when the
shared pattern is real, not just visually similar.

## Language Fit

The corpus contains both object-oriented guidance for Go-like projects
(`structs` and methods, avoid global state) and Rust-oriented guidance that
prefers enums, typed domain structures, and small explicit APIs. The shared
principle is not a single paradigm; it is readable ownership, small surfaces,
and avoiding stringly-typed behavior where the language offers a clearer shape.
