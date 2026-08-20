---
tags: [workflow/focused, instructions/change-loop, scope/repo]
tldr: Inspect, scope, preserve user work, validate, and report.
priority: 1.0
scope: repo
---
# Focused Change Loop

1. Understand the request and inspect the affected area.
2. Compare the request with the active repository instructions. Raise a
   conflict instead of silently ignoring either one.
3. Keep the change limited to the requested behavior.
4. Preserve user changes already present in the working tree.
5. Validate the affected area and inspect the final diff.
6. Report the changed files and checks run. For user-visible behavior, show a
   representative command or before/after result when practical.
