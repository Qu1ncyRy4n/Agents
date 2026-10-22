# Current CLI Command Inventory

Status: inspected implementation baseline for the configuration redesign.
Checked against checkout HEAD `7819114`; documentation changes are uncommitted.
This is an inventory, not approval of replacement commands or new grammar.

Current cutover correction: `init` now creates authored `mogent.build` by
default; `--plan` creates `mogent.plan`. Sources are declarations only;
`--dry-run` performs no writes/fetches. Legacy initialization requires explicit
`--config mogent.hcl` or `--manifest agents.yaml`; `--list-templates` lists legacy
starters. The inventory below otherwise retains its historical baseline (the
new runtime also registers `sequence`, `show`, and `edit`).

Source discovery update: plan/build `source add`, `source list|ls`, user-level
`source register` / `source list --known`, optional Gum initialization, and the
typed catalog tree are implemented. See [SOURCE-DISCOVERY.md](SOURCE-DISCOVERY.md)
for the current contract and supported flags. Older verification/gap notes below
describe the original redesign baseline where explicitly identified.

The [command recovery audit](COMMAND-RECOVERY-AUDIT.md) maps older content
preview/browsing/mutation behavior into native equivalents and records deliberate
format boundaries plus the dogfood decisions still required.

## Registered Commands

`internal/cli/cli.go` supplies dispatch, top-level usage, command-name completion,
and typo suggestions. It has 19 commands plus the special `help` entry point.
`internal/cli/help.go` owns the workflow overview and stage-specific help;
`-h` and `--help` are the same concise overview, `help --all` expands the workflow
with common flags, and `help commands` renders every dispatch-registry signature.

| Command | Implemented role | Format boundary |
|---|---|---|
| `init` | Authored build or `--plan`, source declarations, dry-run; explicit legacy starters | Plan/build default; `--config` selects HCL v2, `--manifest` selects YAML v1 |
| `build` | Sequence plans and assemble builds; legacy render/write outputs | Plan/build with discovery or `--file`; YAML v1 with `--manifest` |
| `sequence` | Compile selection intent into a generated build | New plan/build language |
| `show` | Rendered/physical content preview, heading-oriented result tree, ordered contributions, diagnostics | New plan/build language; `show --tree` groups selected headings by output with an adjacent source column; catalog browsing is `source show` |
| `edit include --localize` | Copy a selected input into an editable local source | New plan/build language |
| `plan` | Resolve, validate, show output diffs/member changes/provenance/diagnostics | HCL v2; does not emit a reusable build file |
| `apply` | Transactional outputs/state and newly resolved commit pins | HCL v2 |
| `update [alias]` | Preview source-pin/output changes; `--accept` updates pins | HCL v2 |
| `lib init/scan/check` | Create, inspect, validate library YAML sidecars | Legacy schema, not v2 HCL-sidecar editing |
| `status` | Output ownership/freshness, source summary, next-action hints | YAML v1 |
| `drift` | Native read-only output-baseline/current-build inspection with optional AGENTS diff; legacy import/reject behavior | Native plan/build by discovery or `--file`; YAML v1 with `--manifest` |
| `coverage` | Source coverage tree convenience entry point | YAML v1 |
| `source add` | Declare a source with dry-run; does not select content | Plan/authored build; explicit legacy YAML |
| `source list` / `source ls` | List project declarations; `--tree` browses native content | New format; legacy list browsing remains separately routed |
| `source register` / `source list --known` | Remember/discover libraries; never active build inputs | User-local YAML registry |
| `source show` | Browse the native source catalog; inspect one node with legacy metadata/content flags | Native catalog; legacy YAML/HCL inspection is retained under explicit legacy configuration |
| `source pin/update` | Pin or preview/accept URL source changes | YAML v1; distinct from v2 top-level `update` |
| `module add` / `module rm` | Validate and append or remove one native include; `--all` explicitly adds opt-in child trees; `--diff` previews project/output changes | Plan/authored build by discovery or `--file`; top-level `add` / `rm` are compatibility shortcuts; YAML v1 uses `add --manifest` |
| `localize` | Copy-on-write override with provenance | YAML v1 |
| `complete` | Print completion candidates | Commands/flags general; source/heading candidates v1 |
| `completion` | Emit bash/zsh scripts | Scripts/candidates need v2 reconciliation |
| `tui` | Manifest editing, final Markdown, selected-source context | YAML v1 |
| `mint` | Proquint handle utility | Separate from assembly semantics |
| `help`, `--help`, `-h` | Workflow overview; focused init/source/select/show/build topics with advanced sections | Legacy/tools topics separate; `help commands` lists all signatures; command `--help` routes to its stage without executing |

