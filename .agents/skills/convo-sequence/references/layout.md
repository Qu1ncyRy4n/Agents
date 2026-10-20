# Conversation Folder And Templates

The owner chose a shared conversation folder, located under `docs/convo/`, with a
salient NOW entry and topics as archives. These templates implement that layout; checker details remain under
review. Preserve pre-existing narrative before migration.

```text
docs/convo/
  head.md
  work.md
  decisions.md
  archive/
    <topic>/
      history.md
```

## File Ownership

| File | Owns | Reading pattern |
|---|---|---|
| `head.md` | Goal, stage, one active ID, one concrete resume action | Read first on every re-entry |
| `work.md` | Stable work IDs, status, dependencies, completion evidence | Read the referenced item and small NEXT queue |
| `decisions.md` | Decision statements, status/evidence, supersession links | Read decisions relevant to the active item |
| `archive/<topic>/` | Detailed narrative, alternatives, examples, rationale, old checkpoints | Read as needed; no competing live authority |

The head references other records. It does not duplicate their complete task or
decision text. A contradiction between its active ID and the work ledger is a
consistency error to resolve, not two equally valid resume options.

## Head Template

```markdown
# NOW — <topic>

- **Goal:** <outcome>
- **Current stage:** S1 — <plain-language definition>
- **One active item:** W-001 — <one bounded item>
- **Next concrete action:** <observable action on W-001>
- **Blocked by:** <specific blocker or none>
- **Do not:** <scope boundaries and deferred work>
- **Last verified:** <check, result, date, revision/environment; or not yet verified>
- **Working-tree/commit state:** <branch/commit, relevant dirty work, commit mode>

---
Everything below is history/reference. Resume only from the NOW panel above.

Records: [work](work.md), [decisions](decisions.md), [archives](archive/README.md).
```

Use descriptive stage names directly when an opaque S-code adds no value. A
reported remote verification can be recorded with its source/revision and the
note "not rerun in this checkout".

## Work Template

```markdown
# Work Records

The current resume action lives in head.md. IDs remain stable when status changes.

| ID | State | Item | Depends on | Evidence |
|---|---|---|---|---|
| W-001 | ACTIVE | Confirm the current example | none | pending |
| W-002 | NEXT | Record the accepted contract | W-001 | pending |
| W-003 | NEXT | Define the bounded implementation slice | W-002 | pending |
| W-004 | PARKED | Broader integration research | none | deferred by owner |
| W-005 | DONE | Capture the original context | none | source reference |
```

Normally keep two or three NEXT items. Permit fewer when only one meaningful
step is known. Keep DONE records and evidence in this ledger; archive rich topic
details without creating separate live heads.

## Decision Template

```markdown
## <topic>-D001 — <short name>

Status: proposed
Statement: <bounded decision>
Evidence: <owner answer, code, or exact checks appropriate to the status>
Scope/reservations: <what remains unsettled>
Supersedes: none
Superseded by: none
```

An accepted replacement creates D002; D001 points to D002 and becomes superseded.
Keep the earlier reasoning in history. Acceptance alone does not imply code or
test evidence.

## Topic Archive Template

```markdown
# History And Rationale

Everything below is history/reference. Historical NEXT/resume text is not current.
Read ../../head.md for the authoritative re-entry point.

## Source preservation

<original source/revision or a complete faithful local capture>

## Dated narrative entries

<questions, examples, owner answers, reasoning, reservations, and evidence>
```

The old per-topic live-head layout is superseded. Link the one workspace
`docs/convo/head.md` prominently from README/handoff navigation. Old topic checkpoints
may remain as archival evidence, with a clear historical notice.

## Checker Review

After the file contract is agreed, define a small read-only checker and its
failure/warning examples. It should validate IDs, the head/ACTIVE match, queue
size, history boundary, and supersession links. Stale live references to finished
or superseded records should be flagged; historical pointers should remain intact.
Do not promise an automatic check of every natural-language contradiction.
