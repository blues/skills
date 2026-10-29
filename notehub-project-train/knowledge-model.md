# Product knowledge and document organization

Use this to organize shared product knowledge and decide where a finding belongs.
The two areas—what people mean and what their data means—can span the existing useful
file set. See the reader-facing output contract for its entrance.

## Choose the relevant reference

Read the rows that apply to the current step; reuse context already established.

| If you need to | Read | Then |
|---|---|---|
| Make the stored set easy for a later agent to enter | [Output contract](output-contract.md) | Provide an index, question routes and honest answerability |
| Define how questions map to evidence and presentation | [Answer shapes and recipes](recipes.md) | Connect the customer wording to meaning, calculation and output |
| Document every Notefile and field | [Field inventory](inventory.md) | Cover all four rungs and preserve the owning definition |
| Orient the trainer with bounded project evidence | [Reconnaissance](reconnaissance.md) | Keep API mechanics in working context rather than product knowledge |
| Interpret native fields | [Platform](../notehub-platform/SKILL.md) | Reuse shared definitions and record product-specific applicability |
| Derive knowledge from supplied documents | [Source handling](sources.md) | Respect sharing restrictions and retain useful provenance |

## What you are doing, and why

Someone is going to ask an agent a question about their product. Not about their data - about their
**product**. *Which units need a service visit before they fail? Is this one being worked harder
than the others? Do we have to send a technician, or can it wait?* They will not know that the
answer lives in a Notefile called `motor.qo`, in a field called `i`, an array of three phase
currents in hundredths of an amp where exactly `0` means that phase's sensor is unplugged rather
than the machine being idle. **Closing that gap is the whole job.** You are running an
interview and writing down **two complementary knowledge areas**; everything below serves one of them, and at
any moment you should know which.

**Knowledge area 1 - the product, what people mean, and how an answer should look.** What the thing *is*:
the physical product, what it does, who has one and where it sits; the mission, the words these
people use, their thresholds and their worries - so an agent can take a question posed in the
product's own language and work out which data answers it, or that nothing here can. Then *how the
answer should look:* map, line over time, ranking, distribution, a sentence with a number in it, or
a refusal - in the units the audience thinks in, at a precision the evidence supports, against the
baseline that makes the number mean something. Answer presentation is that second half;
recipes are where meaning and presentation meet.

**Knowledge area 2 - where the data is, and what it means.** Every Notefile by name, and every field in it:
what the Notefile is for, which direction it travels, who writes it, how often and under what
conditions; then, field by field, the name as it appears on the wire, where it sits in the body, its
type, units and encoding, what a missing or sentinel value means, and how it can be wrong. A schema
shows names and types, never meaning, so without this no query can be aimed at the right data.
The inventory is the most reusable thing the interview makes.

**Knowledge area 1 comes first, and it is not a matter of taste.** The two are not equally recoverable. The
inventory is half machine-derivable and will keep: the schemas, a sample of events and the firmware
are still there next month, and a later run can rebuild much of it. Nothing rebuilds the product.
*"The bottom probe silts up, so never decide on it"*, *"never tell a customer it is safe"*, *"two of
those are demo rigs"* - none of that is written down anywhere, and the person who has it is busy and
may leave mid-session. **The scarce resource in this interview is the human, not the data**, so
spend the first and best minutes on what only they can give you. It also makes knowledge area 2 cheaper:
reading a Notefile called `aer.qo` cold, you are guessing; knowing the product has aeration fans,
you are confirming.

### What you do **not** have to teach

**The agent that will answer people's questions is smart, and it has its own authorized access.**
Given the product in words and the documentation URLs - the pointers in the output contract - it works out its own queries, and needs no pagination, cursor
semantics, filter encoding, SQL dialect, tool signature or call sequence from you. Every line spent
on those is taken from the meaning only this interview can capture. Bounded reconnaissance is the working knowledge *you* need today; it guides the trainer, rather than becoming product knowledge.

