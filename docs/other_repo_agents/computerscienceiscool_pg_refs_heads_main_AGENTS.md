# Repository Guidelines

## Project Structure & Module Organization
- This repository is a standalone static PromiseGrid learning site. Keep user-facing pages at the repository root unless a shared asset or tool clearly belongs in a purpose-named directory.
- Keep shared browser assets under `shared/`.
- Keep local developer tools under `tools/`.
- Keep planning artifacts under `TODO/` and maintain `TODO/TODO.md` as the priority-sorted index.
- Do not commit local state files, generated binaries, browser caches, or scratch output.

## Proquint Coordination IDs
- New TODO files use `TODO/TODO-<handle>-<slug>.md`, where `<handle>` is a proquint-1 minted by `tools/mint-handle`.
- New Decision Intent entries use `ID: DI-<handle>`, where `<handle>` is minted by `tools/mint-handle`.
- Do not create new numbered TODO IDs or timestamp-shaped DI IDs. Existing historical references may remain as history, but new coordination artifacts use proquint handles.
- `tools/mint-handle` is the only supported minting path. It scans the working tree for occupied handle-looking names and exact DI owner lines before choosing a fresh handle.

## Decision-First Specification and Compliance Protocol
- Collect and lock user decisions before behavior-changing code edits.
- Record locked decisions in a `## Decision Intent Log` near the top of the relevant `TODO/TODO-<handle>-<slug>.md` file.
- DI entries are append-only. If intent changes, add a new DI entry with `Supersedes: <old-di-id>` instead of rewriting the old entry.
- Required DI fields:
  - `ID: DI-<handle>`
  - `Date: YYYY-MM-DD HH:MM:SS`
  - `Status: active|superseded`
  - `Decision:`
  - `Intent:`
  - `Constraints:`
  - `Affects:`
  - `Supersedes:` when applicable

## Comment Preservation Protocol
- Do not remove existing code comments unless the same patch replaces them with equal-or-better explanatory comments near the same logic.
- Add short plain-English comments for non-obvious logic, especially parser, encoding, or dynamic rendering code.
- Non-trivial behavior changes must include a behavior-level comment with `Intent:` and `Source: DI-<handle>`.
- After editing, run a comment-delta audit on touched code files:
  `git diff -U0 -- <file> | rg -n '^-\s*//|^-\s*/\*|^\+\s*//|^\+\s*/\*'`

## Coding Style
- Prefer small, focused edits and preserve the standalone-file architecture unless a shared asset removes real duplication.
- Use standard browser APIs and plain JavaScript for static-site behavior.
- Go code must remain `gofmt` formatted and use explicit error handling.
- Use `git mv` for renames.

## Testing
- For `tools/mint-handle`, run `go test ./...` from `tools/mint-handle/`.
- For static pages, open files directly in a browser and verify that relative links, shared scripts, and dynamic rendering work under `file://`.

## Commit Guidelines
- Stage files explicitly; do not use `git add .` or `git add -A`.
- Commit messages should be short, imperative, and capitalized.
- Commit bodies should summarize changes per file.
