# TODO-tugur - V2 Composition Follow-Up

## Purpose

Finish the design and delivery work exposed by the first HCL v2 rollout:
selection-guidance vocabulary, section-selection shape, and composable tree
outputs. This is a follow-up to `TODO-pupab-hcl-v2-rollout.md`; it does not
reopen the completed plan/apply, Git-pinning, or basic tree-output work.

## Current Direction

- The owner-approved first delivery is implemented and documented in
  [Typed Content Composition](../docs/TYPED-CONTENT.md): named physical roots,
  optional parsed heading views, mixed copy/render outputs, explicit replacement
  and append, complete bundle copying, state-backed transactions, and directory
  update previews. Existing selectors/defaults and old output kinds remain valid.
  DI-juvih-first-slice records the approved additive grammar and API-first workflow.
- Keep `exclude` as the consumer-facing operation for suppressing a path from
  a lower tree layer. Do not expose the implementation term `whiteout` in HCL.
- Use `curate` as the v2 library attribute. Keep its broader terminology review
  open for a later deliberate rename, not as alternate accepted grammar.
- A per-child `default` is a candidate replacement for a `choose` branch's
  parent `defaults` map. The branch-level selection-guidance attribute remains:
  it describes how to treat the set of children; a child default records the
  library's recommendation for that one child.
- Earlier selection/composition analysis is recorded in
  [TE-jofit](../docs/thought-experiments/TE-jofit-unified-content-composition.md).
  It compares `content` / `node`, `markdown` / `dir-tree`, authored versus
  imported trees, and common selection. Its collision recommendation is now
  accepted as explicit replacement; lower-layer removal syntax remains open.
  Existing `exclude` filters only its source.
- The typed model refining that baseline is accepted through
  [DR-juvih / DI-juvih](../DR/DR-juvih-typed-content-composition.md), based on
  [TE-gunak](../docs/thought-experiments/TE-gunak-typed-agent-content-composition.md)
  (decided, refined; 2026-10-09) and the
  [ecosystem report](../reports/Universal%20agent%20content%20composition.md).
  Adopt typed physical `dir`/`file` nodes with optional heading views and separate
  role/bundle annotations. Start with intact copy and explicit render operations,
  mixed artifacts, and explicit collision replacement. Naming, exact bundle and
  replacement rules, parser/addressing contracts, and delivery order stay open.
- TE-gunak's dated execution refinement includes SKILL.md heading views and
  twelve concrete tabletop cases. Skills may expose document subtrees without
  losing their intact bundle/file view; extraction, installation, and editing
  are distinct operations. The model/collision direction is approved; remaining
  syntax and staging questions are listed in DR-juvih.

## P0 - Stabilize The Current Rollout

- [x] **V2C-1 Fix macOS tree-output test paths.** Make output containment
  validation accept a real temporary directory whose ancestor `/var` is a
  system symlink, without permitting an output target or source entry to
  traverse a user-controlled symlink. Add a regression test.
- [x] **V2C-2 Verify sibling QMR lookup in `tools/v2-stories`.** The existing
  script already resolves the sibling checkout relative to the Mogent repo.
  The observed failure was missing sidecar metadata, addressed by V2C-3.
- [x] **V2C-3 Reconcile dogfood fixtures.** Update `cdint-demo-lib` and the
  active QMR sidecar to the currently accepted v2 sidecar vocabulary, or revise
  the scenario fixture to make its intended compatibility target explicit.
  Do not promote `library-tmp.mogent.hcl` metadata without review.

## P1 - Thought Experiments

- [x] **V2C-4 Run a selection-surface TE.** Analysis in TE-jofit; owner answers
  to DF-1/DF-2/DF-3 are pending. Compare the implemented object
  expression (`select = { ... }`) with nested `select { section "..." {} }`
  blocks. Test deep explicit selections, `else`, `accept_defaults`, typo
  diagnostics, generated init templates, and HCL readability.
- [x] **V2C-5 Run a tree-composition TE.** Analysis in TE-jofit, refined by
  TE-gunak. Explicit replacement is approved; detailed exclusion, placement,
  and replacement semantics remain open. Exercise a company base, QMR skill
  library, repository overrides, file/file conflicts, file/directory conflicts,
  directory replacement, and an exclusion with no replacement. Decide whether
  `exclude` belongs on a source block, requires a nested block, or needs an
  explicit replace shorthand.
- [x] **V2C-6 Run a section-root/path-layout TE.** Analysis in TE-jofit;
  detailed adapter/path rules remain part of the pending DR. Evaluate named logical
  section roots with inheritable `basedir` against the current single
  `content.markdown_root` model. Define containment, composition, and a
  migration path before changing the sidecar schema.

## P2 - Design Review And Decisions

- [x] **V2C-7 Resolve `DR/DR-nalan-offer-attribute-name.md`.** Record `curate`
  as the current key, its read-aloud form, diagnostics, and relationship to
  each value. Preserve a later rename as a new decision, not parser fallback.
- [ ] **V2C-8 Create a tree-composition DR from V2C-5.** Lock layer order,
  conflict reporting, `exclude` semantics, directory replacement semantics,
  mode preservation, and state/drift behavior. The DR must distinguish the
  HCL vocabulary from any internal overlay/whiteout terminology.
- [ ] **V2C-9 Create a selection-and-defaults DR from V2C-4.** Decide object
  expressions versus nested selection blocks and parent `defaults` versus
  child `default`. Define invalid combinations and migration behavior.
- [ ] **V2C-10 Create a logical-section-roots DR from V2C-6.** Decide whether
  named content roots, directory imports, and inherited `basedir` are worth
  adding; distinguish logical identity from physical/output paths.
