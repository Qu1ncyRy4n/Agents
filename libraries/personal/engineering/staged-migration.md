---
tldr: Replace an implementation in reviewable stages with explicit compatibility, cutover, and rollback evidence.
---
# Staged Migration

## Direction And Decision Gate

Use a staged replacement when an established implementation must remain usable
while a successor reaches reviewed parity. Record the migration direction,
authoritative implementation at each stage, compatibility promise, cutover
criteria, rollback path, and retirement criteria in the repository's design of
record before substantial successor implementation begins.

Do not let a provisional recommendation read as a settled architecture
decision. Mark each migration choice as **Decided**, **Proposed**, or **Open**,
and ask the user to resolve durable alternatives before locking them.

## Compatibility Inventory

Treat the old implementation as behavior evidence, not automatically correct
design. Inventory its user-visible and machine-visible contract: commands,
arguments, responses, errors, configuration and preference fields, persisted
state, files and paths, executables, environment assumptions, service units,
hotkeys, UI controls, dependencies, and lifecycle paths. Compare code,
documentation, packaging, and operational assumptions; record mismatches and
say whether each is preserved, intentionally corrected, or still open.

## Parity Contract

Define compatibility per public surface rather than declaring a whole program
"compatible". For every legacy operation, specify valid input, malformed
input, response, observable state change, error behavior, and compatibility
status. Keep fixtures or contract tests for those cases before claiming parity.

Preserve published interfaces until an approved compatibility change says
otherwise. A newer interface may run beside a legacy one during migration, but
its ownership, selection rule, conversion boundary, and eventual removal must
be explicit.

## Failure, Recovery, And Cutover

Design normal operation, malformed input, corrupted state, concurrent clients,
process crashes, restart, partial upgrades, and rollback before switching the
authoritative implementation. Name the owner of every mutable runtime state and
avoid two implementations silently writing the same state.

Gate cutover on reviewed evidence, not the existence of a new implementation:
contract coverage, migration/rollback rehearsal where practical, supported
activation path, and an explicit decision that the predecessor may lose
authority. Retire the predecessor only in a later, reviewable change after its
support boundary has ended.

## Handoff

Report the migration decisions made, open questions, compatibility deliberately
preserved or changed, runtime paths touched, fixtures/checks run, and what
remains before the next implementation slice is safe.

## See Also

Select `cdint:process/decision-first` for decision records and thought
experiments; `cdint:engineering/test-strategy` for deterministic contract
evidence; and `personal:engineering/local-service-design` when the compatibility
surface includes a daemon or local protocol.
