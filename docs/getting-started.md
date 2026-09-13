---
title: Getting Started
description: What traceknot is, and the local-first idea behind it.
---

# Getting Started

traceknot is a local telemetry collector for AI coding agents. It ties every
agent session to a GitHub, GitLab, or Jira work item, and gives you a
dashboard of real cost and activity per task — not per month, per task.

## The idea

Every agent session pauses early on until you attach it to a work item.
Attribution is enforced at that point, not guessed at afterward by matching
timestamps or branch names. Once a session is claimed, its full cost —
tokens, cache tier, tool calls — rolls up under that task automatically.

## Everything stays on your machine

traceknot runs as a local daemon with a local dashboard. There's no account,
no cloud sync, and no analytics collection. The only outside calls it makes
are read-only lookups against GitHub, GitLab, or Jira to resolve an issue's
key and title — see [Data Collection & Storage](data-collection-storage.md)
for exactly what that means and what's stored where.

Because everything lives on your machine, it's also only on your machine:
there's no server-side copy to restore from if you lose the local data.
Back it up the way you'd back up any other local database if your session
history matters to you.

## Next steps

- [Providers & Hooks](providers-hooks.md) — which agents are supported and
  how sessions get claimed
- [Settings & Configuration](settings-configuration.md) — every toggle and
  what it does
- [Session Tree](session-tree.md) — how to read a session once it's recorded
