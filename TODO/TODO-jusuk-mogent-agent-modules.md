# TODO-jusuk - Mogent: Modular Agent Prompt Manager

A CLI tool for assembling AGENTS.md files from modular, scoped templates.

## Decision Intent Log

ID: DI-jusuk
Date: 2026-07-24 12:00:00
Status: active
Author: 95124070+Qu1ncyRy4n@users.noreply.github.com (Quincy Ryan)
Decision: Build mogent as a Go CLI with tag-based scoping, slash hierarchy, header attribute tags, grouped TOML config, and section-aware semantic diff. Design for promise grid integration but implement later.
Intent: Create a lightweight but flexible tool for managing agent prompts across repos, orgs, and people. Enable cross-org sharing via promise grid in future.
Constraints: Must work with existing AGENTS.md format. Must be simple to adopt incrementally.
Affects: tools/mogent/, AGENTS.toml, .mogent/

ID: DI-jusuk-edit
Date: 2026-07-24 13:00:00
Status: active
Author: 95124070+Qu1ncyRy4n@users.noreply.github.com (Quincy Ryan)
Decision: Edit command opens module in $EDITOR or assembled AGENTS.md. No special diff/merge on save.
Intent: Keep edit simple - use existing editor. Future versions can add merge assistance.
Constraints: Requires $EDITOR to be set.
Affects: tools/mogent/cmd/edit.go

ID: DI-lorad
Date: 2026-07-27 19:21:34
Status: active
Author: 95124070+Qu1ncyRy4n@users.noreply.github.com (Quincy Ryan)
Decision: Redesign mogent's next module model around selectable Markdown heading subtrees at any heading level, referenced by Obsidian-style `[[path#id]]` preset entries. Each selectable block uses an explicit human-authored `agent_module.id` in an HTML comment with YAML-like `agent_module:` metadata placed immediately after the heading. TOML or constructor defaults define render order unless a later explicit override is selected. Diff output compares selected block IDs plus rendered text.
Intent: Keep the authoring surface natural for Markdown and Obsidian-style browsing while avoiding fragile heading-derived or content-hash identifiers. Make files organizational containers and heading subtrees the standard selectable unit so the future selector can flatten choices without forcing every option into its own file. Preserve stable readable preset diffs and leave tags deferred until block IDs, validation, presets, ordering, and diff behavior are stable.
Constraints: Applies to the next mogent redesign before code changes. Metadata comments must be stripped from rendered `AGENTS.md`. A selected heading includes its descendant content until the next heading of the same or higher level. Whole-file selection is not a separate first-class standard; a whole-file module should be represented as a file with one selectable top-level heading. New implementation must validate missing files, missing IDs, duplicate IDs, and ambiguous references loudly instead of silently dropping content. Manual handle allocation remains a user-approved exception while `tools/mint-handle` is unavailable.
Affects: tools/mogent/internal/module/, tools/mogent/internal/assemble/, tools/mogent/internal/diff/, tools/mogent/internal/toml/, tools/mogent/cmd/, AGENTS.toml, .mogent/

ID: DI-soviv
Date: 2026-07-27 19:58:37
Status: active
Author: 95124070+Qu1ncyRy4n@users.noreply.github.com (Quincy Ryan)
Decision: Make the next dogfood iteration favor readable block-first CLI output and initialization: `mogent list` should render a tree with active/inactive markers, `mogent init` should generate block-native AGENTS.toml/module content instead of stale file/tag prompts, and this repo should use the generated AGENTS.md path as a dogfood target.
Intent: Improve day-to-day usability before adding a full TUI. Keep the POC fast and inspectable by making selected blocks obvious, making init produce immediately buildable block metadata, and using this repo's generated AGENTS.md to test the init/build path.
Constraints: This does not settle global-vs-local module storage, manual AGENTS.md drift handling, import/merge behavior, or rich diff semantics; those need a later TE because multiple plausible models remain. The output should stay plain-text friendly while allowing ANSI color where useful later.
Affects: tools/mogent/cmd/list.go, tools/mogent/cmd/init.go, AGENTS.toml, AGENTS.md, .mogent/

