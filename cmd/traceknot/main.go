package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"traceknot/internal/cli"
	"traceknot/internal/server"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "uninstall":
			os.Exit(cli.RunUninstall(os.Args[2:]))
		case "claim":
			os.Exit(cli.RunClaim(os.Args[2:]))
		case "help", "-h", "--help":
			cli.PrintUsage()
			os.Exit(0)
		default:
			fmt.Fprintln(os.Stderr, "unknown command: "+os.Args[1])
			cli.PrintUsage()
			os.Exit(2)
		}
		return
	}

	if cli.IsInteractive() {
		os.Exit(cli.RunMenu())
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	if err := server.Run(defaultDBPath(), "127.0.0.1:4318", logger); err != nil {
		logger.Error("server failed", "error", err)
		os.Exit(1)
	}
}

func defaultDBPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "traceknot.sqlite"
	}
	return filepath.Join(home, ".traceknot", "telemetry.sqlite")
}
