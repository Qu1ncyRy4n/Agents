# Mogent Plan And Build: First-Slice Specification

Status: **initial implementation approved; syntax API implemented/tested, runtime workflow pending**.

This consolidates the configuration discussion into a proposed executable
contract. Accepted product decisions are distinguished from proposed engineering
details below. This document does not retroactively change YAML/HCL behavior.

- Discussion history: [conversation stack](../../TODO/2026-10-14-212422-conversation-stack.md).
- Tests to implement: [acceptance cases](MOGENT-FIRST-SLICE-ACCEPTANCE.md).
- Delivery sequence: [implementation plan](../../TODO/TODO-config-first-slice.md).
- Historical examples: [walkthrough](MOGENT-LANGUAGE-WALKTHROUGH.md).

## 1. Scope And Review Boundary

### Accepted Direction

One small custom language, ordered root actions, explicit values, bracketed lists,
immutable bindings, interpolation, multiline text, private recipe scopes, and one
`include` family. Content kind determines normal inclusion behavior; no required
copy/render/append switches. `replace` is an alias, not another executor operation.

Plans express selections/curation; builds express assembly. Direct build authoring
is independently useful. Every plan passes through the same build representation
and executor. `build` sequences automatically when safe; `sequence` stops after
writing assembly instructions. Authored or hand-edited builds are protected.

Preserve tags, curation, defaults, reasons, logical guidance structure, physical
bundle integrity, local editable sources, provenance, deterministic order, output
ownership, drift protection, and transactions. Root `AGENTS.md`, other initial
outputs under `.agents/`, real files. Git commits and local-node hashes identify
inputs; a hash is not a content store or a stable node identifier.

### Proposed Engineering Details Requiring Review

| Detail | Proposed first-slice choice |
|---|---|
| Files/schema | `mogent.lib`, `mogent.plan`, `mogent.build`; format 1 of this new language, not legacy format 1 |
| Library language (owner chose B) | Design and implement the new library format alongside plan/build; do not make an ongoing legacy-format adapter the first-slice dependency |
| Library schema details | Typed `node` plus dir/file/md-file/md-section aliases, optional role/content, child-local defaults and reusable presets as specified in section 4 |
| Header integrity | Portable generated-file checksum, plus input/decision records; local state is not the only regeneration baseline |
| Binding/string grammar | Strings as bindings, integers for format metadata, double/triple-double strings and literal triple-single strings |
| Resolved metadata | Reserved `resolved` block with input, decision, and ownership records; per-include `values` records declared guidance variables without ambient re-evaluation |
| Provenance storage | `.mogent/provenance.json` for new localization records; migrate rather than overwrite legacy YAML records |
| Direct local build records | `source update local --accept` deliberately refreshes selected-local-node hashes in an authored build; ordinary build never rewrites it |
| Replacement typo handling | Replacing a destination with no prior contribution is an error, not an implicit ordinary include |
| Initial interface | API + text result/build/source views, then interactive viewer; no full TUI rewrite before resolver correctness |

The owner accepted the final combined review, chose `mogent.lib`, and requested
API implementation -> focused API tests -> CLI -> CLI behavior tests. This
authorizes the first implementation slice; finer behavior changes still require
an explicit documented decision rather than silent widening of the language.
No source-library content or existing consumer configuration is deleted implicitly.
The owner explicitly selected new library syntax in the initial redesign. Retain
old config files as migration/reference material that can be interpreted; maintaining
YAML as an active format is not a goal. Any bounded conversion reader is migration
tooling, not an implicit runtime fallback or a second maintained authoring surface.

## 2. Files And Authority

| File | Job | Version-control expectation |
|---|---|---|
| `mogent.plan` | Optional authored intent | Commit when used |
| `mogent.build` | Authored or generated assembly instructions and input records | Commit for reproducible distribution; generated does not mean disposable |
| `.mogent/local/` | Editable localized source subtree | Commit intentional project-local content |
| `.mogent/provenance.json` | Origin/baseline references for localized content | Commit alongside localized sources |
| `.mogent/state.json` | Last installed result, layout, hashes, management history | Machine-local; not an input required to determine expected output |
| `.mogent/sources/` | Retrieved source cache | Machine-local; missing cache can be refetched at an exact pin |
| `AGENTS.md`, selected `.agents/` outputs | Final result | Commit or ignore according to consumer practice; check works either way |

The plan is authoritative for its derived build. A direct authored build is its
own authority. Two independent files never resolve by silently preferring one.
`author` is descriptive; `origin` determines regeneration policy.

