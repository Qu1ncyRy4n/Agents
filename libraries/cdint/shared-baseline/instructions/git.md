---
tags: [workflow/git, safety/external-change, scope/repo]
tldr: Stage explicitly, use clear commits, and avoid external changes unless asked.
priority: 0.8
scope: repo
---
# Git

Stage files explicitly. Use short, imperative, capitalized commit subjects.
Imperative means naming the change as an action, such as `Add parser validation`,
rather than describing history with `Added parser validation`.
Summarize non-trivial changes by file in the commit body. Do not force-push,
open a pull request, or make an external change unless the user explicitly asks.
