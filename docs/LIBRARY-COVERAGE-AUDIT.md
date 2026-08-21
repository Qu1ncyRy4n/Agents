# Agent Library Coverage Audit

Status: first authorized comparison pass, 2026-08-19.

## Scope And Method

This audit compares the 51 current Markdown modules under `libraries/` with the
protected original agent guides and the earlier sanitized research in
`docs/codex_eco/README.md`. It records reusable coverage and omissions without
copying private paths, personal state, study details, credentials, or large
source passages.

This is a category-level pass, not yet the final source-by-source traceability
matrix. `docs/LIBRARY-REVIEW-PLAN.md` defines that later local audit and the
decisions required before source-path or public metadata changes.

Coverage labels:

- **Preserved:** the reusable behavior has a clear current home.
- **Partial:** the main idea survived, but useful conditions, examples, or
  separation were lost.
- **Missing:** a reusable idea has no clear current home.
- **Deliberately local:** the material should not become a shared module.
- **Conflict:** original guides express materially different policies that need
  selection or an explicit portable compromise.

## Coverage Summary

| Area | Coverage | Current home | Main finding |
|---|---|---|---|
| Focused change loop | Preserved | `cdint:shared-baseline/instructions/focused-change-loop` | Good minimal baseline; add expected-usage evidence when behavior is visible. |
| Source of truth | Preserved | `cdint:shared-baseline/identity/source-of-truth` | Keep repository facts local and durable design authoritative. |
| Safe Git behavior | Preserved | `cdint:shared-baseline/instructions/git` | Commit cadence and authorization remain a policy choice, not a universal default. |
| Runtime and generated artifacts | Preserved with overlap | baseline constraints plus `cdint:engineering/runtime-artifacts` | The shallow/deep layering is sensible, but repeated wording needs a clearer selection story. |
| Reviewability and user edits | Partial | `cdint:engineering/reviewability` | Preserves user work and final inspection; lacks a distinct delta-review/classification workflow. |
| Comment intent | Preserved | `cdint:engineering/comment-intent` | Strong portable compromise between exhaustive comments and idiomatic concision. |
| Error handling | Partial | `cdint:engineering/error-handling` and language modules | Core behavior survived; unfamiliar shell and Go terms need concrete examples and language-specific placement. |
| Test strategy | Partial | `cdint:engineering/test-strategy` and language modules | Good principles; examples for offline, success path, structured assertions, and snapshots were lost. |
| Decision-first governance | Partial | `cdint:process/decision-first` subtree | The adaptable core survived, but intake, evidence, stop rules, artifact lifecycle, editing, and examples were heavily compressed. |
| High-consequence review | Preserved | `cdint:process/high-consequence-review` | Correctly generalized; repo-specific danger zones must remain local overlays. |
| Coordination identifiers | Partial | `cdint:engineering/coordination-ids` | Stable handles survived; namespace, minting, history/supersession, and active/archive explanations are too terse. |
| Public technical prose | Preserved with duplication | `cdint:docs/public-prose` and `cdint:cdint-and-promisegrid/public-technical-prose` | Separate universal public prose from PromiseGrid provenance additions. |
| General documentation | Missing/partial | `cdint:shared-baseline/instructions/documentation` | Needs audience, purpose, learning/reference intent, canonical layout boundaries, and inconsistent-style escalation. |
| PromiseGrid terminology and protocol restraint | Preserved | `cdint:cdint-and-promisegrid` | Good domain-specific placement; some broadly useful prose should not be duplicated here. |
| Go practices | Partial and redundant | `cdint:go` plus engineering modules | Current `Code` module is too generic; retain only genuinely Go-specific API, package, tooling, and test guidance. |
| Rust practices | Preserved with extraction residue | `personal:lang/rust` | Strong typed/API and migration guidance; the final corpus-summary section should be decomposed or removed. |
| Python and uv | Preserved with overlap | `personal:lang/python` and `ucd_research:python` | General environment rules and research-specific analysis boundaries need a cleaner ownership split. |
| Nix | Preserved with two mixed variants | `personal:lang/nix` | Host activation safety and project development-shell guidance should be separately selectable. |
| Static sites | Preserved | `personal:domain/static-sites` | Useful conditional module, not a universal browser-stack rule. |
| Research safety and reproducibility | Preserved with duplication | personal research plus `ucd_research` | Shared safety concepts are good; participant, procedure, and lab facts stay local. |
| Media/archive safety | Preserved | `personal:domain/media-archive` | Correctly extracts non-destructive and single-writer patterns without private archive facts. |
| Notes and knowledge work | Preserved/partial | `personal:domain/notes-vault` | Portable notes and skill boundaries survived; detailed vault state and identity material were correctly omitted. |
| Teaching method | Partial | `personal:communication/teacherbot` | Strong method summary; learner ownership, session closeout, and communication style should be independently composable. |
| Personas | Preserved but unmodeled | `personal:communication/personas` | These are mutually exclusive alternatives, not one subtree to render wholesale. |
| Local services and staged migrations | Preserved | `personal:engineering` | Broadly reusable API -> adapters and compatibility rules are candidates for promotion into `cdint`. |
| External-service safety | Missing | no dedicated module | Credential boundaries exist in the baseline, but dry runs, pacing, bounded live tests, and rate-limit behavior need a conditional module. |
| Repository context and safe commands | Deliberately local | consuming repository identity/overlay | Project map, current state, commands, paths, and danger zones should be templates or local modules, not copied universal prose. |

