# Read product skills through the standard Skills API

Project skills and product skills are the same shared knowledge. Reading them is part
of ordinary Project Chat, Q&A, monitors and cloud preparation, as well as Training.
Load the complete set with **`GetSkills` and `includeContent=true`** through the selected
MCP or CLI connection. This includes every document's exact content and `kinds` tags in
one API call. The reading workflow is the same for both connections; only the invocation
syntax differs. Reuse a complete, current snapshot already fetched.

This is the API contract supporting the product-knowledge cookbook. After acquiring or
refreshing content, continue that cookbook's current step; an API download alone does
not establish that the agent has read and understood the skills.

## Read the complete skill set

Use `GET /v1/projects/{projectOrProductUID}/skills?includeContent=true`, supplying the
selected ProjectUID as `projectOrProductUID`. Preserve the complete response and read
all documents, including unlinked skills. Both Q&A and Training need the product
understanding; Training also preserves the snapshot as its immutable baseline.

| Reading need | API operation | What to do with the result |
|---|---|---|
| Initial product knowledge for Q&A or Training | `GetSkills` with `includeContent=true` | Acquire all content and metadata in one consistent snapshot, then read and understand every skill |
| The full-content request returns 413 | `GetSkills` with `includeContent=false`, followed by `GetSkill` for each catalog entry | Acquire the complete set through individual reads, preserving consistent revisions as described below |
| Reuse knowledge in a later question or scheduled run | Reuse the complete snapshot; check `GetSkills` metadata when freshness requires it | Keep unchanged documents, read all new or changed documents and reconcile removals |
| Refresh one changed document or complete the fallback | `GetSkill` with the catalog's skill name | Read its exact content and recheck affected dependencies; a link to already loaded content needs no API call |

Fetching a complete set does not require putting every document into model context at
once. Retain it in private task storage where available and read it in groups until the
whole set is understood. This reading workflow does not publish changes; skill writes
remain part of the dedicated Training workflow.

## Inspect the response and kinds

`GetSkills` returns `project_uid`, `set_revision`, `total`, `skills` and `complete: true`.
There is no list pagination or modified cursor. Each entry includes `name`, `kinds`,
`size_bytes`, `revision`, timestamps and `links`, plus exact UTF-8 `content` when
`includeContent=true`. Omitting `includeContent` or setting it to `false` returns metadata
only, useful for change checks and the 413 fallback. An empty `skills: []` is a valid
untrained set; an API failure is not an empty catalog.

`kinds` are API metadata independent of Markdown front matter. A document may have
several kinds. Inspect them locally to locate an `index`, product context, field
definitions and recipes; do not invent a server-side kind filter. Use the index and
actual content to organize reading and recognize which rules apply to a question.
Missing or unfamiliar kinds do not exclude a document from the full reading, and a tag
alone does not establish what its content means.

For subprojects, reads inherit skills from the parent Project. Retain the returned
`project_uid` separately from the selected ProjectUID. Use the Project's actual
access settings when explaining who can read the knowledge; shared does not mean
private to a particular role.

## Handle individual reads and content errors

`GET .../skills/{skillName}` returns one document's metadata and exact UTF-8 `content`.
Use the catalog's relative `links.self` as an API path through the selected connection,
or encode the **entire** skill name once as one URL path component, including any `/`. For example,
`references/devices.md` becomes `references%2Fdevices.md`. Pass the original unencoded
name as `skillName` to MCP code mode. Resolve document links within this project-owned
catalog, not as files on the public skills site.

`GET .../skills?includeContent=true` includes exact content for every entry in one
consistent snapshot. Preserve the names and contents; empty document content is valid.
On a 413 response, use the metadata catalog and individual reads instead. A 409 indicating
invalid stored content needs an explicit Training repair through the documented API;
report the problem rather than attempting an undocumented legacy fallback during Chat.

## Reuse reads and recognize changes

Retain `set_revision` with the catalog and each document's `revision` with its content.
These opaque values let a reader avoid downloading unchanged skills and recognize when
the meaning behind a question or monitor has changed. Refresh the catalog at the task's
freshness boundary, compare names and revisions, and acquire and read all new or changed
documents. Reconcile removed entries in the private copy and recheck affected query plans
or monitors. Reading or refreshing this knowledge does not publish anything.

When assembling documents from separate calls, require each catalog revision with
`GetSkill`'s optional `If-Match`, or compare the returned revision to the catalog. Recheck
`set_revision` after acquisition; if it changed or a document returned 412, refresh the
catalog and reconcile before treating the assembled set as one snapshot. Do not combine
conflicting versions into an answer. Unchanged cached documents can be reused.

Where the transport exposes response headers, `If-None-Match` can skip unchanged reads.
Metadata-only and full-content `GetSkills` responses have different ETags: use the
validator for the same representation. A 304 reuses that matching cached representation.
When headers are unavailable, compare the returned body revisions instead. No revision
handling requires writing a skill.
