---
tags: [lang/python, workflow/files, cli/dry-run, scope/project]
tldr: Use dry runs for file moves and keep generated paths explicit.
priority: 0.7
scope: project
---
# File And CLI Work

Use a dry run before changing file names, locations, or a user's document
archive. Preserve idempotence where a tool may run more than once. Keep output
paths explicit and do not treat generated data as source material.
