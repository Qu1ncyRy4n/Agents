# Mogent Library And Policy Decision Map

Status: discussion draft; no decisions recorded.

This map separates choices that have become tangled during library review. Work
from the top down: later command behavior depends on the earlier content model.

```text
1. Library shape
   What is one selectable module?
   |
   +-- 2. Metadata surface
   |      Where do requires, conflicts, groups, and discovery tags live?
   |      |
   |      +-- 3. Composition semantics
   |             What must render together, what may coexist, and what is
   |             mutually exclusive?
   |             |
   |             +-- 4. Organization policy
   |                    Who may require a bundle for a repository, and how is
   |                    that authority made visible?
   |                    |
   |                    +-- 5. Update and reconciliation
   |                    |      How are upstream changes previewed, accepted,
   |                    |      preserved locally, or proposed back?
   |                    |
   |                    +-- 6. Proposal delivery
   |                           Patch -> Git branch/commit -> forge pull request
   |
   +-- 7. Later information features
          See-also links, rationale/evidence links, richer provenance, and
          cross-library relationships.
```

## Recommended Decision Order

| Order | Decision | Smallest question to answer now | Why it comes here |
|---|---|---|---|
| 1 | Module boundary | When should one Markdown file be one node versus a compatible bundle? | Every relationship needs stable node identity. |
| 2 | Metadata location | File frontmatter, heading comments, directory YAML, or a hybrid? | Parsers and authoring UX depend on this. |
| 3 | Relationship meanings | What do `requires`, `conflicts_with`, and `exclusive_group` do? | Policy and updates must know what a valid selection is. |
| 4 | Organization requirement | Does the repo manifest visibly opt into a required organization bundle, or can an external policy layer impose it? | This is a trust and ownership boundary. |
| 5 | Reconciliation surface | Keep only `status`, `preserve`, `inherit`, and `propose` public initially? | A small directional API is easier to dogfood. |
| 6 | Required-update review | Does `required` affect selection only, or also update acceptance? | Necessity is not the same as trust or urgency. |
| 7 | Forge delivery | Is a PR-ready branch enough initially, or should Mogent open the PR? | External publication needs separate authority. |

## Current Lean

- Use physical directories for browsing and taxonomy rather than repeating
  `dir:` paths in Markdown metadata.
- Keep a compatible bundle together; split optional or mutually exclusive
  choices into their own files.
- Keep file frontmatter for node metadata. Use a source-level `library.yaml`
  only for source-wide groups, presets, and policy declarations.
- Keep `self:` for exact relationships inside one source. Defer cross-library
  relationships.
- Define `exclusive_group` as "at most one." Express "one is required" in a
  preset or policy separately.
- Make `requires` a validation rule: the target's content must be selected and
  rendered. Never use first-match or source order to choose or suppress it.
- Defer `urgent`, informal provenance, and generalized logical expressions.
- Start reconciliation with four public verbs: `status`, `preserve`, `inherit`,
  and `propose`. Three-way `reconcile` can be an internal operation until a
  distinct public job appears.

## First Owner Question

Approve or revise this module convention before choosing syntax:

> A file may contain a larger compatible bundle. Content that is optional,
> independently selectable, or mutually exclusive with a sibling becomes a
> separate one-top-level-heading file.

The alternatives and worked examples are in
`docs/thought-experiments/TE-mavok-library-module-shape.md`.

