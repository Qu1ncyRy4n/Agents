# Mogent HCL V2 Specification

Status: active design specification for the v2 configuration model. Revised
2026-10-05 from the 2026-09 review with Steve. Supersedes the earlier draft's
`mogent.lock` file, `inclusion` block, and `build --dry-run` command. Current
implementation status is tracked in `TODO/TODO-pupab-hcl-v2-rollout.md`.

V1 (`agents.yaml`, `mogent.lock.yaml`, `mogent build`) is unchanged by this
document. A repository uses one format at a time. V1-to-V2 migration is a later
command, not an implicit parser fallback.

## Purpose

V2 lets a consuming repository select reusable library content without
recreating the library's hierarchy. A user starts from boilerplate, inspects
available content, makes small selection changes, previews the complete result
as a diff, and applies only after every required choice is resolved.

The terms *library side* and *consumer side* name the two roles. They do not
require a network service: a library may be a local directory or a Git
repository.

### Model

Promise Theory (Burgess) is the vocabulary. An agent can only promise its own
behavior; it cannot impose on another agent.

- The library makes `+` promises: it **offers** content and states how strongly
  it stands behind each section.
- The consumer makes `-` promises: it **selects** what it accepts.
- The library never imposes. Even a `foundation` section is an offer the
  consumer may decline, on record, with a reason.

Everything the consumer decides is visible in one committed file. Everything the
library recommends is visible in one committed sidecar. Nothing is merged
silently.

## Files

```text
mogent.hcl            Authored consumer configuration. Committed.
.mogent/state.json    Generated-output state. Ignored.
.mogent/sources/      Cached Git checkouts, keyed by alias and commit. Ignored.
library.mogent.hcl    Library sidecar, at the library root. Committed there.
```

There is no lock file. Git already content-addresses a library; Mogent records
the resolved commit in `mogent.hcl` itself, where the consumer can read it next
to the source declaration. A future Grid source will use its CID the same way.

## Consumer Configuration

The outer shape is sources followed by independent outputs.

```hcl
mogent {
  format = 2
}

sources {
  source "cdint" {
    git    = "https://github.com/example/cdint-lib.git"
    ref    = "main"
    commit = "a1b2c3d4e5f6789012345678901234567890abcd"
  }

  source "personal" {
    local = "../qmr-agents-library"
  }
}

outputs {
  output "agent-instructions" {
    paths = ["AGENTS.md", "CLAUDE.md"]
    kind  = "markdown"

    source "cdint" {
      from = "cdint:agents"

      select = {
        org = {
          foundation = true
          core       = { accept_defaults = true }
          optional   = { review = true }
        }
        else = "exclude"
      }
    }

    source "personal" {
      from     = "personal:agents"
      select   = { all = true }
      tags_any = ["lang/go"]
    }
  }

  output "agent-skills" {
    path = ".agents/skills/"
    kind = "tree"

    source "personal" {
      from   = "personal:skills"
      select = { all = true }
    }
  }
}
```

`path` defines one output. `paths` defines identical, independently managed
outputs. It does not create symlinks.

`source` blocks within an output are evaluated in file order. Each selects from
one library and appends its content as a separate root. Mogent does not merge
headings from different libraries by title.

### Sources And Pinning

A source is `local` or `git`, never both.

| Attribute | Meaning |
|---|---|
| `local` | Directory path, relative to `mogent.hcl`. Read as checked out right now. |
| `git` | HTTP(S) Git URL. |
| `ref` | Branch or tag the consumer intends to follow. Optional; default `HEAD`. |
| `commit` | Full 40-hex commit the consumer is pinned to. Written by Mogent. |
| `subdir` | Library root inside the repository. Optional. |

Pinning follows the Go modules model, not the Nix lock-file model:

- `commit` is the pin. When present, `plan` and `apply` read the library at
  exactly that commit, from `.mogent/sources/<alias>/<commit>/`, fetching it
  once if absent. Output is reproducible from `mogent.hcl` alone.
- When `commit` is absent, `plan` resolves `ref` over the network, prints the
  resolved commit, and warns that the source is unpinned. `apply` writes the
  resolved `commit` into `mogent.hcl`, editing only that attribute and
  preserving the rest of the file byte for byte.
