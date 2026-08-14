import { Link } from "react-router";
import type { SessionMeta } from "../api";

export function SessionHeader({ meta }: { meta: SessionMeta }) {
  return (
    <div className="flex min-w-0 items-center gap-2">
      <Link
        to="/work-items"
        className="shrink-0 text-sm text-zinc-400 transition-colors hover:text-zinc-200 light:text-zinc-500 light:hover:text-zinc-800"
        title="Back to work items"
      >
        ←
      </Link>
      <span className="inline-flex shrink-0 items-center gap-1.5 rounded bg-zinc-800 px-2 py-0.5 font-mono text-xs text-zinc-300 light:bg-zinc-200 light:text-zinc-700">
        {meta.provider}
      </span>
      <code
        className="min-w-0 truncate font-mono text-xs text-zinc-500 light:text-zinc-500"
        title={meta.session_id}
      >
        {meta.session_id}
      </code>
    </div>
  );
}
