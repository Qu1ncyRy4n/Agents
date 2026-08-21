# TE-sulap: Organization Policy Composition

TE ID: TE-sulap
Date: 2026-08-20
Status: open thought experiment

`tools/mint-handle` is unavailable, so `sulap` is manually assigned under the
existing user-approved exception.

## Question

How can an organization say "repositories under this authority must include
our coding practices and safeguards" without reviving hidden ambient scope
composition?

The desired shorthand has been described approximately as:

```yaml
cdint:
  - from: "*:*:*"
```

The exact wildcard axes are not yet defined. This TE tests the policy meaning
before choosing punctuation.

## Recovered Prior Design

The first-generation design used activated scopes such as `org/acme` and a
"most specific wins" rule. It could apply organization guidance broadly, but a
repository's final instructions depended on ambient tags and hidden merge
precedence. The current explicit manifest replaced that model so a reviewer can
read the selected sources, nodes, and output outline directly.

The requirement use case is still valid. The hidden activation mechanism is the
part not worth restoring.

## Authority Boundary

A library must not make itself mandatory merely by declaring:

```yaml
mandatory: true
```

Libraries are content inputs and may be untrusted. The repository manifest or
an explicitly trusted organization-policy binding must grant policy authority.

## Option A: Required Bundle Visible In Every Repository

```yaml
sources:
  cdint: ../agent-libraries/cdint
  local: .agents/library

policy:
  required:
    - cdint:shared-baseline
    - cdint:engineering/safety

doc:
  - Working Agreement: cdint:shared-baseline
  - Engineering: cdint:engineering
  - Repository Notes: local:repo
```

Worked result: deleting `Working Agreement` makes validation fail because the
required source subtree is no longer represented in the rendered document.

Advantages: completely visible and reviewable. Cost: repeated declarations in
every repository and no central proof that all organization repositories have
the rule.

## Option B: Explicit Trusted Organization Policy

Repository manifest:

```yaml
policy_sources:
  - alias: cdint-policy
    url: https://example.com/cdint/agent-policy.git
    revision: 9b77c4d0
    profile: engineering-repository

sources:
  local: .agents/library

doc:
  - Repository Notes: local:repo
```

Pinned policy profile:

```yaml
profiles:
  engineering-repository:
    sources:
      cdint:
        url: https://example.com/cdint/agent-library.git
        revision: 6a31e829
    required:
      - cdint:shared-baseline
      - cdint:engineering/safety
    placement:
      - source: cdint:shared-baseline
        before: "*"
```

Worked result: Mogent loads only the pinned policy named by the repository,
shows the injected required nodes in `status` and preview output, and rejects a
render that omits them. The policy cannot access an undeclared path or silently
float to a new revision.

Advantages: central policy with explicit local trust and immutable review.
Costs: a new public policy format, policy-source trust rules, placement rules,
and update lifecycle.

## Option C: Ambient Organization Detection

```yaml
applies_to: "github.com/cdint/*"
required: [self:shared-baseline]
```

Mogent could infer the organization from the Git remote and apply matching
rules. This is convenient but fragile: remotes can be renamed, forked, absent,
or intentionally point elsewhere. The effective document is no longer fully
explained by `agents.yaml`. Do not use this as the primary model.

## Wildcard Meaning

Before adopting `*:*:*`, each axis would need a stable definition. Plausible
axes include organization, repository, and path, but the notation gives a
reader no clue which is which. Prefer named match fields if matching is needed:

```yaml
applies_to:
  organization: cdint
  repositories: ["*"]
  paths: ["*"]
```

For the first version, avoid matching entirely: the repository explicitly
names a policy profile. Organization-wide enforcement belongs in CI, which can
check that repositories contain the approved binding and lock revision.

## Required Does Not Mean Trusted Update

`required` answers whether content must be present. It does not prove that a new
revision is correct. Even a required safety module can receive a mistaken or
malicious update.

A safe initial review flow is:

```text
mogent status
  required policy update available: 9b77c4d0 -> c2140a66

mogent inherit cdint:shared-baseline --dry-run
  previewed origin revision: c2140a66

mogent inherit cdint:shared-baseline --accept c2140a66
```

`--required-only` can filter a preview or batch, but should not bypass revision
verification. An organization may later approve a signed revision centrally;
repositories can then inherit that already-reviewed pin without individually
reviewing every line.

## Priority, Urgency, And Requirement

Keep these separate:

```yaml
priority: 0.9       # discovery/sort importance, current 0.0-1.0 field
required: true      # policy says it must be selected
risk: high          # possible future review classification
```

No durable enum-to-float scheme for urgency was found in the repository. The
existing float is only `priority`. `urgent` has little operational meaning
without an update feed, last-checked time, or watcher, so defer it.

## Recommendation

Use option A while the library shape is changing. Design option B only after
real repositories demonstrate repeated required bundles. Keep the effective
policy visible in `status`, tree previews, and rendered provenance even when a
trusted profile later supplies it. Use CI—not implicit Git-remote detection—to
enforce organization-wide adoption.

## Owner Decisions

1. Is explicit per-repository `policy.required` sufficient for the first
   dogfood version?
2. Does organization policy require exact placement as well as presence?
3. Should required content render a compact origin/profile annotation, or is a
   text-visible `status` tree plus committed manifest enough?

