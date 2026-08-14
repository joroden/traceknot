import { ChevronRight } from "lucide-react";

export interface RankedBarRow {
  label: string;
  sublabel?: string;
  value: number;
  display: string;
  onSelect?: () => void;
  barClassName?: string;
}

export interface RankedBarsProps {
  rows: RankedBarRow[];
  maxValue?: number;
  initialLimit?: number;
  onShowAll?: () => void;
}

const RANK_SHADES = [
  "bg-violet-500",
  "bg-violet-500/78",
  "bg-violet-500/60",
  "bg-violet-500/45",
  "bg-violet-500/32",
];

function rankShade(index: number): string {
  return RANK_SHADES[Math.min(index, RANK_SHADES.length - 1)];
}

export function RankedBars({ rows, maxValue, initialLimit = 3, onShowAll }: RankedBarsProps) {
  const scale = maxValue ?? rows.reduce((max, row) => Math.max(max, row.value), 0);
  const visible = rows.length > initialLimit ? rows.slice(0, initialLimit) : rows;

  return (
    <div className="flex flex-col gap-1">
      {visible.map((row, index) => {
        const width = scale > 0 ? (row.value / scale) * 100 : 0;
        const content = (
          <>
            <div className="flex min-w-0 items-baseline justify-between gap-3">
              <span className="min-w-0 truncate font-mono text-xs font-semibold text-zinc-100 light:text-zinc-900">
                {row.label}
              </span>
              <span className="flex-shrink-0 font-mono text-xs text-zinc-100 light:text-zinc-900">
                {row.display}
              </span>
            </div>
            {row.sublabel ? (
              <p className="mt-0.5 truncate text-xs text-zinc-400 light:text-zinc-500">
                {row.sublabel}
              </p>
            ) : null}
            <div className="mt-1.5 h-1.5 rounded-full bg-zinc-800 light:bg-zinc-200">
              <div
                className={`h-1.5 rounded-full ${row.barClassName ?? rankShade(index)}`}
                style={{ width: `${width}%` }}
              />
            </div>
          </>
        );

        if (row.onSelect) {
          return (
            <button
              key={row.label}
              type="button"
              onClick={row.onSelect}
              className="cursor-pointer rounded-md border border-transparent px-2 py-2 text-left transition-colors hover:border-zinc-800 hover:bg-zinc-900 light:hover:border-zinc-200 light:hover:bg-zinc-50"
            >
              {content}
            </button>
          );
        }
        return (
          <div key={row.label} className="rounded-md px-2 py-2">
            {content}
          </div>
        );
      })}
      {rows.length > initialLimit && onShowAll && (
        <button
          type="button"
          onClick={onShowAll}
          className="inline-flex cursor-pointer items-center gap-1 self-start rounded px-2 py-1 text-xs text-zinc-500 transition-colors hover:bg-zinc-800 hover:text-zinc-300 light:hover:bg-zinc-100 light:hover:text-zinc-600"
        >
          {`Show all (${rows.length - initialLimit} more)`}
          <ChevronRight className="h-3 w-3" />
        </button>
      )}
    </div>
  );
}
