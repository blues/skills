# Prepare a reusable question

Use this when a calculation will be reused for Q&A or monitoring. Reuse the selected
Project, loaded platform/product knowledge and any answer already verified. A prepared
question connects a human intent to an executable calculation. Keep it as small as the
question needs; a one-off answer does not require a catalog or new framework.

## 1. Recover the meaning and choose the evidence

Find a plan with the same quantity, population, time basis and intended use. Reuse it
when its source meanings and assumptions still hold. Otherwise establish those meanings
and the source that can supply the required evidence before writing the calculation.

| If something is unresolved | Procedure | Result |
|---|---|---|
| Quantity, scope, native or product meaning | [Question meaning](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-chat/querying.md#1-establish-what-the-question-means) | Grounded fields, units, validity and population rules |
| Source, coverage or freshness | [Choose the evidence](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-chat/querying.md#2-choose-the-cache-repository-or-notehub-api) | A suitable source and explicit evidence requirements |

A native platform question does not need a trained product recipe. An unsupported product
interpretation remains a gap; an executable query cannot supply missing meaning.

## 2. Implement or update the plan

Record the information needed to run and explain the calculation again:

| Part | Retain |
|---|---|
| Intent | Real wording, or a clearly proposed example; the intended decision and any necessary clarification |
| Meaning | Applicable platform definitions and product skill names/headings/revisions; fields, units, scaling, validity, reductions, counting unit and version scope |
| Evidence | ProjectUID, entities, population, time basis, required columns, coverage and freshness; how unknown Devices or periods remain visible |
| Execution | Parameterized query or script, inputs and result shape; source selection and any supported fallback |
| Presentation | Units, precision, source/as-of time, unknown counts and a useful output shape |
| Verification | Relevant checks, assumptions, meaning revisions and supported/conditional/blocked/untried status |

Match later wording by meaning. Similar words can require different calculations.
Parameters supply identifiers, filters and windows; do not interpolate arbitrary user
text as executable SQL. For cached data, reuse typed per-Notefile projections backed by
raw JSON. Record queried fields and their roles as the plan executes so cloud maintenance
can optimize real usage. Do not copy private results or executable state into shared skills.

## 3. Verify the answer the plan produces

Execute a bounded case and compare it with source evidence. Check the cases that could
change the answer, such as an invalid latest record, a silent Device, zero versus a
sentinel, equal values under distinct EventUIDs, counter resets, firmware variants or a
window boundary. Use synthetic fixtures only for cases unavailable in real evidence.
For broad Q&A preparation, also check that another natural-language phrasing resolves to
the same intended calculation.

| For calculation checks | Procedure | Result |
|---|---|---|
| Verify identities, validity, reductions and time assumptions | [Calculate and verify](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-chat/querying.md#4-calculate-and-verify-the-answer) | A result checked against its intended meaning, not just SQL that runs |

For example, a **synthetic** product has a valid above-limit temperature in the morning
and a fault at noon. These questions need different plans:

| Question | Calculation |
|---|---|
| Which Devices are above the limit now? | Select the latest observation, then apply validity; this Device is unknown now |
| Which Devices exceeded the limit today? | Select valid observations during the requested day; this Device exceeded earlier, with the later fault separate |
| Was it above the limit continuously for an hour? | Requires evidence supporting duration; two separated samples alone do not establish it |

The example supplies no real threshold or freshness rule. A monitor reusing a plan must
preserve its particular interpretation of “above the limit.”

## 4. Retain the plan and return to the caller

Keep the plan and verification privately in the database on a Cloud Computer, or within
the available foreground session. Reuse the calculation with new parameters and refreshed
evidence; do not serve an earlier numeric result as current data. If a meaning or schema
change invalidates an assumption, review the affected plan before relying on it again.

Return the executable calculation and its supported, conditional or blocked status to
Q&A preparation or monitor creation. That caller owns presentation, scheduling and
notification preferences. The plan adds no independent background task.
