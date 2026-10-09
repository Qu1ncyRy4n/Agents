# DR: Name Of The Sidecar Selection-Guidance Attribute

DR-ID: DR-nalan
Date: 2026-10-05
State: accepted for v2; `curate` is in use
Asked by: Quincy Ryan

## Question

What should the v2 library sidecar call the attribute that states how
strongly the library stands behind a branch's direct children when a consumer
selects the branch broadly? The attribute currently reads
`curate = "foundation" | "choose" | "optional" | "opt_in"` in
`library.mogent.hcl`, with `defaults` and `default` beside it.

The name must survive in three places at once: the sidecar a library author
writes, the diagnostic a consumer reads (`plan` prints it), and the spec. It
replaced the earlier `inclusion { policy = "baseline" | "explicit" | ... }`
block, which was rejected as abstract and unreadable in error messages.

## Candidates

| Key | Read aloud | Assessment |
|---|---|---|
| `curate` | "curate this branch as foundation" | Library-side guidance that explains how the consumer should handle the branch's direct children. Chosen for v2. The existing values remain readable. |
| `offer` | "library offers constraints as foundation" | Historical provisional choice. It reads more abstractly in authoring and diagnostics than `curate`. |
| `choice` | "your choice here is required" | Consumer-side voice. Equally explicit, arguably clearer at a glance. Requires renaming values to `none`, `required`, `optional`, `opt_in` because `choice = "choose"` is redundant. Loses the metaphor and the PT word. The serious alternative. |
| `policy` | "policy is baseline" | Generic abstract noun. Best of the earlier round, still the complaint that started this. |
| `intent` / `use_intent` | "intent is foundation" | Vague; says nothing about children. Hyphen illegal in unquoted HCL names. |
| `enforcement` | "enforcement is optional" | Self-contradiction. The library cannot enforce; it can only promise. |
| `structural_level` / `policy_level` / `tier` | ordinal | Three of four values sit on a strength axis but `choose` is a demand for a decision, not a strength. `level` also collides with Markdown heading level inside a heading-tree tool. |
| `contract` | "contract is optional" | A contract binds both sides; the library writes this alone and the consumer has signed nothing. |
| `disclaim` | "disclaim foundation" | Inverted direction: the library pushes content toward the consumer. |
| `placement` | "placement is foundation" | Says where, not how strongly; collides with manifest placement flags. |
| `purpose` | "purpose is choose" | Vague and nonsensical with half the values. |

## Decision

- Use `curate` now. It describes practical authoring intent without implying
  the library imposes a contract on the consumer.
- Keep `foundation`, `choose`, `optional`, and `opt_in`; they carry the
  behavioral meaning and require no semantic migration.
- A later naming change needs a new decision and explicit migration rather
  than accepting two sidecar grammars indefinitely.

## What Finalizing Requires

The v2 cutover changes the HCL struct tag in `v2/library.go`, CLI labels,
tests, specification, README, demos, and the QMR sidecar. Legacy `offer` is
rejected with a direct migration message. An explicit sidecar migration command
remains future work.

Related open question in the spec: whether `curate` becomes required on every
branch section rather than inherited from the parent's broad selection.
