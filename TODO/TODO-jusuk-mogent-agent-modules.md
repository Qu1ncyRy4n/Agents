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

ID: DI-ralik
Date: 2026-07-29 21:30:00
Status: active
Author: 95124070+Qu1ncyRy4n@users.noreply.github.com (Quincy Ryan)
Decision: Reframe mogent's config as a document manifest, not a delta. The library tree and the document tree are distinct: the manifest (`agents.yaml`) is the output document's nested, ordered outline, where each entry references a source node, pulls a whole subtree (with optional deep excludes), or nests further structure. YAML is the canonical config format with a strict, loudly-validated schema. Named `sources` (local path or URL) replace ambient includes; every reference carries explicit provenance (`shared:`, `grid:`, `local:`). Document order is manifest order. Swap = change one reference; edit = copy-on-write flip to `local:`. One interaction model unifies the CLI — manipulate the manifest, confirm, build — with gh-CLI-style init (choose sources → navigate/toggle/view/edit/swap → confirm read) as the primary flow. Old implementation (`tools/mogent`, `.mogent`) removed for a clean rebuild: milestone 1 parse/resolve/render, milestone 2 navigator, milestone 3 save flow.
Intent: The delta model (flat select/deselect over the library) made the config illegible as a description of the output and conflated library organization with document structure. The manifest makes the config readable as the document's table of contents, makes swap/exclude/reorder natural operations, and resolves the YAML-vs-TOML question by model fit rather than taste.
Constraints: Supersedes the flat select/deselect TOML config portion of DI-modun; the tree model, heading-path identity, id escape hatch, copy-on-write editing, category roles, and fail-loud validation from DI-modun remain in force. Manifest composes from libraries rather than mirroring them — explicit nesting only where a structural choice is made. Exact YAML schema finalizes during rebuild milestone 1. Manual handle allocation remains a user-approved exception while `tools/mint-handle` is unavailable.
Affects: docs/DESIGN.md, agents.yaml (future), AGENTS.toml (to be replaced), rebuilt codebase (location TBD)

ID: DI-sufok
Date: 2026-07-29 22:15:00
Status: active
Author: 95124070+Qu1ncyRy4n@users.noreply.github.com (Quincy Ryan)
Decision: The manifest owns the rendered document headings and outline. Source headings guide init/build suggestions but do not control the output. A manifest may name user-defined sources such as `cdint`, `cclab`, `personal`, or a collaborator name; source aliases have no fixed enum. One manifest node may compose matching source subtrees in declared order. Before confirmation, mogent warns about overlapping descendant paths and offers either one combined section or separate locally named sections. Tree state, collapsed state, and source provenance are always visible in text; color is an early secondary cue, never the sole signal.
Intent: Let each repository express its own document vocabulary and combine useful libraries without hidden source precedence, while retaining clear provenance and a readable monochrome interface.
Constraints: Duplicate YAML keys and duplicate source aliases are errors. Source aliases are manifest-local names and may be deliberately renamed with their references. Exact YAML syntax for multi-source composition is finalized in rebuild milestone 1. URL source pinning remains deferred; changed source content must be shown in confirmation before writing output.
Affects: docs/DESIGN.md, agents.yaml (future), navigator and confirmation UI (future)

ID: DI-venit
Date: 2026-07-29 22:15:00
Status: active
Author: 95124070+Qu1ncyRy4n@users.noreply.github.com (Quincy Ryan)
Decision: Rewrite this repository's root AGENTS.md as a concise adaptive baseline now. Routine work uses a focused change loop and proportionate validation. Architecture, protocol, irreversible, security, public-specification, cross-repository, or otherwise suspicious work escalates to the full decision-first process. Normal handoffs are concise; the complete decision matrix and runtime-path matrix apply only to fully governed work.
Intent: Keep normal work easy to execute while preserving deliberate analysis and evidence where an error would have durable or broad consequences.
Constraints: Documentation-only work receives structural and diff checks. Code changes validate the affected area; Go behavior changes also run gofmt, focused tests, and errcheck from the relevant module. When results, scope, or risk feel uncertain, recommend the wider relevant check before handoff.
Affects: AGENTS.md, docs/codex_eco/README.md, future shared baseline modules

ID: DI-munet
Date: 2026-07-29 23:00:00
Status: superseded by DI-renim
Author: 95124070+Qu1ncyRy4n@users.noreply.github.com (Quincy Ryan)
Decision: Begin rebuild milestone one at `tools/mogent/` with a strict, uniform YAML manifest schema. Entries use `heading` plus exactly one of ordered `from` references or `children`; source paths are local Markdown-library directories only. URL fetching, compact scalar syntax, and mixed source/child entries are deferred.
Intent: Deliver a reliable parse-resolve-render core now while keeping later source fetching and interaction design explicit.
Constraints: This implements the recommended baseline in TE-vorum at the user's direction to defer questions while work proceeds. It does not settle future URL source, compact syntax, or mixed-entry support.
Affects: tools/mogent/, agents.yaml examples and tests

ID: DI-renim
Date: 2026-07-30
Status: superseded by DI-vukam
Author: 95124070+Qu1ncyRy4n@users.noreply.github.com (Quincy Ryan)
Decision: Mogent is a root Go module. Its manifest accepts both the compact
heading-key YAML syntax shown in DESIGN.md and explicit `heading` object syntax,
normalizing them to one internal document tree. Sources may be local paths or
un-pinned HTTP(S) Git URLs now; URL fetches are intentionally fragile until the
later pinning/lockfile design. `from` and `children` remain mutually exclusive.
Intent: Keep the manifest concise for ordinary use while allowing richer
composition, and enable early cross-repository sharing without pretending the
remote-source safety design is complete.
Constraints: URL sources must fail loudly, never use an ambient source search
path, and warn that content may change. Pinning, caching, and changed-content
confirmation remain follow-up work.
Affects: go.mod, cmd/mogent/, internal/, agents.yaml parsing and source resolution

