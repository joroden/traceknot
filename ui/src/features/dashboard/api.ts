import { getJSON } from "../../lib/http";

export interface DashboardRequest {
  [key: string]: string | number | undefined;
  range: "today" | "week" | "month" | "all" | "custom";
  start_unix_ms?: number;
  end_unix_ms?: number;
}

export interface PeriodSummary {
  start_unix_ms: number;
  end_unix_ms: number;
  cost: number;
  cost_delta_pct: number | null;
  tokens: number;
  tokens_delta_pct: number | null;
  input_tokens: number;
  output_tokens: number;
  session_count: number;
  coverage_pct: number | null;
  coverage_delta_pct: number | null;
  unattributed_cost: number;
  unattributed_session_count: number;
}

export interface WorkItemCost {
  key: string;
  title: string;
  provider: string;
  project: string;
  cost: number;
  session_count: number;
}

export interface NamedCost {
  name: string;
  cost: number;
}

export interface TrendBucket {
  start_unix_ms: number;
  total_cost: number;
  models: NamedCost[];
}

export interface RecentSession {
  session_id: string;
  provider: string;
  started_at_unix_ms: number | null;
  cost: number;
  tokens: number;
  status: string;
  node_count: number;
  models: string[];
  title: string;
}

export interface DashboardData {
  first_run: boolean;
  period: PeriodSummary;
  by_work_item: WorkItemCost[];
  by_model: NamedCost[];
  by_agent: NamedCost[];
  over_time: TrendBucket[];
  trend_granularity_ms: number;
  recent_sessions: RecentSession[];
}

export async function getDashboard(request: DashboardRequest): Promise<DashboardData> {
  return getJSON<DashboardData>("/dashboard", {
    range: request.range,
    start_unix_ms: request.start_unix_ms,
    end_unix_ms: request.end_unix_ms,
  });
}
