# Mogent: Reusable, Reviewable Agent Instructions

Quick presentation · 2026-08-20

> Mogent composes repository-specific `AGENTS.md` files from reusable Markdown
> libraries and an explicit, reviewable `agents.yaml` manifest.

Suggested 10-minute route: sections 1–3, 6, 8–10, 15–21, and 24. Use section
25 for a live demo. The library reorganization, organization-policy material,
and appendices are available for questions.

---

## 1. The Problem

Teams maintain agent instructions in many repositories. Those files usually
share the same core material:

- coding and testing practices;
- safety and review rules;
- Git and documentation expectations;
- organization-specific terminology and process;
- language-, project-, or persona-specific guidance.

Copying one large `AGENTS.md` into every repository creates predictable drift:

```text
Repo A fixes a safety rule
       |
       +-- Repo B never receives it
       +-- Repo C rewrites it differently
       `-- Repo D cannot tell what is shared versus local
```

The goal is not to hide this complexity. The goal is to make reuse, local
choice, provenance, and updates visible enough to review.

---

## 2. The Core Idea: Two Trees, One Manifest

```text
Reusable library tree                Repository document tree

libraries/cdint/                     agents.yaml
|-- shared-baseline/                 |-- Identity
|-- engineering/                     |-- Instructions
|-- process/                         |-- Constraints
`-- go/                              `-- Format
          \                              /
           \                            /
            +---- mogent build --------+
                         |
                         v
                     AGENTS.md
```

- The **library** owns reusable Markdown content.
- The **manifest** owns the repository's rendered outline, order, headings, and
  source selections.
- `AGENTS.md` is generated output that remains readable and reviewable without
  requiring Mogent at agent runtime.

The manifest reads like the output table of contents rather than a set of
hidden merge rules.

---

## 3. Design Philosophy

### Markdown-native

A normal Markdown heading and its content form a selectable source node. The
library is still useful with ordinary editors, Git, and Markdown tools.

### Explicit composition

Every source has a user-defined alias such as `cdint`, `personal`, or `local`.
References name their provenance:

```text
cdint:engineering/test-strategy
personal:communication/personas/operator
```

There is no ambient library search path and no first-source-wins behavior.

### The manifest owns the output

Source headings help discovery, but the repository decides where and under
what heading the material renders.

### Fail loudly

Unreadable sources, duplicate YAML keys, unresolved nodes, path escapes,
ambiguous selection, and empty output are errors—not silent omissions.

### Deterministic and review-first

The same explicit libraries, lock data, variables, and manifest produce the
same output. Mutations support previews, complete validation, and atomic writes.

### One product core

The CLI and TUI are clients of the same public Go workspace API. Durable
behavior is not trapped in terminal handlers.

---

## 4. What A Library Looks Like

```text
libraries/cdint/
|-- shared-baseline/
|   |-- identity/
|   |   `-- role.md
|   `-- instructions/
|       `-- focused-change-loop.md
|-- engineering/
|   |-- test-strategy.md
|   `-- error-handling.md
`-- process/
    `-- lightweight-escalation.md
```

Example module:

```markdown
---
tags: [testing/strategy, scope/team]
tldr: Use focused deterministic tests, then broaden with risk.
priority: 0.9
scope: team
---
# Test Strategy

Scale validation to risk. Use focused checks for narrow changes...
```

Its source reference is:

```text
cdint:engineering/test-strategy
```

Directory paths and Markdown heading paths contribute to identity; filenames
do not.

---

## 5. What A Repository Manifest Looks Like

```yaml
sources:
  cdint: ../agent-libraries/cdint
  local: .mogent/library

vars:
  repo_name: example-service
  repo_url: https://example.test/cdint/example-service

output: AGENTS.md

doc:
  - Identity:
      - Role: cdint:shared-baseline/identity/role
  - Instructions:
      - Workflow: cdint:shared-baseline/instructions/focused-change-loop
      - Testing: cdint:engineering/test-strategy
  - Constraints:
      - Safe Defaults: cdint:shared-baseline/constraints/safe-defaults
```

What is shared, what is local, what order it appears in, and what source each
section came from are visible in one file.

---

## 6. Basic Usage Walkthrough

### 1. Preview a starter

```sh
mogent init --template minimal \
  --source cdint=../agent-libraries/cdint \
  --dry-run
