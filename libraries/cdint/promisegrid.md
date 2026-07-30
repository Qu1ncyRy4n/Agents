# CDINT And PromiseGrid

## Terminology

Use Promise Theory terms precisely. A request asks another actor to act. A
promise is a voluntary statement by the promiser. Evidence and trust are local
assessments, not universal facts.

A content identifier names bytes. A protocol content identifier names the
specification that defines a protocol role; it does not name one message or one
payload.

## Protocol Changes

Do not add a new top-level protocol action by default. First ask whether the
meaning can be expressed as a voluntary promise with protocol-defined payload
semantics, local evidence, or implementation-local mechanics. Escalate that
choice through the decision-first process when the answer is not clear.

## Collaboration Boundaries

Keep repository-wide rules in the canonical agent guide. Keep role-specific
environment, credentials, branch, and private logging rules in local overlays.
Do not relax a repository-wide rule in an overlay.

## Public Technical Prose

Write direct, specific prose. Use concrete examples where useful. Keep decision
record identifiers out of public slides; use the repository's reference style
for provenance in durable public specifications.
