# AGENTS.md — CogCtrlLab / Neural optimization project

Instructions for AI coding assistants. The root README is currently empty; this
file summarizes what's present from the code, configs, and vendored repos.

## What this repo is

An early-stage analysis project working with processed **dlPFC (dorsolateral
prefrontal cortex) neural data** (`Datasets/dlPFC_data_processed_ForQuincy...zip`)
in the UC Davis Cognitive Control Lab. Python, managed with **uv**
(`pyproject.toml`, package name `cogctrllab`, `requires-python >=3.14`). Core
scientific stack: numpy, pandas, scipy, scikit-learn, matplotlib, seaborn.

No analysis scripts/notebooks exist at the root yet — this is scaffolding plus
vendored tooling. If asked to build analysis, put it here and start a real
README/devlog.

## Vendored / reference subrepos (not the primary work — treat as external)

- `UC-Davis---Cognitive-Control-Lab/` — the lab's shared GitHub repo (its own uv
  project, plus `pass1_/pass2_for_supercat.bat`). Shared lab code; changes here
  belong upstream, not local-only edits.
- `RNEL_neuroshare/` — MATLAB + Python readers for Ripple Grapevine NIP / Nomad
  neuroshare files (Trellis Neuroshare API; Python lib by mfliu). Use its
  `wrappers/` (`read_continuousData`, `read_digitalEvents`, `read_spikeEvents`,
  `read_stimEvents`) to read raw neural recordings — don't reimplement readers.

## Guidance for assistants

- Use `uv` for dependencies; don't hand-edit the lock.
- Requires Python 3.14+ — be aware some scientific packages may lag that version;
  confirm the interpreter before debugging import/build errors.
- Don't commit the datasets zip or unpacked neural data.
- When you add the first real analysis, replace the empty root README with a
  description of the dataset, the question, and how to run the pipeline.
