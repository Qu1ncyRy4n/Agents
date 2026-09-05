---
status: needs-choice
candidate_outputs: [agent-guide, library-authoring]
tags: [exemplar, anchor, process/decision, workflow/te, engineering/comments, engineering/errors, tools/git, process/artifacts]
tldr: The Promisegrid wire-lab guide is the v2 exemplar anchor; each section carries contrasts from the other source guides and numbered diffs to resolve before org promotion.
sources:
  - docs/other_repo_agents/promisegrid_wire-lab_refs_heads_main_AGENTS.md
  - docs/other_repo_agents/newest_promisegrid.md
  - docs/other_repo_agents/cdint-grid_skills/
requires_all: []
mutually_exclusive: []
see_also: []
---

# Wire-Lab Exemplar: Anchor Guide with Cross-Guide Contrasts

This document is the v2 exemplar. The anchor text is a verbatim copy of
[`promisegrid_wire-lab`](../../../docs/other_repo_agents/promisegrid_wire-lab_refs_heads_main_AGENTS.md),
the most complete captured repository guide. Every other section below it adds
contrast material from the other source guides, so overlap, variants, and
conflicts are visible in one place. The current cdint-grid guide is captured as
[`newest_promisegrid.md`](../../../docs/other_repo_agents/newest_promisegrid.md),
and its 18 progressively loaded operational skills are preserved under
[`cdint-grid_skills/`](../../../docs/other_repo_agents/cdint-grid_skills/). Both
captures are verbatim from cdint-grid commit
`79961cb3147054bda8e94b9703fdece2fbe30a8d` (`AGENTS.md` blob
`63de43205c1439fc4560160b2ccbf05bdd9813b3`; skills tree
`56425b68232c5ed2492b6139383323c68b4d7626`).

## How to Read This Document

- **Anchor** — verbatim wire-lab text, unmodified. It is the content that must
  be integrated into organization libraries first; all other intake content is
  triaged and deferred until this anchor is placed.
- **Contrast blocks** (`#### Contrasts and unresolved diffs`) are editorial
  additions from this review, not source text. Quotes from other guides are
  verbatim excerpts; comparisons and "diff to resolve" statements are
  adaptation/proposed material until reviewed.
- Contrast labels: **Same** (verbatim or near-verbatim family copy), **Variant**
  (same rule, different mechanism), **Stricter/Looser** (stronger or weaker
  form), **Missing from anchor** (other guides have content the anchor lacks),
  **Conflict** (both rules cannot hold as written), **Generated evidence**
  (non-authoritative generated guide supporting or contradicting the rule).
- Numbered diffs (D1, D2, …) are collected in the index at the end. Resolution
  is an owner decision; nothing here is approved by being committed.

---

# Repository Guidelines

## Project Structure & Module Organization
- Keep packages at module root or under purpose-named top-level directories (`contexts/`, `state/`, etc.); avoid `internal/` and `pkg/`.
- Keep planning artifacts in per-protocol `protocols/<slug>.d/TODO/` directories (harness-level under `protocols/wire-lab.d/TODO/`) and maintain the master cross-listed index at `protocols/wire-lab.d/TODO/TODO.md` sorted by priority. Each `protocols/<slug>.d/TODO/` also has its own per-protocol `TODO.md` queue.
- Do not commit local state files (for example `.grok`, `.grok.lock`) or generated binaries.

