# Write product skills during Training

Use this contract only within Project Training's publication workflow. Ordinary Project
Chat, Q&A, monitors and cloud preparation read shared skills; they do not rewrite them.
Use the live OpenAPI `PutSkill` and `DeleteSkill` operations through the selected MCP or
CLI connection. There is no separate upload authentication, v0 `/req` flow, base64
wrapping, timestamp-based filename or manual old-version cleanup.

## Choose the relevant reference

| If you need to | Read | Then |
|---|---|---|
| Read the original baseline or current documents before merging | [Read product skills](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-access/project-skills.md) | Retain exact content, kinds and revisions for every affected document |
| Review meaning, publication scope and the continuing interview | [Training publication](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-train/publishing.md) | Follow the lesson's review, publish, verify and continue cycle |

## Address the Project and skill

For subprojects, writes must target the parent Project that owns the skills. Use its
ProjectUID and make the shared effect clear; queries still use the selected
ProjectUID. Discover actual access settings; do not promise that shared skills are
visible only to a particular role or that a public-readable project is private.

Use the returned `links.self`, or encode the **entire** skill name as one path component,
including any `/`. MCP code mode accepts the original unencoded `skillName`; preserve
header names and quote the exact revision as the live contract requires. When response
headers are unavailable, use the body's revision to construct that condition.

## Write conditionally

| Action | Method and required condition |
|---|---|
| Create | `PUT .../skills/{skillName}`, `If-None-Match: *` |
| Replace | `PUT .../skills/{skillName}`, `If-Match: "EXACT_REVISION"` |
| Delete | `DELETE .../skills/{skillName}`, `If-Match: "EXACT_REVISION"` |

PUT takes JSON `{"content": "exact Markdown", "kinds": ["product", "usage"]}`.
`kinds` are API metadata, independent of Markdown front matter. Omitting them preserves
existing kinds; `[]` clears them. Preserve unfamiliar valid kinds when revising an existing set. Kinds are lowercase strings
matching the live schema; `publish` is reserved (currently at most 32 kinds, 64 characters each).
Do not derive kinds from arbitrary body text or assert that legacy scalar tags are still
required. Content is limited to 1 MiB and the encoded JSON request to 2 MiB plus 64 KiB; verify live limits.

Create/replace returns 201/200 with `outcome` and the resulting `skill`; an identical
write can be unchanged. A 412 means the
precondition failed: fetch the current document, compare it and the draft to the
last verified saved version,
merge nonconflicting changes and resolve substantive conflicts before retrying with its
new revision. A 428 means a required precondition is missing. Never retry unconditionally.
DELETE returns 204 and removes all versions of that name; retaining a local baseline
makes recovery possible but does not make deletion harmless.

Multi-file publication is **not atomic**. Review the semantic change and authorization,
publish dependencies/replacements first, publish the derived index last, and retire old
names only after replacements are read back and inbound references are fixed. On partial
failure, stop dependent deletions and report exactly what reached the service. Read back
names, kinds and exact content after publication; local drafts are not verification.
Retain the resulting verified revisions and content for the next lesson's writes;
the original session snapshot remains an immutable comparison copy, not the write
condition after a replacement. Training publication supplies the meaning-review
and consistency workflow.
