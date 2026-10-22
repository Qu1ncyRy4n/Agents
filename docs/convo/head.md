# NOW

- **Goal:** Dogfood native source browsing and upstream-change visibility without automatic selection rewrites.
- **Current stage:** CS-W024 — native browsing and upstream module reporting implemented; dogfood a fresh recorded source baseline in Teacher_Bot.
- **One active item:** CS-W024 — observe source-change reporting against the QMR/CDINT simulation after one intentional baseline build.
- **Next concrete action:** In Teacher_Bot, run `mogent build`, change one QMR/CDINT module, then inspect `mogent drift --all --diff` for changed/added/removed paths.
- **Blocked by:** Owner must choose when to write Teacher_Bot's currently dirty generated outputs/state.
- **Do not:** Rebuild or alter Teacher_Bot’s committed generated outputs, resolve overlap automatically, touch Seam's active development, discard the inference stash, or treat historical checkpoints as current.
- **Last verified:** 2026-10-21: focused `go test ./assembly ./state ./internal/cli` passed for upstream source baselines, recursive changed/added/removed reporting, and optimize preview/apply. Installed Teacher_Bot verification remains pending because its state must first be refreshed by a normal build.
- **Working-tree/commit state:** Agents implementation is pending an atomic commit. QMR and CDINT demo libraries have pre-existing plus isolated simulation changes; Teacher_Bot has owner-authored dirty generated outputs/state and remains out of scope for automatic edits. No push requested.

---
Everything below is history/reference. Resume only from the NOW panel above.

[Work](work.md) · [Decisions](decisions.md) · [Topic archives](archive/README.md)
