# Mogent Multi-Repo Deploy And Guided Dogfood

Created: 2026-10-09 19:12:38 PDT
Status: planned; project inventory and selections require the guided sessions below
Owner: Quincy Ryan

## Goal

Make a manifest of existing projects, choose existing guidance modules and skills,
author a suitable Mogent configuration for each project, then deploy and dogfood
them. Walk through the configuration with the owner so they can experience it,
predict what it will do, and judge whether it feels intuitive.

Success means useful agent behavior and understandable configuration across real
projects, as well as successful builds. This document records the campaign and
its results; each consuming project's `mogent.hcl` remains its executable recipe.

## Two Manifests With Different Jobs

1. **Project deployment manifest:** the roster of repositories, current files,
   selected modules/skills, intended outputs, experiment status, and UX findings.
   Begin with the table below; turn it into a separate machine-readable registry
   only if coordinating the campaign needs that. It is not a new Mogent grammar.
2. **Per-project configuration:** ordinary `mogent.hcl` consumed by existing
   `source`, `plan`, `apply`, and `update` commands.

The project roster coordinates individual runs. There is no implemented
cross-repository batch-apply command or multi-repository transaction.

## Project Deployment Manifest

First confirm the code roots and intended project set with the owner. Inventory
existing instructions, configs, skills, guides, and relevant repository state;
inspect stack/workflows rather than choosing modules from repository names.

Seed candidates known from earlier work, not a completed filesystem survey:

| Project | Role | Inventory/stack | Existing config/instructions | Guidance / skills | Outputs | Campaign status |
|---|---|---|---|---|---|---|
| `cdint-nix` | Consumer candidate | Verify during inventory | Existing HCL demo config; inspect actual instructions/skills | Decide with owner | Decide | Candidate |
| `todo_app_project` | Consumer candidate | Verify during inventory | Inspect current files and existing work | Decide with owner | Decide | Candidate |
| `Agents` / Mogent | Optional self-hosting consumer | Go tool; verify desired self-use scope | Existing test fixtures are not a deployment | Decide whether to include | Decide | Candidate |
| `qmr-agents-library` | Source library | Existing guidance and skill modules | Live library sidecar | Catalog existing content for consumers | Source exposure only | Source preparation |
| _Additional project_ | _Consumer/source_ | _Verified facts_ | _Observed files_ | _Chosen references_ | _Exact paths_ | _Candidate_ |

For every included consumer, record:

- A confirmed local checkout location or stable nickname, and active branch/worktree.
- The agent tools actually used, stack, stage, and a concrete task for dogfooding.
- Existing `mogent.hcl` / `agents.yaml`, root/nested `AGENTS.md`, `CLAUDE.md`, and
  skill directories; identify generated ownership versus hand-authored material.
- Repository-specific guidance that must be represented in the resulting recipe.
- Exact library/node references and reasons for selecting or omitting them.
- Output paths, source version/local state, latest plan/apply results, and UX notes.

Campaign states: **candidate → inventoried → selections reviewed → configured →
planned → applied → exercised → reviewed**. Record a blocker rather than advancing
a repository when a required step fails.

## Guided Sessions

### 1. Inventory Together

- [ ] Confirm the project roots and full candidate list.
- [ ] Populate the deployment manifest with observed files and verified stacks.
- [ ] Choose one representative pilot consumer before configuring the rest.
- [ ] Choose an actual agent-assisted task to perform after deployment.

The owner sees the inventory before module choices are made. Ask which existing
instructions matter and what the agent should help with in that project.

### 2. Choose Existing Content

- [ ] Browse the live library's available guidance and skills.
- [ ] Make a per-project selection matrix: reference, purpose, selected/omitted,
  explicit choice versus accepted default, and expected output destination.
- [ ] Start with a small useful baseline, then add project-relevant skills.
- [ ] Preserve local distinctions rather than copying one large config to every repo.

Use discovery against the project's config once it is drafted:

```sh
mogent source list qmr --config /path/to/project/mogent.hcl --tldr --tags
mogent source show qmr:workflow/git-commit-mode --config /path/to/project/mogent.hcl
```

The existing QMR authored root and legacy skill tree can be consumed immediately.
Selecting individual skill bundles through the typed API needs a named physical
root. If the live sidecar has no such exposure, review this additive declaration
with the owner before using the typed example below:

```hcl
# In qmr-agents-library/library.mogent.hcl
content "skill-files" {
  directory = "skills"
}
```

Use a unique root name: `skill-files` does not collide with the legacy `skills`
tree. This exposes existing files; it does not require rewriting skill content.
Once available, browse a bundle and its optional document view:

```sh
mogent source show qmr:skill-files/nix-development --config /path/to/project/mogent.hcl
mogent source show qmr:skill-files/nix-development/SKILL.md --headings --config /path/to/project/mogent.hcl
```

### 3. Author The Pilot Configuration With The Owner

Show the whole small config, explain a few lines at a time, then let the owner
change a selection or destination. Before running it, ask them to predict the
result. Record where explanation was necessary rather than smoothing over friction.

Illustrative consumer recipe, conditional on the named library exposure above:

