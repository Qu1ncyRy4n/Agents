# TE-gunak: Typed Agent Content, Explicit Operations

TE ID: TE-gunak
Status: decided, refined
Date: 2026-10-09
Handle: minted with `go run ./cmd/mogent mint -n` in the repository.

Current decision status: model, intact-copy/explicit-render direction, explicit
replacement policy, and deferred integrations accepted under
[DR-juvih / DI-juvih](../../DR/DR-juvih-typed-content-composition.md).
Grammar, parsed-heading contract, and delivery sequencing remain open. The body
below preserves the pre-decision analysis; the owner approval is recorded in the
final dated refinement.

First-delivery update: additive grammar, parsed heading views, core APIs and CLI
usage are implemented under DI-juvih-first-slice. See
[Typed Content Composition](../TYPED-CONTENT.md) for the current contract;
the earlier open questions below remain the historical design analysis.

## Question and baseline

Can Mogent serve instructions, skills, research/docs/guides, arbitrary `.agents`
files, and global/project/nested origins through one composition vocabulary,
without treating output `markdown` versus directory as the only type distinction?

The owner questions output-only typing and asks us to consider typed `dir`,
`file`, and `heading` nodes. This experiment recommends a typed common inventory
over physical and document adapters. It does **not** decide names, selector
grammar, collision defaults, bundle replacement, or scope projection rules.
`curate` and retaining the word `exclude` are existing inputs.

[TE-jofit](TE-jofit-unified-content-composition.md) is the earlier local design
baseline needing refinement, not a decided contract or sourced ecosystem report.
Its authored/imported distinction and select/place/materialize decomposition
survive. Its `markdown`/`dir-tree` output distinction is insufficient for mixed
directory artifacts. Its body remains unchanged. The accompanying
[ecosystem report](../../reports/Universal%20agent%20content%20composition.md)
synthesizes all three [research notes](../../research_notes/Universal%20agent%20content%20composition/).

Today's v2 planner permits one tree source with `all`; authored sections are not
a parsed heading AST. Rendering validates the first `# Title`, removes it, and
normalizes the body. Copying preserves bytes/file permission bits, not every
filesystem metadata property. All syntax below is **proposed**, not executable
current configuration. ([Planner](../../v2/plan.go),
[Renderer](../../v2/render.go), [Copier](../../v2/tree.go))

## Evidence constrains the universal promise

