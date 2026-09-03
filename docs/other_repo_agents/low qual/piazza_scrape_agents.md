# AGENTS.md - piazza scrape

## Project

Small Python 3.12 utility that exports all posts from a Piazza course to one Markdown file using `piazza-api`.

## Structure

- `piazza_scrape.py` - single-file CLI scraper.
- `pyproject.toml` - package metadata and `piazza-api` dependency.
- `piazza_posts.md` - generated export; treat as output data, not source.

## Development

- Keep credentials out of the repo and out of logs.
- The script supports `--password`, but prefer the macOS Keychain path when possible.
- Preserve `--limit`; it is useful for safe smoke tests against a live Piazza course.
- Keep request pacing jittered to avoid hammering Piazza.

## Commands

```bash
uv run python piazza_scrape.py --help
uv run python piazza_scrape.py --email USER --network-id NETWORK --limit 5 --output sample.md
```

The second command requires valid credentials and network access.

## Behavior To Preserve

- Export top-level posts, instructor answers, student answers, and follow-up threads.
- Keep Markdown output UTF-8.
- Do not commit real course exports or credentials.
