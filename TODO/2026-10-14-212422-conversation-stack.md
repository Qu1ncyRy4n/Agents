# Configuration Redesign Conversation Stack

Captured: 2026-10-14 21:24:22 PDT
Status: active owner discussion; not an executable specification
Resume: syntax API implemented/tested; S2 typed catalog/resolver API next; source-editing follow-up remains in S1

Implementation checkpoint: new `syntax` package supplies Parse, ordered/source-
located declarations and comments, exact Original input retention, explicit string
evaluation, and deterministic Format with optional alias expansion. External-package
tests exercise library/plan/build examples, separator rules, arrays, invalid syntax,
multiline/literal/interpolation semantics, Unicode ranges, depth bounds and format
round trips. No new dependencies. `nix develop -c tools/check` passes. New CLI/
catalog/runtime workflow remains unimplemented; do not confuse parsed options with
resolved behavior. API implementation preceded tests, per owner process. No commits
made. Resume S2 API then tests; CLI expected cases already documented in acceptance.

Owner accepted the combined final review and says good to go. Library filename is
`mogent.lib`. Clarification: no semicolon required for a single final statement
before a closing brace; semicolons separate multiple statements on one line;
commas remain array separators, not alternate statement separators. Development
order requested: API implementation -> API tests -> CLI -> CLI behavior tests;
define expected CLI behavior beforehand and optionally encode those tests before
CLI wiring. Begin the bounded syntax API slice; keep the remaining staged tasks
visible and do not claim complete runtime/CLI support prematurely. No commits
are authorized by this approval alone.

Owner requests all remaining pre-implementation items in one pass. The consolidated
spec now proposes `library.mogent`, typed nodes/aliases (dir/file/md-file/md-section),
optional role/content, curation/presets, declared Markdown vars and per-include
resolved values. Acceptance L01-L13 cover these directly. Existing header/source/
localization/transaction proposals remain in the same spec. Present one finite
review covering naming/schema, source-value records, protection, and staged tests;
do not reopen format surveys or start implementing before this review. These
concrete spellings are still draft choices, not owner approval inferred from the
request to see them together.

Latest owner response: suggests `md-section` rather than `md-heading`, leans toward
`md-file` meaning intact Markdown, says the direction looks good, and requests the
stack/readiness checkpoint. Use these as working names, not a new format search:
`dir`/`file` physical payloads; `md-file` intact Markdown; `md-section` logical
guidance or an explicit document subtree. Role and main-file association remain
separate annotations; explicit vars apply to generated guidance, not intact copies.
Before code depending on library syntax, reconcile this typed proposal with spec
sections 3-5 and the acceptance fixture. Do not claim the existing untyped grammar
already specifies it. Parser/catalog work can begin once that reconciliation is
reviewed; runtime implementation has not started. Mechanical source-value recording,
alias normalization, and integrity serialization still need concrete contracts, not
another broad product questionnaire.

Latest owner feedback on library example: consider explicit node types such as
md-file, md-heading, dir, and front aliases normalizing to node declarations;
distinguish intended agent/skill target type; consider listing variables injected
into Markdown; main skill content and its target filename may need association.
These are exploratory, not settled fields. Agent recommendation to discuss:
separate structural type (file/heading/dir), optional semantic role (skill/
instructions/profile), source/main-file locator, and consumer output placement.
`md-file`/`md-heading` can be convenience aliases if their canonical structure and
format/view meaning are explicit. A declared vars list can whitelist required
consumer-provided string inputs for authored guidance rendering only; do not
evaluate templates in intact copied skill/config payloads by default. Main-file
association need not enumerate assets. Renaming SKILL.md to skill.md is an explicit
projection/transform question, not an automatic result of content inclusion.
Keep the next response concrete and small; do not infer a final library grammar
or implement node/target/variable schema without owner review.

New library example proposed in walkthrough section 8: library header; logical
guidance `node` with Markdown `file`; physical `node` with `directory`, annotated
children/tags and child-local defaults; `preset` reusing include with proposed
`self:` alias. No special `skill` keyword or entrypoint field. Library filename
and these field spellings remain proposals. Owner has not accepted this example
yet. Physical discovery must not become exhaustive per-file registration.

Latest owner decision: choose B, designing the new library format in the initial
redesign. Keep existing config files as material to interpret for migration/reference;
the owner has no desire to maintain YAML. This supersedes the draft's proposed
temporary runtime HCL-library bridge. Do not delete existing files or imply a
permanent dual-format parser. Library grammar/filename and metadata field names
still need review; no special skill keyword has been accepted. Spec, acceptance
fixture, and R2/S1/S2 tasks now reflect this decision.

Latest owner answer: accepts one-time localization and asks whether enough is
available to write the spec and acceptance cases. Consolidated drafts now exist:
[spec](../docs/proposals/MOGENT-FIRST-SLICE-SPEC.md),
[acceptance cases](../docs/proposals/MOGENT-FIRST-SLICE-ACCEPTANCE.md), and
[implementation sequence](TODO-config-first-slice.md). Review R1-R3, then start the
first approved slice. Proposed grammar/serialization/adapter choices are explicitly
marked, not silently promoted to owner decisions. Cases are not executed tests.
No runtime implementation or commits have been made. Older paragraphs below are
conversation history and may describe proposals superseded by the consolidated draft.

Latest owner decision: skip the `copy` alias. Initial layout is a real root
`AGENTS.md` and other managed artifacts under `.agents/`. Symlink layout is not
part of this first slice. Owner requests a deferred TODO for per-file protocol
defaults such as `protocol = "dotagents"` or `"openai"`; these are illustrative
names, not an approved schema or compatibility claim. Record target versioning,
destination overrides, recipe scope, and generated-build propagation when this
is revisited. Continue the core walkthrough without designing protocols now.

