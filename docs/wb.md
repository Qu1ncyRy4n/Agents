# First Dogfood Target: Review Space

Rule: Source blocks are authoritative content. Concerns and options are editorial review material and must not be mistaken for source.

## User Decision

**User-provided decision (not source material):** Keep the cdint-grid TODO layout and handle-minting rules verbatim. Replace only the absolute minter executable reference with a documented reference to a universally available external tool, a Mogent-provided implementation, or an explicitly listed dependency. The delivery mechanism remains unresolved.

Minter impact: This affects only five cited cdint-grid source entries which invoke the absolute minter: Project Structure (`/home/qix/dev/cdint/cdint-grid/AGENTS.md:3-17`), Handle Minting (`/home/qix/dev/cdint/cdint-grid/AGENTS.md:18-33`), Design Notes (`/home/qix/dev/cdint/cdint-grid/AGENTS.md:34-48`), Thought Experiment Protocol (`/home/qix/dev/cdint/cdint-grid/AGENTS.md:1361-1531`), and DR Records (`/home/qix/dev/cdint/cdint-grid/AGENTS.md:1612-1642`). It is conditional: each requires resolution only when it mints a new handle; the source material itself is not blocked from being added/reviewed.

**User-provided decision (not source material):** Retain the TE protocol as
organization-canonical and generally reusable. Do not decide its artifact
locations or record-naming details from this direction alone.

**User-provided item (not source material):** Glossary placement is deferred to
a Steve question: preserve the current glossary semantics while deciding root
`GLOSSARY.md`, `docs/GLOSSARY.md`, or an adapter convention.

## Handoff Review Choices

For each remaining cited section, Steve must select one editorial option:

- **A. Canonical verbatim:** Steve accepts it verbatim into a canonical dogfood
  library and supplies scope, interface, and default selection.
- **B. Cited adapter:** Steve permits a precise cited adapter or change and
  records it.
- **C. Optional/local overlay:** Steve marks it optional or local overlay and
  says whether Mogent manages or renders it.
- **D. Source-project documentation:** Steve retains it in source-project
  documentation outside Mogent.

Individual review sections are paused pending this choice.

## Build, Test, and Development Commands

