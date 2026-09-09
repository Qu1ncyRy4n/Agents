# Keep Planning Records In A Root TODO Directory

Output type: guide. Keep separate from always-loaded agent guidance. Tentative
CDINT default; confirmation requested in
[`provenance/questions.md`](../../provenance/questions.md).

## Tentative CDINT Default

Source: [`newest_promisegrid.md`](../../reference/source-guides/newest_promisegrid.md#project-structure--module-organization)

> Keep planning artifacts in the root `TODO/` directory for this repo and maintain
> `TODO/TODO.md` sorted by priority. TODO files use
> `TODO/TODO-<handle>-<slug>.md`, where `<handle>` is minted by
> `/home/stevegt/bin/mint-handle`. The per-protocol `protocols/<slug>.d/TODO/` layout used
> by wire-lab is not the active layout for cdint-grid unless a later DI migrates
> this repo. Source: DI-gigoz; DI-topih; DI-sojum; DI-jufiz
>
> Keep durable worker-attempt records beside their one owning TODO under
> `TODO/TODO-<handle>-<slug>.d/<full-attempt-id>/`. The TODO Markdown file and
> matching `.d` directory have the same full stem and must move together with
> `git mv` when the TODO slug changes. Related TODOs link to the owning copy;
> they do not duplicate it. Source: DI-jifus

## Retained Alternatives

Source: [`promisegrid_wire-lab_refs_heads_main_AGENTS.md`](../../reference/source-guides/promisegrid_wire-lab_refs_heads_main_AGENTS.md#project-structure--module-organization)

Wire-lab stores planning records per protocol:

> Keep planning artifacts in per-protocol `protocols/<slug>.d/TODO/` directories
> (harness-level under `protocols/wire-lab.d/TODO/`) and maintain the master
> cross-listed index at `protocols/wire-lab.d/TODO/TODO.md` sorted by priority.
> Each `protocols/<slug>.d/TODO/` also has its own per-protocol `TODO.md` queue.

Source: [`promisegrid_grid-poc_refs_heads_main_AGENTS.md`](../../reference/source-guides/promisegrid_grid-poc_refs_heads_main_AGENTS.md#task-tracking-todo)

Grid-poc uses a root TODO directory with numeric IDs:

> Track tasks and plans in `TODO/`, with a `TODO/TODO.md` index of small tasks.
> Number TODOs with zero-padded IDs (e.g., `007`), don't renumber, and mark
> completed items by checking them off (e.g., `- [x] 007 ...`).
