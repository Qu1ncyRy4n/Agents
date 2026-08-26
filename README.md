# Mogent

Mogent builds repository-specific `AGENTS.md` files from reusable Markdown
libraries and an `agents.yaml` manifest.

The short version:

- Libraries hold reusable Markdown modules.
- `agents.yaml` chooses, orders, and renames those modules for one repo.
- `mogent build` renders the manifest to `AGENTS.md`.

This is usable for local and immutably pinned URL libraries, manifest-driven
generation, copy-on-write localization, and conservative single-section drift
import. Broader package management and automatic merge workflows remain future
work.

## Install

From this repository:

```sh
tools/install
```

That installs the `mogent` CLI built from your current checkout. It will not
auto-update; run `tools/install` again after pulling or making code
changes.

Mogent contains no C code, but its terminal UI dependency chain reaches the Go
standard library's `os/user` package. With Go's default `CGO_ENABLED=1`, that can
make `go install` invoke `gcc` even though Mogent does not need CGO. The install
helper deliberately builds the portable pure-Go form. The equivalent manual
command is:

```sh
CGO_ENABLED=0 go install ./cmd/mogent
```

If you are already in `cmd/mogent`, `CGO_ENABLED=0 go install .` is equivalent.
From the repository root, `go install .` is not: the root has no Go package.
The Nix development shell includes GCC for diagnostics and dependencies that
genuinely require a C compiler, but Mogent's supported install does not require
it.

The repository check entrypoint defaults to the pure-Go build and keeps its
build cache under `/tmp`:

```sh
nix develop -c tools/check

# or, when go and errcheck are already installed
tools/check
```

## Basic Files

Create an editable starter manifest:

```sh
mogent init --list-templates
mogent init --template minimal --source shared=./libraries/cdint --dry-run
mogent init --template minimal --source shared=./libraries/cdint
```

Templates are ordinary manifest starting points, not hidden presets. `init`
validates their source references and rendered output before writing.
It materializes discovered `repo_name` and `repo_url` values under manifest
`vars`, where they remain reviewable and reproducible. Override either value
with a repeatable `--var name=value` flag.

A source library is a directory of Markdown files:

```text
library/
  identity/
    role.md
  instructions/
    workflow.md
  lang/
    go/
      testing.md
```

Directory paths contribute to source references. Filenames do not. For example:

```markdown
<!-- library/lang/go/testing.md -->
# Testing

Run focused Go tests for changed packages.
```

is referenced as:

```text
shared:lang/go/testing
```

because the source-relative directory is `lang/go` and the heading is
`Testing`.

File-level YAML frontmatter is supported for source browsing and filtering:

```yaml
---
tags: [lang/go, testing/unit]
tldr: Prefer focused Go tests and table tests.
priority: 0.9
scope: project
requires: [shared:instructions/workflow]
conflicts_with: [shared:security/loose-security]
---
```

Metadata is tool-only. It does not render into `AGENTS.md`.

The checked-in example libraries use this atomic style. Each reusable module
usually lives in one file with one top-level heading and frontmatter metadata.
For example, `libraries/cdint/shared-baseline/instructions/focused-change-loop.md`
resolves as `shared:shared-baseline/instructions/focused-change-loop`.

When a document wants several atomic modules under one local heading, name them
explicitly in the manifest:

```yaml
doc:
  - Instructions:
      - Decision First:
          - Decision Intent: shared:process/decision-first/decision-intent
          - Thought Experiment: shared:process/decision-first/thought-experiment
          - Open Questions: shared:process/decision-first/open-questions
```

Mogent does not currently stitch separate files into one implicit cross-file
subtree. That keeps source selection explicit while the library model settles.

## Manifest

`agents.yaml` is the authored document outline:

```yaml
sources:
  shared: library

vars:
  repo_name: example
  repo_url: https://example.test/example

output: AGENTS.md

doc:
  - Identity:
      - Role: shared:identity/role
  - Instructions:
      - Workflow: shared:instructions/workflow
      - Go Testing: shared:lang/go/testing
```

The manifest owns the rendered headings. Source headings help find content, but
the manifest decides where content appears in the final document.

## Organization Adoption

For a team repository, treat the manifest, lock, and rendered agent file as
reviewed project configuration:

- Commit `agents.yaml` so the selected guidance and document structure are
  explicit.
- Commit `mogent.lock.yaml` when URL sources are used. It pins the exact Git
  commit, source subdirectory, and Markdown content hash.
- Commit the generated `AGENTS.md` so agents can use it without installing
  Mogent and reviewers can see policy changes in ordinary diffs.
- Commit `.mogent/library/` and `.mogent/provenance.yaml` only when localized
  modules are intentional, repository-owned policy.

Ignore machine-local state and downloaded source checkouts:

```gitignore
.mogent/sources/
.mogent/state.json
```

