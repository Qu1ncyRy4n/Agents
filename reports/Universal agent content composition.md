# Typed composition preserves content across agents

**Mogent should unify selection and composition through a typed inventory backed by physical-content and document-structure adapters.** Physical `dir` and `file` nodes preserve paths and bytes; structural `heading` nodes organize explicitly renderable guidance. Roles such as instruction, skill, research, and guide, together with bundle policies, belong on separate annotation axes. Output-only `markdown`/directory typing cannot explain a directory containing copied skills, rendered instructions, and opaque configuration files: target shape and payload operation must be independent. The strongest documented interoperability is instruction `AGENTS.md` and Agent Skills bundles, including `.agents/skills` discovery in several tools; agentsstandard.com's broader directory cascade remains a distinct proposal. Global, project, and nested origins require explicit identity and target-aware projection, not an invisible host-content merge. This is a design recommendation for owner review, with naming, collision defaults, and bundle semantics still open.

## The broader directory protocol remains a proposal

The evidence assembled on **2026-10-09** separates three authorities: nbiish's Agents Standard at agentsstandard.com, the separately stewarded AGENTS.md format, and the Agent Skills specification. The first proposes a loading protocol; the others define instruction and skill formats. AGENTS.md permits ordinary Markdown with arbitrary headings and root/nested placement. Agent Skills requires an exact `SKILL.md` containing YAML frontmatter and Markdown, with required `name` and `description`, inside a directory that can include arbitrary supporting files. Its integration guide recommends `.agents/skills` alongside client-specific locations but explicitly says the format does not mandate installation paths. Support for that convention therefore does not establish adoption of the whole Agents Standard protocol. ([Agents Standard](https://agentsstandard.com/), [AGENTS.md](https://agents.md/), [Skill specification](https://agentskills.io/specification), [Integration guidance](https://agentskills.io/integrate-skills))

