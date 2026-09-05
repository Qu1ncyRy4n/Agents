# Repository Guidelines

## Project Structure & Module Organization
- Keep packages at module root or under purpose-named top-level directories (`contexts/`, `state/`, etc.); avoid `internal/` and `pkg/`.
- Keep planning artifacts in the root `TODO/` directory for this repo and maintain
  `TODO/TODO.md` sorted by priority. TODO files use
  `TODO/TODO-<handle>-<slug>.md`, where `<handle>` is minted by
  `/home/stevegt/bin/mint-handle`. The per-protocol `protocols/<slug>.d/TODO/` layout used
  by wire-lab is not the active layout for cdint-grid unless a later DI migrates
  this repo. Source: DI-gigoz; DI-topih; DI-sojum; DI-jufiz
- Keep durable worker-attempt records beside their one owning TODO under
  `TODO/TODO-<handle>-<slug>.d/<full-attempt-id>/`. The TODO Markdown file and
  matching `.d` directory have the same full stem and must move together with
  `git mv` when the TODO slug changes. Related TODOs link to the owning copy;
  they do not duplicate it. Source: DI-jifus
- Do not commit local state files (for example `.grok`, `.grok.lock`) or generated binaries.

## Handle Minting (Required)
- From the repository root, mint a new TODO, TE, DR, DI, or DN handle with
  `/home/stevegt/bin/mint-handle -r .`. The default `-w 1` output is the canonical
  five-character proquint for new coordination artifacts. Source: DI-sojum;
  DI-dodud; DI-jufiz
- The command scans the whole working tree below `-r`, excluding `.git`, and
  prints one currently unused handle to stdout. Persist that handle promptly in
  the artifact filename or DI owner line; minting does not maintain a separate
  reservation registry. Source: DI-sojum
- Use `-w 2` only when a governing decision explicitly requires an
  eleven-character proquint. Do not use `-n` or `-s` for normal allocation;
  they exist for dry-run and deterministic test cases. Source: DI-sojum
- Keep the repo-local `tools/mint-handle` source and tests until a separate
  decision removes or migrates them; it is not the preferred executable.
  Source: DI-jufiz

## Design Notes (Required)
- A design note is a long-lived explanatory synthesis. It is not by itself a
  thought experiment, decision request, locked decision, or protocol
  specification. Cite DIs for settled claims and DRs for unresolved questions.
- Store design notes directly under `docs/` using
  `docs/DN-<handle>-<slug>.md`. Mint `<handle>` with
  `/home/stevegt/bin/mint-handle -r .`; the handle is permanent and globally unique
  across TODO, TE, DR, DI, and DN artifacts, while the descriptive slug may
  change. Source: DI-jufiz
- When assigning an ID to an existing note, mint the handle, rename a tracked
  file with `git mv` or an untracked file with a filesystem move, and update all
  current-pointer references. Preserve historical paths in append-only decision
  records and add a superseding DI when the naming decision changes. Source:
  DI-dodud

## TODO-dobab Update Protocol (Required)
- A request to `update dobab` means reconcile and replan the complete active
  deadline record in `TODO/TODO-dobab-august-osc-qb-production.md`. Merely
  appending evidence, changing a date heading, or listing unfinished work under
  an elapsed date is not a complete update. Source: DI-tijab; DI-nobam;
  DI-zutuh
- Determine the current calendar date in `America/Los_Angeles`, then read the
  component milestone table, dependency matrix, Daily Calendar, checklist,
  current Git history, and relevant test or target evidence before editing.
- Keep a `## Table Of Contents` immediately below the TODO-dobab title. It must
  link every other `##` section in document order and must not include itself,
  the document title, or any `###` subsection except the one required `Current
  Day` link. Nest `Current Day` beneath `Daily Calendar`, keep its target at
  `#current-day`, and place exactly one raw
  `<a id="current-day" name="current-day"></a>` anchor immediately before the
  unique dated `Active` heading, or before the first dated `Forecast` heading in
  document order when no day is active. Move only that anchor when the selected
  day changes; the dated heading retains its generated Gitea fragment as a
  second direct target. More than one active day, a missing or duplicate anchor,
  mismatched `id` and `name` values, or an anchor before the wrong day is an
  error.
  Whenever section order, heading text, or Daily Calendar status changes, update
  the contents list and movable-anchor placement in the same change, then verify
  every top-level anchor and order, the stable `#current-day` target, and the
  selected dated-heading fragment in the Gitea rendering. Retest both anchor
  attributes before accepting a Gitea renderer upgrade. Source: DI-zuhob;
  DI-lokug; DI-juhun; DI-mutag; DI-sovik
- Within each dated Daily Calendar subsection, make the first mention of each
  distinct TODO, TE, or DR artifact a Markdown link to its relevant file. Use a
  repository-relative link for a cdint-grid artifact and the canonical remote
  file URL for an artifact in another repository. Leave later mentions of the
  same artifact within that dated subsection as plain text. Do not apply this
  rule to generic category words, DI or DF records, DNs, section names, or task
  identifiers. When an artifact moves, a link-target-only correction is
  permitted in a closed day, but its historical wording must remain unchanged.
  Source: DI-punuj
- Maintain exactly one accurate task-ledger list under each dated day. While a
  day is `Active`, update its existing bullets in place as work completes, moves,
  or changes estimate; do not freeze a stale list and append a replacement list.
  Every originally promised item remains in that one list and is classified as
  complete, partial, blocked, not started and rescheduled, superseded, or
  explicitly deferred. Source: DI-bugob
- Put one current workload summary immediately below every dated heading. For an
  `Active` day, keep completed time, remaining time, and current expected total
  accurate after every update. For a `Forecast` day, show its planned total. For
  a `Closed` day, show the final attended total. A stale top summary or a second
  task list for the same day is an error. Source: DI-bugob
- Calendar-wide prose may precede the dated subsections. Timestamped day-wide
  prose may record material plan changes, variance, and lessons before the day's
  single task list, but it must not duplicate task bullets or retain obsolete
  totals. When the day closes, its final accurate summary and single list become
  historical evidence. Preserve previously closed days as written unless a
  separately approved correction is required. Source: DI-baluv; DI-bugob
- Mark elapsed dates `Closed`, the current date `Active`, and future dates
  `Forecast`; no elapsed date may remain `Active`.
- Start every elapsed or current completed-work bullet with wall-clock task time
  rounded to the nearest tenth of a decimal hour. Do not label past time
  `observed`; elapsed placement implies that status. Start every unfinished
  current or future bullet with its estimated wall-clock time followed by the
  recommended GPT-5.6 Sol planning/execution reasoning pair, such as
  `2.5h high/xhigh`. Use only supported reasoning levels: `none`, `low`,
  `medium`, `high`, `xhigh`, and `max`. Source: DI-baluv
- Reconstruct elapsed time from the relevant session log and Git history. A task
  clock starts at the first relevant request or explicit resume, includes model
  reasoning, tool execution, tests, reviews, and waits while that task remains
  current, and stops at completion, an explicit task switch, or an explicit
  pause. Only one task clock may run at a time. For work spanning overnight,
  estimate and subtract sleep; also subtract explicit daytime pauses. Assign
  each interval to at most one bullet so elapsed time is not double-counted.
  Source: DI-baluv
- Base future estimates on comparable elapsed work when available. Apply a
  conservative 1.0x multiplier for a direct analogue, 1.5x for a first
  production or cross-module slice, and 2.0x for Windows, live-system, uncertain
  side-effect, or recovery work. Recommend `medium` for mechanical work, `high`
  for bounded design or implementation, `xhigh` for cross-module or recovery
  work, and `max` planning for irreversible protocol, identity, trust, or live
  acceptance decisions. Source: DI-baluv
- Show a total attended wall-clock workload for each day. No future day may
  exceed 10.0 hours in the single sequential task lane. Validated unattended
  read-only scale tests and dump jobs may continue outside that lane only when
  their runtime is logged separately; review and analysis remain attended work,
  and no QuickBooks, OSC, or other external-system mutation may run unattended.
  If a replan would exceed the attended limit,
  redistribute work and synchronize every affected milestone, dependency,
  checklist gate, and dump attempt; do not merely warn. In the active plan,
  preserve any fixed delivery date that the current Goal Contract And Forecast
  actually selects, along with its attended external-system mutation and named
  defect reserve. When that contract selects no fixed date, preserve every gate
  and report a dependency-based earliest-to-conservative forecast instead of
  inventing a deadline. Do not consume named reserve hours with unrelated
  planning or housekeeping. Source: DI-baluv; DI-samur; DI-duvip; DI-mijof
- The 10.0-hour attended limit applies to the repo owner across all concurrent
  Codex sessions and worktrees, not separately to each session. Prompts,
  decisions, reviews, conflict resolution, and integration are attended work.
  Record overlapping model reasoning, tool execution, and unattended test
  runtime separately; they may shorten the critical path but do not create
  additional human capacity and must not be added to the attended total.
  Source: DI-ramar
