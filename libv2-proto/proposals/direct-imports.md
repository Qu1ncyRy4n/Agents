# Make Authored Imports Explicit and Keep Selection Logic in the Manifest

Status: syntax proposal, not implemented.

## Desired Job

A large authored module may have a stable insertion point for one selected
member of a family. The source document should be able to name that slot while
the consuming manifest chooses the concrete module. Mogent should validate the
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

## Provisional Explicit Form

The Markdown document names a slot:

```markdown
## Choose One Communication Persona
<!-- mogent-slot: communication-persona -->
```

The manifest fills it explicitly:

```yaml
slots:
  communication-persona:
    from: personal:communication/personas/surfer
```

The source library or manifest declares the constraint separately:

```yaml
slot_rules:
  communication-persona:
    required: true
    exclusive_group: communication/persona
```

This is intentionally provisional. It establishes several required semantics:

- a missing required slot is an error;
- two values for a single slot are an error;
- imports never search directories or choose the first match;
- the manifest visibly identifies the selected source;
- cycle detection applies if imported content contains slots;
- insertion does not weaken normal `requires` or conflict validation; and
- rendered provenance identifies both the containing module and inserted node
  in tool output, even when metadata is omitted from `AGENTS.md`.

## Simpler Existing Alternative

The manifest can already assemble a parent and chosen child as adjacent outline
entries. A Markdown slot earns new syntax only when preserving an exact
in-document insertion point is a demonstrated authoring need.