Owner accepted the build-header/resolved-source example: generated origin and
generated_from/generated_by fields; authored origin protects direct builds;
optional human/agent author; detect edits before regenerating and be conservative
without a baseline. Git ref in plan resolves to commit in build; complete local
nodes use content hashes. Exact schema/local-hash serialization is not finalized.
Next owner asks for localization syntax and remaining work before implementation.
See walkthrough sections 6/7. Proposed `edit include ... --localize --replace`
creates editable local source + provenance, then normal local include; it is not
a repeated build-time copy. No parser/CLI implementation has started.

Source identity discussion: owner asks what "fingerprint" and "frozen build"
mean and is willing to require Git tracking or per-node content hashes. Explain
fingerprint as an ordinary content checksum, not new identity or storage machinery.
"Frozen build" was agent terminology, not an accepted extra mode; prefer one
resolved-build contract. Agent recommendation pending review: remote inputs pinned
by commit; local inputs either reference a reconstructible committed Git revision
or carry selected-node integrity hashes, with source bytes retained in Git or a
snapshot if later reconstruction is promised. Tracked alone does not cover dirty
or untracked bundle members; hash directory members/paths/permission information
as relevant, not only the main Markdown file. Hashes verify bytes but cannot
retrieve missing content. Source identity remains stable logical path, not hash.
Exact per-node versus source-level serialization and manual-build authoring
workflow need specification; do not invent a new store or silently ignore assets.

Owner accepted the source-verification direction: Git commits where possible,
selected-local-node content hashes including complete bundles, visible dirty inputs,
and retained/committed bytes before claiming fresh-checkout reproducibility.

Owner follow-up: consider a skill, harness integration, or GitHub Action checking
whether a Mogent setup is up to date. Agent recommendation: one non-mutating
`mogent build --check` core used by CI and optional skill/harness wrappers. Compare
current intent's in-memory resolved sequence to the generated build (when plan-
derived), verify recorded input revisions/hashes, validate choices/requirements,
and compare expected managed outputs to disk (bytes/modes/layout as relevant).
Report stale generated build, output drift, and unresolved decisions with nonzero
status. Do not rewrite files, refresh moving refs, or silently repair during check;
fetching an exact missing pin into cache is a separate documented side effect.
Report upstream update availability separately from "reproduces pinned inputs."
Public command contract and CI integration remain proposed, not implemented.

Owner steering: keep CI/skill/harness integration in TODOs and return to the core
task. Git/hash direction is accepted. Next artifact should consolidate the small
end-to-end fixture into a first-slice spec: source declarations and integrity,
plan -> resolved build -> result, canonical include/alias normalization, authored
build protection, and unresolved-choice/output inspection. Separate accepted
decisions from exact grammar/serialization details still needing approval. Do not
restart format research or expand integrations. Follow with a bounded vertical
implementation slice once the executable contract is reviewed.

Latest discussion: owner wants optional `replace` (and possibly other) aliases
that normalize to canonical include/options; explicit expanded build output must
avoid aliases when requested. Alias expansion is independent of compaction and
logging verbosity. Owner also questions `.agents` as the universal output base
and proposes central payloads plus root symlinks. See walkthrough section 5 for
primary-doc evidence, proposed role-specific defaults, and unresolved output-link
semantics. No symlink support or alias parser has been implemented.

## How To Use This Stack

The owner spent about an hour responding to the preceding design discussion.
Reduce the conversational load: discuss one concrete issue, record the answer,
then advance. Put tangents in the parking lot instead of answering every open
question at once. Do not present another large questionnaire or treat exploratory
examples as implementation approval.

`TODO/TODO.md` indexes this work. This file preserves the conversation, questions,
and detailed resume point. Existing specifications still describe the executable
formats. Record durable grammar decisions separately when the owner is ready.

## Owner Direction Already Given

- Improve authoring and inspection while reusing the existing composition engine.
- Preserve metadata, tags, curation, recommendations, and content composition.
- Pursue a hybrid: curated libraries + consumer intent -> resolved ordered actions
  -> materialized result. The exact file roles and grammar remain open.
- Accepted entry-point relationship: a build can be authored and used directly,
  without a plan. If a plan is used, it MUST compile through the same build
  representation and executor; it cannot materialize outputs by a separate path.
- Workspace targets first; automatic home/global installation is deferred.
- Workspace DotAgents configuration is the proposed 1.0 target. Before freezing
  release scope, compare actual harness support with AGENTS.md, Agent Skills,
  agentsstandard.com, and other relevant formats. Do not conflate these projects.
- Provide a small, safe recommended vocabulary, with advanced variations available
  when needed. Broad includes plus exceptions should be easy to read.
- An active/interactive result-tree viewer is central, recovering the spirit of
  the original TUI. Show available choices, resolved reasons/actions, and results.
- Comprehensibility, reproducibility, clarity, and expressiveness with simplicity
  are the evaluation criteria. Test whether the owner can predict a small example.
- No random choice may affect selection, defaults, action order, or output. If a
  future feature genuinely needs randomness, its seed must be explicit,
  replaceable, and recorded in state. No such feature is currently required.
- Prefer Git and established tooling over a new Nix-like integrity/store system.
- Document feature coverage/status and workflow guidance; use existing repositories
  and useful real modules for examples once their selection is agreed.

## Working Metaphor

| Workshop concept | Mogent concept |
|---|---|
| Curated supply library | Content, metadata, tags, recommendations |
| Work order | Consumer's authored intent |
| Assembly sheet | Concrete actions in execution order |
| Finished assembly | Workspace files and generated documents |
| Inspection record | Provenance, choices, omissions, replacements, integrity |

The library explains what is available; the work order records what is wanted;
the assembly sheet explains exactly what will happen. This is a thinking aid,
not a requirement to rename the CLI with workshop terminology.

## Focus Stack

