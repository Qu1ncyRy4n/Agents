---
status: reliable
tags: [process/decision, cognition/scenario-analysis, artifact/thought-experiment]
tldr: Use a thought experiment to reduce a broad design space to a few viable options before asking for the final decision.
sources:
  - docs/other_repo_agents/ciwg_decomk-conf-cswg_refs_heads_main_AGENTS.md#thought-experiment-protocol-required
requires: []
conflicts_with: []
exclusive_group:
see_also: []
---
# Use Thought Experiments to Narrow a Broad Design Space

A thought experiment comes before final decision questions. Start with the
decision being tested and all plausible alternatives that survive basic factual
checks. Use concrete scenarios to expose consequences, reject weak alternatives,
and narrow the field. The output is a few surviving choices and the unresolved
questions needed for a human decision; the thought experiment does not itself
lock that decision.

This corrects a misleading compression in the v1 library: a thought experiment
does not begin only after the choices have already been narrowed to the final
few. Narrowing is its central job.

## Identify the Decision, Alternatives, Assumptions, and Scope
<!-- tldr: Define what is being tested before comparing designs. -->

Verbatim source from
[`TE Intake Requirements`](../../../../docs/other_repo_agents/ciwg_decomk-conf-cswg_refs_heads_main_AGENTS.md#te-intake-requirements):

> Before locking decisions or asking final DF questions, the agent must identify:
>
> - the decision being tested,
> - the candidate alternatives,
> - the assumptions and threat/trust model,
> - the scope and systems affected.
>
> If the TE relates to an existing TODO, the agent must reference the TODO handle and subtask handle (for example, `fogus.5`).

The TODO reference is repository-policy-specific. Retain it when the consuming
repository uses that coordination system; the intake questions stand without
it.

## Compare Every Alternative Under the Same Concrete Scenarios
<!-- tldr: Hold assumptions constant and test normal, failure, evolution, trust, concurrency, and scale behavior. -->

Verbatim source from
[`TE Execution Requirements`](../../../../docs/other_repo_agents/ciwg_decomk-conf-cswg_refs_heads_main_AGENTS.md#te-execution-requirements):

> Each TE must evaluate the same decision across multiple concrete scenarios.
> Scenarios must include, when relevant:
>
> - normal operation,
> - failure/corruption/incomplete writes,
> - concurrent actors or mixed-version nodes,
> - long-horizon evolution and migration,
> - trust-boundary changes,
> - scale effects (storage, bandwidth, CPU, operational complexity).
>
> The agent must compare alternatives under the same assumptions instead of switching assumptions mid-analysis.
> The agent must state what each alternative makes easier, what it makes harder, and what new obligations it creates.

## Narrow the Field Before Asking for the Final Decision
<!-- tldr: Reject weak options, retain viable ones, and expose the exact remaining choices. -->

Verbatim source from
[`TE Output to DF`](../../../../docs/other_repo_agents/ciwg_decomk-conf-cswg_refs_heads_main_AGENTS.md#te-output-to-df):

> After the TE, the agent must identify:
>
> - rejected alternatives,
> - surviving alternatives,
> - unresolved questions that still require user choice,
> - any new naming/path/runtime decisions exposed by the TE.
>
> Final DF questions must be framed from the surviving alternatives identified by the TE. The agent must not ask broad DF questions that ignore TE results.

## Keep Analysis and Decision Authority Separate
<!-- tldr: A recommendation may follow the analysis, but a TE is not a locked decision. -->

Verbatim source from
[`TE Decision Rules`](../../../../docs/other_repo_agents/ciwg_decomk-conf-cswg_refs_heads_main_AGENTS.md#te-decision-rules):

> A TE does not by itself lock a decision.
> After the TE, the agent must either:
>
> - ask the user to choose among the surviving alternatives, or
> - recommend one surviving alternative and clearly state why the others were rejected.
>
> After user choice is resolved, the agent must record the locked result via the existing DI process before implementation.
> If a TE exposes a new ambiguity, dependency, or naming/path decision, the agent must stop and resolve that before implementation.

The requirement to use a particular DI process is a selectable governance
policy. The reusable core is that analysis, recommendation, human authority,
and the durable record of a decision are distinct stages.

## Worked Example: Choose a Library Relationship Address
<!-- tldr: Begin broadly, eliminate unsafe identity schemes through scenarios, then present the viable choices. -->

Proposed example:

1. Decision under test: how reusable modules refer to related modules.
2. Initial alternatives: consumer aliases, bare paths searched across sources,
   stable global library identities, manifest-only relationships, and `self:`
   for same-library references.
3. Scenarios: alias rename, two sources with the same path, a fork retaining
   metadata, an offline build, and a moved source path.
4. Rejections: consumer aliases couple the library to one manifest; bare paths
   become ambiguous; automatic global identity introduces unresolved fork and
   trust rules.
5. Survivors: `self:` for intrinsic same-library relationships and manifest
   aliases for integration between declared sources; stable global identity is
   deferred until cross-library versioning requires it.
6. Final decision question: approve the hybrid now, or accept the cost of a
   stable published-library identity in this milestone?

The thought experiment began with five plausible approaches and produced two
decision-ready directions. It did not merely compare two preselected answers.
