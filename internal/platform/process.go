package platform

import (
	"context"
	"net"
	"strconv"
	"strings"
	"time"

	gopsnet "github.com/shirou/gopsutil/v4/net"
	"github.com/shirou/gopsutil/v4/process"
)

const daemonProcessName = "traceknot"

func isTraceknotProcess(ctx context.Context, proc *process.Process) bool {
	name, err := proc.NameWithContext(ctx)
	if err != nil {
		return false
	}
	return strings.TrimSuffix(strings.ToLower(name), ".exe") == daemonProcessName
}

func killProcess(ctx context.Context, pidStr string) {
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return
	}
	proc, err := process.NewProcessWithContext(ctx, int32(pid))
	if err != nil {
		return
	}
	if !isTraceknotProcess(ctx, proc) {
		return
	}
	_ = proc.KillWithContext(ctx)
}

func freePort(ctx context.Context, port string) {
	portNum, err := strconv.ParseUint(port, 10, 32)
	if err != nil {
		return
	}
	conns, err := gopsnet.ConnectionsWithContext(ctx, "tcp")
	if err == nil {
		killed := make(map[int32]bool)
		for _, conn := range conns {
			if conn.Status != "LISTEN" || conn.Laddr.Port != uint32(portNum) || conn.Pid <= 0 || killed[conn.Pid] {
				continue
			}
			if proc, err := process.NewProcessWithContext(ctx, conn.Pid); err == nil {
				_ = proc.KillWithContext(ctx)
				killed[conn.Pid] = true
			}
		}
	}
	waitPortFree(port)
}

func waitPortFree(port string) {
	deadline := time.Now().Add(2 * time.Second)
	for {
		ln, err := net.Listen("tcp", "127.0.0.1:"+port)
		if err == nil {
			_ = ln.Close()
			return
		}
		if time.Now().After(deadline) {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
}
