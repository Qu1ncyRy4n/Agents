# Protolibrary Provenance Ledger

This is the human-readable source-to-module ledger. It starts coarse and must
become heading-level before the protolibrary can claim exhaustive coverage.

Coverage states are `unreviewed`, `partial`, `represented`, `repo-local`,
`generated-assessment`, `duplicate`, or `rejected-with-reason`.

| Source | Initial classification | Candidate destinations | Coverage |
|---|---|---|---|
| [`ciwg_cswg`](../docs/other_repo_agents/ciwg_cswg_refs_heads_main_AGENTS.md) | captured repository guide | static sites, tools/Git, testing, TODO, handoff | unreviewed |
| [`ciwg_decomk-conf-cswg`](../docs/other_repo_agents/ciwg_decomk-conf-cswg_refs_heads_main_AGENTS.md) | captured repository guide; strict governance choice | process, engineering, tools/Git, TODO, handoff | partial |
| [`ciwg_FAB26-Presentation`](../docs/other_repo_agents/ciwg_FAB26-Presentation_refs_heads_main_AGENTS.md) | authoritative repository guide; governance plus public prose | process, documentation, provenance | unreviewed |
| [`ciwg_grid-examples`](../docs/other_repo_agents/ciwg_grid-examples_refs_heads_main_AGENTS.md) | authoritative repository guide; repeated governance variant | process, PromiseGrid, engineering | unreviewed |
| [`ciwg_mob-sandbox`](../docs/other_repo_agents/ciwg_mob-sandbox_refs_heads_main_AGENTS.md) | authoritative repository guide; repeated governance variant | process, engineering | unreviewed |
| [`computerscienceiscool_llm-runtime`](../docs/other_repo_agents/computerscienceiscool_llm-runtime_refs_heads_audit-sweep_AGENTS.md) | captured repository guide | repository structure, tests, tools/Git, security/local state | unreviewed |
| [`computerscienceiscool_pg`](../docs/other_repo_agents/computerscienceiscool_pg_refs_heads_main_AGENTS.md) | captured repository guide | coordination IDs, decision process, comments, Go, testing | unreviewed |
| [`chroma_terra`](<../docs/other_repo_agents/low qual/chroma_terra_agents.md>) | generated repository-grounded guide | teaching, Rust/Python boundary, architecture decisions | generated-assessment |
| [`color-engine`](<../docs/other_repo_agents/low qual/color-engine_agents.md>) | generated repository-grounded guide | Rust library boundaries, deterministic algorithms, workspace dependencies | generated-assessment |
| [`DriftDiffusionModel`](<../docs/other_repo_agents/low qual/DriftDiffusionModel_agents.md>) | generated repository-grounded guide | exploratory research stages, source-of-truth distinctions | generated-assessment |
| [`eink_screenshare`](<../docs/other_repo_agents/low qual/eink_screenshare_agents.md>) | generated repository-grounded guide | architecture constraints, protocol evolution, platform validation | generated-assessment |
| [`EYNA_Eyetracking`](<../docs/other_repo_agents/low qual/EYNA_Eyetracking_agents.md>) | generated repository-grounded guide | generated-doc ownership, research API reference, no invented build | generated-assessment |
| [`grey_helper`](<../docs/other_repo_agents/low qual/grey_helper_agents.md>) | generated repository-grounded guide | archive mutation, parsing fixtures, generated artifacts, GUI smoke tests | generated-assessment |
| [`NeuralOptimizationProject`](<../docs/other_repo_agents/low qual/NeuralOptimizationProject_agents.md>) | generated repository-grounded guide | research data, vendored ownership, Python/uv, first-analysis documentation | generated-assessment |
| [`NixConfig assessment`](<../docs/other_repo_agents/low qual/NixConfig_AGENTS.md>) | generated assessment and suggested guide | Nix host safety, migration status, comments, non-activating validation | generated-assessment |
| [`pdf_md_extractor`](<../docs/other_repo_agents/low qual/pdf_md_extractor_agents.md>) | generated repository-grounded guide | file/archive dry runs, idempotence, CLI/dependency validation | generated-assessment |
| [`pg_learning`](<../docs/other_repo_agents/low qual/pg_learning_AGENTS.md>) | generated consolidation containing teaching rules and private state | communication, teaching, generated-docs; private state excluded | generated-assessment |
| [`photo_getter`](<../docs/other_repo_agents/low qual/photo_getter_agents.md>) | generated repository-grounded guide | external services, credentials, archive/DB safety, dry runs | generated-assessment |
| [`piazza_scrape`](<../docs/other_repo_agents/low qual/piazza_scrape_agents.md>) | generated repository-grounded guide | scraping, external-service pacing, credentials, offline fixtures | generated-assessment |
| [`RoSE`](<../docs/other_repo_agents/low qual/RoSE_agents.md>) | generated guide with valuable repository-grounded safety rules | research safety, high-consequence review, data protection | generated-assessment |
| [`SEFproject`](<../docs/other_repo_agents/low qual/SEFproject_agents.md>) | generated repository-grounded guide | research pipelines, limitations, no fabricated data | generated-assessment |
| [`SemanticMapProject`](<../docs/other_repo_agents/low qual/SemanticMapProject_agents.md>) | generated repository-grounded guide | Python/Nix environments, research analysis, project structure | generated-assessment |
| [`promisegrid_grid-poc`](../docs/other_repo_agents/promisegrid_grid-poc_refs_heads_main_AGENTS.md) | captured repository guide | Go, tools/Git, testing, TODO, reference maps | unreviewed |
| [`promisegrid_promisegrid`](../docs/other_repo_agents/promisegrid_promisegrid_refs_heads_main_AGENTS.md) | authoritative repository guide; repeated governance variant | process, PromiseGrid | unreviewed |
| [`promisegrid_wire-lab`](../docs/other_repo_agents/promisegrid_wire-lab_refs_heads_main_AGENTS.md) | authoritative repository guide; unique POC and protocol constraints | process, PromiseGrid, engineering | unreviewed |
| [`stevegt_decomk`](../docs/other_repo_agents/stevegt_decomk_refs_heads_main_AGENTS.md) | authoritative repository guide; governance family | process, engineering | unreviewed |
| [`stevegt_godecide`](../docs/other_repo_agents/stevegt_godecide_refs_heads_main_AGENTS.md) | captured repository guide | Go, testing, TODO, tools/Git | unreviewed |
| [`stevegt_grokker`](../docs/other_repo_agents/stevegt_grokker_refs_heads_main_AGENTS.md) | captured repository guide | Go/JavaScript, testing, tools/Git, runtime state | unreviewed |
| [`stevegt_mob-consensus`](../docs/other_repo_agents/stevegt_mob-consensus_refs_heads_main_AGENTS.md) | captured repository guide | Go, tests, TODO, tools/Git, repository-local collaboration | unreviewed |
| [`stevegt_navlog`](../docs/other_repo_agents/stevegt_navlog_refs_heads_main_AGENTS.md) | authoritative repository guide; repeated strict workflow | process, engineering | unreviewed |
| [`nix-dotfiles assessment`](<../docs/other_repo_agents/tbd/nix-dotfiles_agents.md>) | generated assessment and suggested guide | Nix host safety, activation, validation, style | generated-assessment |
| [`Notes Obsidian Skills`](<../docs/other_repo_agents/tbd/Notes_Obsidian_Skills_AGENTS.md>) | generated reference/consolidation | notes style, skills, flashcards, templates | generated-assessment |
| [`Notes Vault Agents`](<../docs/other_repo_agents/tbd/Notes_Vault_Agents_AGENTS.md>) | generated reference containing private state and personas | notes, communication choices, private state | generated-assessment |
| [`openai_codex`](../docs/other_repo_agents/tbd/openai_codex_refs_heads_main_AGENTS.md) | authoritative upstream guide captured for comparison | Rust, testing, API evolution, TUI | unreviewed |
| [`seam_game assessment`](<../docs/other_repo_agents/tbd/seam_game_AGENTS.md>) | generated assessment and suggested guide, not source authority | architecture-law and tutoring candidates | generated-assessment |
| [`todo_app`](<../docs/other_repo_agents/tbd/todo_app_AGENTS.md>) | generated repository-grounded guide | Rust, exploratory design, typed models, dependency restraint | generated-assessment |
| [`Video project`](<../docs/other_repo_agents/tbd/VideoProject_agents.md>) | generated guide with valuable repository-grounded safety rules | research safety, source-of-truth, MATLAB | generated-assessment |

