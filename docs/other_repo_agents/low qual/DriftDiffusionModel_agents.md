# AGENTS.md — DriftDiffusionModel

Instructions for AI coding assistants. This is an early-stage modeling project;
the README (`DriftDiffusionModel/README.md`) is a research-idea scratchpad, not a
spec — treat its later sections as hypotheses, not committed design.

## What this repo is

A computational model of executive function / cognitive control, framed around
the drift-diffusion model (DDM) and extending it. Intended stack: **Go
(test-driven)**, the **Emergent** framework, and **Cogent Core / lab**.
Architecturally inspired by the Necker-cube bistable-perception model.

## Planned architectural progression

1. Replicate the DDM via mathematical modeling.
2. Train an abstract neural network to reproduce the DDM.
3. Build a biologically-plausible NN grounded in brain areas hypothesized to
   support cognitive control.

## Research thesis (context, not yet implemented)

The author questions whether the DDM's drift-rate ("information uptake")
adequately captures EF. Proposed extension: parallel/enmeshed perceptual
processes (memory/regressive, present/instantaneous, future/predictive) feeding a
subsequent EF decision process that manipulates active information — i.e. EF as a
decision rule, not just evidence accumulation. See Weigard for the baseline DDM
parameters (`z`, `Ter`, `v`, decision thresholds). Candidate task: Stroop.

## Guidance for assistants

- This is exploratory. Before writing code, confirm which progression stage is
  actually being worked on — don't assume the biologically-plausible NN when the
  math replication may not exist yet.
- Prefer Go with tests, per the stated design. Match Emergent / Cogent Core
  idioms if extending those.
- The README mixes committed decisions with open questions and a collaborator
  list. Keep new design notes and TODOs there so the thinking stays in one place.
