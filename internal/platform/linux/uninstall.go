//go:build linux

package linux

import "os"

func RemoveBinary(exe string) error {
	return os.Remove(exe)
}
