---
tldr: Handle errors explicitly and preserve useful causal context.
---
# Error Handling

## Shell Commands

Never hide a command failure with `|| true`. Use explicit branching when a
non-fatal command may fail, and record the exit status or relevant diagnostics
instead of dropping the result.

## Cleanup And Diagnostics

For cleanup, probing, or optional diagnostics, make the success/failure handling
visible. A cleanup failure may be non-fatal, but it should not be silent when it
affects user data, generated output, or test confidence.

## Go

In Go code, do not ignore errors with `_ = ...`. Handle, propagate, or report
the error. Use `%w` when adding context that callers can unwrap. Keep
`errcheck ./...` passing when Go behavior changes and the repo expects it.

## Tests

Assert expected failures explicitly. Avoid tests that pass because setup failed
or because an external dependency was unavailable unless the test is deliberately
checking that skip path.
