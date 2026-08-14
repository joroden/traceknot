import { CheckCircle2, Forward } from "lucide-react";
import type { WorkItemRow } from "../../../../types/workItem";

export type Outcome = {
  kind: "claimed" | "skipped";
  item: WorkItemRow | null;
} | null;

export interface OutcomeOverlayProps {
  outcome: Outcome;
}

export function OutcomeOverlay({ outcome }: OutcomeOverlayProps) {
  if (!outcome) {
    return null;
  }
  const claimed = outcome.kind === "claimed";
  return (
    <div className="fixed inset-0 z-40 flex flex-col items-center justify-center gap-2 bg-zinc-950/95 p-6 text-center">
      <span
        className={
          claimed
            ? "mb-1.5 inline-flex rounded-full border border-emerald-500/35 p-3.5 text-emerald-500"
            : "mb-1.5 inline-flex rounded-full border border-zinc-700 p-3.5 text-zinc-400"
        }
      >
        {claimed ? <CheckCircle2 size={26} /> : <Forward size={26} />}
      </span>
      <h2 className="text-lg font-bold">
        {claimed ? `Attached ${outcome.item?.key}` : "Context skipped"}
      </h2>
      <p className="mb-3.5 max-w-[420px] font-mono text-xs text-zinc-400">
        {claimed
          ? outcome.item?.title
          : "Proceeding without attributing this session. You can claim it later from the triage queue."}
      </p>
      <button
        type="button"
        className="cursor-pointer rounded-lg border border-violet-600 bg-violet-600 px-4 py-2 text-sm font-semibold text-white transition-colors hover:bg-violet-500"
        onClick={() => window.close()}
      >
        Close this window
      </button>
      <p className="mt-2.5 max-w-[380px] text-xs text-zinc-500">
        If this window does not close on its own, close it manually — the
        session continues either way.
      </p>
    </div>
  );
}
