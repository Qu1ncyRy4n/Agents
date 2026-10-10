# TE-jofit: One Content Tree, Two Output Modes

TE ID: TE-jofit
Status: needs DF
Date: 2026-10-09
Handle: minted with `go run ./cmd/mogent mint -n` against the repository corpus.

## Question And Owner Intent

Can Mogent use one selection and composition model for both agent instructions
and directories of skills, while keeping the output behavior explicit?

The owner finds `sections` versus `section` confusing and asks whether
`document`, `content`, or another name can cover both cases. The owner prefers
an output name such as `dir-tree` over the ambiguous `tree`, and wants to
understand what separate overlay terminology buys beyond ordered composition.

Already chosen: use `curate` now and retain the word `exclude`. Not yet chosen:
root/child names, selection grammar, output-kind rename, collision defaults,
placement, or lower-layer removal syntax. This TE narrows those decisions;
its examples do not change the implemented HCL contract.

## Starting Point: What Exists Today

| Concern | Implemented v2 behavior |
|---|---|
| Library Markdown | One `content { markdown_root = "agents" }` and a `sections { section ... }` hierarchy |
| Library directories | `trees { tree "skills" { root = "skills" } }`; optional `entry` discovery metadata |
| Consumer selections | Markdown uses sparse `select = { ... }` objects; directory outputs require `{ all = true }` |
| Markdown output | `kind = "markdown"`; several sources append selected content in source order |
| Directory output | `kind = "tree"`; exactly one source, copied transactionally |
| Directory exclusion | Source-level `exclude = ["path"]` filters that source; paths must exist |
| Ownership | Generated state detects direct edits/unmanaged outputs; `apply` requires `--force` to replace them |

Two details matter: directory `entry` metadata does not currently select files,
and matching Markdown heading titles do not cause automatic merging/replacement.
No output mount paths, multi-source directory compositor, or block selection
grammar exist yet.

## Separate Three Jobs

1. **Select:** which library nodes are accepted, including curation/defaults.
2. **Place and compose:** where contributions belong and in which order.
3. **Materialize:** render a document, or preserve selected files/directories.

"Tree" describes the structure used in all three jobs; it does not by itself
tell us whether the final result is a file or a directory. "Overlay" describes
one directory collision strategy, not a separate selection language.

## Naming Experiment

Read each root name with both `agents` and `skills` before considering parser
changes:

| Root name | Fit | Cost |
|---|---|---|
| `sections` | Clear for headings | Implies document sections when children are directories/files; plural/singular confusion persists |
| `document` | Clear for an authored document | A skills directory is not a document; confuses library content with generated output |
| `tree` | Structurally universal | Does not distinguish selectable structure from a filesystem result |
| `collection` | Works for both | Does not communicate hierarchy; could imply an unordered set |
| `content` | Natural for both reusable guidance and skill files | Already used by the current Markdown-root block; requires an explicit schema migration |

For children, compare `section`, `entry`, and `node`. `section` is most readable
for prose but misleading for files. `entry` fits inventories but weakly implies
hierarchy. `node` is a neutral branch/leaf term at the cost of being technical.

**Surviving naming pair:** `content "agents"` / `content "skills"`, with nested
`node` blocks where the author declares structure. Keep `document` as a word
for the rendered result, not the universal source root.

For outputs, `kind = "markdown"` versus `kind = "dir-tree"` states the result
clearly. `directory` is a viable shorter alternative; `dir-tree` makes recursive
copying explicit. A hyphen is valid here because the value is a quoted string.
Neither name means "skills only".

## Broad Model Alternatives

### A. Keep Separate Markdown And Directory Models

Preserve `sections` and `trees`; add mounts, ordered layers, and a directory
selector independently. Lowest migration cost, but two selection surfaces can
diverge, and explaining the system still requires "sections versus trees".
Survives as an incremental implementation strategy, not the preferred end state.

### B. Let The Filesystem Be The Only Tree

Every directory is a branch and every file is a leaf. Use that structure both
for document headings and copied directories. Excellent for skill bundles and
automatic discovery, but reorganizing files changes consumer references and
rendered headings. A library cannot group physically distant guidance under a
logical heading without moving it. Reject as the universal model.

