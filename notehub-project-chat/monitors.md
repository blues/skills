# Monitors that notice, explain and close the loop

This is part of Project Chat. Start from the question already answered or requested,
using its platform definitions, product skills or both. A follow-up such as “watch this,”
“check every hour” or “notify me when the result changes” expresses a request for
recurring work; it does not need the word “monitor.” Once the definition is complete, follow
activation to install, exercise and verify the real
job and delivery path. Reuse existing cloud setup; don't restart authentication.

An explicit request to monitor something at a frequency authorizes that scoped recurring
work where the harness supports it. Reuse the request; ask only for missing details that
change the result. Propose a concrete definition before activation when a threshold,
population, cadence or recipient is ambiguous. Do not turn a one-off question into an
unlimited standing alarm. Use harness scripts; avoid Notehub-native operations whose
operationId contains `Monitors` or `Alerts` for this workflow.

A monitor is the recurring evaluation and its state; an alert is a delivered finding
from that monitor. They share one scope and evidence path. A fast cadence also needs
sufficiently fresh acquisition; the 15-minute default cache is not enough for every
monitor.

## Suggest a useful monitor

After answering, consider whether a periodic version would help the person act: a
repeated question, an important threshold, a fleet rollout, a recovery they are waiting
for, or a changing result they otherwise have to inspect manually. Offer it when relevant,
without attaching a generic monitoring pitch to every answer.

Make the suggestion concrete: which question will run, for which population and window,
how often, what change matters, where a notification would appear and why it is useful.
Use established product/platform rules; expose a missing threshold rather than inventing
one. Keep a suggestion distinct from an active job, and respect a decline or “not now.”
If the person already requested recurring notification, proceed within that request
instead of asking again whether they want monitoring.

## Create or change a monitor

1. Recover the question's exact calculation, subject and filters from the conversation
   and prepared plan. Preserve units, population, geography, validity and interpretation.
   Decide whether “the last two weeks” remains a rolling window at each run; state that
   choice and resolve ambiguity if it changes the intended result.
2. Define “changes” in terms of the answer: membership entering/leaving a result set,
   a state transition, a material numeric difference, a threshold crossing or a requested
   reminder. Ignore row order and incidental fetch timestamps. Use product-defined
   precision/significance; unknown, stale or failed retrieval is not a legitimate empty
   result or a recovery. Establish the first successful run as the comparison baseline
   and make clear whether an already-active condition calls for an initial notification.
3. Reuse cadence and destination from the request or standing preferences. Ask only for
   missing details that materially change the work. Verify the supported in-app or
   messaging-channel destination; keep user-local delivery times distinct from UTC
   analytical windows.
4. Write or reuse a deterministic script/query evaluator that accepts the monitor
   definition, fetches sufficiently fresh evidence, computes the result, compares durable
   state and produces notification candidates. Keep the selected MCP or CLI connection.
5. Follow activation to register the actual periodic job and delivery path. Exercise the
   evaluator, restart/state handling and delivery within scope. Record real job IDs,
   last/next run and delivery verification. A definition or script alone is not an active
   monitor. In a foreground-only environment, preserve a draft and explain the missing
   continuing capability without claiming anything will run.
6. Report the active condition, population, cadence, destination and how to change or
   pause it. For a changed monitor, preserve its identity where appropriate, recheck the
   comparison baseline and supersede obsolete queued alerts.

Each monitor stores: user wording and purpose, project/population, data/recipe and revision,
quantity/units, window and time basis, validity/unknown handling, thresholds or change
mechanism, cadence, severity, deduplication/rearm/closeout rules, recipient/channel,
authorization scope, quiet hours, and pause/delete controls. Mark inferred choices so
the person can correct them. Distinguish a user-local schedule from UTC analytical windows.

## Choose the relevant reference

Read the rows that apply to the current step; reuse context already established.

