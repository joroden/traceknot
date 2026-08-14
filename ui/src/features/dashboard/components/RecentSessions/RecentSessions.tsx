import { Link } from "react-router";
import { Badge } from "../../../../components/Badge";
import { formatProviderLabel, formatRelativeTime, formatUSD } from "../../../../utils/format";
import { providerBadgeTone } from "../../../../utils/providers";
import type { RecentSession } from "../../api";

export interface RecentSessionsProps {
  sessions: RecentSession[];
}

function rowTitle(session: RecentSession): string {
  if (session.title) {
    return session.title;
  }
  return `${formatProviderLabel(session.provider)} session · ${session.node_count} nodes`;
}

export function RecentSessions({ sessions }: RecentSessionsProps) {
  if (sessions.length === 0) {
    return null;
  }
  return (
    <div className="flex flex-col gap-1">
      {sessions.map((session) => (
        <Link
          key={session.session_id}
          to={`/sessions/${encodeURIComponent(session.session_id)}`}
          className="flex items-center gap-3 rounded-md border border-transparent px-2 py-2 no-underline transition-colors hover:border-zinc-800 hover:bg-zinc-900 light:hover:border-zinc-200 light:hover:bg-zinc-50"
        >
          <Badge tone={providerBadgeTone(session.provider)}>{formatProviderLabel(session.provider)}</Badge>
          <span className="min-w-0 flex-1 truncate text-sm text-zinc-100 light:text-zinc-900">
            {rowTitle(session)}
          </span>
          <span className="hidden shrink-0 text-xs text-zinc-500 sm:inline">
            {formatRelativeTime(session.started_at_unix_ms)}
          </span>
          <span className="flex-shrink-0 font-mono text-xs font-semibold text-zinc-100 light:text-zinc-900">
            {formatUSD(session.cost)}
          </span>
        </Link>
      ))}
    </div>
  );
}
