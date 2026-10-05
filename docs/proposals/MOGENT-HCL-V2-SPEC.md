# Mogent HCL V2 Specification

Status: draft design specification. It describes a new configuration model; the
current implementation reads `agents.yaml` and does not implement this format.

## Purpose

V2 makes a user-side repository select reusable server-side library content
without recreating the library's hierarchy. A user can start with boilerplate,
inspect available content, make small selection changes, preview a complete
plan, and build only after every required choice is resolved.

The terms *server side* and *user side* describe the two roles in this document.
They do not require a network service: a server-side library may be local or
Git-pinned.

## Files

```text
mogent.hcl        Authored user-side configuration, committed.
mogent.lock       Generated resolved source pins and content hashes, committed.
.mogent/state     Generated-output state, normally ignored.
```

`mogent.hcl` is the only authored build configuration. `mogent.lock` makes its
external inputs reproducible; it is not an alternate procedural configuration.
The derived build plan is shown by `mogent plan`, not stored as a second config
file.

V1 `agents.yaml` and V2 `mogent.hcl` are separate formats. A repository uses one
format at a time. V1-to-V2 migration is a later command, not an implicit parser
fallback.

## User-Side Configuration

The outer shape preserves the current manifest model: sources followed by
independent outputs.

```hcl
mogent {
  format = 2
}

sources {
  source "cdint" {
    git = "https://github.com/example/cdint-demo-lib.git"
    ref = "main"
  }

  source "personal" {
    local = "../qmr-agents-library"
  }
}

outputs {
  output "agent-instructions" {
    paths = ["AGENTS.md", "CLAUDE.md", "GEMINI.md"]
    kind  = "markdown"

    source "cdint" {
      from = "cdint:agents"

      select = {
        org = {
          baseline = true

          core = {
            accept_defaults = true
          }

          optional = {
            review = true
          }
        }

        else = "exclude"
      }
    }

    source "personal" {
      from = "personal:agents"

      select = {
        accessibility = true
        style         = true
        else          = "exclude"
      }

      tags_all = ["lang/go", "workflow"]
    }
  }

  output "agent-skills" {
    path = ".agents/skills"
    kind = "tree"

    source "personal" {
      from = "personal:skills"

      select = {
        all = true

        exclude = {
          experimental = true
        }
      }
    }
  }
}
```

`path` defines one output. `paths` defines identical independently managed
outputs. It does not create symlinks.

`source` blocks within an output are evaluated in file order. They select content
from the named server-side library and append it as separate source roots. Mogent
does not merge headings from separate server-side trees by title.

## Server-Side Library Sidecar

The server-side sidecar owns the structural tree. Markdown supplies content; the
sidecar gives every structural node a stable identity and attaches metadata where
it belongs. Tags are metadata, not a second hierarchy and not rendered headings.

```hcl
library {
  format = 2
  id     = "cdint/demo"
  name   = "CDINT Demo Library"
}

content {
  markdown_root = "agents"
  tree_root     = "skills"
}

sections {
  section "org" {
    title = "Organization"
    tags  = ["org"]

    section "baseline" {
      title = "Baseline"

      inclusion {
        policy = "baseline"
      }

      section "identity" {
        source = "org/baseline/identity.md"
        tldr   = "Describe the expected working relationship."
        tags   = ["identity"]
      }

      section "security" {
        source = "org/baseline/security.md"
        tldr   = "Treat credentials and irreversible actions carefully."
        tags   = ["security"]
      }
    }

    section "core" {
      title = "Core Practices"

      inclusion {
        policy = "explicit"

        defaults = {
          review            = true
          "release-process" = false
        }
      }

      section "review" {
        source = "org/core/review.md"
        tldr   = "Keep changes reviewable and verify changed behavior."
        tags   = ["workflow", "review"]
      }

      section "release-process" {
        source = "org/core/release-process.md"
        tldr   = "Use the release process when publishing a version."
        tags   = ["workflow", "release"]
      }
    }

    section "optional" {
      title = "Optional Practices"

      inclusion {
        policy  = "optional"
        default = false
      }

      section "issue-triage" {
        source = "org/optional/issue-triage.md"
        tldr   = "Triage incoming issues consistently."
        tags   = ["workflow", "triage"]
      }
    }
  }
}

trees {
  tree "skills" {
    root = "skills"

    entry "review" {
      path = "review"
      tags = ["workflow", "review"]
    }

    entry "experimental" {
      path = "experimental"
      tags = ["experimental"]
    }
  }
}
```

