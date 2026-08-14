package session

import (
	"context"
	"database/sql"
	"fmt"

	"traceknot/internal/model"
	"traceknot/internal/sqlutil"
)

type Querier interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

func ResolveID(ctx context.Context, db Querier, rawID string) (*string, error) {
	if rawID == "" {
		return nil, nil
	}
	var sessionID string
	err := db.QueryRowContext(ctx, `
		SELECT session_id FROM sessions
		WHERE native_session_id = ? OR external_conversation_id = ?
		ORDER BY started_at_unix_ms DESC, inserted_at DESC
		LIMIT 1`,
		rawID,
		rawID,
	).Scan(&sessionID)
	if err == nil {
		return &sessionID, nil
	}
	if err != sql.ErrNoRows {
		return nil, fmt.Errorf("resolve session id: %w", err)
	}
	return nil, nil
}

func Upsert(ctx context.Context, db Querier, seed *model.SessionSeed) error {
	_, err := db.ExecContext(ctx, `
		INSERT INTO sessions (
			session_id, external_conversation_id, session_id_source, native_session_id,
			provider, title, started_at_unix_ms, ended_at_unix_ms,
			service_name, metadata_json
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(session_id) DO UPDATE SET
			external_conversation_id = COALESCE(excluded.external_conversation_id, sessions.external_conversation_id),
			session_id_source = excluded.session_id_source,
			native_session_id = COALESCE(excluded.native_session_id, sessions.native_session_id),
			provider = excluded.provider,
			title = CASE WHEN excluded.title != '' THEN excluded.title ELSE sessions.title END,
			started_at_unix_ms = excluded.started_at_unix_ms,
			ended_at_unix_ms = excluded.ended_at_unix_ms,
			service_name = COALESCE(excluded.service_name, sessions.service_name),
			metadata_json = excluded.metadata_json,
			updated_at = CURRENT_TIMESTAMP`,
		seed.SessionID,
		seed.ExternalConversationID,
		seed.SessionIDSource,
		seed.NativeSessionID,
		seed.Provider,
		seed.Title,
		seed.StartedAtUnixMs,
		seed.EndedAtUnixMs,
		seed.ServiceName,
		sqlutil.MustJSON(seed.Metadata),
	)
	if err != nil {
		return fmt.Errorf("upsert session: %w", err)
	}
	return nil
}
