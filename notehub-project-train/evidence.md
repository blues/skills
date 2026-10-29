# Evidence, anomalies and corrections

Use this when recording a finding, resolving an anomaly or revising a stored claim.
Preserve origin and applicability separately from confidence; then carry corrections
through the owning rule and every dependent recipe before publication.

## Choose the relevant reference

Read the rows that apply to the current step; reuse context already established.

| If you need to | Read | Then |
|---|---|---|
| Record source identity, citations or sharing restrictions | [Source handling](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-train/sources.md) | Use appropriate role/provenance descriptions without copying restricted material |
| Promote confirmed meaning or publish reconciled claims | [Publication](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-train/publishing.md) | Confirm the actual sentences and review every dependent recipe |

## Record origin, confidence and applicability

Every claim carries a marker saying where it came from, how sure you are, and what it applies to:

```markdown
- An excursion means above 8 degrees for fifteen consecutive minutes. `[stated:trainer, operations 2026-09-17 conf:high]`
- The median interval between notes is 15 minutes across 240 devices. `[observed:events 2026-08-18..2026-09-17 n=31,000 has_more=false conf:high]`
- The 8-degree limit may be contractual rather than physical. `[assumed conf:low]`
```

Five origins, and they are not interchangeable:

- **`stated`** - a human told you, in conversation with you, in this session, **and you read the
  sentence back to them and they agreed to it.** Record **which role** said it - `trainer`, or
  `operator`, or `on-call`, never a name (source handling). A document that quotes a human is
  `documented`; your own paraphrase they never heard is `heard` until they have seen the words.
  `stated` is the one origin a later run may never silently overwrite. A read-back promotes exactly
  the sentences read: a batched "yes, fine" promotes each sentence you enumerated and nothing in the
  prose around them, and one promoted sentence may be cited in many places, so demoting a genuine
  claim to make citations match read-backs is data loss.
- **`documented`** - you read it in firmware, a datasheet, a manual, a website. Record the source
  *as source handling has you describe it* - "the V3 firmware", "the published owner's guide" - with an approved public citation where appropriate and its version or review date; an assertion stored by a previous run
  is a prior assertion, not a fresh verification.
- **`observed`** - you inferred it from data. Record the window, the filters, how many events, and
  whether more were available; a large row count is not representativeness.
- **`assumed`** - you guessed. Say so, keep the open question, and never let it reach an answer
  unhedged or become a decision rule.
- **`heard`** - the person said it in this session and has not yet heard it read back. Record the
  role and when. It outranks `assumed`, may be corrected, and becomes `stated` only by read-back - so a
  rule that was not read back is `heard`, and the recipe says so.

Use short source IDs (`F1`, `E3`) from the front matter rather than a long citation on every line,
and confidence is a separate axis from origin: a person can be confidently wrong.
Evaluate mechanism, sampling bias and coverage, not just the number of observations.

**What an update run may do with an existing claim: two axes, and collapsing them is how good
knowledge gets lost.** **Origin governs retraction**: only a human may retract what a human stated,
a `stated` sentence is never silently rewritten or deleted, and everything else you may correct
yourself, saying what changed and why. **Strength and scope govern replacement**: newer is not
stronger, and a weaker or narrower claim never silently replaces a stronger or broader one - record
it beside the existing claim and surface the conflict. contradictions carries both out, and it runs before
you write anything on an update run.

**Applicability travels with the claim.** A firmware checkout describes *that version's* behaviour,
not necessarily the population that produced the data in front of you. Carry version, cohort and
effective period on every consequential claim, and keep unknown-version data unknown. Never write
the word *current* in a rule: where a rule changed for everyone on a date, key it on that date
("notes captured before 2026-03-01 carry `t` in whole degrees; from that date, tenths"), but do not
manufacture a boundary for a rollout that went device by device, or transfer one generation's
formulas onto another. **Verify a load-bearing claim against a primary source before you write
it** - the firmware, the API's specification, a read-only query; a research summary is a pointer,
not a source. And where an effect is smaller than the precision the value is stored at, say so
beside the finding.

Front matter can carry source and review metadata so staleness can be judged quickly.
The API's `kinds` array is authoritative for classification; a Markdown `kinds` list is
optional descriptive metadata, not an upload-tag protocol:

```yaml
---
kinds: [notefiles, derivation, sentinels]
description: Every Notefile and what each field means
updated: 2026-09-17T14:02:00Z
sources:
  - id: F1
    description: host firmware v2
    sharing: restricted
    reviewed: 2026-09-17
  - id: H1
    description: operator interview
    reviewed: 2026-09-17
---
```