- `mogent update [alias]` re-resolves `ref`, shows the content diff between the
  old and new commit, and on `--accept` rewrites `commit`. Nothing else moves a
  pin. A stale pin is never a build error; it is the consumer's choice.
- A `local` source cannot be pinned. When it is inside a Git work tree, `plan`
  prints its HEAD commit and whether the tree is dirty, as information only.
  `local` is for authoring loops; `git` is for reproducible consumption.

Resolution and fetch happen only in `plan`, `apply`, and `update`, and only when
`commit` is missing or its checkout is absent. No other command touches the
network.

## Library Sidecar

The sidecar owns the structural tree. Markdown supplies content; the sidecar
gives every node a stable name and attaches metadata. Tags are metadata, not a
second hierarchy and not rendered headings.

```hcl
library {
  format = 2
  id     = "qmr/agents"
  name   = "QMR Agent Library"
}

content {
  markdown_root = "agents"
}

sections {
  section "intro" {
    title = "Intro"
    tldr  = "Working relationship and instruction-conflict guidance."
    tags  = ["scope/core"]
    offer = "foundation"

    section "assistant-agent" {
      title  = "Assistant Agent"
      source = "intro/assistant-agent.md"
      tldr   = "Establish the developer-agent working relationship."
    }
  }

  section "workflow" {
    title = "Workflow / Process"
    tags  = ["scope/core", "topic/workflow"]
    offer = "choose"

    defaults = {
      "principled-code-and-tool-use" = true
      "git-commit-mode"              = true
    }

    section "principled-code-and-tool-use" {
      title  = "Principled Code And Tool Use"
      source = "workflow/principled-code-and-tool-use.md"
    }

    section "git-commit-mode" {
      title  = "Use An Explicit Git Commit Mode"
      source = "workflow/git-commit-mode.md"
      tags   = ["tools/git"]
    }
  }

  section "presentation" {
    title   = "Presentation"
    offer   = "optional"
    default = false

    section "caveman" {
      title  = "Caveman"
      source = "presentation/caveman.md"
      tags   = ["persona"]
    }
  }
}

trees {
  tree "skills" {
    root = "skills"

    entry "nix-development" {
      path = "nix-development"
      tldr = "Work safely in a Nix development environment."
      tags = ["lang/nix", "tools/nix"]
    }
  }
}
```

A branch section organizes its children. A leaf has `source` and owns Markdown
content. A sidecar must reject duplicate sibling names, duplicate canonical
paths, missing Markdown files, a `source` outside `markdown_root`, and a leaf
with children.

Tags and TLDRs are sidecar metadata. A tag on an ancestor is effective for every
descendant during tag matching; it is not copied onto leaves.

### Offer

`offer` is a flat attribute on a branch section. It states how strongly the
library stands behind the branch's direct children when the consumer selects
the branch broadly. Selecting one exact descendant is always allowed and never
requires resolving unrelated siblings.

| `offer` | Library says | Broad selection of the parent does | Consumer must |
|---|---|---|---|
| `foundation` | "You should not run without this." | Include every direct child. | Give `force_exclude = true` and a non-empty `reason` to drop one. |
| `choose` | "Decide each of these yourself." | Nothing implicit. | Say `true` or `false` per direct child, or `accept_defaults = true`. |
| `optional` | "Take it or leave it; here is my default." | Apply the section's `default`. | Nothing. Omitted children follow `default`. |
| `opt_in` | "Here if you ask; probably not for you." | Exclude every direct child. | Name a child explicitly. Mogent warns. |

A branch without `offer` inherits its parent's broad selection unchanged, as
the current planner does. Whether `offer` should become required is an open
question below.

`defaults` is required and complete for `choose`, forbidden otherwise. `default`
is required for `optional`, forbidden otherwise. A `foundation` or `opt_in`
branch carries neither.

`accept_defaults = true` delegates the current direct-child choices to the
library's `defaults`. `plan` prints each accepted decision as a warning so a
later library revision that changes a default is visible in the plan diff.

### Markdown Heading Contract

A leaf's source Markdown must begin at line one with a level-one heading equal
to the sidecar `title`. Mogent validates the heading, renders the title at its
sidecar-defined depth, and inserts only the body. Source files stay readable
alone; the sidecar remains the one authoritative hierarchy.

## Selection

`select` is a sparse mirror of the library's section tree, never a second
definition of it.

