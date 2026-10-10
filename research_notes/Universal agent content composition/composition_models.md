# Universal agent content composition: model alternatives for Mogent

Analytical notes, 2026-10-09. Provisional recommendations only; no schema, decision intent, or implementation is approved here. This analysis reads the local design and implementation independently of the web-standard research assigned elsewhere. The global/project `.agents` layouts discussed below are design scenarios supplied by the assignment, not verified standards requirements. Incorporate the separate **agentsstandard.com findings** before resolving discovery, scope, naming, or interoperability decisions.

## 1. What should be universal: one typed tree, physical/document models, or semantic collections?

### Takeaway

The strongest recommendation is **one selection and composition vocabulary over two domain-aware adapters: physical directory/file content and structural document content**. A common typed inventory or plan can expose `dir`, `file`, and `heading`, but those types must retain different identity, containment, and materialization rules; bundle/profile semantics belong on a separate axis.

### Cited Findings

- TE-jofit already separates selection, placement/composition, and materialization, and recommends authored and directory-import adapters rather than filesystem-only hierarchy or compulsory per-file declarations. Its names and grammar remain proposals. — [TE-jofit, lines 18–21, 40–48, 76–110, 344–380](../../docs/thought-experiments/TE-jofit-unified-content-composition.md)
- The current library has distinct `Sections`/`ByPath` and `Trees` inventories. `TreeEntry` annotates subdirectories for discovery and explicitly does not alter copying. A section has logical name/path, title, source, metadata, and children. — [library.go, lines 15–64](../../v2/library.go)
- Current section loading prohibits both a source and children, and requires one or the other. It assigns logical paths from section names, independently of source-file paths. Thus today's `Section` is an authored hierarchy node, not a parsed Markdown heading AST. — [library.go, lines 238–295](../../v2/library.go)
- The current renderer validates a source's line-one `# Title`, emits sidecar headings at structural depth, and includes the remaining body. It does not discover/select the body's nested headings. — [render.go, lines 38–94](../../v2/render.go)
- Named logical roots and inherited base directories were proposed specifically to separate public content identity from repository layout. — [October review, lines 3–35](../../docs/reviews/2026-10-05.md)
- The follow-up roadmap leaves logical identity, physical/output paths, layer order, collisions, and migration pending design review. — [TODO-tugur, lines 59–105](../../TODO/TODO-tugur-v2-composition-followup.md)

### Inferences

#### Distinguish three axes, rather than one increasingly overloaded type enum

| Axis | Examples | What it determines | What it does not determine |
|---|---|---|---|
| Physical storage/object | Directory; file with bytes, relative path, permissions | Copy boundary, path containment, source location | Whether content is a skill, instruction, guide, or research |
| Document structure | Authored heading/group; explicitly imported document section | Render order, heading depth, logical section identity | A filesystem directory or output filename |
| Semantic grouping/policy | Skill bundle; research bundle; agent-document role; project profile | Selection unit, completeness expectations, dependencies, target applicability | A new operating-system object or automatic merge algorithm |

- **`dir`** should mean a physical directory, not “a group that happens to have children.” **`file`** should mean an opaque byte-bearing object, even when its extension is `.md`. **`heading`** should mean a structural document node, not “a directory rendered as a title.” These are recommended interpretations, not chosen syntax.
- **Bundle** means a set of related payloads selected/materialized with a coherence policy. A skill may be rooted at a directory, but “every directory is a skill” and “every bundle is a directory” are both unnecessary restrictions. A research note plus its data and figures can also be a bundle.
- **Profile** means a named reusable selection/composition recipe. It can span disconnected directories, several libraries, agent docs, guides, and skills. It need not own payload bytes. It is better modeled as references/recipe metadata than by reparenting all referenced objects under a fabricated directory.
- Tags, semantic roles, scope applicability, and atomicity are independent metadata. A file can be Markdown + an agent document + project-applicable + part of a bundle. Making these mutually exclusive node types creates needless type combinations.
- A universal surface need not imply every content node supports every operation. Selecting a binary asset can be valid while rendering it into instruction prose is invalid. Capability checks should follow selection, before materialization.

