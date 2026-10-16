# First-Slice Acceptance Cases

Status: **behavior specification; syntax API portions now tested, remaining cases pending**. These verify the proposed
[first-slice contract](MOGENT-FIRST-SLICE-SPEC.md). CLI spellings are proposed.
Keep behavior assertions independent of incidental help/diff formatting.

Delivered tests: `syntax/syntax_test.go` covers separators/optional trailing
semicolons, library/plan/build grammar, source ranges/Unicode columns, comments,
bindings/string evaluation, literal/multiline text, invalid expression syntax,
closed option names, duplicate declarations, depth bounds, typed/replace alias
expansion, and deterministic parse/format round trips. A02's source lookup and
binding-use/unused checks remain resolver work. No new CLI acceptance journey is
implemented yet; keep these expected behaviors as the contract before CLI wiring.

## Pilot Fixture

Use temporary consumer/library directories. No tests write to sibling live repos.
Create local Git repositories for pinned-source tests; no public network required.

```text
team-library/
  mogent.lib                 typed schema from spec section 4
  guidance/workflow.md
  skills/review/SKILL.md
  skills/review/scripts/team-check.sh
  skills/review/references/checklist.md
  skills/testing/SKILL.md
  skills/deployment/SKILL.md
consumer/
  mogent.plan
  .mogent/local/skills/review/SKILL.md
```

Logical guidance root contains Workflow, default-included when its offered
collection is accepted. Physical skill root order is review, testing, deployment;
curation/defaults in the normalized test catalog recommend review/testing and
exclude deployment. Owner selected new library syntax in the initial redesign.
Use the complete typed library example in spec section 4, with `md-section`, `dir`,
child annotations/defaults, and the recommended preset. Run these cases through
the new library parser and CLI as well as catalog APIs. Hand-built catalogs alone
do not establish library-format coverage. Retain old HCL/YAML configs as migration
evidence, not the fixture's maintained runtime format.

Team review is a complete valid skill and includes a script/reference; local
review is a complete valid skill with different text and no dependency on those
team assets. Treat capitalization/dialect as this fixture's Agent Skills input,
not DotAgents conformance. Workflow starts with `# Workflow`, contains a `${repo_name}`
placeholder in its body, and declares that variable in the library. Preserve
physical payload bytes, including frontmatter; render Workflow with Example bound.

Proposed plan:

```hcl
plan { format = 1; }
let repo_name = "Example"

source "team" { local = "../team-library"; }
source "local" { local = ".mogent/local"; }

include "team:skills" {
  accept_defaults = true
  exclude = ["deployment"]
}
replace "local:skills/review"

include {
  text = """
    # ${repo_name}

    Follow the repository's development conventions.
    """
}
include "team:guidance/workflow"
```

Expected output:

```text
AGENTS.md                     generated introduction, then Workflow hierarchy/body
.agents/skills/review/
  SKILL.md                    exact local bytes
.agents/skills/testing/
  SKILL.md                    exact team bytes
```

Absent: deployment, team review script/reference. Record that team review was
selected then replaced, not silently deduplicated. Explicit exclusions and accepted
defaults remain in explanations. `.mogent/local` remains unchanged by building.

## Core Vertical Slice

| ID | Given / action | Required observable result |
|---|---|---|
| A01 | Parse multiline and equivalent semicolon one-line instructions | Same ordered AST semantics, correct ranges for each spelling; unknown syntax is diagnosed |
| A01a | Single field before closing brace, multiple same-line fields, comma-separated array | No required trailing semicolon; semicolon needed between same-line fields; commas only separate array values |
| A02 | Duplicate binding/field, unknown alias, forward/unknown binding reference | Blocking, located diagnostics; no files written |
| A03 | Normal interpolation, escaped `${...}`, literal triple-single text, CRLF/indent variants | Exact documented text; no clock/env/CWD influence; unused name warns |
| A04 | Pilot selection through catalog API | Exact member set, action order, reasons and placement; no per-file registry |
| A05 | Explicit child-list selects review/testing with reversed textual list | Library-sidecar member order retained; unknown/duplicate child errors |
| A06 | `replace` alias versus `include { replace = true; }` | Same resolved actions, errors, target manifest, management roots, explanation |
| A07 | Sequence pilot twice with identical inputs | Byte-identical generated build, no timestamps/randomness; canonical output has include/options, exact destinations, pins/hashes, no unresolved defaults |
| A08 | Build pilot after sequencing | Expected files/bytes/modes; no stale team script/reference; source payloads unchanged |
| A09 | Repeat build | No duplicate document text, no extra files, stable expected result |
| A10 | Same concrete instructions in authored build, no plan | Same result as generated build; authored bytes unchanged |
| A11 | `sequence` standalone | Build written, no workspace outputs installed; on invalid plan old build remains unchanged |
| A12 | `build` with plan only | Automatic sequence then assembly; one failed assembly restores generated build/state/outputs |