### 1. Hand-Authored Build As A First-Class Entry Point — ACCEPTED RELATIONSHIP

Owner proposal: a `mogent.build` assembly sequence can be written directly, or
generated from a higher-level `mogent.plan`. File extensions are undecided.
The supervisor likely prefers direct sequence authoring.

Owner answer, 2026-10-14: independently authored builds are a complete entry
point. Plans MUST compile through builds. Only using a plan is optional, not
the plan-to-build compilation stage. The owner also accepted protecting manual
edits to generated builds from silent regeneration.

```text
hand-authored build ---------------------> shared executor -> result
authored plan -> generated build --------> shared executor -> result
```

Owner suggestion: record generated/derived versus authored status and human/agent
authorship in a file header or `state.json`. Exact storage/schema remains open.
Keep derivation and author identity separate: an agent can author a direct build,
and a human can edit a generated build. Do not use identity to determine whether
regeneration may overwrite edits. State is currently ignored/machine-local, so
cross-checkout derivation discovery needs consideration if state is the only record.

Follow-up owner answer: separate `author` and `origin` concepts are accepted;
a build header is accepted as a provenance location. Exact field spelling/schema
is still open. The owner questions ignoring state if it harms reproducibility.
Explain the boundary: reproducibility inputs/pins/derivation must travel with the
build; machine-local installation/drift state has a different job. Do not silently
change `.gitignore` or the current state contract during this discussion.

Remaining details to revisit after the next concrete examples:

- Which metadata/curation checks apply to explicit build instructions? Safety,
  identity, and bundle validation should not disappear; the accepted curation
  boundary still needs a decision.
- If a build was generated, how is that provenance shown and how are manual edits
  protected from regeneration? Never silently overwrite them or let two files
  compete as authoritative input.
- Can one small example be expressed both ways with the same final manifest?

The shared executor and direct-build entry point are accepted direction, not
implemented functionality. Input binding, provenance schema, regeneration
mechanics, and format grammar remain undecided.

### 2. Recommended Includes, Exceptions, And Replacements — REPRESENTATIONS ACCEPTED

Owner's preferred happy paths:

1. Include a team's full tree with a few exceptions.
2. Include personal content selected by tags X/Y/Z.
3. Advanced: include a recommended collection and replace node X with node Y from
   another library.

Current discussion: the owner proposes alternative build representations:

- Verbose: explicit lowest complete nodes/bundles, possibly every physical node.
- Exception-based: whole subtrees with explicit exclusions.
- Smart: a short combination of explicit nodes and subtree exclusions.

Question to answer now: can these yield a deterministic build sequence, and does
minimizing its size introduce an NP-hard optimization problem? Distinguish
execution semantics from serialization styles. All styles should resolve to the
same selected members, operation order, placement, bundle boundaries, and reasons
at pinned inputs. Explicit per-file display need not permit partial skill installs.
Tree-only compression admits bounded recursive algorithms; arbitrary overlapping
tag/preset selection can make global shortest-expression optimization resemble
set cover. Prefer a deterministic heuristic with validated equivalence, not a
promise of globally minimum lines. This is agent analysis pending discussion.

Follow-up owner answer: explicit and subtree-with-exceptions representations of
the same build make sense; use "compact" for the shorter mixed representation.
These are presentation/serialization choices, not different executor semantics.
The owner suggests comparing inclusion/exclusion counts at each node and taking
the smaller description. Evaluate recursively with a defined instruction cost
(including the parent include), not just selected-node counts; preserve ordering,
operations, destinations, bundle boundaries, and provenance. Prefer the explicit
form on equal cost as an initial recommendation, not a finalized rule.

Potential formatting switches suggested by the owner: `--optmise` (as written),
`--unoptimized`, or `--raw-seq`. Exact names and CLI placement are deferred;
none is implemented. A raw/explicit form must not imply bypassing validation.
No globally shortest tag/preset expression is promised. Reverse inference remains
in the parking lot. Public replacement vocabulary still needs the later examples.

Avoid making users learn six verbs before they can do these jobs. Compare a
small `include` family with explicit copy/render options and a `replace` family.
Do not assume the illustrative CLI flags already exist.

Resolve the confusing distinction with examples, not terminology:

- A guidance slot gets different text while retaining its local heading/place.
- A complete installed skill is replaced; obsolete scripts must disappear.
- A local override is an input with provenance, not an edit to generated output.

Decide whether one public replacement operation can infer/validate these cases
from addressed node kinds. They need distinct internal behavior, but do not
necessarily need distinct public verbs. A localization option may fit better.

### 3. Resolver Traversal And Decision Feedback — CORE RULES ACCEPTED

Start with one precedence example: a node matches a tag query but is explicitly
excluded in the same include. Agent recommendation: the explicit exclusion wins
within that include; it must not silently delete another include's contribution.
Discuss later explicit re-inclusion separately and validate final curation and
requirements. Settle this small example before the full severity/default table.

Owner answer: accepted. Explicit exclusions beat tag matches within the same
include. They filter only that include and do not remove an earlier include's
contribution. This does not exempt the final result from curation/requirement
validation. Later explicit re-inclusion and deliberate result removal remain open.

Next single discussion: when a selected curated branch needs a decision, what
should preview show and what blocks execution? Start with a branch having complete
recommended defaults versus one having an unanswered required choice. Preserve
accepted-default reasons and an inspectable unresolved tree.

Owner answer: accepted. Explicitly accepting recommended defaults resolves those
choices without repeated confirmation; preview labels their default-derived
origin. Unanswered required choices show UNDECIDED and block execution. Omitted
opt-in items remain excluded without error.

