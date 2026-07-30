# Mogent Design (Design of Record)

Status: **active** · Last reconciled: 2026-07-29 · Supersedes: `examples/DESIGN-SUMMARY.md`, `thought-experiments/TE-mogent-module-architecture.md`, and the block/tag framing in `TE-tavim` and `DI-lorad`.

This is the single authoritative description of what mogent is and where it is going.
When any other doc disagrees with this one, this one wins.

---

## 1. Problem

You maintain an `AGENTS.md` in many repos. They share most content — identity,
workflow, constraints, style — but each needs local tweaks. Today every repo
hand-copies and re-edits a large file, so improvements never propagate and the files
drift apart.

## 2. Solution

Keep instruction content in a **library organized as a tree**. Markdown headings
*are* the tree. To produce a repo's `AGENTS.md`, **select a subset of the tree**;
selecting a node pulls its subtree. Render the selected nodes in order → `AGENTS.md`.

Author once, compose per repo.

That's the whole model. Everything below is the minimum needed to make it real.

Design principles:

- **Markdown-native.** The library is normal Markdown. No invented "block" concept —
  a heading and its content *is* the unit.
- **The tree is the model.** Selection, inheritance, and render order all come from
  the heading hierarchy. No parallel tag namespace, no wikilink graph.
- **Fail loud.** Missing files, missing paths/ids, empty output are errors, never
  silent drops.
- **Deterministic.** Same library + config ⇒ same `AGENTS.md`. Tool-only metadata
  never appears in the rendered output.

---

## 3. The four challenges (and their resolutions)

The core model surfaces exactly four questions. Nothing else is essential.

| # | Challenge | Resolution |
|---|---|---|
| 1 | How is a node **named** so config can select it, and the name survives edits? | **Heading path** by default (`instructions/testing`). Optional explicit `id` as a rename-proof escape hatch. Never content hashes. |
| 2 | How do **selection + inheritance** work? | Select a node → its subtree is included. Optionally deselect a descendant as an exception. Pure tree behavior. |
| 3 | How is the library **shared across repos** (so "author once" is true)? | `include` a shared library — local path or URL — with last-wins merge. Core, not future. |
| 4 | How does a shared node avoid being **repo-specific**? | Template variables (`{{ .repo_name }}`) filled at build. Core, not future. |

Challenges 3 and 4 are the point of the tool, not add-ons. Without them mogent is
just a file concatenator.

### 3.1 Node identity — path, with `id` as escape hatch

A node is addressed by its **heading path**: the slugified chain of headings from the
file's top down to it, joined with `/`.

```
# Instructions        -> instructions
## Testing            -> instructions/testing
### Flaky retries     -> instructions/testing/flaky-retries
```

Paths are readable and greppable. They break only when you **rename or move** a
heading — a rare, deliberate act, and a plain find-and-replace when it happens.

When you want a name that survives renames (e.g. a widely-referenced node in a shared
library), a node may opt into a stable `id` via a one-line HTML comment right after
its heading:

```markdown
## Testing  <!-- id: strict-testing -->
```

An `id`, when present, is a valid selection target *in addition to* the path. Most
nodes need no `id`. IDs are never forced on every node, and never content hashes — a
hash changes on every prose edit, which is the common case, so it breaks selections
exactly when you're improving content. (A content hash is fine later as a *lockfile*
integrity field; it is not an identity.)

### 3.2 Selection + inheritance

- Selecting a node includes it **and its entire subtree**.
- A descendant pulled in by a selected ancestor is *inherited*, not separately listed.
- You may **deselect** a specific descendant as an exception; its own subtree drops
  with it.

`mogent list` / `mogent tui` show this with three markers:

- `+` explicitly selected
- `|` inherited via a selected ancestor
- `-` inactive

This is the only inheritance concept in mogent. (The old *tag* inheritance —
`org/acme/team` implies `org/acme` — is retired with tags; see §7.)

### 3.3 Include + merge (sharing)

`AGENTS.toml` may `include` other libraries before applying local selection:

- an `include` entry is a **local path or a URL**;
- includes are resolved depth-first; **last definition wins** on conflict;
- local config always overrides included config;
- included libraries are untrusted input: paths are resolved safely (no traversal
  outside the library root), and a failed/unreadable include is an error, not a
  silent skip.

This is how a shared personal or team library propagates to every repo.

### 3.4 Templates

Modules are Go `text/template`s rendered at build time.

- Tool-provided vars: `repo_name`, `repo_url`, etc.
- User-defined vars from `[vars]` in `AGENTS.toml`.
- Example: a shared identity module says `You are working on {{ .repo_name }}` and
  fills in per repo, so one node serves every repo without editing prose.

Rendering happens after include-merge and before selection.

### 3.5 Editing = copy-on-write localization

Editing a module never mutates a shared/included source in place. When you edit a
node (`mogent edit`, or an edit through the TUI), mogent:

1. copies the node's source into the **local** `.mogent` library (under `library`),
   preserving its tree path, with your modification applied;
2. rewrites config so that path now resolves to the **local** copy instead of the
   included/shared original.

