---
name: notehub-cloud-agent
description: Prepare a Cloud Computer to maintain a Notehub data cache for fast Q&A, run monitors that notify users according to their preferences, and create engagement monitors from observed user interests. Configure shared storage, authenticated background execution, maintenance, evaluation and delivery.
---

# Prepare the Cloud Computer

The Cloud Computer serves three purposes:

- **Answer questions from a maintained cache.** Keep key Notehub data locally so most
  questions can query the cache directly, with known freshness and coverage.
- **Run monitors the user requests.** Repeatedly evaluate their questions or conditions
  and alert them periodically, on changes or thresholds, or under another agreed rule,
  through their chosen chat or messaging channel.
- **Run monitors inferred from the user's interests.** Learn what matters from the
  conversation, create useful engagement monitors, and bring relevant findings, status
  or graphs back to the user without waiting for a separate monitoring request.

Use this cookbook only when the agent is running on a Cloud Computer. Access and API
learning come first. Reuse existing setup, preserve pauses and complete only missing
steps. Prepare the shared facilities now; Project Selection will enroll projects, and
Project Chat will supply their questions, requested monitors and inferred interests.
Follow [working style](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub/working-style.md).

## 1. Establish how background work runs

Before creating the cache, verify script locations, available storage and low-space
warnings, authenticated background access, scheduling and restart recovery, and how
workers wake an LLM task. Prefer Python doing routine work without an LLM turn on every
poll. Reuse verified setup and capture the actual paths, limits and hooks.

| To complete this step | Read | Result |
|---|---|---|
| Establish the execution and recovery paths | [Background execution](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/harness.md) | Durable worker locations, authenticated API access, scheduling, restart recovery and agent wake-up |

## 2. Create or reuse the Notehub Cache Database

Use **DuckDB**, installed in the worker environment. Check for an existing shared cache
first. Where bots share storage, keep one database in their standard shared workspace,
independent of any bot. Verify persistence, prepare a minimal project registry, and make
the cache discoverable through a nearby README and shared memory fact.

| To complete this step | Read | Result |
|---|---|---|
| Locate, install and organize the database | [Cache preparation](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/cache-preparation.md) | Shared storage, project-scoped tables, safe concurrent access and verified persistence |

## 3. Start the Notehub Cache Maintenance Task

Register or reuse a task with a refresh goal of **every 15 minutes**. If the Cloud
Computer's background scheduler cannot meet that goal, refresh as often as it permits:
the user may want to be notified of changes. Refresh Devices, Events, full per-device
Sessions, project metadata and Notefile Schemas. Use modified cursors where supported.
Run daily purge, schema refresh and optimization using queried-field
history. Each run also checks available storage and quotas, waking the agent when
space needs attention. It reads the live project registry; new selections join
automatically. With no projects yet, storage checks still run. Apply Step 1's restart
mechanism so periodic work resumes after a restart. Respect existing schedules and pauses.

| To complete this step | Read | Result |
|---|---|---|
| Register and verify recurring maintenance | [Maintenance task setup](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/cache-maintenance.md) | A verified worker, schedule, storage monitoring and restart recovery |
| Poll for changes in Events, Sessions, Devices, and more | [Modified cursors](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-chat/modified-cursors.md) | Batches committed with exact resumable checkpoints |

## 4. Create or verify the monitor facilities

Prepare a routine that checks what a monitor watches and notifies the user according to
its preferences. When Chat supplies a monitor, bind its calculation, checking cadence,
notification rules and destination to that routine. Query the cache when it meets the
monitor's freshness requirements. With no monitors yet, leave the facility ready.

| To complete this step | Read | Result |
|---|---|---|
| Prepare background checking and notification | [Monitor preparation](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/monitor-preparation.md) | A verified routine ready for requested and interest-driven monitors |

## 5. Enable engagement from observed interests

Register or reuse engagement checks, defaulting to three daily local-time opportunities.
As Chat reveals interests, create bounded monitors using the same evaluation/delivery facilities,
without a separate “watch this” request. Record their inferred origin and honor delivery
preferences, quiet hours and pauses. No interests means an idle check; no useful finding
means no message. Requested monitors keep their own cadence.

| To complete this step | Read | Result |
|---|---|---|
| Prepare recurring engagement opportunities | [Engagement preparation](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/engagement-preparation.md) | A verified schedule ready for interests learned in Project Chat |

## 6. Verify readiness and return to the project workflow

Verify database persistence, task execution, schedule recovery and the worker-to-agent
handoff, retaining job IDs, run receipts and next runs. Report cache, monitoring and
engagement readiness separately, including
anything awaiting a project or missing capability. Return to project selection for
enrollment and project cache checks, then product reading and the requested activity.
Useful Q&A and Training can proceed while history fills.

| When continuing or repairing setup | Read | Result |
|---|---|---|
| Continue with the selected project | [Project Selection](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-access/projects.md) | Verified scope, enrollment and cache readiness |
| Repair, pause, resume or remove work | [Lifecycle](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/lifecycle.md) | Verified job state with user choices preserved |
