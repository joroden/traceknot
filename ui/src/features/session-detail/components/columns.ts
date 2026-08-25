export const GRID_TEMPLATE =
  "minmax(200px, 320px) 64px 84px minmax(150px, 1fr) 100px repeat(5, 84px) 90px";

export interface ColumnTitle {
  key: string;
  title: string;
  tooltip: string;
  headerClassName: string;
}

export const COLUMN_TITLES: ColumnTitle[] = [
  { key: "name", title: "Name", tooltip: "Node name", headerClassName: "" },
  { key: "status", title: "Status", tooltip: "Node status", headerClassName: "text-center" },
  {
    key: "subagents",
    title: "Subagents",
    tooltip: "Subagents in this subtree",
    headerClassName: "text-center",
  },
  { key: "model", title: "Model", tooltip: "Model used", headerClassName: "pl-2 pr-2 text-left" },
  {
    key: "duration",
    title: "Duration",
    tooltip: "Duration (subtree where applicable)",
    headerClassName: "pr-3 text-right",
  },
  {
    key: "input",
    title: "Input",
    tooltip: "Input tokens (subtree where applicable)",
    headerClassName: "pr-3 text-right",
  },
  { key: "cached", title: "Cached", tooltip: "Cached input tokens", headerClassName: "pr-3 text-right" },
  { key: "write", title: "Write", tooltip: "Cache write tokens", headerClassName: "pr-3 text-right" },
  { key: "output", title: "Output", tooltip: "Output tokens", headerClassName: "pr-3 text-right" },
  {
    key: "reasoning",
    title: "Reasoning",
    tooltip: "Reasoning tokens",
    headerClassName: "pr-3 text-right",
  },
  { key: "cost", title: "Cost", tooltip: "Cost (subtree where applicable)", headerClassName: "pr-1 text-right" },
];
