package cli

import "traceknot/internal/platform"

func canOpenBrowser() bool {
	return platform.Current.CanOpenBrowser()
}

func openBrowser(target string) error {
	return platform.Current.OpenBrowser(target)
}

func openWithDefaultHandler(target string) error {
	return platform.Current.OpenDefaultHandler(target)
}
