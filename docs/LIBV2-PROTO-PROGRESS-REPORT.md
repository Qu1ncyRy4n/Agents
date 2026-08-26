# Library V2 Protolibrary Progress Report

Date: 2026-08-21
Status: active extraction and owner review

## Why This Work Restarted

The first library under `libraries/` contains useful guidance, but it does not
fulfill the intended job of an exhaustive protolibrary. The earlier extraction
optimized for compact generalized principles before recording every source
idea, condition, example, and competing policy.

That produced several recurring problems:

```text
SOURCE AGENT GUIDES
        │
        ├─ repeated organization policy
        ├─ detailed examples and exceptions
        ├─ strict and lightweight alternatives
        ├─ repository-local facts
        └─ generated assessments
        │
        ▼ premature compression
V1 LIBRARY
        ├─ useful principles
        ├─ missing operational detail
        ├─ unresolved overlap
        ├─ extraction-facing headings such as “Corpus Variants”
        ├─ almost no examples
        └─ no complete source-to-module ledger
```

The correction is not to discard v1 immediately. V1 remains a curated candidate
and comparison target while v2 performs lossless intake.

## Current Strategy

```text
37 PROTECTED SOURCE GUIDES
        │
        ▼
FILE-LEVEL PROVENANCE LEDGER       complete
        │
        ▼
HEADING-LEVEL SEMANTIC INTAKE      active
        ├─ preserve verbatim wording
        ├─ retain examples and exceptions
        ├─ cite duplicate sources
        ├─ mark generated material
        ├─ preserve competing policies
        └─ record possible output destinations
        │
        ▼
OWNER REVIEW
        ├─ keep
        ├─ adapt
        ├─ combine as a compatible bundle
        ├─ split for deliberate selection
        ├─ retain as a choice
        ├─ move to skill/spec/generated documentation
        └─ reject with a recorded reason
        │
        ▼
ORGANIZATION-OWNED LIBRARIES
        │
        ▼
EXAMPLE MANIFESTS AND RENDER VALIDATION
        │
        ▼
V2 ADOPTION AND EVENTUAL V1 RETIREMENT
```

## Decisions and Working Conventions

These are current working conventions for extraction, not all permanent Mogent
product decisions.

### Preserve Before Pruning

Every useful source instruction receives an explicit disposition. Repetition,
specificity, or disagreement is not a reason for silent omission.

### Organize Intake by Meaning Before Ownership

New candidates begin under `libv2-proto/intake/<semantic-area>/`. This makes
overlap visible before deciding whether cdint, personal, UCD Research, or
another organization publishes the canonical result.

Reviewed modules later move to `orgs/<owner>/<semantic-area>/`.

### Keep Intended Bundles Together

A document may remain a substantial multi-heading compatible bundle when its
guidance is intended to travel together. Independently optional or mutually
exclusive content becomes a separate module so selection is deliberate.

### Make Titles State the Direction

Titles should communicate both subject and instruction, such as “Preserve
Decision History Through Supersession,” rather than category-only labels such
as “Decision Intent,” “Code,” or “Corpus Variants.”

### Cite One Verbatim Copy and Every Duplicate

When several source guides repeat identical wording, v2 preserves one verbatim
copy and cites all duplicates. This retains provenance without manufacturing
several prose copies that can drift.

### Separate Editorial Confidence from Composition

`reliable`, `needs-choice`, and `generated-candidate` describe review state.
They do not mean required, optional, or mutually exclusive.

### Use Plain Relationship Terms

Current review vocabulary includes “requires all,” “requires one choice,”
“mutually exclusive,” “exactly one,” “any combination,” “optional,” and “see
also.” Boolean notation may explain a relationship but is not the preferred
user-facing vocabulary.

## Important Choices Still Open

```text
Decision workflow
├─ risk-based routine work + escalation              recommended candidate
└─ strict decisions before every code/path choice    preserved source option

Module composition
├─ compatible multi-heading bundle
└─ separately selectable child modules

Module includes after TE-nufad
├─ manifest composition + separate modules           default survivor
├─ frontmatter + minimal named insertion marker      survivor if needed
└─ restricted literal Go-template include            survivor if needed

Path evolution
├─ sparse old-path -> new-path move map               current provisional lean
├─ permanent ID for every module
├─ breaking rename plus communication
└─ stable global library identity

Output destination
├─ rendered agent guide
├─ skill
├─ generated reference document
├─ product or protocol specification
├─ library-authoring guidance
└─ repository-local instruction
```

