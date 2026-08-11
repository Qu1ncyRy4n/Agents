---
tldr: Protect authoritative media and use non-destructive, single-writer curation workflows.
---
# Media Archive

## No Physical Deletion

Never delete source or archive files, or add physical deletion behavior, without
explicit user approval. Prefer mark-only curation with a reason field when a
repo supports it.

## Authoritative Data

Treat the original media archive as authoritative. Databases, embeddings,
indexes, duplicate groups, and review UI state are derived aids unless the repo
documents a stronger guarantee.

## SQLite

Respect single-writer constraints. Do not run embedding, classification,
near-duplicate, or review processes concurrently when they write the same
SQLite database.

## Mounted Paths

Treat mounted drive paths as significant. Confirm expected mounts before resume
or batch operations; remounting can change paths and break assumptions.

## Validation

For DB-affecting changes, test on a temporary DB or fixture first. Inspect SQL
effects before running on an authoritative database.