No top-level `tree`, `render`, `diff`, `validate`, `compile`, `include`,
`replace`, or `explain` command is registered. New file roles do not automatically
change existing command meanings.

## Inspection Features Worth Reusing

- **V1 coverage:** `source list --coverage --tree` shows included, inherited,
  partial, excluded, and unused states. Supports tag/search filters, content/leaves/
  depth filters, TLDRs, terminal fitting, and ASCII/Unicode presentation. It is a
  source-selection tree, not a mixed filesystem result tree.
- **V1 TUI:** manifest-tree, final-AGENTS.md, and selected-source panes; collapse,
  navigation, add/remove, sibling reordering, draft discard, confirmed save/build.
  Loads a valid workspace first; not an incomplete typed-result authoring viewer.
- **V1 addition preview:** `add --dry-run --preview=tree|patch|full|summary`
  displays placement, inherited subtrees, and related selections.
- **V2 sources:** `source list [alias] --tldr --tags` shows hierarchy, curation,
  child defaults, raw trees, and physical contents. `source show alias:path`
  inspects a node; `--headings` discovers physical Markdown heading addresses.
  V1 coverage/search/tree flags do not have parity here.
- **V2 result preview:** `plan` shows unified Markdown diffs, added/changed/removed
  members, copy/render origins, replacement/append history, and stable diagnostic
  codes. It is not a full interactive result tree.

Only `source list/show` auto-route to v2: explicit `--config`, or `mogent.hcl`
present without `agents.yaml`. Other source subcommands remain v1. V2 browsing
calls `ResolveSources`; do not assume inspection is always network-free.

## Gaps Relevant To The Agreed Direction

1. `v2/build.go` returns before rendering/composition if compile diagnostics contain
   errors. `v2.Diagnostic` holds severity/code/message/suggestion, not structured
   node/range attachments. Visible unresolved result nodes need a partial
   inspection model, not only a new formatter. Source-load failures also need an
   inspection strategy. Invalid previews must never become applicable builds.
2. Show blockers at relevant nodes and summarize them; use an unresolved group
   when the destination is not known. Distinguish defaults/exclusions from errors.
   Do not demand before/after declarations where accepted deterministic order
   already resolves the sequence.
3. Reuse source coverage, presentation, TUI, and typed manifests through one
   inspection model for source/build/result and text/interactive clients.
4. `plan` still previews HCL v2 output. New compilation/assembly is implemented
   through `sequence` and `build`; inspection uses `source show`, `show --build`, and `show --issues`,
   not the proposed `show source/build/result` positional vocabulary.
5. Workflow help and focused stage topics are implemented. Unknown help topics
   now fail with guidance rather than silently displaying the same list. Registry
   signatures are not an exhaustive per-handler flag reference; full config help,
   richer per-command option documentation and newer-format completion remain gaps.

## Verification

Ran these successfully with pure-Go compilation and Go caches under `/tmp`:

```sh
go run ./cmd/mogent help
go run ./cmd/mogent plan --config testing-ground/typed-content-demo/mogent.hcl
go run ./cmd/mogent source list demo --config testing-ground/typed-content-demo/mogent.hcl --tldr --tags
go run ./cmd/mogent source list --manifest testing-ground/basic-org/agents.yaml --coverage --tree --depth 1
```

Typed plan showed five output files, a rendered heading derivative, operation
origins, and local-source context. V1 coverage showed partial/included/unused tree
states. No apply or TUI save was run.
