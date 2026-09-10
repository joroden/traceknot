package cli

import "fmt"

func PrintUsage() {
	fmt.Println("traceknot - local telemetry collector for AI coding agents")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  traceknot            open settings (autostart, hooks, work items, skills)")
	fmt.Println("                       in your browser, or run the daemon otherwise")
	fmt.Println("  traceknot uninstall  remove traceknot (keeps local data in ~/.traceknot)")
	fmt.Println("  traceknot export     write a session or work item to markdown for an agent to read")
	fmt.Println("  traceknot help       show this help")
}
