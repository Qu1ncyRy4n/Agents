# Protolibrary Provenance Ledger

This is the human-readable source-to-module ledger. It starts coarse and must
become heading-level before the protolibrary can claim exhaustive coverage.

Coverage states are `unreviewed`, `partial`, `represented`, `repo-local`,
`generated-assessment`, `duplicate`, or `rejected-with-reason`.

| Source | Initial classification | Candidate destinations | Coverage |
|---|---|---|---|
| [`ciwg_decomk-conf-cswg`](../docs/other_repo_agents/ciwg_decomk-conf-cswg_refs_heads_main_AGENTS.md) | authoritative repository guide; strict governance family | `orgs/cdint/process`, engineering, Git, TODO, handoff | partial |
| [`ciwg_FAB26-Presentation`](../docs/other_repo_agents/ciwg_FAB26-Presentation_refs_heads_main_AGENTS.md) | authoritative repository guide; governance plus public prose | process, documentation, provenance | unreviewed |
| [`ciwg_grid-examples`](../docs/other_repo_agents/ciwg_grid-examples_refs_heads_main_AGENTS.md) | authoritative repository guide; repeated governance variant | process, PromiseGrid, engineering | unreviewed |
| [`ciwg_mob-sandbox`](../docs/other_repo_agents/ciwg_mob-sandbox_refs_heads_main_AGENTS.md) | authoritative repository guide; repeated governance variant | process, engineering | unreviewed |
| [`promisegrid_promisegrid`](../docs/other_repo_agents/promisegrid_promisegrid_refs_heads_main_AGENTS.md) | authoritative repository guide; repeated governance variant | process, PromiseGrid | unreviewed |
| [`promisegrid_wire-lab`](../docs/other_repo_agents/promisegrid_wire-lab_refs_heads_main_AGENTS.md) | authoritative repository guide; unique POC and protocol constraints | process, PromiseGrid, engineering | unreviewed |
| [`stevegt_decomk`](../docs/other_repo_agents/stevegt_decomk_refs_heads_main_AGENTS.md) | authoritative repository guide; governance family | process, engineering | unreviewed |
| [`stevegt_navlog`](../docs/other_repo_agents/stevegt_navlog_refs_heads_main_AGENTS.md) | authoritative repository guide; repeated strict workflow | process, engineering | unreviewed |
| [`openai_codex`](../docs/other_repo_agents/tbd/openai_codex_refs_heads_main_AGENTS.md) | authoritative upstream guide captured for comparison | Rust, testing, API evolution, TUI | unreviewed |
| [`RoSE`](<../docs/other_repo_agents/low qual/RoSE_agents.md>) | generated guide with valuable repository-grounded safety rules | research safety, high-consequence review | unreviewed |
| [`Video project`](<../docs/other_repo_agents/tbd/VideoProject_agents.md>) | generated guide with valuable repository-grounded safety rules | research safety, source-of-truth, MATLAB | unreviewed |
| [`pg_learning`](<../docs/other_repo_agents/low qual/pg_learning_AGENTS.md>) | generated consolidation containing teaching rules and private state | communication, teaching; private state excluded | unreviewed |
| [`seam_game assessment`](<../docs/other_repo_agents/tbd/seam_game_AGENTS.md>) | generated assessment and suggested guide, not source authority | architecture-law and tutoring candidates | generated-assessment |
| [`NixConfig assessment`](<../docs/other_repo_agents/low qual/NixConfig_AGENTS.md>) | generated assessment and suggested guide | Nix safety candidates | generated-assessment |
| [`nix-dotfiles assessment`](<../docs/other_repo_agents/tbd/nix-dotfiles_agents.md>) | generated assessment and suggested guide | Nix safety candidates | generated-assessment |

The remaining source files under `docs/other_repo_agents/` are still
`unreviewed`. Add every file and then every useful heading; do not infer
coverage from its containing directory name.

## First Heading-Level Mapping

| Source heading | Candidate module | Treatment | Coverage |
|---|---|---|---|
| [`Decision-First Specification and Compliance Protocol`](../docs/other_repo_agents/ciwg_decomk-conf-cswg_refs_heads_main_AGENTS.md#decision-first-specification-and-compliance-protocol-required) | `orgs/cdint/process/decision-governance/` | preserve as a strict option; do not make universal | partial |
| [`Thought Experiment Protocol`](../docs/other_repo_agents/ciwg_decomk-conf-cswg_refs_heads_main_AGENTS.md#thought-experiment-protocol-required) | `orgs/cdint/process/decision-governance/use-thought-experiments-to-narrow-a-broad-design-space.md` | preserve workflow and correct the v1 reversal | represented |
| [`TE Intake Requirements`](../docs/other_repo_agents/ciwg_decomk-conf-cswg_refs_heads_main_AGENTS.md#te-intake-requirements) | same module | verbatim evidence | represented |
| [`TE Execution Requirements`](../docs/other_repo_agents/ciwg_decomk-conf-cswg_refs_heads_main_AGENTS.md#te-execution-requirements) | same module | verbatim evidence | represented |
| [`TE Output to DF`](../docs/other_repo_agents/ciwg_decomk-conf-cswg_refs_heads_main_AGENTS.md#te-output-to-df) | same module | verbatim evidence | represented |

