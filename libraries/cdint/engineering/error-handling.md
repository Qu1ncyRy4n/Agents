---
tldr: Handle errors explicitly and preserve useful causal context.
---
# Error Handling

## Shell Commands

Never hide a command failure with `|| true`. Use explicit branching when a
non-fatal command may fail, and record the exit status or relevant diagnostics
instead of dropping the result. `|| true` is called out because it is shell
syntax that converts any failure into an apparently successful command; the
general rule applies to every language.

For example, if a probe is optional, branch on it and say what was unavailable.
If the main build fails, return that failure rather than continuing as though
the build succeeded.

## Cleanup And Diagnostics

For cleanup, probing, or optional diagnostics, make the success/failure handling
visible. A cleanup failure may be non-fatal, but it should not be silent when it
affects user data, generated output, or test confidence.

For example, a successful install followed by failure to remove its staging
directory should report both facts: the install completed, but temporary data
remains at a named path. If the primary operation also failed, preserve that
error and attach the cleanup failure rather than replacing either one.

## Go

In Go code, do not ignore errors with `_ = ...`. Handle, propagate, or report
the error. Use `%w` when adding context that callers can unwrap. Keep
`errcheck ./...` passing when Go behavior changes and the repo expects it.

For example, `fmt.Errorf("load manifest: %w", err)` explains the failed
operation while preserving `err` for `errors.Is` or `errors.As`. `errcheck`
finds returned Go errors that code discarded without handling.

## Tests

Assert expected failures explicitly. Avoid tests that pass because setup failed
or because an external dependency was unavailable unless the test is deliberately
checking that skip path. A parser test is not evidence about invalid input if
its fixture failed to load and the test accepted any error; verify setup first,
then assert the parser's specific failure.
