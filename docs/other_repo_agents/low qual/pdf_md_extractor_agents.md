# AGENTS.md - pdf_md_extractor

## Project

Python 3.12 `uv` CLI for converting PPTX/PDF inputs into clean Markdown. The installed command is `docproc`, implemented in `docproc.__main__:main`.

## Structure

- `docproc/__main__.py` - CLI entry point.
- `docproc/config.py` - CLI/config handling.
- `docproc/processor.py` - conversion and extraction pipeline.
- `docproc/state.py` - MD5/state tracking for skipped and duplicate work.
- `README.md` - user-facing install and command examples.
- `SPEC.md` - design/reference material.
- `piazza scrape/` is a separate nested git repo; treat it independently.

## Development

- Use `uv` for dependency management.
- Install/edit locally with `uv tool install -e .`.
- If dependencies change, sync the isolated tool environment with `uv tool install --force -e .` or `uv tool upgrade pdf-md-extractor`.
- Run the CLI examples with `--dry-run` first when touching file moving, naming, or conversion behavior.

## Behavior To Preserve

- Keep conversion/extraction idempotent through state tracking.
- Preserve the three filename sanitization modes: `strict`, `med`, and `loose`.
- Preserve the space handling modes: `underscore`, `space`, and `dash`.
- Do not remove OCR fallback or alternate extraction engines unless the CLI/docs are updated.
- Be careful around file renames and output paths; this tool processes user document archives.

## Validation

- At minimum, run the affected `docproc` command with `--dry-run`.
- For dependency or CLI changes, run `uv tool install --force -e .` and smoke-test `docproc --help`.