ID: DI-vukam
Date: 2026-07-30
Status: active
Author: 95124070+Qu1ncyRy4n@users.noreply.github.com (Quincy Ryan)
Decision: Milestone one is a root Go module with compact and explicit manifest entries, local source paths only, and mutually exclusive `from` and `children`. Require source-qualified exclusions, strict YAML, non-empty selected nodes, and all template variables. Default output to `AGENTS.md`. Build atomically and retain only a local generated-output hash; refuse to overwrite direct edits unless the user passes `--force`. URL sources, collision choices, and interactive direct-edit import are the next features.
Intent: Deliver a reliable local renderer without silently fetching mutable remote content, overwriting handwritten instructions, or guessing where a direct edit belongs.
Constraints: Local paths may be relative, absolute, or home-relative but references cannot escape a source root. Duplicate keys, YAML anchors or aliases, unknown keys, wrong types, unresolved paths, overlap collisions, and empty output are errors. The exact manifest contract is `docs/IMPLEMENTATION-M1.md`.
Affects: docs/IMPLEMENTATION-M1.md, go.mod, cmd/mogent/, internal/, agents.yaml parsing, source resolution, output state, tests
Supersedes: DI-renim

ID: DI-voraz
Date: 2026-07-29 22:45:00
Status: active
Author: 95124070+Qu1ncyRy4n@users.noreply.github.com (Quincy Ryan)
Decision: Seed plain Markdown source libraries at `libraries/` before the renderer is rebuilt. Add reusable, sanitized modules for CDINT baseline/process/Go/PromiseGrid, UCD research/Python with `uv` under dependency management, and independently selectable Nix safety. Do not materialize a personal library until its privacy and source boundary are resolved.
Intent: Give the new implementation a useful local corpus and keep sensitive learner, vault, and machine details out of reusable source content.
Constraints: Headings are the library tree. Modules must contain no source-specific credentials, absolute paths, identifying research data, or private learner state. `cdint`, `ucd_research`, `personal`, and `nix` boundaries remain provisional under DR-garom.
Affects: libraries/, docs/codex_eco/README.md, future agents.yaml dogfood manifest

ID: DI-tuvim
Date: 2026-07-30
Status: active
Author: 95124070+Qu1ncyRy4n@users.noreply.github.com (Quincy Ryan)
Decision: Build the manifest-based `mogent tui` navigator with Bubble Tea, Lip Gloss, and only the needed Bubbles components. On terminals at least 140 columns wide, its default screen has three panes: manifest tree, complete rendered AGENTS.md, and selected-source context. The selected tree node aligns and highlights its output and source text. On narrower terminals, use a tree pane plus one detail pane; the detail pane explicitly switches between Final and Source. Tab moves keyboard focus between visible panes, while 1, 2, and 3 directly select Tree, Final, and Source views. Escape returns to the tree. A later full-view mode may use the same view tabs and navigation. Draft manifest edits are marked `~`; saving manifest selection/order removes it. `L` is reserved for an explicit copy-on-write local override. Save and build write agents.yaml and AGENTS.md after confirmation, but do not change source-library modules.
Intent: Make provenance, document structure, and rendered output visible together for a convincing and usable first navigator without relying on color or hidden modes.
Constraints: Text labels and symbols remain meaningful without color. Full Final view shows the entire rendered document and aligns to the selected node. Default Source context shows the selected section, its parent heading, and adjacent siblings, with a later full-file toggle. When a selected node maps to multiple output locations, expose a count and provide cycling rather than claiming one location. An invalid draft preserves the last valid preview and shows the current error. The user must explicitly confirm before writes. Source-library changes require a separate local-override or import action.
Affects: docs/IMPLEMENTATION-M2.md, docs/DESIGN.md, cmd/mogent/, internal/, TUI tests

ID: DI-vasel
Date: 2026-08-04
Status: active
Author: 95124070+Qu1ncyRy4n@users.noreply.github.com (Quincy Ryan)
Decision: Keep mogent core-first. The TUI remains an important demonstration and editing surface, but durable behavior belongs in reusable workspace operations that can also back CLI commands and future UI clients. Every TUI action that changes a manifest, local override, source selection, or drift decision should have a matching core operation and a noninteractive command shape when useful.
Intent: Make mogent usable by humans, agents, and scripts without trapping product behavior inside a terminal UI. This also makes future GUI work less risky because it can reuse the same load/draft/preview/save/source/drift operations.
Constraints: Do not roll back the M2 TUI. Extract behavior incrementally: preserve the current navigator behavior while moving save, preview, dirty state, source inspection, source coverage, copy-on-write, and drift/import decisions into reusable packages.
Affects: docs/DESIGN.md, internal/workspace/ or equivalent core package, internal/navigator/, internal/cli/, future M3+

ID: DI-pesun
Date: 2026-08-11
Status: active
Author: user
Decision: Move Nix guidance under the `personal` source as `personal:lang/nix` rather than retaining a separate `nix` source. Keep Nix independently selectable as a manifest node, and retain its activation and machine-safety warnings in the module content.
Intent: Match the source boundary to current ownership and dogfood use while preserving deliberate selection for guidance that can affect an active machine.
Constraints: This resolves only the Nix portion of DR-garom. The broader `cdint`, `ucd_research`, `personal`, and future independent-library boundaries remain provisional and may be revisited after more real manifests are exercised.
Affects: DR/DR-garom-library-source-boundaries.md, libraries/personal/lang/nix.md, libraries/nix/, docs/codex_eco/README.md, testing-ground/personal-go-nix/

