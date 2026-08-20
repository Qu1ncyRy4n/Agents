---
tags: [docs/consistency, instructions/documentation, scope/repo]
tldr: Keep docs and implementation aligned without unrelated rewrites.
priority: 0.7
scope: repo
---
# Documentation

## Audience And Purpose

Identify who the document serves and what job it performs before expanding it.
Distinguish tutorials and learning material from task recipes, reference,
explanation, specifications, and handoffs. Keep agent-facing operational
instructions direct; keep human-facing learning context when it helps the
reader understand rather than merely execute.

## Consistency

Keep design records, implementation, and public documentation consistent. Do
not rewrite unrelated prose for style. Preserve explanatory comments unless the
same change replaces them with a clearer explanation near the same logic.

When local documentation styles conflict, follow the authoritative document
and established convention for the changed area. Raise a consequential
inconsistency instead of silently normalizing the whole repository.

## Usage Evidence

When a user-visible workflow changes, include a representative command,
example, or before/after result where it makes the new behavior easier to
verify. Do not let an example become a second source of truth; update or remove
it when the authoritative interface changes.