#### Model alternatives

| Alternative | Strongest advantage | Main failure mode | Assessment |
|---|---|---|---|
| Markdown output plus untyped raw directory overlays | Small extension of v2; opaque copying is useful across formats | Structural selection remains separate; bundle integrity is invisible | Useful staged route, insufficient universal explanation |
| Filesystem is the universal tree | Automatic discovery, little authoring overhead, straightforward copying | Physical moves change logical references; folder depth becomes heading depth | Strong for file distribution, weak for reusable prose |
| One strictly declared `dir/file/heading` ownership tree | One traversal/selector; explicit source kinds | Requires duplicate file registry or forces logical and physical parentage to coincide | Viable only with imports, references, and precise edge rules |
| One typed facade/plan over physical and document adapters | Shared selection/provenance; preserves each domain's invariants | Adapter rules and projections require explicit contracts | **Strongest provisional recommendation** |
| Independent physical and document products/selectors | Domain semantics remain simple; lower migration cost | Defaults, exclusions, metadata, browsing, and diagnostics can diverge | Credible incremental strategy; avoid gratuitously divergent vocabulary |

The third and fourth alternatives are not necessarily competitors. A discriminated union in the resolved inventory can be a “one typed tree” implementation of adapter-backed domains. The material distinction is whether unification erases domain rules, not how many structs exist.

#### Why a literal ownership tree may be too strong

Suppose `guides/review.md` supplies body text for logical `workflow/review`, is copied whole into `.agents/guides/review.md`, and is referenced by two profiles. The physical file has one storage parent; the logical document contribution has a different structural parent; both profiles reference it. One parent pointer cannot express all three relationships honestly.

Prefer distinct relationships:

1. Physical containment: directory → directory/file.
2. Document hierarchy: heading/group → heading/contribution.
3. Payload reference: document contribution → file or explicitly addressed fragment.
4. Semantic membership/dependency: bundle/profile → referenced nodes.

A per-root tree can remain the user-facing browse/selection view, while cross-view links remain explicit references. This is not a recommendation for a general graph query language or config AST engine.

#### What a shared inventory needs to preserve

Conceptual information, **not proposed HCL or implementation fields**:

- Stable logical identity qualified by library/root, plus physical source locator where applicable.
- Node kind/capabilities, display title distinct from identity, ordered document children versus deterministic imported path order.
- Original byte payload reference and, only when needed, an explicit structural/document view.
- Metadata and curation/defaults without requiring declarations for every imported script or asset.
- Bundle boundary and dependency references where declared.
- Origin scope, source alias/pin, selected node, destination projection, and operation provenance.

The selection surface should distinguish traversal relationships. Selecting a physical file broadly should not automatically mean selecting a derived heading view and losing whole-file preservation. Selecting one document section should not automatically select unrelated files next to its source.

#### Strong counterexamples to this recommendation

- **Opaque distribution only:** A library ships scripts, YAML/TOML/JSON, binary assets, and complete Markdown files; nobody wants document composition. A physical tree with bundle annotations is sufficient. Document adapters add no value to that workload.
- **Prose-only library:** All sources intentionally satisfy v2's first-heading fragment contract and are rendered into one AGENTS.md. A simple authored section model is easier than exposing physical nodes everywhere.
- **Naturally nested authored package:** An author genuinely wants a directory containing a generated Markdown file whose children are headings. A well-defined `dir → file → heading` facade may be clearer than visible dual models. The recommendation should permit that presentation, while the file's generated/preserved status remains explicit.
- **Already coherent logical tree:** An authored hierarchy intentionally mirrors folders and never reuses content across parents. A single tree can work well; dual relationships should not require extra declarations in this simple case.
- **Arbitrary heading extraction:** A user needs one subsection from an existing frontmatter-bearing guide without splitting files. Current v2 sections are not enough. An explicit document-import/fragment-addressing adapter would be needed, or the author must supply a fragment. Merely adding `kind = heading` does not solve parsing or stable addressing.

### Gaps

