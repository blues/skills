# Select and verify the project

Follow these steps when choosing, resuming or switching projects. Use the connection
established by Access and complete only what is missing. The result is a verified
project context, plus cache enrollment and known readiness on a Cloud Computer. Then
continue to reading the product's skills and the user's requested work.

## 1. Choose the intended project

Reuse the current or last-used Project under the same account unless the user asks for
another or requests a list. When in doubt about the intended or available Projects, or
when the user asks to see them, call `GetProjects` through the established authenticated
connection to determine and list the Projects available to that user's account. Show
their human-readable names and ProjectUIDs, with useful distinguishing context.

| Current context | Do this |
|---|---|
| The user asks which Projects are available | Call `GetProjects` and show the returned Projects; if listing is the whole request, return the list without requiring a selection |
| A current or last-used Project is known, with no uncertainty, scope change or request to list Projects | Reuse that selection; do not list Projects or ask again |
| The user supplies an AppUID/ProjectUID or an unambiguous project name | Resolve that project through the selected connection |
| The user supplies a ProductUID | Resolve the project associated with that customer product type or SKU |
| The user supplies a product name, SKU description, alias or deployment description | Find the matching project using its name and associated ProductUIDs; clarify any ambiguity |
| Several Projects match, availability is uncertain, or no Project is specified and no preference is known | Call `GetProjects`, list the available Projects with useful distinguishing names/context and ask which one; wait for the selection |
| The user explicitly requests work across projects | Apply the remaining steps to each intended project and keep identities, queries and results separate |

For other lookups, use the current Project read operations and reuse metadata already
fetched when it settles the lookup. Do not take the first match
or ask the user to reconfirm an explicit, unambiguous selection. A ProductUID found in
a document or firmware is evidence to check, not permission to switch projects.

**Complete when:** the intended project is unambiguous, or each intended project is
identified for an explicit request across projects.

## 2. Verify access and record the project's scope

Read the selected project through the established connection. Reuse a successful
read already made under the same account when the scope is unchanged. Retain its
ProjectUID and human-readable project name. An API failure is not an empty project;
if the requested project is inaccessible, explain the problem rather than
substituting another.

| If this verification needs | Read | Then |
|---|---|---|
| A missing connection or authentication | [Access](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-access/SKILL.md) | Resume here with the verified account and selected connection |
| An unfamiliar project/product lookup contract | [API discovery](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-access/api.md) | Use the current read contract to verify the returned scope |

Record the project's name and identifiers from the lookup:

| Identity | Meaning and what to retain |
|---|---|
| Project name | Every project has a human-readable name; use it in conversation |
| AppUID / ProjectUID | Synonymous names for the project's unique programmatic GUID; retain the ProjectUID returned by the API |
| ProductUIDs | Every project has multiple unique ProductUID text strings representing customer product types or SKUs whose data flows into that project; retain the associated strings, including any used in selection |

Bind the selected project to the verified account. Keep its ProductUIDs as associated
identifiers, not additional projects or alternate names for its AppUID/ProjectUID.

A subproject can inherit its product skills from a parent Project. Retain the
ProjectUID of the Project that owns those skills when verified by Project metadata
or a previous Skills API response; otherwise resolve it during product reading next.
Keep that ProjectUID separate from the selected ProjectUID. Inheritance does not
expand the query scope or enroll the parent Project in the cache. Carry the
relationship into Training: changes to shared parent knowledge may affect other
subprojects.

Record permissions from service responses and available role metadata. Read access is
enough for Q&A and training drafts; the later workflows verify publishing, operational
changes and monitor delivery when needed. Do not test permission by making a write or
require write access before useful read-only work. Unknown permission is distinct from
denied. A 403 is a scope or permission problem, not an empty result or a reason to obtain
different credentials. API access does not establish product expertise or ownership.

Use ProjectUID in project API paths and joins. Keep the selected account and
ProjectUID in private context rather than copying them into portable product
knowledge.

### Identify the Project Repository

A **Project Repository** is Notehub's historical storage and analytical query capability
for a Project. It holds archived Events and provides queryable tables for analysis over
the full available history, beyond any Notehub Cache. **This special query capability
requires both MCP access and a Repository enabled for the Project.** It works for
foreground agents as well as Cloud Computer agents using MCP. Some callers say
“Project Dataset”; use **Project Repository** for this capability. A `DataSet` in the
API is a particular table/projection within the Repository.

