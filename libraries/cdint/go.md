# Go

## Development

Run `gofmt` on changed Go code. Run focused `go test` from the affected module.
Run `errcheck ./...` for Go behavior changes. If a required command is missing
or the environment is broken, report the blocker instead of changing unrelated
files.

## Code

Prefer small, readable packages and explicit dependencies. Keep package names
short and lower-case. Handle errors explicitly and use `%w` when adding useful
context. Do not ignore errors or hide shell-command failures with `|| true`.

## Tests

Use the standard `testing` package. Keep tests deterministic. Prefer fixtures
and table-driven tests when the same behavior has several cases. Add coverage
for new behavior and error paths close to the code they exercise.