```

`init` discovers useful repository values, materializes them in the manifest,
and validates the complete proposed render before writing.

### 2. Explore the library

```sh
mogent source list cdint --tree
mogent source list cdint --search test --tldr
mogent source show cdint:engineering/test-strategy --metadata
```

### 3. Preview a manifest addition

```sh
mogent add cdint:engineering/test-strategy \
  --under Instructions \
  --dry-run \
  --preview=patch
```

Placement can also use stable named anchors:

```sh
mogent add cdint:engineering/error-handling \
  --after Instructions/Testing \
  --dry-run
```

### 4. Render and inspect

```sh
mogent build
mogent status
git diff -- agents.yaml AGENTS.md
```

The generated `AGENTS.md` is normally committed so humans and agents can use
and review it without installing Mogent.

---

## 7. Remote Sources Without Surprise Updates

A repository may declare a narrow subdirectory of a Git source:

```yaml
sources:
  cdint:
    location: https://github.com/example/agent-library.git
    subdir: libraries/cdint
```

```sh
mogent source pin cdint
mogent source update cdint
mogent source update cdint \
  --ref 9b77c4d06c651362a363d195328007f43473df1a \
  --accept
```

The committed lock records the exact Git revision, subdirectory, and Markdown
hash. Ordinary workspace commands remain offline and never silently float to a
new upstream version.

---

## 8. Dogfooding Strategy

Mogent is tested at three levels:

```text
Public Go API tests
    Same workspace operations used by the CLI
            |
            v
Disposable testing-ground fixtures
    Repeatable parsing, rendering, mutation, and failure smoke tests
            |
            v
Real working repositories
    Is the workflow actually understandable and useful?
```

Real-project dogfooding is the important product test. `testing-ground/` is
retained for controlled failure and mutation cases, but it is not evidence that
the authoring experience feels good.

Dogfood work has included:

- building and editing manifests in an active project;
- browsing the real CDINT and personal libraries;
- selecting whole source directories and narrower modules;
- localizing shared content without editing its origin;
- importing a single direct edit conservatively;
- pinning this repository over HTTPS at an immutable revision;
- using the public Go workspace API for the same load/status/discovery/add flow
  exposed by the CLI.

---

## 9. Improvements Driven By Dogfooding

| Observed friction or risk | Improvement made |
|---|---|
| It was hard to discover useful modules from a long flat list. | Added source search, tag search, TLDR display, metadata display, and aligned tree views. |
| Coverage and source browsing looked like separate models. | Unified them around one source tree; `coverage` is now a coverage-oriented presentation of the same inventory. |
| Adding a directory did not explain all inherited descendants. | Directory selections became first-class, with expanded tree previews and explicit inherited state. |
| Placement by index would become stale. | Added `--under`, `--first`, `--last`, `--before`, and `--after` using named manifest paths. |
| Typos and organizational directory references produced poor errors. | Added close-match and descendant suggestions with actionable next commands. |
| Consuming a complete remote repository produced duplicate heading paths. | Added pinned source `subdir`, so a manifest can trust and hash only the intended library root. |
| Builds varied with checkout directory or Git remote. | `init` now writes `repo_name` and `repo_url` explicitly; `build` uses only manifest variables. |
| Local symlinks could read Markdown outside the declared source boundary. | Local and pinned sources now reject symlinked roots, directories, and Markdown files. |
| Reusable behavior lived behind Go's `internal/` boundary. | Published focused top-level packages and made the CLI/TUI use the public workspace core. |
| Atomic writers handled cleanup differently. | Centralized atomic replacement and joined meaningful cleanup failures. |
| Helper hints were useful but could become noisy for agents. | Retained actionable hints and recorded configurable presentation as follow-up work. |

This is the central development loop:

```text
Use Mogent on real work
        -> find repeated friction
        -> record the product question
        -> make one focused change
        -> exercise the same real workflow again
