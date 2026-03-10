//go:build linux

package disk

import (
	"log/slog"
	"os"

	"golang.org/x/sys/unix"
)

// dropPageCache advises the OS to drop page cache for the file (Linux).
func dropPageCache(f *os.File) {
	err := unix.Fadvise(int(f.Fd()), 0, 0, unix.FADV_DONTNEED)
	if err != nil {
		slog.Debug("[disk.io] fadvise FADV_DONTNEED failed", "error", err)
	} else {
		slog.Debug("[disk.io] dropped page cache via fadvise")
	}
}
