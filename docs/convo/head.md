# NOW

- **Goal:** Dogfood native source browsing and upstream-change visibility without automatic selection rewrites.
- **Current stage:** CS-W024 — source-oriented browsing accepted; inspect the changed/added/removed upstream cases before choosing provenance behavior.
- **One active item:** CS-W024 — observe and decide how native status should explain upstream library changes.
- **Next concrete action:** Review the Teacher_Bot `status`/`drift --all --diff` result from the isolated QMR/CDINT simulation; decide whether to design an authored input baseline.
- **Blocked by:** Owner decision on the desired upstream-change explanation and baseline semantics.
- **Do not:** Rebuild or alter Teacher_Bot’s committed generated outputs, resolve overlap automatically, touch Seam's active development, discard the inference stash, or treat historical checkpoints as current.
- **Last verified:** 2026-10-21: focused `go test ./internal/cli ./assembly` passed; installed CLI verified `source show`, `show`, `status`, `drift --all --diff`, and `optimize` in Teacher_Bot. `build --check` correctly failed after the upstream QMR simulation.
- **Working-tree/commit state:** Agents native command implementation is committed; this documentation checkpoint records its active follow-up. QMR and CDINT demo libraries have pre-existing plus isolated simulation changes; Teacher_Bot main at 878de7c remains clean. No push requested.

---
Everything below is history/reference. Resume only from the NOW panel above.

[Work](work.md) · [Decisions](decisions.md) · [Topic archives](archive/README.md)
