---
tags: [lang/go, workflow/development, testing/go, scope/project]
tldr: Run gofmt, focused go test, and errcheck for Go behavior changes.
priority: 0.9
scope: project
requires: [shared:shared-baseline/instructions/focused-change-loop]
---
# Development

Run `gofmt` on changed Go code. Run focused `go test` from the affected module.
Run `errcheck ./...` for Go behavior changes. If a required command is missing
or the environment is broken, report the blocker instead of changing unrelated
files.
