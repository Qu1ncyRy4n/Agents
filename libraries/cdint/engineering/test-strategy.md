---
tldr: Use focused deterministic tests, then broaden validation when risk or results warrant it.
---
# Test Strategy

## Scope

Scale validation to risk. Use focused checks for narrow changes and broader
suites when shared behavior, public interfaces, persistence, or user-facing
workflows change.

## Determinism

Keep tests deterministic and offline by default. Use fixtures, mocks, or local
servers instead of live services unless the task explicitly requires an
integration against a real external system.

## Coverage Shape

For new behavior, cover the main success path, meaningful error paths, and at
least one boundary or migration case. Keep tests close to the code they
exercise unless the repository has an established integration-test location.

## Assertions

Prefer structured assertions over manual string or JSON digging when helpers or
parsers exist. Compare complete objects when that is clearer than asserting
individual fields.

## UI And Text Output

When user-visible UI or text output changes, update or add snapshot, golden, or
before/after coverage when the repository uses that style. Review generated
snapshot changes before accepting them.

## Repository-Specific Evidence

The corpus includes Go `testing`, Rust `cargo`/`just` workflows, static-site
file-browser checks, notebook/data-pipeline smoke checks, and lab-machine smoke
tests. Keep the shared question stable: what evidence would reveal the likely
failure mode of this change?
