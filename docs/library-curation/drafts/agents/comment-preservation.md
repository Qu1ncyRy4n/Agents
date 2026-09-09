# Preserve Comments And Decision Intent

Current profile: strict default, aligned with current `cdint-grid`. Keep the
complete rule in always-loaded guidance for initial dogfooding. The procedural
subheadings below are marked for possible later skill extraction; they do not
weaken the current rule.

Source: [`cdint-grid/AGENTS.md`](../../cdint-grid/AGENTS.md#comment-preservation-protocol-required)

## Preserve Existing Intent

- Never remove existing code comments unless they are replaced in the same patch by equal-or-better explanatory comments near the same logic.
- When rewriting or refactoring code, port old explanatory intent first, then improve wording.
- If a touched non-trivial code block has no comments, add explanatory comments.
- Do not treat shorter comments as better unless they preserve all important intent.
- If a comment must be dropped with no replacement, stop and ask the user before proceeding.
- Before editing a file, review existing comments in that file.
- Do not remove comments or documentation; update them if outdated or incorrect.

## Record Behavior Intent And Provenance

- For any non-trivial behavior change, include a behavior-level comment with:
  - `Intent:` a short, clear rationale (a sentence or a few; no hard cap if more is needed for clarity).
  - `Source:` a DI ID in the format `DI-<handle>`.
  - `<handle>` is minted by `/home/stevegt/bin/mint-handle` and is globally unique across TODO, TE, DR, DI, and DN owners. Source: DI-jufiz
  - Optional: TODO file/section reference for faster lookup.

## Maintain Decision Intent History

Potential future skill extraction: the record-location and append-only procedure
are operational conventions. They remain required here during initial dogfooding.

- Maintain a `## Decision Intent Log` at the top of relevant root
  `TODO/TODO-<handle>-<slug>.md` files. Source: DI-gigoz; DI-topih
- Treat DI logs as append-only history. Do not rewrite or delete prior entries.
- When intent evolves, add a new DI entry and set `Supersedes: <old-di-id>`.
- DI entries must include:
  - `ID: DI-<handle>`
  - `Date: YYYY-MM-DD HH:MM:SS`
  - `Status: active|superseded`
  - `Decision:`
  - `Intent:`
  - `Constraints:`
  - `Affects:`
  - `Supersedes:` (optional)

## Audit Comment And Intent Changes

Potential future skill extraction: this is a repeatable post-edit procedure. It
remains required here during initial dogfooding.

- After editing, run a comment-delta audit on each touched code file using: `git diff -U0 -- <file> | rg -n '^-\\s*//|^-\\s*/\\*|^\\+\\s*//|^\\+\\s*/\\*'`.
- Resolve all removed-comment lines before finalizing unless explicit user approval was given.
- In the final response, include:
  - `Comment audit: PASS/FAIL`, with file list.
  - `Intent provenance audit: PASS/FAIL`, listing files with behavior changes and DI sources.
- Hard gate: behavior-changing work is incomplete unless comments preserve intent and include DI provenance.
