import { getJSON } from "../../lib/http";

export type UnclaimedStatus = "pending" | "skipped" | "unclaimed";

export interface UnclaimedSession {
  session_id: string;
  provider: string;
  title: string;
  started_at_unix_ms: number | null;
  duration_ms: number | null;
  models: string[];
  cost: number;
  claim: { status: UnclaimedStatus };
}

export interface UnclaimedSessionsResponse {
  sessions: UnclaimedSession[];
  total_count: number;
}

export interface UnclaimedFilter {
  [key: string]: string | number | undefined;
  provider?: string;
  q?: string;
  sort?: string;
  offset?: number;
  limit?: number;
}

export function getUnclaimedSessions(filter: UnclaimedFilter): Promise<UnclaimedSessionsResponse> {
  return getJSON<UnclaimedSessionsResponse>("/sessions", { ...filter, unclaimed: "true" });
}
