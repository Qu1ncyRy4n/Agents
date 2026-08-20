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
integration against a real external system. Offline means the default test does
not require the public internet or a live third-party account; an in-process or
loopback test server is still offline for this purpose.

## Coverage Shape

For new behavior, cover the main success path, meaningful error paths, and at
least one boundary or migration case. Keep tests close to the code they
exercise unless the repository has an established integration-test location.

The main success path is the ordinary intended use. For a config loader, that
might mean loading one valid file into the expected complete config object. A
useful boundary case might be an empty optional section; a meaningful error
case might be a duplicate key.

## Assertions

Prefer structured assertions over manual string or JSON digging when helpers or
parsers exist. Compare complete objects when that is clearer than asserting
individual fields. For example, parse JSON into the response type and compare
the full expected value instead of checking that the raw text contains three
field fragments.

## UI And Text Output

When user-visible UI or text output changes, update or add snapshot, golden, or
before/after coverage when the repository uses that style. Review generated
snapshot changes before accepting them.

A snapshot or golden test stores the complete expected output in a reviewed
file and compares future output against it. Update that file only when the full
change is intended, not merely to make a failing test green.

## Repository-Specific Evidence <!-- id: corpus-variants -->

Use the repository's documented language, UI, data-pipeline, or machine checks
when they exercise the affected risk. Keep the shared question stable: what
evidence would reveal the likely failure mode of this change?