- During every `update dobab`, recompute each parallel-activity statistic shown
  in a daily summary from the underlying unrounded durations for that same
  accounting period. First calculate `effective active hours = main-lane active
  hours + sum(worker-lane active hours)`. Then calculate `parallel activity
  percentage = 100 * effective active hours / attended wall-clock hours`. Use
  the day's attended wall-clock total as the denominator; never use the elapsed
  calendar or coordination-window duration, and never use the main lane's active
  subset. Show the operands and exact calculated percentage in the summary, then
  display the percentage rounded to the nearest tenth of a percentage point.
  Keep effective active hours separate from attended hours rather than adding
  them together. If attended wall-clock hours are zero, report the percentage as
  not applicable instead of dividing by zero. Source: DI-susim; DI-kakuh
- After each `update dobab`, total the remaining estimated wall-clock times for
  the current day. In the chat response, reproduce each remaining current-day
  calendar bullet as written, including its leading hour estimate,
  planning/execution reasoning pair, `dobab.*` owner, and human-readable
  description; do not reduce the list to task identifiers. Then state the
  combined remaining time. Obtain the current `America/Los_Angeles` local time
  after completing the update, add the combined remaining wall-clock time, and
  report both that current time and the resulting estimated completion date and
  time. State that this arithmetic assumes uninterrupted sequential work and
  excludes unplanned breaks, sleep, and parallel execution. When the active Goal
  Contract And Forecast contains a fixed delivery date, also report remaining
  schedule slack in hours and percent through that date, inclusive, using the
  existing capacity rules. When it contains no fixed date, do not report slack
  or a completion percentage. Instead report remaining critical-path hours, the
  earliest-to-conservative forecast range, each named defect/retry/canary
  reserve, the assumptions that distinguish the bounds, and the concrete risks
  that can move them. Source: DI-baluv; DI-tatot; DI-muhal; DI-mijof
- For every item promised by an elapsed or current forecast, classify it as
  complete, partially complete, blocked, not started, superseded, or explicitly
  deferred. Record actual evidence, variance, and lessons. A design document,
  TE, or partial prototype does not complete a production-code milestone, and
  work without the required tests and focused commit remains unfinished.
- When the date has rolled over, map every unfinished deadline item from the
  elapsed date to an explicit current or future date. Preserve the missed item in
  the elapsed narrative, but also assign its remaining work to a dated
  destination with its `dobab.*` owner, dependency, and completion gate. Work
  explicitly removed from deadline scope must instead name its post-deadline
  owner or decision record. A vague statement such as "carry forward" or a
  past-only `Remaining gate` does not count as redistribution.
- If a future forecast already contains the unfinished work, reconcile and cite
  that destination instead of duplicating it. Otherwise insert it into the
  earliest realistic date permitted by its dependencies. Recalculate affected
  start dates, completion dates, dump attempts, integration gates, and downstream
  milestones rather than stacking displaced work onto an impossible day.
- Keep the component milestone table, direct dependency matrix, Daily Calendar,
  and checklist synchronized. Status, due dates, and task wording must describe
  the same execution order. Do not silently drop work from one representation or
  leave a completed item open without explaining the remaining scope.
- Maintain the full-calendar Graphviz documentation with every applicable dobab
  change. When a retained Daily Calendar item's text, date, status, worktree,
  subtask ownership, or hours changes, review
  `docs/dobab-calendar/full-calendar-annotations.json`, update its reviewed
  title or direct prerequisite edges when needed, and run
  `make -C tools/dobab-calendar all` in the same change. Commit the resulting
  `docs/dobab-calendar/graphviz-full-calendar.dot`; do not update the preserved
  `x/dobab-calendar-bakeoff/` snapshot. Source: DI-nunun
- Put any current experiment result that controls scheduled work in literal
  plain English near the top of the active day before discussing handback,
  review, or merge state. Distinguish an arbitrary comparison cutoff from a
  justified operational requirement. Draw a red critical-path blocker only for
  an actual unresolved launch safety failure or missed justified requirement;
  do not draw one merely because every candidate missed an arbitrary comparison
  value. Source: DI-vigah
- Preserve the current Goal Contract And Forecast's attended live-system action,
  named retry or defect reserve, final acceptance, and every prerequisite gate.
  Preserve a final delivery day only when that contract explicitly selects one.
  If it selects no fixed date, maintain dependency-based forecast bounds and
  move them when evidence changes rather than silently consuming gates. Source:
  DI-samur; DI-pupug; DI-duvip; DI-mijof
- Cite concrete evidence already present in the repository, including commit
  hashes, passing checks, target-machine observations, dump identifiers, or
  decision records. Do not claim that the uncommitted dobab update itself is
  already committed.
- Before finalizing, audit that every unfinished past item has exactly one visible
  current, future, or named post-deadline destination, all affected dependency
  and checklist statuses agree, no past date is active, no deadline work
  disappeared, and
  `git diff --check` passes. In the handoff, summarize what completed, what moved
  to which date, what remains blocked, and any change to deadline risk. Source:
  DI-tijab; DI-nobam

## Build, Test, and Development Commands
- `go test ./...` runs the test suite.
- `go -C x/cas-native-index test ./...` runs the isolated CAS-native-index
  experiment suite; root-module wildcard commands do not cross the nested module
  boundary. Source: DI-gudop
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

## Agent Instruction Architecture (Required)
- `AGENTS.md` is the canonical home for repo-wide protocol, workflow, and vocabulary rules. Role-specific files such as `AGENTS-codex.md` and `AGENTS-ppx.md` are overlays: they may add stricter role constraints, environment setup, and role-specific procedures, but they must not duplicate or relax canonical repo-wide rules. (DI-034-20260508-060134)
- When a rule applies to every agent or every repo artifact, move it here and replace role-file copies with pointers. When a rule applies only to one agent's runtime environment, identity, credentials, branch lifecycle, or private logging system, keep it in that role overlay. (DI-034-20260508-060134)

## Public Repository Isolation (Required)
- Keep unsettled public-facing architecture, protocol, specification, TE, DF,
  DI, implementation, fixture, and migration work in this private `cdint-grid`
  repository. During that phase, treat the corresponding public repository as
  read-only reference material. Do not create or modify a public worktree,
  branch, commit, issue, pull request, or remote ref merely to continue private
  planning or worker coordination. Source: DI-gimok
- Private worker handoffs, results, merge declarations, scratch paths, tmux
  names, local usernames, private repository names, business names, schedules,
  and application-specific context must not enter a public repository. Public
  repositories receive only the minimal settled design, specification, code,
  fixtures, and migration material needed by that public project. Source:
  DI-gimok
- After the governing private TE, DF, and DI are complete, prepare the public
  change as a separate sanitized copy. Review the complete public diff and Git
  metadata for private information, run the private public-export checker, and
  obtain explicit user approval before any public remote operation. A passing
  sanitizer does not replace full-diff review. Source: DI-gimok
- If the public repository advances while private work is in progress, re-read
  its accepted state before preparing the sanitized copy. Reconcile the settled
  private result with that current public state; do not copy stale worker
  coordination records or overwrite newer public work. Source: DI-gimok

## Deterministic Context Loading (Required)
- Use the matching repo-local skill under `.agents/skills/` for repeated project
  context work. Those skills route to deterministic commands that validate the
  complete source locally and emit exact bounded Markdown with provenance;
  expand by named term, task, decision, section, or date when needed. Do not
  substitute an LLM summary for source validation. Source: DI-sipog; DI-pugid;
  DI-hasod
- A worker handoff records the current Git commit and `AGENTS.md` blob. When both
  match the already-injected guidance, do not read `AGENTS.md` again. On a
  mismatch, use `token-efficiency agents verify` to emit the complete changed
  sections, or `agents section` for a named section. Source: DI-hasod
- Use `dobab-calendar context` for the default operational dobab view. Request
  older closed history only by exact date or literal match, or when rollover
  reconciliation needs it. Actual model-visible calendar PNG review is required
  only for generator/render behavior changes or an explicitly suspected visual
  defect; deterministic freshness and Graphviz checks remain required. Source:
  DI-hasod; DI-kamam

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

## Architecture Precedent And Role Checks (Required)
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

## Promise Action Minimalism (Required)
- Future PromiseGrid protocol, simulation, POC, scoring, generation, and guide work must not invent workflow-specific top-level action kinds by default. The default future-facing top-level semantic action is `promise`. (DI-mosoj)
- Treat observation as a promise that the promiser observed something from its local vantage. Treat refusal as absence of a promise, a promise not to do something, or a promise that the agent does not currently promise the requested behavior. (DI-mosoj)
- Treat repair, offer, counteroffer, acceptance, routing, introduction, redemption, transfer, storage, computation, TCP-link changes, authorization, dispatch, grant, registration, and enforcement as pCID-defined payload semantics, local trust/evidence interpretation, or implementation-local mechanics unless a scoped TE/DI proves a distinct wire-level role. (DI-mosoj)
- Before adding any new top-level action kind, stop and answer: why is this not just an agent's voluntary promise with pCID-defined payload meaning? If that answer is not explicit in a locked TE/DI, do not add the action kind. (DI-mosoj)
- POC7 through POC10 action names, including refusal and observation labels, are historical executable evidence, not the forward naming pattern. Do not rewrite scored/generated artifacts in place; supersede or reframe them when reused. (DI-mosoj; supersedes DI-fitav)

