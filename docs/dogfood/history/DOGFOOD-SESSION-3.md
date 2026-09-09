# Mogent Dogfood Session 3: Directory Trees And Display

Everything in this session is new since Session 2. It exercises selectable
directory subtrees, the unified inventory/coverage tree, alignment, and ASCII
versus Unicode display. Run it in the existing dogfood repository after
installing the current feature checkout.

```sh
cd /home/qix/dev/cdint/Agents
tools/install
cd /home/qix/dev/omnicortex/todo_app_project
```

## A. Compare Inventory And Coverage

```sh
mogent source list personal --tree --depth 3
mogent source list personal --tree --coverage --depth 3 --tldr
mogent coverage personal --depth 3 --tldr
```

The last two commands should be the same view. The key distinguishes physical
directories (`/`), Markdown headings (`#`), and paths that are both (`/#`).
Category directories such as `communication/`, `engineering/`, and `lang/`
should no longer disappear.

## B. Select A Directory Subtree Safely

Preview adding a whole category beneath an existing manifest heading:

```sh
mogent add personal:engineering --under Instructions --heading Engineering --dry-run --preview=tree
mogent add personal:engineering --under Instructions --heading Engineering --dry-run --preview=patch
```

If you want to exercise rendering without changing the active manifest, copy
`agents.yaml` to a temporary filename, add an entry like the following, and
build using `--manifest`:

```yaml
- heading: Engineering
  from: [personal:engineering]
  exclude: [personal:engineering/staged-migration]
```

Expected: every non-excluded descendant is rendered in deterministic source
path order. Selecting both a directory and one of its descendants in the same
entry is rejected as overlapping.

The original Session 3 preview exposed only the manifest insertion. The current
checkout resolves that gap: verify the `Inherited source subtree` and, when
applicable, `Related existing selections` sections.

Exercise relative placement without writing:

```sh
mogent add personal:engineering --under Instructions --first --heading Engineering --dry-run --preview=tree
mogent add personal:engineering --before Constraints --heading Engineering --dry-run --preview=tree
mogent add personal:engineering --after Instructions/Test Strategy --heading Engineering --dry-run --preview=tree
```

## C. Exercise Display Preferences

One-command overrides:

```sh
mogent source list personal --tree --depth 2 --chars ascii
mogent source list personal --tree --depth 2 --chars unicode
mogent source list personal --tree --depth 2 --align=false
mogent source list personal --tree --depth 2 --tldr --fit term
mogent source list personal --tree --depth 2 --tldr --width 88
```

Optional personal default at `~/.config/mogent/config.yaml`:

```yaml
display:
  chars: unicode
  align: true
  fit: term
  width: 0
```

Explicit flags override the file. These preferences must not modify
`agents.yaml`, `mogent.lock.yaml`, or generated output.

## D. Report

```text
Mogent commit:
Command:
Expected:
Actual:
Directory/heading marker wrong:
Alignment problem:
Files changed unexpectedly:
Suggested next hint:
```
