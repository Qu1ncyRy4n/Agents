# CIWG / PromiseGrid Consolidation User Story

## Consolidate agents and skills across PromiseGrid and CIWG repositories

**As a** PromiseGrid/CIWG maintainer, **I want** the duplicated agent guidance
and skills currently copy-pasted across `promisegrid/wire-lab`,
`promisegrid/grid-poc`, `ciwg/grid-examples`, `ciwg/mob-sandbox`,
`ciwg/decomk-conf-cswg`, and their sibling repos consolidated into one
versioned, composable library **so that** a shared rule is authored once, fixed
once, and every repo renders the same `AGENTS.md`/skill set from an
`agents.yaml` manifest instead of drifting apart.

### Scope

- **ciwg org** (https://github.com/ciwg): `grid-examples`, `cswg`,
  `mob-sandbox`, `decomk-conf-cswg`, `FAB26-Presentation`,
  `promisegrid-dev-guide`, plus the workshop template chain.
- **PromiseGrid org**: `wire-lab`, `grid-poc`, `promisegrid`, and Steve's
  `stevegt/*` repos.
- Existing local captures in `docs/other_repo_agents/` already show the
  duplication: `wire-lab` and `grid-examples` both carry near-identical
  *Agent Instruction Architecture*, *Promise Action Minimalism*, *POC Superset
  Discipline*, *Decision-First*, and *Thought Experiment Protocol* sections.

### Acceptance Criteria

- One inventory classifies each captured section/skill as universal,
  org-universal, project-specific, or person/machine-specific
  (`docs/library-curation/open-questions.md:7-9`).
- Shared modules are extracted (e.g. `agents/decision-first.md`,
  `agents/diff-discipline.md`, `guides/promisegrid/promise-action-minimalism.md`,
  `skills/promisegrid-poc/`) with provenance back to the originating repo and
  commit.
- Each source repo gets a manifest selecting those modules and renders an
  `AGENTS.md` that matches its current one modulo intentional edits, proven by
  `mogent build --dry-run`.
- Divergences between repos are surfaced as conflicts for review, not silently
  normalized, and historical DI/TE/DR handles remain resolvable.
- Adding a new CIWG repo requires only a manifest, not a fresh copy of the
  protocols.
