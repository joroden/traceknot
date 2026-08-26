# Skill / subagent / context effectiveness analysis

Goal: judge whether a skill, subagent, or upfront context/instructions
actually helped — and if not, whether too little or too much was given.

1. **Find the invocations.** Agent nodes (`kind: agent`) are subagent
   spawns — each has its own subtree. If the harness records skill
   invocations as tool calls, they'll show up as `tool_call` nodes too;
   check `timeline.md` for anything named after a skill/prompt.

2. **For each subagent, compare what it was given vs. what it needed.**
   Open its node file for `spawn_prompt` (the upfront context), then walk
   its subtree in `timeline.md` (nodes owned by that agent):
   - **Under-provided** — the subagent's first several tool calls are
     exploratory (Read/Grep/ls of things the parent already knew, or
     already had open) before it does any real work, or it asks a
     clarifying question / errors out early for missing information the
     parent could have supplied upfront.
   - **Over-provided** — the `spawn_prompt` includes large content
     (files, examples, prior context) that never gets referenced again in
     the subagent's own tool calls or output — dead weight it paid to
     read.

3. **Weigh cost against deliverable.** Each agent node's file states its
   subtree cost (aggregate over all descendant nodes). If that cost is
   dominated by exploration rather than the actual deliverable (compare
   against `result_summary`), the upfront prompt was probably
   insufficient — the subagent had to rebuild context it should have
   received.

4. **Check for signal the agent itself gave.** A subagent's
   `result_summary` or a chat node's Output sometimes states directly
   that it "didn't have enough information about X" or "assumed Y" —
   treat this as strong, direct evidence, not just an inference from
   token counts.

5. **Answer with the specific gap**, not a verdict alone — e.g. "the
   subagent spent $0.80 of its $1.10 subtree cost re-discovering the
   project's file layout (nodes/009–014 are all exploratory reads); the
   spawn prompt should have named the relevant directories" or "the
   spawn prompt included the full 400-line spec, but the subagent's work
   only ever touched one section of it — a pointer would have been
   enough."