## POC Superset Discipline (Required)
- Future PromiseGrid POCs must be supersets of the previous POC's implemented behavior, architecture lessons, analyzer gates, and documented acceptance criteria unless a scoped DI explicitly declares the new POC non-superset and lists every intentionally dropped feature. (DI-sinur)
- A new POC may add focus, specialization, or new protocol surfaces, but it must not silently regress app/kernel boundaries, local trust semantics, monitor/analyzer gates, pCID routing, Promise Theory vocabulary, or previously proven workflows. (DI-sinur)
- When repairing or extending a POC, the analyzer must include inherited regression gates for the prior POC lineage, and the README or run narrative must state whether the POC is a superset or cite the DI authorizing an exception. (DI-sinur)

## Decision-First Specification and Compliance Protocol (Required)
- Decision-first means decisions must be locked before coding; it does not forbid pre-decision analysis such as required thought experiments.
- The agent must collect and lock user decisions before making any code edits for a task.
- Locked decisions must be recorded as Decision Intent Log entries in the relevant
  root `TODO/TODO-<handle>-<slug>.md` file(s) with clear intent and rationale.
  Source: DI-gigoz; DI-topih
- The agent must ask exactly one user decision question at a time, including
  exactly one question in each `request_user_input` call, and wait for the
  answer before asking the next question. Gather discoverable facts first and
  ask the remaining questions in dependency order. Source: DI-zukup
- Required decision categories are architecture, design/behavior, implementation approach, function naming, variable naming, and file/path decisions.
- The agent must ask these as multiple-choice questions whenever practical.
- When a thought experiment (TE) is required, the agent must complete the TE before asking final DF questions. TEs narrow alternatives; DF questions and answers lock the decision before implementation.
- Thought experiments (TEs) are analysis artifacts; Decision Intent (DI) entries are the separate records that capture the locked decision after DF is resolved.

## Planning Cadence And Executable Progress (Required)

- Limit each planning or DF round to roughly 20 active minutes. Start the clock
  at the first relevant request or explicit resume. Count reasoning, repository
  inspection, tool execution, questions, and answers while the round remains
  active. Source: DI-zokip
- Exclude only pauses that the session log identifies reliably: time waiting for
  a user answer after a question, time after a completed response until the next
  user request, an explicit pause or abort, or an explicit task switch. Do not
  infer a pause merely from a quiet timestamp interval while reasoning or a tool
  may still be running; count ambiguous intervals as active. Source: DI-zokip
- At roughly 15 active minutes, stop opening new topics and consolidate the
  current decisions, evidence, and unanswered questions. By roughly 20 active
  minutes, end the planning round without rushing a decision. Persist settled
  decisions in DIs, unresolved decisions in DRs, implementation work in TODOs,
  deadline scheduling in dobab, and substantial analysis in the applicable TE,
  DN, contract, or specification. Give each item one authoritative owner and
  link to it elsewhere rather than duplicating mutable descriptions. Source:
  DI-zokip
- Alternate planning with executable progress: after a planning round,
  write a bounded code slice, run a focused test, or execute a focused
  experiment before beginning another planning round. Follow "write a
  little, test a little" (WALTAL) and avoid hours of discussion that
  produce no executable evidence. Persist all outstanding decisions
  and discussion in TODO or DR files instead of extending the planning
  round. Source: DI-zokip
- Prefer a focused test or experiment when it can answer a discoverable question
  faster or more precisely than further discussion. Tests and experiments do not
  bypass TE, DF, DI, path-approval, safety, or side-effect requirements. Source:
  DI-zokip
- A later `continue` or explicit resume begins another bounded planning round
  from the persisted state. This planning clock is separate from dobab's attended
  wall-clock accounting, although both use the same session-log evidence.
  Source: DI-zokip

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

## Parallel Worktree Coordination (Required)

### Active Git-Native Peer Workflow

- Main and every linked-worktree worker use the same `commit` and `consensus`
  skills. Main is special only because it is the sole writer of the official
  `main` branch. No peer's acceptance, decline, deferment, queue position, or
  trust estimate is global. Source: DI-pufuk; DI-bisaf
- A new-style worker's stable identity is its suffix-free branch name. Its
  worktree basename, inbox identity, tmux identity, and worker-inventory identity
  match that branch. Worktree-local
  `author.name` and `author.email` set its Git Author to
  `<branch> <branch@local.invalid>`; `user.name` and `user.email` retain the
  normal repository identity as Committer. Do not roll branch/worktree names or
  create attempt IDs for tasks, retries, model changes, Codex restarts, or later
  promises. Do not retain an `-aNN` attempt suffix in a worker identity. Source:
  DI-goziz; DI-pufuk; DI-bibur
- A TODO records work and a DR records an unresolved decision. Before substantial
  work, the worker voluntarily creates one starting-promise commit naming the
  accepted scope, paths, barriers or their absence, and stop conditions or their
  absence. Every new-style commit, including merge and empty commits, follows
  the repo-local `commit` skill and contains exactly one nonempty plain-English
  `Promise:` trailer. Published promise commits are corrected only by successor
  commits. Source: DI-zofug; DI-vomid
- Any peer may read another peer's committed branch, working tree, and available
  session log. Dirty state and session dialogue are provisional observations.
  A worker may translate, reimplement, disagree with, ignore, or merge peer work,
  but writes only its own branch and working tree. Path scope describes what that
  worker promises to change; it is not a global path reservation. Material
  translated committed work uses adjacent full `Source-commit:` and matching
  `Co-authored-by:` trailers. Source: DI-kijuz; DI-bavon; DI-vomid
- Use the repo-local `consensus` skill before depending on peer work. For a
  complete peer-tip import, inspect and disposition the entire selected tip,
  name both exact commits, choose checks from the governing TODO/DR, established
  repository gates, and peer evidence, then use
  `git merge --no-ff --no-commit <exact-peer-commit>`. Resolve conflicts and fix
  failed checks on the receiving branch, rerun the checks, and commit the exact
  tested tree with the `commit` skill. Stop only for an unresolved decision,
  unapproved path or side effect, changed governing scope, or another explicit
  stop condition. Source: DI-bisaf
- Warn, but do not prohibit, when the peer tip is already contained, the merge
  adds no tree change, only reciprocal merge commits are new, or peers repeatedly
  exchange merges without substantive work. An integration worker is an
  ordinary peer that may combine candidates and publish correction successors.
  Main may publish that candidate, another candidate, or none. Source: DI-bisaf
- Every material `kept`, `broken`, or `inconclusive` judgment is the evaluating
  worker's own promise commit about the exact named promise. Concrete repair
  work belongs in a TODO; an unresolved decision may be requested from any named
  agent in a DR. Each worker's current candidate and trust view is disposable,
  local state beneath
  `/tmp/cdint-grid/consensus/<worker-id>/peers/<peer-id>.md`; it must name its
  exact Git frontier and be rebuildable from Git. There is no global trust score.
  Source: DI-dabuv; DI-pufuk; DI-bisaf
- `newtree` uses the repo-local `newtree` skill. It creates or reuses one stable
  worker, enables worktree-local Git configuration, starts Codex in Plan mode,
  and sends a short instruction naming the governing TODO or DR. Runtime launcher
  and prompt files live beneath `/tmp/cdint-grid/workers/<worker-id>/**` and are
  disposable. New-style work creates no `handoff.json`, `handoff.md`,
  `results.md`, `merge-checks.json`, numbered attempt directory, or FIFO entry.
  Source: DI-pufuk
- The merge queue remains available for old-style attempts and optional
  mechanical exact-tree checking. It is not required for new peer work, FIFO
  order does not control peer integration, and queue admission does not create
  acceptance. Separate independent review remains exceptional and requires the
  repo owner's per-case approval under DI-gafom. Source: DI-bisaf
- The activation commit is the first new-style promise on `main`. Validate the
  skills, inventory, and workflow immediately after it; repair failures with
  successor commits before new product work. Before an old-style worker
  transitions, close any live legacy queue entry through its retained-evidence
  path. Preserve its old scope, paths, side-effect permissions, barriers, stop
  conditions, results, and evidence in the governing records and transition
  promise. Reconcile an active attempt-suffixed branch or mismatched worktree
  path at its next safe restart before later peer work. Preserve historical
  identifiers and local-CAS evidence rather than rewriting them. Source:
  DI-dijul; DI-bibur

### Worker Instruction Bundle Upgrades

- `.agents/worker-instructions.manifest` is main's deterministic coordination
  write set. It contains repository rules, repo-local skills, and the tracked
  tool sources needed to execute them. It deliberately excludes TODOs, design
  notes, live dobab state, reviewed calendar annotations, and generated calendar
  artifacts. The write set does not grant product authority, settle a design,
  broaden a task, or compel acceptance of peer product work. Product and design
  candidates still use `consensus`. Source: DI-kolat; DI-hohuh
- A direct user message consisting exactly of `upgrade` invokes the repo-local
  `upgrade` skill. The skill runs the installed
  `cdint-grid-skills-sync` command. From a worker it selects that stable worker;
  from main it selects every linked non-main worktree. Routine instruction
  distribution does not use the inbox and does not require per-worker model
  reasoning. The unique local `main` worktree, with clean committed
  manifest-covered paths, is the source; do not fetch or substitute a remote
  tip. Source: DI-kolat
