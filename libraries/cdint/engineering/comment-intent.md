---
tldr: Preserve comments that explain intent and replace them only with equal or better context.
---
# Comment Intent

## Preservation

Never remove existing explanatory comments unless the same change replaces them
with equal-or-better explanation near the same logic. When refactoring, port the
old intent first, then improve wording.

## Non-Obvious Logic

Add short plain-English comments for non-obvious parser, encoding, rendering,
concurrency, storage, protocol, or recovery logic. Avoid comments that merely
repeat the code.

## Behavior Provenance

For non-trivial behavior changes in repos that use Decision Intent records, add
a behavior-level comment with `Intent:` and `Source: DI-<handle>` where it will
help future maintainers understand why the behavior exists.

## Audit

After editing code in strict repos, audit comment deltas with a focused diff,
for example:

```sh
git diff -U0 -- <file> | rg -n '^-\s*//|^-\s*/\*|^\+\s*//|^\+\s*/\*'
```

Resolve removed-comment lines before finalizing unless the user explicitly
approved dropping the comment.
