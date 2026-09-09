# Use The Basic Go Workflow

For any Go project:

- `go test ./...` runs the current module's test suite.
- `gofmt -w .` (or `go fmt ./...`) formats Go code.

Run these commands from the relevant module root. The `./...` pattern covers
the current module; it does not cross into nested modules.

If the local `go` binary and compiled standard-library objects are from
different Go versions, report the environment blocker instead of changing
repository files.

For project-specific build systems, Make targets, nested modules, release
steps, or extended validation, follow the repository's build guide, such as
`build.md`, when present.

Sources:

- [`promisegrid_wire-lab_refs_heads_main_AGENTS.md`](../reference/source-guides/promisegrid_wire-lab_refs_heads_main_AGENTS.md#build-test-and-development-commands)
- [`promisegrid_promisegrid_refs_heads_main_AGENTS.md`](../reference/source-guides/promisegrid_promisegrid_refs_heads_main_AGENTS.md#build-test-and-development-commands)
- [`newest_promisegrid.md`](../reference/source-guides/newest_promisegrid.md#build-test-and-development-commands)
