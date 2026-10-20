# Concise Conversation And Validation

## Observed Failure

The owner supplied a response labelled "one item now" that actually prescribed:

1. Implement native carry instrumentation.
2. Validate its read-only API and clearing behavior.
3. Run a 240-tick cross-backend capture schedule with four streams.
4. Diagnose divergence from those streams.
5. Perform a normal-window comparison afterward.

Those are multiple stages of a workstream, not one actionable item. Their detail
can remain useful in the work records; announcing "one item" does not make the
full sequence a concise answer. The response also spent words restating the
skill's intent instead of proceeding directly to the requested validation.

## Better One-Step Response

Illustrative only, based on the owner's example. Assume the referenced evidence
API is approved and implemented; this is not independent authority for Spec33h or
evidence that these checks have been run. Keep any role label the repo requires.

> First validate the read-only carry-evidence API, before recording a long trace.
>
> Check that:
> - No hand returns `None`.
> - Repeated reads neither change evidence nor wake physics.
> - Drop/reset clears evidence.
> - An absent or deliberately constant evidence source makes the oracle fail.
>
> Pass means those checks behave correctly in both backends. It does not mean
> the rifle feels good. Keep controller and tuning parameters unchanged; surface
> any check that needs such changes. Bring back these results before moving on
> to the held-rifle trace.

The larger trace schedule and diagnosis matrix stay in NEXT/history until this
step's results are available. If the API is not implemented, the next item is
clarifying or completing that prerequisite, not issuing a fictional test run.

## Review Cases

These are review scenarios, not executed agent benchmarks:

- **Next validation only:** With instrumentation available, request validation
  directions. Expect one bounded behavior, a short check set, a pass condition,
  and no bundled implementation/trace/tuning campaign.
- **Routine continuation:** Ask whether a proposed name fits the current item.
  Expect the recommendation and brief reason; no repeated NOW panel or queue dump.
- **Full plan explicitly requested:** Ask for the entire validation strategy.
  Expect the full staged plan, clearly distinguished from the single active step;
  do not truncate required analysis merely to meet the ordinary chat word target.
- **Required TE:** A repository requires concrete multi-scenario analysis before
  a durable design decision. Preserve that analysis in the TE and supply a concise
  chat orientation. The one-step rule is not permission to substitute an opinion
  for the required artifact.
