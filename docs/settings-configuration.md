---
title: Settings & Configuration
description: Every traceknot setting and what it actually changes.
---

# Settings & Configuration

Run `traceknot` with no arguments to open Settings in your browser — the
same page setup uses, and reachable from the dashboard's sidebar at any
time.

## Start the daemon on login

Off by default. When on, traceknot starts automatically when you log in, so
the dashboard and hooks are ready without running anything manually first.
When off, you start it yourself before your first session of the day.

## Require a work item

Off by default: a session pauses to ask for a work item, but you can skip
and it's tracked and costed the same either way. Turn this on to remove the
skip option entirely — the agent itself is blocked and won't start working
on that prompt until something is claimed. See [Work Items](work-items.md)
for what claiming looks like, including the free-text custom option.

## Agent hooks

One switch per supported agent (Claude Code, Codex CLI, Copilot CLI, VSCode
Copilot Chat). Turning a hook off stops traceknot from tracking that agent
at all — no pause, no telemetry, nothing recorded — without affecting the
others.

On Windows, the Copilot CLI integration adds a function to your PowerShell
profile. If profile scripts are disabled by the effective `Restricted`
execution policy, installation saves the current-user policy and changes it to
`RemoteSigned`. `traceknot uninstall` restores the saved policy when it is still
`RemoteSigned`; if you changed it afterward, your choice is left in place.
Launching Copilot CLI through that profile function also grants traceknot's
Copilot extension permission for the current repository without another prompt.
Turning off the Copilot hook or uninstalling traceknot removes those extension
permissions.

## Analysis skill

One switch per agent, off by default. Turning it on lets that agent read its
own session and work-item telemetry when you ask it to. See
[Analysis Skill](analysis-skill.md) for what it does and how to use it well.

## Other commands

```
traceknot            # open Settings in your browser
traceknot uninstall  # remove traceknot (keeps your recorded data — see
                      # Data Collection & Storage)
traceknot help       # usage
```
