# Decision rationale analysis

Goal: find the moment a specific decision, deviation, or dropped feature
actually happened, and say whose call it was.

1. **Locate the moment.** Grep across `nodes/*.md` (all sessions, if it's
   a work item) for terms from the decision in question — the feature
   name, the approach that was rejected, the file/function involved.
   `timeline.md` previews are a faster first pass if the terms are common.

2. **Classify what you find into one of three cases** — this is the
   actual answer the user wants:
   - **User-directed** — an explicit instruction or approval in a chat
     node's Prompt ("do it this way", "yes, drop that", answering a
     clarifying question the agent asked).
   - **Agent-initiated** — the agent chose it unprompted; look at the
     Output/Reasoning of the chat node right before the change for its
     stated justification (or lack of one).
   - **Forced by evidence** — a tool result left no real choice (a
     failing test, a linter error, an API/library constraint, a file
     that didn't exist). Check the tool_call node immediately before the
     decision.

3. **If a subagent was involved**, open its agent node file: the
   `spawn_prompt` shows what it was told to do (and whether the decision
   was actually delegated), and `result_summary` shows what it reported
   back to the parent — the parent's next chat node shows whether that
   recommendation was accepted as-is or overridden.

4. **Check the sequence, not just the moment.** Was the decision made
   early (following the original ask) and never revisited, or made once
   and then changed later? A later change usually has its own trigger —
   find it the same way (grep + classify).

5. **Answer with the classification plus the citation** — e.g. "this was
   agent-initiated: node_9af... reasons that the simpler approach avoided
   a new dependency; the user was never asked" or "user-directed: the
   user's message in node_002 explicitly asked to drop the feature."
