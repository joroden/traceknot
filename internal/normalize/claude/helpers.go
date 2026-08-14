package claude

import (
	"traceknot/internal/normalize/shared"
	"traceknot/internal/tokenize"
)

func estimateToolTokens(args string, result string, estimator *tokenize.Estimator) (int64, int64, *string) {
	return shared.EstimateToolTokens(args, result, estimator)
}