**The proposal's published versions conflict.** Live `agents.json` reports v1.4.0, scoped concatenation, and project `mcp-settings.json`; the pinned website repository reports v2.0.0 and unified cascade, while the pinned specification makes root `.mcp.json` the runtime MCP file and the home catalog non-runtime. V2 proposes global `~/.agents/AGENTS.md`, root `llms.txt` requirements, project `.agents/AGENTS.md`, root active `AGENTS.md`, and nested instruction files. It also specifies global/project skill replacement and a global `providers.txt` reference. These are requirements within nbiish's proposal, not verified requirements of every receiving tool. ([Live registry](https://agentsstandard.com/agents.json), [Pinned website v2 registry](https://raw.githubusercontent.com/nbiish/agentsstandard-dot-com/bb838b1853513bed0b6b5a897b407d3324737c00/agents.json), [Pinned v2 specification](https://raw.githubusercontent.com/nbiish/agents-standard/63e166a76acb61e6a197e2049ba1e87bea0ddf1a/llms.txt))

Even proposal text and implementation differ. The live site describes a launch-directory ceiling; v2 text and the pinned loader use marker-based upward root discovery. Advertised depth 1 means global-only, but the loader still loads project scope when a root exists; only folder rules are depth-gated. Its skill helper scans immediate children rather than recursively implementing the advertised `**` paths. It concatenates rule text and records bridge diagnostics without implementing the proposed agent-specific MCP merge. Mogent should identify a proposal projection by source/version rather than promise one reconstructed “standard cascade.” ([Live discovery rules](https://agentsstandard.com/#spec), [Pinned loader](https://raw.githubusercontent.com/nbiish/agents-standard/63e166a76acb61e6a197e2049ba1e87bea0ddf1a/cli/lib/loader.js), [Pinned skill helper](https://raw.githubusercontent.com/nbiish/agents-standard/63e166a76acb61e6a197e2049ba1e87bea0ddf1a/cli/lib/fs-helpers.js))

No collected source establishes universal discovery or configuration support for arbitrary `.agents/agents`, `.agents/research`, `.agents/guides`, `.agents/hooks`, or `.agents/settings.*`. These remain useful storage conventions. Research and guides are legitimate skill resources, or explicit references/configured instructions, without becoming skills themselves. The website repository's `advisory-council/RESEARCH.md` illustrates a companion resource, not a newly standardized entrypoint. ([Skill resources](https://agentskills.io/specification#optional-directories), [Research companion](https://github.com/nbiish/agentsstandard-dot-com/blob/bb838b1853513bed0b6b5a897b407d3324737c00/.agents/skills/advisory-council/RESEARCH.md), [OpenCode configured instructions](https://opencode.ai/docs/rules/))

## Typed nodes unify selection without erasing boundaries

The earlier [TE-jofit](../docs/thought-experiments/TE-jofit-unified-content-composition.md) is **local design analysis needing refinement, not an ecosystem finding or decided contract**. Its separation of selection, placement/composition, and materialization remains useful. The refinement is to distinguish node kind from content role, and target shape from the operation producing bytes. Today's v2 planner supports authored Markdown selections and single-source trees requiring `all`; its renderer consumes first-heading fragments rather than discovering an arbitrary Markdown heading AST. The proposed model extends those concepts rather than claiming current support. ([Current planner](../v2/plan.go), [Current renderer](../v2/render.go))

| Distinction | Recommended meaning | Boundary |
|---|---|---|
| `dir` | Physical containment of directories/files | No implicit Markdown, heading depth, skill status, or parser inheritance |
| `file` | Byte-bearing object with source path and permission information | `.md` does not automatically mean render; TOML/JSON remain whole files |
| `heading` | Authored document structure, or a separately specified imported document view | Owns document contributions, not arbitrary physical files |
| Role annotation | Instruction, host-specific agent profile, skill, guide, research, configuration | Does not change physical kind or automatically activate content |
| Bundle annotation | Membership/coherence policy over referenced payloads | Selection, replacement, and dependencies need their own decisions |
| Profile/recipe | Reusable selections and placements | References content rather than manufacturing ownership |

One typed browsing view can expose physical containment and document structure, but it need not pretend that every relationship is the same parent edge. A guide file can be copied intact, supply text to a logical heading, and be referenced by multiple profiles. Those are storage, payload-reference, and recipe relationships. A `dir → generated file → heading` presentation is reasonable when the file explicitly hosts a document view; **a heading cannot own an asset directory**. Nearby files are dependencies only when included by a declared bundle boundary or explicit reference.

Unified selection should use qualified root/node identity and kind distinctions, whether ultimately expressed as `node { kind = ... }`, typed blocks, or object selectors. Selecting the physical `review.md` means selecting its whole payload; selecting `workflow/review` in a document view means selecting a structural contribution. Matching display headings are not shared override identities. Curation, defaults, tags, exclusions, and diagnostics should operate consistently across adapters, with applicability checked after selection. A binary can be selected for copying without being renderable as prose. Arbitrary heading extraction requires a later parsing/addressing contract for repeated headings, frontmatter, introductory text, and fenced code; the type label alone provides none.

## Mixed outputs require explicit copy and render operations

**A directory manifest can contain both rendered and preserved files.** Conversely, a standalone file can be generated Markdown or an opaque copy. The following HCL is illustrative proposed syntax, not accepted grammar. Source aliases and typed inventories are assumed; `agent-export` is a bounded export snapshot whose internal paths describe a future project layout, not an instruction to manage the entire repository.

```hcl
output "project-package" {
  path = "agent-export"
  kind = "dir" # naming remains open

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

The result includes the complete review skill, rendered root `AGENTS.md`, a TOML storage payload, and a whole JSON settings file. The TOML storage location has **no claimed automatic Codex discovery**; its native global target would require a separately authorized context projection. Claude's settings path is documented when placed in an actual project, whereas the export snapshot itself activates nothing. The expanded [TE-gunak](../docs/thought-experiments/TE-gunak-typed-agent-content-composition.md) includes research and agent-profile copies and concrete typed declarations. ([Codex configuration payload](https://developers.openai.com/codex/skills), [Claude settings](https://code.claude.com/docs/en/settings))

Copying preserves complete bytes, including opening frontmatter, active placeholders, sidecars, and bundle-relative layout; provenance belongs outside preserved payloads. Rendering explicitly authored fragments can rehead and normalize according to its declared contract. **Semantic TOML/YAML/JSON merges are deferred: every selected configuration is a complete file, not a key patch.** Native OpenCode configuration merging or Claude settings precedence remains receiving-tool behavior and must not become an invisible Mogent merge. Exact byte preservation also does not promise timestamps, extended attributes, or semantic portability after relocation. Today's copier preserves file permission bits, creates directories with `0755`, and rejects source symlinks; a future symlink import policy needs review despite host support for some linked skills. ([Frontmatter sensitivity](https://code.claude.com/docs/en/skills#frontmatter-reference), [OpenCode native config](https://opencode.ai/docs/config/), [Current copier](../v2/tree.go))

Bundle atomicity needs four separate meanings: selection completeness, coherent source version, whole-member-set replacement, and transactional installation. If company review contains `scripts/old.sh` and the replacement contains only `SKILL.md`, per-file overlay retains the script while bundle replacement removes it. Either can be intentional; the former needs an explicit derivative policy rather than an unannounced mixed-version skill. Ordered sources alone cannot choose between these operations. File collisions, file/subtree conflicts, source filtering, and removal of earlier contributions similarly need distinct treatment. The October review favors automatic later-wins, while TE-jofit prefers explicit conflict resolution; **the owner has not settled that difference**. ([October review](../docs/reviews/2026-10-05.md), [Earlier bundle analysis](../docs/thought-experiments/TE-jofit-unified-content-composition.md))

Plan provenance must show source alias/pin or local snapshot, origin context, qualified selection and kind, bundle membership/dependency reasons, operation, destination, exclusions, shadowed contributions, winner, removed members, and output ownership. Equal final bytes do not erase different histories. All rendered and copied members inside one managed directory need one final manifest and coordinated writer, preserving drift detection and rollback rather than overlapping independent targets.

## Tool discovery determines useful import and export targets

These priorities reflect **documented support and integration value, not popularity rankings**. Whole skill bundles and instruction files provide the strongest first imports. Vendor rules, agents, and commands add value without pretending that all frontmatter schemas are interchangeable.

| Import priority | Worthwhile documented targets | Export restriction |
|---|---|---|
| First: skills | `.agents/skills` in Codex, OpenCode, Copilot, Cursor; native `.claude/skills`, `.opencode/skills`, `.github/skills`, `.cursor/skills`, and documented personal counterparts | Preserve whole bundle and host extensions; direct Claude `.agents/skills` discovery is not established here |
| First: instructions | Root/nested `AGENTS.md`; Codex override files; Claude instruction variants; `.github/copilot-instructions.md`; documented user instruction homes | Preserve scope and receiving-host discovery; not a universal cascade |
| Next: rules/profiles | `.claude/rules/**/*.md`, `.cursor/rules/**/*.mdc`, `.github/instructions/**/*.instructions.md`; `.claude/agents`, `.opencode/agents`, `.github/agents/*.agent.md` | Preserve suffix/frontmatter/dialect; copying a profile is not cross-host conversion |
| Compatibility: commands/prompts | `.claude/commands`, `.opencode/commands`, `.github/prompts/*.prompt.md`, legacy `.cursorrules` | Preserve active syntax and lifecycle restrictions; do not choose legacy formats as universal outputs |
| Opaque: configuration/hooks | Claude settings/hook JSON and scripts, OpenCode JSON/JSONC, Codex TOML; skill YAML sidecars travel immediately with bundles | Whole-file import only; copying hook scripts does not register hooks |

The paths above follow the researched primary documentation. Skill interoperability spans a common bundle specification, but host extensions still affect permissions and invocation. Rules differ particularly sharply: Cursor requires `.mdc`, Copilot uses `applyTo`, and Claude uses `paths`. Agent identity also differs: Claude uses frontmatter name while OpenCode uses filename. These differences justify tool-qualified role annotations rather than new physical kinds or a generic YAML-frontmatter merger. ([Codex skills](https://developers.openai.com/codex/skills), [OpenCode skills](https://opencode.ai/docs/skills/), [Copilot skills](https://docs.github.com/en/copilot/concepts/agents/about-agent-skills), [Cursor skills](https://cursor.com/docs/context/skills), [Claude skills](https://code.claude.com/docs/en/skills), [Cursor rules](https://cursor.com/docs/context/rules), [Copilot instructions](https://docs.github.com/en/copilot/how-tos/configure-custom-instructions/add-repository-instructions), [Claude agents](https://code.claude.com/docs/en/sub-agents), [OpenCode agents](https://opencode.ai/docs/agents/))

Commands and profiles also retain host-specific operational syntax and schemas; hooks require registration, not merely a copied executable. Preserve these inputs in their native dialect before considering separately reviewed conversions. ([OpenCode commands](https://opencode.ai/docs/commands/), [Copilot profiles](https://docs.github.com/en/copilot/reference/custom-agents-configuration), [Claude hook locations](https://code.claude.com/docs/en/hooks#hook-locations))

Version qualifications matter. Current Claude docs add conditional AGENTS support from **v2.1.277**, replacing the older proposal registry's symlink-only account; project/ancestor CLAUDE files suppress default fallback, and instruction loading excludes `.agents` contents. OpenCode findings describe the fetched `/docs/` series, which advertises a separate v2 site; equivalence was not verified. GitHub documents prompt-file availability, while current VS Code says Agent Host sessions no longer load them, though Local agent still does. These are host/version differences, not reasons to erase evidence. ([Claude AGENTS compatibility](https://code.claude.com/docs/en/memory#agents-md), [OpenCode documentation](https://opencode.ai/docs/config/), [GitHub prompt tutorial](https://docs.github.com/en/copilot/tutorials/customization-library/prompt-files/your-first-prompt-file), [VS Code prompt reference](https://code.visualstudio.com/docs/copilot/customization/prompt-files))

## Scope identity must precede target projection

**Origin scope, applicability, Mogent composition order, and host runtime precedence are independent.** Record user-global, repository/worktree, and nested context boundaries in identity/provenance so equal relative paths do not collapse accidentally. A global library can contribute to a project export; that does not make it a lower build layer. Host-global content must enter a build through declared discovery/selection, never ambient reads silently merged during planning. Flattening global content into project instructions can duplicate it if the host also loads the global file.

Codex uses `$CODEX_HOME` instructions, root-to-CWD concatenation, and `.agents/skills` at ancestors; same-name skills can coexist. OpenCode uses `~/.config/opencode/AGENTS.md`, first-match local instruction categories, and skill discovery bounded by the Git worktree. Claude accumulates ancestor instructions and loads descendants on demand rather than enforcing deterministic prose conflict resolution. Cursor supports nested skills despite nbiish v2's global/root-only skill proposal. Consequently, global/project/nested is not one portable precedence ladder, and a main checkout supplies no automatic cross-worktree cascade. ([Codex instructions](https://developers.openai.com/codex/guides/agents-md), [Codex skills](https://developers.openai.com/codex/skills), [OpenCode rules](https://opencode.ai/docs/rules/), [OpenCode skills](https://opencode.ai/docs/skills/), [Claude memory](https://code.claude.com/docs/en/memory), [Cursor skills](https://cursor.com/docs/context/skills))

Exports should therefore declare whether they retain separate contexts, flatten selected contributions, or generate a particular tool's recognized entrypoints. Unknown `.agents` members remain preservable without claimed auto-loading. Existing project-relative output validation should not be bypassed to reach home: global targets need their own roots, ownership, and state lifecycle. Projection loss includes curation, logical identity, and provenance not represented in target formats; re-importing generated AGENTS.md cannot reconstruct the source recipe by bytes alone.

## Conclusion

The useful universal boundary is an **explainable artifact plan**, not a universally interpreted directory. That boundary lets Mogent distribute unfamiliar content intact while composing selected prose precisely, without becoming a configuration-schema engine or imitating every host's runtime loader.

The next decision should settle typed adapter relationships and explicit payload operations before names or override defaults. Mixed export and scope-projection tabletop cases can then expose the actual tradeoffs; the accompanying needs-DF experiment keeps those choices open while turning the evidence into reviewable design questions.
