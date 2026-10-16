# Small-Language Walkthrough

Status: discussion walkthrough, not executable configuration. Consolidated proposed
contract: [first-slice spec](MOGENT-FIRST-SLICE-SPEC.md), with
[acceptance cases](MOGENT-FIRST-SLICE-ACCEPTANCE.md). Later sections below preserve
the evolution of the examples; the consolidated draft identifies review choices.
Code fences use HCL highlighting for readability.

## 1. Inputs And Intended Result

```text
team-library/
  guidance/workflow.md
  skills/review/SKILL.md
  skills/review/scripts/team-check.sh
  skills/testing/SKILL.md
  skills/deployment/SKILL.md
project/
  .mogent/local/skills/review/SKILL.md
```

The library recommends workflow, review, and testing; deployment is opt-in. Those
recommendations do not choose workspace destinations. Desired result: a project
introduction plus workflow guidance in root `AGENTS.md`, testing copied from
team, review replaced completely by local review, deployment omitted. The old
review script must not survive whole-bundle replacement.

## 2. Bundle Behavior And Naming Reservations

Owner accepts exploring complete-bundle behavior, remains skeptical, and requests
validation through examples. This is not final approval of metadata vocabulary.

- Including a recognized bundle copies its main document AND supporting scripts,
  references, and assets by default. No per-file `include SKILL.md` boilerplate or
  exhaustive additional-file registry is needed. New support files travel with it.
- Excluding the bundle excludes its contents. Excluding only its main document
  while installing the remaining bundle is an incomplete-installation error.
- The viewer may open the main document when its containing node is selected;
  physical file identity and optional heading views remain distinct internally.
- Extracting/rendering a heading creates a separate document contribution; it does
  not edit or partially install the original bundle.
- Metadata can identify a main document for nonstandard layouts; automatic
  discovery is appropriate for recognized layouts. No schema spelling is final.
- No generic non-Markdown warning for expected copied scripts/assets/native JSON
  or TOML. Operation/type mismatch, incomplete bundles, and missing local payload
  dependencies need actionable diagnostics. Copying a script never executes it.
- Existing implementation recognizes exact `SKILL.md`; lowercase DotAgents
  `skill.md` and any target conversion require a deliberate adapter decision.

Naming candidates preserved from discussion:

| Candidate | Open question |
|---|---|
| `skill "review" { ... }` | Owner dislikes making skill a special declaration keyword |
| `is-skill` / equivalent boolean annotation | Role marker, but does not itself locate the main document |
| `entrypoint` | Agent suggestion; owner is unsure of this term |
| `node-file` | Owner candidate for associating a main file with a logical node |
| `content = "skill.md"` | Owner's original example; distinguish primary file from all bundle content |
| `additional/scripts = ...` | Prefer discovered members initially; explicit external relationships may be separate work |

Do not invent a required role/keyword or flatten heading/file relationships merely
to settle these names. Test the behavior first, then choose metadata vocabulary.

## 3. Consumer Plan — Proposed Next Walkthrough Step

This uses source roots and named selections; it does not depend on choosing a
special `skill` keyword. Library metadata schema is deliberately not finalized.
Owner-selected initial layout: root `AGENTS.md`; other managed artifacts under
`.agents/`. Do not introduce symlinks or protocol-specific placement in this slice.
Owner correction: use one `include` contribution family with explicit options,
not separate `replace`/`document` constructs or a `mode` enum. Earlier examples
using those constructs are superseded exploration, not accepted vocabulary.
Further owner correction: ordinary inclusion should not require copy/render/append
flags either. Logical guidance trees produce Markdown through normal traversal;
complete physical files/bundles retain their payload. Sequential document
contributions accumulate in written order. The engine may still distinguish these
operations internally. This behavior follows content kind, not `.md` extension
alone; copied skill documents remain intact rather than being flattened.

```hcl
let repo_name = "Example"

source "team" {
  local = "../team-library"
}

source "local" {
  local = ".mogent/local"
}

include "team:skills" {
  into            = ".agents/skills"
  accept_defaults = true
  exclude         = ["deployment"]
}

include "local:skills/review" {
  into    = ".agents/skills/review"
  replace = true
}

include {
  into = "AGENTS.md"
  text = """
    # ${repo_name}

    Follow the repository's development conventions.
    """
}

include "team:guidance/workflow" {
  into = "AGENTS.md"
}
```

Expected explanations: team review/testing accepted by defaults; deployment
explicitly excluded; local review replaces complete team review and removes the
old script; document contributions retain declaration order. Review default
acceptance remains visible even though it is superseded later.

Option spelling is still proposed. Inline `text` supplies content without a
source reference. Inline text and logical guidance assigned to one generated
document contribute in include order without `append = true`. They form one
document, not independent writers appending to the on-disk file. Whole physical
files/bundles colliding at a destination still need explicit `replace = true`;
sequence alone does not decide replacement versus combining their bytes. Exact
kind/applicability rules, full-file Markdown reuse, and nested include grouping
need specification before implementation.

