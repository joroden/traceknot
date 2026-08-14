package model

type NodeKind string

const (
	NodeKindChat     NodeKind = "chat"
	NodeKindToolCall NodeKind = "tool_call"
	NodeKindAgent    NodeKind = "agent"
)

type SessionSeed struct {
	SessionID              string
	ExternalConversationID *string
	SessionIDSource        string
	NativeSessionID        *string
	Provider               string
	Title                  string
	ServiceName            *string
	StartedAtUnixMs        *int64
	EndedAtUnixMs          *int64
	Metadata               map[string]any
}

type NodeSeed struct {
	NodeID                string
	ParentNodeID          *string
	Kind                  NodeKind
	Name                  string
	Model                 *string
	Status                *string
	StartedAtUnixMs       *int64
	EndedAtUnixMs         *int64
	DurationMs            *float64
	PreviewText           string
	InputTokens           int64
	CachedInputTokens     int64
	CacheWriteTokens      int64
	OutputTokens          int64
	ReasoningTokens       int64
	Cost                  float64
	EstimatedInputTokens  int64
	EstimatedOutputTokens int64
	TokenEstimateMethod   *string
	MetadataJSON          string
}

type ChatSeed struct {
	NodeSeed
	SystemText    string
	PromptText    string
	OutputText    string
	ReasoningText string
}

type ToolCallSeed struct {
	NodeSeed
	ToolName         string
	ToolCallID       string
	ArgumentsJSON    string
	ResultText       string
	ErrorDetailsJSON string
	ApprovalDecision *string
	ApprovalSource   *string
}

type AgentSeed struct {
	NodeSeed
	AgentName       string
	AgentType       string
	SpawnPrompt     string
	SpawnToolCallID string
	ResultSummary   string
}

type SessionContent struct {
	Chats     []*ChatSeed
	ToolCalls []*ToolCallSeed
	Agents    []*AgentSeed
}
