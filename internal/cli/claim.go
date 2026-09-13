package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"traceknot/internal/config"
)

func RunClaim(args []string) int {
	flags := flag.NewFlagSet("claim", flag.ExitOnError)
	server := flags.String("server", "http://127.0.0.1:4318", "daemon base URL")
	agent := flags.String("agent", "", "calling agent: claude, codex, or copilot")
	_ = flags.Parse(args)
	ctx := context.Background()

	payload, err := io.ReadAll(os.Stdin)
	if err != nil {
		logClaim("claim invoked, no stdin payload")
		return 0
	}
	sessionID, ok := hookSessionID(payload)
	if !ok {
		logClaim("claim invoked, no session payload: " + string(payload))
		return 0
	}
	prompt := hookPrompt(payload)
	promptPreview := prompt
	if len(promptPreview) > 40 {
		promptPreview = promptPreview[:40]
	}
	turnKey := hookTurnKey(payload)
	if !daemonHealthy(ctx, *server) {
		logClaim("claim: session " + sessionID + ", daemon not running, skipped")
		return 0
	}
	status, err := offerPicker(ctx, *server, sessionID, turnKey, prompt)
	if err != nil {
		logClaim("claim: session " + sessionID + ", offer failed: " + err.Error())
		return 0
	}
	silent := *agent == "copilot" && turnKey != ""
	switch status {
	case "claimed":
		return 0
	case "offered":
		logClaim("claim: opening picker for session " + sessionID + ", prompt=" + fmt.Sprintf("%q", promptPreview))
		outcome, err := runSelectFlow(ctx, *server, sessionID)
		if err != nil {
			logClaim("claim: session " + sessionID + ", picker failed: " + err.Error())
		}
		return decideBlock(sessionID, *agent, err == nil && outcome.Status == "claimed", silent)
	case "pending", "skipped":
		logClaim("claim: session " + sessionID + " reusing existing offer (" + status + "), prompt=" + fmt.Sprintf("%q", promptPreview))
		outcome, err := waitForOutcome(ctx, *server, sessionID)
		if err != nil {
			logClaim("claim: session " + sessionID + ", wait for outcome failed: " + err.Error())
		}
		return decideBlock(sessionID, *agent, err == nil && outcome.Status == "claimed", silent)
	default:
		return 0
	}
}

func decideBlock(sessionID string, agent string, claimed bool, silent bool) int {
	if claimed || agent == "" {
		return 0
	}
	if !config.Load().RequireWorkItem {
		return 0
	}
	logClaim("claim: session " + sessionID + " blocked (" + agent + "), required mode on")
	if silent {
		return 2
	}
	reason := "No work item was assigned to this session and assignment is required. The run was not allowed to continue."
	if agent == "codex" {
		decision, err := json.Marshal(map[string]string{"decision": "block", "reason": reason})
		if err == nil {
			fmt.Println(string(decision))
		}
		return 0
	}
	fmt.Fprintln(os.Stderr, reason)
	return 2
}

func offerPicker(ctx context.Context, server string, sessionID string, turnKey string, prompt string) (string, error) {
	body, err := json.Marshal(map[string]string{"session_id": sessionID, "turn_key": turnKey, "prompt": prompt})
	if err != nil {
		return "", err
	}
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		strings.TrimRight(server, "/")+"/api/v1/picker/offer",
		bytes.NewReader(body),
	)
	if err != nil {
		return "", err
	}
	request.Header.Set("content-type", "application/json")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	var result struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.Status, nil
}

func logClaim(message string) {
	logPath := filepath.Join(mustHome(), ".traceknot", "claim.log")
	if err := os.MkdirAll(filepath.Dir(logPath), 0o755); err != nil {
		return
	}
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return
	}
	defer file.Close()
	_, _ = file.WriteString(time.Now().Format("2006-01-02 15:04:05") + "  " + message + "\n")
}
