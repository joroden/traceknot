import { Link } from "react-router";
import { Badge } from "../../../components/Badge";
import { formatProviderLabel } from "../../../utils/format";
import { providerBadgeTone } from "../../../utils/providers";
import type { SessionMeta } from "../api";

export function SessionHeader({ meta }: { meta: SessionMeta }) {
  return (
    <div className="flex min-w-0 shrink-0 items-center gap-3">
      <Link
        to="/work-items"
        className="shrink-0 rounded-md p-1 text-zinc-400 transition-colors hover:bg-zinc-800 hover:text-zinc-100 light:text-zinc-500 light:hover:bg-zinc-200 light:hover:text-zinc-800"
        title="Back to work items"
      >
        ←
      </Link>
      <Badge tone={providerBadgeTone(meta.provider)}>{formatProviderLabel(meta.provider)}</Badge>
      <code
        className="min-w-0 truncate font-mono text-sm text-zinc-300 light:text-zinc-600"
        title={meta.session_id}
      >
        {meta.session_id}
      </code>
    </div>
  );
}
