# Mogent

**Conversation NOW: [docs/convo/head.md](docs/convo/head.md).** Start here to resume;
older discussion and checkpoints are separated under `docs/convo/archive/`.

Mogent composes repository-specific agent instructions and file bundles from
reusable libraries. It supports the new plan/build language alongside the
existing YAML manifest and HCL v2 workflows during cutover.

The short version:

- `mogent.lib` describes reusable typed content.
- `mogent init` creates an authored `mogent.build`; `mogent init --plan` creates
  selection intent in `mogent.plan`. Neither selects content or installs outputs.
- `mogent source show` inspects sources; add includes before `mogent build`.
- The following YAML and HCL workflows remain available through explicit flags:
- Libraries hold reusable Markdown modules.
- `agents.yaml` chooses, orders, and renames those modules for one repo.
- `mogent build` renders the manifest to `AGENTS.md`.
- `mogent build --dry-run` validates the render without writing outputs or state.
- HCL v2 uses `mogent.hcl` and `plan` / `apply` / `update`, including mixed
  copied files, complete skill bundles, and rendered Markdown.

## Feature Coverage And Release Planning

Status is specific to the format shown; a v1 capability is not automatically
available in v2. Dates below are recorded delivery checkpoints, not a claim about
the earliest introducing commit. A dash means no delivery date is asserted.

| Capability | Delivered scope | Recorded checkpoint | Redesign / release treatment |
|---|---|---|---|
| Authored headings, nesting, renaming, ordered multi-source composition | YAML v1 `doc` | — | Preserve consumer composition intent |
| Explicit content template variables | YAML v1 `vars`; not HCL v2 | 2026-08-18 | Compatibility decision required |
| Source discovery, metadata/TLDR browsing, coverage | v1; v2 has `source list` / `show`, not full v1 parity | v2 browsing: 2026-10-05 | Preserve; unify views and selection explanations |
| Copy-on-write localization and conservative single-section drift import | YAML v1 workspace API / CLI | — | Preserve workflow; v2 parity is not delivered |
| Library identity, hierarchy, metadata, and relationship validation | v1 sidecar; v2 sidecar has its own schema | — | Preserve; relationship validation is not a general solver |
| Sparse selections, curated choices/defaults and reasoned foundation exclusions | HCL v2 authored guidance | 2026-10-05 | Preserve curation; simplify choice/exception interaction |
| Tag-based content selection | v1 output selectors; v2 additive authored-leaf queries | — | Preserve; define tag/default/exclusion precedence |
| Named init templates | v1 and v2; editable starting configurations | v2 `qmr-core`: 2026-10-05 | Distinguish templates from live includable presets |
| Pinned Git sources, update previews | v1 lock; v2 inline commit pins | v2: 2026-10-05 | Reuse Git; frozen local inputs still need a contract |
| Plan/apply, output drift protection, transactional writes | HCL v2 | 2026-10-05 | Reuse engine and state lifecycle |
| Single-source raw directory outputs and source exclusions | HCL v2 `tree` | — | Consolidation with richer composition remains open |
| Physical inventory and intact file/bundle copies | Typed `content` API and HCL v2 `dir-tree` | 2026-10-09 | Preserve bytes, frontmatter, modes, bundle boundaries |
| Parsed Markdown heading views and explicit rendering | Typed API / HCL v2 | 2026-10-09 | Preserve; extraction is not editing an installed skill |
| Mixed ordered contributions, subtree replacement, rendered append | Typed API / HCL v2 | 2026-10-09 | Reuse; simplify recommended authoring vocabulary |
| Member/rendered-document update diffs and operation provenance | Typed HCL v2 | 2026-10-09 | Preserve reasons, omissions, replacements, origin context |
| Byte, permission, and empty-directory drift protection | Typed HCL v2 | 2026-10-09 | Preserve |
| Interactive manifest navigator | YAML v1 TUI | — | Typed-result/decision viewer is proposed, not delivered |
| Small-language syntax API | Public `syntax` parser/formatter, typed `mogent.lib` catalog, and API tests | Current redesign slice | Source ranges, comments, literals, aliases, logical/physical nodes |
| Directly authored assembly file and intent-to-build compilation | Local/Git sequencing, checksummed builds, assembly, localize CLI | Current redesign slice | Advanced curation, full rollback/state, cutover still pending |
| Includable recipes, curated cross-branch presets, build-to-plan suggestions | Not implemented as proposed | — | Distinguish from existing branch selection/init templates |
| Lower-layer deletion, editable heading derivatives, semantic JSON/TOML merges | Not implemented | — | Separate design work; whole-file copying already works |
| DotAgents workspace layout/format compatibility | Opaque files can be copied; protocol conformance not established | — | Candidate 1.0 scope; compare native harness support first |
| Automatic global installation and vendor-format conversion | Not implemented | — | Deferred; workspace-first direction |
| Guided multi-repository deployment | Campaign documented; no batch transaction | 2026-10-09 plan | Agree pilot repositories and useful modules with owner |

