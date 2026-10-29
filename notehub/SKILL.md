---
name: notehub
description: Start or resume Notehub work by establishing secure access, selecting a project and reading its product knowledge, then continue with Project Chat or Project Training. Use when the user asks to load Notehub or needs the overall workflow.
---

# Start working with Notehub

> Establish access, choose a project and read its product skills. Cloud Computers prepare
> a Notehub Cache Database and Notehub Cache Maintenance Task during Access, then enroll
> projects on first use. Then chat or train: Chat covers questions, operations,
> monitors, alerts and engagement.

Keep the user's actual request in view. “Load {{.DOMAIN}}” means become ready to help.
A question should lead to an answer; a request to configure continuing assistance should
lead to working, verified jobs; a training request should lead to shared product knowledge.
Read the complete pages linked for your path, rather than fetching the whole corpus.
Use the routing tables to select the relevant path, not as a checklist to read every page.
Complete the prerequisites in order, reusing any connection, environment or project
already verified for this task.
Follow [working style](working-style.md) throughout.

## 1. Establish access

Access owns connection setup, authentication and preparation of the execution environment,
including a Notehub Cache Database and Notehub Cache Maintenance Task on the user's
Cloud Computer. DuckDB is the current database of choice.

| If you need to | Read | Then |
|---|---|---|
| Establish missing connection or environment prerequisites | [Access](../notehub-access/SKILL.md) | Carry the verified access context into project selection |

## 2. Choose the project

Project Selection owns which project's data and shared skills the task uses. On a Cloud
Computer, each project the user first references is enrolled in the existing cache
database and maintenance task, including projects used during Training.

| If you need to | Read | Then |
|---|---|---|
| Establish or change the project binding | [Project selection](../notehub-access/projects.md) | Carry the verified project context into product reading and the requested activity |

## 3. Read the product's skills

Understand the selected product before working with its data. Reuse its current catalog
and knowledge if already loaded; otherwise discover the stored skills, read the index
and product overview, and follow the material relevant to the task.

| If you need to | Read | Then |
|---|---|---|
| Understand the selected product and find its knowledge | [Reading product knowledge](../notehub-project-chat/product-knowledge.md) | List all skill names and kinds, then fetch needed documents, or read a full-content snapshot; retain the product context |

An empty skill set is a knowledge gap, not a reason to invent product meaning or start
an unsolicited interview. Platform questions can still use verified native knowledge.
Training reads the complete baseline, including unlinked files, in its own workflow.

## 4. Check cloud cache readiness

| If you need to | Read | Then |
|---|---|---|
| Check a newly enrolled project's cache or repair missing infrastructure | [Cache preparation](../notehub-cloud-agent/cache-preparation.md) | Verify the database, maintenance task, project enrollment and actual readiness; complete missing setup |
| Prepare useful questions and requested recurring work for the project | [Cloud project setup](../notehub-cloud-agent/setup.md) | Build on the maintained cache, product knowledge and actual coverage |

Reuse the working maintenance task and respect pauses. Begin supported Q&A or the training interview
while initial history loads; user-facing monitors follow the requested Chat workflow.
Foreground agents use bounded reads without persistent cache jobs.

## 5. Do the requested work

There are two high-level activities. Access, Platform and Cloud Computer are supporting
skills used within these workflows.

| The person wants | Continue here | A useful completion looks like |
|---|---|---|
| Ask questions, investigate, act, monitor answers, receive alerts or get ongoing updates | [Project chat](../notehub-project-chat/SKILL.md) | Grounded answers using platform and/or product knowledge; requested operations and monitors completed and verified |
| To teach, resume training or revise shared product meaning | [Project training](../notehub-project-train/SKILL.md) | Preserved knowledge, reviewed changes, verified publication and a resumable handoff |

For a bare “load” request, say which project is ready and, when helpful, give a few example
questions grounded in its actual knowledge. State cloud readiness accurately. Do not
invent product capabilities, start an unsolicited interview, or create alerts from a
question alone.

## Reading and installation

Fetch raw text with a tool that returns the whole page, for example:

```sh
curl -fsSL https://{{.DOMAIN}}/notehub-access/SKILL.md
```

The server qualifies public relative links. Fetch only the references your current step
requires; completing a prerequisite returns you to the task that called it. Installed
skills are live pointers, so they follow the same journey without reinstalling every update:

```sh
npx skills add https://{{.DOMAIN}}
```

This top-level `notehub` skill routes to two activity skills and three supporting skills.
An older project-access entry remains for compatibility and routes to Access or Project Selection.

| If you need | Use | Purpose |
|---|---|---|
| Compatibility with an older project-access installation | [Project-access entry](../notehub-project-access/SKILL.md) | Follow its route to the missing access or project prerequisite |
| A discovery index | [llms.txt](../llms.txt) | Find the relevant entrypoint |
| The complete corpus, expressly needed for the task | [llms-full.txt](../llms-full.txt) | Read the aggregated public skills |
| Machine-readable installation metadata | [Skill manifest](../.well-known/skills/index.json) | Discover installable skills |

If a required page is unavailable, identify it rather than improvising its
instructions. Public documentation downloads do not require Notehub credentials.
