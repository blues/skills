# Review, publish and verify

Keep an immutable baseline of every original name, content, kinds and revision, alongside
working drafts. Record the semantic diff as work proceeds. Preserve existing useful
knowledge and file organization unless restructuring improves retrieval; a shorter file
is not automatically a better skill. Drafts are not saved to Notehub until a verified API
write. Explain the actual persistence of local drafts without promising automatic deletion.

## Choose the relevant reference

Read the rows that apply to the current step; reuse context already established.

| If you need to | Read | Then |
|---|---|---|
| Review citations and sharing restrictions | [Source handling](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-train/sources.md) | Resolve source restrictions before publishing derived knowledge |
| Write, replace or delete stored skills | [Training skill writes](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-train/skill-writes.md) | Use conditional revisions, dependencies first, index last and exact readback |

## Confirm meaning and consequence

**Use this rhythm at the end of every lesson: update drafts and handoff, read back,
review, publish within scope, verify, then return to the interview.** Do not accumulate
the whole interview for one final write when the person is already leaving. If there
is no new supported finding to publish, keep the unresolved question visible and
continue to the next useful step without inventing a change.

At each lesson boundary, read back the two or three load-bearing sentences: those
that change what somebody would conclude or do. Check clarity yourself; ask the expert
whether they mean what you wrote. Show the answer consequence too: a rule may be true
as phrased yet support a misleading conclusion. For example, detecting a reporting gap
does not establish how long a machine was down.

Keep a small record of the enumerated sentences confirmed and promoted to `stated`.
The two lists must match. A general acknowledgment does not promote adjacent unreviewed
claims; an actual confirmed sentence must not be left as `heard`. Preserve provenance
when it is cited elsewhere. Confirmation of a fact is not by itself authorization to
publish it, and authorization to publish does not confirm every assumed fact.

Ask whether the sentences express the person's meaning; do not delegate your clarity
check by asking whether a stranger would understand them. Incorporate corrections in
their vocabulary. Make the confirmed sentences and their promotion to `stated` an
explicit check now rather than deferring it to a later draft review.

## Prepare a reviewable change

Read recipes alongside their owning rules. Check units, thresholds, per-note reductions,
windows, precision, population, vocabulary and constraints. Do not duplicate an owning
rule in a recipe; link it. Every unresolved load-bearing dependency makes the recipe
conditional or blocked, including a rule retracted without a replacement. Preserve
version/cohort history needed to interpret older events.

Read each recipe against the audience, constraints and presentation rules as well as
its field definitions. A mathematically valid result can still make a forbidden claim,
use the wrong recipient's vocabulary or suggest an action the product rules exclude.
No high-confidence conclusion may silently depend on an unanswered question.

Review source handling, citations and sharing restrictions. Explain each
new/replaced/retired file and what changed in meaning, not just a list of filenames.
Include shared parent scope where applicable. Use the user's existing authorization for
this concrete scope. If it does not cover publication or a newly consequential change,
ask once after the diff is ready. Don't impose a separate confirmation for every write
when the user has authorized publishing the defined work. Silence is never approval.

Draft `index.md` last. Derive its answerability list from actual recipe statuses, including
the same unresolved question IDs. A conditional recipe's title must not become an
unconditional promise in the index. Include training status, the last useful lesson,
solid/thin/unreached areas, and the first actionable question for next time. Reconcile the
question worklist against answers already given; don't fabricate unfinished work.

## Publish with revisions

Use the standard Skills API:
conditional creates, revision-checked replacements/deletions, exact UTF-8 content and
explicit intended `kinds`. Re-read and compare against the preserved baseline when
concurrent edits are possible; the server's precondition closes the race after the read.
A conflict is a merge task, not a reason for an unconditional overwrite.

Publish dependencies and replacements first, read them back, then publish the derived
index. Retire old names only when replacements are verified and no references depend on
them. Deleting a local draft does not delete a stored skill. Deletion removes all versions
of a name; keep the baseline and honor its authorized scope.

There is no multi-file transaction. On partial failure, stop dependent removals, inspect
the server, and report exactly which changes landed. Recovery is a new revision-aware
write, not a blind replay of the whole batch. Verify names, exact contents and kinds from
Notehub, including unexpected/missing documents. Report what is published and what is
still a draft, and where the resumable handoff lives. Complete already authorized writes
when the user ends the interview; do not open a new investigation or ask new lesson
questions after a request to stop.

Keep the person informed at each stage: what changed in the drafts, what the proposed
publication changes in meaning, and what the service actually saved. If only part was
authorized, publish only that part. Explain who can read the stored skills using the
project's actual access settings; do not call a local draft “saved to the project.”

Unless the person has paused or ended the task, return to the next lesson after the
readback. A successful write is an improvement checkpoint, not completion of training.

| If you need to | Read | Then |
|---|---|---|
| Continue after review and publication | [Lesson loop](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-train/lessons.md) | Ask the next useful question and build on what was just learned |