- [x] **V2C-18 Resolve typed content and mixed-artifact design direction.**
  Accepted in DR-juvih / DI-juvih: physical nodes plus optional heading views,
  independent output shape and copy/render operations, and orthogonal roles/bundles.
  Headings do not own arbitrary files and directories do not inherit hidden format
  rules. Exact grammar and adapter sequencing remain V2C-21 work, not completed
  implementation or approval of example spellings.
- [x] **V2C-19 Deliver first opaque-import and bundle policies.** Complete
  file/bundle copying preserves bytes and frontmatter; partial skill installs and
  exclusions are rejected. Explicit replacement removes stale members, and plan
  reports operation/source/replacement history. Parsed views render derivatives
  and check explicit local-link dependencies. Richer annotations/derivatives and
  semantic configuration merging remain deferred.
- [ ] **V2C-20 Define context identity and tool export projections.** Use
  TE-gunak DF-6/DF-7 to retain global/project/nested/worktree origins, applicability,
  explicit build order, and receiving-host discovery as distinct concepts.
  Avoid invisible host merges and implicit main-worktree inheritance. Prioritize
  documented skill/instruction paths, qualify vendor-native targets by version,
  and keep the broader `.agents` protocol an optional proposal projection;
  arbitrary `.agents` storage does not imply standard configuration support.
- [x] **V2C-21 Specify and verify the first typed-composition delivery.**
  DR-juvih, DI-juvih-first-slice, and docs/TYPED-CONTENT.md distinguish delivered
  grammar from future proposals. Unit/usage cases cover mixed manifests,
  copy/render applicability, bundle replacement, explicit provenance, rollback,
  byte/mode/directory drift, heading ambiguity, and directory-update previews.
  Global/tool scope projections and generalized migration remain backlogged.
- [x] **V2C-22 Execute typed-content TE with skill heading subtrees.** Recorded
  in TE-gunak's 2026-10-09 refinement (R1..R12). These are reasoning cases;
  no proposed grammar or imported-heading implementation was executed.

## P3 - Specification And Documentation

- [ ] **V2C-11 Update the v2 specification after the P2 decisions.** Keep one
  canonical grammar, semantics, examples, validation table, and migration
  rules. Remove superseded syntax from examples rather than supporting it
  indefinitely.
- [ ] **V2C-12 Update user-facing documentation and templates.** Cover
  `plan`/`apply`/`update`, pinning in `mogent.hcl`, the final selection syntax,
  `curate` or its decided replacement, and tree composition with `exclude`.
  Correct the presentation draft separately; it is user-authored and currently
  describes superseded lock-file and build/dry-run behavior.
- [ ] **V2C-13 Document the library-author migration.** Include sidecar
  validation, QMR metadata promotion, source browsing output, and how existing
  `offer`/`defaults` files move to the finalized grammar.

## P4 - Implementation

- [ ] **V2C-14 Implement the finalized selection and defaults grammar.** Make
  loading, planning, diagnostics, source browsing, init templates, and tests
  agree. Reject only grammar that the migration command can safely rewrite.
- [x] **V2C-15 Deliver ordered typed directory composition.** `dir-tree`
  outputs accept multiple explicit contributions and `into` paths, disjoint
  directory union, requested subtree replacement, and explicit rendered append.
  Plan shows member and generated-document changes plus provenance; captured
  manifests stage transactionally and track complete bytes/modes/directory layout.
- [ ] **V2C-16 Extend safe tree exclusions to overlays.** Single-source
  exclusions are implemented relative to the declared source tree: exact files
  and recursive directories are filtered from the managed output. Mount-relative
  lower-layer removal and originating-layer reporting still need V2C-8/V2C-15.
- [ ] **V2C-17 Add an explicit v2 migration command.** Only after the grammar
  is final: rewrite `offer` and parent `defaults` sidecars/configurations with a
  preview and explicit acceptance. Do not add parser fallbacks for abandoned
  v2 syntax.

## Backlog - Secondary Integrations

Owner direction: defer the ecosystem report's "next" and compatibility
integrations while the common content model is resolved. Complete native files
may still be imported/copied opaquely; adapter development and conversion wait.

- [ ] **V2C-23 Vendor rule adapters.** Claude rules, Cursor `.mdc`, and Copilot
  scoped instructions; preserve each host's applicability and frontmatter dialect.
- [ ] **V2C-24 Custom-agent/profile adapters.** Claude, OpenCode, and Copilot
  profiles; do not infer cross-tool conversion from common Markdown extensions.
- [ ] **V2C-25 Commands and prompt integrations.** Preserve host syntax and
  support/version limitations before introducing export profiles.
- [ ] **V2C-26 Automated hooks/settings registration.** Scripts and complete
  configuration files remain opaque copies; copying does not activate a hook.
  Semantic TOML/YAML/JSON/JSONC merging remains deferred.
- [ ] **V2C-27 Global installation and cascade projections.** Keep origin/context
  identity in the core design; defer automatic home writes, host-specific scope
  loaders, and versioned agentsstandard cascade emulation.
- [ ] **V2C-28 Editable heading derivatives.** Explicitly modify selected
  headings inside a skill/document while preserving frontmatter, remaining content,
  dependencies, and provenance; separate from read-only heading extraction.

## Verification Gates

- [ ] `tools/check` passes on macOS and Linux. macOS passed for this delivery;
  Linux execution remains a CI follow-up.
- [x] `tools/v2-stories` runs against the sibling QMR library discovered from
  the repository layout.
- [x] Focused API and CLI usage tests cover the first-delivery grammar,
  collision/exclusion behavior, parsed headings, metadata drift, and rollback.
- [x] A temporary consumer can `plan`, review an exact tree/output diff, and
  `apply` without unmanaged files or direct edits being overwritten.
