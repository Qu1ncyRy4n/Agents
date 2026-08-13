# Mogent Dogfooding Guide

This guide exercises Mogent in stages. Stop at the first stage that behaves
unexpectedly and record the command, output, manifest, and expected behavior.
Do not work around a failure by editing generated `AGENTS.md` unless the stage
specifically tests drift handling.

For only the newly implemented directory-tree, unified-coverage, alignment,
and display-character work, start with
[`DOGFOOD-SESSION-3.md`](DOGFOOD-SESSION-3.md). The longer staged guide below
also includes older regression exercises.

## Current Dogfood Target

The current real-project dogfood target is:

```text
/home/qix/dev/omnicortex/todo_app_project
```

Use these variables throughout the guide:

```sh
mogent_repo=/home/qix/dev/cdint/Agents
mogent_target=/home/qix/dev/omnicortex/todo_app_project
mogent_manifest="$mogent_target/agents.yaml"
```

If the terminal restarts, the simpler recovery path is to enter the target and
use Mogent's default manifest name; no variables are required:

```sh
cd /home/qix/dev/omnicortex/todo_app_project
mogent status
```

Target-specific examples should ultimately use this form. Variables below are
retained only where the same command is intentionally switched between the real
project and optional disposable fixture during this dogfood pass.

Install the current checkout from the Mogent repository:

```sh
(cd "$mogent_repo" && tools/install)
```

Before allowing Mogent to write in the target, check its Git state and inspect
any existing `AGENTS.md`. Do not use `--force` in the real project merely to get
past an unexpected status:

```sh
git -C "$mogent_target" status --short
test ! -e "$mogent_target/AGENTS.md" || sed -n '1,120p' "$mogent_target/AGENTS.md"
```

Start with **Stage -1** against the real target. The first command is a dry run.
Create `agents.yaml` only after its preview looks appropriate, then continue
through the read-only discovery and coverage stages. Stage 4 is the first point
that intentionally rebuilds `AGENTS.md`.

The disposable fixture below remains useful for testing mutation, localization,
drift, and rejection behavior without risking the real project. Commands in the
numbered stages use `$mogent_manifest`, so confirm whether it points to the real
target or disposable fixture before each writing command.

## Optional Disposable Fixture — Skip For The Current Real-Project Run

You do **not** need this section to dogfood Mogent on
`/home/qix/dev/omnicortex/todo_app_project`. Continue directly to Stage -1 with
`mogent_manifest` still pointing at that project's `agents.yaml`.

This optional block creates a temporary copy of Mogent's bundled example and
libraries under `/tmp`. It is useful later when you want to test destructive or
failure behavior without touching the real project. It can be run from any
directory because the copy sources are anchored at `$mogent_repo`.

Create a disposable copy of the fixture and its source libraries:

```sh
mogent_dogfood_dir=$(mktemp -d)
mkdir -p "$mogent_dogfood_dir/testing-ground"
cp -R "$mogent_repo/testing-ground/personal-go-nix" "$mogent_dogfood_dir/testing-ground/"
cp -R "$mogent_repo/libraries" "$mogent_dogfood_dir/"
mogent_manifest="$mogent_dogfood_dir/testing-ground/personal-go-nix/agents.yaml"
```

The last command deliberately changes `mogent_manifest` to the temporary
fixture. Reset it before returning to the real project:

```sh
mogent_manifest="$mogent_target/agents.yaml"
```

The fixture's relative source paths resolve in the copied layout. Remove the
temporary directory when the exercise is complete.

## Stage -1: Guided Init

Purpose: create a reviewable starter manifest instead of beginning from blank
YAML.

For the current real target, begin with the minimal template and the shared
CDINT library. Preview first:

```sh
mogent init --template minimal \
  --source shared="$mogent_repo/libraries/cdint" \
  --manifest "$mogent_manifest" \
  --dry-run
```

If the preview fits the project, create only the manifest (not `AGENTS.md`):

```sh
mogent init --template minimal \
  --source shared="$mogent_repo/libraries/cdint" \
  --manifest "$mogent_manifest"
```

Review the new manifest and target Git diff before proceeding. If the project
is Go-based, also preview the `go` template; if it is a Python research project,
preview `research-python`. Do not select a language template solely from the
project name.

To preview `personal-go-nix` for the current real target, use the real library
paths—not `$mogent_dogfood_dir`, which exists only after creating the optional
disposable fixture:

```sh
mogent init --list-templates
mogent init --template personal-go-nix \
  --source shared="$mogent_repo/libraries/cdint" \
  --source go="$mogent_repo/libraries/cdint" \
  --source personal="$mogent_repo/libraries/personal" \
  --manifest "$mogent_manifest" \
  --dry-run
```

Expected: the template list explains required aliases. Dry-run prints ordinary
editable YAML and writes nothing. Remove `--dry-run` to write the starter; add
`--build` only when its output path is safe.

## Stage 0: Build A Known Manifest

Purpose: verify parsing, source resolution, rendering, and overwrite safety.

