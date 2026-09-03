# AGENTS.md — CogCtrlLab / SemanticMapProject

Instructions for AI coding assistants. The substantive analysis and its full
docs live in the `Meaning-Map_Analysis/` subdirectory — read
`Meaning-Map_Analysis/README.md` before touching analysis code.

## What this repo is

A container for meaning-map / salience analysis of NHP free-viewing eye-tracking
data. The repo root is a near-empty uv scaffold (`main.py` is a stub
"Hello from semanticmapproject!"); the real work is in subprojects.

## Layout

- `Meaning-Map_Analysis/` — **the real project: "PPC Inactivation GLMM."**
  Bayesian GLMM (Bambi/PyMC/ArviZ) of how PPC (posterior parietal cortex)
  inactivation affects fixation behavior, modeling P(fixate) from scene
  meaning, low-level salience, and center bias. Authors: O. Soyuhos (CClab
  functions/notebook), T. R. Hayes (VClab functions, original pipeline). Its
  README documents the full pipeline, data schema, module reference
  (`fun_VClab.py`, `fun_CClab.py`), per-monkey parameters, and outputs.
- `data_analysis_pipeline/` — empty uv scaffold (no deps, empty README); not yet
  built out.
- `images 2/` — zipped stimulus / salience image sets (1680×1050, 1920×1080).
  Large; keep out of commits.
- `main.py` — stub; ignore unless building a root entry point.

## Environments (important gotcha)

Two environments exist and they differ:
- **Conda (original, Windows):** `Meaning-Map_Analysis/code/environment.yml`,
  env `freeviewing_GLMM`, Python 3.9, `bambi 0.13` / `pymc 5.12` / `arviz 0.17`.
- **uv (this repo):** `pyproject.toml`, Python **3.14+**, but only a minimal
  Jupyter kernel — the scientific stack is **not** pinned here and must come from
  the conda env or be installed separately.

Confirm which environment the user is running before assuming package versions;
GLMM results can depend on the PyMC/Bambi version.

## Guidance for assistants

- Fixation-data CSV schema (per `Meaning-Map_Analysis/README.md`): `Subject,
  Scene, FixationNumber, x, y, Duration, Task, Visual_Field, Condition`.
  Predictor pipeline drops the first (cross) fixation and samples matched
  non-fixated control locations — preserve that logic.
- Per-monkey image size and fixation-window diameter differ (Quito 1920×1080 /
  59 px; Jimmy 1680×1050) — never hardcode one monkey's geometry globally.
- Outputs go to `output/<monkey>/`. Don't commit stimulus images, data CSVs, or
  `.nc` trace files.
