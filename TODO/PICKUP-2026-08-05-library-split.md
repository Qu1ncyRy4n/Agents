# PICKUP: Library V2 Protolibrary

Date updated: 2026-08-21

Current cross-project head: `~/dev/Notes/inbox/2026-08-W34.md`.

## Current Work

The August 5 atomic-library split below is historical context. The active work
is now a separate exhaustive protolibrary under `libv2-proto/`. Do not delete or
broadly rewrite `libraries/` while v2 intake and owner review are in progress.

Current model:

```text
protected source guides
        ↓
PROVENANCE.md file and heading ledger
        ↓
intake/<semantic-area>/ source-grounded candidates
        ↓
owner review and explicit decisions
        ↓
orgs/<owner>/<semantic-area>/ reviewed modules
        ↓
representative manifests and rendered-output validation
        ↓
eventual v1 retirement decision
```

Current status:

- all 37 source-guide files are listed in the file-level ledger;
- heading-level exhaustive mapping is active but incomplete;
- workflow and decision-governance are the first review area;
- the current candidate batch and progress report are intentionally uncommitted
  pending owner review;
- v1 remains a comparison target, not an authoritative source for v2 unless its
  source-guide derivation is recovered;
- generated repository assessments remain visibly distinguished from captured
  authoritative guides.

Read in this order:

1. `docs/LIBV2-PROTO-PROGRESS-REPORT.md`
2. `libv2-proto/WORKLIST.md`
3. `libv2-proto/SPEC.md`
4. `libv2-proto/PROVENANCE.md`
5. candidates under `libv2-proto/intake/workflow/`
6. candidates under `libv2-proto/intake/process/decision-governance/`

## Resume Instructions

1. Wait for owner notes on the current workflow/decision candidate batch.
2. Revise source coverage, wording, module boundaries, and review concerns
   without prematurely resolving strict-versus-lightweight policy.
3. Commit the reviewed batch only after an explicit `good to go, commit` or
   equivalent instruction.
4. Continue heading-level extraction in broadly applicable order: workflow,
   process, safety, engineering, documentation, tools, languages,
   communication, domains, repository-local material, then generated
   candidates.
5. For each source heading, record `represented`, `partial`, `duplicate`,
   `repository-local`, `generated-assessment`, or `rejected-with-reason`.
6. Schedule focused TE/DR decisions only after real extracted examples exist.

## Current Decision Queue

- edit-scope plus user-work preservation: one compatible bundle or two atomic
  modules;
- risk-based escalation versus the complete strict decision-first workflow;
- compatible documents versus separately selectable directory children;
- whether exact mid-document insertion is needed beyond manifest composition,
  then minimal frontmatter-named marker versus restricted Go-template include;
- warning versus error levels for relationships, imports, and migrations;
- sparse source-path move maps versus permanent identity machinery;
- output targets including `AGENTS.md`, `CLAUDE.md`, skills, specifications,
  generated docs, tool-specific guidance, and repository-local references;
- eventual ownership of canonical modules after semantic intake.

Include/composition narrowing has now been completed in
`docs/thought-experiments/TE-nufad-module-includes.md`. Rich YAML embedded in
HTML comments is rejected as the default. The surviving directions are:

1. manifest composition plus separate modules;
2. frontmatter plus a minimal named insertion marker if exact placement is
   required; and
3. a restricted literal Go-template include if one template surface is
   preferred.

Owner decisions remain open. Do not implement an include syntax yet.

## Relevant Commits

- `86780a0 Preserve v1 library review state`
- `1eebacb Start lossless v2 agent protolibrary`
- `895fc96 Define provisional v2 library semantics`
- `9f2f6a4 Reorganize v2 candidates by semantic intake`

---

# Historical Checkpoint: Initial Atomic Library Split

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
