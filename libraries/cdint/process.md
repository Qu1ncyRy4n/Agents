---
tldr: Scale planning and decision ceremony to the durability and risk of the change.
---
# Process

## Lightweight Escalation

Use the focused change loop for routine fixes, small documentation edits, and
local improvements. Pause and ask before making a durable choice when the work
changes architecture, a public or wire format, persistent data, a security
boundary, an irreversible operation, a public specification, or another
repository's ownership boundary.

Recommend a wider check when the scope, result, or risk looks uncertain.

## Decision First

### Decision Intent

For a change that needs a durable decision, identify the decision before code
or behavior changes. Ask the user to resolve real alternatives, then record the
result in the repository's decision system before implementation.

### Thought Experiment

Run a thought experiment before locking a decision when several plausible,
durable designs remain. Compare the same alternatives under normal operation,
failure or corruption, concurrent actors, long-term evolution, trust-boundary
changes, and scale. State what each alternative makes easier, harder, and newly
required.

### Open Questions

Record a Decision Request when an unresolved question blocks safe progress.
Keep filed analysis and decision history intact; add a new record when intent
materially changes.

## High-Consequence Review

Require human review for changes that affect human data, experiment validity,
production or destructive data operations, security, protocol compatibility, or
an active machine. A passing build alone is not sufficient evidence for these
changes.

## Governed Handoff

For decision-first work, report the decision IDs, evidence that implements each
decision, runtime paths touched, checks run, and only user-approved exceptions.
