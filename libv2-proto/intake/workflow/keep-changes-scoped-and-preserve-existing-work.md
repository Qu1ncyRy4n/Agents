---
status: reliable
candidate_outputs: [agent-guide]
tags: [workflow/scope, workflow/user-work, review/diff, safety/non-destructive]
tldr: Change only what the task requires, preserve work already in the tree, and raise conflicts instead of silently rewriting either.
sources:
  - docs/other_repo_agents/ciwg_FAB26-Presentation_refs_heads_main_AGENTS.md#diff-discipline-required
  - docs/other_repo_agents/promisegrid_promisegrid_refs_heads_main_AGENTS.md#diff-discipline-required
  - docs/other_repo_agents/tbd/todo_app_AGENTS.md#agent-behavior
requires_all: []
mutually_exclusive: []
see_also: []
---
# Keep Changes Scoped and Preserve Existing Work

Make the smallest coherent change that satisfies the request or an explicitly
settled decision. Preserve pre-existing and concurrent work in the working tree.
When the requested change and existing work cannot both be preserved, show the
conflict and ask for direction rather than silently choosing an owner or
rewriting a broader area.

<!-- review-note: This candidate combines task scope and preservation of user
work because both govern the edit boundary. Review whether they always travel
together or should become two atomic modules. -->

<!-- review-note: Likely overlapping modules include preservation of explanatory
comments and classification of pre-existing changes. Add machine-readable
see-also relationships only after those modules exist at stable intake paths. -->

## Tie Every Change to the Request or a Settled Decision
<!-- tldr: Do not turn a focused task into an unrequested cleanup. -->

Verbatim source from
[`Diff Discipline`](../../../docs/other_repo_agents/ciwg_FAB26-Presentation_refs_heads_main_AGENTS.md#diff-discipline-required):

> Keep changes minimal and directly tied to the user's request or locked decision.
> Do not rewrite files from scratch, arbitrarily rewrap lines, reorder unrelated sections, or normalize prose style unless the user explicitly asks for that cleanup.

The same wording appears in
[`promisegrid/promisegrid`](../../../docs/other_repo_agents/promisegrid_promisegrid_refs_heads_main_AGENTS.md#diff-discipline-required).

Several source guides state a shorter version:

> Prefer small, focused edits and avoid rearranging files without a clear need.

Sources:
[`ciwg/cswg`](../../../docs/other_repo_agents/ciwg_cswg_refs_heads_main_AGENTS.md#agent-specific-notes),
[`ciwg/decomk-conf-cswg`](../../../docs/other_repo_agents/ciwg_decomk-conf-cswg_refs_heads_main_AGENTS.md#editing--documentation-standards),
[`ciwg/mob-sandbox`](../../../docs/other_repo_agents/ciwg_mob-sandbox_refs_heads_main_AGENTS.md#editing--documentation-standards),
[`stevegt/grokker`](../../../docs/other_repo_agents/stevegt_grokker_refs_heads_main_AGENTS.md#agent-specific-notes), and
[`stevegt/navlog`](../../../docs/other_repo_agents/stevegt_navlog_refs_heads_main_AGENTS.md#editing--documentation-standards).

## Preserve Changes Already in the Working Tree
<!-- tldr: Existing work is not cleanup material and does not become agent-owned merely because it is uncommitted. -->

Verbatim source from
[`promisegrid/promisegrid Diff Discipline`](../../../docs/other_repo_agents/promisegrid_promisegrid_refs_heads_main_AGENTS.md#diff-discipline-required):

> When user edits appear during a task, preserve them. Do not revert,
> overwrite, stage, unstage, commit, or "clean up" user edits unless the
> user explicitly asks for that exact action. If the edits affect the
> current task, stop and ask how to proceed.

Corroborating verbatim source from the generated
[`todo_app Agent Behavior`](<../../../docs/other_repo_agents/tbd/todo_app_AGENTS.md#agent-behavior>):

> Keep edits scoped to the request. Avoid broad rewrites while the project is in
> this exploratory stage.
>
> Do not overwrite user changes in the working tree.

The shorter `todo_app` guide was generated from repository inspection rather
than captured as authoritative upstream instruction. It supports discovery but
is not needed as the primary authority now that the stronger captured wording
has been found.

## Raise an Unavoidable Conflict Before Editing Through It
<!-- tldr: Preserve both sides when possible; otherwise identify the exact collision and request a choice. -->

Adaptation from the two source rules above:

If an existing change overlaps the requested edit, first determine whether both
can be retained. If they cannot, identify the affected file or behavior, explain
what each version is trying to preserve, and ask which direction owns the
result. Do not discard the existing change, abandon the requested behavior, or
silently combine incompatible intent.

This conflict procedure is not stated this explicitly in the cited source
guides. It is a proposed completion of their preservation rule and should remain
marked as adaptation until reviewed.

## Worked Examples
<!-- tldr: Scope is determined by behavioral necessity, not merely by the number of files changed. -->

Proposed example—unrelated formatting:

```text
Request: Fix configuration precedence.
Existing diff: A documentation draft changes headings in README.md.
Action: Change and test configuration loading; leave the README draft intact.
Do not rewrap or normalize README.md while nearby.
```

Proposed example—overlapping user edit:

```text
Request: Rename the `source add` flag.
Existing diff: The same parser branch contains an unfinished user change to
accept multiple sources.
Action: Inspect whether both behaviors can coexist. If the rename would erase
or reinterpret the unfinished behavior, show the collision and ask before
editing that branch.
```

Proposed counterexample—small is not the same as incomplete:

```text
Bad: Change one implementation file but leave its public documentation and
behavioral test incorrect because touching three files feels less minimal.

Better: Change the smallest complete set of implementation, tests, and
authoritative documentation required for one coherent behavior.
```

## Review Questions

1. Should task scope and preservation of existing work remain one compatible
   bundle, or become two atomic modules that usually select together?
2. Is “settled decision” general enough, or should the canonical module avoid
   assuming any particular Decision Intent system?
3. Should staging, unstaging, reverting, and committing user-owned changes live
   here, or in a separate tools/Git module that requires this workflow rule?
4. Is the proposed conflict procedure correct, or should this intake module
   preserve only source-verbatim instructions until a later editorial pass?
