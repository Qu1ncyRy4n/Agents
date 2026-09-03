# Repository Guidelines

## Project Structure & Module Organization
- Root Go module lives at the repo root (`go.mod`) with core sources like `interfaces.go` and tests like `interfaces_test.go`.
- `x/` contains experimental prototypes and notes; many subfolders are standalone Go modules (`x/<name>/go.mod`).
- Documentation and exploratory writeups are primarily in `README.md` and `x/**/*.md`.
- `llms-full.txt` is a dense, llms.txt-style PromiseGrid description intended for LLM ingestion (includes links to repo docs plus local-only references).
- Local Grokker state (e.g., `.grok`) and generated binaries should stay uncommitted.

## Build, Test, and Development Commands
- `go test ./...` runs the root module test suite.
- `go test -v ./...` runs tests verbosely (matches the default in `devloop.sh`).
- `./devloop.sh` watches for file changes, runs tests, and emits audio cues for pass/fail.
- For experiments, run tests within the module directory, e.g. `cd x/ipld-path && go test ./...`.

## Coding Style & Naming Conventions
- Go code follows standard `gofmt` formatting (tabs, canonical import grouping).
- Package names are short and lower-case.
- Keep edits small and focused; do not remove comments or docs—update them when they drift.
- Use 'git mv' for file renames to preserve history.

## Testing Guidelines
- Use Go’s standard `testing` package with `*_test.go` filenames.
- Prefer deterministic tests; avoid network calls unless explicitly required.
- Table-driven tests are encouraged when validating multiple cases.

## Task Tracking (TODO)
- Track tasks and plans in `TODO/`, with a `TODO/TODO.md` index of small tasks.
- Number TODOs with zero-padded IDs (e.g., `007`), don’t renumber, and mark completed items by checking them off (e.g., `- [x] 007 ...`).

## Commit & Pull Request Guidelines
- Recent history uses short, imperative one-line subjects (often lowercase); use concise, imperative, capitalized subjects going forward.
- Commit bodies include a section per changed file with bullet summaries.
- PRs should include a concise summary, relevant test commands run, and linked issues; include before/after notes when behavior or output changes.

## References

Known sources of PromiseGrid information and code used during this work.

### `grid-poc` (this repo)
- `AGENTS.md`
- `README.md`
- `llms-full.txt`
- `x/rfc/draft-promisegrid.md`
- `x/wire/wire.md`
- `x/sim3/message-format.md`
- `x/sim3/design.md`
- `x/references.md`
- Root module code: `go.mod`, `interfaces.go`, `interfaces_test.go`
- Experimental modules and prototypes: `x/**/go.mod`, `x/**/*.md`, `x/**/*.go`

### `~/lab/promisegrid/promisegrid`
- `~/lab/promisegrid/promisegrid/README.md`
- `~/lab/promisegrid/promisegrid/slides/intro/README.md`

### `~/lab/promisegrid/paper-ism`
- `~/lab/promisegrid/paper-ism/README.adoc`
- `~/lab/promisegrid/paper-ism/local/PromiseGrid, PT, and SST.md`

### `~/lab/grid-cli`
- `~/lab/grid-cli/README.md`
- `~/lab/grid-cli/v2/doc/330-messages.md`
- `~/lab/grid-cli/v2/doc/150-capability-tokens.md`
- `~/lab/grid-cli/v2/doc/420-routing.md`
- Versioned code modules: `~/lab/grid-cli/v2/go.mod`, `~/lab/grid-cli/v1-grid/go.mod`, `~/lab/grid-cli/v0-multibuild/go.mod`

### Other local notes and integration docs
- `~/lab/x/promisegrid/README.md`
- `~/lab/cswg/workshop-2025-03-11-ipfs-and-promisegrid/slides.md`
- `~/lab/collab-editor/docs/promisegrid-collab-editor.md`
- `~/lab/collab-editor/docs/promisegrid-integration.md`
- `~/lab/angela/promisegrid-angela.md`

### Background / theory
- `~/lab/mark_burgess/pt-book/BookOfPromises.pdf`
- `~/lab/mark_burgess/MoneyBook.pdf`
- `~/core/language/README.adoc`
- `~/core/language/requirements.md`

### Missing docs (restic)
- `TODO/001-fetch-missing-design-docs.md` (restore `15000.md` and `13130.md`)