ID: DI-ravam
Date: 2026-08-11
Status: active
Author: user-directed goal; recorded by Codex
Decision: Implement localization and drift core-first. Store local Markdown under `.mogent/library`, expose it through the explicit `local` source alias, and store versioned origin records in `.mogent/provenance.yaml`. Address localization by manifest heading path, require an explicit source reference for composed entries, and provide non-mutating dry runs. Drift import is conservative and may change one unambiguously mapped manifest section; ambiguous edits remain unresolved.
Intent: Let CLI, TUI, agents, and future interfaces share one safe copy-on-write and drift workflow without modifying shared sources, polluting Markdown with workspace history, or guessing how direct edits map back to the manifest.
Constraints: Validate the complete resulting render before committing writes. Do not overwrite an existing different local artifact. Shared sources are never transaction targets. Rejecting direct edits requires explicit force. Exact behavior is defined by `docs/IMPLEMENTATION-M3.md` and TE-ravam.
Affects: docs/IMPLEMENTATION-M3.md, internal/workspace/, internal/cli/, .mogent/library/, .mogent/provenance.yaml

ID: DI-fipam
Date: 2026-08-11
Status: active
Author: user-directed goal; recorded by Codex
Decision: Keep remote Git URLs in `agents.yaml`, commit immutable URL/commit/Markdown-hash records in `mogent.lock.yaml`, and keep verified checkouts under ignored `.mogent/sources/<alias>/<commit>`. Ordinary workspace loads are offline-only and fail on missing or mismatched lock/cache state. Initial pinning is explicit; updates preview changed Markdown paths and require `--accept` before changing the lock or installed cache.
Intent: Make remote sources reproducible and reviewable without turning build or browsing into hidden network operations or silently consuming a moving branch.
Constraints: HTTP(S) Git only initially. Full commit IDs and deterministic Markdown content hashes are required. URL/lock mismatch, cache hash mismatch, unsafe path input, duplicate keys, and unreadable content are errors. Exact behavior is `docs/IMPLEMENTATION-M4.md` and TE-fipam.
Affects: agents.yaml URL sources, mogent.lock.yaml, .mogent/sources/, internal/sourcecache/, internal/render/, internal/cli/

ID: DI-vurap
Date: 2026-08-11
Status: active
Author: user-directed dogfood decision; recorded by Codex
Decision: Preserve compact scalar sources and add an explicit normalized source form with `location` plus optional `subdir`. Bind URL subdirectories into the lock identity; scope Markdown hashing, update review, resolution, and browsing to the selected root while retaining the full immutable checkout in the ignored cache.
Intent: Let one Git repository publish several clean libraries without changing short source references or indexing unrelated repository documentation, while leaving repository splits to ownership and release decisions.
Constraints: Subdirectories are relative slash paths with no `.`, `..`, empty components, backslashes, absolute form, or symlinked components. Ordinary loads remain offline-only. Exact behavior is recorded in TE-vurap and the updated M4 contract.
Affects: agents.yaml source values, mogent.lock.yaml, internal/manifest/, internal/sourcecache/, internal/render/, docs/IMPLEMENTATION-M4.md

ID: DI-folar
Date: 2026-08-12
Status: active
Author: user
Decision: Make source-relative directories first-class selectable subtree nodes
and present source inventory plus manifest coverage through one aligned tree
model. Keep `mogent coverage` as the compatibility preset for
`mogent source list --coverage --tree`. Every tree uses explicit textual node
kind and state indicators. Use portable ASCII by default and allow presentation
choice through `display.chars: ascii|unicode`, with explicit CLI flags taking
precedence over user-local XDG YAML.
Intent: Preserve the actual library hierarchy during discovery and selection,
avoid separately evolving overlapping list/coverage commands, and improve
readability without making color or Unicode carry required meaning.
Constraints: Selecting a directory includes every descendant in deterministic
path order; `exclude` narrows that selection. Directory/heading identity
collisions are represented explicitly rather than silently choosing one.
Presentation settings cannot affect manifest resolution or rendered output.
Source overlap remains an error. ASCII references and state words remain stable
regardless of character-set preference.
Affects: docs/DESIGN.md, internal/library/, internal/workspace/, internal/cli/,
internal/presentation/, README.md, docs/DOGFOOD-SESSION-3.md

ID: DI-norim
Date: 2026-08-12
Status: active
Author: user
Decision: Refine source-tree and authoring previews with terminal-width-aware
wrapping, expanded inherited source subtrees, related-selection review, and
stable named placement. `--under` accepts `--first|--last`; `--before` and
`--after` use manifest heading paths and infer their parent. Presentation uses
`display.fit: term|none` plus an optional explicit width; redirected output is
unbounded unless that width is configured.
Intent: Keep large inventories readable in real terminals and make directory
additions reveal their actual effect and likely redundancy before any write.
Constraints: Fitting changes presentation only. Numeric manifest indexes are
not a public interface. Exact source ancestry is labeled overlap; cross-source
path similarity is labeled for review rather than asserted as a conflict.
Affects: internal/presentation/, internal/cli/, internal/workspace/add.go,
docs/DESIGN.md, docs/AUTHORING-PLAN.md, docs/DOGFOOD-SESSION-3.md

