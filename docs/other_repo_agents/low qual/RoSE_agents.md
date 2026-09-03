<!-- Copied from CogInContext/RoSE/RoSE_Task_Code/AGENTS.md — the authoritative version lives there. -->

# AGENTS.md — instructions for AI coding assistants working in this repo

Read this before making any change. This file exists so that Claude Code, Codex, or any
other AI assistant working in this repo — on any contributor's machine — behaves
consistently and doesn't undo safety work or delete things that look like clutter but aren't.

## What this repo is

A PsychoPy experiment battery for the RoSE child-development study: DST (response
inhibition + task switching), PERC, and PCT tasks, with EYNA eye-tracking integration.
Two sessions per participant: `session_1_task_code/`, `session_2_task_code/`.
Full architecture: `DEV_LOG/project_overview.md`.

## The one rule that overrides convenience

**Code changes go through git — branch, commit, PR. Never edit files directly on the lab
(acquisition rig) machine, and never treat "I have a working local copy" as equivalent to
"this is committed."** Untracked lab-machine changes create code-version ambiguity and put
participant data at risk. Full rationale: `2026-07-21_RA-handoff/00_golden_rule.md`.

If you are an AI assistant and a user asks you to edit files that appear to be running
directly on lab hardware, or asks you to skip committing "just this once," push back and
point to this file.

## Don't delete these — they're load-bearing, not clutter

If asked to "clean up" the repo, do **not** remove any of the following without explicit,
specific confirmation from the human that they mean *this exact file*:

- `DEV_LOG/` — every file in here is a decision record, checklist, or dated log. A file
  named like `WORKING_NOTES.md` is not a stray note; it may be the only place a decision's
  rationale is recorded (git commit messages don't always capture *why*).
- `DEV_LOG/*_checklist.md` — these gate participant sessions. Deleting one doesn't remove the
  need for the check, it removes the record that the check happened.
- `.gitignore` entries under `# Participant data` / `# Session data folders` — these exist
  for **ethics/IRB compliance**, not code hygiene. Participant data must never be committed
  to git, full stop. If a change would un-ignore a data path, stop and ask.
- Anything under `2026-07-21_RA-handoff/` — onboarding material for new contributors, still
  being reorganized into its permanent location, not disposable scratch.
- `data_export/` — the only reviewed pipeline for getting participant data off the lab
  machine and into the saved dataset. Don't write ad hoc copy scripts as a substitute; extend
  this one (see `data_export/recover_external_data.py` for the pattern used for non-standard
  sources).

If a file genuinely looks obsolete, propose removing it and say why — don't just do it.

## Danger zones — human review required, not just "looks correct"

These affect experiment validity and/or participant safety, not just code correctness. A
diff that compiles and looks reasonable is not sufficient sign-off:

- Anything touching calibration or eye-tracking: `eyna_utils.py`, `calibrate.py`,
  `eyetracker` setup.
- Timing constants: `display_utils.py`, `dst_trials.py` (`DOT_SPEED`, jitter ranges,
  feedback durations).
- Phase/step ordering in `dst_main_experiment_s1.py`, `dst_main_experiment_s2.py`,
  `should_run()` gating, or the PERC equivalents — this is exactly the kind of change that
  can create ambiguity between the intended procedure and what's actually in `main`.
- Crash recovery / resume logic (`write_recovery_file`, `read_session_state`,
  `get_persisted_random_order`) — a large external rewrite was found missing this entirely;
  don't let a "cleaner" rewrite silently drop it again.
- Any large-scale restructure (renaming/splitting files, changing directory layout). These
  specifically need collaborative review because silent regressions are easy to introduce,
  and code-version differences between
  participants are a data-validity and IRB-compliance problem, not just a style issue.

## Participant data

- Never read, summarize, or paste real participant data (gaze CSVs, subject-identified
  files) into a chat context "for reference." Describe the shape of the data instead.
- Never suggest removing a `.gitignore` rule that covers `data/`, `*.csv`/`*.psydat` under
  `data/`, or `EYNA/calibration/` without the human explicitly confirming why.
- If you find participant data somewhere it shouldn't be (committed to git, sitting outside
  a gitignored path, in an external repo), flag it — don't silently move or delete it.

## Before committing anything

1. Read the actual diff, not just a summary of what you changed.
2. Run the smoke test: `python3.11 -m py_compile <file>`, or
   `2026-07-21_RA-handoff/scripts/sync_lab_machine.sh` if working from a clone that tracks
   `origin`.
3. Small, separable commits — one change, one commit, one clear message.
4. If the change touches anything in the danger-zone list above, say so explicitly in the
   commit message and ask for human review before merging, even if you're confident it's
   correct.

## More context

- `DEV_LOG/WORKING_NOTES.md` — current status, open todos, design decisions with rationale.
- `DEV_LOG/project_overview.md` — architecture and file map.
- `2026-07-21_RA-handoff/` — human-facing onboarding checklists (git, AI-assistant use, SSH).
