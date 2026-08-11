# TE: URL Source Pinning And Change Review

TE ID: TE-fipam
Status: concluded by DI-fipam

## Decision Under Test

How can a manifest name a remote Git source without silently consuming mutable
network content during build, status, discovery, coverage, localization, or TUI
use?

## Alternatives

### Fetch the URL during every load

This is convenient but non-deterministic, network-dependent, and unsafe. A
branch can change between inspection and build.

### Put a commit fragment directly in every source URL

This is explicit but makes routine updates rewrite the manifest and provides no
content-integrity check or local review workflow.

### Commit a lockfile and ignore verified checkouts

Keep the user-defined URL in `agents.yaml`, record the resolved commit and
Markdown-content hash in `mogent.lock.yaml`, and store its checkout under
ignored `.mogent/sources/`. Normal workspace loads never fetch; they require the
lock and verify the checkout before indexing it.

## Change Review

Initial `source pin` is an explicit network and lockfile operation. Later
`source update` fetches a candidate into temporary storage and reports old/new
commit plus changed Markdown content. Without `--accept`, it writes nothing.
Acceptance must repeat the full previewed commit, so a branch cannot move
between review and installation. The verified candidate then updates the cache
and lock.

Builds do not automatically refresh missing caches. The user can restore one
with the explicit pin command, which verifies it against the existing lock when
the lock already names an immutable commit.

## Conclusion

Use `mogent.lock.yaml` beside the manifest and `.mogent/sources/<alias>/<commit>`
for ignored checkouts. Lock URL, commit, and a deterministic hash over Markdown
paths and bytes. Reject URL/lock mismatches, missing caches, hash mismatches,
unsafe aliases, and moving content during every ordinary load.
