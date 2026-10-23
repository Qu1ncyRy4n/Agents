# Plan Files

Use an authored `mogent.plan` when the project wants to record selection intent
separately from the reproducible, resolved build inputs. `mogent sequence`
compiles the plan into `mogent.build`; `mogent build` accepts either entry file
and sequences a plan before installing outputs.

## Minimal plan

```text
plan {
  format = 1
}

source "qmr" {
  local = "../../cdint/qmr-agents-library"
}

include "qmr:intro"
include "qmr:skills" ["rust-development", "maintain-documentation"]
```

Create this shape with `mogent init --plan`, then connect a library with
`mogent source add qmr ../qmr-agents-library` and select modules with
`mogent module include qmr:intro`.

`source` makes a library available. `include` selects content from it; source
declarations alone do not install anything. Includes are ordered.

## Dynamic tag selection

An include can resolve matching descendants by tag:

```text
include "qmr:workflow" {
  tags_any = ["topic/workflow"]
  tags_none = ["experimental"]
}
```

`tags_all` requires every listed tag; `tags_any` requires at least one; and
`tags_none` excludes matches carrying any listed tag. `exclude` is a later,
path-scoped filter. Inspect a prospective query before writing it with
`mogent source search qmr:workflow --tag-any topic/workflow --add-config`.

Dynamic membership/audit policy (`auto_inherit`) is not implemented. See
[SOURCE-DISCOVERY.md](SOURCE-DISCOVERY.md) for the current tag-search and
selection behavior.

## Lifecycle

```sh
mogent init --plan
mogent source add qmr ../qmr-agents-library
mogent module include qmr:intro
mogent sequence --file mogent.plan
mogent build --file mogent.plan
```

`sequence` writes the generated `mogent.build` but does not install workspace
outputs. The generated build records resolved input hashes and must not be
hand-edited. Edit `mogent.plan`, then sequence again. Output baselines live in
`.mogent/state.json`, separate from both authored intent and the generated
build.

For a project that does not need separate intent and generated inputs, use an
authored `mogent.build` directly. Do not create both an authored build and a
plan in the same project: a generated build is owned by its plan, while an
authored build is the project entry file.

The detailed language contract and acceptance cases remain in
[MOGENT-FIRST-SLICE-SPEC.md](proposals/MOGENT-FIRST-SLICE-SPEC.md) and
[MOGENT-FIRST-SLICE-ACCEPTANCE.md](proposals/MOGENT-FIRST-SLICE-ACCEPTANCE.md).
