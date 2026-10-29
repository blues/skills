# What a trained project gives the next agent

The short index explains the physical product, its usage and variants; points to the
knowledge documents; lists answerable, conditional and blocked question families; and
records the next useful training question. These are two areas of knowledge—what people
mean and what their data means—not a requirement to produce exactly two files.

## Choose the relevant reference

Read the rows that apply to the current step; reuse context already established.

| If you need to | Read | Then |
|---|---|---|
| Test entry into the set as a later answering agent | [Reading product knowledge](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-chat/product-knowledge.md) | Find a recipe and follow its meaning dependencies without interview context |
| Include citations or material derived from sources | [Source handling](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-train/sources.md) | Preserve approved public citations and applicable sharing restrictions |

## Make the product set usable

A later agent follows Reading product knowledge
from the public entry guide into this set. Make that transition usable without the
training conversation: a readable entrance, recognizable questions, exact dependencies
and honest answerability. A cloud assistant will derive private query plans and monitors
from these semantics; the shared set does not embed its scheduler, cache or credentials.

Keep an existing useful layout. A new project starts with `index.md`, `product.md`,
`notefiles.md`, `operations.md`, `recipes.md` and `questions.md`, split further only when
independent retrieval warrants it. Each fact has one owner. API `kinds` make the content
discoverable independently of the Markdown's front matter. Preserve source, confidence,
version/cohort/time applicability and unresolved questions at the owning claim.

Near the top of the index, include a documentation routing table rather than copied API mechanics:

| If you need to | Read | Then |
|---|---|---|
| Resolve an API operation or response schema | [Notehub OpenAPI](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-access/openapi.yaml) | Apply the current endpoint contract through the selected transport |
| Read endpoint explanations and examples | [Notehub API reference](https://dev.blues.io/api-reference/notehub-api/) | Check the applicable operation and behavior |
| Interpret native fields or system Notefiles | [Shared native-field meanings](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-platform/SKILL.md) | Combine the relevant platform interpretation with product-specific meaning |

Approved public product citations may also be useful; apply source handling. The
reader already holds the project it fetched these skills from. Runtime ProjectUIDs
and DeviceUIDs, access instructions, credentials, repository availability and cache
readiness belong to runtime context. Product rules should survive a new deployment;
label unavoidable deployment scope explicitly.

For an unmatched question, say how to reuse a recipe's meaning only when its quantity,
population, time and assumptions match; when to ask for clarification; and when to derive
a bounded answer from the inventory. Never stretch the nearest recipe silently. Index
claims must survive every dependency gate in the recipes, including device identity,
population membership, thresholds and valid measurement definitions.

The set is useful when a stranger can explain the machine in use, decode the relevant
records, distinguish unknown from zero, apply the right cohort/version, and know which
answer would be misleading. An inventory with no recipes cannot bridge to a decision;
recipes with incomplete field meaning hide guesses. Mark thin or unreached areas openly.

## Give the index a usable reading order

Use an index document (kind `index`, commonly named `index.md`) that leads with:

1. A short product/use description and the vocabulary a new reader needs.
2. A routing table showing when to read the product overview, population/configuration
   and field definitions, linking to the document and heading that owns each meaning.
3. A question table: actual or clearly proposed wording, supported/conditional/blocked
   status, recipe name/heading and the specific missing dependency if any.
4. How to handle unmatched wording and when to ask a focused clarification.
5. Training handoff: solid/thin/unreached areas and the next actionable question ID.

Keep this handoff current after every lesson, including after the initial coverage floor
is met. Record the next improvement to pursue and why it changes a useful answer. A set
can support today's questions while still having thinner areas, untested cases or new
evidence worth investigating. Its index should enable the next pass of the interview,
not announce a permanent end to training.

Use explicit routing tables wherever a product document branches into other skills,
recipes or topics: **if you need to do or understand X, read Y, then continue with Z**.
Place the table at that decision point; a narrative hyperlink alone is not sufficient
navigation. Keep supporting citations inline. After reading and internalizing the whole
set, use these routes to revisit the material needed for the task; resolve all
consequential dependencies of the chosen recipe.
For a recipe, a dependency table can name the required meaning, its owning document/heading
and how it affects the answer. Preserve existing useful organization and confirmed content.

Use relative document/heading links within the stored skill set, resolving from each
skill's own name. A nested recipe can refer to its sibling or parent document. Those
links refer to Skills API names, not public site paths. Use full URLs for external/shared
documentation. Verify every internal name and heading in the actual publication set.

For example, a synthetic product's question table might distinguish “Which units are
above the limit now?” (current-state recipe), “Which exceeded it today?” (period recipe),
and “For how long continuously?” (blocked on a measurement-mechanism question). Link each
to its own recipe and dependencies rather than promising all three from one threshold.
Derive the statuses from the actual recipes; the example is not evidence that a real
project supports those questions.

Before publication, enter the set as a fresh reader: use a new phrasing, locate its
recipe, follow all meaning dependencies, resolve private relative links, and identify
what evidence/freshness the answer needs. Repair a missing route or ambiguous owner at
its source. A long training transcript is not a usable entrance.
