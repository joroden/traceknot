import { request } from "../../lib/http";

export interface SessionMeta {
  session_id: string;
  provider: string;
  service_name: string | null;
  status: string | null;
  started_at_unix_ms: number | null;
  ended_at_unix_ms: number | null;
  duration_ms: number | null;
  turn_count: number;
  node_count: number;
  tool_call_count: number;
  agent_run_count: number;
  input_tokens: number;
  cached_input_tokens: number;
  output_tokens: number;
  reasoning_tokens: number;
  cache_write_tokens: number;
  cost: number;
  models: string[];
}

export interface TreeNodeRow {
  node_id: string;
  kind: string;
  name: string | null;
  agent_name: string | null;
  tool_name: string | null;
  tool_call_id: string | null;
  model: string | null;
  status: string | null;
  started_at_unix_ms: number | null;
  duration_ms: number | null;
  input_tokens: number;
  cached_input_tokens: number;
  cache_write_tokens: number;
  output_tokens: number;
  reasoning_tokens: number;
  cost: number;
  estimated_input_tokens: number;
  estimated_output_tokens: number;
  owning_agent_id: string | null;
  owning_agent_name: string | null;
  parent_node_id: string | null;
  has_content: boolean;
  agg_input_tokens: number;
  agg_cached_input_tokens: number;
  agg_cache_write_tokens: number;
  agg_output_tokens: number;
  agg_reasoning_tokens: number;
  agg_cost: number;
  descendant_count: number;
  subagent_count: number;
}

export interface TreeResponse {
  session: SessionMeta;
  nodes: TreeNodeRow[];
}

export interface NodeDetail {
  node_id: string;
  session_id: string;
  kind: string;
  name: string | null;
  agent_name: string | null;
  tool_name: string | null;
  tool_call_id: string | null;
  model: string | null;
  status: string | null;
  started_at_unix_ms: number | null;
  ended_at_unix_ms: number | null;
  duration_ms: number | null;
  input_tokens: number;
  cached_input_tokens: number;
  cache_write_tokens: number;
  output_tokens: number;
  reasoning_tokens: number;
  cost: number;
  estimated_input_tokens: number;
  estimated_output_tokens: number;
  token_estimate_method: string | null;
  prompt_text: string | null;
  output_text: string | null;
  reasoning_text: string | null;
  arguments_json: string | null;
  result_text: string | null;
  error_details_json: string | null;
  owning_agent_id: string | null;
  owning_agent_name: string | null;
  metadata_json: string;
}

export function fetchSessionTree(sessionId: string): Promise<TreeResponse> {
  return request(`/sessions/${encodeURIComponent(sessionId)}/tree`);
}

export function fetchNodeDetail(nodeId: string): Promise<NodeDetail> {
  return request<{ node: NodeDetail }>(`/nodes/${encodeURIComponent(nodeId)}`).then(
    (payload) => payload.node,
  );
}