- Whether `heading` denotes only authored structure, imported Markdown structure, or both is unresolved. An imported heading view needs rules for introductory text, nested body headings, fenced code, duplicate headings, frontmatter, and edits that invalidate addresses.
- Whether semantic bundles/profiles need first-class declarations now or can begin as annotation/recipe conventions is an owner/product question. The code does not currently implement these abstractions.
- Discovery conventions and portable semantics for skills/agent docs across `.agents` layouts must incorporate the separately researched agentsstandard.com evidence.

## 2. How should rendering, exact copying, collisions, overrides, and bundle dependencies interact?

### Takeaway

Treat **rendering and copying as explicit payload operations**, independent of whether the final target is a file or directory. Favor whole-bundle preservation/replacement for skills and whole-file copying for YAML/TOML/JSON; permit structural prose overrides only through stable, explicitly mapped logical identities, never heading-title coincidence or general configuration AST merges.

### Cited Findings

- Current Markdown rendering strips the first source heading, converts CRLF to LF, trims the body and final result, emits structural headings, and appends a final newline. A frontmatter-first Markdown file fails the first-line heading check. — [render.go, lines 29–33, 49–65, 85–94](../../v2/render.go)
- Directory copying reads file bytes and writes them using source file permission bits. It creates directories with `0755`; this is byte/file-permission preservation, not preservation of all directory metadata, timestamps, extended attributes, or symlinks. — [tree.go, lines 220–269](../../v2/tree.go)
- The planner currently allows exactly one source for a tree output and requires `select = { all = true }`. Markdown sources resolve separately and are rendered in source order. — [plan.go, lines 60–146](../../v2/plan.go); [render.go, lines 13–28](../../v2/render.go)
- The spec says matching titles do not merge headings across libraries; library section paths must be unique. — [v2 spec, lines 110–121, 234–237, 392–402](../../docs/proposals/MOGENT-HCL-V2-SPEC.md)
- The October review recommends recursive opaque overlays and automatic later-wins, including type changes. TE-jofit instead recommends error-by-default with explicit later-wins and explicit type replacement; the roadmap records these as pending decisions. — [October review, lines 96–143](../../docs/reviews/2026-10-05.md); [TE-jofit, lines 267–278, 369–377](../../docs/thought-experiments/TE-jofit-unified-content-composition.md); [TODO-tugur, lines 64–73](../../TODO/TODO-tugur-v2-composition-followup.md)
- TE-jofit demonstrates that recursively merging a replacement skill retaining only SKILL.md leaves an old supporting script behind, whereas whole-directory replacement removes it. — [TE-jofit, lines 267–272](../../docs/thought-experiments/TE-jofit-unified-content-composition.md)
- Current source exclusions filter exact paths or descendant prefixes; they do not remove earlier contributions. The code verifies source entries before filtering. — [tree.go, lines 45–58, 297–313](../../v2/tree.go); [TE-jofit, lines 280–314](../../docs/thought-experiments/TE-jofit-unified-content-composition.md)
- Current apply stages tree replacements and integrates rollback/state with file writes; update comparison uses rendered output strings and does not inspect directory source manifests. — [tree.go, lines 157–218](../../v2/tree.go); [build.go, lines 158–215, 285–315](../../v2/build.go)

### Inferences

#### Resolve destination artifacts before choosing how to produce their bytes

“Markdown versus directory output” mixes content processing with target shape. A directory can contain both rendered documents and byte-preserved Markdown, configs, scripts, and assets. A single file can likewise be copied exactly or generated by a renderer.

Recommended conceptual pipeline:

`discover/adapt → select → account for bundles/dependencies → place/compose → materialize payloads → final path manifest → transactional installation`

This is an analytical decomposition, not a locked phase order: structural composition must precede rendering, while final path conflicts must be caught before installation.

| Selected payload | Operation | Result | Important boundary |
|---|---|---|---|
| Complete Markdown file, including frontmatter | Copy/preserve | Identical bytes at mapped path | Do not strip, rewrite headings, inject notices, or serialize frontmatter |
| Authored guidance hierarchy/fragments | Render | New Markdown artifact | Reheading/normalization must be part of the declared renderer contract |
| YAML/TOML/JSON file | Copy/replace whole file | Identical selected winner's bytes | No key-level merges, schema AST overlays, or generated parse/serialize round trips |
| Script/binary/asset | Copy/preserve | File bytes and specified permission behavior | Not an implicit Markdown contribution |
| Skill bundle | Preserve coherent member set | Complete directory/file set | Internal generation or granular alteration requires explicit derivative policy |
| Profile | Resolve recipe | Selection/placement contributions | Not intrinsically a renderable body or copyable file |

