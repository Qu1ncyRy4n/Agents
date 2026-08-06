---
tags: [runtime/hygiene, testing/deterministic, scope/repo]
tldr: Keep temp files and caches out of the repo and avoid accidental network tests.
priority: 0.8
scope: repo
---
# Runtime Hygiene

Put temporary files and build caches under `/tmp`, not in the repository. Keep
tests deterministic and avoid network calls unless the task explicitly needs
them.
