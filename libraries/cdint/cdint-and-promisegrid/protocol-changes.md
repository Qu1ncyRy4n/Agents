---
tags: [domain/promisegrid, protocol/design, process/escalation, scope/team]
tldr: Avoid new protocol actions unless the meaning cannot be expressed with existing promise semantics.
priority: 0.9
scope: team
requires: [shared:process/decision-first]
---
# Protocol Changes

Do not add a new top-level protocol action by default. First ask whether the
meaning can be expressed as a voluntary promise with protocol-defined payload
semantics, local evidence, or implementation-local mechanics. Escalate that
choice through the decision-first process when the answer is not clear.
