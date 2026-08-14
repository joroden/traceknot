export const GRID_TEMPLATE =
  "minmax(280px, 1fr) 64px 76px 140px 110px repeat(5, 84px) 90px";

export interface ColumnTitle {
  key: string;
  title: string;
  tooltip: string;
}

export const COLUMN_TITLES: ColumnTitle[] = [
  { key: "name", title: "Name", tooltip: "Node name" },
  { key: "status", title: "Status", tooltip: "Node status" },
  { key: "subagents", title: "Subagents", tooltip: "Subagents in this subtree" },
  { key: "model", title: "Model", tooltip: "Model used" },
  { key: "duration", title: "Duration", tooltip: "Duration (subtree where applicable)" },
  { key: "input", title: "Input", tooltip: "Input tokens (subtree where applicable)" },
  { key: "cached", title: "Cached", tooltip: "Cached input tokens" },
  { key: "write", title: "Write", tooltip: "Cache write tokens" },
  { key: "output", title: "Output", tooltip: "Output tokens" },
  { key: "reasoning", title: "Reasoning", tooltip: "Reasoning tokens" },
  { key: "cost", title: "Cost", tooltip: "Cost (subtree where applicable)" },
];
