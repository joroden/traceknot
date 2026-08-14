import { getJSON } from "../../lib/http";

export interface TokenBreakdown {
  total: number;
  raw?: number;
  cached?: number;
  write?: number;
  reasoning?: number;
}

export type ClaimStatus = "claimed" | "pending" | "skipped" | "unclaimed";

export interface ClaimState {
  status: ClaimStatus;
  work_item_key?: string;
  work_item_title?: string;
}

export interface SessionRow {
  [key: string]: unknown;
  session_id: string;
  provider: string;
  title: string;
  started_at_unix_ms: number | null;
  ended_at_unix_ms: number | null;
  duration_ms: number | null;
  models: string[];
  cost: number;
  input_tokens: TokenBreakdown;
  output_tokens: TokenBreakdown;
  claim: ClaimState;
}

export interface SessionsResponse {
  sessions: SessionRow[];
  total_count: number;
}

export type SortKey = "cost" | "input_tokens" | "output_tokens" | "started" | "duration" | "last_active";
export type SortDir = "asc" | "desc";

export interface SessionsFilter {
  [key: string]: string | number | undefined;
  provider?: string;
  model?: string;
  q?: string;
  start_unix_ms?: number;
  end_unix_ms?: number;
  work_item_key?: string;
  work_item_provider?: string;
  unclaimed?: "true";
  sort?: string;
  offset?: number;
  limit?: number;
}

export function encodeSort(sorts: { key: SortKey; dir: SortDir }[]): string | undefined {
  if (sorts.length === 0) {
    return undefined;
  }
  return sorts.map((entry) => `${entry.key}:${entry.dir}`).join(",");
}

export function getSessions(filter: SessionsFilter): Promise<SessionsResponse> {
  return getJSON<SessionsResponse>("/sessions", filter);
}

export interface WorkItemGroup {
  work_item_key: string;
  work_item_provider: string;
  title: string;
  is_unclaimed: boolean;
  session_count: number;
  cost: number;
  duration_ms: number | null;
  input_tokens: number;
  output_tokens: number;
}

export interface WorkItemGroupsResponse {
  groups: WorkItemGroup[];
  total_count: number;
}

export type GroupSortKey = SortKey | "session_count" | "name";

export interface WorkItemGroupsFilter {
  [key: string]: string | number | undefined;
  provider?: string;
  model?: string;
  q?: string;
  start_unix_ms?: number;
  end_unix_ms?: number;
  work_item_key?: string;
  work_item_provider?: string;
  sort?: string;
  offset?: number;
  limit?: number;
}

export function getWorkItemGroups(filter: WorkItemGroupsFilter): Promise<WorkItemGroupsResponse> {
  return getJSON<WorkItemGroupsResponse>("/work-items", filter);
}

export function groupIdentity(group: Pick<WorkItemGroup, "work_item_key" | "work_item_provider" | "is_unclaimed">): string {
  return group.is_unclaimed ? "unclaimed" : `${group.work_item_provider}:${group.work_item_key}`;
}
