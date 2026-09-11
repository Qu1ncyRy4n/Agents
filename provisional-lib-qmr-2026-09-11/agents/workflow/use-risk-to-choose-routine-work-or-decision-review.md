tags: [workflow/risk, process/escalation, process/decision]
tldr: Handle routine reversible work directly, but pause for review before durable, risky, or ownership-changing decisions.
sources:
  - libraries/cdint/process/lightweight-escalation.md
  - docs/other_repo_agents/ciwg_decomk-conf-cswg_refs_heads_main_AGENTS.md#decision-first-specification-and-compliance-protocol-required
requires_all: []
mutually_exclusive: []
see_also: []
---
# Use Risk to Choose Routine Work or Decision Review

Use an ordinary focused workflow for routine, reversible changes with clear
requirements. Escalate before making a durable choice when the work changes an
architecture, public interface, persistent state, trust boundary, irreversible
operation, public specification, or another repository's ownership boundary.


## Use the Focused Loop for Routine Work
Use the focused change loop for routine fixes, small documentation edits, and
local improvements. Pause and ask before making a durable choice when the work
changes architecture, a public or wire format, persistent data, a security
boundary, an irreversible operation, a public specification, or another
repository's ownership boundary.

Recommend a wider check when the scope, result, or risk looks uncertain.

## Preserve the Strict Decision-First Option

Decision-first means decisions must be locked before coding; it does not forbid pre-decision analysis such as required thought experiments.
The agent must collect and lock user decisions before making any code edits for a task.
Locked decisions must be recorded as Decision Intent Log entries in the relevant `TODO/*.md` file(s) with clear intent and rationale.