```sh
mogent status --manifest "$mogent_manifest"
mogent build --manifest "$mogent_manifest" --force
mogent status --manifest "$mogent_manifest"
```

Expected: the manifest and sources resolve, the output builds, and the final
status is `up-to-date`. If Mogent reports direct edits or an untracked output,
review the diff before considering `--force`. The command above uses `--force`
only because the disposable fixture contains a copied, untracked output.

For `/home/qix/dev/omnicortex/todo_app_project`, omit `--force` on the first
build:

```sh
mogent status --manifest "$mogent_manifest"
mogent build --manifest "$mogent_manifest"
git -C "$mogent_target" diff -- agents.yaml AGENTS.md
mogent status --manifest "$mogent_manifest"
```

If the build refuses because `AGENTS.md` already exists or has direct edits,
stop and review that file rather than retrying with `--force`.

Record issues about YAML validation, source paths, templates, output content,
or overwrite protection.

## Stage 1: Discover Modules

Purpose: find useful material without reading every source file.

`source list` is the searchable inventory of available source nodes and their
metadata. It does not currently have a `--tree` view. `coverage` in Stage 2
looks at the same sources but overlays how the current manifest uses them.

There is not yet a `mogent source add` command for extending an existing
manifest. To add the personal library after initializing with `minimal`, edit
the manifest's existing `sources` mapping:

```yaml
sources:
  shared: /home/qix/dev/cdint/Agents/libraries/cdint
  personal: /home/qix/dev/cdint/Agents/libraries/personal
```

Then verify it before selecting a node:

```sh
mogent status --manifest "$mogent_manifest"
mogent source list personal --manifest "$mogent_manifest"
mogent add personal:lang/python --manifest "$mogent_manifest" \
  --under Instructions --dry-run --preview=patch
```

Do not run `source add ...`: `source` is a shell builtin for loading a shell
script, not the Mogent CLI.

For the current real target initialized with `minimal`:

```sh
mogent source list shared --manifest "$mogent_manifest"
mogent source list shared --manifest "$mogent_manifest" --search test --tldr
mogent source show shared:engineering/test-strategy \
  --manifest "$mogent_manifest" --metadata --content=snippet
```

The disposable personal Go/Nix fixture can additionally exercise:

```sh
mogent source list --manifest "$mogent_manifest"
mogent source list personal --manifest "$mogent_manifest"
mogent source list --manifest "$mogent_manifest" --search nix --tldr
mogent source show personal:lang/nix --manifest "$mogent_manifest" --metadata --content=snippet
```

Expected: source aliases narrow the list, text search finds references and
content, and TLDR output gives a compact selection clue when metadata exists.

Record missing or misleading search results, unhelpful TLDRs, noisy output,
and source paths that are hard to predict.

## Stage 2: Understand Selection And Coverage

Purpose: compare available modules with what the manifest already selects.

Coverage does not simulate importing every source. It classifies each available
node as directly included, inherited through a selected parent, excluded,
partially selected, or unused by the current manifest.

```sh
mogent coverage --manifest "$mogent_manifest"
mogent coverage --manifest "$mogent_manifest" --tree
mogent coverage --manifest "$mogent_manifest" --tree --tldr
mogent coverage --manifest "$mogent_manifest" --unused-only
```

Expected: tree output distinguishes included, inherited, excluded, partial,
and unused nodes in text. Counts should agree with the visible source tree.
`--tldr` shows each file-level summary once rather than repeating it for every
heading that inherits the same metadata.

Record confusing parent/child state, incorrect counts, or modules that appear
selected merely because a related heading is selected.

## Stage 3: Preview A Manifest Change

Purpose: review a durable operation before writing it.

Choose an unused reference from Stage 2, then run:

```sh
mogent add shared:engineering/test-strategy --manifest "$mogent_manifest" --under Instructions --dry-run --preview=summary
mogent add shared:engineering/test-strategy --manifest "$mogent_manifest" --under Instructions --dry-run --preview=tree
mogent add shared:engineering/test-strategy --manifest "$mogent_manifest" --under Instructions --dry-run --preview=patch
```

Those commands fit the current target's `minimal` starter. For the disposable
personal Go/Nix fixture, use:

```sh
mogent add personal:lang/python --manifest "$mogent_manifest" --under Instructions --dry-run --preview=summary
mogent add personal:lang/python --manifest "$mogent_manifest" --under Instructions --dry-run --preview=tree
mogent add personal:lang/python --manifest "$mogent_manifest" --under Instructions --dry-run --preview=patch
```

Expected: no file changes. Each view should describe the same insertion at a
different level of detail. Use a disposable manifest because the example
heading may not exist in every repository.

Record unclear insertion points, source/manifest heading confusion, or preview
views that disagree. For another manifest, choose a heading printed by
`mogent complete manifest-headings --manifest "$mogent_manifest"`.

## Stage 4: Save And Rebuild

Purpose: exercise the complete noninteractive authoring path.