So the local repo always wins, edits are explicit and diffable, and the shared library
stays pristine. This is the concrete answer to "local vs global storage": shared
libraries are read-only inputs; any change becomes a local override that config points
at. Promoting a local override back up to a shared library is a separate, later action.

---

## 4. Configuration: `AGENTS.toml`

TOML is canonical. (YAML was explored in the brainstorm; the *valuable* idea there was
`include`, which TOML expresses fine — see §7.)

```toml
[config]
library = ".mogent/library"      # local module tree root

# shared libraries merged in before local selection (local path or URL; last wins)
include = [
  "~/.mogent/library",
]

# render order of top-level categories
order = ["identity", "instructions", "constraints", "format"]

# selected nodes (subtree-inclusive). Address by path, or by id for anchored nodes.
select = [
  "identity",
  "instructions/testing",
  "constraints/security",
]

# exceptions: drop a descendant that a selected ancestor pulled in
deselect = [
  "instructions/testing/flaky-retries",
]

# template variables
[vars]
project_name = "Agents"

[output]
path = "AGENTS.md"
```

Render order = `order` across top-level categories, then document order within each
subtree. No `[activate] scopes`, no tags, no wikilink/dotted-anchor syntax — those
belonged to the retired model.

---

## 5. Default category tree

`mogent init` scaffolds this tree by default:

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

These are the same *role* axis as the first four (what kind of content), just finer.
Weight (light vs heavy process), domain (grid), and language (go/rust) are **not**
categories — they are handled by `include` + `select`, and later `tags`.

---

## 6. CLI surface

| Command | Purpose |
|---|---|
| `mogent init` | Scaffold `AGENTS.toml` + the starter category **tree**; discover existing modules. |
| `mogent build` | include-merge → template → select subtree → render → validate → write `AGENTS.md`. |
| `mogent list` | Render the library tree with `+` / `\|` / `-` markers. Width-aware. |
| `mogent diff` | Compare selected node sets + rendered text. |
| `mogent edit` | Edit a node; copy-on-write localizes it into the local `.mogent` library and repoints config (§3.5). |
| `mogent tui` | Bubble Tea selector: browse the tree, toggle selection in memory. |

**Scheduled for removal (tag model retired):** `mogent tags`, the `--tags` flags on
`build`/`list`, and `internal/scope/`. Tracked in §8.

---

## 7. Deferred / future

Real directions, parked so they stop masquerading as current design:

- **Tags** — a *search/discovery convenience* over a large library (find nodes by
  topic), added on top of the tree once the core is stable. Not a selection model, not
  an inheritance mechanism.
- **Presets** — saved reusable selections (fast-iteration, design-heavy, session-log,
  etc.), expressed as named `select` sets.
- **Generated index** — internal index over the tree to power fast list/diff/search at
  scale.
- **Lockfile** — optional content hashes recorded for reproducible builds (integrity,
  not identity).
- **Local-vs-global storage policy**, **manual `AGENTS.md` drift detection**, and
  **import of hand-edited output** — each needs its own TE.
- **Promise Grid** — modules addressable by CID; remote/grid `include` sources. The
  `include` seam is already shaped to allow this.

---

## 8. Cleanup follow-ups

Debt to reconcile code/docs with this design (this pass is docs-only):

- [ ] Remove tag code: `cmd/tags.go`, `--tags` flags, `internal/scope/` (§6).
- [ ] Rework `AGENTS.toml` from `[category.*]/blocks` to `library` + `select`/`deselect` + `include` + `[vars]` (§4).
- [ ] Replace the heavy `agent_module:` YAML comment in modules with the optional
  one-line `<!-- id: … -->` form (§3.1); most nodes drop the comment entirely.
- [ ] Reconcile the YAML/templating section of `brainstorn.md` against §3.3–§3.4 / §7.
- [ ] Rename `brainstorn.md` → `brainstorm.md`.
- [x] Removed superseded Gen-1 docs (`examples/`, `TE-mogent-module-architecture.md`), stray `AGENTS.md.bak`, and empty `thought-experiments/archive/`.

---

## Document Map

Treat this file as the entry point.

| Doc | Role | Status |
|---|---|---|
| `docs/DESIGN.md` | **This file** — design of record | active |
| `TODO/TODO-jusuk-mogent-agent-modules.md` | Task tracking + DI log | active |
| `docs/thought-experiments/TE-bakom-...md` | TUI-first rationale | active (basis) |
| `docs/thought-experiments/TE-tavim-...md` | Reference-model exploration | partial — block/id-per-node framing superseded by §3.1 |
| `docs/brainstorn.md` | Raw brainstorm; taxonomy + fork-import notes | partial — YAML config folded into §3.3–§3.4 |
| `docs/other_repo_agents/` | Real-world `AGENTS.md` samples (dev-process + teaching + personal-project genres) | reference corpus |

Gen-1 tag-model docs (`examples/DESIGN-SUMMARY.md`, `examples/AGENTS-style-*.md`,
`thought-experiments/TE-mogent-module-architecture.md`) were removed as
stale/conflicting; recover from git history if ever needed.
