# Multiple Agent Outputs Plan

Status: planning only; do not implement from this document without a recorded
Decision Intent.

## Goal

Let one manifest produce `AGENTS.md` plus explicit reviewable tool entrypoints
such as `CLAUDE.md`, `GEMINI.md`, or `.codex/AGENTS.md`. Preserve one authored
document model and make every tool-specific difference visible in YAML.

## Proposed Manifest Evolution

Keep `output` as the primary generated document for backward compatibility.
Add an optional ordered `outputs` list:

```yaml
output: AGENTS.md

outputs:
  - path: CLAUDE.md
    mode: render
    extra:
      - Claude Strictness: shared:tools/claude/strictness
  - path: GEMINI.md
    mode: symlink
    target: AGENTS.md
```

`render` starts from the base `doc` and appends or inserts explicitly declared
manifest entries. The first implementation should allow only one clear extra
placement rule; arbitrary patching of the base tree would be too difficult to
review. `symlink` creates a relative link to another declared output and never
silently copies when linking fails.

An exact rendered mirror is `mode: render` with no extras. This differs from a
symlink operationally and must remain explicit.

## Core Model

Add a reusable build-plan operation that resolves all output paths, validates
every render and link target, checks path collisions and escapes, and returns a
complete set of proposed writes before changing disk. CLI, TUI, and future
clients consume the same plan.

The operation must reject:

- duplicate or overlapping output paths;
- absolute paths or relative paths escaping the manifest directory;
- symlink cycles, missing targets, and targets not declared by the manifest;
- output paths inside source libraries or `.mogent` runtime state;
- a partially valid set where any render fails.

## State, Drift, And Transactions

Evolve generated-output state from one hash to a versioned map keyed by output
path. Preserve migration from the current `.mogent/state.json` shape.

Build is one transaction: validate all content, snapshot every target and state
file, write temporary files/links, then install the complete set. Failure rolls
back every output. Status and drift report each output independently while also
returning a nonzero aggregate state when any output is missing, stale,
untracked, or directly edited.

Direct-edit import initially applies only to rendered outputs and must map the
edited section back to the shared base document or that output's explicit
extras. Symlink outputs have no independent import path.

## CLI And Review

Planned command behavior:

- `mogent build` validates and writes the full output set.
- `mogent build --output <path>` may be considered only for diagnostics; normal
  writes should stay transactional across all declared outputs.
- `mogent status` lists one row per output and an aggregate result.
- `mogent drift [--output <path>]` inspects a chosen rendered output.
- `mogent init` templates can offer optional output entries but must print them
  in dry-run YAML.

Dry-run/review should show output mode, target, extra manifest nodes, rendered
diff summary, and every filesystem path affected.

## Implementation Sequence

1. Run a thought experiment for extra-entry placement and cross-platform
   symlink policy; record the chosen schema in a DI.
2. Extend strict manifest parsing and normalized serialization.
3. Introduce the multi-output build plan without changing writes.
4. Version and migrate generated-output state.
5. Add transactional file and symlink installation.
6. Extend status, drift, previews, init, completion, and dogfood fixtures.
7. Add failure-injection tests proving rollback across several targets.

## Decisions Required Before Implementation

1. Are output extras append-only, placed under a named manifest heading, or a
   complete per-output document tree? Append-only is the smallest clear first
   contract but may not fit tool-specific identity rules.
2. On platforms or filesystems without symlink support, should symlink mode fail
   loudly or allow an explicitly different `copy` mode? It must never silently
   change semantics.
3. May outputs live in nested directories such as `.codex/AGENTS.md`, and which
   directories are forbidden beyond `.mogent` and source roots?
4. Does direct-edit import from a tool-specific extra create a normal local
   override, or an output-scoped local override?

No multiple-output code should land until these choices are resolved and the
decision-to-test matrix is recorded.
