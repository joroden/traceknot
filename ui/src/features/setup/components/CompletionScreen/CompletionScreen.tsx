import { CheckCircle2, ExternalLink } from "lucide-react";
import { CopilotRestartNotice } from "../../../settings/components/CopilotRestartNotice";

export interface CompletionScreenProps {
  copilotEnabled: boolean;
}

export function CompletionScreen({ copilotEnabled }: CompletionScreenProps) {
  const openDashboard = () => {
    window.open(window.location.origin + "/dashboard", "_blank");
    window.close();
  };

  return (
    <div className="fixed inset-0 z-40 flex flex-col items-center justify-center gap-2 bg-zinc-950/95 p-6 text-center">
      <span className="mb-1.5 inline-flex rounded-full border border-emerald-500/35 p-3.5 text-emerald-500">
        <CheckCircle2 size={26} />
      </span>
      <h2 className="text-lg font-bold">You're all set</h2>
      <p className="mb-3.5 max-w-[380px] text-sm text-zinc-400">
        traceknot is running in the background and ready to collect telemetry from your agents.
      </p>
      {copilotEnabled ? (
        <div className="mb-3.5 w-full max-w-[420px]">
          <CopilotRestartNotice />
        </div>
      ) : null}
      <div className="flex items-center gap-2.5">
        <button
          type="button"
          className="cursor-pointer rounded-lg border border-zinc-700 bg-zinc-800 px-4 py-2 text-sm font-semibold text-zinc-100 transition-colors hover:border-violet-500 light:border-zinc-300 light:bg-zinc-100 light:text-zinc-900"
          onClick={() => window.close()}
        >
          Close
        </button>
        <button
          type="button"
          className="inline-flex cursor-pointer items-center gap-1.5 rounded-lg border border-violet-600 bg-violet-600 px-4 py-2 text-sm font-semibold text-white transition-colors hover:bg-violet-500"
          onClick={openDashboard}
        >
          <ExternalLink size={14} />
          Open Dashboard
        </button>
      </div>
    </div>
  );
}
