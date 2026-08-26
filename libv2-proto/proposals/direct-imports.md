# Compare Minimal Includes Without Embedding Configuration in Markdown

Status: syntax proposal, not implemented.

## Desired Job

A large authored module may have a stable insertion point for one selected
member of a choice group. The source document should define that import point
and its allowed choices while the consuming manifest records the selected
module. Mogent should validate the
choice before rendering.

## Avoid a Logic Expression Inside the Path

This is compact but combines importing, Boolean expressions, globbing, and
selection in an opaque string:

```markdown
[[import:./personas/{(surfer^caveman^alien)&(codeassistant|librarian|senior_dev)}]]
```

It is difficult to explain whether `^`, `&`, and `|` mean exclusive choice,
conjunction, fallback order, or text concatenation. It also hides the selected
source from an ordinary manifest review.

## Previous Inline YAML Form

The first v2 sketch placed a rich YAML block directly beneath the affected
heading. `TE-nufad` now rejects this as the default because it creates a broad
second YAML surface inside Markdown and makes parsing, scope, formatting, and
diagnostics harder.

Keep the import definition immediately beneath the heading it affects:

```markdown
### Adopt a Communication Style

<!-- yaml:mogent:start
imports:
  communication/persona:
    choose: exactly_one
    from:
      - none
      - self:personas/alien
      - self:personas/surfer

  communication/roles:
    choose: any
    from:
      - self:roles/librarian
      - self:roles/code-assistant
      - self:roles/senior-developer
yaml:mogent:end -->
```

The manifest records selections rather than repeating allowed values:

```yaml
imports:
  communication/persona: personal:personas/surfer
  communication/roles:
    - personal:roles/librarian
    - personal:roles/code-assistant
```

This historical sketch established several required semantics:

- an unanswered `exactly_one` choice is an error;
- `none` is valid only when explicitly listed;
- `any` permits zero, one, or several listed modules;
- imports never search directories or choose the first match;
- the manifest visibly identifies the selected source;
- cycle detection applies if imported content contains slots;
- insertion does not weaken normal `requires` or conflict validation; and
- rendered provenance identifies both the containing module and inserted node
  in tool output, even when metadata is omitted from `AGENTS.md`.

Do not implement this syntax without a later owner decision reversing the
`TE-nufad` narrowing result.

## Current Surviving Directions

1. Use manifest composition and separate modules as the default.
2. If exact mid-document insertion is required, keep YAML in frontmatter and
   place only a minimal named marker in the body.
3. Alternatively, prototype a restricted literal Go-template include so Mogent
   reuses its existing template surface.

Choice rules remain outside the body marker. See
`docs/thought-experiments/TE-nufad-module-includes.md` for scenario analysis and
the decisions still needed.

## Simpler Existing Alternative

The manifest can already assemble a parent and chosen child as adjacent outline
entries. A Markdown slot earns new syntax only when preserving an exact
in-document insertion point is a demonstrated authoring need.

## Decision Gate

`TE-nufad` completed the first broad narrowing pass. Do not implement includes
until the owner resolves whether exact insertion is required for the first
libv2 cutover and chooses which surviving directions deserve a prototype.
