---
name: notehub-project-train
description: Train or revise a Notehub project's shared product knowledge through interviews, source review and bounded data investigation. Use when the person explicitly asks to train, resume training, teach product meaning or publish a correction; not for ordinary project Q&A.
---

# Train product knowledge

Learn **what people mean and what the data means**, preserving the hard-won knowledge
already stored. Product skills bridge a customer's real question to the records that
answer it. API mechanics belong in these shared workflow skills, not in each product's
knowledge. Follow [working style](../notehub/working-style.md).

Training is a continuing interview that improves the shared skills on every pass.
Choose a useful gap, investigate, ask, draft, review and save within the authorized
scope, then continue to the next question. A completed lesson or initial baseline is a
checkpoint, not the end of training. Keep going until the person pauses, stops or changes
the task; preserve a concrete handoff so the next session can resume. Honor an explicitly
limited correction without expanding it into an unsolicited interview.

## Choose the current training step

Read the references for the current step in full, including their decision rules and
examples; a summary is not a substitute for the interview procedures. Reuse context
already established.

| If you need to | Read | Then |
|---|---|---|
| Establish a missing connection or execution environment | [Access](../notehub-access/SKILL.md) | Retain the verified connection and environment capabilities |
| Select or verify the project and its skill owner | [Project Selection](../notehub-access/projects.md) | Retain the data project, inherited skill owner and relevant permissions |
| Read the complete immutable training baseline | [Read product skills](../notehub-access/project-skills.md) | Save all content, kinds and revisions, including unlinked files |
| Start, resume or revise an interview | [Interview workflow](interview.md) | Choose the opening that matches the request and saved handoff |
| Orient an initial interview with a small data scan | [Bounded reconnaissance](reconnaissance.md) | Return with observations and hypotheses clearly separated |
| Choose the owner and organization of a finding | [Knowledge model](knowledge-model.md) | Preserve useful layout and one owner per fact |
| Record evidence or reconcile conflicting claims | [Evidence and contradictions](evidence.md) | Retain provenance, confidence and applicability |
| Choose and conduct the next useful lesson | [Lesson loop](lessons.md) | Ask focused questions and update the resumable handoff |
| Document fields, configuration and validity | [Field inventory](inventory.md) | Cover meaning, validity, significance and action |
| Turn real questions into grounded answers | [Answer shapes and recipes](recipes.md) | Record dependencies, reductions, coverage and limits |
| Interpret native Notehub or Notecard data | [Platform meanings](../notehub-platform/SKILL.md) | Reuse shared definitions and record only product-specific interpretation |
| Request or use restricted documents | [Source handling](sources.md) | Explain actual handling before requesting the material |
| Review and publish a lesson or correction | [Publication workflow](publishing.md) | Apply conditional revisions and verify exact readback |
| Close a session or check the resulting skill set | [Output contract](output-contract.md) | Report solid, thin and unreached knowledge and the next actionable question |
| Check that another agent can use the published knowledge | [Reader workflow](../notehub-project-chat/product-knowledge.md) | Follow its questions and dependencies without the training conversation |

## Carry out the training workflow

1. Reuse access and project selection. Verify the data
   project and the parent owning its skills. Read the full Skills API snapshot,
   including unlinked files, and preserve content, kinds and revisions as an immutable
   baseline beside private working drafts. Read access supports drafting even if the
   person lacks publication permission.
2. Read the interview workflow. Choose initial training,
   continuation or revision from the request and saved handoff. A continuation resumes
   its next useful unresolved question; it does not repeat generic onboarding. An
   initial run starts with the machine, how it is used and who needs which answers,
   informed by a brief bounded reconnaissance.
3. Organize findings with the knowledge model. Keep
   existing useful layout and one owner per fact. Record evidence and contradictions
   with source, confidence and applicability. `heard` becomes `stated` only after
   the person confirms the actual sentence and consequence. Newer evidence is not
   automatically stronger; preserve confirmed claims and historical scope.
4. Follow the lesson loop: product/use/audiences first, then field inventory,
   answer shapes and recipes, configuration and reconciliation.
   Ask no more than three related questions at a time. Preserve unresolved questions
   with stable IDs, impact, settling role/evidence, interim behavior and prerequisites.
   Use platform meanings for native data rather than
   copying its dictionary into every project.
5. Before requesting restricted documents, explain source handling.
   Draft meaningful rules without copying private sources. Don't invent confidentiality
   or draft-retention guarantees. Treat unconfirmed thresholds, formulas and population
   rules as gaps, not convenient defaults.
6. At the end of each lesson, update drafts and the resumable handoff. Follow
   review, publication and verification: semantic diff,
   confirmed meaning, existing authorization, conditional revisions, dependencies first,
   index last, and exact readback. Stop dependent deletions on partial failure.
7. After review and any authorized publication, return to the lesson loop with the next
   useful question. Revisit thin meanings, untested recipe boundaries and newly relevant
   evidence; do not declare the product fully trained. When the person pauses or stops,
   close against the output contract: what is solid,
   thin or unreached, which questions work or remain conditional, what was published,
   what remains a draft and the next actionable question. Stop asking when the person
   stops the interview. Do not imply that a set is fully trained because files exist.

Read only the references needed for the current lesson. The output contract
shows how a later agent enters the saved product knowledge; the
reader workflow is the test of
whether it can find and use what you taught. Training observes device behavior; it does
not send test commands or configure monitors.
