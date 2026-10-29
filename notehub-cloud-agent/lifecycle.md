# Lifecycle, recovery and correction

Begin with the assistant's database records and the actual scheduler state.
Use job contracts to resume work, and Setup only for missing
capabilities or a new scope; preserve functioning jobs and checkpoints.

## Choose the relevant reference

Read the rows that apply to the current step; reuse context already established.

| If you need to | Read | Then |
|---|---|---|
| Recover the assistant's configuration, jobs and receipts | [Database state](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/workspace.md) | Reconcile durable records with actual scheduler state |
| Resume or repair scheduled execution | [Job contracts](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/jobs.md) | Preserve valid checkpoints and verify the repaired run |
| Enroll another project or repair the cache infrastructure | [Cache preparation](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/cache-preparation.md) | Preserve existing workers and checkpoints while adding or repairing the affected scope |
| Add requested project questions or monitoring | [Setup](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/setup.md) | Complete only the missing setup and preserve functioning jobs |
| Make an authorized correction to shared product meaning | [Training](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-train/SKILL.md) | Review and publish the shared correction, then recheck affected private plans |

Use stored state and the actual harness to determine what runs, its scope, progress and
how to stop it. Add recovery state only as needed; distinguish pending deliveries from
successful receipts.

Reconcile the Notehub Cache Database and Notehub Cache Maintenance Task created by Access and every enrolled
project. Each run loads the current registry; first-use enrollments are not lost after a
restart. Pause/removal of a project excludes it from future fill and maintenance until
the user changes that state. Switching the active conversation leaves other enrolled
projects running. A failure or loss of permission in one project does not erase another
project's cache or reset its cursors.

At startup, verify the selected connection and project binding, acquire the writer lock,
load checkpoints and schema versions, and resume committed work. Failed or expired auth
pauses acquisition and reports the specific problem without repeatedly initiating login.
A changed hub/account/filter creates a new scope, not reuse of an old cursor. Restore
from committed state after a crash; never skip a batch because it was fetched before the
crash. Rebuild corrupt or incompatible source caches/derived projections in staging and
verify them before switching readers. Preserve valid operational state. If that state is
also damaged, recover it through the supported consistent backup mechanism and reconcile
actual jobs and deliveries; do not silently reset pauses, incident baselines or receipts.

Poller failures leave last-good data available with its actual freshness/coverage. Rate
limit failure notices, but escalate a sustained loss of coverage relevant to active
monitors. Recovery can clear a coverage incident after a complete successful scan; it
does not clear an unrelated product incident without fresh product evidence.

When a product rule or user correction changes, identify affected recipes, views, caches,
monitors and already delivered conclusions. Preserve raw evidence, rebuild derivations,
re-evaluate dependent monitors and correct consequential prior messages. Record what
changed and the scope/version of the new rule. Don't promote a user's private preference
into shared product knowledge; use training only for an authorized shared correction.

Pause means no new monitoring delivery; cancel scheduled jobs when asked to stop them
and verify cancellation. Removing one interest need not destroy the useful Q&A cache.
Remove credentials only on an explicit logout/revocation request through their owning
tool. Data deletion is a separate choice; document retention of cache, drafts and delivery
history using the harness's actual controls. On handoff, preserve nonsecret configuration,
job IDs, revisions, open incidents, acknowledged deliveries and unresolved coverage gaps.
