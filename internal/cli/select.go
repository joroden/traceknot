package cli

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func RunSelect(args []string) int {
	var (
		serverAddr = flag.NewFlagSet("select", flag.ExitOnError)
		server     = serverAddr.String("server", "http://127.0.0.1:4318", "daemon base URL")
		sessionID  = serverAddr.String("session-id", "", "session id to attach the claim to")
		prompt     = serverAddr.String("prompt", "", "user prompt for the picker context banner")
	)
	_ = serverAddr.Parse(args)
	ctx := context.Background()

	if *sessionID != "" && *prompt != "" {
		if err := postPrompt(ctx, *server, *sessionID, *prompt); err != nil {
			fmt.Fprintln(os.Stderr, "select: cannot store prompt context:", err)
			return 2
		}
	}

	query := url.Values{}
	if *sessionID != "" {
		query.Set("session", *sessionID)
	}
	pickerURL := strings.TrimRight(*server, "/") + "/select"
	if len(query) > 0 {
		pickerURL += "?" + query.Encode()
	}

	if err := openBrowser(pickerURL); err != nil {
		fmt.Fprintln(os.Stderr, "select: cannot open browser:", err)
		return 2
	}

	if *sessionID == "" {
		return 0
	}

	outcome, err := waitForOutcome(ctx, *server, *sessionID)
	if err != nil {
		fmt.Fprintln(os.Stderr, "select:", err)
		return 2
	}
	switch outcome.Status {
	case "claimed":
		return 0
	case "skipped":
		return 1
	default:
		return 2
	}
}

func postPrompt(ctx context.Context, server string, sessionID string, prompt string) error {
	body, err := json.Marshal(map[string]string{
		"session_id": sessionID,
		"prompt":     prompt,
	})
	if err != nil {
		return err
	}
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		strings.TrimRight(server, "/")+"/api/v1/picker/context",
		strings.NewReader(string(body)),
	)
	if err != nil {
		return err
	}
	request.Header.Set("content-type", "application/json")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("daemon returned %s", response.Status)
	}
	return nil
}

type selectOutcome struct {
	Status      string `json:"status"`
	WorkItemKey string `json:"work_item_key,omitempty"`
}

func waitForOutcome(ctx context.Context, server string, sessionID string) (*selectOutcome, error) {
	endpoint := strings.TrimRight(server, "/") + "/api/v1/picker/outcome?session_id=" + url.QueryEscape(sessionID)
	client := &http.Client{Timeout: 5 * time.Second}
	deadline := time.Now().Add(5 * time.Minute)

	for time.Now().Before(deadline) {
		attemptCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		request, err := http.NewRequestWithContext(attemptCtx, http.MethodGet, endpoint, nil)
		if err != nil {
			cancel()
			return nil, err
		}
		response, err := client.Do(request)
		if err == nil {
			var outcome selectOutcome
			decodeErr := json.NewDecoder(response.Body).Decode(&outcome)
			_ = response.Body.Close()
			if decodeErr == nil && outcome.Status != "pending" {
				cancel()
				return &outcome, nil
			}
		}
		cancel()
		time.Sleep(2 * time.Second)
	}
	return nil, fmt.Errorf("timed out waiting for a picker outcome")
}
