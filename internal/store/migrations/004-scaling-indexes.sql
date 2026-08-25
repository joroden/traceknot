CREATE INDEX IF NOT EXISTS idx_nodes_started ON nodes (started_at_unix_ms);

CREATE INDEX IF NOT EXISTS idx_sessions_cost ON sessions (cost);
CREATE INDEX IF NOT EXISTS idx_sessions_duration ON sessions (duration_ms);
CREATE INDEX IF NOT EXISTS idx_sessions_ended ON sessions (ended_at_unix_ms);
CREATE INDEX IF NOT EXISTS idx_sessions_input_tokens ON sessions (input_tokens);
CREATE INDEX IF NOT EXISTS idx_sessions_output_tokens ON sessions (output_tokens);

CREATE TABLE IF NOT EXISTS conversation_roots (
  provider TEXT NOT NULL,
  native_id TEXT NOT NULL,
  root_native_id TEXT NOT NULL,
  PRIMARY KEY (provider, native_id)
);

CREATE INDEX IF NOT EXISTS idx_conversation_roots_root ON conversation_roots (provider, root_native_id);
