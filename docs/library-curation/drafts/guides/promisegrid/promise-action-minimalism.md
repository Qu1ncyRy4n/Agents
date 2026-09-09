# Keep PromiseGrid Wire Actions Minimal

Output type: PromiseGrid guide. Default for CDINT selection; may be explicitly
excluded. Selection policy and final domain path remain review items.

Source: [`promisegrid_wire-lab_refs_heads_main_AGENTS.md`](../../reference/source-guides/promisegrid_wire-lab_refs_heads_main_AGENTS.md#promise-action-minimalism-required)

> Future PromiseGrid protocol, simulation, POC, scoring, generation, and guide
> work must not invent workflow-specific top-level action kinds by default. The
> default future-facing top-level semantic action is `promise`. (DI-mosoj)
>
> Treat observation as a promise that the promiser observed something from its
> local vantage. Treat refusal as absence of a promise, a promise not to do
> something, or a promise that the agent does not currently promise the
> requested behavior. (DI-mosoj)
>
> Treat repair, offer, counteroffer, acceptance, routing, introduction,
> redemption, transfer, storage, computation, TCP-link changes, authorization,
> dispatch, grant, registration, and enforcement as pCID-defined payload
> semantics, local trust/evidence interpretation, or implementation-local
> mechanics unless a scoped TE/DI proves a distinct wire-level role. (DI-mosoj)
>
> Before adding any new top-level action kind, stop and answer: why is this not
> just an agent's voluntary promise with pCID-defined payload meaning? If that
> answer is not explicit in a locked TE/DI, do not add the action kind.
> (DI-mosoj)
>
> POC7 through POC10 action names, including refusal and observation labels,
> are historical executable evidence, not the forward naming pattern. Do not
> rewrite scored/generated artifacts in place; supersede or reframe them when
> reused. (DI-mosoj; supersedes DI-fitav)
