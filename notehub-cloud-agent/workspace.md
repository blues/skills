# Keep only the state needed between runs

Use the Notehub Cache Database for structured state that later runs need. Let the agent
choose its tables, fields and organization; there is no required assistant-state schema
or catalogue of records to create in advance.

Start with enough state to identify enrolled projects and their active/paused status,
resume each acquisition stream from its exact cursor, assess cache coverage and freshness,
and find the existing maintenance task. Scheduled work must be able to resume without
conversation memory. Add other state when a feature needs it: for example, a requested
monitor needs its definition, comparison baseline and delivery progress. Reuse information
already available reliably from the harness instead of maintaining duplicate records.

## Choose the relevant reference

| If you need to | Read | Then |
|---|---|---|
| Prepare the database and its maintenance task | [Cache preparation](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/cache-preparation.md) | Create and verify the two deliverables |
| Store cached data and acquisition progress | [Cache database](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/cache.md) | Preserve project boundaries, exact checkpoints and honest coverage |
| Run or resume scheduled work | [Job contracts](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/jobs.md) | Persist the inputs and progress that the work actually needs |
| Implement a requested monitor or alert | [Monitors](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-chat/monitors.md) | Store the state needed for its comparison and delivery behavior |

Keep credentials in the selected MCP or CLI. Commit fetched records with their cursor;
commit monitor transitions with pending deliveries. Preserve user choices and delivery
history when rebuilding cached source data. Stored configuration is not proof that a
task ran or a message was delivered; verify those through the actual harness.

The harness needs enough bootstrap information to open the database and invoke the
worker safely. Files are appropriate for executable code, engine storage or temporary
transfers when required; structured state does not need a parallel directory layout.
Keep a short README beside the shared cache database and a shared memory fact pointing
all bots to that folder, as described in Cache preparation. These explain how to find
and use the cache; live checkpoints and freshness remain in the database.
