package providers

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

const maxCommandOutputBytes = 1 << 20

type cliEnv struct {
	lookPath func(ctx context.Context, name string) (string, error)
	run      func(ctx context.Context, name string, args ...string) (string, error)
}

func (environment cliEnv) require(ctx context.Context, name string) error {
	_, err := environment.lookPath(ctx, name)
	return err
}

type limitedBuffer struct {
	buffer bytes.Buffer
	limit  int
}

func (output *limitedBuffer) Write(data []byte) (int, error) {
	if output.buffer.Len()+len(data) > output.limit {
		return 0, errors.New("command output exceeds the 1 MiB limit")
	}
	return output.buffer.Write(data)
}

func defaultCLIEnv() cliEnv {
	return cliEnv{
		lookPath: func(ctx context.Context, name string) (string, error) {
			return exec.LookPath(name)
		},
		run: func(ctx context.Context, name string, args ...string) (string, error) {
			runCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
			defer cancel()
			command := exec.CommandContext(runCtx, name, args...)
			var stdout limitedBuffer
			stdout.limit = maxCommandOutputBytes
			var stderr strings.Builder
			command.Stdout = &stdout
			command.Stderr = &stderr
			if err := command.Run(); err != nil {
				detail := strings.TrimSpace(stderr.String())
				if detail == "" {
					detail = err.Error()
				}
				return "", fmt.Errorf("run %s: %s", name, detail)
			}
			return stdout.buffer.String(), nil
		},
	}
}