Owner follow-up: reason-required mandatory exclusions and opt-in/experimental
inclusions may be useful. Discuss this as a proposed extension, not already
implemented parity: current v2 foundation exclusions require a reason; opt-in
inclusions warn but do not require one. "Experimental" is not a current separate
curation mode. Candidate rule: the consumer records a nonempty reason for either
exception, the viewer displays it, and missing reasons block execution; no repeated
interactive approval once the reason is recorded. Decide whether reason-required
opt-in inclusion applies universally or only to explicitly marked branches before
finalizing grammar. A recorded exception is not proof of semantic equivalence or
an exemption from unrelated hard requirements.

Final owner answer: accept the recommended distinction. Excluding foundational
content requires a recorded reason. Ordinary opt-in inclusion requires explicit
selection but no reason. Experimental/reason-required inclusion requires a recorded
reason. Missing required reasons block execution; recorded reasons are visible in
the authored configuration and build explanation, with no repeated approval prompt.
Exact experimental metadata/schema remains open. These are redesign decisions,
not implemented changes to the existing v2 parser or planner.

Owner's proposed algorithm, preserved for evaluation:

1. Enter a source root and traverse its tree.
2. Include children according to defaults/curation; flag choices needing attention.
3. Collect exceptions (excluded nodes and alternative inclusions).
4. Traverse those nodes and recurse until selections are accounted for.
5. Order by section/source and declared ordering.
6. Produce a readable concrete list.
7. Optionally notice equivalent tag/include patterns and recommend a simpler form.

Separate selection, validation, ordering, and optional simplification. Check
requirements against the final selected set; a requirement may be encountered
later in traversal. Exclusions must not silently reappear through tags/defaults.
Replacement can invalidate a previously satisfied requirement. Define these
interactions before adopting a traversal implementation.

Owner wants defaults shown as accepted/flagged decisions, unresolved choices
shown near their nodes (`UNDEFINED` or similar), and errors when action is needed.
Exact severity rules are open. Preview should retain a partial diagnostic tree;
an invalid selection must not become an applicable build. Prefer collecting
independent diagnostics rather than losing the tree at the first error.

Clarify `opt_in`: omission normally means excluded, not necessarily an error.
An undecided `choose` branch differs from an intentionally omitted opt-in item.
Hard requirements versus reasoned exceptions also need explicit semantics.

### 4. Ordering In One Thread — ACCEPTED DIRECTION

Owner accepted: consumer include/action
order first; each expanded collection retains library-sidecar order; otherwise
canonical source-relative path lexical order, independent of locale. Tag queries
select membership, not a new ordering. No length tie-breaker or public dependency
ordering language is needed initially. Recipe expansion is inline, not a global
sort. Explicit replacements remain required; order is not overwrite permission.
Tag-selected members retain library-sidecar order, rather than query order or a
separate tag ranking. Exact lexical comparison details remain a specification task.

- Owner prefers library-sidecar order, not a mandatory A-before-B relationship on
  every item. Relative relationships and see-also links may have separate uses.
- In a direct build, authored action order means the sequence in that build file.
- For collection members, library order governs unless the consumer deliberately
  requests reordering. An included recipe's actions participate at its position.
- Accepted precedence: consumer composition order, then library sibling order;
  do not globally sort by source after resolving replacements or appends.
- `requires` means presence, not automatically execution-before. `see_also` is
  navigation, not selection or order.
- Topological sorting is only useful for genuine prerequisite edges. It cannot
  select winners, resolve curation, or infer tree replacement versus directory union.
- Use deterministic lexical path tie-breaking where no order is declared. Earlier
  discussion considered length; the accepted direction does not need it. Define
  case/Unicode/path handling and avoid arbitrary tie-breakers.
- Explain "expand": replace a collection/recipe reference with its concrete
  selected members/actions, retaining its origin and reasons.

### 5. Library Contract, Curation, Tags, And Presets — ACCEPTED DIRECTION

Owner accepted: retain the library's
hierarchy and recommendations without forcing its hierarchy to be the consumer's
only output outline. A named preset is a reusable selection with explicit default
adoption/choices, not a second copy of content. Resolve it with consumer exceptions
and preserve both sources of decisions. A tag query is computed membership; a
preset records a maintained choice. Consumer destinations remain explicit; any
preset placement must be visible in an accepted recipe, not hidden in tags.

Owner wording: "tags should probably just default to sidecar library format."
Tag-selected ordering follows the sidecar as above. Do not infer a new metadata
migration/override rule from that phrase; clarify only if tag storage/format comes
up, since current v1 and v2 have different metadata ownership conventions.

Preserve library identity, logical grouping, titles/order, tags/descriptions,
provenance, bundles, recommendations, relationships, and compatibility.
Consumer acceptance, exceptions, and placement remain visible choices.

Answer the owner's concrete questions:

- Does the existing library have responsibilities it should not have?
- What changes from the v1 consumer-owned outline and v2 library-owned hierarchy?
- How can broad includes and tags compose without bypassing curated decisions?
- Are named directly includable collections already present?

Current boundary: selecting a branch already provides a named subtree. V2
`qmr-core` is an init template, not a live includable preset. A cross-branch
curated selection/recipe is not the same thing as either.

Illustrative preset (not grammar): `team:recommended-go` selects foundational
guidance, curated workflow defaults, Go review/testing skills, and a documented
choice of commit behavior. A tag query such as `lang/go` only computes matching
membership; it does not inherently encode those decisions or placement.

Explain the earlier authority qualification simply: replacing a security rule
with "skip security review" is an explicit replacement, but explicitness alone
does not prove that the organization rule was preserved. Do not introduce an
automated semantic-proof promise or silently remove curation authority.

### 6. Minimal Syntax And File Roles — NEXT

**Current owner decision:** pursue a small custom Mogent language. This supersedes
the format-search preferences below, retained as discussion history. Root-level
ordered actions, nested blocks, explicit assignments/booleans, and bracketed
comma-separated lists are the intended surface. Support for optional single-line
actions is desired. No `steps = [...]` wrapper is required. File roles establish
plan/build behavior; a mode switch is not currently selected.