### C. Declare Every File In One Logical Tree

Author all content with `content` and `node`, including every skill script and
asset. Both outputs select nodes uniformly. This preserves layout independence
but duplicates the filesystem and makes ordinary skills expensive to maintain.
Reject mandatory per-file declaration; retain explicit logical trees for prose.

### D. One Content Model With Authored And Imported Trees

Expose named content roots. A root can carry an authored logical hierarchy or
import a directory hierarchy. Both resolve to selectable nodes with stable
root-qualified identity; nodes retain enough type/payload information to render
or copy without guessing. Physical import exposes directories/files automatically;
metadata can annotate them without becoming a compulsory file registry.

This reconciles directory levels with the existing logical tree instead of
forcing all trees to be physical. It is the recommended survivor. The cost is a
clear contract for adapters, metadata overrides, and output applicability.

## Comparative Examples

All HCL below is proposed syntax. It is deliberately small enough to evaluate;
validation details require the DFs below.

### Library: An Authored Guidance Root

```hcl
content "agents" {
  basedir = "guidance"

  node "workflow" {
    title  = "Workflow"
    curate = "choose"

    node "review" {
      title   = "Review"
      source  = "process/review.md"
      default = true
    }

    node "testing" {
      title   = "Testing"
      source  = "quality/testing.md"
      default = false
    }
  }
}
```

The logical node `workflow/review` maps to `guidance/process/review.md`.
`curate` governs a branch's direct children; a child's `default` is a
recommendation only when its parent permits defaults. Moving defaults does not
move or remove the parent's curation rule.

Proposed path rule: nested `basedir` values compose relative to the inherited
base. With root `basedir = "guidance"`, child `basedir = "process"`, and leaf
`source = "review.md"`, the physical file is `guidance/process/review.md`.
Reject absolute paths, traversal, symlink escapes, and missing files rather than
resetting to an ambient directory. A logical node without `basedir` inherits the
base; its name does not implicitly become a directory component. Root names and
canonical sibling paths must be unique. These are proposed DR constraints.

### Library: Import A Directory Root

```hcl
content "skills" {
  directory = "skills"
}
```

Given `skills/review/SKILL.md` and `skills/review/scripts/check.sh`, import
produces a `review` directory node containing both files. No `node` declaration
is required for each script. Relative imported paths provide node identity.
The exact metadata-annotation syntax remains open; skill tags/TLDRs and curation
must remain possible without hiding undeclared files.

### Consumer: The Same Selection Vocabulary

```hcl
output "instructions" {
  path = "AGENTS.md"
  kind = "markdown"

  source "qmr" {
    from = "qmr:agents"
    select {
      node "workflow" {
        accept_defaults = true
      }
    }
  }
}

output "skills" {
  path = ".agents/skills"
  kind = "dir-tree"

  source "qmr" {
    from = "qmr:skills"
    select {
      node "review" {}
    }
  }
}
```

Results:

```text
AGENTS.md                         .agents/skills/
  # Workflow                        review/
    ## Review                         SKILL.md
      review body                     scripts/check.sh
```

The grammar for accepting a branch is the same. The payload operation differs:
render the authored Markdown hierarchy versus preserve the selected directory
bundle. An empty selected branch block means accept that branch subject to its
curation rule, not override all descendants regardless of curation.

### Consumer: Placement And Ordered Contributions

```hcl
output "agent-files" {
  path      = ".agents"
  kind      = "dir-tree"
  collision = "later_wins"

  source "company" {
    from   = "company:skills"
    into   = "skills"
    select = { all = true }
  }

  source "repo" {
    from = "repo:skills"
    into = "skills"
    select = { all = true }
  }
}
```

This alternative keeps object selection to show that unified nodes do not
require switching to blocks. In a finalized grammar, use one selection form
consistently; do not adopt two spellings merely because this TE compares them.

If company supplies `review/SKILL.md` and `rust/SKILL.md`, and repo supplies
`review/SKILL.md`, the final tree retains Rust and uses repo's review file.
`into = "skills"` is destination placement, not a source lookup or network mount.
Proposed plan output names both contributions and the winner. The spelling
`collision = "later_wins"` is illustrative, not an approved new attribute.