ID: DI-zunap
Date: 2026-08-18
Status: active
Author: 95124070+Qu1ncyRy4n@users.noreply.github.com (Quincy Ryan)
Decision: Publish a supported Go workspace API and make the CLI a thin client
of that same API. Move reusable implementation out of Go's `internal/`
boundary, but export lower-level packages selectively rather than treating
every existing symbol as a compatibility promise. Keep CLI and TUI presentation
as adapters unless a concrete external use justifies their own public API.
Reject symlinked local source roots, directories, and Markdown files by default,
matching pinned-source safety. Make repository template values reproducible by
materializing discovered values in the manifest rather than silently deriving
build output from the checkout directory or Git remote.
Intent: Let agents, scripts, the CLI, and future clients use one safe,
dogfooded implementation without exposing accidental internals or allowing
machine-local filesystem and Git state to change trusted source content or
rendered output.
Constraints: Public operations preserve strict validation, preview/dry-run,
offline ordinary commands, atomic apply, rollback, and overwrite protection.
Symlink rejection must produce actionable diagnostics. Public types and package
paths require compatibility review before release. The detailed decision and
cleanup sequence is `docs/PUBLIC-API-AND-CLEANUP-PLAN.md`.
Affects: go.mod, cmd/mogent/, internal/, future public Go packages,
docs/PUBLIC-API-AND-CLEANUP-PLAN.md, README.md

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
- [x] jusuk.11 Lock the milestone-one manifest schema, source scope, and rebuilt Go module location (DI-vukam; TE-vorum).
- [x] jusuk.12 Resolve the first source-boundary review: move Nix under `personal` while leaving broader boundaries provisional (DI-pesun; DR-garom).
- [x] Reconcile the post-merge Nix conflict by retaining the DI-pesun canonical
  module at `personal:lang/nix` and removing the redundant top-level `nix`
  source after user confirmation.
- [x] jusuk.13 Seed sanitized CDINT, UCD research, and Nix source libraries (DI-voraz).

## Feature Backlog

### User Review Queue

- [ ] **Recurring owner reminder:** Quincy must read every checked-in module
  under `libraries/` and record content corrections, metadata adjustments,
  overlap, and source-boundary notes. Mention this in substantive project
  handoffs until complete, and offer the next bounded library section rather
  than letting the review disappear into the general backlog.

Rebuild milestones now proceed core-first: 1. manifest parse/resolve/render, 2.
navigator proof of concept, 3. extract reusable workspace operations, 4. add
copy-on-write localization, source browsing, drift/import, and URL pinning on top
of that shared core.

For agent- and script-facing capabilities, preserve the smaller progression that
motivated the core extraction: `status` -> show/inspect -> sources/discovery ->
coverage -> localization. The first four now have CLI surfaces and should be
stabilized against the shared workspace core before copy-on-write localization
is added. Repository docs and current Git state are authoritative over older chat
summaries; see `docs/HANDOFF.md` for the compact resume point.

## Phase Sequence

### Phase 0 - Stabilize Current Dogfood

- Keep `build`, `status`, `coverage`, `source list`, `source show`, and `add`
  reliable for local Markdown sources.
- Fix README/install accuracy and keep examples runnable from a fresh checkout.
- Keep Vroca and the testing-ground manifests as smoke-test fixtures for real
  authoring friction.

### Phase 1 - CLI Discoverability And Pleasant Dogfood

- Add shell completion for commands, flags, source aliases, source references,
  manifest heading paths, and `--under` targets.
- Add typo suggestions and close-match diagnostics for source refs, command
  names, flags, source aliases, and manifest heading paths.
- Improve empty states, especially when tag search finds nothing but source
  paths or headings would match the user's text.
- Add compact/TLDR browsing modes and source/path/heading search so users do not
  need two terminals just to explore a library.
- Prioritize these "pleasant first" improvements before deeper metadata work so
  Vroca-style dogfooding remains fast while the library model evolves.
- First pass complete: source refs, source aliases, source subcommands, command
  names, flag names, and manifest heading paths now have close-match
  suggestions.

### Phase 2 - Coverage And Status UX

- Make `coverage --tree` show included, inherited, excluded, partial, and
  unused state in one tree with text markers and optional color.
- Add placement-aware grouping once library metadata exists.
- Improve `status` with clear missing/stale/clean output and actionable hints
  such as `run: mogent build`.
- Add better diagnostics for organizational-directory refs such as
  `cdint:engineering`, including descendant suggestions.
- Resolved first pass: `mogent complete <kind>` exposes candidate lists and
  `mogent completion bash|zsh` emits shell wrappers without committing to Cobra
  or another command framework.

### Phase 3 - Library Metadata And Recommended Order

- Add `library.yaml` or equivalent source-level metadata for suggested
  placement, recommended order, presets, update/review dates, requirements,
  urgency/risk, and related nodes.
- Use that metadata to drive recommended order in init, coverage, source
  browsing, and template generation.
- Add "see also" / related-node support so adjacent styles and modules can point
  at one another.
- Add OR/XOR choice-group metadata for headings whose children are alternatives
  rather than additive modules. XOR groups should require one deliberate choice
  and warn when multiple contradictory children are selected; OR groups should
  make optional compatible choices explicit.
- Reorganize library source boundaries and top-level groups after dogfooding:
  core/code-principles, process/decision, stack/langs, communication/style,
  domains, and CDINT/PromiseGrid as non-universal domain material.

Initial dogfood extraction: `personal:engineering/staged-migration` and
`personal:engineering/local-service-design` are intentionally separate heavy,
opt-in modules. They generalize the Vroca Rust-design handoff without making
ordinary Rust CLI projects inherit daemon, protocol, or parity-gate policy.
They overlap deliberately with `cdint:process/decision-first`,
`cdint:engineering/test-strategy`, and `personal:lang/rust`; Phase 3 metadata
should eventually represent those relationships as selectable `see also` and
preset guidance rather than rendering library-maintenance prose into prompts.

### Phase 4 - Init, Templates, And Repo Prose

- Add a guided `mogent init` / builder walkthrough that starts from repo-shape
  templates and produces an editable starter `agents.yaml`.
- Add templates for transitional Python-to-Rust apps, heavy CDINT design-first
  repos, personal projects, research repos, static sites, and Nix-managed
  projects.
- Add repository-specific prose support through inline YAML block scalars, local
  source snippets, or dedicated repo overlay files.
- Add over-composition diagnostics that warn when a small repo pulls in broad
  policy subtrees.

### Phase 5 - Editing, Localization, And Drift