- Compare each covered regular file by type, mode, modification time, and
  SHA-256 bytes. Identical content and mode are current regardless of timestamp.
  Replace differing worker content only when it is older than main. Skip that
  worker and continue others when differing worker content is newer or has the
  same timestamp, a worker-only entry exists inside a covered directory, a
  covered path is dirty, branch and worktree identities disagree, the worker is
  detached, or a pre-write race is observed. Never delete or force those paths.
  Source: DI-kolat
- Run `rsync -av` without `--delete`. Repeat source and worker snapshots around
  the write. An rsync failure, post-write race, verification failure, or Git
  reference-update failure stops the entire pass because a worker may require
  repair. The synchronizer's lock excludes only another synchronizer; it cannot
  freeze an independently active editor. Source: DI-kolat
- After copying, create one worker-authored promise commit from a private
  temporary Git index. Name the exact main commit, manifest blob, changed paths,
  checks, and limitations; use adjacent `Source-commit` and `Co-authored-by`
  trailers and exactly one `Promise:` trailer. Advance the worker branch with
  compare-and-swap and refresh only manifest-covered paths in its normal index.
  Preserve every unrelated staged and unstaged change. The mechanical promise
  does not claim product or design acceptance. Source: DI-kolat
- Exit `0` means all selected workers were already current or synchronized.
  Exit `2` means at least one worker was safely skipped; investigate each plain
  reason before retrying. Exit `1` means the pass stopped on an operational
  failure. `--dry-run` or `-n` performs no mutation, including no report, lock,
  commit, or file copy. Private run evidence and locks stay beneath
  `/tmp/cdint-grid/skills-sync/`. Source: DI-kolat
- Before another bounded work cycle, `resume` invokes `upgrade`, then rereads
  the resulting local instructions and drains one finite inbox snapshot. If
  synchronization skips that worker or fails, notify main once and stop before
  inbox intake or product work. Another worker skipped by an all-worker pass
  does not invalidate a worker that synchronized successfully. Source:
  DI-kolat; DI-kumof
- Main builds, verifies, and atomically installs shared coordination executables
  beneath `/home/stevegt/bin`; workers consume them without rebuilding private
  copies. The synchronizer copies tracked source only: it never builds, installs,
  or removes binaries. After a successful worker sync, separate cleanup may
  remove only known ignored duplicate executables from that successful
  worktree. A skipped worker retains its existing binaries until investigated.
  Shared CLI changes must remain backward-compatible while skipped workers may
  still run earlier instructions. Source: DI-kolat
- A main promise commit that changes a path covered by the parent or resulting
  manifest verifies required shared executables and immediately runs the
  all-worker synchronizer. A peer promise changing product or design work is
  offered to named peers through ordinary `consensus`; do not reinterpret
  deterministic instruction distribution as global agreement. Source:
  DI-kolat; DI-bisaf

### Legacy Attempt Compatibility

- Every agent that receives an old-style handoff or result file must read the
  entire file
  before acting on it. Review every section, requirement, finding,
  recommendation, limitation, deferred item, runtime path, check, and stated
  follow-up; do not review only the summary, disposition, or changed-file list.
  A handoff is complete only when every requirement has an explicit outcome. A
  result is fully integrated only when every substantive item has an explicit
  disposition tied to durable evidence: implemented and merged, already
  satisfied by a named commit or record, assigned to a durable TODO/DR with an
  owner and completion gate, or explicitly rejected, superseded, or deferred
  with rationale and authority. Retaining a report under `/tmp`, merging its
  code while ignoring its findings, or recording only its headline does not
  count as integration. Mutually exclusive alternatives and intentionally
  deferred work must be dispositioned rather than blindly implemented. A worker
  must map every handoff requirement to its outcome in `results.md` and must
  report every result finding or follow-up that it cannot integrate within its
  authority; it must not silently leave those items for main to rediscover. A
  receiving agent must produce or update an integration matrix that maps every
  substantive item to its disposition and verifies named commits and repository
  paths against its assigned base; main additionally verifies final integration
  against current main. This rule also applies when a handoff or result path
  arrives through the main work queue. Source: DI-tolib

- For a still-active old-style attempt, treat its original `newtree` request as
  an instruction to keep the main
  `cdint-grid` Codex session available for coordination. The main session does
  only enough planning to establish the exact branch point, worktree and branch
  names, disjoint path ownership, protected shared paths, known shared
  constraints, handoff attempt, and collision checks. It then creates or verifies
  the linked worktree and runs the generated launcher. Before launch, main
  commits the strict `handoff.json` machine assignment and the linked human
  `handoff.md` instructions in the worker's TODO attempt directory. The launcher
  starts the tmux and Codex session in Plan mode, sends the absolute `handoff.md` path, and
  the main session verifies that the worker begins reading it. Source: DI-gikut;
  DI-zupol; DI-zinir
- Such an old-style `newtree` worker performs the remaining task-specific exploration, TE and DF
  work, path and naming approvals, DI recording, implementation, and testing in
  its own Codex session. The bootstrap handoff carries decisions already made
  by the user and identifies unresolved categories; it is not required to
  contain a complete implementation plan. The worker must still satisfy every
  decision-first, WALTAL, path, comment, test, commit, and result-report rule
  before editing code or data. Source: DI-gikut; DI-zupol
- For old-style packet assignments, the main checkout remains the coordination owner for `AGENTS.md`, the active
  dobab deadline ledger, shared TODO/DR/DI records, common-interface integration,
  and final acceptance. Secondary linked worktrees must not edit those files.
  A `newtree` handoff may assign one exact, narrow hunk in an otherwise
  main-owned file when the hunk cannot overlap ongoing coordination edits;
  every other hunk in that file remains read-only. Source: DI-ramar; DI-gikut;
  DI-zupol
- Give each old-style secondary worktree one disjoint path ownership set and one focused
  branch. A worktree may read shared records, but it commits only its assigned
  files. If its implementation requires a shared contract or shared path to
  change, stop that worktree and return the decision to the main checkout before
  continuing. Source: DI-ramar
- Freeze durable wire and storage protocols, external side-effect rules,
  security rules, and shared runtime paths before dependent parallel coding.
  When a cross-package Go interface remains uncertain, compare small working
  examples before locking it. Package-private internals may evolve independently.
  Do not build complete competing implementations unless an explicit bakeoff
  requires them. Source: DI-ramar
- When an old-style worker actually needs a shared decision or contract added after its
  branch point, main may direct that clean worker to merge one exact main commit,
  review the merge, resolve only assigned-path conflicts, and rerun affected
  checks. Do not merge main into every worker merely because main advanced, and
  do not require a final worker-owned synchronization before completion. Source:
  DI-gipob; supersedes the routine synchronization requirements of DI-lajig and
  DI-huhon
- Before reporting completion, an old-style non-dobab worker commits and self-reviews its
  focused work, runs its checks, records its exact worker base commit, merge
  base commit, and focused work commits, supplies an ordered
  `merge-checks.json`, and commits `results.md` in
  its TODO attempt directory. `results.md` must not claim the hash of the commit
  that contains itself. Git and the merge queue bind the tracked record bytes to
  the submitted immutable worker tip. Main movement after that point does not
  change that submitted worker tip.
  The standing dobab worker follows the direct-integration exception below and
  does not need `merge-checks.json`. Source: DI-gipob; DI-bijaj; DI-jifus;
  DI-barir
- Except for standing dobab updates governed by DI-bijaj, integrate old-style workers one
  at a time through the local FIFO merge queue. Main adds
  the self-reviewed immutable merge base and tip plus the tracked TODO attempt
  directory in critical-path order. The queue also verifies the worker base
  recorded for the assignment and copies the exact tracked result
  and check bytes from that tip into its private run evidence. `cdint-grid-merge-queue run`
  tests the merge
  against current main in a fresh detached worktree and records the exact tested
  Git tree. Main performs bounded final acceptance: it reads and dispositions
  the complete result, verifies the exact commits, declared checks, queue
  evidence, tested and staged trees, and unresolved risks, but does not repeat a
  full code review. Main then uses `cdint-grid-merge-queue stage`; the tool refuses if main
  moved or the staged tree differs. Main creates the accepted non-fast-forward
  merge commit outside the tool, and
  `cdint-grid-merge-queue finish` verifies its parents and tree before running after-merge
  checks. If main moved, rerun the queue attempt against current main without
  returning the unchanged worker tip for synchronization. A conflict, failed
  check, or rejected review preserves the evidence and returns only actual
  worker defects to the same worker through a fresh numbered attempt. Do not
  resolve worker conflicts or repair rejected work in main, and do not rebase,
  cherry-pick, create an integration-candidate branch, or fast-forward main to a
  worker. Remove a linked worktree only when it is clean and every required
  commit is integrated. Preserve worker branches and commits as evidence.
  A rejected, blocked, crashed, or review-only attempt uses the queue's
  record-only mode to stage and verify only its TODO attempt directory; rejected
  worker code and worker-code checks must never enter or execute during that
  retention pass. Rejecting an ordinary merge candidate must
  preserve an as-yet-untracked TODO attempt directory through that mode before
  the queue entry closes. Source: DI-gipob; DI-ramar; DI-mikil; DI-jifus;
  DI-barir