These should be decided from real extracted examples rather than from invented
schemas alone.

## Module Include Direction After TE-nufad

The first sketch placed strict YAML inside an HTML comment beside the affected
heading. The Steve meeting raised rendering, parsing, and maintainability
problems with that approach. `TE-nufad` now rejects rich embedded YAML as the
default rather than treating it as the leading syntax.

The constrained happy path is manifest composition with intrinsic metadata in
file frontmatter. If exact mid-document insertion is required, the remaining
prototype choices are a frontmatter-named minimal marker or a restricted literal
Go-template include. Choice logic does not belong in the body marker.

### Historical Inline YAML Block

This is retained as rejected design evidence, not the current recommendation:

```markdown
### Adopt a Communication Style

<!-- yaml:mogent:start
imports:
  communication/persona:
    choose: exactly_one
    from:
      - none
      - self:personas/alien
      - self:personas/surfer

  communication/roles:
    choose: any
    from:
      - self:roles/librarian
      - self:roles/code-assistant
      - self:roles/senior-developer
yaml:mogent:end -->
```

The later design slice should instead answer:

- Is exact mid-document insertion needed for the first libv2 cutover?
- If so, should a prototype compare a minimal marker with a restricted literal
  Go-template include?
- Does inclusion insert a full source node and descendants, or only body text?
- Can included content include another module, and what depth is allowed?
- Which choice rules belong in frontmatter, a source descriptor, or the
  consuming manifest?

### Example Review Note

Provisional author notes may use ordinary HTML comments so they do not appear
in rendered Markdown:

```markdown
<!-- review-note: This source rule may belong in a skill rather than AGENTS.md.
Compare it with the repository-local procedure before promotion. -->
```

Open question: should review notes remain unconstrained author prose, or use a
small schema so tooling can list every unresolved note?

### Example Section TLDR

```markdown
## Preserve Existing Work
<!-- tldr: Never overwrite concurrent user changes without explicit direction. -->
```

Open question: should section TLDR comments remain authoring-only, appear in
source browsing, or optionally render into generated output?

## Provisional Diagnostic Levels

The final warning/error boundary needs its own decision. Current candidates:

| Condition | Provisional level | Reason |
|---|---|---|
| malformed frontmatter or duplicate YAML keys | error | Meaning cannot be determined safely. |
| path escapes a source root or traverses a rejected symlink | error | Violates the source trust boundary. |
| direct or indirect import cycle | error | Rendering cannot produce a finite deterministic result. |
| unresolved `exactly_one` import | error | Required authored content is missing. |
| selected mutually exclusive modules | error | Output would contain known incompatible instructions. |
| missing required module | error after migration | Selected guidance is known to be incomplete. |
| old path with a valid compatibility move | warning | Rendering may continue, but the manifest is stale. |
| move mapping past its deprecation date | open | Could warn, error, or follow an explicitly configured policy. |
| overlapping modules with no declared incompatibility | review notice | Similarity alone does not prove conflict. |
| unselected `see also` module | information | Discovery only. |
| generated-source candidate | review label | Source authority differs; composition may still be valid. |

Questions include whether builds may opt into warnings-as-errors, whether
organization policy can raise severity, and whether legacy repositories receive
a migration period before new relationship errors become blocking.

## Proposed Thought Experiments and Decision Records

These are needed decision slices, not approved TE/DR artifacts or assigned IDs.

### Module Includes and Choice Rules

`TE-nufad` completed the broad narrowing pass. Owner decisions now determine
whether exact insertion belongs in the first cutover and whether the two
syntax-bearing survivors deserve prototypes. Test any prototype with nested
includes, no selection, multiple selection, cycles, source updates, and
readable interactive prompts.

Decision output: syntax ownership, permitted choice rules, insertion behavior,
and warning/error boundaries.

### Compatible Bundles and Directory Inheritance

