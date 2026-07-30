# Python

## Dependency Management

### uv

Use `uv` for dependency management when the repository uses `pyproject.toml`.
Do not hand-edit a lockfile. Confirm the supported Python version before
debugging dependency or import failures.

## File And CLI Work

Use a dry run before changing file names, locations, or a user's document
archive. Preserve idempotence where a tool may run more than once. Keep output
paths explicit and do not treat generated data as source material.

## Scientific Analysis

Pin or record the analysis environment when results depend on package versions.
Keep data-loading, transformation, and model assumptions visible. Do not
silently generalize constants across subjects, instruments, or experiments.
