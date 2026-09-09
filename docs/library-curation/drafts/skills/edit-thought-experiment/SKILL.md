---
name: edit-thought-experiment
description: Edit an already-filed thought experiment while preserving historical analysis, classifying mechanical changes, appending refinements, and superseding substantive conclusions. Use whenever changing a filed TE.
---

# Edit A Filed Thought Experiment

## Sources

- [`promisegrid_wire-lab_refs_heads_main_AGENTS.md`](../../reference/source-guides/promisegrid_wire-lab_refs_heads_main_AGENTS.md#te-editing-policy-required)
- [`newest_promisegrid.md`](../../reference/source-guides/newest_promisegrid.md#te-editing-policy-required)

## Read The Relevant Record

Read the affected TE, locally present TEs it directly cites or that directly
cite it, and locally present editing-policy refinements relevant to the edit.
For an obviously mechanical typo or demonstrably moved current path, one-TE
reading is sufficient. Missing historical or external policy artifacts do not
block the edit when the locally captured policy is complete.

## Classify Each Change

- **Cat-1a, current-pointer path:** update a path naming a file's current
  location in place. No top-of-file note is required.
- **Cat-1b, historical-quotation path:** preserve a path that records an earlier
  corpus state. This includes quotations, attributed statements, past-tense
  descriptions, refinements, supersedence notes, and decision-status history.
  When uncertain, classify the path as historical.
- **Cat-2, vocabulary with unchanged meaning:** update in place, but preserve
  historical quotations. Add a top-of-file note naming the driving TE or task
  and every DI in the file, explicitly stating that their meaning is unchanged.
  State that the file has no DIs when applicable. Search the corpus for quoted
  or attributed uses and classify each before changing it.
- **Cat-3, navigational pointer:** append a dated `## Refinements` entry telling
  readers where to look now. Do not rewrite the historical body or file a DI
  solely for this pointer.
- **Cat-4, resolved implication:** append a dated `## Refinements` entry when a
  listed implication resolves through a task, decision request, or downstream
  TE. Do not rewrite the historical body.
- **Cat-5, Cat-6, or Cat-7, substantive change:** create a new TE when a locked
  decision's meaning, scope, or applicability changes. The new TE carries its
  own decisions. Update only the older TE's status and decision-status pointers
  to name the superseding TE and decision.

## Preserve Structure

Every TE has a top-of-file `## Status` field immediately after its TE ID. Prefer
these canonical values for new TEs:

- `needs DF`
- `decided`
- `decided, refined`
- `superseded by TE-<id> / DI-<id>`
- `withdrawn`

Preserve established legacy statuses while retrofitting old records. Keep one
append-only `## Refinements` section after `## Decision status`; order entries
chronologically and title each `### YYYY-MM-DD - <title>`.

## Verify

Confirm that historical claims remain historical, every changed pointer really
describes current state, Cat-2 notes account for all embedded decisions, and
substantive changes live in a superseding TE rather than a rewritten body.