| If you need to | Read | Then |
|---|---|---|
| Ground the monitor calculation in platform/product meaning | [Prepared questions](../notehub-cloud-agent/questions.md) | Use the verified plan and its coverage/freshness requirements |
| Install or change an authorized monitor | [Activation](../notehub-cloud-agent/setup.md#activate-recurring-work) | Register, exercise and verify its real job and delivery path |
| Implement scheduled evaluation and delivery | [Job contracts](../notehub-cloud-agent/jobs.md) | Supply durable inputs, incident state, outbox and execution receipts |
| Pause, resume, remove or repair a monitor | [Lifecycle](../notehub-cloud-agent/lifecycle.md) | Verify the actual scheduler and delivery state |
| Add periodic status or contextual updates around monitors | [Engagement](engagement.md) | Respect the separate status cadence and user interests |

## Deterministic checks, thoughtful delivery

Run scripts for the routine sequence: acquire or query cache → validate coverage and
measurement quality → compute condition/change → compare persisted state → build candidate
and evidence → decide delivery → compose only if useful → record receipt and follow-up.
Do not wake the model for every unchanged 15-minute poll. Persist monitor state and a
pending delivery/outbox entry transactionally. Use stable incident/delivery IDs and an
idempotency key if the delivery provider supports it. A successful compose is not a
successful send: acknowledge only a delivery receipt. Retry failures without generating
a new incident; do not promise exactly-once messages without provider support.

For a result-change subscription, compare each valid result with the stored comparison
baseline and emit a distinct transition for each meaningful change, including changes
on consecutive polls. Commit the new baseline and its pending notification together.
Unknown evidence leaves the last valid baseline intact. Test unchanged results,
consecutive changes, reordered rows, unknown inputs and restart before activation.

Use three states for evidence: true, false and unknown. Unknown includes stale data,
incomplete coverage, bad measurements, failed acquisition and unsupported semantics.
Unknown must not close an incident or produce an “all clear.” Report impaired coverage
when it matters, using its own deduplication/rate limit. Pausing a monitor suppresses
its evaluation/delivery; resuming performs a fresh check instead of replaying stale news.

Open on a verified condition; suppress unchanged repeats; notify on a meaningful
worsening or explicitly requested reminder; close only on fresh evidence that the
condition cleared under the agreed rule. Use hysteresis, dwell or consecutive observations
when product semantics require them, not invented universal constants. A lower reading
can mean worse conditions (e.g. headroom or voltage). Single spikes may be real; repeated
values may be legitimate; movement can invalidate a fixed-site baseline. Preserve suspect
raw evidence with flags and exclude it from a derived alarm only under a justified rule.

Before delivering a queued candidate, recheck the latest relevant evidence and incident
state. A closeout queued before a new excursion must be cancelled or superseded, not
sent as recovery. Recompute every stated statistic using the same population, threshold,
window and units as the claim. Expire candidates past their relevance window; a failed
send is not a reason to forward stale news later. Record why a candidate was suppressed.

## Evaluation and delivery pseudocode

Implement this with the cloud agent's own runtime and scheduling tools. Names below
describe logical operations, not a required schema or API. Serialize evaluations of a
monitor and recheck its definition, pause state and comparison state before committing.

```text
load the monitor's current definition and saved state
if paused or not due: stop
evaluate its prepared question and record evidence time, coverage and meaning revision
if this is a periodic status report:
    create one candidate for the due interval, even if the result is unchanged
    report unknown/stale evidence honestly; it does not prove stability or recovery
else if evidence is unknown, stale or incomplete:
    preserve the last valid baseline and incident state
    record the coverage gap; apply its separate notice policy; stop
else if this is a condition monitor:
    apply the confirmed threshold, dwell/hysteresis and severity rules
    open on a verified active condition
    while open, suppress unchanged repeats; notify on agreed worsening or reminders
    close only on a fresh verified clear condition; rearm for a later occurrence
    apply the agreed initial-notification policy; no initial "all clear"
else if this is a result-change subscription:
    on the first valid result, establish the baseline and apply its initial-send policy
    otherwise compare with the last valid baseline using the agreed significance rule
    create a distinct candidate for each meaningful change, even on consecutive runs

in one database transaction:
    verify the definition and saved state have not changed, and the monitor is not paused
    save the new comparison/incident state and evaluation progress
    save any candidate as pending, with stable project/monitor/delivery identity and evidence
on conflict or failure: leave both unchanged and re-evaluate from current state
```

Advancing comparison state must not depend on a successful send: the pending delivery
survives alongside it. Conversely, an uncommitted evaluation must not send a notification.
Keep any product-specific reminder or escalation policy explicit.

```text
for each due pending delivery, without overlapping sends for that delivery:
    recheck pause, authorization, quiet hours, expiry and current relevant evidence/state
    cancel unauthorized/expired candidates; supersede obsolete ones
    defer during quiet hours or pause
    defer an alert whose premise cannot be verified; a status report can explain a gap
    compose the message from the verified evidence
    send through the configured destination; use the stable idempotency key if supported
    on confirmed receipt: record the receipt and mark delivered
    on known failure: retain pending state and retry with bounded backoff
    on uncertain timeout: reconcile provider status before retrying; keep the same identity
```

Verify the implementation with synthetic evaluations and a fake delivery destination
before exercising any authorized real delivery:

| Case | Expected result |
|---|---|
| Initial clear, open, unchanged repeat, worsening, clear, then open again | No initial all-clear; agreed transitions fire once and the monitor rearms |
| Reordered results, consecutive meaningful changes or unknown input | Row order alone is ignored; each real change is distinct; unknown preserves the baseline |
| Paused monitor or changed rule during evaluation | No stale commit or send; resume by checking fresh evidence |
| Crash between state calculation, commit and delivery | No partial state/candidate commit; pending work survives restart |
| Failed/uncertain send, or closeout followed by a new excursion | Reconcile using the same delivery identity; suppress obsolete closeouts |
| Due periodic status with unchanged or unavailable results | Produce the requested status once for that interval and report any evidence gap honestly |

## Alert the user and close the loop

An alert belongs to its monitor and its evaluated condition. Deliver it in the app or
through the configured messaging channel, using the authorized recipient and delivery
policy. Preserve a reference to the monitor/incident so a reply such as “pause this” or
“show me the underlying data” can be resolved without guessing.

A useful message says what changed, which population, compared with what, why it matters,
what the person can do, and how fresh/reliable the evidence is. Keep it proportional.
A closeout says what is now verified and any residual uncertainty; polling observes at
intervals, not continuously between checks. If the original alert was wrong, correct it
and its downstream rules, not just the message wording.

Delivery to someone other than the requesting user needs explicit authorization for the
recipient and scope, or an existing standing authorization that covers it. Do not invent
contacts or silently forward customer data. Keep recipients and private interests in
harness state, not shared product skills.