Next: validate one complete tiny build (source, binding, copy/include, exception,
inline rendered text) and its exact resulting tree before specifying/implementing
the parser. Immutable string bindings/interpolation and triple-quoted multiline
strings are requested features to evaluate, with exact grammar, indentation,
escaping, scope, and declaration-order rules still open. Explicit values for
repo names/dates/developer names/paths should not depend on ambient discovery or
the clock during execution. Parser reuse remains an implementation option, not
permission to accept a larger language or begin implementation now.

Owner accepted the tiny build example and proposed baseline: `let` bindings,
`${name}` interpolation, local `source` declarations, `include` with explicit
mode/destination/exclusions, and a `document` block with triple-quoted inline text.
The presented rules were immutable declare-before-use bindings, duplicate/unknown
name errors, multiline framing-newline/common-indent removal, top-to-bottom action
order, order-insensitive option fields, build-relative local sources, and
workspace-relative destinations. Single-line blocks can use semicolon-separated
fields. This approves design direction/example, not a finished grammar or parser
implementation. Resolve exact escapes, trailing newline behavior, scope across
included recipes, and local-input freezing before implementation.

Current question: immediately useful additions to this baseline. Prioritize
comments and versioning, pinned Git sources using existing machinery, reusable
build inclusion, explicit replacement, and ordered document contributions using
existing render/append capabilities. Keep language features distinct from action
capabilities; avoid adding general functions/conditions/loops without a story.

Owner accepted the proposed additions as direction: comments/versioning, pinned
Git sources, reusable build inclusion, explicit replacement, mixed inline/source
document contributions, and a literal-text escape capability. Exact spellings
such as `include-build`, document child syntax, and triple-single quotes remain
candidate grammar, not implemented features. Retain tags/curation/presets.

Owner accepted private bindings and source aliases for included builds: importing
a reusable recipe must not overwrite caller bindings/aliases. Explicit parameter
passing, interpolation scope, and source lookup rules still need specification.
Owner suggests undefined/unused alias diagnostics. Recommendation to discuss:
undefined binding or alias is an error; declared but unused binding/alias is a
warning, not a blocker (sources may be declared for discovery). Check usage in
the declaring scope, including interpolation references; avoid treating accepted
default choices as unused merely because they originated in a library.

Owner accepted diagnostic severities: undefined bindings/aliases are errors;
unused declarations are warnings. This is redesign direction, not existing
parser behavior for the hypothetical language.

Next stack item is 7: inspection/workflow. Start with the desired source ->
selection/decision -> result viewing loop and a small CLI navigation sketch.
Existing commands already use `plan` and `build` for different v1/v2 jobs; proposed
file-role names do not automatically authorize repurposing those CLI commands.

#### Format Exploration History

Owner prefers evaluating KDL as a potentially simpler fit than HCL. One message
said "tdk"; context suggests KDL, but confirm before treating that as a distinct
format request. Compare equivalent small examples before selecting a parser.

Owner response to equivalent KDL/HCL examples: they read equally well; the owner
slightly prefers HCL's bracketed comma-separated lists, but wants distance from
the "HCL flavor." This is not approval to migrate to KDL or to finalize HCL.
Separate surface syntax/parser reuse from Mogent semantics and product vocabulary.
Next recommendation to discuss: familiar braces/attributes/lists with a small
instruction-oriented Mogent grammar, potentially parsed as literal-only HCL;
no inherited Terraform model, expression language, or old v2 block vocabulary.
Do not invent a custom parser merely for different punctuation. Exact syntax,
supported literal forms, filenames/extensions, and tooling remain undecided.

Latest owner steering: not convinced by retaining reduced HCL; continue examining
KDL with fuller side-by-side examples. Concern is semantic expectations and a
smaller language, not merely cosmetic difference. No format choice is accepted.
Primary references checked: https://kdl.dev/ and https://kdl.dev/spec/ (KDL 2.0),
https://hcl.readthedocs.io/en/latest/language_design.html. Both leave application
semantics to the host; HCL need not expose variables/functions, but includes an
expression language. KDL has literal arguments/properties/children without an
expression evaluator or native array value. Sequence belongs in child/sibling
nodes, not property order. KDL v2 booleans are `#true`/`#false`; duplicate property
keys resolve rightmost by spec, unlike HCL duplicate attributes. Go parser/version
and round-trip tooling suitability still need verification. Examples are grammar
exploration, not implemented configuration or final schema.

Latest syntax feedback: owner dislikes KDL positional multi-value lists, implicit
marker booleans, and crowded inline options. Prefers document-like nesting,
explicit `key = value`, `true`/`false`, bracket/comma arrays, and meaningful action
order. Requests a better-fitting alternative soon; leaning toward a tiny custom
HCL-like language rather than forcing KDL or inheriting HCL2 expectations.
"Named properties" means literal options on a KDL node, not mutable variables.
Standalone `mode="copy"` inside a KDL children block is not valid KDL syntax.

Bounded alternatives checked: HOCON (Lightbend spec) brings substitution/merge
semantics; UCL (vstakhov/libucl README) has named blocks/arrays but implicit arrays,
macros, and variable/include machinery; JSONC can represent ordered actions in an
explicit array but lacks desired labeled blocks. `alecthomas/hcl` is a closer
Go implementation candidate: README describes HCL1 compatibility, ordered/commented
AST, and no HCL2-specific features; inspected parser has literal lists/maps/blocks
without general expression evaluation. This is not a vetted dependency or a
format decision. If custom syntax is pursued, first specify a small ordered
block/literal grammar and evaluate parser reuse, diagnostics, comment preservation,
and formatter costs; do not implement during this conversation without agreement.
Sources: https://github.com/alecthomas/hcl and its parser.go;
https://raw.githubusercontent.com/lightbend/config/main/HOCON.md;
https://github.com/vstakhov/libucl.