All 37 source-guide files are now present in this file-level ledger. File-level
presence is not exhaustive conceptual coverage. The next pass must add every
useful heading and give it an explicit disposition; do not infer coverage from
its containing directory name or initial quality label.

## First Heading-Level Mapping

| Source heading | Candidate module | Treatment | Coverage |
|---|---|---|---|
| [`Decision-First Specification and Compliance Protocol`](../docs/other_repo_agents/ciwg_decomk-conf-cswg_refs_heads_main_AGENTS.md#decision-first-specification-and-compliance-protocol-required) | `orgs/cdint/process/decision-governance/` | preserve as a strict option; do not make universal | partial |
| [`Thought Experiment Protocol`](../docs/other_repo_agents/ciwg_decomk-conf-cswg_refs_heads_main_AGENTS.md#thought-experiment-protocol-required) | `intake/process/decision-governance/use-thought-experiments-to-narrow-a-broad-design-space.md` | preserve workflow and correct the v1 reversal | represented |
| [`TE Intake Requirements`](../docs/other_repo_agents/ciwg_decomk-conf-cswg_refs_heads_main_AGENTS.md#te-intake-requirements) | same module | verbatim evidence | represented |
| [`TE Execution Requirements`](../docs/other_repo_agents/ciwg_decomk-conf-cswg_refs_heads_main_AGENTS.md#te-execution-requirements) | same module | verbatim evidence | represented |
| [`TE Output to DF`](../docs/other_repo_agents/ciwg_decomk-conf-cswg_refs_heads_main_AGENTS.md#te-output-to-df) | same module | verbatim evidence | represented |
| [`Diff Discipline`](../docs/other_repo_agents/ciwg_FAB26-Presentation_refs_heads_main_AGENTS.md#diff-discipline-required) | `intake/workflow/keep-changes-scoped-and-preserve-existing-work.md` | verbatim evidence plus proposed examples | represented |
| [`Diff Discipline`](../docs/other_repo_agents/promisegrid_promisegrid_refs_heads_main_AGENTS.md#diff-discipline-required) | same module | duplicate wording cited rather than copied twice | represented |
| [`Agent Behavior`](<../docs/other_repo_agents/tbd/todo_app_AGENTS.md#agent-behavior>) | same module | generated-source evidence; preserve scope and user-work rules for review | partial |
| [`Decision-First Specification and Compliance Protocol`](../docs/other_repo_agents/ciwg_decomk-conf-cswg_refs_heads_main_AGENTS.md#decision-first-specification-and-compliance-protocol-required) | `intake/workflow/use-risk-to-choose-routine-work-or-decision-review.md` and `intake/process/decision-governance/lock-durable-decisions-before-implementation.md` | preserve strict option and compare with risk-based v1 adaptation | partial |
| [`Decision Lock and Stop Rule`](../docs/other_repo_agents/ciwg_decomk-conf-cswg_refs_heads_main_AGENTS.md#decision-lock-and-stop-rule) | `intake/process/decision-governance/lock-durable-decisions-before-implementation.md` | verbatim evidence; strictness needs review | represented |
| [`Comment Preservation Protocol`](../docs/other_repo_agents/ciwg_decomk-conf-cswg_refs_heads_main_AGENTS.md#comment-preservation-protocol-required) | `intake/process/decision-governance/preserve-decision-history-through-supersession.md` | extract DI history and supersession separately from comment policy | partial |
| [`Required final handoff artifacts`](../docs/other_repo_agents/ciwg_decomk-conf-cswg_refs_heads_main_AGENTS.md#required-final-handoff-artifacts) | `intake/process/decision-governance/handoff-decisions-with-implementation-evidence.md` | verbatim strict option; compare with ordinary handoff | represented |

## Exemplar Anchor

The owner redirected review order on 2026-09-02: `promisegrid_wire-lab` is the
exemplar anchor. Its full text is preserved verbatim in
`intake/exemplar/wire-lab-exemplar.md` with per-section contrasts against the
other 36 guides and numbered unresolved diffs (D1–D12 in that document). The
exemplar's disposition sketch, once reviewed, drives promotion into
`orgs/<owner>/`; all other intake candidates are triaged and deferred until the
anchor is placed. Per-section dispositions for wire-lab headings are recorded
in the exemplar document and must be mirrored here as they are resolved.

| Source heading | Candidate module | Treatment | Coverage |
|---|---|---|---|
| [`promisegrid_wire-lab` (complete guide)](../docs/other_repo_agents/promisegrid_wire-lab_refs_heads_main_AGENTS.md) | `intake/exemplar/wire-lab-exemplar.md` | verbatim anchor with cross-guide contrasts and diff index | partial |
