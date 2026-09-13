---
title: Work Items
description: Connecting a provider, claiming with a custom id, and making claiming mandatory.
---

# Work Items

A work item is the GitHub, GitLab, or Jira issue — or a free-text label — a
session gets claimed against, so cost rolls up under real work instead of
sitting as a flat list of sessions.

## Connecting a provider

The picker can search your real issues instead of typing a label by hand,
using whichever CLI you already have installed and signed in — there's no
separate traceknot-specific login:

- **GitHub** — [GitHub CLI (`gh`)](https://cli.github.com), signed in via
  `gh auth login`
- **GitLab** — [GitLab CLI (`glab`)](https://gitlab.com/gitlab-org/cli),
  signed in via `glab auth login`
- **Jira** — [Atlassian CLI (`acli`)](https://developer.atlassian.com/cloud/acli/guides/install-acli/),
  signed in via `acli jira auth login`

traceknot only reads through these CLIs to fetch an issue's key and title —
it never creates, edits, or comments on anything in GitHub, GitLab, or Jira.
A provider tab you haven't signed into shows as unavailable rather than
blocking the picker; the other tabs and Custom still work.

## Claiming without a real issue

There's always a **Custom** tab in the picker, regardless of which
providers you've connected. It takes any id and title you type — nothing is
checked against GitHub, GitLab, or Jira. Use it for work that doesn't have a
tracked issue yet, or when you'd rather not look one up mid-session.

## Making a claim mandatory

By default, a session pauses to ask for a work item but lets you skip —
there's a **Skip Context** button, and pressing Escape does the same thing.

Turning on **Require a work item** in
[Settings & Configuration](settings-configuration.md) removes that option
entirely: there's no Skip button and no Escape shortcut. This is a hard
block, not just a UI nag — the agent itself is refused and cannot start
working on that prompt until something is claimed, so it's worth turning on
only if you actually want every session to stop dead without a work item
attached. A Custom entry still satisfies this — "required" means *some*
claim, not specifically a real tracked issue.
