# NOW

- **Goal:** Dogfood native source browsing and upstream-change visibility without automatic selection rewrites.
- **Current stage:** CS-W024 — Teacher_Bot upstream source reporting is live; the next dogfood slice is read-only tag discovery.
- **One active item:** CS-W024 — specify tag search, `tags_none`, and explicit add-command rendering before a tag CLI.
- **Next concrete action:** Define expected QMR/CDINT source-search results and generated add commands, then implement that read-only surface.
- **Blocked by:** none.
- **Do not:** Rebuild or alter Teacher_Bot’s committed generated outputs, resolve overlap automatically, touch Seam's active development, discard the inference stash, or treat historical checkpoints as current.
- **Last verified:** 2026-10-21: Teacher_Bot's installed CLI reported CDINT's removed `org/optional/project-defaults` and QMR's changed `intro/assistant-agent` plus added `intro/drift-check`; `drift --diff` exposed one overly broad output hunk. Focused `go test ./textdiff ./internal/cli` passed after compact three-line-context hunk rendering. Resolver tests verify `tags_all`/`tags_any` selection plus scoped path exclusion; no tag-search CLI exists. `tags_none` is the recommended negative predicate pending read-only dogfood.
- **Working-tree/commit state:** Compact diff rendering (`b0ea529`), focused localization help (`18d3731`), and corrected native-doc guidance (`0f9b752`) are committed. QMR/CDINT retain the intentional source-change simulation; Teacher_Bot has the corresponding authorized generated output/state updates. No push requested.

---
Everything below is history/reference. Resume only from the NOW panel above.

[Work](work.md) · [Decisions](decisions.md) · [Topic archives](archive/README.md)
