import { postJSON, request } from "../lib/http";

export type ProviderStatus = "available" | "cli_missing" | "not_authenticated" | "error";

export interface ProviderProbe {
  provider: string;
  status: ProviderStatus;
  hint?: string;
  install_docs_url?: string;
  auth_docs_url?: string;
}

export interface WorkItem {
  key: string;
  title: string;
  url?: string;
  status?: string;
  type?: string;
  project?: string;
  updated_at_unix_ms?: number;
}

export interface RecentWorkItem {
  key: string;
  provider: string;
  title: string;
  project: string;
  last_attributed_at_unix_ms: number;
  attribution_count: number;
}

export interface WorkItemRow {
  key: string;
  title: string;
  project: string;
  provider: string;
  updatedAtUnixMs: number | null;
}

export function toRow(item: WorkItem, provider: string): WorkItemRow {
  return {
    key: item.key,
    title: item.title,
    project: item.project ?? "",
    provider,
    updatedAtUnixMs: item.updated_at_unix_ms ?? null,
  };
}

export function recentToRow(item: RecentWorkItem): WorkItemRow {
  return {
    key: item.key,
    title: item.title,
    project: item.project ?? "",
    provider: item.provider,
    updatedAtUnixMs: item.last_attributed_at_unix_ms,
  };
}

export async function getProviders(): Promise<ProviderProbe[]> {
  const payload = await request<{ providers: ProviderProbe[] }>("/providers");
  return payload.providers;
}

export async function getRecent(provider: string | null, limit = 50): Promise<RecentWorkItem[]> {
  const query = new URLSearchParams({ limit: String(limit) });
  if (provider) {
    query.set("provider", provider);
  }
  const payload = await request<{ items: RecentWorkItem[] }>(`/picker/recent?${query}`);
  return payload.items;
}

export async function searchItems(provider: string, query: string, limit = 20): Promise<WorkItem[]> {
  const params = new URLSearchParams({ q: query, limit: String(limit) });
  const payload = await request<{ items: WorkItem[] }>(`/providers/${encodeURIComponent(provider)}/search?${params}`);
  return payload.items;
}

export async function postClaim(sessionID: string | null, item: WorkItemRow): Promise<void> {
  const body: Record<string, unknown> = {
    source: "hook",
    work_item: {
      key: item.key,
      title: item.title,
      provider: item.provider,
      project: item.project,
    },
  };
  if (sessionID) {
    body.session_id = sessionID;
  }
  await postJSON("/claims", body);
}

export async function postSkip(sessionID: string | null): Promise<void> {
  if (!sessionID) {
    return;
  }
  await postJSON("/picker/outcome", { session_id: sessionID });
}