```

---

## 10. Current Product Boundary

### Implemented now

- strict local and pinned-Git sources;
- explicit `agents.yaml` composition;
- deterministic rendering and overwrite protection;
- source browsing, search, metadata, TLDRs, and tree coverage;
- dry-run manifest additions and placement previews;
- selectable directory subtrees with exclusions;
- copy-on-write localization with provenance;
- conservative single-section direct-edit import;
- immutable URL locks and reviewed source updates;
- public Go workspace operations used by CLI and TUI.

### Not yet complete

- full localized-origin inheritance and three-way merging;
- enforced `requires`, conflicts, and exclusive choice groups;
- organization policy profiles;
- a settled library taxonomy and complete module-content review;
- source-module creation, move, and reorder workflows;
- local-to-upstream proposal/PR delivery;
- multiple outputs such as `CLAUDE.md` and `GEMINI.md`;
- skills composition across agent ecosystems.

---

## 11. The Library Reorganization Problem

The implementation is now capable enough to expose content-model problems:

- some modules are broad bundles while others are atomic choices;
- compatible and mutually exclusive content sometimes lives in the same file;
- several sections are missing, thin, duplicated, or overly repository-specific;
- file-level metadata applies to every heading in a file;
- personas and alternative workflows need explicit choice semantics;
- moving or renaming a heading can break source references;
- `requires` and `conflicts_with` are currently display-only metadata;
- cross-library relationships introduce identity and trust complexity.

This is why the library review comes before adding a large relationship
language. A better-organized library removes many apparent metadata problems.

---

## 12. Proposed Library Reorganization

Use a simple authoring convention:

> Keep content together when it is always selected as one unit. Split content
> into a separate one-top-level-heading file when it is optional, independently
> selectable, or mutually exclusive with a sibling.

Example:

```text
communication/
|-- styling.md                 compatible baseline bundle
|-- styling/
|   |-- tts-chat.md            optional
|   |-- tts-output-stream.md   optional
|   `-- whiteboard.md          optional
|-- personas/
|   |-- surfer.md              exclusive persona choice
|   |-- robot.md
|   `-- alien.md
`-- roles/
    |-- librarian.md           compatible roles
    |-- teacher.md
    `-- reading-coach.md
```

Recommended metadata ownership:

- file frontmatter describes one selectable node;
- `library.yaml` describes groups, presets, and relationships spanning nodes;
- the filesystem supplies directory organization;
- heading-level relationship comments wait until dogfooding demonstrates a
  real need.

---

## 13. Proposed Relationship Vocabulary

Keep the first vocabulary small and concrete:

```yaml
---
requires:
  - self:safety/data-handling
conflicts_with:
  - self:workflow/unsafe-fast-path
exclusive_group: communication/persona
see_also:
  - self:roles/reading-coach
---
```

| Relationship | Meaning | Effect |
|---|---|---|
| `requires` | This module is incomplete or unsafe without another same-source module. | Both contents must be selected and rendered. |
| `conflicts_with` | Two exact modules cannot form a coherent instruction set. | Selecting both is an error. |
| `exclusive_group` | Members are choices for one slot. | At most one may be selected. |
| `see_also` | Related material may help discovery. | Informational; does not select content. |

`self:` means the same source instance, independent of whether the consuming
manifest calls that source `cdint`, `org`, or something else. Cross-library
relationships are deferred.

Mogent should never use first-match or source order to choose between
relationships. It reports the exact problem and lets the author decide.

---

## 14. Organization Policy Direction

The motivating requirement is straightforward:

> If a repository belongs to this organization, its generated instructions
> must include the organization's coding practices and safeguards.

The first version can keep that requirement explicit in each repository:

```yaml
policy:
  required:
    - cdint:shared-baseline
    - cdint:engineering/safety
```

If required content is absent from the rendered document, validation fails.

A later organization policy profile could centralize this, but the repository
must explicitly trust a pinned policy source. Mogent should not silently infer
authority from a Git remote or let a library declare itself mandatory.

Organization-wide adoption can be enforced in CI while the effective policy
remains visible in each repository's manifest and status output.

---

## 15. Why Inheritance Is The Next Important Problem

Copy-on-write localization is implemented:

```text
Shared module
     |
     | mogent localize
     v