A branch section organizes its children. A leaf normally has `source` and owns
Markdown content. A sidecar must reject duplicate sibling names, duplicate
canonical paths, missing Markdown files, and a `source` outside `markdown_root`.

Tags and TLDRs are sidecar metadata. Tags on an ancestor are effective for all
descendants during tag matching; they do not need to be copied onto leaves.

### Markdown Heading Contract

The sidecar owns the rendered hierarchy. A leaf section's source Markdown must
begin at line one with an exact level-one heading matching the sidecar `title`:

```hcl
section "review" {
  title  = "Review"
  source = "org/core/review.md"
}
```

```markdown
# Review

Keep changes reviewable and verify changed behavior.
```

Mogent validates the heading, renders `Review` at its sidecar-defined depth, and
inserts only the source body. This keeps source files readable alone while
preventing duplicate headings and preserving one authoritative hierarchy.

## Inclusion Policies

An `inclusion` block is server-side guidance applied only when the user side
selects that branch broadly. Selecting one exact descendant is always allowed
and never requires resolving unrelated siblings.

| Policy | Default behavior when its parent is selected | User-side behavior |
|---|---|---|
| `baseline` | Include direct children. | Removing a child needs `force_exclude = true` and a non-empty `reason`. |
| `explicit` | No implicit child choice. | Every current direct child needs `true` or `false`, unless `accept_defaults = true`. |
| `optional` | Use the section's `default`. | Omitted children do not need decisions. |
| `opt_in` | Exclude direct children. | A child must be explicitly selected; broad selection does not include it. |

`accept_defaults = true` delegates the current direct-child choices to the
server-side defaults. `mogent plan` must show each accepted decision. A source
revision update that changes resolved content remains visible through its lock
and generated plan/output diff.

## Selection Overlay

`select` is a sparse mirror of the server-side section tree. It is not a second
definition of that tree.

```hcl
select = {
  org = {
    baseline = {
      identity = true

      security = {
        force_exclude = true
        reason        = "This isolated local fixture contains no credentials or network access."
      }
    }

    core = {
      review            = true
      "release-process" = false
    }
  }

  else = "exclude"
}
```

Within a selection object:

- `true` includes a node and all normally includable descendants.
- `false` excludes a node and descendants.
- `else = "include"` or `else = "exclude"` supplies the ordinary default for
  unmentioned siblings at that level.
- `all = true` is shorthand for broad inclusion at that level.
- `exclude = { name = true }` removes descendants from a broad `all` selection.
- `accept_defaults = true` is valid only for a server-side `explicit` section.
- `force_exclude` and `reason` are valid only when removing a `baseline` node.

Reserved selection keys are `all`, `else`, `exclude`, `accept_defaults`,
`force_exclude`, and `reason`. A server-side library cannot use those names as
section identifiers.

Tag selection is additive to the structural overlay:

```hcl
tags_all = ["lang/go", "workflow"]
tags_any = ["tools/git", "tools/nix"]
```

`tags_all` requires every query and `tags_any` requires at least one query. An
exact tag query also matches descendants: `lang/rust` matches `lang/rust/strict`.
V2 does not support wildcard tag expressions such as `*/nix`; users select known
prefixes or list exact roots until a real use case establishes safe wildcard
semantics.

## Locking

`mogent lock` resolves every declared source and writes `mogent.lock` in a
canonical generated HCL form:

