import type { ReactNode } from "react";
import type { WorkItemRow } from "../../types/workItem";
import { formatRelativeTime } from "../../utils/format";
import { EmptyState } from "../EmptyState";
import { BrandIcon } from "../BrandIcon";
import { LoadingList } from "./LoadingList";

export interface ResultListProps {
  rows: WorkItemRow[];
  loading?: boolean;
  error?: string | null;
  highlighted: number;
  onHighlight: (index: number) => void;
  onSelect: (item: WorkItemRow) => void;
  emptyTitle: string;
  emptyAction?: ReactNode;
}

export function ResultList({
  rows,
  loading = false,
  error = null,
  highlighted,
  onHighlight,
  onSelect,
  emptyTitle,
  emptyAction,
}: ResultListProps) {
  if (loading) {
    return <LoadingList />;
  }
  if (error) {
    return (
      <EmptyState
        title="Search failed"
        action={<span className="text-xs text-red-500">{error}</span>}
      />
    );
  }
  if (rows.length === 0) {
    return <EmptyState title={emptyTitle} action={emptyAction} />;
  }

  return (
    <div className="flex flex-1 flex-col gap-1.5 overflow-y-auto pr-1">
      {rows.map((item, index) => (
        <button
          key={`${item.provider}:${item.key}`}
          type="button"
          className={
            index === highlighted
              ? "flex w-full cursor-pointer items-center gap-2.5 rounded-lg border border-violet-500 bg-zinc-800 px-3.5 py-2.5 text-sm text-zinc-100 transition-colors light:bg-zinc-100 light:text-zinc-900"
              : "flex w-full cursor-pointer items-center gap-2.5 rounded-lg border border-transparent bg-transparent px-3.5 py-2.5 text-sm text-zinc-100 transition-colors hover:border-zinc-700 hover:bg-zinc-800 light:text-zinc-900 light:hover:border-zinc-300 light:hover:bg-zinc-100"
          }
          onMouseEnter={() => onHighlight(index)}
          onClick={() => onSelect(item)}
        >
          <span className="inline-flex flex-shrink-0 text-zinc-400 light:text-zinc-500">
            <BrandIcon provider={item.provider} />
          </span>
          <span className="flex-shrink-0 font-mono text-xs font-semibold text-violet-400 light:text-violet-600">
            {item.key}
          </span>
          <span className="min-w-0 flex-1 truncate">{item.title}</span>
          <span className="flex-shrink-0 text-xs text-zinc-500 light:text-zinc-400">
            {formatRelativeTime(item.updatedAtUnixMs)}
          </span>
        </button>
      ))}
    </div>
  );
}
