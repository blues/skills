# Run the next useful training lesson

Choose the lesson that most changes a useful answer, using the saved handoff and what
has just been learned. See opening/resumption only when establishing the
run; use the inventory, evidence and recipe references for the lesson's detailed work.

## Choose the relevant reference

Read the rows that apply to the current step; reuse context already established.

| If you need to | Read | Then |
|---|---|---|
| Establish or resume the interview | [Opening and resumption](interview.md) | Choose the appropriate opening only when establishing the run |
| Capture Notefiles, fields and configuration meaning | [Field inventory](inventory.md) | Account for meaning, validity, significance and action |
| Define useful questions or their answer presentation | [Answer shapes and recipes](recipes.md) | Ground current-state and period answers in their distinct semantics |
| Decide where a finding belongs or mark an inapplicable dimension | [Knowledge organization](knowledge-model.md) | Preserve one owner per fact and explicit applicability |
| Request or interpret source documents | [Source handling](sources.md) | Confirm handling and sharing constraints before requesting restricted material |
| Investigate an anomaly or reconcile contradictory claims | [Evidence and contradictions](evidence.md) | Preserve provenance and scope; surface unresolved conflicts |
| Read back or publish a completed lesson | [Publication](publishing.md) | Use existing authorization and verify the saved result |

## Keep improving through the interview loop

**Run this loop until the person pauses, stops or redirects it, and make each pass deepen
the skills.** Choose the next useful lesson from what was just learned and the saved
worklist; do not hand the person a menu and make them plan the interview. A completed
lesson is worth revisiting for a new exception, audience, boundary case or consequence.
The initial coverage baseline establishes usefulness, not a final state called “trained.”

**Every active interview turn leaves a clear next ask.** A clarification, read-back,
publication or summary returns to the loop; it does not silently end the session.

- When the person answers, record it, explain briefly which meaning or answer changed,
  and continue in the same turn with the next focused question or relevant finding.
- At each lesson boundary, update the owning claims, affected recipes and handoff;
  read back the consequential sentences and their answer implications, then publish
  within existing authorization or ask for the missing approval. Continue after the
  response instead of ending with “let me know if you need anything else.”
- When meaning confirmation or publication approval is the pending ask, make it
  unambiguous. Keep the next lesson queued and resume it as soon as that answer arrives;
  a bare “yes” must not be interpreted as answering several different questions.
- If the worklist appears exhausted, inspect the thinnest knowledge area, an untested
  recipe boundary, an unexplained exception or a real question the skills cannot yet
  answer well. Propose a specific next improvement; don't repeat a settled question or
  invent a fact or a customer request to manufacture progress.
- A missing source blocks its dependent question, not the interview. Retain the gap and
  move to another useful question the person can answer. If all useful work is blocked,
  explain the dependency and save the next concrete step for resumption.

Respect an explicit limit such as “correct this one rule.” Otherwise, finishing a file,
publishing a batch or reaching a useful checkpoint is a reason to choose the next lesson.

Meaning confirmation and permission to publish are separate. Do not bundle them with an
unrelated question so that a single “yes” is ambiguous. Use publication
and carry forward the user's existing authorization.

**On an initial run, early passes go to knowledge area 1; a continuation follows the saved worklist.** A to D are
the things only this person can tell you - what the product is, how it is worked, who asks about it
and in what words. E to G are the machine-readable sources, and they will still be there next month.
Work down the list rather than jumping to the firmware: until you can say what the product is and
how it is used, you are building an inventory nobody can aim. Spend the person's early attention there and read the machine-readable
sources around it, because those keep and the person does not.

Each pass: **choose the uncertainty that most changes an answer somebody actually wants**; **read
one bounded piece of evidence**, reusing what has been supplied, and inspect something independent
while they go and find a document; **say what you learned in their language**, separating
established behaviour from hypothesis; **ask one focused question, or a small group of related
ones** - *"does this gap mean a delayed upload or a missed measurement?"* beats *"tell me everything
about cadence"*; **write it down now**, with its scope and evidence, updating every recipe the
answer changes; then **read back and publish within the agreed scope** (publication), and propose the next useful lesson. Never
re-ask for anything already in this conversation, never make a non-engineer answer a firmware
question you could answer by reading code, and build up rather than handing over homework.

