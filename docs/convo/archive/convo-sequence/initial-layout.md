# Initial Conversation Layout Proposal

Status: historical agent proposal; superseded by CS-D001.
Everything below is history/reference. Resume from [head.md](../../head.md).

The first draft suggested one folder per topic:

```text
docs/convo-<topic>/
  head.md
  work.md
  decisions.md
  history.md
```

The intended improvement was an authoritative NOW panel, one active item, a short
NEXT queue, stable decision IDs and supersession links, and a preserved detailed
narrative. The owner then preferred a shared `convo/` root with topics as archives.

The accepted arrangement keeps those responsibilities but removes competing live
heads: `convo/head.md`, `work.md`, and `decisions.md` form the current layer;
`convo/archive/<topic>/` stores the rich topic history. The original conversation
stack remains linked in the configuration-redesign archive.

Later location refinement (CS-D004): these records now live under `docs/convo/`;
the reusable skill remains under `.agents/skills/`. The earlier root paths above
describe the historical layout choice, not current navigation.
