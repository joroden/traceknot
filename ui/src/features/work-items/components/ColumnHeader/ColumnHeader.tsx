import type { MouseEvent, ReactNode } from "react";
import { ArrowDown, ArrowUp, ChevronsUpDown } from "lucide-react";
import type { SortDir } from "../../api";

export interface ColumnHeaderProps {
  label: string;
  sortActive?: boolean;
  sortDir?: SortDir;
  sortPriority?: number;
  onSortClick?: (event: MouseEvent) => void;
  filter?: ReactNode;
}

export function ColumnHeader({
  label,
  sortActive,
  sortDir,
  sortPriority,
  onSortClick,
  filter,
}: ColumnHeaderProps) {
  return (
    <div className="flex items-center gap-1.5">
      {onSortClick ? (
        <button
          type="button"
          onClick={onSortClick}
          title="Click to sort, shift-click to add to multi-sort"
          className={
            sortActive
              ? "flex cursor-pointer items-center gap-1 text-zinc-100 light:text-zinc-900"
              : "flex cursor-pointer items-center gap-1 text-zinc-400 hover:text-zinc-200 light:text-zinc-500 light:hover:text-zinc-700"
          }
        >
          {label}
          {sortActive ? (
            sortDir === "asc" ? (
              <ArrowUp className="h-3 w-3" />
            ) : (
              <ArrowDown className="h-3 w-3" />
            )
          ) : (
            <ChevronsUpDown className="h-3 w-3 opacity-40" />
          )}
          {sortPriority ? <span className="text-zinc-500">{sortPriority}</span> : null}
        </button>
      ) : (
        <span>{label}</span>
      )}
      {filter}
    </div>
  );
}