**One turn is not the whole lesson.** The lesson lists below are a coverage checklist across the
session, not a questionnaire to deliver in a paragraph: **at most three closely related questions
per turn**, and a turn carrying eight is homework however politely it is worded. **Count the asks,
not the numbers you put in front of them** - "are all 214 of those real units, what do your
technicians call this thing, and what does the fleet screen show them?" is three questions wearing
one bullet, and the person feels all three. **Count them before you send:** each question mark, and
each further thing you need answered inside one, is an ask - the two opening questions of opening and resumption
count as one each, as written - and a turn over three sends the three that matter most and parks
the rest. Start a lesson from one concrete thing the person just said, go deep on it, and collect the rest over later
passes; park what you did not ask in `questions` so the next pass and the next run can find it.

**The initial coverage floor.** The first baseline requires **knowledge area 2 to have an entry for every Notefile carrying
data anyone might ask about** (inventory); until `product`, `usage`, `vocabulary`, `audiences`,
`presentation`, `constraints`, `population`, `scenarios` and `config` each hold something a stranger
could act on; and until at
least one `recipe` exists - including four things every run misses: **one current-state recipe**
and **one fleet-or-period recipe** (recipes), or a statement of why this project has neither; **the
receiving surface** (lesson E); and **the dimensions this product does not have** (knowledge organization, rule 5).
You will usually not reach that floor; what *is* failure is ending without saying which parts you
did not reach, and leaving thin files looking finished rather than marked `[rungs:1 only]`.
Reaching that floor starts deeper passes: test new question wording, investigate edge
cases, improve explanations and reconcile new evidence. It does not end the interview.

**Keep a resumable handoff throughout the session.** At the end of each lesson, update
the `index` training status with where the interview stopped and link its first question
for next time to the matching ID in `questions`. Keep the remaining questions concrete and
prioritized, including topics not yet reached; reconcile them against answers already
recorded. Include this handoff in the normal review and push, so resumption does not depend
on the last chat message or on getting one more turn after the person says stop.

**How a session pauses.** The person controls when to pause, stop or change direction;
do not end merely because a lesson or baseline is complete. Honor an explicitly limited
task. When they say
stop, **stop asking**: open no new investigation and write no new claims. Recording answers and
confirmations they have *already* given is not new work - promote what they just agreed, write the
closing status and the handoff from gaps already identified, and complete a push they approved.
Only a verified publication is available to other agents from Notehub; draft persistence depends on the harness. If the latest handoff remains an
unapproved draft, say that it was not saved to Notehub; do not imply a closing chat summary updated
the project. Then say plainly where you got to:

- **what is solid** - and name it, so they can disagree;
- **what is thin** - which fields sit at rung 1, which files would read as finished and are not;
- **what you never reached at all** - by name, and that means the lessons: say which of A to G
  were never run, so "we never asked what gets done to the hardware in the field" appears in the
  closing rather than only as a question ID nobody reads;
- **the one question you would ask first next time.**

Never let a session end with a silent implication that the project is now fully trained.

### A. Product materials
 *(Knowledge area 1.)* Ask for a website, datasheets, manuals,
support articles, case studies - and the internal ones too: requirements documents, design docs,
test plans, usually the only written record of *why* the product does what it does. This is where
the product, the mission, the audiences and above all the **thresholds** come from. Ask for the
dashboard the customer already sees, because that is the shape they expect an answer in (answer presentation),
and record each source's date and what it does *not* establish: a case study settles neither which
generation is deployed nor this customer's limit. Feeds `product`, `mission`, `audiences`,
`constraints`, `vocabulary`, `presentation`.

### B. Audiences
 *(Knowledge area 1.)* The subject is the person, not
the data, and it is the lesson every run skips. Skip any question already answered; ask the rest,
not the script:

- *Who asks you about this - and by what route: a call, a ticket, a screen they watch themselves?
  What do they want back, a number, a yes or no, or a list?*
- *What words do they use for these things?* (Write them down as they say them.)
- *What must an answer never claim?*
- *What here is contractual - what carries a fine, a warranty claim, or a regulator?* Ask to see or
  hear the clause; record the rule in the document's own terms, having confirmed under source handling that
  quoting this particular clause is allowed - a contractual definition the customer already shares
  with a counterparty usually is, an internal one is not - with which role holds it, and which record and
  statistic it is decided on. An answer turning on a clause nobody has read says so in the sentence.
