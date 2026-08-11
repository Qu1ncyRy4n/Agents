---
tldr: Keep caches, generated output, and temporary runtime state out of source control.
---
# Runtime Artifacts

## Repository Cleanliness

Do not commit generated binaries, build output, local state, caches, audit logs,
database files, embeddings, coverage output, or other runtime artifacts. Extend
ignore rules when a new generated or sensitive path appears.

## Temporary Paths

Put temporary files, tool caches, build caches, and disposable test data under
`/tmp` or another repo-approved temporary root. Configure tools with
explicit cache paths when they would otherwise write into the repository.

## Lockfiles And Generated Schemas

Do not update lockfiles, generated schemas, or generated documentation unless
the requested change requires it. When they do change, summarize why and include
the command or source change that produced them.

## Local State

Treat files such as `.grok`, `.grok.lock`, local browser caches, local model
state, and machine-specific runtime files as non-source state. Do not use them
as authoritative project data unless the repository explicitly says so.
