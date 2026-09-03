# seam_game Agent Prompt Assessment

## Source Files Found

Primary prompt file:

- `CLAUDE.md`

Related context:

- `LEARNING.md`
- `.claude/settings.local.json`
- `docs/decisions.md`
- `docs/vision.md`
- `docs/traits.md`
- `docs/core-loop.md`

## Qualitative Assessment

`CLAUDE.md` is unusually strong. It does what most agent prompt files fail to do: it gives both architectural law and collaboration law. The technical constraints are concrete enough to prevent broad, damaging changes, and the tutor-mode contract is explicit about which parts the agent may write versus which parts Q should implement for learning.

The file is dense, opinionated, and specific. That is a strength for this project. A generic assistant would otherwise be likely to "helpfully" collapse the core design into conventional Godot gameplay code, overuse LLMs in the simulation layer, hand-author affordances, or skip the learning contract by implementing the interesting Rust systems itself. `CLAUDE.md` blocks those failure modes directly.

The main weakness is portability. The file is written for Claude, but most of its content is agent-agnostic. Converting or duplicating it as `AGENTS.md` would make the contract visible to Codex and other tools that preferentially read `AGENTS.md`. It also mixes permanent project law, current roadmap, and session-process rules in one file; that is workable now, but it may become harder to maintain as the repo grows.

## Strengths

- Clear project identity: solo-dev immersive sim/action RPG, Godot 4 plus Rust, vertical-slice target.
- Strong architecture laws:
  - Rust owns authoritative sim state.
  - Godot is a rendering terminal.
  - Physics is Godot/Jolt; Rust consumes semantic consequences.
  - Affordances are derived relational confidence values, not hand-authored booleans.
  - LLMs sit at the surface and never mutate world state directly.
  - Death/return is diegetic and world-remembered.
- Explicit stack and roadmap.
- Strong anti-scope-creep guidance.
- Excellent "Working With Q" section: directness, typo tolerance, fatigue-aware interaction, and preference for opinionated recommendations.
- Tutor mode is unusually actionable:
  - Agent writes scaffolding and glue.
  - Q writes the "brain" code.
  - The agent reviews and teaches instead of silently replacing core logic.
  - `LEARNING.md` is maintained as portable tutoring state.
  - Micro-goals and commits are used to keep progress bounded.
- Repo structure and toy rules are clear.

## Risks For Agents

- The prompt is named `CLAUDE.md`, so tools that look for `AGENTS.md` may miss it.
- There is no concise quick-start section for non-Claude agents that says "read this first, then read `LEARNING.md`, then `docs/decisions.md`."
- Validation commands are not listed in the prompt. It names Nix, Rust, Godot, and gdext but does not give a standard test/build ladder.
- Branch, commit, and formatting policies are implicit.
- The tutor contract says one proven step equals one commit, but the prompt does not specify whether agents should actually commit autonomously or wait for user confirmation in all contexts.
- The architecture laws are strong, but there is no "common bad changes" checklist. That would help agents catch tempting violations before editing.
- `.claude/settings.local.json` contains machine-specific allowed commands and paths under `/home/qix/dev/seam_game`, while this checkout is under `/Users/q/Documents/Code/seam_game`. That file should not be treated as portable project instruction.

## Recommended Additions

Add a root `AGENTS.md` that either copies the relevant content from `CLAUDE.md` or points to it explicitly. Best practical option: create `AGENTS.md` as the canonical cross-agent file and let `CLAUDE.md` either duplicate it or refer to it.

Specific improvements:

1. Add a short "Read order" section:
   - `AGENTS.md` / `CLAUDE.md`
   - `LEARNING.md`
   - `docs/decisions.md`
   - relevant design docs under `docs/`
2. Add validation commands:
   - `nix develop`
   - `cargo check` or `cargo test` for Rust crates when present
   - Godot headless/editor commands once stable
3. Clarify commit behavior:
   - Micro-goal commits are desired, but agents should ask before committing unless the user has already authorized commits for the session.
4. Add "common violations to reject":
   - Do not put authoritative state in Godot nodes.
   - Do not hand-author affordance booleans.
   - Do not let LLM output mutate ECS state directly.
   - Do not copy toy code into `sim/`; rewrite it.
   - Do not implement Q-owned "brain code" unless explicitly asked.
5. Split mutable session process from durable architecture:
   - Keep core laws in `AGENTS.md`.
   - Keep tutoring state in `LEARNING.md`.
   - Keep volatile current task notes in `PICKUP.md` or TODO docs.
6. Add expected style for docs:
   - Dense decisions, not transcripts.
   - Mark fact versus opinion.
   - Record open questions in `docs/decisions.md`.

## Suggested Skeleton

```markdown
# AGENTS.md

This project is SEAM, a Godot 4 + Rust immersive sim/action RPG. Rust owns authoritative simulation state. Godot renders dumb puppet nodes.

## Read First

1. `LEARNING.md`
2. `docs/decisions.md`
3. Relevant design docs in `docs/`

## Architecture Laws

- Rust owns authoritative ECS state.
- Godot nodes do not own gameplay state.
- Godot/Jolt owns physics; Rust consumes semantic consequences.
- Store causes, derive effects.
- Affordances are relational confidence values, not booleans.
- LLMs stay at the surface and must pass through validated command grammar.
- Toys are fossils and practice grounds; rewrite into `sim/`, do not copy-paste.

## Tutor Contract

- Agent may write scaffolding, glue, config, repetitive data, tests, and harnesses.
- Q writes core sim logic, derivers, ECS systems, affordance math, and hot-path code.
- Explain shape and data flow first, then review Q's implementation honestly.
- Update `LEARNING.md` after each tutoring session.

## Validation

- Use `nix develop` for the project shell.
- Run relevant Rust checks/tests before reporting done.
- Use Godot headless checks when the touched toy or bridge supports them.
```
