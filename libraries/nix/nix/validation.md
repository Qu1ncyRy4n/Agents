---
tags: [nix/validation, testing/non-activating, scope/project]
tldr: Prefer non-activating Nix checks and summarize lockfile changes.
priority: 0.8
scope: project
requires: [nix:nix/safety]
---
# Validation

Prefer non-activating checks and builds before any activation. Run `nix flake
check` when the repository supports it. Build the affected host configuration
when practical. Do not update `flake.lock` unless dependency updates are part of
the requested change; summarize lockfile changes when it is updated.
