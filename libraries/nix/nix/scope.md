---
tags: [nix/scope, config/host-boundaries, scope/project]
tldr: Keep Nix changes behind existing host and feature boundaries.
priority: 0.8
scope: project
---
# Scope

Read the repository's host layout and feature flags before structural edits.
Keep shared behavior behind existing host boundaries. Put host-specific behavior
in the appropriate host configuration rather than hard-coding it into a shared
module.
