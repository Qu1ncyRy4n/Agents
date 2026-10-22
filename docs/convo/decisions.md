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

## CS-D005 — Native catalog browsing belongs to `source show`

Status: implemented
Statement: In native plan/build projects, `mogent source show [alias:path]`
owns the source-catalog tree. Plain `mogent show` continues to show selected
contributions, with `show --build`, `show --issues`, and content preview as its
related views. `show --source` remains a compatibility alias.
Evidence: Teacher_Bot dogfooding made `show --source` and the legacy `source
show` visibly confusing. The owner proposed flipping the vocabulary; focused
CLI/assembly tests passed and the installed CLI was verified in Teacher_Bot on
2026-10-21.
Scope/reservations: Legacy YAML/HCL `source show` behavior remains available
when explicitly routed through those formats. Native upstream input provenance
is still undecided; this naming choice does not add it.
Supersedes: none
Superseded by: none

## CS-D006 — Module-level upstream-change detail

Status: implemented
Statement: Native upstream-change reporting must identify changed, added, and
removed resolved modules within each library, rather than only reporting that a
managed output differs. Each library has deterministic content fingerprints for
every node and for the library root. A consumer records observed fingerprints;
comparison starts at the root and recursively descends only through mismatched
nodes until the detailed differences are found.
Evidence: During the 2026-10-21 Teacher_Bot QMR/CDINT simulation, the owner
asked that module/within-library changes be detailed.
Scope/reservations: Content hashes are change fingerprints, not node identity:
stable source paths remain references. The consumer state stores its prior
comparison baseline; the library's fingerprints are deterministic from current
semantic node metadata, selected payload bytes/physical-tree layout, ordered
children, and library name. Mismatched ancestor hashes are suppressed when a
more-specific changed, added, or removed descendant is found. Older state has
no upstream baseline and requires a normal build to create one. The focused
assembly/state/CLI tests passed on 2026-10-21; Teacher_Bot remains to be
dogfooded without overriding owner-authored outputs.
Supersedes: none
Superseded by: none

## CS-D007 — Tag query conversion is not ordinary optimization

Status: proposed
Statement: Keep `module optimize` as a selection-editing, semantics-preserving
operation. A conversion from explicit module selections to a `tags_all` or
`tags_any` include query must be previewed and separately accepted, because the
query can select future upstream additions.
Evidence: During command-tree review, the owner proposed tag-search output that
can show runnable module-add commands and asked whether it could prove an
optimization-by-tagging method. Resolver tests already establish tag matching
and scoped path exclusion; no tag-search CLI exists.
Scope/reservations: Read-only source tag inspection and command rendering are
the next candidate slice. Negative tag matching, command spelling, and whether
conversion belongs in an optimizer or an explicit migration command remain open.
Supersedes: none
Superseded by: none
