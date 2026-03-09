package sysinfo

import (
	"log/slog"
	"os"
	"runtime"
)

// SystemInfo содержит информацию о системе.
type SystemInfo struct {
	OS       string // "linux", "darwin", "windows"
	Arch     string // "amd64", "arm64"
	Hostname string
	CPUModel string // TODO: реальный детект
	CPUCores int
	RAMTotal uint64 // в байтах, TODO: реальный детект
	Location string // TODO: определение через GeoIP
}

// Detect определяет характеристики системы.
func Detect() (SystemInfo, error) {
	slog.Info("[sysinfo] detecting system information")

	hostname, err := os.Hostname()
	if err != nil {
		slog.Error("[sysinfo] failed to get hostname", "error", err)
		hostname = "unknown"
	}

	info := SystemInfo{
		OS:       runtime.GOOS,
		Arch:     runtime.GOARCH,
		Hostname: hostname,
		CPUModel: "Detection pending", // TODO: читать /proc/cpuinfo или sysctl
		CPUCores: runtime.NumCPU(),
		RAMTotal: 0,        // TODO: читать из системы
		Location: "Unknown", // TODO: GeoIP
	}

	slog.Info("[sysinfo] detected",
		"os", info.OS,
		"arch", info.Arch,
		"hostname", info.Hostname,
		"cpu_cores", info.CPUCores,
	)
	slog.Debug("[sysinfo] full info", "cpu_model", info.CPUModel, "ram_total", info.RAMTotal, "location", info.Location)

	return info, nil
}
