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
mogent coverage --manifest "$mogent_manifest" --unused-only
```

Expected: tree output distinguishes included, inherited, excluded, partial,
and unused nodes in text. Counts should agree with the visible source tree.

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

## Stage 6: Future Editing And Drift

This stage is not implemented yet. It will cover copy-on-write localization,
source-to-local provenance, direct-edit detection and import, reconciliation,
and promotion back to a trusted shared source. Do not treat `--force` as a
substitute for that workflow.

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
