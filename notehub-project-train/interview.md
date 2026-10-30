# Open or resume training

Use this when starting, resuming or revising product training. Preserve
the existing set and choose the opening that matches the request. Then continue with
the relevant lesson, rather than repeating onboarding.

## Choose the relevant reference

Read the rows that apply to the current step; reuse context already established.

| If you need to | Read | Then |
|---|---|---|
| Recover the overall training workflow or missing prerequisites | [Training](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-train/SKILL.md) | Resume at the relevant step without restarting completed setup |
| Acquire the full baseline, including unlinked files | [Read product skills](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-access/project-skills.md) | Retain exact contents, kinds and revisions beside drafts |
| Make the initial, brief data scan | [Bounded reconnaissance](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-train/reconnaissance.md) | Use schemas and a small sample to frame questions, keeping observed shapes separate from assumed meanings |
| Continue after the opening or resume a saved lesson | [Lesson loop](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-train/lessons.md) | Ask the next useful unresolved question |
| Address a reported change, correction or explicit review | [Revisit stored knowledge](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-train/lessons.md#h-revisit-stored-knowledge) | Trace the specific change through owning claims and recipes |
| Request source material or explain how it will be handled | [Source handling](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-train/sources.md) | Explain actual handling before asking for restricted material |
| Confirm meaning or publish the resulting changes | [Publication](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-train/publishing.md) | Separate fact confirmation from authorization and verify exact writes |

## 1. Verify the baseline and private drafts

Reuse the complete baseline and drafts prepared by the Training setup. If missing,
read every skill the Project holds, contents and all, in one request. Make a
**drafts folder** in your harness's own private scratch or temporary space, outside
any supplied source tree or repository somebody might commit. Keep the original
content, kinds and revisions untouched beside working drafts: this is how you know
exactly what changed. Account for every document, including files nothing links to.

Explain where drafts live and what persistence the harness actually provides;
never promise automatic deletion or session-only retention.
Say what changes in meaning as drafts evolve. Before publishing, use the conditional
revision workflow in publication; preserve concurrent edits.
Name the affected drafts and the change in meaning briefly as each lesson updates them,
so the person can distinguish working knowledge from what Notehub already holds.

## 2. Choose the opening

Use the person's request and the saved state. Having read the whole set, review
the `index` training status and the whole `questions` worklist to choose the opening.
Existing skills mean there is something to build on; they do not mean the product changed.

| Situation | Follow | Then |
|---|---|---|
| Existing skills, with no reported change | [Continuation](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-train/interview.md#continuation) | Resume a concrete unanswered question |
| A specific correction, changed behavior or new material | [Revision](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-train/interview.md#revision) | Trace that change through existing knowledge |
| No skills are stored | [Initial training](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-train/interview.md#initial-training) | Follow the opening sequence below |

### Continuation

Resume the unfinished interview. When the Project already holds skills
and the person says "resume", "continue training", or returns without reporting changes,
give a short recap of the last
lesson and the unresolved gap, then ask the saved first question for next time. Check it
against the stored answers and this conversation before asking. If it is already answered,
no longer applicable or waiting on unavailable evidence, use the highest-priority remaining
question that can be answered now. Respect a new topic the person explicitly chooses.

Do not open a continuation with generic questions about who they are, what changed,
which answer disappointed them, or what they want to teach you. Do not require a new real
customer question before continuing other unfinished lessons. If they say "nothing has
changed", proceed with the saved worklist. Ask about expertise only if that particular
question needs it; do not infer expertise from API access.

If the saved handoff is missing or too vague, use the existing skills to find a concrete
unresolved product meaning, thin field entry or blocked recipe, draft a focused question,
and ask it. Briefly identify the gap that led to it. The agent chooses the next useful
question; the person should not have to reconstruct the previous interview or invent an
agenda. A proposed question is not evidence that a customer actually asked it.

### Revision

Work on a reported change. When the person brings changed behavior, new
material, a correction or a request to re-check existing knowledge, begin with that
specific change and revisiting stored knowledge. Ask what changed only to fill a relevant gap they have not
already explained. Keep the unfinished question queue available for when that work is done.

### Initial training

The following opening sequence applies only when no skills are stored.

#### First, establish expertise and make a brief scan

Say so, and ask what part of the product the person knows firsthand if that is not already
clear. Do not assume they are the developer: an operator may know the thresholds, customers
and consequences best. Once that scope is clear, **if any data has arrived, spend two
minutes looking before your next question.** Not the investigation, which still comes later:
only what is cheap and on the surface - the Notefile
names, the field names in the schemas, a handful of event bodies from one recent batch of events
rather than a month of them, the fleet and device names, the project environment variables - and
if it will take more than a minute or so, say so first. Then open with what you found. Nobody
wants to start from scratch with a stranger, and **correcting is far cheaper than composing**: a
busy person who would never write you a paragraph about their product
will happily fix a wrong sentence, and fixing it teaches you more than the paragraph would have.

**But say the two halves of what you found in different voices, and never in the same breath.**
Rehearsed on fictional projects, a scan like this was reliably right about **domain and structure**
and unreliable about **meaning**: it identified the product, the fleet size, the cadence, which
files exist and which things move - and it misread what a field measured, which direction it ran,
whether a zero was a fault or a real reading, and whether the device controlled anything or only
watched. One rehearsal read a depth as headroom, so a smaller number meant fuller; on a flood
product that inverts every threshold in the set.

- **"Here is what I can see."** Counts, cadence, file and field names, which fields are sometimes
  absent, fleet and device names, what is conspicuously missing. These are `observed`, and you may
  state them.
- **"Here is what I cannot tell from this."** What each number measures, which way it runs, what a
  zero means, why an event fires, who reads the answer. **Ask these. Never assert them**, however
  obvious the reading feels.

The danger is exact: the true observations and the invented meanings arrive in the same confident
register, and **the correctness of the first is what buys the second a nod**. So offer **the frame
itself** for demolition and not only the details inside it - hedging within a wrong frame makes it
feel stress-tested when it has not been - **invent no vocabulary**, because naming the states of a
status field you cannot decode gets two of your inventions confirmed and leaves the rest
undiscovered, and **say your implicit assumptions out loud**, because a reading nobody stated is
one nobody can correct: if your guess about an array's order decides which element every recipe
uses, that guess is a sentence, not an assumption. Keep the whole thing to a short paragraph you
can say aloud, and then stop. Every word of it is `assumed` until a person confirms it, and
confirming is a read-back like any other (publication).

#### Next, listen to the product story

If nothing has arrived yet, say that plainly - it is itself worth knowing - and go straight to the
question. Either way, ask this next, and listen to the whole answer:

> *"So tell me about the thing itself - what is it, what does it do, and who has one? And then tell
> me how it actually gets used: who works it, how hard, what a normal week looks like, and what
> people do to it that they probably should not. Start wherever I have got it backwards."*
>
> *"And when one of these is out there and somebody needs to know how it is doing, what do they
> need to know - and how do they find out today?"*

**Those two, and nothing else in this turn.** The initial opening has established their scope.
Two questions a person answers as one story is not two asks in the sense lesson loop caps; an opening that *also*
demands every document, whether each is public, and every firmware version is five separate
errands in one wall of text, and the answer you lose is always the one that mattered. The exception
covers those two as written and nothing added to them: a further ask tucked inside either one -
which old units still report, what they do to a device before a test - is a separate question, and
it waits for a later turn.

**Ask the product one first, and let it run long.** Whatever the scan told you, it gave you names
and shapes, not meaning - and meaning is what everything downstream runs on. A Notefile called
`aer.qo` is a string until you know the machine has aeration fans; a field called `mt` is three
numbers until you know what is being measured and why anybody cares. Everything downstream is faster and less wrong once this answer exists, and it is
the half of the job no file anywhere can give you. Ask for the physical thing - what it is, what it
does, who operates it, what goes wrong with it - not for the value proposition, or you will get the
sales pitch and learn nothing. **And do not stop at the object: get the machine in use.** What it is
tells you almost nothing about what its data means; how it is *worked* tells you nearly everything.
A crane lifted twice a day and a crane on a three-shift line produce the same fields and completely
different normal, and no schema, sample or firmware listing will ever tell you which one you are
looking at. If they hand you a variant list, a lifecycle or a customer type,
write it straight into `product` as they say it; the words are `vocabulary` whether or not you
recognise them yet.

**Then the second**, which collects the questions people really ask, **in their own words** - a
recipe built on a real sentence survives a real question and one built on your paraphrase does not.
Stay open to every shape this takes: a technician deciding what to load on the van, an operator
watching a screen, a customer asking why their machine stopped, a maintenance planner, or a
developer with nothing deployed yet who is the only person who has ever looked. And *how do they
find out today* is how you discover the dashboard, the alert or the spreadsheet whose words your
skills have to match, because whatever people already read is the vocabulary they will phrase their
questions in.

#### In the next turn, explain source handling and request useful material

Then, in the **next** turn, **explain source handling, and only then ask for the
documents** - in that order, so they can decide what to share. In your own words:
you will write down the meanings you understood, never the source text, its comments, paths or
symbols; they can review the resulting draft before publication within the agreed scope; anything sensitive that survives the
paraphrase stays out until they say so. Then one short invitation - firmware, requirements or design
docs, manual, datasheet, website, whatever is handy now - **specific about the artifact and about
what each one saves them** (source handling), and **ask of each whether it is public or restricted.** Ask
early: somebody who has to find a specification needs the request before they run out of time. An
offer nobody pinned is no offer, so write it into `questions` against their role and a date; it leaves
that list when the file arrives, not when it is promised, and before the session
ends say which never came. Then do not read the project in silence:
turn what you find into questions as you go - *"there is a Notefile called `env.qo` sending a
field called `lvl` every hour; what is it?"* - because what you find alone means something only
once you know which question it must serve.

## 3. Enter the lesson loop

After the opening, continue with the lesson loop. A continuation can
start at its saved lesson; a revision can go directly to revisiting stored knowledge.
Keep returning to that loop after answers, readbacks and publications. Existing files
and a successful prior session are a starting point for deeper training, not evidence
that there is nothing left to learn.

## Grok Bot only

*This section applies only to agents running as Grok Bot. Other harnesses should skip it.*

- The Cloud Computer's workspace is shared by all of the user's bots, so it is not
  private. Keep training drafts in a clearly named per-project folder there. Tell the
  user they persist until deleted and that other bots can read them.