- Separate independent review is not part of the normal worker path. Main may
  recommend one only for a credible risk of irreversible data loss, secret or
  privacy exposure, incorrect authorization or identity, failed recovery, or a
  first live external-system mutation. Before assigning a reviewer, main states
  the exact candidate, concrete harm, review scope, reviewer, and expected delay
  and obtains the repo owner's explicit permission for that case. An approved
  review gates only that candidate and may cover corrections caused by its
  findings while task, risk, paths, and reviewer remain unchanged; a scope or
  risk change requires new permission. Unrelated work continues. Refusal of an
  extra reviewer does not waive any correctness or safety gate, and main may
  still reject insufficient evidence. Source: DI-gafom
- When receiving an experiment result, inspect and report its literal measured
  outcome before spending time on handback mechanics or merge eligibility. A
  packet-only defect may receive one bounded correction by the same worker while
  unrelated diagnosis continues. If the corrected packet fails again for a
  packet-only reason, stop the correction loop and identify the missing
  deterministic preflight check before authorizing another attempt. Source:
  DI-vigah
- Retire an eligible old-style ordinary worker in process-first order: close its exact
  tmux/Codex session, verify that session is absent, and only then remove its
  clean linked worktree. If the session cannot be stopped or its absence cannot
  be verified, leave the worktree in place. Before stopping it, read its entire
  latest `results.md` and every associated result artifact, then verify that
  every substantive result item is fully integrated under the item-level
  dispositions required above or is explicitly rejected with durable rationale.
  Also verify that required accounting is recorded, no assignment or merge
  remains active, and its Git state is clean. A summary, clean worktree,
  integrated commit, or passing check does not satisfy this result-integration
  gate by itself. Preserve the worker branch, commits, and all retained handoff,
  result, review, queue, and run evidence. Do not apply this ordinary-worker
  cleanup rule to the dobab or overall-review standing workers. Source:
  DI-tobih; DI-sohah
- Give every old-style worker attempt exactly one owning TODO before launch. Main creates
  that TODO on `main` before the worker base commit is selected; the owning TODO
  must
  remain a regular tracked file on the exact main commit tested by the queue.
  Main then creates
  `TODO/TODO-<handle>-<slug>.d/<full-attempt-id>/handoff.json` and `handoff.md`
  on the worker branch and commits both as that branch's first attempt-specific
  commit before Codex starts. `handoff.json` is the machine authority for the
  assignment; `handoff.md` must contain exactly
  `Assignment metadata: [handoff.json](handoff.json)` near its top. The worker
  commits `results.md` and the applicable
  `merge-checks.json` in that same TODO attempt directory before ending,
  including when partial or blocked. Use
  `/tmp/cdint-grid-handoffs/<full-attempt-id>/` only as the attempt scratch
  directory for `start.sh`, logs, caches, and other reproducible or untracked
  files, as specified by
  [DN-lovug](docs/DN-lovug-parallel-codex-session-handoff-files.md).
  `handoff.json` and tracked human handoffs are each limited to 64 KiB; tracked results and merge
  declarations are limited to 1 MiB each. Before staging, reject credentials,
  private keys, unsafe Git modes or links, path escapes, malformed records, and
  content prohibited by the owning handoff. Human review remains mandatory.
  Refuse dirty worktrees, the wrong branch, an unavailable or non-exact base
  commit, existing results, and, for a new attempt, existing worktree-derived
  tmux sessions. Do not automatically delete or overwrite an attempt. A retained
  same-attempt retry follows the stricter verification rule below. Source:
  DI-baduf; DI-hahoh; DI-zupol; DI-jifus; DI-barir

### Shared Worker Inspection And Runtime Rules

- The main session may check workers whenever current worker state is useful;
  separate repo-owner authorization is not required for each check. Each
  `check workers` operation means exactly one invocation of the verified
  `/home/stevegt/bin/cdint-grid-workers check`, which discovers live workers
  from Git and tmux, atomically writes the generated private current inventory at
  `/tmp/cdint-grid/worker-inventory/WORKERS.json`, and gathers one bounded
  private snapshot. The inventory is output and need not exist before the
  invocation; it does not mean continued monitoring. The generated inventory and
  retained snapshots use schema version 3 and keep `todo_attempt_dir` separate
  from `attempt_scratch_dir`. The reader accepts versions 1 and 2, but emits
  only version 3. Version 3 also reports each worker's first-parent promise
  history after its transition commit. It searches tracked
  `handoff.json` records first and hashes the exact `handoff.json`, `handoff.md`,
  `results.md`, and applicable `merge-checks.json` without retaining their
  bodies. When a reused worker has multiple valid tracked attempts, its exact
  current Codex rollout may select one of those attempts through either a tracked
  TODO handoff path or a legacy scratch handoff path. If that evidence is
  unavailable or does not select an attempt, exactly one unfinished tracked
  attempt wins; multiple unfinished attempts remain ambiguous. Unresolved
  malformed tracked metadata is always an error and prevents fallback. Historical
  committed handoffs whose owner handle no longer matches a moved TODO path are
  read compatibly without rewriting their bytes; new handoff creation retains
  strict owner matching. Merge-queue evidence accepts snapshot schema versions 1
  and 2 and treats `entries: null` as an empty queue, while unknown versions and
  duplicate or ambiguous entries remain errors. Source: DI-muloj
  The helper uses deterministic
  local code, makes no LLM, API, or network calls, and never sends input to,
  focuses, raises, creates, kills, resizes, or otherwise changes a worker
  session. In particular, it never invokes `tmux send-keys`. If the helper
  cannot complete, stop other coordination and report its exact failure
  immediately. Until the helper works again, replace each needed check with one
  bounded, read-only manual inventory using tmux current-visible-pane captures
  and Git worktree inspection. Do not send input to or otherwise mutate workers
  as part of that fallback. After reading that
  one report, main may use its own LLM, networking, guarded `send-keys`, or other
  coordination actions when the report and all existing pane-capture, approval,
  path, merge, and safety rules warrant them. After assigning work, continue
  independent coordination work when available, then return control. Main may
  perform another finite check whenever needed, including after other work or a
  known state change. Never implement checking as a polling loop, fixed or
  repeated sleeps, `tail`, `watch`, filesystem watchers, background monitoring,
  or any other process that monopolizes the terminal while waiting for progress,
  readiness, or completion.
  After the invocation, list every discovered worker exactly once, grouped by
  current pane activity. Show each worker's Git containment, result, queue, and
  notable warning state without conflating those states. Then read the current
  dobab calendar, name the deadline-critical tasks and their assigned workers,
  and discuss whether the critical path is active, idle, blocked, or unstaffed.
  State when dobab is stale or has no active day. Source: DI-jizuj
  During that one invocation, notification-log inspection is limited to the
  final 65,536 bytes of `/tmp/codex-notify-2030.log`. When the tail begins after
  byte zero, discard the first possibly partial line and parse only complete
  remaining records. Fix the read extent to the size observed at open, allow
  append-only growth after that point, and reject replacement or truncation.
  Older notifications are absent evidence, not mismatches or errors. Git and
  tmux remain authoritative, generated inventory retains no notification prose,
  and this bound does not create a continuous monitor or authorize an unbounded
  log read. Source: DI-vahon; DI-gavoh
  Reading a specific result path supplied by the repo owner is not polling. The
  required pane capture immediately before user-directed worker input remains
  mandatory, but it does not authorize readiness monitoring before that action
  is requested. Source: DI-ludag; DI-gopoh; DI-gavoh
- Before sending any ad hoc input to a worker, run `tmux capture-pane` for that
  exact pane, inspect the current prompt, mode, work, and approval state, and
  quote the complete visible capture in chat. Do not send a new task, steering
  message, or other prompt while the pane says `Working`, `Running`, `Waiting`,
  or `Waiting for background terminal`; an input prompt or suggestion visible
  during that state does not make the worker idle. Wait for the current turn to
  visibly complete, capture the idle pane again, and only then use
  `tools/codex-send-prompt <tmux-pane-target> <literal-prompt>`. Source:
  DI-hafiv; DI-tabab
- Never send a multiline instruction or content that the Codex TUI would render
  as `[Pasted Content ...]` as literal worker input. Use the tracked
  `handoff.md` when it is the assignment authority. For other approved steering,
  write the complete text to a mode-`0600` file under
  `/tmp/cdint-grid-handoffs/<attempt-id>/prompts/<UTC-timestamp>-<slug>.md` for
  an old-style attempt or
  `/tmp/cdint-grid/workers/<worker-id>/prompts/<UTC-timestamp>-<slug>.md` for a
  new-style worker. Use only the applicable approved scratch pattern.
  Send one short line naming the absolute file, for example: read the named
  absolute path and continue the assigned work. Keep short one-line commands and exact
  approval or refusal responses inline. Prompt files contain no credentials or
  secrets, remain through attempt integration, and follow the attempt's normal
  cleanup or archive decision. Source: DI-zuzov