- HCL, if retained, should expose a deliberately reduced documented subset.
- Working names: `mogent.plan` = higher-level wanted result;
  `mogent.build` = concrete sequence. `.kdl`, `.hcl`, `.conf`, or JSON remain open.
- Explain whether the generated build also holds pinning information or whether
  a distinct lock has a necessary separate responsibility.
- A frozen build needs exact inputs, choices, operations, destinations, and order;
  a sequence that reruns floating tag queries is not frozen.
- Templates create editable starting points; includes compose reusable input.
  Owner likes that distinction. Final CLI/config naming remains undecided.

### 7. Inspection And Workflow

Desired loop:

```text
initialize -> declare source -> inspect choices -> include/modify
           -> inspect result and unresolved nodes -> repeat -> review/apply
```

Provide CLI and editor examples for the same operation. Source declaration is
distinct from selecting it; recipe inclusion is another job currently sharing
the word "import" in discussion.

Potential help affordances: `--help`, `--help-workflow`, `--help-config-guide`,
`--help-cli`, or structured help topics. None of these new help flags is approved
or implemented. Prefer one discoverable navigation scheme after the workflow is
tested. Do not design a new flag per documentation page prematurely.

Existing TUI is v1-specific; it is not already a complete interactive typed-result
viewer. Reuse its navigation/preview ideas and common APIs. Show source coverage,
unresolved decisions, concrete actions, final files, and rendered heading views
without confusing document headings with filesystem directories.

Owner explicitly requests clear undefined/underdefined callouts/flags in the
result tree so missing declarations and unresolved sequence/composition decisions
are visually discoverable. Preserve a partial inspection tree on resolution
failure; do not present placeholders as executable output. Accepted deterministic
ordering should not itself require extra sequencing declarations.

Owner requested inspection of the actual command set before proceeding. Findings
and executed read-only CLI examples are in
[`docs/CLI-COMMAND-INVENTORY.md`](../docs/CLI-COMMAND-INVENTORY.md). Key gaps:
v1 owns coverage/editing/TUI/status/localize; v2 owns plan/apply/update and typed
source browsing/provenance. No existing top-level show/tree command. V2 blocks
rendering on selection errors, so node-attached unresolved-result inspection is
not just a new formatter. Preserve existing command meanings until a deliberate
migration/compile-command contract is agreed.

Latest owner direction: obsolete commands and configuration formats may be removed.
Do not let historical command names or format compatibility constrain the new CLI.
This supersedes compatibility-first recommendations above for the redesign; it is
not an instruction to delete code or user configurations before agreeing the new
contract. Owner requests a restricted category/action/flag vocabulary matching the
new config language. Keep recommendations distinct from implemented commands.

Current agent proposal, pending owner review:

- `init`: editable build by default; `--plan` creates higher-level intent instead.
- `source add/update`: declare sources; preview updates by default and accept only
  with `--accept`. Browse through `show`, not duplicated source list/show verbs.
- `edit include/exclude/replace/set`: edit authored configuration through the same
  concepts as the language. `--mode copy|render` matches the include `mode` field;
  `--into`, `--reason`, and selection options mirror config. Mutations preview with
  `--dry-run`; they do not apply outputs. Exact edit-target defaults need design.
- `show`: result view by default; exclusive `--source`/`--build`/`--result` views,
  with `--interactive`, `--issues`, `--diff`, and node explanation as view options.
- `build`: check/preview by default; exclusive `--check`/`--emit`/`--apply` effects.
  Plan input always compiles through build representation; direct build needs no
  plan. `--emit` serializes a generated build, protecting manual changes; `--apply`
  materializes outputs through the same executor. `--format explicit|compact`
  selects representation, not validation or execution semantics.
- `library check/scan`: explicit library-author category; keep maintenance apart
  from consumer workflow. Scanning previews; any write behavior needs an explicit
  contract. Do not assume exhaustive library registries.
- `help [workflow|config|command]`: structured guidance rather than separate help
  flags. Machine completion can remain an auxiliary help facility if needed.

Shared candidate options: `--file` for an explicit input, `--dry-run` on authored
config mutations. No automatic choice between competing plan/build authorities.
Default discovery/edit target, recipe parameters, source pin storage, library
mutation, localization UX, and machine-output/completion placement remain open.
No legacy commands have been removed and no new CLI has been implemented.

CLI refinement under discussion: owner questions `emit`/`apply`, suggests build
installation as the default, and asks for explicit vocabulary for declarative
intent -> ordered sequence. Agent recommendation superseding the earlier proposal
if accepted: `sequence` compiles a plan and writes the generated build (no output
installation); `build` executes a build and materializes outputs by default;
`show` previews/explains and `build --check` validates without installing.
Remove `emit`/`apply` modes from the recommended vocabulary. Use consistent `--file`
input selection. `sequence` must detect stale/manual edits before overwriting a
generated build; `build` must reject a stale generated build relative to known
plan inputs rather than silently rebuild an old intent. Exact freshness/source
update and regeneration mechanics still need specification. Direct authored
builds require no plan or sequence command. Inspecting a plan resolves through
the same build representation in memory without necessarily writing a build file.
This is a recommendation pending owner answer, not implemented CLI behavior.

Owner refinement: plain `mogent build` should probably perform sequencing when
needed, or ask whether to build a sequence from the plan (abort as the negative
choice). Agent recommendation pending confirmation: one-command end-to-end build
when input authority is unambiguous. A plan-only workspace sequences then builds;
a direct-build-only workspace builds directly; a plan with its unchanged derived
build refreshes the sequence when stale, then builds. Announce the stages instead
of prompting on ordinary operation. If both files are independent, provenance is
unknown, or the generated build has manual edits, stop with actionable choices
rather than overwrite/infer authority; `--file` selects intent explicitly but
does not authorize erasing edits. Keep `sequence` as a separately callable stage
and `show`/`build --check` non-installing. This refines the earlier proposal to
require manually rerunning sequence on every stale plan; no runtime change yet.

