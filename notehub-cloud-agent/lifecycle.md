# Verify, recover or stop existing work

Use the procedure matching the requested change or observed failure. Start with the
saved scope, definitions, checkpoints and actual scheduler state. Preserve working
components outside that scope, including other projects' caches and monitors.

Choose one procedure; these are independent responses to different situations.

| Situation | Procedure | Complete by |
|---|---|---|
| The user asks to pause, resume or remove a monitor, interest or maintenance scope | [Pause, resume or remove work](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/lifecycle.md#pause-resume-or-remove-work) | Verifying the change without affecting unrelated work |
| The runtime restarted or registrations disappeared | [Restore schedules](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/lifecycle.md#restore-schedules-after-a-restart) | Recovering each active schedule once and preserving pauses |
| Acquisition failed or its saved state is damaged | [Resume acquisition](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/lifecycle.md#resume-acquisition-after-a-failure) | Resuming committed progress with coverage gaps explicit |
| Product meaning or an earlier interpretation changed | [Reconcile meaning](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/lifecycle.md#reconcile-a-changed-product-rule-or-correction) | Rechecking affected calculations and correcting consequential conclusions |

## Restore schedules after a restart

1. Recover the durable worker/runtime, task definitions and database locator. Use the
   persistent scheduler or startup hook established during harness preparation. The
   hook must run even if all periodic registrations disappeared; it needs only enough
   bootstrap information to locate the database and invoke recovery.
2. Compare definitions with the actual scheduler. Reuse matching registrations, restore
   missing active tasks and repair outdated invocations. Preserve paused/cancelled work;
   remove only its dedicated stale registrations. Keep shared workers needed elsewhere.
   Coordinate recovery across bots so concurrent starts cannot duplicate schedules.
3. Check the runtime and authenticated invocation. Resume committed progress and retain
   pending findings/deliveries. Record next runs and any failures. A failed project does
   not erase another project's cache or reset its cursors.
4. Verify with a disposable task: remove its registration, invoke recovery and confirm
   exactly one registration returns. Repeat recovery and test paused/cancelled cases.
   Verify the startup-hook registration or scheduler persistence from the harness;
   do not restart the user's whole Cloud Computer merely to test this.
5. Save the recovery result. If automatic startup is unavailable, retain the explicit
   recovery procedure and report that limitation.

**Complete when:** expected active schedules exist once, pauses remain intact and recovery evidence is saved.

## Resume acquisition after a failure

1. Verify the selected account/project, database access mechanism and stream definitions.
   Failed or expired authentication pauses that scope and reports the problem without
   repeated login. A changed account/filter needs its own scope and checkpoint.
2. Load committed checkpoints and schema versions. Resume from those checkpoints, never
   from a batch merely fetched before the crash. Read the current registry so newer
   enrollments are included. Keep paused/excluded projects inactive.
3. Preserve last-good data with its real freshness and coverage. If source data or a
   projection is corrupt, rebuild in staging and verify it before switching readers.
   Preserve operational state. If that state is damaged too, use the supported backup
   and reconcile actual jobs/deliveries; do not reset pauses, baselines or receipts.
4. Verify a successful bounded run and its continuation/next run. Rate-limit failure
   notices while surfacing sustained loss of coverage. A recovered poller does not
   prove a monitored condition has cleared; that needs fresh product evidence.

**Complete when:** work resumes from verified state, with remaining coverage gaps explicit.

## Reconcile a changed product rule or correction

1. Identify affected recipes, projections, queries, monitors and delivered conclusions.
   Preserve raw evidence and the scope/version of the changed meaning.
2. Rebuild and verify dependent calculations, re-evaluate monitors and correct consequential
   prior messages. Preserve unrelated working plans.
3. Keep private preferences separate from shared product knowledge. If the user requested
   a shared correction, complete it through Training before relying on its publication.

| For a requested shared correction | Resource | Then |
|---|---|---|
| Review and publish changed meaning | [Training](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-train/SKILL.md) | Return with verified revisions and reconcile dependent private work |

**Complete when:** affected work uses verified meaning and prior consequential errors are addressed.

## Pause, resume or remove work

1. Resolve the requested project, monitor or interest and record the desired state.
   Ending a conversation or switching projects alone does not cancel background work.
2. On pause/removal, suppress the affected evaluation and delivery; cancel dedicated
   schedules when requested and verify cancellation. Keep shared dispatchers needed by
   other active definitions. Paused projects stay out of acquisition and maintenance.
3. On resume, verify the current definition, permissions and fresh evidence, then restore
   only missing execution. Re-evaluate rather than replaying stale queued news. A removed
   interest must not recreate itself from old configuration.
4. Verify the resulting scheduler and delivery state. Removing an interest need not
   remove the Q&A cache. Data deletion is a separate choice; describe actual retention
   controls. Remove credentials only for an explicit logout/revocation through their tool.

**Complete when:** the requested change is reflected in both durable state and actual execution.

## Return or hand off

Return a brief account of what changed, what runs next and any remaining dependency.
For handoff, preserve nonsecret configuration, job IDs, meaning revisions, monitor
preferences and progress, and unresolved coverage gaps. Keep only what the next
run or agent needs; a recorded configuration is not proof that execution succeeded.

## Grok Bot only

*This section applies only to agents running as Grok Bot. Other harnesses should skip it.*

- The pure-Python scheduler loop does not survive a box restart, and there is no startup
  hook. Recovery means relaunching the loop script named in the cache README, guarded by
  its lock so it cannot run twice. Report this limitation instead of searching for a hook.
- Grok Bot routines (the agent-turn fetch and any monitor or engagement routines) persist
  across restarts. Pause, resume or remove them with the routine tool, and keep the cache
  task record in step.
