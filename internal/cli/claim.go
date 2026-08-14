package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func RunClaim(args []string) int {
	flags := flag.NewFlagSet("claim", flag.ExitOnError)
	server := flags.String("server", "http://127.0.0.1:4318", "daemon base URL")
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
	if !daemonHealthy(ctx, *server) {
		logClaim("claim: session " + sessionID + ", daemon not running, skipped")
		return 0
	}
	status, err := offerPicker(ctx, *server, sessionID)
	if err != nil {
		logClaim("claim: session " + sessionID + ", offer failed: " + err.Error())
		return 0
	}
	if status != "offered" {
		logClaim("claim: session " + sessionID + " already has an outcome (" + status + "), skipped")
		return 0
	}
	logClaim("claim: opening picker for session " + sessionID)
	_ = RunSelect([]string{"--server", *server, "--session-id", sessionID})
	return 0
}

func hookSessionID(payload []byte) (string, bool) {
	var fields map[string]any
	if err := json.Unmarshal(payload, &fields); err != nil {
		return "", false
	}
	sessionID, _ := fields["session_id"].(string)
	if sessionID == "" {
		sessionID, _ = fields["sessionId"].(string)
	}
	return sessionID, sessionID != ""
}

func offerPicker(ctx context.Context, server string, sessionID string) (string, error) {
	body, err := json.Marshal(map[string]string{"session_id": sessionID})
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
