# Mogent Design (Design of Record)

Status: **active** · Last reconciled: 2026-07-30 · Supersedes: the Gen-1 tag model, the block/id-per-node framing of `TE-tavim`/`DI-lorad`, and the flat select/deselect TOML config of `DI-modun`.

This is the single authoritative description of what mogent is and where it is going.
When any other doc disagrees with this one, this one wins.

---

## 1. Problem

You maintain an `AGENTS.md` in many repos. They share most content — identity,
workflow, constraints, style — but each needs local tweaks. Today every repo
hand-copies and re-edits a large file, so improvements never propagate and the files
drift apart.

## 2. Solution

Two trees, one manifest:

- The **library** is where content lives: Markdown files whose headings form a tree.
- The **document** (`AGENTS.md`) has its own outline — the manifest.
- The **manifest** is the config: the document's outline written down, where each
  node points at a source node in a library. Reading the config *is* reading the
  table of contents of your `AGENTS.md`.

Mogent is **a manifest editor with a renderer** — not a bundle of commands. Author
content once in shared libraries; compose each repo's document by manifest.

Design principles:

- **Markdown-native.** Libraries are normal Markdown. A heading and its content *is*
  the unit — no invented "block" concept.
- **The manifest is the config.** Document structure, order, and provenance are all
  visible in one readable file. No hidden merge logic to mentally replay.
- **Compose, don't mirror.** A one-line subtree reference inherits upstream structure
  dynamically; explicit nesting appears only where you've made a structural choice.
- **Fail loud.** Missing sources, missing nodes, ambiguous references, empty output
  are errors, never silent drops.
- **Deterministic.** Same libraries + manifest ⇒ same `AGENTS.md`. Tool-only metadata
  never appears in rendered output.

---

## 3. Core model

### 3.1 Library: heading trees

- A library is a directory of Markdown files. Files are organizational containers;
  each heading plus its descendant content (up to the next heading of equal-or-higher
  level) is a **node**.
- A node is addressed by its **heading path** — slugified headings joined with `/`:

```
# Instructions        -> instructions
## Testing            -> instructions/testing
### Flaky retries     -> instructions/testing/flaky-retries
```

- Paths are readable and greppable; they break only on rename/move (rare, deliberate,
  find-and-replaceable). A node may opt into a rename-proof anchor with a one-line
  comment — `## Testing  <!-- id: strict-testing -->` — useful for widely-referenced
  nodes in shared libraries. Most nodes need no id.
- Never content hashes as identity: a hash changes on every prose edit — the common
  case — so it breaks references exactly when you improve content. (A hash is fine
  later as *lockfile* integrity data.)

### 3.2 Manifest: the document outline

The manifest (`agents.yaml`) declares the output document as a nested, ordered
outline. Each entry either:

- **references a single node** (`shared:identity/role`),
- **pulls a whole subtree** in one line, optionally with `exclude` for deep nodes to
  drop, or
- **nests further entries**, defining document structure explicitly.

The manifest owns the rendered headings. A source node's heading is a useful
suggestion while selecting content, but a repository may deliberately give the
same content a different place or name in its document. This avoids treating a
library's organizational vocabulary as the output document's vocabulary.

A manifest node may also compose multiple source subtrees in a declared order.
Each source pull remains subtree-inclusive. Before writing output, mogent warns
when composed sources have overlapping descendant paths and lets the user keep
one combined section or place the content in separate, locally named sections.
The exact YAML shape for multi-source composition belongs to milestone 1.
Source: DI-sufok.

Every user operation is a manifest operation:

| Operation | Manifest meaning |
|---|---|
| navigate | walk the outline |
| include / exclude | add / remove an entry (or an `exclude` line under a subtree pull) |
| **swap** | change one entry's source reference |
| edit | copy-on-write: node copies into the local library, reference flips to `local:` |
| reorder | reorder manifest entries — document order *is* manifest order |

Because the manifest composes from libraries rather than mirroring them, upstream
additions inside a pulled subtree flow through automatically, and the manifest stays
small: it only spells out where you've made choices.

### 3.3 Sources

The manifest names its libraries in a `sources` map: a user-defined short name
bound to a local path or a URL. Names such as `cdint`, `cclab`, `personal`, and a
collaborator name are valid; there is no fixed source-type enum. References are
always explicit about provenance (`shared:…`, `grid:…`, `local:…`) — no ambient
search path, no last-wins guessing. A source alias is local to its manifest and
may be deliberately renamed with every reference.

Sources are untrusted input: safe path resolution (nothing outside the library
root), duplicate aliases/keys, and a failed or unreadable source are errors, not
silent skips. Source: DI-sufok.

### 3.4 Templates

Nodes are Go `text/template`s rendered at build:

- tool-provided vars (`repo_name`, `repo_url`, …),
- user vars from the manifest's `vars` section,
- e.g. a shared identity node says `You are working on {{ .repo_name }}` and serves
  every repo unedited.

### 3.5 Editing = copy-on-write localization