Describe restricted sources without private paths or identifying details. Approved public
sources may retain useful citations. Preserve source version, scope and review time;
paraphrasing does not make a sensitive formula or threshold safe to share.

## Turn discrepancies into useful questions

When firmware, existing knowledge and observed data disagree, first make the cheap
checks: generation, cohort, clock, filter and coverage. A real mismatch often has a short
human explanation—“those units were reflashed and never got the new template”—that is
more useful than a long speculative investigation. Ask the person who knows the history
when that can settle it. Retain resolved anomalies and their mechanisms in a compact
record; `questions` is the live worklist and shrinks as answers arrive.

If implementation evidence and general documentation disagree, record which producer,
field and firmware/configuration the evidence actually establishes. Keep a confirmed
exception at its owning field with that scope; retain an unresolved conflict in
`anomalies` and `questions`. Capture the useful rule and provenance without copying raw
project data or internal source into the published skills.

1. You notice something that does not fit. **Write the anomaly** with its evidence, and **open a
   question** with a stable ID that names it. Each points at the other.
2. Put it to the person while they are in the room, the way a colleague would: *"the firmware sends
   these seven fields and I can't find any of them in a year of events - did something change?"*
   Order by how much each would change an answer.
3. **Remove the question** - a worklist that keeps finished items stops being read - and **update
   the anomaly with their explanation**, against their role and the date, marked `heard` until the publication read-back promotes it like any other claim; an answer given in passing is not a read-back.
   The next run can find the resolution rather than repeating the question. Reopen it only
   when new evidence changes its applicability. Continue with the next useful lesson if
   scope and time remain; do not invent a question merely to prolong the interview.
4. If the explanation turned the anomaly into a mechanism, **also put it where a reader will use
   it** - beside the field in the inventory, or in `scenarios` or `population` - and leave a
   pointer.

Three rules keep the pair honest. **Before you park a question for somebody who is not here, ask its
weaker form to the person who is**: *"what causes it"* is for the engineer, but *"what do you do when
you see it"* and *"how often does it bite you"* are for whoever is in the room, and that is usually
what a skill needs - park only the residue, then identify the role or source that could settle the remaining gap. **Never remove
a question because you worked around it**: it leaves the list when it is *answered*, and if you
proceeded on an assumption the question stays and the assumption is marked `assumed`. And **an unresolved anomaly needs a question or an explicit reason no further action is warranted.**

## Reconcile corrections without losing knowledge

When new material arrives for a project already trained, check it against what is stored before
writing anything. For each claim that disagrees with one already there, ask first whether the two
describe different versions, cohorts, clocks or periods - if so, bound both and keep both.
Otherwise apply evidence: you may correct a stored `observed`, `heard` or `assumed` claim, or a
`documented` one your newer source supersedes for the same population, when your evidence is
stronger - saying that you did; if it is weaker or narrower it goes beside the stored claim, not
over it. If the stored claim is `stated`, **stop and ask**: show the person their own words, what
now contradicts them, and let them decide. If they are not available, append the contradiction
beneath the claim rather than resolving it, the original untouched and the contradiction carrying
its own origin. Contradictions are findings, not errors: a claim from six months ago often means a
fleet has been reconfigured since.

When they decide, there is **one retraction form**, and you do not invent another. In the owning
file only: strike the old sentence through, keep its origin, follow it with
`retracted|bounded|superseded by <who> <date>` and the replacement with its own origin. Every other
file points at the owner, `anomalies` gets an entry, and the open question closes carrying the
answer into the owning rule. History needed to read older events is never erased.

**A retraction that removes a decision rule is not finished until the rule is replaced.** Striking
out "use the mean of the three probes" deletes the reduction every recipe depending on it was
using, and "the replacement must be re-chosen" is a note to yourself, not a rule a reader can
follow. Either write the replacement in the same pass, or **mark every dependent recipe blocked on
the named open question** and say so in its status line. A set that retracts a rule and leaves its
dependents reading as verified is worse than one that never retracted it, because every reader now
invents a different replacement - which is the failure the retraction existed to prevent.

```markdown
- ~~Devices report every 15 minutes.~~ `[stated:trainer 2026-09-17]` superseded 2026-09-18:
  before 2026-03-01 they are 15 minutes apart, from that date 30. `[stated:trainer 2026-09-18]`
```
