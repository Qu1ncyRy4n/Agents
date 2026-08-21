# Preserve Source-Path Moves Without Giving Every Module a Permanent ID

Status: provisional behavior proposal; no command or configuration name is
approved.

Pinning makes one library revision reproducible. It does not by itself preserve
a manifest reference after a later pinned revision renames that source path.

Prefer a sparse move map written only when a published path changes:

```yaml
moves:
  process/thought-experiment:
    to: process/narrow-options-with-thought-experiments
    deprecated_after: 2027-01-01
```

Desired behavior:

- an update preview reports the old and new path as a probable move;
- an explicit map makes the move deterministic rather than heuristic;
- the old path may act as a temporary compatibility pointer with a warning;
- the manifest can be updated automatically only under an explicitly reviewed
  policy;
- migration remains reversible while the old compatibility entry exists;
- after the deprecation date, behavior follows the configured policy rather
  than silently changing; and
- the update view shows rendered consequences before accepting the new pin.

A future user setting might permit reviewed automatic migration, but neither
its name nor its default belongs in this proposal. Do not introduce a new
command until existing status, update, localization, and drift operations have
been evaluated as possible homes for the workflow.

Run a later thought experiment covering pinned updates, forks, move chains,
cycles, expired mappings, rollback, and two old paths converging on one new
module. Record an owner decision before defining public schema or commands.