A Markdown extension should enable an optional document adapter, not force rendering. A structural heading should require a renderer/materialization mapping, not force a destination path guessed from its title. This broadens TE-jofit's recommendation: output shape alone cannot decide payload handling.

#### Frontmatter and exactness

- Copying SKILL.md, an agent descriptor, or a guide should preserve frontmatter along with the body exactly. Reading metadata for discovery can be observational; it does not authorize rewriting bytes.
- Rendering fragments from frontmatter-bearing files needs an explicit contract for whether frontmatter is rejected, omitted, retained in a designated output envelope, or otherwise handled. Combining several YAML headers into a valid new document is not an automatic side effect of selection.
- Strong initial scope: preserve frontmatter-bearing complete documents; render explicitly authored guidance fragments; defer arbitrary frontmatter composition. An independently authored complete replacement is still allowed.
- Exact-byte copying does not mean the entire copied package remains semantically unchanged after renaming/mounting. Relative paths inside files can depend on placement; preserve bundle-relative layout and report transformations separately.
- Semantic heading overrides can coexist with preservation of untouched spans, but that requires a source-span-preserving document adapter. Current render.go supplies neither arbitrary heading parsing nor byte-round-trip guarantees. Do not promise exact whole-file bytes for a structurally edited document.

#### Distinguish source identity from destination collisions

1. **Duplicate source declaration identity:** Two nodes with the same qualified identity in one root are invalid. The same physical file may be referenced intentionally by distinct logical contributions; that is not necessarily a duplicate declaration.
2. **Duplicate display title:** Two headings titled “Testing” may be valid under different parents or even repeat intentionally. Display titles are not override keys.
3. **Duplicate destination path:** Two contributions map to `.agents/guides/review.md`. Resolve by a chosen path collision policy, irrespective of whether source IDs differ.
4. **Explicit shared document slot:** Two sources are deliberately mapped to the same stable logical `workflow/testing` slot. This permits a defined replacement/composition policy before rendering; identical source labels alone do not imply shared identity.
5. **Ancestor/type collision:** A contribution supplies file `skills/review`, another supplies `skills/review/SKILL.md`. This is a file/subtree incompatibility, not ordinary file/file shadowing.

For imported headings, title-derived slugs and occurrence numbers are fragile identities: adding an earlier “Examples” shifts the second occurrence. If stable authored IDs are absent, an importer should surface ambiguous addresses rather than silently replace the wrong heading. Exact parsing/address syntax remains open.

#### Override granularity ladder

| Unit | Plausible operation | Example | Tradeoff |
|---|---|---|---|
| Source contribution | Filter selection | Do not take repo's experimental guide | Earlier sources remain unaffected |
| Destination file | Replace whole file | Local `settings.toml` wins over company file | Simple, preserves winner's bytes; loses earlier file entirely |
| Physical directory | Recursive union or replace subtree | Merge guide directories versus replace review bundle | These are different operations; layer order alone cannot choose |
| Logical document section | Append or replace explicit slot/subtree | Replace company testing guidance | Requires stable identity, order, and subtree semantics |
| Semantic bundle | Replace coherent member set | Select repo's complete review skill | Avoids stale assets and mixed versions |
| Earlier destination/logical slot | Explicit removal | Suppress inherited skill/section | Distinct from source filtering and from absence in a later source |

Prefer conflict diagnostics until an override policy is declared, particularly for type changes and atomic bundles. This favors TE-jofit's cautious recommendation over the October review's blanket later-wins policy; it remains an analytical preference, not a decision. Automatic later-wins remains a compelling counterexample for intentionally authored ordered distribution overlays.

