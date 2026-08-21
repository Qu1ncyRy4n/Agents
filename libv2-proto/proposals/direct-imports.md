# Keep Import Choices Beside Their Heading and Selection Visible in the Manifest

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

## Provisional Inline YAML Form

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

This is intentionally provisional. It establishes several required semantics:

- an unanswered `exactly_one` choice is an error;
- `none` is valid only when explicitly listed;
- `any` permits zero, one, or several listed modules;
- imports never search directories or choose the first match;
- the manifest visibly identifies the selected source;
- cycle detection applies if imported content contains slots;
- insertion does not weaken normal `requires` or conflict validation; and
- rendered provenance identifies both the containing module and inserted node
  in tool output, even when metadata is omitted from `AGENTS.md`.

The parser must also reject malformed YAML and duplicate keys, enforce source
root and symlink safety, remain offline during ordinary rendering, and detect
direct or indirect import cycles. These are requirements for a later design
slice, not implemented syntax.

## Simpler Existing Alternative

The manifest can already assemble a parent and chosen child as adjacent outline
entries. A Markdown slot earns new syntax only when preserving an exact
in-document insertion point is a demonstrated authoring need.

## Decision Gate

Run a source-grounded thought experiment after extraction supplies several real
choice groups. Compare inline YAML, manifest-only composition, named import
points, and a compact Boolean expression using the same persona, role, strict
governance, and commit-policy examples. Do not implement this proposal before
that decision.