```hcl
build {
  format         = 1
  origin         = "generated"
  generated_from = "mogent.plan"
  generated_by   = "mogent"
  plan_sha256    = "<SHA-256 of exact plan bytes>"
  content_sha256 = "<generated-file checksum>"
}
```

A plan uses `plan { format = 1; }`. A direct build uses
`build { format = 1; origin = "authored"; }`; optional `author` is `human`,
`agent`, or `unspecified`. No timestamps, machine paths, or random IDs are emitted
into deterministic generated files. Source content may contain user-supplied dates.

**Proposed checksum rule:** SHA-256 of the complete UTF-8 generated file with its
single `content_sha256` value replaced by 64 ASCII zeroes. This detects comment,
format, or semantic edits across checkouts without requiring state.json. It is
integrity/change detection, not authentication. Missing/malformed/mismatched
checksum means regeneration is blocked. A user can deliberately maintain a direct
build by changing its origin to authored; that is not permission to overwrite it.

## 3. Minimal Grammar

Code fences use HCL highlighting only; this is not HCL and does not accept its
expression language. A parser may reuse infrastructure only if it implements
this contract faithfully.

```text
file        := header separator { declaration separator }
declaration := binding | source | include | replace | recipe | resolved
binding     := "let" identifier "=" string
source      := "source" string "{" fields "}"
include     := "include" [string [string_list]] ["{" include_body "}"]
replace     := "replace" string [string_list] ["{" include_body "}"]
recipe      := "include-build" string
field       := identifier "=" literal
literal     := string | boolean | integer | string_list
string_list := "[" [string {"," string} [","]] "]"
separator   := newline | ";" | EOF

library_file := library_header separator { library_item separator }
library_item := node | typed_node | preset | group
node         := "node" string "{" node_body "}"
typed_node   := ("dir" | "file" | "md-file" | "md-section") string "{" node_body "}"
preset       := "preset" string "{" { (include | replace) separator } "}"
group        := "group" string "{" fields "}"
```

This is structural shorthand, not a complete lexer grammar. Braces can close the
last field without an extra separator. Multiline lists allow whitespace/comments;
commas between elements are required. Semicolons separate fields/actions on the
same line; no trailing semicolon is required for a single field/action before
`}`, newline, or EOF. Commas separate array values only, not statements/options.
`include` without a subject requires inline `text`; empty actions are errors.
`include_body` holds fields plus the named `exclude` exception blocks below and
an optional `values { name = "literal"; ... }` block for declared guidance inputs.
`source`/header/resolved bodies have closed schemas. Nested arbitrary executable
blocks and object-valued `= { ... }` assignments are not part of this first slice.

- UTF-8, case-sensitive identifiers and paths. Identifiers use ASCII letters or
  underscore initially, then letters/digits/underscore/hyphen.
- Comments: `//` to end of line. Comment tokens inside strings are text.
- Strings are quoted; paths/tags cannot be unquoted expressions.
- Boolean literals are `true` and `false`; only decimal integers where a schema
  calls for them. No arithmetic, null coercion, functions, interpolation expressions,
  loops, conditionals, environment expansion, or implicit merges.
- Duplicate fields, source aliases, or bindings in one scope and unknown fields
  are errors. Repeated include actions are allowed and checked for collisions.
- Bindings and source declarations must precede their use in that file's scope.
- Every node/field retains a source range and comments for diagnostics/editing.

### Bindings, Interpolation, Multiline Text

`let repo_name = "Example"` declares an immutable string before use. `${repo_name}`
in normal strings inserts that value; the inserted value is not evaluated again.
Undefined names are errors; unused bindings/aliases are warnings. A literal
`\${name}` escapes the interpolation opener in a normal string.

Normal strings support `\"`, `\\`, `\n`, `\r`, `\t`, and escaped interpolation;
unsupported escapes are errors. `"""` strings must open immediately before a
newline; only indentation may precede the closing delimiter on its line (normal
separators or a closing brace may follow it). Normalize line endings
to LF, remove the framing newlines, and strip the closing delimiter's indentation
from each nonblank line (insufficient indentation is an error), then process
escapes/interpolation. Blank lines remain blank. Literal `'''` multiline strings
use the same framing/indentation rules but no escapes or interpolation. Embedded
unescaped matching triple delimiters terminate a string.

