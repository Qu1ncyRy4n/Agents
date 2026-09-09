---
name: promisegrid-poc
description: Create, repair, or extend a PromiseGrid proof of concept while preserving inherited behavior, architecture lessons, analyzer gates, and acceptance criteria. Use for PromiseGrid POC work.
---

# Preserve PromiseGrid POC Behavior

## Source Rules

Source: [`promisegrid_wire-lab_refs_heads_main_AGENTS.md`](../../reference/source-guides/promisegrid_wire-lab_refs_heads_main_AGENTS.md#poc-superset-discipline-required)

> Future PromiseGrid POCs must be supersets of the previous POC's implemented
> behavior, architecture lessons, analyzer gates, and documented acceptance
> criteria unless a scoped DI explicitly declares the new POC non-superset and
> lists every intentionally dropped feature. (DI-sinur)
>
> A new POC may add focus, specialization, or new protocol surfaces, but it must
> not silently regress app/kernel boundaries, local trust semantics,
> monitor/analyzer gates, pCID routing, Promise Theory vocabulary, or previously
> proven workflows. (DI-sinur)
>
> When repairing or extending a POC, the analyzer must include inherited
> regression gates for the prior POC lineage, and the README or run narrative
> must state whether the POC is a superset or cite the DI authorizing an
> exception. (DI-sinur)

## Procedure

Adaptation of the source rules:

1. Identify the previous POC and its implemented behavior, architecture
   lessons, analyzer gates, documented acceptance criteria, and proven
   workflows.
2. Map each inherited item to the proposed POC. Do not treat new focus or
   specialization as permission to drop inherited behavior.
3. Add inherited regression gates to the analyzer.
4. State in the README or run narrative that the POC is a superset. If it is
   not, stop until a scoped DI explicitly lists every dropped feature.
5. Run the inherited and new gates. Report any regression directly.
