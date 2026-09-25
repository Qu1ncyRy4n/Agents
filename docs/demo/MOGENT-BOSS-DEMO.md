# Mogent: Current Usage, Safety, and Direction

Audience: technical leadership familiar with the project

Target duration: 25-30 minutes

Format: off-screen presenter notes for a live terminal demonstration

Talk track:
- This is an update on what Mogent has become, not a first-principles sales pitch.
- The central change is that Mogent now manages reusable agent policy and related raw tooling as explicit, reviewable repository configuration.
- The demo starts from a deliberately messy copy of the todo app, recovers it safely, then shows controlled customization and reproducible remote inputs.

Demo:

```sh
export MOGENT_REPO=/home/qix/dev/cdint/Agents
cd "$MOGENT_REPO"
tools/install
tools/prepare-boss-demo

export QMR_LIBRARY=/home/qix/dev/cdint/qmr-agents-library
export DEMO_ROOT=/tmp/mogent-boss-demo/todo_app_project
export OUTLINE_MANIFEST="$DEMO_ROOT/demo-outline/agents.yaml"
export REMOTE_MANIFEST="$DEMO_ROOT/demo-remote/agents.yaml"
```

Safety/recovery:
- The helper refuses to overwrite an existing demo directory.
- All demo mutations occur in the disposable clone, not the working todo repository.
- If setup already ran, use the prepared directory rather than rerunning the helper.

Timing: 45 seconds.

---

# 1. Mogent Today

Talk track:
- Mogent builds repository-specific agent instructions and supporting tool trees from reusable libraries.
- It is now more than a Markdown concatenator: it provides source discovery, deterministic composition, drift handling, copy-on-write localization, raw directory outputs, and immutable Git pinning.
- It remains pre-v1. The reliable surface is the CLI; the TUI remains a proof of concept.
- No public license is granted yet, so this is currently an internal/project-owned tool rather than an externally adoptable package.

Key outcome:

```text
Central policy can be reused without hiding local selection,
local ownership, or ordinary Git review.
```

Timing: 60 seconds.

---

# 2. Architecture and Code Quality

Talk track:
- Package boundaries separate parsing, source indexing, rendering, filesystem writes, state, remote cache, and workspace orchestration.
- YAML is strict: unknown fields, duplicate keys, anchors, multiple documents, path traversal, overlapping outputs, and unsafe symlinks are rejected.
- Multi-output writes are transactional. Mogent snapshots files, directories, and state, then rolls everything back on failure.
- CI runs formatting, tests, vet, module-tidy checks, errcheck, dead-code analysis, a binary build, and diff checks.
- The CLI adapts the same `workspace` facade available to Go callers.

Demo:

```sh
cd "$MOGENT_REPO"
nix develop -c tools/check
```

Expected:
- All packages pass.
- `deadcode -test ./...` emits nothing.
- The command exits successfully.

Skip if short:
- State that the complete pinned check passed during rehearsal rather than running it live.

Timing: 60-90 seconds if run, 30 seconds if summarized.

---

# 3. Mental Model

Talk track:
- Libraries own reusable content and may publish both Markdown modules and raw directory trees.
- `agents.yaml` is local repository policy: aliases, selections, outputs, and variables are explicit.
- `mogent.lock.yaml` is present only for remote Git sources and pins immutable input.
- Generated output is committed for ordinary review; downloaded caches and machine-local state can remain ignored.

```text
library.mogent.yaml + Markdown + skills
                  |
                  v
       agents.yaml + optional lock
                  |
                  v
      AGENTS.md + .agents/skills/
```

Files:
- Commit: `agents.yaml`, remote lock, intentional local overrides, generated policy/tooling.
- Ignore: `.mogent/sources/` and usually `.mogent/state.json`.

Timing: 60 seconds.

---

# 4. The Todo App as a Real Consumer

Talk track:
- This is the actual todo app manifest copied into the demo repository.
- `qmr_repo` is a physical root alias. `qmr` is the logical source alias used in references.
- One source produces a rendered Markdown policy and a raw skill tree.
- The QMR root sidecar publishes only `agents/` as Markdown and `skills/` as a raw directory. Archive and reference material are excluded.

Demo:

```sh
cd "$DEMO_ROOT"
sed -n '1,120p' agents.yaml
mogent lib check --source qmr --manifest agents.yaml
```

Manifest to point out:

```yaml
roots:
  qmr_repo: /home/qix/dev/cdint/qmr-agents-library

sources:
  qmr:
    path:
      root: qmr_repo

outputs:
  - path: AGENTS.md
    include:
      - all: qmr
  - path: .agents/skills/
    from: qmr:skills
```

