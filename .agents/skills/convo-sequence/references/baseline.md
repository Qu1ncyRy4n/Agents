# Baseline Workflow And Source

## Primary Source

The owner supplied `TODO/2026-10-14-212422-conversation-stack.md` from the Mogent
repository. It was not present in the local checkout during this extraction.
The full 871-line record was retrieved read-only through GitHub and read before
drafting the skill.

- Repository revision observed: `d1235e13506b21ef1ebc15717b17eb58091822b9`.
- Source blob: `1e7167f862cf22042493f558adead7d8a2618e9d`.
- [Full preserved source at the observed revision](https://github.com/Qu1ncyRy4n/Agents/blob/d1235e13506b21ef1ebc15717b17eb58091822b9/TODO/2026-10-14-212422-conversation-stack.md).

This reference extracts the reusable process. The source remains the rich
Mogent-specific narrative; it has not been replaced, truncated, or migrated by
creating this skill. Import it into a local `history.md` only when the owner
chooses that migration and its target checkout. The current local archive index
is `docs/convo/archive/config-redesign/README.md`.

## What Made The Stack Useful

The source's "How To Use This Stack" section asks for low conversational load:
one concrete issue, record the answer, then advance. Tangents go to the parking
lot. Large questionnaires and implementation inferred from exploratory examples
are explicitly discouraged.

The record preserves:

1. Owner goals, preferences, constraints, and terminology reservations.
2. Small worked examples and predictions about their observable results.
3. Separate owner answers, recommendations, alternatives, and unresolved details.
4. A focus stack and detailed resume point.
5. A parking lot for research, integrations, and broader questions.
6. Evidence links to specs, acceptance cases, and actual implementation boundaries.
7. Staged implementation/testing instructions and explicit commit authorization.

Keep those properties. The improvement is separating live authority from the
historical discussion, not deleting the discussion or flattening it into a TODO.

## Observed Authority Conflicts

- Source line 5 resumes syntax-complete work at S2; lines 868-870 still instruct
  resuming item 3 and approving grammar before implementation.
- Lines 17-25 accept `mogent.lib`; lines 27-35 retain an earlier `library.mogent`
  proposal and pre-review instructions. Those earlier paragraphs need an explicit
  historical/supersession interpretation, not another "Latest owner" insertion.
- The opening reports syntax tests and `nix develop -c tools/check` passing, while
  older paragraphs describe no implementation. That is chronology, not permission
  to claim the same tests passed in the older local checkout.

These examples demonstrate record management. They are not general skill rules
about Mogent filenames or authorization to implement S2 from this workspace.

## Requested Improvements Under Review

- One-screen NOW panel with goal, stage, one active item/action, blocker, boundaries,
  verification, and worktree/commit state.
- A hard history/reference boundary.
- Stable decision IDs/statuses and explicit supersession links.
- One ACTIVE item, normally two or three NEXT items, the rest parked.
- Plain-language definitions beside stage labels such as S1/S2.
- A read-only structural check for competing live pointers and stale NEXT references.

Layout refinement: the owner preferred one shared `convo/` folder with `head.md`,
`work.md`, and `decisions.md`, with topics as archives. This supersedes the draft's
per-topic live heads. The primary goal is one clear, salient, accessible NOW entry
point and easy separation of outdated material.

Location refinement: authored records now use `docs/convo/`; the reusable skill
stays under `.agents/skills/`. The earlier root-layout discussion above is preserved
as history; current navigation points to `docs/convo/head.md`.

## Small Review Scenarios

These are proposed human-review cases, not executed skill evaluations:

1. **Resume amid contradictory history.** The head points to S2, while historical
   notes still say item 3. Expected: follow the verified head, define S2, retain
   the old text as history, and identify any evidence gap before acting.
2. **Supersede a name.** An accepted filename replaces an earlier proposal.
   Expected: stable decision IDs with supersession links; no rewriting the whole
   narrative or repeated competing "Latest owner" paragraphs.
3. **Partial implementation.** Syntax API tests passed; runtime/CLI are pending.
   Expected: scope verification precisely, keep one next stage active, and do not
   report the complete redesign delivered or tests rerun in another checkout.
