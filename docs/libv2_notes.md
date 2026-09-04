# exemplar

## Repo guidelines


### project structure and mod org

issue: todo loc? base or not?

no x/. Stick with avoid internal. use more universal than `.grok`?
I'm guessing proqint in `TODO` is the most up to date

### build test, dev commands
pg: no root go modules; commands per module. “If Go commands fail because the local go binary and compiled standard-library objects are from different Go versions, report the environment blocker instead of changing repo files.” from promise-grid

alt: makefiles,

this might make more sense in spec.

### Agent instruction arch

relevant now? since managed by mogent? basically replace by: multi file targets: if global, define in general, add additions per file target.

user level / role level overlayers?

pg and fab26 pres: mentions of public artifacts: maybe not relevant here.

### Promise action minimalism

Should this be in a skill / spec doc?

Think this might need updating from steve

more generalized rule: iterations / addative progression instead of rewrite? Again seems like a project choice, not a generalized thing.

### Dev guide resources

again, repo specific

### Decision first:
great,

Diff noted / options:
1. Strick decision first
2. Risk-based escalation -> decision
3. 'lightweight' lock before behavior change.

1 by default; manual opt out if one needs (localization)

the minting thing: is that existing code that must be imported? might be nice to have a gh link. Maybe a mogent tool?