Repository-local copy + provenance
```

That protects shared sources and makes local divergence explicit. But after the
origin improves, Mogent must answer:

- What exact version did the local copy start from?
- What changed upstream?
- What changed locally?
- Can the upstream update be applied safely?
- If both changed, should we preserve, merge, replace, or propose upstream?

A two-way diff cannot answer those questions reliably. Inheritance needs a
three-way model.

---

## 16. The Three-Way Inheritance Model

For every localized module, preserve:

```text
                         current upstream
                         ORIGIN
                            ^
                            |
                         changes
                            |
BASE -----------------------+
  |
  | changes
  v
LOCAL
current repository-owned version
```

- **Base**: exact content originally localized or last reconciled.
- **Origin**: current content at the pinned source revision.
- **Local**: repository-owned customized content.

Mogent compares `Base -> Origin` and `Base -> Local`, then classifies the
relationship:

```text
unchanged
fast-forwardable
same change on both sides
independent changes
overlapping changes
missing or mismatched provenance
```

Classification is safe to automate. Choosing what content to discard is not.

---

## 17. Worked Inheritance Example

Original localized content:

```text
BASE
Run focused Go tests for changed packages.
```

The organization improves the shared module:

```text
ORIGIN
Run focused Go tests and errcheck for changed packages.
```

The repository has also customized it:

```text
LOCAL
Run focused Go tests through `just test-api`.
```

Mogent reports both directions:

```text
Base -> Origin: adds errcheck
Base -> Local:  replaces the test command with just test-api
```

The first inheritance implementation should preview this and decline automatic
application because Local changed. A later reviewed merge could propose:

```text
Run focused Go tests through `just test-api`, then run errcheck for changed
packages.
```

The user still reviews that result before any write.

---

## 18. Directional Commands

The proposed public model uses verbs that say where content is moving:

```text
mogent status
    Read-only view of output, localization, and origin state.

mogent preserve <manifest-heading> --dry-run
    Generated AGENTS.md edit -> repository-local module.

mogent inherit <manifest-heading> --dry-run
    Newer origin -> localized repository module.

mogent propose <manifest-heading> --to <source> --dry-run
    Local improvement -> upstream patch or change proposal.
```

`reconcile` is the internal comparison operation. It does not need to become a
public command until it represents a distinct user job.

`status` should become the canonical read-only report. The current `drift`
command mixes reporting and mutation; its mutations can migrate to the
directional commands instead of growing another command vocabulary.

---

## 19. Safe Inheritance And Exact Acceptance

Every inheritance mutation starts with a preview:

```sh
mogent inherit Instructions/Testing --dry-run
```

Example output:

```text
Current localized revision:
  a184a144da09d44853f308eea564d83c78ce7237

Previewed origin revision:
  9b77c4d06c651362a363d195328007f43473df1a

Change:
  + Require errcheck after Go behavior changes.
```

Acceptance repeats the exact immutable revision:

```sh
mogent inherit Instructions/Testing \
  --accept 9b77c4d06c651362a363d195328007f43473df1a
```

This means “apply the bytes I reviewed,” never “apply whatever the branch points
to now.” A changed origin identity, concurrent local edit, missing base, or
ambiguous heading mapping causes the operation to stop.

---

## 20. Inheritance Strategies

Implement the low-risk behavior first:

### Phase 1: fast-forward only

If Local is unchanged from Base, preview and accept the new Origin. This is
simple, deterministic, and does not discard customization.

### Phase 2: classify all three-way cases

Show independent and overlapping changes without writing. This tests whether
the provenance and comparison model is understandable.

### Phase 3: reviewed non-overlapping merge

Prepare a merged result only when changes do not overlap. The complete output
still requires review and explicit acceptance.

### Explicit replacement

Replacing customized Local content with Origin is a separate destructive choice,
not a fallback hidden inside inheritance.

`--required-only` may filter organization-required updates, but `required` does
not mean every new revision is automatically trusted. A centrally reviewed and
pinned organization revision can reduce duplicate local review safely.

---

## 21. From Local Improvement To Pull Request

Inheritance is only one direction. Useful local changes should also be able to
flow back to the shared library:

```text
Local module
    |
    | mogent propose --dry-run
    v
Reviewable patch
    |
    | optional Git adapter
    v
Local branch + commit + prepared PR description
    |
    | explicit forge authority
    v
