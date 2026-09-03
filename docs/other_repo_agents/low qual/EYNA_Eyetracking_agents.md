# AGENTS.md — CogInContext / EYNA_Eyetracking

Instructions for AI coding assistants. This directory is **documentation, not
code** — the authoritative artifact is `EYNA_guide.html` ("EYNA Knowledge Base").

## What this is

EYNA is the lab's Python eye-tracking framework used by the experiment code in
sibling projects (notably `CogInContext/RoSE`). This folder holds the knowledge
base / install-and-usage guide for it, not the EYNA package itself. When RoSE (or
another task) says "see EYNA docs," this is where they point.

`EYNA_guide.html` covers:
- **Setup / install** — macOS, Windows, Ubuntu Linux recommended steps; technical
  notes.
- **Tutorial experiment** — Home Mode (no tracker) vs. Lab Mode (full eyetracking).
- **Architecture** — `experiment_utils.py → experiment.py` step mapping;
  key components: `CameraManager`, `Config`, `Eyetracker`, `Calibrate`,
  `VerifyCheck`, and their method-usage patterns.
- **Data files** — experiment data file, EYNA calibration data file, EYNA
  recording data file formats.

## Guidance for assistants

- Treat this as the reference for EYNA APIs and calibration/recording data
  formats when working on any task that integrates EYNA (e.g. RoSE's
  `eyna_utils.py`, `calibrate.py`).
- Do not edit the generated `EYNA_guide.html` by hand as if it were source; if it
  has a source (Markdown/generator) elsewhere, change that instead and confirm
  with the maintainer.
- Nothing here runs; there is no build step. Don't invent one.
