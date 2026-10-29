# Answer shapes and semantic recipes

Use this to turn real customer questions into reusable product meaning and suitable
answers. A future reader follows these recipes through the product knowledge workflow;
a cloud assistant derives executable question plans from them.

## Choose the relevant reference

Read the rows that apply to the current step; reuse context already established.

| If you need to | Read | Then |
|---|---|---|
| Check how a later agent finds and uses a recipe | [Reader workflow](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-chat/product-knowledge.md) | Verify that all consequential meaning dependencies are reachable |
| Turn a semantic recipe into a private cloud query plan | [Prepared questions](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/questions.md) | Keep executable plans and runtime state outside shared product knowledge |
| Learn from an existing customer screen or project evidence | [Project evidence lesson](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-train/lessons.md#e-project-evidence) | Record the expected answer shape and its actual data source |
| Check consistency and publish recipes | [Publication](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-train/publishing.md) | Verify owning rules, dependency gates and index statuses together |

## Choose the answer shape

Half of knowledge area 1, and the half almost every session skips: an agent that computes the right number
and delivers it in the wrong form has still failed the person who asked. **Ask for the form as part
of asking for the question**, in lessons B and C - *"if I got you that, what would you do with it,
and where would it end up - a screen somebody watches, an alert to whoever is on call, a work order,
a line in the report that goes to the customer?"* - and record in `presentation` which shape each
family of question wants, and where two audiences want the same question in different shapes or
units, record both:

- **Spatial - a map.** Where is it, which sites, what is our coverage, which region is worst. Say
  what the points are, where the position comes from, how stale it may be, and whether position is
  reliable enough to draw at all.
- **Temporal - a line or a series.** Is it getting worse, when did it change, what did last month
  look like. Say the interval it is bucketed at, what an empty bucket means - no measurement is not
  zero - and whether gaps are drawn as gaps.
- **Comparison - a ranking or a bar.** Which is worst, who consumes most, the top ten sites. Say
  what is ranked, over what window, and that the ranking names those it could not speak for.
- **Distribution - a histogram or percentiles.** How bad does it usually get, what is the tail. Say
  which percentile these people use: operations often care about the worst case, not the mean.
- **A single state - a sentence with a number in it.** Is this one all right, how many hours are on
  it, when did it last report. Most real questions are this one, and a chart is the wrong answer to it.
- **A decision - go or do not go, with the reason attached.** *Does this need a visit? Can it wait
  until the scheduled service? Do we send a part with the van?* This is the shape a lot of these
  products exist to produce, it is the one nobody thinks to write down, and it is not a number: it
  is a recommendation, what it rests on, and what would change it. Say which way to fail when the
  evidence is thin - an unnecessary visit and a missed failure are not equally expensive, and only
  these people know which way round it is here.
- **A refusal, or a qualified narrative.** When the data cannot support the claim, the answer is
  prose: what *is* known, what is not, and what would settle it - a legitimate answer shape the
  skills must name, or an agent will draw a confident chart instead.

For each shape write four more things, because they are what make an answer usable. **The units the
audience thinks in** - run hours against the service interval, cycles against the rated life,
degrees or degree-hours, events per machine-week; if a service manager asks "how long before this is
due", a raw counter reading is not an answer. **The precision
the evidence supports** - one decimal because the sensor stores tenths, whole hours because the
cadence is fifteen minutes, none at all on a figure from three samples: **never report a number more
precise than its evidence**. **The baseline that makes it meaningful** - last month, the fleet
median, the contractual limit, this site's own history; a bare number with nothing to compare it to
is the commonest useless answer. And **when a visual would mislead**: a map drawn from three located
devices out of two hundred, a trend line over a period containing a sensor swap, a ranking whose
denominator differs by row, a series drawn through a gap as though it had been measured - in each
case say what to draw instead, or to say it in words. Note also **where the answer is going**,
because the destination decides the form: a dashboard can carry a chart, an alert to somebody on
call is one line and has to survive being read on a phone screen at the roadside, a work order needs
the part number and what the technician will find, and a report to the equipment's owner wants the
definition printed beside the figure. And if there is already a screen people watch (the project evidence lesson), use
its words and its statistics - somebody checking your answer against that screen will believe the
screen. **A shape you chose rather
than heard is `assumed`**, and goes in `questions` like any other.

## Build the semantic recipe

A recipe teaches **what an asker's question means here, what data answers it, what may not be
concluded from it, and what the answer should look like**. It is a semantic bridge, not a runbook:
it must stand alone for an agent holding the skills and its own access - no trainer, no firmware,
no dependency on the trainer's machine, scripts or credentials - and that agent derives its own calls. All six parts, or it is not done:

1. **The question, in the customer's words**, as somebody actually said it, who is asking, and the
   status: verified within stated bounds, conditional on something named, or blocked.
2. **The ambiguities that must be resolved before it can be answered at all**, each either citing a
   rule in the set by file and heading or reading *"ask the asker X first"*. **Identity**: how this
   asker names the thing - short forms, partial numbers, their own names for sites and units - and
   the normalisation that resolves it (case, punctuation, abbreviations, leading zeroes), matched on
   whole tokens and never a substring, with a no-match branch that asks rather than guesses, and the
   matched name read back. **Window**: what this asker's time words mean, using UTC unless the
   user explicitly instructs otherwise, which of the Notefile's two clocks the question wants,
   and the coverage required, which the answering agent must verify at runtime. **Population**: which devices count and
   which must be excluded, by an evaluable rule, citing
   `population`; and where the question is spatial, how the set establishes that a device is where
   the asker thinks it is - position, station, fleet membership, an environment variable - asking
   which when two places share a word. **Audience**: what this asker may be given, citing
   `audiences`, and their word for the answer.
3. **What the question means in this project's data**: which Notefile, which fields, which
   conditions conceptually - "notes from `motor.qo` where the highest of the three phase currents
   exceeds the nameplate rating, over units in the rental fleet" - naming fields by their wire
   names and pointing at their inventory entries. Not a call sequence: the *semantics* of the
   selection. **Name two reductions; neither follows from the other.** *Within a note*, which
   element or statistic of a multi-valued field the rule cuts on - noting that over a non-empty
   set of valid readings "the maximum exceeds the limit" and "any one of them exceeds it" are the
   *same* rule, while the mean, "all of them" and one named probe genuinely differ. *Across a
   period*, how one device's many observations become one result per bucket - *any exceedance in
   the day* (which is the same rule as *the day's maximum*), *the last observation of the day*,
   *the day's mean*, *every observation*, and *time spent above* are genuinely different
   questions, and a counting unit and denominator do not choose between them. Then the counting
   unit (events, devices or device-days), the denominator, duplicate handling, version branches
   and reset handling.
4. **What the data does and does not support for this question.** Any recipe whose question involves
   a duration ("for N minutes", "how long"), a total from a rate, a volume from a level or a count
   from a counter carries one line labelled *Inference*, stating what the measurement mechanism
   supports and what it does not. **If the mechanism does not support the inference the question
   needs, the recipe's status is conditional or blocked - never verified.** Two above-limit readings
   fifteen minutes apart do not establish fifteen continuous minutes above the limit, because the
   value may have fallen between them - so *Inference: instantaneous samples at a nominal 15-minute
   cadence; supports "above the limit at these sample times", not "continuously above for any
   span"*, the status is conditional, and the answer names the intervals nobody measured.
5. **The shape of the answer** (answer presentation), chosen and named: map, series, ranking, distribution,
   sentence or narrative refusal - with which axes or dimension, in which units, at what precision,
   against which baseline, for which audience. Then the **answer template** in the customer's words,
   read back against `constraints` and every caveat the owning derivation attaches to the number, in
   which **every figure is a placeholder** - `<N>`, `<X.X>`, `<date>` - so no reader mistakes a
   specimen computed during training for a result. And **what this answer must never claim**.
6. **Dependencies** - the files and headings a reader must load: Notefile entries, units, version
   cohorts, sentinels, population rules, lifecycle breaks, thresholds, audience limits, and any open
   question that bounds the answer. A link is not evidence the reader loaded it, so list them, and
   say what to do if one is missing or conflicts.

**Inferences a recipe may not make.** Each produces a confident wrong answer, and each needs a
mechanism or a person first:

- **A span between samples is not a duration.** Claim it only where the device reports the extreme
  that **bounds the whole interval** on the side the question asks - an interval minimum for an
  above-limit claim, a maximum for a below-limit one - with its coverage established. **A mean or
  an integral does not establish it:** 9 then 9 has the same mean as 7 then 11 and a different
  time above 8. A person-confirmed interpolation is an agreed estimate and says so, with its gaps.
  Otherwise report what was sampled and name the intervals nobody observed. The same governs dose
  from a rate, volume from a level and a total from a counter, and **every such recipe carries
  part 4's inference line**.
- **Missing data is not a negative result.** A gap is not "no excursion", not zero, not healthy; an
  event shorter than one interval falls wholly between notes, so a "no" at any cadence says which
  intervals were unobserved.
- **A trend is not a forecast.** *Is it rising? What about overnight? Where will it be worst
  tomorrow?* get asked of every product that measures anything, and no measurement supports them
  by itself. A direction over past samples is a statement about those samples; projecting it
  forward needs a mechanism somebody has confirmed - a wear rate, a duty cycle, a model - and
  where there is none the recipe is blocked and `constraints` carries the refusal in the audience's
  words. **Record the forward-looking questions its askers actually have**, and whether they are
  supported or blocked. If none apply, say so; do not invent a customer utterance.
- **A shared name or serial is not one asset.** Join history across identifiers only through an
  established mapping with effective periods; without one, report per identifier and say the history
  is not established. Likewise a configured value is not an applied one, and an upload time is not a
  measurement time - use the two timestamp fields the inventory names.

**A recipe inherits every gate**, named among its dependencies - the `population` rules, `sentinels`
and `scenarios` windows it applies - and any recipe returning a list or a count also returns **the
set it could not speak for**, by name and with the reason: excluded, silent, sentinel-faulted, or
missing configuration. **For "latest" and "current" questions, take the latest event first and judge
validity second**: a newer failed measurement must never disappear behind an older good one - that
turns a broken sensor into a healthy reading - and a device in the requested population with nothing
to report stays visible as unknown. **This gate belongs to current-state questions and does not
travel to period questions.** A device reading a valid 9 against a limit of 8 in the morning and
faulting to a sentinel at noon is *unknown now* and *did exceed the limit today*; applying the
latest-only gate to the daily count erases the exceedance. A recipe states which of the two it is,
and a shared "common gates" section says that this one is scoped rather than universal.

**Write the current-state recipe, because it is the commonest call and the easiest to skip.**
*Is this one all right? How many hours are on it? When did it last report?* A set that scopes the
latest-before-validity gate but never writes a recipe that uses it has left the rule as a note
about itself, and the next reader has to quote a disclaimer as though it were an instruction.
That recipe says in its own words: take this device's latest note first, judge its validity
second, and **where the latest note is invalid, the answer is unknown - never fall back to an
earlier good one.** Any "skip the invalid notes" wording belongs to period recipes only, and a
recipe carrying it says so, because read on a "right now" question it means exactly the walk-back
this gate forbids.

**A fleet-or-period question is a different question.** *How many across the whole fleet last
quarter? Is it getting better or worse? Which sites are worst? Where are they?* are a distinct
family from *what is this device doing now*, and a set that only describes single-device lookups
has trained half its job, so **include a recipe in this family or explain why it does not apply.** Name the
Notefiles and fields, the counting unit, the numerator and an **observed-coverage denominator** -
two agents counting different units return different numbers for the same quarter, and any
grouping drops silent devices by construction - and keep excluded and unknown devices visible in
the answer rather than dropping them. **Write the denominator as a test a reader can evaluate, not
as a phrase.** "The devices that reported" is not one: a device that reported hourly with a dead
sensor both reported and cannot be spoken for, so two readers following you faithfully publish
different numbers. Say which side it falls on - *"counts only devices with at least one valid
reading in the window; the failed-sensor ones are reported beside it and are not in the
denominator"* - and the answer template names both counts. Say what reading raw events would cost at the scale `config`
records. Where the *data* cannot answer - the measurement was never taken, the population cannot
be identified, the retention window is shorter than the question - the recipe is **blocked**, and
names the missing fact, the narrower answer that *is* supported, its refusal wording and the
question to ask. Blocked means blocked by meaning, never by which interface happens to be switched
on today; the reader sees that for itself. **A blocked recipe is still a recipe**, and silence in
its place leaves the reader to invent the meaning.

**Before you call a recipe verified, try it.** Where access permits, run the bounded read-only reads
it needs and check that field names, null behaviour, row counts and coverage are what the inventory
says; work one representative case, one edge case that would change the answer, and a row exactly on
the window boundary. publication's consistency pass checks the citations. **Work the awkward cases
explicitly into the recipe**: where the interpretation branches - a version test, a sentinel, a real
reading equal to the sentinel, a short array from old firmware, a counter reset, a duplicate, a
delayed upload - list the cases and the expected result for each, labelling synthetic inputs as
synthetic, because prose describing a branch is not testable and a table of cases is. **Then read it
as a stranger**: take a realistic *new* wording and answer it using only the index, the recipe and
its declared dependencies. Can you find the population, decode the fields, see which Notefile to
read, notice missing coverage, and phrase the answer in the right shape? Fix what is missing, record
what you checked, and mark an untried plan unverified rather than holding publication hostage.
