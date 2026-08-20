---
tldr: Use stable identifiers when work must coordinate across records or actors.
---
# Coordination IDs

## Record Types

Use the repository's vocabulary consistently. A TODO tracks work; a Thought
Experiment (`TE`) compares durable alternatives; a Decision Request (`DR`)
records an unresolved choice that needs an owner; and a Decision Intent (`DI`)
records the chosen direction and why.

## Proquint Handles

Use proquint handles for durable coordination artifacts when the repository has
adopted them. New TODO, TE, DR, and DI records should share one handle
namespace so one handle cannot accidentally name two different coordination
threads. The stable word-like handle lets people mint records independently
without coordinating a global integer counter or relying on a timestamp as
identity.

## Minting

Mint new handles with the repository's supported tool, usually
`tools/mint-handle`. Do not invent a handle by inspection when a minting tool
exists. Minting means generating a new candidate and checking it against the
existing namespace before assigning it. Existing numeric or timestamp IDs may
remain as historical records.

## TODO Files

When using handle-based TODOs, name files `TODO/TODO-<handle>-<slug>.md` and
keep `TODO/TODO.md` as the priority-sorted index. Inside TODO files, use
handle-prefixed subtasks where helpful.

## Decision History

Do not rewrite a filed decision so it appears that the current intent was always
the rule. If intent changes materially, add a new `DI-<handle>` that supersedes
the old one. An active index may separate current records from an archive or
mark status clearly, while the historical record remains available.

## Thought Experiments

Name thought experiments `docs/thought-experiments/TE-<handle>-<slug>.md`.
Treat filed TE documents as durable analysis records; use refinements or a
superseding TE for material changes. A TE informs a decision; it does not become
the decision merely because it recommends an option.
