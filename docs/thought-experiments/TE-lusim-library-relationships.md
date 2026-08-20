# TE-lusim: Library Relationships And Composition Diagnostics

TE ID: TE-lusim
Date: 2026-08-20
Status: complete analysis; awaiting owner decisions in `DR-lusim`

`tools/mint-handle` remains unavailable, so `lusim` is manually assigned under
the repository's existing user-approved exception.

## Decision Under Test

How should Mogent represent and validate prerequisites, discovery links, exact
conflicts, mutually exclusive alternatives, and provenance without coupling a
library to aliases chosen by a consuming manifest?

The same decision must say what happens when source paths move and whether a
relationship produces information, a warning, or an error.

## Constraints Already In Force

- Manifest source aliases are user-defined and local to one manifest.
- Source provenance must remain visible; Mogent never silently chooses a value.
- Missing or ambiguous source nodes are errors.
- Metadata is tool-only and does not render into `AGENTS.md`.
- The manifest owns the output outline. Relationship metadata may validate a
  selection, but must not silently add headings or choose placement.
- Sources are untrusted. Metadata must not grant access to undeclared sources,
  paths outside a source root, a network, or a writable location.
- The CLI and future TUI must use the same public workspace diagnostics.
- Existing `requires` and `conflicts_with` values are displayed only. Enforcing
  them is a public behavior change and needs a compatible transition.

## Relationship Meanings

The source review identified five meanings that should not be conflated:

| Relationship | Meaning | Selection effect |
|---|---|---|
| `requires` | The module is incomplete or unsafe without another module. | Selection is invalid while the requirement is absent. |
| `see_also` | Related material may help discovery or understanding. | Informational only. |
| `conflicts_with` | Two known exact modules cannot form one coherent instruction set. | Selection is invalid while both are active. |
| `alternative_family` | Several modules are choices for one role, such as communication persona. | More than one selected family member is invalid. |
| `provenance` | A rationale, decision, or source record explains why the module exists. | Informational only; never selects content. |

A soft tension is not `conflicts_with`. If authors cannot name the incompatible
module exactly or cannot guarantee incompatibility, they should explain the
relationship through `see_also`, prose, or an advisory diagnostic.

## Reference Alternatives

### A. Store Manifest Aliases In Library Metadata

```yaml
requires: [shared:process/decision-first]
```

This is readable inside one manifest but not portable. A consumer may name the
same library `cdint`, `org`, or `policy`. Renaming an alias makes valid library
metadata appear broken even though the source bytes did not change.

This also reverses ownership: a shared library would dictate a name that the
design assigns to the consumer.

### B. Use Bare Paths And Search Every Source

```yaml
requires: [process/decision-first]
```

This is concise but ambiguous as soon as two declared sources contain the same
path. First-wins, last-wins, source-order, and closest-match behavior all violate
Mogent's fail-loud rule. Asking the user on every build would make deterministic
rendering interactive.

### C. Give Every Library A Stable Global Identity

```yaml
requires: [library:cdint/process/decision-first]
```

A durable library identity could support cross-source relationships, versioned
dependencies, and registries. It also introduces a new public descriptor,
identity ownership, collision policy, fork behavior, version compatibility, and
trust rules. A copied or forked library must decide whether it retains the
identity, derives a new one, or impersonates the original.

This may become useful with remote registries or Promise Grid, but it is more
machinery than current local and pinned-Git sources require.

### D. Keep Every Relationship In The Consuming Manifest

The manifest can safely use its own aliases:

```yaml
relationships:
  conflicts:
    - [cdint:process/full-governance, personal:process/fast-iteration]
```

This is the right ownership boundary for a relationship between independent
sources. It is repetitive for relationships intrinsic to one library and lets
consumers accidentally omit a module's own hard prerequisite.

### E. Hybrid: Reserved Same-Source References Plus Manifest Cross-Source Rules

Library metadata uses a reserved pseudo-alias for paths in the same source:

```yaml
requires: [self:process/decision-first]
see_also: [self:engineering/comment-intent]
conflicts_with: [self:process/fast-only]
alternative_family: communication/persona
provenance: [DI-venit]
```

`self:` is interpreted only inside source metadata. It means the source
instance containing the module; it is not a manifest alias and never searches
other sources. A manifest remains free to bind that source as `cdint`, `org`, or
another name.

Relationships between independent sources are declared by the consuming
manifest, where aliases are meaningful. Cross-source library dependencies are
deferred until Mogent has evidence that stable library identities are worth
their trust and versioning cost.

This hybrid keeps intrinsic knowledge with the library and integration
knowledge with the integrator.

## Alternative-Family Shape

Two shapes were considered:

```yaml
family: communication/persona
variant: operator
```

and:

```yaml
alternative_family: communication/persona
```

The second is sufficient. The selected node is already the variant and has a
path, heading, TLDR, and provenance. Repeating a variant name creates another
identity that can drift.