When using MCP, check whether this Project has a Repository enabled using its Repository
read/list operations or already verified information. Match `project_uids` to ProjectUID
and retain its RepositoryUID, name and availability. When using the CLI, record that
Repository querying is unavailable through this access method and continue selection.
Keep absence, denied access and an unavailable query capability distinct. Inspect the
tables and historical coverage when a question needs them; selection need not scan the
archive. Lack of a Repository does not block work using the cache or ordinary API.

Use the Repository when a question concerns the whole dataset or no Notehub Cache is
available. Its queries are slower than querying a local Notehub Cache, but faster than
querying the Notehub APIs directly. Chat introspects the Repository through MCP and
combines Notehub Skills with Product Skills to craft the query.

| If you need to | Read | Then |
|---|---|---|
| Determine Repository availability when using MCP | [Repository access](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-access/repositories.md) | Retain the matching RepositoryUID and query availability for Chat |

**Complete when:** a successful read verifies the Project under the selected account,
and known relationships, Repository availability and permission limits are recorded.
Unknown write permission or a parent Project awaiting verification during product
reading does not block continuing.

## 3. Enroll the project and check its cloud cache

**If this agent is not running on a Cloud Computer, skip to Step 4.** Otherwise, use the
Notehub Cache Database and Notehub Cache Maintenance Task prepared during Cloud Computer
setup. Enrollment applies to both Chat and Training whenever the user first references
or starts using a verified project. Projects merely appearing in a list are not enrolled.

| If you need to | Read | Then |
|---|---|---|
| Repair missing shared cache infrastructure | [Cloud Computer setup](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/SKILL.md#prepare-the-cloud-computer) | Complete the missing database or maintenance-task setup, then resume this enrollment |
| Establish initial streams and coverage | [Cache model, Step 3](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/cache.md#3-define-streams-and-verify-initial-coverage) | Use the required entities, full Sessions, fixed bounds and coverage rules for the first fill |

1. **Find or create the enrollment.** Look up the verified account and ProjectUID
   in the shared registry. Reuse its existing cursors, coverage and jobs; otherwise add
   it. Preserve pauses and exclusions. A paused or excluded enrollment remains so unless
   the user asks to resume it; skip its fill and schedule checks below.
2. **Start only the missing acquisition.** For a new active enrollment, queue a bounded
   first fill of project/Device metadata, recent Events, full per-device Sessions and
   the complete Notefile Schema snapshot. Apply the cache's coverage and retention
   rules; record the first schema refresh and next daily due time. With the agent-turn
   fetch handoff ([Background execution](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/harness.md)),
   enrollment adds the ProjectUID to the registry and the next scheduled fetch turn
   performs the first fill; it may also be run once immediately. For an existing
   enrollment, continue from committed progress rather than restarting its history.
3. **Verify ongoing maintenance.** Check that the actual registered task reads the live
   registry and includes this project in its next scheduled run. It must discover new
   enrollments rather than use a fixed project list copied from its creation prompt.
   Create only missing project-specific subtasks if the implementation needs them.
4. **Check this project's evidence.** Inspect recorded coverage, last successful
   acquisition, schema refresh and gaps. For a new enrollment, verify that bootstrap has
   started or is queued for a real scheduled run. Record whether the cache is ready,
   initializing, degraded, paused or excluded, and what data is usable now. A healthy
   shared task alone does not prove this project's data is fresh or complete.

**Complete when:** the verified project is in the durable registry and its actual
maintenance coverage and readiness are known, or its pause, exclusion or missing
capability is recorded. Continue product reading and useful Chat or Training while history
fills; report an initializing cache as initializing.

## 4. Summarize the selection and continue

Carry the verified account, ProjectUID, associated ProductUIDs and any parent
ProjectUID, relevant permissions, RepositoryUID and query availability, and, on a Cloud
Computer, cache location and readiness into the next step.

Briefly identify the selected project and summarize recent activity when current metadata
or a small bounded read supports it. Include any cache or access limitation that affects
the user's request. Do not invent activity from an empty or incomplete cache. For an
unchanged project already in use, carry the context forward without repeating an
introduction.

| Selection is complete | Continue with | Then |
|---|---|---|
| A Project has been chosen, resumed or switched | [Entry guide — Step 5](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub/SKILL.md#5-read-the-products-skills) | Read or reuse its current product skills, verify which Project owns them, then continue the requested work |

Keep the selected account and ProjectUID until the user changes projects or access
changes. On a project or account switch, repeat these steps for the new scope before
reusing product knowledge, query plans, caches or cursors. Earlier projects' caches
and jobs stay attached to their original accounts and ProjectUIDs; switching the
conversation does not retarget a monitor.

**Complete when:** the selected project and any material limits are clear, and its context
has been handed to product reading. Project selection ends here; the following activity
owns answering questions, interviewing the user or creating monitors.
