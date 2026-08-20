# Agent Library Review Plan

Status: active planning and content-audit plan.

## Goal

Turn the owner's review in `docs/notes_on_lib_mods.md` into a coherent,
well-explained library without losing useful detail from the original agent
guides. Keep three concerns separate:

1. what an agent should be instructed to do;
2. how reusable modules should be organized and related;
3. what Mogent should validate or explain while composing them.

The owner notes are a working review artifact. Preserve them as user-owned work
and do not stage or rewrite them unless explicitly requested.

## Current Assessment

The library contains useful material, but its first extraction pass left visible
provenance residue and compressed some rules too aggressively:

- `Corpus Variants` describes how source repositories differed, not an
  instruction that belongs in a rendered `AGENTS.md`. Replace it with concrete
  language, domain, or repository-local modules, or omit it after its useful
  distinction is captured.
- `shared-baseline`, `engineering`, `go`, `personal`, and `ucd_research` contain
  real overlap. Some duplication is intentional layering; other duplication
  currently gives a user no clear reason to select one module instead of
  another.
- Several compact directives need one concrete example or a clearer name:
  offline tests, main success paths, structured assertions, snapshots/golden
  files, imperative commit subjects, shell failure handling, Go `%w`, stable
  coordination handles, and decision-record lifecycle.
- `requires` and `conflicts_with` are parsed and displayed, but Mogent does not
  enforce them. Existing values also assume source aliases such as `shared` or
  `ucd`, even though aliases are chosen by each consuming manifest.
- Personas are alternatives, while related process or language modules are
  often complementary. Those are different relationships and should not be
  represented by source-tree order or a silent "first one wins" rule.
- Some personal and research modules contain broadly reusable engineering
  guidance. Moving that guidance requires a deliberate ownership boundary;
  project-specific safety and procedure must remain in their narrower sources.

## Issue Tracks

### 1. Clarity And Examples

These changes can be made as focused prose slices once the source comparison
has captured the current baseline:

- remove every rendered `Corpus Variants` section after relocating any useful
  distinction;
- rename vague headings such as Go `Code`, `Shared Baseline`, `Minimal Diff`,
  and `Lightweight Escalation` where a more literal name improves selection;
- add short examples only where a term changes behavior;
- explain TE, DR, DI, minting, stable namespaces, supersession, and archival
  status without requiring the reader to know this repository's history;
- add an explicit instruction to raise a user request that conflicts with the
  active agent guide instead of silently choosing one;
- add "show expected usage or before/after behavior when applicable" to the
  focused change and handoff guidance;
- expand documentation guidance around audience, learning goals,
  agent-reference docs, human-facing docs, and inconsistent local style.

Examples should stay small. A reusable instruction module is not a tutorial;
longer explanations belong in authoring documentation or `source show` help.

### 2. Taxonomy And Deduplication

Review these boundaries as groups rather than file-by-file rewrites:

- **Shared baseline versus engineering:** keep the baseline as the smallest
  safe workflow; keep engineering modules as selectable depth. Remove repeated
  sentences that add no stronger rule or example.
- **Go:** split generic error and test policy from genuinely Go-specific
  tooling, package/API shape, table tests, formatting, and error wrapping.
- **Decision-first:** retain a selectable parent summary plus its atomic child
  modules if subtree selection remains useful. Decide whether the directory
  name is the concept and whether the parent summary earns a separate file.
- **Communication:** make persona choice an explicit alternative family. Break
  Teacherbot into reusable teaching/learning/chat methods only if those methods
  are useful without the Teacherbot persona.
- **Personal engineering:** consider promoting the API -> CLI -> GUI/local API
  shared-operation rule and other broadly useful material into `cdint`, while
  retaining migration- or project-specific overlays.
- **Research:** keep participant-data, experiment-validity, and lab-procedure
  rules in research sources. Move general Python dependency guidance to a
  language source only when it is not tied to UCD procedure.

Use `git mv` for accepted moves. Record a mapping from every old source path to
its new path or deliberate retirement so consuming manifests can be updated.

### 3. Module Relationships And Composition Diagnostics

Before changing the public metadata format, run a decision-first design slice
for these relationships:

- `requires`: a hard prerequisite without which a module is incomplete or
  unsafe;
