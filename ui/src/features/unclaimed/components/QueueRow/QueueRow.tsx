import { Badge } from "../../../../components/Badge";
import { formatDuration, formatProviderLabel, formatTimestamp, formatUSD } from "../../../../utils/format";
import { providerBadgeTone } from "../../../../utils/providers";
import type { UnclaimedSession } from "../../api";
import { ModelBadges } from "../../../../components/ModelBadges";
import { StatusBadge } from "../StatusBadge";

export interface QueueRowProps {
  session: UnclaimedSession;
  onOpenSession: (sessionId: string) => void;
  onFindWorkItem: () => void;
}

export function QueueRow({ session, onOpenSession, onFindWorkItem }: QueueRowProps) {
  return (
    <div className="flex items-center gap-3.5 rounded-lg border border-zinc-800 bg-zinc-900 p-3.5 light:border-zinc-200 light:bg-white">
      <div className="flex min-w-0 flex-1 flex-col gap-1.5">
        <button
          type="button"
          onClick={() => onOpenSession(session.session_id)}
          className="min-w-0 cursor-pointer truncate bg-transparent p-0 text-left text-sm font-medium text-zinc-100 hover:text-violet-400 light:text-zinc-900 light:hover:text-violet-600"
        >
          {session.title || "Untitled"}
        </button>
        <div className="flex flex-wrap items-center gap-1.5">
          <Badge tone={providerBadgeTone(session.provider)}>{formatProviderLabel(session.provider)}</Badge>
          <StatusBadge status={session.claim.status} />
        </div>
        <div className="flex flex-wrap items-center gap-3">
          <ModelBadges models={session.models} />
          <div className="flex items-center gap-3 text-xs text-zinc-400 light:text-zinc-500">
            <span>{formatTimestamp(session.started_at_unix_ms)}</span>
            <span className="font-mono">{formatDuration(session.duration_ms)}</span>
            <span className="font-mono font-semibold text-zinc-100 light:text-zinc-900">
              {formatUSD(session.cost)}
            </span>
          </div>
        </div>
      </div>

      <button
        type="button"
        onClick={onFindWorkItem}
        className="flex-shrink-0 cursor-pointer rounded-md border border-violet-600 bg-violet-600 px-3.5 py-1.5 text-sm font-semibold text-white transition-colors hover:bg-violet-500"
      >
        Claim
      </button>
    </div>
  );
}
