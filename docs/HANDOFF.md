# Mogent Handoff

Updated: 2026-08-31

This file is the stable navigation entry point for an agent handoff. It does
not duplicate task state, design decisions, or a workstream's detailed resume
instructions.

## Sources of Authority

- [`TODO/TODO.md`](../TODO/TODO.md) owns current tasks, owner decisions,
  statuses, and blockers.
- [`TODO/PICKUP-2026-08-25-libv2-intake.md`](../TODO/PICKUP-2026-08-25-libv2-intake.md)
  is the detailed resume checkpoint for the active libv2 workstream.
- [`DESIGN.md`](DESIGN.md) is the product design of record.
- [`LIBV2-PROTO-PROGRESS-REPORT.md`](LIBV2-PROTO-PROGRESS-REPORT.md) explains
  the libv2 problem, approach, organization, and open decisions.
- [`../libv2-proto/PROVENANCE.md`](../libv2-proto/PROVENANCE.md) and
  [`../libv2-proto/WORKLIST.md`](../libv2-proto/WORKLIST.md) track source intake.

If a chat recap or this navigation file conflicts with Git or an authority
above, inspect the authority and current working tree.

## Current Resume Point

Libv2 semantic intake is active. The first workflow and decision-governance
candidate modules are committed but pending owner review. Start with the dated
pickup, then review the candidate batch in its listed order.

The root `agents.yaml` and generated `AGENTS.md` are intentionally absent until
reviewed libv2 modules exist. Do not restore the deleted v1 guide or create the
new root manifest against an incomplete intake library.

Module include syntax remains undecided. `TE-nufad` narrowed the design space;
it did not authorize implementation. Relationship validity and merge
precedence also remain separate open questions.

## Document Convention

- `TODO/TODO.md`: durable, canonical work state.
- `TODO/PICKUP-<date>-<workstream>.md`: volatile, detailed checkpoint for one
  active workstream; replace or archive it when its resume point changes.
- `docs/HANDOFF.md`: short, stable map to the current authorities and pickup.