```hcl
lock {
  format = 1
}

source "cdint" {
  declared = {
    git = "https://github.com/example/cdint-demo-lib.git"
    ref = "main"
  }

  resolved = {
    commit         = "a1b2c3d4e5f6789012345678901234567890abcd"
    content_sha256 = "sha256-..."
  }
}
```

`mogent build` verifies that each declaration matches its lock entry, resolves a
remote Git source only at the locked commit, and verifies the selected content
hash. Local sources record their selected-content hash and require `mogent lock`
after intentional library changes. Missing, stale, or mismatched lock data is a
build error, never a network fallback.

## Guided Commands

```text
mogent init
  Creates minimal mogent.hcl boilerplate and offers declared source choices.

mogent lock
  Resolves declared sources and updates mogent.lock.

mogent source list|show|explain
  Browses the server-side tree with section paths, TLDRs, tags, inclusion
  policies, defaults, and source provenance.

mogent plan
  Resolves the user-side overlay without writing outputs. Shows selected,
  excluded, defaulted, unresolved, and tag-selected content plus target paths.

mogent build
  Runs the same plan and writes only when no errors remain.
```

`mogent init` should produce a small valid configuration. A user adding one exact
section from a server-side library must not be forced to resolve unrelated
`explicit`, `optional`, or `opt_in` branches.

## Diagnostics

Every diagnostic has a stable code, severity, user-side source location,
server-side canonical section path when relevant, consequence, and an actionable
fix. `mogent plan` prints suggested HCL patches for incomplete selections.

Example:

```text
error[MOGENT102]: incomplete selection for cdint:agents > org > core

The server-side section requires a decision for each direct child.
Suggested user-side configuration:

  core = {
    review            = true
    "release-process" = false
  }

Use `accept_defaults = true` instead to adopt these current defaults.
```

Build errors include invalid HCL, unknown fields, duplicate names, invalid
canonical paths, missing sidecar content, missing/stale lock entries, lock hash
mismatch, selection references to missing nodes, unresolved explicit choices,
invalid forced exclusions, unsafe paths, output overlap, and unmanaged output
replacement without explicit force.

Warnings include zero-match tag queries, no-op selections, accepted defaults,
forced baseline exclusions, selection of opt-in content, and output drift. A
warning never prevents a user from running `mogent plan`; build behavior for
warnings remains non-blocking unless a later strict mode is explicitly added.

## Output And Safety Rules

- Markdown and raw-tree outputs are planned together and written transactionally.
- Output paths must be repository-relative, non-overlapping, and safe from path
  traversal and symlink escapes.
- Existing unmanaged or directly edited outputs require explicit force after
  inspection.
- Output state remains bound to normalized output paths and content hashes.
- Rendered Markdown preserves each selected server-side hierarchy. Separate
  server-side roots do not merge by matching heading title.

## Deferred

- Config imports and recursive composition.
- Variables, templates, conditionals, heredocs, and arbitrary expressions.
- Procedural file or Markdown-section replacement.
- Wildcard tag queries and general query expressions.
- Automatic merging of same-named headings from different server-side libraries.
- A TUI editing workflow. The command-line `init`, `source`, `plan`, and `build`
  flow must be complete before a new interactive editor is designed.
- Automatic V1-to-V2 conversion.

## Implementation Slices

1. Parse and validate HCL configuration and lock files with source locations.
2. Load and validate the server-side hierarchical sidecar and inline metadata.
3. Compile user-side selection into a source-ranged plan and diagnostics.
4. Render Markdown and tree outputs from the plan.
5. Implement locking and content verification.
6. Integrate transactional writes, drift, status, and generated state.
7. Implement `init`, `source`, `plan`, and migration support.

Each slice requires focused fixtures that cover a successful minimal user-side
configuration, exact single-section borrowing, explicit-choice completion,
accepted defaults, forced baseline exclusion, tag matching, lock mismatch, and
multi-output rollback.
