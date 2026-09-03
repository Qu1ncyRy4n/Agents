# AGENTS.md - photo_getter

## Project

Family archive recovery and organization pipeline. Rust `photosync` handles scan/hash/dedupe/stage/reconcile/report operations. Python `photocluster` handles embeddings, NSFW scoring, near-duplicate grouping, and the local review UI.

## Hard Rules

- Never delete source/archive files or add physical deletion behavior without explicit user approval.
- Curation is mark-only: `files.purge_marked=1` plus `files.purge_reason`.
- SQLite is single-writer. Do not run `embed.py`, `classify.py`, `neardup.py`, or `serve.py` as concurrent writers.
- Keep `photosync.db` on the internal disk. It is an index, not authoritative data.
- Treat mounted drive paths as significant; remounting may change `/Volumes/...` names and break resume assumptions.

## Structure

- `src/main.rs` - Rust CLI commands.
- `src/lib.rs` - scan, hash, dedupe, stage/reconcile, status/report, verify-covered.
- `tests/pipeline.rs` - Rust integration tests.
- `photocluster/embed.py` - CLIP embedding writer.
- `photocluster/classify.py` - NSFW scoring/quarantine marking.
- `photocluster/neardup.py` - near-duplicate grouping.
- `photocluster/serve.py` - local curation UI.
- `video_health.py` - video probe/repair/livecheck/codeccheck/transcode support.
- `SPEC.md` and `HANDOFF-2026-07-22.md` - operational context and current state.

## Commands

```bash
cargo test
cargo check
cargo run -- status
cargo run -- report
```

Python commands should usually use the existing venv directly:

```bash
photocluster/.venv/bin/python photocluster/embed.py --db photosync.db
photocluster/.venv/bin/python photocluster/classify.py score --db photosync.db
photocluster/.venv/bin/python photocluster/neardup.py --db photosync.db --threshold 0.93
photocluster/.venv/bin/python photocluster/serve.py --db photosync.db
```

Run those Python writers serially.

## Resume Context

As of `HANDOFF-2026-07-22.md`, embedding coverage was the blocker for photo curation. Re-run `embed.py` only after confirming the expected drive mount path exists.

## Validation

- Rust changes: run `cargo test`.
- DB-affecting changes: test on a temporary DB or fixture first.
- Review or purge-marking changes: inspect SQL effects before running on `photosync.db`.
- Pre-wipe coverage gates must remain strict: `verify-covered` should fail if any in-scope source file is not represented outside the source root.
