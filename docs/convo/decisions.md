# Conversation Decisions

These records own the conversation-workflow choices. Read [head.md](head.md) for
the current action. Product DR/DI/spec records retain their own responsibilities;
this ledger does not redefine the Mogent language or promote remote test claims
to locally verified implementation.

## CS-D001 — One shared convo root and salient NOW entry

Status: superseded
Statement: Use convo/head.md as the sole current conversation entry point, with
work.md and decisions.md beside it. Keep topic narratives and superseded
checkpoints under convo/archive/, rather than a separate live head per topic.
Evidence: Owner requested "convo/ ... head work decisions ... with topics as archives"
and "a clear and salient and accessible NOW entry point" with outdated material
easily separated.
Scope/reservations: The files are now drafted. The owner review of the populated
panel and the automatic checker remains active; file creation is not proof of a
completed re-entry evaluation.
Supersedes: CS-D000
Superseded by: CS-D004

## CS-D000 — Per-topic live heads

Status: superseded
Statement: The initial draft proposed docs/convo-<topic>/head.md plus work,
decisions, and history files per topic.
Evidence: Initial skill layout proposal, preserved in archive/convo-sequence/initial-layout.md.
Scope/reservations: This was an agent proposal, not a previously accepted owner policy.
Supersedes: none
Superseded by: CS-D001

## CS-D002 — Consistency checking

Status: proposed
Statement: Add a read-only check for one current head/ACTIVE match, a small NEXT
queue, valid IDs, stale live pointers, history boundaries, and supersession links.
Evidence: Owner's supplied criticisms request automatic checks; the skill records
a candidate structural contract.
Scope/reservations: No checker is implemented yet. Decide exact errors/warnings
and closure behavior before claiming automated enforcement. Historical NEXT text
is allowed in archives and must not be treated as live control.
Supersedes: none
Superseded by: none

## CS-D003 — Concise chat, rich records, one actionable step

Status: implemented
Statement: Keep routine chat short and result-first, with one independently
actionable step rather than a multi-phase workstream labelled as one item. Keep
the larger roadmap and detailed reasoning in records. Validation directions name
one behavior, a short check set, a pass condition, and the relevant stop condition.
Evidence: Owner supplied the overlong Spec33h/V1 example and requested an update;
skill instructions and references/response-style.md were updated and deployed to Seam.
Scope/reservations: Metadata validation and build checks passed; these do not prove
response quality. Full plans, required TEs, and requested explanations can remain
substantial. Review the next real replies rather than claiming a completed agent benchmark.
Supersedes: none
Superseded by: none

## CS-D004 — Author conversation records under docs/convo

Status: implemented
Statement: Keep the single NOW entry at docs/convo/head.md, with work.md,
decisions.md, and topic archives beside it. The reusable skill stays in
.agents/skills/convo-sequence/. Keep a prominent README link to NOW.
Evidence: Owner requested conversation records under docs or .agents; docs/convo
was recommended for human-authored records and implemented during Git cleanup.
Scope/reservations: This changes the folder location, not the accepted authority
split. Automatic semantic consistency checking remains a separate proposed item.
Supersedes: CS-D001
Superseded by: none
