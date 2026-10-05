# HCL Configuration Spike

Status: superseded as the detailed design by
[`MOGENT-HCL-V2-SPEC.md`](MOGENT-HCL-V2-SPEC.md). Current Mogent reads
`agents.yaml`; it does not read or execute `mogent.hcl`.

The companion demo is
[`testing-ground/hcl-v2-demo/mogent.hcl`](../../testing-ground/hcl-v2-demo/mogent.hcl).
It retains the existing manifest's hierarchy: `sources`, `outputs`, and an
ordered `include` list for every output. Its small organization source is the
sibling `../cdint-demo-lib/` repository directory.

## Recommendation

Rename the product repository to `mogent` when its repository ownership and
remote migration are separately approved. Use `mogent.hcl` as the v2 consumer
configuration name. Do not rename the existing `agents.yaml` in place: it is a
material grammar and execution-model change, so v1 and v2 should remain
distinct until migration tooling exists.

HCL is worth evaluating because named blocks communicate the existing hierarchy
without YAML's indentation and collection edge cases. It does not by itself make
the system procedural: HCL parsing still produces syntax objects. The proposed
v2 remains a declarative build plan, preserving the current source and
multi-output model rather than introducing a general command language.

## Minimal Surface

`sources` declares explicit local paths or immutable Git commits. `outputs`
declares independent Markdown and directory targets. Each output has an ordered
`include` list. An entry selects one module (`source`), a named group/category
(`all`), or metadata (`tags { source, all, any }`). Existing safety rules remain:
refuse unsafe paths and mutable remote references, plan all outputs before
writing, then write outputs and state transactionally.

The QMR commit in the demo is real, but that library has not yet adopted the v2
module metadata used by the tag and accessibility/style examples. Those clauses
are deliberate target behavior, not claims about the current QMR tree. Its
existing `agents/workflow`, `agents/constraints`, and `skills` content is enough
to demonstrate pinned source retrieval once the v2 reader exists.

Do not add config imports, heredocs, variables, conditionals, arbitrary
expression evaluation, or imperative file/section replacement in the first
version. Those features create a different configuration language and should
need a demonstrated user story rather than inheriting authority from notes.

## Distilled Observations

- The consuming repository should own final composition. A source supplies
  modules and metadata, not a mandatory hidden configuration.
- Ordering must be visible when it changes rendered order or precedence. The
  `include` blocks above are ordered, but that does not require a procedural
  command language.
- Git commits are a useful reproducibility boundary for remote sources. Content
  hashes can still detect drift, but need not be a second versioning scheme.
- Existing-repository adoption is inherently iterative: build a plan, inspect
  the diff, revise selection, and only then take ownership of generated files.
- Tags and trees should solve concrete selection problems. Keep both narrow
  until real library users show that more metadata or graph behavior is needed.

## Open Decisions

- Define tag vocabulary and whether tag selection needs all-of, any-of, or both.
  The demo shows both forms but needs no implicit defaults.
- Define the exact Markdown section matching rules, including repeated headings
  and nested subsections.
- Define git-clean/adoption workflow before allowing the single-pass executor to
  manage an existing repository.
- Specify v1-to-v2 migration rather than accepting both grammars indefinitely.
