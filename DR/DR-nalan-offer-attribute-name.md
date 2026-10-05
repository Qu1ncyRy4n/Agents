# DR: Name Of The Sidecar Selection-Guidance Attribute

DR-ID: DR-nalan
Date: 2026-10-05
State: open; `offer` is in use provisionally
Asked by: Quincy Ryan

## Question

What should the v2 library sidecar call the attribute that states how
strongly the library stands behind a branch's direct children when a consumer
selects the branch broadly? The attribute currently reads
`offer = "foundation" | "choose" | "optional" | "opt_in"` in
`library.mogent.hcl`, with `defaults` and `default` beside it.

The name must survive in three places at once: the sidecar a library author
writes, the diagnostic a consumer reads (`plan` prints it), and the spec. It
replaced the earlier `inclusion { policy = "baseline" | "explicit" | ... }`
block, which was rejected as abstract and unreadable in error messages.

## Candidates

| Key | Read aloud | Assessment |
|---|---|---|
| `offer` | "library offers constraints as foundation" | Library-side voice. Matches Promise Theory: a `+` promise is an offer; the consumer's `select` is the `-` promise; the library never imposes. Keeps the building-metaphor values. Provisional choice. |
| `choice` | "your choice here is required" | Consumer-side voice. Equally explicit, arguably clearer at a glance. Requires renaming values to `none`, `required`, `optional`, `opt_in` because `choice = "choose"` is redundant. Loses the metaphor and the PT word. The serious alternative. |
| `policy` | "policy is baseline" | Generic abstract noun. Best of the earlier round, still the complaint that started this. |
| `intent` / `use_intent` | "intent is foundation" | Vague; says nothing about children. Hyphen illegal in unquoted HCL names. |
| `enforcement` | "enforcement is optional" | Self-contradiction. The library cannot enforce; it can only promise. |
| `structural_level` / `policy_level` / `tier` | ordinal | Three of four values sit on a strength axis but `choose` is a demand for a decision, not a strength. `level` also collides with Markdown heading level inside a heading-tree tool. |
| `contract` | "contract is optional" | A contract binds both sides; the library writes this alone and the consumer has signed nothing. |
| `disclaim` | "disclaim foundation" | Inverted direction: the library pushes content toward the consumer. |
| `placement` | "placement is foundation" | Says where, not how strongly; collides with manifest placement flags. |
| `purpose` | "purpose is choose" | Vague and nonsensical with half the values. |

## Reasoning Behind The Provisional Choice

- The sidecar is the library's file and should speak in the library's voice.
  `offer` does; `choice` makes the library narrate the consumer's situation.
- Promise Theory gives one correct word for a `+` promise, and the
  architecture already is that model: offer on one side, select on the other,
  `force_exclude` with `reason` as an on-record refusal of a strong offer.
- The values carry most of the meaning in diagnostics. `foundation` is
  self-explanatory ("you do not remove a foundation without a reason");
  `choose` says exactly what the consumer must do; `optional` and `opt_in`
  are ordinary English. These values survive under `offer` and not under
  `choice`.
- Explicitness, the owner's stated priority, is nearly equal between `offer`
  and `choice`. The tie breaks on voice and on keeping the values.

## What Finalizing Requires

Pick `offer` or `choice`. If `choice`, rename the values as above. Either way
the change is one commit touching: the `Offer*` constants and the HCL struct
tag in `v2/library.go`, the `offerLabel` and `Offer:` lines in
`internal/cli/v2_source.go`, `v2/offer_test.go`, `docs/proposals/MOGENT-HCL-V2-SPEC.md`,
`README.md`, and the QMR `library.mogent.hcl`. Record the decision as a DI
in `TODO/TODO.md` and set this DR's state to resolved.

Related open question in the spec: whether `offer` becomes required on every
branch section rather than inherited from the parent's broad selection.