- Add safe source/local editing workflows, including copy-on-write local
  overrides and explicit shared-source editing for trusted libraries.
- Add drift detection and import/merge flows for hand-edited `AGENTS.md`.
- Add source-vs-local diffs, localization provenance, and promote-local-to-shared
  support.

### Phase 6 - Multi-Output And Broader Project Guidance

- Design multiple agent outputs such as `CLAUDE.md`, `GEMINI.md`, and
  `.codex/AGENTS.md`.
- Design multi-artifact project guidance generation beyond `AGENTS.md`: usage
  notes, dev docs, repo-specific reference indexes, command cheat sheets, and
  external doc refs.
- Add URL source pinning, lockfiles, source cache design, and changed-content
  review after local workflows are stable.

- [x] Extract M2 draft/save behavior from `internal/navigator` into a reusable workspace/session package.
- [x] Add first CLI parity command: `mogent status`.
- [x] Add first source coverage command: `mogent coverage`.
- [x] Add compact unused source output: `mogent coverage --unused-only`.
- [x] Add coverage filters: `--source`, `--content-only`, and `--tree`.
- [x] Add remaining structural coverage filters: `--leaves-only` and `--depth`.
- [x] Add source inspection command: `mogent source show <ref>`.
- [x] Add first manifest mutation command: `mogent add <ref>` with `--under`, `--append`, `--heading`, `--dry-run`, and `--rebuild`.
- [x] Add add-preview modes: `--preview=summary`, `--preview=patch`, `--preview=tree`, and `--preview=full`.
- [x] Add first file-level source metadata: YAML frontmatter with `tags`, `tldr`, `priority`, `scope`, `requires`, and soft `conflicts_with`.
- [x] Add metadata-aware source browsing and tag filtering: `source show --metadata` and `coverage --tag`.
- [x] Add source browsing command: `source list` with source, exact tag, tag-search, metadata, and sort options.
- [x] Add CLI command shapes for core-first localization and conservative drift handling (`localize`, `drift`; DI-ravam).
- [ ] Add optional force/logging for manifest mutations and rebuilds, likely under `.mogent`.
- [ ] Add richer source display: source-vs-local diffs and optional metadata-driven summaries.
- [ ] Decide raw Markdown/no-manifest command behavior; default commands remain manifest-based and use `agents.yaml` unless `--manifest` points elsewhere.
- [ ] Add diagnostics panel: surface missing sources, unresolved references, empty nodes, and duplicate ids in one place.
- [x] Add drift detection: regenerate from manifest, report output state, conservatively import one unambiguous section, or explicitly reject edits.
- [x] Add immutable URL source pinning with committed locks, verified ignored caches, and explicit changed-Markdown review (DI-fipam; TE-fipam).
- [x] Add explicit, pinned URL-source subdirectories so one repository can expose `libraries/cdint` or `libraries/personal` with short references (DI-vurap; TE-vurap). Live combined HTTPS pin/build verification remains pending while the approval service blocks network execution.
- [ ] After URL subdirectories work, decide repository topology by ownership and release cadence rather than tooling limitations: keep code plus small fixtures here; consider a public reusable-library repository; keep personal/private and institution-owned material in separately governed repositories. Preserve immutable URL/commit/subdirectory identity in locks.
- [x] Design and implement localization provenance with `local:` plus `.mogent/provenance.yaml` (DI-ravam; TE-ravam).
- [ ] Add promote-local-to-shared: push a localized override back up to its source library.
- [ ] Resolve TE-kavam: hierarchical tags, declared alternative families, meaningful conflict warnings, and atomic-file library shape.
- [x] Decide whether atomic library directory paths remain organizational only or contribute to source reference paths: directories contribute, filenames do not.
- [ ] Add split/import workflow: turn a complete hand-written `AGENTS.md` into a draft atomic module library for review.
- [ ] Design global source cache for init: remember known local/remote sources for visible reuse without ambient source resolution.
- [x] Plan multiple agent outputs without implementing them: manifest evolution, build transaction, state migration, drift behavior, CLI review, sequencing, and open decisions (`docs/MULTIPLE-OUTPUTS-PLAN.md`).
- [ ] Design multi-artifact project guidance generation beyond `AGENTS.md`: usage notes, dev docs, repo-specific reference indexes, command cheat sheets, and links to external docs. This needs substantial design because those files have different audiences, update cadence, visibility, and source-of-truth rules from agent prompts.
- [ ] Far-future note: explore pinned reference-doc support for project/API/design/dependency docs once prompt composition is stable.
- [x] Split initial checked-in libraries into atomic metadata-bearing modules.
- [ ] Add library expansion: extract Tier 1/Tier 2 corpus nodes; add tutor mode, TTS-friendly communication, architecture laws, strict testing, commit cadence, docs/session logs, and developer involvement levels.
- [ ] Design cross-file subtree behavior, or keep explicit manifest grouping as the standard for atomic modules.
- [ ] Add import/merge workflow: help convert manually edited `AGENTS.md` changes into local overrides, shared nodes, or rejected drift.
- [ ] Add tags as search/discovery after the core stabilizes: searchable tags, swap-alternative groups (XOR), and conflict warnings for incompatible styles.
- [x] Local-vs-global storage: resolved by copy-on-write localization (DESIGN.md §3.5) — shared libraries read-only, local overrides explicit in the manifest.
- [x] Add shell completion for source references, especially `mogent source show <ref>` and `mogent add <ref>`: implemented `mogent complete <kind>` candidate lists plus `mogent completion bash|zsh` wrappers.
- [x] Add a guided, agent-readable init flow with template listing, explicit source bindings, dry-run preview, and validated writes.
- [x] Add the first editable starter templates: minimal, Go, personal Go/Nix, and research Python. Expand to migration/static-site shapes after dogfooding.
- [ ] Add library-level metadata manifests, likely `library.yaml`, for suggested placement, recommended order, status, last updated/reviewed dates, requirements, urgency/risk, and preset membership.
- [ ] Add suggested-placement-aware browsing and coverage output, so unused modules can be grouped by where they likely belong in the rendered document rather than only by source path.
- [x] Improve `mogent coverage --tree` to show included, inherited, excluded, and unused state in one tree with text markers and optional color.
- [x] Improve `mogent status` with clearer stale/missing/clean status output, actionable hints such as `run: mogent build`, and optional color that is not the only signal.
- [x] Improve source-list empty results. When `--tag-search` finds nothing but matching source paths/headings exist, explain that the current filter searches metadata tags only and suggest `source list | rg <term>` or a future path/heading search.
- [x] Add source-list search over source reference paths, headings, TLDRs, and content snippets; keep tag search available as a precise metadata filter.
- [x] Add `source list <source-prefix>` shorthand for filtering one or more sources, for example `mogent source list cdint` instead of requiring `--source cdint`.
- [ ] Add filters for metadata fields beyond tags: `requires`, `conflicts_with`, `scope`, urgency/risk/status, and update/review age once those fields exist.
- [x] Add TLDR source browsing and coverage modes; coverage shows each file-level summary once rather than repeating inherited metadata on every heading.
- [x] Add concise file-level TLDR metadata across the seed libraries.
- [ ] Add heading-level TLDR metadata using a one-line
  `<!-- tldr: One-sentence summary. -->` comment immediately after the heading
  (and after any inline `<!-- id: ... -->`). Parse it as tool-only metadata,
  strip it from rendered output bytes, let it override file-level TLDR for that
  heading, and keep inherited file metadata as the fallback. Validate duplicate,
  misplaced, empty, and malformed TLDR comments loudly before encouraging
  library authors to depend on the syntax.
