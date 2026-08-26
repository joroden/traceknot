# Cost & efficiency analysis

Goal: explain where the money/time actually went, and whether any of it
was avoidable.

1. **Start with `flags.md`.** Read every `cost_outlier` and
   `duration_outlier` entry first — these are the turns/tool calls that
   are statistically abnormal for this session, not just the biggest
   numbers.

2. **For each flagged chat node**, open its file and check the Prompt:
   was a large chunk of prior tool output or file content pasted back in
   rather than referenced? Repeated large prompts across several
   consecutive chat nodes usually means context was being re-sent instead
   of reused.

3. **For each `repeated_call` flag**, this is wasted cost by definition —
   the same tool with the same arguments fired 3+ times. Read the chat
   node(s) around it: did the agent get no useful error/feedback after
   the first call, so it had no signal it was repeating itself? Was the
   instruction ambiguous enough to cause retries?

4. **For each `error` flag** (failed tool call), check what happened
   next — did the agent retry blindly, or did it read the error and
   adjust? Repeated cost after an error with no visible change in
   approach is a concrete inefficiency.

5. **Check the cache ratio** in `summary.md`/`overview.md`
   (`cached` vs. total input tokens). A low cache-hit ratio across many
   turns in one session usually means something large (a system prompt,
   a big file) was being resent instead of staying cached — look at
   whether the session's shape (many small back-and-forth turns vs. long
   uninterrupted tool sequences) explains why.

6. **Count dedup pointers** (`_(identical to ...)_`) across the export —
   each one is a tool result that was fetched again unchanged. A cluster
   of these around one file/command is direct evidence of redundant
   re-fetching.

7. **Conclude concretely.** Don't just say "it was expensive" — name the
   specific loop, re-fetch, or oversized prompt, cite its node(s), and
   say what a leaner path would have looked like (e.g. "should have
   grepped for the function instead of reading the whole 2,000-line file
   three times — see nodes/012, 034, 057").
