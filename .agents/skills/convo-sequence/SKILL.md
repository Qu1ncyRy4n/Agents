---
name: convo-sequence
description: Maintain a resumable, decision-aware conversation sequence with one salient docs/convo/head.md NOW entry point and rich topic archives. Use when a user asks to manage a conversation stack, walk through design decisions piece by piece, resume a long discussion, separate current work from historical notes, or reconcile stale next/resume pointers. Also use when creating or updating docs/convo/head.md, work.md, decisions.md, and topic archives.
---
# Convo Sequence

Draft workflow with an owner-accepted shared conversation layout under `docs/`. The existing
conversation stack is the baseline; checker details and bulk narrative migration
still need review. See [the baseline and source evidence](references/baseline.md).

## Purpose

Help the owner re-enter a long-running conversation without rereading its whole
history, while retaining that history for recovering rationale. Keep live control
state and narrative distinct. A shorter document is not the goal: clear authority
and low conversational load are.

## Chat Is Small; Records Can Be Rich

Lead with the answer, recommendation, or next action. In ordinary chat, aim for
roughly 80–180 words, fewer when enough. Use short paragraphs and one small example
when it clarifies the issue. This is a default rhythm, not a cap on a requested
full TE, specification, substantive explanation, or consequential disagreement.

One item means **one independently actionable step with an observable result**.
It does not mean naming an entire workstream "V1" and putting implementation,
instrumentation tests, trace capture, diagnosis, and windowed comparison beneath
it. If an answer contains several stages, deliver the first stage and keep the
rest in NEXT/history unless the owner requested the whole plan.

Keep detailed rationale, alternate branches, and the longer roadmap in files.
Do not narrate ledger maintenance, repeat the review agenda, recite NOW fields,
or explain this skill on every turn. Update the records quietly and mention a
change only when it affects the owner's next action or answers their question.
Respect required repository role labels and actual milestone reports without
turning every small reply into a stage handoff.

Ask at most one focused question, and only when the answer changes what happens
next. An actionable answer does not need a closing question. After giving the
current step, wait for its result or the owner's steering rather than walking
through the remaining queue preemptively.

### Validation Directions

Give the next validation step, not the whole validation campaign:

1. Name the specific existing behavior or deliverable being checked.
2. Give a short set of checks for that one behavior (normally three or four).
   Use exact known commands when available; do not invent project commands.
3. State the observable pass condition and the relevant scope stop condition.
4. Ask for the result if needed, then stop. Future trace runs, tuning, and other
   stages belong in NEXT, not nested beneath a purported "one item" response.

If a prerequisite is missing or unknown, resolve that one prerequisite first.
Do not commission unapproved implementation under a request for validation.
Keep outcomes factual: a test definition is not a test run, and green core tests
do not prove a visual/runtime deliverable that they never inspect.

See [the response example and review cases](references/response-style.md).

## Preserve The Established Flow

1. Read the current checkpoint and enough relevant history to recover context.
2. On re-entry or a material scope change, briefly state the goal/stage and one
   active issue. On normal continuation, proceed directly to that issue.
3. Discuss one concrete question using a small example and an observable result.
   Ask the owner to predict the example when testing comprehensibility.
4. Record the answer, its scope, the reasoning, reservations, and open details.
   Keep owner direction separate from agent recommendations and examples.
5. Advance one item. Put tangents in the parking lot rather than expanding into
   another questionnaire or restarting settled format/product research.
6. At interruption or handoff, update the checkpoint, task states, and changed
   decision references together. Retain the detailed narrative.

Agree the decision-involvement level and implementation stage when relevant.
The baseline's software sequence is API implementation -> API tests -> CLI ->
CLI behavior tests, with expected CLI cases defined before wiring. Apply the
repository's actual agreed process; do not invent implementation stages for a
discussion-only task or infer authorization to commit from implementation approval.

## Re-entry And Evidence

Read `head.md` first, then the referenced active work item and decision records.
Open history on demand to explain a decision or resolve a contradiction.

Compare the checkpoint with the actual branch/worktree before acting. Record
committed, uncommitted, implemented, and verified separately. A checkpoint from
another revision or machine is evidence of what was reported there, not proof
that the current checkout has the same code or passing checks.

If the records disagree, identify the conflict and reconcile it from owner
answers and concrete evidence. Do not choose whichever "Latest owner" paragraph
looks newest, or resume an old narrative instruction as though it were current.

## Two-Layer Layout

Use one authored conversation folder per workspace, normally `docs/convo/`:

```text
head.md       current control panel and the only resume pointer
work.md       stable task records: ACTIVE, NEXT, PARKED, DONE
decisions.md  stable decision IDs, current status, evidence, supersession
archive/      topic narratives, examples, rationale, and old checkpoints
```