- [ ] Add paging or compact layouts for large command output, or first-class hints for piping through `less`, `rg`, and `fzf`.
- [x] Add terminal-width-aware inventory formatting. Keep the tree label and
  source reference together when practical, then wrap a long TLDR onto an
  indented continuation line instead of creating a very wide sparse row.
  Provide an explicit CLI/config preference (for example a width of `term`, a
  numeric width, or unbounded output), and keep redirected/machine output
  deterministic when no terminal width is available.
- [ ] Finish CLI presentation from real-project dogfooding: source-list fields are aligned; still indent template requirements beneath each template and consider restrained color or separators without making color the only semantic cue.
- [x] Keep normal list/tree output useful ASCII by default, including an inline key where states need explanation. Add `display.chars: ascii|unicode` and `--chars` without a vague `--pretty-print` mode.
- [x] Unify source inventory and coverage around one listing/tree presenter. `source list` supplies available nodes and metadata; `source list --coverage` overlays included, inherited, excluded, partial, and unused state. `mogent coverage` is the compatibility preset for the same tree.
- [ ] Design user-level presentation config with explicit precedence for color, field alignment, TLDR display, hints, and preferred preview mode. Respect `NO_COLOR`; keep durable project composition out of personal display preferences.
- [ ] Extend the XDG YAML presentation config beyond implemented `display.chars` and `display.align`; wire defaults into Nix/Home Manager and design color, TLDR, hints, preview, legend, and `auto` behavior.
- [ ] Revisit the CLI inspection surface: make `status` the concise aggregate workspace view and keep drift-specific mutation under an explicit resolution command or subcommand rather than maintaining two overlapping read-only reports.
- [ ] Make read-only `drift` a compatibility alias for the direct-edit view of
  `status`; move import/reject/inherit/propose behavior to directional mutation
  commands with migration guidance.
- [ ] Clarify `help`, `complete`, and `completion`: human help explains commands; the machine-readable candidate backend should be internal or clearly documented; add a safe shell-specific installation path instead of only printing a completion script.
- [ ] Package zsh/bash completion through Nix in the shells' normal completion directories. An explicit install command may help non-Nix users; generated completion text should remain available for package managers without encouraging users to paste it into shell startup files.
- [ ] Generate human help, shell candidates, and future agent-readable command descriptions from one command schema rather than maintaining separate human and LLM documentation. Prefer explicit output modes or aliases where the same information differs only in presentation.
- [x] Add `source list <alias> --tree` using the same presenter as the coverage overlay.
- [x] Make source directories first-class selectable subtree nodes. `personal:engineering` selects every descendant module, individual headings remain selectable, and exclusions narrow a selected directory.
- [ ] Add `mogent source add <alias> <path-or-url>` for extending an existing manifest, with dry-run preview, duplicate-alias validation, URL pinning guidance, and no implicit node selection.
- [ ] When the first CLI argument looks like `alias:path`, suggest `mogent add alias:path`; if the alias is undeclared, explain how to declare it rather than only reporting an unknown command.
- [ ] Make actionable hints configurable later (`--no-hints` and/or a persisted hint setting) while keeping hints enabled by default.
- [ ] Add a core-first manifest reorder command with explicit relative placement (`--before`, `--after`, or `--under`) and dry-run previews; resolve the exact heading-path and nesting semantics before implementation.
- [x] Extend `add` placement beyond parent selection: use `--first` or `--last`
  within `--under`, and `--before <manifest-heading-path>` or
  `--after <manifest-heading-path>` for stable relative placement. Prefer named
  anchors over fragile numeric indexes.
- [x] Make directory additions informative in `add --preview=tree`: expand the
  selected source subtree beneath the proposed manifest node, distinguish
  inherited source headings from authored manifest headings, and identify
  already-selected descendants or likely semantic overlap before writing.
