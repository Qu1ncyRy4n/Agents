# Active Mogent Dogfood Repositories

Status: active registry; maintained by the repository owner during real-project
use.

Use this file to track working repositories where Mogent is actually being
used. `testing-ground/` remains a disposable smoke-test collection for now; it
is not evidence that a workflow is pleasant or useful in a real repository.

Do not record credentials, private source content, or sensitive absolute paths.
A local-only path may be described with a stable nickname when the repository
itself is private.

## Active Repositories

Add one row per active experiment:

| Repository / nickname | Stack and stage | Mogent features in use | Questions or friction | Last exercised |
|---|---|---|---|---|
| _Add repository_ | _Language, maturity, migration state_ | _Sources, init, build, localization, etc._ | _What are you trying to learn?_ | _YYYY-MM-DD_ |

## Per-Repository Notes

For each row, record only enough context to reproduce the product question:

- the intended prompt size and audience;
- selected source aliases and broad module groups;
- commands or UI paths being exercised;
- what felt natural, confusing, repetitive, or unsafe;
- whether the result improved actual agent work;
- the next bounded experiment.

Prefer notes about Mogent behavior over copying the repository's private prompt
or source content here. Put detailed project-specific facts in that repository.

## Review Loop

1. Check the target repository's Git and generated-output state.
2. Exercise one bounded Mogent workflow.
3. Record expected versus actual behavior here or in a linked issue.
4. Turn repeated friction into a design/TODO item.
5. Keep one-off preference differences as local presentation configuration
   rather than project composition.
