# Semantic Intake

This is the lossless comparison workspace for source-guide content. Organize
new candidates by what they instruct an agent to do, not by the organization
that supplied the source or may eventually publish the result.

## Provisional Review Order

```text
1. workflow                 broadly relevant on nearly every task
2. process                  decisions, escalation, records, and handoff
3. safety                   destructive actions, data, credentials, boundaries
4. engineering              changes, tests, errors, comments, interfaces
5. documentation            audience, authority, examples, public prose
6. tools                    Git, shell, package managers, build systems
7. languages                Go, Rust, Python, Nix, MATLAB, JavaScript
8. communication            response style, accessibility, teaching, personas
9. domains                  research, PromiseGrid, notes, media, static sites
10. repository-local        facts and procedures retained for disposition
11. generated-candidates    useful but non-authoritative generated material
```

This is a review order, not a claim that the categories are mutually exclusive.
Use tags and `see_also` notes when a candidate spans subjects. Do not duplicate
source prose merely to place it in two directories.

## Promotion

After source coverage and owner review, promote a candidate to
`orgs/<owner>/<semantic-area>/`. Record the old intake path and final path in
`PROVENANCE.md`. Promotion is an ownership decision, not automatic cleanup.
