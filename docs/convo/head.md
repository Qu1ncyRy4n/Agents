# NOW

- **Goal:** Dogfood native source browsing and upstream-change visibility without automatic selection rewrites.
- **Current stage:** CS-W024 — Teacher_Bot upstream source reporting is live; refine the next native workflow from the observed report.
- **One active item:** CS-W024 — review the source-change report and choose the next native feature to dogfood.
- **Next concrete action:** In Teacher_Bot, run `mogent build`, change one QMR/CDINT module, then inspect `mogent drift --all --diff` for changed/added/removed paths.
- **Blocked by:** none.
- **Do not:** Rebuild or alter Teacher_Bot’s committed generated outputs, resolve overlap automatically, touch Seam's active development, discard the inference stash, or treat historical checkpoints as current.
- **Last verified:** 2026-10-21: focused `go test ./assembly ./state ./internal/cli` passed for upstream source baselines, recursive changed/added/removed reporting, and optimize preview/apply. Installed Teacher_Bot verification remains pending because its state must first be refreshed by a normal build.
- **Working-tree/commit state:** Native-only cleanup and source-report refinement are committed. QMR/CDINT retain the intentional source-change simulation; Teacher_Bot has the corresponding authorized generated output/state updates. No push requested.

---
Everything below is history/reference. Resume only from the NOW panel above.

[Work](work.md) · [Decisions](decisions.md) · [Topic archives](archive/README.md)
