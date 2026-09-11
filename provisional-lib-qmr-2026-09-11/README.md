# Provisional QMR Library

Status: draft review library, version 0.1.0.

This tree is the first dogfood library for QMR's own repository. Every module
is a manual-review candidate: its source text is preserved verbatim in a
blockquote and is not approved, selected, or rendered by Mogent yet.

## Layout

- `agents/` holds candidate always-loaded agent guidance.
- `guides/` holds human-readable development guidance.
- `skills/` holds triggered procedures using the Agent Skills shape.
- `specs/` holds candidate record or protocol truth.
- `intake/` holds unreviewed extraction work and review instructions.
- `status/` indexes every candidate and its review state.

The source-guide corpus remains in `../docs/other_repo_agents/`. The v2
wire-lab exemplar remains in `../libv2-proto/intake/exemplar/`. They are
evidence and intake tooling, not candidate modules in this provisional
library. PromiseGrid-specific material is deliberately deferred.

## Review Rule

Do not turn a review file into a selectable module by accident. To approve a
candidate, create a separate promoted module containing the approved prose and
frontmatter, retain its source link, and update `status/index.md`.

## Composition Rule

The current default is manifest composition of separate files. `[[...]]`
imports are neither approved nor implemented. If exact in-document insertion
is later required, choose and validate one of the deferred designs recorded in
`../libv2-proto/SPEC.md`; do not invent an include syntax here.

Split a promoted module whenever it needs different per-node metadata or YAML,
including status, tags, provenance, requirements, conflicts, or an exclusive
group. `library.yaml` owns cross-module groups, presets, and relationships.
