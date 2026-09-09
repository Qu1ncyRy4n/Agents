# Lock Decisions Before Code Edits

Current profile: strict default. Keep this complete workflow in agent guidance
for initial dogfooding. Future lighter profiles may be mutually exclusive
alternatives. Skill extraction is deferred and must remain reversible.

## Source Rules

Source: [`promisegrid_wire-lab_refs_heads_main_AGENTS.md`](../reference/source-guides/promisegrid_wire-lab_refs_heads_main_AGENTS.md#decision-first-specification-and-compliance-protocol-required)

> Decision-first means decisions must be locked before coding; it does not
> forbid pre-decision analysis such as required thought experiments.
>
> The agent must collect and lock user decisions before making any code edits
> for a task.
>
> Locked decisions must be recorded as Decision Intent Log entries in the
> relevant `protocols/<slug>.d/TODO/TODO-<handle>-<slug>.md` file(s) with clear
> intent and rationale.
>
> The agent must ask decision questions up front in a single intake round
> whenever possible.
>
> Required decision categories are architecture, design/behavior,
> implementation approach, function naming, variable naming, and file/path
> decisions.
>
> The agent must ask these as multiple-choice questions whenever practical.
>
> When a thought experiment (TE) is required, the agent must complete the TE
> before asking final DF questions. TEs narrow alternatives; DF questions and
> answers lock the decision before implementation.
>
> Thought experiments (TEs) are analysis artifacts; Decision Intent (DI)
> entries are the separate records that capture the locked decision after DF is
> resolved.

## Current Adaptation

1. Gather discoverable facts before asking questions.
2. Use a small related question batch when that is easier to answer. Ask one at
   a time for complex or dependent choices. A configurable `n`-question limit
   is deferred.
3. Check all six decision categories on every invocation. This is semantic
   coverage, not a formatting requirement and not a demand for six questions.
   Mark a category settled or not applicable when no real choice exists. Ask
   and lock only unresolved choices.
4. Approve bounded path sets for ordinary work; stop if work leaves them.
   Require exact approval for risky writes, protected data, external side
   effects, destructive actions, and other high-consequence paths. Final path
   policy remains a question for Steve.
5. Present a Decision Lock before editing. Stop if a required choice remains
   missing, ambiguous, or conflicting.
6. Record implementation evidence and any deviation from the locked direction.
