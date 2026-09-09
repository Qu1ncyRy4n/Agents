# Handle Errors Explicitly

Source: [`promisegrid_promisegrid_refs_heads_main_AGENTS.md`](../reference/source-guides/promisegrid_promisegrid_refs_heads_main_AGENTS.md#error-handling-policy-required)

> - Never use `|| true` in scripts, templates, or make recipes. Always inspect
>   command exit codes explicitly with `if/else` branches and handle each outcome.
> - For non-fatal cleanup/diagnostics steps, record command status (exit code and
>   logs) explicitly; do not fail silently.
> - In Go code, never ignore errors with `_ = ...`; handle, propagate, or report
>   errors explicitly.
> - Run `errcheck ./...` from the relevant Go module and keep it passing for Go changes.
