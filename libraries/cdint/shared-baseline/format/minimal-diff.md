---
tags: [format/diff, workflow/scope-control, scope/repo]
tldr: Keep diffs focused and avoid unrelated reorganization.
priority: 0.8
scope: repo
---
# Minimal Diff

Keep changes tied to the request or a locked decision. Use `git mv` for moves
and renames. Do not reorganize files or normalize unrelated formatting without
an explicit reason.
