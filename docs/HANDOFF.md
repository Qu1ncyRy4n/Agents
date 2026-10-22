# Mogent Handoff

**Current conversation: [NOW](convo/head.md).** This handoff is product
navigation/reference; older resume language below is not a competing live checkpoint.

Updated: 2026-10-21; source-oriented native browsing and upstream-change simulation recorded.

This file is the stable navigation entry point for an agent handoff. It does
not duplicate task state, design decisions, or a workstream's detailed resume
instructions.

## Current Dogfood Checkpoint

The detailed [Teacher_Bot native walkthrough handoff](dogfood/2026-10-20-teacher-bot-native-walkthrough-handoff.md)
records the current consumer flow, installed outputs/state, exact verification,
working-tree distinctions, deferred decisions, and safe re-entry. Read it after
the NOW panel when resuming the native plan/build dogfood work.

The current follow-up uses `mogent source tree` for the native catalog and
`mogent show tree` for selected contributions. An isolated QMR/CDINT
change/add/remove simulation confirms that `drift` identifies changed, added,
and removed upstream module paths as well as output state. [NOW](convo/head.md)
owns the remaining command-vocabulary review.

## Sources of Authority

- [`convo/head.md`](convo/head.md) owns the live conversation resume point;
  its linked work and decision records own the conversation queue and choices.
- [`TODO/TODO.md`](../TODO/TODO.md) owns the product backlog and related owner
  decisions, statuses, and blockers.
- [`TODO/PICKUP-2026-09-04-agents-library-dogfood.md`](../TODO/PICKUP-2026-09-04-agents-library-dogfood.md)
  is historical checkpoint context. It is not the current curation entry point.
- [`DESIGN.md`](DESIGN.md) is the product design of record.
- [`library-curation/archive/LIBV2-PROTO-PROGRESS-REPORT.md`](library-curation/archive/LIBV2-PROTO-PROGRESS-REPORT.md)
  preserves the historic libv2 problem, approach, organization, and open
  decisions.
- [`../libv2-proto/PROVENANCE.md`](../libv2-proto/PROVENANCE.md) and
  [`../libv2-proto/WORKLIST.md`](../libv2-proto/WORKLIST.md) track source intake.
- Git history records authorship from human developers only; commit messages
  cite no co-authors except human developers.

If a chat recap or this navigation file conflicts with Git or an authority
above, inspect the authority and current working tree.

## Native Product Baseline

Mogent's supported project interface is the native `mogent.build`/`mogent.plan`
language and `mogent.lib` libraries. The current commands and their supported
semantics are recorded in [CLI-COMMAND-INVENTORY.md](CLI-COMMAND-INVENTORY.md);
the active dogfooding state is in [NOW](convo/head.md).

The implemented core composes typed Markdown sections and complete physical
bundles from explicit source declarations and includes. Builds protect managed
`AGENTS.md` and selected `.agents/skills` outputs; `status` and `drift` inspect
output state and saved upstream source fingerprints. Localizing a selected
physical input makes one editable project-local copy without later automatic
recopying.

Richer curation/relationship CLI workflows, provenance reconciliation,
interactive UI work, global targets, and reverse inference remain future work.
The first-slice specification and walkthrough are design/history references,
not alternate live command references.

The language resolves `tags_all`/`tags_any`/`tags_none` and applies scoped path
exclusions after tag filtering. `source search` provides read-only inspection
and an explicit `--add-commands` rendering mode; `module add` accepts the same
tag predicates. Do not treat
rewriting explicit selections into a tag query as ordinary optimization: it
changes how future upstream additions are selected. `optimize` remains a
selection-editing command, never a build flag. `replace` is specified only as
an include-language alias/property and is intentionally deferred from this CLI
review.

## Near-Term Native Roadmap

1. Dogfood source tag search, `tags_none`, and `--add-commands` against QMR/CDINT.
   Confirm matching, result readability, and the distinction between commands for
   explicit leaves and a dynamic tag-query include.
2. Dogfood tag-aware `module add` with `--dry-run --diff`; do not add automatic
   explicit-to-tag compaction.
3. Dogfood the consolidated catalog/result/module command tree, then decide which
   compatibility aliases remain useful.
4. Complete localization provenance and upstream-reconciliation UX, then dogfood
   a localized physical module through an upstream change.
5. Dogfood a real output collision and replacement. Only then decide whether the
   CLI needs a replacement option on `module add` or a separate command.
6. Return to broader curation/defaults, relationship diagnostics, output adoption
   and recovery, source freshness/pinning, and richer interactive inspection.

Vendor adapters, global installation, protocol profiles, hooks, reverse
inference, and a full interactive UI follow the core consumer workflow; they do
not block this sequence.

CDINT source review remains paused pending Steve curation. For a request such
as "I'm Steve, I'm here to do the content curation that Quincy requested.
Please guide me through it.", read the [Steve library curation
handoff](STEVE-LIBRARY-CURATION-HANDOFF.md) first, then the [local curation
landing page](library-curation/README.md), [open questions](library-curation/open-questions.md),
and [wb review material](wb.md). Pause before rendering or promoting drafts.

Personal live-library dogfooding is owned by the sibling QMR Agents Library.
Its design workbench, templates, archive, and provenance are not CDINT policy.

The root `agents.yaml` and generated `AGENTS.md` remain intentionally absent:
this repository owns Mogent product code, not a selectable policy library. Do
not restore the deleted v1 guide or create a product-root manifest. Use a
reviewed consumer repository or the QMR library template for dogfood work.

Module include syntax remains undecided. `TE-nufad` narrowed the design space;
it did not authorize implementation. Relationship validity and merge
precedence also remain separate open questions.

## Document Convention

- `TODO/TODO.md`: durable, canonical work state.
- `TODO/<datetime>-conversation-stack.md`: indexed owner discussion, single-item
  resume point, preferences, proposals, and parked questions; not a grammar spec.
- `TODO/PICKUP-<date>-<workstream>.md`: volatile checkpoint for one active
  workstream; it may point to a detailed checkpoint in the repository that owns
  the active work.
- `docs/HANDOFF.md`: short, stable map to the current authorities and pickup.
