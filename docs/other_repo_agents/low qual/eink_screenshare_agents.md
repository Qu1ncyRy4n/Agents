# AGENTS.md - eink_screenshare

## Project

Early design-stage project for turning an Android e-ink tablet into a real extended macOS monitor over USB. The Mac side is planned in Rust; the Android side is planned in Kotlin.

## Current State

- Design only. `README.md` and `SPEC.md` are the authoritative documents.
- Development environment is defined by `flake.nix`.
- No source implementation is present yet.

## Development

```bash
nix develop
```

Use the Nix shell for toolchain work. macOS is required for the Mac virtual-display runtime because it depends on private CoreGraphics APIs; Linux is only suitable for editing/building non-macOS pieces.

## Architecture Constraints

- v1 is wired USB via `adb reverse`; do not add Wi-Fi/Bluetooth behavior unless explicitly requested.
- v1 is read-only display output; pointer/touch input is phase 2.
- Keep `CGVirtualDisplay` and other private CoreGraphics FFI isolated behind one module when implementation starts.
- Optimize for e-ink text use: 1-bit dithered frames, tile-diff updates, low refresh, periodic full refresh.
- Do not optimize for video or high frame rates.

## Protocol

Follow `SPEC.md`: length-prefixed little-endian messages with `HELLO`, `TILE_BATCH`, and `FULL_REFRESH` for v1. Add compression only after measuring USB throughput.

## Validation

- For design changes, update `SPEC.md` and `README.md` together.
- For future Rust implementation, run checks inside `nix develop`.
