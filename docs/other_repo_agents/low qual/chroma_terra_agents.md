# AGENTS.md - chroma_terra

## Project

Container repo for `biome-theme`, a planned CLI/GUI tool that generates base16 color schemes from local ecology using iNaturalist photos and perceptual color sampling.

## Important Context

- `biome-theme-claude-code-prompt.md` is the main project brief and collaboration contract.
- `biome-theme/` contains a Rust workspace.
- `biome-theme/crates/color-engine/` is also its own nested git repo; treat changes there as belonging to that repo too.

## Collaboration Style

- Prioritize the user's Rust learning over speed.
- When asked how to do something, offer 2-3 approaches with tradeoffs before choosing.
- Flag architectural choices that will be hard to reverse.
- Suggest tests before implementation when an interface is unclear.
- Keep code changes small enough to explain clearly.

## Planned Stack

- Rust color engine for median cut, OKLCH, and orthogonality sampling.
- Future PyO3 bridge to Python.
- Python 3.12+ with `uv` for iNaturalist, geospatial lookup, scheme generation, and CLI.
- Base16/tinted templates for output.

## Commands

```bash
cd biome-theme
cargo check
cargo test
```

## Project Hygiene

- Add or update `DECISIONS.md` for major architecture or design choices.
- Reference real crate docs for API details.
- Do not invent the Python package tree until the specific phase needs it.
