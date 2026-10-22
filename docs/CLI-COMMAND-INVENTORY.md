# CLI Command Inventory

Mogent has one project language: `mogent.build` for an authored selection or
`mogent.plan` for selection intent that is sequenced into a build. Sources are
declared in that file; they do not select modules.

| Command | Role |
|---|---|
| `init` | Create an authored build or a plan |
| `source add`, `list`, `tree`, `search`, `content`, `register`, `update local` | Connect, browse, search, remember, and refresh libraries |
| `module add`, `module rm`, `module exclude`, `module localize`, `module optimize` | Change or refine direct selected modules |
| `add`, `rm` | Short compatibility forms for native module changes |
| `except`, `exclude`, `reorder`, `optimize` | Compatibility aliases or project-wide selection utilities |
| `show tree`, `show build`, `show issues` | Inspect selected output and contributions |
| `sequence`, `build` | Create a generated build and install managed outputs |
| `status`, `drift` | Inspect output ownership, current composition, and upstream source changes |
| `completion`, `mint` | Shell completion and handle utilities |

`show tree` groups selected contributions by result output, with the source
shown as a right-side annotation. `drift` compares current source libraries with
the fingerprints saved by the last build, reporting changed, added, and removed
module paths.

The language supports `tags_all`, `tags_any`, and `tags_none` on an include,
plus scoped path `exclude` after tag matching. `source search` exposes these
predicates read-only; `--add-commands` prints explicit module-add commands.
`module add` accepts the same tag predicates for a new include. See
[SOURCE-DISCOVERY.md](SOURCE-DISCOVERY.md).
