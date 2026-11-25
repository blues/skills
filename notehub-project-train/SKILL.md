---
name: notehub-project-train
description: Train or revise a Notehub project's shared product knowledge through interviews, source review and bounded data investigation. Use when the person explicitly asks to train, resume training, teach product meaning or publish a correction; not for ordinary project Q&A.
---

# Train product knowledge

Learn **what people mean and what the data means**, building on the knowledge
already stored. Product skills connect a customer's real question to the records
that answer it: the physical product and its use, the customer's vocabulary, the
meaning of fields, and the rules for drawing useful conclusions.

**Training is a continuing interview.** Set up once, repeat the lesson loop, and
save useful checkpoints throughout. A successful publication is the beginning of
the next lesson. Keep improving until the person pauses, stops or changes the
task; honor an explicitly limited correction without expanding its scope.

This cookbook organizes the work. Read the relevant references in full, including
their examples and decision rules; this overview does not replace the detailed
interview, field inventory, evidence or publication procedures.

## 1. Prepare the training baseline

Reuse the access, API knowledge, Project selection and complete product knowledge
already established by the entry workflow. Do not repeat that setup or provision
Cloud Computer infrastructure during an interview.

Retain the selected ProjectUID and, when skills are inherited, the parent
ProjectUID that owns them. Reads and drafts can proceed without write permission;
publication must use the owning Project and its actual access settings.

Read or reuse a complete current Skills API snapshot, including every document's
content, `kinds` and revision, even documents nothing links to. Internalize the
whole set. Preserve an untouched baseline beside private working drafts in the
harness's own workspace, outside any supplied source repository. Keep the existing
useful layout and knowledge; restructuring is not a reason to discard detail.

| If you need to | Read | Complete by |
|---|---|---|
| Recover missing access or Project context | [Entry workflow](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub/SKILL.md) | Complete only the missing prerequisites and return with the selected Project |
| Obtain all stored content, kinds and revisions | [Read product skills](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-access/project-skills.md) | Hold a complete baseline, including unlinked documents and the owning ProjectUID |
| Organize findings and tag skills | [Knowledge model](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-train/knowledge-model.md) | Preserve one owner per fact and all applicable kinds; a document may have several kinds |
| Understand what the resulting set must provide | [Output contract](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-train/output-contract.md) | Know how a later agent will read, use and resume the saved knowledge |

**Ready when:** you understand the existing product knowledge, its unresolved
questions and source restrictions, and can distinguish stored content from drafts.
Explain where drafts live and their actual persistence; do not promise automatic
deletion or call a local draft saved to Notehub.

## 2. Open or resume the interview

Read the interview workflow and choose the opening from the person's request and
the saved handoff. Existing skills are a foundation, not evidence of a change.

| Situation | Procedure | Complete by |
|---|---|---|
| No product skills are stored | [Initial training](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-train/interview.md#initial-training) | Make a brief reconnaissance, then hear the product story and how people use its answers |
| Continue improving an existing set | [Continuation](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-train/interview.md#continuation) | Resume the first useful unanswered question without repeating onboarding |
| Correct a rule, review new material or address changed behavior | [Revision](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-train/interview.md#revision) | Start with the specific change and preserve the remaining interview worklist |

**Ready when:** you have a concrete lesson to pursue. The agent chooses useful
questions; the person should not have to design the interview or reconstruct the
previous session.

## 3. Repeat the lesson loop

Read the lesson procedure, then repeat these steps. A lesson can span several
turns. During initial training, learn the product, its use, audiences and their
questions before concentrating on machine evidence. A continuation follows its
saved worklist and whatever the person now wants to teach.

1. **Choose a consequential gap.** Use the saved questions and what was just
   learned. Prefer a gap that changes a useful answer; skip settled questions and
   park those waiting on unavailable evidence.
2. **Investigate and ask.** Read a bounded piece of relevant evidence and ask one
   focused question, or at most three closely related questions in a turn. Keep
   observations separate from guesses about meaning. Before requesting restricted
   material, explain actual source handling and what the artifact would settle.
3. **Improve the drafts immediately.** Record the finding at its owning rule with
   provenance, scope and applicable kinds. Update affected field entries, recipes
   and their dependencies. Preserve unresolved meanings with stable question IDs,
   their impact, who or what could settle them, and what an agent does meanwhile.
4. **Read back consequential meanings.** Confirm the actual sentences and what an
   answer based on them would claim. Promote `heard` to `stated` only for confirmed
   statements. Preserve version differences and established knowledge when new
   evidence conflicts; do not silently replace them with guesses.
5. **Save a lesson checkpoint.** Update the index and question worklist, review the
   semantic change, publish within existing authorization, and verify the saved
   result. The checkpoint includes the next actionable training question. Use the
   publication procedure below; do not leave all saving until the person stops.
6. **Continue with the next useful question.** A readback, summary or successful
   write returns here. If confirmation is pending, keep that ask unambiguous and
   resume after the answer. If the worklist is empty, examine thin knowledge,
   untested recipe boundaries or new evidence without manufacturing gaps.

| During the loop, when you need to | Procedure | Complete by |
|---|---|---|
| Choose questions and conduct a lesson | [Lesson loop and topics](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-train/lessons.md) | Deepen a useful meaning and leave a concrete next ask |
| Examine schemas, samples and Project configuration | [Bounded reconnaissance](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-train/reconnaissance.md) | Gather evidence that helps the interview without turning it into a bulk data investigation |
| Read firmware, documents or other supplied material | [Source handling](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-train/sources.md) | Derive useful meanings while honoring the source's handling and sharing restrictions |
| Map Notefiles and fields to product meaning | [Field inventory](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-train/inventory.md) | Account for every relevant name and field, including validity, significance, action and remaining gaps |
| Interpret native Notehub or Notecard data | [Platform meanings](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-platform/SKILL.md) | Reuse shared definitions and teach the product-specific interpretation |
| Record confidence, provenance or a contradiction | [Evidence](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-train/evidence.md) | Preserve confirmed claims and resolve or explicitly scope conflicts |
| Teach how to answer useful questions | [Answer shapes and recipes](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-train/recipes.md) | Connect customer wording to fields, calculations, presentation and limits; test what can be tested |
| Save a checkpoint or correction | [Review, publish and verify](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-train/publishing.md) | Review meaning and kinds, use conditional revisions, publish dependencies before the index, and verify exact readback |

Training observes Device behavior; it does not send test commands or configure
monitors. API mechanics and private execution state belong in their existing
workflow and runtime locations, not in the product's shared knowledge.

## 4. Pause, stop or switch Projects

When the person says stop, **stop asking and investigating**. Record answers and
confirmations already given, finish already authorized publication unless they
also cancel saving, and preserve the handoff from known gaps. Stopping supplies
no new publication permission. If saving is unavailable or unauthorized, identify
the unpublished draft and its actual location rather than implying Notehub has it.

Briefly report what is solid, thin or unreached; which answers work or remain
conditional; what was saved versus left in draft; and the first question for next
time. Never imply that the product is fully trained because files exist.

| If the person | Continue here | Then |
|---|---|---|
| Pauses, finishes, says stop or wants another Project | [Project selection — Step 4](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub/SKILL.md#4-choose-the-project) | Leave the handoff; wait after a stop, or select the requested Project and read its skills before continuing |

Existing caches and monitors keep their assigned ProjectUIDs. A later training
session resumes the saved question rather than restarting the product interview.