**Shared platform meaning already has a home.** When you reach native Notecard or
Notehub fields, read notehub-platform and its relevant
references. It explains how firmware, health, sessions, signal, usage, tracking,
location and power map to API objects, Event metadata and system Notefile bodies.
The answering agent uses it too. Do not spend the first interview minutes re-deriving
these shared definitions; preserve Knowledge area 1's priority. Training contributes the
product's physical meaning, applicable hardware/configuration, custom fields,
thresholds and verified deviations. An API type or familiar field name is not enough
to establish a native field's meaning.

### Two tests, and the shapes that defeat an agent

Two tests apply to every sentence you write. **Would this help an agent answer a question asked by
someone who has never seen the data?** If it only helps a reader who already understands the schema,
it is documentation. And **is it shape rather than data?** "Three phase currents in hundredths of
an amp, and a sustained reading above the nameplate rating is an overload" is shape; "on Tuesday
unit 4021 drew 9.1 amps" is neither. Hunt these shapes deliberately, and write what each means in the asker's words: **names
that lie**; **opaque names**; **arrays where position is meaning**; **integers that are really
enums**; **implied units**; **templated encodings**, where an omitted field arrives as a zero so
absence and a real zero are alike on the wire; **the same quantity under two names** in two
Notefiles, where a join on one name drops half the fleet; and **absence as the signal**, where the
meaningful event is the note that never arrived. Each is a hole in knowledge area 2; the four rungs and the
`sentinels` and `derivation` kinds are how you close them.

## Classify knowledge and choose its owner

Classify each document with the Skills API's `kinds` array. Front matter may repeat
those kinds for human readability, but it does not control API metadata. Keep the file
set small and the kinds rich. The Area column below identifies a knowledge area, not
a separate required file.

These knowledge categories classify knowledge stored in a project; they are not the
discovery mechanism for shared skills on {{.DOMAIN}}. No new kind is needed merely
because a finding concerns platform data. Use `notefiles`, `derivation`, `sentinels`,
`scenarios`, `config` and `recipe` for the project's application of that data. `usage`
continues to mean how the customer's machine is worked, not cellular byte accounting.

| Kind | Area | What it holds |
|---|---|---|
| `product` | 1 | The physical thing, its variants, what is inside it, who operates it |
| `usage` | 1 | **How the machine is worked**: duty cycle and what a normal week looks like against an abnormal one; operating modes and what switches between them; season and weather; who touches it day to day and what they do to it, including the shortcuts they take; what overload, misuse or running-through-a-fault looks like; and how installations differ - a basement, a rooftop, a metal enclosure, a moving vehicle, a rural site - because that is what decides whether silence is a fault |
| `mission` | 1 | What the organization is trying to achieve; the measure of success that is not a sensor reading |
| `audiences` | 1 | Who asks, in what words and units, and what they may be told - always with a row for anyone not listed, and a relay rule |
| `vocabulary` | 1 | The customer's words mapped to fields, both directions, tagged by who uses them |
| `presentation` | 1 | How an answer should arrive - a sentence, a decision, a map, a series, a ranking or a refusal - with units, precision and baseline (answer presentation) |
| `constraints` | 1 | What an answer must never claim, and who must be escalated to |
| `notefiles` | 2 | Every Notefile by name: role, direction, writer, cadence, its field-by-field account (inventory), and which field carries which concept |
| `derivation` | 2 | Units, formulas, and what is computed rather than measured |
| `sentinels` | 2 | What absence, zero and placeholders mean - absent is not zero - and whether a sentinel can ever be a real reading |
| `scenarios` | 2 | What situation produces what data: which conditions make a note appear, at what cadence, carrying which fields, and what its absence means. This is the **data** side of `usage` - the machine being worked is knowledge area 1, the notes that result are this |
| `population` | 2 | Which devices are real, which are test rigs, gateways, or long dead |
| `config` | 2 | Project identity, fleets, environment variables, smart fleet rules, routes: what each controls, which are dead letters, and this project's rough scale (devices, notes per day, history available), so a plan can be sized before it is run |
| `recipe` | both | A question in the customer's words, the semantic bridge to the data that answers it, and how to prove an answer right (recipes) |
| `anomalies` | both | Everything that does not fit - where what you were trained and what you observed disagree - each written as a question for the person who can explain it |
| `questions` | both | Prioritized unanswered questions for this and later sessions: a live worklist that shrinks as they answer it |
| `index` | both | Start here: what this project is in a paragraph, what the other skills cover, the pointers a reader needs to query, and in one line which questions this set can answer and which it cannot |