Avoid ignoring all of `.mogent/` if the repository tracks localized modules.
Never commit credentials, private source contents, or other personal material
just because Mogent can resolve them locally.

Keep personal guidance in a separate, explicitly named source such as
`personal`. That source can be a local path or its own access-controlled Git
repository. A repository should only select personal modules that are suitable
for every intended collaborator; broadly applicable additions should go
through the organization's normal review process before becoming shared
policy.

For organization-owned libraries, prefer stable source aliases, narrow
repository `subdir` roots, immutable lock updates, and code review of both the
lock diff and changed rendered output. Updating a URL source is an explicit
`mogent source update` operation; ordinary builds remain offline and never
silently consume upstream changes.

## Common Commands

Build `AGENTS.md`:

```sh
mogent build
```

Use a non-default manifest path:

```sh
mogent build --manifest testing-ground/metadata-library/agents.yaml
```

Show workspace state:

```sh
mogent status
```

`status` prints an actionable hint. For example, missing or stale output points
back to `mogent build`, while direct edits warn before `--force`.

List available source modules:

```sh
mogent source list
mogent source list cdint
mogent source list personal --tree
mogent source list --search python --tldr
mogent source list --tag-search go
mogent source list --sort priority --metadata
```

Source directories are selectable subtrees. For example,
`personal:engineering` selects every module beneath `engineering/`; a manifest
entry may use `exclude` to remove a narrower directory or heading subtree.

`--tag-search` searches metadata tags only. Use `--search` to search source
references, headings, TLDRs, tags, and direct body text.

Inspect one source module:

```sh
mogent source show shared:lang/go/testing --metadata --content=snippet --lines=8
mogent source show shared:lang/go/testing --align-source --under Instructions
```

Add a local source without selecting any modules:

```sh
mogent source add personal ../agent-libraries/personal --dry-run
mogent source add personal ../agent-libraries/personal
```

Declare, pin, and later review a URL source:

```sh
mogent source add shared https://github.com/Qu1ncyRy4n/Agents.git \
  --subdir libraries/cdint --dry-run
mogent source add shared https://github.com/Qu1ncyRy4n/Agents.git \
  --subdir libraries/cdint
mogent source pin shared
mogent source update shared
mogent source update shared --ref <full-previewed-commit> --accept
```

Adding a source changes only `sources:`; it never selects content. Normal
commands never fetch URL sources. They require the committed
`mogent.lock.yaml` entry and verify the ignored `.mogent/sources/` checkout.
Compact scalar source values remain valid. An explicit `subdir` is locked and
makes source references relative to that selected repository directory.

Show included and unused source modules:

```sh
mogent source list --coverage --tree
mogent source list --coverage --tree --unused-only
mogent source list --coverage --tree --tldr
mogent coverage
```

`coverage` is a compatibility alias for `source list --coverage --tree`.
The unified tree distinguishes directories, headings, and combined
directory/headings, and shows text states such as `[included]`, `[inherited]`,
`[excluded]`, `[partial]`, and `[unused]`. Add `--tldr` to show each source
file's compact summary once beside its first visible node.

Source inventory output uses aligned ASCII by default. Personal presentation preferences can
be placed in `$XDG_CONFIG_HOME/mogent/config.yaml` (normally
`~/.config/mogent/config.yaml`):

```yaml
display:
  chars: unicode
  align: true
  fit: term
  width: 0
```

Use `--chars ascii|unicode`, `--align=false`, `--fit term|none`, or an explicit
`--width <columns>` for a one-command override. Terminal fitting wraps long
references and TLDRs onto labeled continuation lines; redirected output stays
unbounded unless a width is explicitly configured.
These settings only affect display; they never change the manifest or rendered
agent document.

Add a source module to the manifest:

```sh
mogent add shared:lang/go/testing --under Instructions
```

Preview before writing:

```sh
mogent add shared:lang/go/testing --under Instructions --dry-run --preview=tree
mogent add shared:lang/go/testing --under Instructions --dry-run --preview=patch
mogent add shared:lang/go/testing --before Instructions/Workflow --dry-run --preview=tree
```

`--under` inserts last by default and accepts `--first` or `--last`.
`--before` and `--after` use full manifest heading paths and infer the parent.
Directory tree previews list inherited source headings separately and flag
exact overlap or related existing selections for review.

Write the manifest and rebuild `AGENTS.md`:

```sh
mogent add shared:lang/go/testing --under Instructions --rebuild
```

Create a copy-on-write local override for one manifest entry:

```sh
mogent localize Instructions/Testing --dry-run
mogent localize Instructions/Testing --rebuild
```

Inspect or explicitly resolve direct edits to generated output:

```sh
mogent drift
mogent drift --import Instructions/Testing
mogent drift --reject --force
```

Localization writes ordinary Markdown under `.mogent/library`, records its
origin in `.mogent/provenance.yaml`, and changes the selected manifest reference
to `local:`. Drift import refuses edits outside the selected section.