## Tabletop Scenarios And Failure Cases

These are reasoning experiments, not claims of executable acceptance tests.

| Scenario | A: separate models | B: filesystem only | C: fully declared | D: hybrid content |
|---|---|---|---|---|
| Logical workflow spans physical directories | Works for Markdown only | Forces moves or wrong headings | Works | Works through authored nodes |
| Skill gains a new supporting script | Existing raw copy works | Discovered automatically | Registry must be edited | Imported automatically |
| Select one skill with every required asset | Needs new directory selector | Directory subtree works | Depends on registry completeness | Directory subtree works |
| Same selections produce different output forms | Separate contracts | Easy, but headings tied to paths | Possible with explicit payloads | Common selector; materializer checks applicability |
| New child under a `choose` branch | Existing incomplete-choice rule | Needs metadata/default rules | Can require a declared default | Must expose new child and missing default, not silently accept |
| Two sources supply one destination file | Needs collision rule | Needs collision rule | Needs collision rule | Still needs collision rule; unification does not answer it |

### Logical Branches Are Not Automatically Directories

Selecting `workflow` does not tell the copier whether review belongs at
`process/review.md`, `workflow/review.md`, or a new name. For the first cut,
imported roots preserve their relative paths. Copying authored logical roots
requires an explicit path policy; reject it until specified rather than invent
filenames from titles. Similarly, a binary file cannot be rendered as Markdown.
One selector does not mean every payload supports every output kind.

### Directory Merge Is Different From Whole-Bundle Replacement

Company review has `SKILL.md` and `scripts/old.sh`. Repo review has only
`SKILL.md`. A recursive merge keeps `old.sh`. Replacing the entire review
directory removes it. The plan must distinguish these results; order alone
does not tell us which one the consumer meant.

Two different files at one destination can either be an error by default or
implicitly later-wins as recommended in the October review. Both remain viable.
Prefer error-by-default plus an explicit later-wins choice for a first cut;
this is a TE recommendation, not an owner decision. Directory/file type changes
should require explicit replacement rather than silently erasing a subtree.

### Exclusion Has Two Scopes

Current `exclude = ["review"]` means "filter review out of this source". If
company already contributed review, excluding review from repo does not remove
company's version. Preserve that behavior during migration.

Lower-layer removal is a distinct operation. Keep the word `exclude`, but
compare explicit scope alternatives before introducing it:

```hcl
# Alternative 1: retain the list and add a separate lower-layer block.
exclude = ["experimental"]
exclude {
  scope = "earlier"
  paths = ["review"]
}
```

```hcl
# Alternative 2: require a block for every exclusion, including source filters.
exclude {
  scope = "source"
  paths = ["experimental"]
}
exclude {
  scope = "earlier"
  paths = ["review"]
}
```

Each path would be relative to that contribution's `into` destination when
scope is `earlier`, and relative to its source root when scope is `source`.
Earlier-path removal occurs before copying the current contribution, so the
current contribution may add a replacement. Neither syntax is decided. Reject
silently changing today's exclusion list into lower-layer deletion.

### Ordering, Discovery, Safety, And Updates

- Authored siblings retain sidecar order; imported children need deterministic
  path order. Consumer source-block order governs contributions. A block selector
  must explicitly state whether its declaration order affects rendering; do not
  accidentally change today's sidecar-order rendering during a syntax migration.
- Curation applies to both authored and imported branches. Define how an
  imported `choose` branch annotates every direct child's default, and diagnose
  newly discovered children. Do not omit unannotated files from discovery.
- A common selector must carry `else`, incomplete-choice diagnostics, and
  foundation refusal with `force_exclude`/`reason` across both adapters.
  Avoid a new path filter that silently bypasses an imported branch's declared
  curation; decide that interaction in the DR before adding curation to imports.
- Tags remain metadata. An opaque skill bundle's tags should select that bundle,
  not render its scripts as guidance or invent directory headings.
