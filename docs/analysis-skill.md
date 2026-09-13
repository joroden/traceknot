---
title: Analysis Skill
description: What the traceknot-analyze skill does and when to reach for it.
---

# Analysis Skill

## What is it?

An opt-in skill your own coding agent can use to answer questions about a
session or work item by reading its real telemetry — instead of guessing
from memory of what it did. It's off by default; turn it on per agent in
[Settings & Configuration](settings-configuration.md).

## What kinds of questions is it good for?

- **Cost and time** — why did this cost so much, why did it take so long,
  did the agent stall or loop, was tool use wasteful?
- **Decisions** — why was something built one way instead of another, was a
  deviation from the spec agreed with you or the agent's own call, why was a
  feature dropped?
- **Skill and subagent use** — was a subagent used well, should it have had
  more context upfront, why didn't it reach for the obvious tool?

## How do I use it?

Just ask, in your normal agent session — "why did this session cost so
much," or "why did we end up building it this way." With the skill enabled,
your agent recognizes these as telemetry questions, pulls the real session
or work-item data, and answers by citing the specific step its answer rests
on, so you can check it yourself rather than take the summary on faith.

If it's not obvious which session or work item you mean, or you're asking
about a whole ticket rather than one run, expect your agent to ask which one
before it looks anything up.

## Session vs. work item

- **A session** is one continuous agent run — use this for "why did this
  particular run behave the way it did."
- **A work item** covers every session claimed against the same issue or
  ticket — use this for questions that span multiple runs, like "did we go
  back and forth on this across sessions."
