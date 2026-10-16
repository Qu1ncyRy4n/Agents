# Current CLI Command Inventory

Status: inspected implementation baseline for the configuration redesign.
Checked against checkout HEAD `7819114`; documentation changes are uncommitted.
This is an inventory, not approval of replacement commands or new grammar.

## Registered Commands

`internal/cli/cli.go` supplies dispatch, top-level usage, command-name completion,
and typo suggestions. It has 16 commands plus the special `help` entry point.

| Command | Implemented role | Format boundary |
|---|---|---|
| `init` | Editable starters, template listing, dry-run | Defaults to HCL v2; explicit `--manifest` selects YAML v1 setup |
| `build` | Render/write outputs; dry-run/force | YAML v1; not a plan-to-build compiler |
| `plan` | Resolve, validate, show output diffs/member changes/provenance/diagnostics | HCL v2; does not emit a reusable build file |
| `apply` | Transactional outputs/state and newly resolved commit pins | HCL v2 |
| `update [alias]` | Preview source-pin/output changes; `--accept` updates pins | HCL v2 |
| `lib init/scan/check` | Create, inspect, validate library YAML sidecars | Legacy schema, not v2 HCL-sidecar editing |
| `status` | Output ownership/freshness, source summary, next-action hints | YAML v1 |
| `drift` | Inspect/diff direct edits; `--import` localizes; `--reject --force` rebuilds | YAML v1; includes mutation |
| `coverage` | Source coverage tree convenience entry point | YAML v1 |
| `source add` | Declare a source with dry-run; does not select content | YAML v1 |
| `source list/show` | Browse sources or inspect one node | Both, through different handlers |
| `source pin/update` | Pin or preview/accept URL source changes | YAML v1; distinct from v2 top-level `update` |
| `add` | Add selection with heading/placement, previews, optional rebuild | YAML v1 |
| `localize` | Copy-on-write override with provenance | YAML v1 |
| `complete` | Print completion candidates | Commands/flags general; source/heading candidates v1 |
| `completion` | Emit bash/zsh scripts | Scripts/candidates need v2 reconciliation |
| `tui` | Manifest editing, final Markdown, selected-source context | YAML v1 |
| `mint` | Proquint handle utility | Separate from assembly semantics |
| `help`, `--help`, `-h` | Top-level usage | No workflow/config help topics |

No top-level `show`, `tree`, `render`, `diff`, `validate`, `compile`, `include`,
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
4. Preserve existing command meanings until explicit migration: `plan` currently
   previews v2 output; `build` writes v1 output. Neither currently implements the
   new compilation stage. Proposed `show source/build/result` is not implemented.
5. Top-level `help <anything>` prints the same usage list. Structured workflow,
   config, and command help would be new functionality. Top-level usage also omits
   many existing flags such as drift import/reject and add preview/rebuild.

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
