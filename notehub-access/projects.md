# Select and verify the project

Use the connection established by Access. This step owns project resolution, read
verification, data-versus-skill-owner scope and the binding carried into Chat or Training.
Complete only what is missing, then return to the caller with that project context.

## Resolve the intended project

| Current context | Do this | Result |
|---|---|---|
| Project already verified under the same hub/account and the user has not changed scope | Reuse the selection; do not list projects or ask again | Continue with the existing binding |
| The user supplies a project UID or an unambiguous project name | Resolve and verify that project through the selected connection | Retain the canonical UID and human-readable name |
| The user supplies a product name, product UID, alias or deployment description | Resolve the candidate projects and their actual relationship | Choose the intended data scope; a product can span deployments |
| Several projects match, or no project is specified | List accessible candidates with useful distinguishing names/context and ask which one | Wait for the user's selection rather than taking the first result |
| The user explicitly requests work across projects | Verify each intended project and retain separate bindings | Keep identities, queries and results partitioned by project |

Use the current project read/list operations (`GetProjects` lists accessible projects)
when resolution is needed. Reuse metadata already fetched during Access where it settles
the lookup. Once resolved, verify the selected project through a read; do not confirm an
explicit, unambiguous selection with the user again. Do not treat API failure as an empty
project or silently substitute another project when the requested one is inaccessible.

## Keep the scopes distinct

| Scope | Meaning | What to preserve |
|---|---|---|
| Product | The product identity/grouping, potentially spanning deployments | The user's product name, verified aliases and any product UID used to resolve the project |
| Data project | The concrete scope whose devices, events and configuration answer the request | Canonical project UID and name, bound to the verified hub/account |
| Skill-owning project | The project that owns the product/project skills being read or maintained | The verified owner and any parent/subproject inheritance relationship |

Product skills are also called project skills because they are stored in the product's
project. A subproject can inherit them from a parent. Resolve that owner when the task
uses stored skills, and carry it separately from the data project. An inherited skill
owner does not expand the scope of queries or device actions. Record whether a shared
parent means a knowledge change could affect other subprojects; the calling activity
uses that scope when reviewing its work.

A product UID found in a document or firmware is evidence to check, not permission to
switch the selected project. Use canonical UIDs in API paths and joins, human names in
conversation. Keep runtime bindings in private context rather than copying them into
portable product knowledge. API access does not establish the person's product expertise
or ownership.

## Verify the permissions the task needs

Use service responses and available role metadata. Read access can support Q&A and
training drafts; publishing or an operational change needs the corresponding permission.
Do not probe permission by making a write or demand write access before beginning useful
read-only work. Unknown permission is distinct from denied. A 403 identifies a scope or
permission problem, not an empty result or an invitation to obtain different credentials.

Repository availability and monitor delivery are checked by the workflows that need them.
On a Cloud Computer, a verified selection also triggers the cache enrollment below;
waiting for the entire cache to fill is not a prerequisite for useful work.

## Enroll the project in the Cloud Computer cache

For each project the user first references or starts using, enroll the verified
hub/account/project binding in the Notehub Cache Database and Notehub Cache Maintenance
Task prepared by Access. This applies
to Chat and Training. Reuse existing enrollment and preserve pauses/exclusions. An inherited
skill owner or a project merely appearing in a list does not expand the data scope.

| If you need to | Read | Then |
|---|---|---|
| Add a newly referenced project or recover its existing cache binding | [First-use cache enrollment](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/cache-preparation.md#enroll-a-project-on-first-use) | Add it to the maintenance task's registry, queue a bounded initial fill, and retain its readiness state |

Continue product reading and the requested task while initial history loads. A foreground
agent carries the verified project context without creating persistent cache jobs.

## Return and preserve the project context

| Context | What to retain |
|---|---|
| Account binding | Hub/account from Access under which the project was verified |
| Data binding | Canonical project UID/name and any resolved product/alias context |
| Skill binding | Verified owning project/inheritance when relevant; unresolved ownership stays explicit |
| Permission | Observed read/write capabilities and any limitation relevant to the requested task |
| Cloud cache | Enrollment/scope, database location and initializing/ready/degraded/paused state, or the precise missing capability |

Return to the calling step; the entry guide places product reading after selection,
checks cloud cache readiness, and routes into Chat or Training.
Keep the binding until the user changes it or evidence shows it no longer applies.
On a project or hub/account switch, resolve the new scope before reusing product knowledge,
query plans, caches or cursors. Keep existing projects' jobs attached to their original
bindings; switching the current conversation does not retarget a monitor.

## Resolve missing prerequisites or contracts

| If you need to | Read | Then |
|---|---|---|
| Establish a missing connection or execution environment | [Access](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-access/SKILL.md) | Return here with the verified hub/account and selected transport |
| Resolve an unfamiliar project/product lookup operation | [API discovery](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-access/api.md) | Use the current read contract and verify the actual returned scope |
| Read inherited product skills and identify their owner | [Read product skills](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-access/project-skills.md) | Retain the verified skill owner separately from the data project |
| Understand the selected product after Project Selection was invoked alone | [Entry guide](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub/SKILL.md#3-read-the-products-skills) | Read its product skills, then complete the preparation for the requested activity |