Latest owner confirmation/refinement: preserve the authored-build protection flag
and treat writing a build by hand as optional; otherwise generate it from the
plan. Desired default workflow: `build` sequences when needed and then
materializes; standalone `sequence` stops after build generation. Authored build
origin protects against regeneration regardless of human/agent author. Generated
builds edited after generation also need change detection; a generated header
alone cannot authorize overwriting edits. Exact marker/schema and missing-state
recovery remain specification work. Optional authored file does not mean optional
build representation: every plan still passes through it. Recommend writing the
generated build for inspection; persistence/commit/input-pin ownership still needs
finalization.

Readiness checkpoint: owner asks whether the spec is nearly complete and when
implementation can start. Broad model is coherent; discussion notes are not yet
an implementation-ready spec. Next is item 8: one tiny fixture showing plan,
equivalent direct/generated build, library metadata, result tree, and blockers.
Use it to settle first-slice grammar, generated-file/pin ownership, and curation
boundaries before parser/CLI implementation. Prioritize a vertical slice reusing
content.Compose/Git/state rather than wholesale deletion or a full TUI rewrite.
Standards/harness comparison remains required before 1.0 scope/conformance claims,
but need not block a generic workspace composition prototype.

### 8. Small Orthogonal Stories On Real Repositories

Walkthrough in progress: team library has workflow guidance and review/testing/
deployment skills; project has a local review replacement. Intended result renders
intro + workflow, installs testing and local review, omits deployment, and removes
the replaced team's stale review script. Viewer accounts for defaults, omissions,
and replacement removals. This is a hypothetical fixture, not a deployed repository.

Current owner question: should a skill node declare `content = "skill.md"` and
additional scripts, can defaults avoid per-file inclusion, can its main file
represent the node, and should non-Markdown imports warn? Existing typed inventory
already detects exact SKILL.md and copies its complete enclosing bundle including
scripts/resources without per-file declarations. Proposed recommendation: keep
physical discovery; skill node can expose an entrypoint/document view without
collapsing file identity or enumerating every member. Optional entrypoint metadata
only for adapters/nonstandard layouts, exact schema unresolved. Bundle inclusion
includes all members; excluding bundle excludes them all; do not interpret omission
of entrypoint in a selector as implicit asset deletion or allow accidental partial
installation. Separate heading extraction from installing the bundle. Warn/error
on operation/type mismatch or declared dependency issues, not on every expected
non-MD bundle member. Owner has not accepted this refinement yet.

Follow-up owner response: "sounds good" with skepticism; document all behavior
and naming reservations and move on. Complete bundles/automatic assets remain
the working direction to validate. Owner questions a special `skill` keyword and
the name `entrypoint`; candidates include `is-skill`, `node-file`, or only a main
file association. No metadata keyword is settled. The full behavior, candidate
names, source tree, intended result, and proposed next consumer plan are preserved
in [MOGENT-LANGUAGE-WALKTHROUGH.md](../docs/proposals/MOGENT-LANGUAGE-WALKTHROUGH.md).
Resume at its consumer plan, then the equivalent generated build and diagnostic
variant. Do not reopen role naming unless the example exposes a concrete need.

Latest owner correction: minimize vocabulary as previously requested. Replacement
is an `include` with `replace = true`, not a separate `replace` construct. Inline
and rendered documents should also use `include`, not a separate `document`
construct. Prefer explicit operation options to a `mode` enum; current walkthrough
uses proposed `copy = true`, `render = true`, `replace = true`, and `append = true`
with `into` and inline `text`. Contradictory options need validation, not precedence
guessing. Exact shorthand/nesting is not settled. Earlier separate verbs/mode enum
examples and `edit replace` CLI proposals must be reconciled with this direction.

Owner proposes `include ... localize = true` or a destination path for external
content. Preserve localization in the include vocabulary but clarify editable
source derivation versus final output placement before defining behavior; ordinary
assembly must not silently mutate source files or repoint authored configuration.

Documentation/commit checkpoint: owner asks whether updates and commits have been
made. Docs, README feature matrix, handoff, CLI inventory, and walkthrough have
been updated; NO commits have been made during this conversation. Changes remain
in the working tree. The question alone does not instruct committing.

Further owner correction: copy/render flags still obscure ordinary inclusion.
Logical tree inclusion should traverse its Markdown leaves and produce guidance
by default; sequential contribution order should make `append` unnecessary for
one composed document. Agent recommendation following this steering: content-kind
directed defaults, not extension guessing. Logical guidance + inline text compose
one generated document in declared order; physical files/bundles preserve bytes;
conflicting whole-file/bundle destinations still need explicit replacement.
Internal copy/render operations can remain in the existing engine without public
switches. Exact full-file Markdown extraction/view selection and mixed-kind
collision validation still need examples; do not silently flatten installed skills.
The walkthrough now omits copy/render/append flags. Earlier switch-heavy sketches
are superseded for the recommended happy path.

Owner clarification: localize means an editable source in a workspace-local
library/sublibrary. It is not `into` output placement. Preserve upstream provenance
and use the local source after creation; subsequent builds must not overwrite
local edits by recopying upstream. Existing v1 copy-on-write behavior is the model;
exact include option/authoring command/path remains open.

Latest owner reduction to explore: `include "team:skills" ["review", "testing"]`
followed by `include "local:skills/review"`, with workspace target default `.agents`
and optional library-declared destination defaults; `into` is explicit placement
or an override. Owner asks about block `{}` versus object assignment `= {}` and
one-line forms, and requests node inspection showing library versus local/effective
declarations. Do not treat this rough syntax as finalized grammar.