The path is `docs/convo/head.md`; link it prominently from the repository entry page.
Keep topics under `docs/convo/archive/<topic>/`, not separate live topic heads.
Keep the reusable skill under `.agents/skills/`; its files are not the workspace's
conversation records.
See [the layout and templates](references/layout.md). Reuse existing records and
links rather than duplicating DR/DI/product backlogs; adapt an existing workspace
layout with the owner before migrating it.

`head.md` owns the current goal/stage/resume action. `work.md` owns task status;
the head's active ID references that record and must agree with it. Decision
records own decisions. Topic archives supply rationale, never a competing current
resume instruction. These are distinct responsibilities, not several competing
summaries of the same control state.

## NOW Panel

Keep the live panel to about one screen:

```markdown
# NOW — <topic>

- **Goal:** <outcome>
- **Current stage:** <ID — plain-language meaning>
- **One active item:** <stable work ID — one bounded item>
- **Next concrete action:** <one observable action on that item>
- **Blocked by:** <specific dependency/owner question, or none>
- **Do not:** <current boundaries and explicitly deferred work>
- **Last verified:** <check/result, date, revision or working-tree scope>
- **Working-tree/commit state:** <branch, commit, relevant uncommitted work>

---
Everything below is history/reference. Resume only from the NOW panel above.
```

Define opaque stage labels beside the resume point: write "S2 — typed catalog
and resolver API", not just "S2". Treat reported checks as reported until actually
verified; do not invent dates, commit hashes, or pass claims.

Long narrative belongs in `archive/<topic>/`, with an explicit history/reference
notice so copied excerpts do not become accidental instructions. A retained old
single-file stack can remain intact as an archived snapshot or pinned reference.

## Work Queue

Use stable IDs that do not change when an item moves or its display order changes.
Keep exactly one ACTIVE item while the sequence is open, a short NEXT queue
(normally two or three, never more than three), and everything else PARKED.
Do not manufacture NEXT work to fill a quota. A blocked active item stays visible
with its blocker; do not silently start another active thread.

Record completed work as DONE with evidence. "Discussed", "accepted", "implemented",
and "verified" are different achievements. When closing a sequence explicitly,
record CLOSED and no active item rather than leaving a false live resume point.
Completion/closure conventions remain part of the layout review.

Keep one `work.md` ledger initially. Completed items retain stable IDs and
evidence in DONE; superseded topic detail moves to the archive. A later split
must not introduce competing current status or break record identity.

## Decision Records

Give each decision a stable topic-qualified ID and one current status:

- **proposed:** under consideration, not authorization to implement.
- **accepted:** owner approved the stated scope; exact unsettled details remain open.
- **implemented:** code/artifact evidence exists for the accepted scope.
- **verified:** named checks passed for a specified revision/environment/scope.
- **superseded:** a named later decision replaces this one.

Record the statement, owner/evidence reference, date, scope/reservations, and
`supersedes` / `superseded by` links. A change of meaning gets a new decision ID;
progress from accepted to implemented/verified retains the ID and adds evidence.
Reuse an existing DR/DI identifier when available rather than inventing a rival.

Do not let "looks good" approve every field in an exploratory example. Record
what the answer actually settled. Failed or stale verification must be visible,
not covered by an old verified label. For partial delivery, state the delivered
scope rather than calling the complete feature implemented.

## Update And Migration Discipline

After an owner answer or completed action:

1. Update the affected work/decision record and its evidence.
2. Update `head.md` to reference the one current active item and next action.
3. Append the detailed narrative or a dated refinement, preserving earlier meaning.
4. Mark obsolete live pointers historical or superseded; do not keep appending
   additional "Latest owner" control paragraphs.
5. Check consistency before reporting the checkpoint as current.

For an existing stack, preserve a complete original record or immutable source
reference before restructuring it. Resolve competing pointers from evidence;
don't guess. A useful migration adds authority without deleting the rationale.

## Proposed Automatic Check

Once the layout is agreed, provide a read-only check after sequence updates:

- One NOW panel/current resume location; no second live `Resume:` outside history.
- One ACTIVE work item, matching the head's active ID.
- At most three NEXT items; unique task/decision IDs and valid references.
- Stale live NEXT/action references to DONE, missing, or superseded items flagged.
- Every short stage label defined; history/reference boundary present.
- Supersession links point to existing decisions and do not form cycles.

Ignore obsolete instructions inside explicitly historical material as live
commands. Report them as history, not as an excuse to erase the narrative.
Mechanical checks catch structural contradictions, not every semantic conflict.
No checker, hook, or CI integration is claimed implemented by this draft.

## Walk Through Changes With The Owner

Start by reviewing the file ownership and NOW panel on one concrete example.
Then review the work ledger, decision statuses/supersession, and checker behavior
one at a time. Keep one active review question. Do not bulk-restructure their
conversation or resume the paused product implementation during the review.