Pull request / merge request
```

The core operation does not require GitHub, GitLab, credentials, or a messaging
system. It can stop at a patch or PR-ready branch. Pushing and opening a change
request are separate, explicitly authorized effects.

---

## 22. In-Progress And Planned Features

Near-term work should stay ordered by dependency:

```text
Library module boundaries
        |
        v
Relationship and metadata semantics
        |
        +--------> Organization-required bundles
        |
        v
Provenance v2: retrievable Base + exact Origin identity
        |
        v
status -> preserve -> inherit -> propose
        |
        v
Git branch / forge proposal adapters
```

Other planned surfaces:

- guided source addition, module creation, move, and reorder;
- multiple explicit outputs such as `CLAUDE.md` and `GEMINI.md`;
- optional presentation controls for hints, color, width, TLDRs, and TTS;
- stale-document and handoff generation from explicit repository state;
- later skills support, after comparing how each agent ecosystem handles
  discovery, trust, pinning, precedence, and installation.

---

## 23. Decisions Currently In Front Of Us

| Decision | Current lean |
|---|---|
| What is one selectable module? | Keep compatible bundles together; split optional or exclusive choices. |
| Where does metadata live? | Node facts in file frontmatter; spanning groups/presets in `library.yaml`. |
| How are intrinsic relationships named? | Exact same-source `self:` references; defer cross-library semantics. |
| How is exclusivity represented? | `exclusive_group` means at most one; “one required” is separate. |
| Who can mandate organization content? | The repo manifest or an explicitly trusted pinned policy—not the library itself. |
| What is the read command? | `status` is canonical. |
| What is the first inheritance mode? | Dry-run classification plus fast-forward-only acceptance. |
| How are accepted updates identified? | Exact previewed Git revision or content digest. |
| How far does proposal delivery go initially? | Patch or PR-ready local branch before forge automation. |

---

## 24. What Success Looks Like

For a repository owner:

```text
I can see exactly which organization, personal, and local instructions apply.
I can preview changes before the manifest or output changes.
I can customize shared guidance without silently forking it forever.
I can inherit reviewed upstream improvements without losing local work.
I can propose useful local improvements back to the shared library.
```

For an organization:

```text
Shared safeguards can be adopted and audited across repositories.
Policy changes appear as ordinary manifest, lock, and AGENTS.md diffs.
Repositories retain explicit local structure and controlled customization.
No runtime service is required for agents to read the result.
```

The larger idea is a visible lifecycle for agent guidance:

```text
author -> compose -> review -> render -> customize -> inherit -> contribute
```

---

## 25. Suggested Live Demo

If time is short, show these five commands:

```sh
# 1. What is the workspace state?
mogent status

# 2. What reusable guidance exists?
mogent source list cdint --tree --tldr

# 3. What would one module contribute?
mogent source show cdint:engineering/test-strategy --metadata

# 4. What would adding it change?
mogent add cdint:engineering/test-strategy \
  --under Instructions --dry-run --preview=patch

# 5. Build and review ordinary project files.
mogent build
git diff -- agents.yaml AGENTS.md
```

Then use the three-way Base/Origin/Local diagram to explain where inheritance is
going next.

---

## Appendix A: Vocabulary

| Term | Definition |
|---|---|
| Source/library | A declared directory or pinned Git subdirectory containing reusable Markdown. |
| Source node | One addressable Markdown heading and its content. |
| Source reference | Explicit alias plus source path, such as `cdint:engineering/test-strategy`. |
| Manifest | The repository-authored `agents.yaml` output outline and selection. |
| Rendered output | Generated `AGENTS.md`. |
| Localization | Copy-on-write conversion of shared guidance into repository-owned content. |
| Provenance | The recorded origin identity and original content needed to explain a localized module. |
| Base | Exact content at the last localization or reconciliation point. |
| Origin | Current upstream content at an exact source revision. |
| Local | Current repository-owned customized content. |
| Preserve | Turn a generated-output edit into durable local source. |
| Inherit | Bring a reviewed origin change into localized content. |
| Propose | Prepare a local improvement for an upstream source. |

## Appendix B: Presentation Status Note

The basic composition, browsing, rendering, localization, drift-import, public
API, and immutable source-pin workflows are implemented. The library
relationship syntax, organization policy format, provenance-v2 storage,
directional reconciliation commands, automatic merge behavior, and forge
delivery are active design work, not completed product behavior.
