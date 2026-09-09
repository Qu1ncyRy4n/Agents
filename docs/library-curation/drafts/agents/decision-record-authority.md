# Keep Decision Records Authoritative

Current profile: always-loaded for initial dogfooding. Keep this authority and
provenance rule visible in agent guidance; revisit skill extraction after
dogfooding.

Sources: [`promisegrid_wire-lab_refs_heads_main_AGENTS.md`](../reference/source-guides/promisegrid_wire-lab_refs_heads_main_AGENTS.md#drdi-source-of-truth-protocol-required),
[`newest_promisegrid.md`](../reference/source-guides/newest_promisegrid.md#drdi-source-of-truth-protocol-required),
[`ciwg_grid-examples_refs_heads_main_AGENTS.md`](../reference/source-guides/ciwg_grid-examples_refs_heads_main_AGENTS.md#drdi-source-of-truth-protocol-required),
[`promisegrid_promisegrid_refs_heads_main_AGENTS.md`](../reference/source-guides/promisegrid_promisegrid_refs_heads_main_AGENTS.md#drdi-source-of-truth-protocol-required),
and [`ciwg_FAB26-Presentation_refs_heads_main_AGENTS.md`](../reference/source-guides/ciwg_FAB26-Presentation_refs_heads_main_AGENTS.md#drdi-source-of-truth-protocol-required).

- Treat DR and DI records as the source of truth for decisions and open
  questions. Documents and code are outputs of that process and link back to
  the relevant record.
- Cite at least one DI ID for a settled statement in documentation or critical
  code-comment logic. Cite at least one DR ID for an unresolved question or
  uncertainty. Create a DR before finalizing a change with an unresolved
  question.
- A DI author is the decision-maker, not merely the recorder. An agent is the
  author only when decision authority was explicitly delegated to that agent.
- Use one globally unique proquint handle namespace for TODO, TE, DR, and DI
  records. Defer adding DN until CDINT adopts a DN record type.