Equal bytes can avoid a content diff, but two origins can still matter for future updates, permissions, bundle boundaries, and explanation. Plan should retain both contributions and explain the effective ownership/winner, even if a collision policy deduplicates equal payloads.

#### Skill integrity: four meanings of “atomic”

- **Selection atomicity:** Taking a skill takes SKILL.md and its associated scripts/references/assets, not only the matched metadata leaf.
- **Replacement atomicity:** Choosing a newer/different bundle replaces the earlier member set so obsolete members do not survive by recursive union.
- **Version coherence:** A bundle's members should come from a coherent pinned snapshot unless an explicit derivative composition is intended.
- **Installation atomicity:** A failure must restore the previous output and state. This is a transaction property, distinct from the first three.

The current transaction engine addresses installation rollback; it does not establish the other meanings. Even a directory-shaped skill can depend on content outside that subtree or on external executables, so “copy this directory” is not a universal dependency solution.

For declared dependencies, distinguish:

- **Contained members:** Included by the bundle boundary without per-file registration.
- **Other content bundles/files:** Referenced by qualified identity; chosen selection must account for them without silently making every adjacent file a dependency.
- **External runtime prerequisites:** A command, service, or package requirement can be reported as metadata; Mogent need not become an environment installer.

A first coherent approach could expand explicitly declared content dependencies and show why each was added, or diagnose missing requirements and ask for explicit selections. Which behavior best fits curation is unresolved. In either case, replacing/excluding a depended-on object must reconsider dependent bundles; do not automatically resurrect a deliberately removed requirement. Multiple incompatible dependency versions and cycles require a defined diagnostic/closure policy, not an assumed package solver.

Granular overrides are a legitimate counterexample to whole-bundle replacement: a team may intentionally patch one script in an otherwise unchanged skill. That should be an explicit derivative bundle or explicit non-atomic composition policy, with inherited member provenance visible. Atomic preservation should be a coherent default/capability for declared skills, not a prohibition on intentional derivatives.

#### Tabletop failure cases

- Company `review` has SKILL.md + `scripts/old.sh`; repo `review` has a new SKILL.md only. Recursive union keeps the obsolete script. Atomic replacement removes it, but may leave the new skill incomplete if it actually expected inherited assets. Intent/completeness must be explicit.
- A tag query matches only SKILL.md. Copying that leaf alone can break relative references. If it identifies a declared bundle, promote selection to the bundle boundary and explain that promotion.
- Repo replaces only `workflow/testing` in a rendered document. Replacing AGENTS.md wholesale is too coarse; appending another “Testing” heading does not implement an override.
- Project supplies `config.json` containing only one key. Treating it as a patch would violate whole-file semantics. It is a complete replacement; authors must provide the complete desired file.
- Two independent guides both contain “Examples.” Title matching merges unrelated prose. Qualified identity/explicit slot mapping avoids that.
- Selecting a heading fragment while also copying the same file to the same destination creates two writers. A single final artifact declaration must resolve this before filesystem writes.

### Gaps

- Bundle discovery, boundaries, dependency metadata, and completeness checks need concrete examples and later standards evidence; this analysis does not claim arbitrary script dependencies can be inferred.
- Collision defaults, file/directory replacement, heading-slot overrides, and explicit derivative semantics remain owner decisions. The existing review and TE differ and should not be presented as one settled contract.
- Repeated-heading addressing, body heading rebasing, link rewriting, and frontmatter behavior for structural rendering need a focused follow-up before promising arbitrary document imports.
- Directory update review remains incomplete in current v2; a universal design should compare final manifests and bundle membership, not only rendered strings.

## 3. How do global/project scope, Mogent layer order, and import/export projections fit together?

### Takeaway

Keep **where content lives and applies** separate from **the order Mogent composes contributions into a target**. Treat `.agents` and other tool paths as explicit discovery/materialization projections with honest lossiness reporting, while preserving arbitrary files whole and leaving consumer runtime precedence to later verified evidence.

### Cited Findings

