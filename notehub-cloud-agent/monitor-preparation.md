# Prepare or verify background monitoring

This is Cloud Computer cookbook Step 4. Prepare a reusable routine that watches a
question and alerts the user according to its monitoring preferences. Project Chat
supplies the Project, calculation and preferences later. Reuse the database, runtime,
scheduler and wake-up path already prepared, preserving existing monitors and pauses.

## 1. Locate or prepare the monitoring routine

Reuse an existing working routine. Otherwise choose its durable code location and how
it reads monitor definitions from the database. Retain only the progress needed to resume:
pause state, last valid comparison, last check and notification progress. Preserve this
state when rebuilding cached Notehub data. One routine can serve many monitors.

**Complete when:** background execution can load a monitor and its saved progress.

## 2. Implement or verify checking and notification rules

Run the prepared query in code, verify coverage and freshness, and apply the monitor's
preferences: periodic notification, meaningful change, threshold or another agreed rule.
Support combinations of these preferences. Preserve the last valid comparison when data
is unknown; a due status update can explain the gap. Keep pending notifications with
saved progress so a restart does not lose a result or repeat a completed notification.

| For implementation | Resource | Then |
|---|---|---|
| Calculate suitable evidence, apply preferences and resume reliably | [Monitoring routine](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/jobs.md#evaluation) | Implement the shared pseudocode and verify its relevant cases |

**Complete when:** a supplied monitor can calculate its result and identify due notifications.

## 3. Connect the notification path

Use the harness's supported chat stream or messaging channel, such as WhatsApp, SMS or
Slack, according to the user's preferences. Verify which destinations are available.
Use the same approach as cache acquisition. Prefer no model involvement: evaluate
monitors in Python against the cache, save due findings to a pending store and deliver
prepared messages from code wherever the harness offers a script-callable delivery
path. Fall back to scheduled agent turns only where code cannot reach the channel or a
finding needs interpretation, composition or user interaction. Set this up without
asking the user about model cost. Record actual send outcomes.

| For implementation | Resource | Then |
|---|---|---|
| Wake an LLM when needed | [Worker-to-agent handoff](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/jobs.md#wake-the-agent-only-for-work-that-needs-it) | Retain work if a wake-up fails |
| Send through the configured channel | [Notification delivery](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/jobs.md#delivery) | Respect preferences, record outcomes and reconcile retries |

These can be parts of the same routine; use separate workers only when the harness needs
them. A channel awaiting user configuration remains an explicit limitation.

**Complete when:** a due notification has a usable send path or a known missing capability.

## 4. Verify through the background runtime

Use disposable examples and a fake destination to check changed and unchanged results,
periodic updates, unknown evidence, pause, restart and failed/uncertain sends. Invoke the
actual background path. Send a real test only within the setup request or established
preferences; otherwise retain the destination's unverified status.

**Complete when:** the routine passes the relevant monitor checks through unattended execution.

## 5. Record readiness and return

Save the invocation and verification result and apply the established restart mechanism.
Project Chat supplies each monitor's question, cadence and preferences; with none yet,
report the routine ready for monitors. Preserve existing registrations and pauses.
Return to the cookbook for engagement setup, which uses this same monitoring routine.

## Grok Bot only

*This section applies only to agents running as Grok Bot. Other harnesses should skip it.*

- Box scripts cannot start an agent turn or post a message. Evaluate monitors either in
  Python from the scheduler loop, writing findings to the cache's pending-findings folder,
  or inside a scheduled Grok Bot routine. In both cases a routine reads the pending
  findings and delivers them to the user's chosen destination (this chat, Slack and so
  on). Stay silent when there is nothing due.
- Do not use internal agent sockets or local service ports as a wake-up path.
