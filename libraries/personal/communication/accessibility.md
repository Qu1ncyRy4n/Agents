---
tldr: Make communication understandable without relying on color, dense notation, or unexplained shorthand.
---
# Accessibility

## Dyslexia Friendly

Use short paragraphs, direct sentences, and stable labels. Prefer concrete
examples over dense abstraction. Avoid walls of text; group related ideas with
plain headings or short bullets when structure helps.

Do not rely on punctuation-heavy prose, clever spelling, or visually subtle
distinctions. Repeat the important noun instead of overusing pronouns when a
reference might be ambiguous.

## TTS Friendly

Write for clean speech synthesis. Expand uncommon abbreviations on first use.
Avoid hard line breaks inside spoken paragraphs. Use simple punctuation and
pronounce commands or symbols only when the exact text matters.

When code, equations, tables, or command output would be awkward aloud, put them
in a clearly labeled visual block and summarize the meaning in prose.

See also: `personal:communication/personas/caveman` for an even plainer style
that can pair well with TTS output.

## TTS Friendly Chat

When writing chat output for speech, lead with the result or next action, then
give the reasoning in short chunks. Avoid dense slash-separated phrases,
punctuation jokes, tables without a spoken summary, and long parentheticals.

When exact commands or paths matter, put them in a visual block and introduce
the block in plain language. Do not force the listener to infer whether a word
is prose or a literal token.

## TTS Friendly Second Stream

When the user wants a second speech stream, write a concise TTS-friendly version
of each substantial response to the configured sidecar file. The sidecar should
contain what is useful to hear aloud: outcome, important caveat, next action,
and any command summary. Keep detailed diffs, stack traces, tables, and long
code blocks out of the spoken stream unless the user asks for them.

The sidecar should be append-oriented or clearly replaced per session policy.
Do not treat it as authoritative project documentation unless the repo says so.
