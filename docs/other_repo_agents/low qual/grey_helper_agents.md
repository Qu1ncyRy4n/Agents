# AGENTS.md - grey_helper

## Project

Rust 2024 native macOS receipt processor. It reads UPS shipping receipt PDFs, extracts shipment data, writes `shipments.csv`, and moves processed PDFs to a `processed/` subfolder.

## Structure

- `src/main.rs` - egui app entry point.
- `src/app.rs` - UI layout, buttons, log, and folder picker.
- `src/processor.rs` - pipeline orchestration.
- `src/extractor.rs` - UPS/FedEx-style receipt parsing.
- `src/validator.rs` - duplicate detection and date-gap checks.
- `src/csv_writer.rs` - CSV append and existing-record loading.
- `src/pdf.rs` - PDF text extraction and carrier detection.
- `src/shipment.rs` - shipment data model.
- `build_app.sh` - builds `dist/Receipt Processor.app`.

## Development

```bash
cargo run
cargo test
cargo check
bash build_app.sh
```

Use `cargo run` for UI smoke tests on macOS. Use `build_app.sh` only when validating app bundle behavior.

## Behavior To Preserve

- Default output location is `~/Desktop/ShipmentProcessor/` unless the user changes it.
- Processed PDFs are moved to a `processed/` subfolder.
- CSV output should remain append-friendly and duplicate-aware.
- Keep parsing changes covered by focused fixtures/tests when touching receipt extraction.
- Do not edit generated `dist/` artifacts unless explicitly packaging a release.

## Notes

First launch of an unsigned app may require right-click then Open. Keep that user-facing caveat in docs if packaging instructions change.
