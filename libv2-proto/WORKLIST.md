# Protolibrary Worklist

This connects the provenance ledger to the order of extraction and owner
review. It is not an `agents.yaml` manifest and does not define rendered output.

Current cross-project head: `~/dev/Notes/inbox/2026-08-W34.md`.

## Coverage Gates

- [x] Preserve v1 while v2 is developed separately.
- [x] Define provisional intake vocabulary and metadata.
- [x] Add all 37 source-guide files to the provenance ledger.
- [ ] Map every useful source-guide heading to an explicit disposition.
- [ ] Extract every useful instruction before pruning or deduplication.
- [ ] Record every duplicate, conflict, option, and possible output destination.
- [ ] Review each semantic area with the owner.
- [ ] Promote reviewed candidates into organization-owned libraries.
- [ ] Validate representative manifests and rendered guides.
- [ ] Decide whether and how to retire v1.

## Review Queue

```text
NOW
├─ workflow/edit boundaries
│  └─ Keep Changes Scoped and Preserve Existing Work
├─ workflow/decision threshold
│  ├─ Use Risk to Choose Routine Work or Decision Review
│  └─ Require Human Review for High-Consequence Changes
├─ process/decision governance
│  ├─ Use Thought Experiments to Narrow a Broad Design Space
│  ├─ Lock Durable Decisions Before Implementation
│  ├─ Preserve Decision History Through Supersession
│  └─ Handoff Decisions with Implementation Evidence
├─ safety
├─ engineering
├─ documentation
├─ tools
├─ languages
├─ communication
├─ domains
├─ repository-local material
└─ generated candidates
```

Within each area, review broadly applicable and every-prompt guidance first,
then conditional overlays, organization-specific policy, and repository-local
facts.

## Current Decision Queue

- [ ] Decide whether edit scope and preservation of existing work form one
  compatible bundle or two atomic modules.
- [ ] Decide whether risk-based escalation is the canonical workflow default.
- [ ] Decide whether strict decision-first remains one selectable complete
  bundle or contributes independently selectable stages.
- [ ] Decide which decision-record requirements belong in agent guidance versus
  a specification or generated reference.
- [ ] Decide the ordinary versus governed handoff boundary.
- [x] Run the broad include/composition narrowing pass in `TE-nufad`.
- [ ] Decide whether exact mid-document insertion is required for the first
  libv2 cutover.
- [ ] If exact insertion is required, prototype frontmatter plus a minimal
  marker and a restricted literal Go-template include.
- [ ] Decide whether choice rules live in module frontmatter, a source-level
  descriptor, or the consuming manifest according to who owns the relationship.
- [ ] Separately decide YAML-map merge precedence; do not conflate inheritance
  overrides with module requirements or conflicts.
- [ ] Later: decide source-path move compatibility and migration behavior.

## Per-Candidate Review

For every candidate, answer:

1. Is all useful source material present and cited?
2. Is source wording clearly separated from adaptation and proposed examples?
3. Does the title state both subject and direction?
4. Is this one compatible selection unit?
5. What overlaps with it, and may the overlapping modules coexist?
6. Is its best output an agent guide, skill, generated document,
   specification, library-authoring guide, or repository-local instruction?
7. Who should own the reviewed canonical version?
8. Keep, adapt, combine, split, retain as a choice, or reject?