Expected:
- `Library sidecar check passed`.

Timing: 75 seconds.

---

# 5. Starting With a Messy Existing Repository

Talk track:
- The clone is intentionally staged like an existing repository entering Mogent adoption.
- It already contains `AGENTS.md` and skills, but generation state is absent.
- Both outputs also contain accidental direct edits.
- Mogent must diagnose this without assuming it may replace anything.

Demo:

```sh
cd "$DEMO_ROOT"
git status --short
git diff -- AGENTS.md .agents/skills/rust-development/SKILL.md
mogent status --manifest agents.yaml
```

Expected:
- Git shows deleted `.mogent/state.json` plus modified generated files.
- Mogent reports existing outputs as `untracked`, not clean.
- The hint recommends review before force.

Safety:
- Do not run `build --force` until the existing differences have been reviewed.

Timing: 75 seconds.

---

# 6. Inspect Before Writing

Talk track:
- Status explains ownership state; drift compares expected Markdown with disk.
- `drift --diff` is read-only even when no state exists.
- Raw directory output drift is reported by `status`; Markdown drift gets a unified patch.

Demo:

```sh
mogent drift --manifest agents.yaml --diff
mogent build --manifest agents.yaml
```

Expected:
- The diff exposes the accidental policy comment.
- Normal build refuses to overwrite unmanaged output.
- No files change.

Files changed:
- None.

Safety/recovery:
- This refusal is the desired adoption behavior, not an error to bypass automatically.

Timing: 60 seconds.

---

# 7. Restore a Reviewed Baseline

Talk track:
- We reviewed both differences and classified them as accidental.
- `--force` is now an explicit, informed ownership decision.
- The build writes all configured outputs and state as one transaction.

Demo:

```sh
mogent build --manifest agents.yaml --dry-run
mogent build --manifest agents.yaml --force
mogent status --manifest agents.yaml
git diff -- AGENTS.md .agents/skills/rust-development/SKILL.md
```

Expected:
- Dry run says `Dry run: no files written`.
- Real build says `Wrote generated outputs`.
- Both outputs become `up to date`.
- The two accidental edits disappear.

Important limitation:
- Current `build --dry-run` validates the primary Markdown render without writing; it does not simulate a complete directory synchronization plan.

Timing: 75 seconds.

---

# 8. Explore Sources Before Selecting Them

Talk track:
- Users can inspect source content without opening or knowing its file layout.
- References use `source-alias:directory/heading/path`.
- Source-relative directories contribute to references; filenames do not.
- Coverage distinguishes included, inherited, excluded, partial, and unused nodes.

Demo:

```sh
mogent source list qmr --manifest agents.yaml --tree --tldr
mogent source list --manifest agents.yaml --search security --tldr
mogent source show qmr:agents/constraints/security \
  --manifest agents.yaml --metadata --content=snippet --lines=8
mogent coverage --manifest "$OUTLINE_MANIFEST"
```

Optional display demo:

```sh
mogent source list qmr --manifest agents.yaml --tree \
  --chars unicode --fit term
```

Files changed:
- None.

Current boundary:
- Source browsing understands canonical output selectors.
- Coverage currently follows authored `doc:` selections, so demonstrate coverage
  against `demo-outline`; canonical selector coverage is planned follow-up work.

Timing: 90 seconds.

---

# 9. Modify Local Configuration

Talk track:
- The manifest is ordinary reviewed YAML, not hidden CLI state.
- We will temporarily exclude QMR's Security section from the Markdown output.
- A manifest or source change makes an intact generated output `stale`, distinct from direct edits.

Demo:

```sh
$EDITOR agents.yaml
```

Add beneath the Markdown output's `include` block:

```yaml
    exclude:
      - source: qmr:agents/constraints/security
```

Then run:

```sh
mogent status --manifest agents.yaml
mogent drift --manifest agents.yaml --diff
mogent build --manifest agents.yaml --dry-run
mogent build --manifest agents.yaml
```

Expected:
- Status reports `stale`, because disk still matches prior generated state.
- Drift shows the Security section being removed.
- Build applies the reviewed configuration.

Timing: 90 seconds.

---

# 10. Restore the Current Todo Configuration

Talk track:
- The previous change demonstrated a legitimate repository-local selection.
- For the rest of the demo, restore the todo app's checked-in composition.

Demo:

```sh
git restore agents.yaml
mogent build --manifest agents.yaml
mogent status --manifest agents.yaml
```

Expected:
- Security returns.
- Both outputs are clean.

