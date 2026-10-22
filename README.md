# Mogent

Mogent composes project instructions and skills from reusable libraries.

Use mogent.build for an authored selection, or mogent.plan when selection intent
should be sequenced into a generated build. Sources are declarations; select
modules explicitly before building.

Start with mogent init, then use source add, source show, module add, and build.
Inspect a project with show --tree, status, drift --all --diff, and optimize.

Drift compares managed outputs with disk and the current build. It also compares
connected libraries with source fingerprints saved by the last build, reporting
changed, added, and removed modules.

See docs/CLI-COMMAND-INVENTORY.md, docs/SOURCE-DISCOVERY.md, and
docs/convo/head.md.