- A visible command-approval dialog is the narrow exception to the idle-worker
  rule. After capturing and quoting that complete screen, use the helper to send
  only the exact approval or refusal response displayed by the dialog. The
  helper sends literal text, silently captures the pane until the display has
  changed and stabilized, and sends `Enter` exactly once. It then requires a
  visible response before succeeding. Its internal captures are timing evidence
  and are not printed; they do not replace the caller's complete visible capture.
  A readiness or response timeout exits nonzero and must not automatically
  resend text or `Enter`. Do not call direct `tmux send-keys` commands for worker
  input. If the helper is interrupted, aborted, or times out, capture and inspect
  the pane again for partially entered or submitted text before retrying. A
  generated launcher satisfies the pre-send inspection requirement through its
  bounded readiness, idle, or mode check immediately before invoking the helper.
  Source: DI-hafiv; DI-tabab; DI-tirap
- Immediately after every `tools/codex-send-prompt` invocation, run an explicit
  visible `tmux capture-pane` for that exact worker, reproduce and inspect the
  complete visible pane, and verify that the instruction appears as submitted
  transcript and that the worker has begun reading it or doing the work. The
  helper's internal captures, a changed pane, or text still visible in the
  composer do not satisfy this check. If intake is not visible, stop and repair
  that worker before calling it active; do not resend blindly. This is one
  mandatory post-send verification, not permission to poll. Source: DI-zuzov
- A generated `start.sh` creates the worktree-derived tmux session separately
  from opening its visible terminal. When the session is absent, the script
  changes to the assigned worktree and runs
  `CODEX_TMUX_CODEX_COMMAND='codex -c model_reasoning_effort=high -c
  plan_mode_reasoning_effort=xhigh' codex-tmux --create-only`.
  `codex-tmux` creates the usual Codex, `nv`, and shell windows; the environment
  variable changes only the command started in the Codex window. It sets `high`
  reasoning for Default mode and `xhigh` reasoning for Plan mode without opening
  the model menu. This applies to ordinary, standing, source, correction, and
  independent-review workers; reviews do not use a separate `max` requirement.
  Before sending any input, it uses a bounded pane check to require the
  initialized Codex TUI's visible input prompt and status line; tmux-session
  existence, pane-path match, process existence, a shell prompt, or a fixed sleep
  is not TUI readiness. It
  then invokes `tools/codex-send-prompt` exactly once for literal `/plan`,
  verifies that the pane displays `Plan mode`, invokes
  `tools/codex-send-prompt` once for the legacy handoff instruction or the
  new-style short TODO/DR instruction plus `Enter`, and
  verifies intake. The helper must keep
  literal text and `Enter` as separate tmux operations, with bounded
  changed-and-stable pane evidence that the Codex composer processed the text.
  Never use `C-m` to
  submit a Codex TUI prompt; its
  existing use by `codex-tmux` to launch shell commands is outside this rule. A
  retry of the same retained attempt may reuse that session only after verifying
  the exact clean worktree, branch, pinned commit, session name, pane path, and
  prior prompt intake. It does not resend accepted prompts, and any mismatch or
  ambiguous partial startup stops the launcher. Existing active workers adopt
  these reasoning levels only at their next clean idle assignment boundary.
  Source: DI-sifis; DI-zinir; DI-talif; DI-hafiv; DI-tirap; DI-vumir
- After the prompt checks, `start.sh` runs
  `codex-tmux --attach-in-new-terminal` exactly once from the worker worktree.
  The main session then verifies that the expected worker session is attached
  and that the worker began reading the governing instruction. Do not combine an explicit mode
  with legacy `-t`. The launch delay may reduce polling but does not prove
  initialization or replace the required TUI-readiness, session, mode, prompt,
  pane, or attachment checks. Source: DI-sifis; DI-zinir; DI-talif
- Starting in Plan mode does not replace the decision-first protocol. The fresh
  worker performs task-specific exploration and any required neutral TE, then
  reviews the surviving alternatives and proposed final DF questions against
  Ninik before asking the first question. After the user answers the questions,
  it records the decision lock and reviews the decision-driven changes in the
  final implementation plan before implementation. It also obtains path and
  naming approvals required by its governing TODO/DR or legacy handoff. The Plan-mode response records both
  reviews when DF was needed; user approval of the final plan releases code
  generation. A new-style peer publishes promise commits and uses `consensus`;
  an old-style worker still returns `results.md` for compatibility. Source:
  DI-baduf; DI-gikut; DI-zupol; DI-zinir; DI-gafom; DI-latag
- Old-style tracked `results.md` files do not repeat the legacy fixed `Worker
  Identity` section. Their assignment identity comes from adjacent
  `handoff.json`, while the result records outcomes, focused commits, checks,
  evidence, and compliance. Existing `/tmp` attempts without `handoff.json`
  retain the fixed `Worker Identity` section and exact legacy matching rules.
  Source: DI-zopaz
- `newtree` creates or reuses a stable branch-based worker for a governing TODO
  or DR. It uses a fresh Codex process when starting or restarting, but no fresh
  attempt ID or packet. Resume an existing conversation only through an explicit
  instruction naming that worker. Do not treat `fork` as an alias. Source:
  DI-zupol; DI-libot; DI-pufuk
- Always return related implementation, correction, review-response, and
  integration work to the existing worker that owns the work. Preserve that
  worker's stable branch, worktree, task context, and accepted scope. Do not
  launch a replacement implementation worker merely to get a clean context or
  a newer main branch; peers can inspect and merge exact immutable commits. Use a different implementation
  worker only when the original worker is unavailable, cannot continue safely,
  has incompatible path ownership or unrelated changes that cannot be isolated,
  or the user explicitly directs otherwise; record the exact reason in the
  governing TODO or promise. A separately approved independent review still belongs to a different
  worker. Source: DI-botuz; DI-gafom
- For old-style compatibility, identify an attempt by its full attempt
  ID or by its owning TODO together with its complete attempt-directory name.
  Never use a bare short suffix such as `A01` or `A13`; retain that suffix only
  inside the complete attempt identity. Historical text need not be rewritten
  solely to expand an old suffix. Source: DI-donut
- When reusing an existing worker for a fresh assignment, choose an eligible
  clean worktree whose directory name describes the work being assigned. Do not
  reuse a technically available but unrelated worktree merely to avoid creating
  or restarting a better-matched worker. If no related existing worktree can
  continue safely, create a descriptively named `newtree` or obtain an explicit
  user decision for another choice. Source: DI-torut
- Every new-style starting plan and every old-style handoff, Plan-mode response,
  and result must address the Ninik review. For work that needs DF, the
  worker completes any required neutral TE and then reviews the surviving
  alternatives and proposed final DF questions before asking the first question.
  Required corrections stop the question sequence until corrected and reviewed;
  recommended corrections are incorporated or explicitly dispositioned and
  remain visible in the question context. If an earlier answer materially
  changes a later question or its alternatives, the worker reviews that affected
  question again before asking it. After the user answers and the DIs record the
  decisions, the worker reviews only the resulting changes in the final
  implementation plan unless the selected design materially differs from the
  reviewed alternatives. Work that needs no DF runs only the final pre-code
  review. Each review states which touchstone sections apply and classifies the
  reviewed material as having no design effect, aligning with the touchstone,
  making an explicitly approved departure, or requiring a user decision for an
  unresolved departure. User approval of the final plan releases implementation.
  The new-style promise history, or old-style result, records the review
  chronology, approval, and any implementation deviation; it does not perform a
  new post-code design review. Each receiving peer verifies that chronology
  before depending on the candidate.
  Until the byte-identical design note
  is integrated, use the exact A02 source at
  `/tmp/cdint-grid-handoffs/ninik-promisegrid-idiom-review-a02/results.md` with
  SHA-256
  `94186fd3bddae05e4d0b961c81ad23c6cf2941db8599caf62345d424fdb9105e`;
  afterward use the normal tracked note path. The touchstone is advisory and
  does not override user decisions, DIs, protocols, or tested requirements. The
  Ninik source, its byte-identical design-note promotion, and the
  `review-against-ninik` skill do not review themselves; verify those artifacts
  through source fidelity, exact hashes, functional tests, and bounded main
  acceptance unless the repo owner approves a separate review for a named
  concrete-harm case. Apply `Ninik Review` to the non-Ninik design portions of
  mixed work. Source: DI-zijob; DI-nutan; DI-bihud; DI-gafom; DI-latag
- A Ninik Review reports one bullet for each applicable existing mechanism. Each
  mechanism is the top-level bullet. Under it, render separate nested bullets in
  this order for `Current treatment`, `Severity`, `Plan location`, `Proposed
  edit`, `Disposition`, and `Reason`; do not combine those fields into one
  wrapped paragraph. An already-correct plan uses proposed edit `none` and
  disposition `already satisfied`. If no mechanism applies, include one explicit
  no-applicable-mechanism bullet with the same nested fields. Do not use
  standalone `Simpler existing mechanism` or `Plan changes` fields. Preserve the
  existing classifications, required/recommended dispositions, and TE-before-DF
  chronology. Source: DI-zaniv; DI-pataf
