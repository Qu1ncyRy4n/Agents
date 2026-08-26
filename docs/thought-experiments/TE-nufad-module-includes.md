# TE-nufad: Narrow the Module Include and Composition Surface

TE ID: TE-nufad
Date: 2026-08-25
Status: complete analysis; needs owner decisions

`tools/mint-handle` is unavailable, so `nufad` is manually assigned under the
repository's existing user-approved exception.

## Decision Under Test

How should an authored Markdown module incorporate another selectable module,
especially when the included content must appear at a particular point inside
a larger compatible document?

This thought experiment starts with the broad design space. Its job is to reject
or defer weak approaches and reduce the field to a few decision-ready options.
It does not choose the final syntax on the owner's behalf.

## Distinguish Three Jobs Before Choosing Syntax

“Include” currently hides three different jobs:

1. **Document composition:** select several modules and order their headings in
   the output. The manifest already does this.
2. **Exact insertion:** place one module's content at a particular location
   inside another authored document.
3. **Template evaluation:** substitute variables, choose conditional text, or
   repeat data while rendering one module.

One mechanism does not necessarily need to perform all three jobs. A small
composition model is safer if exact insertion and template evaluation are rare.

## Required Properties

Any surviving design should preserve:

- ordinary Markdown readability without Mogent;
- one obvious place to inspect and update a module's structural metadata;
- explicit source provenance and deterministic rendering;
- module references rather than unrestricted filesystem paths;
- local-source root and symlink safety;
- offline ordinary builds;
- missing-target and ambiguity errors;
- direct and indirect cycle detection;
- clear heading behavior and insertion order;
- compatibility with pinned source updates and renamed modules;
- useful diagnostics with source module and include chain; and
- a simple happy path for modules that require no special composition.

## Broad Alternatives

### A. Manifest Composition Only

Split independently selectable content into files and compose it through
`agents.yaml`. Use frontmatter only for facts intrinsic to each module.

```text
communication/
  baseline.md
  personas/
    surfer.md
    alien.md
  roles/
    librarian.md
    code-assistant.md
```

```yaml
doc:
  - Communication:
      - Baseline: personal:communication/baseline
      - Persona: personal:communication/personas/surfer
      - Roles:
          - Librarian: personal:communication/roles/librarian
```

This requires no Markdown include syntax. It makes selection, order, and
provenance visible in the manifest.

Tradeoff: even a small independently selected subsection becomes a separate
file. It cannot be injected between two paragraphs of one source node unless
the containing document is split at that boundary or the manifest gains a more
explicit composition shape.

### B. Frontmatter-Declared Includes

Keep all module YAML in the file's frontmatter:

```yaml
---
includes:
  - self:communication/personas/surfer
---
# Communication
```

This provides one constrained metadata location and uses the existing YAML
parser. It works when includes have a fixed placement rule, such as before or
after the containing body.

Tradeoff: frontmatter alone cannot identify an arbitrary insertion point inside
the body. Adding placement selectors that refer to headings or markers can
recreate a distant two-surface editing problem.

### C. Frontmatter Plus a Minimal Body Marker

Declare allowed or selected imports in frontmatter, then use a small body marker
only for placement:

```yaml
---
includes:
  persona: self:communication/personas/surfer
---
# Communication
```

```markdown
## Adopt a Communication Style

.include persona
```

The body marker contains no embedded YAML or choice logic. Frontmatter remains
the single structural metadata location.

Tradeoff: `.include` is invented syntax, is visible to ordinary Markdown
readers, and needs escaping rules. The exact marker could instead be a narrowly
defined HTML comment or template action; the important distinction is that it
names a frontmatter declaration rather than embedding a second configuration
document.

### D. Direct Minimal Include Directive

Put a single source-module reference at the insertion point:

```markdown
## Adopt a Communication Style

.include self:communication/personas/surfer
```

This is concise and easy to parse. It avoids repeating the reference in
frontmatter.

