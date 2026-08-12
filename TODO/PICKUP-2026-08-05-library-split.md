# PICKUP: Library Split Checkpoint

Date: 2026-08-05

## Current Git State

Latest commit:

```text
fb8f007 Split libraries into atomic modules
```

`git status --short` was clean immediately after the commit.

Ignored local-only paths may exist and should stay uncommitted:

- `.mogent/`
- `docs/other_repo_agents/`
- `testing-ground/*/.mogent/`

`docs/other_repo_agents/` is intentionally ignored because it may contain
private or proprietary material.

## What Was Just Done

The initial checked-in libraries were split from a few large Markdown files into
atomic metadata-bearing modules.

Old large files removed:

- `libraries/cdint/go.md`
- `libraries/cdint/process.md`
- `libraries/cdint/promisegrid.md`
- `libraries/cdint/shared-baseline.md`
- `libraries/nix/system.md`
- `libraries/ucd_research/python.md`
- `libraries/ucd_research/research.md`

New module directories:

- `libraries/cdint/shared-baseline/`
- `libraries/cdint/process/`
- `libraries/cdint/go/`
- `libraries/cdint/cdint-and-promisegrid/`
- `libraries/nix/nix/`
- `libraries/ucd_research/python/`
- `libraries/ucd_research/ucd-research/`

Each module generally has one top-level heading plus YAML frontmatter:

- `tags`
- `tldr`
- `priority`
- `scope`
- `requires` where useful

## Important Design Detail

Directory paths contribute to source references. Filenames do not.

Example:

```text
libraries/cdint/shared-baseline/instructions/focused-change-loop.md
```

with heading:

```markdown
# Focused Change Loop
```

resolves as:

```text
shared:shared-baseline/instructions/focused-change-loop
```

Mogent does not currently stitch separate atomic files into one implicit
cross-file subtree. If a document wants several atomic modules under one local
heading, group them explicitly in `agents.yaml`.

Example from `testing-ground/foss-fork/agents.yaml`:

```yaml
- Decision First:
    - Decision Intent: shared:process/decision-first/decision-intent
    - Thought Experiment: shared:process/decision-first/thought-experiment
    - Open Questions: shared:process/decision-first/open-questions
```

## Validation Already Run

These passed before the commit:

```sh
env GOCACHE=/tmp/mogent-gocache GOMODCACHE=/tmp/mogent-gomod go test ./...
env GOCACHE=/tmp/mogent-gocache GOMODCACHE=/tmp/mogent-gomod go vet ./...
env GOCACHE=/tmp/mogent-gocache GOMODCACHE=/tmp/mogent-gomod errcheck ./...
env GOCACHE=/tmp/mogent-gocache GOMODCACHE=/tmp/mogent-gomod go build -o /tmp/mogent ./cmd/mogent
git diff --check
```

Smoke tests:

```sh
/tmp/mogent source list --manifest testing-ground/basic-org/agents.yaml --tag-search workflow --sort priority
/tmp/mogent build --manifest testing-ground/foss-fork/agents.yaml --force
```

The forced build generated ignored `.mogent` state under testing ground folders;
that state is intentionally ignored.

## Suggested Next Work

1. Compare module overlap and tag consistency.
   Use `mogent source list --metadata --sort priority` against the testing
   manifests and look for near-duplicates, weak TLDRs, or inconsistent tag
   families.

2. Decide whether explicit grouping remains the standard.
   Current behavior is explicit manifest grouping for multi-module sections.
   A future alternative is cross-file subtree stitching, but that could hide
   document structure and needs a deliberate design decision.

3. Add more library content.
   Candidate areas already discussed:
   tutor/teacher mode, TTS-friendly communication, architecture laws, strict
   testing, commit cadence, docs/session logs, developer involvement levels,
   Claude/Gemini/Codex-specific output modules, and UCD research workflows.

4. Design multi-output support.
   Future output examples include `CLAUDE.md`, `GEMINI.md`, and
   `.codex/AGENTS.md`. These can be symlinks or tool-specific rendered files,
   but tool-specific differences should stay visible in YAML.

5. Keep sensitive corpus material out of commits.
   Do not add `docs/other_repo_agents/` unless the user explicitly reviews and
   approves it for git history.
