package providers

import "context"

type Status string

const (
	StatusAvailable        Status = "available"
	StatusCLIMissing       Status = "cli_missing"
	StatusNotAuthenticated Status = "not_authenticated"
	StatusError            Status = "error"
)

type Probe struct {
	Provider       string `json:"provider"`
	Status         Status `json:"status"`
	Hint           string `json:"hint,omitempty"`
	InstallDocsURL string `json:"install_docs_url,omitempty"`
	AuthDocsURL    string `json:"auth_docs_url,omitempty"`
}

type WorkItem struct {
	Key             string `json:"key"`
	Title           string `json:"title"`
	URL             string `json:"url,omitempty"`
	Status          string `json:"status,omitempty"`
	Type            string `json:"type,omitempty"`
	Project         string `json:"project,omitempty"`
	UpdatedAtUnixMs int64  `json:"updated_at_unix_ms,omitempty"`
}

type Provider interface {
	ID() string
	Probe(ctx context.Context) Probe
	Search(ctx context.Context, query string, limit int) ([]WorkItem, error)
}
