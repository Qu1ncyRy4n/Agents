---
status: reliable
candidate_outputs: [agent-guide]
tags: [workflow/review, safety/high-consequence, process/escalation]
tldr: Require human review when changes affect sensitive data, safety, security, compatibility, irreversible operations, or an active machine.
sources:
  - libraries/cdint/process/high-consequence-review.md
requires_all: []
mutually_exclusive: []
see_also: []
---
# Require Human Review for High-Consequence Changes

Require human review when changes affect human data, experiment validity,
production or destructive data operations, security, protocol compatibility,
irreversible operations, or an active machine. Passing automated checks is
necessary evidence when available, but it is not sufficient authorization or
validation for these changes.

Verbatim v1 library source from
[`High-Consequence Review`](../../../libraries/cdint/process/high-consequence-review.md):

> Require human review for changes that affect human data, experiment validity,
> production or destructive data operations, security, protocol compatibility,
> or an active machine. A passing build alone is not sufficient evidence for
> these changes.

## Review Concerns

- Trace this adaptation back to every source guide that contributes a danger
  zone; the current v1 module citation is not sufficient final provenance.
- Decide whether “human review” means approval before implementation, before a
  destructive action, before commit, or before deployment in each risk class.
- Consider whether irreversible operations should be named explicitly in the
  final version.
