//go:build darwin

package darwin

import "os"

func RemoveBinary(exe string) error {
	return os.Remove(exe)
}