ID: DI-nasot
Date: 2026-07-27 20:17:27
Status: active
Author: 95124070+Qu1ncyRy4n@users.noreply.github.com (Quincy Ryan)
Decision: Extend the block-first CLI with width-aware tree output, explicit/inherited/inactive markers, init-time module discovery, simple module view/edit actions, and build validation for empty selected blocks or empty rendered output. Use `+` for explicitly selected blocks, `|` for blocks included through a selected parent subtree, and `-` for inactive blocks.
Intent: Make mogent easier to dogfood in real terminals before a full TUI exists. The list view should stay readable on narrow terminals, init should show what modules already exist instead of pretending only defaults exist, and build should surface weak outputs instead of silently writing empty or accidental files.
Constraints: Width awareness may use plain terminal width detection and truncation instead of a full layout engine. View/edit in init remains a simple prompt-driven action, not a yazi-style TUI. Function naming for this slice is locked as `terminalWidth`, `truncateToWidth`, `blockSelectionState`, `selectionMarker`, `discoverModuleSources`, `promptModuleAction`, `viewModule`, `editModule`, and `ValidateOutput`. Manual AGENTS.md drift handling remains deferred to a future TE.
Affects: tools/mogent/cmd/list.go, tools/mogent/cmd/init.go, tools/mogent/cmd/build.go, tools/mogent/internal/assemble/engine.go, tests under tools/mogent/

ID: DI-bakom
Date: 2026-07-27 20:30:47
Status: active
Author: 95124070+Qu1ncyRy4n@users.noreply.github.com (Quincy Ryan)
Decision: Add a Bubble Tea-backed `mogent tui` command as the next selector proof of concept before expanding the module library. The first TUI scope is config/module loading, yazi-like block browsing, cursor navigation, visual selected/inherited/inactive states, and in-memory block toggling. Config writes, presets, global/local module storage, tags, drift detection, and import/merge are deferred.
Intent: The user cannot productively dogfood larger module sets through plain TOML and static tree output. A small TUI shell addresses the immediate navigation blocker while keeping dangerous write semantics out of the first iteration.
Constraints: The implementation may add Bubble Tea and Lip Gloss dependencies. Keep the first command isolated from existing build/list behavior. Use `tui` as the command name, `newTUIModel` for model construction, `tuiItem` for selectable rows, and `toggleCurrentItem` for in-memory state changes. Manual handle allocation remains a user-approved exception while `tools/mint-handle` is unavailable.
Affects: tools/mogent/cmd/root.go, tools/mogent/cmd/tui.go, tools/mogent/go.mod, tools/mogent/go.sum, tests under tools/mogent/

ID: DI-modun
Date: 2026-07-29 20:20:00
Status: active
Author: 95124070+Qu1ncyRy4n@users.noreply.github.com (Quincy Ryan)
Decision: Reorient mogent onto one tree model, recorded in `docs/DESIGN.md` as the single design of record. Markdown headings are the tree; the selectable unit is a heading plus its subtree (drop the "block" term). Nodes are addressed by heading path with an optional `<!-- id: … -->` escape hatch, never content hashes, never an id forced on every node, never Obsidian `[[path#id]]` syntax. Selection pulls subtrees; a descendant can be deselected as an exception. `include` (local path/URL, last-wins merge) and Go text/template vars are core sharing mechanisms, not future work. Editing is copy-on-write localization into the local `.mogent` library with config repointed. Config is TOML (`library` + `include` + `order` + `select`/`deselect` + `[vars]`). Categories classify by content role only; weight/domain/language are handled by include/select and later tags. Tags are deferred to a future search convenience.
Intent: Collapse two conflicting doc generations (tag model vs block model) into one coherent, pragmatic model so the library, selector, and diff share one selection unit and the design stops feeling fragmented.
Constraints: Supersedes the block/id-per-node and wikilink framing of DI-lorad and TE-tavim, and the tag-based DESIGN-SUMMARY/TE-mogent-module-architecture. Docs-only so far; code still reflects the old model and must be reconciled (see Cleanup follow-ups in DESIGN.md §8). Manual handle allocation remains a user-approved exception while `tools/mint-handle` is unavailable.
Affects: docs/DESIGN.md, tools/mogent/ (reconciliation pending), AGENTS.toml, .mogent/, docs/other_repo_agents/

## Subtasks

