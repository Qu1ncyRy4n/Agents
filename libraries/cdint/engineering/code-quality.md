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

## Language Fit <!-- id: corpus-variants -->

Follow the language and repository's established design idioms rather than
forcing one paradigm across stacks. Prefer readable ownership, small surfaces,
and typed domain concepts over stringly-typed behavior when the language offers
a clearer shape.
