# NOW

- **Goal:** Dogfood native source browsing and upstream-change visibility without automatic selection rewrites.
- **Current stage:** CS-W024 — Tag search, tag-aware selection, and the consolidated command tree are implemented; dogfood them in Teacher_Bot.
- **One active item:** CS-W024 — dogfood catalog tag search, tag-query selection previews, and canonical command names.
- **Next concrete action:** Review the recorded Teacher_Bot tag-search and tag-query preview output; decide whether the wording/commands need refinement before a write.
- **Blocked by:** none.
- **Do not:** Rebuild or alter Teacher_Bot’s committed generated outputs, resolve overlap automatically, touch Seam's active development, discard the inference stash, or treat historical checkpoints as current.
- **Last verified:** 2026-10-21: full `go test ./...` passed. Installed Teacher_Bot search for `qmr --tag-any topic/workflow --add-commands` showed matching group/leaf nodes and leaf-only add commands; a dry-run/diff `module add cdint:org/core --tag-all review` proposed the expected dynamic include and Review output, with no Teacher_Bot files written.
- **Working-tree/commit state:** Tag-filtering/search/add/tree implementation is committed as `4989165`, `358a7b8`, `205452d`, `fb3cc97`, and leaf-only suggestion refinement `33ae133`. QMR/CDINT retain the intentional source-change simulation; Teacher_Bot has the corresponding authorized generated output/state updates. No push requested.

---
Everything below is history/reference. Resume only from the NOW panel above.

[Work](work.md) · [Decisions](decisions.md) · [Topic archives](archive/README.md)