Inline text and rendered document fragments are normalized to a single final
newline in the generated document, with one blank line between contributions.
Physical copies, including copied Markdown, are byte-preserving.

## 4. Sources, Catalogs, And Recipes

```hcl
source "team" {
  git = "https://example.org/team/library.git"
  ref = "main"
}

source "local" {
  local = ".mogent/local"
}
```

A source has exactly one of `local` or `git`, with optional safe `subdir`.
Git accepts the existing supported HTTP(S)/file URLs and full commits. Local
paths resolve relative to the declaring file; write destinations resolve relative
to the workspace containing the entry file. Source roots may be sibling libraries
but selected content may not escape its declared root. No process-CWD dependence.

The normalized catalog exposes stable qualified paths, physical/logical kinds,
library order, tags/TLDRs/titles, curation/defaults, relationships, bundle boundaries,
and suggested placement. Metadata is inspected independently of payload rewriting.

The initial implementation loads `mogent.lib`, using the same small syntax as
plan/build. Existing HCL/YAML declarations remain migration references, not runtime
fallback formats. The following is the final-review library schema proposal.

```hcl
library {
  format = 1
  name   = "Team Library"
}

md-section "guidance" {
  curate = "foundation"

  md-section "workflow" {
    source = "guidance/workflow.md"
    vars   = ["repo_name"]
    tags   = ["workflow"]
  }
}

dir "skills" {
  source = "skills"
  curate = "choose"

  node "review" {
    role    = "skill"
    content = "SKILL.md"
    default = true
    tags    = ["workflow", "quality"]
  }
  node "testing" {
    default = true
    tags    = ["lang/go", "quality"]
  }
  node "deployment" {
    default = false
  }
}

preset "recommended" {
  include "self:guidance"
  include "self:skills" { accept_defaults = true; }
}
```

Typed aliases normalize to `node "name" { type = "<alias>"; ... }`. Contradictory
explicit type is an error. Aliases are notation, not independent loaders. No
special skill/agent keyword or copy alias is required.

| Type | Structure and normal inclusion |
|---|---|
| `dir` | Physical directory inventory; include selected complete payloads |
| `file` | Physical regular file, extension-independent, byte-preserving |
| `md-file` | Physical Markdown file, intact frontmatter/bytes, optional heading inspection |
| `md-section` | Logical guidance group or Markdown document/heading view, rendered into a generated document |

- A physical root requires `source`, relative to the library file. Physical child
  `node` blocks annotate discovered paths and may omit type/source when discovery
  supplies them. Explicit type must agree with the discovered kind; missing paths
  are errors. Undeclared assets remain discovered, not hidden by the metadata list.
- An `md-section` branch can have logical child nodes without a source. An
  `md-section` leaf requires a Markdown `source` and optional `heading` array of
  exact title components. No array means the complete document view. Duplicate/
  ambiguous headings are errors. `md-file` plus `heading` is invalid; declare a
  separate section view without changing the intact file/bundle's identity.
- A source-free section group with no `title` is virtual and contributes no wrapper
  heading. A supplied `title` renders a wrapper. Source-backed sections use the
  selected source heading unless `title` explicitly overrides it. Same-title nodes
  from different libraries do not implicitly merge.
- Metadata fields: `type`, `source`, `title`, `tldr`, `tags`, `role`, `content`,
  `curate`, `default`, `reason_required`, `requires`, `conflicts_with`,
  `exclusive_group`, `into`, `vars`, `heading`, as applicable. Unknown fields,
  duplicate sibling names, incompatible kind/field combinations, or escaping paths
  are errors. Metadata remains tool-only, never inserted into intact copied bytes.
- Top-level nodes and presets share one library-visible name namespace; duplicate
  public names are errors. A preset keeps its members' original qualified node
  identities. Exclusions on a preset use those library-root paths (for example
  `skills/deployment`), not fictional files beneath a preset-named directory.
- `role` records intended use, not structural type or automatic harness conversion.
  Initial known roles are `instructions`, `skill`, and `agent`; unknown roles can be
  retained as annotations but do not gain implicit validation or placement rules.
- `content` is an optional primary file path relative to a physical directory;
  it must exist within it. It changes the viewer's default document, not the copied
  member set or output filename. Exact `SKILL.md` supplies the initial conventional
  bundle boundary/default content; explicit role=skill requires a valid main file.
  Lowercase DotAgents or vendor-format conversion remains a target-profile question.
- Named child declaration order is authoritative, followed by lexical discovered
  children not explicitly ordered. Tags inherit from ancestors by union. No implicit
  role, destination, or default is inferred from tag text.
