# Configuration Redesign: First Implementation Slice

Status: **syntax API implemented and API-tested; catalog/resolver and CLI pending**.

- [Spec](../docs/proposals/MOGENT-FIRST-SLICE-SPEC.md)
- [Acceptance cases](../docs/proposals/MOGENT-FIRST-SLICE-ACCEPTANCE.md)
- [Conversation history](2026-10-14-212422-conversation-stack.md)
- [Current CLI inventory](../docs/CLI-COMMAND-INVENTORY.md)

## Objective

Prove the agreed small-language workflow on temporary consumers: authored plan
or direct build -> shared ordered sequence -> inspectable complete/partial result
-> safe installation. Preserve metadata/curation and editable local sources.
Reuse the compositor, Markdown parser, Git retrieval, and state/transaction rules.
Do not expand protocols, integrations, or language features during this slice.

## Review Gates

- [x] **R1 Contract review.** Owner accepted the combined review and starting the
  first slice, chose `mogent.lib`, and clarified optional trailing semicolons.
  The original review covers the proposed engineering-detail table
  in spec section 1: exact grammar, header/source records, generation checksum,
  provenance storage, and direct-local-build record workflow. Accepted conversation
  direction is not blanket approval of these new spellings.
- [x] **R2 New library grammar.** Owner chose B: include library syntax in the
  initial redesign, not an ongoing legacy-sidecar bridge. Review spec section 4's
  typed node/alias schema, optional role/content, physical annotations/curation,
  named presets, explicit guidance vars/values, reason requirements, relationships,
  and metadata ownership. Preserve old config files to interpret
  during migration; ongoing YAML maintenance is not a goal. Reuse content types
  and semantics, not a requirement to retain their old serialization.
- [ ] **R3 Generated/input records implementation verification.** Verify the
  reviewed serialization and saved-baseline recovery proposal
  before claiming fresh-checkout reproducibility or safe cross-checkout regeneration.
  The generated checksum proposal is integrity, not authenticity. Do not put the
  only pins/decisions into ignored state.

Use API implementation -> focused API tests -> CLI -> CLI behavior tests. Write
expected CLI behavior before wiring it; the acceptance cases are that initial
contract and can be encoded as executable tests in advance. Parser work is the
first approved slice. The owner has permitted retiring legacy commands/formats, not uncontrolled
deletion of live consumer files or unrelated implementation work.

## Delivery Sequence

### S1 — Parser, Syntax Tree, And Authoring Preservation

- [x] Implement syntax for format-identified plan/build headers, bindings, quoted/
  multiline strings, child lists, fields, exclude blocks, and replace alias.
- [x] Implement the reviewed new library declaration grammar using the same syntax
  machinery; source metadata is part of the initial redesign, not a later HCL bridge.
- [x] Normalize typed aliases to nodes and preserve explicit section-input values
  in formatted syntax; interpretation of physical/section payloads waits for the catalog.
- [x] Preserve source spans, comments, and ordered actions; distinguish option
  field order from action order. Unknown/duplicate syntax fails loudly.
- [x] Evaluate parser reuse against this grammar first. If custom parser is needed,
  keep the lexer/grammar small and avoid HCL2 evaluation or a new general language.
- [x] Retain exact original source for no-op preservation and implement deterministic
  canonical formatting with alias expansion and semantic string round trips.
- [ ] Implement actual source-range edits for authored mutation commands; formatter
  can move inline comments and is not an authoring editor.
- [x] Verify syntax portions of A01/A01a/A02/A03 and alias/comment/string round trips.
  Resolve undefined/unused names and binding declaration order in the resolver.

Implementation: `syntax/` has no extra dependencies or I/O. API implementation
preceded the external-package API tests; no new CLI behavior has been wired.
Parser is a bounded recursive-descent implementation because child-list shorthand,
immutable bindings, and the agreed string syntax do not match existing HCL/KDL
grammars. Reader reuse remains useful for legacy migration only, not parsing a
different syntax under a misleading format name.

### S2 — Catalog, Decisions, And Partial Inspection

- [ ] Load new library declarations and physical inventories into a shared catalog with
  qualified identities, order, metadata, bundles, and operation capabilities.
- [ ] Resolve exact/child/tag selections, scoped exclusions, curation/defaults,
  reasons, destinations, and requirements on the final set. Keep omissions and
  superseded choices, not just selected leaves.
- [ ] Produce a structured partial result for errors with node/range attachment,
  UNDECIDED and unplaced nodes. No materialization is possible with blockers.
- [ ] Preserve library/local/effective values and per-decision reasons for text/TUI.
- [ ] Verify A04-A06 and B01-B12 at API level and through new library parsing;
  track syntax-loader coverage independently from hand-built catalog tests.
- [ ] Verify L01-L13 for typed library parsing, private/whitelisted variable inputs,
  presets, content roles, unresolved defaults, and retained old-config isolation.

### S3 — Sequence Serialization And Integrity

- [ ] Resolve/capture pins, local node contents and complete metadata/recipe inputs.
  Reuse existing source retrieval; avoid implicit remote pin advancement.
