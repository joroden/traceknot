import { useEffect, useState } from "react";
import { X } from "lucide-react";
import { useNodeDetail } from "../hooks/useNodeDetail";
import type { MergedRow } from "../hooks/useSessionDetail";
import { formatDuration, formatTimestamp, formatTokens, formatUSD } from "../../../utils/format";
import { AgentChip } from "./AgentChip";
import { buildTabs } from "./detailTabs";
import { DetailTabContent } from "./DetailTabContent";
import { KindBadge } from "./KindBadge";
import { MetaRow } from "./MetaRow";

interface DetailDrawerProps {
  nodeId: string;
  launchNodeId: string | null;
  treeRow: MergedRow | null;
  onClose: () => void;
}

export function DetailDrawer({ nodeId, launchNodeId, treeRow, onClose }: DetailDrawerProps) {
  const { detail, launchDetail, loading, error } = useNodeDetail(nodeId, launchNodeId);
  const [activeTab, setActiveTab] = useState("");

  useEffect(() => {
    setActiveTab("");
  }, [nodeId]);

  if (!nodeId) {
    return null;
  }

  const inputTokens =
    treeRow && treeRow.agg_input_tokens > 0 ? treeRow.agg_input_tokens : (detail?.input_tokens ?? 0);
  const outputTokens =
    treeRow && treeRow.agg_output_tokens > 0
      ? treeRow.agg_output_tokens
      : (detail?.output_tokens ?? 0);
  const cachedTokens =
    treeRow && treeRow.agg_cached_input_tokens > 0
      ? treeRow.agg_cached_input_tokens
      : (detail?.cached_input_tokens ?? 0);

  const isSubagent =
    launchNodeId !== null || (detail?.kind === "agent" && detail.tool_name === "Agent");
  const isPrompt = detail?.kind === "chat" && detail.name === "user";
  const isMeta = detail?.kind === "chat" && detail.name === "meta";
  const displayKind = isSubagent
    ? "subagent"
    : isPrompt
      ? "prompt"
      : isMeta
        ? "meta"
        : (detail?.kind ?? "");

  let title = "Node";
  if (detail) {
    if (detail.kind === "tool_call") {
      title = detail.tool_name ?? detail.name ?? "Tool call";
    } else if (detail.kind === "agent") {
      title = detail.agent_name ?? detail.name ?? "Agent";
    } else {
      title = detail.name ?? detail.kind;
    }
  }

  const tabs = detail ? buildTabs(detail, isSubagent, launchDetail) : [];
  const active = activeTab || tabs[0]?.key || "";

  return (
    <aside className="flex w-[460px] shrink-0 flex-col overflow-hidden rounded-lg border border-zinc-800 bg-zinc-900 light:border-zinc-200 light:bg-white">
      <div className="flex items-center gap-2 border-b border-zinc-800 p-3 light:border-zinc-200">
        <KindBadge kind={displayKind} />
        <h2 className="min-w-0 flex-1 truncate font-mono text-sm font-semibold text-zinc-100 light:text-zinc-900">
          {title}
        </h2>
        <button
          type="button"
          onClick={onClose}
          className="cursor-pointer rounded p-1 text-zinc-500 transition-colors hover:bg-zinc-800 hover:text-zinc-200 light:hover:bg-zinc-100 light:hover:text-zinc-800"
          aria-label="Close detail"
        >
          <X className="h-4 w-4" />
        </button>
      </div>

      <div className="flex min-h-0 flex-1 flex-col overflow-hidden p-3">
        {loading && <p className="text-xs text-zinc-500">Loading…</p>}
        {error && <p className="text-xs text-red-400">{error}</p>}
        {detail && (
          <div className="flex min-h-0 flex-1 flex-col gap-3">
            <div className="shrink-0 space-y-1.5">
              {detail.kind === "tool_call" && !detail.model ? (
                <MetaRow label="Est. tokens">
                  {detail.estimated_input_tokens > 0 || detail.estimated_output_tokens > 0 ? (
                    <span title={detail.token_estimate_method ?? undefined}>
                      ~{formatTokens(detail.estimated_input_tokens)} in · ~
                      {formatTokens(detail.estimated_output_tokens)} out
                    </span>
                  ) : (
                    "—"
                  )}
                </MetaRow>
              ) : (
                <>
                  <MetaRow label="Model">{detail.model ?? "—"}</MetaRow>
                  <MetaRow label="Cost">
                    {treeRow && (treeRow.agg_cost > 0 || treeRow.cost > 0)
                      ? formatUSD(treeRow.agg_cost > 0 ? treeRow.agg_cost : treeRow.cost)
                      : "—"}
                  </MetaRow>
                  <MetaRow label="Tokens">
                    {formatTokens(inputTokens)} in · {formatTokens(outputTokens)} out
                    {cachedTokens > 0 ? ` · ${formatTokens(cachedTokens)} cached` : ""}
                  </MetaRow>
                </>
              )}
              <MetaRow label="Status">{detail.status ?? "—"}</MetaRow>
              <MetaRow label="Duration">{formatDuration(detail.duration_ms)}</MetaRow>
              <MetaRow label="Started">{formatTimestamp(detail.started_at_unix_ms)}</MetaRow>
              {detail.owning_agent_id && (
                <MetaRow label="Agent">
                  <AgentChip name={detail.owning_agent_name} />
                </MetaRow>
              )}
            </div>

            {tabs.length > 0 ? (
              <>
                <div className="flex shrink-0 gap-4 border-b border-zinc-800 light:border-zinc-200">
                  {tabs.map((tab) => (
                    <button
                      key={tab.key}
                      type="button"
                      onClick={() => setActiveTab(tab.key)}
                      className={`-mb-px cursor-pointer border-b-2 pb-1.5 text-xs font-medium transition-colors ${
                        tab.key === active
                          ? "border-violet-500 text-zinc-100 light:text-zinc-900"
                          : "border-transparent text-zinc-500 hover:text-zinc-300 light:hover:text-zinc-600"
                      }`}
                    >
                      {tab.label}
                    </button>
                  ))}
                </div>
                <div className="flex min-h-0 flex-1 flex-col overflow-hidden">
                  <DetailTabContent
                    detail={detail}
                    active={active}
                    launchDetail={launchDetail}
                    isSubagent={isSubagent}
                  />
                </div>
              </>
            ) : (
              <p className="shrink-0 text-xs text-zinc-500 light:text-zinc-500">
                No content captured for this node — the provider may not export it.
              </p>
            )}

            <p className="shrink-0 break-all font-mono text-xs text-zinc-600 light:text-zinc-400">
              {detail.node_id}
            </p>
          </div>
        )}
      </div>
    </aside>
  );
}
