# Repository Guidelines

## Project Structure & Module Organization
- `content/` holds Markdown pages and posts with YAML front matter (`---`) like `title`, `date`, `linktitle`, and tags.
- `layouts/` and `themes/hugo-book/` define templates, shortcodes, and theme overrides.
- `static/` stores assets copied as-is; `data/` holds structured data; `archetypes/` stores Hugo templates.
- `config.toml` contains site-wide configuration and theme settings.
- `public/` and `resources/` are generated outputs (GitHub Pages + Hugo cache); avoid manual edits.

## Build, Test, and Development Commands
- `hugo server`: run the local dev server at `http://localhost:1313` with live reload.
- `make build`: generate the site into `public/`.
- `make clean`: recreate `public/` from the `gh-pages` branch before publishing.
- `make graphs`: regenerate diagrams in `content/docs/change/` and `static/docs/roadmap/` (needs `dot` and related tools).
- `make ship` / `make upstream`: publish to `gh-pages` (maintainers).

## Coding Style & Naming Conventions
- Write content in Markdown and keep YAML front matter minimal and consistent.
- Use kebab-case for content filenames, e.g. `content/posts/git-and-github-collaboration.md`.
- Prefer spaces over tabs; keep lists and headings clean and scannable.

## Testing Guidelines
- No automated test suite today. Validate with `hugo server` and a clean `make build`.
- If you change diagrams, run `make graphs` and confirm generated `.svg`/`.pdf` outputs.

## TODO Tracking
- Track outstanding work in `TODO.md` and keep entries categorized.
- Number TODOs sequentially with 3 digits (e.g., `001`, `002`) and refer to them by number.
- Mark completed items by adding `DONE` after the number.

## Commit & Pull Request Guidelines
- Commit messages are short, imperative, and capitalized (e.g., "Fix typo in pull request instructions").
- In commit bodies, include a section per changed file with bullet points summarizing edits.
- Include issue references like `Fix #23` in the commit message or PR description to auto-close issues.
- If behavior/output changes, include before/after notes or example output.
- Branch names follow `NN_short_description` (e.g., `23_fix_typo_in_index`).
- PRs should include a concise summary, tests run, linked issues, and manual checks (pages reviewed, commands run).

## Agent-Specific Notes
- Prefer small, focused edits; avoid rearranging files without a clear need.
- Do not remove comments or documentation; update them if outdated.
- If you add files or directories, document them here and explain the rationale in the PR.
- Before major changes and after significant milestones, prompt the user to commit.
- Summarize `git diff` output when generating commit messages, not just chat context.
- Treat a line containing only `commit` as: add and commit all changes with an AGENTS.md-compliant message.
- When running `git add`, list individual files, not `git add .` or `git add -A`.
