package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	gitlabInstallDocs   = "https://gitlab.com/gitlab-org/cli#installation"
	gitlabAuthDocs      = "https://docs.gitlab.com/cli/authentication/"
	gitlabGroupCacheTTL = 5 * time.Minute
)

var gitlabPersonalScopes = []string{"created_by_me", "assigned_to_me"}

type GitLab struct {
	env        cliEnv
	groupCache gitlabGroupCache
}

type gitlabGroupCache struct {
	expiresAt time.Time
	groupIDs  []int
}

func NewGitLab() *GitLab {
	return &GitLab{env: defaultCLIEnv()}
}

func (gitLab *GitLab) ID() string {
	return "gitlab"
}

func (gitLab *GitLab) Probe(ctx context.Context) Probe {
	if err := gitLab.env.require(ctx, "glab"); err != nil {
		return Probe{
			Provider:       gitLab.ID(),
			Status:         StatusCLIMissing,
			Hint:           "GitLab CLI not found. Install from https://gitlab.com/gitlab-org/cli",
			InstallDocsURL: gitlabInstallDocs,
			AuthDocsURL:    gitlabAuthDocs,
		}
	}
	if _, err := gitLab.env.run(ctx, "glab", "--version"); err != nil {
		return Probe{
			Provider:       gitLab.ID(),
			Status:         StatusCLIMissing,
			Hint:           "GitLab CLI is not usable. Reinstall from https://gitlab.com/gitlab-org/cli",
			InstallDocsURL: gitlabInstallDocs,
			AuthDocsURL:    gitlabAuthDocs,
		}
	}
	if _, err := gitLab.env.run(ctx, "glab", "auth", "status", "--hostname", "gitlab.com"); err != nil {
		return Probe{
			Provider:       gitLab.ID(),
			Status:         StatusNotAuthenticated,
			Hint:           "Not signed in to GitLab CLI. Run: glab auth login",
			InstallDocsURL: gitlabInstallDocs,
			AuthDocsURL:    gitlabAuthDocs,
		}
	}
	return Probe{
		Provider:       gitLab.ID(),
		Status:         StatusAvailable,
		InstallDocsURL: gitlabInstallDocs,
		AuthDocsURL:    gitlabAuthDocs,
	}
}

type gitlabIssueRecord struct {
	State      string           `json:"state"`
	Title      string           `json:"title"`
	WebURL     string           `json:"web_url"`
	UpdatedAt  string           `json:"updated_at"`
	References gitlabReferences `json:"references"`
}

type gitlabReferences struct {
	Full string `json:"full"`
}

type gitlabGroupRecord struct {
	ID int `json:"id"`
}

func (gitLab *GitLab) Search(ctx context.Context, query string, limit int) ([]WorkItem, error) {
	trimmedQuery := strings.TrimSpace(query)
	items := make(map[string]WorkItem)

	for _, scope := range gitlabPersonalScopes {
		records, err := gitLab.fetchIssues(ctx, "issues", trimmedQuery, limit, "scope="+scope)
		if err != nil {
			return nil, err
		}
		gitLab.mergeIssues(items, records)
	}

	for _, groupID := range gitLab.myGroupIDs(ctx) {
		endpoint := "groups/" + strconv.Itoa(groupID) + "/issues"
		if records, err := gitLab.fetchIssues(ctx, endpoint, trimmedQuery, limit); err == nil {
			gitLab.mergeIssues(items, records)
		}
	}

	results := make([]WorkItem, 0, len(items))
	for _, item := range items {
		results = append(results, item)
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].UpdatedAtUnixMs > results[j].UpdatedAtUnixMs
	})
	if len(results) > limit {
		results = results[:limit]
	}
	return results, nil
}

func (gitLab *GitLab) myGroupIDs(ctx context.Context) []int {
	now := time.Now()
	if gitLab.groupCache.expiresAt.After(now) {
		return gitLab.groupCache.groupIDs
	}

	output, err := gitLab.env.run(
		ctx, "glab", "api", "-X", "GET", "groups",
		"-f", "all_available=false",
		"-f", "top_level_only=true",
		"-f", "per_page=100",
		"--paginate",
	)
	if err != nil {
		return nil
	}
	var groups []gitlabGroupRecord
	if err := json.Unmarshal([]byte(output), &groups); err != nil {
		return nil
	}

	groupIDs := make([]int, 0, len(groups))
	for _, group := range groups {
		groupIDs = append(groupIDs, group.ID)
	}
	gitLab.groupCache = gitlabGroupCache{expiresAt: now.Add(gitlabGroupCacheTTL), groupIDs: groupIDs}
	return groupIDs
}

func (gitLab *GitLab) fetchIssues(ctx context.Context, endpoint string, query string, limit int, extraFilters ...string) ([]gitlabIssueRecord, error) {
	args := []string{
		"api", "-X", "GET", endpoint,
		"-f", "order_by=updated_at",
		"-f", "sort=desc",
		"-f", "per_page=" + strconv.Itoa(limit),
	}
	for _, filter := range extraFilters {
		args = append(args, "-f", filter)
	}
	if query != "" {
		args = append(args, "-f", "search="+query)
	}

	output, err := gitLab.env.run(ctx, "glab", args...)
	if err != nil {
		return nil, fmt.Errorf("glab api %s: %w", endpoint, err)
	}
	var records []gitlabIssueRecord
	if err := json.Unmarshal([]byte(output), &records); err != nil {
		return nil, fmt.Errorf("parse glab api %s output: %w", endpoint, err)
	}
	return records, nil
}

func (gitLab *GitLab) mergeIssues(items map[string]WorkItem, records []gitlabIssueRecord) {
	for _, issue := range records {
		if issue.References.Full == "" {
			continue
		}
		items[issue.References.Full] = WorkItem{
			Key:             issue.References.Full,
			Title:           issue.Title,
			URL:             issue.WebURL,
			Status:          strings.ToLower(issue.State),
			Type:            "issue",
			Project:         gitLab.projectName(issue),
			UpdatedAtUnixMs: parseRFC3339Ms(issue.UpdatedAt),
		}
	}
}

func (gitLab *GitLab) projectName(issue gitlabIssueRecord) string {
	full := issue.References.Full
	if index := strings.LastIndex(full, "#"); index > 0 {
		return full[:index]
	}
	return "gitlab"
}
