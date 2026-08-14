CREATE TABLE IF NOT EXISTS sessions (
  session_id TEXT PRIMARY KEY,
  external_conversation_id TEXT,
  session_id_source TEXT NOT NULL,
  native_session_id TEXT,
  provider TEXT NOT NULL,
  title TEXT NOT NULL DEFAULT '',
  started_at_unix_ms INTEGER,
  ended_at_unix_ms INTEGER,
  duration_ms REAL,
  status TEXT,
  service_name TEXT,
  turn_count INTEGER NOT NULL DEFAULT 0,
  node_count INTEGER NOT NULL DEFAULT 0,
  tool_call_count INTEGER NOT NULL DEFAULT 0,
  agent_run_count INTEGER NOT NULL DEFAULT 0,
  input_tokens INTEGER NOT NULL DEFAULT 0,
  cached_input_tokens INTEGER NOT NULL DEFAULT 0,
  non_cached_input_tokens INTEGER NOT NULL DEFAULT 0,
  output_tokens INTEGER NOT NULL DEFAULT 0,
  reasoning_tokens INTEGER NOT NULL DEFAULT 0,
  cache_write_tokens INTEGER NOT NULL DEFAULT 0,
  cost REAL NOT NULL DEFAULT 0,
  models_json TEXT NOT NULL DEFAULT '[]',
  metadata_json TEXT NOT NULL DEFAULT '{}',
  inserted_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_sessions_started ON sessions (started_at_unix_ms);
CREATE INDEX IF NOT EXISTS idx_sessions_provider ON sessions (provider, started_at_unix_ms);

CREATE TABLE IF NOT EXISTS nodes (
  node_id TEXT PRIMARY KEY,
  session_id TEXT NOT NULL REFERENCES sessions(session_id) ON DELETE CASCADE,
  parent_node_id TEXT REFERENCES nodes(node_id) ON DELETE CASCADE,
  kind TEXT NOT NULL,
  name TEXT,
  status TEXT,
  model TEXT,
  started_at_unix_ms INTEGER,
  ended_at_unix_ms INTEGER,
  duration_ms REAL,
  preview_text TEXT,
  input_tokens INTEGER NOT NULL DEFAULT 0,
  cached_input_tokens INTEGER NOT NULL DEFAULT 0,
  cache_write_tokens INTEGER NOT NULL DEFAULT 0,
  output_tokens INTEGER NOT NULL DEFAULT 0,
  reasoning_tokens INTEGER NOT NULL DEFAULT 0,
  cost REAL NOT NULL DEFAULT 0,
  estimated_input_tokens INTEGER NOT NULL DEFAULT 0,
  estimated_output_tokens INTEGER NOT NULL DEFAULT 0,
  token_estimate_method TEXT,
  owning_agent_id TEXT,
  owning_agent_name TEXT,
  has_content INTEGER NOT NULL DEFAULT 0,
  metadata_json TEXT NOT NULL DEFAULT '{}',
  inserted_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_nodes_session ON nodes (session_id, started_at_unix_ms);
CREATE INDEX IF NOT EXISTS idx_nodes_parent ON nodes (session_id, parent_node_id);

CREATE TABLE IF NOT EXISTS chat_nodes (
  node_id TEXT PRIMARY KEY REFERENCES nodes(node_id) ON DELETE CASCADE,
  system_text TEXT,
  prompt_text TEXT,
  output_text TEXT,
  reasoning_text TEXT
);

CREATE TABLE IF NOT EXISTS tool_call_nodes (
  node_id TEXT PRIMARY KEY REFERENCES nodes(node_id) ON DELETE CASCADE,
  tool_name TEXT,
  tool_call_id TEXT,
  arguments_json TEXT,
  result_text TEXT,
  error_details_json TEXT,
  approval_decision TEXT,
  approval_source TEXT
);

CREATE TABLE IF NOT EXISTS agent_nodes (
  node_id TEXT PRIMARY KEY REFERENCES nodes(node_id) ON DELETE CASCADE,
  agent_name TEXT,
  agent_type TEXT,
  spawn_prompt TEXT,
  spawn_tool_call_id TEXT,
  result_summary TEXT
);

CREATE TABLE IF NOT EXISTS app_meta (
  id INTEGER PRIMARY KEY CHECK (id = 1),
  install_ts INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS raw_signal (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  native_id TEXT NOT NULL,
  provider TEXT NOT NULL,
  signal TEXT NOT NULL,
  dedup_key TEXT NOT NULL,
  timestamp_unix_ms INTEGER NOT NULL,
  payload_json TEXT NOT NULL,
  received_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_raw_signal_dedup ON raw_signal (native_id, dedup_key);
CREATE INDEX IF NOT EXISTS idx_raw_signal_native ON raw_signal (native_id, timestamp_unix_ms, id);
CREATE INDEX IF NOT EXISTS idx_raw_signal_provider ON raw_signal (provider, native_id);