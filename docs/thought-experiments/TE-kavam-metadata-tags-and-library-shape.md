# TE-kavam: Metadata Tags And Library Shape

## Status

Open thought experiment.

## Question

How should mogent use source metadata without pretending it understands more
than it does?

This covers three related choices:

- whether tags should be flat strings or hierarchical paths,
- whether `conflicts_with` can produce meaningful warnings across multiple
  libraries,
- whether libraries should prefer one atomic module per file or complete
  `AGENTS.md`-style documents.

## Context

Mogent now reads file-level YAML frontmatter from source Markdown files. The
first supported fields are `tags`, `tldr`, `priority`, `scope`, `requires`, and
`conflicts_with`. Metadata is tool-only and does not render into `AGENTS.md`.

The first implementation uses exact tag matching. That is deliberately simple,
but it raises a real design question: a tag such as `go` may be too vague once
libraries grow, while a tag such as `lang/go/testing` can be browsed, filtered,
and grouped more naturally.

## Option A: Flat Tags

Example:

```yaml
tags: [go, testing, security]
```

This is easy to type and easy to filter. It is also easy to misuse. Two libraries
can both use `security` while meaning very different things: one might mean
secure coding practices, another might mean prompt-injection hardening, and a
third might mean deployment secrets.

Flat tags are fine for small personal libraries. They age poorly for shared org
libraries.

## Option B: Hierarchical Tags

Example:

```yaml
tags: [lang/go, testing/unit, scope/org, risk/security]
```

This keeps tags readable while adding enough shape for better browsing:

- `lang/go` means Go-specific content.
- `testing/unit` means unit-test behavior.
- `scope/org` means broadly applicable organization guidance.
- `risk/security` means the node affects a safety or security boundary.

The key advantage is that hierarchy supports both exact and prefix-style
questions later:

- exact: "show `lang/go` modules"
- prefix: "show everything under `testing/`"

The risk is false precision. If tags become too elaborate, authors spend more
time taxonomy-gardening than writing useful instructions.

## Option C: Compatibility Families

Some modules are alternatives. For example:

```text
security/loose
security/tight
```

Those sibling modules may be meaningful alternatives because they live in the
same family and were written with each other in mind. In that case, mogent could
warn when both are selected.

This works best when conflict meaning is local to a known family:

```yaml
tags: [policy/security]
variant: tight
```

or:

```yaml
family: security
variant: tight
```

This does not work well as a universal, cross-source claim. A file in one source
cannot reliably know whether a file from another source truly conflicts with it.
It can only say "I conflict with this known reference" or "I belong to this
alternative family."

## Conflict Warning Thought Experiment

Imagine a manifest that pulls these two nodes:

```yaml
sources:
  org: ./org-library
  personal: ./my-library

doc:
  - Security: org:security/tight
  - Workflow: personal:workflow/fast-iteration
```

If `personal:workflow/fast-iteration` says:

```yaml
conflicts_with: [org:security/tight]
```

then a warning is meaningful because it names the exact source and path.

But if it says:

```yaml
conflicts_with: [security/tight]
```

that is ambiguous. Which source owns `security/tight`? Is it a path, a tag, or a
semantic idea? Mogent should not guess.

A safer future rule:

- exact source references can produce warnings,
- declared families can produce warnings inside the same family,
- loose semantic conflicts remain human-visible metadata only.

## Library Shape Thought Experiment

Atomic module files may be the best authoring shape for shared libraries:

```text
library/
  instructions/
    workflow.md
    testing.md
  security/
    loose.md
    tight.md
  lang/
    go/
      testing.md
      errors.md
```

Each file can hold:

- one top-level heading,
- frontmatter metadata,
- the content for one reusable module.

That makes browsing and selection easier. It also makes source provenance easier:
file path, heading path, tags, TLDR, and priority all describe one thing.

Complete `AGENTS.md`-style documents are still useful as imports, examples, and
human-authored starting points. They are less ideal as shared libraries because
one file may contain many unrelated modules but only one file-level metadata
block.

This suggests two complementary workflows:

- split: take a large `AGENTS.md` and produce an atomic library directory,
- join/render: take a manifest and render a complete `AGENTS.md`.

The renderer does not need to require atomic files. The library authoring tools
can strongly recommend them.

## Current Lean

Prefer hierarchical tags, using slash paths such as `lang/go`, `scope/org`, and
`risk/security`.

Keep `priority` strict: values outside `0.0` through `1.0` are errors.

Keep `conflicts_with` as visible metadata for now. Later, warnings should be
limited to exact source references and declared alternative families. Avoid
global semantic conflict guessing.

Prefer atomic module files for shared libraries, organized by directory. Add a
future split/import command that can turn a complete `AGENTS.md` into a draft
library directory for review.