#### Contrasts and unresolved diffs
- **Same** ([`ciwg/grid-examples`](../../../docs/other_repo_agents/ciwg_grid-examples_refs_heads_main_AGENTS.md#project-structure--module-organization)): identical package-placement rule.
- **Variant** ([`promisegrid/promisegrid`](../../../docs/other_repo_agents/promisegrid_promisegrid_refs_heads_main_AGENTS.md#project-structure--module-organization)): planning artifacts live in a root `TODO/` directory, with optional root `DR/` and `docs/thought-experiments/` — not per-protocol directories.
- **Current cdint-grid variant** ([`newest_promisegrid`](../../../docs/other_repo_agents/newest_promisegrid.md#project-structure--module-organization)): root `TODO/` is now explicitly authoritative for this repository, durable worker-attempt records live beside their owning TODO, and design notes have a separate `DN-<handle>` identity. This is newer evidence that the anchor's per-protocol layout is repository-local rather than a family-wide default.
- **Variant** ([`promisegrid/grid-poc`](../../../docs/other_repo_agents/promisegrid_grid-poc_refs_heads_main_AGENTS.md#project-structure--module-organization)): root Go module plus `x/` standalone experiment modules; ([`stevegt/grokker`](../../../docs/other_repo_agents/stevegt_grokker_refs_heads_main_AGENTS.md#project-structure--module-organization)) uses a versioned `v3/` tree; ([`stevegt/godecide`](../../../docs/other_repo_agents/stevegt_godecide_refs_heads_main_AGENTS.md#project-structure--module-organization)) keeps core code at the root with `cmd/` entry.
- **Conflict** ([`computerscienceiscool/llm-runtime`](../../../docs/other_repo_agents/computerscienceiscool_llm-runtime_refs_heads_audit-sweep_AGENTS.md#project-structure--module-organization)): this guide deliberately uses `pkg/` for public packages and `internal/core/` for private code — the opposite of the anchor's "avoid `internal/` and `pkg/`". **D1 — resolve which layout rule is canonical and whether language- or tool-specific overlays may relax it.**
- **Same** (all captured guides): local state (`.grok`, `.grok.lock`) and generated binaries stay uncommitted.
- **Diff to resolve** — **D2:** planning-artifact placement and ID scheme split across the corpus: per-protocol directories with proquint handles (anchor, [`stevegt/decomk`](../../../docs/other_repo_agents/stevegt_decomk_refs_heads_main_AGENTS.md#project-structure--module-organization)), root `TODO/` with proquint handles ([`promisegrid/promisegrid`](../../../docs/other_repo_agents/promisegrid_promisegrid_refs_heads_main_AGENTS.md#todo-tracking)), and root `TODO/` with zero-padded numbers ([`promisegrid/grid-poc`](../../../docs/other_repo_agents/promisegrid_grid-poc_refs_heads_main_AGENTS.md#task-tracking-todo), [`stevegt/navlog`](../../../docs/other_repo_agents/stevegt_navlog_refs_heads_main_AGENTS.md#todo-tracking)). Decide whether this is canonical or repository-local.

## Build, Test, and Development Commands
- `go test ./...` runs the test suite.
- `gofmt -w .` (or `go fmt ./...`) formats Go code.

#### Contrasts and unresolved diffs
- **Variant** ([`promisegrid/promisegrid`](../../../docs/other_repo_agents/promisegrid_promisegrid_refs_heads_main_AGENTS.md#build-test-and-development-commands)): no root Go module; commands run per module. Adds a general rule the anchor lacks — **Missing from anchor:** "If Go commands fail because the local `go` binary and compiled standard-library objects are from different Go versions, report the environment blocker instead of changing repo files." Candidate for a canonical environment-blocker module.
- **Variant** ([`computerscienceiscool/llm-runtime`](../../../docs/other_repo_agents/computerscienceiscool_llm-runtime_refs_heads_audit-sweep_AGENTS.md#build-test-and-development-commands)): full Makefile surface (`make build/test/test-coverage/bench/quality`); relevant precedent for Mogent's open Makefile question (TODO M-P4).
- **Diff to resolve:** this whole section is repository-local command fact, not portable guidance. Proposed disposition: extract only generalizable rules (report environment blockers; keep caches out of the tree) and classify command lists as `repository-local`.

## Agent Instruction Architecture (Required)
- `AGENTS.md` is the canonical home for repo-wide protocol, workflow, and vocabulary rules. Role-specific files such as `AGENTS-codex.md` and `AGENTS-ppx.md` are overlays: they may add stricter role constraints, environment setup, and role-specific procedures, but they must not duplicate or relax canonical repo-wide rules. (DI-034-20260508-060134)
- When a rule applies to every agent or every repo artifact, move it here and replace role-file copies with pointers. When a rule applies only to one agent's runtime environment, identity, credentials, branch lifecycle, or private logging system, keep it in that role overlay. (DI-034-20260508-060134)

#### Contrasts and unresolved diffs
- **Looser** ([`promisegrid/promisegrid`](../../../docs/other_repo_agents/promisegrid_promisegrid_refs_heads_main_AGENTS.md#agent-instruction-architecture-required), [`ciwg/FAB26-Presentation`](../../../docs/other_repo_agents/ciwg_FAB26-Presentation_refs_heads_main_AGENTS.md#agent-instruction-architecture-required)): same canonical-home principle, two lines, no overlay mechanism.
- **Missing from anchor** ([`ciwg/decomk-conf-cswg`](../../../docs/other_repo_agents/ciwg_decomk-conf-cswg_refs_heads_main_AGENTS.md#currency-of-information) family): a *user-level* layer above repo guides — "Frequently check `~/.codex/AGENTS.md` for updates" plus a maintained cross-repo `~/.codex/meta-context.md`. The anchor models canonical-vs-role overlays but not user-level overlays. **Diff to resolve:** does the v2 model include user-level and role-level overlay layers above the canonical manifest? (Directly informs Mogent's manifest/overlay design; see `docs/MULTIPLE-OUTPUTS-PLAN.md`.)
- **Missing from anchor** ([`promisegrid/promisegrid`](../../../docs/other_repo_agents/promisegrid_promisegrid_refs_heads_main_AGENTS.md#public-artifact-provenance-required), [`ciwg/FAB26-Presentation`](../../../docs/other_repo_agents/ciwg_FAB26-Presentation_refs_heads_main_AGENTS.md#public-artifact-provenance-required)): public-artifact provenance rules (no DI/DR/TE references in slides; footnote tags in white papers with a `## References` section). These are output-target-specific rules. **Diff to resolve:** where artifact-type rules live when one library serves several output targets.
- This section is the closest source analogue to Mogent's own model (canonical manifest plus stricter overlays, promotion of shared rules upward). Proposed destination: `library-authoring`.
- **New execution layer** ([`cdint-grid skills`](../../../docs/other_repo_agents/cdint-grid_skills/)): the current guide remains canonical, but 18 tracked skills provide progressively loaded procedures, command routing, exact-once constraints, stop conditions, and output schemas. Skills repeatedly disclaim independent task, path, side-effect, decision, and acceptance authority. **D13 — decide whether Mogent models policy authority, operational procedure, discovery metadata, and executable tooling as separate typed layers, and how it checks drift between them.**
- **Composed workflow variant:** skills call or require other skills (`newtree` → `resume` → `upgrade` → `drain-inbox`; decision work may add `run-thought-experiment`, `review-against-ninik`, and `commit`; peer integration uses `consensus` and `commit`). The anchor describes rules but has no machine-readable dependency graph. **D14 — decide whether skill dependencies, compatibility state, preconditions, side effects, and output contracts are explicit metadata or remain prose.**

## Promise Action Minimalism (Required)
- Future PromiseGrid protocol, simulation, POC, scoring, generation, and guide work must not invent workflow-specific top-level action kinds by default. The default future-facing top-level semantic action is `promise`. (DI-mosoj)
- Treat observation as a promise that the promiser observed something from its local vantage. Treat refusal as absence of a promise, a promise not to do something, or a promise that the agent does not currently promise the requested behavior. (DI-mosoj)
- Treat repair, offer, counteroffer, acceptance, routing, introduction, redemption, transfer, storage, computation, TCP-link changes, authorization, dispatch, grant, registration, and enforcement as pCID-defined payload semantics, local trust/evidence interpretation, or implementation-local mechanics unless a scoped TE/DI proves a distinct wire-level role. (DI-mosoj)
- Before adding any new top-level action kind, stop and answer: why is this not just an agent's voluntary promise with pCID-defined payload meaning? If that answer is not explicit in a locked TE/DI, do not add the action kind. (DI-mosoj)
- POC7 through POC10 action names, including refusal and observation labels, are historical executable evidence, not the forward naming pattern. Do not rewrite scored/generated artifacts in place; supersede or reframe them when reused. (DI-mosoj; supersedes DI-fitav)

#### Contrasts and unresolved diffs
- **Same minus history** ([`ciwg/grid-examples`](../../../docs/other_repo_agents/ciwg_grid-examples_refs_heads_main_AGENTS.md#promise-action-minimalism-required)): identical rules without the POC7–10 historical clause.
- **Absent** ([`promisegrid/promisegrid`](../../../docs/other_repo_agents/promisegrid_promisegrid_refs_heads_main_AGENTS.md)): the parent organization guide carries none of this.
- **Diff to resolve:** domain content. Proposed disposition: `orgs/promisegrid/` (or a specification target), never the general library. No generalization attempt without a PromiseGrid-domain owner.

## POC Superset Discipline (Required)
- Future PromiseGrid POCs must be supersets of the previous POC's implemented behavior, architecture lessons, analyzer gates, and documented acceptance criteria unless a scoped DI explicitly declares the new POC non-superset and lists every intentionally dropped feature. (DI-sinur)
- A new POC may add focus, specialization, or new protocol surfaces, but it must not silently regress app/kernel boundaries, local trust semantics, monitor/analyzer gates, pCID routing, Promise Theory vocabulary, or previously proven workflows. (DI-sinur)
- When repairing or extending a POC, the analyzer must include inherited regression gates for the prior POC lineage, and the README or run narrative must state whether the POC is a superset or cite the DI authorizing an exception. (DI-sinur)

#### Contrasts and unresolved diffs
- **Unique to anchor:** no other source guide states a superset rule. The only structural analogue is the TE Editing Policy's supersedence discipline (nothing is rewritten in place; changes supersede).
- **Generated evidence** ([`RoSE`](<../../../docs/other_repo_agents/low qual/RoSE_agents.md>#danger-zones--human-review-required-not-just-looks-correct>): a generated guide reports the same failure mode from practice — "a large external rewrite was found missing [crash recovery] entirely; don't let a 'cleaner' rewrite silently drop it again." This independently supports a general no-silent-regression rule.
- **Diff to resolve — D11:** keep superset discipline PromiseGrid-local, or extract the general principle (rewrites must enumerate and preserve proven behavior and inherited regression gates; exceptions must be explicit) as a canonical module?

## DEV-GUIDE-RESOURCES.md

- Update DEV-GUIDE-RESOURCES.md file when a cited or relevant DR, DI, TE, or TODO is updated, when a cited draft spec freezes, or when the PromiseGrid Development Guide (https://github.com/ciwg/promisegrid-dev-guide) settles prose that supersedes a wire-lab note. Source: `DI-nunut`.
- Whenever `DEV-GUIDE-RESOURCES.md` is updated, also regenerate its top
  `## Current Design State` section by Codex LLM analysis of the current
  simulations and root `results/` evidence. The section must read like a
  concise protocol design document for PromiseGrid kernel and app-guide
  developers, not a wall of prose. It must derive examples from consensus
  across near-contender simulations rather than from only the top score, prefer
  simple positional pCID-selected wire shapes and deterministic pCID-defined
  payload contracts, use Burgess Promise Theory vocabulary, preserve
  small-device and 100-year durability constraints, avoid capability-shopping
  maps or large general-purpose claim-card maps as the apparent consensus, call
  out weak or missing coverage, and must not present simulations, proposal
  children, or score results as final PromiseGrid APIs. Source: `DI-baral`.

#### Contrasts and unresolved diffs
- **Unique to anchor:** no other guide maintains a generated design-state document.
- **Diff to resolve — D12:** disposition is clearly `repository-local` for wire-lab, but the underlying pattern (a generated summary derived from *consensus across near-contender evidence*, explicitly non-normative, with weak coverage called out) may be a reusable `generated-doc` pattern for v2. Decide whether to extract the pattern or defer entirely.

## Decision-First Specification and Compliance Protocol (Required)
- Decision-first means decisions must be locked before coding; it does not forbid pre-decision analysis such as required thought experiments.
- The agent must collect and lock user decisions before making any code edits for a task.
- Locked decisions must be recorded as Decision Intent Log entries in the relevant `protocols/<slug>.d/TODO/TODO-<handle>-<slug>.md` file(s) with clear intent and rationale.
- The agent must ask decision questions up front in a single intake round whenever possible.
- Required decision categories are architecture, design/behavior, implementation approach, function naming, variable naming, and file/path decisions.
- The agent must ask these as multiple-choice questions whenever practical.
- When a thought experiment (TE) is required, the agent must complete the TE before asking final DF questions. TEs narrow alternatives; DF questions and answers lock the decision before implementation.
- Thought experiments (TEs) are analysis artifacts; Decision Intent (DI) entries are the separate records that capture the locked decision after DF is resolved.

#### Contrasts and unresolved diffs
- **Same** (verbatim family with path and legacy-notation variants): [`promisegrid/promisegrid`](../../../docs/other_repo_agents/promisegrid_promisegrid_refs_heads_main_AGENTS.md#decision-first-specification-and-compliance-protocol-required), [`ciwg/FAB26-Presentation`](../../../docs/other_repo_agents/ciwg_FAB26-Presentation_refs_heads_main_AGENTS.md#decision-first-specification-and-compliance-protocol-required), [`ciwg/grid-examples`](../../../docs/other_repo_agents/ciwg_grid-examples_refs_heads_main_AGENTS.md#decision-first-specification-and-compliance-protocol-required), [`ciwg/decomk-conf-cswg`](../../../docs/other_repo_agents/ciwg_decomk-conf-cswg_refs_heads_main_AGENTS.md#decision-first-specification-and-compliance-protocol-required), [`ciwg/mob-sandbox`](../../../docs/other_repo_agents/ciwg_mob-sandbox_refs_heads_main_AGENTS.md#decision-first-specification-and-compliance-protocol-required), [`stevegt/navlog`](../../../docs/other_repo_agents/stevegt_navlog_refs_heads_main_AGENTS.md#decision-first-specification-and-compliance-protocol-required), [`stevegt/decomk`](../../../docs/other_repo_agents/stevegt_decomk_refs_heads_main_AGENTS.md#decision-first-specification-and-compliance-protocol-required).
- **Lighter** ([`computerscienceiscool/pg`](../../../docs/other_repo_agents/computerscienceiscool_pg_refs_heads_main_AGENTS.md#decision-first-specification-and-compliance-protocol)): "Collect and lock user decisions before behavior-changing code edits" plus the append-only DI log — no DF intake round, no multiple-choice requirement, no TE trigger.
- **Absent** (no decision protocol at all): [`ciwg/cswg`](../../../docs/other_repo_agents/ciwg_cswg_refs_heads_main_AGENTS.md), [`promisegrid/grid-poc`](../../../docs/other_repo_agents/promisegrid_grid-poc_refs_heads_main_AGENTS.md), [`stevegt/grokker`](../../../docs/other_repo_agents/stevegt_grokker_refs_heads_main_AGENTS.md), [`stevegt/godecide`](../../../docs/other_repo_agents/stevegt_godecide_refs_heads_main_AGENTS.md), [`stevegt/mob-consensus`](../../../docs/other_repo_agents/stevegt_mob-consensus_refs_heads_main_AGENTS.md), [`computerscienceiscool/llm-runtime`](../../../docs/other_repo_agents/computerscienceiscool_llm-runtime_refs_heads_audit-sweep_AGENTS.md). These guides substitute ordinary review/PR gates; none of them is "wrong" — they are lower-ceremony repos.
- **Diff to resolve — D3 (the strongest policy question):** is the canonical v2 workflow (a) strict decision-first for every repo, (b) risk-based escalation with strict decision-first as a selectable overlay, or (c) lightweight lock-before-behavior-change? Existing intake candidates `intake/workflow/use-risk-to-choose-routine-work-or-decision-review.md` and `intake/process/decision-governance/lock-durable-decisions-before-implementation.md` already frame (b) vs (a); the corpus evidence above shows real repos operating happily at all three levels, which supports making strictness a choice rather than a baseline.
- **Current cdint-grid conflict** ([`newest_promisegrid`](../../../docs/other_repo_agents/newest_promisegrid.md#decision-first-specification-and-compliance-protocol-required)): final decision questions are now required one at a time, while the anchor says to ask them up front in one intake round whenever possible. This is a direct workflow conflict rather than a repository-path variant. **D15 — choose whether question batching is canonical, selectable, or left to repository policy.**

## Thought Experiment Protocol (Required)
- Before locking any non-trivial decision that will require DF questions and answers, the agent must run a thought experiment (TE) if multiple plausible designs remain.
- A TE happens before final DF questions. Its purpose is to narrow the design space so DF questions and answers are informed by explicit scenario analysis.
- The agent must not collapse a TE into a short opinion or recommendation. The agent must explicitly model concrete scenarios and consequences.
- Each new TE must have a unique proquint handle in the format `TE-<handle>`, where `<handle>` is minted by `tools/mint-handle` from the global TODO/TE/DR/DI handle namespace. Existing pre-upgrade TE handles and prior aliases remain historical records.
- The TE doc filename must be `TE-<handle>-<slug>.md` and live under `docs/thought-experiments/`, for example: `docs/thought-experiments/TE-mumuv-naming-reconciliation.md`. The slug is informational and may be edited; the proquint handle is permanent.

### TE Intake Requirements
- Before locking decisions or asking final DF questions, the agent must identify:
  - the decision being tested,
  - the candidate alternatives,
  - the assumptions and threat/trust model,
  - the scope and systems affected.
- If the TE relates to an existing TODO, the agent must reference the TODO handle and subtask handle (for example, `fonuz.1`).

### TE Execution Requirements
- Each TE must evaluate the same decision across multiple concrete scenarios.
- Scenarios must include, when relevant:
  - normal operation,
  - failure/corruption/incomplete writes,
  - concurrent actors or mixed-version nodes,
  - long-horizon evolution and migration,
  - trust-boundary changes,
  - scale effects (storage, bandwidth, CPU, operational complexity).
- The agent must compare alternatives under the same assumptions instead of switching assumptions mid-analysis.
- The agent must state what each alternative makes easier, what it makes harder, and what new obligations it creates.

### TE Authoring Conventions
- Named actors follow the cryptography-literature alphabetical convention. Use Alice, Bob, Carol, Dave, Ellen, Frank, and so on for cooperative actors; use Mallory for adversaries; name Steve explicitly only when his repo-owner role is load-bearing in the scenario. Use this convention in TEs, scenario analyses, tabletop simulations, DR/DI prose, and worked examples in specs or docs when named actors are useful. Do not invent ad-hoc names when the convention fits. (DI-034-20260508-060134)

### TE Output to DF
- After the TE, the agent must identify:
  - rejected alternatives,
  - surviving alternatives,
  - unresolved questions that still require user choice,
  - any new naming/path/runtime decisions exposed by the TE.
- Final DF questions must be framed from the surviving alternatives identified by the TE. The agent must not ask broad DF questions that ignore TE results.

### TE Artifacts
- The agent must track required TEs in the relevant `protocols/<slug>.d/TODO/TODO-<handle>-<slug>.md`.
- For each completed TE, the agent must write a verbatim copy of the thought experiment into a standalone file under `docs/thought-experiments/`.
- The doc filename must begin with the TE ID and then use a descriptive suffix.
- The doc must stand on its own and include:
  - title,
  - TE ID,
  - decision under test,
  - assumptions,
  - alternatives,
  - scenario analysis,
  - conclusions,
  - implications for the repo's open TODOs and pending DIs.

### TE Decision Rules
- A TE does not by itself lock a decision.
- After the TE, the agent must either:
  - ask the user to choose among the surviving alternatives, or
  - recommend one surviving alternative and clearly state why the others were rejected.
- After user choice is resolved, the agent must record the locked result via the existing DI process before implementation.
- If a TE exposes a new ambiguity, dependency, or naming/path decision, the agent must stop and resolve that before implementation.

### TE Final Handoff Requirements
- In the final response for TE work, the agent must include:
  - which TE was completed,
  - the TE ID,
  - the doc path under `docs/thought-experiments/`,
  - the surviving alternatives,
  - the recommended conclusion or the exact DF question that remains for the user.
- Hard gate: for decisions that require a TE, work is incomplete until the TE doc exists and the resulting decision status is explicit (`needs DF`, `locked`, or `deferred`).

### TE Editing Policy (Required)

Once a TE is filed in `docs/thought-experiments/`, edits to it follow a categorized policy. The policy is locked in `DI-020-20260502-213103` (categorized editing regimes), `DI-020-20260502-213104` (uniform applicability across all TE corpora; this rule applies wherever TEs are stored, not just under `docs/thought-experiments/`), and `DI-020-20260502-213105` (holistic reading by default for substantive questions, single-TE reading allowed for obviously mechanical ones). The Cat-1 clause of `DI-020-20260502-213103` was superseded on 2026-05-02 by `DI-020-20260502-232651` (Cat-1a / Cat-1b split). Four further Cat-3 navigational refinements are appended to TE-dabol (`docs/thought-experiments/TE-dabol-te-editing-policy-and-holistic-corpus.md`) and to TE-vudaf (`docs/thought-experiments/TE-vudaf-editing-policy-tabletop.md`); the agent must read both TEs and their `## Refinements` sections before performing any TE edit.

The seven categories are:

- **Cat-1a (current-pointer paths).** A path reference that names the current location of a file. Mechanical sweep in place; no top-of-file note required.
- **Cat-1b (historical-quotation paths).** A path reference that quotes an earlier corpus state — inside a markdown blockquote, attributed to another TE ("TE-N states ..."), in past tense ("TE-magup used the path ..."), inside a `## Refinements` section, supersedence note, or `Decision status` line. Left untouched; rewriting would falsify the historical record. Per-match classification with five heuristics (quotation context; Refinements / supersedence framing; past tense; default Cat-1a; when-in-doubt-Cat-1b). Sweep tools may emit matches with surrounding context for human review but must not auto-rewrite.
- **Cat-2 (vocabulary updates).** A rename of a term whose meaning is unchanged (typo fixes, terminology consolidation). In place, with a top-of-file note pointing at the driving TE or TODO. The note must enumerate by ID every DI that lives in the affected TE, paired with an explicit promise that the rewrite preserves each DI's meaning. A TE without DIs gets a one-line `no DIs in this file` note. Form: `Cat-2 vocabulary update per <driving TE or TODO>: '<old term>' -> '<new term>'. The following DIs in this file are unchanged in meaning: DI-XXX-..., DI-YYY-..., DI-ZZZ-... .` Mandatory pre-step: grep the entire corpus for the old term inside quotation contexts (markdown blockquotes; fenced code blocks presented as citations; single/double-quoted phrases attributed to another TE via `TE-N states`, `TE-N reads`, `originally said`, `as of TE-N`, `the corpus showed`); each match is classified Cat-2 (sweep) or Cat-2-historical (leave) per the same heuristics as Cat-1a/Cat-1b.
- **Cat-3 (navigational forward pointers).** Append a dated entry to the TE's `## Refinements` section (created if absent, placed after `## Decision status`) describing where the affected reader should now look. The TE body above is unchanged. No DI is filed for a Cat-3 entry. Procedural tightenings of an existing category's how-to are Cat-3.
- **Cat-4 (resolved-implication forward pointers).** Same shape as Cat-3, used when an item from the TE's `Implications and future work` list has resolved (a TODO filed; a DR opened; a downstream TE landed). Append-only; no body edit.
- **Cat-5 / Cat-6 / Cat-7 (substantive supersedence).** A material change to a locked DI's meaning, scope, or applicability requires a new TE that supersedes the affected one. The new TE carries its own DFs and DIs; the older TE's `## Decision status` is updated to `superseded by TE-<id>` and its top-of-file `## Status` field is updated to `superseded by TE-<id> / DI-<id>`. The older TE's body is otherwise untouched.

Every TE in the corpus carries a top-of-file `## Status` field placed immediately after the TE ID line. Canonical values: `needs DF`, `decided`, `decided, refined`, `superseded by TE-<id> / DI-<id>`, `withdrawn`. Legacy values preserved during retrofit: `stub`, `open`, `recommended for immediate adoption`, `locked for the <protocol>`. New TEs prefer canonical values; the field is updated by Cat-1a sweep when the TE's state changes.

The `## Refinements` section is the single append-only home for Cat-3 / Cat-4 entries on a TE. Entries are dated (`### YYYY-MM-DD — <title>`) and ordered chronologically. The body of the TE above the `## Refinements` section is treated as historical evidence: a Cat-1a path-rename or Cat-2 vocabulary sweep on the body is permitted under its category rules; a Cat-3 / Cat-4 forward-pointer is appended to `## Refinements` rather than rewriting the body; a Cat-5 / Cat-6 / Cat-7 substantive change is filed as a new superseding TE rather than as an edit. The four Cat-3 Refinements on TE-dabol are themselves examples of this shape: each one tightens a category's procedure without changing the locked policy, and is filed as a Refinement entry rather than as a new DI.

Reading default: holistic. When deciding whether an edit is mechanical or substantive, when interpreting a single TE's claims, or when reasoning about whether a refinement is Cat-3 or Cat-5–7, the agent must read the corpus holistically (the affected TE plus the corpus's editing-policy chain: TE-dabol, TE-vudaf, and any other TEs they cite or that cite them). Single-TE reading is reserved for obviously mechanical questions (a single typo; a path that has demonstrably moved; a Status field retrofit) and only after the holistic read has confirmed the question is mechanical. When in doubt, read holistically.

Applicability: this policy applies uniformly to every TE corpus in this repository, regardless of which protocol or harness directory it lives in. Per-protocol corpora may add stricter rules but may not relax these rules.

### Naming Decisions (Required)
- The agent must not invent function names or variable names that are not already covered by locked naming decisions.
- If naming is not covered, the agent must stop and ask multiple-choice naming options before continuing.

### File/Path Decisions (Required)
- Path approvals are mandatory for all touched paths:
  - repo-changed files (create/rename/move/delete),
  - runtime touched paths (read/write/delete), including input files, output files, DB files, caches, fixtures, and temporary test files.
- The agent must ask path approvals one path at a time via multiple-choice questions.
- Path-question order must be dependency order.
- Each path question must include: action, exact path (or approved dynamic pattern ID), purpose, class (`prod-code | prod-data | test | temp`), and lifecycle intent.
- Temporary test paths require explicit approval and an explicit cleanup plan before handoff.
- Dynamic/runtime-generated paths must be approved by pattern, with:
  - allowed root bounds,
  - allowed actions,
  - concrete examples.
- The agent must ask one multiple-choice approval per dynamic path pattern.
- If any unapproved runtime path appears, the agent must stop and ask before continuing.

### Decision Lock and Stop Rule
- The agent must produce a Decision Lock summary with decision IDs before code edits begin.
- The agent must not proceed if any required decision is missing, ambiguous, or conflicting.
- The agent must stop and ask immediately if a new decision need appears during implementation.
- The agent must not assume defaults for locked categories unless the user explicitly approves defaults.

### Compliance Ownership (Agent)
- The agent must treat user decisions as authoritative and implement to those decisions.
- The agent must run a compliance self-review before finalizing and must fix all non-compliance before handoff.
- Hard gate: work is incomplete until compliance is PASS, or the user explicitly approves an exception.
- The user should not need to manually inspect diffs to determine compliance.

### Required final handoff artifacts
- `Decision Compliance: PASS/FAIL`
- Decision Matrix mapping each locked decision ID to implementation evidence.
- Inline diff annotations in the form `path:line -> decision_id -> rationale`.
- Runtime Path Touch Matrix listing each approved runtime path/pattern, action used, and where it is implemented/validated.
- `Exceptions:` listing only user-approved deviations.
- Every non-trivial behavior change must include intent provenance per existing DI requirements.

#### Contrasts and unresolved diffs
- **Same** (verbatim family): the intake/execution/output/decision/handoff subsections repeat in [`promisegrid/promisegrid`](../../../docs/other_repo_agents/promisegrid_promisegrid_refs_heads_main_AGENTS.md#thought-experiment-protocol-required), [`ciwg/FAB26-Presentation`](../../../docs/other_repo_agents/ciwg_FAB26-Presentation_refs_heads_main_AGENTS.md#thought-experiment-protocol-required), [`ciwg/grid-examples`](../../../docs/other_repo_agents/ciwg_grid-examples_refs_heads_main_AGENTS.md#thought-experiment-protocol-required), [`ciwg/decomk-conf-cswg`](../../../docs/other_repo_agents/ciwg_decomk-conf-cswg_refs_heads_main_AGENTS.md#thought-experiment-protocol-required), [`ciwg/mob-sandbox`](../../../docs/other_repo_agents/ciwg_mob-sandbox_refs_heads_main_AGENTS.md#thought-experiment-protocol-required), [`stevegt/navlog`](../../../docs/other_repo_agents/stevegt_navlog_refs_heads_main_AGENTS.md#thought-experiment-protocol-required), and [`stevegt/decomk`](../../../docs/other_repo_agents/stevegt_decomk_refs_heads_main_AGENTS.md#thought-experiment-protocol-required). Two historical ID variants appear: timestamp IDs `TE-YYYYMMDD-HHMMSS` (mob-sandbox, navlog) and proquint IDs (the rest) — this is the corpus's own evidence for a handle-migration policy.
- **TE Authoring Conventions** — the anchor's cryptography-alphabet convention has been replaced in current cdint-grid by stable names whose initials suggest explicitly stated roles, with Mallory retained for malicious actors. Related **Missing from anchor**: [`ciwg/FAB26-Presentation`](../../../docs/other_repo_agents/ciwg_FAB26-Presentation_refs_heads_main_AGENTS.md#public-prose-style) extends Alice/Bob conventions into public prose style. **D16 — choose a reusable actor-naming rule or classify both conventions as selectable style modules.**
- **TE Editing Policy** — three granularities exist in the corpus: the full seven-category regime (anchor, grid-examples), a three-line durable-records rule ("do not rewrite filed TE history for style cleanup; use `## Refinements` for navigational updates; use a superseding TE for material changes" — promisegrid, FAB26), and none (everywhere else). **Diff to resolve — D4:** the three-line version is a strong candidate for the canonical v2 module, with the seven-category system as a strict overlay; the anchor itself demonstrates the cost of the full regime.
- **Hard gates** (TE handoff gate, Compliance PASS gate, Runtime Path Touch Matrix): present verbatim across the whole governance family, absent elsewhere. **Generated evidence** ([`RoSE`](<../../../docs/other_repo_agents/low qual/RoSE_agents.md>#danger-zones--human-review-required-not-just-looks-correct)) reaches the same idea from the safety side: some diffs need named human review, "not just 'looks correct'". Existing intake candidate `intake/workflow/require-human-review-for-high-consequence-changes.md` covers the generalized form. **Diff to resolve — D5:** keep all hard gates canonical, or reserve hard gates for defined high-consequence classes and make the ceremony (matrix artifacts, inline annotations) optional per organization?
- **File/Path Decisions** — no other guide has per-path approval intake; [`RoSE`](<../../../docs/other_repo_agents/low qual/RoSE_agents.md>#participant-data) instead names standing no-go zones (participant data paths) that require stopping. The pattern generalizes as "declare protected path classes; runtime surprises stop work," which is closer to Mogent's own source-root safety rules than per-path Q&A. Proposed: extract the generalized boundary, keep per-path approval as a strict overlay.

## Coding Style & Naming Conventions
- Use object-oriented design with structs and methods; avoid large functions and global state.
- Follow generally accepted object oriented design patterns.
- Keep Go code `gofmt`-clean; package names should be short and lower-case.
- Prefer focused edits over broad refactors unless required.
- Add and maintain explanatory comments for non-obvious logic.
- Use `git mv` for file moves/renames to preserve history.

#### Contrasts and unresolved diffs
- **Stricter** ([`ciwg/decomk-conf-cswg`](../../../docs/other_repo_agents/ciwg_decomk-conf-cswg_refs_heads_main_AGENTS.md#coding-style-quality--naming-conventions), [`ciwg/mob-sandbox`](../../../docs/other_repo_agents/ciwg_mob-sandbox_refs_heads_main_AGENTS.md#coding-style-quality--naming-conventions), [`stevegt/navlog`](../../../docs/other_repo_agents/stevegt_navlog_refs_heads_main_AGENTS.md#coding-style-quality--naming-conventions)): add mandatory detailed plain-English comments, duplication-detection refactoring, a dependency allowlist (`stevegt/*, ciwg/*, promisegrid/*, cdint/*, t7a/*`), and standing orders to scan for code smells, bad architecture, and bad structure.
- **Diff to resolve — D6:** proactive-quality mandates ("always look for smells and recommend improvements") conflict in spirit with the anchor's own "prefer focused edits over broad refactors unless required" and with Diff Discipline in the promisegrid/FAB26 guides. The corpus holds both policies in different repos. This is the same tension already recorded in `intake/workflow/keep-changes-scoped-and-preserve-existing-work.md` (review question 1) — resolve together.
- **Diff to resolve — D7:** dependency policy appears only in the decomk-family guides (allowlist + "prefer standard library"); the anchor is silent; the generated [`todo_app`](../../../docs/other_repo_agents/tbd/todo_app_AGENTS.md) guide independently argues dependency restraint. Decide whether a dependency-restraint module is canonical, allowlisted per organization, or deferred.
- **Same** (partial): table-driven tests encouraged ([`promisegrid/grid-poc`](../../../docs/other_repo_agents/promisegrid_grid-poc_refs_heads_main_AGENTS.md#coding-style--naming-conventions), decomk family); `gofmt` + `go vet` ([`computerscienceiscool/llm-runtime`](../../../docs/other_repo_agents/computerscienceiscool_llm-runtime_refs_heads_audit-sweep_AGENTS.md#coding-style--naming-conventions)); `git mv` for renames (universal in captured guides).

## Error Handling Policy (Required)
- Never use `|| true` in scripts, templates, or make recipes. Always inspect
  command exit codes explicitly with `if/else` branches and handle each outcome.
- For non-fatal cleanup/diagnostics steps, record command status (exit code and
  logs) explicitly; do not fail silently.
- In Go code, never ignore errors with `_ = ...`; handle, propagate, or report
  errors explicitly.
- Run `errcheck ./...` and keep it passing for Go changes.

#### Contrasts and unresolved diffs
- **Same** — this is the strongest cross-repo consensus in the entire corpus, verbatim in [`promisegrid/promisegrid`](../../../docs/other_repo_agents/promisegrid_promisegrid_refs_heads_main_AGENTS.md#error-handling-policy-required), [`ciwg/FAB26-Presentation`](../../../docs/other_repo_agents/ciwg_FAB26-Presentation_refs_heads_main_AGENTS.md#error-handling-policy-required), [`ciwg/grid-examples`](../../../docs/other_repo_agents/ciwg_grid-examples_refs_heads_main_AGENTS.md#error-handling-policy-required), [`ciwg/decomk-conf-cswg`](../../../docs/other_repo_agents/ciwg_decomk-conf-cswg_refs_heads_main_AGENTS.md#error-handling-policy-required), [`ciwg/mob-sandbox`](../../../docs/other_repo_agents/ciwg_mob-sandbox_refs_heads_main_AGENTS.md#error-handling-policy-required), [`stevegt/navlog`](../../../docs/other_repo_agents/stevegt_navlog_refs_heads_main_AGENTS.md#error-handling-policy-required), and [`stevegt/decomk`](../../../docs/other_repo_agents/stevegt_decomk_refs_heads_main_AGENTS.md#error-handling-policy-required); the generic "explicit error handling" clause also appears in [`computerscienceiscool/pg`](../../../docs/other_repo_agents/computerscienceiscool_pg_refs_heads_main_AGENTS.md#coding-style).
- **No conflict exists.** Proposed: promote this section to one of the first canonical `orgs/cdint/engineering/` modules during the exemplar cutover; it needs no editing beyond genericizing the wording.

## DR/DI Source-of-Truth Protocol (Required)
- In this repo, DR and DI logs are the primary source of truth for decisions and open questions.
- Documents and code are outputs of that process and must link back to DR/DI records.
- Person identity in DR/DI records must use full email with label format: `user@example.com (FirstName)`.
- In DRs, `Asked by` and person-valued `Waiting on` fields must use that format.
- In DIs, `Author` must use that format.
- A DI `Author` is the decision-maker, not merely the recorder. If an agent records a DI for Steve's chat decision, Steve is the `Author`; an agent may be `Author` only when Steve explicitly delegates that decision authority to that agent. (DI-034-20260508-060134)
- A settled statement in docs (or critical logic in code comments) must cite at least one DI ID.
- An unresolved question or uncertainty must cite at least one DR ID.
- If an unresolved question has no DR yet, create a DR before finalizing the change.
- During TODO-kugod migration, apply these rules incrementally as sections/files are brought under DR/DI tracking.
- Intent: New coordination artifacts use a single proquint handle namespace so TODO, TE, DR, and DI references do not depend on timestamp or integer allocation. Source: DI-nisam

#### Contrasts and unresolved diffs
- **Variant** ([`promisegrid/promisegrid`](../../../docs/other_repo_agents/promisegrid_promisegrid_refs_heads_main_AGENTS.md#drdi-source-of-truth-protocol-required), [`ciwg/FAB26-Presentation`](../../../docs/other_repo_agents/ciwg_FAB26-Presentation_refs_heads_main_AGENTS.md#drdi-source-of-truth-protocol-required)): same rules, but author identity is resolved mechanically — "See the git config for this repo to detect the author's name" — instead of naming Steve. **Diff to resolve:** the anchor hard-codes a person; the variant generalizes. The generalized form is the better module; the explicit-delegation rule is the portable part.
- **Stricter** ([`computerscienceiscool/pg`](../../../docs/other_repo_agents/computerscienceiscool_pg_refs_heads_main_AGENTS.md#proquint-coordination-ids)): "`tools/mint-handle` is the only supported minting path" with working-tree scanning — the anchor describes minting but does not forbid alternatives.
- **Variant** ([`stevegt/decomk`](../../../docs/other_repo_agents/stevegt_decomk_refs_heads_main_AGENTS.md#comment-preservation-protocol-required)): a DI handle "must be owned by a Decision Intent Log entry in the relevant TODO file" — handles are not free-floating.
- **Diff to resolve — D9:** the corpus solves artifact identity with minted proquint handles plus explicit legacy-retention rules ("Existing numeric TODO and DI records remain valid legacy records until migrated"). This is the same problem as Mogent's library identity and source-path moves (`libv2-proto/proposals/source-path-moves.md`): stable identity for referenced things across renames. Resolve the library-identity question with this evidence in view.

## Comment Preservation Protocol (Required)
- Never remove existing code comments unless they are replaced in the same patch by equal-or-better explanatory comments near the same logic.
- When rewriting or refactoring code, port old explanatory intent first, then improve wording.
- If a touched non-trivial code block has no comments, add explanatory comments.
- Do not treat shorter comments as better unless they preserve all important intent.
- For any non-trivial behavior change, include a behavior-level comment with:
  - `Intent:` a short, clear rationale (a sentence or a few; no hard cap if more is needed for clarity).
  - `Source:` a DI ID in the format `DI-<handle>`.
  - `<handle>` is minted by `tools/mint-handle` and is globally unique across TODO, TE, DR, and DI owners.
  - Optional: TODO file/section reference for faster lookup.
- If a comment must be dropped with no replacement, stop and ask the user before proceeding.
- Before editing a file, review existing comments in that file.
- Maintain a `## Decision Intent Log` at the top of relevant `protocols/<slug>.d/TODO/TODO-<handle>-<slug>.md` files.
- Treat DI logs as append-only history. Do not rewrite or delete prior entries.
- When intent evolves, add a new DI entry and set `Supersedes: <old-di-id>`.
- DI entries must include:
  - `ID: DI-<handle>`
  - `Date: YYYY-MM-DD HH:MM:SS`
  - `Status: active|superseded`
  - `Decision:`
  - `Intent:`
  - `Constraints:`
  - `Affects:`
  - `Supersedes:` (optional)
- After editing, run a comment-delta audit on each touched code file using: `git diff -U0 -- <file> | rg -n '^-\\s*//|^-\\s*/\\*|^\\+\\s*//|^\\+\\s*/\\*'`.
- Resolve all removed-comment lines before finalizing unless explicit user approval was given.
- In the final response, include:
  - `Comment audit: PASS/FAIL`, with file list.
  - `Intent provenance audit: PASS/FAIL`, listing files with behavior changes and DI sources.
- Hard gate: behavior-changing work is incomplete unless comments preserve intent and include DI provenance.
- Do not remove comments or documentation; update them if outdated or incorrect.

### Comment + DI Examples
- Comment format example:
  - `// Intent: Keep context resolution stable across workspace scans to avoid target drift between plan and run. Source: DI-vapoj`
- Decision Intent Log entry template (for TODO files):
  - `ID: DI-<handle>`
  - `Date: YYYY-MM-DD HH:MM:SS`
  - `Status: active`
  - `Decision: <what was decided>`
  - `Intent: <short clear rationale>`
  - `Constraints: <hard limits, dependencies, assumptions>`
  - `Affects: <paths, modules, commands, docs>`
  - `Supersedes: <old DI ID, optional>`

#### Contrasts and unresolved diffs
- **Same** (verbatim family): [`promisegrid/promisegrid`](../../../docs/other_repo_agents/promisegrid_promisegrid_refs_heads_main_AGENTS.md#comment-preservation-protocol-required), [`ciwg/FAB26-Presentation`](../../../docs/other_repo_agents/ciwg_FAB26-Presentation_refs_heads_main_AGENTS.md#comment-preservation-protocol-required), [`ciwg/grid-examples`](../../../docs/other_repo_agents/ciwg_grid-examples_refs_heads_main_AGENTS.md#comment-preservation-protocol-required), [`ciwg/decomk-conf-cswg`](../../../docs/other_repo_agents/ciwg_decomk-conf-cswg_refs_heads_main_AGENTS.md#comment-preservation-protocol-required), [`ciwg/mob-sandbox`](../../../docs/other_repo_agents/ciwg_mob-sandbox_refs_heads_main_AGENTS.md#comment-preservation-protocol-required), [`stevegt/navlog`](../../../docs/other_repo_agents/stevegt_navlog_refs_heads_main_AGENTS.md#comment-preservation-protocol-required), [`stevegt/decomk`](../../../docs/other_repo_agents/stevegt_decomk_refs_heads_main_AGENTS.md#comment-preservation-protocol-required).
- **Lighter** ([`computerscienceiscool/pg`](../../../docs/other_repo_agents/computerscienceiscool_pg_refs_heads_main_AGENTS.md#comment-preservation-protocol)): same replacement rule and comment-delta audit, without the DI provenance fields or PASS/FAIL handoff artifacts.
- **Lightest** ([`ciwg/cswg`](../../../docs/other_repo_agents/ciwg_cswg_refs_heads_main_AGENTS.md#agent-specific-notes), [`stevegt/grokker`](../../../docs/other_repo_agents/stevegt_grokker_refs_heads_main_AGENTS.md#agent-specific-notes)): "Do not remove comments or documentation; update them if outdated."
- **Diff to resolve — D5 (shared):** the replacement-and-intent rule is universal and should be canonical; the DI-provenance requirement and PASS/FAIL audit artifacts are ceremony that only the governance repos run. Split the module: canonical comment-preservation, with intent-provenance as a required overlay where a decision-record system exists. Overlaps `intake/process/decision-governance/preserve-decision-history-through-supersession.md` (which already separates DI history from comment policy).

# DR Records

The DR/ directory stores Decision Request (DR) records for coordination work.

Rules:
- One DR per file.
- DR files are append-only event logs.
- Keep TODO files as snapshots; link TODOs to DR files for open questions.
- Person identity format: `user@example.com (FirstName)`.

Recommended file naming:
- `DR-<handle>-<slug>.md`, where `<handle>` is minted by `tools/mint-handle`.

Required DR fields:
- `DR-ID`
- `Date`
- `Asked by` (person identity format above)
- `State` (`open | decided | blocked | implemented | closed`)
- `Question`
- `Why this blocks progress`
- `Affects` (repos/files/components)
- `Unblocks` (TODO IDs/tasks)
- `Waiting on` (person identity format above, or DI ID)
- `Decision` (filled when decided)
- `Linked DI`
- `Related commits`
- `Last updated`

Reference pattern:
- From TODO files: `../DR/<filename>.md`

#### Contrasts and unresolved diffs
- **Same** (verbatim family): [`promisegrid/promisegrid`](../../../docs/other_repo_agents/promisegrid_promisegrid_refs_heads_main_AGENTS.md#dr-records), [`ciwg/FAB26-Presentation`](../../../docs/other_repo_agents/ciwg_FAB26-Presentation_refs_heads_main_AGENTS.md#dr-records), [`ciwg/grid-examples`](../../../docs/other_repo_agents/ciwg_grid-examples_refs_heads_main_AGENTS.md#dr-records).
- **Absent** (no DR records): every other guide. Open questions live in TODO files there instead.
- **Diff to resolve:** DR is the governed end of the handoff spectrum (structured, field-complete, blocking) while other guides resolve questions conversationally. This is exactly the "ordinary versus governed handoff boundary" decision in `libv2-proto/WORKLIST.md`; these four guides are the evidence base for the governed side. Note also `DR/DR-lusim-library-relationships.md` in this repo uses the same record shape.

## Testing Guidelines
- Use Go's standard `testing` package with deterministic tests.
- Avoid network calls in tests unless explicitly required and documented.
- When changing `plan/run` behavior, add coverage for both command paths when possible.

#### Contrasts and unresolved diffs
- **Same** (deterministic, offline, standard library): universal across Go guides. **Stricter** additions: table-driven tests encouraged ([`promisegrid/grid-poc`](../../../docs/other_repo_agents/promisegrid_grid-poc_refs_heads_main_AGENTS.md#testing-guidelines), decomk family); mocking preferred ([`ciwg/decomk-conf-cswg`](../../../docs/other_repo_agents/ciwg_decomk-conf-cswg_refs_heads_main_AGENTS.md#testing-guidelines)); test caches under `/tmp` ([`promisegrid/promisegrid`](../../../docs/other_repo_agents/promisegrid_promisegrid_refs_heads_main_AGENTS.md#testing-guidelines), [`ciwg/FAB26-Presentation`](../../../docs/other_repo_agents/ciwg_FAB26-Presentation_refs_heads_main_AGENTS.md#testing-guidelines)).
- **Variant** ([`openai/codex`](../../../docs/other_repo_agents/tbd/openai_codex_refs_heads_main_AGENTS.md#code-review-rules), authoritative upstream): "For agent changes prefer integration tests over unit tests" and features changing agent logic **must** add integration tests. Also imposes numeric review caps: "the total number of changed lines should not exceed 800 lines... under 500 lines" for complex logic, with split-into-stages guidance. **Diff to resolve:** the anchor has no change-size rule at all; its Diff Discipline is qualitative ("smallest coherent change"). Decide whether numeric caps are adopted, rejected, or offered as an optional strict overlay for review-heavy repos.
- **Diff to resolve:** proposed canonical extraction: "deterministic, offline tests; coverage alongside behavior changes; test caches outside the working tree." Everything else (integration-first, table-driven, watch loops like grid-poc's `devloop.sh`) is language- or repo-specific overlay.

## Commit & Pull Request Guidelines
- Treat a line containing only `commit` as: add and commit all changes with an AGENTS-compliant message.
- Use short, imperative, capitalized commit subjects.
- Summarize changes per file in commit bodies.
- Stage files explicitly (avoid `git add .` / `git add -A`).
- Do not open GitHub pull requests for normal wire-lab convergence. Steve explicitly dropped the require-PR merge rule; merge by the role-specific direct-push workflow instead. (DI-001-20260428-195702; DI-034-20260508-060134)
- Do not force-push repo branches unless a scoped DI/DR explicitly authorizes the exception. Role overlays may add stricter no-force-push rules for their branches or may document a narrow private-remote exception, but the repo-wide default is no history rewrites. (DI-034-20260508-060134)
- Do not commit local state files, generated binaries, credentials, tokens, signing keys, or other secrets. (DI-034-20260508-060134)

#### Contrasts and unresolved diffs
- **Variant** (PR policy spans the corpus): "Do not open GitHub pull requests unless the user explicitly asks" ([`promisegrid/promisegrid`](../../../docs/other_repo_agents/promisegrid_promisegrid_refs_heads_main_AGENTS.md#commit--pull-request-guidelines)) vs "PRs should include a concise summary, tests run, linked issues" (encouraged: [`ciwg/cswg`](../../../docs/other_repo_agents/ciwg_cswg_refs_heads_main_AGENTS.md#commit--pull-request-guidelines), [`stevegt/godecide`](../../../docs/other_repo_agents/stevegt_godecide_refs_heads_main_AGENTS.md#commit--pull-request-guidelines), [`promisegrid/grid-poc`](../../../docs/other_repo_agents/promisegrid_grid-poc_refs_heads_main_AGENTS.md#commit--pull-request-guidelines)) vs the anchor's org-locked direct-push. **Diff to resolve — D8:** merge/PR policy is organization policy, not universal guidance; it belongs in org-owned overlays, not the canonical module. (Matches TODO M-O2's explicit-vs-trusted-profile question.)
- **Missing from anchor** ([`ciwg/grid-examples`](../../../docs/other_repo_agents/ciwg_grid-examples_refs_heads_main_AGENTS.md#commit--pull-request-guidelines), decomk family): here-doc commit bodies ("use a here-doc (`git commit -F -`)... Do not use -m flags"). Minor tooling-hygiene rule; candidate for a tools/Git module.
- **Same** (universal): explicit staging (no `git add .` / `-A`), imperative capitalized subjects, per-file commit bodies, the `commit` keyword, and the no-secrets/no-state rule (extended with container/audit outputs in [`computerscienceiscool/llm-runtime`](../../../docs/other_repo_agents/computerscienceiscool_llm-runtime_refs_heads_audit-sweep_AGENTS.md#security--local-state)).
- **Variant** ([`stevegt/mob-consensus`](../../../docs/other_repo_agents/stevegt_mob-consensus_refs_heads_main_AGENTS.md#build-test-and-development-commands)): a tool that pushes after commits by default — direct-push taken to its mechanical conclusion; useful evidence that push policy is per-repo configuration.

## Glossary
- **TE**: Thought Experiment. Analysis doc under `docs/thought-experiments/TE-<handle>-<slug>.md` or an approved per-protocol TE corpus. The handle is a proquint minted by `tools/mint-handle`. (DI-034-20260508-060134)
- **DR**: Decision Request. Open question or decision-tracking record under `DR/DR-<handle>-<slug>.md`, where `<handle>` is minted by `tools/mint-handle`. (DI-nisam; DI-034-20260508-060134)
- **DI**: Decision Intent. Locked decision record inside a `## Decision Intent Log` in a protocol TODO file. New DI ID format is `DI-<handle>`, where `<handle>` is minted by `tools/mint-handle`. (DI-nisam; DI-034-20260508-060134)
- **DF**: Decision Framing. The multiple-choice intake round used to lock a decision after any required TE narrows the alternatives. (DI-034-20260508-060134)
- **TODO**: Task-tracking file under `protocols/<slug>.d/TODO/TODO-<handle>-<slug>.md`, where `<handle>` is minted by `tools/mint-handle`. The master cross-listed index is `protocols/wire-lab.d/TODO/TODO.md`; per-protocol queues live at `protocols/<slug>.d/TODO/TODO.md`. (DI-nisam; DI-034-20260508-060134)
- **twig**: Short kebab-case task name used in branch names, usually as `<user>/<twig>` such as `ppx/<twig>` or `stevegt/<twig>`. (DI-034-20260508-060134)
- **pCID**: Protocol CID. A pCID is the content hash of a spec document that defines a wire protocol; it is analogous to a TCP/UDP port number with no central registry because the spec hash is the port number. A pCID is not the hash of a particular message, payload, or promise body. (DI-009-20260429-173359; DI-034-20260508-060134)

#### Contrasts and unresolved diffs
- **Same** (process vocabulary): TE/DR/DI/DF/TODO definitions repeat across the governance family; twig exists only in wire-lab and grid-examples.
- **Diff to resolve:** split destinations. Process vocabulary (TE, DR, DI, DF) is canonical `agent-guide` material for any org adopting decision governance. Domain vocabulary (pCID, twig, Burgess references) is PromiseGrid-owned `specification`/`repository-local` material and must not leak into the general library.

## Current cdint-grid Skill Corpus

The current source pairs a 1,555-line canonical guide with 18 focused skills.
Each skill has YAML frontmatter (`name`, `description`) and an adjacent
`agents/openai.yaml` interface record (`display_name`, `short_description`,
`default_prompt`). The snapshot is intake evidence, not a claim that these
repository-specific procedures should be promoted unchanged.

- **Context, decisions, and design:** `load-project-context`, `decision-first-change`, `explain-architecture`, `run-thought-experiment`, and `review-against-ninik` load bounded evidence and preserve the distinction between analysis, recommendation, and locked authority.
- **Worker lifecycle and instruction distribution:** `newtree`, `resume`, `upgrade`, and `retire` compose stable-worker startup, finite work cycles, deterministic instruction synchronization, and evidence-preserving shutdown.
- **Worker communication:** `putq`, `getq`, and `drain-inbox` separate one voluntary send, one-message receipt, and one frozen-snapshot classification pass. Receipt and queue position never imply acceptance.
- **Git promises and peer convergence:** `commit` records one exact local promise per commit; `consensus` evaluates immutable peer commits without creating a global verdict or trust score.
- **Legacy compatibility:** `handoff` and `handback` isolate old packet-based attempts from the current Git-native workflow instead of burdening every current operation with both models.
- **Planning and operational state:** `update-dobab` loads a validated bounded deadline view, while `next10` reconciles evidence into ten priorities without mutating the queue.

#### Contrasts and unresolved diffs

- **Progressive disclosure:** the skill set moves detailed procedures out of the always-loaded path while retaining canonical authority in `AGENTS.md`. This is concrete evidence for Mogent's static-library versus runtime-loaded-skill distinction (M-G5), but it also shows that context reduction depends on deterministic selectors and disciplined cross-references.
- **Authority and capability remain distinct:** possessing a skill or executable does not authorize its represented action. User decisions and durable TODO/DR/DI records govern; skills route procedure; installed tools implement mechanics; Git objects and hashes carry evidence. A reusable library needs to preserve these distinctions rather than treating skill installation as permission.
- **Durable versus disposable state is explicit:** Git promises and coordination records are durable evidence; local consensus views, carrier timestamps, pane observations, and working trees are mutable or rebuildable. This distinction is useful candidate metadata, not merely cdint-grid workflow detail.
- **Finite-operation safety recurs across skills:** exact-once reads, frozen snapshots, no polling, no blind retries, immutable commit selection, and successor corrections bound races and accidental duplicate effects. These are potential reusable procedure primitives.
- **Portability is unresolved:** skills embed repository names, absolute `/home/stevegt/bin` commands, `/tmp/cdint-grid` layouts, named documents, and cdint-grid-specific DIs. Mogent must distinguish portable procedure templates from repository-bound instances before composing or installing them.
- **Metadata can drift from procedure text:** `upgrade` interface metadata says “review and adopt” although the skill requires deterministic synchronization without model review; `handoff` metadata does not foreground its legacy-only status; and `retire` metadata describes one worker while its procedure inventories all workers and seeks per-worker approval. Display capitalization and `$skill-name` use are also inconsistent. Discovery metadata therefore needs linting against normative content.
- **Missing structured fields:** the current interface metadata has no explicit version, authority class, dependencies, side effects, compatibility/deprecation state, required binaries, output schema, exact-message trigger, mutability, retention, or portability fields. Adding all of them blindly would overfit this source; D13 and D14 should decide the minimum model after cross-ecosystem research M-S6/M-S7.
- **Source-versus-runtime identity:** tracked skill text does not prove that a corresponding installed binary was built, qualified, installed, or cut over. Mogent's source model may need pins and provenance for authored skills while leaving executable attestation to an integration layer.
- **Diff to resolve:** decide how conflicting claims among canonical guide text, skill procedure, interface metadata, and installed executable behavior are detected and which layer wins. The source itself treats the guide and durable decisions as authority, but that repository-specific answer should not be silently universalized.

---

## Unresolved Diff Index

Diffs marked **D** are owner decisions. Everything else in the contrast blocks
is triage already settled by corpus evidence (labels Same/Variant/Stricter/
Looser with a clear majority).

- **D1 — Layout rule:** canonical "no `internal/`/`pkg/`" vs llm-runtime's deliberate `pkg/`+`internal/` layout. Decide canonical rule and whether overlays may relax it.
- **D2 — Planning-artifact location and IDs:** per-protocol TODO trees (wire-lab) vs root `TODO/` (majority); proquint handles vs zero-padded numbers; legacy-retention stances differ ("remain valid until migrated" vs "do not create new"). Decide canonical-vs-repository-local, and whether v2 needs a move-map for ID-scheme migrations.
- **D3 — Decision workflow baseline:** strict decision-first (8 guides) vs lock-only (cscool/pg) vs none+review gates (6 guides) vs risk-based escalation (v1 adaptation, pending). The anchor's own section is the strict candidate; corpus evidence supports strictness as a choice, not a universal default. Worklist item: "strict decision-first bundle vs independently selectable stages."
- **D4 — TE editing policy granularity:** full seven-category regime (2 guides) vs three-line durable-records rule (2 guides) vs none. Proposed canonical: three-line rule; proposed overlay: category regime.
- **D5 — Hard-gate ceremony:** Decision Compliance PASS/FAIL, comment audit artifacts, runtime-path matrices, per-path approvals — canonical, strict-overlay, or reserved for high-consequence classes? Generated evidence (RoSE danger zones) supports class-based human review; the governance family applies gates universally.
- **D6 — Coding-style posture:** restrained focused edits (anchor, promisegrid, FAB26) vs proactive quality mandates (decomk family). Conflicts in spirit; resolve with `intake/workflow/keep-changes-scoped-and-preserve-existing-work.md` review question 1.
- **D7 — Dependency policy:** allowlist + standard-library preference (decomk family) vs silence (anchor) vs restraint argument (generated todo_app guide).
- **D8 — PR vs direct-push:** organization policy, not universal guidance; evidence spans PR-encouraged (3 guides), user-gated (promisegrid), and org-locked direct-push (wire-lab, grid-examples).
- **D9 — Artifact identity:** proquint handles + legacy retention vs paths-as-identity (Mogent's model); informs the library-identity and source-path-move decisions.
- **D11 — Superset/no-silent-regression:** extract as a general rewrite-preservation rule, or keep protocol-local?
- **D12 — Consensus-derived generated summary:** extract the DEV-GUIDE-RESOURCES pattern as a reusable `generated-doc` module, or classify wholly repository-local?
- **D13 — Skill architecture layers:** model policy authority, operational procedure, discovery metadata, and executable tooling separately, or keep a simpler source-node model with conventions?
- **D14 — Skill composition contracts:** make dependencies, compatibility state, preconditions, side effects, exact-once rules, and output schemas machine-readable, or retain some/all in prose?
- **D15 — Decision-question cadence:** batch questions up front, require one question at a time, offer both as selectable policy, or leave cadence repository-local?
- **D16 — Scenario actor naming:** role-initial stable names vs cryptography-alphabet names; choose a reusable convention or expose style variants.

*(Numbering skips D8→D11 to leave room for splits during owner review.)*

## Disposition Sketch (Exemplar Cutover Triage)

Where each anchor section should land once its diffs are resolved. This is the
"integrate with org repos asap" map; everything else in `intake/` stays triaged
and deferred.

| Anchor section | Generality | Proposed destination |
|---|---|---|
| Project Structure & Module Organization | family + repo-local mix | rules → `orgs/<owner>/engineering/`; command/protocol specifics → repository-local |
| Build, Test, and Development Commands | repo-local | repository-local; extract "report environment blockers" as general |
| Agent Instruction Architecture | general | `library-authoring` (it describes Mogent's own model) |
| Current cdint-grid Skill Corpus | cross-ecosystem evidence + repo-local procedures | skills research M-S6/M-S7; architecture pending D13 + D14 |
| Promise Action Minimalism | domain | `orgs/promisegrid/` or specification target |
| POC Superset Discipline | generalizable kernel | pending D11 |
| DEV-GUIDE-RESOURCES.md | repo-local + pattern | repository-local; pattern candidate pending D12 |
| Decision-First Protocol | general, strictness-contested | pending D3; feeds `intake/process/decision-governance/` |
| Thought Experiment Protocol | general (minus editing-policy weight) | `intake/process/decision-governance/use-thought-experiments-to-narrow-a-broad-design-space.md` |
| TE Editing Policy | contested granularity | pending D4 |
| Naming / File-Path / Lock / Compliance / Handoff | family-strict | pending D3 + D5 |
| Coding Style & Naming Conventions | language-specific | `orgs/<owner>/go/`; proactive mandates pending D6 |
| Error Handling Policy | universal consensus | first canonical `orgs/cdint/engineering/` module |
| DR/DI Source-of-Truth Protocol | general with identity questions | pending D9 |
| Comment Preservation Protocol | universal core + ceremony | pending D5 split |
| DR Records | general (governed orgs) | pending handoff-boundary decision |
| Testing Guidelines | general core + overlays | `orgs/<owner>/go/` core module |
| Commit & Pull Request Guidelines | universal core + org policy | core → shared; PR policy → org overlays per D8 |
| Glossary | mixed | process vocabulary → shared; domain vocabulary → `orgs/promisegrid/` |

<!-- review-note: The disposition sketch is proposed, not approved. Promotion
into orgs/ still requires the owner review per TODO/PICKUP-2026-08-25 and
PROVENANCE.md updates per heading. -->
