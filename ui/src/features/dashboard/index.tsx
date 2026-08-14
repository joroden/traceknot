import { useCallback, useState } from "react";
import { useNavigate } from "react-router";
import { RankedBars } from "../../components/RankedBars";
import { StatCard } from "../../components/StatCard";
import { DaemonUnreachable } from "../../components/DaemonUnreachable";
import { EmptyState } from "../../components/EmptyState";
import { formatCompactCount, formatPct, formatUSD } from "../../utils/format";
import { providerBarClass } from "../../utils/providers";
import type { DashboardRequest, NamedCost, WorkItemCost } from "./api";
import { useDashboard } from "./hooks/useDashboard";
import { FirstRunHero } from "./components/FirstRunHero";
import { RangeSelector, type DashboardRangeRequest } from "./components/RangeSelector";
import { RecentSessions } from "./components/RecentSessions";
import { TrendChart } from "./components/TrendChart";

const PERIOD_LABELS: Record<DashboardRequest["range"], string> = {
  today: "today",
  week: "this week",
  month: "this month",
  all: "all time",
  custom: "in the selected period",
};

export function DashboardPage() {
  const [rangeRequest, setRangeRequest] = useState<DashboardRequest>({
    range: "week",
  });
  const { data, error } = useDashboard(rangeRequest);
  const navigate = useNavigate();

  const handleRange = useCallback((request: DashboardRangeRequest) => {
    setRangeRequest({
      range: request.key,
      ...(request.startMs !== null ? { start_unix_ms: request.startMs } : {}),
      ...(request.endMs !== null ? { end_unix_ms: request.endMs } : {}),
    });
  }, []);

  const openWorkItems = useCallback(
    (filter: {
      provider?: string;
      model?: string;
      workItemKey?: string;
      workItemProvider?: string;
      groupBy?: "work_item" | "none";
    }) => {
      navigate("/work-items", { state: { filter } });
    },
    [navigate],
  );

  const workItemRow = (item: WorkItemCost) => ({
    label: item.key,
    sublabel: item.title,
    value: item.cost,
    display: formatUSD(item.cost),
    onSelect: () => openWorkItems({ workItemKey: item.key, workItemProvider: item.provider }),
  });

  const namedRow = (item: NamedCost) => ({
    label: item.name,
    value: item.cost,
    display: formatUSD(item.cost),
    onSelect: () => openWorkItems({ model: item.name, groupBy: "none" }),
  });

  const agentRow = (item: NamedCost) => ({
    label: item.name,
    value: item.cost,
    display: formatUSD(item.cost),
    onSelect: () => openWorkItems({ provider: item.name, groupBy: "none" }),
    barClassName: providerBarClass(item.name),
  });

  if (error) {
    return <DaemonUnreachable error={error} showStartHint />;
  }

  if (!data) {
    return (
      <div className="flex h-full flex-col gap-6">
        <div className="grid grid-cols-2 gap-4 lg:grid-cols-4">
          {[0, 1, 2, 3].map((index) => (
            <div key={index} className="h-[88px] animate-pulse rounded-lg bg-zinc-900 light:bg-zinc-100" />
          ))}
        </div>
        <div className="h-48 animate-pulse rounded-lg bg-zinc-900 light:bg-zinc-100" />
        <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
          <div className="h-40 animate-pulse rounded-lg bg-zinc-900 light:bg-zinc-100" />
          <div className="h-40 animate-pulse rounded-lg bg-zinc-900 light:bg-zinc-100" />
        </div>
      </div>
    );
  }

  if (data.first_run) {
    return (
      <div className="flex flex-col gap-6">
        <RangeSelector onRequest={handleRange} />
        <FirstRunHero />
      </div>
    );
  }

  const period = data.period;
  const periodLabel = PERIOD_LABELS[rangeRequest.range];

  return (
    <div className="flex flex-col gap-6">
      <RangeSelector onRequest={handleRange} />

      <section className="grid grid-cols-2 gap-4 lg:grid-cols-4">
        <StatCard
          label={`Spent ${periodLabel}`}
          value={formatUSD(period.cost)}
          deltaPct={period.cost_delta_pct}
          goodWhenDown
          note={`${period.session_count} sessions`}
        />
        <StatCard
          label="Tokens"
          value={formatCompactCount(period.tokens)}
          deltaPct={period.tokens_delta_pct}
          goodWhenDown
          note={`${formatCompactCount(period.input_tokens)} in · ${formatCompactCount(period.output_tokens)} out`}
        />
        <StatCard
          label="Coverage"
          value={formatPct(period.coverage_pct)}
          deltaPct={period.coverage_delta_pct}
          note="share of sessions claimed"
        />
        <StatCard
          label="Unattributed"
          value={formatUSD(period.unattributed_cost)}
          tone={period.unattributed_session_count > 0 ? "attention" : "neutral"}
          note={
            period.unattributed_session_count > 0
              ? `${period.unattributed_session_count} sessions to triage`
              : "everything is claimed"
          }
          onSelect={() => navigate("/unclaimed")}
        />
      </section>

      <section className="rounded-lg border border-zinc-800 bg-zinc-900 p-4 light:border-zinc-200 light:bg-white">
        <h2 className="mb-3 text-sm font-semibold">Top work items by cost</h2>
        {data.by_work_item.length > 0 ? (
          <RankedBars
            rows={data.by_work_item.map(workItemRow)}
            maxValue={data.by_work_item[0].cost}
            initialLimit={5}
            onShowAll={() => openWorkItems({ groupBy: "work_item" })}
          />
        ) : (
          <EmptyState title={`No cost attributed to work items ${periodLabel}`} />
        )}
      </section>

      <section className="grid grid-cols-1 gap-4 lg:grid-cols-2">
        <div className="rounded-lg border border-zinc-800 bg-zinc-900 p-4 light:border-zinc-200 light:bg-white">
          <h2 className="mb-3 text-sm font-semibold">Top models by cost</h2>
          {data.by_model.length > 0 ? (
            <RankedBars
              rows={data.by_model.map(namedRow)}
              maxValue={data.by_model[0].cost}
              initialLimit={5}
              onShowAll={() => openWorkItems({ groupBy: "none" })}
            />
          ) : (
            <EmptyState title={`No model costs ${periodLabel}`} />
          )}
        </div>
        <div className="rounded-lg border border-zinc-800 bg-zinc-900 p-4 light:border-zinc-200 light:bg-white">
          <h2 className="mb-3 text-sm font-semibold">Top agents by cost</h2>
          {data.by_agent.length > 0 ? (
            <RankedBars
              rows={data.by_agent.map(agentRow)}
              maxValue={data.by_agent[0].cost}
              initialLimit={5}
              onShowAll={() => openWorkItems({ groupBy: "none" })}
            />
          ) : (
            <EmptyState title={`No agent costs ${periodLabel}`} />
          )}
        </div>
      </section>

      <section className="rounded-lg border border-zinc-800 bg-zinc-900 p-4 light:border-zinc-200 light:bg-white">
        <h2 className="mb-3 text-sm font-semibold">Cost over time</h2>
        <TrendChart buckets={data.over_time} granularityMs={data.trend_granularity_ms} />
      </section>

      <section className="rounded-lg border border-zinc-800 bg-zinc-900 p-4 light:border-zinc-200 light:bg-white">
        <h2 className="mb-3 text-sm font-semibold">Recent sessions</h2>
        <RecentSessions sessions={data.recent_sessions} />
      </section>
    </div>
  );
}
