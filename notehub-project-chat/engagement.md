# Follow the user's ongoing interests

Recognizing and following up on ongoing interests is part of Project Chat. On a Cloud
Computer, use the prepared engagement schedule and monitoring routine to bring back
useful findings, status or graphs without waiting for a separate “watch this” request.
Reuse existing monitors and the user's notification preferences.

An interest-driven monitor is the same kind of monitor as one requested directly. This
procedure decides what deserves continued attention and what makes an update useful.
The monitor procedure owns its calculation, checking cadence and notifications.

## 1. Recognize and retain the interest

Notice recurring questions, specific concerns, goals, unresolved issues and outcomes the
user cares about. Use the conversation and existing private interest records. Record
what matters and why, the Project and population, and whether the interest was expressed
or inferred. Keep only information useful for follow-up.

Use established preferences for timing, destination and what to send or exclude. Keep
interests user-editable and preserve corrections, declines and pauses. A useful continuing
interest warrants action; a passing mention need not become a standing monitor.

**Complete when:** the interest has a clear purpose and scope, or no continuing action
is warranted. In a foreground session, retain useful context without promising background
follow-up.

## 2. Create or update the monitor

Within an enrolled Project, reuse a matching monitor or create one for the useful
interest. Keep the inferred origin visible. Reuse an existing requested monitor's results
without changing its cadence or rules. Add new interest-driven monitors to the prepared
engagement task, which reads current interests on each run.

| What is needed | Procedure | Result |
|---|---|---|
| Turn the interest into a running check | [Create or change a monitor](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-chat/monitors.md#create-or-change-a-monitor) | A verified calculation and scheduled monitoring under the established preferences |
| Engagement scheduling is missing or has failed | [Engagement preparation](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/engagement-preparation.md) | Repair the missing setup, then return to this interest |

Act within the enrolled scope and established notification preferences. Resolve missing
consequential thresholds, a wider scope or a new recipient/channel before relying on
them. No separate monitoring request is needed for useful engagement within that scope.

**Complete when:** the interest has a verified running monitor, or a specific missing
meaning or capability is recorded. Return to the conversation; Steps 3–4 govern later
opportunities and feedback.

## 3. Decide what to bring back at each opportunity

Read current interests, fresh monitor results and recent notification/feedback history.
Use the monitor's rules and the following distinctions:

| Monitoring preference | Bring back | Apply |
|---|---|---|
| Notify on a change or condition | The verified result when its rule fires | Preserve the monitor's requested urgency and freshness |
| Send regular status, including when unchanged | Current result, trend, coverage and a concise explanation or graph | Deliver at the agreed interval even when the result is unchanged |
| Follow an interest without a fixed report | A relevant finding or occasional useful status update | Send when it helps the user's understanding or next action; a scheduled opportunity need not produce a message |

A regular status preference can be explicit (“send me a daily graph even if it is stable”)
or established in continuing-assistance settings. Without that preference, unchanged
status can still be useful during a watched rollout or after a problem. Missing or stale
evidence must not appear as stable performance.

### Choose relevant findings

Consider magnitude, persistence, baseline and product significance together. Example
thresholds such as 24 hours or 40 percent are not universal rules. Bundle related
findings, suppress repeated or weak speculation, and preserve uncertainty. Connect to a
live concern: Device recovery, a rollout reaching its target, a change in a population
the user follows, or new evidence that settles an earlier question.

Outside context such as weather, an advisory or news is useful when it relates to that
concern. Retrieve current authoritative evidence, date and cite it, and distinguish an
inferred connection from a demonstrated cause. Leave unrelated news out.

### Present the useful update

Lead with the development or current status, its significance and a useful next action
or focused question. Choose a short message, a graph or both for the destination. A graph
shows the relevant window, units, confirmed baseline or threshold, and evidence gaps.
Explain what to notice. Avoid exaggerated urgency, guilt, false familiarity and claims
of continuous watching. Silence is a valid result when no useful update is due.

| To send the update | Procedure | Result |
|---|---|---|
| Apply the monitor's timing, destination and notification progress | [Monitor notifications](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-chat/monitors.md#notify-the-user-and-handle-replies) | An evidence-based update under the user's preferences, with its actual send outcome recorded |

**Complete when:** a due or useful update is delivered through its monitor, or the reason
for silence, deferral or failed delivery is retained. Routine checks run in code; use an
LLM turn when interpretation, composition or user interaction adds value.

## 4. Learn from replies and keep interests current

Record useful feedback so later checks reflect what the user values. A reply asking for
explanation returns to the Chat question loop with the monitor's calculation and evidence.
Adjust the same monitor when the user changes timing, scope, format or destination.

Close useful loops with verified recovery or completion. Review whether an interest still
matters and retire it when resolved, declined or expired; do not keep sending reminders
about a finished issue or recreate a removed interest.

| For a control request | Procedure | Result |
|---|---|---|
| Pause, resume or remove an interest and its monitoring | [Manage existing work](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/lifecycle.md#pause-resume-or-remove-work) | Verify the change while preserving unrelated monitors and cache maintenance |

Return to Project Chat with the updated context. Personal interests, preferences and
notification history remain private; product Training owns changes to shared knowledge.
