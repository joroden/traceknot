import { useCallback, useMemo, useState } from "react";
import { useLocation, useNavigate } from "react-router";
import { DaemonUnreachable } from "../../components/DaemonUnreachable";
import { EmptyState } from "../../components/EmptyState";
import { SegmentedControl, type SegmentedOption } from "../../components/SegmentedControl";
import { Skeleton } from "../../components/Skeleton";
import { Table } from "../../components/Table";
import { encodeSort, type GroupSortKey, type SessionsFilter, type SortDir } from "./api";
import { type SessionsTableMeta, type WorkItemsRow, workItemsColumns } from "./columns";
import { useDebouncedValue } from "../../hooks/useDebouncedValue";
import { useGroupedRows, type WorkItemScope } from "./hooks/useGroupedRows";
import { useMultiSort } from "./hooks/useMultiSort";
import { useSessions } from "./hooks/useSessions";

type GroupBy = "work_item" | "none";
type GroupSortField = Extract<GroupSortKey, "cost" | "session_count" | "name">;

interface LocationFilter {
  provider?: string;
  model?: string;
  workItemKey?: string;
  workItemProvider?: string;
  groupBy?: GroupBy;
}

const GROUP_BY_OPTIONS: SegmentedOption<GroupBy>[] = [
  { value: "work_item", label: "Work item" },
  { value: "none", label: "None" },
];

const GROUP_SORT_OPTIONS: SegmentedOption<GroupSortField>[] = [
  { value: "cost", label: "Cost" },
  { value: "session_count", label: "Sessions" },
  { value: "name", label: "Name" },
];

const GROW_COLUMN_IDS = ["prompt"];
const EMPTY_EXPANDED = new Set<string>();

function initialGroupBy(filter: LocationFilter): GroupBy {
  if (filter.groupBy) {
    return filter.groupBy;
  }
  if (filter.workItemKey) {
    return "work_item";
  }
  if (filter.provider || filter.model) {
    return "none";
  }
  return "work_item";
}

function groupRowClassName(row: WorkItemsRow): string {
  return row.kind === "group" ? "bg-zinc-900 font-medium light:bg-zinc-100" : "";
}

