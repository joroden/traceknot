package cli

import "context"

func RunStop() int {
	stopDaemon(context.Background(), defaultServerURL)
	return 0
}
