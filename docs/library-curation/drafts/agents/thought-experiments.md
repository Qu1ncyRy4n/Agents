# Use Thought Experiments Before Design Decisions

Current output: agent guidance for initial dogfooding. Skill extraction remains
deferred with the rest of decision-first workflow.

## Default Trigger

Source: [`promisegrid_wire-lab_refs_heads_main_AGENTS.md`](../reference/source-guides/promisegrid_wire-lab_refs_heads_main_AGENTS.md#thought-experiment-protocol-required)

> Before locking any non-trivial decision that will require DF questions and
> answers, the agent must run a thought experiment (TE) if multiple plausible
> designs remain.
>
> A TE happens before final DF questions. Its purpose is to narrow the design
> space so DF questions and answers are informed by explicit scenario analysis.
>
> The agent must not collapse a TE into a short opinion or recommendation. The
> agent must explicitly model concrete scenarios and consequences.

Default: require a TE when multiple plausible designs remain. Future mutually
exclusive profiles may propose a TE for every non-trivial decision or require
one for every non-trivial decision.

## Neutrality

Source: [`newest_promisegrid.md`](../reference/source-guides/newest_promisegrid.md#thought-experiment-protocol-required)

> The agent MUST NOT prejudice a thought experiment outcome when planning a
> thought experiment -- the agent must not pre-select a preferred alternative
> or design, and must not bias the TE toward a particular outcome.

Compare alternatives under the same assumptions.

## Intake

Before locking decisions or asking final decision questions, identify:

- the decision being tested;
- the candidate alternatives;
- the assumptions and threat/trust model; and
- the scope and systems affected.

If the TE relates to existing work, reference its governing task.

## Scenarios

Default: include every relevant source scenario:

- normal operation;
- failure, corruption, or incomplete writes;
- concurrent actors or mixed-version nodes;
- long-horizon evolution and migration;
- trust-boundary changes; and
- scale effects across storage, bandwidth, CPU, and operational complexity.

State what each alternative makes easier, what it makes harder, and what new
obligations it creates. A reduced risk-based scenario profile remains a future
option; the full relevant list is the current default.

## Output To Final Decisions

After the TE, identify:

- rejected alternatives;
- surviving alternatives;
- unresolved questions that still require user choice; and
- new naming, path, or runtime decisions exposed by the analysis.

Frame final decision questions from the surviving alternatives. A TE does not
lock a decision. After the user chooses, record the decision before
implementation.

## Artifact

Always write a standalone TE before the final decision. It must stand on its
own and include title, TE ID, decision under test, assumptions, alternatives,
scenario analysis, conclusions, and implications for open work and pending
decisions.

Track the TE from its governing work item. Report its path, surviving
alternatives, and decision state. A future profile may require standalone
artifacts only for durable or high-consequence decisions; that is not the
current default.

## Editing Filed Thought Experiments

Treat a filed TE as durable analysis evidence. Update mechanical current
pointers in place while preserving historical quotations. Put navigational
updates and resolved follow-ups in its append-only `## Refinements` section.
Write a superseding TE when a conclusion, scope, or decision impact changes.

Use the [`edit-thought-experiment`](../skills/edit-thought-experiment/SKILL.md)
skill for the complete category, status, reading-scope, and supersedence
procedure.
