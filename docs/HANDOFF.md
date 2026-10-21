# Mogent Handoff

**Current conversation: [NOW](convo/head.md).** This handoff is product
navigation/reference; older resume language below is not a competing live checkpoint.

Updated: 2026-10-21; source-oriented native browsing and upstream-change simulation recorded.

This file is the stable navigation entry point for an agent handoff. It does
not duplicate task state, design decisions, or a workstream's detailed resume
instructions.

## Current Dogfood Checkpoint

The detailed [Teacher_Bot native walkthrough handoff](dogfood/2026-10-20-teacher-bot-native-walkthrough-handoff.md)
records the current consumer flow, installed outputs/state, exact verification,
working-tree distinctions, deferred decisions, and safe re-entry. Read it after
the NOW panel when resuming the native plan/build dogfood work.

The current follow-up uses `mogent source show` for the native catalog and
`mogent show` for selected contributions. An isolated QMR/CDINT change/add/remove
simulation confirms that output state detects changed effective outputs but does
not yet identify authored upstream inputs; [NOW](convo/head.md) owns that active
design decision.

## Sources of Authority

- [`convo/head.md`](convo/head.md) owns the live conversation resume point;
  its linked work and decision records own the conversation queue and choices.
- [`TODO/TODO.md`](../TODO/TODO.md) owns the product backlog and related owner
  decisions, statuses, and blockers.
- [`TODO/PICKUP-2026-09-04-agents-library-dogfood.md`](../TODO/PICKUP-2026-09-04-agents-library-dogfood.md)
  is historical checkpoint context. It is not the current curation entry point.
- [`DESIGN.md`](DESIGN.md) is the product design of record.
- [`library-curation/archive/LIBV2-PROTO-PROGRESS-REPORT.md`](library-curation/archive/LIBV2-PROTO-PROGRESS-REPORT.md)
  preserves the historic libv2 problem, approach, organization, and open
  decisions.
- [`../libv2-proto/PROVENANCE.md`](../libv2-proto/PROVENANCE.md) and
  [`../libv2-proto/WORKLIST.md`](../libv2-proto/WORKLIST.md) track source intake.
- Git history records authorship from human developers only; commit messages
  cite no co-authors except human developers.

If a chat recap or this navigation file conflicts with Git or an authority
above, inspect the authority and current working tree.

## Product Workstream Reference

The owner approved initial implementation, chose `mogent.lib`, clarified
separators, and requested API implementation -> API tests -> CLI -> CLI behavior
tests. The first bounded syntax API slice is now implemented in `syntax/`:
Parse, source-located ordered AST/comments, explicit string substitution, and
deterministic formatting/alias expansion. External-package API tests pass.
`nix develop -c tools/check` passes, including full tests/vet/deadcode/errcheck.

The catalog/resolver, generated builds, local/Git source pinning, assembly, basic
CLI (`sequence`, `show`, `build --file`, `edit include --localize`) and focused
behavior tests are now delivered. `build --check` compares managed local outputs.

**Next: finish S2/S4/S5 gaps**: optional/opt-in broad curation and retained omission
provenance, complete output/state rollback, source-update acceptance/provenance baseline
records, and richer result inspection/help. Scoped tag selection, include-local
exclusions, reason enforcement, final requires/conflicts/exclusive-group validation,
and mode/empty-directory output drift checks are delivered. Do not describe the current CLI as a replacement
for legacy default commands until the cutover slice verifies those workflows.

Delivered no-argument behavior: when invoked without `--manifest`/`--file`,
`mogent build` chooses a local `mogent.plan` first, then `mogent.build`; otherwise
it preserves legacy manifest behavior. `sequence` refuses to replace an authored
or manually edited generated build. `build --file` remains the unambiguous route
for scripted workflows. These behaviors have focused CLI tests.

1. [First-slice spec](proposals/MOGENT-FIRST-SLICE-SPEC.md): accepted direction plus
   a clearly marked proposed-engineering-detail table.
2. [Acceptance cases](proposals/MOGENT-FIRST-SLICE-ACCEPTANCE.md): pilot fixture,
   concrete success/failure/rollback cases; syntax portions are tested, others pending.
3. [Implementation sequence](../TODO/TODO-config-first-slice.md): R1-R3 review gates
   and S1-S6 slices mapping onto reusable code.

