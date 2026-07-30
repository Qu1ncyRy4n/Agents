# Mogent Module Ecosystem Research Memo

Status: research memo. This document records corpus observations and proposed
library work. It does not lock any new module selection, wording, or preset.

Scope: all files in `docs/other_repo_agents/` and the active mogent design
material. `docs/DESIGN.md` is authoritative where earlier material conflicts.

Decision context: DI-ralik establishes the manifest, named-source, and
content-role category model used below. DI-sufok establishes manifest-owned
headings, user-defined source aliases, and visible source provenance. Process
weight, domain, and language are selectable library content, not top-level
categories.

## Corpus Summary

The corpus contains four distinct genres. They should not be flattened into one
default prompt.

1. Engineering repository guides: concise project map, commands, tests, code
   conventions, artifact hygiene, and commit rules. `godecide`, `grokker`,
   `llm-runtime`, `color-engine`, and `pdf_md_extractor` are representative.
2. Decision-first governance guides: a shared CIWG/PromiseGrid family with
   decision records, TEs, provenance, comment preservation, and detailed
   handoff requirements. `promisegrid`, `wire-lab`, `grid-examples`, and
   `decomk` are representative.
3. High-consequence operational/research guides: experiment validity, human or
   scientific data, hardware/system activation, irreversible filesystem work,
   and external-service boundaries. `RoSE`, the eye-tracking projects,
   `photo_getter`, and the Nix assessments are representative.
4. Instructional and knowledge-work guides: learner model, teaching loop,
   durable session artifacts, Obsidian conventions, and named personas.
   `pg_learning` and the Notes vault documents are representative.

The first two genres supply broadly reusable text. The latter two supply
patterns and templates, but their project facts and personal state must remain
local overlays.

## Patterns That Repeat

### High-value universal content

- A short repository identity, source-of-truth map, and only the commands that
  an agent can actually run safely.
- Deterministic validation close to the change. Language-specific test guidance
  works best as a subtree, not a global claim.
- Explicit generated-artifact and secret exclusions.
- Minimal diffs, preservation of explanatory comments, explicit staging, and
  clear commit behavior.
- A boundary between repository-wide rules and role/environment overlays. The
  `grid-examples`/`wire-lab` statement of that boundary is especially useful.

### High-value conditional content

- Decision-first/TE/DR/DI rules are coherent and useful for protocol, public
  specification, and long-lived architecture work. They are too costly as the
  default for a typo, small repair, or early prototype. Keep them as a full
  governance subtree, alongside a lightweight decision-escalation subtree.
- Destructive-data and system-activation gates belong in dedicated safety
  modules. They should activate for archive recovery, migrations, lab rigs,
  Nix, deployment, and production data work.
- Tutor mode is valuable when its ownership boundary is explicit: agent writes
  scaffolding/tests and teaches; learner writes selected core logic. It must not
  include a specific learner's private abilities, session state, or persona.
- Research guides need data lineage, reproducibility, uncertainty, and
  calibration constraints. Specific subjects, hardware geometry, and study
  protocols are local project context.

### Common failures to avoid

- Copying a 200-300 line governance prompt into every repository. This creates
  drift, burns context, and hides the few rules relevant to small changes.
- Turning an assessment of a missing prompt into an AGENTS module unchanged.
  The Nix and SEAM files are useful analyses and skeletons, not project rules.
- Treating machine paths, local agent files, active lab state, or current
  learner state as shareable library material.
- Shipping stale configuration examples. The removed `.mogent` examples and
  `AGENTS.toml` use retired models; use the current `agents.yaml` manifest from
  `docs/DESIGN.md`.
- Encoding broad technology preferences as absolute rules. A language, package
  layout, or commit cadence is a stack/profile choice unless the repo has
  already locked it.

## Fit To Mogent Categories

| Content | Category | Library or overlay | Notes |
| --- | --- | --- | --- |
| Agent role, repo map, source documents, stack | Identity | local overlay plus small templates | Facts change per repo. |
| Workflow, testing, commits, TODO, decision process | Instructions | shared selectable subtrees | Most reusable corpus content. |
| Secret/data rules, destructive-action gates, artifact hygiene | Constraints | shared safety subtrees plus local paths | Keep high-risk rules explicit. |
| Language style, diff/comment rules, handoff format, glossary | Format | shared stack/profile subtrees | Do not force Go rules into Rust/Python repos. |
| Deliberate research, rapid repair, pedagogic loop | Cognition/process | optional profile library | A mode, not a category replacement. |
| Direct/TTS/Socratic communication and learner ownership | Communication/style | optional profile library | Separate voice from teaching method. |
| Go/Rust/Python/Nix/static-site practices | Code | stack libraries | Select with the relevant stack. |
| Obsidian, session logs, public-prose, decision-record conventions | Notes/docs | document-work libraries | Preserve privacy boundary. |

