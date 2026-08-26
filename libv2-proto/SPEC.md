# Provisional Protolibrary Vocabulary and Authoring Rules

Status: provisional. These terms make extraction consistent; they are not yet
an approved Mogent file format.

## Source Guide

A **source guide** is an existing agent-instruction document being examined for
useful material. In this project, the protected source guides live under
`docs/other_repo_agents/`.

Do not use “corpus” in user-facing module prose. It obscures whether the writer
means source guides, the codebase, documentation, or the emerging library.

## Module

A **module** is guidance that Mogent may select and render as a unit.

An **atomic module** has one top-level heading and the content belonging to that
heading. Atomic describes its selection boundary, not its length or quality.

A **canonical module** is the organization-owned, recommended home for a
subject. It may be a substantial multi-heading document with examples and links
when those sections are intended to travel together. “Canonical” describes
authority, not file size.

A **compatible bundle** is a multi-heading module whose sections are intended
to be selected together. Do not split it merely to make every heading a file.

A **choice group** is a named set of modules for which selection is constrained,
such as choosing one communication persona. It is not an organization and does
not imply that its members live in separate libraries.

An **overlapping module** covers substantially similar behavior to another
module. Overlap alone does not say whether both may be selected. Record the
actual relationship separately.

## Plain-Language Selection Relationships

Use these terms in review prose and diagnostics. Machine-field names remain
provisional.

| Plain-language term | Meaning | Logical analogy | Provisional field |
|---|---|---|---|
| requires all | Selecting A also requires every named module. | AND | `requires_all` |
| requires one choice | Selecting A requires one named option. If `none` is listed, it is an explicit valid choice. | one-of | `requires_choice` |
| mutually exclusive | The named modules must not be selected together. | NOT-AND / NAND constraint | `mutually_exclusive` |
| exactly one | Select one listed option, including `none` only when it is explicitly listed. | XOR for two choices; one-of for many | `choose: exactly_one` |
| any combination | Select zero, one, or several listed modules. | unconstrained subset | `choose: any` |
| optional | This module is not required by the containing bundle or preset. | zero-or-one | `optional: true` |
| see also | Related for discovery; selection has no effect. | none | `see_also` |

Avoid using XOR as the canonical term because its meaning becomes unclear with
more than two choices. Avoid using “family” when “choice group” or “overlapping
modules” says what is actually meant.

Compatibility is the default. Two modules need no relationship metadata merely
because both may be selected.

## Editorial Status Is Separate from Selection

Use editorial status to track review confidence:

- `reliable`: strong source-supported candidate;
- `needs-choice`: competing policies require owner selection;
- `needs-rewrite`: useful source material is not yet reusable prose;
- `repository-local`: useful only with repository-specific facts;
- `generated-candidate`: proposed by a generated assessment rather than an
  authoritative source guide;
- `proposed`: invented during v2 design;
- `rejected`: retained in the ledger with a reason but not a module candidate.

These statuses do not make a module optional, required, or incompatible.

## Possible Output Destinations

Extraction records useful material even when `AGENTS.md` may not be its best
eventual output:

- `agent-guide`: instructions that should render into `AGENTS.md`;
- `skill`: a triggered operational workflow with inputs, outputs, and tools;
- `generated-doc`: reference material better generated for humans or agents;
- `specification`: product or protocol truth that instructions should point to;
- `repository-local`: project facts, commands, paths, or active state;
- `library-authoring`: guidance for building and maintaining Mogent libraries.

Use `candidate_outputs` metadata during intake when more than one destination
is plausible. Do not discard useful content merely because it may later become
a skill or specification.

## Module Includes

`TE-nufad` rejects rich YAML embedded in Markdown comments as the default
authoring model. It retains three decision-ready directions:

1. manifest composition plus separate modules as the default;
2. frontmatter declarations plus a minimal named body marker when exact
   insertion is necessary; or
3. a restricted literal Go-template include function.

Keep choice logic, relationships, and rich mappings out of body markers. The
owner must decide whether exact mid-document insertion is required for the first
libv2 cutover before either marker or template syntax is prototyped.

Required safety properties:

- reject malformed frontmatter, duplicate keys, duplicate import names, and
  unknown fields;
- resolve only declared sources and paths beneath their validated roots;
- retain existing symlink rejection for local and pinned source content;
- perform no network access during ordinary render;
- detect direct and indirect import cycles before rendering;
- impose a defensible maximum nesting depth and return the complete cycle or
  chain in the diagnostic;
- reject ambiguous matches rather than choosing the first;
- validate requirements and mutual exclusions across imported content; and
- keep the selected source visible in manifest and diagnostic output.

## Directories and Documents Carry Selection Meaning

Provisional rule: content intended to travel together stays in one compatible
bundle. Content intended for deliberate independent selection becomes a
separate module under a directory. Directory selection may inherit descendants;
selecting one document selects its complete compatible contents.

This makes physical organization behaviorally meaningful and requires a future
thought experiment against the current directory-plus-heading model before it
becomes product behavior.

## Deferred Decision Work

Schedule focused thought experiments and decisions for:

1. whether exact mid-document insertion is needed beyond manifest composition,
   then minimal marker versus restricted Go-template include if it is;
2. the smallest sufficient choice vocabulary and whether richer Boolean
   expressions are ever necessary;
3. directory inheritance and compatible-bundle boundaries;
4. source-path moves, compatibility pointers, and automatic migration policy;
5. stable published-library identity if real cross-library relationships need
   more than manifest aliases; and
6. whether metadata belongs in frontmatter, heading-local blocks, a library
   descriptor, or a narrowly defined combination.