Owner clarified localization: create an editable source in the workspace's local
library/sublibrary; it is NOT final output placement. Preserve upstream provenance
and subsequently select the local source so edits survive builds. An authoring
operation can create/repoint this local source; ordinary builds must not keep
recopying upstream over the user's local edits. Exact `localize = true` syntax,
authoring command, and path remain undecided; existing v1 localization is relevant.

This is a local authoring example, not a frozen/reproducible build: pin/snapshot
ownership remains to be finalized. Next show the generated concrete build,
source identity and decision provenance, result tree, and an unresolved-choice
variant before writing the executable specification.

## 4. More Concise Inclusion — Under Discussion

Owner proposes child-list shorthand and default destinations:

```hcl
include "team:skills" ["review", "testing"]
include "local:skills/review"
```

This is rough custom syntax, not valid current Mogent HCL grammar. It selects the
same source bundles if aliases/layout match, but destination defaults and collision
handling must be specified before claiming equivalent output. Proposed baseline:
workspace target `.agents`, optional library-relative destination recommendations,
and explicit `into` overrides. Show where each effective value came from. Library
defaults must not authorize writes outside the consumer's workspace target.

Earlier accepted replacement policy still requires explicit intent; concise form
could retain an options block only on the replacement include:

```hcl
include "local:skills/review" {
  replace = true
}
```

`{ ... }` is proposed instruction options/nesting; `= { ... }` would assign an
object-valued field and need not exist in the minimal grammar. One-line versus
multiline layout must preserve identical semantics. Proposed inspection through
the `show` family should distinguish library-declared values, consumer overrides,
and effective values, including inferred/default placement and replacement history.

## 5. Aliases And Destination Conventions — Current Discussion

Owner proposes `replace` as surface sugar for `include` with `replace = true`,
not a distinct semantic operation. Other action spellings such as `copy` may be
considered only if their canonical meaning is precise; ordinary includes still
need no copy flag. Canonical/expanded output should use `include` plus explicit
options. Separate alias expansion from subtree compaction and diagnostic verbosity.
Proposed `--no-command-aliases`/verbose-output flags are not implemented.

Owner asks whether all payloads should live in `.agents` with root symlinks. The
earlier `.agents` default is not a universal discovery contract. Primary docs:

- https://agents.md/ recommends root `AGENTS.md` and nested subproject files.
- https://developers.openai.com/codex/guides/agents-md documents project-root-to-CWD
  discovery, not automatic `.agents/agents.md` discovery.
- https://agentskills.io/integrate-skills describes `.agents/skills/` as a widely
  adopted convention but does not mandate installation paths in the format.
- DotAgents' `.agents/agents.md` is a distinct protocol projection, not equivalent
  evidence of native AGENTS.md discovery.

Agent recommendation pending decision: workspace containment root plus per-role/
target defaults (root AGENTS.md for conventional instructions, .agents/skills for
compatible skills). Keep .agents-only output as a target profile, not universal.
Central storage plus relative root symlinks can be an explicit alternative, not
the default assumption. Supporting generated symlinks needs an explicit state/
transaction/containment contract; current selected source symlinks are rejected.
Do not conflate output links with recipe aliases or localized editable sources.

Owner resolution: omit the `copy` alias. Keep the `replace` alias proposal and
canonical `include ... replace = true` representation. Initial output layout is
root `AGENTS.md`, with other managed artifacts under `.agents/`; use real files.
Earlier `.agents/agents.md` examples are superseded for this initial layout.

Deferred TODO: a file header such as `protocol = "dotagents"` or `"openai"`
could select destination/format defaults. These names are illustrative and are
not supported protocol identifiers or conformance claims. Before implementing,
define versioned target profiles, harness discovery/format support, explicit
`into` overrides, propagation to generated builds, and scope across included
recipes. Return to this after the core workspace workflow; do not block it on a
protocol-selection system.

## 6. Build Header And Resolved Sources — Accepted Example Direction

Owner accepted the header/source walkthrough. Exact schema remains proposed:

```hcl
build {
  format         = 1
  origin         = "generated"
  generated_from = "mogent.plan"
  generated_by   = "mogent"
}
```

Direct builds instead use `origin = "authored"`, with optional human/agent author.
Authored origin prevents automatic regeneration. Generated origin is not overwrite
permission: compare against the last generated content; if that baseline is absent,
do not guess that regeneration is safe.

Plans may follow a Git ref; the resolved build records its exact commit. Local
selected nodes carry whole-bundle content hashes; exact serialization is still
open. Editing local sources causes plan sequencing to refresh those records;
checking the old recorded build reports the mismatch. Hashing does not preserve
uncommitted source bytes or authorize rewriting localized content.

## 7. Localization — Proposed Authoring Syntax

Owner accepted the one-time localization behavior below. Exact provenance schema
and editing grammar remain proposed, not implemented new CLI behavior:

```sh
mogent edit include team:skills/review --localize --replace
```

