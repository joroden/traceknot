import { useRef } from "react";
import { useVirtualizer } from "@tanstack/react-virtual";
import type { TreeDerived, VisibleRow } from "../hooks/useSessionDetail";
import { EmptyState } from "../../../components/EmptyState";
import { COLUMN_TITLES, GRID_TEMPLATE } from "./columns";
import { NodeRow } from "./NodeRow";

interface TreeTableProps {
  rows: VisibleRow[];
  derived: TreeDerived;
  expanded: Set<string>;
  selectedId: string | null;
  onToggle: (nodeId: string) => void;
  onSelect: (nodeId: string) => void;
}

const ROW_HEIGHT = 36;

export function TreeTable({
  rows,
  derived,
  expanded,
  selectedId,
  onToggle,
  onSelect,
}: TreeTableProps) {
  const scrollRef = useRef<HTMLDivElement>(null);
  const virtualizer = useVirtualizer({
    count: rows.length,
    getScrollElement: () => scrollRef.current,
    estimateSize: () => ROW_HEIGHT,
    overscan: 20,
  });

  if (rows.length === 0) {
    return (
      <div className="flex h-full min-h-64 min-w-0 flex-1 items-center justify-center rounded-lg border border-zinc-800 light:border-zinc-200">
        <EmptyState title="No nodes match — adjust the search or clear the focus." />
      </div>
    );
  }

  return (
    <div className="flex min-h-0 min-w-0 flex-1 flex-col overflow-hidden rounded-lg border border-zinc-800 light:border-zinc-200">
      <div ref={scrollRef} className="min-h-0 flex-1 overflow-auto">
        <div
          className="sticky top-0 z-10 grid h-9 shrink-0 items-center gap-x-2 border-b border-zinc-800 bg-zinc-900 pr-2 text-xs font-semibold uppercase tracking-wide text-zinc-500 light:border-zinc-200 light:bg-zinc-100 light:text-zinc-500"
          style={{ gridTemplateColumns: GRID_TEMPLATE }}
        >
          {COLUMN_TITLES.map((column) => (
            <div
              key={column.key}
              title={column.tooltip}
              className={
                column.key === "name"
                  ? "sticky left-0 z-10 bg-zinc-900 pl-3 light:bg-zinc-100"
                  : column.headerClassName
              }
            >
              {column.title}
            </div>
          ))}
        </div>
        <div style={{ height: virtualizer.getTotalSize(), position: "relative" }}>
          {virtualizer.getVirtualItems().map((item) => {
            const entry = rows[item.index];
            const row = entry.row;
            const hasChildren = (derived.childrenOf.get(row.node_id) ?? []).length > 0;
            return (
              <div
                key={row.node_id}
                style={{
                  position: "absolute",
                  top: 0,
                  left: 0,
                  width: "100%",
                  transform: `translateY(${item.start}px)`,
                }}
              >
                <NodeRow
                  row={row}
                  depth={entry.depth}
                  hasChildren={hasChildren}
                  expanded={expanded.has(row.node_id)}
                  selected={selectedId === row.node_id}
                  onToggle={onToggle}
                  onSelect={onSelect}
                />
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
}
