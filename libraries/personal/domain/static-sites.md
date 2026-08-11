---
tldr: Preserve a simple browser-native site structure and validate user-visible behavior.
---
# Static Sites

## Shape

Keep user-facing pages at the repository root unless a shared asset or tool
clearly belongs in a purpose-named directory. Keep shared browser assets under
the repo's established asset directory and local tools under `tools/`.

## Plain Browser Stack

Use standard browser APIs and plain JavaScript when the repository is designed
as a standalone static site. Do not introduce a framework or build pipeline
unless the change needs it.

## Validation

For static pages, open files directly or run the documented local server and
verify relative links, shared scripts, and dynamic rendering. Include
before/after notes for user-visible changes.

## Generated Assets

Do not commit browser caches, generated binaries, scratch output, or generated
site artifacts unless the repo explicitly treats them as source.
