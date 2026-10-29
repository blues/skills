---
name: notehub-project-chat
description: Continue work with a selected Notehub Project through an ongoing loop of questions, investigations, operations and monitors that keep the user informed. Reuses established access and internalized product knowledge, answers from the Notehub Cache when it covers the question, and carries follow-ups into verified actions or recurring work. Shared product-knowledge training is a separate activity.
---

# Chat about the Project and keep the user informed

Arrive here after the entry cookbook has established access, taught the API, prepared
any Cloud Computer, selected ProjectUID and loaded the complete product skill set.
Reuse that connection, product understanding and known cache readiness. Follow
[working style](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub/working-style.md).

A **monitor watches something and alerts the user according to their monitoring
preferences**: periodically, when a result changes, when a threshold is crossed, or under
another agreed rule. Its notifications can be messages, graphs or both, delivered in the
chat stream, WhatsApp, SMS, Slack or another supported channel. These preferences belong
to the monitor, whether requested directly or inferred from an ongoing interest.

**For ordinary data questions, use the Notehub Cache when it has the required fields,
population, history and freshness.** If the question concerns the whole dataset or no
cache is available, use the **Project Repository when access is through MCP and the
Project has a Repository enabled**. Repository queries are slower than the local cache
but faster than direct Notehub API queries. The question and Repository procedures below
apply Notehub and Product Skills to the chosen source; direct API reads supply missing
evidence or current service state.

| Only if preparation is missing | Read | Then |
|---|---|---|
| This skill was entered without the required context, or a prerequisite has become unavailable | [Entry cookbook](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub/SKILL.md) | Complete or repair the first missing step, preserve working context, then resume Chat |

## Repeat this conversation loop

Follow Steps 1–4 for each user turn or monitor finding brought into the conversation.
Follow one matching procedure in Step 2, reading its supporting references only where
needed. A supporting procedure returns to its caller; after the work is complete, return
here to explain the result and carry it forward. Apply the **along-the-way** actions as their
triggers arise. Completing one answer does not end Project Chat.

### 1. Understand the next request

Use the conversation to resolve what the person means, including “this,” “those Devices,”
or “notify me when that changes.” Carry forward the question's population, filters, units
and time window. Ask only about ambiguity that changes the work. A notification request
can reuse the calculation just answered; it does not require repeating setup or inventing
a different question.

Handle a change of activity immediately:

