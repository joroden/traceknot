CREATE TABLE IF NOT EXISTS work_item_recent (
  key TEXT NOT NULL,
  provider TEXT NOT NULL,
  title TEXT NOT NULL,
  project TEXT NOT NULL DEFAULT '',
  last_attributed_at_unix_ms INTEGER NOT NULL,
  attribution_count INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (provider, key)
);

CREATE INDEX IF NOT EXISTS idx_work_item_recent_attributed
  ON work_item_recent (last_attributed_at_unix_ms DESC);

CREATE TABLE IF NOT EXISTS claims (
  claim_id TEXT PRIMARY KEY,
  session_id TEXT,
  work_item_key TEXT NOT NULL,
  work_item_title TEXT NOT NULL,
  provider TEXT NOT NULL,
  project TEXT NOT NULL DEFAULT '',
  source TEXT NOT NULL DEFAULT 'hook',
  status TEXT NOT NULL DEFAULT 'claimed',
  claimed_at_unix_ms INTEGER NOT NULL,
  updated_at_unix_ms INTEGER NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_claims_session_primary
  ON claims (session_id) WHERE session_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_claims_work_item
  ON claims (provider, work_item_key, claimed_at_unix_ms DESC);