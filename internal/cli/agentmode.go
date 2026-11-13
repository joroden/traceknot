package cli

import (
	"context"
	"os"
	"path"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v4/process"
)

// Hook stdin and stdout are pipes even in interactive sessions, so inspect the
// calling agent's arguments instead of the hook's terminal status.
func nonInteractiveAgent(ctx context.Context, agent string) bool {
	if agent != "claude" && agent != "codex" && agent != "copilot" {
		return false
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	return nonInteractiveAncestor(agent, int32(os.Getppid()), func(pid int32) ([]string, int32, error) {
		proc, err := process.NewProcessWithContext(ctx, pid)
		if err != nil {
			return nil, 0, err
		}
		args, err := proc.CmdlineSliceWithContext(ctx)
		if err != nil {
			return nil, 0, err
		}
		parent, err := proc.PpidWithContext(ctx)
		return args, parent, err
	})
}

func nonInteractiveAncestor(agent string, pid int32, readProcess func(int32) ([]string, int32, error)) bool {
	for depth := 0; pid > 0 && depth < 32; depth++ {
		args, parent, err := readProcess(pid)
		if err != nil {
			return false
		}
		if agentArgs, ok := agentCommandArgs(agent, args); ok {
			// Stop at the nearest agent. An interactive child of a scripted run
			// must not inherit the outer agent's non-interactive mode.
			return nonInteractiveArgs(agent, agentArgs)
		}
		if parent == pid {
			return false
		}
		pid = parent
	}
	return false
}

func agentCommandArgs(agent string, args []string) ([]string, bool) {
	if len(args) == 0 {
		return nil, false
	}
	// Copilot also uses its executable as the Node host for extensions. This
	// bootstrap process has no CLI mode flags; keep walking to its parent CLI.
	if agent == "copilot" && len(args) > 1 && commandBase(args[1]) == "extension_bootstrap.mjs" {
		return nil, false
	}
	command := strings.TrimSuffix(commandBase(args[0]), ".exe")
	if command == agent {
		return args[1:], true
	}
	// npm installations run a package entry point under Node or Bun.
	if (command == "node" || command == "nodejs" || command == "bun") && len(args) > 1 {
		script := strings.ToLower(strings.ReplaceAll(args[1], "\\", "/"))
		if commandBase(script) == agent {
			return args[2:], true
		}
		packages := map[string]string{
			"claude":  "/@anthropic-ai/claude-code/",
			"codex":   "/@openai/codex/",
			"copilot": "/@github/copilot/",
		}
		if pkg := packages[agent]; pkg != "" && strings.Contains(script, pkg) {
			return args[2:], true
		}
	}
	return nil, false
}

func commandBase(command string) string {
	return strings.ToLower(path.Base(strings.ReplaceAll(command, "\\", "/")))
}

func nonInteractiveArgs(agent string, args []string) bool {
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--" {
			return false
		}
		switch agent {
		case "claude":
			if arg == "--print" || claudePrintFlag(arg) {
				return true
			}
		case "copilot":
			if arg == "--prompt" || strings.HasPrefix(arg, "--prompt=") ||
				(strings.HasPrefix(arg, "-p") && !strings.HasPrefix(arg, "--")) {
				return true
			}
		case "codex":
			if !strings.HasPrefix(arg, "-") {
				// Only the subcommand selects exec mode. Prompt text and
				// arguments to interactive resume/fork commands do not.
				return arg == "exec" || arg == "e" || arg == "review"
			}
		}
		if agentOptionTakesValue(agent, arg) {
			i++
		}
	}
	return false
}

func claudePrintFlag(arg string) bool {
	if !strings.HasPrefix(arg, "-") || strings.HasPrefix(arg, "--") {
		return false
	}
	for _, flag := range arg[1:] {
		switch flag {
		case 'p':
			return true
		case 'c', 'v': // Boolean flags can be combined, e.g. -cp.
		default:
			return false
		}
	}
	return false
}

func agentOptionTakesValue(agent string, arg string) bool {
	// Values can themselves look like flags or subcommands, e.g.
	// claude --system-prompt "-p" or codex --profile review.
	if strings.Contains(arg, "=") {
		return false
	}
	switch arg {
	case "--model", "--agent", "--add-dir":
		return true
	}
	switch agent {
	case "codex":
		switch arg {
		case "-c", "--config", "--enable", "--disable", "-i", "--image", "-m",
			"-p", "--profile", "-s", "--sandbox", "-C", "--cd", "-a", "--ask-for-approval",
			"--local-provider", "--remote", "--remote-auth-token-env":
			return true
		}
	case "claude":
		switch arg {
		case "--system-prompt", "--append-system-prompt", "--system-prompt-file", "--append-system-prompt-file",
			"--settings", "--setting-sources", "--permission-mode", "--permission-prompt-tool",
			"--allowedTools", "--allowed-tools", "--disallowedTools", "--disallowed-tools", "--tools",
			"--mcp-config", "--agents", "--plugin-dir", "--betas", "--file", "--debug-file",
			"--effort", "--fallback-model", "--input-format", "--output-format", "--json-schema",
			"--max-budget-usd", "--max-turns", "--session-id", "-n", "--name":
			return true
		}
	case "copilot":
		switch arg {
		case "-i", "--interactive", "--additional-mcp-config", "--add-github-mcp-tool", "--add-github-mcp-toolset",
			"--allow-tool", "--allow-url", "--deny-tool", "--deny-url", "--available-tools", "--excluded-tools",
			"--attachment", "--config-dir", "--disable-mcp-server", "--enable-mcp-server", "--log-dir", "--log-level",
			"--max-autopilot-continues", "--mode", "--output-format", "--reasoning-effort",
			"--stream", "--usage-output-file":
			return true
		}
	}
	return false
}