- [ ] Extend localization provenance for pinned URL sources with the immutable source revision/lock identity. Keep one canonical sidecar record unless a later export feature deliberately embeds a portable copy in Markdown; do not add redundant `localized: true` metadata.
- [ ] Design promote-local-to-shared separately from localization: preview the local-vs-upstream diff, require an explicit writable source, and preserve Git review rather than silently editing a shared library.
- [ ] Later integrate `propose` with Git or another version-control adapter so a reviewed local-to-shared change can become a branch/commit/change request; keep the first reconciliation model independent of any one forge.
- [ ] Replace or broaden the `drift` command vocabulary after a workflow thought experiment. The command family must express direction: preserve a hand edit locally, compare with origin, inherit upstream changes into a customized node, or propose/promote a local change to a shared source. Consider Promise Theory/Grid vocabulary, but prefer terms that remain understandable without that background.
- [ ] Add origin freshness to status: pinned revision, last explicit upstream check, available reviewed update when known, and whether a localized node has diverged from its recorded origin. Never perform a network check as a side effect of status.
- [ ] Design explicit local/upstream reconciliation using the provenance base: three-way compare original content, current upstream, and local content; support an inherit/rebase-like flow and a reviewable change-request/promote flow without silently overwriting either side.
- [ ] Add a guided module-creation command independent of init: prompt or accept flags for heading, Markdown content, TLDR, destination/source, and optional metadata; preview the resulting file/reference before writing.
- [ ] Run a privacy/security thought experiment for cross-agent activity and memory interoperability. Prefer explicit workspace event/handoff records; do not scrape private Claude, Codex, or other product histories by default.
- [ ] Revisit the product name before a broader public release. Explore respectful psychology/neuroscience references involving plurality, integration, memory, or perspective; avoid stigmatizing "multiple personality" framing and check project/package-name availability before choosing.
- [ ] Keep `Mosagent` in the naming candidates: it preserves the current sound while making the multi-agent association more visible. Revisit alongside Mosaic, Engram, Connectome, and availability research later.
- [x] Add `source show --align-source` or equivalent rendered-preview mode that shows how a source node would align under a manifest heading, including shifted heading levels and optional `--under`/`--heading`.
- [x] Add source-directory expansion for path-like refs such as `cdint:engineering`, `cdint:engineering/`, and `cdint:engineering/*`, or provide a better diagnostic that lists available descendant refs: implemented descendant diagnostics, not wildcard expansion.
- [x] Improve diagnostics for failed source refs by suggesting close matches and descendant headings when the user names an organizational directory rather than a heading path.
- [ ] Add manifest shorthand for reusable base/source aliases inside one subtree, while preserving the readable compact `- Heading: source:path` and `- Heading: [children...]` forms.
- [ ] Explore relative child references under a declared base, for example a future explicit form that can reference `./role` and `./source-of-truth` without repeating the source prefix.
- [ ] Add direct source/local editing workflows: safe `source edit`, copy-on-write local override editing, and explicit shared-source editing for trusted libraries.
- [ ] Reorganize source libraries into clearer top-level groups: core/code-principles, process/decision, stack/langs, communication/public-prose, domains, and CDINT/PromiseGrid as a non-universal domain source.
- [ ] Clarify or rename `Corpus Variants` sections. They currently mean "differences observed in the captured repo-agent corpus"; the label is unclear and appears too widely in rendered output.
- [ ] Add hierarchy-aware section spacing in rendered output or preview output when dense heading transitions make generated AGENTS.md hard to scan.
- [ ] Add accessibility support for a TTS-friendly second stream: optionally write a concise spoken version of substantial chat/tool output to a configured sidecar file.
- [ ] Add related-node metadata or "see also" support so adjacent styles and modules can point at each other without forcing users to discover relationships manually.
- [ ] Add OR/XOR heading choice groups for sibling modules that are alternatives instead of additive guidance. Use this for mutually exclusive styles, dependency managers, migration strategies, validation levels, or repo ownership models where the generator should force a deliberate choice and warn on contradictions.
- [x] Fix coverage tree rendering/markers where source boundaries and nested trees are visually confusing, especially around transitions between sources.
- [ ] Decide whether communication styles should live under `communication/personas`, `communication/style`, or a separate source. Preserve persona examples so users can understand the expected voice quickly.
- [ ] Add repository-specific prose support in manifests, either as inline YAML block scalars, local source snippets, or a dedicated repo overlay file, so prompts can include concise project identity and invariants without forcing every repo fact into shared libraries.
- [ ] Add prompt-size/over-composition diagnostics that warn when a generated AGENTS.md pulls broad policy subtrees into a small repo and suggest narrower child nodes.
- [ ] Add formatting preferences for generated Markdown, including sentence/paragraph-oriented source lines and avoiding arbitrary hard wraps when the user wants display wrapping to be handled by the viewer.
- [ ] Add stale-doc review support: compare docs claims against code-visible surfaces where possible, or at least provide a checklist for public-surface docs such as sockets, commands, prefs/state files, service lifecycle, and implemented roadmap items.
- [ ] Add `mogent handoff`: generate a compact, agent-readable project-state summary from explicit repository and workspace data. Include current manifest/output state, source and coverage summaries, Git state, maintained TODO/design pointers, and next actions; never scrape chat history or protected/private corpora.
- [ ] Design and publish Mogent's Go API so the CLI is a thin client of the same
  supported workspace operations available to other callers. Resolve package,
  compatibility, effects, error, and transaction contracts before moving code
  out of `internal/`; see `docs/PUBLIC-API-AND-CLEANUP-PLAN.md`.
- [ ] Execute the focused code-cleanup backlog in
  `docs/PUBLIC-API-AND-CLEANUP-PLAN.md`: CLI decomposition, one atomic writer,
  ignored-error audit, local-source symlink safety, deterministic tool vars,
  public examples, and stale dogfood/doc cleanup.
- [ ] Research how each major agent ecosystem manages skills before designing
  `.agents/skills/` support. Compare discovery scopes, formats, trust,
  dependencies, pinning, precedence, and interoperability using current primary
  documentation and small fixtures.