- `curate` describes the branch's direct children. `default` is child-local: a
  choose branch requires a default on every discovered direct child before defaults
  can be accepted. An optional branch uses each child's default, with omitted
  defaults meaning false. Foundation and opt-in branches do not carry child defaults.
  A missing choose default is visible during inspection; it blocks broad/default
  adoption, not inspection or an exact independently resolved descendant selection.
- `reason_required = true` on a node makes explicit inclusion require a reason,
  independently of ordinary opt-in curation. This is the initial annotation for
  experimental/reason-required content; no extra curation enum or warning-per-asset.
- Relationship targets are library-local qualified node paths. `group "name"`
  has `mode = "at-most-one"`; nodes name it via `exclusive_group`. Missing targets,
  unknown groups, and hard-requirement cycles are diagnostics. Full enforcement
  uses the final effective result, not traversal order or first-match selection.
- Presets contain selection-only includes/replacements over `self:` node/preset
  references; they do not declare sources, mutate files, run scripts, or localize.
  Include cycles are errors. Consumer exceptions filter the selected qualified
  nodes before composition; preset default/reason decisions retain their origin.
  Cross-source recipe composition belongs to include-build, not hidden preset fetches.

Local raw directories remain usable without a library config; discovery provides
physical identity and recognized bundles, but no fabricated logical prose tree,
curation, or presets. Source aliases qualify normalized catalog nodes in all views.

### Guidance Variables

`vars = ["repo_name"]` on a source-backed `md-section` explicitly declares required
string inputs. Only these names are substituted in its source's `${name}` slots;
missing declared values are located errors. Other `${...}` text in Markdown stays
literal (for example shell examples); without vars, source text is not a template.
An explicitly supplied value not declared by selected guidance is an error. Literal
`\${name}` in source Markdown escapes a placeholder and renders `${name}` literally;
normal Markdown backslashes otherwise retain their meaning. Literal triple-single
inline text remains available to avoid interpolation. Physical file,
md-file, and dir payloads never receive template substitution; vars there is invalid.

For plan/direct authored includes, declared names bind from that include's private
file scope, optionally overridden by explicit `values` entries. Included recipes
do not inherit caller bindings. Presets may select templated sections but do not
read arbitrary caller values; a leaf's explicit vars contract is the only bridge.
Unused explicit values or names outside the selected leaves' declared vars are
errors rather than hidden configuration. Library identity/path metadata is not a
template. In referenced Markdown, template substitution is deliberately limited
to named strings and a documented literal-placeholder escape, not Go templates.
A declared source variable with no matching placeholder may warn as unused; it
must not trigger broad substitution of otherwise literal `${...}` text.

Generated builds record each section's values explicitly so execution does not
depend on removed plan bindings or substitute a different caller environment:

```hcl
include "team:guidance/workflow" {
  into = "AGENTS.md"
  values {
    repo_name = "Example"
  }
}
```

Source files and chosen values are integrity inputs. Missing or undeclared values
still fail when an authored build is executed directly. Parsed heading extraction
is applied after string substitution; diagnose newly ambiguous/invalid outlines
rather than selecting a different heading silently. Multiline binding values do
not receive another template pass.

`include-build "team:builds/common.mogent.build"` includes a source-qualified recipe
at its exact position. Raw local relative recipe paths may also be used. Each
recipe has private source aliases/bindings and an explicit header; recipe input
is pinned/captured like other source input. Recipes cannot implicitly access
caller bindings, fetch undeclared libraries, or write above the workspace.
Explicit parameter passing is deferred. Detect cycles with the full include chain.
Repeated recipe inclusion is not silently deduplicated. A generated build flattens
recipe actions with uniquely qualified source aliases and retained recipe provenance.

## 5. Include And Replacement

```hcl
include "team:skills" ["review", "testing"]

replace "local:skills/review"

include {
  into = "AGENTS.md"
  text = """
    # ${repo_name}

    Follow the repository's development conventions.
    """
}

include "team:guidance/workflow"
```

The child list selects named children/subpaths below the referenced root. It is
membership, not a consumer reordering instruction: library order still applies.
Unknown/duplicate/overlapping selections are diagnosed instead of double-emitting.
The default operation depends on catalog kind:

| Selected input | Natural inclusion |
|---|---|
| Logical guidance branch/leaves | Traverse selected leaves; render its hierarchy into the generated document |
| Physical file | Preserve complete file bytes/mode |
| Physical directory/bundle | Preserve selected relative layout and complete bundles |
| Inline text | Contribute to a generated document |

