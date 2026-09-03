# nix-dotfiles Agent Prompt Assessment

## Source Files Found

No `AGENTS.md`, `CLAUDE.md`, `GEMINI.md`, `.cursorrules`, `.windsurfrules`, or GitHub Copilot instruction file was found in this repo.

Useful repo context was inferred from:

- `README.md`
- `flake.nix`
- `.claude/settings.local.json`

## Qualitative Assessment

This repo has strong human-facing documentation. The README explains the flake shape, host matrix, feature flags, common workflows, module layout, and operational commands. That gives an agent enough project context to avoid many obvious mistakes.

What is missing is an explicit agent contract. The repo is a system configuration repo, so the consequences of incorrect edits are higher than normal application code: an agent could break boot, networking, display manager behavior, home-manager activation, or a user environment. The README documents how the system works, but it does not clearly say what an automated coding agent must avoid, how to validate changes, or which host-specific boundaries are sacred.

The repo appears mature enough to benefit from a short, strict `AGENTS.md`. The most valuable addition would not be more architectural explanation; the README already covers that. The missing layer is operational discipline.

## Strengths

- Clear top-level description: NixOS desktop, ThinkPad, WSL, and macOS support share one flake.
- Good directory map with module ownership and host-specific notes.
- Feature flags are documented, which helps agents avoid hard-coding host behavior.
- Day-to-day commands are concrete and copyable.
- The flake itself is reasonably organized around `mkHost`, host modules, home-manager roots, and templates.

## Risks For Agents

- No explicit warning to avoid applying host rebuilds without user confirmation.
- No stated validation ladder, such as `nix flake check`, host-specific dry builds, or `home-manager` checks.
- No guardrails around secrets, hostnames, usernames, hardware configs, VPN/firewall, storage mounts, bootloader, or focus-blocking modules.
- No instruction to preserve host flag boundaries when editing shared modules.
- No clear rule for how to handle lockfile updates.
- `.claude/settings.local.json` only permits a web fetch domain; it does not convey repo workflow or safety constraints.

## Recommended Additions

Add `AGENTS.md` at repo root with:

1. A short project summary: "Declarative NixOS, WSL, and nix-darwin configuration; prefer minimal, host-scoped changes."
2. A safety section:
   - Do not run `sudo nixos-rebuild`, `darwin-rebuild switch`, destructive storage commands, VPN/firewall changes, or bootloader-changing commands without explicit user approval.
   - Do not edit `hardware-configuration.nix` unless the user asks.
   - Do not introduce secrets into tracked files.
3. Validation commands:
   - `nix flake check`
   - `nix build .#nixosConfigurations.<host>.config.system.build.toplevel` for touched NixOS hosts
   - `nix build .#darwinConfigurations.Q-MBP.system` for macOS changes, if supported
4. Change-boundary rules:
   - Shared modules must stay host-flagged.
   - Host-specific behavior belongs under `hosts/<host>/`.
   - Home-manager changes shared by all machines go through `home/base.nix` or the relevant `home/dev/*.nix`.
5. Lockfile policy:
   - Do not update `flake.lock` unless dependency updates are the task.
   - Summarize lockfile diffs if it changes.
6. Style guidance:
   - Prefer existing module/flag patterns.
   - Keep packages in the appropriate package module instead of adding them directly to top-level files.
   - Use `lib.mkDefault`, `mkEnableOption`, and existing option names consistently.

## Suggested Skeleton

```markdown
# AGENTS.md

This repo is Q's declarative system configuration for NixOS desktop, ThinkPad, WSL, and macOS.

## Safety

- Do not run system activation commands without explicit approval.
- Do not edit hardware configs, bootloader, storage mounts, firewall, VPN, or secrets unless requested.
- Do not commit secrets or machine-local credentials.

## Workflow

- Read `README.md` and `flake.nix` before structural edits.
- Preserve host flag boundaries.
- Keep shared behavior in shared modules and host-specific behavior under `hosts/<host>/`.

## Validation

- Run `nix flake check` when possible.
- For NixOS host changes, build the touched host toplevel.
- For macOS changes, build the darwin configuration when possible.

## Style

- Follow the existing module layout.
- Prefer minimal, reviewable Nix changes.
- Do not update `flake.lock` unless dependency updates are in scope.
```
