# NOW

- **Goal:** Dogfood native source browsing and upstream-change visibility without automatic selection rewrites.
- **Current stage:** CS-W024 — source-oriented browsing accepted; inspect the changed/added/removed upstream cases before choosing provenance behavior.
- **One active item:** CS-W024 — observe and decide how native status should explain upstream library changes.
- **Next concrete action:** Specify the deterministic per-node and library-root content-hash payloads, then implement comparison/descent against the consumer's recorded baseline.
- **Blocked by:** Exact canonical bytes/metadata and physical-tree coverage for each node type.
- **Do not:** Rebuild or alter Teacher_Bot’s committed generated outputs, resolve overlap automatically, touch Seam's active development, discard the inference stash, or treat historical checkpoints as current.
- **Last verified:** 2026-10-21: focused `go test ./internal/cli ./assembly` passed; installed CLI verified `source show`, `show`, `status`, `drift --all --diff`, and `optimize` in Teacher_Bot. `build --check` correctly failed after the upstream QMR simulation. Owner requires module-level upstream-change detail.
- **Working-tree/commit state:** Agents native command implementation is committed; this documentation checkpoint records its active follow-up. QMR and CDINT demo libraries have pre-existing plus isolated simulation changes; Teacher_Bot main at 878de7c remains clean. No push requested.

---
Everything below is history/reference. Resume only from the NOW panel above.

[Work](work.md) · [Decisions](decisions.md) · [Topic archives](archive/README.md)