The [NOW panel](docs/convo/head.md) owns the current conversation resume point.
The [configuration redesign conversation stack](TODO/2026-10-14-212422-conversation-stack.md)
preserves the historical discussion and rationale. The consolidated
[first-slice spec](docs/proposals/MOGENT-FIRST-SLICE-SPEC.md) and
[acceptance cases](docs/proposals/MOGENT-FIRST-SLICE-ACCEPTANCE.md) define the approved
initial implementation, with a [staged implementation plan](TODO/TODO-config-first-slice.md).
The syntax API is now implemented/tested; catalog/resolver and CLI integration
are now available for the initial local/Git workflow. The legacy commands remain
available while advanced curation, output transaction/state parity, and cutover
work continue.
They do not replace the executable
[HCL v2 specification](docs/proposals/MOGENT-HCL-V2-SPEC.md) or
[typed-content contract](docs/TYPED-CONTENT.md).

### Current Plan/Build Workflow

`mogent -h` and `mogent --help` show the same concise overview, grouped by
**project initialization, browsing, selection, preview, and building**.
`mogent help --all` expands that workflow with common flags under each command.
Each group links to focused help with its own advanced section.

Help uses bold headings, command names, and flags in terminals. Piped or
redirected output stays plain; `NO_COLOR` and `TERM=dumb` disable styling.

```sh
mogent help --all
mogent help source
mogent help select
mogent build --help
```

The new format selects modules through `mogent add <source:path>` (also available
as `mogent select` or `mogent module add`), which appends an `include` declaration
to the authored build or plan. `mogent add --manifest agents.yaml` retains the
legacy YAML placement workflow. Legacy commands are grouped under `mogent help
legacy`, utilities under `mogent help tools`, and all registered command
signatures remain available with `mogent help commands`.

For an opt-in tree such as QMR skills, select its direct children explicitly:

```sh
mogent add qmr:skills --all
```

Without `--all`, a bare opt-in tree selects no members and therefore installs
nothing.

Exclude a descendant from a selected parent with `mogent except <source:path>`
(`exclude` is an alias). It edits that parent include's `exclude` list; it does
not remove a separate direct include or automatically optimize overlap.

Run inside a chosen consumer workspace:

```sh
mogent init
mogent source add qmr ../qmr-agents-library
mogent source list
mogent source show
mogent show qmr:intro/assistant-agent --content
# Edit mogent.build to include the content appropriate for this project.
mogent show --build --file mogent.build
mogent build --file mogent.build
```

Initialization refuses existing files. Sources are declarations, not automatic
selections. For intent-first authoring, start with `mogent init --plan`; edit
`mogent.plan`, then use `mogent sequence` to emit a generated build, or
`mogent build` to sequence and assemble it. Use `--file path` for a custom entry
path. Dry-run initialization writes nothing and does not fetch sources.

`source show` renders a tree with aligned path/type columns; `--tree` is
accepted explicitly, `--flat` prints full paths, and `--chars ascii|unicode`
overrides the existing user display preference. Empty source lists explain how
to add a local library or Git repository.

Use `mogent show source:path --content` to preview a source node without
selecting it. `--lines N` limits text output (`0` is full content) and
`--headings` lists heading addresses for physical Markdown files.