```hcl
select = {
  org = {
    foundation = {
      identity = true
      security = {
        force_exclude = true
        reason        = "Isolated fixture; no credentials or network."
      }
    }
    core = { review = true, "release-process" = false }
  }
  else = "exclude"
}
```

Within a selection object:

- `true` includes a node and its normally includable descendants.
- `false` excludes a node and its descendants.
- `else = "include"` or `else = "exclude"` is the default for unmentioned
  siblings at that level.
- `all = true` is shorthand for broad inclusion at that level.
- `exclude = { name = true }` removes descendants from a broad selection.
- `accept_defaults = true` is valid only under a `choose` branch.
- `force_exclude` and `reason` are valid only when dropping a `foundation`
  child.

Reserved keys: `all`, `else`, `exclude`, `accept_defaults`, `force_exclude`,
`reason`. A library cannot use them as section names.

### Tag Selection

Tag selection is additive to the structural overlay and applies to leaves.

```hcl
tags_all = ["lang/go", "workflow"]   # every query must match
tags_any = ["tools/git", "tools/nix"] # at least one must match
```

A query matches a tag exactly or as a slash-prefix: `lang/rust` matches
`lang/rust/strict`. Both forms may appear together; a leaf must satisfy
`tags_all` and `tags_any`. A query that matches nothing is a warning. V2 has no
wildcard tag expressions.

## Commands

The command model borrows Terraform: show, then write.

```text
mogent init [--config mogent.hcl] --source alias=local-path|url
  Write minimal boilerplate. Never overwrite.

mogent plan [--config mogent.hcl]
  Resolve sources, compile selection, render in memory, and print a unified
  diff of every output path against its current content. Print diagnostics.
  Write nothing except a cache fetch for a missing pinned commit. Exit non-zero
  on any error diagnostic.

mogent apply [--config mogent.hcl] [--force]
  Run plan. Refuse on errors. Write every output and its state transactionally.
  Write resolved commits into mogent.hcl for sources that lacked one. Refuse to
  replace an output that Mogent does not manage, or that was edited by hand,
  unless --force is given.

mogent update [alias] [--accept]
  Re-resolve ref for one or all git sources. Show commit and content diff.
  Rewrite commit only with --accept.

mogent source list [alias] [--tags] [--tldr]
mogent source show alias:path
  Browse the library tree with section paths, offers, defaults, tags, TLDRs,
  and provenance.
```

There is no `build` and no `--dry-run` in v2. `plan` is the dry run. V1 keeps
`mogent build`; the two formats never share a command name with different
semantics.

## Diagnostics

Every diagnostic has a stable code, a severity, the consumer-side location, the
library-side section path when relevant, the consequence, and an actionable
fix. `plan` prints suggested HCL for incomplete selections.

| Code | Severity | Meaning |
|---|---|---|
| MOGENT100 | error | `select` is not a valid selection object. |
| MOGENT101 | error | Output names a source that did not load. |
| MOGENT102 | error | Broad selection of a `choose` branch with undecided children. |
| MOGENT104 | error | `from` does not name the library's markdown root or a declared tree. |
| MOGENT105 | error | Tree output selection is not `{ all = true }`. |
| MOGENT106 | error | Selection names a child the library does not have. |
| MOGENT107 | error | `force_exclude` on a non-`foundation` child, or without `reason`. |
| MOGENT108 | error | `accept_defaults` under a branch that is not `choose`. |
| MOGENT109 | error | `commit` is not a full 40-hex SHA, or does not exist at `git`. |
| MOGENT201 | warning | A tag query matched no leaf. |
| MOGENT202 | warning | A source selects no content. |
| MOGENT203 | warning | A default was accepted for a `choose` child. |
| MOGENT204 | warning | An `opt_in` child was selected. |
| MOGENT205 | warning | A `foundation` child was force-excluded. |
| MOGENT206 | warning | A `git` source has no `commit`; `apply` will write one. |
| MOGENT207 | info | A `local` source's Git HEAD and dirty state. |
| MOGENT208 | warning | An output on disk differs from its recorded state (drift). |

Sidecar validation failures (unknown `offer`, incomplete `defaults`, missing
source file, bad heading) are load errors with the sidecar path and line, not
plan diagnostics.

A warning never prevents `plan`. `apply` ignores warnings unless a later
`--strict` flag is added.

