# Notefile and field inventories

Use this during a field or configuration lesson to account for all relevant evidence,
including conditional and system fields. Read product meaning alongside the appropriate
platform definitions; mark unresolved meaning explicitly.

## Choose the relevant reference

Read the rows that apply to the current step; reuse context already established.

| If you need to | Read | Then |
|---|---|---|
| Interpret native objects or system fields in the inventory | [Platform definitions](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-platform/SKILL.md) | Record product applicability without duplicating shared meanings |

## Cover meaning, validity, significance and action

The first answer is the top of a ladder, not the end of it. Every field, every fleet and every
threshold has four rungs, and you do not leave it until all four are filled, written down as open,
or named as inapplicable - named, never quietly invented:

1. **What does it mean, in units?**
2. **What does it look like when it is wrong?** What do I see when the sensor has failed, the
   reading is stale, or the thing is disconnected - and how do I tell that from a real measurement?
   Can the sentinel value ever be a real reading here - a machine genuinely drawing 0 amps because
   it is switched off, a counter legitimately at 0 after a reset? If yes, write the rule that tells
   them apart; if no, write why, as a `stated`
   sentence.
3. **What number divides normal from notable from alarming?**
4. **What do you do about it today?** Who do you tell, how fast, and what must an answer never
   claim?

Rung two is where sentinels come from, rung three the thresholds that exist nowhere in the data,
rung four most of `audiences` and `constraints`. A rung that does not apply is explicitly marked not applicable with a reason; a relevant
rung you cannot fill becomes a question naming who can; a rung you skip becomes a confident wrong answer six months from now. It applies to every
number a recipe cuts on, not only to fields - a service interval, a forecast horizon, a silence
threshold, the anchor of "this week" - and one whose origin is `assumed` puts its own value in the
answer sentence ("due in 40 running hours, against an interval nobody has confirmed").

**So mark it, on the field, where somebody will act on it.** Every field carries its rung state the
way every claim carries its origin: `[rungs:1-4]`, `[rungs:1,2]`, `[rungs:1 only]`, on the field and
not once per file. Without the marker a file one question deep is indistinguishable from a finished
one: you write "bit 0 is the main contactor and bit 2 is the overload trip", and the file reads as
a complete description of a bitfield whose other six bits you never asked about.

Two habits get you up the ladder faster. **Ask for the exception** - "when is that not true?" and
"when does that break?" produce more than any question about the normal case. And **a method you
invented is not a rung**: any formula, threshold or gap rule you devised is worked through on a representative case (label synthetic evidence) and read back
before a recipe may rely on it, and until then the answer sentence says
"estimate by a method nobody has confirmed".

## Inventory the Notefiles and their fields

This is the inventory, and it is the deliverable a later agent cannot do without.

**Write the roster before you write a single entry.** The first thing in this file is a plain list
of every Notefile name you have seen anywhere, with nothing beside it but the source that saw it.
Build it **source by source, naming each source as you exhaust it** - the schemas, then the sampled
events, then the firmware, then any existing skills - because a roster built from whichever source
was cheapest looks complete and is not. **A source you read after writing the roster sends you back
to the roster first**: the firmware is usually read late and may reveal a Notefile that never
appeared in the sample, so its names are especially easy to miss. Build the list
mechanically, before you know what any of them mean. Then work down it, and **the file is not
finished until every name on the roster has its own entry below.** Do this in the other order and
you will write up the three you understood and lose the rest, which is the single commonest way
this step fails.

**Every name on the roster gets an entry** - including the system files the Notecard and Notehub
maintain, whose names usually begin with an underscore, and the ones you found in the firmware but
never saw arrive, marked as such. **Those are exactly the ones that get dropped**, and a set that
covers the application Notefiles and skips `_session.qo` cannot answer "has this site gone dark".
If you judge one irrelevant, that judgement is itself its entry: one line, naming it and why.
**A coverage warning at the top of the file is not an entry, and neither is an open question** -
both say a name is missing without letting anybody look it up. A Notefile absent from the roster's
entries is a question nobody can answer.

**For native system files, inherit common definitions explicitly.** Read the relevant
`notehub-platform` reference, then keep this project's entry: exact file name, observed
field roster, producer/direction, applicability, product interpretation, configuration
dependencies, firmware differences, cadence and silence meaning, custom additions and
exceptions. The common field dictionary may
be referenced by its named section through the platform documentation pointer in
`index.md`; do not duplicate it in each project. Mark shared field meanings as covered
only after verifying that the producer, template and generation match. A blanket
"see platform docs" does not account for unlisted or changed fields: give each such
field its own definition or explicit unresolved entry. Shared knowledge does not
excuse missing a system Notefile from the roster.

