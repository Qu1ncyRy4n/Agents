# CLI Command Inventory

Mogent has one project language: `mogent.build` for an authored selection or
`mogent.plan` for selection intent that is sequenced into a build. Sources are
declared in that file; they do not select modules.

| Command | Role |
|---|---|
| `init` | Create an authored build or a plan |
| `source add`, `list`, `show`, `register`, `update local` | Connect, browse, remember, and refresh libraries |
| `module add`, `module rm` | Change direct selected modules |
| `add`, `rm` | Short compatibility forms for native module changes |
| `except`, `exclude`, `reorder`, `optimize` | Refine an authored selection |
| `show` | Inspect selected output, contributions, issues, and content |
| `sequence`, `build` | Create a generated build and install managed outputs |
| `status`, `drift` | Inspect output ownership, current composition, and upstream source changes |
| `edit include --localize` | Copy a selected input into an editable local source |
| `completion`, `mint` | Shell completion and handle utilities |

`show --tree` groups selected contributions by result output, with the source
shown as a right-side annotation. `drift` compares current source libraries with
the fingerprints saved by the last build, reporting changed, added, and removed
module paths.

The language supports `tags_all` and `tags_any` on an include, plus scoped path
`exclude` after tag matching. Tag-search and command-suggestion CLI surfaces are
not implemented yet; see [SOURCE-DISCOVERY.md](SOURCE-DISCOVERY.md).