```hcl
mogent { format = 2 }

sources {
  source "qmr" {
    local = "../qmr-agents-library"
  }
}

outputs {
  output "instructions" {
    path = "AGENTS.md"
    kind = "markdown"

    source "qmr" {
      from = "qmr:agents"
      select = {
        intro       = true
        constraints = true
        workflow    = { accept_defaults = true }
      }
    }
  }

  output "skills" {
    path = ".agents/skills"
    kind = "dir-tree"

    source "qmr" {
      from      = "qmr:skill-files"
      node      = "nix-development"
      operation = "copy"
      into      = "nix-development"
    }
  }
}
```

This is an example to edit together, not a preselected deployment for every
project. Verify `local` relative to each consumer config; the sibling path suits
only matching checkout layouts. Verify the selected skill fits the actual task.

Walk through these distinctions:

- Source location versus library reference in `from`.
- Authored `select` choices versus physical `node` selection.
- Rendered `AGENTS.md` versus a complete copied skill bundle.
- Output `path` versus contribution `into`.
- Accepted library defaults versus explicit per-project choices.
- Source filtering versus replacement of an earlier contribution.

Generated `CLAUDE.md` or other destinations are chosen per actual tool usage.
Identical instruction copies can use `paths`, but nested instructions and
tool-specific discovery require deliberate per-project decisions.

### 4. Plan, Predict, And Review

- [ ] Record the owner's predicted files/headings before running the plan.
- [ ] Run the plan for this consumer and compare it with that prediction.
- [ ] Inspect additions, removals, accepted defaults, provenance, and collisions.
- [ ] Account for existing hand-authored content before assuming generated ownership.
- [ ] Edit the config and re-plan until the result fits the project.

```sh
mogent plan --config /path/to/project/mogent.hcl
```

`plan` does not write outputs, but a missing Git source checkout may be cached.
Use local sources during the initial authoring loop or explicitly chosen Git
revisions for repeatable consumption. Do not use `--force` to bypass unexplained
loss of existing instructions or skill content.

### 5. Apply The Reviewed Pilot And Perform Real Work

- [ ] Apply the reviewed config; reconcile unmanaged outputs if adoption is needed.
- [ ] Re-plan and expect no remaining output changes.
- [ ] Confirm the receiving agent discovers the intended instructions and skills.
- [ ] Perform the chosen real task with the installed guidance.
- [ ] Record what helped, what was ignored, what was missing, and what was excessive.

```sh
mogent apply --config /path/to/project/mogent.hcl
mogent plan --config /path/to/project/mogent.hcl
```

A successful copy is not proof that the receiving tool loaded it. Keep deployment
results and actual agent-use observations separate. Record generated ownership
and the project's decision about committing outputs; ignore `.mogent/` state/cache.

### 6. Configure The Remaining Manifest Entries

- [ ] Apply lessons from the pilot to each project's own recipe.
- [ ] Walk through materially different choices with the owner.
- [ ] Plan and apply one reviewed repository at a time; record per-repo status.
- [ ] Exercise a concrete task in each deployed project.
- [ ] Test one library update/replacement after the initial deployments settle.

For a pinned Git source:

```sh
mogent update qmr --config /path/to/project/mogent.hcl
mogent update qmr --config /path/to/project/mogent.hcl --accept
mogent plan --config /path/to/project/mogent.hcl
mogent apply --config /path/to/project/mogent.hcl
```

Update acceptance moves pins, not generated outputs. A source change can affect
several consumers, but each consumer still has its own recipe, state, and review.

## Usability Record

For each walkthrough, fill a short record:

| Project / date | Operation/config change | Owner expected | Actual result | Felt natural or confusing | Follow-up |
|---|---|---|---|---|---|
| _Pilot_ | _Choose skill / accept defaults / move destination_ | _Prediction_ | _Observed_ | _Exact keyword or step_ | _Local config change / product issue_ |

Ask the owner:

1. Could you predict what this config would generate without my explanation?
2. Which names or repeated references made you hesitate?
3. Did `path`, `from`, `node`, and `into` feel distinct and necessary?
4. Did the plan make defaults, copies, renders, and replacements clear?
5. Did the installed guidance improve real work, or only enlarge the context?

Keep repeated friction as a product/design issue. Keep legitimate project-specific
preferences in that project's config. Link active consumer experiments from
`docs/ACTIVE-DOGFOOD.md`; library-source smoke tests and test fixtures are not
substitutes for real consumer results.

## Completion Criteria

- [ ] All agreed projects are inventoried in the deployment manifest.
- [ ] Every included consumer has reviewed guidance/skill choices and its own config.
- [ ] Each deployment has an inspected plan, recorded apply, and unchanged re-plan.
- [ ] The owner has directly experienced configuration edits and at least one update.
- [ ] Each deployed project has a real agent-use observation and UX record.
- [ ] Remaining blockers and repeated friction are linked to concrete follow-up tasks.

## References

- [Typed content grammar and runnable demo](../TYPED-CONTENT.md)
- [Current v2 specification](../proposals/MOGENT-HCL-V2-SPEC.md)
- [Active dogfood registry](../ACTIVE-DOGFOOD.md)
- [Composition roadmap](../../TODO/TODO-tugur-v2-composition-followup.md)
- [Accepted typed-content design](../../DR/DR-juvih-typed-content-composition.md)
