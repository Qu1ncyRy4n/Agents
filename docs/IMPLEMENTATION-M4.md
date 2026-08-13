# Mogent URL Source Pinning Contract

Status: active implementation contract.

URL sources are explicit Git repositories over HTTP or HTTPS. Source aliases
remain manifest-owned. A URL is never fetched as a side effect of build,
status, browsing, coverage, localization, drift, init, or TUI use. Source:
DI-fipam.

## Durable And Runtime Files

- `agents.yaml` keeps the compact `alias: https://...` form or uses an explicit
  `location` plus optional `subdir` mapping.
- `mogent.lock.yaml` is a strict, deterministic, reviewable file intended for
  version control.
- `.mogent/sources/<alias>/<commit>/` is an ignored checkout/cache.

Each lock entry records URL, normalized subdirectory, full Git commit, and
SHA-256 over sorted Markdown paths and bytes below the selected source root.
Ordinary loads require all four to match.

Example:

```yaml
sources:
  shared:
    location: https://github.com/Qu1ncyRy4n/Agents.git
    subdir: libraries/cdint
```

Scalar source values normalize to the same representation with an empty
subdirectory. Subdirectories are relative slash paths and may not contain
empty, `.`, or `..` components, backslashes, or symlinked components.

## Commands

- `mogent source pin <alias> [--ref <ref>]` resolves and installs the initial
  immutable source. If a lock already exists, it restores exactly that commit
  unless an explicit ref requests a candidate change.
- `mogent source update <alias> [--ref <ref>]` fetches a candidate and prints
  old/new commit plus old/new Markdown content. It writes nothing without
  `--accept`.
- `mogent source update <alias> --ref <full-previewed-commit> --accept` installs
  exactly the reviewed candidate and atomically updates the lock. Moving branch
  or tag names are rejected for acceptance.

Refs beginning with `-`, non-HTTP(S) URLs, non-commit lock revisions, cache path
escapes, duplicate YAML keys, and hash mismatches are errors.

## Git Server Compatibility

Pinning uses ordinary Git smart HTTP(S), not a GitHub-specific API. Public
GitLab, Gitea, Forgejo, and comparable servers should work when they support a
shallow fetch of the requested ref and present a certificate valid for the URL
hostname.

Authenticated/private sources are not yet a supported contract. Mogent rejects
credentials embedded in URLs and does not manage tokens, SSH keys, interactive
prompts, or forge-specific login flows. A preconfigured noninteractive Git
credential helper may happen to work because the Git subprocess inherits its
environment, but users must not rely on that behavior until authentication is
designed and tested explicitly.

TLS hostname, trust-chain, DNS, proxy, VPN, and server-side failures occur below
Mogent. Mogent reports Git's failure and must never recommend disabling TLS
verification as a workaround.

## Acceptance Checks

- Local-source behavior is unchanged.
- URL loads perform no network access.
- Missing lock/cache and changed cache content fail loudly.
- A verified cache builds while offline.
- Update preview writes neither lock nor cache.
- Update preview exposes changed Markdown content, not only file names.
- Accepted updates change the lock only after the candidate checkout and hash
  validate, and only when the accepted full commit matches the fetched commit.
- Changed-content review lists Markdown additions, modifications, and removals.