- [x] jusuk.1 Project scaffolding - Go module, CLI skeleton, basic build
- [x] jusuk.2 TOML parser - Parse AGENTS.toml with categories, tags, sources
- [x] jusuk.3 Module parser - Parse markdown with header attribute tags {#tag1 #tag2}
- [x] jusuk.4 Tag resolver - Hierarchical tag matching with parent-implies-children
- [x] jusuk.5 Assembly engine - Combine modules based on active scopes
- [x] jusuk.6 Build command - mogent build outputs AGENTS.md
- [x] jusuk.7 Edit command - mogent edit <module> and mogent edit (assembled)
- [x] jusuk.8 Diff command - Section-aware diff between scopes
- [x] jusuk.9 List command - List available and active modules
- [x] jusuk.10 Dogfood - Use mogent on this repo's AGENTS.md

## Feature Backlog

- [ ] Add TUI save support: persist selected blocks back to `AGENTS.toml`, show dirty/saved state, and avoid silent config rewrites.
- [ ] Add TUI preview support: show rendered `AGENTS.md` output for the current in-memory selection before saving.
- [ ] Add TUI diagnostics panel: surface missing files, duplicate IDs, empty selected blocks, unresolved references, and no-extension source hints in one place.
- [ ] Add generated-output drift detection: warn when `AGENTS.md` was manually changed since the last build and offer an explicit handling path.
- [ ] Add a richer diff model: compare selected block sets, rendered text, and similar blocks inside the same module instead of only section-level status.
- [ ] Add presets/templates: save and load reusable selections for learning-heavy, fast-iteration, design-heavy, Nix/devshell, Rust/Godot, Rust/Python ML, and Obsidian/session-log workflows.
- [ ] Add module library expansion: create reusable modules for tutor mode, TTS-friendly communication, architecture laws, strict testing, commit cadence, docs/session logs, and developer involvement levels.
- [ ] Decide local-vs-global module storage: define whether project-specific changes stay in-repo, sync to a shared module library, or become explicit overrides.
- [ ] Add import/merge workflow: help convert manually edited `AGENTS.md` changes into local modules, shared modules, or rejected drift.
- [ ] Add tag/conflict support after the block model stabilizes: searchable tags, optional XOR groups, and conflict warnings for incompatible thinking/communication styles.

## Module Extraction Plan (from corpus)

Empirical basis for building the library, derived from the 18 samples in
`docs/other_repo_agents/`. The corpus splits into tiers that all map onto the four
content-role categories:

- Tier 1 — universal core (in ~every sample): Project Structure (Identity); Build/Test,
  TODO, Commit/PR, Workflow (Instructions); Coding Style (Format); Testing
  (Instructions).
- Tier 2 — heavy process (ciwg/promisegrid/decomk family, near-verbatim shared text):
  Decision-First, Thought Experiment Protocol (11-node subtree), DR/DI, Change Review,
  Comment Preservation (Instructions); Diff Discipline, Error Handling, Glossary
  (Format); Runtime Artifact Hygiene (Constraints). Heavy repos `select` these; light
  repos omit them.
- Tier 3 — domain/genre: Promise Action Minimalism (grid); openai/codex Rust+TUI+lang
  best-practices; pg_learning teaching agent (Teaching Method/Assessment → Cognition;
  Communication Contract → Communication; Session Protocol/Obsidian → Notes/docs).
  These are separate includable libraries / later tags, not categories.

Next build steps: extract Tier 1 + Tier 2 as real modules from the corpus text;
reconcile `AGENTS.toml` and `.mogent/modules/` to the `library`/`select`/`include`
shape; remove retired tag code.

## Design References

- docs/DESIGN.md - **Design of record** (DI-modun). Read first; supersedes the docs below where they disagree.
- docs/brainstorn.md - Conversation-derived brainstorm notes for category taxonomy, TUI selection/save behavior, and fork import workflow (YAML config section folded into DESIGN.md §3.3-3.4).
- docs/examples/DESIGN-SUMMARY.md - **Superseded** Gen-1 tag-model summary (kept for history).
- docs/examples/AGENTS-style-e-hierarchical.md - TOML style reference
- docs/examples/AGENTS-style-f-slashes.md - Slash notation reference
- docs/thought-experiments/TE-tavim-mogent-module-reference-model.md - Reference-model TE; block/id-per-node/wikilink framing superseded by DESIGN.md §3.1. Manual handle allocation approved by user because `tools/mint-handle` was unavailable at the approved paths.
- docs/thought-experiments/TE-bakom-mogent-tui-first-selector.md - TUI-first selector TE. Recommends a Bubble Tea `mogent tui` browser before expanding the module library.
- docs/thought-experiments/TE-tavim-mogent-module-reference-model.md - Follow-up TE for block references, presets, metadata, render order, and diff model. Manual handle allocation approved by user because `tools/mint-handle` was unavailable at the approved paths.
- docs/thought-experiments/TE-bakom-mogent-tui-first-selector.md - TUI-first selector TE. Recommends a Bubble Tea `mogent tui` browser before expanding the module library.
