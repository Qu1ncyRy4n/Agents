# DR: Source Stack Metadata And Coverage

DR-ID: DR-source-stack-metadata-coverage
Date: 2026-08-04
State: open

## Question

Should mogent add a source-stack mode that can include all matching modules
from ordered sources, filtered by metadata and tags, while showing source
coverage and conflicts before writing `AGENTS.md`?

## Why This Is Open

The current manifest model is precise: each document node names the source node
it wants. That is safe and readable, but it becomes tedious for org-wide
libraries where the common case may be "take the org, language, team, project,
repo, and maintainer layers in order."

An include-all source stack would make org adoption easier, but it introduces
harder questions:

- how to mark required, recommended, and optional sections;
- how to include new upstream sections without silently changing repo behavior;
- how to show source content that is available but not included;
- how to handle duplicate paths and semantic conflicts;
- how to let a repo apply a stable adjustment to a remote source it cannot edit;
- how to keep generated context visible instead of depending on global or
  invisible prompt files.

## Current Direction

- Support both complete-document sources and one-module-per-file sources.
- Recommend one-module-per-file for org and shared libraries because it makes
  review, diffing, coverage, and copy-on-write localization easier.
- Keep manifest headings authoritative. Source modules may suggest a title and
  path, but the manifest can place the content under a local heading.
- Add per-module metadata with fields such as `path`, `title`, `summary`,
  `requirement`, `scope`, `tags`, `priority`, and `conflicts_with`.
- Prefer `requirement: mandatory | recommended | optional` for adoption level.
- Add source coverage views that show included, overridden, excluded,
  available-but-unused, and new-since-last-lock modules.
- Treat URL/shared updates as reviewable changes. Fetching an update must not
  silently alter `agents.yaml` or `AGENTS.md`.
- Add an adjustment/patch layer for cases where the user cannot edit the remote
  source and does not want to create a user-wide policy. The adjustment should
  be visible in the manifest or in a manifest-referenced local file.

## Candidate Syntax

Avoid bare YAML `*` syntax because YAML already uses `*` for aliases, and
mogent rejects aliases.

Possible include-all shape:

```yaml
doc:
  - include_all:
      from: [org, go, team, repo]
      where:
        requirement: [mandatory, recommended]
        tags: [security, testing]
      merge: later_sources_override
```

Possible adjustment shape:

```yaml
adjustments:
  - path: ./agents.adjust.yaml

doc:
  - include_all:
      from: [org, go, team, repo]
      where:
        requirement: [mandatory, recommended]
```

The exact adjustment language is not settled. Options include replacement by
source path, structured patch operations, or copy-on-write localization with a
recorded reason.

## Metadata Placement Options

YAML front matter at the beginning of a Markdown file is conventional and easy
to parse, but it only describes the file as a whole. If a single file contains
many headings, mogent needs heading-level metadata.

Possible heading-level metadata:

```md
## Testing
<!-- mogent:
path: instructions/testing
title: Testing
summary: Validation expectations before handoff.
requirement: recommended
tags: [testing, deterministic]
priority: 60
-->

Run relevant deterministic tests before handoff.
```

The heading-level form keeps ordinary Markdown readable, supports complete
document sources, and avoids forcing every shared library into one-section files.
The one-section-per-file form remains simpler for org libraries.

## Conflict Handling

Some conflicts are structural and can be checked deterministically:

- duplicate paths in the same source;
- duplicate source aliases;
- incompatible metadata;
- include-all sources defining the same path under an override policy.

Semantic conflicts need LLM-assisted review with user approval. Generated
`AGENTS.md` files may also include a visible instruction asking future agents to
flag semantic conflicts before acting.

## Affects

- `docs/DESIGN.md`
- future include-all manifest grammar
- source library format
- source browser and coverage UI
- URL source pinning and update review
- copy-on-write localization and adjustment files

## Unblocks

This does not block M3 copy-on-write editing. It should be resolved before
building URL-source scale workflows, include-all source stacks, or large org
library adoption.

## Waiting On

- Dogfooding one-module-per-file libraries versus complete-document libraries.
- A small source browser prototype that shows available-but-unused modules.
- A concrete adjustment-file example.
- Empirical examples from existing org/team/repo `AGENTS.md` files.

## Linked DI

None yet. A Decision Intent should supersede this request once source-stack
metadata and coverage behavior are locked.