Remember a library once for discovery across projects:

```sh
mogent source register qmr ~/dev/cdint/qmr-agents-library
mogent source list --known
# Inside another initialized project:
mogent source add qmr
```

With [Gum](https://github.com/charmbracelet/gum) installed, interactive `init`
offers known libraries, a local directory, a Git repository, or skipping. Gum is
included in the Nix development shell. `--no-interactive`, explicit `--source`,
dry runs, and nonterminal execution bypass the picker. Without Gum, flags and
manual source addition remain available. `source ls` aliases `source list`.

The discovery registry uses YAML beside the existing user preferences; it is
never a build input. See [Source discovery and browsing](docs/SOURCE-DISCOVERY.md)
for storage, project authority, and compatibility details.

### Legacy HCL V2 Workflow

The current v2 loop is **initialize -> inspect sources -> edit configuration ->
plan -> revise -> apply -> review updates**. For example, with the sibling QMR
library checked out, run these inside a chosen consumer workspace:

```sh
mogent init --config mogent.hcl --template qmr-core --source qmr=../qmr-agents-library
mogent source list qmr --tldr --tags
mogent source show qmr:constraints/security
# Edit mogent.hcl to change sources, selections, or destinations.
mogent plan
mogent apply
```

For a Git-backed source, `mogent update qmr` previews a pin update; a local
source is read as checked out. V2
selection changes currently use the editor; v1 `source add` / `add` commands are
not a promised v2 editing interface. For a self-contained typed-content inspection
and plan/apply example, see [Typed Content Composition](docs/TYPED-CONTENT.md).

The new language accepts `mogent.lib`, `mogent.plan`, and `mogent.build`, with
`sequence`, `show`, `build --file`, and `edit include --localize` available for
the initial workflow. Richer inspection/help and cutover remain staged work;
KDL is historical exploration, not an accepted input format. See the
[first-slice spec](docs/proposals/MOGENT-FIRST-SLICE-SPEC.md) and
[implementation sequence](TODO/TODO-config-first-slice.md) for scope and gaps.

## License

No license is granted at this time. Do not copy, redistribute, or adopt Mogent
outside this repository without explicit permission from its owner.

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

## Handle Minting

`mogent mint` generates an unused proquint handle after scanning the repository
corpus. Use `-n` to verify a deterministic candidate without reserving it:

```sh
mogent mint -r . -w 1 -n
```

Use `-w 2` when a two-word handle is required. The command prints the handle; it
does not create a file or modify the repository.

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
mogent init --manifest agents.yaml --template minimal --source shared=/path/to/your/library --dry-run
mogent init --manifest agents.yaml --template minimal --source shared=/path/to/your/library
```

Mogent does not ship a selectable policy library. Use an explicit local library
path, such as the sibling QMR personal library, or a source owned by the target
organization.

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

The checked-in metadata fixture uses this atomic style. Each reusable module
usually lives in one file with one top-level heading and frontmatter metadata.
For example, `testing-ground/metadata-library/library/lang/go/testing.md`
resolves as `shared:lang/go/testing`.

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

## HCL v2 Configuration

`mogent.hcl` is the v2 consumer configuration. It is a separate format from
`agents.yaml`; a repository uses one or the other. The design of record is
[`docs/proposals/MOGENT-HCL-V2-SPEC.md`](docs/proposals/MOGENT-HCL-V2-SPEC.md).

Start from a library and look at what it offers:

```sh
mogent init --config mogent.hcl --template qmr-core --source qmr=../qmr-agents-library
mogent source list qmr --tldr
mogent source show qmr:constraints/security
```

Show what would change, then write:

```sh
mogent plan    # unified diff per output path, diagnostics, no writes
mogent apply   # writes outputs, state, and any newly resolved commit pin
```

A library is `local` (read as checked out) or `git` (pinned by `commit`,
which `apply` writes into `mogent.hcl`). Move a pin forward with review:

```sh
mogent update            # show new commit and the output diff it causes
mogent update qmr --accept
```

The library sidecar `library.mogent.hcl` declares how each branch should
`curate` its children:
`foundation` (in unless dropped with a reason), `choose` (decide per child
or `accept_defaults`), `optional` (library default applies), `opt_in`
(excluded unless named). `plan` reports every accepted default, forced
exclusion, unpinned source, and hand-edited output with a stable `MOGENT`
code. `tools/v2-stories` runs the spec's acceptance scenarios end to end.

Typed packages use `kind = "dir-tree"` with explicit copy/render contributions.
Named physical roots can preserve complete skill bundles and configuration files,
while optional heading views supply generated Markdown in the same managed output.
Collisions require explicit replacement; plan reports operations and provenance.
See [Typed Content Composition](docs/TYPED-CONTENT.md) for the grammar, API, and
the runnable `testing-ground/typed-content-demo/` workflow.

## Manifest

New manifests use output-scoped selectors, so every target is independent:

```yaml
roots:
  qmr: /path/to/qmr-library
sources:
  qmr_agents:
    path: {root: qmr, subdir: agents}
  qmr_skills:
    path: {root: qmr, subdir: skills}
outputs:
  - path: AGENTS.md
    include:
      - all: qmr_agents
    exclude:
      - tags:
          any: [lang/go, os/windows]
  - path: .agents/skills/
    include:
      - all: qmr_skills
```

`all` selects a source root, `source` selects `alias:heading/path`, and
`tags.any` matches indexed library metadata. Markdown outputs preserve selected
source headings; directory outputs copy selected contents. Tag filtering fails
for raw skill trees whose metadata cannot be indexed. Roots and subdirectories
are literal paths only: Mogent does not expand shell or environment variables.
Local rooted sources may safely use `subdir`.

Use block YAML for readable manifests. YAML flow syntax such as
`include: [{all: qmr_agents}]` remains valid input, but is not the documented
authoring style.

`output` plus `doc` remains fully supported legacy input. Migrate by replacing
the document outline with one or more `outputs` entries and moving each output's
source selection into `include`/`exclude`; retain the old form when its custom
heading placement is needed.

Legacy `agents.yaml` is an authored document outline:

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

### Additional Outputs

In the canonical model every `outputs:` entry is a peer target. A `.md` path is
Markdown; a path ending in `/` (including `.d/`) is a directory. Bare paths are
rejected so output kind is never guessed. The older `output` plus `doc` form
continues to produce its primary document and may retain legacy additions.

```yaml
outputs:
  - path: docs/skills-index.md
    include:
      - all: qmr
  - path: .agents/skills/
    include:
      - all: qmr_skills
```

The directory example copies the *contents* of the selected `qmr_skills`
source directory recursively, including `SKILL.md` and non-Markdown assets. It
does not copy the source root directory itself. Legacy directory outputs may
continue to use `from`; canonical Markdown outputs use `include`.
Mogent records every copied file, rejects symlinked sources or targets, and
refuses untracked or edited output content unless `--force` is supplied.

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

## Library Sidecars

An optional `library.mogent.yaml` at a library root records library identity,
the exhaustive selected directory/heading tree, display titles, ordering, and
validated relationships. A central repository sidecar can explicitly select
publishable Markdown and raw-copy roots while leaving unrelated material out:

```yaml
content:
  markdown_roots: [agents]
  directory_roots: [skills]
```

Markdown source paths remain repository-root-relative, such as
`agents/intro`. The configured Markdown root is virtual, so selecting `all` of
the repository renders its children without an `Agents` wrapper. `archive/`,
`reference/`, and `proposals/` outside `content` are not parsed or indexed.
`directory_roots` are declarative raw-copy trees and their Markdown is not
parsed. A sidecar without `content` retains full-library behavior. When present,
its complete valid tree controls rendered source ordering and headings without
changing source references or Markdown bodies.
Create or inspect a deterministic candidate and validate it with:

```sh
mogent lib init /path/to/library
mogent lib scan /path/to/library --dry-run
mogent lib check /path/to/library
mogent lib check --source shared --manifest agents.yaml
```

The manifest source alias (for example `qmr`) is a local manifest name, while a
library root is the repository directory containing `library.mogent.yaml`.
Aliases may still use `subdir` to address a narrow source root, but a central
root sidecar always stores `tree.source` paths relative to its library root.

`init` never overwrites an existing sidecar. Rendered Markdown strips HTML
comments by default; use `mogent build --preserve-html-comments` for an audit
render. Raw directory outputs always copy source files unchanged.

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
mogent source add personal /path/to/personal-library --dry-run
mogent source add personal /path/to/personal-library
```

Declare, pin, and later review a URL source:

```sh
mogent source add qmr https://github.com/Qu1ncyRy4n/qmr-agents-library.git \
  --subdir agents --dry-run
mogent source add qmr https://github.com/Qu1ncyRy4n/qmr-agents-library.git \
  --subdir agents
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
Use `/` between heading levels. Escape a literal slash in a heading as `\/` and
a literal backslash as `\\`; for example,
`--under 'Constraints \/ Safety'` addresses one heading containing a slash.
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
mogent drift --diff
mogent drift --import Instructions/Testing
mogent drift --reject --force
```

Localization writes ordinary Markdown under `.mogent/library`, records its
origin in `.mogent/provenance.yaml`, and changes the selected manifest reference
to `local:`. Drift import refuses edits outside the selected section.
`drift --diff` compares the expected Markdown render with the local generated
Markdown without writing files. Source-version drift and pin comparison are
future work; they are not `drift` behavior.

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
- `drift --diff` is read-only and prints no patch when generated Markdown matches.
- Shared source libraries are read-only inputs during normal build/add flows.

## Testing Ground

For a feature-by-feature exercise sequence and issue template, see
[`docs/DOGFOOD.md`](docs/DOGFOOD.md). For only the newest selectable-directory
and unified-tree behavior, use
[`docs/dogfood/history/DOGFOOD-SESSION-3.md`](docs/dogfood/history/DOGFOOD-SESSION-3.md).

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
- tool-specific Markdown renders or template inheritance,
- split/import from a full `AGENTS.md` into an atomic library.

Future multi-output support should allow both exact mirrors and tool-specific
renders. For example, `CLAUDE.md` could include extra visible manifest entries
for stricter Claude behavior, while `GEMINI.md` might simply symlink to
`AGENTS.md`.

## Design Docs

- [docs/DESIGN.md](docs/DESIGN.md) is the design of record.
- [docs/proposals/MOGENT-HCL-V2-SPEC.md](docs/proposals/MOGENT-HCL-V2-SPEC.md)
  is the v2 configuration specification; `TODO/TODO-pupab-hcl-v2-rollout.md`
  tracks its implementation.
- [docs/user-stories.md](docs/user-stories.md) holds the product user stories.
- [docs/DOGFOOD.md](docs/DOGFOOD.md) stages current features for dogfooding and issue reporting.
- [docs/demo/MOGENT-BOSS-DEMO.md](docs/demo/MOGENT-BOSS-DEMO.md) is a
  slide-separated live-demo runbook using a disposable todo-app clone.
- [docs/dogfood/history/DOGFOOD-SESSION-3.md](docs/dogfood/history/DOGFOOD-SESSION-3.md) exercises only the new directory-tree and display work.
- [docs/HANDOFF.md](docs/HANDOFF.md) is the stable handoff navigation entry;
  its linked dated pickup is the detailed current resume point.
- [docs/AUTHORING-PLAN.md](docs/AUTHORING-PLAN.md) stages source declaration, move/reorder, and guided module creation.
- [docs/ORIGIN-RECONCILIATION-PLAN.md](docs/ORIGIN-RECONCILIATION-PLAN.md) stages preservation, inheritance, and proposal workflows.
- [TODO/TODO-jusuk-mogent-agent-modules.md](TODO/TODO-jusuk-mogent-agent-modules.md)
  tracks implementation and open decisions.
- [docs/thought-experiments/TE-kavam-metadata-tags-and-library-shape.md](docs/thought-experiments/TE-kavam-metadata-tags-and-library-shape.md)
  covers metadata, tags, conflict warnings, and atomic library shape.
