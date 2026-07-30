# AGENTS.md

## Identity

This repository designs mogent, a tool that composes repository-specific
`AGENTS.md` files from Markdown libraries and an `agents.yaml` manifest.

`docs/DESIGN.md` is the design of record. `TODO/` holds decisions and planned
work. The current repository is documentation-first while the new implementation
is rebuilt. Source: DI-ralik.

Keep prose, slides, and small tools in purpose-named top-level directories.
Keep planning records in `TODO/`, decision requests in `DR/` when needed, and
thought experiments in `docs/thought-experiments/`.

## Workflow

Start by understanding the request and checking the relevant design, decision,
and source files. Keep the change focused. Preserve user changes in the working
tree. Validate the changed area, inspect the diff, then report what changed and
what you checked.

Use the lightweight workflow for routine fixes, small documentation edits, and
local improvements. Escalate before editing when work changes architecture,
wire or public formats, persistent data, security boundaries, irreversible
operations, public specifications, or another repository's ownership boundary.
Escalate too when the scope or test result looks uncertain. Source: DI-venit.

For escalated work, identify the decisions first. Run a thought experiment when
several durable designs remain. Ask the user to choose when needed, then record
the result as a Decision Intent before implementation. Record unresolved
questions as Decision Requests. Use the full decision-first handoff only for
that escalated work.

## Changes And Validation

- Make small changes that directly answer the request. Do not reorganize or
  rewrite unrelated material.
- Do not delete or overwrite work unless the user has explicitly requested that
  action.
- Preserve explanatory comments. Replace a comment only with an equally useful
  or clearer explanation near the same logic.
- Use `git mv` for moves and renames.
- Documentation-only changes: check structure, links where practical, and the
  final diff.
- Go behavior changes: run `gofmt`, focused `go test`, and `errcheck ./...`
  from the affected Go module. Keep tests deterministic and avoid network calls
  unless the task explicitly needs them.
- Handle Go errors explicitly and use `%w` when adding context. Do not hide
  shell-command failures with `|| true`.
- Keep temporary files and Go caches under `/tmp`, never in this repository.
- If results look suspicious or the change has broad impact, recommend the
  wider relevant test suite before handoff. Source: DI-venit.

## Mogent Rules

- The manifest owns the rendered document outline and headings. Source headings
  guide selection but do not control the output. Source: DI-sufok.
- Use explicit, user-defined source aliases such as `cdint`, `cclab`, or
  `personal`. Do not rely on an ambient library search path. Source: DI-sufok.
- Source subtrees may be composed in manifest order. Warn before confirmation
  when their descendant paths overlap, and let the user combine or separate
  them. Source: DI-sufok.
- Show selection state, tree structure, collapsed state, and source provenance
  in text. Color may reinforce those states but must not be the only signal.
  Source: DI-sufok.
- Treat unreadable sources, duplicate YAML keys, duplicate source aliases,
  unresolved nodes, and empty output as errors. Never silently choose a value.

## Safety And Git

- Do not commit secrets, credentials, signing keys, generated binaries, local
  state, caches, or other runtime artifacts.
- Do not force-push, open a pull request, or make an external change unless the
  user explicitly asks.
- Stage files explicitly. Use short, imperative, capitalized commit subjects
  and summarize non-trivial changes by file in the commit body.
- A user message containing only `commit` authorizes staging and committing the
  current requested changes with an AGENTS-compliant message.

## Handoff

For routine work, report the changed files and checks run. For escalated,
decision-first work, also report decision compliance, the decision-to-evidence
matrix, runtime paths, and user-approved exceptions. Source: DI-venit.

Use clear, direct language. Explain a decision with a concrete example when it
would otherwise be hard to understand.
