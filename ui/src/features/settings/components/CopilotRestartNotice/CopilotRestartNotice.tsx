import { useState } from "react";
import { RotateCw } from "lucide-react";
import { restartSessions } from "../../api";

export interface CopilotRestartNoticeProps {
  actionable?: boolean;
}

export function CopilotRestartNotice({ actionable = true }: CopilotRestartNoticeProps) {
  const [running, setRunning] = useState(false);
  const [result, setResult] = useState<string | null>(null);

  const restart = async () => {
    setRunning(true);
    setResult(null);
    try {
      const { closed } = await restartSessions();
      setResult(closed === 0 ? "Nothing running to close." : `Closed ${closed} running session(s).`);
    } catch (reason) {
      setResult(reason instanceof Error ? reason.message : String(reason));
    } finally {
      setRunning(false);
    }
  };

  if (!actionable) {
    return (
      <div className="rounded-lg border border-amber-500/40 bg-amber-500/10 px-4 py-3 text-sm text-zinc-100 light:text-zinc-900">
        Restart VS Code and Copilot CLI after finishing setup for this to take effect.
      </div>
    );
  }

  return (
    <div className="flex items-center justify-between gap-4 rounded-lg border border-amber-500/40 bg-amber-500/10 px-4 py-3 text-sm">
      <p className="text-zinc-100 light:text-zinc-900">
        {result ?? "Restart VS Code and Copilot CLI for this to take effect."}
      </p>
      <button
        type="button"
        className="inline-flex flex-shrink-0 cursor-pointer items-center gap-1.5 rounded-lg border border-amber-500 px-3 py-1.5 text-xs font-semibold text-amber-500 transition-colors hover:bg-amber-500/10 disabled:cursor-default disabled:opacity-60"
        onClick={() => void restart()}
        disabled={running}
      >
        <RotateCw size={12} />
        {running ? "Restarting…" : "Restart now"}
      </button>
    </div>
  );
}
