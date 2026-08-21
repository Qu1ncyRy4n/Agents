---
tldr: Use stable identifiers when work must coordinate across records or actors.
---
# Coordination IDs

## Proquint Handles

Use proquint handles for durable coordination artifacts when the repository has
adopted them. New TODO, TE, DR, and DI records should share one handle
namespace so references do not depend on timestamps or integer allocation.

## Minting

Mint new handles with the repository's supported tool, usually
`tools/mint-handle`. Do not invent a handle by inspection when a minting tool
exists. Existing numeric or timestamp IDs may remain as historical records.

## TODO Files

When using handle-based TODOs, name files `TODO/TODO-<handle>-<slug>.md` and
keep `TODO/TODO.md` as the priority-sorted index. Inside TODO files, use
handle-prefixed subtasks where helpful.

## Decision Intent

Record locked decisions as append-only `DI-<handle>` entries in the relevant
TODO or decision log. If intent changes, add a new DI that supersedes the old
one instead of rewriting history.

## Thought Experiments

Name thought experiments `docs/thought-experiments/TE-<handle>-<slug>.md`.
Treat filed TE documents as durable analysis records; use refinements or a
superseding TE for material changes.
