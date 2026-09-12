//go:build windows

package windows

import (
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows/registry"
)

func ApplyEnv(vars map[string]string, binDir string) error {
	key, err := registry.OpenKey(registry.CURRENT_USER, `Environment`, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer key.Close()

	for name, value := range vars {
		if strings.EqualFold(name, "Path") {
			continue
		}
		if err := key.SetExpandStringValue(name, value); err != nil {
			return err
		}
	}

	path, err := readPath(key)
	if err != nil {
		return err
	}
	if !pathContains(path, binDir) {
		if err := key.SetExpandStringValue("Path", appendPathEntry(path, binDir)); err != nil {
			return err
		}
	}

	broadcastEnvChange()
	return nil
}

func RemoveEnv(names []string, binDir string) error {
	key, err := registry.OpenKey(registry.CURRENT_USER, `Environment`, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer key.Close()

	path, err := readPath(key)
	if err != nil {
		return err
	}
	if newPath := removePathEntry(path, binDir); newPath != path {
		if err := key.SetExpandStringValue("Path", newPath); err != nil {
			return err
		}
	}

	for _, name := range names {
		if strings.EqualFold(name, "Path") {
			continue
		}
		if err := key.DeleteValue(name); err != nil && err != registry.ErrNotExist {
			return err
		}
	}

	broadcastEnvChange()
	return nil
}

func readPath(key registry.Key) (string, error) {
	path, _, err := key.GetStringValue("Path")
	if err != nil {
		if err == registry.ErrNotExist {
			return "", nil
		}
		return "", err
	}
	return path, nil
}

func pathContains(path, entry string) bool {
	for _, part := range strings.Split(path, ";") {
		if strings.EqualFold(strings.TrimSpace(part), entry) {
			return true
		}
	}
	return false
}

func appendPathEntry(path, entry string) string {
	path = strings.TrimRight(strings.TrimSpace(path), ";")
	if path == "" {
		return entry
	}
	return path + ";" + entry
}

func removePathEntry(path, entry string) string {
	var kept []string
	for _, part := range strings.Split(path, ";") {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" || strings.EqualFold(trimmed, entry) {
			continue
		}
		kept = append(kept, part)
	}
	return strings.Join(kept, ";")
}

var (
	user32                  = syscall.NewLazyDLL("user32.dll")
	procSendMessageTimeoutW = user32.NewProc("SendMessageTimeoutW")
)

const (
	hwndBroadcast   = 0xffff
	wmSettingChange = 0x1a
	smtoAbortIfHung = 0x0002
)

func broadcastEnvChange() {
	envPtr, err := syscall.UTF16PtrFromString("Environment")
	if err != nil {
		return
	}
	var result uintptr
	_, _, _ = procSendMessageTimeoutW.Call(
		uintptr(hwndBroadcast),
		uintptr(wmSettingChange),
		0,
		uintptr(unsafe.Pointer(envPtr)),
		uintptr(smtoAbortIfHung),
		uintptr(5000),
		uintptr(unsafe.Pointer(&result)),
	)
}
