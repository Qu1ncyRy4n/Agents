# DR: Library Relationships And Composition Diagnostics

DR-ID: DR-lusim
Date: 2026-08-20
State: open
Asked by: agent-library review

## Question

Which relationship and enforcement model should Mogent adopt for reusable
library modules?

## Recommendation

Approve the hybrid in `TE-lusim`:

1. reserve `self:` for same-source references in library metadata;
2. keep cross-source relationships in the consuming manifest;
3. add `see_also`, source-scoped `alternative_family`, and informational
   `provenance` metadata;
4. make recognized `self:` requirements/conflicts and multiple selected family
   members strict after migrating the checked-in libraries;
5. keep legacy alias-bound strings advisory during that migration;
6. preserve heading paths with explicit IDs and defer source-path-changing
   directory moves to a separate source-evolution decision.

Because metadata is currently file-level, the first alternative-family
dogfood would split persona children into atomic files while preserving their
existing source paths. It would not require heading-level metadata.

## Owner Decisions Requested

1. Approve `self:` as the reserved same-source reference syntax, or request a
   different spelling.
2. Approve `see_also`, `alternative_family`, and `provenance` as the next
   metadata fields.
3. Approve errors for recognized missing requirements, exact conflicts, and
   multiple alternatives after migration, while keeping soft links and legacy
   strings non-blocking.
4. Approve deferring source-path-changing directory moves until Mogent has a
   source-path migration contract.

## Why This Blocks Implementation

These choices change the public source metadata format and whether existing
manifests can build. Implementing validation first would either bind libraries
to consumer aliases or turn previously informational strings into unannounced
errors.

## Affects

- `library.Metadata` and source validation
- public workspace diagnostic types
- `source show`, `status`, `add`, build validation, and future TUI previews
- checked-in library frontmatter
- future manifest-owned cross-source relationship syntax
- library taxonomy moves and source-path compatibility

## Unblocks

- hard requirement and exact conflict diagnostics;
- mutually exclusive persona selection;
- soft module discovery links;
- migration of the nine current alias-bound `requires` values;
- the next bounded taxonomy and deduplication work.

## Related Work

- `docs/thought-experiments/TE-lusim-library-relationships.md`
- `docs/thought-experiments/TE-kavam-metadata-tags-and-library-shape.md`
- `docs/LIBRARY-REVIEW-PLAN.md`
- `docs/LIBRARY-COVERAGE-AUDIT.md`