- Current sources can be local directories or pinned Git libraries; local paths are relative to mogent.hcl, and output paths are validated as relative paths. Planned filesystem targets are rooted at the configuration directory. There is no implemented implicit global/project context resolver in these paths. — [config.go, lines 154–197, 232–263, 300–309](../../v2/config.go); [build.go, lines 67–81](../../v2/build.go)
- The source model distinguishes Git `ref`, pinned `commit`, and unpinned local content. Scope is not an attribute of the present `Source` or `Output` structs. — [config.go, lines 26–52](../../v2/config.go)
- Current Markdown sources contribute in declared order without title-based merging. Multiple directory mounts/layers are proposed, not implemented. — [v2 spec, lines 110–121](../../docs/proposals/MOGENT-HCL-V2-SPEC.md); [plan.go, lines 64–85](../../v2/plan.go); [TODO-tugur, lines 95–102](../../TODO/TODO-tugur-v2-composition-followup.md)
- Current `paths` means identical independently managed outputs, not symlinks or tool-specific translation. — [v2 spec, lines 110–111](../../docs/proposals/MOGENT-HCL-V2-SPEC.md)
- TE-jofit distinguishes destination `into` placement from source lookup, and retains direct-edit protection, complete-output state, and non-overlapping managed targets. — [TE-jofit, lines 239–243, 330–342](../../docs/thought-experiments/TE-jofit-unified-content-composition.md)

### Inferences

#### Four independent kinds of context/order

| Dimension | Question | Example |
|---|---|---|
| Origin/discovery scope | Where was the source found? | User-global content root versus repository-local root |
| Applicability context | Where is it intended to be used? | All projects; one repository; one subtree or task |
| Build composition order | Which contribution wins at a declared destination/slot? | Organization → personal → project, explicitly configured |
| Consumer runtime behavior | What does a receiving tool load and prioritize? | It may discover both a global and a project document |

- “Global” does not mean an intrinsically lower Mogent layer, and “project” does not mean universally higher runtime instruction authority. Both conclusions require explicit policy/evidence.
- A company baseline can be the last enforcing build layer even when a personal source lives globally. Conversely, a global library can be used solely to generate project-local output. Origin is provenance, not precedence.
- If a receiving tool reads both global and project outputs, flattening all global content into the project artifact can duplicate it at runtime. If it reads only project paths, flattening may be the desired export. This is target-aware projection policy, not one universal default.
- Within-document order is another dimension: later prose is not automatically a replacement or an instruction-precedence mechanism. Structural slot replacement must be explicit if intended.
- Global/project independence favors explicit context roots and target ownership boundaries. It does not justify relaxing current safe relative path validation to permit arbitrary `../` or absolute outputs. A future global target would need its own defined root and state lifecycle; current v2 does not provide that simply by naming `.agents`.

#### Illustrative `.agents` content, without asserting a standard layout

```text
<global-context>/.agents/           <project-context>/.agents/
  skills/review/                     skills/project-review/
    SKILL.md                         agents/reviewer.md
    scripts/check.sh                 guides/testing.md
  agents/reviewer.md                 research/design-notes.md
  guides/workflow.md                 config/tool-settings.toml
```

This example deliberately mixes semantic roles and opaque file formats. Discovery can inventory physical dirs/files across either context. An authored document view may reuse workflow/testing fragments. A profile may choose a skill plus a guide and an agent descriptor across both roots. None of these relationships requires “everything is rendered Markdown” or “a directory automatically creates a heading.” The exact subdirectory names remain hypothetical until the separate findings are incorporated.

#### Import and export are projections, not declarations of universal tool equivalence

**Import projection:** Map a known source layout to physical inventory plus optional semantic annotations/document views. Keep source identity and raw bytes; mark inferred bundle/role metadata as inferred. Do not manufacture missing scripts, silently discard unknown files, or assume every Markdown document satisfies the v2 fragment contract.

**Export projection:** Map selected content to target paths and declared artifact operations. Some exports preserve whole files/bundles; some explicitly render structural guidance. A target profile can describe supported roles and path placements without changing universal node types.

Illustrative tool paths such as `.claude/skills`, `.claude/agents`, `.codex/skills`, AGENTS.md, or CLAUDE.md are candidate export/import destinations only in this analysis; their supported scopes, formats, and behavior have not been verified here. Do not infer interchangeability merely from similar names.

