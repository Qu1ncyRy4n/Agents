---
tags: [safety/defaults, constraints/secrets, scope/repo]
tldr: Do not commit sensitive artifacts or overwrite work without explicit approval.
priority: 1.0
scope: repo
---
# Safe Defaults

Do not commit secrets, credentials, signing keys, generated binaries, local
state, caches, or other runtime artifacts. Do not delete or overwrite work
unless the user explicitly requests that action.
