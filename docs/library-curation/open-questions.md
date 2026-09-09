# Open Questions

## Questions For Steve

### Scope And Placement Boundary

- Define and approve a source-disposition boundary for four scopes: universal,
  organization-universal, project-specific, and person/machine-specific. Source
  provenance must not itself classify a rule.
- Decide the intended output surface for Go and Git material: language/tool
  guide, triggered skill, specification, organization policy, or another
  defined layer. The question includes the exact cdint-grid build/test
  (`/home/qix/dev/cdint/cdint-grid/AGENTS.md:230-262`) and commit
  (`/home/qix/dev/cdint/cdint-grid/AGENTS.md:1687-1695`) sources, and must
  separate portable behavior from Go/project command facts.
- Confirm the placement and naming convention for organization vocabulary:
  retain `GLOSSARY.md` at repo root, move it under `docs/`, or define an adapter
  rule. Cite current cdint-grid `AGENTS.md:308-376` and `AGENTS.md:1696-1701`.
- Define whether `AGENTS.md` is an exact required filename or an adapter-level
  example for a canonical instruction layer plus non-relaxing role overlays.
  Cite cdint-grid `AGENTS.md:263-266`. This question is about replacing literal
  filenames with a target's selected instruction-file convention while
  retaining the core/overlay rule; it is not proposing removal of the rule.
- Review the cdint-grid source as a coupled body: determine who has authority to
  split organization policy, project protocol, tool/language guidance, and
  personal/machine state. Cite current source ranges for Architecture Precedent
  `AGENTS.md:377-410`, Testing `AGENTS.md:1644-1686`, and Commit
  `AGENTS.md:1687-1695`.
- Preserve TE protocol as organization-canonical and generally reusable.
  Resolve separately which TE artifact locations and record naming details are
  organization conventions versus project configuration. Cite `AGENTS.md:1361-1531`.

### Existing Evidence

- `/home/qix/dev/cdint/agents_library/intake/exemplar/wire-lab-exemplar.md:57-72,95-106,278-293`
- `/home/qix/dev/cdint/Agents/docs/LIBRARY-REVIEW-PLAN.md:44-87`
- `/home/qix/dev/cdint/Agents/docs/LIBRARY-COVERAGE-AUDIT.md:25-57,122-137`
- `/home/qix/dev/cdint/Agents/docs/thought-experiments/TE-kavam-metadata-tags-and-library-shape.md:57-59`

- Planning records: confirm root `TODO/` plus proquint handles as the CDINT-wide
  default, replacing wire-lab's per-protocol layout for new work. Current
  treatment: tentative default in `guides/development/planning-records.md`.
- TE editing categories: confirm whether the historical Cat-1a through Cat-7
  labels should remain the CDINT-wide interface or eventually gain clearer
  names. Current treatment: retain the functional category policy unchanged;
  the numbering is unusual but not currently harmful.
- Diff discipline: confirm the full no-unrelated-rewrites, rewrapping,
  reordering, or prose-normalization rule as a CDINT-wide default. Current
  treatment: retain the complete PromiseGrid source rule in
  `agents/diff-discipline.md`.

## Deferred Work

- Define a cleanup skill that inventories local state and generated artifacts,
  reports findings, and asks before destructive cleanup. It supplements, never
  replaces, `agents/repository-hygiene.md`.
- Revisit whether every skill repeats its authority limit or one global
  `agents/skill-use.md` rule is sufficient.
- Design tool-specific output adapters for AGENTS.md, CLAUDE.md, GEMINI.md,
  imports, nested/path-scoped instructions, and override behavior.
- Revisit optional generated skill summaries, for example a Mogent
  `include_skill_summary = true` setting that emits one-line descriptions.
- Revisit whether PromiseGrid material remains under `guides/promisegrid/` and
  `skills/promisegrid-*` after dogfooding and discussion with Steve.
- Decide whether `guides/promisegrid/promise-action-minimalism.md` is selected by
  default for all CDINT members, and define its explicit-exclusion mechanism.
- Consider generalizing `skills/promisegrid-poc/` into a non-PromiseGrid POC
  preservation skill after the domain-specific version is tested.
- Reconsider DEV-GUIDE-RESOURCES maintenance as a generated-document skill.
  Current treatment: retained only in reference source guides.
- Ask Steve whether ordinary work may use approved bounded path sets while
  exact per-path approval is reserved for risky writes, protected data,
  destructive actions, and external side effects.
- Define mutually exclusive decision-first profiles after dogfooding: strict
  default, behavior-change lock, and risk-triggered workflow.
- Consider a configurable decision-question batch limit, such as `n` questions
  per round; current behavior uses small batches and one-at-a-time dependent
  questions.
- After dogfooding `agents/decision-first.md`, try extracting its ordered
  procedure into a skill. Keep the agent version recoverable if routing or
  adherence becomes worse.
- Define mutually exclusive TE trigger profiles: default only when multiple
  plausible designs remain, optional proposal for every non-trivial decision,
  and mandatory TE for every non-trivial decision.
- Revisit a reduced risk-based TE scenario profile. Current default uses the
  full source scenario list whenever relevant.
- Revisit standalone TE artifact scope. Current default always writes one;
  possible later profile requires one only for durable or high-consequence
  decisions.
- Revisit the coding-style split among specifications, guides, and skills.
  Current treatment preserves the complete source under
  `guides/development/go-coding-style.md`; the global-state warning remains
  useful, while the object-oriented and Go-specific parts may later be trimmed
  or moved.
- Revisit a portable error-handling variant that states the invariant without
  Go- and shell-specific syntax. Current treatment retains explicit shell exit
  statuses and the prohibition on ignored Go errors.
- Revisit moving `gofmt` and `errcheck` requirements from always-loaded agent
  guidance into a Go-specific guide. Current treatment retains both for initial
  dogfooding.
