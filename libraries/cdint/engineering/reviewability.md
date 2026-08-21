---
tldr: Keep changes small, inspectable, and supported by concrete evidence.
---
# Reviewability

## Diff Discipline

Keep changes directly tied to the user's request or a locked decision. Do not
rewrite files from scratch, arbitrarily rewrap lines, reorder unrelated
sections, or normalize prose style unless that cleanup is the requested change.

## Change Size

Prefer small, reviewable stages for complex work. If a non-mechanical change is
large, identify the smallest coherent stage that can land first and explain what
remains.

## User Edits

When user edits appear during a task, preserve them. Do not revert, overwrite,
stage, unstage, commit, or clean up user edits unless the user explicitly asks.
If those edits conflict with the current task, stop and ask how to proceed.

## Breaking Surfaces

Search for and call out changes to integration surfaces: public APIs, protocol
or wire formats, CLI arguments, config loading, persisted data, resuming
existing state, and user-visible output.

## Final Inspection

Read the actual diff before handoff. Report the changed files, checks run, and
known residual risk. Include before/after output for user-visible behavior when
that makes review easier.

## Local Size Policies

Some repos set soft change-size targets, such as keeping complex diffs under a
few hundred lines. Others rely on explicit decision records instead of numeric
limits. Treat numbers as repo-local policy; keep the reusable rule focused on
reviewable stages and visible consequences.