- *Who gets escalated to and what happens out of hours, and what is the worst thing an agent could
  say to a customer about this product?*
- *Who else will ask an agent about this that never reaches you today - an insurer, a regulator or
  auditor, finance, an engineer, the company that owns the machine, the dealer who installed it -
  what may each be given, and who signs off
  before anything leaves the company?* `audiences` always carries a row for anyone not listed
  ("route through <who>; give nothing directly") and a relay rule: when a listed asker says the
  answer is for someone not listed, it is written to the highest-consequence recipient's rules.
- *For each of them, what word will they use for the answer, and what does it mean here?* Write
  these into `vocabulary` as false friends per audience: "resolved" to an auditor is not "cleared"
  on the controller, and a vocabulary that maps the asker's words for the data but never their word
  for the answer is half a vocabulary.
- *And in what form does each one want it?* That is answer presentation, and it belongs in this conversation.

Feeds `audiences`, `constraints`, `vocabulary`, `presentation`. **A recipe written before this
lesson is written for nobody** - correct, pitched at the wrong person, with no idea what it may not
say.

### C. Real customer questions
 *(Knowledge area 1.)* **Ask first how a question even reaches
them** - an alarm at two in the morning, a dispatcher deciding what to put on the van, a technician
on the radio from a rooftop, a customer asking why their machine stopped, somebody's Monday look at
the fleet, or a unit that just went quiet and nobody noticed for a week. Then take **one** real
recent example, in their exact words, verbatim: what was asked, what a useful answer would have
been, and what was done with it. Work that one into a recipe before asking for another - and if the
inventory does not exist yet, write it as far as it goes and mark it conditional on the fields it
will need, rather than inventing them. Then gather the rest over later passes, and ask explicitly
for the families a single example never surfaces - fleet-wide, spatial, trend and forward-looking:
*how many last quarter, is it getting better or worse, which sites are worst, where are they, when
will this one be due*. Run each recipe's clarifications against the
real utterance: if the identity rule cannot resolve something an asker actually said, the rule is
wrong, not the asker. Feeds `recipe`, `presentation`, `audiences`, `vocabulary`.

### D. Operator knowledge
 *(Knowledge area 1.)* You ask the questions, and three of them are always asked. If the
person announces their last minutes, ask the first in one sentence before anything else is written.
First: *who works this machine day to day, and what do they do to it?* Not the service crew - the
operators. Ask what a normal week looks like and what an abnormal one looks like; what the modes are
and what switches between them; whether it is seasonal; what people do when they are in a hurry, and
what the machine looks like when it is being worked too hard, run through a fault, or used for
something it was not meant for. **Every one of those is invisible in a schema and obvious to them**,
and each is the difference between an alarming number and a Tuesday. Ask too **how installations
differ** - a basement, a rooftop, a metal enclosure, a moving vehicle, a rural site with one bar of
signal - because that is what decides whether silence means a fault or a wall. All of it goes in
`usage`, and what the data does about it goes in `scenarios`.

Second: *what gets done to the hardware in the field?* For each thing they name - a unit swapped
out, a part replaced, a sensor recalibrated, a machine re-sited, firmware reflashed - ask which
counter it resets, which trend it breaks, and whether it is logged anywhere. **Any recipe that spans time is
wrong across every one of these events unless it knows about them**, and none is visible in the data
as anything but an unexplained step.

Third: *what did you learn the hard way, what have you had to explain to a customer, and what
changed afterwards?* Each incident becomes a `usage` or `scenarios` entry with its data signature
and a `constraints` entry with its wording. This is the highest-yield question in the document: ask it in every run, whatever else is
skipped, and if it was never asked the closing summary says so by name. Then hunt the gaps, because
anything marked `assumed` is a question waiting to be asked.

### E. Project evidence
 *(Knowledge area 2.)* Read the schemas, a bounded
sample of recent events covering the variants that matter, the fleets and how devices are
distributed among them, the routes, the environment variables at every level, and the earliest event
the interface still returns. "Exploring the live project" says how, and what will mislead you.
Classify streams by behaviour, not name, and write what you find into the inventory inventory.

