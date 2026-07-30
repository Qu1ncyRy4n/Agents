# DR: Library Source Boundaries

DR-ID: DR-garom
Date: 2026-07-29 22:30:00
State: open
Asked by: 95124070+Qu1ncyRy4n@users.noreply.github.com (Quincy Ryan)

## Question

Should the first shared library sources be organized around `cdint`,
`ucd_research`, `personal`, and independently selectable `nix`, or should some
of those become narrower source libraries as the content grows?

## Why This Is Open

The proposed grouping gives useful ownership and trust boundaries, but real
composition may show that a module belongs in more than one group. For example,
generic Go guidance may serve CDINT and personal projects; research safety may
serve more than one lab; Nix may be shared across personal and research work.

## Current Direction

- Keep `cdint`, `ucd_research`, and `personal` as provisional source-library
  groupings.
- Keep `nix` as its own independently selected source because it can affect an
  active machine.
- Place `uv` under Python dependency management.
- Revisit the grouping after the first real manifests and source libraries have
  been used.

## Affects

- `docs/codex_eco/README.md`
- Future library layout and `agents.yaml` source aliases
- `jusuk.11` milestone-one dogfood manifest

## Unblocks

The immediate schema and renderer work. This question does not block milestone
one because that milestone can use a small local fixture library.

## Waiting On

Real use of the first manifests and source libraries.

## Linked DI

None yet. A Decision Intent will supersede this request when the source
boundaries are locked.

## Related Work

- `docs/codex_eco/README.md`
- `docs/thought-experiments/TE-vorum-manifest-schema-milestone-one.md`