No public `mode`, `copy`, `render`, `append`, `document`, or localization-at-build
switch is needed. `replace "x"` normalizes to `include "x" { replace = true; }`;
`replace = false` on the alias is contradictory and an error. No `copy` alias.

Fields: `into`, `exclude`, `tags_all`, `tags_any`, `accept_defaults`, `replace`,
`reason`, `text` as applicable, with optional `values` for section inputs. A source
subject and `text` cannot be combined.
`exclude` is a list of source-relative subpaths, or named exception blocks when
a reason is needed. No separate global delete operation in this slice.

```hcl
include "team:guidance" {
  accept_defaults = true
  exclude "security" {
    reason = "Using the repository's reviewed security guidance."
  }
}
```

### Destinations And Composition

- Workspace-relative destinations only. Initial policy permits root `AGENTS.md`
  and paths beneath `.agents/`; never home paths, `.mogent`, source/config files,
  traversal, or symlink escapes. Do not manage the whole repository root.
- Logical guidance/inline text defaults to `AGENTS.md`. Physical paths exposed as
  `skills/review` default to `.agents/skills/review`. Other physical content needs
  a visible library-suggested workspace path or explicit `into`; unknown placement
  is an unresolved diagnostic, not a guessed filename. A physical root named
  `skills` has that conventional prefix; arbitrary root names do not gain it.
- Explicit `into` wins over the suggested/default destination. Directory `into`
  is the exact destination for its children, not an implied extra basename.
- Logical contributions to one document compose in action order, starting from an
  empty in-memory document each run. They do not append to yesterday's output.
- Complete physical files are not document fragments: collisions with each other
  or a generated document require explicit replacement, even for equal bytes.
- Disjoint physical directories may union. Replacing a destination removes its
  complete prior subtree, including stale bundle assets and document fragments.
  Replacing a missing destination is an error in this draft (makes typos visible).
- Installing only part of a recognized skill bundle or overlaying loose files
  inside one is an error. Bundle inclusion discovers supporting files automatically.
  Main-file display never collapses physical identity into heading ownership.
- Validate links from derived document views against the final output manifest;
  do not infer assets or rewrite links automatically. Copying never executes scripts.

Track explicit ownership roots separately from selected files. Compacting an
include must not broaden/narrow which on-disk directories Mogent owns. For the
pilot, root AGENTS.md and the included `.agents/skills` directory are separate,
non-overlapping managed roots. Never replace all `.agents` merely because one
skill is installed; unrelated out-of-scope artifacts remain untouched.

## 6. Selection, Curation, And Deterministic Ordering

Resolve candidates, make decisions, validate the final effective result, then
serialize/compose. Do not validate a requirement only when first encountered.

1. Candidate entrypoints are named nodes, child lists, tag queries, or named
   presets. A preset references existing selections and accepted choices; it is
   not another content copy. The referenced root bounds query scope; it is not an
   extra broad selection when a child list or tag query is supplied. With neither,
   include the root broadly under curation. With tags only, select matching units;
   with a child list only, select those children/subpaths. If both are supplied,
   explicit children and tag matches union with deduplication. `tags_all` and
   `tags_any` combine as AND-of-all and OR-of-any; queries use exact/slash-prefix
   matching as current v2 does. Physical bundle tags select complete bundles,
   not only the main Markdown file or scripts.
2. Explicit exclusions suppress matches/defaults within that include only. They
   never remove an earlier include's output. A later include may deliberately
   select the same node again; actual destination conflicts still need replacement.
3. Apply curation to acceptance from broad roots; exact descendant selection need
   not decide unrelated siblings. Tag selection may not bypass curation or reasons.
4. `accept_defaults = true` explicitly adopts a choose branch's complete defaults.
   Named child lists decide include for named direct children and exclude for
   omitted siblings; nested choices beneath selected children still need resolution.
   A sparse exclude-only selection is not automatically a complete choose decision.
5. Render children in library order; generated actions preserve consumer order,
   then library order, then bytewise UTF-8 lexical source-relative path order when
   no library order exists. Tags never reorder. No random ordering or mandatory
   before/after graph language.