For custom tracking or a renamed system file, retain the actual file name and its
producer alongside the shared definition it uses. Before publishing, check that each
observed system file has an entry and that custom fields, generation differences and
unresolved meanings remain visible. Preserve the existing publication review workflow.

Include the native metadata and API objects needed for the product's questions, even
when they are outside a Notefile body: for example Event capture/location metadata,
Session usage and high-power counters, Device freshness, and DFU reported versus
requested state. Link them to the affected recipe using exact paths and the shared
definitions, and record project-specific interpretation alongside them. Do not invent
a body field just because the training inventory is organized by Notefile.

If the inventory outgrows one file, split it per Notefile and let `index.md` route by Notefile
name; that is the one split that needs no justification. Each entry carries, in one compact block:

- **The exact name on the wire**, spelled as it appears - `motor.qo`, `_session.qo` - never a
  prettified version, and the suffix convention it follows.
- **Its role in the product**, in one sentence a stranger would understand: what it is for, not what
  it contains.
- **Direction** - outbound from device, inbound to device, or maintained by the service - and **who
  writes it**: the host firmware, the Notecard itself, Notehub, or a route.
- **Cadence and trigger**: what makes a note appear at all - a timer, a threshold crossing, a state
  change, a button - what governs the rate, and what changes it. **What its silence means** is part
  of this: a file that appears only on alarm and one that should appear hourly imply opposite things
  when nothing arrives.
- **Whether it is templated**, because that decides whether an omitted field arrives as a zero, and
  **which generations send it** where they differ.
- **Which timestamp means what**: the field carrying measurement time and the field carrying arrival
  time, **by name**; whether the device clock can be unset or wrong on this file; and which of the
  two a time-bounded question about this Notefile should use. Every time-bounded question in this
  document depends on that answer, and no schema states it.

Then **field by field**, and do not stop at the ones you understand. The field list is the union of
the schema, the sampled events, existing knowledge, and every body the firmware builds - templated, untemplated and
conditional alike, since a field only sent on one branch is absent from the template and from most
samples. No one source is complete, and a field present in only one is recorded with which
source saw it:

- **The name exactly as it appears on the wire**, and **where it sits** - top level of the body,
  nested under a named object, an element of an array - as the path a reader would follow.
- **What it means, in the product's terms**, in a sentence that does not use the field name as its
  own explanation. `lvl` is not "the level".
- **Type and encoding on the wire**: number, string, boolean, array, object; integer or real; the
  wire type the template declares, which decides what a test on the field can detect. For an array,
  **what position means** and whether the length is fixed - `[12, 14, 11]` may be three samples in
  time, three sensors, or a value with its minimum and maximum. For an integer that is really an
  enum, **every value and what each means**. For a bitfield, every bit, or a note of which you never
  asked about.
- **Units and scale**, always - tenths of a degree, millivolts, seconds since boot, a percentage of
  what - whether the value is measured, derived or configured, and its **range**, so that a reader
  can recognise a broken one.
- **Sentinels and absence**: what `0`, `-1`, an empty string, a null and a missing key each mean
  here, whether the field is ever omitted, and **whether the sentinel value can also be a real
  reading** - with the rule that tells them apart, or the reason there is none.
- **Validity**: how the field can be wrong while looking right - a stale cache, a sensor that reads
  a fixed value when disconnected, a clock that never set, a value carried over from the last note.
- **The rung marker** (four rungs), on the field, and **aliases**: where the same quantity appears
  under another name in another Notefile, and which one a cross-file question should use.

Two rules keep the inventory usable. **Names go both ways**: the product word for a field is
recorded beside it and `vocabulary` maps it back, so an agent can arrive from either direction. And
**write what you do not know as a gap, not a silence**: "`rs`: unknown, seen only on v1 units, Q12"
is a useful line; an absent field is a reader's confident mistake.

**Then reconcile before you leave, in both directions.** First walk each *source* back to the
roster - schemas, samples, firmware, existing skills, one at a time - because a roster that is
short by a Notefile still reconciles perfectly against its own entries, and "6 on the roster, 6
entries" is a number that can be true and worthless. Then walk the roster to the entries, and do
the same for the field paths. Report both counts out loud - *"4 sources walked; 7 names on the
roster; 7 entries"* - because a number you have to say is one you have to check. Anything still unticked is either written now or written as its one-line
entry saying why not. This sweep is the difference between an inventory and a sample of one.