Authoring operation: copy the complete selected bundle into the editable local
sublibrary (candidate `.mogent/local/skills/review`), record exact upstream source/
revision and baseline provenance, and write a normal local-source include:

```hcl
include "local:skills/review" {
  replace = true
}
```

Reuse/declare `source "local" { local = ".mogent/local" }` explicitly. Existing
local sources must not be overwritten. Retain a generated include's enclosing
collection and add the explicit local replacement when necessary; localization
must not weaken curation/reason requirements. Subsequent builds read local edits
without repeating localization. Never localize scripts/assets separately from a
required whole bundle.

The `--localize` authoring flag is not a permanent executor side effect. If a
future `localize = true` config shorthand is accepted, define its one-time lowering
and write-back behavior before implementing it. Record provenance without altering
copied skill frontmatter or payload bytes. Provenance schema and missing-baseline
handling remain implementation-spec details, not reasons to create another verb.

## 8. New Library Example — Proposed For Review

Owner chose new library syntax as part of the initial redesign. The following
uses the same small block/literal language. Filename `library.mogent`, `node`,
`file`, `directory`, child-local `default`, and `preset` are proposed spellings,
not approved parser schema. No special skill declaration or main-file override
keyword is introduced.

```hcl
library {
  format = 1
  name   = "Team Library"
}

node "guidance" {
  curate = "foundation"

  node "workflow" {
    file = "guidance/workflow.md"
    tags = ["workflow"]
  }
}

node "skills" {
  directory = "skills"
  curate    = "choose"

  node "review" {
    default = true
    tags    = ["workflow", "quality"]
  }

  node "testing" {
    default = true
    tags    = ["lang/go", "quality"]
  }

  node "deployment" {
    default = false
  }
}

preset "recommended" {
  include "self:guidance"
  include "self:skills" {
    accept_defaults = true
  }
}
```

Proposed meanings:

- A top-level node without `directory` declares logical guidance. `file` supplies
  a Markdown leaf; its own heading supplies the title unless deliberately overridden.
  Logical names/placement are independent of physical file locations.
- `directory` imports physical members. Child `node` declarations annotate matching
  discovered paths; they do not enumerate every asset or create new filesystem nodes.
  Scripts/references/frontmatter are preserved on intact inclusion. A generic
  physical Markdown file remains a physical file; `file` here identifies authored
  guidance, not a main-file override inside a copied bundle.
- Library child declaration order controls named siblings; unannotated physical
  children remain discovered and use lexical fallback. A choose branch requires
  complete defaults before accepting them; a newly discovered undecided child is
  reported rather than hidden or silently accepted. Exact library-load versus
  consumer-diagnostic treatment needs review with partial inspection.
- `default = false` is a recommendation to omit deployment, not a new experimental
  or reason-required marker. Curation applies to the enclosing collection.
- `self:` is a proposed library-local alias. A preset is a reusable selection of
  existing nodes, not a content copy or an installer. Consumer includes remap its
  references to the pinned library source and retain selection provenance.
- Consumer `include "team:recommended"` expands guidance then selected skills.
  Defaults put guidance in root AGENTS.md and these bundles under .agents/skills;
  consumer `into` can override valid destinations. This example does not hide
  extra deployment actions, fetches, or global writes in metadata.

The next review should check whether one `node` form is clear for both authored
guidance and physical annotations, and whether `file` versus `directory` is enough
to express the distinction without reintroducing public copy/render switches.

### Owner Feedback: Typed Nodes And Input Variables

The owner suggests explicit types (`md-file`, `md-heading`, `dir`, etc.) with
surface aliases normalizing to generic node declarations, an intended agent/skill
target type, listed Markdown variables, and main content-to-target association.
The untyped example above is not accepted as final grammar.

Candidate separation, pending review:

- Structural type: physical file/directory versus a document-heading view.
- Format: Markdown when document parsing is intended; `.md` alone does not force
  intact copied payloads to become templates or generated documents.
- Role: skill, instructions, or tool-native agent profile; independent of structure.
- Locator: original payload and optional main-file association such as `content`.
- Output: existing consumer `into` or visible destination defaults, not role-driven
  hidden conversion. Main-file renaming is not implied by bundle inclusion.
- Variables: an explicit required-name whitelist for rendered authored guidance;
  consumer `let` values bind those declared names. Missing values are located errors;
  unrelated ambient/env/caller names are not implicitly injected. Copied skills and
  opaque files remain byte-preserving. Cross-recipe parameter passing stays separate.

Exact field names and the meaning of md-file alias versus physical file plus
document view need review before the consolidated spec adds these constructs.

Latest owner preference: call logical guidance `md-section`; `md-file` should mean
intact Markdown. Working distinction: a section contributes generated document
content, whereas a file preserves the original complete payload. Keep optional
heading inspection available without implicitly changing installation behavior.
This refines the earlier `md-heading` spelling, not a new feature. The typed
library proposal is now reconciled into spec section 4 and acceptance cases L01-L13.
Use those for the final combined review; the untyped example above remains the
historical discussion step, not the current proposed executable contract.
