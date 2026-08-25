import { Search } from "lucide-react";

interface TreeToolbarProps {
  query: string;
  onQueryChange: (query: string) => void;
  visibleCount: number;
  onExpandAll: () => void;
  onCollapseAll: () => void;
}

export function TreeToolbar({
  query,
  onQueryChange,
  visibleCount,
  onExpandAll,
  onCollapseAll,
}: TreeToolbarProps) {
  return (
    <div className="flex min-w-0 flex-1 flex-wrap items-center gap-3">
      {!query.trim() && (
        <div className="flex items-center gap-1 text-xs">
          <button
            type="button"
            onClick={onExpandAll}
            className="cursor-pointer rounded px-2 py-1 text-zinc-400 transition-colors hover:bg-zinc-800 hover:text-zinc-200 light:text-zinc-500 light:hover:bg-zinc-200 light:hover:text-zinc-800"
          >
            Expand all
          </button>
          <button
            type="button"
            onClick={onCollapseAll}
            className="cursor-pointer rounded px-2 py-1 text-zinc-400 transition-colors hover:bg-zinc-800 hover:text-zinc-200 light:text-zinc-500 light:hover:bg-zinc-200 light:hover:text-zinc-800"
          >
            Collapse all
          </button>
        </div>
      )}

      <div className="relative ml-auto min-w-52 flex-1 sm:flex-none">
        <Search className="pointer-events-none absolute left-2 top-1/2 h-3.5 w-3.5 -translate-y-1/2 text-zinc-500" />
        <input
          type="search"
          value={query}
          onChange={(event) => onQueryChange(event.target.value)}
          placeholder="Search tools, agents, names…"
          className="w-full rounded-lg border border-zinc-700 bg-zinc-900 py-1.5 pl-7 pr-2 text-xs text-zinc-200 outline-none transition-colors placeholder:text-zinc-500 focus:border-violet-500 light:border-zinc-300 light:bg-white light:text-zinc-800"
        />
      </div>

      <span className="text-xs tabular-nums text-zinc-500 light:text-zinc-500">
        {visibleCount} rows
      </span>
    </div>
  );
}
