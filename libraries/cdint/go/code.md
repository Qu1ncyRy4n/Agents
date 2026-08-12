---
tags: [lang/go, code/style, errors/explicit, scope/project]
tldr: Prefer small Go packages, explicit dependencies, and wrapped errors.
priority: 0.8
scope: project
---
# Code

Prefer small, readable packages and explicit dependencies. Keep package names
short and lower-case. Handle errors explicitly and use `%w` when adding useful
context. Do not ignore errors or hide shell-command failures with `|| true`.