Agent recommendation pending review: allow child-list shorthand and content
destination defaults with visible provenance, explicit `into` winning. Distinguish
library-relative suggested path from consumer-authorized target root; avoid
automatic global writes. The two-line example has equivalent selection but only
equivalent replacement behavior if the second include declares replacement or an
explicitly accepted policy supplies it. Do not silently turn source order/local
origin into overwrite permission. Consider `replace = true` in an optional include
block while leaving destinations implicit. Use `{}` for instruction options or
nested declarations; defer `= {}` object values unless a genuine data-map need
exists. Reuse `show <node> --explain` rather than adding another inspection verb:
show library declaration, local override/default acceptance, effective destination,
selection reason, origin, and replacement history. Exact addressing/CLI remains open.

Choose the pilot with the owner. Existing candidates in the deployment plan are
`todo_app_project` and `cdint-nix`; these are not newly inspected deployments.
The sibling QMR repository supplies content, not automatically organization policy.

Exercise these independently before combining them:

1. Team baseline with one deep exclusion: all other content accounted for.
2. Personal Go tag selection: show exact bundle membership and curation decisions.
3. Replace one complete skill: stale lower assets removed, provenance visible.
4. Render a guidance section with a local alternative under a stable local heading.
5. Library update changes a tag/default: show selection and output diffs before
   accepting a new resolved build.
6. Copy a native JSON/Markdown payload unchanged: inspect metadata but preserve bytes.

Owner is willing to build useful modules. Capture missing modules as source work;
do not manufacture a canonical organization library or modify a sibling repo
without agreeing the content/ownership boundary.

## Parking Lot

### Standards And Harness Support — Needed Before Release Scope Freeze

Research AGENTS.md, Agent Skills, DotAgents, agentsstandard.com, and relevant
vendor-native formats. Compare versioned primary documentation for discovery
paths, file formats, global/workspace scope, native support, and conversion needs.
Do not equate shared `.agents/` paths with support for a complete protocol.
Give a recommendation grounded in harness support and useful interoperability.
Earlier research in `reports/` concerns agentsstandard.com, not DotAgents.
Keep this research out of the first file-role discussion; return before claiming
1.0 compatibility. DotAgents lowercase `skill.md` and its frontmatter dialect
differ from Mogent's current exact `SKILL.md` bundle detection.

### Reverse Generation / Suggested Generalization

Owner idea: derive a plan from a hand-authored build, perhaps through
`apply --generate-plan`; notice tag queries or broad includes that express it
more concisely. This is not a promised lossless bidirectional compiler.

Many intentions produce the same current files. A tag query that matches today's
explicit list may select extra items after an update. Preserve exact lists as
the baseline reverse representation; semantic compression should be an optional
previewed suggestion, never silently inferred during apply. Compare compiled
manifests AND provenance/ordering before declaring a rewrite equivalent.

### Symlinks And Aliases

Owner asks whether verbatim copy could also mean symlink/alias. Existing selected
source symlinks are rejected. Copy preserves an independent payload; a symlink
references another location and has different portability, integrity, and update
semantics. Source aliases are names, not filesystem symlinks. Defer a link
operation unless a concrete workspace story requires it.

### Remote Inputs And Git

A raw HTTP URL returns a resource/file; it does not enumerate and download a
repository directory. Complete remote skills need pinned Git or a defined archive.
Prefer existing Git retrieval/pinning. Git pins describe committed inputs, not
uncommitted local bytes. Decide snapshot or verified-local-input handling before
promising reproducibility; using standard hashes is distinct from inventing a
custom content-addressed store.

## Evidence And Existing Implementation Boundaries

- `docs/DESIGN.md`: v1 authored outline, swap/reorder/localize, explicit variables.
- `docs/user-stories.md`: shared baselines, borrowing/local changes, adoption,
  curation, drift, and authority stories; some are aspirations, not features.
- `docs/proposals/MOGENT-HCL-V2-SPEC.md`: current v2 curation and sparse selection.
- `docs/thought-experiments/TE-jofit-unified-content-composition.md`: select,
  place/compose, materialize; proposed syntax is not current grammar.
- `docs/thought-experiments/TE-gunak-typed-agent-content-composition.md` and
  `docs/TYPED-CONTENT.md`: physical identity, optional document views, mixed outputs.
- `content/compose.go`: existing ordered contributions, complete manifest,
  copy/render, explicit replacement/append, bundles, link dependency checks.
- `v2/plan.go`: structural/curation resolution followed by additive tag selection;
  exception precedence needs explicit redesign, not accidental inheritance.
- `library/sidecar.go`: validates legacy relationship declarations and cycles;
  this is not general runtime enforcement or a dependency solver.
- `workspace/localize.go`, `workspace/drift.go`: implemented v1 localization and
  conservative drift import; no automatic v2 feature parity.
- `docs/dogfood/2026-10-09-191238-mogent-multirepo-deploy.md`: pilot workflow and
  owner prediction exercise; cross-repository orchestration remains unimplemented.

## Conversation Progress

- [x] Preserve owner response, priorities, unresolved questions, and tangents.
- [x] Add a README feature/status matrix and point handoff/index here.
- [x] Confirm direct builds are independent and plans must compile through builds.
- [x] Record that answer; resume at item 2 with one small include example.
- [x] Accept explicit/subtree-exception build representations and optional compact formatting.
- [x] Accept explicit exclusion over a tag match within one include, without deleting other contributions.
- [x] Accept explicit default adoption, UNDECIDED required choices blocking execution, and omitted opt-in exclusion.
- [ ] Resume at item 3: reason-required mandatory exclusions and opt-in/experimental inclusions.
- [ ] Compare both entry points using a tiny example; settle build-level selection granularity.
- [ ] Approve examples, contract, migration, and grammar before implementation.
