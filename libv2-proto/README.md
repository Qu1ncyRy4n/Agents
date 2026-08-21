# Library V2 Protolibrary

Status: lossless intake and review workspace. This tree does not replace
`libraries/` and is not yet a stable Mogent source.

## Preserve Before Pruning

The protolibrary exists to retain every useful, generalizable instruction from
the protected guides under `docs/other_repo_agents/` before editorial cleanup.
Do not summarize away a condition, example, exception, competing opinion, or
workflow step merely because a shorter shared principle exists.

Every source instruction must end in one of these recorded states:

- represented verbatim or with an explicitly marked adaptation;
- retained as a competing option;
- classified as repository-local;
- classified as generated speculation rather than source authority;
- rejected with a recorded reason; or
- awaiting review.

Silence is not a coverage state.

## Provisional Physical Organization

```text
libv2-proto/
  README.md
  PROVENANCE.md
  proposals/
  orgs/
    <owning-library>/
      <semantic-area>/
        <compatible-bundle>.md
        <swappable-family>/
          <option>.md
```

The first directory axis is the organization or library that owns and
publishes the guidance. Semantic areas below it use a loose recommended
vocabulary such as `identity`, `workflow`, `process`, `engineering`,
`documentation`, `communication`, `security`, `language`, and `domain`.

Source provenance is independent of destination ownership. A module owned by
`cdint` may have evidence from several repositories or organizations, and each
source remains cited.

## Module Boundary

Keep guidance in one large Markdown document when its sections are compatible
and intended to be selected, reviewed, and updated together. Split a child
into its own file or directory when it is:

- optional independently of the parent;
- useful in another combination;
- mutually exclusive with a sibling;
- governed by a different owner or stability promise; or
- large enough to have a meaningful lifecycle of its own.

File size alone is not a reason to split. A heading is not automatically an
independently selectable module.

## Titles State the Subject and Direction

Prefer an instructive claim that can stand alone:

- `Use Thought Experiments to Narrow a Broad Design Space`
- `Keep External-Service Tests Offline by Default`
- `Preserve Explanatory Comments Across Refactors`

Avoid category-only or extraction-facing titles such as `Code`, `Testing`,
`Corpus Variants`, or `Workspace Fit`. Parent directories and tags provide the
category; the visible title tells a selector what opinion the content adds.

Subheadings should follow the same rule when practical. A short label is fine
when the parent already supplies the missing claim and the subsection is not
independently selectable.

## Required Metadata From Intake

Every candidate begins with frontmatter, even while fields remain provisional:

```yaml
---
status: reliable | needs-choice | needs-rewrite | repo-local | proposed
tags: [process/decision, artifact/thought-experiment]
tldr: Use scenario analysis to reduce a broad design space to a few viable options.
sources:
  - docs/other_repo_agents/example_AGENTS.md#thought-experiment-protocol
requires: []
conflicts_with: []
exclusive_group:
see_also: []
---
```

`status` describes editorial confidence, not whether Mogent should render the
module. `optional`, hard requirements, conflicts, exclusive groups, and presets
are composition facts and must not be inferred from editorial confidence.

Tags are added during intake. Tags aid discovery; they are not node identity
and do not silently select or suppress content.

## TLDRs at Both Selection and Reading Boundaries

File frontmatter contains the selection-time `tldr`. Long internal sections
may also start with a short HTML comment:

```markdown
## Compare the Surviving Designs Under the Same Scenarios
<!-- tldr: Hold assumptions constant and expose each option's obligations. -->
```

The comment is authoring and reading metadata. Whether it renders into the
final document remains a separate Mogent decision. Do not hide instructions
only inside the comment.

## Examples Are Normal Content

Use worked examples whenever a rule is abstract, easy to misread, or changes a
tool invocation, artifact format, failure path, or user interaction. Prefer a
small realistic example over another paragraph of abstraction. Preserve useful
source examples verbatim with citations before rewriting them.

Each substantial module should normally include at least one of:

- a before/after example;
- a success and failure example;
- a concrete command and expected result;
- a worked decision or interaction;
- a counterexample showing the tempting wrong interpretation.

Do not invent an example and present it as source-derived. Label new examples
`Proposed example` until reviewed.

## Provenance and Verbatim Material

Use ordinary relative Markdown links to the protected source guides. Each
verbatim block or closely adapted section must identify its source immediately
before or after the material. When several guides repeat the same text, cite a
primary captured instance and list the duplicates in `PROVENANCE.md`.

Obsidian wiki links may be offered as a generated view later. Ordinary Markdown
links are canonical because they work in GitHub, editors, and rendered docs.

Source text and editorial text must be visibly distinguishable:

- `Verbatim source` means wording is preserved apart from Markdown-level
  relocation.
- `Adaptation` means wording changed and the original remains linked.
- `Proposed` means the content was invented during library design.

## Relationship Semantics

The protolibrary records these meanings without relying on source order:

- `requires`: every named target must also render;
- `requires_one_of`: at least one named target must render;
- `conflicts_with`: the exact pair cannot render together coherently;
- `exclusive_group`: at most one selected member of the group may render;
- `see_also`: discovery only;
- `sources`: evidence and provenance only.

Compatible content is the default and needs no relationship tag. “Optional” is
a discoverability or preset property, not the opposite of `requires`.

No first-match or first-file-wins behavior is allowed. A violated hard
relationship should produce a warning during migration and an error once its
syntax is stable and the library has been validated.

## Identity and Cross-Library References

A node needs an address that survives use in different consuming manifests.
Manifest aliases such as `org2:` are local bindings, not the permanent identity
of a library. The provisional distinction is:

- inside one library, `self:path/to/node` identifies a same-library target;
- inside a consuming manifest, `org2:path/to/node` uses the alias declared by
  that manifest;
- a source declaration binds `org2` to a local path or immutable Git source.

For example:

```yaml
sources:
  org2:
    location: https://github.com/example/agent-library.git
    subdir: libv2
```

```yaml
doc:
  - Persona: org2:communication/personas/surfer
```

The unresolved identity question is whether published libraries also need a
stable name independent of both manifest alias and URL. That affects forks,
renames, and cross-library relationships and remains a proposal, not a settled
format.

## Direct Includes Are Explicit Composition

A Markdown-level import can be useful when an authored bundle has a fixed
insertion point but delegates a choice to the manifest. It must not embed a
logic language in a path string or silently pick an option.

The desired behavior is recorded in `proposals/direct-imports.md`. Until Mogent
implements and validates such imports, candidate modules use metadata and
examples rather than executable `[[import:...]]` directives.

## Useful Ideas That Are Not Source-Derived

Do not mix invented modules into corpus coverage. Record potentially useful
future content under `proposals/upcoming-content.md`; promote it only after
review and label its origin when it becomes a module.

