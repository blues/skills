# Turn platform and product knowledge into prepared questions

Use this to implement Project Chat's reusable calculations during cloud setup or when
a new question or meaning change requires a plan. Use the relevant platform definitions,
product skills or both.
A prepared question connects a human intent to a verified calculation; it is neither a
stored answer nor a generic dashboard query.

## Choose the relevant reference

Read the rows that apply to the current step; reuse context already established.

| If you need to | Read | Then |
|---|---|---|
| Interpret native quantities, records, location or time | [Platform](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-platform/SKILL.md) | Resolve the native meanings used by the calculation |
| Load or refresh the product meaning behind a question | [Reading product knowledge](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-chat/product-knowledge.md) | Resolve the recipe and all consequential dependencies |
| Turn a requested question into a monitor | [Monitors](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-chat/monitors.md) | Define its predicate, freshness, state and delivery semantics |
| Make an authorized monitor run | [Activation](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/setup.md#activate-recurring-work) | Install, exercise and verify the actual scheduled work |
| Optimize the plan from real usage | [Query performance](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/query-performance.md) | Record fields and measured performance for daily schema and layout review |

## From recipe to executable plan

For each useful question family, record a private plan with these parts:

1. **Intent:** real wording or an explicitly labeled proposed question, audience and
   decision. Match future wording by meaning, not exact strings. Preserve required
   clarification about device identity, population, period or interpretation.
2. **Meaning:** owning platform definitions and/or product documents/headings/revisions, field paths, units, scaling,
   validity/sentinels, version branches, population definition, counting unit, reductions
   and constraints. Resolve dependencies before writing SQL. A platform-only question
   does not require a product recipe; preserve the native definitions it uses.
3. **Evidence:** required entities and columns, measurement/modification time distinction,
   window and freshness budget. Specify what complete coverage means and how unknown
   devices/periods remain visible. Choose cache, repository or live fallback explicitly.
4. **Execution:** a parameterized SQL query or script plus its input/result contract.
   Use typed per-Notefile views over retained raw JSON. Parameters supply identifiers,
   filters and windows; never interpolate an arbitrary user string as executable SQL.
   Record field roles and a query fingerprint so executions can update the private
   usage/performance history without retaining unnecessary literal values.
5. **Presentation:** the result shape, units, precision, source/as-of time, unknown counts
   and the answer template. Keep actual result values out of reusable product knowledge.
6. **Checks/status:** representative and boundary cases, source comparison, a new phrasing,
   validated skill revisions, and supported/conditional/blocked/untried status.

Don't mark a plan verified merely because its SQL runs. Native platform definitions
matter too: use full Sessions for session accounting, distinguish DFU desired/reported
state, and treat timestamps and nullable fields according to their documented meaning.
A small product can need only a few plans; don't enumerate every possible question.

## Example of the semantic distinction

Suppose a **synthetic** product recipe defines a valid temperature in degrees Celsius,
with an explicit fault marker and a confirmed limit. A device has a valid above-limit
reading in the morning and a fault at noon. These are different prepared questions:

| Question | Calculation and output |
|---|---|
| Which devices are above the limit now? | Latest observation per device, then validity; the faulty device is unknown now |
| Which devices exceeded the limit today? | Valid observations in the UTC day (or requested timezone); that device did exceed earlier, with its later fault reported separately |
| Was it above the limit continuously for an hour? | Requires evidence of duration, such as a documented interval-bound measurement; two separated samples alone do not establish it |

The same facts can support the first two plans and block the third. Turning “above the
limit” into a monitor must keep the exact intended calculation, not just reuse a field
and threshold. The example supplies no real product threshold or freshness rule.

## Choose and validate a backend

For recent, covered questions, use the cache and record its stream/coverage requirements
beside the query. Derive typed fields only from supported product meanings. If repository
history is used, inspect its actual schema and validate equivalent semantics; a similarly
named column may be transformed or incomplete. A live fallback should be a bounded query
through the selected connection, not a second credential path.

Test the cases that change the answer: invalid latest record, silent device, legitimate
zero versus sentinel, duplicate identity versus equal value, counter reset, firmware
branch, window boundary and partial population. Use the relevant subset, not a ritual
list for every trivial question. Label synthetic fixtures. Compare one bounded case with
source evidence and retain the check result. A changed skill revision makes dependent
plans pending review until meaning changes are reconciled; unrelated queries can continue.

## Reuse the plan for questions and monitors

Foreground answers and scheduled evaluators should share the same calculation and result
contract. A result includes value/list, units, scope/window, source/as-of time, coverage,
unknowns and caveats. Apply the user's new parameters, check freshness and execute; do not
serve yesterday's numeric answer from the plan. If conditions differ from those validated,
recheck or narrow the claim.

A monitor adds a condition/change rule, cadence, freshness, incident state, destination
and delivery policy to a prepared question. It does not add missing meaning. Follow
Monitors and activation to make it run.
An interesting recurring question can become an engagement candidate within the user's
standing scope; a catalog entry alone is not authorization to send alerts.
