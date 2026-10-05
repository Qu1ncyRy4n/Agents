# YAML pseudo / sandbox

Scratch space for the two YAML surfaces Mogent reads. These examples are
accurate to the current schema (`manifest/manifest.go`, `library/sidecar.go`) so
they parse, then a clearly-marked **sketch** section to play with redesigns.

Terminology used below:

- **in-repo YAML** = `agents.yaml`, the repository-local build spec. Chooses,
  orders, and renames modules. Roles: `sources`, `vars`, `output`, `outputs`,
  `doc`.
- **library YAML** = `library.mogent.yaml`, the optional sidecar at a library
  root. Owns library identity, tree order, and relationships. Roles: `schema`,
  `library`, `content`, `tree`, `groups`.
- A **reference** is `alias:heading/path`. The alias names a source; the rest is
  a heading path inside that source.
- An **entry** in `doc` is one output heading. Compact form is
  `Heading: alias:path`; composed/nested form is
  `heading:` + `from:` / `exclude:` / `children:`.

---

## 1. In-repo YAML — current, accurate

`agents.yaml`

```yaml
# Named remote/local roots reusable across sources (optional).
roots:
  org: ../org-agents-library
  pinned: https://github.com/example/agents-library.git

sources:
  # Compact form: alias -> location string.
  shared: ../org-agents-library
  # Options form: location + subdir (index only part of a repo).
  go:
    location: ../org-agents-library
    subdir: go
  # Pinned remote (immutable revision required for URLs).
  personal: https://github.com/me/personal-agents.git@v1.2.0
  # Root-relative form keeps the root map separate from the subdir.
  research:
    path:
      root: ../org-agents-library
      subdir: research

vars:
  repo_name: my-project
  repo_url: https://github.com/me/my-project

# Primary rendered document. Defaults to AGENTS.md when omitted.
output: AGENTS.md

doc:
  - Identity:
      - Role: shared:shared-baseline/identity/role
      - Source Of Truth: shared:shared-baseline/identity/source-of-truth
  - Instructions:
      - Focused Change Loop: shared:shared-baseline/instructions/focused-change-loop
      - Go Development: go:go/development
      # Composed entry: rename the output heading and merge several sources.
      - Decision First:
          from:
            - shared:process/decision-first/decision-intent
            - shared:process/decision-first/thought-experiment
          exclude:
            - shared:process/decision-first/legacy-variant
          children:
            - Open Questions: shared:process/decision-first/open-questions
  - Constraints:
      - Safe Defaults: shared:shared-baseline/constraints/safe-defaults
      - Research Data Protection: research:ucd-research/data-protection

# Additional generated targets. Markdown (.md) renders the doc; a trailing /
# copies a source tree raw.
outputs:
  - path: .agents/skills/
    include:
      - tags:
          any: [skill]
```

Things to notice / that tend to confuse:

- `doc` is both the **document outline** and the **selection list**. There is no
  separate "include" layer.
- `Heading: alias:path` means "render this source under this heading". The
  source's own headings are not the output structure; the manifest is.
- The source alias (`shared`) and the module path (`process/...`) are glued into
  one string with a colon.
- `from`/`exclude`/`children` are the escape hatch when the compact form is too
  small; `heading:` is the explicit form of the same thing.
- `outputs` reuses a different reference vocabulary (`all`, `source`, `tags`).

---

## 2. Library YAML — current, accurate

`library.mogent.yaml` (library root)

```yaml
schema:
  id: mogent/1        # optional; defaults to mogent/1

library:
  name: cdint-shared
  version: 0.3.0
  description: Shared agent modules for CDINT projects
  maintainers:
    - Quincy Ryan
  organization: CDINT
  repository: https://github.com/example/agents-library
  license: UNLICENSED

# Limit this root to explicitly publishable subtrees. Directory roots are
# copied raw by directory outputs and are not parsed as Markdown.
content:
  markdown_roots:
    - shared-baseline
    - process
    - lang/go
  directory_roots:
    - skills

# The authoritative tree: order, nesting, rename, relationships.
# Every discovered Markdown heading must appear exactly once.
tree:
  - source: shared-baseline
    title: Baseline
    children:
      - source: shared-baseline/identity/role
      - source: shared-baseline/identity/source-of-truth
      - source: shared-baseline/instructions/focused-change-loop
  - source: process/decision-first
    children:
      - source: process/decision-first/decision-intent
      - source: process/decision-first/thought-experiment
  - source: process/decision-first/thought-experiment
    requires:
      - process/decision-first/decision-intent
  - source: lang/go
    children:
      - source: lang/go/development
      - source: lang/go/testing
  - source: lang/go/testing
    conflicts_with:
      - lang/go/testing-legacy
    exclusive_group: go-test-strategy

groups:
  go-test-strategy:
    mode: at-most-one
```

Things to notice / that tend to confuse:

- `library.mogent.yaml` and `agents.yaml` are two different documents with
  different shapes for overlapping ideas (order, selection, relationships).
- `tree[].source` is a **path** (e.g. `lang/go/testing`), not `alias:path`.
- `tags` and `tldr` are accepted here only so `check` can flag them; those
  belong in Markdown frontmatter, not the sidecar.
- `requires` / `conflicts_with` / `exclusive_group` exist here but not (as
  first-class fields) in `agents.yaml`.

---

## 3. Sketch — idealized sandbox (NOT a spec)

Deliberately different, for playing with wording and shape. Do not implement
from this section.

`agents.yaml` (sketch)

```yaml
version: 1
name: my-project

sources:
  org: ../org-agents-library
  me: ../personal-agents

# Pull a named starting point, then add or drop from it.
from: org/shared-baseline
add:
  - org/identity/role
  - org/instructions/focused-change-loop
  - me/style/accessibility
  - when: [go]
    add: org/lang/go/development
drop:
  - org/process/decision-first/legacy-variant

render:
  AGENTS.md: default
  CLAUDE.md: default
```

`library.mogent.yaml` (sketch)

```yaml
version: 1
name: cdint-shared
publish:
  markdown: [shared-baseline, process, lang/go]
  raw: [skills]

modules:
  - shared-baseline/identity/role
  - shared-baseline/instructions/focused-change-loop
  - lang/go/development

rules:
  - module: process/decision-first/thought-experiment
    needs: process/decision-first/decision-intent
  - module: lang/go/testing
    one-of-group: go-test-strategy
```

Open design questions to argue about while playing:

1. Should selection live in the manifest only, and the library sidecar only
   describe the library? (Today both can express structure.)
2. Should the source alias stay in the reference string, or become a separate
   field (`source: org` + `module: lang/go/development`)?
3. Should there be one reference syntax for `doc` and `outputs`, or is the split
   intentional?
4. Should "use everything under this heading" and "use this one heading" be the
   same keyword?
5. Should `requires` / `conflicts_with` / groups be resolvable at render time
   (blocking/auto-include), or advisory only?