- A coordination barrier is a trusted-worker dependency stop, not security
  enforcement or a PromiseGrid wire protocol. A governing TODO/DR or legacy
  handoff either states that it
  has no coordination barrier or names the work allowed before stopping, the
  exact point the worker must not cross, the prerequisite and acceptable
  evidence, one named confirmer, local verification, and what happens if the
  prerequisite fails or changes. When `main` is the named confirmer or must
  perform the action that releases the stop, the worker sends one bounded
  worker-inbox request to `main`, records its CID, and also asks once locally.
  The worker does not poll or resend an unchanged request. Its result records
  the ordinary chronology of stopping, notification, confirmation,
  verification, and resumption without claiming proof that the worker could not
  have crossed early. The same unfinished attempt may resume only while its
  scope, approved paths, worker base commit, prerequisite, and named confirmer
  remain unchanged. A material change requires a new promise and, when a
  decision is unresolved, a DR. An old-style completed handback still requires
  a fresh legacy attempt for later packet work. Do not add a semaphore, handoff
  schema field, validator, daemon, or new wire protocol until an observed
  coordination defect and a later DF/DI justify it. Source: DI-gumiv; DI-rufim;
  DI-barir; DI-fobit
- Repo-local skills operationalize these worker rules but do not replace
  `AGENTS.md`, TODO and DI authority, path approval, or any separately approved
  review. `commit`, `consensus`, `newtree`, and `resume` are the active peer
  skills; `handoff` and `handback` are compatibility skills. Source: DI-jagol;
  DI-gumiv; DI-gafom; DI-pufuk; DI-bisaf; DI-fivat

- `/home/stevegt/bin/codex-tmux` has three standalone explicit modes.
  `--create-only` creates and initializes a missing worktree-derived session
  without attaching; `--attach-only` attaches or switches the current terminal
  to an existing session; and `--attach-in-new-terminal` opens one GNOME
  Terminal attached to an existing session. Explicit attachment modes never
  create. Pairing any explicit mode with `-t`, combining or duplicating modes,
  duplicate `-t`, positional arguments, unknown options, or help plus another
  argument is a usage error before any tmux command. Preserve no-argument and
  legacy `-t` behavior. `CODEX_TMUX_CODEX_COMMAND` may replace the plain
  `codex` command started in window 0; unset or empty preserves plain `codex`.
  Source: DI-sifis; DI-vumir

## Worker Inboxes (Required)

- A cdint-grid Codex session recognizes `putq ` only when those five characters
  begin a direct user message. Invoke the repo-local `putq` skill, which sends
  the exact nonblank remainder to `main`, reports the returned CID, and does not
  process the deferred body in the same turn. Explicit peer use names a stable
  recipient worker. Queueing grants no authority. Source: DI-zutot; DI-jirob;
  DI-nunan; DI-josov
- A worker that must stop because `main` must make a decision, provide evidence,
  integrate work, or perform another named action invokes the repo-local `putq`
  skill once with recipient `main`. The request names the worker, governing TODO
  or DR, exact requested action, blocking point, and required evidence. Record
  the request CID, ask once in the local session, and remain stopped. Do not
  resend an unchanged request, poll `main`, or infer acceptance from a receipt.
  If sending fails, record the exact failure and remain stopped. Source:
  DI-fobit
- Every cdint-grid session recognizes a direct user message consisting exactly
  of `getq`. Invoke the repo-local `getq` skill exactly once and process at most
  the one returned request under current instructions, decisions, paths,
  barriers, and side-effect authority. Source: DI-jirob; DI-nunan; DI-josov
- Every cdint-grid session recognizes a direct user message consisting exactly
  of `drain-inbox`. Invoke the repo-local `drain-inbox` skill for one frozen
  read-only list snapshot and repeated exact single-CID picks. Map each
  actionable request, completed result, and blocker to exactly one governing
  TODO or DR, finish the complete frozen snapshot, and prioritize the mapped work
  before beginning any of it. Continue past deferred material decisions after
  recording them; stop only when an operational error or ambiguity prevents safe
  classification. Do not chase messages that arrive after the snapshot. Do not
  treat completed or blocked notices as passive evidence, and do not treat TODO
  capture as request completion. End with one list item per frozen inbound
  message showing sender and resulting TODO or DR, followed by categorical
  counts for states, assignments, and unresolved mappings. Source: DI-bidoz;
  DI-vilih; DI-rulak; DI-pidoz; DI-vunuj
- When `consensus` delegates an instruction-bearing candidate scope to
  `upgrade`, a successful upgrade or exact validated no-op must be followed by
  one completed `drain-inbox` pass in the current bounded cycle before consensus
  finishes. Reuse a drain already completed in that cycle instead of taking a
  second snapshot. If upgrade blocks, declines, or leaves partial state, stop
  without draining. A newly required drain only maps and prioritizes its frozen
  snapshot; do not begin unrelated mapped work before finishing the selected
  peer evaluation. Source: DI-javur
- `list`, `list -v`, `status CID`, and `cat CID` are explicit local inspection
  commands. Their output may contain private request bodies; protect redirected
  or captured output. The view is rebuilt from `<worktree>/.grid/cas` and
  pending directional connection files and is not durable authority. Source:
  DI-jirob; DI-nunan
- The `getq` skill owns request-selection and work-state procedure. Ordinary
  `getq` processes unread supported carriers by recipient-side carrier
  modification time, oldest first, before eligible local-ready requests. A
  missing parent is reported but does not make the request ineligible; explicit
  `--lifo` retains newest-message-first access. Carrier time is disposable
  local scheduling input, not event order or durable evidence. Use explicit
  `pick CID` only to choose another exact ready request or opaquely store an
  unknown-pCID message; never execute unknown meaning. Completion records that
  the recipient finished handling the request, not that a TODO was merely
  created or that it accepted every proposal. Source: DI-jirob; DI-nunan;
  DI-josov; DI-fufav; DI-ruhog; DI-sulad; DI-bidoz
- Each worker writes durable inbox evidence only to its own ignored
  `<worktree>/.grid/cas`. Temporary connection files under
  `/tmp/cdint-grid/connections/<recipient>/<sender>/` are transport, not a CAS.
  Do not delete local CAS blocks automatically, use inbox messages to poll or
  control workers, or infer that a receipt means work acceptance. Source:
  DI-jirob; DI-nunan
- The old `/home/stevegt/bin/cdint-grid-main-work-queue` and retained
  `/tmp/cdint-grid/main-work-queue/**` are version-1 migration evidence, not the
  active mechanism after verified cutover. Existing workers finish their current
  atomic action, synchronize their instruction bundle, and adopt the inbox skills
  before another task or later promise; do not interrupt active work merely to
  accelerate adoption. Source: DI-jirob; DI-nunan; DI-josov

## Idle Worker Resume (Required)

- Every cdint-grid session recognizes a direct user message consisting exactly
  of `resume`. To resume an idle worker, invoke the current `resume` skill
  directly. `resume` invokes `upgrade`, which deterministically synchronizes the
  worker against clean committed local `main`; only after that succeeds does the
  worker reread and run its resulting local `resume` skill. A stale worker may
  bootstrap by reading the current skill from local `main`. Do not send routine
  instruction synchronization through the inbox. The resulting `resume` skill completes one finite
  inbox drain, works the highest-priority unblocked items within approved scope,
  uses `consensus` when selected work needs peer evidence or the worker is behind
  `main`, and notifies `main` once about each material completion or blocker. Do
  not poll, repeatedly drain, cross a barrier, use `consensus` as an instruction
  synchronization substitute, or start work after a skipped or failed
  synchronization. Do not add a second drain during consensus in the same
  bounded cycle.
  Starting or resuming an idle peer through `newtree` or mainloop invokes this
  procedure instead of an inbox offer or ad hoc procedure. Source: DI-fivat;
  DI-kumof; DI-javur; DI-kolat

## Standing Worker Coordination And Accounting (Required)

- The only standing Codex workers are the dobab worker and overall-review
  worker. Each keeps one named tmux/Codex session and persistent linked worktree,
  and each new task is governed by a TODO or DR plus a new starting promise.
  Ordinary peers may also persist under their stable branch identity when
  related work continues.
  These are Codex coordination roles, not PromiseGrid application agents,
  administrative principals, or agent identities. Source: DI-libot; TE-judog
- Before assigning a standing worker, verify the expected tmux session, pane
  path, exact Codex thread identity, idle visible TUI with no active turn or
  approval, clean expected branch and worktree, absent merge state, closed prior
  task, current governing instructions, and renewed path and side-effect
  authority. Long idle time alone is not failure. Return a completed standing
  assignment to that same clean idle state. Source: DI-libot
- After a standing process crash, resume the exact thread in the same attempt
  only after proving the thread and rollout identity, tmux pane and worktree,
  Git state, and known side effects. If proof is missing or side effects are
  ambiguous, stop for operator review. A fresh process in the approved standing
  worktree is an explicit recovery action, not a silent replacement. Source:
  DI-libot; TE-judog
- Every dobab update, reconciliation, replan, calendar-graph change, and daily
  close must be performed by the standing dobab worker. Main must not make those
  edits itself. Once the worker commits and self-reviews its promises and returns
  a clean worktree, main applies the same `consensus` process used for every peer
  and promptly merges an accepted exact tip. Dobab work does not wait behind a
  legacy FIFO entry or unrelated review. Source: DI-bijaj; DI-bisaf
- Main repeats this coordination loop: inspect current goals, hard failures,
  dependencies, peer promises, disposable candidate views, and standing health;
  start or resume stable peers for ready work through the DI-kumof self-upgrading
  `resume` procedure; observe panes without blind prompting; use `consensus` to
  inspect and combine exact candidate commits;
  repair failures on main when main chooses to own the correction, or publish a
  TODO/DR/evaluation for another peer; commit main's exact accepted tree with its
  own promise; request separate review only after case-specific user approval;
  account for the outcome; reconcile dobab; and retire only workers whose durable
  promises and obligations are dispositioned. Source: DI-gipob; DI-mikil;
  DI-pufuk; DI-bisaf
