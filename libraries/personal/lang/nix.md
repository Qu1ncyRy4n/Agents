---
tldr: Use non-activating Nix checks first and protect machine-level configuration boundaries.
---
# Nix

## Safety

Do not run an activating system command without explicit user approval. Do not
change hardware configuration, bootloader, storage, firewall, VPN, DNS,
networking, shell startup behavior, or secrets unless the request explicitly
covers that area.

## Scope

Read the repository's host layout, module layout, overlays, and feature flags
before structural edits. Keep shared behavior behind existing host boundaries.
Put host-specific behavior in the appropriate host configuration rather than
hard-coding it into a shared module.

## Flakes

When the repository uses flakes, inspect `flake.nix`, `flake.lock`, and the
affected module path before editing. Do not update `flake.lock` unless the task
includes dependency updates or the check/build requires a lock refresh. Summarize
lockfile changes when they happen.

## nix develop

Prefer `nix develop` or the repository's documented development shell before
debugging missing tools. Do not assume the ambient shell represents the intended
toolchain.

## Validation

Prefer non-activating checks and builds before any activation. Run `nix flake
check` when the repository supports it. Build the affected host, package, or
home configuration when practical.

## Selection Boundary <!-- id: corpus-variants -->

Distinguish machine configuration from projects that use Nix only as a
development environment. Select host-activation safeguards when work can affect
the active machine. A project that only uses `nix develop` for tools needs the
environment and validation rules, not unrelated host policy.
