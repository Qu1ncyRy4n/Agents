---
tags: [lang/go, testing/unit, testing/deterministic, scope/project]
tldr: Use deterministic Go tests with fixtures and table cases.
priority: 0.9
scope: project
---
# Tests

Use the standard `testing` package. Keep tests deterministic. Prefer fixtures
and table-driven tests when the same behavior has several cases. Add coverage
for new behavior and error paths close to the code they exercise.
