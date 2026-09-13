---
title: Providers & Hooks
description: Which agents traceknot supports and how sessions get claimed.
---

# Providers & Hooks

## Which agents are supported?

- Claude Code
- Codex CLI
- GitHub Copilot CLI
- VSCode Copilot Chat

## How does traceknot know a session started?

Each supported agent has a hook installed during setup that fires at the
start of a session. That hook is what pauses the session and asks you to
claim a work item — see [Settings & Configuration](settings-configuration.md)
to make claiming optional instead of required.

If you skip claiming, the session still gets recorded — it just shows up as
unclaimed in the dashboard rather than rolled up under a task.

## Do I need to configure anything for a supported agent?

No. Installing traceknot configures the environment variables each agent
needs to call its hook — you don't edit any agent config yourself.

One thing to know: an agent already running when you install or reconfigure
traceknot won't pick up the new environment variables. Launch it from a new
terminal window afterward, and for VSCode, restart every open window.

## What if I only want some agents tracked?

Each provider's hook can be turned off independently without uninstalling
traceknot — see [Settings & Configuration](settings-configuration.md).

## Claiming a work item

Once a session pauses, you pick the work item it's for — see
[Work Items](work-items.md) for connecting GitHub, GitLab, or Jira,
claiming with a free-text custom entry, and making a claim mandatory.
