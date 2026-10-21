# Command Recovery Audit

Status: 2026-10-20 dogfooding audit. This maps older user-facing behavior into
the native plan/build workflow; it does not silently revive competing formats.

## Recovered Now

| Native behavior | Prior evidence | Current command |
|---|---|---|
| Source content preview | YAML `source show --content=none|snippet|full --lines`; HCL v2 source inspector body preview | `show source:path --content [--lines N]` |
| Physical Markdown heading addresses | HCL v2 `source show alias:path --headings` | `show source:path --content --headings` |
| Source hierarchy and display controls | YAML coverage/list tree and `display.chars`/alignment presentation | `source show [--tree|--flat] [--chars ascii|unicode]` |
| Add a selected item | YAML `add`, native include semantics | `add`, `select`, or `module add source:path` |
| Complete opt-in tree selection | Native child-selection syntax / physical bundles | `add source:path --all` |
| Authored output baseline | Generated-plan state writer / typed output manifest | `.mogent/state.json` after every native build |
| Output drift inspection | YAML drift expected-vs-disk comparison and state classification | `drift [--diff]` |
| Exclude from an ancestor selection | Native include `exclude` filter and older selection constraints | `except` or `exclude source:path` |
| Authoring declaration order | Native ordered syntax / formatter | `reorder [--dry-run]` |
| Source declaration/discovery | YAML `source add/list`; HCL v2 source list | `source add`, `source list|ls`, `source register`, `source list --known` |
| All managed-output status | Typed state plus `build --check` | `status` or `drift --all [--diff]` |
| Remove a direct include | Native declaration editing (previously absent from CLI) | `rm` or `module rm source:path [--diff]` |
| Overlap inspection | Explicit selection review (previously absent) | `optimize` (`optimise` alias); `show --build` warning |

Native content preview is read-only. It renders logical Markdown sections as a
build would, using entry-file bindings, while physical files keep their bytes.
Text can be limited with `--lines`; binary content is reported as opaque. No
preview selects content, writes an output, or creates state.

## Deliberately Not Aliased

| Capability | Existing behavior | Why it remains separate |
|---|---|---|
| YAML `add --under/--before/--after/--append` | Places rendered sections in `agents.yaml` | Native inclusion order is declaration order, not a manifest-heading placement model. `add --manifest` stays explicit. |
| HCL v2 `plan` / `apply` | HCL compile diff then transactional apply | Native plan/build has generated-build authority and different state/checksum rules. Do not conflate `plan` with `sequence`. |
| YAML `source show --metadata` / `--align-source` | Manifest-node metadata and rendered-heading alignment | Native catalog metadata surface is not yet equivalent; adding partial flags would promise unavailable semantics. |
| YAML TUI mutation | Edits a YAML manifest and builds after confirmation | Navigation/presentation is reusable, but its mutators are format-specific. Native interactive editing needs a typed result model first. |
| `exclude` as removal | Filters descendants from one include | It must not be repurposed to remove an earlier include or installed bundle. |

## Gaps To Decide At Dogfood Wrap-up

1. **Command vocabulary:** `add`, `select`, and `module add` are currently
   equivalent dogfooding aliases. Choose one canonical spelling before 1.0 and
   retain aliases only if they reduce real user ambiguity.
2. **Native metadata/search/TLDR/coverage:** older YAML/HCL browsing has richer
   filters and source metadata. Add through a shared native source/build/result
   inspector rather than routing native files through legacy manifests.
3. **Completion:** current zsh completion is basic commands/flags. Native source
   aliases and module references need project-aware completion after its safe
   discovery contract is designed.

## Evidence

- Legacy YAML source preview: `internal/cli/source_commands.go` (`runSourceShow`).
- HCL physical preview/headings: `internal/cli/v2_source.go`, `v2_content.go`.
- Existing tree presentation: `internal/cli/inventory.go` and
  `internal/presentation/config.go`.
- Native source/build inspection: `internal/cli/new_commands.go`.
- Native selection: `assembly/select.go`, `internal/cli/select_command.go`.
- Format boundaries and earlier gaps: `docs/CLI-COMMAND-INVENTORY.md`,
  `docs/proposals/MOGENT-FIRST-SLICE-SPEC.md`, and
  `TODO/TODO-config-first-slice.md`.