## Decisions, Curation, And Inspection

| ID | Given / action | Required observable result |
|---|---|---|
| B01 | Include a choose branch without full decisions or default adoption | UNDECIDED placeholders at affected nodes; independent resolvable nodes still visible; no sequence/build installation |
| B02 | Add `accept_defaults = true` with complete defaults | Choices resolved, each default origin visible; no repeated prompt |
| B03 | Tag matches review, same include explicitly excludes review | Review absent from this contribution; tag/default cannot resurrect it |
| B03a | Tag-only include scopes to skills with review/testing matching Go, deployment not matching | Only review/testing candidates; root scope does not silently add deployment; bundles remain complete |
| B04 | Earlier include contributes review, later include excludes its own review | Earlier review remains; source filtering is not global deletion |
| B05 | Exclude foundation with empty/no reason; repeat with reason | First blocks; second records reason and proceeds if other constraints are satisfied |
| B06 | Ordinary opt-in omitted, explicitly selected; experimental selected without/with reason | Omission excluded, ordinary explicit allowed; experimental blocks without reason, records supplied reason otherwise |
| B07 | Requirement selected later, required node excluded, conflicting group members | Later selection satisfies presence; absent/conflicting final set blocks; no auto-inclusion or dependency-order inference |
| B08 | Replace a source identity required elsewhere with unrelated local identity | Placement does not falsely satisfy identity requirement; show affected requirement and replacement |
| B09 | Source unavailable or destination unknown | Partial source/result tree retains errors; unknown placement under Unplaced contributions; no invented file |
| B10 | Inspect selected review node | Distinguish library metadata, accepted default, local replacement, effective destination/source and stale members removed |
| B11 | Two includes choose same physical destination without replacement, even equal bytes | Located collision, no writes; explicit replacement resolves it |
| B12 | Select physical SKILL.md or exclude script inside an installed bundle | Incomplete bundle diagnosed; whole-bundle include/exclusion works |

## Sources, Generation Protection, And Installation

| ID | Given / action | Required observable result |
|---|---|---|
| C01 | Git ref moves after initial sequencing | Ordinary resequence/build uses recorded commit; update preview shows change; only explicit accept advances pin |
| C02 | Local bundle member changes bytes/path/mode or adds/removes an empty directory | Recorded node integrity changes; old build check fails; plan build can refresh unchanged generated file |
| C03 | Untracked supporting script or dirty tracked file | Included and hashed; source status explains dirty/uncommitted inputs; no false clean-checkout reproducibility claim |
| C04 | Move equivalent checkout to another absolute directory | Same sequence/output; source roots and input identities are portable |
| C05 | Edit generated build comment, option, or origin to authored | Comment/option edit blocks regeneration; authored origin blocks regardless of author identity; no flag implicitly discards edits |
| C06 | Fresh checkout with no state and valid generated-file checksum | Expected result reconstructible from committed inputs; missing state alone does not mark generated file edited |
| C07 | Generated file missing/malformed integrity record, or unrelated plan + authored build | No silent authority choice/regeneration; actionable explicit-file/recovery guidance |
| C08 | Existing output differs from installed state | Drift reported; refuse overwrite without explicit output force; force does not bypass invalid source/choice/paths |
| C09 | Root AGENTS.md and `.agents/skills` managed; unrelated `.agents/settings.json` exists | Build preserves settings; replacing review only removes review descendants |
| C10 | Source symlink, special file, traversal, or output symlink escape | Refuse before installation; external files unchanged |
| C11 | Fail staging/write/state persistence after one target installed | Restore every touched target, generated build, and state; surface rollback errors |
| C12 | Change source after manifest capture, before materialization | Materialize captured/verified payloads, not newly reread bytes |
| C13 | Authored build has local includes without records; preview/accept source update | Check explains missing records; preview writes nothing; accept writes reviewed input hashes; later normal builds never rewrite authored bytes |

## Localization And Recipe Reuse

| ID | Given / action | Required observable result |
|---|---|---|
| D01 | Localize team review with dry-run | Show whole local bundle path, include/reference edits, provenance; zero writes |
| D02 | Localize team review, then edit local SKILL.md and build twice | Local bytes preserved; output uses local version; upstream unchanged; no repeated localization |
| D03 | Local target exists or local alias points elsewhere | Error; no clobber; authoring transaction leaves source/config/provenance consistent |
| D04 | Simulate failure while writing localization config/provenance | Restore all newly created/modified authoring files; no output installation |
| D05 | Localize source lacking retrievable baseline | Retain captured baseline or explicitly report limited future comparison; a hash alone cannot recover content |
| D06 | Included recipe reuses caller's variable/source names | Private scope; caller values unchanged; declaration-relative paths resolve correctly |
| D07 | Direct/indirect recipe cycle or repeated independent inclusion | Full cycle diagnostic; repetition remains visible and conflicts normally, not silent dedup |
| D08 | Included recipe's file or relevant library metadata changes | Recorded inputs invalidate sequence; metadata/default updates appear in explanations |