Recovery:
- If the manifest edit was malformed, `git restore agents.yaml` returns to the known configuration before any build.

Timing: 30 seconds.

---

# 11. Generated Skills Have the Same Safety Boundary

Talk track:
- Directory outputs copy complete trees, including non-Markdown assets.
- Every generated file is hashed.
- Edited, added, or removed target files are detected.
- `drift` concerns the primary Markdown output; `status` covers every configured output.

Demo:

```sh
printf '\n<!-- demo local edit -->\n' >> \
  .agents/skills/rust-development/SKILL.md

mogent status --manifest agents.yaml
mogent build --manifest agents.yaml
```

Expected:
- `.agents/skills` reports `direct edits`.
- Normal build refuses to overwrite it.

Resolve after review:

```sh
mogent build --manifest agents.yaml --force
mogent status --manifest agents.yaml
```

Timing: 60 seconds.

---

# 12. Two Composition Modes

Talk track:
- Canonical `outputs:` entries independently select content for each target.
- The older `output` plus `doc` form remains supported when the repository must rename, group, or precisely place headings.
- Today, `mogent add`, `localize`, drift import, and the TUI edit that authored outline.
- This is a current interface boundary, not something to conceal.

Canonical selector output:

```yaml
outputs:
  - path: AGENTS.md
    include:
      - all: qmr
```

Authored outline:

```yaml
output: AGENTS.md
doc:
  - Instructions:
      - Security: qmr:agents/constraints/security
```

Demo:

```sh
sed -n '1,120p' "$OUTLINE_MANIFEST"
mogent status --manifest "$OUTLINE_MANIFEST"
```

Timing: 75 seconds.

---

# 13. Add a Source, Then Add Content

Talk track:
- Declaring a source and selecting content are separate decisions.
- `source add` changes only `sources:`.
- `add` previews placement and output impact before changing the authored outline.

Demo:

```sh
mogent source add fixture \
  "$MOGENT_REPO/testing-ground/fixture-library/cdint" \
  --manifest "$OUTLINE_MANIFEST" --dry-run

mogent source add fixture \
  "$MOGENT_REPO/testing-ground/fixture-library/cdint" \
  --manifest "$OUTLINE_MANIFEST"

mogent add fixture:go/tests \
  --manifest "$OUTLINE_MANIFEST" \
  --under Instructions \
  --heading 'Go Test Discipline' \
  --dry-run --preview=patch

mogent add fixture:go/tests \
  --manifest "$OUTLINE_MANIFEST" \
  --under Instructions \
  --heading 'Go Test Discipline' \
  --rebuild
```

Expected:
- Dry runs write nothing.
- Source declaration alone does not change generated output.
- `add --rebuild` updates manifest, Markdown, and state together.

Timing: 2 minutes.

---

# 14. Localize Shared Content

Talk track:
- Localization is copy-on-write ownership transfer for one selected module.
- Mogent copies Markdown into the consumer repository, records provenance, and repoints exactly one manifest reference to `local:`.
- It never writes to the shared QMR library.

Demo:

```sh
mogent localize 'Instructions/Security' \
  --manifest "$OUTLINE_MANIFEST" --dry-run

mogent localize 'Instructions/Security' \
  --manifest "$OUTLINE_MANIFEST" --rebuild

git -C "$DEMO_ROOT" status --short demo-outline
sed -n '1,160p' "$DEMO_ROOT/demo-outline/.mogent/provenance.yaml"
```

Expected new repository-owned files:

```text
demo-outline/.mogent/library/agents/constraints/security.md
demo-outline/.mogent/provenance.yaml
```

Expected manifest change:

```yaml
sources:
  local: .mogent/library

doc:
  # ...
  - Security: local:agents/constraints/security
```

Timing: 90 seconds.

---

# 15. Import an Intentional Direct Edit

Talk track:
- Sometimes the useful change begins in generated output during real use.
- Mogent can preserve one unambiguous edited section as a local module.
- It refuses edits outside the selected section or multi-section automatic import.

Demo:

```sh
$EDITOR "$DEMO_ROOT/demo-outline/AGENTS.md"
```

Under `## Focused Change Loop`, add:

```text
For this repository, run ./dev check before handing work back.
```

Then run:

```sh
mogent drift --manifest "$OUTLINE_MANIFEST" --diff
mogent drift --manifest "$OUTLINE_MANIFEST" \
  --import 'Instructions/Focused Change Loop'

mogent status --manifest "$OUTLINE_MANIFEST"
sed -n '1,160p' "$DEMO_ROOT/demo-outline/.mogent/provenance.yaml"
```

