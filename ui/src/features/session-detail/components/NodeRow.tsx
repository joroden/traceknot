import { ChevronRight, ChevronDown } from "lucide-react";
import type { TreeNodeRow } from "../api";
import type { MergedRow } from "../hooks/useSessionDetail";
import { formatDuration, formatUSD } from "../../../utils/format";
import { GRID_TEMPLATE } from "./columns";
import { kindStyle } from "./KindBadge";
import { StatusDot } from "./StatusDot";
import { TokenValue } from "./TokenValue";

interface NodeRowProps {
  row: MergedRow;
  depth: number;
  hasChildren: boolean;
  expanded: boolean;
  selected: boolean;
  onToggle: (nodeId: string) => void;
  onSelect: (nodeId: string) => void;
}

function labelFor(row: TreeNodeRow): string {
  if (row.kind === "tool_call") {
    return row.tool_name ?? row.name ?? "Tool call";
  }
  if (row.kind === "agent") {
    return row.agent_name ?? row.name ?? "Agent";
  }
  if (row.kind === "chat") {
    return row.name === "user" ? "Prompt" : "Chat";
  }
  return row.name ?? row.kind;
}

export function NodeRow({
  row,
  depth,
  hasChildren,
  expanded,
  selected,
  onToggle,
  onSelect,
}: NodeRowProps) {
  const isSubagent = row.launchNodeId !== null || (row.kind === "agent" && row.tool_name === "Agent");
  const isPrompt = row.kind === "chat" && row.name === "user";
  const displayKind = isSubagent ? "subagent" : isPrompt ? "prompt" : row.kind;
  const style = kindStyle(displayKind);
  const Icon = style.icon;
  const indent = depth * 16;

  return (
    <div
      className={`grid h-9 cursor-pointer items-center border-b border-zinc-800/60 pr-2 text-sm transition-colors last:border-b-0 hover:bg-zinc-900/70 light:border-zinc-200/60 light:hover:bg-zinc-100/70 ${selected ? "bg-violet-500/10 hover:bg-violet-500/10 light:bg-violet-500/10" : ""}`}
      style={{ gridTemplateColumns: GRID_TEMPLATE }}
      onClick={() => onSelect(row.node_id)}
      role="button"
      tabIndex={-1}
    >
      <div
        className="sticky left-0 z-[1] flex min-w-0 items-center gap-1.5 bg-zinc-950 py-1 pr-2 light:bg-zinc-50"
        style={{ paddingLeft: 8 + indent }}
      >
        {hasChildren ? (          <button
            type="button"
            className="shrink-0 cursor-pointer text-zinc-500 hover:text-zinc-200 light:hover:text-zinc-800"
            onClick={(event) => {
              event.stopPropagation();
              onToggle(row.node_id);
            }}
            aria-label={expanded ? "Collapse" : "Expand"}
          >
            {expanded ? (
              <ChevronDown className="h-3.5 w-3.5" />
            ) : (
              <ChevronRight className="h-3.5 w-3.5" />
            )}
          </button>
        ) : (
          <span className="w-3.5 shrink-0" />
        )}
        <span
          className={`inline-flex shrink-0 items-center rounded px-1 py-0.5 ${style.className}`}
          title={style.label}
        >
          <Icon className="h-3 w-3" />
        </span>
        <span className="min-w-0 truncate text-zinc-200 light:text-zinc-800">
          {labelFor(row)}
        </span>
        {row.has_content && (
          <span
            className="h-1.5 w-1.5 shrink-0 rounded-full bg-emerald-500/70"
            title="Content captured for this node"
          />
        )}
      </div>
      <div className="flex items-center justify-center">
        <StatusDot status={row.status} />
      </div>
      <div className="flex items-center justify-end pr-3 font-mono text-xs tabular-nums text-zinc-400 light:text-zinc-500">
        {row.subagent_count > 0 ? row.subagent_count : "—"}
      </div>
      <div className="flex items-center pr-2 text-xs text-zinc-400 light:text-zinc-500">
        <span className="truncate">{row.model ?? "—"}</span>
      </div>
      <div className="flex items-center justify-end pr-3 font-mono text-xs tabular-nums text-zinc-400 light:text-zinc-500">
        {formatDuration(row.duration_ms)}
      </div>
      <TokenValue
        aggregated={row.agg_input_tokens}
        self={row.input_tokens}
        estimated={row.estimated_input_tokens}
        title="Input"
      />
      <TokenValue
        aggregated={row.agg_cached_input_tokens}
        self={row.cached_input_tokens}
        estimated={0}
        title="Cached"
      />
      <TokenValue
        aggregated={row.agg_cache_write_tokens}
        self={row.cache_write_tokens}
        estimated={0}
        title="Cache write"
      />
      <TokenValue
        aggregated={row.agg_output_tokens}
        self={row.output_tokens}
        estimated={row.estimated_output_tokens}
        title="Output"
      />
      <TokenValue
        aggregated={row.agg_reasoning_tokens}
        self={row.reasoning_tokens}
        estimated={0}
        title="Reasoning"
      />
      <div className="flex items-center justify-end pr-1 font-mono text-xs tabular-nums text-zinc-300 light:text-zinc-600">
        {row.agg_cost > 0 || row.cost > 0 ? formatUSD(row.agg_cost > 0 ? row.agg_cost : row.cost) : "—"}
      </div>
    </div>
  );
}
