# HCL V2 Dogfood Demo

This is a runnable local v2 fixture. It uses the sibling
`../../../../cdint-demo-lib/` server-side library and does not require a remote
lock.

Run from the Mogent repository root:

```sh
go run ./cmd/mogent plan --config testing-ground/hcl-v2-demo/demo.mogent.hcl
go run ./cmd/mogent build --config testing-ground/hcl-v2-demo/demo.mogent.hcl
```

The build writes identical `AGENTS.md` and `CLAUDE.md` outputs in this directory.
Use `--dry-run` to validate without writing and `--force` only after reviewing a
direct edit to a generated output.
