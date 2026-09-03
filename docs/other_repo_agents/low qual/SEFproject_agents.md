# AGENTS.md — CogCtrlLab / SEF project

Instructions for AI coding assistants. Read `Notes.md` for the running TODO and
Orhan's data-quality notes before changing analysis logic.

## What this repo is

Analysis of an **SEF (supplementary eye field) stimulation** free-viewing
eye-tracking experiment in a monkey (Vennie). The pipeline converts raw EyeLink
recordings into a fixations table, then fits a Bayesian GLMM of fixation
behavior as a function of scene meaning, low-level salience, and center bias.
Python, managed with **uv** (`pyproject.toml`, `requires-python >=3.12`).

## Pipeline

1. `preprocess_edf_to_fixations.py` (+ `.ipynb`) — raw EyeLink → fixations CSV.
   - Reads `Behavioral Data - Fullscreen Freeviewing/` (`*all.mat` workspace,
     EDF binary, 1 kHz `*.dat` gaze samples).
   - Maps trial → image/condition, extracts `FreeviewStart` timestamps, slices
     each 5 s free-view window, runs **I-DT** fixation detection.
   - Writes `data/Vennie_fixations_data.csv` with columns: `Subject, Scene,
     FixationNumber, x, y, Duration, Condition`.
   - Run: `uv run --with scipy,numpy,pandas preprocess_edf_to_fixations.py`
2. `SEF_stimulation_GLMM_stage1.ipynb` — Bayesian GLMM (PyMC / Bambi / ArviZ /
   nutpie). Writes to `output/` (`GLMM_data_stage1.csv`, `stage1_meaning.nc`,
   coefficient/marginal/heatmap figures).

`template code/` holds the reference pipeline this project adapts
(`fun_VClab.py` = T.R. Hayes core utils, `fun_CClab.py` = O. Soyuhos extensions)
plus related notebooks. `task_scripts/` holds the original MATLAB EyeLink task
(`fixationReward_final_withImage.m`, `fixationWrapper_freeview.m`) — reference
for what each EDF event/message means.

## Known data limitations (don't "fix" by fabricating data)

- **Block-1 stimulation trials are unrecoverable.** The MATLAB workspace is saved
  at session end, so block 2's ordering arrays overwrite block 1's. Block-1 BLANK
  trials survive via `all_no_trigger_trials`; block-1 STIM trials do not — the
  image at each position is lost. These ~25 trials/session are correctly excluded.
- Orhan flagged (2026-04-29) that the fixation window may be off; expected result
  is a normal distribution of fixations. Verify data processing against the task
  script before trusting downstream stats.

## Guidance for assistants

- Use `uv` for dependencies; don't hand-edit the lock.
- EDF parsing relies on `edf2asc -miss -1.0` conventions and the screen/geometry
  constants at the top of `preprocess_edf_to_fixations.py` (1920×1080, ~36 ppd,
  I-DT dispersion 72 px ≈ 2 dva). Keep these in sync if the rig geometry changes.
- Don't commit raw EDF / `.mat` / `.dat` data or the behavioral-data zip.