Tradeoff: module structure now has two possible homes—frontmatter relationships
and body includes. Adding conditions or choice expressions to `.include` would
quickly create a new language.

### E. Go Template Include Function

Extend the existing Go-template rendering model:

```gotemplate
{{ include "self:communication/personas/surfer" }}
```

Or pass a manifest-selected reference through the template data map:

```gotemplate
{{ include .communication_persona }}
```

This reuses familiar delimiters and can support variables and conditionals.
Mogent would need to preload or safely resolve included templates, restrict the
function surface, track provenance, and detect cycles.

Tradeoff: Go templates are more powerful and less Markdown-native than a single
include directive. Dynamic include names make static validation and source
browsing harder. Template errors can become more obscure than schema errors.

### F. External YAML Paired with Markdown

Keep Markdown free of structural configuration and define composition in a
sidecar or source-level YAML file:

```yaml
modules:
  communication:
    file: communication.md
    insertions:
      communication-persona:
        after_heading: Adopt a Communication Style
```

This centralizes machine configuration and can validate a source library as a
whole.

Tradeoff: editing one module requires updating a distant file. Heading-based
placement can break on prose renames unless another stable marker is introduced.
The sidecar can become a registry duplicating the physical tree.

### G. Embedded YAML in an HTML Comment

```markdown
<!-- yaml:mogent:start
imports:
  communication/persona:
    choose: exactly_one
    from: [none, self:personas/alien, self:personas/surfer]
yaml:mogent:end -->
```

This keeps rich configuration near the affected heading and remains invisible
in normal rendered Markdown.

Tradeoff: it creates YAML inside Markdown inside HTML comments, adds another
parser surface, permits schema growth anywhere in the document, complicates
diagnostics and formatting, and leaves unclear how heading scope is determined.
The Steve meeting independently raised the same maintainability concern.

### H. Obsidian or Wiki-Style Transclusion

```markdown
![[communication/personas/surfer]]
```

This is compact and familiar in Obsidian.

Tradeoff: GitHub-Flavored Markdown does not define wiki-link transclusion.
Paths, headings, aliases, and embed behavior would follow an external product's
semantics unless Mogent invented a similar but subtly different dialect.

### I. C-Preprocessor-Style Include

```text
#include "communication/personas/surfer.md"
```

This is familiar from programming languages and naturally suggests conditions.

Tradeoff: it is filesystem-oriented, not source-module-oriented; it invites a
larger preprocessor language; and it fits Markdown authoring less naturally than
Mogent's existing manifest and Go-template model.

## Scenario Analysis

### Ordinary Compatible Bundle

A baseline communication document and its always-required subsections should
remain one Markdown file. Every include design is unnecessary overhead here.
This favors no include syntax on the happy path.

### Small Optional Subsection

An optional two-paragraph TTS section illustrates the owner's concern: splitting
tiny content into a file adds navigation and maintenance cost. Manifest-only
composition handles selection cleanly but may make the source tree feel
fragmented. A minimal marker handles exact placement but introduces syntax.

This is the strongest real case for preserving a narrowly scoped include
primitive after module extraction supplies examples.

### Mutually Exclusive Persona and Compatible Roles

Choice rules should not be encoded in an include path. Frontmatter, a source
descriptor, or the consuming manifest can declare the allowed set and selection
cardinality. The insertion marker should identify only what content belongs at
that position.

This rejects compact Boolean path expressions as the default authoring model.

### Exact Mid-Document Placement

Frontmatter-only includes cannot express “insert here” without a placement
convention. Manifest composition can express the same output if the containing
document is split at the insertion boundary. Minimal markers and Go-template
includes express exact placement directly.

The product decision therefore depends on whether exact mid-document insertion
is a required early job or a convenience that can wait.

### Variables and Conditional Text

Mogent already uses Go templates for manifest variables. If include conditions
depend on those values, a restricted Go-template function has a coherent model.
If choices are manifest selections rather than value-dependent text, ordinary
manifest composition remains more visible and auditable.

