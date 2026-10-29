# Explore or change Project state

Explore or investigate state of Devices, Fleets, or other Project entities. Make changes
if authorized. Reuse the selected Project, established MCP or CLI connection and product
knowledge. Follow Steps 1–2 for investigation; continue through Steps 3–5 for a change.

| If the task is specialized | Procedure | Result |
|---|---|---|
| Explore versions, investigate update progress or manage firmware | [Explore and manage Firmware Updates](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-firmware/SKILL.md) | Apply the Device and Fleet firmware cookbook |
| Resolve an unfamiliar operation's parameters or response | [OpenAPI contract](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-access/openapi.yaml) | Use the current contract through the already selected tool |

## 1. Identify and inspect the target

Resolve the entity and actual identifier: ProjectUID, DeviceUID, FleetUID, RouteUID or
the key the operation requires. Use the cache for exploration when it covers the question;
read current API state for gaps or before a change. Check relationships and effective
configuration where settings can be inherited.

For a bulk action, enumerate or count the intended population and verify filters. An
empty filter must not expand a Device or Fleet request to the whole Project. ProductUID
filters and inherited product skills do not change which ProjectUID owns the operation.

## 2. Explain the state and decide whether a change is requested

Explain the observations, discrepancies, freshness and missing evidence. For a question
or investigation, return these findings to Chat. A suspected problem does not itself
request a repair.

For a change, make the target, proposed modification and expected effect concrete.
Reuse existing authorization; ask only about unresolved decisions or actions outside
that scope. Preserve unrelated settings and working entities.

## 3. Apply the operation's contract

Use the selected proxy and current request definition. Apply dry runs, conditional
revisions or a dedicated single-Device operation when available and appropriate. Do not
invent a confirmation parameter or assume that a bulk request is a transaction.

Execute the bounded change and retain identifiers, returned state and any partial
failures. After a timeout, inspect actual state before retrying so a repeated request
does not produce an unintended second effect.

## 4. Verify the effect

Read back the changed state. For a Device command, distinguish accepted, queued, delivered,
acknowledged and effective according to the operation's actual contract. A service
acceptance does not prove that the Device applied it. If effect depends on a later
contact, report what is verified and what remains pending.

## 5. Return to the conversation

Give Chat the findings or verified change, affected identifiers and remaining limitation.
Return to Step 3 of the Chat loop. That loop handles follow-up questions, requested
monitoring and ongoing interests; an operation does not start a second conversation loop.