export function WorkItemsPage() {
  const location = useLocation();
  const navigate = useNavigate();
  const initialFilter = (location.state as { filter?: LocationFilter } | null)?.filter ?? {};

  const [groupBy, setGroupBy] = useState<GroupBy>(() => initialGroupBy(initialFilter));
  const [provider, setProvider] = useState(initialFilter.provider ?? "");
  const [model, setModel] = useState(initialFilter.model ?? "");
  const [promptQuery, setPromptQuery] = useState("");
  const [startUnixMs, setStartUnixMs] = useState<number | undefined>(undefined);
  const [endUnixMs, setEndUnixMs] = useState<number | undefined>(undefined);
  const [scope, setScope] = useState<WorkItemScope | null>(
    initialFilter.workItemKey && initialFilter.workItemProvider
      ? { key: initialFilter.workItemKey, provider: initialFilter.workItemProvider }
      : null,
  );
  const [sorts, onSortToggle] = useMultiSort(
    initialFilter.groupBy === "none" ? [{ key: "cost", dir: "desc" }] : [],
  );
  const [groupSortKey, setGroupSortKey] = useState<GroupSortField>("cost");
  const [groupSortDir, setGroupSortDir] = useState<SortDir>("desc");

  const debouncedModel = useDebouncedValue(model, 250);
  const debouncedPromptQuery = useDebouncedValue(promptQuery, 250);

  const onDateRangeChange = useCallback((start: number | undefined, end: number | undefined) => {
    setStartUnixMs(start);
    setEndUnixMs(end);
  }, []);

  const baseFilter = useMemo(
    () => ({
      provider: provider || undefined,
      model: debouncedModel || undefined,
      q: debouncedPromptQuery || undefined,
      start_unix_ms: startUnixMs,
      end_unix_ms: endUnixMs,
    }),
    [provider, debouncedModel, debouncedPromptQuery, startUnixMs, endUnixMs],
  );

  const flatFilter = useMemo<SessionsFilter>(
    () => ({ ...baseFilter, sort: encodeSort(sorts) }),
    [baseFilter, sorts],
  );
  const flat = useSessions(flatFilter);
  const groupSort = useMemo(() => ({ key: groupSortKey, dir: groupSortDir }), [groupSortKey, groupSortDir]);
  const grouped = useGroupedRows({ baseFilter, sorts, groupSort, scope });

  const rows: WorkItemsRow[] = useMemo(() => {
    if (groupBy === "work_item") {
      return grouped.rows;
    }
    return flat.rows.map((session): WorkItemsRow => ({ kind: "session", session }));
  }, [groupBy, grouped.rows, flat.rows]);

  const columns = useMemo(
    () => workItemsColumns.filter((column) => column.id !== (groupBy === "work_item" ? "claim" : "key")),
    [groupBy],
  );

  const loading = groupBy === "work_item" ? grouped.loading : flat.loading;
  const error = groupBy === "work_item" ? grouped.error : flat.error;

  const onRowClick = useCallback(
    (row: WorkItemsRow) => {
      if (row.kind === "group") {
        grouped.toggleGroup(row.group);
        return;
      }
      navigate(`/sessions/${encodeURIComponent(row.session.session_id)}`);
    },
    [grouped, navigate],
  );

  const onEndReached = useCallback(
    (row: WorkItemsRow) => {
      if (groupBy === "work_item") {
        grouped.onEndReached(row);
      } else {
        flat.loadMore();
      }
    },
    [groupBy, grouped, flat],
  );

  const meta: SessionsTableMeta = useMemo(
    () => ({
      sorts,
      onSortToggle,
      provider,
      onProviderChange: setProvider,
      model,
      onModelChange: setModel,
      promptQuery,
      onPromptQueryChange: setPromptQuery,
      startUnixMs,
      endUnixMs,
      onDateRangeChange,
      expandedGroups: groupBy === "work_item" ? grouped.expandedGroups : EMPTY_EXPANDED,
    }),
    [
      sorts,
      onSortToggle,
      provider,
      model,
      promptQuery,
      startUnixMs,
      endUnixMs,
      onDateRangeChange,
      groupBy,
      grouped.expandedGroups,
    ],
  );

  if (error) {
    return <DaemonUnreachable error={error} />;
  }

  return (
    <div className="flex h-full min-h-0 flex-col gap-4">
      <div className="flex flex-wrap items-center justify-between gap-2">
        <div className="flex flex-wrap items-center gap-4">
          <div className="flex items-center gap-2">
            <span className="text-xs text-zinc-500">Group by</span>
            <SegmentedControl options={GROUP_BY_OPTIONS} value={groupBy} onChange={setGroupBy} />
          </div>
          {groupBy === "work_item" ? (
            <div className="flex items-center gap-2">
              <span className="text-xs text-zinc-500">Sort groups by</span>
              <SegmentedControl options={GROUP_SORT_OPTIONS} value={groupSortKey} onChange={setGroupSortKey} />
              <button
                type="button"
                onClick={() => setGroupSortDir((current) => (current === "desc" ? "asc" : "desc"))}
                title={groupSortDir === "desc" ? "Descending" : "Ascending"}
                className="cursor-pointer rounded-md border border-zinc-700 px-2 py-1.5 text-xs text-zinc-400 hover:text-zinc-100 light:border-zinc-300 light:text-zinc-500 light:hover:text-zinc-900"
              >
                {groupSortDir === "desc" ? "↓" : "↑"}
              </button>
            </div>
          ) : null}
        </div>
        {scope ? (
          <div className="flex items-center gap-2 rounded-md border border-violet-700 bg-violet-900/30 px-2.5 py-1.5 text-xs text-violet-300 light:border-violet-300 light:bg-violet-100 light:text-violet-700">
            <span className="font-mono font-semibold">{scope.key}</span>
            <span className="text-violet-400 light:text-violet-600">scoped to this work item</span>
            <button
              type="button"
              onClick={() => setScope(null)}
              className="cursor-pointer font-semibold text-violet-300 hover:text-violet-100 light:text-violet-700 light:hover:text-violet-900"
            >
              ×
            </button>
          </div>
        ) : null}
      </div>

      {loading && rows.length === 0 ? (
        <div className="flex flex-col gap-2">
          {Array.from({ length: 10 }, (_, index) => (
            <Skeleton key={index} className="h-11 w-full" />
          ))}
        </div>
      ) : (
        <div className={loading ? "min-h-0 flex-1 opacity-60" : "min-h-0 flex-1"}>
          <Table
            columns={columns}
            data={rows}
            getRowId={(row: WorkItemsRow) =>
              row.kind === "group" ? `group:${row.group.work_item_provider}:${row.group.work_item_key}:${row.group.is_unclaimed}` : row.session.session_id
            }
            onRowClick={onRowClick}
            onEndReached={onEndReached}
            meta={meta}
            growColumnIds={GROW_COLUMN_IDS}
            rowHeightPx={52}
            rowClassName={groupBy === "work_item" ? groupRowClassName : undefined}
            emptyState={<EmptyState title="No sessions match these filters." />}
          />
        </div>
      )}
    </div>
  );
}