For the first implementation, a family is scoped to one source instance. Two
members of `communication/persona` in that source conflict. A same-named family
in an unrelated source is not assumed to be compatible or incompatible.
Cross-source alternatives can later be declared in the manifest.

Current metadata is file-level and copied onto every heading in a file. An
alternative family therefore cannot be applied naively to today's multi-persona
file: it would label the parent, every persona, and the customization note as
separate alternatives. The first family dogfood needs one of two prerequisites:

1. split each persona into an atomic metadata-bearing file; or
2. implement heading-level metadata first.

Atomic files are the smaller and already preferred library-authoring change.
For this field, Mogent should treat the file's single top-level heading as the
variant root; its descendant headings inherit that variant identity rather than
counting as additional alternatives. A file declaring `alternative_family`
with zero or multiple top-level headings is invalid until heading-level metadata
exists.

The source paths can remain stable: keep the `communication/personas` parent
heading in its current file, create a matching `communication/personas/`
directory, and give each child file its existing persona heading. Mogent's
directory-plus-heading identity model preserves paths such as
`communication/personas/operator`, while each persona may still contain
subsections beneath its one variant root.

## Enforcement Alternatives

### Advisory Forever

Mogent would display every relationship but never reject a composition. This is
backward compatible, but `requires` and `conflicts_with` would overpromise. A
manifest could knowingly omit a safety prerequisite and still build cleanly.

### Strict Immediately

Mogent would reinterpret every current string as enforceable. Existing
alias-bound metadata would become invalid, and external libraries may already
use those fields as human notes. This would turn a documented display field
into a breaking build rule without migration.

### Recognized Strict Syntax With A Migration Period

Only unambiguous new relationships are strict:

- `self:` requirements and conflicts are validated;
- `alternative_family` is validated;
- `see_also` and provenance remain informational;
- legacy alias-bound relationship strings remain visible and receive a
  deprecation warning until migrated;
- unknown or malformed reserved syntax is an error, never silently advisory.

This gives authors an explicit opt-in to real semantics and permits the checked-
in libraries to migrate in focused commits.

## Diagnostic Behavior

The workspace API should return structured diagnostics rather than CLI prose.
At minimum:

| Code | Severity after migration | Example meaning |
|---|---|---|
| `missing_requirement` | error | A selected node requires an unselected node. |
| `selected_conflict` | error | Two selected nodes declare an exact conflict. |
| `multiple_alternatives` | error | Two selected nodes belong to one source-scoped family. |
| `legacy_relationship` | warning | Metadata still uses a manifest-style alias. |
| `deprecated_source_path` | warning | A manifest uses a compatibility path scheduled for removal. |
| `related_module` | information | An unselected `see_also` target is available. |

Each human diagnostic should include:

- the selected manifest heading or source reference;
- the related source reference as bound in the current manifest;
- why the relationship matters;
- whether the operation was blocked;
- the safest next command or explicit placement choice.

Hints remain presentation and may be disabled. Diagnostic codes, severity, and
machine-relevant fields must remain available regardless of prose-hint settings.

## Scenario Analysis

### Normal Selection

A manifest binds the same library as `org`. Selecting
`org:go/development` finds `requires: self:shared-baseline/...` in that source.
The workspace reports the requirement using the manifest's actual `org:` alias.
The library remains unchanged if another manifest calls it `cdint`.

`mogent add` previews the missing requirement before changing the manifest. It
must not guess where the required module belongs in the output outline. It can
offer `--with-requirements` only when placement is already determined by an
explicit parent or a reviewed manifest rule; otherwise it prints the exact
missing reference and asks the user to choose placement.

### Subtree Selection And Exclusion

A requirement is satisfied whether its node is explicit or inherited through a
selected subtree. An excluded descendant is not selected. If another active
node requires that excluded descendant, validation reports the exclusion as the
reason the requirement is missing.

Conflicts and alternative families apply to every active descendant, not just
explicit manifest entries. Exact conflicts are symmetric at validation time:
one module naming another is sufficient. Selecting a parent subtree that
contains two alternatives is invalid and should be explained during preview.

### Alias Rename

Changing a manifest alias does not affect `self:` relationships. Manifest-owned
cross-source rules must be rewritten with the alias because they are part of the
same manifest transaction. Validation catches any stale reference before write.

### Missing Or Renamed Target

A malformed `self:` path, missing target, self-conflict, empty family, or family
on a file without exactly one top-level variant root is a source validation
error. A soft link does not block selection, but a misspelled `see_also` target
still makes the source invalid; "soft" describes selection effect, not
permission to carry a broken reference.

### Cycles

Requirements may form a cycle when every member is a co-requisite. The cycle is
valid if every requirement is selected; it does not need to block a build.
Source inspection should show the cycle because authors may prefer one composed
subtree when the modules are never useful separately. An authoring convenience
that computes requirement closure must track visited nodes, add each member
once, and never recurse forever.

`see_also` cycles are harmless. Conflict and alternative relationships are
evaluated as sets rather than traversed as dependencies.

### Concurrent Local Source Edit