Two things here are skipped almost every time. **The system Notefiles** - the ones the Notecard and
Notehub maintain themselves, whose names usually begin with an underscore - carry the connectivity,
power and session story behind every "has this site gone dark" question; they go in the inventory by
name like any other, remembering that such a file can carry application data too, and that a
server-generated event does not prove the product measured anything. And **the routes**: before
writing a recipe that detects a condition, ask whether something downstream already detects it and
whether anybody acts on it. For each route write **the receiving surface** - what it shows a person
per note, which statistic of a multi-sample field, how a zero and a missing note render, what words
it uses, whether it alerts and to whom - because people ask about what the dashboard shows, so
those words go in `vocabulary` and its layout in `presentation`. Fleet rules repay the same
attention - a fleet defined by `$exists(body.t) = false or body.t = 0` teaches you more than the
schema will - but a fleet name is a clue, not a definition, and today's membership is not
historical. Feeds `notefiles`, `population`, `config`.

### F. Host firmware
 *(Knowledge area 2.)* Ask for the firmware that drives the Notecard, and
its version - and ask for it the way source handling says to, naming the file you want and what it buys,
because this is the lesson most often lost to a request nobody bothered to make twice. Read every
`note.add`, `note.template`, `note.get`, environment-variable read and Notefile name - but what
you are really after is **when each fires, and under what conditions**: power tiers, transport
choice, modes, thresholds, the state machine. That is where `scenarios`
comes from, and it turns an unexplained gap into "this device is on its winter battery tier,
behaving as designed". Firmware is often the best source for **`derivation`** - the constants
and formulas behind a field - and much of **`sentinels`**; a template declares each field's *wire
type*, which decides what a test on it can detect. Check the version you read against the versions
actually reporting, and where the data is older than the firmware you were handed, use the firmware
to generate candidate explanations rather than discounting it. Read the comments: often the only
written record of what a field means. Learn the **inbound** side too - what the host consumes, its
defaults and precedence, and how to tell requested from delivered from acknowledged from effective.
Training observes mechanisms; it never sends test commands to devices. Feeds `notefiles`,
`scenarios`, `derivation`, `sentinels`, `config`.

### G. Reconcile evidence
 *(Both.)* The mission says what the
product is for, the firmware says how and above all *when* it senses, the data says which of those
situations really occur, and the schema says what Notehub inferred. Where they agree you have
learned a mechanism, and it belongs in `scenarios` and the inventory; only a genuine conflict
belongs in `anomalies`, and a mismatch earns a cheap check of generation, time filter, paging and
cohort first, because a schema entry missing from a small sample is "not observed in this sample",
not "never sent". A cadence you observed should correspond to a tier you read in the firmware. So
once you have both A and C, compare them in a small table - concept, documented mechanism, what was
observed and over what coverage, agree or conflict, effect on a query, next check - and take the
conflicts to the person (anomalies), reporting without resolving anything yourself: fields the
firmware sends that never arrive, fields arriving that no current firmware produces, Notefiles whose
names no longer describe their contents, devices behaving unlike their siblings. Do not spend an
hour on what they can explain in a sentence. Feeds `anomalies`.

### H. Revisit stored knowledge
 *(Both.)* Use this lesson when new material, a reported
change or an explicit review calls for revisiting stored knowledge (contradictions). It is not a
mandatory opening for every continuation. Work from the specific change and relevant
`questions`, highest value first; update linked `anomalies` when applicable. Ask for a
source only when it is needed to settle that issue, making the source handling promise before
requesting it. A previously promised but unavailable source remains queued; it need not
block other interview questions. Preserve stored provenance and scope: a new session alone
does not demote a `documented` claim to `heard`. Do not re-ask a settled question unless
new evidence makes its scope doubtful, and then say what changed first.

**Review the existing set against source handling too.** Preserve legitimate
product knowledge and approved citations. Remove credentials and unneeded private
identifiers from proposed drafts, replacing identifying details with roles or reusable
rules where possible. Do not silently remove a consequential fact or an approved public
source citation. Include privacy corrections in the semantic publication diff.
Report what you found by count and kind, because a set that had them once will have them again
wherever the next run is careless.
