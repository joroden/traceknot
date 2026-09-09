import { joinSession } from "@github/copilot-sdk/extension";
import { spawn } from "node:child_process";

const TRACEKNOT_BIN = "TRACEKNOT_BIN_PLACEHOLDER";
const ABORT_DELAY_MS = 500;

function runClaim(sessionId) {
  return new Promise((resolve) => {
    const child = spawn(TRACEKNOT_BIN, ["claim", "--agent", "copilot"], {
      stdio: ["pipe", "ignore", "ignore"],
    });
    child.on("error", () => resolve(0));
    child.on("close", (code) => resolve(code ?? 0));
    child.stdin.write(JSON.stringify({ session_id: sessionId ?? "" }));
    child.stdin.end();
  });
}

const session = await joinSession({
  hooks: {
    onUserPromptSubmitted: async (input, invocation) => {
      const exitCode = await runClaim(invocation.sessionId);
      if (exitCode === 2) {
        await session.log(
          "No work item was assigned to this session and assignment is required. The run was not allowed to continue.",
          { level: "error" },
        );
        setTimeout(() => {
          session.abort().catch(() => {});
        }, ABORT_DELAY_MS);
      }
    },
  },
});