Agent Skills defines exact `SKILL.md`, required frontmatter, and arbitrary bundle
resources, not mandatory installation roots. `.agents/skills` has documented
support in Codex, OpenCode, Copilot, and Cursor; Claude has documented native
`.claude/skills`. No examined source establishes automatic universal discovery of
`.agents/research`, `.agents/guides`, agent profiles, hooks, or settings there.
Those files must remain useful opaque content without invented loading behavior.
([Specification](https://agentskills.io/specification),
[Integration guide](https://agentskills.io/integrate-skills),
[Codex](https://developers.openai.com/codex/skills),
[OpenCode](https://opencode.ai/docs/skills/),
[Copilot](https://docs.github.com/en/copilot/concepts/agents/about-agent-skills),
[Cursor](https://cursor.com/docs/context/skills),
[Claude](https://code.claude.com/docs/en/skills))

agentsstandard.com's broader protocol is independently proposed. Live registry
v1.4.0 and pinned repository v2.0.0 disagree about loading and MCP placement;
v2 uses project-root `.mcp.json` and a non-runtime home catalog. Its global/root
skill replacement differs from Codex's documented coexistence, and its root/depth
prose differs from its loader. It is an optional versioned projection candidate,
not the composition engine's inherited scope law.
([Live registry](https://agentsstandard.com/agents.json),
[Pinned v2 spec](https://raw.githubusercontent.com/nbiish/agents-standard/63e166a76acb61e6a197e2049ba1e87bea0ddf1a/llms.txt),
[Pinned loader](https://raw.githubusercontent.com/nbiish/agents-standard/63e166a76acb61e6a197e2049ba1e87bea0ddf1a/cli/lib/loader.js))

## Broad alternatives

| Model | Strongest case | Failing case / cost | Assessment |
|---|---|---|---|
| A: current prose model plus opaque overlays | Small migration; whole-file distribution | Shared selection and bundle integrity remain separate | Credible staged route |
| B: filesystem-only tree | Automatic file discovery and copying | Logical headings become folders; reusable prose follows physical moves | Strong physical-only subset |
| C: declare one strict ownership tree | A uniform traversal looks simple | Duplicate asset registry; reuse and physical/document parents conflict | Reject mandatory registry/parent conflation |
| D: typed facade over physical/document adapters | Common selection, provenance, and mixed manifests | Requires explicit capabilities and edge rules | Preferred hypothesis |
| E: wholly separate products/selectors | Prose-only and distribution-only interfaces stay simple | Curation and diagnostics can diverge | Viable if common interface costs outweigh reuse |
| F: generic semantic AST composition | Fine-grained settings/heading patches | Tool schemas, serialization, and precedence become the product | Deferred; raw config import instead |

D can internally use a discriminated union or separate structures. The decision
is about preserving domain rules, not prescribing structs. A physical-only user
should not author document views; a prose-only library should not enumerate
physical nodes merely to render fragments. Simple trees remain simple.

## Types distinguish relationships and operations

`dir` means physical containment; `file` means a byte-bearing payload; `heading`
means document structure. A directory does not inherit hidden Markdown rules,
skill identity, or parsing based on its name. A heading owns subordinate document
contributions, **never arbitrary files**. It can reference a payload; a bundle
can separately reference an asset. A generated file can expose heading children
in a document view, but that view must be explicit.

Roles and bundles are separate annotations. A Markdown file can be both a guide
and an instruction source; a directory can be a declared skill bundle or simply
a folder. A profile is a recipe referencing content, not a fabricated ownership
parent. One physical file can supply two logical contributions without acquiring
two storage parents. Qualified logical IDs, physical locators, display titles,
and destination paths are distinct.

Selection should use one vocabulary with kind-aware addresses. Copying a file
and selecting its derived heading view are different requests. Current authored
fragments can back initial headings; arbitrary extraction needs later rules for
frontmatter, repeated headings, body headings, code fences, and stable addresses.
No heading-title coincidence implies an override.

## Concrete typed inventory and mixed export

Assume source aliases `team` and `local` are already configured. Imports discover
all directories/files; annotations do not restrict inventory to declared members.
The following candidate library sidecar illustrates kinds and orthogonal policy:

```hcl
content "physical" {
  kind   = "dir"
  source = "payloads"

  annotate "skills/review" {
    role   = "skill"
    bundle = "review" # identifies a boundary; atomicity remains a DF
  }
  annotate "research/design.md" { role = "research" }
  annotate "agents/reviewer.md" { role = "claude-agent" }
}

content "guidance" {
  kind = "heading"

  node "workflow" {
    kind  = "heading"
    title = "Workflow"
    node "review" {
      kind   = "heading"
      title  = "Review"
      source = "guidance/review-fragment.md"
    }
  }
}
```

`local:physical` separately imports local payloads containing complete
`config/codex.toml` and `config/claude-settings.json`. The review skill contains
`SKILL.md`, `scripts/check.sh`, and `references/checklist.md`; no per-file sidecar
registration is necessary. `guidance/review-fragment.md` satisfies the explicitly
authored first-heading fragment contract, not the copied skill contract.

```hcl
output "project-package" {
  path = "agent-export"
  kind = "dir"

  source "team" {
    from = "team:physical"
    select {
      node "skills/review" { kind = "dir" }
    }
    operation = "copy"
    into      = ".agents/skills/review"
  }
  source "team" {
    from = "team:guidance"
    select {
      node "workflow" { kind = "heading" }
    }
    operation = "render-markdown"
    into      = "AGENTS.md"
  }
  source "team" {
    from = "team:physical"
    select {
      node "research/design.md" { kind = "file" }
    }
    operation = "copy"
    into      = ".agents/research/design.md"
  }
  source "team" {
    from = "team:physical"
    select {
      node "agents/reviewer.md" { kind = "file" }
    }
    operation = "copy"
    into      = ".claude/agents/reviewer.md"
  }
  source "local" {
    from = "local:physical"
    select {
      node "config/codex.toml" { kind = "file" }
    }
    operation = "copy"
    into      = ".agents/config/codex.toml"
  }
  source "local" {
    from = "local:physical"
    select {
      node "config/claude-settings.json" { kind = "file" }
    }
    operation = "copy"
    into      = ".claude/settings.json"
  }
}
```

Expected bounded snapshot:

```text
agent-export/
  AGENTS.md                           # rendered Workflow → Review
  .agents/skills/review/SKILL.md        # complete copied bundle
  .agents/skills/review/scripts/check.sh
  .agents/skills/review/references/checklist.md
  .agents/research/design.md           # exact copy, not auto-loaded
  .agents/config/codex.toml            # storage copy, not native config path
  .claude/agents/reviewer.md            # intact Claude dialect
  .claude/settings.json                # complete JSON file
```

The output is one directory artifact, not a wholesale managed repository root.
Internal paths describe a project export; installation/discovery is a separate
projection. Selecting the bundle root preserves its interior layout. Rendering
explicitly produces AGENTS.md; files with `.md` remain opaque unless requested
through a renderable document view. A standalone copied file would use the same
file selection and copy operation. `kind`, `operation`, `into`, annotations,
root/child names, and block selectors are candidate spellings, not approvals.

Copied frontmatter stays at the first line; no provenance banner is inserted.
TOML, JSON/JSONC, YAML, and sidecars are entire files, with no parse/serialize or
key/table/array merging. Native host merges can occur after explicit export to
recognized paths, but Mogent does not simulate them. Copying cannot automatically
repair external paths or translate Claude agent identity into OpenCode identity.
([Claude parsing](https://code.claude.com/docs/en/skills#frontmatter-reference),
[Claude settings](https://code.claude.com/docs/en/settings),
[OpenCode configuration](https://opencode.ai/docs/config/),
[OpenCode agents](https://opencode.ai/docs/agents/))

## Tabletop tests expose independent decisions

These are reasoning tests, not implemented acceptance tests.

| Scenario | Required observable result | Decision pressure |
|---|---|---|
| Skill gains a new script | Import discovers and preserves it | No mandatory file registry |
| Tag matches SKILL.md only | Bundle promotion or explicit incomplete-selection diagnosis is explained | Selection atomicity |
| Company review has old script; repo has SKILL.md only | Union retains script; replacement removes it; plan distinguishes both | Bundle replacement versus derivative overlay |
| One-script patch is intentional | Explicit derivative preserves other members and exposes mixed origins | Atomic policy must not prohibit intended customization |
| Two guides contain “Examples” | No implicit shared heading slot | Identity is not display title |
| A logical heading references a guide and figure | Figure remains a file reference/bundle member, not heading-owned storage | Edge rules |
| Project JSON contains one setting | It replaces the complete file, never acts as a patch | Opaque import |
| Directory happens to be named `markdown` | Files remain physical payloads | No hidden format inheritance |
| Same path receives copy and render | Resolve before writes through one final artifact owner | Destination collision, not source identity |
| Later source excludes its review | Earlier review remains unless explicit earlier-layer removal is requested | Source filter versus removal |
| File replaces directory with descendants | Diagnose or explicitly replace the subtree under chosen policy | Ancestor/type conflict |
| Equal bytes have two origins | Plan retains contribution/winner history | Provenance beyond content diff |
| Global guide is flattened and also host-loaded | Duplication is shown or avoided by chosen projection | Runtime discovery differs from build order |
| Nested worktree content is imported | Retain actual boundary; do not inherit main-checkout files invisibly | Context identity |

Layer order cannot decide recursive union versus replacement. Selection
completeness, version coherence, member-set replacement, and transactional
installation are four different atomicity properties. Dependencies outside the
bundle require explicit accounting: choose explained expansion or missing-member
diagnostics, not an assumed package solver. Excluding a requirement must not
silently resurrect it. External executables can remain declared prerequisites.

The [October review](../reviews/2026-10-05.md) recommends implicit later-wins and
type replacement; TE-jofit recommends explicit conflict/type handling. Both
remain alternatives. A coherent skill default is desirable, but neither atomic
replacement nor per-file overlay is an owner decision yet.

## Scopes and projections stay explicit

Keep origin/discovery scope, applicability, build contribution order, document
order, and host runtime behavior separate. Identity should retain context root
and node/root qualification so global and nested equal paths are not collapsed.
The exact serialization remains open. No ambient host merge enters planning.
Global export needs an explicit authorized root/state lifecycle, not relaxed
`../` output containment or an assumed alias for project output.

Codex instructions use `$CODEX_HOME`, not proposed `~/.agents/AGENTS.md`; its
skills traverse project ancestors and can preserve duplicate names. OpenCode
skills stop at the Git worktree and its global instructions use the XDG-style
OpenCode home. Claude accumulates instructions and conditionally supports AGENTS
from v2.1.277; its instruction mechanism excludes `.agents` contents. Export
profiles must state separate-context versus flattening behavior and report
unrepresentable metadata. ([Codex instructions](https://developers.openai.com/codex/guides/agents-md),
[Codex skills](https://developers.openai.com/codex/skills),
[OpenCode rules](https://opencode.ai/docs/rules/),
[OpenCode skills](https://opencode.ai/docs/skills/),
[Claude memory](https://code.claude.com/docs/en/memory#agents-md))

Discovery-supported skill/instruction paths deserve first adapters; native
rules/profiles next, commands/prompts as qualified compatibility inputs, settings
and hooks as whole-file payloads. Cursor requires `.mdc`; copying hook scripts
does not register them. VS Code Agent Host no longer loads prompt files despite
broader GitHub availability documentation. OpenCode evidence covers fetched
`/docs/`, not independently checked separate v2 documentation. These are support
qualifications, not popularity rankings.
([Cursor rules](https://cursor.com/docs/context/rules),
[Claude hooks](https://code.claude.com/docs/en/hooks),
[VS Code prompts](https://code.visualstudio.com/docs/copilot/customization/prompt-files),
[GitHub prompts](https://docs.github.com/en/copilot/tutorials/customization-library/prompt-files/your-first-prompt-file),
[OpenCode docs](https://opencode.ai/docs/config/))

## Plan and ownership expectations

Explain every contribution with source alias/pin or local snapshot, context,
qualified node/kind, role/bundle, selection/default/dependency reason, payload
operation, placement, exclusions, replacements, and effective winner. Removed
bundle members and equal-byte collisions remain visible. Build one final manifest
for the mixed target; independently owning AGENTS.md inside that same managed
directory would create overlapping writers. Preserve direct-edit protection,
mode handling, containment, rollback, and complete-result state. Directory update
review needs manifests/member-set changes, not only rendered string diffs.

## Decision-first questions

1. **DF-1 Universal boundary:** choose D's typed facade, A's staged extension,
   or intentionally separate products? Which physical/document relationships
   must be visible to users?
2. **DF-2 Types and grammar:** should `dir/file/heading` be explicit declaration
   kinds, adapter-resolved inventory kinds, or both? Compare typed blocks,
   `node` with kind attributes, and objects; settle names separately.
3. **DF-3 Mixed artifacts:** accept independent output shape and copy/render
   operations with one manifest? Define placement/path preservation without
   guessing names from titles.
4. **DF-4 Bundles and collisions:** choose selection promotion versus diagnostics,
   replacement versus per-file derivatives, and conflict defaults/type changes
   independently. How are external content requirements accounted for?
5. **DF-5 Document scope:** initially render authored fragments only, or also
   specify an imported-heading adapter? Define stable slots before section
   overrides; preserve whole-file frontmatter without semantic config merging.
6. **DF-6 Context projection:** how are global/project/nested/worktree roots
   identified, owned, and projected? Which exports remain separate, which flatten,
   and how does plan reveal receiving-host duplication or unsupported discovery?
7. **DF-7 Import sequence:** prioritize documented bundles/instructions, then
   vendor-native files? Keep the versioned nbiish proposal opt-in rather than
   treating its adoption claims as native support.

No decision intents, schema changes, or implementation approvals are recorded.
After owner answers, revise DR/specification and migration guidance, then design
focused acceptance fixtures for the tabletop failures. The
[follow-up roadmap](../../TODO/TODO-tugur-v2-composition-followup.md) tracks this
work without changing existing completion statuses.

## Refinements

### 2026-10-09 — Tabletop execution: skills have document views

The owner asked to run the TE, pointed out heading subtrees inside skills, and
approved backlogging the secondary integrations. The following is the execution
record of concrete reasoning cases, not runtime tests of the proposed grammar.
Current Mogent cannot execute these candidate typed-node operations. Analysis is
complete for this pass; the TE remains `needs DF` for the narrowed choices below.
No DIs occur in this file; this entry extends the analysis without locking a
schema or changing the earlier alternatives into owner decisions.

#### Input fixture

Use two declared local libraries and one explicitly declared global library:

```text
team/
  skills/review/
    SKILL.md
    scripts/check.sh
    references/checklist.md
  config/tool.toml
  guidance/project.md
repo/
  skills/review/
    SKILL.md
  guides/testing.md
global/
  guidance/project.md
```

Team's `SKILL.md` is a normal complete skill, conceptually containing:

````markdown
---
name: review
description: Review a change using the project checklist.
---
# Review

## Before starting
Read [the checklist](references/checklist.md).
Run `scripts/check.sh` from the skill directory.

## Procedure
Inspect the change against the project requirements.

### Verification
Record the check result.

## Examples
```text
## This is example text, not a selectable heading
```
````

Repo's version contains valid skill frontmatter and its own procedure, but no
script or reference files. Its relative references do not require team assets.
The TOML file is a complete configuration payload, not a patch.

Physical and document views share one payload identity:

```text
review/                         dir, role=skill, bundle boundary
├── SKILL.md                     file, original bytes/frontmatter
│   └── document view           structural view of SKILL.md
│       └── Review              heading
│           ├── Before starting heading
│           ├── Procedure       heading
│           │   └── Verification heading
│           └── Examples        heading
├── scripts/check.sh             file
└── references/checklist.md      file
```

The `document view` edge is not an extra filesystem directory. Selecting it
does not duplicate SKILL.md or add heading-named directories. A heading can
reference another payload, but its structural children are headings, not assets.
This makes a skill heading subtree possible without redefining skill storage.

#### Executed cases

| Case | Request and predicted observable result | Alternative eliminated / narrowed |
|---|---|---|
| R1: intact installation | Copy team's review directory. All three files arrive with original bytes and relative layout, including frontmatter. | Output-wide rendering or treating every `.md` as a fragment fails. |
| R2: skill heading inventory | Open SKILL.md's document view. Procedure contains Verification; the fenced example is not a heading. | File-leaf-only inventory is insufficient; naive line splitting on `#` fails. |
| R3: reuse a skill subsection | Render Procedure plus Verification into a separate AGENTS.md. Team's installed SKILL.md and other bundle members remain unchanged. The copied skill's metadata is not pasted as document body. | Selecting a heading must not mean deleting the file's other headings. |
| R4: extract a reference-bearing section | Render Before starting. Its checklist link needs an explicit destination/reference policy. Record the dependency and reject an unresolved output link instead of silently moving broken prose. | Extraction cannot imply automatic asset copying or arbitrary link rewriting. |
| R5: replace the whole skill | Choose repo review as an explicit bundle replacement. The final review has repo SKILL.md; team scripts/check.sh and references/checklist.md are absent. Plan lists both removals. | Recursive per-file union is not equivalent to bundle replacement. |
| R6: deliberate derivative | Keep team bundle and request a heading edit from repo. A derived SKILL.md must retain required metadata, untouched sections/assets, and both origins; ambiguity/conflicting edits block the plan. | Heading editing needs a declared transform, not ordinary last-wins file copying. Defer that transform. |
| R7: source exclusion | Exclude review from the repo contribution. Team review remains. | Existing source filters cannot silently become lower-layer deletion. |
| R8: exact config import | Copy team's complete config/tool.toml to an explicit destination. Preserve it byte-for-byte; if another contribution targets that file, resolve a whole-file collision. The same whole-file rule applies to JSON/YAML imports. | General AST/key merging remains outside the first cut. |
| R9: mixed result | Produce a copied review bundle, rendered AGENTS.md, copied guide, and opaque config under a bounded export directory. One final manifest owns each destination once. | One operation inherited from the output directory cannot express the result. |
| R10: same heading/title | Two files contain Procedure, or one file repeats Examples. IDs remain qualified by source, file/view, and declared selector; ambiguous selection blocks rather than picking the first match. | Heading title alone is not universal identity or an override slot. |
| R11: global/local reuse | Select global and repo project guidance explicitly. Plan reports both contexts and declared order. No ambient host files join the build; exporting flattened text warns about potential duplicate host loading. | Source scope cannot be assumed to be build precedence. |
| R12: target ownership | An opaque copy and a render both target AGENTS.md. Resolve before writes; independent output writers cannot both own it. Direct edits still require the existing force workflow. | Mixed composition must share destination ownership/state, not bypass safety. |

R4 and R6 are limits, not successful prototype features: they expose operations
that require more specification. R11 can identify a flattening risk from declared
profiles, not prove what an arbitrary host will load at runtime.

#### Results by model

- **Filesystem-only** handles R1/R5/R8 but cannot express R2/R3 without gaining
  a document adapter. Keep it as the simple import path, not the universal model.
- **Heading-only** handles prose reuse but loses complete bundles, opaque files,
  and frontmatter boundaries in R1/R8. Reject as the universal model.
- **One strict parent tree** can display R2 but cannot give one physical file
  several logical document placements without duplication. Do not make display
  navigation enforce exclusive storage/document ownership.
- **Output-only typing** fails R9. The destination directory does not determine
  whether each contained file is copied or generated.
- **Typed physical inventory plus document views** handles the distinguishing
  cases while diagnosing ambiguous/dependency-bearing requests. This is the
  surviving recommendation. Expose a tree-shaped browser, not a generic graph
  language; references and views carry the few additional relationships.

#### Narrowed operating contract

Recommend the following for owner review:

1. A directory contains files/directories. A Markdown file may expose an optional
   heading view; a heading contains subordinate document contributions. Roles
   and bundle boundaries remain separate from structural kind.
2. Copy means preserve the whole selected physical payload. Rendering from a
   document view creates a separate file and never mutates the source or an
   installed skill implicitly.
3. Selecting a skill root includes its entire bundle. A partial installed skill
   should be diagnosed unless explicitly requested as a derivative. Browsing or
   extracting a heading is not partial skill installation.
4. Permit mixed copy/render contributions with one final output manifest.
   Use explicit destinations and report selection, operation, origin, and
   replacement. Preserve whole-file opaque config imports.
5. Default to unresolved-collision diagnostics in the first cut; support
   explicitly chosen whole-file/bundle replacement. This remains a recommendation
   against the October review's implicit later-wins choice, not an approved change.
6. Preserve source `exclude` semantics. Lower-layer deletion must declare its
   scope. Keep context roots explicit and project output containment intact.

An imported Markdown adapter needs a parsing/address contract before R2/R3 become
implementation acceptance tests: frontmatter, heading body boundaries, ATX and
Setext headings, fenced/indented code, HTML blocks, repeated headings, introductory
text, rebasing heading levels, and provenance back to the physical file. Prefer
an established Markdown parser over a second invented regex grammar. Stable
authored IDs or unambiguous qualified heading paths remain candidate addresses;
source updates must diagnose invalidated/ambiguous selectors. This does not
authorize rewriting frontmatter or extracting arbitrary configuration ASTs.

#### First implementation versus backlog

**Core candidate:** physical import/copy, explicit authored or parsed document
views, independent copy/render operations, whole-file/bundle replacement,
source exclusions, provenance, and one managed output manifest. Deliver adapters
in slices; until an imported heading adapter exists, SKILL.md can still be copied
intact but its headings cannot be advertised as selectable by current Mogent.

**Backlog now:** vendor-specific rules, custom-agent/profile conversion, command
and prompt integrations, registered hooks/settings automation, automatic global
installation/projection, versioned agentsstandard cascade emulation, and editable
heading derivatives. Arbitrary native files remain eligible for an explicit
opaque copy; backlogging an adapter does not prohibit copying its whole file.
Semantic TOML/YAML/JSON merging stays deferred as already requested.

#### Remaining owner choices after execution

The broad field is reduced to four decision groups rather than a new list of
loosely defined features:

- **DF-R1 Model and first cut:** accept physical `dir`/`file` inventory with
  optional heading views and explicit operations? Should parsed SKILL.md headings
  be in the first delivery slice, or follow the existing authored-fragment adapter?
- **DF-R2 Collision granularity:** accept explicit whole-file/bundle replacement
  with diagnostic-by-default collisions, or automatic later-wins? Keep directory
  union, bundle replacement, and future heading editing distinct.
- **DF-R3 Selection/placement syntax:** settle candidate names/blocks after the
  above semantics; avoid inferring destinations from heading titles or introducing
  a second selector just for skills.
- **DF-R4 Exclusion scope:** keep source-filter lists and add explicit earlier
  scope, or migrate to scoped blocks? Never silently reinterpret existing lists.

The owner's approval to backlog secondary integrations is recorded as a scope
direction, not approval of the above recommendation or illustrative HCL grammar.

### 2026-10-09 — Owner approval: typed views and explicit operations

After the tabletop results were presented, the owner asked to document them and
approved the recommendation. Record the accepted direction in
[DR-juvih](../../DR/DR-juvih-typed-content-composition.md) and DI-juvih in
`TODO/TODO.md`:

- Typed physical directory/file inventory with optional heading views.
- Intact copying and explicit rendering; extraction does not implicitly edit
  an installed skill.
- Mixed artifacts with one final manifest and explicit replacement for
  collisions, rather than automatic later-wins.
- Whole-file opaque configuration imports and the secondary-integration backlog.

DF-R1's model choice and DF-R2's collision default are settled at that level.
DF-R1's adapter delivery order, DF-R2's exact granularity/syntax, DF-R3's grammar,
and DF-R4's exclusion-scope syntax still need specification. The approval does
not adopt the illustrative HCL or claim current imported-heading support.

### 2026-10-09 — First delivery: approved additive grammar and API-first implementation

The owner subsequently requested implementation with atomic commits, following
API/unit tests before CLI/usage tests, and explicitly approved the proposed
first-delivery interface. DI-juvih-first-slice in TODO/TODO.md records that answer.

Delivered: named physical `content` directory roots, output `dir-tree`, explicit
`node`/`operation`/`into`, optional `heading` title arrays, whole-subtree `replace`,
and rendered `append`. Existing authored `select` objects and legacy output kinds
remain supported. Goldmark-backed heading views preserve physical payload identity;
copying stays opaque and complete, while rendering creates a separate derivative.
Captured manifests stage under one output owner; permissions and empty directories
are state tracked alongside bytes. Plans show operations, source context,
replacement history, member changes, and generated Markdown diffs. Update previews
compare old/new directory manifests and still write only pins on acceptance.

Core unit tests preceded CLI usage tests. The first-delivery cases cover intact
skills, extracted subtrees, ambiguity/fenced/quoted/Setext handling, opaque configs,
explicit replacements, link dependencies, source filtering, drift, rollback, and
update review. The runnable typed demo and exact grammar are in docs/TYPED-CONTENT.md.

Generalized nested selectors and typed sidecar declarations, custom annotations,
scoped earlier-layer removal, stable authored heading IDs, editing derivatives,
vendor adapters, and global installation remain follow-up items. No retrospective
approval of the earlier illustrative grammar is implied.
