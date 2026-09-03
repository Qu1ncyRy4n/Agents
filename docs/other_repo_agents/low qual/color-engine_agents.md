# AGENTS.md - color-engine

## Project

Nested Rust crate intended to become the pure color-math engine for `biome-theme`.

## Current State

- Crate name: `color-engine`.
- Rust edition: 2024.
- Current `src/lib.rs` is a rough scaffold and does not compile as written.
- It is also included as a member of the parent `biome-theme` Cargo workspace.

## Intended Scope

- Pure color/image algorithms only.
- No network, CLI, iNaturalist API, shapefile, template, or GUI logic.
- Planned responsibilities include image decoding, RGB to OKLCH conversion, median cut quantization, and diversity/orthogonality sampling.

## Development Guidance

- Keep APIs small and testable before wiring Python or CLI layers.
- Prefer deterministic algorithms and fixtures for tests.
- When adding dependencies, update the crate manifest deliberately and check whether the parent workspace manifest also needs changes.
- Avoid PyO3 or FFI in this crate unless the project phase explicitly moves to bindings.

## Commands

From this crate:

```bash
cargo check
cargo test
```

From the parent workspace:

```bash
cd /Users/q/Documents/Code/misc_tools/chroma_terra/biome-theme
cargo check
cargo test
```
