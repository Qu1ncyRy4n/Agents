# Organize Go Packages by Purpose

Output type: guide. Keep separate from always-loaded agent guidance.

## CDINT Default

Source: [`promisegrid_wire-lab_refs_heads_main_AGENTS.md`](../../reference/source-guides/promisegrid_wire-lab_refs_heads_main_AGENTS.md#project-structure--module-organization)

> Keep packages at module root or under purpose-named top-level directories
> (`contexts/`, `state/`, etc.); avoid `internal/` and `pkg/`.

## Documented Alternative

Source: [`computerscienceiscool_llm-runtime_refs_heads_audit-sweep_AGENTS.md`](../../reference/source-guides/computerscienceiscool_llm-runtime_refs_heads_audit-sweep_AGENTS.md#project-structure--module-organization)

> `cmd/llm-runtime/` is the CLI entry; `make build` emits `./llm-runtime`.
>
> `pkg/` holds public packages (`app`, `cli`, `config`, `evaluator`, `sandbox`,
> `scanner`, `search`, `session`); `internal/core/` is private.

Use the CDINT default unless the repository explicitly selects and documents
the alternative layout.
