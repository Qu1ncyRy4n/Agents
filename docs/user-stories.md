# Mogent User Stories

The working set of user stories for Mogent, in rough priority order. Core
stories (US-1..US-4) drive near-term direction; review-later stories
(US-5..US-8) are documented but not scheduled. The CIWG/PromiseGrid
consolidation story lives separately in
[`user-stories-ciwg.md`](user-stories-ciwg.md).

Each story carries short scope notes capturing the criticisms and refinements
raised during review, so the framing is not lost.

---

## Core stories

### US-1 — Shared core template for personal and professional use

**As** someone setting up agent guidance for both personal and professional
projects, **I want** a small core template that both can share, with the option
of context templates that share most but not all modules, **so that** I get a
useful `AGENTS.md` quickly without learning the full model.

Acceptance criteria:
- `mogent init` offers a minimal, legible manifest and explains every field.
- A personal repo and a professional repo can both start from the shared core
  and then include or drop individual modules.
- One command renders a usable document from local or pinned sources.

Scope notes:
- Personal and professional are divergent audiences; the answer is a shared
  core plus per-context deltas, not one generic template.
- "No unexplained concepts" applies to the starter surface, not to the tool's
  full capability.

### US-2 — Continuous borrowing between personal and professional

**As** a developer working across personal and professional repos, **I want** to
borrow org-style practice, modify it, and layer on accessibility, style, and
workflow additions with minimal friction, **so that** improvements flow between
contexts without manual copy/paste.

Acceptance criteria:
- Localize a module from a source alias, edit it, and keep provenance to the
  origin revision.
- Add local modules alongside borrowed ones and render them together.
- Pull upstream updates later; conflicts are surfaced rather than clobbered.
- Friction budget: one command, no manual file moves, diff shown before write.

Scope notes:
- Org practice is emulated for now by a second local library; no real org
  library is required to develop this.
- Publishing or redistributing org-derived practice needs an explicit
  authority/licensing boundary before it is treated as shareable.

### US-3 — Multi-agent prompting (manager)

**As** a manager coordinating guidance for multiple agents or roles, **I want**
shared org policy plus per-role overlays, **so that** many agents follow the
same constraints while retaining role-specific instructions.

Acceptance criteria:
- A core instruction layer plus role overlays that may tighten but never relax
  the core (the wire-lab / grid-examples model).
- Defined precedence across personal, org, and repo sources.

Scope notes:
- Runtime multi-agent coordination is out of scope. This story covers
  multi-agent **prompting** only.
- This is expected to be the fusion point for a lot of responsibility; the exact
  boundary is TBD and should be revisited.

### US-4 — Adoption into a large existing codebase

**As** a developer adding Mogent to an existing repo full of inconsistent
`AGENTS.md`, `CLAUDE.md`, and `.agents/skills/` files, **I want** to capture and
reconcile that guidance into a manifest and library with provenance, **so that**
I converge on one source of truth incrementally without a big-bang rewrite.

Acceptance criteria:
- Read-only capture and inventory of root, nested, role-specific, and skill
  files.
- Classify every rule by scope and flag duplicates and conflicts.
- Import into modules with origin revision; `status` / `--dry-run` shows the
  diff before any write.
- Adopt one repo and one section at a time; nothing is silently overwritten.

Scope notes:
- Capture (read-only) is separate from reconcile (writes, which need owner
  decisions). The hard part is authority when sources disagree, not parsing.
- `AGENTS.md` and skills are different artifact types and may need separate
  ingestion models.
- "Converged" needs a defined target state to be measurable.

---

## Review later

Documented for completeness; not scheduled.

### US-5 — Library curator

**As** a library maintainer (e.g. a CIWG/org curator), **I want** to receive
proposals, review non-overlapping merges, publish revisions, and deprecate
modules with provenance intact, **so that** the shared library stays coherent as
many repos contribute.

### US-6 — CI drift gate

**As** a CI automation maintainer, **I want** the build to fail when a checked-in
`AGENTS.md` no longer reproduces from its manifest, **so that** generated output
is never hand-edited and drift is caught at the PR.

### US-7 — Runtime retrieval and token cost

**As** an agent operating in a repository, **I want** to discover the right
module or skill by tag on demand without loading everything, **so that** context
stays small and relevant.

### US-8 — Authority auditor

**As** an organization reviewer, **I want** to prove that no repository relaxed a
required org rule and to see which policy revision each agent acted under, **so
that** shared policy is enforceable and auditable.
