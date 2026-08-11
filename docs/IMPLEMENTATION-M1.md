# Mogent Milestone-One Contract

Status: completed milestone contract; URL-source clauses are superseded by
`docs/IMPLEMENTATION-M4.md`.

This document resolves the milestone-one choices left open by the design of
record. It is the implementation target for the first parser and renderer.
Source: DI-vukam.

## Scope

Milestone one provides a root Go module and one public command:

```text
mogent build
```

It reads a local `agents.yaml`, resolves local Markdown libraries, renders an
`AGENTS.md`, validates all input before writing, and writes the output safely.
The milestone did not fetch URL sources, provide a navigator, or import direct
edits into local overrides. Later contracts add those capabilities without
changing the local renderer rules recorded here.

## Manifest Forms

Both forms below are valid and normalize to the same internal document tree.

```yaml
doc:
  - Identity: cdint:shared-baseline/identity

  - heading: Testing
    from:
      - cdint:shared-baseline/instructions
      - ucd_research:research
    exclude:
      - cdint:shared-baseline/instructions/flaky-retries

  - heading: Instructions
    children:
      - Workflow: cdint:shared-baseline/instructions/focused-change-loop
      - Go: cdint:go
```

Rules:

- A compact entry is one manifest heading mapped to one `source:path` reference.
- An explicit entry has `heading` plus exactly one of `from` or `children`.
- `from` is an ordered list of one or more `source:path` references.
- `children` is an ordered list of compact or explicit entries.
- `from` and `children` cannot appear together.
- `exclude` is valid only with `from`.
- Each exclusion is source-qualified and names a descendant of one selected
  source subtree. Relative exclusions are deferred.
- Milestone one does not support authored body text in the manifest.

## Sources And References

```yaml
sources:
  cdint: ./libraries/cdint
  ucd_research: ~/agent-libraries/ucd_research
```

- Source aliases are user-defined and unique.
- Local source paths may be relative to `agents.yaml`, absolute, or begin with
  `~/`.
- A source reference is always `alias:heading/path`.
- References cannot escape their selected source root.
- HTTP(S) URL values require the immutable lock/cache and changed-content review
  behavior in `docs/IMPLEMENTATION-M4.md`.

## Rendering

The manifest owns output headings. The selected source root heading is removed;
its body is rendered below the manifest heading. Descendant headings are shifted
to remain below that manifest heading.

When `from` selects several sources, their bodies and descendant subtrees render
in listed order. If selected subtrees contain the same descendant heading path,
`mogent build` fails and names the collision. The later navigator will let the
user combine or separate these sections interactively.

An empty selected source node is an error. A missing template variable is an
error that names the variable.

## Validation

The YAML schema is strict. Reject:

- duplicate keys or source aliases;
- YAML anchors and aliases;
- unknown keys;
- wrong value types;
- empty `from` lists;
- simultaneous `from` and `children`;
- unresolved source aliases or heading paths;
- source references or exclusions outside their source root;
- empty rendered output.

## Output And Direct Edits

`output` is optional and defaults to `AGENTS.md` beside `agents.yaml`. An
explicit relative output path is also resolved beside `agents.yaml`.

Build writes to a temporary file and atomically replaces the output only after
successful parsing, resolution, templating, and validation.

After a successful build, Mogent writes an ignored local state file at
`.mogent/state.json`. It records the hash of the generated output. On a later
build:

- if the on-disk output still matches that hash, Mogent replaces it normally;
- if it differs, Mogent refuses to overwrite it and tells the user to inspect a
  three-way diff: prior generated output, current direct edits, and new output;
- `--force` explicitly replaces an untracked or modified output and records the
  new generated hash.

The next feature adds an interactive diff/import flow. It will let a user choose
a manifest node and create a copy-on-write local override from direct edits.

## Deferred From Milestone One

- URL source fetching, caching, pinning, and changed-content confirmation
  (implemented later by `docs/IMPLEMENTATION-M4.md`).
- Interactive source-collision choices.
- Relative exclusion syntax.
- Entries that combine `from` with `children`.
- Interactive direct-edit import into local overrides (a conservative core/CLI
  form is implemented by `docs/IMPLEMENTATION-M3.md`).
- Navigator, init, edit, swap, and save flow.