The four core categories remain a good compact render order. The richer
categories are useful roots only if they help browsing; they must not become a
second classification system for domains or process weight.

## Proposed Initial Ecosystem

Start with few, composable libraries. Each item below is a proposed heading
subtree, not a mandate to create a file per bullet.

### 1. `core/`

- `identity/repository-context-template`: templated role, repo name, source-of-
  truth pointers, and explicit project status.
- `instructions/focused-change-loop`: understand, scope, edit, validate, report.
- `instructions/commit-discipline`: explicit staging, imperative subjects, no
  autonomous PR/force-push unless selected.
- `constraints/secret-and-generated-artifacts`: no credentials, generated
  binaries, caches, or private runtime state.
- `format/minimal-diff-and-comment-preservation`: retain explanation without
  imposing an inappropriate language style.

This should be the smallest useful baseline. It should not include full
decision-first procedure by default.

### 2. `process/`

- `decision-first/lightweight-escalation`: ask only where a real design fork or
  irreversible behavior change exists.
- `decision-first/full-te-di-dr`: the mature PromiseGrid-style workflow,
  including provenance and compliance evidence.
- `change-review/delta-review`: review uncommitted changes and classify user
  edits without reverting them.
- `risk/high-consequence-change-gate`: human review for protocol, calibration,
  production, migration, or externally visible behavior changes.

The full process module should be split into independently selectable
subtrees. The current corpus repeatedly shows that a single all-or-nothing
governance block is too heavy for general use.

### 3. `code/`

- `go/cli-and-library`: `gofmt`, `go test`, deterministic fixtures, error
  wrapping, and module-local command execution.
- `rust/cargo-workspace`: `cargo fmt/check/test`, crate boundaries, feature and
  lockfile discipline. Take the concise project-specific parts from the Rust
  examples, not the whole OpenAI monorepo policy.
- `python/uv-cli`: `uv`, dry-run-first file processing, dependency-lock policy,
  and deterministic tests.
- `static-site/hugo`: generated-output discipline and preview/build validation.
- `shell/explicit-failure-handling`: applicable only where shell automation is
  in scope.

### 4. `safety/`

- `data-pipeline/non-destructive-archive`: mark-only deletion, source remains
  authoritative, serial database writers, temporary database validation.
- `research/human-data-and-irb`: never expose identified data, retain ignore
  rules, require review for experiment-validity changes.
- `research/scientific-reproducibility`: record environment, preserve raw data,
  distinguish known limitations from repairable errors, validate against task
  protocol before analysis claims.
- `systems/nix-nonactivating-validation`: no activation, hardware, boot,
  networking, firewall, or secrets work without explicit scope; prefer checks
  and builds.
- `external-service/credential-and-rate-limit`: secrets out of logs, dry runs,
  bounded live smoke tests, pacing.

### 5. `teacher/`

- `method/active-learning-loop`: prediction, observation, explanation,
  teach-back, and transfer.
- `session/sustainable-closeout`: time boundary, explicit thread disposition,
  precise resume pointer, and no unrequested homework.
- `communication/tts-friendly`: short spoken paragraphs, expanded uncommon
  abbreviations, visual material in an associated artifact.
- `ownership/learner-implements-core`: an explicit division between scaffolding
  and the learner's intended implementation work.
- `artifacts/lesson-session-template`: generic plan, board, worksheet, and log
  template without a learner profile.

This is a strong candidate for a `teacher_bot` preset: `core` + selected
`teacher` nodes + the applicable language/domain library.

### 6. `research/`

- `orientation/exploratory-stage`: confirm current stage before coding; label
  hypotheses separately from committed design.
- `analysis/reproducible-pipeline`: data inputs, transforms, output locations,
  environment boundaries, and fixture-first validation.
- `experiment/calibration-and-timing-review`: human review gate and update of
  authoritative procedure docs for behavior changes.
- `data/subject-and-media-exclusions`: raw data, traces, stimulus media, and
  derived sensitive outputs excluded from commits.

Use this for a `research_sci` preset. Lab-specific hardware, task names, subject
identifiers, and parameter values belong in the consuming repository's local
Identity/Constraints nodes.

### 7. `cdint/`

- `promisegrid/terminology`: precise promise, request, evidence, trust, CID,
  and pCID vocabulary.
