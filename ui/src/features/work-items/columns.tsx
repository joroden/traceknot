import type { MouseEvent } from "react";
import type { TableColumn } from "../../components/Table";
import { Badge } from "../../components/Badge";
import { SearchInput } from "../../components/SearchInput";
import { formatDuration, formatProviderLabel, formatTimestamp, formatUSD } from "../../utils/format";
import { providerBadgeTone } from "../../utils/providers";
import type { SessionRow, SortDir, SortKey, WorkItemGroup } from "./api";
import { groupIdentity } from "./api";
import { ClaimCell } from "./components/ClaimCell";
import { ColumnHeader } from "./components/ColumnHeader";
import { DateRangeFilter } from "./components/DateRangeFilter";
import { FilterPopover } from "./components/FilterPopover";
import { GroupRowLabel } from "./components/GroupRowLabel";
import { ModelBadges } from "../../components/ModelBadges";
import { ProviderFilterList } from "./components/ProviderFilterList";
import { TokenCell } from "./components/TokenCell";

export interface SortEntry {
  key: SortKey;
  dir: SortDir;
}

export type WorkItemsRow =
  | { [key: string]: unknown; kind: "session"; session: SessionRow; groupIdentity?: string }
  | { [key: string]: unknown; kind: "group"; group: WorkItemGroup };

export interface SessionsTableMeta {
  [key: string]: unknown;
  sorts: SortEntry[];
  onSortToggle: (key: SortKey, additive: boolean) => void;
  provider: string;
  onProviderChange: (value: string) => void;
  model: string;
  onModelChange: (value: string) => void;
  promptQuery: string;
  onPromptQueryChange: (value: string) => void;
  startUnixMs: number | undefined;
  endUnixMs: number | undefined;
  onDateRangeChange: (startUnixMs: number | undefined, endUnixMs: number | undefined) => void;
  expandedGroups: Set<string>;
}

interface HeaderContext {
  table: { options: { meta?: unknown } };
}

interface RowContext extends HeaderContext {
  row: { original: WorkItemsRow };
}

const PROVIDER_OPTIONS = [
  { value: "claude", label: "Claude" },
  { value: "codex", label: "Codex" },
  { value: "copilot", label: "Copilot" },
];

function requireMeta(context: HeaderContext): SessionsTableMeta {
  const meta = context.table.options.meta as SessionsTableMeta | undefined;
  if (!meta) {
    throw new Error("work items table meta is required");
  }
  return meta;
}

function sortHeaderProps(meta: SessionsTableMeta, key: SortKey, label: string) {
  const index = meta.sorts.findIndex((entry) => entry.key === key);
  return {
    label,
    sortActive: index !== -1,
    sortDir: index !== -1 ? meta.sorts[index].dir : ("desc" as SortDir),
    sortPriority: meta.sorts.length > 1 && index !== -1 ? index + 1 : undefined,
    onSortClick: (event: MouseEvent) => meta.onSortToggle(key, event.shiftKey),
  };
}

function sortableHeader(key: SortKey, label: string) {
  return (context: HeaderContext) => {
    const meta = requireMeta(context);
    return <ColumnHeader {...sortHeaderProps(meta, key, label)} />;
  };
}

function PromptHeader(context: HeaderContext) {
  const meta = requireMeta(context);
  return (
    <ColumnHeader
      label="Title"
      filter={
        <FilterPopover active={meta.promptQuery.length > 0}>
          <SearchInput
            value={meta.promptQuery}
            onChange={meta.onPromptQueryChange}
            placeholder="Search titles…"
            autoFocus
          />
        </FilterPopover>
      }
    />
  );
}

function StartedHeader(context: HeaderContext) {
  const meta = requireMeta(context);
  return (
    <ColumnHeader
      {...sortHeaderProps(meta, "started", "Started")}
      filter={
        <FilterPopover active={meta.startUnixMs !== undefined || meta.endUnixMs !== undefined}>
          <DateRangeFilter
            startUnixMs={meta.startUnixMs}
            endUnixMs={meta.endUnixMs}
            onChange={meta.onDateRangeChange}
          />
        </FilterPopover>
      }
    />
  );
}

function ProviderHeader(context: HeaderContext) {
  const meta = requireMeta(context);
  return (
    <ColumnHeader
      label="Provider"
      filter={
        <FilterPopover active={meta.provider.length > 0}>
          <ProviderFilterList
            value={meta.provider}
            onChange={meta.onProviderChange}
            options={PROVIDER_OPTIONS}
            placeholder="All providers"
          />
        </FilterPopover>
      }
    />
  );
}

function ModelsHeader(context: HeaderContext) {
  const meta = requireMeta(context);
  return (
    <ColumnHeader
      label="Models"
      filter={
        <FilterPopover active={meta.model.length > 0}>
          <SearchInput value={meta.model} onChange={meta.onModelChange} placeholder="Search models…" />
        </FilterPopover>
      }
    />
  );
}

function keyCell(info: RowContext) {
  const row = info.row.original;
  if (row.kind === "group") {
    const meta = requireMeta(info);
    return <GroupRowLabel group={row.group} expanded={meta.expandedGroups.has(groupIdentity(row.group))} />;
  }
  const key = row.session.claim.work_item_key;
  if (!key) {
    return <span className="text-xs text-zinc-600">—</span>;
  }
  return <span className="truncate font-mono text-xs text-zinc-400 light:text-zinc-500">{key}</span>;
}

function titleOf(row: WorkItemsRow): string {
  const title = row.kind === "group" ? row.group.title : row.session.title;
  return title && title.trim().length > 0 ? title : "Untitled";
}

