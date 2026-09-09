# Mogent Handoff

Updated: 2026-09-08

This file is the stable navigation entry point for an agent handoff. It does
not duplicate task state, design decisions, or a workstream's detailed resume
instructions.

## Sources of Authority

- [`TODO/TODO.md`](../TODO/TODO.md) owns current tasks, owner decisions,
  statuses, and blockers.
- [`TODO/PICKUP-2026-09-04-agents-library-dogfood.md`](../TODO/PICKUP-2026-09-04-agents-library-dogfood.md)
  is historical checkpoint context. It is not the current curation entry point.
- [`DESIGN.md`](DESIGN.md) is the product design of record.
- [`library-curation/archive/LIBV2-PROTO-PROGRESS-REPORT.md`](library-curation/archive/LIBV2-PROTO-PROGRESS-REPORT.md)
  preserves the historic libv2 problem, approach, organization, and open
  decisions.
- [`../libv2-proto/PROVENANCE.md`](../libv2-proto/PROVENANCE.md) and
  [`../libv2-proto/WORKLIST.md`](../libv2-proto/WORKLIST.md) track source intake.

If a chat recap or this navigation file conflicts with Git or an authority
above, inspect the authority and current working tree.

## Current Resume Point

The agents-library source review is paused pending Steve curation. For a request
such as "I'm Steve, I'm here to do the content curation that Quincy requested.
Please guide me through it.", read the [Steve library curation
handoff](STEVE-LIBRARY-CURATION-HANDOFF.md) first, then the [local curation
landing page](library-curation/README.md), [open questions](library-curation/open-questions.md),
and [wb review material](wb.md). Pause before rendering or promoting drafts.

The root `agents.yaml` and generated `AGENTS.md` are intentionally absent until
reviewed libv2 modules exist. Do not restore the deleted v1 guide or create the
new root manifest against an incomplete intake library.

Module include syntax remains undecided. `TE-nufad` narrowed the design space;
it did not authorize implementation. Relationship validity and merge
precedence also remain separate open questions.

## Document Convention

- `TODO/TODO.md`: durable, canonical work state.
- `TODO/PICKUP-<date>-<workstream>.md`: volatile checkpoint for one active
  workstream; it may point to a detailed checkpoint in the repository that owns
  the active work.
- `docs/HANDOFF.md`: short, stable map to the current authorities and pickup.