- `promisegrid/action-minimalism`: the `promise` top-level semantic-action rule.
- `governance/decision-record-provenance`: the CDINT-specific TE/DI/DR details
  and public-artifact rules.
- `collaboration/repo-and-role-boundaries`: canonical AGENTS content versus
  overlays, plus cross-repository handoff ownership.
- `docs/rfc-like-public-prose`: only for public specs/papers, not general chat.

This library should select the full process profile deliberately, not cause it
to leak into unrelated applications.

### 8. `notes/` (private or carefully sanitized)

- `obsidian/portable-markdown-and-links`.
- `zk/atomic-note-and-processing-workflow`.
- `docs/decision-record-maintenance`.
- `privacy/knowledge-vault-boundary`.

Do not publish the personal vault's memory schema, identity/persona parameters,
or biographical context as shared modules. A private include can retain them if
they are genuinely desired.

## Proposed Source Boundaries

The library tree above describes content. The following is a separate,
provisional grouping of that content into source libraries. It is not a locked
taxonomy; the open boundary question is tracked in `DR-garom`.

```text
sources/
|-- cdint/
|   |-- shared-baseline
|   |-- process
|   |-- go
|   |-- promisegrid
|   `-- public-prose
|
|-- ucd_research/
|   |-- research-orientation
|   |-- research-safety
|   |-- human-data-and-irb
|   |-- reproducible-analysis
|   `-- python/
|       |-- dependency-management/
|       |   `-- uv
|       |-- cli-and-file-processing
|       `-- scientific-analysis
|
|-- personal/
|   |-- teacher
|   |-- notes-and-obsidian
|   |-- archive-safety
|   |-- rust
|   `-- web-and-static-sites
|
`-- nix/
    |-- nonactivating-validation
    |-- host-boundary-rules
    `-- lockfile-and-package-discipline