- Plan must show exclusions and replacements even when final bytes happen to
  equal existing output. File hashes alone cannot explain contribution history.
- Retain full source/path validation, symlink rejection within selected roots,
  file modes, transactional staging, rollback, and direct-edit protection.
- State covers the complete final output. Ownership cannot depend on Git status
  or on only the highest-priority source.
- The current `update` implementation previews Markdown outputs, not directory
  source-to-source changes. Add that preview before claiming equivalent review
  support for multi-source directory outputs; `update --accept` still only moves
  pins, and `apply` materializes outputs.
- Overlapping managed output targets remain invalid, including `AGENTS.md`
  inside a separately managed directory output containing that path.

## Narrowed Result

Recommend D: one content model with authored and directory-import adapters.
Retain A as a staged implementation route. Eliminate filesystem-only hierarchy
and compulsory per-file sidecar registries as universal solutions.

Recommend `content` + `node`, and `markdown` + `dir-tree`, as a coherent naming
set for owner review. These names distinguish reusable input from output form
without making skills a special type. Prefer shared block selection if avoiding
reserved-name collisions outweighs its verbosity; retaining objects remains a
viable lower-migration-cost alternative.

Composition remains ordered in either design. Placement and collision handling
are capabilities attached to that model; they do not justify a second selector.
Preserve source-filtering exclusions. Decide scoped lower-layer removal explicitly.

## Decision-First Questions For Owner Review

1. **DF-1 Model:** accept the hybrid content model, or keep separate models and
   add directory selection incrementally?
2. **DF-2 Names:** adopt `content` / `node` and output `dir-tree`, or retain
   current names? `directory` is the surviving alternative to `dir-tree`.
3. **DF-3 Selection/defaults:** adopt one nested `select { node ... }` surface
   with defaults beside children, or retain object selection/parent maps?
   Preserve library order unless consumer reordering is separately requested.
4. **DF-4 Collisions:** error-by-default with explicit later-wins, or implicit
   later-wins with mandatory plan reporting? Decide whole-directory replacement
   and type changes independently from ordinary file collisions.
5. **DF-5 Exclusion:** retain the source list plus an explicitly scoped block,
   or migrate all exclusions to scoped blocks? Keep source filtering distinct
   from removal of earlier contributions.
6. **DF-6 First cut:** limit directory materialization to imported roots and
   Markdown rendering to authored roots, then add cross-kind payload rules only
   when examples justify them? This avoids a speculative filename/asset renderer.

No DIs or final schema decisions are recorded by this TE. `curate` and the
preference for `exclude` are inputs from the owner, not new decisions made here.

## Delivery After Decisions

1. Record owner answers in DRs and decision intents; reconcile the v2 spec.
2. Add named content adapters, node inventory, common selection/default rules,
   applicability diagnostics, and deterministic order tests.
3. Add directory contribution placement, collision reporting, and the chosen
   scoped exclusion/replacement behavior on the existing transaction engine.
4. Extend pin-update review to directory content; verify state and rollback
   behavior across mixed output kinds.
5. Provide an explicit preview/accept migration for existing `sections`,
   `trees`, `content.markdown_root`, kind values, defaults, and selectors that
   actually change. Do not infer a directory import from a heading tree.
6. Exercise company base + QMR skills + repo overrides, direct edits, source
   updates, and absent/unsafe paths in temporary consumer fixtures.

## References

- [V2 specification](../proposals/MOGENT-HCL-V2-SPEC.md)
- [October composition review](../reviews/2026-10-05.md)
- [Composition follow-up roadmap](../../TODO/TODO-tugur-v2-composition-followup.md)
- [Curation naming DR](../../DR/DR-nalan-offer-attribute-name.md)
- [Include/composition narrowing TE](TE-nufad-module-includes.md): distinguishes
  composition from exact insertion and template evaluation; those remain separate.
- [Source-subdirectory TE](TE-vurap-url-source-subdirectories.md): containment
  reasoning remains relevant; its YAML lock contract is historical v1 behavior.
- Current code: `v2/library.go`, `v2/config.go`, `v2/plan.go`, `v2/tree.go`,
  `v2/build.go`.
