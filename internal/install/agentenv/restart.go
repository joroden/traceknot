package agentenv

import (
	"context"

	"traceknot/internal/platform"
)

func CloseInteractiveSessions(ctx context.Context) (int, error) {
	return platform.Current.CloseInteractiveSessions(ctx)
}