```

`nix` remains an independent source because changes can affect the active
machine. It should be selected deliberately, even when a personal project uses
Nix. `uv` belongs under Python dependency management rather than being a
top-level language or domain library.

## Seed Library

The first local source files now live under `libraries/`:

- `libraries/cdint/` provides a shared baseline, adaptive process rules, Go
  rules, and CDINT/PromiseGrid vocabulary.
- `libraries/ucd_research/` provides sanitized research and Python/uv rules.
- `libraries/nix/` provides standalone system-configuration safety rules.

`libraries/personal/` is intentionally deferred. Its useful source examples
contain private learner, vault, or machine context. Source: DI-voraz; DR-garom.

## Recommended Presets

Presets should be saved selections, not categories:

- `engineering-go`: core + Go + lightweight escalation.
- `engineering-rust`: core + Rust + lightweight escalation.
- `python-data-tool`: core + Python/uv + external-service or data-pipeline safety.
- `research-sci`: core + research + scientific reproducibility + a relevant
  stack node.
- `teacher-bot`: core + teacher method/session/communication/ownership + a
  chosen content domain.
- `cdint-protocol`: core + CDINT + full decision-first + Go where applicable.
- `nix-admin`: core + Nix safety + lightweight escalation.
- `docs-vault-private`: core + notes + private local include.

## Templates And Base Configuration

Useful reusable templates:

- An `agents.yaml` base using explicit `sources`, `vars`, `output`, and the
  nested `doc` manifest model from `docs/DESIGN.md`.
- A minimal local identity override template with `repo_name`, project status,
  source-of-truth documents, stack, safe commands, generated paths, and local
  danger zones.
- A high-consequence overlay template: protected data, irreversible actions,
  required human review, approved validation, and prohibited runtime paths.
- A teaching session artifact template, separated from private learner records.
- A research pipeline template: raw inputs, transformations, outputs, pinned
  environment, known limitations, and validation evidence.

The most useful corpus template is the role-overlay boundary in `grid-examples`
and `wire-lab`: canonical shared rules are pointers; overlays may only add
stricter role/environment-specific rules. The RoSE danger-zone list and
`photo_getter` hard rules are excellent templates for high-risk overlays.

## Conflicts And Deliberate Compromises

- **Decision rigor versus throughput:** the CIWG/PromiseGrid approach requires
  complete pre-code decisions, paths, and compliance artifacts. Productive for
  protocol and durable architecture; excessive for narrow repairs. Offer both
  process levels and make the repository select one.
- **Universal Go architecture rule versus language idiom:** "object-oriented
  structs and methods" and "avoid `internal/`/`pkg/`" appear in one lineage.
  These are repo conventions, not generic Go truths. Keep them in a CDINT/Go
  overlay, never the universal core.
- **Detailed comments versus local readability:** some prompts demand comments
  for every touched non-trivial block; others prefer idiomatic concise code.
  The portable rule is preserve explanatory intent and comment non-obvious
  invariants. Degree of detail is a stack/repo policy.
- **Autonomous commit prompting versus user control:** several samples request
  frequent commit prompts or interpret a bare `commit` command. The portable
  rule is explicit commit authorization and staging; prompting cadence stays
  local.
- **Central canonical rules versus private context:** the corpus correctly says
  private runtime notes stay out of committed AGENTS files, but some samples
  embed machine paths and personal learning state. Mogent includes support a
  clean compromise: shared canonical library + local ignored/private overlay.
- **Current manifest model versus corpus's old module metadata:** design has
  superseded tag/block selection. Preserve corpus wording as source material,
  but author new libraries as normal heading trees and compose output through
  `agents.yaml`.

## Technology Assessment

- **Go:** best fit for mogent itself: fast static CLI, simple distribution,
  `text/template`, existing Cobra/Bubble Tea direction, and an already-Go
  project. The Go corpus consistently has clear deterministic test practices.
- **Rust:** best for robust local CLIs, native utilities, and high-integrity
  processing where ownership and type modelling help. Cargo provides a coherent
  command set. The OpenAI Codex sample is excellent as a source of narrow Rust,
  API, TUI, and snapshot-test subtrees, but unsuitable as a generic preset due
  to its monorepo/Bazel-specific requirements.
- **Python with uv:** best for scientific analysis, extraction, and API-backed
  tooling. The successful pattern is explicit environments, immutable/raw data
  exclusions, dry runs, and carefully bounded live behavior. Do not generalize
  unpinned scientific environments or a particular Python version.
- **Nix:** strong reproducibility and host composition, but the highest need for
  activation and hardware/network safety rules. A dedicated safety overlay is
  mandatory; it is not a default developer stack module.
- **MATLAB/Psychtoolbox and lab frameworks:** appropriate where lab hardware and
  existing experiment tooling determine the choice. They need protocol and data
  safety overlays, not replacement-stack recommendations.
- **Static Markdown/Hugo/Obsidian:** good for human-facing knowledge and sites;
  validate rendered output and preserve portable Markdown. Keep generator- or
  vault-specific features selectable.

There is no globally "best" stack. The strongest pattern is a small baseline,
a conventional stack subtree, and a domain-specific safety overlay.

## Sensitive Material Inventory

Do not copy these items into shared modules, public examples, template vars, or
generated output. This memo intentionally describes categories rather than
repeating the contents.

- `Notes_Vault_Agents_AGENTS.md` contains a local cloud-vault path, personal
  biographical/academic details, memory instructions, and relationship/persona
  parameters including intimate fields. Treat as private; do not use as a
  shared include.
- `pg_learning_AGENTS.md` contains private learner progress, session resume
  state, strengths/gaps, and longitudinal learning records. Keep it private and
  reduce any reusable teacher module to method only.
- `RoSE_agents.md`, `SEFproject_agents.md`, `SemanticMapProject_agents.md`, and
  `VideoProject_agents.md` contain research/experiment context, participant or
  subject-data boundaries, lab procedures, and potentially identifiable local
  artifacts. Preserve safety rules but never copy study facts or data paths.
- `photo_getter_agents.md` concerns a family archive and mounted-drive state.
  It is sensitive personal media context; extract only generic non-destructive
  archive rules.
- `piazza_scrape_agents.md` involves course data and credentials. Reuse only
  credential hygiene, request pacing, and dry-run patterns.
- `color-engine_agents.md`, `seam_game_AGENTS.md`, the Nix assessments, and
  CIWG/PromiseGrid samples include machine-local absolute paths, local agent
  references, host names, or user-environment information. Remove those from
  portable modules.
- `stevegt_grokker_refs_heads_main_AGENTS.md` names an environment variable for
  an API credential. Keep the generic "never expose credentials" rule, never
  reproduce credential setup in a shared prompt.

## Suggested Next Work

1. Decide the first ecosystem slice: recommended `core`, `code/go`,
   `process/lightweight-escalation`, and `safety/secret-and-generated-artifacts`.
2. Run a focused TE before selecting the default process profile, because the
   full decision-first protocol and lightweight workflow impose materially
   different operational costs.
3. Create the first `agents.yaml` and source-library examples after rebuild
   milestone 1 makes the manifest parser and renderer available.
4. Add `teacher`, `research`, and `cdint` libraries as separate increments,
   with a privacy review before importing any source wording.