- `see_also`: a soft discovery link that does not cause selection;
- `conflicts_with`: a known exact incompatibility;
- `alternative_family` plus a variant name: mutually exclusive choices such as
  personas;
- provenance/reference: a pointer to rationale or a Decision Intent, not a
  selection dependency.

That decision must resolve how references survive user-chosen source aliases.
Candidate approaches include source-local paths for same-library relationships,
stable library identities distinct from manifest aliases, or delaying
machine-enforced cross-source references. Do not make bare paths globally
meaningful and do not let manifest order silently resolve conflicts.

Once the format is decided, Mogent should:

1. show relevant relationships during `source show` and selection preview;
2. reject missing hard requirements or offer the exact addition that fixes
   them;
3. warn on exact conflicts and multiple selected alternatives before writing;
4. say what conflicts, why it matters, and the safest next command;
5. keep soft links informational;
6. remain explicit when it cannot infer a semantic conflict.

The rendered manager guidance should separately tell an agent to raise a
conflict between the user's request and active repository instructions. Mogent
can validate declared composition relationships; it cannot generally determine
whether arbitrary prose or a new chat request is semantically incompatible.

### 4. Source-Guide Comparison

The original guides under `docs/other_repo_agents/` are protected evidence, not
normal library input. The user has authorized a later comparison against the
existing library. Perform that comparison before broad content rewrites and do
not edit the originals.

For each original guide, record locally:

| Field | Question |
|---|---|
| Source location | Which original guide and heading supplied the idea? |
| Current module | Where is the idea represented now? |
| Coverage | Preserved, compressed, missing, or deliberately omitted? |
| Fit | Reusable, language-specific, domain-specific, or repo-local? |
| Relationship | Duplicate, prerequisite, soft reference, alternative, or conflict? |
| Conflict | Does it disagree with the design of record or another active rule? |
| Action | Keep, clarify, split, move, merge, retire, or request a decision? |

Keep the first detailed matrix local and untracked because filenames, paths,
and prose may reveal private repository context. Promote only reviewed findings
that are safe and useful, and describe reusable ideas rather than quoting large
source passages. The design of record wins when historical source material
conflicts with current Mogent decisions.

## Implementation Sequence

1. Snapshot the current source tree, metadata, and source paths as the audit
   baseline.
2. Compare the protected original guides using the local matrix above.
3. Classify findings into clarity fixes, taxonomy changes, relationship-format
   decisions, Mogent behavior, and intentionally repo-local material.
4. Make a focused clarity pass: remove extraction residue, explain unfamiliar
   terms, and add small examples without moving paths.
5. Run a thought experiment for taxonomy and relationship metadata; ask the
   user to resolve durable alternatives and record the resulting Decision
   Intent before changing source paths or public metadata.
6. Reorganize one bounded module family at a time, with an old-to-new path map
   and manifest compatibility review.
7. Implement relationship diagnostics through the public workspace API, then
   expose the same behavior in the CLI.
8. Dogfood the revised modules and diagnostics in active repositories before
   declaring the review complete.

## Validation And Evidence

For prose-only slices:

- inspect the exact diff and source-tree output;
- validate frontmatter and source loading;
- render a fixture and confirm tool-only metadata and explanations do not leak
  unexpectedly into output;
- show a representative before/after selection or rendered section.

For metadata or behavior changes:

- add success, missing-requirement, conflict, alternative-family, cross-source,
  and alias-renaming coverage;
- exercise the public `workspace` API and the equivalent CLI path;
- run focused Go tests, `gofmt`, `go test ./...`, `go vet ./...`, and
  `errcheck ./...` with caches under `/tmp`;
- verify existing manifests either remain valid or receive an explicit,
  reviewed migration.

## Decisions Still Needed

1. Whether relationship references use source-local paths, stable library
   identities, or another alias-independent form.
2. Whether `see_also`, `alternative_family`, and provenance become first-class
   metadata fields in the next slice.
3. Whether missing hard requirements are errors immediately or warnings until
   authoring commands can offer an atomic fix.
4. The target top-level taxonomy and the compatibility policy for moved source
   paths.
5. Which broadly reusable personal/research rules the owners want to promote
   into `cdint`.
6. Whether detailed source-comparison notes stay permanently local or receive a
   sanitized, committed summary after owner review.
