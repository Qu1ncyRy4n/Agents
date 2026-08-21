# TE-mavok: Library Module Shape And Relationship Surface

TE ID: TE-mavok
Date: 2026-08-20
Status: open thought experiment

`tools/mint-handle` is unavailable, so `mavok` is manually assigned under the
existing user-approved exception.

## Question

How should a library organize compatible bundles, optional modules, and
mutually exclusive choices without spreading the same relationship model over
Markdown comments, paths, and several YAML files?

## Worked Library

Suppose a teaching-agent library contains styling tools, personas, and roles.
Some styling tools can be combined. A persona is a voice choice. A librarian
and teacher can be compatible roles, while `senior-dev` and
`coding-assistant` may be redundant rather than logically incompatible.

The recommended physical shape is:

```text
teaching/
  library.yaml
  styling.md                 # one compatible baseline bundle
  styling/
    tts-chat.md              # optional
    tts-output-stream.md     # optional
    whiteboard.md            # optional
  personas/
    surfer.md                # at most one persona
    robot.md
    alien.md
    caveman.md
  roles/
    librarian.md             # compatible with teacher
    teacher.md
    senior-dev.md
    coding-assistant.md
    reading-coach.md
```

The directory is already the browsing structure, so a comment such as
`dir: ./personas` would duplicate a fact that can drift after a move.

## Recommended Metadata Split

Use file frontmatter for facts intrinsic to one selectable node:

```yaml
---
tags: [communication/persona, discovery/voice]
tldr: Speak like a relaxed surfer while keeping instructions precise.
priority: 0.5
exclusive_group: communication/persona
requires: [self:styling]
conflicts_with: []
see_also: [self:roles/reading-coach]
---
# Surfer
```

Use one optional source-level `library.yaml` for facts spanning nodes:

```yaml
groups:
  communication/persona:
    mode: at_most_one
    members:
      - self:personas/surfer
      - self:personas/robot
      - self:personas/alien
      - self:personas/caveman

presets:
  reading-tutor:
    include:
      - self:styling
      - self:roles/teacher
      - self:roles/reading-coach
    requires_one_of:
      - self:personas/surfer
      - self:personas/robot
      - self:personas/alien
      - self:personas/caveman
```

This example deliberately separates two meanings:

- `exclusive_group` means selecting more than one member is invalid.
- `requires_one_of` means the preset is incomplete until one member is chosen.

They should not be collapsed into an "exclusive tree." A tree describes
containment; exclusivity describes a selection constraint.

## Worked Relationship Examples

### Hard Requirement

```yaml
---
requires: [self:safety/data-handling]
---
# Human-Subject Analysis
```

Manifest before:

```yaml
doc:
  - Research: research:human-subject-analysis
```

Result:

```text
error: research:human-subject-analysis requires
       research:safety/data-handling
hint:  mogent add research:safety/data-handling --before Research --dry-run
```

Manifest after review:

```yaml
doc:
  - Safety: research:safety/data-handling
  - Research: research:human-subject-analysis
```

Both contents render. `requires` is not a hyperlink and does not mean "choose
whichever appears first." If the target does not render, the generated agent
cannot obey it.

### Exact Conflict

```yaml
---
conflicts_with: [self:python/dependency-management/poetry]
---
# uv
```

Selecting both `uv` and `poetry` is an error. Their order does not resolve the
conflict.

### Exclusive Group

```yaml
---
exclusive_group: communication/persona
---
# Robot
```

If `Surfer` is already active, adding `Robot` produces a diagnostic naming both
members and the group. The user removes one; Mogent does not pick the first.

### Compatible Options

```yaml
# tts-chat.md
---
tags: [interface/tts, output/chat]
---
# Chat Speech
```

```yaml
# whiteboard.md
---
tags: [interface/visual, collaboration/whiteboard]
---
# Whiteboard
```

No relationship is needed. Both may be selected. Compatibility is the default,
not another metadata field.

### Discovery Link

```yaml
---
see_also: [self:roles/reading-coach]
---
# Teacher
```

`source show` may display the related node, but it is not automatically selected
or rendered.

### Canonical Node And Discovery Tags

```yaml
---
tags: [lang/python, dependency/management, tool/uv]
---
# uv
```

The canonical identity remains `self:python/dependency-management/uv`. Tags
support search; they do not create aliases or alternate identities.

## Heading Comments

Heading-local metadata would be useful only when one file intentionally holds
several independently selectable headings. A possible future form is:

```markdown
## Teacher
<!-- mogent:
requires: [self:roles/reading-coach]
-->
```

It should not be the first relationship surface. Multiline YAML hidden in HTML
comments needs another parser and makes relationships harder to audit. First
dogfood atomic files plus frontmatter. Add heading-local metadata only when a
real compatible large document cannot be split cleanly.

Short presentational metadata such as a subsection TLDR remains a different,
smaller use of an HTML comment:

```markdown
## Validation
<!-- tldr: Run the focused check, then broaden it when risk warrants. -->
```

## Logical Operators

Do not add `||`, `&&`, or a mini expression language initially. Use named,
typed fields:

```yaml
requires: [self:a, self:b]       # all
requires_one_of: [self:c, self:d] # one or more, only where supported
```

This is less flexible but produces clearer validation and error messages. Add
richer expressions only after real libraries contain a relationship that these
two forms cannot represent.

## Options Considered

1. Put everything in heading comments. Fine-grained, but parser-heavy and hard
   to audit.
2. Put everything in one `library.yaml`. Centralized, but node edits require a
   second distant edit and the file becomes a registry.
3. Put all relationships in file frontmatter. Simple for atomic nodes, but
   awkward for groups and presets spanning several nodes.
4. Use the hybrid above: node facts in frontmatter; spanning facts in
   `library.yaml`. This has two surfaces, but each fact has one clear owner.

## Recommendation

Adopt option 4 provisionally. Keep compatible guidance together when it is
always selected as one unit. Split optional or mutually exclusive guidance into
atomic files. Dogfood that structure during the upcoming library reorganization
before adding heading-level relationship metadata.

