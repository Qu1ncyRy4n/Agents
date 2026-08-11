---
tldr: Follow the repository's Python environment and validate data work on safe fixtures or subsets.
---
# Python

## uv Dependency Management

Use the repository's declared environment manager. For `pyproject.toml` projects
that use `uv`, run commands through `uv` and do not hand-edit lockfiles. For
virtualenv-based projects, use the existing venv path when the repo documents
one.

## Nix Dependency Management

When Python dependencies are provided by a Nix flake or shell, enter the
documented `nix develop` environment before debugging imports or tool paths.
Treat `flake.nix` and Python project metadata as separate dependency surfaces:
do not update one to paper over drift in the other without explaining which
environment is authoritative for the repo.

## Scientific And Data Work

Keep data loading, transformations, constants, model assumptions, and output
paths visible. Do not silently generalize constants across subjects,
instruments, screen geometry, or experiments.

## Notebooks

Treat notebooks as source when they are the primary analysis artifact. Keep
paired scripts, exported data, and notebook outputs consistent with the
repository's convention. Do not overwrite exploratory notebooks for style.

## Validation

Use focused checks such as `python -m py_compile <file>`, `uv run ...`,
`pytest`, or a documented smoke command. For data-affecting changes, validate on
a fixture, temporary copy, or safe subset before touching authoritative data.

## Corpus Variants

The corpus includes research pipelines, local review UIs, scrape tools, and
experiment frameworks. Some use `uv`; others use an existing `.venv`. The module
should preserve that difference rather than forcing one Python workflow.