Owner selected **B: include the new library format in the initial redesign**.
Typed-node refinements are now in spec section 4 and acceptance L01-L13:
`md-section` for guidance and `md-file` for intact Markdown, with `dir`/`file`,
optional role/content, curation/presets, and explicit guidance vars/values.
The owner accepted the final review and starting the first slice. Keep existing config files for interpretation during
migration/reference, but do not carry YAML as a maintained runtime format. This
supersedes the earlier proposed HCL-library bridge. Library filename is `mogent.lib`;
no special `skill` keyword is needed. No commits have been made during this work.

Key accepted direction: custom small language; one `include` family; `replace`
alias normalizes to `include` with a flag; no copy alias or ordinary copy/render/
append switches. Plans compile through the same build representation/executor as
direct authored builds. Default build sequences if needed, protects authored and
edited generated files, and materializes root AGENTS.md plus selected `.agents/`
artifacts. Localize once into an editable local sublibrary; later builds preserve
edits. Git commits/local complete-node hashes identify inputs. Tags/curation,
reasons, defaults, library order, and unresolved result callouts remain core.

The [conversation stack](../TODO/2026-10-14-212422-conversation-stack.md) preserves
the historical refinements and parked questions; its old resume pointers are
reference only. Use [NOW](convo/head.md) for current work.
Keep responses short, handle one issue at a
time, and do not restart vocabulary/format research. The
[walkthrough](proposals/MOGENT-LANGUAGE-WALKTHROUGH.md) preserves example evolution;
the [CLI inventory](CLI-COMMAND-INVENTORY.md) distinguishes delivered behavior.
Legacy commands/formats may be retired after equivalent core workflows are verified.
Protocol headers, standards/harness comparison, CI/skill wrappers, symlinks/global
targets, reverse inference, and the richer interactive UI remain follow-ups.

### Existing Implementation Baseline

The delivered HCL v2 rollout is tracked in
[`TODO/TODO-pupab-hcl-v2-rollout.md`](../TODO/TODO-pupab-hcl-v2-rollout.md)
against [`proposals/MOGENT-HCL-V2-SPEC.md`](proposals/MOGENT-HCL-V2-SPEC.md).
`plan`, `apply`, `update`, commit pinning, curation, tree outputs, and v2
`source list`/`show` are implemented; `tools/v2-stories` exercises the
acceptance scenarios against the sibling QMR library. Open owner decisions:
whether `curate` becomes required
on every branch.

Typed composition is now delivered under DR-juvih / DI-juvih-first-slice:
named physical roots, optional parsed heading views, mixed `dir-tree` copy/render
outputs, explicit replacement/append, complete bundle copying, directory-update
previews, and byte/mode/directory drift protection. Core APIs/unit tests precede
CLI usage tests in the atomic commit sequence. See [TYPED-CONTENT.md](TYPED-CONTENT.md)
and the runnable `testing-ground/typed-content-demo/` fixture. Generalized selectors,
annotations, scoped earlier-layer deletion, global/vendor adapters, and editable
heading derivatives remain follow-up work.

CDINT source review remains paused pending Steve curation. For a request such
as "I'm Steve, I'm here to do the content curation that Quincy requested.
Please guide me through it.", read the [Steve library curation
handoff](STEVE-LIBRARY-CURATION-HANDOFF.md) first, then the [local curation
landing page](library-curation/README.md), [open questions](library-curation/open-questions.md),
and [wb review material](wb.md). Pause before rendering or promoting drafts.

Personal live-library dogfooding is owned by the sibling QMR Agents Library.
Its design workbench, templates, archive, and provenance are not CDINT policy.

The root `agents.yaml` and generated `AGENTS.md` remain intentionally absent:
this repository owns Mogent product code, not a selectable policy library. Do
not restore the deleted v1 guide or create a product-root manifest. Use a
reviewed consumer repository or the QMR library template for dogfood work.

Module include syntax remains undecided. `TE-nufad` narrowed the design space;
it did not authorize implementation. Relationship validity and merge
precedence also remain separate open questions.

## Document Convention

- `TODO/TODO.md`: durable, canonical work state.
- `TODO/<datetime>-conversation-stack.md`: indexed owner discussion, single-item
  resume point, preferences, proposals, and parked questions; not a grammar spec.
- `TODO/PICKUP-<date>-<workstream>.md`: volatile checkpoint for one active
  workstream; it may point to a detailed checkpoint in the repository that owns
  the active work.
- `docs/HANDOFF.md`: short, stable map to the current authorities and pickup.
