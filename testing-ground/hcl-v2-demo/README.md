# HCL V2 Dogfood Demo

This is a runnable local v2 fixture. It uses the sibling
`../../../../cdint-demo-lib/` server-side library.

Run from the Mogent repository root:

```sh
go run ./cmd/mogent plan  --config testing-ground/hcl-v2-demo/demo.mogent.hcl
go run ./cmd/mogent apply --config testing-ground/hcl-v2-demo/demo.mogent.hcl
```

`plan` prints a unified diff per output path and writes nothing. `apply`
writes identical `AGENTS.md` and `CLAUDE.md` outputs in this directory. Use
`apply --force` only after reviewing a direct edit to a generated output that
`plan` reported as MOGENT208.
