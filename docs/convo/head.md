# NOW

- **Goal:** Dogfood native source browsing and upstream-change visibility without automatic selection rewrites.
- **Current stage:** CS-W024 — Teacher_Bot upstream source reporting is live; refine command names and choose the next tag-discovery dogfood slice.
- **One active item:** CS-W024 — settle the catalog/result/localization command tree and tag-search contract.
- **Next concrete action:** Review read-only source tag search and explicit add-command rendering before implementing a tag CLI.
- **Blocked by:** none.
- **Do not:** Rebuild or alter Teacher_Bot’s committed generated outputs, resolve overlap automatically, touch Seam's active development, discard the inference stash, or treat historical checkpoints as current.
- **Last verified:** 2026-10-21: Teacher_Bot's installed CLI reported CDINT's removed `org/optional/project-defaults` and QMR's changed `intro/assistant-agent` plus added `intro/drift-check`; `drift --diff` exposed one overly broad output hunk. Focused `go test ./textdiff ./internal/cli` passed after compact three-line-context hunk rendering. Resolver tests verify `tags_all`/`tags_any` selection plus scoped path exclusion; no tag-search CLI exists.
- **Working-tree/commit state:** Compact diff rendering (`b0ea529`), focused localization help (`18d3731`), and corrected native-doc guidance (`0f9b752`) are committed. QMR/CDINT retain the intentional source-change simulation; Teacher_Bot has the corresponding authorized generated output/state updates. No push requested.

---
Everything below is history/reference. Resume only from the NOW panel above.

[Work](work.md) · [Decisions](decisions.md) · [Topic archives](archive/README.md)