## Output And Safety

- Markdown and tree outputs are planned together and written transactionally.
  A failed write restores every touched file and the state file.
- Output paths are repository-relative, non-overlapping, and safe from path
  traversal and symlink escape.
- Output state binds normalized path to content hash. A hand-edited output
  needs `--force` after inspection.
- Rendered Markdown preserves each library's hierarchy. Separate libraries do
  not merge by matching heading title.
- Rendered output never contains sidecar metadata.

## Acceptance Scenarios

The product-level user stories live in `docs/user-stories.md` (US-1..US-8)
and `docs/user-stories-ciwg.md`. The scenarios below are the v2 command-level
walkthroughs that the rollout must pass against the QMR library. Each names
the story it serves.

**New repository (US-1).** Run `mogent init --source qmr=../qmr-agents-library`. Run
`mogent source list qmr --tldr` to see what the library offers, with each
branch's `offer` visible. Edit `select`. Run `mogent plan`; it shows the full
`AGENTS.md` as an addition diff plus any `choose` branches left undecided, with
suggested HCL. Fix, re-plan, `apply`.

**Adopt an existing repository (US-4, first slice).** A repository already has a hand-written
`AGENTS.md`. `init` and `plan` run without touching it. `plan` shows the diff
between the current file and the generated one, so the author sees what they
would lose and gain. `apply` refuses without `--force`. The author moves
unique local content into a `local` overlay library or accepts the loss, then
applies with `--force`. Mogent now owns the file.

**Update a working repository (US-2).** The library moved ahead. `mogent update qmr`
shows the commit range and a diff of the rendered output under the new commit.
`--accept` rewrites `commit`. `plan` is now clean; `apply` writes. Review is
one ordinary Git diff: `mogent.hcl` changed one line, `AGENTS.md` changed
content.

**Share a library (US-5, CIWG story).** A library author adds `offer`, `tldr`, and `tags` to the
sidecar, pushes, and hands a consumer a URL. The consumer's `init` plus `plan`
produces a useful baseline with no further reading. See the threshold below.

## Sharing Threshold

Draft; owner to confirm. A library is ready to hand to another person when:

1. `git` sources with `commit` pinning work, so the consumer gets what the
   author tested.
2. Every branch in the sidecar has `offer`, and every leaf has `tldr`, so
   `source list` is self-explanatory.
3. Tags exist on at least language, tool, and scope axes, so `tags_any` is
   useful.
4. The four user stories above pass against the QMR library.
5. The spec and `mogent --help` agree.

Tag-based selection and a polished core library are necessary; a TUI is not.

## Deferred

- Template variables inside content (`{{ .repo_name }}`). Steve has asked for
  this. Scheduled for the feature set after this rollout, not before.
- Config imports and recursive composition. `init --template` covers per-org
  starting points for now.
- Conditionals, heredocs, arbitrary expressions.
- Procedural file or Markdown-section replacement.
- Wildcard tag queries.
- Automatic merging of same-named headings across libraries.
- A TUI. The command-line `init`, `source`, `plan`, `apply`, `update` flow must
  be complete first.
- Automatic V1-to-V2 conversion.
- Non-AGENTS outputs such as slide decks from Mogent-managed content. Raised
  in review; scope not yet defined.

## Open Questions

- Behavior of a branch with no `offer`. Options: treat as `optional` with
  `default = true` (library must say less), or require `offer` on every branch
  (library must say more, consumer never guesses). Leaning: require it, since
  the sidecar is small and the whole point is explicitness.
- Whether `plan` should accept `--out` to write the rendered result somewhere
  other than the configured paths, for review tooling.

## Supersession

- The earlier draft's `mogent.lock` file and `mogent lock` command are removed.
  V1's `mogent.lock.yaml` under `docs/IMPLEMENTATION-M4.md` and DI-fipam is a
  V1 contract and is untouched.
- The `inclusion { policy = ... }` block is replaced by the `offer` attribute.
  `baseline` becomes `foundation`; `explicit` becomes `choose`.
- `mogent build --dry-run` for v2 is replaced by `mogent plan`; `mogent build`
  for v2 is replaced by `mogent apply`.
- `docs/proposals/HCL-CONFIG-SPIKE.md` remains the historical rationale for
  choosing HCL.
