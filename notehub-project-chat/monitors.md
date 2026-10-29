# Create, change and run a monitor

A **monitor watches something and alerts the user according to their monitoring
preferences**. It can notify periodically, when a result changes, when a threshold is
crossed, or under another agreed rule. Notifications can appear in the chat stream or
through WhatsApp, SMS, Slack or another channel supported by the harness. Timing,
message format and destination are all part of the same monitor.

This is part of Project Chat. Start from the question already answered or requested,
using its platform definitions, product skills or both. “Watch this,” “check every hour,”
“notify me when the result changes” and “send me a daily graph” are monitoring requests.
Reuse the existing Cloud Computer setup and selected authenticated connection.

Engagement creates the same kind of monitor from an observed interest. Record that
origin and use the established engagement timing and notification preferences.
Engagement owns how interests are inferred, reviewed and retired.

## Suggest a useful monitor

After answering, consider whether watching the question would help the person act: a
repeated question, an important threshold, a rollout, a recovery they are waiting for,
or a changing result they otherwise inspect manually. Make a relevant suggestion
concrete: what to watch, how often to check, when to notify, where and why it is useful.
Respect a decline. If recurring notification was already requested, proceed within
that request instead of asking again whether the user wants monitoring.

## Create or change a monitor

Follow these steps for a requested monitor or one inferred through engagement. Reuse
an existing matching monitor and change only what is needed. A request authorizes the
scoped recurring work where the harness supports it; ask only about missing decisions
that affect the result or delivery.

### 1. Recover the question and its calculation

Preserve the existing monitor's identity, ProjectUID and pause state. Record the request
or inferred interest and purpose. Recover the question's calculation, population, filters,
units, geography and validity rules. Decide whether a window such as “the last two weeks”
rolls forward with each check.

| If the calculation is not ready | Resource | Then |
|---|---|---|
| Ground and verify the question | [Prepared questions](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-chat/prepared-questions.md) | Retain the calculation, meaning revision and evidence requirements |

**Complete when:** what the monitor watches and how it calculates the answer are clear.

### 2. Set the monitoring preferences

Reuse the user's request and established preferences. A monitor can combine rules, such
as “notify me when the limit is crossed, and send a daily graph even if nothing changes.”

| Preference | Define or verify |
|---|---|
| When to check | Evaluation cadence and required data freshness |
| When to alert | Periodic schedule, meaningful result change, threshold crossing or another agreed condition; include reminders or recovery notifications when wanted |
| First check | Whether to send the initial result or only establish a comparison baseline |
| Where to alert | Chat stream, WhatsApp, SMS, Slack or another available channel, with the intended recipient |
| What to send | Concise result, explanation, graph or a useful combination |
| Timing and controls | Timezone, quiet hours, expiry if useful, and how to change, pause or remove the monitor |

Checking frequency and notification frequency can differ. A monitor might check every
15 minutes and send one daily update. A requested periodic update is due even when the
result is unchanged. Use the user's local timezone for delivery and the question's
established time basis for calculations.

For change or threshold rules, define what counts as significant using product meaning:
ignore row order and fetch timestamps, retain units and precision, and apply persistence,
hysteresis or escalation only when the product or request requires it. Decide when a
repeated condition warrants another notification. Missing, invalid or stale evidence is
unknown, not zero or recovery. A missing necessary threshold leaves this monitor a draft;
it does not prevent unrelated work.

Respect product-specific direction and validity: a lower voltage may be worse, repeated
values may be legitimate, and movement may invalidate a fixed-site baseline. Preserve
suspect evidence with flags and exclude it from a calculation only under a justified rule.

**Complete when:** the monitor's checking and notification behavior is unambiguous.

### 3. Implement or reuse the scheduled check

Bind the prepared query and preferences to the harness's existing monitoring routine.
Use code for ordinary queries and comparisons; wake the agent when interpretation or
user interaction is needed. Query the Notehub Cache when its coverage and freshness
suffice, otherwise acquire the necessary bounded evidence through the selected connection.

