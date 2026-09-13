---
title: Data Collection & Storage
description: Exactly what traceknot records, where it lives, and what leaves your machine.
---

# Data Collection & Storage

## What gets recorded?

For every tracked agent session:

- Every prompt, chat turn, and tool call, in order
- Every subagent spawned within a session
- The model used for each step
- Token counts (input, output, and cache-tier breakdown) for each step
- The work item the session is claimed against, if any

See [Session Tree](session-tree.md) for how this is structured once you open
a session, and [Pricing & Cost Calculation](pricing-cost-calculation.md) for
how token counts turn into a cost.

## Where is it stored?

Everything is written to a local SQLite database at
`~/.traceknot/telemetry.sqlite`, served by a local dashboard at
`http://127.0.0.1:4318`. There's no remote database, no sync, and no
account tied to it.

## What leaves my machine?

Nothing, except one thing: if you claim sessions against real GitHub,
GitLab, or Jira issues, traceknot makes a read-only lookup through your own
signed-in CLI (`gh`, `glab`, or `acli`) to fetch that issue's key and title.
It never creates, edits, or comments on anything in those systems. See
[Providers & Hooks](providers-hooks.md) for setup. Skip that entirely and
traceknot makes no outside calls at all.

## Is my history backed up anywhere?

No. `~/.traceknot/telemetry.sqlite` is the only copy of your recorded
history. If that file is lost, the history is gone — there's no
account-linked copy to recover it from. Back it up the same way you'd back
up any other local file that matters to you.

## What happens if I uninstall?

`traceknot uninstall` removes the binary and every agent hook, but leaves
`~/.traceknot` in place, so reinstalling later doesn't lose your history.
Delete that folder yourself if you want the data gone too.
