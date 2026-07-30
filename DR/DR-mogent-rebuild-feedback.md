# DR: Mogent Rebuild Feedback And Deferred Questions

DR-ID: DR-mogent-rebuild-feedback
Date: 2026-07-29
State: open

## Purpose

This is the running feedback section for questions intentionally deferred while
the milestone-one renderer is implemented. None of these block the local-path
parse, resolve, and render path selected in DI-vukam.

## Manifest And Rendering

1. Should a later formatter choose compact or explicit syntax as its canonical
   emitted representation, given that both authoring forms are accepted?
2. Should an entry eventually support both `from` and `children`? A concrete
   example is needed to define whether pulled source content comes before,
   after, or between explicit children.
3. For composed sources with overlapping descendant paths, what exact
   confirmation interaction should choose combined versus separate sections?
4. Should `exclude` identify nodes with source-qualified paths only, or gain a
   concise relative-path form once its ambiguity rules are specified?
5. Are source root bodies without descendant headings a valid long-term library
   style, and should empty source nodes be errors or simply render nothing?

## Sources And Local Overrides

6. What cache location, pin/lockfile format, and changed-source review behavior
   should be used when URL sources are added after milestone one?
7. After dogfooding, should `cdint`, `ucd_research`, `personal`, and `nix`
   remain the initial source boundaries (see DR-garom)?
8. What exact local-library path and copy-on-write file layout should `edit`
   create under `.mogent/`?
9. Should source paths be restricted to paths within the consuming repository,
   or may a manifest intentionally point to an arbitrary readable local library?

## CLI And Product Flow

10. Is `mogent build` sufficient as the first public command, or should
    `validate`/`render --stdout` be introduced before the navigator?
11. What should `init` generate when no library is available yet: an empty
    skeleton, a local identity template, or guided source discovery?
12. Which manifest diagnostics and provenance views are necessary before the
    read-only navigator milestone begins?

## Compatibility And Adoption

13. When should root `AGENTS.toml` be replaced by a real `agents.yaml` dogfood
    manifest, and should the previous file be retained temporarily as a sample?
14. How should manually edited generated `AGENTS.md` be reported before formal
    drift detection exists?
