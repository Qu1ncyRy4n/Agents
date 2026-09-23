# Library Sidecar Design

Status: first-slice implementation.

## Goal

`library.mogent.yaml` is an optional, versioned library-root sidecar. It owns
library identity, exhaustive source-tree accounting, display titles, ordering,
and node relationships that do not fit cleanly in individual Markdown files.

The sidecar is not a second content format. Markdown remains the source of
rendered prose; frontmatter remains the home for small intrinsic module facts.
When a valid sidecar is present, its tree order and optional titles control
source-tree presentation during rendering without changing source paths or
Markdown bodies.

## Schema Identity

Use an explicit schema ID:

```yaml
schema:
  id: mogent/1
```

Do not use `current_mogent`: a sidecar must retain its meaning across Mogent
upgrades. When omitted during the first transition, Mogent may assume `mogent/1`
with a warning; `mogent lib init` writes the explicit form.

Future custom schemas require a registered adapter, for example
`example.org/acme-library/1`. Mogent must not accept arbitrary per-library field
ownership declarations: it could no longer validate or safely edit the library.

## Example

```yaml
schema:
  id: mogent/1

library:
  name: QMR Mogent Library
  version: 0.1.0
  description: Personal reusable Mogent library for dogfooding and future forks.
  maintainers:
    - name: Quincy Ryan
      email: quincymryan@gmail.com
  organization:
    kind: personal
    name: QMR
  repository:
    url: null
    default_branch: main
  license: null

content:
  markdown_roots: [agents]
  directory_roots: [skills]

tree:
  - source: agents/intro
    title: Intro

  - source: agents/workflow
    title: Workflow / Process
    children:
      - source: agents/workflow/decision-first
        title: Decision First
        children:
          - source: agents/workflow/decision-first/plan-ahead
            title: Plan Ahead

  - source: agents/constraints
    title: Constraints and Safety

groups:
  thought-experiment-trigger:
    mode: at-most-one
```

`tree` is a sequence, so its order is deliberate YAML sequence order rather than
mapping iteration order. `source` paths are library-root-relative. A source
alias may select a repository subdirectory, but that alias is manifest-local and
is distinct from the library root that owns the sidecar.

`content.markdown_roots` declares the only Markdown subtrees this root sidecar
loads and validates. `content.directory_roots` declares raw directory-copy
subtrees; Mogent does not parse their Markdown. Roots must be distinct, existing
ordinary directories and may not overlap. The configured Markdown root is
virtual: `agents` is not a tree node, so an `all` selection renders `Intro`,
`Workflow / Process`, and `Constraints and Safety` without an `Agents` wrapper.
Archives or reference material outside these roots are ignored by sidecar-backed
loading. A sidecar without `content` retains exhaustive whole-library behavior.

Every discovered directory and Markdown heading within declared Markdown roots
must appear once in the tree. With `content`, the root directories themselves
are virtual and cannot appear in the tree.
A node may be represented only by `source`; `title`, relationships, and other
properties are optional. This makes the sidecar exhaustive without requiring
metadata for every ordinary module.

## Field Ownership

The `mogent/1` schema fixes field ownership. It is not configurable per library.

| Field | Authoritative location |
|---|---|
| rendered prose and Markdown headings | Markdown source |
| `tags` | Markdown frontmatter |
| `tldr` | Markdown frontmatter |
| library identity/contact/repository/license | sidecar `library` |
| source-tree accounting/title/order | sidecar `tree` |
| `requires`, `conflicts_with`, `exclusive_group` | sidecar node |
| group cardinality | sidecar `groups` |

`priority` is deferred. Existing priority frontmatter remains readable for
compatibility but is not recommended for new modules.

The former free-form `scope` field is also deferred. Its intended meaning is
being refined into applicability and selection policy.

## Applicability And Selection

Applicability describes the intended audience or portability boundary of a node:

```yaml
applicability:
  level: universal # universal | organization | project | person | machine
```

It does not itself select a node. Selection strength is target policy:

- `required`: target policy requires selection.
- `default`: target policy selects it unless deliberately removed.
- `optional`: target policy selects it only explicitly.
- `forbidden`: target policy rejects it.

The first sidecar implementation does not add applicability or selection
strength. Future policy must keep library guidance separate from target
authority; an organization-wide module is not automatically required by every
consumer.

## Relationships

Relationships are node-centric:

```yaml
- source: agents/workflow/thought-experiment-risk
  requires:
    - agents/workflow/developer-decision-involvement-level
  exclusive_group: thought-experiment-trigger
```

`groups` holds only group-wide rules. `at-most-one` means members may not be
selected together. A target requirement to choose one member is separate policy,
not an inference from the group.

The first implementation validates paths, duplicate declarations, cycles, and
group membership. It does not automatically include requirements or enforce
relationships during rendering.

## HTML Comments

Keep HTML comments in library source for editorial/history notes. Markdown
rendering strips them by default. `--preserve-html-comments` is the future audit
override. Directory outputs copy files unchanged, including skill comments.

Structured metadata must not be encoded as YAML inside HTML comments.

## Commands

First slice:

```sh
mogent lib init .
mogent lib scan . --dry-run
mogent lib check .
mogent lib check --source qmr_agents
```

- `init` creates a previewed sidecar from a scan.
- `scan` inventories directories/headings/frontmatter and proposes missing tree
  nodes without writing.
- `check` validates sidecar/source consistency and field ownership.

Future mutation commands include `mogent lib module new`, `mogent lib section
add`, and `mogent lib move`. They must preview both source and sidecar changes.

## I/O Debugging

Future command:

```sh
mogent output explain AGENTS.md --io-lines
```

It reports renderer accounting for a selected source section: input lines,
frontmatter omitted, HTML comments stripped, heading changes, content emitted,
and output lines. It does not claim editorial provenance between an upstream
source and an independently adapted module.

## Migration

`mogent lib check` initially warns when relationship fields appear in
frontmatter or tags/TLDR appear in the sidecar. It does not rewrite existing
metadata or fail legacy libraries during the first release. A future migration
command may make those moves explicit and reviewable.