function promptCell(info: RowContext) {
  const row = info.row.original;
  if (row.kind === "group") {
    return (
      <span className="flex min-w-0 items-center gap-2">
        <span className="truncate text-sm text-zinc-100 light:text-zinc-900">{titleOf(row)}</span>
        <span className="shrink-0 text-xs text-zinc-500">
          ({row.group.session_count} session{row.group.session_count === 1 ? "" : "s"})
        </span>
      </span>
    );
  }
  return <span className="truncate text-zinc-100 light:text-zinc-900">{titleOf(row)}</span>;
}

function startedCell(info: RowContext) {
  const row = info.row.original;
  if (row.kind === "group") {
    return <span className="text-xs text-zinc-600">—</span>;
  }
  return (
    <span className="text-xs text-zinc-400 light:text-zinc-500">{formatTimestamp(row.session.started_at_unix_ms)}</span>
  );
}

function durationCell(info: RowContext) {
  const row = info.row.original;
  const durationMs = row.kind === "group" ? row.group.duration_ms : row.session.duration_ms;
  return (
    <span
      className={
        row.kind === "group"
          ? "font-mono text-xs font-semibold text-zinc-200 light:text-zinc-800"
          : "font-mono text-xs text-zinc-400 light:text-zinc-500"
      }
    >
      {formatDuration(durationMs)}
    </span>
  );
}

function providerCell(info: RowContext) {
  const row = info.row.original;
  if (row.kind === "group") {
    return <span className="text-xs text-zinc-600">—</span>;
  }
  return (
    <Badge tone={providerBadgeTone(row.session.provider)}>{formatProviderLabel(row.session.provider)}</Badge>
  );
}

function modelsCell(info: RowContext) {
  const row = info.row.original;
  if (row.kind === "group") {
    return <span className="text-xs text-zinc-600">—</span>;
  }
  return <ModelBadges models={row.session.models} />;
}

function inputTokensCell(info: RowContext) {
  const row = info.row.original;
  const tokens = row.kind === "group" ? { total: row.group.input_tokens } : row.session.input_tokens;
  return <TokenCell tokens={tokens} kind="input" />;
}

function outputTokensCell(info: RowContext) {
  const row = info.row.original;
  const tokens = row.kind === "group" ? { total: row.group.output_tokens } : row.session.output_tokens;
  return <TokenCell tokens={tokens} kind="output" />;
}

function costCell(info: RowContext) {
  const row = info.row.original;
  const cost = row.kind === "group" ? row.group.cost : row.session.cost;
  return (
    <span
      className={
        row.kind === "group"
          ? "font-mono text-xs font-bold text-zinc-100 light:text-zinc-900"
          : "font-mono text-xs font-semibold text-zinc-100 light:text-zinc-900"
      }
    >
      {formatUSD(cost)}
    </span>
  );
}

function claimCell(info: RowContext) {
  const row = info.row.original;
  if (row.kind === "group") {
    return <span className="text-xs text-zinc-600">—</span>;
  }
  return <ClaimCell claim={row.session.claim} />;
}

export const workItemsColumns: TableColumn<WorkItemsRow>[] = [
  {
    id: "key",
    accessorFn: (row: WorkItemsRow) => (row.kind === "group" ? row.group.work_item_key : row.session.claim.work_item_key ?? ""),
    size: 160,
    header: () => <ColumnHeader label="Key" />,
    cell: keyCell,
  },
  {
    id: "prompt",
    accessorFn: (row: WorkItemsRow) => titleOf(row),
    size: 270,
    header: PromptHeader,
    cell: promptCell,
  },
  {
    id: "cost",
    accessorFn: (row: WorkItemsRow) => (row.kind === "group" ? row.group.cost : row.session.cost),
    size: 90,
    header: sortableHeader("cost", "Cost"),
    cell: costCell,
  },
  {
    id: "duration",
    accessorFn: (row: WorkItemsRow) => (row.kind === "group" ? row.group.duration_ms : row.session.duration_ms),
    size: 90,
    header: sortableHeader("duration", "Duration"),
    cell: durationCell,
  },
  {
    id: "provider",
    accessorFn: (row: WorkItemsRow) => (row.kind === "group" ? "" : row.session.provider),
    size: 130,
    header: ProviderHeader,
    cell: providerCell,
  },
  {
    id: "models",
    accessorFn: (row: WorkItemsRow) => (row.kind === "group" ? [] : row.session.models),
    size: 200,
    header: ModelsHeader,
    cell: modelsCell,
  },
  {
    id: "input_tokens",
    accessorFn: (row: WorkItemsRow) => (row.kind === "group" ? row.group.input_tokens : row.session.input_tokens.total),
    size: 220,
    header: sortableHeader("input_tokens", "Input"),
    cell: inputTokensCell,
  },
  {
    id: "output_tokens",
    accessorFn: (row: WorkItemsRow) =>
      row.kind === "group" ? row.group.output_tokens : row.session.output_tokens.total,
    size: 170,
    header: sortableHeader("output_tokens", "Output"),
    cell: outputTokensCell,
  },
  {
    id: "started",
    accessorFn: (row: WorkItemsRow) => (row.kind === "group" ? null : row.session.started_at_unix_ms),
    size: 130,
    header: StartedHeader,
    cell: startedCell,
  },
  {
    id: "claim",
    accessorFn: (row: WorkItemsRow) => (row.kind === "group" ? "" : row.session.claim.status),
    size: 190,
    header: () => <ColumnHeader label="Work item" />,
    cell: claimCell,
  },
];
