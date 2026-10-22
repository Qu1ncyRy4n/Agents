# NOW

- **Goal:** Dogfood native source browsing and upstream-change visibility without automatic selection rewrites.
- **Current stage:** CS-W024 — Tag search, tag-aware selection, and the consolidated command tree are implemented; dogfood them in Teacher_Bot.
- **One active item:** CS-W024 — dogfood catalog tag search, tag-query selection previews, and canonical command names.
- **Next concrete action:** Install Mogent, run the three tag/command-tree checks in Teacher_Bot, and report the observed output/wording.
- **Blocked by:** none.
- **Do not:** Rebuild or alter Teacher_Bot’s committed generated outputs, resolve overlap automatically, touch Seam's active development, discard the inference stash, or treat historical checkpoints as current.
- **Last verified:** 2026-10-21: focused syntax/assembly/CLI tests pass for `tags_none`, source tag search with explicit add commands, tag-aware `module add`, and canonical `source tree`, `show tree`, and module subcommands. Teacher_Bot's installed CLI still needs this new slice dogfooded.
- **Working-tree/commit state:** Tag-filtering/search/add/tree implementation is committed as `4989165`, `358a7b8`, `205452d`, and `fb3cc97`. QMR/CDINT retain the intentional source-change simulation; Teacher_Bot has the corresponding authorized generated output/state updates. No push requested.

---
Everything below is history/reference. Resume only from the NOW panel above.

[Work](work.md) · [Decisions](decisions.md) · [Topic archives](archive/README.md)
