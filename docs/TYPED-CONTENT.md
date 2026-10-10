# Typed Content Composition

Status: implemented first delivery under DR-juvih. The HCL extension was approved
by the owner before CLI integration. Existing v2 Markdown/raw-tree configurations
and v1 YAML remain supported.

## Start With The Runnable Demo

From the Mogent repository:

```sh
go run ./cmd/mogent source list demo --config testing-ground/typed-content-demo/mogent.hcl
go run ./cmd/mogent source show demo:physical/review/SKILL.md --headings --config testing-ground/typed-content-demo/mogent.hcl
go run ./cmd/mogent plan --config testing-ground/typed-content-demo/mogent.hcl
go run ./cmd/mogent apply --config testing-ground/typed-content-demo/mogent.hcl
go run ./cmd/mogent plan --config testing-ground/typed-content-demo/mogent.hcl
```

The final plan is unchanged. The generated `agent-export/` and `.mogent/` remain
ignored. Once installed with `tools/install`, replace `go run ./cmd/mogent` with
`mogent`.

## Named Physical Roots

A library may expose physical content alongside its existing authored guidance:

```hcl
library {
  format = 2
  id     = "example/content"
  name   = "Example Content"
}

content "physical" {
  directory = "payloads"
}
```

Mogent inventories every regular file and directory below that root. Named roots
must be unique and cannot duplicate the authored Markdown root or a legacy raw
tree name. Root/subdirectory traversal, symlinks, and special files are rejected.

Directories containing exact `SKILL.md` are inferred skill bundles. Other files
remain opaque regardless of extension. This first delivery does not add sidecar
`annotate`, typed `node` declarations, or vendor-specific metadata conversion.
The existing unlabeled `content { markdown_root = "agents" }` and `sections`
continue to describe authored guidance and its curation/defaults.

## Mixed Directory Outputs

```hcl
mogent { format = 2 }

sources {
  source "team" { local = "library" }
}

outputs {
  output "package" {
    path = "agent-export"
    kind = "dir-tree"

    source "team" {
      from      = "team:physical"
      node      = "review"
      operation = "copy"
      into      = ".agents/skills/review"
    }

    source "team" {
      from      = "team:physical"
      node      = "review/SKILL.md"
      heading   = ["Review", "Procedure"]
      operation = "render-markdown"
      into      = "AGENTS.md"
    }

    source "team" {
      from      = "team:physical"
      node      = "settings.toml"
      operation = "copy"
      into      = "config/tool.toml"
    }
  }
}
```

The review bundle is preserved, AGENTS.md is generated from one heading subtree,
and TOML is copied intact. Heading extraction does not edit installed SKILL.md.

| Attribute | Contract |
|---|---|
| `from` | Source alias and named root; authored roots remain available for rendering |
| `node` | Safe physical path relative to that root; defaults to `"."` for copy |
| `operation` | Required `copy` or `render-markdown` |
| `into` | Required safe path relative to the output directory; copied directories map directly here, files name their exact destination |
| `heading` | Optional array of exact parsed title components for a physical Markdown file |
| `replace` | Optional boolean; explicitly remove the previous contribution at this destination, including its descendants |
| `append` | Optional boolean; explicitly append to an existing rendered document, never an opaque copied file |
| `exclude` | Copy-only source-relative filter; does not delete earlier contributions |
| `select` | Existing sparse selection object for authored guidance; cannot combine with physical `node` |

Selecting directory `review` with `into = "skills/review"` puts its children in
`skills/review/`. `into = "skills"` would put those children directly in `skills/`;
Mogent does not infer a destination basename. Output `paths` can produce identical
independently managed packages. Output roots cannot overlap each other, the
consumer configuration, `.mogent/`, or the entire workspace root.

## Heading Views

`source show alias:root/path.md --headings` prints pasteable `heading` arrays.
The file keeps its original physical identity while a separate view exposes its
document outline. CommonMark ATX and Setext headings are parsed with Goldmark;
fenced examples and headings embedded in quotations/lists are not promoted into
the file outline. Duplicate sibling title matches are errors, not occurrence
guesses. Arrays permit titles containing slashes without ambiguous path splitting.

Rendering removes a leading delimited YAML frontmatter block from the derivative
only, rebases actual outline headings starting at level one, preserves body text,
and normalizes the result to a final newline. Frontmatter must terminate when the
file begins with `---`; copied files are never parsed/rewritten by this operation.
Heading levels above six after rebasing, non-UTF-8 document views, missing headings,
and ambiguous addresses are errors. Original copied Markdown/frontmatter remains
byte-for-byte intact.

Extracted link/image references to local files must resolve to explicitly copied
payloads or directories in the final manifest. To place a skill reference separately,
declare another named root over its references directory and copy that root; do not
silently install a partial skill. Mogent does not rewrite links, infer requirements
from inline command text, or provision external tools.

## Explicit Replacement And Append

Directories can combine at disjoint member paths. Same-file collisions fail,
including identical bytes. A complete skill bundle cannot silently absorb loose
files or be patched by an unrelated contribution. Select its whole directory;
exclude whole bundles rather than their required members.

Use `replace = true` on the winning contribution to replace a file or complete
subtree. Obsolete lower members disappear. `append = true` instead combines rendered
contributions to one document in declaration order; `replace` and `append` cannot
be combined. Granular edits inside an installed skill remain backlogged.

## Review And Ownership

`plan` shows member additions/changes/removals, rendered Markdown diffs, operation,
qualified source and local/Git context, exclusions, and replacement history.
Equal-byte replacements remain visible. Typed state tracks bytes, permission bits,
and all directories, including empty ones. Hand edits, chmod edits, and additional
empty directories require review and `apply --force`.

`apply` captures the complete manifest, materializes it into staging, installs the
directory, and updates state/pins through the existing rollback transaction. Source
payloads are not reread during manifest materialization. A partial failure restores
the previous output and state.

`update` compares old/new source manifests for typed packages and legacy raw trees,
including member and rendered-document changes. `update --accept` moves Git pins
only; `apply` writes the reviewed result afterward.

## Core API

The reusable `content` package supplies `Discover`, `Inventory.Document`,
`ParseDocument`, `Document.Resolve`/`Render`, `Compose`, and `Manifest.Materialize`.
`v2.LoadLibrary`, `Load`, `PlanConfig`, `Apply`, and `Update` adapt the HCL/workspace
workflow. Core unit tests precede CLI usage tests in this delivery.

Generalized nested selectors, custom role/bundle annotations, scoped lower-layer
exclusion blocks, vendor adapters, global installation, and configuration AST
merges are separate follow-up work. This extension does not claim a universal
agent runtime cascade or automatic loading of arbitrary `.agents` files.
