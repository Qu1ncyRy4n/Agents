# Mogent Dogfood Session 2

This session resumes after init, local-source selection, localization, and basic
drift handling were exercised in
`/home/qix/dev/omnicortex/todo_app_project`. It is safe to run while URL-source
subdirectory support is developed separately: use the installed stable binary
until the guide explicitly says to install the feature build.

Start without shell variables:

```sh
cd /home/qix/dev/omnicortex/todo_app_project
command -v mogent
mogent status
mogent drift
git status --short
```

Expected baseline: Mogent reports `up to date`; Git may still show the manifest,
generated output, local library, provenance, and state changes from Session 1.
Do not discard those files merely to make Git status empty.

## A. Source Shape And Hierarchy

Compare the source inventory with the actual atomic Markdown layout:

```sh
mogent source list shared
mogent source list personal
mogent coverage --tree --depth 2
mogent coverage --tree --depth 4
```

Then inspect representative files in the Mogent checkout:

```sh
find /home/qix/dev/cdint/Agents/libraries/cdint -type f -name '*.md' | sort
find /home/qix/dev/cdint/Agents/libraries/personal -type f -name '*.md' | sort
```

Record where CLI nesting is shallower than the file-directory and Markdown
heading hierarchy, and whether the lost levels are organizational directories,
file roots, or descendant headings. The desired view should expose real source
shape without inventing selectable nodes that do not exist.

Observed baseline: file roots and their descendant headings render correctly,
but category directories such as `communication`, `domain`, `engineering`, and
`lang` disappear. The intended follow-up view inserts those as visibly
non-selectable groups with aggregate state. For example:

```text
personal
|-- communication/  [group, partial]
|   |-- Accessibility  [included]
|   `-- Personas       [unused]
`-- lang/           [group, partial]
    |-- Python         [unused]
    `-- Rust           [included]
```

## B. Stable Dry-Run Authoring

Exercise suggestions and previews without writing:

```sh
mogent source list shared --search error --tldr
mogent source show shared:engineering/error-handling --content=snippet
mogent add shared:engineering/error-handling \
  --under Instructions --dry-run --preview=tree
mogent add shared:engineering/error-handling \
  --under Instructions --dry-run --preview=patch
```

Expected: typo/descendant suggestions are actionable; both previews describe
the same insertion; no target files change.

## C. URL Subdirectory Pinning

Run this section only after the feature build has landed and been installed.
Keep existing local aliases so the current generated document remains stable.
Add a temporary, unused source to `agents.yaml`:

```yaml
sources:
  remote_shared:
    location: https://github.com/Qu1ncyRy4n/Agents.git
    subdir: libraries/cdint
```

Then pin and inspect it:

```sh
mogent source pin remote_shared
mogent source list remote_shared --tldr
mogent source show remote_shared:engineering/test-strategy --content=snippet
mogent coverage --source remote_shared --tree --depth 3
```

Expected: `mogent.lock.yaml` records URL, full commit, subdirectory, and content
hash. Browsing indexes only `libraries/cdint`, not repository documentation or
other libraries. Ordinary commands do not fetch after pinning.

Temporarily disconnecting the network should not change these results:

```sh
mogent source list remote_shared --search test --tldr
mogent status
```

## D. Report

```text
Mogent commit:
Stage:
Command:
Expected:
Actual:
Hierarchy level lost or invented:
Files changed:
Suggested next hint:
```
