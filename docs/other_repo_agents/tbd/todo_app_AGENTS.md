# Agent Guide for `todo_app_project`

This repository is an early Rust CLI todo app that is also a design seed for a
larger omnicortex task substrate. Keep changes small, explicit, and aligned with
the documented direction in `SPEC.md` and `DESIGN-QUESTIONS.md`.

## Project Shape

- Language: Rust, edition 2024.
- Entry point: `src/main.rs`.
- Current dependency: `clap` in `Cargo.toml`, though the current binary does not
  yet use it.
- Current behavior: interactive CLI prompts for task descriptions until `exit`,
  stores tasks in memory, and prints them.
- Design target: a content-addressed graph of Nodes, Tags, and typed Links/Edges
  that can eventually share data formats with other omnicortex tools.

## Core Design Direction

Treat `SPEC.md` as the high-level source of truth:

- Nodes are the shared primitive for tasks, notes, roles, goals, projects,
  values, sources, and related content.
- Tags classify and cluster nodes.
- Links should evolve toward typed edges such as `blocks`, `part-of`,
  `depends-on`, `promised-to`, or `derived-from`.
- Prefer immutable/content-addressed records plus mutable human-readable pointers
  when designing persistence.
- Avoid baking in a single global `done` truth as the long-term model. The design
  notes point toward evidence, observations, and local assessments.

The current `Task { description, done, id }` struct is a short-term prototype,
not the durable schema.

## Development Workflow

Use Cargo for local checks:

```sh
cargo fmt
cargo check
cargo test
```

There are no tests yet, so add focused tests when introducing parsing, storage,
graph behavior, ID generation, or command handling. At minimum, run
`cargo check` after Rust changes.

## Implementation Guidelines

- Preserve the CLI-first path unless the task explicitly moves toward a GUI.
- Keep data-model changes compatible with the Nodes/Tags/Links direction.
- Prefer typed Rust structures and enums for domain concepts instead of strings
  spread through command handling.
- Keep user-facing CLI output plain and predictable.
- Do not introduce a database, async runtime, GUI framework, or broad dependency
  unless the change clearly needs it.
- If adding persistence, prefer an inspectable local format first and document
  where files are written.
- Keep integration with other omnicortex tools format-based rather than relying
  on a shared runtime.

## Known Rough Edges

- `src/main.rs` currently has an indexing bug in the non-empty task case:
  `all_tasks[all_tasks.len()]` is out of bounds. Use `last()` if this path is
  touched.
- `clap` is declared but unused.
- The README contains rough project notes rather than polished user
  documentation.
- Spelling and naming in docs are still draft-quality; preserve meaning when
  cleaning them up.

## Agent Behavior

- Read `SPEC.md` and `DESIGN-QUESTIONS.md` before making schema or architecture
  changes.
- Keep edits scoped to the request. Avoid broad rewrites while the project is in
  this exploratory stage.
- Call out design tradeoffs when a change would commit the project to a long-term
  representation.
- Do not overwrite user changes in the working tree.
