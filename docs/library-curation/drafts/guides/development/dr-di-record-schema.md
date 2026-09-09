# DR/DI Record Schema

Output type: development guide. This documents the current organization-record
schema; it complements the always-loaded portable authority rule in
[`agents/decision-record-authority.md`](../../agents/decision-record-authority.md).

Sources: [`promisegrid_wire-lab_refs_heads_main_AGENTS.md`](../../reference/source-guides/promisegrid_wire-lab_refs_heads_main_AGENTS.md#drdi-source-of-truth-protocol-required),
[`newest_promisegrid.md`](../../reference/source-guides/newest_promisegrid.md#drdi-source-of-truth-protocol-required),
[`ciwg_grid-examples_refs_heads_main_AGENTS.md`](../../reference/source-guides/ciwg_grid-examples_refs_heads_main_AGENTS.md#drdi-source-of-truth-protocol-required),
[`promisegrid_promisegrid_refs_heads_main_AGENTS.md`](../../reference/source-guides/promisegrid_promisegrid_refs_heads_main_AGENTS.md#drdi-source-of-truth-protocol-required),
and [`ciwg_FAB26-Presentation_refs_heads_main_AGENTS.md`](../../reference/source-guides/ciwg_FAB26-Presentation_refs_heads_main_AGENTS.md#drdi-source-of-truth-protocol-required).

- Use `user@example.com (FirstName)` for person identity in DR and DI records.
- In DRs, use that format for `Asked by` and person-valued `Waiting on` fields.
  In DIs, use it for `Author`.
- Keep legacy numeric TODO and DI records valid until they are migrated.
- Apply DR/DI tracking rules incrementally as files are brought under tracking.
