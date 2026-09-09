# Steve Library Curation Handoff

## Status And Boundary

Assessment: after stable curated input, multi-file targets and rendering are likely
a bounded one-to-two-day engineering task. Splitting, classifying, and safely
adapting the intertwined cdint-grid corpus is unbounded without owner direction.
This is an assessment, not an estimate guarantee.

The current cdint-grid `AGENTS.md` is coupled to PromiseGrid, Go, repository
layouts, named tools, and personal or machine facts. An opinionated rewrite
would risk destructive or inaccurate inference. The existing
[`library-curation/drafts/`](library-curation/drafts/) output is draft evidence
only, not a promoted canonical library. During this pause, do not promote or
render newly inferred content.

Start Here: read this handoff, then [`library-curation/README.md`](library-curation/README.md),
[`library-curation/open-questions.md`](library-curation/open-questions.md), and
[`wb.md`](wb.md), in that order; pause before rendering or promoting drafts.

The exact user decisions are recorded in [the wb review space](wb.md): retain the
TE protocol as organization-canonical and generally reusable; TODO, DR, DI,
GLOSSARY, and other named repository documentation or task locations are
permitted conventions, but their placement and scope require explicit curation;
and replace the absolute minter path only with a documented portable external
tool, Mogent implementation, or dependency reference.

## Requested Curation Input

Before content adaptation, please provide:

1. One canonical curated library or source document for dogfooding.
2. Provisional per-section scope labels: `universal`, `organization-universal`,
   `project-specific`, or `person/machine-specific`.
3. A per-section placement and default-selection decision.
4. An explicit record of permitted source changes or adapters.
5. Ownership and authority for resolving conflicts and changes.
6. A pinned source revision or explicit curation window, so a changing corpus
   cannot silently change the dogfood input mid-review.
7. A dogfood acceptance and promotion gate: what evidence permits a curated
   draft to become a default or organization policy, and who approves it.

Source provenance alone does not assign scope.

## Proposed Organization (Requires Steve Approval)

This proposal separates scope from interface:

- `agents/`: always-loaded, language/tool-neutral rules with no private or
  project-local references.
- `skills/`: triggered procedures.
- `guides/`: explanatory or selectable operational guidance.
- `specs/`: normative organization or project protocols that are not
  automatically rendered.
- `overlays/`: optional project, language/tool, or person/machine content.
- `reference/`: non-rendered evidence.

Language, domain, and scope should be metadata or source-selection axes, not
content-category assumptions. Existing design evidence: [`DESIGN.md:196-224`](DESIGN.md#5-categories)
and [`TE-kavam-metadata-tags-and-library-shape.md:47-101`](thought-experiments/TE-kavam-metadata-tags-and-library-shape.md#option-b-hierarchical-tags).

## Proposed Sequence (Requires Steve Approval)

1. Curate verbatim source.
2. Annotate scope, interface, and default selection.
3. Select a small bounded dogfood set.
4. Render a multi-file target.
5. Dogfood real projects.
6. Only then extract or adapt based on evidence.

Do not make automatic content reductions.

## Questions For Steve

Answer the questions in [`library-curation/open-questions.md`](library-curation/open-questions.md), including:

- Go and Git placement.
- Glossary location and adapter.
- Exact `AGENTS.md` filename versus adapter.
- Which record conventions are organization versus project conventions.
- Canonical curation ownership.
- Whether Mogent manages local/project overlays or merely renders them.

## Deferred But Not Lost

Code-quality work is deferred by the curation and dogfood block. Reconcile the
durable Mogent task authority in [`TODO/TODO.md`](../TODO/TODO.md) before
resuming product work.

## Navigation

- [wb review space](wb.md)
- [library curation landing page](library-curation/README.md)
- [questions](library-curation/open-questions.md)
- [drafts](library-curation/drafts/)
- [DESIGN](DESIGN.md)
- [LIBRARY-REVIEW-PLAN](LIBRARY-REVIEW-PLAN.md)
- [LIBRARY-COVERAGE-AUDIT](LIBRARY-COVERAGE-AUDIT.md)
