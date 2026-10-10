# DR-juvih: Typed Content Views And Explicit Composition

DR-ID: DR-juvih
Date: 2026-10-09
State: accepted design direction; grammar and delivery sequencing remain open
Decision intent: DI-juvih in `TODO/TODO.md`
Owner: Quincy Ryan
Basis: TE-gunak, including its twelve-case tabletop execution and subsequent
owner approval of the recommendation.

## Decision

Adopt a common typed content inventory with physical directory/file structure,
optional document heading views, and explicit copy/render operations. Begin with
intact payload copying and explicit rendering. Require explicit replacement for
destination collisions rather than treating source order as implicit permission
to overwrite content.

This is the approved next design direction, not a claim that the current HCL
parser or build engine implements it. Candidate names and HCL examples in the
TEs remain illustrative until the grammar is specified.

## Content Model

| Concept | Meaning |
|---|---|
| Directory | Physical container of files and directories |
| File | Original byte-bearing payload, with path and permission information |
| Heading view | Optional document structure associated with a file or explicitly authored contributions |
| Role | Instructions, skill, research, guide, configuration, or tool-native profile; separate from structural kind |
| Bundle | A coherent set of payloads, such as a skill's SKILL.md, scripts, references, and assets |

One physical Markdown file can have a heading subtree without losing its
whole-file identity. Document-view relationships do not manufacture filesystem
directories, duplicate the file, or give a heading ownership of arbitrary assets.
Logical addresses, physical locators, display titles, and output destinations
are distinct. Source/context qualification must distinguish otherwise identical
names from different files or libraries.

## Skill Example

```text
review/                           directory; skill bundle
├── SKILL.md                       file; original bytes/frontmatter
│   └── document view
│       └── Review                 heading
│           └── Procedure          heading
│               └── Verification   heading
├── scripts/check.sh               file
└── references/checklist.md        file
```

Three different requests have different meanings:

1. **Install review:** copy the complete selected skill bundle intact.
2. **Reuse Procedure in AGENTS.md:** explicitly render that contribution into a
   separate document. Do not implicitly edit SKILL.md or discard its other
   headings/assets.
3. **Customize Procedure inside the skill:** derive a modified skill while
   retaining metadata, remaining content, dependencies, and provenance. This
   editing operation is backlogged, not an implicit side effect of selection.

Parsed imported-heading support is part of the intended model, but its parser,
addressing contract, and delivery order relative to the existing authored-fragment
adapter remain open. Until that adapter exists, installed skills remain intact
copies and their headings are not selectable by current Mogent.

## Composition And Ownership

- Output target shape and per-contribution operation are independent. A managed
  directory may contain copied bundles, rendered Markdown files, and opaque
  configuration payloads.
- Source declaration order remains visible. Order determines contribution order;
  it does not authorize an unrequested replacement.
- Unresolved destination collisions block composition. Explicit whole-file or
  bundle replacement identifies what is being replaced and the winning source.
  Report collisions and replacements even when the final bytes are identical.
- Directory union is not whole-bundle replacement: replacing a bundle removes
  obsolete lower members, whereas a union may retain them. Do not silently
  produce a partial or mixed-origin installed skill.
- One final manifest owns each destination once. Mixed artifacts retain current
  generated-state tracking, direct-edit protection, containment, and rollback.
- Plan must explain source/context, selected node/view, operation, placement,
  replacement, and removed bundle members. Provenance remains outside preserved
  copied payloads.

## Preservation And Scope Constraints

- Copy complete file bytes and preserve frontmatter and bundle-relative layout.
  Rendering is an explicit transformation, not an automatic interpretation of
  every file ending in `.md`.
- TOML/YAML/JSON/JSONC remain whole-file opaque imports/replacements. Defer
  semantic configuration merging and tool-schema conversion.
- Preserve current source-filtering `exclude`. Removing earlier contributions
  must have an explicit scope; its exact syntax is still open.
- Heading extraction can carry file/link dependencies. Specify how they are
  retained, placed, or diagnosed before promising safe extraction. No automatic
  asset discovery or link rewriting is implied by a heading selection.
- Origin scope, build order, destination scope, and agent runtime loading remain
  separate. No invisible global-content merge or implicit home-directory write.
- Unknown `.agents` files may be copied without claiming universal auto-loading.
  The wider agentsstandard.com cascade remains a versioned proposal rather than
  an inherited law of the composition engine.

## Backlog

Defer vendor rules, custom-agent/profile conversions, command/prompt adapters,
hook/settings registration, automatic global installation/cascade projections,
and editable heading derivatives. Complete native files can still be explicitly
copied; backlogging a semantic adapter does not prohibit raw file import.

## Still Open Before Implementation

1. Root/node names and HCL grammar: `content`, `node`, typed blocks, and output
   `dir-tree` remain candidates, not accepted spellings.
2. Parsed heading addresses, ambiguity handling, dependencies, and delivery order.
3. Precise directory union, file/subtree replacement, and explicit replacement
   syntax under the accepted collision policy.
4. Scope syntax for source filters versus earlier-contribution exclusions.
5. Migration, acceptance fixtures, and the staged implementation sequence.

These detailed specifications refine the accepted model; they do not reopen
automatic later-wins as the default or reintroduce semantic config merging.

## References

- [TE-gunak and execution record](../docs/thought-experiments/TE-gunak-typed-agent-content-composition.md)
- [Earlier TE-jofit](../docs/thought-experiments/TE-jofit-unified-content-composition.md)
- [Ecosystem report](../reports/Universal%20agent%20content%20composition.md)
- [Composition roadmap](../TODO/TODO-tugur-v2-composition-followup.md)
- [Current v2 specification](../docs/proposals/MOGENT-HCL-V2-SPEC.md)