| Curation/choice | Consequence |
|---|---|
| Foundation broad selection | Include offered children; exclusion needs nonblank reason |
| Choose branch | Complete explicit choices or explicit adoption of complete defaults; otherwise UNDECIDED |
| Optional branch | Apply published child defaults (false if absent) and show that origin unless explicitly overridden |
| Ordinary opt-in | Omitted = excluded; explicit selection sufficient, no reason required |
| Experimental/reason-required | Inclusion needs a nonblank recorded reason |

`reason_required` is the proposed inclusion-reason annotation in the new library
schema. It remains a final-review spelling, not an implemented feature.

Validate requirements, conflicts, and exclusive-group cardinality on effective
content after replacement. A reason for declining foundation does not waive an
unrelated hard requirement. A different source at the same output path does not
automatically satisfy a requirement on the original source identity. Replacement
does not prove semantic policy equivalence. Report blockers near affected nodes.
Relationship enforcement is new work, not something the legacy validator already
guarantees. Do not automatically resurrect excluded requirements or choose winners.

## 7. Sequencing And Input Integrity

Generated output is canonical `include` actions, no unresolved tags, curation,
bindings, recipe includes, or placement defaults. Exact source identities,
effective destinations, replacement flags, and resolved text are recorded. Physical
subtree shorthand is allowed when its complete selected content is fixed by the
recorded pin/hash and it preserves ordering/ownership/provenance.

Reserved metadata proposal:

```hcl
resolved {
  input "local:skills/review" {
    sha256 = "<node digest>"
  }
  decision "team:skills/deployment" {
    action = "exclude"
    reason = "explicit consumer exclusion"
  }
  owner ".agents/skills" {
    kind = "directory"
  }
}
```

Input records also cover sidecar/recipe bytes and rendered source files affecting
the result. Decisions identify their originating plan/preset/recipe action and
selection/default reason; source ranges and replacement history remain attached
in the shared inspection model. Do not hide superseded steps merely because final
bytes match. Generated builds must retain enough declared curation context to
check their accepted exceptions without re-running floating queries.
Preserve original guidance contribution groups and qualified ancestor paths so
expansion does not repeat wrapper headings for every leaf. Rendering selected
logical siblings once is not automatic same-title merging between libraries.
Input records for Git sources retain the requested URL/ref/subdir as well as the
resolved commit so pin reuse can distinguish a changed request from a stale plan.

- Initial sequencing can resolve a Git ref and record its commit in the build.
  Thereafter reuse that pin while source URL/ref/subdir is unchanged. Resequencing
  a plan must not silently advance remote pins. `source update --accept` is the
  deliberate update path. A missing build with floating sources is a bootstrap,
  not proof of historical reproducibility; check fails until records exist.
- Direct authored builds with Git sources require full commit pins. Local authored
  inputs may be captured during an in-memory resolution, but a reproducible saved
  sequence requires recorded local input hashes. Inspection distinguishes captured
  authoring inputs from recorded/verifiable inputs; no extra "frozen" mode.
  Proposed authoring path: `source update local` previews selected-node hash changes
  and `source update local --accept` writes only those reviewed input records into
  the authored build. Check reports missing records with this hint; ordinary build
  may use live local content but does not claim verified recorded inputs. This
  command uses the existing source-update category, not an extra build mode.
- Node hashes use standard SHA-256. Proposed directory recipe: hash UTF-8 JSON of
  an array of `[kind, relative_path, permission_octal_string, file_sha256_or_empty]`
  entries sorted by relative path, including root/empty directories and all selected
  members. No absolute paths, mtimes, UID/GID, or unselected sibling content.
  Sidecar selection-affecting data is recorded separately, not omitted from inputs.
- Hashing must cover bytes captured for assembly, not a later source reread. Detect
  unsafe/changed entries during capture; materialization uses captured payloads.
- Local Git HEAD and dirty status are informative alongside hashes. Dirty/untracked
  payload is visible; selected assets never disappear merely because untracked.
  Reproducibility from fresh checkout is claimed only with retained source bytes
  (committed files/exact Git sources or an explicitly retained snapshot).
- File-hash mismatch in a generated build is an input change error. With an
  unmodified generated build and its plan, normal build may resequence local input
  changes and write new records. `--check` only reports differences. Direct authored
  builds with recorded hashes never silently rewrite those records.

Explicit formatting emits lowest complete units, not partial skill files. Compact
formatting can compare parent-plus-exclusions against explicit children using a
deterministic instruction cost; ties prefer explicit. It must verify equivalence
of selection, order, destinations, ownership, and provenance. Do not infer new tag
intent or promise global minimum expressions. Aliases and formatting are independent
of diagnostic verbosity and validation.
When serializing resolved text, preserve literal `${...}`, backslashes, and triple
delimiters without re-interpolating them on the next parse; choose literal strings
or the documented escapes as needed. Text equivalence is byte-tested, not assumed.

