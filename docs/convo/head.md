# NOW

- **Goal:** Dogfood native source browsing and upstream-change visibility without automatic selection rewrites.
- **Current stage:** CS-W024 — Source/tag discovery is dogfooded; resolved catalog filters now cover selected and available content. Canonical module vocabulary and a compact plan-file guide now match the implementation. The next presentation refinement is tree-first text search with optional flat rows and content context.
- **One active item:** CS-W024 — Resume tree/snippet search presentation, then dogfood it before a real selection write.
- **Next concrete action:** Define and implement `source search <term>` tree output, `--flat`, and content/context display. Keep `auto_inherit` as a separately logged selection-policy slice.
- **Blocked by:** none.
- **Do not:** Rebuild or alter Teacher_Bot’s committed generated outputs, resolve overlap automatically, implement `auto_inherit` before its membership/audit persistence is designed, touch Seam's active development, discard the inference stash, or treat historical checkpoints as current.
- **Last verified:** 2026-10-23: `go test ./internal/cli` and full `go test ./...` passed for canonical `module include`/`module remove`, help/reference output, and plan-guide changes. 2026-10-22 installed Teacher_Bot `source show qmr --included` rendered selected branches and `--included --flat` listed exact resolved QMR paths; no Teacher_Bot file was written. Earlier installed Seam/Teacher_Bot checks confirmed the complementary unselected, source/tree, tag, and text-search behavior.
- **Working-tree/commit state:** Tag-filtering/search/add/tree implementation is committed through `c1e3f04`; the proposed `auto_inherit` policy is recorded but unimplemented. The current verified slice implements canonical `module include`/`module remove` aliases, updates help/reference output, and adds `docs/PLAN-FILES.md`; it awaits its atomic commit. QMR/CDINT retain the intentional source-change simulation; Teacher_Bot has the corresponding authorized generated output/state updates. No push requested.

---
Everything below is history/reference. Resume only from the NOW panel above.

[Work](work.md) · [Decisions](decisions.md) · [Topic archives](archive/README.md)
