# Be useful between questions

Engagement continues Project Chat between questions, using the user's monitors, question
history and interests. Configure this through cloud Setup; scheduled checks follow
the engagement job contract. Reuse prepared questions and the
private interest profile, rather than inventing a new investigation at every wake-up.

## Choose the relevant reference

Read the rows that apply to the current step; reuse context already established.

| If you need to | Read | Then |
|---|---|---|
| Configure authorized engagement opportunities and delivery | [Activation](../notehub-cloud-agent/setup.md#activate-recurring-work) | Register the real schedule and verify the delivery path |
| Run a scheduled engagement opportunity | [Engagement job contract](../notehub-cloud-agent/jobs.md#engagement) | Evaluate prepared questions, record the outcome and send only useful findings |
| Define or change the monitor behind an update | [Monitors](monitors.md) | Keep the question, evaluator, alert rule and status preferences connected |

The purpose is to help the person notice relevant changes and make better decisions,
not to fill a notification quota. Maintain a small, user-editable interest record:
what they care about, why, affected projects/populations, explicit versus inferred origin,
what to tell them, what never to send, when/how to deliver and when an interest expires.
Repeated or specific questions can reveal an interest; one casual query does not authorize
an intrusive standing monitor. Use inferred interests for proportionate, low-cost checks
and helpful suggestions. Confirm consequential thresholds, expanded scope or external sends.
Respect “not now,” declines, quiet hours and pauses without arguing or repeatedly proposing
the same thing. Give the person a simple way to inspect, correct or remove interests.

Unless they specify a schedule, register three daily **evaluation opportunities** in
their timezone: morning, early afternoon and evening, at reasonable local hours. Resolve
an unknown timezone before scheduling local-time delivery; handle daylight saving through
the scheduler's timezone support. These are not three mandatory messages. Query scripts
first; call the model only for a meaningful candidate that needs interpretation or writing.
Explicit monitors keep the requested cadence and urgency instead of being delayed to
these engagement windows.

## Choose the useful kind of update

| User need | Deliver | Keep distinct |
|---|---|---|
| Know when a monitored condition changes | An alert when its condition, escalation or recovery rule fires | The monitor's requested cadence and freshness govern it |
| Receive regular status, including when unchanged | A fresh periodic report at the agreed cadence: current result, trend, coverage and a concise explanation | “No material change” can be useful; a status report is not a new incident |
| Follow an interest without requesting a fixed report | A relevant finding or occasional useful status check at an engagement opportunity | Evaluate usefulness and recent delivery history; do not send a message merely to fill a slot |

A regular status preference can be explicit (“send me a daily graph even if it is stable”)
or already established in the user's continuing-assistance settings. Respect that cadence:
do not suppress a due report just because the monitor's alert state is unchanged. Without
such a preference, send unchanged status only when it adds value, such as reassuring the
user during a watched rollout or showing sustained stability after an incident.

Use a short message, a graph or both according to the question and delivery surface.
A graph should show the relevant time window, units, baseline/threshold when established,
and gaps or stale evidence. Explain what the person should notice. A stable graph with
fresh complete coverage and a flat line caused by missing data are different reports.

For each opportunity, read the interest state and recent delivery/feedback history, inspect
fresh results, and check whether a requested status report is due or an update usefully
informs the person's view or next action. Use baseline, magnitude, persistence and product significance together;
24 hours, 40 percent, two times normal and similar example thresholds are not universal
rules. Rank and bundle related findings; suppress repeats and weak speculation. Carry
uncertainty forward rather than laundering an estimate into certainty with repetition.

Connect to a live thread: a device recovery after a discussed outage, a rollout reaching
its agreed target, a change in a repeatedly asked-about population, or newly available
evidence that settles a question. Outside context such as weather or an advisory is
useful only when relevant to that thread; retrieve current authoritative evidence, date
and cite it, and explain the connection as an inference where appropriate. Do not insert
unrelated news or a speculative cause just to sound proactive.

Write as an attentive colleague: the development or current status first, its significance, and a useful
next step or focused question. Avoid dashboards in prose, exaggerated urgency, guilt,
false familiarity and claims of continuous watching. Silence is a valid result. Record
what was delivered and the response so future checks can adapt. Track periodic reports by
their reporting interval so alert deduplication does not suppress a requested status
report, and retries do not send the same report twice. Close useful loops with
verified recovery or completion; don't keep reminding someone about a resolved issue.
