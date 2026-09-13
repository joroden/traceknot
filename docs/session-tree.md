---
title: Session Tree
description: How a recorded session is structured, and what its numbers mean.
---

# Session Tree

If you've used an LLM observability tool like Langfuse before, the idea is
the same: every session is a tree you can drill into, not just a total.
traceknot's version is purpose-built for coding agents specifically, so the
tree is shaped around prompts and tool calls rather than generic LLM calls.

## The structure

- **Prompts** — each message you send to the agent
- **Turns** — the agent's responses within a prompt, including any back and
  forth with tools before it replies
- **Tool calls** — every tool the agent invoked (reading a file, running a
  command, searching, and so on)
- **Subagents** — any subagent the main agent spawned, shown nested under
  the tool call that launched it, with its own full tree underneath

Every node — prompt, turn, tool call, or subagent — has its own model,
token, and cost breakdown, so you can see exactly which step drove a
session's total rather than only the total itself.

## Estimated vs. exact token counts

Two different things are labeled "tokens" in a session, and they're not
measured the same way:

- **LLM response usage** (the input, output, and cached tokens for a turn)
  comes directly from the provider's own reported usage — exact, not
  estimated.
- **Tool call token counts** are estimated locally, because providers don't
  report token counts for tool input/output. traceknot still calculates
  them so you can spot which tools are burning through context — a file
  read that pulls in a huge blob, a search with results too broad to be
  useful — and change how you use them accordingly. Treat these as a close
  approximation for that purpose, not as billing-grade numbers: only the
  LLM response usage feeds into the session's cost.
