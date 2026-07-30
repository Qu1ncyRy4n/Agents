# TE: Manifest Schema For Milestone One

TE ID: TE-vorum
Status: concluded by DI-vukam
Related TODO: `jusuk.11`

## Decision Under Test

What YAML entry shape should mogent use for its first manifest parser and
renderer, which source kinds should that first version accept, and where should
the rebuilt Go module live?

This TE does not revisit the settled rules that the manifest owns rendered
headings, sources have explicit user-defined aliases, and source subtrees may be
composed. Source: DI-sufok.

## Assumptions

- `agents.yaml` is the manifest and `AGENTS.md` is generated output.
- A library is Markdown whose headings form a tree.
- A manifest entry can reference one node, pull a subtree, compose several
  source subtrees in order, or explicitly contain child entries.
- Manifest headings may differ from source headings.
- The first milestone must fail clearly for invalid manifests and write output
  only after successful validation.
- URL pinning, caching, and the interactive navigator are later milestones.

## Alternatives

### A. Compact, Polymorphic Entries

Use short scalar entries for simple references and mappings only when options
are needed.

```yaml
doc:
  - testing: cdint:instructions/testing
  - research:
      from:
        - cclab:instructions/testing
        - shared:instructions/testing
      exclude: [flaky-retries]
```

### B. Uniform Entry Objects

Every entry is a mapping with explicit fields. `heading` always controls the
rendered heading; `from`, `exclude`, and `children` describe composition.

```yaml
doc:
  - heading: Testing
    from:
      - cdint:instructions/testing
      - cclab:instructions/testing
    exclude:
      - cdint:flaky-retries
```

An entry can instead have `children` when the manifest explicitly defines its
own lower-level outline.

### C. Manifest Nodes As a Source-Like Tree

Make every manifest node resemble a library node, with a heading, body, source
reference, and nested content in one YAML shape.

```yaml
doc:
  heading: Instructions
  children:
    - heading: Testing
      source: cdint:instructions/testing
```

## Scenario Analysis

### Normal Use

Alice needs a small manifest that reads like her final `AGENTS.md`. Alternative
A is shortest for the common case. B is longer, but each entry has one obvious
shape. C makes the document tree explicit but adds structure even for a simple
one-line source pull.

### Two Sources In One Section

Bob combines general testing rules with lab testing rules under one manifest
heading. A supports this, but readers must learn both scalar and mapping forms.
B naturally makes `from` an ordered list. C can express the same idea, but the
source property competes with the manifest-node structure and invites unclear
rules for several sources.

### Invalid Input

Carol accidentally repeats a YAML key, gives `from` a scalar in a list-only
field, or combines `children` with an incompatible source pull. A has more
valid-looking shapes and therefore more validation branches. B gives the
validator one entry type and clear mutually-exclusive-field rules. C is also
regular, but has more fields that can be mistaken for authored document text.

### Local Overrides And Editing

Dave edits a shared source. Mogent creates a local copy and changes one source
reference. B has one predictable place to replace the reference. A must locate
either a scalar or a mapping. C has a predictable field too, but it blurs the
line between library content and the manifest, which the design deliberately
separates.

### Long-Term Evolution

Ellen later adds per-source exclusions, source-change warnings, comments, and
display options. A becomes harder to extend without adding further special
forms. B has room for optional fields while retaining a strict schema. C also
has room, but risks re-creating a second content authoring format in YAML.

### Trust Boundaries And Reproducibility

Frank starts with local source paths. URL support later adds fetching, caching,
pinning, and changed-source confirmation. All three schemas can name URL
sources. Restricting milestone one to local paths keeps failures, source
inspection, and test fixtures simple while the renderer becomes reliable.

### Scale And Navigation

At many sources and nested sections, source provenance and document shape must
be easy for a navigator to display. B directly maps to a tree row with a
heading, sources, exclusions, and children. A requires the navigator to first
normalize several shapes. C is easy to walk, but it duplicates the library's
tree model in the manifest and makes content-vs-composition harder to explain.

## Conclusions

Reject C. It makes the manifest look too much like a second library and weakens
the design's clear separation between reusable content and a document outline.

Alternative B best supports strict validation, source composition, copy-on-
write edits, and a future navigator. Its extra YAML is acceptable because the
manifest is a durable, reviewable description of the generated document.

Alternative A remains attractive only if real use shows the uniform form too
verbose. A serializer or `mogent fmt` could later offer compact display without
making compact syntax part of the parser contract.

For milestone one, accept local source paths only. URL sources need their own
fetch, cache, pinning, and confirmation design; pretending that fetching is a
small parser feature would make the first renderer less reliable.

## Questions For Decision Framing

1. Use uniform entry objects (B), or retain compact scalar shorthand (A)?
2. In milestone one, allow only local source paths, or also allow unpinned URLs?
3. Place the rebuilt Go module at `tools/mogent/`, or make the repository root
   the Go module?
4. May an entry have both `from` and `children`? The recommended first rule is
   no: use one or the other until a concrete composition case proves both are
   needed.

## Implications

After these choices are locked in a DI, implementation can create the Go module,
the manifest parser, source resolver, template renderer, validator, fixtures,
and a local dogfood manifest. Source composition and URL support beyond this
scope remain explicit follow-up work.

## Decision Result

DI-vukam selects a root Go module; accepts both compact and explicit entry
forms; accepts local source paths only; and keeps `from` exclusive with
`children`. Compact forms normalize to the explicit internal tree, preserving a
single renderer model. `docs/IMPLEMENTATION-M1.md` defines the full milestone-
one contract. URL sources require pinning and changed-content review before they
are enabled.
