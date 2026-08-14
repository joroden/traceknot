package providers

import "regexp"

var (
	repositoryURLPattern = regexp.MustCompile(`/repos/([^/?#]+/[^/?#]+)/?$`)
	htmlURLPattern       = regexp.MustCompile(`^https?://[^/]+/([^/]+/[^/]+)/(?:issues|pull)/\d+`)
	jiraKeyPattern       = regexp.MustCompile(`^[A-Z][A-Z0-9]+-\d+$`)
	jiraSelfPattern      = regexp.MustCompile(`^(https?://[^/]+)/rest/api/\d+/issue/([^/?#]+)`)
)