A local source can change between preview and apply. Mutation operations must
reload and revalidate the manifest and sources immediately before commit. If the
relationship graph changed, the operation stops and returns a stale-preview
diagnostic. Pinned URL sources retain their existing lock/hash guarantees.

### Untrusted Source

Metadata cannot declare a new source, fetch a dependency, or read outside its
root. `self:` resolves only through the already loaded source index. A manifest-
owned cross-source rule can mention only aliases already declared by that
manifest.

### Duplicate Mounts And Forks

If identical bytes are mounted twice under different aliases, each binding is a
separate source instance in the first implementation. Same-source relationships
resolve within the binding that supplied the selected node. Mogent does not
infer equivalence from paths, URLs, or content hashes.

A fork can retain its internal `self:` relationships without claiming to be the
original library. This is simpler and safer than a premature global identity.

### Scale

Source-load validation indexes paths, families, and relationships once.
Composition validation then walks active nodes and their declared edges. The
expected cost is linear in nodes plus relationship edges; no semantic text
comparison or global registry search is required.

## Source-Path Evolution

Heading renames can preserve identity with the existing inline `id` escape
hatch. The first clarity cleanup demonstrates this: visible headings changed
while their `corpus-variants` source path segments remained stable.

Directory moves change source paths. Four approaches were considered:

1. break consumers and publish a manual old-to-new table;
2. duplicate old and new modules temporarily, risking prose drift;
3. add a library-level alias/migration map;
4. avoid moves until a versioned source descriptor and migration contract exist.

The near-term recommendation is option 4. Improve names with explicit IDs, add
new modules where needed, and do not move existing directories during the
relationship-metadata slice. A later source-evolution decision should compare a
library descriptor with source-authored path aliases and a `mogent migrate`
operation. Until then, a directory move is a breaking change that needs explicit
owner approval and a reviewed manifest migration.

Relationship metadata should use canonical current paths. It must not make an
old path silently resolve to an arbitrary new node.

## Warning Narrative

Warnings should explain a problem, consequence, and next action. For example:

```text
Cannot add personal:communication/personas/operator.

The selected source already contributes
personal:communication/personas/research-tutor, and both belong to the
communication/persona alternative family. Rendering both would give the agent
two competing voices.

Keep the current persona, or remove it and retry the add. No files changed.
```

For a legacy alias-bound relationship:

```text
cdint:go/development declares requires: shared:shared-baseline/...

Library metadata cannot rely on the consuming manifest's alias "shared".
The relationship is shown for review but is not enforced. The library author
should migrate it to self:shared-baseline/....
```

## Conclusions

Rejected:

- manifest aliases embedded in reusable library metadata;
- bare path lookup across all sources;
- global library identities in the next slice;
- silently selecting prerequisites or resolving alternatives by source order;
- permanently advisory semantics for fields named `requires` or
  `conflicts_with`;
- immediate enforcement of every legacy metadata string;
- source-path-changing reorganization in the same change as relationship
  semantics.

Surviving recommendation:

1. Reserve `self:` for same-source metadata references.
2. Keep cross-source integration relationships in the consuming manifest.
3. Add `see_also`, `alternative_family`, and informational `provenance` fields.
4. Treat the family-bearing module root as the variant identity; do not add a
   duplicate variant identifier, and treat its descendants as the same variant.
5. Make recognized `self:` requirements, exact conflicts, and source-scoped
   alternative families strict after a migration period.
6. Keep `see_also`, provenance, and uncertain semantic tensions non-blocking.
7. Return structured public-API diagnostics with optional human hints.
8. Preserve heading paths with explicit IDs and defer source-path-changing
   directory moves until a separate source-evolution contract exists.

## Implementation Sequence If Approved

1. Record a Decision Intent resolving `DR-lusim`.
2. Extend metadata parsing for `see_also`, `alternative_family`, and
   `provenance`; define and validate reserved `self:` references.
3. Add source-graph validation and structured diagnostic types to the public
   library/workspace API.
4. Show relationships and legacy warnings in `source show`, source selection,
   and status without blocking existing manifests.
5. Migrate checked-in alias-bound metadata to `self:` in focused library
   commits.
6. Enable strict composition validation for recognized hard relationships.
7. Add CLI previews and actionable placement hints; do not auto-place required
   modules without an explicit rule.
8. Split persona children into atomic files without changing their source
   paths, mark them as the first alternative family, and dogfood conflict
   diagnostics.
9. Design manifest-owned cross-source relationships only after a real
   composition needs them.

## Evidence Required

- the same source mounted under two different aliases;
- requirements satisfied explicitly and through subtree inheritance;
- a required node removed by `exclude`;
- missing, malformed, cyclic co-requisites, and self-conflicting relationships;
- two selected family members;
- file-level metadata eligibility and path-preserving atomic persona files;
- legacy alias-bound metadata during migration;
- local-source mutation between preview and apply;
- public workspace and equivalent CLI diagnostics;
- unchanged rendered bytes when only tool metadata changes;
- no directory-source-path moves in the implementation slice.
