# Read product skills through the standard Skills API

Project skills and product skills are the same shared knowledge. Reading them is part
of ordinary Project Chat, Q&A, monitors and cloud preparation. Use the live OpenAPI
`GetSkills` and `GetSkill` operations through the selected MCP or CLI connection.
The reading workflow is the same for both; only the invocation syntax differs.
Choose between a complete metadata catalog followed by relevant documents, or the
whole collection with content in one call. Reuse current material already fetched.

## Choose the relevant reference

| If you need to | Read | Then |
|---|---|---|
| Use stored knowledge to answer questions or prepare monitors | [Reading product knowledge](../notehub-project-chat/product-knowledge.md) | Follow the product index, relevant recipes and their meaning dependencies |
| Teach or revise shared knowledge at the user's request | [Project Training](../notehub-project-train/SKILL.md) | Enter the dedicated interview and publication workflow; skill writes belong there |

## Choose how to read

| Reading need | API operation | What to do with the result |
|---|---|---|
| Discover the set and select relevant documents | `GetSkills`, omitting `includeContent` or setting it to `false` | List all names and inspect `kinds`, sizes and links; fetch the index and needed documents with `GetSkill` |
| Read one selected document or follow a product link | `GetSkill` with the catalog's skill name | Read its exact content and follow consequential dependencies in the same set |
| Get a modest set in one call, prepare broad Q&A or acquire a complete training baseline | `GetSkills` with `includeContent=true` | Retain the complete snapshot, then read the relevant documents; Training examines all documents, including unlinked ones |
| Reuse knowledge in a later question or scheduled run | Reuse fetched content; check `GetSkills` metadata when freshness requires it | Keep unchanged documents and fetch new or changed content only as needed |

Metadata-first reading is a normal choice for a focused question or a large collection,
not merely an error fallback. A full-content call is simpler when the set is small or
most documents will be useful. Fetching a complete set does not require putting every
document into model context at once; retain it in private task storage when available.

## List the catalog and inspect kinds

`GET /v1/projects/{projectOrProductUID}/skills` returns the entire metadata catalog:
`project_uid`, `set_revision`, `total`, `skills` and `complete: true`. There is no list
pagination or modified cursor. Each entry includes `name`, `kinds`, `size_bytes`,
`revision`, timestamps and `links`. An empty `skills: []` is a valid untrained set;
an API failure is not an empty catalog.

`kinds` are API metadata independent of Markdown front matter. A document may have
several kinds. Inspect them locally to locate an `index`, product context, field
definitions and recipes; do not invent a server-side kind filter. Use the index and
actual content to decide relevance. Missing or unfamiliar kinds do not make a document
irrelevant, and a tag alone does not establish what its content means.

For subprojects, reads inherit parent-owned skills. Retain the returned owner separately
from the selected data project. Use the project's actual access settings when explaining
who can read the knowledge; shared does not mean private to a particular role.

## Read selected documents or the complete set

`GET .../skills/{skillName}` returns one document's metadata and exact UTF-8 `content`.
Use the catalog's relative `links.self` as an API path on the selected hub, or encode
the **entire** skill name once as one URL path component, including any `/`. For example,
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
freshness boundary, compare names and revisions, and acquire only the needed new or
changed documents. Reconcile removed entries in the private copy and recheck affected
query plans or monitors. Reading or refreshing this knowledge does not publish anything.

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