## 8. CLI Contract And File Selection

All commands below are proposed new behavior. Obsolete consumers/commands can be
retired after replacement workflows pass; no compatibility requirement determines
this vocabulary.

| Command | Effect |
|---|---|
| `init [--plan]` | Create authored build by default or a plan; never overwrite |
| `source add <alias> <location>` | Add declaration to the chosen authored file, not select content |
| `source update [alias] [--accept]` | Preview pin + selection/output changes; accept updates records deliberately |
| `edit include <ref> [options]` | Add/edit a consumer include; defaults to normal include behavior |
| `edit include <ref> --localize [--replace]` | One-time local editable source derivation, never ordinary build execution |
| `edit set <name> <string>` | Set an authored binding |
| `show [node]` | Result tree by default, including unresolved callouts |
| `show --source [alias]` / `show --build` | Catalog/coverage or ordered action view |
| `show <node> --explain` | Library declaration, local override, effective values and origin |
| `sequence` | Compile plan to inspectable build file; no final output installation |
| `build` | Sequence if needed, then materialize through shared executor |
| `build --check` | Verify saved sequence/input freshness and compare expected managed output; no repairs |
| `library check/scan` | Validate/discover library content; scanning previews |
| `help [workflow|config|command]` | Structured help |

`--file` explicitly selects the entry file; mutations support `--dry-run`.
`edit include` options mirror config fields (`--into`, repeated `--exclude`,
`--tags-all`, `--tags-any`, `--accept-defaults`, `--replace`, `--reason`). CLI
`--localize` is authoring-only. No `--mode`, required `--copy`, or `--append`.

`sequence --format explicit|compact --no-command-aliases` controls written
representation. Explicit canonical output is the initial serializer; compaction
is a later slice. `--verbose` increases explanations, not semantic changes.

Default discovery is in the current directory, not an ambient ancestor search:

| Inputs | Default build behavior |
|---|---|
| Plan only | Resolve, sequence, capture result, then transactionally install build + outputs/state |
| Authored build only | Resolve/capture and install; do not rewrite authored file |
| Plan + unchanged derived build | Reuse pins; resequence if needed, then install |
| Plan + independent authored build | Ambiguous; choose explicitly with `--file` |
| Edited/unknown-origin generated build | Stop; do not overwrite even when `--file` names the plan |
| Neither file | Explain `init`; no writes |

`sequence` requires a plan. `show` may resolve a plan through an in-memory build
without writing it. Edits choose the plan when the only build is its unchanged
derivative; otherwise choose the single authored input or require `--file`.

Ordinary operation announces stages, not approval prompts. Ambiguity is an error
with recovery hints; any later interactive recovery must default to abort.
`--force` may permit reviewed output adoption/drift replacement, but NEVER bypass
invalid selection, path safety, source integrity, authored-build protection, or
localization overwrite refusal. It is not a generated-file regeneration flag.

`build --check` reconstructs using recorded pins, not latest remote refs. It
reports stale/missing generated build, changed inputs, unresolved choices, or
output differences with nonzero status. It may fetch an exact missing pin into
cache but does not mutate plan/build/output/state or repair files. Ignore absent
machine state when comparing a fresh checkout's expected/actual bytes/layout;
installation ownership still needs its own adoption contract.

## 9. Localization

```sh
mogent edit include team:skills/review --localize --replace
```

Capture the selected complete node into `.mogent/local/skills/review`, preserve
payload bytes/layout/modes, and record source identity, commit (if applicable),
baseline digest, and original source location. Store retained baseline bytes or
an exact retrievable Git reference if future upstream comparison is promised.
Without a retained baseline, report that limitation rather than pretending a hash
can reconstruct content. Do not insert tool provenance into skill frontmatter.

Declare/reuse a compatible `local` source; conflicting existing alias/path or
existing localized destination is an error. Rewrite the exact selected include,
or add a local replacement after its enclosing broad include. Preserve exclusion
and reason semantics; ambiguous target matches require explicit selection.
Preview source/config/provenance changes with `--dry-run`; commit those authoring
changes transactionally. No output installation until `build`. Later builds read
the local source; they never repeat localization or overwrite edits from upstream.

## 10. Inspection, Ownership, And Failure Behavior

