# NixConfig Agent Prompt Assessment

## Source Files Found

No `AGENTS.md`, `CLAUDE.md`, `GEMINI.md`, `.cursorrules`, `.windsurfrules`, or GitHub Copilot instruction file was found in this repo.

Useful repo context was inferred from:

- `README.md`
- `flake.nix`

## Qualitative Assessment

This repo looks like an older or narrower macOS/nix-darwin configuration. The README is mostly a TODO list, and the flake is a single-file nix-darwin plus home-manager setup with several commented sections. Compared with `nix-dotfiles`, this repo gives an agent much less structure to rely on.

The biggest problem is ambiguity. It is not obvious whether this repo is still active, superseded by `nix-dotfiles`, or maintained as a separate macOS-specific config. An agent working here needs an instruction file that establishes status, scope, and migration intent before it starts reorganizing things.

The repo would benefit from an `AGENTS.md` that is more defensive than expansive. The main goal should be preventing accidental large refactors and clarifying whether fixes should be made here or ported to the newer `nix-dotfiles` repo.

## Strengths

- The README captures the user-visible priorities: Yabai, SKHD, blocky, modularization, Linux usability, scheduler behavior, SSH, Git, uninstall/reinstall, and secrets.
- The flake is compact enough for an agent to understand quickly.
- The TODO list gives a rough roadmap.

## Risks For Agents

- No status marker saying whether this repo is active, archival, experimental, or superseded.
- No explicit validation command.
- Large commented blocks make it easy for an agent to "clean up" code that may be preserved context.
- No guidance about whether to modularize in place or migrate into `nix-dotfiles`.
- Uses macOS system settings, Homebrew casks, launchd agents, shell aliases, and user paths; these can affect the active machine directly.
- Secrets/passwords are explicitly mentioned as future work, but there is no rule preventing secret material from being committed.

## Recommended Additions

Add `AGENTS.md` at repo root with:

1. Repo status:
   - State whether this is active, legacy, or a staging area for `nix-dotfiles`.
   - Say whether changes should be ported elsewhere.
2. Safety guardrails:
   - Do not run `darwin-rebuild switch` without approval.
   - Do not alter launchd agents, Homebrew casks, shell startup files, networking, DNS, or secrets without explicit scope.
3. Refactor policy:
   - Do not modularize broadly unless the task is specifically modularization.
   - Preserve commented sections unless the user asks to delete them.
4. Validation:
   - Use `nix flake check` if supported.
   - Prefer `darwin-rebuild build --flake .` or equivalent non-activating build before activation.
5. Migration notes:
   - If `nix-dotfiles` is the successor, compare patterns there before editing this repo.

## Suggested Skeleton

```markdown
# AGENTS.md

This repo contains Q's nix-darwin/home-manager configuration. Confirm whether it is active or legacy before doing broad work.

## Safety

- Do not run `darwin-rebuild switch` without explicit approval.
- Do not change launchd agents, DNS/networking, shell startup behavior, Homebrew app lists, or secrets unless requested.
- Never commit secrets or machine-local credentials.

## Scope

- Prefer small fixes over broad modularization.
- Preserve commented sections unless cleanup is the task.
- If this repo is being migrated into `nix-dotfiles`, check that repo's current structure before editing.

## Validation

- Run `nix flake check` when possible.
- Use a non-activating darwin build before any activation command.
```
