//go:build darwin

package disk

import (
	"log/slog"
	"os"
	"syscall"
)

// F_NOCACHE disables caching for the file descriptor on macOS.
const fNocache = 48 // F_NOCACHE

// dropPageCache disables caching for the file via F_NOCACHE (macOS).
func dropPageCache(f *os.File) {
	_, _, errno := syscall.Syscall(syscall.SYS_FCNTL, f.Fd(), uintptr(fNocache), 1)
	if errno != 0 {
		slog.Debug("[disk.io] F_NOCACHE failed", "error", errno)
	} else {
		slog.Debug("[disk.io] disabled cache via F_NOCACHE")
	}
}