Compare one large multi-heading document with separate child modules. Test
selecting a parent, excluding a descendant, swapping one project-specific
section, browsing the tree, and reviewing upstream additions.

Decision output: when a file is one selection unit, what directory selection
inherits, and whether headings inside a compatible bundle are addressable.

### Source-Path Moves and Compatibility

Compare sparse move maps, permanent node IDs, breaking renames, and stable
published-library identities. Test pinned updates, move chains, forks, expired
mappings, rollback, and two prior paths converging on one target.

Provisional syntax:

```yaml
moves:
  process/thought-experiment:
    to: process/narrow-options-with-thought-experiments
    deprecated_after: 2027-01-01
```

Decision output: compatibility duration, default behavior, reversibility, and
whether migration belongs inside an existing operation rather than a new
command.

### Strict and Risk-Based Decision Workflows

Compare the complete strict source-guide protocol with risk-based escalation.
Test typo fixes, local bug fixes, public API changes, persistence migrations,
security boundaries, research validity, and cross-repository ownership.

Decision output: canonical baseline, selectable overlays or alternatives, and
which decision artifacts are mandatory at each risk level.

### Multiple Output Targets

Determine which content Mogent should render, alias, synchronize, generate, or
only reference. Candidate targets include:

```text
AGENTS.md
├─ aliases or synchronized equivalents such as CLAUDE.md
├─ organization and repository overlays
└─ tool-specific instruction formats

other module outputs
├─ skills
├─ SPEC.md and protocol specifications
├─ generated human documentation
├─ generated agent-reference documentation
├─ TODO/worklists and decision-record templates
├─ library-authoring guidance
└─ repository-local procedures and active state
```

Questions include whether aliases are real symlinks or separately rendered
files, how divergence is detected, whether one manifest owns several outputs,
and which metadata is shared versus output-specific.

## Smaller Questions to Preserve

- Should common language packages be discoverable modules, generated reference
  docs, or skills rather than permanent prompt content?
- Does Nix belong under `languages`, `tools`, or both through one canonical
  module plus cross-links?
- Should Git guidance live under `tools/git` even when it expresses general
  workflow policy?
- How are Markdown anchors validated when headings change?
- Can provenance links use a source-root shorthand without sacrificing ordinary
  Markdown compatibility?
- Should generated source guides be quarantined physically or identified only
  through metadata and review labels?
- When two modules overlap, must authors declare whether they are compatible,
  mutually exclusive, or one supersedes the other?
- Should repository facts and active state render into agent instructions, a
  generated reference, or remain repository-local source material?
- Which proposed modules are really operational skills with explicit triggers,
  inputs, outputs, and allowed writes?
- Should an organization publish recommended presets separately from required
  policy?

## Current Review Batch

```text
intake/
├── workflow/
│   ├── keep-changes-scoped-and-preserve-existing-work.md
│   ├── use-risk-to-choose-routine-work-or-decision-review.md
│   └── require-human-review-for-high-consequence-changes.md
└── process/decision-governance/
    ├── use-thought-experiments-to-narrow-a-broad-design-space.md
    ├── lock-durable-decisions-before-implementation.md
    ├── preserve-decision-history-through-supersession.md
    └── handoff-decisions-with-implementation-evidence.md
```

The strongest immediate policy question is whether strict decision-first is a
complete selectable governance bundle or whether some of its stages should be
promoted independently into a lighter canonical workflow.

## Repository History Created During This Review

- `86780a0 Preserve v1 library review state`
- `1eebacb Start lossless v2 agent protolibrary`
- `895fc96 Define provisional v2 library semantics`
- `9f2f6a4 Reorganize v2 candidates by semantic intake`
- `f19efb2 Add pending-review v2 workflow candidates`
- `7112e90 Add pending-review v2 decision candidates`
- `4c6feaf Track pending v2 intake review`
- `6c8575f Document pending v2 review checkpoint`
- `5a55b8a Narrow v2 module include design`
- `c701667 Migrate canonical Mogent tasks into repository`

The review batch is committed for durable inspection but remains pending owner
feedback. Committed does not mean approved or promoted into an organization
library. Canonical detailed task state now lives in `TODO/TODO.md`; use
`TODO/PICKUP-2026-08-25-libv2-intake.md` as the current resume point.
