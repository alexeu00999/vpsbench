//go:build !linux && !darwin

package disk

import (
	"log/slog"
	"os"
)

// dropPageCache is a no-op on unsupported platforms.
func dropPageCache(f *os.File) {
	slog.Debug("[disk.io] dropPageCache not supported on this platform")
}
