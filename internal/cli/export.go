package cli

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"traceknot/internal/export"
	"traceknot/internal/export/render"
	"traceknot/internal/store"
)

func RunExport(args []string) int {
	if len(args) == 0 {
		printExportUsage()
		return 2
	}
	switch args[0] {
	case "session":
		return runExportSession(args[1:])
	case "work-item":
		return runExportWorkItem(args[1:])
	case "help", "-h", "--help":
		printExportUsage()
		return 0
	default:
		fmt.Fprintln(os.Stderr, "export: unknown subcommand: "+args[0])
		printExportUsage()
		return 2
	}
}

func printExportUsage() {
	fmt.Println("Usage:")
	fmt.Println("  traceknot export session --id SESSION_ID [--out DIR]")
	fmt.Println("  traceknot export work-item --provider PROVIDER --key KEY [--out DIR]")
	fmt.Println()
	fmt.Println("Writes a session (or every session claimed against a work item) to a")
	fmt.Println("directory of markdown files for an agent to read through. Without --out,")
	fmt.Println("a fresh OS temp directory is created and its path is printed on success.")
}

func runExportSession(args []string) int {
	flags := flag.NewFlagSet("export session", flag.ExitOnError)
	sessionID := flags.String("id", "", "session id to export")
	out := flags.String("out", "", "output directory (default: a new temp dir)")
	_ = flags.Parse(args)

	if *sessionID == "" {
		fmt.Fprintln(os.Stderr, "export session: --id is required")
		return 2
	}

	st, err := store.Open(exportDBPath())
	if err != nil {
		fmt.Fprintln(os.Stderr, "export session: open database:", err)
		return 1
	}
	defer st.Close()

	sess, err := export.LoadSession(context.Background(), st, *sessionID)
	if err != nil {
		fmt.Fprintln(os.Stderr, "export session:", err)
		return 1
	}

	outDir, err := resolveExportDir(*out, "traceknot-export-session-*")
	if err != nil {
		fmt.Fprintln(os.Stderr, "export session:", err)
		return 1
	}

	if err := render.WriteSession(sess, outDir, nil); err != nil {
		fmt.Fprintln(os.Stderr, "export session:", err)
		return 1
	}

	fmt.Println(outDir)
	return 0
}

func runExportWorkItem(args []string) int {
	flags := flag.NewFlagSet("export work-item", flag.ExitOnError)
	provider := flags.String("provider", "", "work item provider (github, gitlab, jira, custom)")
	key := flags.String("key", "", "work item key")
	out := flags.String("out", "", "output directory (default: a new temp dir)")
	_ = flags.Parse(args)

	if *provider == "" || *key == "" {
		fmt.Fprintln(os.Stderr, "export work-item: --provider and --key are required")
		return 2
	}

	st, err := store.Open(exportDBPath())
	if err != nil {
		fmt.Fprintln(os.Stderr, "export work-item: open database:", err)
		return 1
	}
	defer st.Close()

	workItem, err := export.LoadWorkItem(context.Background(), st, *key, *provider)
	if err != nil {
		fmt.Fprintln(os.Stderr, "export work-item:", err)
		return 1
	}

	outDir, err := resolveExportDir(*out, "traceknot-export-workitem-*")
	if err != nil {
		fmt.Fprintln(os.Stderr, "export work-item:", err)
		return 1
	}

	if err := render.WriteWorkItem(workItem, outDir); err != nil {
		fmt.Fprintln(os.Stderr, "export work-item:", err)
		return 1
	}

	fmt.Println(outDir)
	return 0
}

func resolveExportDir(out, tempPattern string) (string, error) {
	if out == "" {
		return os.MkdirTemp("", tempPattern)
	}
	if err := os.MkdirAll(out, 0o755); err != nil {
		return "", fmt.Errorf("create output directory: %w", err)
	}
	return out, nil
}

func exportDBPath() string {
	return filepath.Join(mustHome(), ".traceknot", "telemetry.sqlite")
}