Open the TUI:

```sh
mogent tui
```

Print completion candidates for shell wrappers:

```sh
mogent complete commands
mogent complete source-aliases
mogent complete source-refs --prefix shared:lang
mogent complete manifest-headings
mogent complete flags
```

Install shell completion by evaluating or saving the generated script:

```sh
mogent completion bash
mogent completion zsh
```

The TUI is useful for browsing and proof-of-concept editing, but the CLI is the
more reliable surface for basic usage today.

## Go API

Mogent's reusable core is available through public Go packages. The CLI uses
the same `workspace` operations available to other callers:

```go
session, err := workspace.New("agents.yaml")
if err != nil {
    return err
}
status, err := session.Status()
if err != nil {
    return err
}
preview, err := session.AddSource(workspace.AddOptions{
    Reference: "shared:lang/go/testing",
    Under:     "Instructions",
}, true) // true means dry-run
```

`workspace` is the supported orchestration facade. Packages such as `manifest`,
`library`, `render`, `sourcecache`, and `state` expose the lower-level model used
by that facade and the CLI. Mogent is still pre-v1; changes to exported package
contracts must be deliberate and documented.

## Safety Behavior

Mogent writes `agents.yaml` and `AGENTS.md` conservatively.

- `build` refuses to overwrite an output file that does not look generated by
  mogent unless `--force` is used.
- `add` writes `agents.yaml` by default.
- `add --dry-run` writes nothing.
- `add --rebuild` also writes `AGENTS.md` through the normal overwrite checks.
- `localize --dry-run` writes nothing; localization never changes a shared source.
- `drift --import` accepts only one unambiguous edited section.
- `drift --reject` requires `--force` before discarding direct edits.
- Shared source libraries are read-only inputs during normal build/add flows.

## Testing Ground

For a feature-by-feature exercise sequence and issue template, see
[`docs/DOGFOOD.md`](docs/DOGFOOD.md). For only the newest selectable-directory
and unified-tree behavior, use
[`docs/DOGFOOD-SESSION-3.md`](docs/DOGFOOD-SESSION-3.md).

Try the metadata-oriented sandbox:

```sh
mogent source list --manifest testing-ground/metadata-library/agents.yaml --tag-search security --sort priority --metadata
mogent build --manifest testing-ground/metadata-library/agents.yaml
```

If the checked-in sandbox `AGENTS.md` has no local mogent state yet, `build` may
ask for `--force`. That is the overwrite guard working as intended.

## Current Status

Ready for basic local usage:

- local Markdown source directories,
- strict `agents.yaml` manifests,
- deterministic `AGENTS.md` rendering,
- atomic module libraries with file-level metadata,
- source browsing,
- selectable directory subtrees and unified inventory/coverage trees,
- metadata/tag filtering,
- source coverage,
- simple manifest additions,
- dry-run previews,
- copy-on-write local overrides and conservative single-section drift import,
- immutable pinned HTTP(S) Git sources with explicit update review.

Not ready yet:

- authenticated or non-Git remote sources,
- global source cache,
- automatic or multi-section import from hand-edited `AGENTS.md`,
- conflict warnings from `conflicts_with`,
- multiple agent outputs such as `CLAUDE.md`, `GEMINI.md`, or `.codex/AGENTS.md`,
- split/import from a full `AGENTS.md` into an atomic library.

Future multi-output support should allow both exact mirrors and tool-specific
renders. For example, `CLAUDE.md` could include extra visible manifest entries
for stricter Claude behavior, while `GEMINI.md` might simply symlink to
`AGENTS.md`.

## Design Docs

- [docs/DESIGN.md](docs/DESIGN.md) is the design of record.
- [docs/DOGFOOD.md](docs/DOGFOOD.md) stages current features for dogfooding and issue reporting.
- [docs/DOGFOOD-SESSION-3.md](docs/DOGFOOD-SESSION-3.md) exercises only the new directory-tree and display work.
- [docs/HANDOFF.md](docs/HANDOFF.md) is the stable handoff navigation entry;
  its linked dated pickup is the detailed current resume point.
- [docs/AUTHORING-PLAN.md](docs/AUTHORING-PLAN.md) stages source declaration, move/reorder, and guided module creation.
- [docs/ORIGIN-RECONCILIATION-PLAN.md](docs/ORIGIN-RECONCILIATION-PLAN.md) stages preservation, inheritance, and proposal workflows.
- [TODO/TODO-jusuk-mogent-agent-modules.md](TODO/TODO-jusuk-mogent-agent-modules.md)
  tracks implementation and open decisions.
- [docs/thought-experiments/TE-kavam-metadata-tags-and-library-shape.md](docs/thought-experiments/TE-kavam-metadata-tags-and-library-shape.md)
  covers metadata, tags, conflict warnings, and atomic library shape.