## Dogfood Feedback - Vroca setup, 2026-08-06

Context: user is setting up `/home/qix/dev/omnicortex/vroca_tts`, a Nix-backed
Python prototype intended to become a Rust CLI plus GUI app. The user is trying
to build `agents.yaml` by browsing local `cdint` and `personal` libraries from
the CLI.

Observed friction:

- `source list --tag-search lang`, `--tag-search python`, and
  `--tag-search python --metadata` returned `No source nodes matched` even
  though refs such as `personal:lang/python` exist. The command searches tags
  only; most seed library files currently have no frontmatter tags. This is
  correct by implementation but misleading for a newcomer.
- `source list --sort priority --metadata` prints a large dense stream of
  `Metadata: none`, which makes it hard to discover useful modules.
- At the time, `coverage --tree` focused on unused refs and did not clearly show
  included, inherited, or excluded nodes. DI-folar resolved this with the
  unified state-overlay tree.
- The user tried `mogent add cdint:engineering` and `mogent add
  cdint:engineering/`. Organizational directories were not selectable then;
  DI-folar now makes `cdint:engineering` a subtree reference.
- `--manifest` is unclear to new users. It means "use this agents.yaml instead
  of `./agents.yaml`", but the help text does not explain when or why to use it.
- The current source-library structure makes `CDINT And PromiseGrid` appear as a
  universal peer of general engineering modules. That is misleading; CDINT and
  PromiseGrid should be a domain/source, while engineering/process/language
  modules should be easier to browse as general reusable material.
- `Corpus Variants` sections are unclear in rendered output. They preserve
  observed differences from the captured source-agent corpus, but the label and
  ubiquity make them feel like noise.
- The compact manifest form is readable:

  ```yaml
  doc:
    - Identity:
        - Role: cdint:shared-baseline/identity/role
        - Source Of Truth: cdint:shared-baseline/identity/source-of-truth
  ```

  But it becomes repetitive when several child refs share the same source base.
  Need a shorthand that preserves readability without inventing a confusing
  mini-language.

Requested / candidate improvements:

- Add a beginner-friendly `init` or builder walkthrough with recommended
  ordering and templates. For a Vroca-like repo, a template might suggest:
  Identity, Project Direction, Instructions, Constraints, Stack,
  Communication, Format.
- Add presets/templates rather than making users discover everything from a huge
  source list. Presets should be editable manifests, not opaque hidden config.
- Add `library.yaml` or equivalent source-level metadata for recommended
  placement/order, update/review date, requirements, urgency/risk, and preset
  membership.
- Add source-search over refs/headings/TLDR/content, not only tags.
- Add useful empty-state messages and per-command help hints.
- Add tab completion for source refs.
- Add `source show --align-source` to preview render alignment.
- Add path/glob expansion or descendant suggestions for organizational paths
  such as `cdint:engineering/*`.
- Add coverage states and color/text markers.
- Add local/source edit flows.
- Revisit library top-level hierarchy and the naming of `Corpus Variants`.

Follow-up review from generated Vroca `AGENTS.md`:

- The generated prompt was about 485 lines for a small Python/Nix repo and read
  as a strong policy library but a weak repo-specific prompt. The right dogfood
  target for Vroca is closer to 60-100 lines.
- Useful retained ingredients: preserve user changes, narrow diffs, read
  architecture before behavior changes, risk-scaled validation, runtime-artifact
  hygiene, separate Nix and Python dependency surfaces, TTS-friendly
  communication, and public-surface caution around socket protocol,
  preferences, and command-line interface.
- Poor fit: Go guidance, premature Rust guidance while `rust_impl/` is only a
  placeholder, Decision Intent/Decision Request/proquint/TODO infrastructure
  that does not exist in Vroca, generator refs rendered as meaningful prose,
  corpus-library commentary, persona nodes without selection rules, and
  unrelated white paper / slide / experiment / database / data-pipeline rules.
- Missing Vroca-specific prompt material: `docs/vroca.md` should be named as the
  design of record; prompt should identify Vroca as a TTS/assistive reading
  framework with a Python implementation and imminent Rust refactor; prompt
  should name public surfaces and lifecycle boundaries.
- Missing Vroca docs/code topics surfaced by review: daemon singleton
  ownership, stale socket recovery, malformed command behavior, systemd restart
  semantics, child mpv cleanup, and client/daemon version compatibility.

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

Next build steps: extract Tier 1 + Tier 2 as real library nodes from the corpus
text; implement rebuild milestone 1 (parse `agents.yaml` manifest → resolve sources →
render → validate) per DI-ralik; dogfood this repo's `AGENTS.md` through it.

## Design References

- docs/DESIGN.md - **Design of record** (DI-modun). Read first; supersedes the docs below where they disagree.
- docs/brainstorn.md - Conversation-derived brainstorm notes for category taxonomy, TUI selection/save behavior, and fork import workflow (YAML config section folded into DESIGN.md §3.3-3.4).
- docs/thought-experiments/TE-tavim-mogent-module-reference-model.md - Reference-model TE; block/id-per-node/wikilink framing superseded by DESIGN.md §3.1. Manual handle allocation approved by user because `tools/mint-handle` was unavailable at the approved paths.
- docs/thought-experiments/TE-bakom-mogent-tui-first-selector.md - TUI-first selector TE. Recommends a Bubble Tea `mogent tui` browser before expanding the module library.
- docs/thought-experiments/TE-tavim-mogent-module-reference-model.md - Follow-up TE for block references, presets, metadata, render order, and diff model. Manual handle allocation approved by user because `tools/mint-handle` was unavailable at the approved paths.
- docs/thought-experiments/TE-bakom-mogent-tui-first-selector.md - TUI-first selector TE. Recommends a Bubble Tea `mogent tui` browser before expanding the module library.