Expected:
- The edit becomes a new local Markdown source.
- The manifest reference changes from `qmr:` to `local:`.
- Output and state are rebuilt cleanly.

Safety:
- Automatic multi-section import is intentionally not implemented.

Timing: 2 minutes.

---

# 16. Library Sidecars Define the Published Surface

Talk track:
- A library sidecar separates source-repository organization from consumer-visible content.
- QMR publishes `agents/` as indexed Markdown and `skills/` as raw tooling.
- Archive, reference, proposals, and unrelated repository files remain outside the consumer surface.
- Sidecars also define identity, version, display order, titles, and validated relationships.

Demo:

```sh
sed -n '1,140p' "$QMR_LIBRARY/library.mogent.yaml"
mogent lib check "$QMR_LIBRARY"
mogent lib scan \
  "$MOGENT_REPO/testing-ground/fixture-library/cdint" --dry-run
```

Optional fresh-library example:

```sh
mogent lib init /path/to/new/library --dry-run
```

Safety:
- `lib init` refuses to overwrite an existing sidecar.
- `lib scan` inventories a full unsidecarred tree; `lib check` respects an
  existing sidecar's published content roots.

Timing: 90 seconds.

---

# 17. Remote Git Pinning

Talk track:
- Remote sources are explicit and immutable during ordinary use.
- Adding a URL does not fetch or select content.
- Pinning records the exact commit, source subdirectory, and content hash.
- Normal build, list, and show commands never silently update upstream content.

Demo:

```sh
mogent source add qmr_remote \
  https://github.com/Qu1ncyRy4n/qmr-agents-library.git \
  --subdir agents \
  --manifest "$REMOTE_MANIFEST" --dry-run

mogent source add qmr_remote \
  https://github.com/Qu1ncyRy4n/qmr-agents-library.git \
  --subdir agents \
  --manifest "$REMOTE_MANIFEST"

mogent source pin qmr_remote --manifest "$REMOTE_MANIFEST"

sed -n '1,180p' "$DEMO_ROOT/demo-remote/mogent.lock.yaml"
mogent source list qmr_remote \
  --manifest "$REMOTE_MANIFEST" --tree --depth 2
```

Expected:
- Pin prints the URL, full Git commit, and content SHA-256.
- The lock records `agents` as the narrow source subtree.
- The lock and ignored cache are created beside the remote demo manifest.

Current boundary:
- Remote pinning currently packages the selected source subtree, not a
  repository-root sidecar plus all of its declared content roots. Pinning QMR's
  narrow `agents` subtree avoids indexing unrelated archive material; preserving
  root sidecars in remote caches remains follow-up work.

Network recovery:
- If GitHub is unavailable, show the verified rehearsal lock output and continue. Do not troubleshoot credentials during the presentation.

Timing: 2 minutes.

---

# 18. Updating a Pinned Source Is a Review Step

Talk track:
- Updates are separate from builds.
- Preview shows the old and candidate commits plus Markdown file changes.
- Acceptance requires the exact full commit shown by preview.

Demo:

```sh
mogent source update qmr_remote --manifest "$REMOTE_MANIFEST"
```

Expected when already at current HEAD:
- Old and new commits may match.
- Markdown changes may be `none`.
- Mogent still prints the exact acceptance command shape.

Acceptance shape, but do not run unless the preview has a new commit:

```sh
mogent source update qmr_remote \
  --manifest "$REMOTE_MANIFEST" \
  --ref <full-previewed-commit> --accept
```

Timing: 60 seconds.

---

# 19. Team Repository Shape

Talk track:
- This is the intended review boundary for shared adoption.

Commit:
- `agents.yaml` for local composition.
- `mogent.lock.yaml` for remote source identity.
- Generated `AGENTS.md` and generated skills for ordinary code review.
- `.mogent/library/` and provenance only for intentional repository-owned overrides.

Ignore:

```gitignore
.mogent/sources/
.mogent/state.json
```

Important:
- Do not ignore all of `.mogent/`; that would hide intentional local policy.
- Do not place credentials or private personal content into a shared manifest merely because Mogent can resolve it.

Timing: 60 seconds.

---

# 20. Previous Shape 1: Tag-Filtered TOML

Talk track:
- The first implementation grouped modules under fixed categories and activated them through hierarchical tags/scopes.
- It worked, but the configuration described filtering rules rather than the final policy document.
- A reader had to replay ambient activation logic to understand the output.

Historical shape:

