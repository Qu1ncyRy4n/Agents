# Conversation Work

The one current resume action lives in [head.md](head.md). This ledger owns task
status; its ACTIVE ID must match the head. IDs do not change when items move.
Product backlogs remain in the repository's `TODO/`; do not reproduce them all here.

| ID | State | Item | Depends on | Evidence |
|---|---|---|---|---|
| CS-W001 | DONE | Find and read the supplied conversation stack | none | Pinned source linked in archive/config-redesign/README.md; all 871 lines read |
| CS-W002 | DONE | Draft the reusable convo-sequence skill | CS-W001 | .agents/skills/convo-sequence/SKILL.md; metadata validator passed |
| CS-W003 | DONE | Review the salient NOW entry point and shared folder layout | CS-W002 | Owner said "great" and moved on to Seam configuration |
| CS-W004 | PARKED | Agree and implement the small read-only consistency checker | CS-W003 | One active ID, queue/pointer/status/supersession checks; exact contract still under review |
| CS-W005 | PARKED | Record the accepted workflow and validate a re-entry example | CS-W004 | Proposed scenarios in the skill's references/baseline.md |
| CS-W006 | PARKED | Return to the Mogent demo/config-design discussion | none | Owner said to come back to the demo later; recover product context from the pinned archive |
| CS-W007 | PARKED | Continue the three-repo deployment campaign | CS-W006 | docs/dogfood/2026-10-09-191238-mogent-multirepo-deploy.md |
| CS-W008 | PARKED | Resume, revise, or discard the operation-inference patch after design review | CS-W006 | Preserved in dedicated stash 5e826857a3a64633be5bf2625cd22e827c2a30af; not adopted into main |
| CS-W009 | DONE | Update seam_game configuration and include convo-sequence | CS-W003 | Owner approved local guide source; new plan/lib/build, complete skill installed, build --check passed; guide and skill hashes preserved |
| CS-W010 | DONE | Review Seam's completed deployment with the owner | CS-W009 | Owner proceeded to request relevant QMR modules |
| CS-W011 | DONE | Import relevant QMR guidance/skills into Seam | CS-W010 | Owner requested localization; 4 shared sections and 7 localized QMR skills deployed alongside convo-sequence; source preservation and build --check passed |
| CS-W012 | DONE | Make convo-sequence chat concise and one independently actionable step at a time | CS-W002 | Added explicit chat/validation guidance and response-style example; Seam skill rebuilt; metadata and build checks passed |
| CS-W013 | PARKED | Review Seam's localized shared-plus-project deployment with the owner | CS-W011 | mogent.plan, .mogent/local/mogent.lib, agent-library/LOCALIZATION.md; no automatic QMR recopy during future builds |
| CS-W014 | DONE | Pull and reconcile Mogent/QMR/Seam while preserving local work | CS-W011 | Mogent at d1235e1; QMR merge 2380b3c; local work restored; tools/check, QMR stories, receipt resolution and Seam build check passed; recovery stashes retained |
| CS-W015 | DONE | Review reconciled state and next commit scope with the owner | CS-W014 | Owner chose to commit docs/skill and park the unapproved code |
| CS-W016 | DONE | Move conversation records under docs/convo and clean this worktree | CS-W015 | Owner requested docs or .agents; docs/convo holds authored records, .agents holds the reusable skill; paused code stashed |
| CS-W017 | ACTIVE | Review the clean checkpoint and choose the next task | CS-W016 | Resume from docs/convo/head.md; no automatic restart of parked implementation |

Exactly one item is ACTIVE. Keep the NEXT queue small; all other open topics are
PARKED. A blocked active item stays visible rather than silently spawning a second
active thread. Keep completion evidence when moving items to DONE.
