# Notes Obsidian Style And Agent Skills Reference

Reference compiled from the Notes vault `Agents/skills/`, `Templates/`, and `ZK/ZK and obsidian style guide and synatx.md`.

This file is not a replacement for the original vault files. It is portable behavior/style context for agents working with the Notes Obsidian vault.

## Obsidian Style

Use Obsidian-native Markdown:

- YAML frontmatter at top for metadata.
- Tags in YAML and/or top-level body tags, depending on note type.
- Wikilinks for internal references: `[[note name]]`, `[[note name#Heading]]`.
- Embedded references where useful: `![[Tasks_day#Today]]`.
- PDF/page links are valid citation targets: `[[file.pdf#page=7]]`.
- Prefer heading links for precise references.

## ZK Style

Atomic ZK notes:

- Title should be the claim itself, not a broad topic label.
- Body should elaborate the claim in 3-10 sentences.
- Include links to related notes in `ZK/` and `Reference_notes/`.
- Use tags appropriate to domain.
- Save under `ZK/[kebab-slug-of-title].md`.

Atomic note template:

```yaml
---
tags:
  - atomic
date/time: "{{date:YYYYMMDD}}{{time:HHmm}}"
---
```

Then body content, followed by:

```markdown
## References
```

Fleeting note template:

```yaml
---
tags:
  - "#fleeting"
  - "#ZK"
date/time: "{{date:YYYYMMDD}}{{time:HHmm}}"
---
```

```markdown
> [!QUOTE] Original Capture
> {{VALUE}}

## References

- See `[[this note here]]` for xyz
```

Daily note template:

```yaml
---
tags:
  - periodic/daily
date/time: "{{date:YYYYMMDD}}{{time:HHmm}}"
---
```

```markdown
# TODO
![[Tasks_day#Today]]

weekly note: [[{{date:GGGG-MM-[W]WW}}]]

# Daily
```

## Flashcard Syntax

From the vault style guide and flashcards skill:

Single-line one-direction card:

```markdown
First side :: second side
```

Single-line bidirectional card:

```markdown
Side 1 ::: Side 2
```

Multi-line card:

```markdown
Front
?
answer line 1
answer line 2
```

Multi-line bidirectional/synthesis card:

```markdown
Front
??
answer line 1
answer line 2
```

Rules:

- One card equals one retrievable fact.
- Do not put blank lines inside answer blocks.
- Do not separate cards with `---`.
- Use `?` for multi-point cards.
- Use `??` for bidirectional synthesis cards.
- No answer should give away another card's answer.
- Markdown formatting and Obsidian links are valid inside cards.

Inline cloze:

```markdown
The hippocampus is a ==medial temporal lobe== structure.
```

Cloze rules:

- Highlight only the retrieval target.
- Keep sentence readable.
- Use cloze for terms and short facts.
- Use cards for multi-point answers.

## Codex/Flashcards Document Style

Use when building reviewable technical reference documents or study decks.

Document structure:

```markdown
---
title: [Title]
created: YYYY-MM-DD
updated: YYYY-MM-DD
status: draft | active | stable
codex_domain: [domain/topic]
course: [course code if relevant]
---

#flashcards/[semester]/[course]/[unit]
#codex/[field]/[topic]

## TODO:

## Index:

## Connections:

---

## [Topic Section]
```

Every document should have:

- YAML frontmatter.
- Flashcard tag hierarchy.
- Codex tag hierarchy.
- TODO block for missing images/audio, memorization warnings, uncertain content, gaps.
- Index with heading wikilinks.
- Connections block.
- Topic sections with outline, explainer prose, cloze, and cards.
- Specific citations to lecture files, readings, PDFs, or headings.

Connections block format:

```markdown
## Connections

**Links to:** [[codex/dopamine.md]], [[codex/memory-systems.md]]
**Referenced by notes:** [[why-hippocampal-replay-needs-pattern-sep-and-completion.md]]
**See also:** [[codex/amygdala.md]] (parallel emotional memory system)
```

Callouts:

```markdown
> [!caution] Check this: [what seems off and why]

> [!warning] Heavy memorization: [what needs drilling]

> [!note] Gap: [what is missing]
```

## Skill: new-zk

Purpose: create a new atomic ZK note.

Behavior:

1. If argument is provided, use it as topic or claim.
2. If no argument, ask: `What's the claim?` Wait for answer.
3. Check `ZK/` and `Reference_notes/` for related notes.
4. Scaffold frontmatter:

```yaml
---
title: [claim as title]
type: atomic
date: [today's date]
tags:
  - [domain/subtopic]
links: []
---
```

5. Save to `ZK/[kebab-slug-of-title].md`.
6. Report path and claim. Ask whether content needs filling in or body should be drafted.

## Skill: process-note

Purpose: analyze one note and propose a destination/action.

Input: path or filename.

Analyze:

1. Type: fleeting capture, literature note, reference material, already-atomic claim, lab/research note, periodic note.
2. Destination: `ZK/`, `Reference_notes/[subdir]`, `inbox/periodic archive/`, or current location.
3. Split needed: whether multiple atomic ZK notes should be created.
4. Links: related notes in `ZK/` and `Reference_notes/`.

Then propose specific action and wait for confirmation before moving, splitting, or rewriting.

Consult `Agents/Agent_vault_map.md` for current structure when needed.

## Skill: youare

Purpose: activate a vault agent persona.

Argument: agent name or partial match, e.g. `og`, `claudia`, `vera`.

Steps:

1. List subdirectories in `Agents/`.
2. Find best partial match.
3. Read matching personality file.
4. Read `Agents/Agent_instructions.md`.
5. Load memories with:

```bash
grep -rl "agents/$ARGUMENTS" Agents/Memories/
grep -rl "agents/shared" Agents/Memories/
```

6. Adopt persona from that point forward.

If no match, list available agents and ask. If ambiguous, list matches and ask.

## Skill: vault-status

Purpose: give current vault state.

Steps:

1. Read `Agents/Agent_vault_map.md`.
2. Read `Agents/Memories/_index.md`.
3. Grep `Agents/Memories/` for recent entries, last 5 files by date in filename.
4. Report:
   - Current vault priorities.
   - Recent memory notes.
   - One suggested next action.

Be concise. No preamble.

## Skill: obsidian-flashcards-codex

Purpose: generate Obsidian Spaced Repetition study documents that double as Codex entries. Use for flashcards, study decks, spaced repetition notes, lecture/readings conversion, or technical Codex entries.

Core idea: same file is both reference and review deck.

Delivery checklist:

- YAML frontmatter: title, created, updated, status, codex_domain, course.
- Both tag hierarchies at top.
- TODO block.
- Index with heading wikilinks.
- Connections block.
- Each section has outline, prose with cloze, and cards.
- Lecture files linked above content they cover.
- Inline citations are specific.
- Cloze only on key terms.
- Cards are atomic.
- Multi-point answers use `?` or `??`.
- No blank lines inside answer blocks.
- No `---` between cards.
- Callouts for uncertain content, memorization warnings, and gaps.
- Summary `??` card for each major concept.

