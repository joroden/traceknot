---
name: traceknot-analyze
description: Use this skill whenever the user asks you to analyze a coding agent session or work item. If they want to know why something happened, or why it cost so much, or why a particular decision was made, this skill will help you find the answer in the traceknot telemetry.
---

# Analyzing a traceknot session or work item

Never answer from memory or guesswork. traceknot records real telemetry
for every tracked agent run — export it and read it before answering.

## 1. Scope: one session, or a whole work item?

- **Session** — one continuous agent run, identified by a `session_id`.
- **Work item** — every session claimed against the same issue/ticket,
  identified by a `provider` + `key` (e.g. `github` + `org/repo#12`). Use
  this when the question spans multiple runs ("across everything on this
  ticket...", "did we go back and forth on this?").

If it's not obvious which one applies, or you don't have the id/key, ask
the user rather than guessing.

## 2. What do they actually want to learn?

Ask if it isn't already clear — the answer decides which reference below
you read, and the three kinds of question need different evidence:

| If they're asking...                                                                                                          | Read                                 |
|---------------------------------------------------------------------------------------------------------------------------------|---------------------------------------|
| Why did this cost so much / take so long? Was tool use wasteful? Did it stall, loop, or redo work?                              | `references/cost-drivers.md`         |
| Why was X built as Y and not Z? Why did we deviate from the spec? Was this agreed with the user, or the agent's own call? Why was a feature dropped? | `references/decision-rationale.md`   |
| Was a skill or subagent used well? Should more/less context have been given upfront? Why didn't it use the obvious tool?         | `references/skill-effectiveness.md`  |

More than one can apply — read more than one reference file if so, but
still ask first rather than exporting and hoping the answer falls out.

## 3. Export

```
traceknot export session --id <SESSION_ID>
traceknot export work-item --provider <PROVIDER> --key <WORK_ITEM_KEY>
```

Prints a directory path (a temp dir unless you pass `--out DIR`).

## 4. Read the reference file for the category from step 2

Each one names exactly what to check and in what order — don't skip this
and read the export blind.

## 5. Navigate top-down, not everything at once

- `summary.md` (session) / `overview.md` (work item) — the numbers.
- `flags.md` — statistical cost/duration outliers, failed tool calls,
  repeated-call loops. A starting point, not the full answer.
- `timeline.md` — every turn/tool call, in order, one line each, linking
  to its full file.
- `nodes/NNN-<kind>-<name>.md` — full content, opened only once you know
  it's relevant.
- Grep across `nodes/*.md` for a keyword, error, or file path instead of
  reading everything.
- A section reading `_(identical to <path>#field)_` means that exact
  content already appeared earlier — follow the pointer, it isn't missing.
- A work item's export is `overview.md` plus one `NN-session_xxx/`
  subdirectory per session, each shaped like a single session's export.

## 6. Answer with citations

Name the specific node ids / file paths your conclusion rests on, so the
user can check it themselves.
