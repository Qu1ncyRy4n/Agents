# PICKUP: Libv2 Semantic Intake and Include Decisions

Date: 2026-08-25
Status: active; content candidates are tracked but pending owner review

## Canonical Task Authority

`TODO/TODO.md` is the canonical detailed Mogent task list, owner-decision queue,
and dependency graph. Detailed tasks were migrated from W34 section 1.1 in
commit `c701667`; the Mogent-specific decisions followed on 2026-08-31.

The weekly note remains the cross-project head. Do not copy detailed M-* tasks
or Mogent decisions back into it. Update their state in `TODO/TODO.md` and leave
the weekly note as a short status link.

## Current State

```text
v1 libraries/                         preserved comparison baseline
        │
        ▼
37 protected source guides           file-level inventory complete
        │
        ▼
libv2-proto/intake/                   heading-level semantic intake active
        │
        ├─ workflow candidates        tracked, pending owner review
        └─ decision-governance        tracked, pending owner review
        │
        ▼
owner review and explicit choices    active gate
        │
        ▼
orgs/<owner>/                         not promoted yet
        │
        ▼
libv2 cutover                         blocked on completed intake/review
        │
        ▼
root agents.yaml + AGENTS.md          intentionally deferred/clean slate
```

The root manifest and generated guide are absent intentionally. Author the new
root `agents.yaml` only after reviewed libv2 modules exist; then build
`AGENTS.md` from it.

## Include and Composition Narrowing

`docs/thought-experiments/TE-nufad-module-includes.md` is the current analysis.
It begins with nine approaches and narrows them to three:

1. manifest composition plus separate modules as the default;
2. frontmatter plus a minimal named insertion marker if exact mid-document
   placement is required; and
3. a restricted literal Go-template include if one template surface is
   preferred.

Rejected as defaults: rich YAML embedded in HTML comments, Obsidian-only
transclusion, C-preprocessor syntax, and Boolean expressions in include paths.

Owner decisions remain open:

- whether exact mid-document insertion is required for initial libv2;
- which syntax-bearing survivor deserves a prototype if it is required;
- whether GitHub-Flavored Markdown plus YAML frontmatter is canonical;
- where intrinsic, source-spanning, and repository-specific choice rules live.

Do not implement include syntax before those decisions.

## Relationship and Merge Boundary

Do not recombine these two questions:

- **C5a relationship validity:** requirements, conflicts, and selection
  validity. `TE-lusim` analyzes this, but its DR is still unapproved.
- **C5b merge precedence:** how parent, child, or included YAML/template data
  combine. This needs a focused TE after the include direction is chosen.

Module order must not silently resolve a declared hard requirement or exact
conflict.

## Review Batch

Review these first, in order:

1. `libv2-proto/intake/workflow/keep-changes-scoped-and-preserve-existing-work.md`
2. `libv2-proto/intake/workflow/use-risk-to-choose-routine-work-or-decision-review.md`
3. `libv2-proto/intake/workflow/require-human-review-for-high-consequence-changes.md`
4. `libv2-proto/intake/process/decision-governance/use-thought-experiments-to-narrow-a-broad-design-space.md`
5. `libv2-proto/intake/process/decision-governance/lock-durable-decisions-before-implementation.md`
6. `libv2-proto/intake/process/decision-governance/preserve-decision-history-through-supersession.md`
7. `libv2-proto/intake/process/decision-governance/handoff-decisions-with-implementation-evidence.md`

For each candidate, verify source completeness, source/adaptation/proposal
labels, module boundary, overlaps, relationship to other modules, likely output
target, and eventual owner.

## Resume Sequence

1. Read `TODO/TODO.md` for current task states and blockers.
2. Read `docs/LIBV2-PROTO-PROGRESS-REPORT.md` for the shared model.
3. Review the candidate batch with the owner; do not infer approval from its
   committed state or `reliable` editorial label.
4. Update `libv2-proto/PROVENANCE.md` as each source heading receives an
   explicit disposition.
5. Continue semantic intake in order: workflow, process, safety, engineering,
   documentation, tools, languages, communication, domains, repository-local,
   generated candidates.
6. Run focused TEs only when real extracted examples expose a decision.
7. Promote reviewed candidates into `orgs/<owner>/` only after ownership and
   module boundaries are approved.

## Recent Commits

- `5a55b8a Narrow v2 module include design`
- `c701667 Migrate canonical Mogent tasks into repository`
- `6c8575f Document pending v2 review checkpoint`
- `4c6feaf Track pending v2 intake review`
- `7112e90 Add pending-review v2 decision candidates`
- `f19efb2 Add pending-review v2 workflow candidates`