Keep only the state needed to resume: the definition, pause state, last valid comparison,
last check, and notification progress. Retain pending notifications when retries are
needed. The agent chooses the organization; a shared routine can run many monitors.
When changing a rule, recheck its baseline and discard pending messages made obsolete
by that change. Use harness jobs for these monitors; Notehub-native API operations
containing `Monitors` or `Alerts` are not this workflow.

| If execution is missing or unverified | Resource | Then |
|---|---|---|
| Prepare background checking and notification | [Monitor preparation](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/monitor-preparation.md) | Repair the missing capability, then return to this monitor |

**Complete when:** the monitor can evaluate its question and resume from saved progress.

### 4. Verify the calculation and notification path

Reuse the background routine and delivery-path checks completed during Cloud Computer
preparation. Verify this monitor's calculation and changed preferences using the relevant
cases below; use synthetic inputs and a fake destination where needed. Preserve real
data and production settings.

| Case | Verify |
|---|---|
| First check, unchanged results and consecutive meaningful changes | Follow the initial-send preference; notify for each relevant change without repeated messages for the same check |
| A threshold crossing, repeated condition and later recovery | Apply the agreed significance, repeat and recovery rules |
| A periodic update is due, with unchanged or unavailable data | Send the requested update once for that interval; explain any evidence gap |
| Stale evidence, pause or a changed definition | Preserve the last valid comparison; suppress claims or messages that no longer apply |
| Restart or failed/uncertain send | Retain needed progress, reconcile retries and avoid duplicate notifications |

Reuse a verified destination. For a new or changed destination, verify it through the
harness. Send a real test only when the request or established preferences cover it;
otherwise retain an explicit unverified-delivery status. A composed message is not a send.

**Complete when:** calculation and notification behavior pass the relevant checks.

### 5. Register or verify recurring execution

Register at the agreed checking cadence, or add the monitor to an existing shared task.
An interest-driven monitor can use the engagement task. Compare required freshness with
the cache's actual refresh interval: a faster check needs fresher inputs, not repeated
queries of unchanged data. Acquire bounded live inputs or adjust acquisition within scope
when necessary. Report scheduler or source limits that prevent the requested behavior.

Run this definition through the established unattended path. Verify its registration,
last/next run and saved progress, preserve pauses, and reuse restart recovery. In a
foreground-only environment, retain the draft and explain the missing execution capability.

**Complete when:** this monitor's scheduled execution is verified, or its limitation is recorded.

### 6. Report the result and return

State what is watched, how often it is checked, when and where the user will be notified,
and how to change or pause it. Distinguish active monitors from drafts or blocked work.
Keep its identity available so “pause this” or “send it to Slack instead” resolves to the
same monitor. Return the monitor and its verified status to Chat or the engagement
procedure that called this one. Personal preferences and notification history stay private.

| For subsequent work | Resource | Then |
|---|---|---|
| Pause, resume, remove or repair a monitor | [Lifecycle](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/lifecycle.md) | Verify the requested change in execution and notification behavior |
| Infer and maintain monitors from user interests | [Engagement](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-chat/engagement.md) | Apply the same monitor procedure within the established interest and delivery scope |

## Notify the user and handle replies

Use the monitor's configured timing, format and destination. A useful notification says
what changed or what the current status is, which population and time window it covers,
why it matters, and how fresh the evidence is. Include a graph when useful. Keep the
message proportional, and preserve its connection to the monitor and calculation.

A reply asking for explanation or underlying data returns to the normal Chat question
procedure. A request to change timing, threshold or destination changes this same monitor.
When recovery notifications are wanted, report only what fresh evidence verifies; periodic
checks do not prove what happened between checks. Correct a mistaken notification and any
dependent calculation or rule.

Deliver within the user's requested or established recipient/channel scope. Do not invent
contacts or forward customer data to someone else without authorization. Keep recipients
and private interests in harness state, not shared product skills.