- [ ] Emit deterministic canonical include steps, resolved text/destinations,
  management roots and decision provenance. No alias/control query remains unresolved.
- [ ] Detect generated-file modifications from portable records; authored files
  are never overwritten. Unknown provenance is a blocker, not permission to guess.
- [ ] Implement source/plan freshness comparison and direct-build authority rules.
- [ ] Verify A07/A10/A11 and C01-C07/C13, including moving checkout, missing state, and
  dirty/untracked bundle members. Do not use random IDs/times to identify steps.

### S4 — Assemble Through Existing Core And Transaction

- [ ] Lower physical inputs to existing copy contributions and logical/text inputs
  to rendered contributions, using explicit internal append where appropriate.
  Do not expose these engine switches as normal user syntax.
- [ ] Compose one virtual manifest without managing the whole workspace. Partition
  it into non-overlapping output owners such as root AGENTS.md and `.agents/skills`.
- [ ] Preserve complete bundle replacement/removals and immutable captured payloads.
  No symlinks, source escapes, script execution, or semantic config merges.
- [ ] Stage/install outputs + generated build + state with rollback; preserve
  authored sources, localized sources, and files outside managed roots.
- [ ] Verify A08-A12 and C08-C12, failure injection at each write/install boundary,
  existing state drift behavior, modes, empty directories, and document dependencies.

### S5 — CLI Workflow, Text Inspection, And Localization

- [ ] Introduce `--file`, `sequence`, default sequence-if-needed `build`, and
  source/build/result `show` views after shared APIs exist.
- [ ] Add result-default tree, node explanation, issue filter, and generated-build
  freshness messages. A partial view remains available even when build would fail.
- [ ] Implement normal initialization/source/config editing and shared localized
  include operation with dry-run; source/config/provenance authoring is transactional.
- [ ] Support proposed local source-update preview/accept for explicit input-hash
  refresh in authored builds; ordinary build never edits authored input records.
- [ ] Add `build --check` from the same resolver/result inspector; do not implement
  a separate CI-specific checker. CI wrapper/skill remains a later TODO.
- [ ] Verify D01-D05/E03 and the pilot CLI journey. Help explains plan -> sequence
  -> build, direct builds, editing versus installation, and recovery from blockers.

### S6 — Library Coverage, Recipes, Compact Output, And Cutover

- [ ] Finish R2 new-library coverage and executable fixtures for all curation/tag/
  preset/relationship behavior promised by the redesigned surface.
- [ ] Add scoped recipe inclusion and cycle diagnostics (D06-D08).
- [ ] Add compact serializer/alias output choices only after explicit form is
  correct; verify full semantic/ownership/provenance equivalence (E01/E02).
- [ ] Retire obsolete CLI/config paths after replacement capability mapping and
  fixture verification. Give actionable format errors; no silent parser fallback.
  Keep historical docs clearly historical and update README/help/completions.
- [ ] Preserve useful localization/coverage/editing behavior and public API
  boundaries; permitted legacy removal is not permission to lose intended features.

## Reuse And Necessary Adaptation

| Existing area | Reuse | Necessary adaptation |
|---|---|---|
| `content/` | Inventory, document views, bundle checks, captured manifests | Catalog annotations; user-natural include lowering; managed-root projection |
| `v2/plan.go` | Curation/default validation concepts | One precedence process, final-set validation, structured partial decisions |
| `v2/sources.go`, `sourcecache/` | Git pins/fetch/cache mechanics | Plan/build pin authority; local integrity records |
| `v2/build.go`, `v2/tree.go`, `state/` | Drift/layout/transaction behavior | Generated-build protection and root-file + directory transaction without whole-root ownership |
| `workspace/localize.go` | Copy-on-write/provenance/rollback concepts | Complete bundles, new config edits, retained baseline, source alias scoping |
| `internal/presentation`, `internal/navigator` | Tree symbols, navigation, pane ideas | Shared source/build/result inspection model, unresolved nodes; interactive UI follows text |
| `internal/cli` | Thin adapters and command-registry pattern | Restricted new vocabulary; no v1/v2 dispatch confusion |

Prefer focused public core APIs over putting syntax-independent behavior into CLI
handlers or binding the new executor directly to legacy config structs. Exact
package names and exported types belong to implementation review, not this task list.

## Verification And Stop Conditions

- Focused meaningful tests per slice, then required `tools/check` at integration
  gates. No need to run full Go checks for the documentation-only specification.
- Keep real external repo testing read-only until the pilot and source boundaries
  are chosen with the owner. Temporary fixtures are sufficient for these cases.
- Stop a dependent implementation if R1-R3 is unanswered or a fixture requires
  unsupported metadata; record the blocker rather than inventing semantics.
- Initial implementation is authorized. Commit only when requested.

## Deferred

Interactive result viewer (required later, not discarded), protocol/target-profile
headers, harness compatibility comparison, CI/skill wrappers, global installation,
output symlinks, source-to-source semantic merges, editable heading derivatives,
and build-to-plan inference. Preserve these in the conversation stack/TODO index.