## Representation And Deferred Integrations

- E01: explicit and compact builds compile to equal payloads, order, ownership,
  reasons, exclusions, and replacement history. A shorter line count is not enough.
- E02: permute map iteration, filesystem discovery order, and machine locale;
  result/sequence remain identical. No randomness/seed enters semantics.
- E03: non-mutating check detects stale plan/build/input/output without rewriting
  any of them; pinned cache retrieval is allowed only as documented. This is a
  core verification test, not a GitHub Actions workflow implementation.
- E04: target paths are root AGENTS.md and `.agents/` only; unsupported protocol
  header, output symlink request, semantic JSON merge, or build-time localization
  produces an actionable unsupported-feature error, not partial execution.
- E05: render two selected siblings beneath the same guidance ancestor. Explicit
  sequence expansion retains their contribution group and emits the ancestor once;
  equal display titles from different libraries are not implicitly merged.
- E06: resolved inline text contains literal `${...}`, backslashes, and quote
  delimiters. Parse/sequence/reparse preserves intended text without a second
  interpolation pass or formatter corruption.

## Typed Library And Guidance Inputs

| ID | Given / action | Required observable result |
|---|---|---|
| L01 | Typed alias and equivalent `node { type = "..."; }` for each proposed type | Equal canonical catalog/kind/behavior; contradictory explicit type errors |
| L02 | Include `md-file` with frontmatter, template-like text and comments | Exact bytes preserved, no interpolation or metadata stripping |
| L03 | Include `md-section` backed by a document or exact heading array | Expected section body/hierarchy; source file/bundle unchanged; ambiguity is an error |
| L04 | Physical directory has annotated children and an undeclared supporting script | Script discovered, included and hashed; no exhaustive asset registry required |
| L05 | role=skill/content associates a main file; unrelated physical role=agent file | Whole skill preserved, main file shown by default; role does not secretly rename/convert file or select a protocol |
| L06 | Virtual guidance group versus explicit title | No wrapper for virtual group; one explicit wrapper when requested; source title used unless overridden |
| L07 | Source declares vars; consumer supplies binding or explicit include values; source also contains shell `${PATH}` | Only declared names bound, override precedence explicit, shell placeholder retained literally, expected rendered bytes |
| L08 | Missing declared input, undeclared/unused explicit value, or vars on intact file/directory | Located error with node/source context; no silent ambient injection or byte rewriting |
| L09 | Sequence templated guidance, remove plan-only bindings from generated file, execute recorded build | Recorded per-include values reproduce identical text; no missing caller scope or second interpolation |
| L10 | Preset selects guidance plus defaults; consumer excludes/replaces one member | Selection/reason provenance preserved, no content duplicate; defaults and final collision rules identical to explicit equivalent |
| L11 | Preset cycle, duplicate library IDs, missing role main file, or type/source mismatch | Actionable loader/resolver diagnostic; partial source inspection survives independent errors |
| L12 | Discover new choose child without default; select unrelated exact descendant | Broad default adoption shows undecided child and blocks; exact independently resolved selection remains usable |
| L13 | New format library beside retained old HCL/YAML configs | New loader uses only declared new-format input; no implicit merging, deletion, or YAML-runtime dependency |

## Completion Gates

1. API-level cases precede thin CLI wrappers. Test actual byte/layout/ownership
   behavior and rollback, not just whether expected strings are printed.
2. A01-A12, B01-B05 (including B03a)/B09-B12, C01-C13 and core D localization cases define the
   useful vertical slice; it is not the whole promised release.
3. Full curation B06-B08 and physical metadata/preset CLI coverage require the
   reviewed new library schema/parser. Do not label them delivered based only on
   hand-built catalogs or legacy format readers.
   L01-L13 cover this explicit typed-library contract and guidance-value recording.
4. D06-D08 recipes and E01 compaction may follow the first end-to-end slice. Retain
   their contracts and report unsupported syntax rather than ignoring it initially.
5. Run repository `tools/check` once implementation lands, focused tests per slice,
   and the exact pilot CLI journey. Protocol/CI wrappers and interactive UI are
   separately tracked follow-ups, not reasons to skip core inspection diagnostics.
