---
tldr: Apply focused Go formatting, testing, error-handling, and module hygiene.
---
# Go

## Development

Run `gofmt` on changed Go code. Run focused `go test` from the affected module.
Run `go vet` when the repository uses it directly or through `make quality`.
Run `errcheck ./...` for Go behavior changes. If a required command is missing
or the environment is broken, report the blocker instead of changing unrelated
files.

When a repository has several Go modules, run commands from the module that owns
the changed package. Common layouts include a root module, `cmd/<tool>` CLI
entry points, and experimental modules under `x/**/go.mod`.

For CLI tools, validate at least the help path or the smallest affected workflow
after behavior changes, for example `go run ./cmd/<tool> -h` or the repo's
documented smoke command.

## Code

Prefer small, readable packages and explicit dependencies. Keep package names
short and lower-case. Handle errors explicitly and use `%w` when adding useful
context. Do not ignore errors or hide shell-command failures with `|| true`.

Prefer the standard library and well-known dependencies. Add a dependency only
when it removes real complexity or matches the repository's established stack.
Avoid large functions, unnecessary global state, and package-level mutable
state. Extract helpers when they remove meaningful duplication, not merely to
name one short expression.

Keep import grouping canonical and let `gofmt` own whitespace. Keep generated
artifacts, local state, coverage output, databases, and compiled binaries out of
commits.

## Tests

Use the standard `testing` package. Keep tests deterministic. Prefer fixtures
and table-driven tests when the same behavior has several cases. Add coverage
for new behavior and error paths close to the code they exercise.

Co-locate tests with the code they cover unless the repository already uses a
separate integration-test tree. Mock external calls by default. Avoid network
tests unless the task explicitly needs them and the repository documents the
expected external dependency.

When adding behavior, cover the success path, important error paths, and at
least one boundary case. Prefer asserting whole structured values when that
keeps the test clearer than checking fields one by one.

## Repository Shape

Common Go repository shapes in the corpus:

- root module with core files and `*_test.go` beside them;
- CLI entry under `cmd/<name>/`;
- public packages under purpose-named directories or `pkg/`;
- private packages under `internal/` when the repository already uses that
  boundary;
- experimental prototypes under `x/`, often with their own `go.mod`.

Follow the local shape. Do not introduce `internal/`, `pkg/`, `cmd/`, or `x/`
as a style preference when the repository's instructions choose another layout.

## Corpus Variants

Several captured repos require packages at the module root or under
purpose-named directories and explicitly avoid `internal/` or `pkg/`. Others use
`pkg/` for public packages and `internal/` for private core packages. Treat this
as a repo-local architecture decision, not a universal Go rule.

Some repos wrap Go checks in `make` targets such as `make test`,
`make test-coverage`, `make quality`, or `make build`. Use those wrappers when
they encode repo policy; use raw `go` commands when the repo documents raw Go
commands.
