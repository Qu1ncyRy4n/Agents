---
tldr: Protect sensitive research data and require evidence for validity-affecting changes.
---
# Research Data

## Protection

Do not read, summarize, paste, commit, or expose identified participant,
subject, gaze, behavioral, calibration, or raw experiment data. Describe the
shape of the data instead of copying records into chat or docs.

## Ignore Rules

Preserve ignore rules that protect raw recordings, exports, stimulus media,
participant/session folders, calibration data, and derived sensitive outputs.
Do not unignore data paths without explicit human confirmation.

## Experiment Validity

Require human review for calibration, timing constants, trial ordering,
condition mapping, crash recovery, resume logic, participant flow, and any
change that can affect experiment validity or safety.

## Known Limitations

Do not fabricate or infer missing data to make an analysis complete. Preserve
known exclusions and limitations in the analysis notes and output narrative.

## Validation

Validate data-processing changes against the task procedure, fixture data, or a
safe copy before trusting downstream statistics.