- Treat the user directive `mainloop` as an instruction for the main
  `cdint-grid` Codex session to execute one bounded pass through that complete
  loop, with critical-path work first. Honor a stated time box by deferring and
  reporting unfinished work; never shorten a required proof, send blind worker
  input, force a merge, delete retained evidence, or leave an unfinished merge
  at cutoff. Source: DI-kobik
- Within `mainloop`, prioritize a real launch safety failure or a missed
  operational requirement first, followed by the diagnostic or implementation
  work that can remove it, then critical-path integration, then housekeeping.
  Immediately report a missed arbitrary experiment cutoff in literal terms, but
  do not promote it into that blocker tier, optimize toward it, or repeatedly
  rerun it without a new decision. Source: DI-lajul
- The dobab worker reviews the complete remaining schedule for planning errors,
  parallelization opportunities, missing dependencies, stale estimates,
  bottlenecks, and other improvements. Keep the overall-review worker clean and
  idle unless the repo owner explicitly approves one concrete-harm review case.
  An approved review gates only its candidate while unrelated work continues;
  it does not replace origin self-review or receiving-peer consensus. Do not run
  an automatic daily overall review. Source: DI-mikil;
  DI-gafom
- Account one logical worker lane by joining attempt boundaries, thread and
  rollout events, tmux/worktree identity, explicit coordinator annotations, Git
  evidence, and final disposition. Treat `state_5.sqlite` as discovery metadata,
  rollout JSONL as the detailed thread audit, and `logs_2.sqlite` as supplemental
  telemetry. Query SQLite through a normal read-only WAL-aware connection with
  `query_only=ON`; never use immutable mode, checkpoint Codex state, or copy
  message bodies, tool arguments or output, or reasoning summaries into repo
  records. Source: DI-susim
- `effective hours/day` is the sum of active hours across main and all worker
  lanes. `parallel activity percentage` is `100 * effective active hours /
  attended wall-clock hours` for the same accounting period. The coordination
  window bounds extraction, but neither its elapsed duration nor the main
  lane's active subset is the denominator. Union intervals within each lane;
  the aggregate, or one worker lane's contribution, may exceed 100 percent.
  Count foreground tools once as part of active time, record detached unattended
  runtime separately, and classify active intervals as accepted, rejected,
  blocked, or rework. Keep token dimensions separate; monetary cost is out of
  scope. Source: DI-susim; DI-kakuh
- Remove an attempt only after integration or explicit rejection, complete
  accounting, dobab reconciliation, and review closure. Preserve worker branches
  and commits and any experiment evidence with separate retention rules. Source:
  DI-mikil

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
  compiled packages. Use `--isolated-go-cache PURPOSE` only for a test whose
  subject is empty-cache or missing-dependency behavior. Shared-cache cleanup
  requires separate approval. Source: DI-subaz; DI-dogas
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

## Coding Style & Naming Conventions
- Use object-oriented design with structs and methods; avoid large functions and global state.
- Follow generally accepted object oriented design patterns.
- Keep Go code `gofmt`-clean; package names should be short and lower-case.
- Prefer focused edits over broad refactors unless required.
- Add and maintain explanatory comments for non-obvious logic.
- Use `git mv` for file moves/renames to preserve history.

## Error Handling Policy (Required)
- Never use `|| true` in scripts, templates, or make recipes. Always inspect
  command exit codes explicitly with `if/else` branches and handle each outcome.
- For non-fatal cleanup/diagnostics steps, record command status (exit code and
  logs) explicitly; do not fail silently.
- In Go code, never ignore errors with `_ = ...`; handle, propagate, or report
  errors explicitly.
- Run `errcheck ./...` and keep it passing for Go changes.

## DR/DI Source-of-Truth Protocol (Required)
- In this repo, DR and DI logs are the primary source of truth for decisions and open questions.
- Documents and code are outputs of that process and must link back to DR/DI records.
- Person identity in DR/DI records must use full email with label format: `user@example.com (FirstName)`.
- In DRs, `Asked by` and person-valued `Waiting on` fields must use that format.
- In DIs, `Author` must use that format.
- A DI `Author` is the decision-maker, not merely the recorder. If an agent records a DI for Steve's chat decision, Steve is the `Author`; an agent may be `Author` only when Steve explicitly delegates that decision authority to that agent. (DI-034-20260508-060134)
- A settled statement in docs (or critical logic in code comments) must cite at least one DI ID.
- An unresolved question or uncertainty must cite at least one DR ID.
- If an unresolved question has no DR yet, create a DR before finalizing the change.
- During any DR/DI tracking migration, apply these rules incrementally as sections
  and files are brought under DR/DI tracking.
- Intent: New coordination artifacts use a single proquint handle namespace so
  TODO, TE, DR, DI, and DN references do not depend on timestamp or integer
  allocation. Source: DI-sojum; DI-topih; DI-dodud


## Comment Preservation Protocol (Required)
- Never remove existing code comments unless they are replaced in the same patch by equal-or-better explanatory comments near the same logic.
- When rewriting or refactoring code, port old explanatory intent first, then improve wording.
- If a touched non-trivial code block has no comments, add explanatory comments.
- Do not treat shorter comments as better unless they preserve all important intent.
- For any non-trivial behavior change, include a behavior-level comment with:
  - `Intent:` a short, clear rationale (a sentence or a few; no hard cap if more is needed for clarity).
  - `Source:` a DI ID in the format `DI-<handle>`.
  - `<handle>` is minted by `/home/stevegt/bin/mint-handle` and is globally unique across TODO, TE, DR, DI, and DN owners. Source: DI-jufiz
  - Optional: TODO file/section reference for faster lookup.
- If a comment must be dropped with no replacement, stop and ask the user before proceeding.
- Before editing a file, review existing comments in that file.
- Maintain a `## Decision Intent Log` at the top of relevant root
  `TODO/TODO-<handle>-<slug>.md` files. Source: DI-gigoz; DI-topih
- Treat DI logs as append-only history. Do not rewrite or delete prior entries.
- When intent evolves, add a new DI entry and set `Supersedes: <old-di-id>`.
- DI entries must include:
  - `ID: DI-<handle>`
  - `Date: YYYY-MM-DD HH:MM:SS`
  - `Status: active|superseded`
  - `Decision:`
  - `Intent:`
  - `Constraints:`
  - `Affects:`
  - `Supersedes:` (optional)
- After editing, run a comment-delta audit on each touched code file using: `git diff -U0 -- <file> | rg -n '^-\\s*//|^-\\s*/\\*|^\\+\\s*//|^\\+\\s*/\\*'`.
- Resolve all removed-comment lines before finalizing unless explicit user approval was given.
- In the final response, include:
  - `Comment audit: PASS/FAIL`, with file list.
  - `Intent provenance audit: PASS/FAIL`, listing files with behavior changes and DI sources.
- Hard gate: behavior-changing work is incomplete unless comments preserve intent and include DI provenance.
- Do not remove comments or documentation; update them if outdated or incorrect.

### Comment + DI Examples
- Comment format example:
  - `// Intent: Keep context resolution stable across workspace scans to avoid target drift between plan and run. Source: DI-vapoj`
- Decision Intent Log entry template (for TODO files):
  - `ID: DI-<handle>`
  - `Date: YYYY-MM-DD HH:MM:SS`
  - `Status: active`
  - `Decision: <what was decided>`
  - `Intent: <short clear rationale>`
  - `Constraints: <hard limits, dependencies, assumptions>`
  - `Affects: <paths, modules, commands, docs>`
  - `Supersedes: <old DI ID, optional>`

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


## Testing Guidelines
- Use Go's standard `testing` package with deterministic tests.
- Avoid network calls in tests unless explicitly required and documented.
- When changing `plan/run` behavior, add coverage for both command paths when possible.

## Commit & Pull Request Guidelines
- Treat a line containing only `commit` as: add and commit all changes with an AGENTS-compliant message.
- Use short, imperative, capitalized commit subjects.
- Summarize changes per file in commit bodies.
- Stage files explicitly (avoid `git add .` / `git add -A`).
- Do not open GitHub pull requests for normal wire-lab convergence. Steve explicitly dropped the require-PR merge rule; merge by the role-specific direct-push workflow instead. (DI-001-20260428-195702; DI-034-20260508-060134)
- Do not force-push repo branches unless a scoped DI/DR explicitly authorizes the exception. Role overlays may add stricter no-force-push rules for their branches or may document a narrow private-remote exception, but the repo-wide default is no history rewrites. (DI-034-20260508-060134)
- Do not commit local state files, generated binaries, credentials, tokens, signing keys, or other secrets. (DI-034-20260508-060134)

## Vocabulary Source
- `GLOSSARY.md` contains agreed terms, explicitly unresolved `TBD-<slug>`
  placeholders, external standard terms used by this repository, and historical
  terms that must not be mistaken for current vocabulary. Read and maintain it
  according to the Architecture Question Context And Vocabulary rules above.
  Source: DI-lodos