Recommended boundaries:

- **Path/layout adaptation:** Rename/mount an intact bundle or copy a whole document to a declared target when its format is compatible.
- **Instruction rendering:** Produce a target document from explicitly selected structural contributions under a specified renderer contract.
- **Metadata conversion, if eventually needed:** Make it a narrowly named format conversion with its own constraints and preservation/loss report. It is not implied by “copy this Markdown file” and is not proposed as part of this first scope.
- **Unsupported content:** Reject a requested incompatible projection or leave it unselected with explicit explanation. A guide can be retained as an opaque file in a general content target even if a particular receiving tool does not consume it automatically.
- **Configuration files:** Copy the entire selected YAML/TOML/JSON file, or replace it as a whole. No general key/AST merge or tool-config synchronization engine.
- **Ownership:** Build one final artifact manifest per target. A directory managed wholesale and an independently managed rendered file inside it need a coordinated single writer, not overlapping independent outputs. Existing state/rollback is a useful base, not a complete multi-target projection design.

No lossless export claim should be made when logical IDs, curation, profile membership, dependency information, origin scope, or per-section provenance have no representation in the target format. Those can remain in Mogent's authored configuration/state/provenance; they should not be injected into byte-preserved files. Likewise, re-importing a generated AGENTS.md cannot reconstruct its full source graph without additional provenance. “Round trip” must specify whether it promises bytes, layout, metadata, or semantics.

#### Strongest recommendation, stated as an owner-review hypothesis

1. Retain TE-jofit's shared selection direction, but sharpen it to **physical/document adapters with a typed common view**, rather than an unqualified single ownership tree.
2. Treat directory/file as storage kinds, heading as document structure, and bundle/profile as separate semantic policy/recipes. Avoid mandatory per-file declarations.
3. Separate target artifact shape from payload operation: directory manifests may contain exact copies and explicit rendered documents; standalone files may also be copied whole.
4. Prefer atomic preservation/replacement for declared skill bundles, with visible dependency accounting and explicitly intentional derivative overrides.
5. Use qualified identities and explicit destination/slot mapping; never override solely by heading title. Keep source filtering, earlier-layer removal, recursive union, and subtree/bundle replacement distinct.
6. Keep origin scope, applicability, build layer order, document order, and consumer runtime precedence separate. Support interoperability through bounded import/export projections, not config AST merges.

This recommendation is strongest because it accommodates skills, research, agent documents, guides, arbitrary opaque files, and logical prose reuse without making file paths into heading identities or making semantic roles into storage types. Its main cost is explicit adapter/projection contracts; a physical-only subset is the stronger choice when the product intentionally only distributes complete files.

#### Counterexamples that should change or narrow the recommendation

- If verified target formats require exact files everywhere and no user wants structural prose selection, choose physical composition plus bundle policies; retain rendering as a separate optional feature.
- If target guidance is always authored from independent fragments and no intact document/bundle distribution is needed, the existing document model can remain the primary interface.
- If users intentionally rely on per-file later-wins overlays to customize a skill, rigid atomic replacement would be counterproductive. Offer explicit derivative composition rather than silently treating all directories as immutable bundles.
- If a target requires transformed frontmatter, path copying alone cannot produce a compatible export. That warrants a separately reviewed conversion adapter or an incompatible-projection diagnostic, not silent normalization.
- If separate global and project outputs are both loaded by the consumer, materializing a fully flattened project profile can duplicate or alter runtime context. Export them separately or flatten only under an explicit target policy.

### Gaps

- **Required follow-up:** Incorporate agentsstandard.com findings and other assigned standards research before deciding canonical `.agents` paths, scope resolution, bundle discovery, heading conventions, or tool export compatibility. No external standards claims were researched or verified in this note.
- Validate target-specific discovery, runtime precedence, format/frontmatter compatibility, relative-path assumptions, and global/project behavior with concrete tools. Build layer order cannot substitute for that evidence.
- Owner review should decide the universal facade and semantic boundaries before locking names, grammar, override defaults, dependency expansion, or migration. These notes intentionally do not record a decision or amend the TE/TODO/spec/code.