## Important Missing Detail

### Decision-First Is Too Compressed

The original governance family contains substantially more operational detail
than the current five small process modules. Useful candidates to recover are:

- what evidence is required before escalating a question;
- how to compare alternatives under normal, failure, trust, concurrency, and
  scale conditions;
- when analysis stops and a human decision is required;
- how a new decision supersedes rather than silently rewrites old intent;
- how decision, finding, open-question, and handoff artifacts relate;
- what may be edited in place and what is historical evidence;
- a small worked example from question through decision to implementation
  evidence.

Do not restore the old governance text as one mandatory 200-line module. Keep a
short adaptive parent and offer deeper selectable children.

### Documentation Is Underdeveloped

The current docs modules emphasize consistency and public prose. The source
guides also imply reusable questions that are absent:

- Who is this document for: a user, maintainer, agent, learner, or operator?
- Is it a tutorial, task recipe, reference, explanation, specification, or
  handoff?
- Which file is canonical when examples and implementation disagree?
- Should inconsistent local style be preserved, repaired in scope, or raised?
- What usage and expected result should accompany a user-visible change?

These should become small documentation modules rather than one universal
canonical layout.

### Change Review Needs Its Own Shape

`reviewability` says to preserve user edits and inspect the final diff. A useful
missing workflow is to classify a pre-existing or incoming delta before acting:

1. requested change;
2. user-owned concurrent change;
3. generated or runtime artifact;
4. unrelated cleanup opportunity;
5. conflicting edit that requires coordination.

That module would explain what to stage, leave alone, report, or escalate. It
should complement minimal-change guidance instead of repeating it.

### External Services Need Conditional Safety

Several source guides contain reusable behavior for API-backed or scraping
tools. A dedicated module should cover:

- credentials through documented secret mechanisms, never logs or fixtures;
- offline fixtures as the default test path;
- explicit approval before a live or billable action when the task did not
  already authorize it;
- dry runs and bounded smoke tests;
- request pacing, retry/backoff, rate limits, and partial-failure reporting;
- no claim of success merely because a remote dependency was unavailable and a
  test skipped.

## Conflicts And Portable Resolutions

| Conflict in source guidance | Recommended treatment |
|---|---|
| Full decision protocol for every change versus fast local iteration | Keep lightweight risk-based escalation as the baseline and deeper governance as a deliberate selectable subtree. |
| Frequent autonomous commits versus explicit user review before commit | Model commit policy as alternatives; keep explicit authorization in the safe baseline. |
| Exhaustive comments versus idiomatic concise code | Preserve explanatory intent and require comments for non-obvious invariants; let repos add stricter comment density. |
| Object-oriented Go conventions versus language-idiomatic typed designs | Keep repository-specific architecture local; shared Go guidance should describe behavior, API clarity, and tooling rather than one paradigm. |
| Strict host Nix safety versus using Nix only as a project dev shell | Split host activation and dev-environment modules so selection expresses the actual risk boundary. |
| Shared canonical guide versus private machine, learner, or research state | Keep shared modules portable and put private facts in local ignored overlays or other explicitly protected state. |
| Multiple personas rendered together | Declare an alternative family and reject or warn on multiple selections; never use source order as "first wins." |

## Relationship And Metadata Findings

The current library has nine `requires` entries. They are displayed but not
validated, and their source aliases are not portable because each manifest
chooses its own aliases. Two modules use rendered `See Also` sections instead
of tool-only discovery metadata. Personas have no declared alternative family.

Recommended semantic split:

- hard prerequisite: `requires`;
- optional discovery: future `see_also`;
- known exact incompatibility: `conflicts_with`;
- mutually exclusive choices: future alternative family plus variant;
- historical rationale: provenance/Decision Intent reference.

The reference syntax must be decided before these fields drive warnings or
automatic fixes. Same-library relationships need an alias-independent form;
cross-library relationships need stable identities or must remain advisory.

## Intentionally Omitted Material

The following source content should not be treated as missing shared coverage:

- private learner progress, biography, persona parameters, or session state;
- research subjects, raw data, study parameters, lab paths, or active-machine
  state;
- credentials, credential names, mounted paths, and host-specific details;
- repository-specific architecture, commands, current status, and generated
  paths, except as local-module templates;
- assessments or suggested skeletons presented as if they were active rules;
- monorepo-specific build, snapshot, or package policy generalized to every
  project using the same language;
- obsolete block/tag configuration superseded by the current manifest model.

## Recommended First Slices

1. **Clarity without path changes:** replace the six `Corpus Variants` sections
   with concrete instructions, splits, or omission; add short examples for the
   unclear testing/error/Git/coordination terms.
2. **Documentation depth:** add audience/purpose and usage-evidence modules.
3. **Decision process depth:** recover selected workflow and lifecycle detail as
   atomic children under the existing decision-first subtree.
4. **Missing conditional modules:** add delta review and external-service safety.
5. **Decision gate:** resolve alias-independent relationships, alternatives,
   taxonomy, and path compatibility before moves or metadata enforcement.
6. **Bounded reorganization:** split Nix host safety from dev-shell guidance,
   separate Teacherbot methods, and then review Go and Python ownership.

Each slice should have its own focused diff and example render. Do not combine
the source-path reorganization with the public metadata-format change.
