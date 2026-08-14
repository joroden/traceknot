import { useEffect, useRef } from "react";
import type { ReactNode } from "react";
import { flexRender } from "@tanstack/react-table";
import { getCoreRowModel, useLegacyTable, type LegacyColumnDef } from "@tanstack/react-table/legacy";
import { useVirtualizer } from "@tanstack/react-virtual";

export type TableColumn<T extends Record<string, unknown>> = LegacyColumnDef<T>;

export interface TableProps<T extends Record<string, unknown>> {
  columns: TableColumn<T>[];
  data: T[];
  getRowId: (row: T) => string;
  onRowClick?: (row: T) => void;
  onEndReached?: (row: T) => void;
  rowHeightPx?: number;
  emptyState?: ReactNode;
  meta?: Record<string, unknown>;
  growColumnIds?: string[];
  rowClassName?: (row: T) => string;
}

const END_REACHED_THRESHOLD = 20;

const DEFAULT_ROW_CLASS =
  "flex items-center border-b border-zinc-900 light:border-zinc-100";
const CLICKABLE_ROW_CLASS =
  "flex cursor-pointer items-center border-b border-zinc-900 transition-colors hover:bg-zinc-900 light:border-zinc-100 light:hover:bg-zinc-50";

export function Table<T extends Record<string, unknown>>({
  columns,
  data,
  getRowId,
  onRowClick,
  onEndReached,
  rowHeightPx = 44,
  emptyState,
  meta,
  growColumnIds,
  rowClassName,
}: TableProps<T>) {
  const table = useLegacyTable({
    data,
    columns,
    getRowId,
    getCoreRowModel: getCoreRowModel(),
    meta,
  });

  const rows = table.getRowModel().rows;
  const scrollRef = useRef<HTMLDivElement>(null);

  const virtualizer = useVirtualizer({
    count: rows.length,
    getScrollElement: () => scrollRef.current,
    estimateSize: () => rowHeightPx,
    overscan: 12,
  });
  const virtualRows = virtualizer.getVirtualItems();

  useEffect(() => {
    const lastVisible = virtualRows[virtualRows.length - 1];
    if (!lastVisible || !onEndReached) {
      return;
    }
    if (lastVisible.index >= rows.length - END_REACHED_THRESHOLD) {
      onEndReached(rows[lastVisible.index].original);
    }
  }, [virtualRows, rows.length, onEndReached]);

  if (rows.length === 0 && emptyState) {
    return <>{emptyState}</>;
  }

  const totalWidth = table.getTotalSize();
  const cellFlex = (columnId: string, size: number) =>
    growColumnIds?.includes(columnId) ? `1 1 ${size}px` : `0 0 ${size}px`;

  return (
    <div
      ref={scrollRef}
      role="table"
      className="h-full overflow-auto rounded-lg border border-zinc-800 light:border-zinc-200"
    >
      <div style={{ minWidth: totalWidth }}>
        <div role="rowgroup" className="sticky top-0 z-10 bg-zinc-950 light:bg-zinc-50">
          {table.getHeaderGroups().map((headerGroup) => (
            <div role="row" key={headerGroup.id} className="flex w-full">
              {headerGroup.headers.map((header) => (
                <div
                  role="columnheader"
                  key={header.id}
                  style={{ flex: cellFlex(header.column.id, header.getSize()) }}
                  className="min-w-0 border-b border-zinc-800 px-3 py-2 text-left align-bottom text-xs font-medium text-zinc-400 light:border-zinc-200 light:text-zinc-500"
                >
                  {header.isPlaceholder
                    ? null
                    : flexRender(header.column.columnDef.header, header.getContext())}
                </div>
              ))}
            </div>
          ))}
        </div>

        <div role="rowgroup" style={{ height: virtualizer.getTotalSize(), position: "relative" }}>
          {virtualRows.map((virtualRow) => {
            const row = rows[virtualRow.index];
            return (
              <div
                role="row"
                key={row.id}
                onClick={onRowClick ? () => onRowClick(row.original) : undefined}
                style={{
                  position: "absolute",
                  top: 0,
                  left: 0,
                  width: "100%",
                  height: virtualRow.size,
                  transform: `translateY(${virtualRow.start}px)`,
                }}
                className={`${onRowClick ? CLICKABLE_ROW_CLASS : DEFAULT_ROW_CLASS} ${rowClassName?.(row.original) ?? ""}`}
              >
                {row.getVisibleCells().map((cell) => (
                  <div
                    role="cell"
                    key={cell.id}
                    style={{ flex: cellFlex(cell.column.id, cell.column.getSize()) }}
                    className="min-w-0 truncate px-3 py-2 text-zinc-200 light:text-zinc-800"
                  >
                    {flexRender(cell.column.columnDef.cell, cell.getContext())}
                  </div>
                ))}
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
}