| The user says | Continue with | Result |
|---|---|---|
| Done or stop, referring to the conversation | [Entry cookbook — Step 4](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub/SKILL.md#4-choose-the-project) | End this loop, preserve a brief resumable context and wait for the user's next request |
| Switch Projects | [Entry cookbook — Step 4](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub/SKILL.md#4-choose-the-project) | Select the requested Project, read its product skills and enter the requested activity |
| Teach or revise shared product knowledge | [Project Training](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-train/SKILL.md) | Carry the current knowledge and relevant evidence into the dedicated interview and publication workflow |

A request to stop a particular monitor applies to that monitor in Step 2. Ending Chat
or switching Projects keeps existing caches and monitors attached to their ProjectUIDs;
it does not cancel or retarget background work.

### 2. Complete the requested work

Choose the matching procedure and follow its steps. For questions, the order is:
understand the meaning → choose the cache, Repository or API → acquire any missing evidence →
calculate and verify. A follow-up reuses completed work where it still applies.

| The person wants to | Procedure | Complete by |
|---|---|---|
| Ask a question or investigate a finding | [Answer a question](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-chat/querying.md) | Calculating a grounded result from suitable evidence, using the cache first when it covers the question |
| Get ready for fast Q&A or improve repeated questions | [Optimize Q&A](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-chat/optimizing-qa.md) | Preparing or reusing verified queries against the existing data sources |
| Query about 100% of historical data (beyond any Notehub Cache), or answer when no cache is available | [Query the Project Repository](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-repository/SKILL.md) | Introspecting and querying the Repository through MCP using Notehub Skills and Product Skills; requires MCP access and a Repository enabled for this Project |
| Watch something, receive periodic updates or graphs, or be alerted when a result changes or a threshold is crossed | [Create or change a monitor](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-chat/monitors.md#create-or-change-a-monitor) | Reusing the question's calculation, setting when and where to notify, and verifying scheduled checks and notification delivery |
| Explore or investigate state of Devices, Fleets, or other Project entities. Make changes if authorized. | [Operations](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-chat/operations.md) | Explaining the observed state and, when requested and authorized, making changes and verifying their actual effect |
| Explore and manage Firmware Updates for Devices and Fleets | [Explore and manage Firmware Updates](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-firmware/SKILL.md) | Explaining reported and desired firmware, investigating update progress, and carrying authorized updates through verification |
| Pause, resume or remove a monitor or interest | [Manage existing work](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/lifecycle.md#pause-resume-or-remove-work) | Verifying the requested change in saved state, actual execution and delivery |

Questions may need platform meanings, the product knowledge already internalized, or
both. Use the product's terminology and its mappings from NotefileID and Event `body`
fields to database columns. Consult Platform for native metadata or system Notefiles;
the question procedure places these resources at the point they are needed.

A question about a monitor's notification uses the same question procedure and its saved
calculation and evidence. A request to change its timing, threshold or channel changes
that monitor through the create/change procedure.

Recurring execution requires verified Cloud Computer capabilities. Foreground agents can
answer questions and perform supported operations through the established connection;
distinguish a prepared monitor from one that is registered and running.

### 3. Explain the verified result

Give the answer or outcome first, its practical meaning, and any useful next action.
Choose a concise explanation, table, graph or other output that fits the question. State
the population, time window, evidence freshness and limits that affect the decision.
Distinguish unknown from zero and incomplete coverage from an empty result.

For an operation, distinguish accepted, queued and effective state; an accepted request
is not proof that a Device applied it. For recurring work, say what will run, how often,
what triggers a notification and where it will arrive. Name any remaining limitation;
a saved definition alone is not an active monitor.

For a readiness request, report actual readiness and offer a few questions grounded in
the loaded platform or product knowledge. An initializing cache can still support useful
answers within its verified coverage.

### 4. Carry the context into the next turn

Keep the question's calculation, filters, units, time semantics, evidence coverage and
meaning revisions available. Retain relevant operation results and monitor IDs so a
follow-up such as “show the data” or “pause this” resolves to the right work. Reuse query
plans while refreshing the evidence; do not reuse an old numeric answer as current data.

Carry forward what the conversation teaches you about the user's ongoing interests.
Apply the engagement action below when an interest warrants continued attention.

**This turn is complete when** the requested answer or action has a verified result, or
its specific missing information or capability is explained. Wait for the next user turn
or a finding from registered background work, then repeat from Step 1. Do not manufacture
more questions to keep the conversation going.

## Along the way in the loop

**Recognizing and following up on the user's ongoing interests is a core responsibility
of Project Chat.** Throughout the conversation, actively notice recurring questions,
specific concerns, goals and outcomes the user cares about. On a Cloud Computer, turn
useful interests into monitors that bring back relevant findings, status or graphs under
the established notification preferences. **Do not wait for an explicit “watch this”
request.** Keep interests current as the conversation develops, and respect corrections,
declines and pauses. The engagement procedure below makes this part of the ongoing loop.

Apply these actions as their triggers arise:

| When this happens | Do this | Supporting procedure |
|---|---|---|
| **The conversation reveals an ongoing interest** through repeated questions, a specific concern or a goal | **Act on the interest.** On a Cloud Computer, record what matters and why, reuse or create a monitor within the established scope and notification preferences, and verify it is running. Use it to bring back useful updates without waiting for another request | [Interest-driven engagement](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-chat/engagement.md) |
| A query uses fields or exposes a slow or repeated calculation | Record the fields, query plan and observed performance in private workload history; reuse a verified plan where meaning still matches | [Optimize Q&A](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-chat/optimizing-qa.md) |
| Live reads are needed to fill an evidence gap or follow changes | Use modified cursors where supported, with the exact checkpoint and batch contract; keep the cache's coverage honest | [Poll for changes in Events, Sessions, Devices, and more](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-chat/modified-cursors.md) |
| An answer reveals a useful recurring question | Suggest the concrete question, benefit, cadence and notification behavior; respect a decline. An explicit monitoring request goes straight to Step 2 | [Suggest a useful monitor](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-chat/monitors.md#suggest-a-useful-monitor) |
| A monitor's notification rule or periodic schedule is due | Send a relevant explanation, graph or both to its configured destination; follow its preferences, including unchanged status when requested | [Monitor notifications and replies](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-chat/monitors.md#notify-the-user-and-handle-replies) |
| Product knowledge changes or evidence reveals an incorrect interpretation | Reconcile affected calculations, monitors and prior conclusions; shared skill changes go through Training when requested | [Reconcile changed meaning](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/lifecycle.md#reconcile-a-changed-product-rule-or-correction) |

Cache maintenance continues on its existing schedule. Preserve user-requested pauses and
use available evidence while missing history fills. Personal interests, query plans,
monitor state and delivery history remain private harness state, separate from shared
product skills. Useful engagement can include unchanged status or relevant outside
context; an isolated question is not an unlimited standing alarm.