Expose one structured inspection result for CLI/API/TUI: source nodes, selection
decisions, resolved steps, target files, heading views, shadowed contributions,
management roots, and diagnostics. Attach diagnostics to source ranges/node IDs;
message strings are not the machine API.

- Undefined alias/binding/selector: error. Unused declaration: warning.
- UNDECIDED required choice, missing reason/destination, conflict, input mismatch,
  missing requirement, malformed source: blocking error.
- Accepted defaults, opt-in omission, reasoned exception, dirty-source status:
  explained states, not automatically blockers.
- Continue independent inspection after errors where possible. Keep unresolved
  placeholders beside known targets or in an `Unplaced contributions` group.
  Do not invent output bytes or executable placeholders for unresolved steps.
- Show source/build/result views and library/local/effective values, including
  default destination provenance. `--issues` filters while retaining ancestor context.
- No generic warning for non-Markdown assets in normal copying. Do not manufacture
  sequence errors when authored/library/lexical order already resolves it.
- Every final destination has one owner. Capture/stage complete managed roots;
  write AGENTS.md, directory roots, generated build, and state with rollback.
  Failure restores touched outputs/build/state. Independent `sequence` does not
  install outputs and keeps the old build on failure.
- Do not silently delete an orphaned old management root after configuration
  changes. Show it as retained/orphaned and require a separately designed cleanup
  action. Within a still-managed directory, absent old members are removed.
- Existing hand edits, chmod changes, or extra directories inside managed roots
  require explicit output adoption/drift handling. Unrelated files outside those
  roots are untouched. Source symlinks/special files and output symlink escapes
  remain rejected. Missing state does not permit overwriting unrelated files.

No interactive TUI, protocol headers, global installation, symlink outputs, JSON
semantic merging, hooks activation, automatic upstream publishing, reverse plan
inference, or CI wrappers are required for the first end-to-end slice. Their
deferred status must not erase preserved metadata or future API needs.

## 11. Review Before Code

Review the proposed-detail table in section 1 and the acceptance matrix together.
The finite remaining review set is:

1. **Files and library schema:** mogent.lib/plan/build, typed node aliases,
   role/content annotations, child-local defaults, self-qualified presets. No
   new separate skill keyword or protocol choice is required.
2. **Literal/value contract:** immutable before-use string bindings, multiline
   indentation/escaping, explicit guidance vars, generated per-include values.
   Intact Markdown/files/bundles are not templates.
3. **Generated records:** format/version, origin/author, portable content checksum,
   plan/input hashes, decision and output-ownership records. Records required for
   fresh-checkout verification travel with the generated build, not only state.
4. **Lifecycle/recovery:** build sequences automatically only when authority is
   unambiguous; edited/authored builds block regeneration. Source update explicitly
   refreshes remote pins or authored local input hashes. Missing replacement target
   errors; output force does not authorize source/config overwrite.
5. **Localization:** one-time edit include --localize, editable local library plus
   provenance, normal local include afterward; no repeated copy during build.
6. **Implementation gate:** typed parser/catalog and meaningful API cases first,
   then shared composition/transaction, thin CLI/text inspection, localization,
   recipe/compact/cutover coverage. Old configs retained as migration evidence;
   new runtime need not maintain YAML.

The owner has now approved starting the implementation with `mogent.lib` naming
and the separator clarification above. No further product questionnaire is needed
before the syntax API slice. Do not claim the full spec delivered until the staged
acceptance cases pass. API implementation precedes focused API tests; CLI follows
the tested API, then CLI behavior tests. Expected CLI cases are documented before
implementation and may be encoded as executable contracts before CLI wiring.

## 12. Delivered Syntax API Boundary

The new public `syntax` package implements Parse, located ordered syntax, original
source/comment retention, explicit string substitution, and deterministic Format
with optional alias expansion. It checks format/header identity, syntax/closed
option names, duplicate declarations, and contradictory aliases. It performs no
filesystem I/O or network access and does not execute includes.

Catalog/type applicability, binding-use ordering, unused declarations, guidance
variable whitelists, library curation/relationships, node hashes, generated-file
checksums, transactions, and new CLI routing remain later slices. The presence of
recognized field spellings in the parser is not implementation of those semantics.
The formatter preserves all comments but can move comments from within an array
to their containing statement's end; Original returns exact input for a no-op
authoring edit. A source editor must use source ranges, not pretend canonical
formatting is byte-preserving. New syntax is not yet accepted by the executable.
