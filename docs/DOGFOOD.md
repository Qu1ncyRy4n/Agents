# Mogent Dogfooding Guide

This guide exercises Mogent in stages. Stop at the first stage that behaves
unexpectedly and record the command, output, manifest, and expected behavior.
Do not work around a failure by editing generated `AGENTS.md` unless the stage
specifically tests drift handling.

Run commands from the repository root. Install the current checkout with:

```sh
go install ./cmd/mogent
```

Create a disposable copy of the fixture and its source libraries:

```sh
mogent_dogfood_dir=$(mktemp -d)
mkdir -p "$mogent_dogfood_dir/testing-ground"
cp -R testing-ground/personal-go-nix "$mogent_dogfood_dir/testing-ground/"
cp -R libraries "$mogent_dogfood_dir/"
mogent_manifest="$mogent_dogfood_dir/testing-ground/personal-go-nix/agents.yaml"
```

The relative source paths still resolve in this copied layout. Remove the
temporary directory when the exercise is complete.

## Stage -1: Guided Init

Purpose: create a reviewable starter manifest instead of beginning from blank
YAML.

```sh
mogent init --list-templates
mogent init --template personal-go-nix \
  --source shared="$mogent_dogfood_dir/libraries/cdint" \
  --source go="$mogent_dogfood_dir/libraries/cdint" \
  --source personal="$mogent_dogfood_dir/libraries/personal" \
  --manifest "$mogent_dogfood_dir/starter-agents.yaml" \
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

Record issues about YAML validation, source paths, templates, output content,
or overwrite protection.

## Stage 1: Discover Modules

Purpose: find useful material without reading every source file.

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

Use a disposable manifest with one HTTP(S) Git URL source, then run:

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
