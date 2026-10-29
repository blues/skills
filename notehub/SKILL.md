---
name: notehub
description: Start or resume Notehub work by establishing secure access, learning the API, checking Cloud Computer setup, selecting a project and reading its product knowledge. Continue in Project Chat or Project Training; return to project selection when the user finishes or switches projects.
---

# Start working with Notehub

Establish access, understand how to use Notehub, prepare the Cloud Computer when available,
then choose a project, read its product skills and do the requested work.
Keep the user's actual request in view. A question should lead to an answer; a request
for continuing assistance should lead to working, verified jobs; training should improve
shared product knowledge.

Complete the prerequisites in order, reusing verified context. Read the instructions
for your path in full; use the routing tables to choose relevant references rather than
fetching the whole corpus. Follow [working style](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub/working-style.md) throughout.

## 1. Establish access

Use one authenticated request proxy: prefer Notehub MCP, or use the `notehub` CLI when
MCP is unavailable or the user chooses it. Access owns installation/connection,
authentication and verification. Reuse an already verified connection.

| If you need to | Read | Then |
|---|---|---|
| Install or connect a tool, authenticate and verify access | [Access](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-access/SKILL.md) | Retain the selected connection, verified account and invocation method |

## 2. Read Notehub's skills

**All agents must complete API discovery**, regardless of how they accessed Notehub.
Read the complete OpenAPI specification and survey its available operations. Then read
**only one** of the MCP or CLI guides, matching the connection established in Access:
the MCP guide for MCP access, or the CLI guide for CLI access. Reuse what you learned
during authentication. Resolve detailed request and response contracts when needed.
Load Platform's native data meanings later according to the work.

| If you are using or need to understand | Read | Then |
|---|---|---|
| Notehub's API capabilities through either connection | [API discovery](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-access/api.md) | Read and retain the current OpenAPI contract before learning the selected tool's invocation syntax |
| Notehub MCP: an authenticated request proxy | [MCP guide](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-access/mcp.md) | Use 'code mode' to map OperationId and parameters to an 'execute' MCP tool |
| Notehub CLI: an authenticated request proxy | [CLI guide](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-access/cli.md) | Use the tool's CURL-like syntax as an API proxy; read `notehub --help-agent` |

## 3. Check Cloud Computer setup

A subset of agent harnesses run on Cloud Computers with persistent storage, such as
Grok Bot, Muse, Dots, or 'claws'. If and only if this agent is running on a Cloud Computer,
ensure that a **Notehub Cache Database** and **Notehub Cache Maintenance Task** are
configured properly and that the maintenance task is running.

| Operating capability | Read | Then |
|---|---|---|
| This agent is running on a Cloud Computer | [Cloud Computer setup](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/SKILL.md#prepare-the-cloud-computer) | Verify the setup, then continue to project selection |
| This agent is not running on a Cloud Computer | No cloud setup | Continue to project selection |

## 4. Choose the project

A Notehub project is the place where customers manage their fleets of physical
hardware products as well as the data flowing out of those devices. Users and their
agents only work on one project at a time, and if possible prefer to be brought
immediately back into the context of the last project they worked on. If their
preferred project is not yet known, users prefer to see a list of their projects
from which to choose.

| If you need to | Read | Then |
|---|---|---|
| Resume or choose a Project | [Project selection](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-access/projects.md) | Verify the ProjectUID and summarize the Project's name, associated ProductUIDs and recent activity |

## 5. Read the product's skills

Understand the hardware products whose data flows into the selected Project, their
terminology and how product concepts map to Event `body` fields or database columns.
Read all skills and their `kinds` tags in one `GetSkills` call with `includeContent=true`,
or reuse a complete, current snapshot already read. Internalize the entire set and
carry that understanding into both Q&A and Training.

| If you need to | Read | Then |
|---|---|---|
| Load and understand the selected Project's product knowledge | [Product knowledge cookbook](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-chat/product-knowledge.md) | Acquire the complete skill set with content and kinds, read every skill, and connect product terminology to the data |

An empty skill set is a knowledge gap, not a reason to invent product meaning or
start an unsolicited interview. Even without product skills, useful platform
questions can still be asked - and sometimes the nature of products can be
determined by notefile and field names in the JSON data.

## 6. Do the requested work

**Both activities are ongoing loops.** The user may spend a long time asking questions,
refining monitors, or improving product skills. Keep working within the selected activity
and project, reusing context across turns; completing an answer or lesson is a checkpoint.
When the user says they are done, says “stop,” or wants to switch projects, end the current
loop and return to **Step 4**. Preserve any unfinished work and a brief resumable handoff;
in Training, keep confirmed knowledge and unpublished drafts distinct. A stop directed
at a particular monitor or task applies to that work; ending the conversation does not
cancel existing monitors or cache maintenance.

| The person wants | Continue here | Work within the loop |
|---|---|---|
| Ask questions, investigate, act, or monitor answers with notifications and ongoing updates | [Project chat](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-chat/SKILL.md) | Answer using platform and/or product knowledge; carry follow-ups into operations and monitors that notify according to the user's preferences |
| Teach, resume training or revise shared product meaning | [Project training](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-train/SKILL.md) | Interview, investigate, review, publish verified improvements and continue with the next useful lesson |
| Finish, stop or switch projects | [Project selection — Step 4](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub/SKILL.md#4-choose-the-project) | End the current activity; wait after a stop, or select the requested project and continue through Steps 5 and 6 |

For a bare “load” request, say which project is ready and, when helpful, give a few example
questions grounded in its actual knowledge. State cloud readiness accurately and wait
for the user's work. Do not invent product capabilities, start an unsolicited interview,
or activate a monitor from a bare load request.
