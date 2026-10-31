# Prepare or verify engagement scheduling

This is Cloud Computer cookbook Step 5. Prepare recurring opportunities for the monitors
that Project Chat creates from the user's interests. Reuse the monitoring routine from
Step 4, the shared database and existing registrations. Project Chat supplies the interests;
setup does not need an invented Project or question. This schedule does nothing until
Project Chat registers interests via [engagement](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-chat/engagement.md); load that procedure
when entering Project Chat.

## 1. Recover the current schedule and preferences

Find existing interests, monitor definitions, notification preferences and task records.
Preserve pauses, declined interests and established schedules. With no interests yet,
keep an empty set and allow the task to finish quietly.

**Complete when:** the task can locate the current interests and monitoring routine.

## 2. Register or verify recurring opportunities

Reuse the existing schedule. Otherwise use three daily local-time opportunities by
default: morning, early afternoon and evening. Apply the user's preferences and timezone;
resolve an unknown timezone before scheduling and use the scheduler's daylight-saving
support. These opportunities do not change the cadence or urgency of explicitly
requested monitors.

The task reads current interests and monitors on every run, including those added later.
Routine checks run in code. Wake the agent only for findings that need interpretation,
composition or user interaction. Reuse the established restart mechanism.

| To implement the task | Procedure | Result |
|---|---|---|
| Run due checks and use current interests | [Engagement worker](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/jobs.md#engagement) | A scheduled invocation that reads live state and stays quiet when no useful update is due |

**Complete when:** the actual scheduler reports the task and its next run.

## 3. Verify unattended execution

Run through the actual background path. Use disposable cases to verify that newly added
interests are found without recreating the schedule, pauses remain intact, empty or quiet
results cause no message, and pending work survives restart. Reuse notification-path
verification from monitor preparation; a fake destination can check these cases without
sending an unsolicited message.

**Complete when:** the task runs correctly through the unattended path and preserves
user choices and notification progress.

## 4. Record readiness and return

Retain the invocation, task ID, last/next run and any limitation. With no interests yet,
report scheduling ready. Return to the Cloud Computer cookbook.

| When Chat supplies an interest or an opportunity runs | Procedure | Result |
|---|---|---|
| Recognize interests and choose useful updates | [Follow ongoing interests](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-chat/engagement.md) | Reuse or create grounded monitors and adapt follow-up to the user |