`questions` deserves discipline, because a later run reads it to know what to ask: every claim
marked `assumed`, every unfinished lesson's concrete knowledge gap, and every question the
person was not there to answer belongs in it. Give each a stable ID (`Q7`, not a list
position), the actual question to ask, why it matters, which rule it affects, which role
or evidence could settle it, and what an agent does meanwhile. Order by value to an answer,
mark prerequisites or unavailable evidence, and identify the first question ready for the
next session. Remove answered items and do not invent gaps merely to fill the worklist.
This queue covers product knowledge and unfinished lessons as well as anomalies.

**The file set, and its size.** **Keep an existing layout**: read what is there and extend it. For a
project that holds nothing yet, start with six files - `index.md` (`index`), `notefiles.md`
(`notefiles,derivation,sentinels`), `product.md`
(`product,usage,mission,audiences,vocabulary,presentation,constraints`), `operations.md`
(`population,config,scenarios`), `recipes.md` (`recipe`) and `questions.md` (`questions,anomalies`)
- splitting one only when a topic becomes independently useful. **Readers load the index and relevant documents, and training reads the whole set, so it has a budget:** the index fits on a screen, and a set past
roughly a hundred kilobytes has become a research log. Consolidate duplicates, keep the fact and
drop the narrative that produced it, never store transcripts or raw dumps, and record in `index.md` the source-sharing restrictions relevant to the published knowledge (source handling).

### Where a finding goes

Some findings fit several kinds at once. Route by asking **which question would be answered wrongly
if nobody knew this?**, and file it where that question will be read, not where you found it. Then:

1. A **disagreement between two sources**: what is true goes at the owning rule, bounded by the
   version, cohort and period it holds for; the disagreement goes in `anomalies`.
2. **Whether a device's data can be trusted**: `population`.
3. **What a value means**: beside its field in `notefiles`, or in `sentinels` or `derivation` -
   not buried in unrelated narrative, because the agent that needs it is interpreting a number.
4. **This project's configuration** rather than its data: `config`, including anything configured
   but inert - a fleet rule that never fires is a dead letter there, and `population` says in one
   line that membership carries no signal.
5. **The dimensions this product does not have** - a work order, a customer or contract, a site,
   a shift, a cost, a technician's visit: for each one people ask about and this project cannot
   supply, say so in `product` or `config`, and name where it does live if you know.

**A category is not a population rule.** When someone says "two of them are test rigs", "the
trade-show ones", "a couple of customers are on the old units", you have a count, not knowledge. Ask
for the identifiers, or for a rule a query can actually evaluate - a name pattern, a fleet, an
environment variable - with an owner and a way to refresh it. If you get neither, it goes in
`questions`, it is named beside every recipe it affects, and it **never** goes in `population`. Do
not invent identifiers, or quietly report a fleet-wide denominator in its place.

**One file owns each fact.** Write the full statement once, under a stable heading, in the file whose
reader needs it most; every other file carries a one-line pointer to that heading, never a
restatement, and **a pointer carries no origin tag of its own** - write `see notefiles.md,
Conversions` and nothing else, because the origin stays with the sentence in the owning file.
