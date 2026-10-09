# TODO-tugur - V2 Composition Follow-Up

## Purpose

Finish the design and delivery work exposed by the first HCL v2 rollout:
selection-guidance vocabulary, section-selection shape, and composable tree
outputs. This is a follow-up to `TODO-pupab-hcl-v2-rollout.md`; it does not
reopen the completed plan/apply, Git-pinning, or basic tree-output work.

## Current Direction

- Keep `exclude` as the consumer-facing operation for suppressing a path from
  a lower tree layer. Do not expose the implementation term `whiteout` in HCL.
- Use `curate` as the v2 library attribute. Keep its broader terminology review
  open for a later deliberate rename, not as alternate accepted grammar.
- A per-child `default` is a candidate replacement for a `choose` branch's
  parent `defaults` map. The branch-level selection-guidance attribute remains:
  it describes how to treat the set of children; a child default records the
  library's recommendation for that one child.
- Tree composition is an ordered overlay: each source mounts opaque files into
  an output tree, later sources win path conflicts, and `exclude` removes a
  lower-layer path at that mount before later files are applied.

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

- [ ] **V2C-4 Run a selection-surface TE.** Compare the implemented object
  expression (`select = { ... }`) with nested `select { section "..." {} }`
  blocks. Test deep explicit selections, `else`, `accept_defaults`, typo
  diagnostics, generated init templates, and HCL readability. File a new
  `docs/thought-experiments/TE-<handle>-selection-surface.md` before deciding.
- [ ] **V2C-5 Run a tree-composition TE.** Exercise a company base, QMR skill
  library, repository overrides, file/file conflicts, file/directory conflicts,
  directory replacement, and an exclusion with no replacement. Decide whether
  `exclude` belongs on a source block, requires a nested block, or needs an
  explicit replace shorthand. File a new
  `docs/thought-experiments/TE-<handle>-tree-composition.md`.
- [ ] **V2C-6 Run a section-root/path-layout TE.** Evaluate named logical
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
  named `sections "agents"` roots and inherited `basedir` are worth adding.

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
- [ ] **V2C-15 Implement finalized tree overlays.** Support multiple ordered
  tree sources and mount paths; report every shadowed, excluded, added, changed,
  and removed path in `plan`; materialize the complete result transactionally;
  retain managed-state hashes for the final tree.
- [ ] **V2C-16 Extend safe tree exclusions to overlays.** Single-source
  exclusions are implemented relative to the declared source tree: exact files
  and recursive directories are filtered from the managed output. Mount-relative
  lower-layer removal and originating-layer reporting still need V2C-8/V2C-15.
- [ ] **V2C-17 Add an explicit v2 migration command.** Only after the grammar
  is final: rewrite `offer` and parent `defaults` sidecars/configurations with a
  preview and explicit acceptance. Do not add parser fallbacks for abandoned
  v2 syntax.

## Verification Gates

- [ ] `tools/check` passes on macOS and Linux.
- [ ] `tools/v2-stories` runs against the sibling QMR library discovered from
  the repository layout.
- [ ] New focused tests cover every finalized grammar rule, diagnostic, tree
  conflict, exclusion, state/drift case, and transaction rollback.
- [ ] A temporary consumer can `plan`, review an exact tree/output diff, and
  `apply` without unmanaged files or direct edits being overwritten.
