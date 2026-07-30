# Nix

## Safety

Do not run an activating system command without explicit user approval. Do not
change hardware configuration, bootloader, storage, firewall, VPN, DNS,
networking, shell startup behavior, or secrets unless the request explicitly
covers that area.

## Scope

Read the repository's host layout and feature flags before structural edits.
Keep shared behavior behind existing host boundaries. Put host-specific behavior
in the appropriate host configuration rather than hard-coding it into a shared
module.

## Validation

Prefer non-activating checks and builds before any activation. Run `nix flake
check` when the repository supports it. Build the affected host configuration
when practical. Do not update `flake.lock` unless dependency updates are part of
the requested change; summarize lockfile changes when it is updated.