For the current real target, save the change previewed in Stage 3:

```sh
mogent add shared:engineering/test-strategy --manifest "$mogent_manifest" --under Instructions
mogent status --manifest "$mogent_manifest"
mogent build --manifest "$mogent_manifest"
git -C "$mogent_target" diff -- agents.yaml AGENTS.md
mogent status --manifest "$mogent_manifest"
```

For the disposable personal Go/Nix fixture, use:

```sh
mogent add personal:lang/python --manifest "$mogent_manifest" --under Instructions
mogent status --manifest "$mogent_manifest"
mogent build --manifest "$mogent_manifest"
mogent status --manifest "$mogent_manifest"
```

Expected: `add` changes only the manifest, status reports stale output, build
updates the generated file, and status becomes clean. Add `--rebuild` to the
`add` command in a separate trial to exercise the combined operation.

Record partial writes, poor errors, formatting surprises, and any case where a
generated file is overwritten without the documented guard.

## Stage 5: Shell And Agent Use

Purpose: check discoverability for repeated CLI and scripted use.

```sh
mogent complete commands
mogent complete source-refs --manifest "$mogent_manifest" --prefix personal:
mogent complete manifest-headings --manifest "$mogent_manifest"
mogent completion zsh
```

Expected: candidate output is stable, one item per line, and useful without
parsing presentation text. Generated completion scripts should suggest commands,
source aliases, references, flags, and manifest headings.

Record missing candidates, incorrect prefixes, shell errors, or output that is
unsafe or awkward for scripts.

## Stage 6: Localization

Purpose: create an explicit local override without modifying a shared source.

For the current real target's `minimal` manifest, localize the focused change
loop. Keep the dry run as the first attempt:

```sh
mogent localize Instructions/Focused\ Change\ Loop --manifest "$mogent_manifest" --dry-run
mogent localize Instructions/Focused\ Change\ Loop --manifest "$mogent_manifest"
mogent source show local:shared-baseline/instructions/focused-change-loop \
  --manifest "$mogent_manifest" --metadata
```

The disposable personal Go/Nix fixture uses:

```sh
mogent localize Instructions/Go\ Tests --manifest "$mogent_manifest" --dry-run
mogent localize Instructions/Go\ Tests --manifest "$mogent_manifest"
mogent source show local:go/tests --manifest "$mogent_manifest" --metadata
```

Expected: dry-run writes nothing. The real operation creates ordinary Markdown
under `.mogent/library`, records origin data in `.mogent/provenance.yaml`, and
changes only the chosen manifest entry to `local:`. The generated output becomes
stale until rebuilt unless `--rebuild` is supplied.

Record unexpected source changes, path collisions, missing provenance, or a
manifest that points at a nonexistent local file.

## Stage 7: Drift Detection And Import

Purpose: exercise explicit handling of a directly edited generated output.

First build the disposable fixture, edit prose inside exactly one rendered
manifest section, then run:

```sh
mogent drift --manifest "$mogent_manifest"
mogent drift --manifest "$mogent_manifest" --import Instructions/Go\ Tests
```

On the current real target, use `Instructions/Focused Change Loop` in place of
`Instructions/Go Tests`, and inspect the target Git diff before and after the
import. Do not use `--reject --force` there until the direct edit is safely
recorded elsewhere or is intentionally disposable.

Expected: drift reports direct edits. Import succeeds only when every change is
inside the selected unambiguous section; it creates a local override and rebuilds
without losing that edit. An edit elsewhere in the document makes import fail
without writes. To deliberately discard edits instead:

```sh
mogent drift --manifest "$mogent_manifest" --reject --force
```

Promotion back to a trusted shared source remains future work.

## Stage 8: Pinned URL Sources

Purpose: verify that remote content is fetched only by explicit commands and is
then consumed immutably offline.

Use a disposable manifest with one HTTP(S) Git URL source. To exercise a
library within a larger repository, use the explicit source form:

```yaml
sources:
  shared:
    location: https://github.com/Qu1ncyRy4n/Agents.git
    subdir: libraries/cdint
```

Then run:

```sh
mogent source pin <alias> --manifest path/to/agents.yaml
mogent build --manifest path/to/agents.yaml
mogent source update <alias> --manifest path/to/agents.yaml
# after reviewing the printed commit and Markdown content:
mogent source update <alias> --manifest path/to/agents.yaml \
  --ref <full-previewed-commit> --accept
```

Expected: pin writes reviewable `mogent.lock.yaml` and an ignored verified cache.
Build performs no fetch. Update reports old/new commits and old/new Markdown
content but writes nothing until repeated with the full previewed commit and
`--accept`. Disconnecting the network after pinning must not prevent build,
source browsing, or coverage.

Record authentication prompts, mutable content consumption, missing change
paths, integrity failures, or any ordinary command that unexpectedly fetches.

## Issue Report Template

```text
Stage:
Mogent commit:
Command:
Manifest path:
Expected:
Actual:
Output status before/after:
Minimal reproduction:
```
