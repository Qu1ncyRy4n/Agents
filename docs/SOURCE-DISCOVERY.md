# Source Discovery And Browsing

## Project Sources

`mogent init` creates authored `mogent.build`; `--plan` creates `mogent.plan`.
Sources declared there are available inputs, not selections:

```text
source "qmr" {
  local = "../../cdint/qmr-agents-library"
}
```

Use the CLI instead of editing the declaration manually:

```sh
mogent source add qmr ../../cdint/qmr-agents-library
mogent source add team https://example.org/team/library.git
mogent source add preview ../library --dry-run
mogent source list
mogent source ls
```

`add` appends a validated declaration while preserving existing bytes/comments
and file permissions. Duplicate aliases and edits to generated builds are
rejected; edit the owning plan instead. Listing declarations never opens local
libraries or fetches Git sources, so unavailable sources can still be listed.

Paths are relative to the project entry file. Absolute local paths (including
`~/...`) are converted to relative paths when written. Git declarations support
HTTP(S) and `file://`, without embedded credentials, query strings, or fragments.

`show` and new source commands discover `mogent.plan` or `mogent.build`. A plan
paired with its generated build selects the plan. Competing plan/authored-build
files require `--file`.

## Catalog View

```sh
mogent source tree
mogent source tree --chars unicode
mogent source tree --flat
mogent source tree qmr:intro
mogent source list --tree
```

The default source view is a hierarchy with aligned PATH and TYPE columns.
Library node types remain `dir`, `file`, `md-file`, and `md-section`; presentation
does not change their semantics. Tree structure follows catalog parent/child
relationships and declaration order. When every payload in a top-level logical
branch lives below one declared physical directory, the tree nests that branch
under the directory; mixed-source branches remain logical roots rather than
being guessed into a directory.
Flat output uses complete paths sorted lexically. A node address limits the
view to that subtree. Presets and recipes are not typed content nodes and are
not included in this initial tree view.

The connectors and alignment reuse existing conventions from
[DESIGN.md](DESIGN.md): `display.chars: ascii|unicode`, `display.align`, and CLI
`--chars` override. ASCII remains the portable default. The catalog adapter does
not yet implement coverage, TLDR, width, or search flags.

An empty project reports no sources and gives local/Git addition and known-library
discovery commands. Catalog inspection can load/fetch libraries; declaration
listing does not. No selection or installation occurs during source addition.

Preview a module without selecting it:

```sh
mogent source content qmr:intro/assistant-agent
mogent source content qmr:intro --lines 12
mogent source content qmr:skills/review
mogent source content qmr:assets/readme.md --headings
```

Section previews use the same rendered-section behavior as builds, including
declared project variable bindings. Physical text files preserve their payload
bytes; binary files report their byte count. `--lines 0` is the default full
text view. Heading addresses are only meaningful for physical Markdown files.
All previews are read-only.

## Module Selection

```sh
mogent add qmr:intro/author-ownership
mogent module add qmr:intro/author-ownership
mogent add qmr:intro/author-ownership --dry-run
```

These native aliases validate the source and catalog node, then append one
`include` declaration to an authored build or plan. They preserve existing text
and do not build/install outputs. Selecting the same direct include twice is
rejected. Generated builds remain read-only: select content in their owning plan.

Opt-in directories, such as QMR's skills tree, need an explicit child selection:

```sh
mogent add qmr:skills --all
```

`--all` creates or upgrades a bare direct include to a list of its current direct
children. It is not implicit because skills are opt-in content. The command
refuses to silently replace an existing explicit selection.

When a bare opt-in include resolves to nothing, `show --issues` now warns with
the affected reference and recommends explicit child selection or `--all`.

## Output State

Every successful native build records its managed-output baseline in
`.mogent/state.json`, whether the build file is authored or generated. The state
contains SHA-256 hashes for `AGENTS.md` and every file beneath `.agents/skills/`.
It lets later status/drift work distinguish Mogent's last output from a direct
edit. It is an output baseline, not a replacement for the authored build file or
a source-selection lock; generated plan builds additionally retain their build
and resolved-input hashes.

Inspect drift without writing:

```sh
mogent drift
mogent drift --diff
mogent drift --all --diff
mogent status
mogent build --check
```

Native `drift` reports whether `AGENTS.md` matches its last output baseline and
whether it matches the current configuration/source composition. `--all` adds
the managed `.agents/skills` tree, including missing, modified, and unexpected
files; `status` is the concise all-output form. `--diff` prints text output
diffs and changed skill paths. `build --check` verifies all managed outputs,
including skill bytes, modes, and directory layout. Native drift/status are
read-only. Native builds also record deterministic library-root and per-node fingerprints. Drift
compares those saved baselines with each connected library, descending a changed
tree until it can name changed, added, or removed module paths. State written
before this feature has no upstream baseline; run `mogent build` to record one.

Exclude one descendant from a direct selected ancestor:

```sh
mogent add qmr:workflow
mogent except qmr:workflow/decision-first
# `mogent exclude ...` is an alias.
```

`except` finds the narrowest direct selected ancestor, then adds the descendant's
relative path to that include's `exclude` list. It previews with `--dry-run` and
does not build outputs. For a curated choose branch, it also records
`accept_defaults = true`: the resulting meaning is "the defaults except this
item." It deliberately does not remove another direct include. Use explicit
review commands instead:

