# AGENTS.md — CogCtrlLab / Video project

Instructions for AI coding assistants. Read before making changes. Derived from
the project README (`CogCtrlLab/Video project/README.md`), which is the
authoritative, continuously-updated source — consult it for full detail.

## What this repo is

A monkey movie free-viewing eye-tracking task: the monkey holds fixation, watches
a movie, gets a juice reward. Built on **MATLAB + Psychtoolbox (PTB) + EyeLink**.
Not Python/PsychoPy — you run it inside MATLAB, from the `Code/` folder, not a
terminal. Started as S. Novik's rotation project (based on O. Soyuhos' 2025
fixation-training code).

## Key entry points

- `Code/RUN_freeviewingTraining_movie.m` — run this (trial state machine).
- `Code/CONFI_freeviewingTraining_movie.m` — all params. Set `computer_name` at
  the top; a `switch` block derives `dummymode` and `filepath` per machine.
- `Code/pseudorandomization.m` — movie-order picker; reads
  `video_ebm_dataset/MANIFEST.csv`.
- `Code/cclab-matlab-tools/` — lab MATLAB utilities (git submodule); provides
  `cclabReward`, `cclabInitDIO`, `cclabPulse` for real reward/TTL.

## Ground rules

- **Movies are NOT in this repo** and are gitignored. They live on the lab NAS.
  `MANIFEST.csv` is read from the repo, not from the video path — don't "fix"
  code to copy it alongside videos.
- **Never enable TTL sync pulses (`cclabPulse`) blind** — they're intentionally
  commented out until line assignment is confirmed against the Neuropixel setup.
- `SkipSyncTests=1` disables PTB timing verification — acceptable for mouse
  testing (`dummymode=1`), **never silently on the real rig**; it's a pending
  decision for real data collection.
- Fixation window is a **square** (half-width `windowSize` deg) despite "radius"
  naming; the dot is drawn as an oval. Don't "correct" this without checking.
- `TrialSuccess=1` means fixation held through **phase 1 only** (dot-on-movie);
  phase 2 is genuine free-viewing with no gaze requirement. Keep this in mind
  for any viewing-time analysis code.
- Reward `.png` images must be in the working directory (`rewardImagePath = pwd`)
  — always assume the task runs from `Code/`.
- `practiceBlockSize` must stay `0` — the practice path uses old `dir()` structs
  incompatible with the current trial-struct format and will crash if enabled.

## Before changing behavior

- The README's **Config params**, **Trial state machine**, **Output data**, and
  **Open issues / Devlog** sections are kept current — read the relevant one
  before editing timing, the Results table schema, or EyeLink messaging.
- When you change the Results table, EyeLink message markers, config params, or
  the trial flow, **update the README** (tables + Devlog entry, dated) in the
  same change. That devlog is how the maintainer tracks the project.
- Distinguish the two "sync" concepts the devlog warns about: EyeLink/TTL
  cross-device alignment vs. PTB's display-timing (`Screen('Flip')` VBL)
  self-test. They are unrelated layers; don't conflate them.

## Data / safety

- `.mat` Results tables are MATLAB opaque objects — Python/scipy cannot read
  them; open in MATLAB. Config structs are Python-readable.
- Don't commit videos, EDF files, or `.mat` session data.
