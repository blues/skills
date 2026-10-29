---
name: notehub-project-chat
description: Work with a Notehub project through questions, investigations, operations, periodic monitors, alerts and ongoing engagement. Uses platform knowledge, stored product/project skills or both; optimizes Q&A and turns requested follow-ups into verified recurring work. Shared-knowledge training is a separate dedicated activity.
---

# Chat about the project, act and keep the user informed

Project Chat carries the conversation from a question to its answer, a requested action,
a monitor of that answer, and useful follow-up. Access establishes the connection and
execution environment first; the project is selected before working with its data.
Reuse that context. Follow [working style](../notehub/working-style.md).

## Complete only missing prerequisites

| If you need to | Read | Then |
|---|---|---|
| Install/connect a tool, authenticate or prepare a Cloud Computer | [Access](../notehub-access/SKILL.md) | Return with one secure connection and the environment's actual capabilities |
| Choose or disambiguate the project | [Project selection](../notehub-access/projects.md) | Retain the verified data project and skill-owning parent |
| Understand the selected product if its knowledge is not already loaded | [Reading product knowledge](product-knowledge.md) | Discover its skills and read the index and product overview before choosing task-specific material |
| Reuse an established connection, project and environment | No repeated setup | Continue with the user's question or follow-up |

## Choose the knowledge the question needs

Platform knowledge explains native Notecard/Notehub data. Product skills (also called
project skills) explain the customer's product, its firmware-generated Event bodies,
vocabulary and rules. Load the relevant knowledge for the question; use both when needed.

| The question concerns | Read | Use it to establish |
|---|---|---|
| Devices, location, time, coverage, battery, firmware, sessions or other native data | [Platform](../notehub-platform/SKILL.md) | Which native records and fields answer the question and what they mean |
| Customer measurements, custom Event bodies, product context or business meaning | [Reading product knowledge](product-knowledge.md) | The stored product index, relevant recipes, field definitions and all consequential dependencies |
| A product measurement constrained by native location, time, identity, health or other metadata | [Platform](../notehub-platform/SKILL.md) and [Reading product knowledge](product-knowledge.md) | Combine both meanings in one query over the same records and scope |

The producer determines which knowledge applies: native system Notefiles also have bodies,
and a familiar field name does not establish a customer's measurement semantics.

For example, “radiation data in the last two weeks within 100 km of Kyiv” needs product
knowledge for the radiation field, units, validity and interpretation, plus platform
knowledge for capture time, location provenance and geographic filtering. Establish the
requested spatial reference, select the two-week window and associate the appropriate
location with each measurement. Do not substitute today's Device location for historical
Event location or assume every radiation field is a comparable dose-rate measurement.

## Prepare data access for the selected project

On the user's Cloud Computer, Access prepares the Notehub Cache Database and Notehub
Cache Maintenance Task, currently using DuckDB;
Project Selection enrolls each project on first use. Reuse that infrastructure and
verify the project's cache state. Honor user-requested pauses and schedules. This baseline
supports ordinary Q&A as well as later monitors. Answer from verified available evidence
while initial history loads, and record fields queried so daily maintenance can optimize
the workload.

| If you need to | Read | Then |
|---|---|---|
| Create missing cache infrastructure or enroll a newly referenced project | [Cache preparation](../notehub-cloud-agent/cache-preparation.md) | Reuse the database and maintenance task and verify this project's enrollment |
| Prepare queries and check the project's cache coverage | [Cloud project setup](../notehub-cloud-agent/setup.md) | Use actual acquisition runs and product knowledge to prepare grounded answers |
| Poll for changes in Events, Sessions, Devices, and more | [Modified cursors](modified-cursors.md) | Read successive batches through MCP or CLI; cloud workers resume the saved checkpoint on each poll |
| Answer in a foreground session or while cache coverage is incomplete | [Query planning](querying.md) | Use a suitable repository or bounded live read for the missing evidence |

## Choose the work within Chat

| If the person wants to | Read | Complete the work by |
|---|---|---|
| Ask or investigate a question | [Query planning](querying.md) | Returning a grounded answer with the relevant scope, freshness and uncertainty |
| Get ready for fast Q&A or improve repeated questions | [Optimize Q&A](optimizing-qa.md) | Reusing validated plans and suitable evidence; preparing maintained caches on a capable Cloud Computer |
| Consider useful monitors after a question | [Suggest monitors](monitors.md#suggest-a-useful-monitor) | Proposing a concrete recurring question, benefit and notification behavior |
| “Keep checking,” “watch this,” or “notify me when the result changes” | [Create a monitor](monitors.md#create-or-change-a-monitor) | Writing the evaluator, registering the periodic job and verifying the requested delivery path |
| Receive alerts or understand an existing alert | [Alerts from monitors](monitors.md#alert-the-user-and-close-the-loop) | Explaining the monitor finding with current evidence and handling delivery, recovery or correction |
| Receive periodic status, graphs or relevant context | [Engagement](engagement.md) | Giving useful updates tied to monitors and interests, including unchanged status when appropriate |
| Change project, device or other operational state | [Operations](operations.md) | Performing the authorized change and verifying its actual effect |
| Explicitly teach or revise shared product knowledge | [Project training](../notehub-project-train/SKILL.md) | Switching to the dedicated training activity with the existing context |

## Answer, investigate or act

1. Identify the quantity, population, time basis and intended decision. Reuse explicit
   context; ask only about ambiguity that changes the answer. Analytical times default
   to UTC unless the user specifies otherwise. Use stable IDs for joins and names/serials
   for display; partition identities and results by project.
2. Apply the relevant platform definitions and product recipe/dependencies. For an unmatched
   product question, derive a bounded interpretation from its inventory. Preserve status,
   units, validity, version/cohort rules, reductions and population definition. A blocked
   recipe can support a narrower answer or focused clarification without starting training.
3. Choose evidence: fresh complete cache when it covers the question; a suitable repository
   for missing history or analytics; live API for current/missing state. Keep the selected
   authenticated transport and use modified cursors for change scans where supported.
   Inspect representative rows and edge cases before relying on
   an unfamiliar aggregate. Answer a supported first question while larger history initializes.
4. Give the result and practical meaning in the right form, with scope/window, freshness
   and limits that affect the decision. Distinguish full Sessions from `_session.qo`,
   requested from effective state, and unknown from zero. Do not invent thresholds or
   hide unobserved populations behind a fleet-wide conclusion.
5. For a requested change, resolve the target, use the current operation and relevant
   guards, perform the authorized change, then verify its actual effect. An accepted
   request is not proof a device applied it.
6. Carry the question's meaning into follow-ups. A request for notification can turn the
   same calculation into a monitor; a useful monitor can support later status updates.
   Preserve its population, filters, units and time semantics rather than starting over.

For a bare readiness request, offer useful questions supported by actual platform or
product knowledge and report cloud preparation accurately. Foreground agents answer
through repositories or bounded live calls; recurring execution requires a verified
Cloud Computer capability. Personal interests, query plans, monitor state and delivery
history remain private harness state, separate from shared product skills.
