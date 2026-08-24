package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

const (
	jiraInstallDocs = "https://developer.atlassian.com/cloud/acli/guides/install-acli/"
	jiraAuthDocs    = "https://developer.atlassian.com/cloud/acli/guides/how-to-get-started/"
)

type Jira struct {
	env cliEnv
}

func NewJira() *Jira {
	return &Jira{env: defaultCLIEnv()}
}

func (jira *Jira) ID() string {
	return "jira"
}

func (jira *Jira) Probe(ctx context.Context) Probe {
	if err := jira.env.require(ctx, "acli"); err != nil {
		return Probe{
			Provider:       jira.ID(),
			Status:         StatusCLIMissing,
			Hint:           "Atlassian CLI (acli) not found. Install from " + jiraInstallDocs,
			InstallDocsURL: jiraInstallDocs,
			AuthDocsURL:    jiraAuthDocs,
		}
	}
	if _, err := jira.env.run(ctx, "acli", "--version"); err != nil {
		return Probe{
			Provider:       jira.ID(),
			Status:         StatusCLIMissing,
			Hint:           "Atlassian CLI (acli) is not usable. Reinstall from " + jiraInstallDocs,
			InstallDocsURL: jiraInstallDocs,
			AuthDocsURL:    jiraAuthDocs,
		}
	}
	if _, err := jira.env.run(ctx, "acli", "jira", "auth", "status"); err != nil {
		return Probe{
			Provider:       jira.ID(),
			Status:         StatusNotAuthenticated,
			Hint:           "Not authenticated with Jira. Run: acli jira auth login",
			InstallDocsURL: jiraInstallDocs,
			AuthDocsURL:    jiraAuthDocs,
		}
	}
	return Probe{
		Provider:       jira.ID(),
		Status:         StatusAvailable,
		InstallDocsURL: jiraInstallDocs,
		AuthDocsURL:    jiraAuthDocs,
	}
}

func (jira *Jira) Search(ctx context.Context, query string, limit int) ([]WorkItem, error) {
	output, err := jira.env.run(
		ctx,
		"acli",
		"jira",
		"workitem",
		"search",
		"--jql",
		buildJiraJQL(query),
		"--limit",
		fmt.Sprintf("%d", limit),
		"--fields",
		"key,summary,status",
		"--json",
	)
	if err != nil {
		return nil, fmt.Errorf("acli jira workitem search: %w", err)
	}

	records, err := extractRecords(output)
	if err != nil {
		return nil, fmt.Errorf("parse acli jira workitem search output: %w", err)
	}

	items := make([]WorkItem, 0, 16)
	for _, record := range records {
		item := jira.normalize(record)
		if item == nil {
			continue
		}
		items = append(items, *item)
	}
	return items, nil
}

func buildJiraJQL(query string) string {
	trimmed := strings.TrimSpace(query)
	upper := strings.ToUpper(trimmed)
	if upper != "" && jiraKeyPattern.MatchString(upper) {
		return fmt.Sprintf("project IS NOT EMPTY AND key = %s ORDER BY updated DESC", upper)
	}
	if upper == "" {
		return "project IS NOT EMPTY ORDER BY updated DESC"
	}
	return fmt.Sprintf(
		"project IS NOT EMPTY AND text ~ %s ORDER BY updated DESC",
		quoteJiraJQL(trimmed),
	)
}

func quoteJiraJQL(value string) string {
	return `"` + strings.NewReplacer(`\`, `\\`, `"`, `\"`).Replace(value) + `"`
}

func extractRecords(output string) ([]map[string]any, error) {
	var value any
	if err := json.Unmarshal([]byte(output), &value); err != nil {
		return nil, err
	}
	return collectRecords(value), nil
}

func collectRecords(value any) []map[string]any {
	switch typed := value.(type) {
	case []any:
		output := make([]map[string]any, 0, len(typed))
		for _, entry := range typed {
			output = append(output, collectRecords(entry)...)
		}
		return output
	case map[string]any:
		if hasRecordShape(typed) {
			return []map[string]any{typed}
		}
		output := make([]map[string]any, 0, 8)
		for _, key := range []string{"issues", "workItems", "items", "values", "data"} {
			output = append(output, collectRecords(typed[key])...)
		}
		return output
	default:
		return nil
	}
}

func hasRecordShape(record map[string]any) bool {
	_, hasKey := record["key"].(string)
	if !hasKey {
		if _, ok := record["issueKey"].(string); !ok {
			return false
		}
	}
	return fieldString(record, "summary", "title") != ""
}

func (jira *Jira) normalize(record map[string]any) *WorkItem {
	key := firstString(record, "key", "issueKey")
	if key == "" {
		return nil
	}
	summary := fieldString(record, "summary", "title")
	status := firstString(record, "status")
	if status == "" {
		status = nestedString(record, "fields", "status", "name")
	}
	url := firstString(record, "url", "browseUrl", "self")
	updated := firstString(record, "updated", "updatedAt")
	item := &WorkItem{
		Key:             key,
		Title:           summary,
		Status:          status,
		URL:             jiraBrowseURL(url, key),
		Project:         projectFromKey(key),
		UpdatedAtUnixMs: parseRFC3339Ms(updated),
	}
	if item.Title == "" {
		item.Title = key
	}
	return item
}

func jiraBrowseURL(value string, key string) string {
	if matches := jiraSelfPattern.FindStringSubmatch(strings.TrimSpace(value)); len(matches) == 3 {
		return matches[1] + "/browse/" + matches[2]
	}
	return value
}

func nestedString(record map[string]any, outer string, inner string, leaf string) string {
	outerValue, ok := record[outer].(map[string]any)
	if !ok {
		return ""
	}
	innerValue, ok := outerValue[inner].(map[string]any)
	if !ok {
		return ""
	}
	return asString(innerValue[leaf])
}

func projectFromKey(key string) string {
	parts := strings.SplitN(key, "-", 2)
	if len(parts) != 2 || parts[0] == "" {
		return ""
	}
	return parts[0]
}

func asString(value any) string {
	if text, ok := value.(string); ok {
		return strings.TrimSpace(text)
	}
	return ""
}

func firstString(record map[string]any, keys ...string) string {
	for _, key := range keys {
		if text := asString(record[key]); text != "" {
			return text
		}
	}
	return ""
}

func fieldString(record map[string]any, keys ...string) string {
	if text := firstString(record, keys...); text != "" {
		return text
	}
	fields, ok := record["fields"].(map[string]any)
	if !ok {
		return ""
	}
	return firstString(fields, keys...)
}