Source (verbatim): /home/qix/dev/cdint/cdint-grid/AGENTS.md:230-262 (## Build, Test, and Development Commands).

```markdown
## Build, Test, and Development Commands
- `go test ./...` runs the test suite.
- `go -C x/cas-native-index test -timeout=30m ./...` runs the isolated
  CAS-native-index experiment suite with an explicit 30-minute Go test timeout;
  root-module wildcard commands do not cross the nested module boundary. Do not
  rely on Go's 10-minute default for this suite. Source: DI-gudop; DI-homis
- `gofmt -w .` (or `go fmt ./...`) formats Go code.
- Runtime or agent process-flow execution that exercises `grid` stage0/stage1
  behavior on this development machine should run through the POC16-style Docker
  Compose harness. Unit tests, static checks, code generation, handle minting,
  Git operations, and static file inspection may run natively unless they launch
  `grid` runtime or agent process flows. Later production deployment may run
  native binaries on target machines. Source: DI-vojid
- `docker compose build grid` builds the local runtime image for containerized
  stage0/stage1 process-flow checks. `docker compose run --rm grid help` runs the
  scaffold CLI inside that image. Source: DI-zoran
- Every shipped or production-shaped Go executable must compile with
  `CGO_ENABLED=0`. Dependencies may contain optional CGO implementations only
  when the selected production build succeeds without them. Dedicated
  `go test -race` checks may enable CGO, but their binaries are test-only and
  must not be treated as production artifacts. Seeing a C compiler during a race
  check does not by itself prove a production CGO dependency. Source: DI-pihav
- For the CAS-native-index module, use
  `CGO_ENABLED=0 go -C x/cas-native-index test -run '^$' ./...` as the bounded
  compile gate when no test execution is needed. Source: DI-pihav
- Before a `GOPROXY=off` Go gate whose shared module-cache readiness is not
  already proved, preload the same approved `GOMODCACHE` with `go mod download`
  while network access is explicitly allowed. If a worker may not use the
  network, main must preload the cache before launch or the handoff must identify
  evidence that the required modules are present. Keep the validation command
  offline and do not apply this rule to tests intentionally exercising an empty
  or missing-dependency cache. Source: DI-punon

```

Concerns

- Review concern (editorial): The named local Docker Compose harness and `docker compose` commands require this development environment's container tooling. Source: /home/qix/dev/cdint/cdint-grid/AGENTS.md:237-245.

## Agent Instruction Architecture

Source (verbatim): /home/qix/dev/cdint/cdint-grid/AGENTS.md:263-266 (## Agent Instruction Architecture (Required)).

```markdown
## Agent Instruction Architecture (Required)
- `AGENTS.md` is the canonical home for repo-wide protocol, workflow, and vocabulary rules. Role-specific files such as `AGENTS-codex.md` and `AGENTS-ppx.md` are overlays: they may add stricter role constraints, environment setup, and role-specific procedures, but they must not duplicate or relax canonical repo-wide rules. (DI-034-20260508-060134)
- When a rule applies to every agent or every repo artifact, move it here and replace role-file copies with pointers. When a rule applies only to one agent's runtime environment, identity, credentials, branch lifecycle, or private logging system, keep it in that role overlay. (DI-034-20260508-060134)

```

Concerns

- Review concern (editorial): The canonical and overlay filenames couple the rules to this repository agent-role layout. Source: /home/qix/dev/cdint/cdint-grid/AGENTS.md:263-265.

## Architecture Question Context And Vocabulary

Source (verbatim): /home/qix/dev/cdint/cdint-grid/AGENTS.md:308-376 (## Architecture Question Context And Vocabulary (Required)).

```markdown
## Architecture Question Context And Vocabulary (Required)
- Read `GLOSSARY.md` before asking an architecture question or comparing files,
  data structures, databases, indexes, messages, or other named concepts. Use
  its agreed terms exactly and preserve the distinctions it records. Source:
  DI-lodos
- Before asking such a question, provide enough plain-English context for the
  user to evaluate the choice. State all of the following:
  - the exact purpose or question the thing answers;
  - what it contains, with at least one concrete example;
  - whether it is durable evidence, mutable local decision state, or
    disposable and rebuildable data;
  - who writes it, who reads it, and when it is used during startup or normal
    operation;
  - how it differs from nearby files or data structures;
  - the alternatives actually under consideration;
  - the full category being discussed before naming one example from that
    category; and
  - the result of re-reading the governing docs or code, including whether each
    premise is a locked decision, a proposal, an experiment, or an unresolved
    question.
  Source: DI-lodos
- Before asking any question whose answer depends on how agents interact, first
  explain one concrete example in ordinary language. When named actors add
  clarity, prefer stable names whose first letters suggest their roles, such as
  Olivia for order processing, Sam for synchronization, and Quinn for
  QuickBooks. Do not rename an already clear actor or add an alphabetical alias
  merely to satisfy a naming convention. Introduce each person's role explicitly
  before describing actions, preferably in a separate sentence such as "Olivia
  is the order-processing role." Do not rely on later verbs to imply the role,
  and keep each person's role stable throughout the example. Describe what each
  person does, who
  controls each decision, what remains stored, and what happens in the
  relevant failure case. Do this before introducing specialized terms.
  Apply this rule to questions about promises, trust, authority,
  messages or data moving between agents, stored state, ownership, and
  failures. Do not force named people into a mechanical question where
  they add no clarity. Source: DI-sohub; DI-rupas; DI-jadut
- An abstract list of purpose, contents, authority, readers, writers, and
  alternatives does not replace the concrete named-actor explanation. After the
  example, map only the necessary agreed technical terms to what the named
  actors did, then ask one short question using the same plain language.
  Source: DI-sohub; DI-jadut
- Treat `ABC` from the user as an instruction to repeat the immediately
  preceding question, explanation, or both in plain English using stable named
  actors whose names suggest their roles when natural. `ABC` names the form of
  the explanation; it does not require the literal names Alice, Bob, and Carol.
  State each person's role explicitly before the actions. Preserve the meaning
  and choices. Do not answer the question,
  advance the discussion, change an alternative, or introduce a new design.
  Source: DI-sohub; DI-rupas; DI-jadut
- Before sending the question, check it for a false premise, an unexplained or
  overloaded term, a proposal described as an existing decision, confusion
  between durable evidence and a disposable index, and an example incorrectly
  presented as the whole category. Correct those problems before asking the
  user. Source: DI-lodos
- Do not invent jargon or assign a new conceptual label merely to shorten a
  discussion. A new term requires DF and a DI before it becomes agreed
  vocabulary. Until then, describe the thing by its purpose or use a
  `TBD-<slug>` placeholder and add that placeholder to `GLOSSARY.md` with its
  unresolved status and governing DR or TE. Source: DI-lodos
- Maintain `GLOSSARY.md` as the repository vocabulary source. Each specialized
  term must state its purpose, a concrete example when useful, its authority or
  lifecycle, related terms it must not be confused with, its decision status,
  and its governing DI, DR, TE, specification, or external standard. A word's
  appearance in historical prose does not make it agreed vocabulary. Preserve
  historical text under the applicable editing policy and mark obsolete or
  disputed terms in the glossary instead of silently legitimizing them. Source:
  DI-lodos

```

Concerns

- Review concern (editorial): The required `GLOSSARY.md` file is repository-layout coupling. Source: /home/qix/dev/cdint/cdint-grid/AGENTS.md:308-312.

## Architecture Precedent And Role Checks

Source (verbatim): /home/qix/dev/cdint/cdint-grid/AGENTS.md:377-410 (## Architecture Precedent And Role Checks (Required)).

```markdown
## Architecture Precedent And Role Checks (Required)
- For architecture analysis, implementation planning, or review that depends
  on how PromiseGrid's VCS, projects, timelines, physical storage, filesystems,
  execution, identity, trust, economics, or governance fit together, consult
  the relevant sections of
  `docs/DN-lazit-grid-as-vcs-mental-model.md` before changing the design. Keep
  narrower protocols, TEs, DNs, DRs, and TODOs authoritative for their own
  scopes. If DN-lazit and a narrower source appear to conflict, stop only the
  work that depends on the conflict, preserve both statements, and obtain a
  human decision recorded by a new DI. Do not choose automatically by filename,
  date, apparent scope, or worker preference. Source: DI-hikab
- During architecture analysis, thought experiments, implementation planning,
  and review, actively look for a cleaner model that reuses established CAS/VCS
  mechanisms and existing agreed PromiseGrid concepts. Prefer that model when
  it preserves the required semantics and failure behavior. If a distinct
  mechanism remains necessary, state the concrete mismatch that prevents reuse.
  Surface this comparison before final DF or implementation rather than waiting
  for the repo owner to identify it. Source: DI-naliv
- Before making an architecture claim or asking a DF question about identity,
  boot, deployment, updates, storage, or security, distinguish the
  administrative principal, agent identity, software process, host, credential
  or key, and managed resource. Do not transfer a constraint on one of those
  things to another without stating and supporting the relationship. Source:
  DI-jujus
- When established practice can materially change the answer, consult primary
  documentation for at least two applicable real systems before presenting the
  claim or question. State whether each conclusion is established practice, a
  PromiseGrid-specific requirement, or an analogy; outside systems inform the
  comparison but do not select its outcome. Source: DI-jujus
- Never turn a per-process, per-writer, per-agent-identity, or per-key constraint
  into a limit on how many machines or resources an administrative principal may
  manage. In particular, do not assume one sysadmin agent or sysadmin private key
  per managed host. Source: DI-monop; DI-jujus

```

Concerns

- Review concern (editorial): The required PromiseGrid design note is repository- and project-specific. Source: /home/qix/dev/cdint/cdint-grid/AGENTS.md:377-387.

## Experiment Cutoffs And Production Requirements

Source (verbatim): /home/qix/dev/cdint/cdint-grid/AGENTS.md:482-526 (## Experiment Cutoffs And Production Requirements (Required)).

```markdown
## Experiment Cutoffs And Production Requirements (Required)

- When selecting or reporting a numeric performance threshold, state whether it
  is an experiment comparison cutoff or a production requirement. An experiment
  comparison cutoff is a convenient value used to compare or bound a test. A
  production requirement must instead name the real operating scenario, the
  user or system impact being limited, the supporting measurements, the
  rationale for the numeric value, and the owner DF/DI that approved it. Do not
  silently promote a comparison cutoff into a launch gate or production
  requirement. Source: DI-lajul
- Report a missed cutoff literally before using any summary label. State how
  many candidates were tested, the measured duration or range, the exact cutoff,
  and whether that cutoff was arbitrary or operationally justified. For example,
  say "all four 100k candidates ran longer than the arbitrary three-minute
  comparison cutoff." Do not replace that statement with jargon or euphemisms
  such as `zero-passer`, `qualification failure`, `no-go`, `hard-gate failure`,
  or `unusable`. Source: DI-lajul
- Keep three conclusions separate: whether the run executed correctly, how its
  measurements compare with the experiment's cutoff, and whether the design
  meets justified production needs. A valid run may miss an arbitrary cutoff
  without proving anything about production suitability. Only a missed
  operational requirement may block the dependent production milestone, and
  its report must name the resulting user or system impact. Source: DI-lajul
- When every candidate misses an experiment comparison cutoff, tell the user the
  literal result immediately, before unrelated coordination continues, and stop
  automatic unchanged retries. Preserve the evidence. A bounded profile may
  identify where time or resources were spent, but do not optimize solely to
  beat the cutoff, invent a replacement target, or classify the result as launch-
  blocking without a new owner decision. Repeat an unchanged expensive run only
  for an explicitly approved reproducibility purpose. Source: DI-lajul; DI-vigah
- Keep an actual launch-blocking correctness or safety failure, data-loss or
  live-system risk, or missed justified operational requirement visibly open
  until the repo owner acknowledges the literal result and approves, rejects, or
  changes the corrective direction. Record that acknowledgment in the governing
  DI, DR, or TODO before clearing the alert. Do not demand launch-blocker
  acknowledgment for an arbitrary comparison cutoff; record its interpretation
  instead. Source: DI-vigah
- Before changing code for performance, derive the relevant requirements from
  actual workflows and measure the stages that can affect them. For cdint-grid
  this includes foreground customer and invoice lookup, OSC-order duplicate
  checks, incremental observation ingest, activity-sensitive background
  QuickBooks scanning, restart, complete index rebuild, and recovery. Optimize
  the measured stage only after the applicable requirement is locked. Source:
  DI-lajul

```

Concerns

- Review concern (editorial): The cited workflows name cdint-grid, OSC, and QuickBooks, coupling performance work to this repository and external system domain. Source: /home/qix/dev/cdint/cdint-grid/AGENTS.md:519-525.

## Thought Experiment Protocol

Source (verbatim): /home/qix/dev/cdint/cdint-grid/AGENTS.md:1361-1531 (## Thought Experiment Protocol (Required)).

```markdown
## Thought Experiment Protocol (Required)
- Before locking any non-trivial decision that will require DF questions and answers, the agent must run a thought experiment (TE) if multiple plausible designs remain.
- The agent MUST NOT prejudice a thought experiment outcome when planning a thought experiment -- the agent must not pre-select a preferred alternative or design, and must not bias the TE toward a particular outcome.
- A TE happens before final DF questions. Its purpose is to narrow the design space so DF questions and answers are informed by explicit scenario analysis.
- The agent must not collapse a TE into a short opinion or recommendation. The agent must explicitly model concrete scenarios and consequences.
- Each new TE must have a unique proquint handle in the format `TE-<handle>`, where `<handle>` is minted by `/home/stevegt/bin/mint-handle` from the global TODO/TE/DR/DI/DN handle namespace. Existing pre-upgrade TE handles and prior aliases remain historical records. Source: DI-jufiz
- The TE doc filename must be `TE-<handle>-<slug>.md` and live under `docs/thought-experiments/`, for example: `docs/thought-experiments/TE-mumuv-naming-reconciliation.md`. The slug is informational and may be edited; the proquint handle is permanent.

### TE Intake Requirements
- Before locking decisions or asking final DF questions, the agent must identify:
  - the decision being tested,
  - the candidate alternatives,
  - the assumptions and threat/trust model,
  - the scope and systems affected.
- If the TE relates to an existing TODO, the agent must reference the TODO handle and subtask handle (for example, `fonuz.1`).

### TE Execution Requirements
- Each TE must evaluate the same decision across multiple concrete scenarios.
- Scenarios must include, when relevant:
  - normal operation,
  - failure/corruption/incomplete writes,
  - concurrent actors or mixed-version nodes,
  - long-horizon evolution and migration,
  - trust-boundary changes,
  - scale effects (storage, bandwidth, CPU, operational complexity).
- The agent must compare alternatives under the same assumptions instead of switching assumptions mid-analysis.
- The agent must state what each alternative makes easier, what it makes harder, and what new obligations it creates.

### TE Authoring Conventions
- When named actors are useful, prefer stable human names whose first letters
  naturally suggest their roles, and state each role explicitly before its
  actions. Reuse established role-matching names instead of adding alphabetical
  aliases. Mallory remains suitable when a malicious actor is relevant. Name
  Steve explicitly only when his repo-owner role is necessary to the scenario.
  Apply this convention in TEs, scenario analyses, tabletop simulations, DR/DI
  prose, and worked examples in specs or docs. (DI-jadut; supersedes
  DI-034-20260508-060134 only for actor naming)

### TE Output to DF
- After the TE, the agent must identify:
  - rejected alternatives,
  - surviving alternatives,
  - unresolved questions that still require user choice,
  - any new naming/path/runtime decisions exposed by the TE.
- Final DF questions must be framed from the surviving alternatives identified by the TE. The agent must not ask broad DF questions that ignore TE results.

### TE Artifacts
- The agent must track required TEs in the relevant root
  `TODO/TODO-<handle>-<slug>.md` file.
- For each completed TE, the agent must write a verbatim copy of the thought experiment into a standalone file under `docs/thought-experiments/`.
- The doc filename must begin with the TE ID and then use a descriptive suffix.
- The doc must stand on its own and include:
  - title,
  - TE ID,
  - decision under test,
  - assumptions,
  - alternatives,
  - scenario analysis,
  - conclusions,
  - implications for the repo's open TODOs and pending DIs.

### TE Decision Rules
- A TE does not by itself lock a decision.
- After the TE, the agent must either:
  - ask the user to choose among the surviving alternatives, or
  - recommend one surviving alternative and clearly state why the others were rejected.
- After user choice is resolved, the agent must record the locked result via the existing DI process before implementation.
- If a TE exposes a new ambiguity, dependency, or naming/path decision, the agent must stop and resolve that before implementation.

### TE Final Handoff Requirements
- In the final response for TE work, the agent must include:
  - which TE was completed,
  - the TE ID,
  - the doc path under `docs/thought-experiments/`,
  - the surviving alternatives,
  - the recommended conclusion or the exact DF question that remains for the user.
- Hard gate: for decisions that require a TE, work is incomplete until the TE doc exists and the resulting decision status is explicit (`needs DF`, `locked`, or `deferred`).

### TE Editing Policy (Required)

Once a TE is filed in `docs/thought-experiments/`, edits to it follow a categorized policy. The policy originated in `DI-020-20260502-213103` (categorized editing regimes), `DI-020-20260502-213104` (uniform applicability across all TE corpora; this rule applies wherever TEs are stored, not just under `docs/thought-experiments/`), and `DI-020-20260502-213105` (holistic reading by default for substantive questions, single-TE reading allowed for obviously mechanical ones). The Cat-1 clause of `DI-020-20260502-213103` was superseded on 2026-05-02 by `DI-020-20260502-232651` (Cat-1a / Cat-1b split). Those identifiers remain historical provenance; the complete operative local policy is the text below, and no absent external policy artifact is a prerequisite for editing a TE. Source: DI-tavil

The seven categories are:

- **Cat-1a (current-pointer paths).** A path reference that names the current location of a file. Mechanical sweep in place; no top-of-file note required.
- **Cat-1b (historical-quotation paths).** A path reference that quotes an earlier corpus state — inside a markdown blockquote, attributed to another TE ("TE-N states ..."), in past tense ("TE-magup used the path ..."), inside a `## Refinements` section, supersedence note, or `Decision status` line. Left untouched; rewriting would falsify the historical record. Per-match classification with five heuristics (quotation context; Refinements / supersedence framing; past tense; default Cat-1a; when-in-doubt-Cat-1b). Sweep tools may emit matches with surrounding context for human review but must not auto-rewrite.
- **Cat-2 (vocabulary updates).** A rename of a term whose meaning is unchanged (typo fixes, terminology consolidation). In place, with a top-of-file note pointing at the driving TE or TODO. The note must enumerate by ID every DI that lives in the affected TE, paired with an explicit promise that the rewrite preserves each DI's meaning. A TE without DIs gets a one-line `no DIs in this file` note. Form: `Cat-2 vocabulary update per <driving TE or TODO>: '<old term>' -> '<new term>'. The following DIs in this file are unchanged in meaning: DI-XXX-..., DI-YYY-..., DI-ZZZ-... .` Mandatory pre-step: grep the entire corpus for the old term inside quotation contexts (markdown blockquotes; fenced code blocks presented as citations; single/double-quoted phrases attributed to another TE via `TE-N states`, `TE-N reads`, `originally said`, `as of TE-N`, `the corpus showed`); each match is classified Cat-2 (sweep) or Cat-2-historical (leave) per the same heuristics as Cat-1a/Cat-1b.
- **Cat-3 (navigational forward pointers).** Append a dated entry to the TE's `## Refinements` section (created if absent, placed after `## Decision status`) describing where the affected reader should now look. The TE body above is unchanged. No DI is filed for a Cat-3 entry. Procedural tightenings of an existing category's how-to are Cat-3.
- **Cat-4 (resolved-implication forward pointers).** Same shape as Cat-3, used when an item from the TE's `Implications and future work` list has resolved (a TODO filed; a DR opened; a downstream TE landed). Append-only; no body edit.
- **Cat-5 / Cat-6 / Cat-7 (substantive supersedence).** A material change to a locked DI's meaning, scope, or applicability requires a new TE that supersedes the affected one. The new TE carries its own DFs and DIs; the older TE's `## Decision status` is updated to `superseded by TE-<id>` and its top-of-file `## Status` field is updated to `superseded by TE-<id> / DI-<id>`. The older TE's body is otherwise untouched.

Every TE in the corpus carries a top-of-file `## Status` field placed immediately after the TE ID line. Canonical values: `needs DF`, `decided`, `decided, refined`, `superseded by TE-<id> / DI-<id>`, `withdrawn`. Legacy values preserved during retrofit: `stub`, `open`, `recommended for immediate adoption`, `locked for the <protocol>`. New TEs prefer canonical values; the field is updated by Cat-1a sweep when the TE's state changes.

The `## Refinements` section is the single append-only home for Cat-3 / Cat-4 entries on a TE. Entries are dated (`### YYYY-MM-DD — <title>`) and ordered chronologically. The body of the TE above the `## Refinements` section is treated as historical evidence: a Cat-1a path-rename or Cat-2 vocabulary sweep on the body is permitted under its category rules; a Cat-3 / Cat-4 forward-pointer is appended to `## Refinements` rather than rewriting the body; a Cat-5 / Cat-6 / Cat-7 substantive change is filed as a new superseding TE rather than as an edit. A procedural tightening that preserves locked meaning belongs in `## Refinements`; it does not require a new DI. Source: DI-tavil

Reading default: holistic. When deciding whether an edit is mechanical or substantive, when interpreting a single TE's claims, or when reasoning about whether a refinement is Cat-3 or Cat-5–7, the agent must read the affected TE, locally present TEs it cites or that cite it, and locally present editing-policy refinements relevant to the question. An absent external or historical artifact does not block the edit; apply the complete local policy and record any material limitation. Single-TE reading is reserved for obviously mechanical questions (a single typo; a path that has demonstrably moved; a Status field retrofit). When in doubt, read holistically. Source: DI-tavil

Applicability: this policy applies uniformly to every TE corpus in this repository, regardless of which protocol or harness directory it lives in. Per-protocol corpora may add stricter rules but may not relax these rules.

### Naming Decisions (Required)
- The agent must not invent function names or variable names that are not already covered by locked naming decisions.
- If naming is not covered, the agent must stop and ask multiple-choice naming options before continuing.

### File/Path Decisions (Required)
- Path approvals are mandatory for all touched paths:
  - repo-changed files (create/rename/move/delete),
  - runtime touched paths (read/write/delete), including input files, output files, DB files, caches, fixtures, and temporary test files.
- The agent must ask path approvals one path at a time via multiple-choice questions.
- Path-question order must be dependency order.
- Each path question must include: action, exact path (or approved dynamic pattern ID), purpose, class (`prod-code | prod-data | test | temp`), and lifecycle intent.
- Temporary test paths require explicit approval and an explicit cleanup plan before handoff.
- Dynamic/runtime-generated paths must be approved by pattern, with:
  - allowed root bounds,
  - allowed actions,
  - concrete examples.
- The agent must ask one multiple-choice approval per dynamic path pattern.
- A path matched by a tracked repository `.gitignore` is standing-approved for
  local create, read, update, and delete operations and does not require an
  individual path question. Before relying on this approval, use
  `git check-ignore -v --no-index -- <path>` and confirm that the reported ignore
  source is itself tracked by this repository. The artifact must remain
  untracked and must never be committed. A handoff may still name the path and
  lifecycle for clarity, but omission does not make an ignored artifact a path
  violation. Ignore status does not authorize secrets, production data,
  external side effects, unreasonable disk usage, or bypass of an independent
  privacy, security, retention, cleanup, or live-system rule. Local
  `.git/info/exclude`, global ignore files, and untracked `.gitignore` files do
  not grant this approval. Source: DI-sohig
- Run normal Go commands through `tools/attempt-env -- COMMAND ...`. It sets
  `GOCACHE=/tmp/cdint-grid-gocache/build` and
  `GOMODCACHE=/tmp/cdint-grid-gocache/modules` so workers share downloads and
  compiled packages. It clears inherited `GOROOT` and `GOPATH` before launching
  the command so a long-running process cannot pair the current Go command with
  compiler tools or installation paths from an older version. Use
  `--isolated-go-cache PURPOSE` only for a test whose subject is empty-cache or
  missing-dependency behavior. Shared-cache cleanup requires separate approval.
  Source: DI-subaz; DI-dogas; DI-soluj
- Main reproducibly builds and atomically installs the shared coordination
  executables `/home/stevegt/bin/cdint-grid-workers`,
  `/home/stevegt/bin/cdint-grid-worker-inbox`,
  `/home/stevegt/bin/cdint-grid-skills-sync`,
  `/home/stevegt/bin/cdint-grid-token-efficiency`,
  `/home/stevegt/bin/cdint-grid-dobab-calendar`, and
  `/home/stevegt/bin/cdint-grid-merge-queue` from tracked source. Ordinary
  workers use those installed commands and do not rebuild private copies merely
  to inventory workers, load context, update instructions, inspect the calendar,
  or coordinate a merge. The instruction synchronizer copies tracked source but
  never builds, installs, or removes an executable. Source: DI-dogas;
  DI-fotur; DI-kolat
- If any unapproved runtime path appears, the agent must stop and ask before continuing.

### Decision Lock and Stop Rule
- The agent must produce a Decision Lock summary with decision IDs before code edits begin.
- The agent must not proceed if any required decision is missing, ambiguous, or conflicting.
- The agent must stop and ask immediately if a new decision need appears during implementation.
- The agent must not assume defaults for locked categories unless the user explicitly approves defaults.

### Compliance Ownership (Agent)
- The agent must treat user decisions as authoritative and implement to those decisions.
- The agent must run a compliance self-review before finalizing and must fix all non-compliance before handoff.
- Hard gate: work is incomplete until compliance is PASS, or the user explicitly approves an exception.
- The user should not need to manually inspect diffs to determine compliance.

### Required final handoff artifacts
- `Decision Compliance: PASS/FAIL`
- Decision Matrix mapping each locked decision ID to implementation evidence.
- Inline diff annotations in the form `path:line -> decision_id -> rationale`.
- Runtime Path Touch Matrix listing each approved runtime path/pattern, action used, and where it is implemented/validated.
- `Exceptions:` listing only user-approved deviations.
- Every non-trivial behavior change must include intent provenance per existing DI requirements.

```

Concerns

- Review concern (editorial): The handle executable and `docs/thought-experiments/` location couple TE artifacts to a host and repository layout. Source: /home/qix/dev/cdint/cdint-grid/AGENTS.md:1366-1367.

## DR Records

Source (verbatim): /home/qix/dev/cdint/cdint-grid/AGENTS.md:1612-1642 (# DR Records; section text before Testing Guidelines).

```markdown
# DR Records

The DR/ directory stores Decision Request (DR) records for coordination work.

Rules:
- One DR per file.
- DR files are append-only event logs.
- Keep TODO files as snapshots; link TODOs to DR files for open questions.
- Person identity format: `user@example.com (FirstName)`.

Recommended file naming:
- `DR-<handle>-<slug>.md`, where `<handle>` is minted by `/home/stevegt/bin/mint-handle`. Source: DI-jufiz

Required DR fields:
- `DR-ID`
- `Date`
- `Asked by` (person identity format above)
- `State` (`open | decided | blocked | implemented | closed`)
- `Question`
- `Why this blocks progress`
- `Affects` (repos/files/components)
- `Unblocks` (TODO handles/tasks)
- `Waiting on` (person identity format above, or DI ID)
- `Decision` (filled when decided)
- `Linked DI`
- `Related commits`
- `Last updated`

Reference pattern:
- From TODO files: `../DR/<filename>.md`

```

**User decision (not source material):** DR records should stay as close to cdint-grid as possible; document and sketch future extraction into guides and skills. This is a user decision, not an implementation command.

Concerns

- Review concern (editorial): The `DR/` directory and TODO-relative reference pattern couple records to this repository layout. Source: /home/qix/dev/cdint/cdint-grid/AGENTS.md:1614-1642.

## Testing Guidelines

Source (verbatim): /home/qix/dev/cdint/cdint-grid/AGENTS.md:1644-1686 (## Testing Guidelines).

```markdown
## Testing Guidelines
- Use Go's standard `testing` package with deterministic tests.
- Avoid network calls in tests unless explicitly required and documented.
- When changing `plan/run` behavior, add coverage for both command paths when possible.
- Run a command with a pseudo-terminal when that command, its wrapper, or the
  login-shell startup path is known to depend on terminal detection. After any
  invocation emits `not a tty`, main and every worker must request a PTY for all
  future invocations through the same execution path unless the test
  intentionally exercises non-TTY behavior. Never count output containing
  `not a tty` as a successful check, regardless of exit status. With a PTY,
  distinguish any login-shell terminal-device line from the target program's
  output and verify the executable when uncertain. Pure file inspection and Git
  plumbing do not require a PTY solely because one is available. Source:
  DI-basap
- When a command runner returns a session or continuation identifier without a
  terminal exit code, the command is still running. Retain the exact
  command-to-session association and use that runner's continuation mechanism
  until the same command returns a terminal exit status. Parallel execution is
  permitted only when the orchestrator retains every association and does not
  finish while a child session remains unhandled. A missing exit code is not
  success, failure, or a test result. Prefer serial execution when complete
  tracking cannot be guaranteed. If the identifier is lost, preserve the
  available output and apply DI-kogaf's human-in-the-loop stop. Source:
  DI-zusur
- Use the Go parser built into `branch-accounting` as the sole authoritative
  parser for its TSV ledger. Quoted fields may contain tabs and newlines, so do
  not use `awk` or another line-oriented command to count, summarize,
  transform, or validate ledger records. `csvtool -t TAB` may be used for
  record-aware inspection, but it does not replace `branch-accounting` semantic
  validation against Git and recorded evidence. Ordinary text search is not
  record parsing and must not support a row-level conclusion. Source: DI-gutog
- Treat any inconclusive, unexpected, timed-out, contradictory, malformed, or
  otherwise uninterpretable test, experiment, validation, or diagnostic result
  as an immediate human-in-the-loop stop. Preserve the raw output and current
  process state, then tell the repository owner in chat, in plain English, what
  command ran, what result was expected, what actually happened, and why the
  result cannot be accepted. Do this before rerunning, changing a timeout,
  weakening a check, trying a workaround, committing, recording a disposition
  in documentation, or continuing the current or dependent work. Do not make a
  commit message, TODO, design note, worker result, or inbox notice the first or
  only disclosure. Wait indefinitely for the repository owner's direction;
  resume only after that direction is explicit. Source: DI-kogaf

```

Concerns

- Review concern (editorial): The Go-specific testing requirement and `branch-accounting` parser couple these instructions to this repository language and toolchain. Source: /home/qix/dev/cdint/cdint-grid/AGENTS.md:1644-1645,1668-1674.

## Commit & Pull Request Guidelines

Source (verbatim): /home/qix/dev/cdint/cdint-grid/AGENTS.md:1687-1695 (## Commit & Pull Request Guidelines).

```markdown
## Commit & Pull Request Guidelines
- Treat a line containing only `commit` as: add and commit all changes with an AGENTS-compliant message.
- Use short, imperative, capitalized commit subjects.
- Summarize changes per file in commit bodies.
- Stage files explicitly (avoid `git add .` / `git add -A`).
- Do not open GitHub pull requests for normal wire-lab convergence. Steve explicitly dropped the require-PR merge rule; merge by the role-specific direct-push workflow instead. (DI-001-20260428-195702; DI-034-20260508-060134)
- Do not force-push repo branches unless a scoped DI/DR explicitly authorizes the exception. Role overlays may add stricter no-force-push rules for their branches or may document a narrow private-remote exception, but the repo-wide default is no history rewrites. (DI-034-20260508-060134)
- Do not commit local state files, generated binaries, credentials, tokens, signing keys, or other secrets. (DI-034-20260508-060134)

```

Concerns

- Review concern (editorial): The wire-lab convergence, Steve, GitHub, and role-specific workflow references are repository and role coupling. Source: /home/qix/dev/cdint/cdint-grid/AGENTS.md:1692-1693.

## Vocabulary Source

Source (verbatim): /home/qix/dev/cdint/cdint-grid/AGENTS.md:1696-1701 (## Vocabulary Source).

```markdown
## Vocabulary Source
- `GLOSSARY.md` contains agreed terms, explicitly unresolved `TBD-<slug>`
  placeholders, external standard terms used by this repository, and historical
  terms that must not be mistaken for current vocabulary. Read and maintain it
  according to the Architecture Question Context And Vocabulary rules above.
  Source: DI-lodos
```

Concerns

- Review concern (editorial): The `GLOSSARY.md` dependency couples the vocabulary rule to this repository file layout. Source: /home/qix/dev/cdint/cdint-grid/AGENTS.md:1696-1701.