```toml
[config]
module_dir = ".mogent/modules"

[[category.lang.module]]
name = "go-style"
source = "go-style.md"
tags = ["lang/go"]

[activate]
scopes = ["person/quincy"]
```

Improvement:
- Current manifests attach explicit selection to explicit outputs.

Timing: 45 seconds.

---

# 21. Previous Shape 2: ID-Bearing Blocks

Talk track:
- A later design made sections selectable through IDs and metadata comments embedded in Markdown.
- That provided stable handles but imposed another abstraction and extra syntax on ordinary source documents.

Historical shape:

```markdown
# Instructions

<!--
agent_module:
  id: instructions
  tldr: Defines repository process rules.
-->
```

Improvement:
- Current libraries use Markdown headings and source-relative directories as the tree.
- File frontmatter supplies optional discovery metadata without rendering.
- Mandatory IDs and Obsidian-style block references are no longer required.

Timing: 45 seconds.

---

# 22. Previous Shape 3: Authored YAML Outline

Talk track:
- The first YAML model was a major improvement: explicit sources and a readable final document outline.
- The manifest owned output heading names and placement.
- This form remains useful and supported; it powers the localization portion of this demo.

```yaml
sources:
  shared: library

output: AGENTS.md

doc:
  - Identity:
      - Role: shared:identity/role
  - Instructions:
      - Testing: shared:lang/go/testing
```

Limitation:
- It centered one authored document and made multiple independent output targets awkward.

Timing: 45 seconds.

---

# 23. Current Shape: Output-Scoped Selectors and Sidecars

Talk track:
- Outputs are peers, each with visible selection.
- Root aliases avoid repeating a physical repository path.
- Sidecars publish a bounded library surface and allow raw directory trees beside indexed Markdown.
- The older authored outline remains available where custom heading placement is genuinely needed.

```yaml
roots:
  qmr_repo: /path/to/qmr-agents-library

sources:
  qmr:
    path:
      root: qmr_repo

outputs:
  - path: AGENTS.md
    include:
      - all: qmr
  - path: .agents/skills/
    from: qmr:skills
```

Key evolution:

```text
ambient filtering
  -> embedded block identity
  -> explicit authored document
  -> explicit independent generated targets
```

Timing: 60 seconds.

---

# 24. Near-Term Planned Behavior

Talk track:
- The next productivity work is about making authoring operations as strong as build and inspection.

Planned interfaces:

```sh
mogent move <heading> --before <heading> --dry-run
mogent module new --source <alias> --path <path> --heading <text>
mogent include normalize <output> --style explicit|all-except --dry-run
```

Goals:
- Move and reorder by stable heading path rather than numeric index.
- Create modules only in explicitly writable local sources.
- Preview manifest, coverage, and rendered differences together.
- Normalize equivalent selector strategies without changing output order.
- Extend authoring operations to the canonical selector model rather than leaving them `doc:`-only.

Status:
- These commands are planned, not implemented.

Timing: 60 seconds.

---

# 25. Longer-Term Direction

Talk track:
- Local customization eventually needs better origin reconciliation than a one-time copied file.

Planned model:

```text
BASE   = content when localized
ORIGIN = current shared source
LOCAL  = repository-owned customization
```

Candidate direction:

```sh
mogent preserve <heading> --dry-run
mogent inherit <heading> --ff-only
mogent propose <heading> --to <writable-source> --dry-run
```

Other planned areas:
- Tool-specific rendered outputs and exact mirrors/symlinks.
- Enforced requirement/conflict relationships after semantics are approved.
- Transport-neutral proposals for upstreaming useful local policy.
- Better source-version drift visibility.

Deliberate limits today:
- No automatic multi-section import.
- No automatic three-way merge.
- No authenticated or non-Git remote transport.

Timing: 75 seconds.

---

# 26. Closing Takeaway

Talk track:
- The project has moved from renderer experimentation toward a practical agent-tooling platform.
- The useful interface is the complete loop: discover, preview, compose, generate, detect drift, preserve intentional differences, and review ordinary files.
- Safety is not a separate feature. Refusal, immutable pins, explicit force, bounded import, and transactional writes are part of normal usage.
- The next challenge is making canonical authoring and origin reconciliation as polished as current inspection and generation.

```text
Reusable central policy
        +
explicit local ownership
        +
ordinary Git review
        =
safe customization across repositories
```

Optional final demo:

```sh
mogent status --manifest "$DEMO_ROOT/agents.yaml"
mogent status --manifest "$OUTLINE_MANIFEST"
mogent status --manifest "$REMOTE_MANIFEST"
```

Timing: 60 seconds plus questions.
