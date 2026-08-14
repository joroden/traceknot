package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	githubScopeCacheTTL = 5 * time.Minute
	githubInstallDocs   = "https://cli.github.com"
	githubAuthDocs      = "https://cli.github.com/manual/gh_auth_login"
)

type GitHub struct {
	env        cliEnv
	scopeCache githubScopeCache
}

type githubScopeCache struct {
	expiresAt time.Time
	value     githubSearchScope
}

type githubSearchScope struct {
	login      string
	qualifiers []string
}

func NewGitHub() *GitHub {
	return &GitHub{env: defaultCLIEnv()}
}

func (gitHub *GitHub) ID() string {
	return "github"
}

func (gitHub *GitHub) Probe(ctx context.Context) Probe {
	if err := gitHub.env.require(ctx, "gh"); err != nil {
		return Probe{
			Provider:       gitHub.ID(),
			Status:         StatusCLIMissing,
			Hint:           "GitHub CLI not found. Install from https://cli.github.com",
			InstallDocsURL: githubInstallDocs,
			AuthDocsURL:    githubAuthDocs,
		}
	}
	if _, err := gitHub.env.run(ctx, "gh", "--version"); err != nil {
		return Probe{
			Provider:       gitHub.ID(),
			Status:         StatusCLIMissing,
			Hint:           "GitHub CLI is not usable. Reinstall from https://cli.github.com",
			InstallDocsURL: githubInstallDocs,
			AuthDocsURL:    githubAuthDocs,
		}
	}
	if _, err := gitHub.env.run(ctx, "gh", "auth", "status", "--hostname", "github.com"); err != nil {
		return Probe{
			Provider:       gitHub.ID(),
			Status:         StatusNotAuthenticated,
			Hint:           "Not signed in to GitHub CLI. Run: gh auth login",
			InstallDocsURL: githubInstallDocs,
			AuthDocsURL:    githubAuthDocs,
		}
	}
	return Probe{
		Provider:       gitHub.ID(),
		Status:         StatusAvailable,
		InstallDocsURL: githubInstallDocs,
		AuthDocsURL:    githubAuthDocs,
	}
}

type githubSearchResponse struct {
	TotalCount *int                `json:"total_count"`
	Items      []githubIssueRecord `json:"items"`
}

type githubIssueRecord struct {
	Number        int        `json:"number"`
	State         string     `json:"state"`
	Title         string     `json:"title"`
	HTMLURL       string     `json:"html_url"`
	UpdatedAt     string     `json:"updated_at"`
	Repository    githubRepo `json:"repository"`
	RepositoryURL string     `json:"repository_url"`
}

type githubRepo struct {
	FullName string `json:"full_name"`
}

type githubUser struct {
	Login string `json:"login"`
}

type githubOrganization struct {
	Login string `json:"login"`
}

func (gitHub *GitHub) Search(ctx context.Context, query string, limit int) ([]WorkItem, error) {
	scope, err := gitHub.searchScope(ctx)
	if err != nil {
		return nil, fmt.Errorf("resolve github search scope: %w", err)
	}

	terms := []string{"is:issue"}
	terms = append(terms, scope.qualifiers...)
	if trimmed := strings.TrimSpace(query); trimmed != "" {
		terms = append(terms, trimmed)
	}

	output, err := gitHub.env.run(
		ctx,
		"gh",
		"api",
		"-X",
		"GET",
		"search/issues",
		"-f",
		"q="+strings.Join(terms, " "),
		"-f",
		"sort=updated",
		"-f",
		"order=desc",
		"-f",
		"per_page="+strconv.Itoa(limit),
	)
	if err != nil {
		return nil, fmt.Errorf("gh api search/issues: %w", err)
	}

	var response githubSearchResponse
	if err := json.Unmarshal([]byte(output), &response); err != nil {
		return nil, fmt.Errorf("parse gh api output: %w", err)
	}

	items := make([]WorkItem, 0, len(response.Items))
	for _, issue := range response.Items {
		items = append(items, WorkItem{
			Key:             gitHub.issueKey(issue),
			Title:           issue.Title,
			URL:             issue.HTMLURL,
			Status:          strings.ToLower(issue.State),
			Type:            "issue",
			Project:         gitHub.repositoryName(issue),
			UpdatedAtUnixMs: parseRFC3339Ms(issue.UpdatedAt),
		})
	}
	return items, nil
}

func (gitHub *GitHub) searchScope(ctx context.Context) (githubSearchScope, error) {
	now := time.Now()
	if gitHub.scopeCache.expiresAt.After(now) {
		return gitHub.scopeCache.value, nil
	}

	userOutput, err := gitHub.env.run(ctx, "gh", "api", "user")
	if err != nil {
		return githubSearchScope{}, fmt.Errorf("gh api user: %w", err)
	}
	var user githubUser
	if err := json.Unmarshal([]byte(userOutput), &user); err != nil {
		return githubSearchScope{}, fmt.Errorf("parse gh user: %w", err)
	}
	if strings.TrimSpace(user.Login) == "" {
		return githubSearchScope{}, fmt.Errorf("gh api user returned no login")
	}

	orgOutput, err := gitHub.env.run(ctx, "gh", "api", "--paginate", "--slurp", "user/orgs")
	if err != nil {
		return githubSearchScope{}, fmt.Errorf("gh api user/orgs: %w", err)
	}
	var organizationPages []json.RawMessage
	if err := json.Unmarshal([]byte(orgOutput), &organizationPages); err != nil {
		return githubSearchScope{}, fmt.Errorf("parse gh orgs: %w", err)
	}

	qualifiers := []string{"user:" + user.Login}
	for _, page := range organizationPages {
		var organizations []githubOrganization
		if err := json.Unmarshal(page, &organizations); err != nil {
			continue
		}
		for _, organization := range organizations {
			if login := strings.TrimSpace(organization.Login); login != "" {
				qualifiers = append(qualifiers, "org:"+login)
			}
		}
	}

	scope := githubSearchScope{login: user.Login, qualifiers: qualifiers}
	gitHub.scopeCache = githubScopeCache{expiresAt: now.Add(githubScopeCacheTTL), value: scope}
	return scope, nil
}

func (gitHub *GitHub) issueKey(issue githubIssueRecord) string {
	name := gitHub.repositoryName(issue)
	return fmt.Sprintf("%s#%d", name, issue.Number)
}

func (gitHub *GitHub) repositoryName(issue githubIssueRecord) string {
	if name := strings.TrimSpace(issue.Repository.FullName); name != "" {
		return name
	}
	if matches := repositoryURLPattern.FindStringSubmatch(issue.RepositoryURL); len(matches) == 2 {
		return matches[1]
	}
	if matches := htmlURLPattern.FindStringSubmatch(issue.HTMLURL); len(matches) == 2 {
		return matches[1]
	}
	return "github"
}

func parseRFC3339Ms(value string) int64 {
	parsed, err := time.Parse(time.RFC3339, strings.TrimSpace(value))
	if err != nil {
		return 0
	}
	return parsed.UnixMilli()
}
