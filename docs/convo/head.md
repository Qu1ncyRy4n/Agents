# NOW

- **Goal:** Dogfood native source browsing and upstream-change visibility without automatic selection rewrites.
- **Current stage:** CS-W024 — Tag search, tag-aware selection, and the consolidated command tree are implemented; dogfood their dynamic-versus-explicit alternatives in Teacher_Bot.
- **One active item:** CS-W024 — dogfood catalog nesting plus dynamic tag-query config and explicit leaf-command previews.
- **Next concrete action:** Run the read-only Teacher_Bot previews, then review the two forms for clarity before any write.
- **Blocked by:** none.
- **Do not:** Rebuild or alter Teacher_Bot’s committed generated outputs, resolve overlap automatically, touch Seam's active development, discard the inference stash, or treat historical checkpoints as current.
- **Last verified:** 2026-10-21: focused `go test ./internal/cli ./assembly` passed after adding conservative physical-directory nesting, `source search --add-config`, and `module add --explicit-commands`. Teacher_Bot read-only verification remains next.
- **Working-tree/commit state:** Tag-filtering/search/add/tree implementation is committed as `4989165`, `358a7b8`, `205452d`, `fb3cc97`, and leaf-only suggestion refinement `33ae133`. QMR/CDINT retain the intentional source-change simulation; Teacher_Bot has the corresponding authorized generated output/state updates. No push requested.

---
Everything below is history/reference. Resume only from the NOW panel above.

[Work](work.md) · [Decisions](decisions.md) · [Topic archives](archive/README.md)
