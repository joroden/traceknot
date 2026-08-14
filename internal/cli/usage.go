package cli

import "fmt"

func PrintUsage() {
	fmt.Println("traceknot - local telemetry collector for AI coding agents")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  traceknot            open the interactive menu (server, autostart, hooks)")
	fmt.Println("                       from a terminal, or run the daemon otherwise")
	fmt.Println("  traceknot uninstall  remove traceknot (keeps local data in ~/.traceknot)")
	fmt.Println("  traceknot help       show this help")
}
