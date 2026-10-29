# Documentation server features

The [Notehub workflow](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub/SKILL.md) and all its references can be read directly from
GitHub Raw. Installation and generated indexes are optional conveniences provided
by a documentation server. If your agent can only read GitHub Raw, follow the
workflow's routing tables and fetch the linked Markdown files directly.

## Install a live pointer

When using a documentation server, its base URL is the scheme and host, including
any port, of the document URL you fetched. Use that base URL with `npx skills add`;
for example, a local development server can be installed with:

```sh
npx skills add http://localhost:8888
```

Use the actual accessible documentation server's base URL for a hosted installation.
Do not pass the GitHub Raw repository prefix to this command: it serves source
files and does not provide the generated installer manifest.

Skills installed from a documentation server are live pointers. Before using one,
fetch the current skill from its stated URL; it follows the same reading journey
without reinstalling every documentation update.

## Use generated indexes

These paths are relative to the documentation server's base URL. They are generated
by that server and are not files in the GitHub repository. An agent restricted to
GitHub Raw should use the [workflow's routing tables](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub/SKILL.md) for discovery instead.

| If you need | Request from the documentation server | Purpose |
|---|---|---|
| A discovery index | `/llms.txt` | Find the relevant entrypoint |
| The complete corpus, expressly needed for the task | `/llms-full.txt` | Read the aggregated public skills |
| Machine-readable installation metadata | `/.well-known/skills/index.json` | Discover installable skills |
