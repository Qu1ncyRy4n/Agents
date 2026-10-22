# NOW

- **Goal:** Dogfood native source browsing and upstream-change visibility without automatic selection rewrites.
- **Current stage:** CS-W024 — Source/tag discovery is dogfooded; the next presentation refinement is a tree-first text-search result with optional flat rows and content context.
- **One active item:** CS-W024 — Resume with the tree/snippet search presentation, then dogfood it before a real selection write.
- **Next concrete action:** On resume, define and implement `source search <term>` tree output, `--flat`, and content/context display; use `source show <source> --not-included` to choose unselected source content. Keep `auto_inherit` as a separately logged selection-policy slice.
- **Blocked by:** none.
- **Do not:** Rebuild or alter Teacher_Bot’s committed generated outputs, resolve overlap automatically, implement `auto_inherit` before its membership/audit persistence is designed, touch Seam's active development, discard the inference stash, or treat historical checkpoints as current.
- **Last verified:** 2026-10-21: full `go test ./...` passed before `d2c531d`; its focused CLI tests also passed. Installed Seam `source show qmr --not-included` hid selected QMR skills and showed the remaining QMR guidance tree; no Seam file was written. Earlier installed Teacher_Bot checks confirmed source/tree/tag/text-search behavior.
- **Working-tree/commit state:** Tag-filtering/search/add/tree implementation is committed through `c1e3f04`; the proposed `auto_inherit` policy is recorded but unimplemented. QMR/CDINT retain the intentional source-change simulation; Teacher_Bot has the corresponding authorized generated output/state updates. No push requested.

---
Everything below is history/reference. Resume only from the NOW panel above.

[Work](work.md) · [Decisions](decisions.md) · [Topic archives](archive/README.md)
