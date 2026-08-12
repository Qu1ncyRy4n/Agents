# TE: Pinned URL Source Subdirectories

TE ID: TE-vurap
Status: concluded by DI-vurap

## Question

How should one pinned Git repository expose a library rooted below the checkout
without changing established short source references or indexing unrelated
Markdown?

## Alternatives

### Encode the subdirectory in the URL

Fragments or query parameters make Git identity and library-root identity hard
to distinguish, interact poorly with ordinary Git URLs, and are easy to omit
from lock validation.

### Require one repository per library

This gives a clean root but forces cross-repository ownership and release
decisions merely to work around a resolver limitation. Separate repositories
may still be desirable when ownership or release cadence differs.

### Add an explicit normalized source object

Keep existing scalar sources as the compact form. Accept an explicit mapping
with `location` and optional `subdir`. Bind the normalized subdirectory into the
lock, hash only Markdown below that root, and return that root to the library
indexer.

## Safety Walkthrough

- Reject absolute subdirectories, empty path segments, `.`, `..`, and
  backslashes.
- Reject a missing/non-directory root and symlinks in any subdirectory
  component.
- Keep the full verified checkout in the ignored cache, but scope hashing,
  browsing, update review, and source resolution to the locked subdirectory.
- A changed subdirectory is a changed source identity and requires explicit
  repinning; ordinary commands never reinterpret an existing lock.
- Compact local paths and URLs continue to parse and serialize as before.

## Conclusion

Use the explicit source form:

```yaml
sources:
  shared:
    location: https://example.com/agents.git
    subdir: libraries/cdint
```

Normalize scalar and explicit sources to the same internal representation. Add
`subdir` to each remote lock entry and include it in URL/lock identity checks.
Source references remain relative to the selected library root.
