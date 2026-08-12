---
tags: [lang/python, dependency/uv, research/reproducibility, scope/project]
tldr: Use uv with pyproject projects and avoid hand-editing lockfiles.
priority: 0.8
scope: project
---
# uv

Use `uv` for dependency management when the repository uses `pyproject.toml`.
Do not hand-edit a lockfile. Confirm the supported Python version before
debugging dependency or import failures.
