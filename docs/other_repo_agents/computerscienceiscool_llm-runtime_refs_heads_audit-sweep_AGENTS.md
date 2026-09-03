# Repository Guidelines

## Project Structure & Module Organization
- `cmd/llm-runtime/` is the CLI entry; `make build` emits `./llm-runtime`.
- `pkg/` holds public packages (`app`, `cli`, `config`, `evaluator`, `sandbox`, `scanner`, `search`, `session`); `internal/core/` is private.
- `configs/` and root `llm-runtime*.yaml` are sample configs; `embeddings.db` is generated.
- Supporting assets: `docs/` (incl. `docs/TODO.md`), `scripts/`, `tests/integration/`, `Dockerfile.io`, and the `Makefile`.

## Build, Test, and Development Commands
- Build: `make build` or `go build ./cmd/llm-runtime`; interactive run via `make run`.
- Tests: `make test`; coverage `make test-coverage`; benchmarks `make bench`; formatting/vet `make quality`.
- Containers: `make build-io-image`, `make check-io-image`, `make test-io-container`; sweep `make test-suite`; release tarball `make release`.

## Coding Style & Naming Conventions
- Go 1.21+; run `gofmt` and `go vet`.
- Package names stay short and lower-case; configs use lower-kebab `*.yaml`.
- Keep generated artifacts (binaries, `embeddings.db`, audit logs, `.grok`/container state) out of commits; favor table-driven `*_test.go` cases.

## Testing Guidelines
- Default: `go test ./...` or `make test`; co-locate cases with the code they cover.
- Use `make test-coverage` when adding behavior; hit new branches and error paths.
- Integration lives in `tests/integration/`; docker-heavy flows should run `make check-docker` or `make test-io-container`. Keep tests deterministic/offline; mock external calls.

## TODO Tracking
- Track work in `TODO/TODO.md` (create `TODO/` if missing); number items with zero-padded IDs (001, 002, …) without renumbering.
- Keep active items under High/Medium/Low sections, sorted by priority; reference IDs in commits/PRs; mark completion with checkboxes.
- When done, move the line to the `DONE` section at the bottom (checked) instead of deleting it. `docs/TODO.md` is legacy; prefer `TODO/TODO.md`.

## Commit & Pull Request Guidelines
- Commits are short, imperative, and capitalized; add per-file bullets in the body when several areas change.
- Stage files explicitly (e.g., `git add AGENTS.md cmd/llm-runtime/main.go pkg/config/...`) instead of `git add .`.
- Run `make test` (and `make build` for CLI changes) before committing; PRs list executed commands, linked issues, configs touched, and before/after output for user-visible changes. A message containing only `commit` means: stage and commit all changes with an AGENTS.md-compliant message based on `git diff`.
- `make commit` generates a message via `grok commit` and pushes; review before use.

## Security & Local State
- Docker is required for exec and containerized I/O; pull `python-go` and run `make build-io-image` before exec workflows; keep whitelists/limits in `llm-runtime.yaml` / `llm-runtime.config.yaml`.
- Do not commit secrets or generated state (audit logs, binaries, embeddings, `.grok`, container outputs); extend excludes if new sensitive paths appear.