Editing never mutates a shared library in place. Editing a node:

1. copies it into the **local** library (under `.mogent/`), preserving its path, with
   your modification applied;
2. flips the manifest reference from the shared source to `local:`.

Shared libraries stay read-only inputs; the local library holds exactly your
overrides; the manifest makes every override visible. Promoting a local override back
up to a shared library is a separate, later action.

---

## 4. Configuration: `agents.yaml`

**YAML is canonical.** Format follows model: a manifest is a nested, ordered,
human-read outline — YAML's home turf, and where TOML cannot stay legible. YAML's
footguns are contained by a **strict schema with loud validation**, quoting values,
and owning every key. `docs/IMPLEMENTATION-M1.md` defines the exact local-source
grammar and validation rules for the first renderer. Source: DI-vukam.

Sketch (shape is settled; exact schema finalizes during the rebuild):

```yaml
sources:
  shared: ~/.mogent/library
  grid:   https://github.com/ciwg/agents

vars:
  project_name: Agents

output: AGENTS.md

doc:
  - identity:
      - role:     shared:identity/role
      - overview: local:identity/overview.md      # copy-on-write override
  - instructions:
      - workflow: shared:instructions/workflow
      - testing:
          from: [shared:instructions/testing]      # subtree pulled in…
          exclude: [shared:instructions/testing/flaky-retries] # drops deep node
  - constraints: shared:constraints                # whole subtree, one line
  - format:
      - coding-style: grid:lang/go/style           # swapped in from another source
```

The navigator and confirmation view always show the document tree, collapsed or
expanded state, and the source of each entry in text. Color distinguishes state
and provenance early, but symbols and source labels remain sufficient in a
monochrome terminal. When a URL source changes before pinning exists, the
confirmation view shows the changed content before it writes `AGENTS.md`.
Source: DI-sufok.

---

## 5. Categories

Categories are conventional top-level manifest headings — content, not machinery.
The default four:

1. **Identity** — agent role, project overview, tech stack, project structure.
2. **Instructions** — workflow, code changes, testing, commits, decision protocol
   (TE/DI/DF), comment preservation, TODO tracking.
3. **Constraints** — never-do / always-do rules, security, runtime hygiene.
4. **Format** — coding style, diff discipline, error format, response/handoff, glossary.

Additional categories we have in mind, added as the library grows:

5. **Cognition / process** — thinking depth, TE/DI/DF behavior, fast-iteration vs
   deliberate research, learning-focused mode.
6. **Communication / style** — directness, Socratic mode, TTS-friendly output, humor,
   simplicity level.
7. **Code** — per-language style, stack rules, testing strategy, commit cadence.
8. **Notes / docs** — README/changelog conventions, dev logs, session notes, human- vs
   LLM-facing docs.

These are the same *role* axis (what kind of content), just finer. Weight (light vs
heavy process), domain (grid), and language (go/rust) are **not** categories — they
are separate libraries in `sources`, and later tags.

---

## 6. Workflow

One interaction model everywhere: **manipulate the manifest → confirm → build.**

`mogent init` (gh-CLI style, the primary flow):

1. **Choose sources** — detected defaults (`~/.mogent/library`, org URL) plus custom.
   Writes `sources:`.
2. **Navigate the tree** — one selector over the merged source trees: toggle
   include/exclude, `v` view a node's text, `e` edit (copy-on-write, reference
   repointed automatically), `s` swap (pick a different source node for this slot).
3. **Confirm read** — show the manifest (optionally the rendered preview). Confirm →
   write `agents.yaml` → build `AGENTS.md`.

Other entry points are the same model:

| Command | Role in the model |
|---|---|
| `mogent tui` | re-enter step 2 on an existing manifest |
| `mogent build` | render the manifest → `AGENTS.md` (validates, fails loud) |
| `mogent status` | show manifest/output paths, generated-output state, and source counts |
| `mogent coverage` | show included and unused source nodes for the manifest; `--unused-only` prints a compact unused list |
| `mogent source show <ref>` | inspect one source node with file, heading line, and optional content |
| `mogent add <ref>` | add a source node to the manifest; writes the manifest by default, `--dry-run` previews only, `--rebuild` also rebuilds output through normal overwrite protection |
| `mogent diff` | manifest vs rendered output; later, drift vs on-disk `AGENTS.md` |
| `mogent edit <node>` | direct shortcut to the copy-on-write edit action |

Commands default to manifest mode. If no manifest flag is provided, mogent looks
for `agents.yaml` in the current directory. `--manifest <path>` only points the
command at a non-default manifest location; it does not opt into a separate mode.
Future commands that intentionally inspect or transform raw Markdown without a
manifest should use explicit command names or a clear `--no-manifest` style flag
only when that behavior is genuinely useful.

Source discovery should support several output shapes without changing the
underlying model:

- `coverage --source <alias>` filters to one source.
- `coverage --unused-only` prints a compact list of unused references.
- `coverage --content-only` hides empty organizational headings.
- `coverage --leaves-only` hides parent nodes and shows only terminal nodes.
- `coverage --depth <n>` / `--level <n>` limits the displayed tree depth in a
  later slice.