Do not adopt a general template engine merely to solve static module selection.

### Static Validation and Browsing

Manifest references, frontmatter references, and literal minimal directives can
be indexed without executing templates. Dynamic Go-template includes cannot
always be enumerated in advance. Embedded YAML can be indexed but requires a
new scoped-parser model. External sidecars are indexable but may drift from the
Markdown they describe.

### Nested Includes and Cycles

Every include-capable survivor requires a directed graph over canonical source
module references. Validation must reject a cycle with its complete chain and
enforce a maximum nesting depth. Included files must not expand the source root,
fetch undeclared libraries, or bypass pinned-cache verification.

### GitHub and Obsidian Reading

Plain Markdown plus frontmatter reads acceptably in both. Go-template actions
and `.include` directives remain visible as authoring syntax. HTML comments are
hidden but difficult to inspect in rendered views. Obsidian transclusion works
best in Obsidian but is not portable to GitHub.

### Rename and Migration

All survivors should reference canonical source nodes rather than raw file
paths. A source-path move map can update or warn on literal references.
Heading-text selectors and unrestricted filesystem includes are substantially
more fragile.

## Narrowing Results

### Reject as the Default

- **Embedded YAML in HTML comments:** too open and parser-heavy for the default.
- **Obsidian transclusion:** insufficiently portable as Mogent's canonical
  syntax.
- **C-preprocessor syntax:** introduces the wrong filesystem/preprocessor model.
- **Boolean expressions inside include paths:** obscure selection meaning and
  diagnostics.

These may remain comparison references; they should not lead implementation.

### Defer Unless Evidence Requires It

- **External YAML placement registry:** useful for library-wide groups or
  validation, but too distant for ordinary module insertion.
- **Dynamic Go-template includes and loops:** powerful, but static selection does
  not yet justify the complexity.

### Surviving Directions

1. **Manifest composition plus atomic modules as the default.** Keep intrinsic
   metadata in frontmatter. Split content only when it genuinely needs an
   independent selection boundary.
2. **Frontmatter plus a minimal named insertion marker** if exact placement is a
   required near-term job. Keep choice rules out of the marker.
3. **A restricted literal Go-template include function** if Mogent wants one
   template surface for exact insertion and existing variable substitution.
   Dynamic include names, arbitrary functions, loops, and conditions need not
   be enabled initially.

Direct `.include self:path` remains a syntactic variant of directions 2 or 3,
not a separate architecture decision.

## Tradeoffs for Owner Decision

### Direction 1: Manifest Composition Only

Best constraint, easiest validation, and one visible composition source. Costs
more files and cannot preserve arbitrary insertion inside a large source
document without splitting it.

### Direction 2: Frontmatter Plus Minimal Marker

Keeps YAML in one place and exact insertion readable. Adds a small Mogent
directive and requires the marker name to stay synchronized with frontmatter.

### Direction 3: Restricted Go-Template Include

Reuses an existing rendering language and avoids a separate marker grammar.
Exposes more template complexity and requires careful restrictions for static
validation, provenance, and safety.

## Decisions Needed

1. Is exact mid-document insertion required for the first libv2 cutover, or may
   v2 begin with manifest composition and split modules?
2. If exact insertion is required, should the next prototype compare a minimal
   frontmatter-named marker against a restricted literal Go-template include?
3. Should choice-group rules live in module frontmatter when intrinsic, in a
   source-level descriptor when spanning modules, and in the manifest when the
   relationship exists only for one consuming repository?
4. Is GitHub-Flavored Markdown plus YAML frontmatter the canonical portable
   authoring format, with Obsidian features treated as optional views rather
   than library syntax?

## Follow-Up Boundaries

This TE does not settle:

- merge precedence between inherited YAML maps;
- requirement and conflict validation;
- source-path move compatibility;
- multiple output formats;
- general-purpose commands embedded in Markdown; or
- whether section TLDRs need heading-local metadata.

Those should remain separate decisions so include syntax does not become a
container for every unresolved feature.