```sh
mogent rm qmr:workflow/decision-first --dry-run --diff # preview: writes nothing
mogent rm qmr:workflow/decision-first        # remove one exact direct include
mogent module rm qmr:workflow/decision-first # equivalent explicit form
mogent optimize                              # default preview of overlap cleanup
mogent optimize --apply                      # remove redundant direct descendants
```

`rm` never removes an ancestor selection or guesses among repeated exact
includes. Combine `--dry-run --diff` to compose the before/after managed
outputs without writing; bare `--diff` applies and reports its patch.
`optimize` (with `optimise` retained as an alias) shows that same combined
preview by default. Its `--apply` removes redundant direct descendants, but
never reorders or silently changes selections.

## Tag Filtering And Discovery

Libraries may tag nodes. An authored include can select direct children by tag:

```text
include "qmr:skills" {
  tags_all = ["language/go"]
  tags_any = ["review", "security"]
  exclude = ["review"]
}
```

`tags_all` requires every named tag and `tags_any` requires one or more. A
qualified tag prefix matches its descendants, so `language/go` also matches
`language/go/tests`; an unqualified term matches a complete tag component, so
`workflow` matches `topic/workflow`. Neither form is fuzzy text matching.
Explicit child selections are unioned with tag matches. The normal scoped
`exclude` filter runs afterward; it excludes a source-relative path, not a tag.
`tags_none` retains nodes carrying none of its tags.

Use `mogent source search [source:path] --tag-all <tag> --tag-any <tag>
--tag-none <tag>` to inspect matching catalog nodes. `--add-commands` prints
runnable explicit `mogent module add` commands for matching leaf modules only;
it never suggests a broad parent that would include unrelated descendants. For a
`source:path` scope, `--add-config` instead prints the pasteable dynamic
`include` block that expresses the query. The same repeated `--tag-all`,
`--tag-any`, and `--tag-none` flags are accepted by `mogent module add` for a
new dynamic tag-query include. `module add --explicit-commands` previews the
same stable leaf-by-leaf alternative without writing. Tag predicates cannot
combine with `--all`, because an all-child selection would make the tag query
redundant.

`source search` currently searches catalog tags only. Title/path and rendered
content search are separate future modes; they must not silently change the
selection-oriented tag-query meaning. Interactive fuzzy matching is also
deferred until those mode boundaries are specified.

`module optimize` belongs to selection editing, not `build`: builds never
rewrite authored choices. Its current overlap removal is semantics-preserving.
Replacing explicit modules with a tag query is not automatically equivalent,
because a future upstream tagged module would then become selected. If offered,
that conversion must be a separately previewed, explicitly accepted change—not
an ordinary optimize rewrite.

Reorder authored configuration without changing selection semantics:

```sh
mogent reorder --dry-run
mogent reorder
```

This explicit formatter puts bindings first, then source declarations, followed
by includes in their existing order. Each include's options stay attached. It
does not deduplicate content, remove modules, or alter rendering behavior.

## User-Level Address Book

```sh
mogent source register qmr ~/dev/cdint/qmr-agents-library
mogent source register team https://example.org/team/library.git
mogent source list --known
# Inside a new initialized project, materialize a remembered declaration:
mogent source add qmr
```

The registry is `$XDG_CONFIG_HOME/mogent/registry.yaml`, or `mogent/registry.yaml`
under the platform user config directory when XDG_CONFIG_HOME is unset. It is
separate from `config.yaml` display preferences:

```yaml
format: 1
sources:
  - alias: qmr
    name: QMR Agent Library
    location: /absolute/path/to/qmr-agents-library
```

YAML was selected for this small user-local record because existing preferences
already use it, the dependency is present, and this file has no composition
semantics. TOML would add a parser; extending the native plan/build language with
a fourth registry role would enlarge its contract for a simple address book.
Project/library content continues to use the native language.

Registration is explicit. Local locations are made absolute and inspected for
their library name; Git URLs are stored without fetching. Repeating the same
registration is idempotent. Reusing a name for another location fails rather
than silently replacing it. The current CLI has no registry remove/replace
command; edit this user-local file to manage stale entries.

Choosing an entry copies its location into the project's declaration. The
project uses no implicit registry lookup afterward: moving, deleting, or corrupting
the registry cannot alter builds. Known entries are not active instructions,
installed skill directories, or a global runtime configuration.

## Interactive Init With Gum

In a terminal with Gum on PATH, plain `mogent init` offers registered libraries
(name and location), Local directory, Git repository, and Skip for now. New
locations prompt for a source name and path/URL. They are project declarations;
use `source register` separately to remember them for future projects.

Gum only handles selection/input; the Go APIs own validation and file writes.
The target is checked before prompting, and cancellation leaves no project file.
Gum's UI uses stderr; its selected value is captured separately.

Gum is optional and included in `nix develop`. It can also be installed with
`brew install gum`. Without Gum, init creates the minimal file and explains
source addition. No picker runs with `--no-interactive`, `--dry-run`, explicit
`--source` or nonterminal input/output.

## Verification

- API tests: source preservation, permissions, duplicate/generated-build refusal,
  no-fetch listing, registry validation and XDG isolation.
- CLI tests: native tree/flat/subtree output, source add/list/ls/register, empty
  help, malformed-registry-independent browsing/builds, picker known/new/skip,
  cancellation, existing-target preflight, noninteractive bypass, all-output
  status/drift, direct-include removal/output preview, and overlap diagnostics.
- Real Gum v0.17.0 pseudo-terminal smoke: selecting a registered library creates
  a project declaration; Ctrl-C exits without creating the project file.
- Full repository validation uses `tools/check`.