- `coverage --tree` renders an ASCII hierarchy rather than a flat list.
- `coverage --tag <tag>` filters once metadata tags exist.

Source inspection should stay separate from coverage. `source show <ref>` owns
file paths, heading line numbers, metadata/tags, TLDR fields, first-N-line
snippets, full content, and later source-vs-local diffs.

Manifest mutation commands write `agents.yaml` by default because that is the
primary authored state. `--dry-run` previews without writing. `--rebuild` also
writes `AGENTS.md`, using the same generated-output overwrite protection as
`build`; direct edits must not be overwritten silently. `--force` and optional
`.mogent` operation logs are future additions.

Dry-run output should offer several review views:

- `--preview=summary` shows the intended manifest operation in compact prose.
- `--preview=patch` shows the manifest insertion and rendered section addition.
- `--preview=tree` shows the document tree with the new node marked in place.
- `--preview=full` shows the complete rendered `AGENTS.md`.

The default should be compact enough for repeated CLI use; full output remains
available for debugging and careful review.

Localization should preserve upstream provenance without pretending the local
copy still inherits changes automatically. The exact names remain unsettled:
avoid over-committing to `local` as a universal concept if `source`, `scope`, or
location handles describe the model better. A localized node should record where
it came from, such as source alias, source reference, source URL or path, source
commit/hash when available, source file, heading path, localization time, and
original content hash. Updating or reconciling a localized node with upstream
changes should be a separate explicit workflow.

The TUI is a client of the model, not the product core. Durable behavior lives in
plain core operations that the CLI, TUI, and future GUI can all call:

- load a workspace from `agents.yaml`,
- list sources and source nodes,
- draft manifest changes,
- rebuild preview text,
- save `agents.yaml` and `AGENTS.md` through the same safety checks,
- inspect drift and import or reject local edits.

Rule of thumb: if the TUI can make a durable change, there should be a matching
noninteractive operation or command shape for agents and scripts. A polished TUI
still matters, but it should present decisions rather than hide the actual
workflow inside key bindings. Source: DI-vasel.

Rebuild milestones (from scratch; old code removed):

1. Parse manifest → resolve sources → render → validate. (No UI.)
2. The navigator (init step 2 / `tui`): tree, rendered-output, and source
   provenance previews; in-memory draft changes; confirmed atomic save/build.
   The exact interaction contract is `docs/IMPLEMENTATION-M2.md`.
3. Extract reusable workspace operations from the navigator and add CLI-facing
   command shapes for the same actions.
4. Copy-on-write editing, direct-edit import, source browsing, and richer
   save/history flows.

---

## 7. Deferred / future

- **Pinning / lockfile** — URL sources pinned to commit/tag; optional content hashes
  as integrity data. Not needed for local-path POC.
- **Swap alternatives via tags** — tags as a search/discovery layer; mark nodes as
  alternatives for a slot (the old XOR-group idea). After the core is stable.
- **Drift detection** — regenerate from manifest, diff against `AGENTS.md` on disk,
  offer explicit handling.
- **Promote local → shared** — push a copy-on-write override back up to its library.
- **Generated index** — fast search over large libraries for the navigator.
- **Pinned reference docs** — far-future support for making relevant project,
  API, design, or dependency docs easy for both humans and agents to find,
  pin, and cite. This may share source/manifest ideas with AGENTS composition,
  but should remain light until core prompt composition is stable.
- **Promise Grid** — libraries addressable by CID; grid `sources`. The named-source
  seam is already shaped for this.

---

## 8. Cleanup follow-ups

- [ ] Replace root `AGENTS.toml` + regenerate `AGENTS.md` via the new manifest path
  once rebuild milestone 1 lands (they remain as stale dogfood artifacts until then).
- [ ] Reconcile the YAML section of `brainstorn.md` (its include idea landed here as
  `sources`; its config sketch is superseded by §4).
- [ ] Rename `brainstorn.md` → `brainstorm.md`.
- [x] Removed old code (`tools/mogent`, `.mogent`), superseded Gen-1 docs, stray files.

---

## Document Map

Treat this file as the entry point.

| Doc | Role | Status |
|---|---|---|
| `docs/DESIGN.md` | **This file** — design of record | active |
| `TODO/TODO-jusuk-mogent-agent-modules.md` | Task tracking + DI log | active |
| `docs/thought-experiments/TE-bakom-...md` | TUI-first rationale | active (basis) |
| `docs/thought-experiments/TE-tavim-...md` | Reference-model exploration | historical — block/id framing superseded by §3 |
| `docs/brainstorn.md` | Raw brainstorm; taxonomy + fork-import notes | partial — config section superseded by §4 |
| `docs/other_repo_agents/` | Real-world `AGENTS.md` samples (untracked, pending review) | reference corpus |

Gen-1 tag-model docs and the old implementation were removed; recover from git
history if ever needed.
