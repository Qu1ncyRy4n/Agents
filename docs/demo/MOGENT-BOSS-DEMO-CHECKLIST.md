# Mogent Boss Demo Checklist

Use this checklist with [`MOGENT-BOSS-DEMO.md`](MOGENT-BOSS-DEMO.md). The demo
uses a disposable clone of the todo app and never needs to modify its working
repository.

## One Day Before

- Confirm the QMR remote is reachable:

  ```sh
  git ls-remote https://github.com/Qu1ncyRy4n/qmr-agents-library.git HEAD
  ```

- Install the current Mogent checkout:

  ```sh
  cd /home/qix/dev/cdint/Agents
  tools/install
  mogent help
  ```

- Confirm `mogent help` lists `build [--manifest agents.yaml] [--dry-run]`.
- Run the pinned quality gate:

  ```sh
  nix develop -c tools/check
  ```

- Confirm the real repositories are understood before preparing the clone:

  ```sh
  git -C /home/qix/dev/cdint/Agents status --short
  git -C /home/qix/dev/exocortex/todo_app_project status --short
  git -C /home/qix/dev/cdint/qmr-agents-library status --short
  ```

## Prepare a Fresh Demo

The preparation helper refuses to replace an existing target. If a prior demo
directory exists, inspect the path and remove it explicitly before continuing.

```sh
cd /home/qix/dev/cdint/Agents
tools/prepare-boss-demo

export MOGENT_REPO=/home/qix/dev/cdint/Agents
export QMR_LIBRARY=/home/qix/dev/cdint/qmr-agents-library
export DEMO_ROOT=/tmp/mogent-boss-demo/todo_app_project
export OUTLINE_MANIFEST="$DEMO_ROOT/demo-outline/agents.yaml"
export REMOTE_MANIFEST="$DEMO_ROOT/demo-remote/agents.yaml"
```

Expected starting state:

- The root todo workspace has no `.mogent/state.json`.
- Root `AGENTS.md` contains one demo-only accidental comment.
- Root `.agents/skills/rust-development/SKILL.md` contains one demo-only edit.
- `demo-outline` has a clean authored-outline build.
- `demo-remote` has a clean local fixture build and no remote declaration yet.

Verify:

```sh
git -C "$DEMO_ROOT" status --short
mogent status --manifest "$DEMO_ROOT/agents.yaml"
mogent status --manifest "$OUTLINE_MANIFEST"
mogent status --manifest "$REMOTE_MANIFEST"
```

## Terminal Layout

- Terminal 1: live commands in `$DEMO_ROOT`.
- Terminal 2: this speaker deck, kept off screen.
- Editor: open only files named by the current slide.
- Browser: optional; no browser step is required.
- Increase terminal height enough to show source trees without paging.
- Disable shell commands or plugins that automatically modify Git repositories.

## Rehearsal Gates

Confirm these transitions work exactly once before the presentation:

1. Root status begins `untracked`.
2. Normal root build refuses overwrite.
3. `build --force` restores clean root outputs.
4. Manifest exclusion changes root status to `stale`.
5. Skill edit changes only the directory output to `direct edits`.
6. `source add` and `add --rebuild` work in `demo-outline`.
7. Localizing `Instructions/Security` creates local Markdown and provenance.
8. Importing one edit under `Instructions/Focused Change Loop` succeeds.
9. Remote source pin with `--subdir agents` writes
   `demo-remote/mogent.lock.yaml` and source browsing succeeds.
10. All three workspaces finish clean according to `mogent status`.

## Fast Route if Time Is Short

Keep these slides:

1. Mogent Today.
2. Architecture and Code Quality.
3. Todo App as a Real Consumer.
4. Messy Existing Repository.
5. Inspect Before Writing.
6. Restore a Reviewed Baseline.
7. Explore Sources.
8. Generated-Output Safety.
9. Add and Localize Content.
10. Remote Git Pinning.
11. Four-step interface evolution, summarized in one minute.
12. Roadmap and closing takeaway.

Skip:

- Running `tools/check` live.
- Display customization.
- The source-update preview.
- Detailed library-sidecar scan output.

## Failure Recovery

If a root command behaves unexpectedly:

```sh
git -C "$DEMO_ROOT" diff
mogent status --manifest "$DEMO_ROOT/agents.yaml"
```

If the canonical root manifest is malformed:

```sh
git -C "$DEMO_ROOT" restore agents.yaml
mogent build --manifest "$DEMO_ROOT/agents.yaml" --force
```

If the authored-outline sequence was already completed, prepare a fresh demo
rather than manually deleting localization/provenance files during the meeting.

If GitHub is unavailable:

- Do not troubleshoot credentials live.
- State that `git ls-remote` and `source pin` passed in rehearsal.
- Show the pin command and explain the lock/cache outputs from the deck.
- Continue to the interface-evolution section.

## Cleanup

After confirming the path is exactly the disposable demo directory:

```sh
rm -rf /tmp/mogent-boss-demo
```

This cleanup is optional. Nothing in the real todo app, QMR library, or Mogent
repository should have been changed by the demo workflow.
