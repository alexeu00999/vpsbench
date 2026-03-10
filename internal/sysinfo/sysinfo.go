package sysinfo

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
)

// DiskInfo содержит информацию об одном диске.
type DiskInfo struct {
	Device string // "/dev/sda", "/dev/nvme0n1"
	Model  string // "Samsung SSD 970 EVO Plus"
	Size   uint64 // байты
	Type   string // "SSD", "HDD", "NVMe"
}

// SystemInfo содержит информацию о системе.
type SystemInfo struct {
	OS        string // "linux", "darwin", "windows"
	Arch      string // "amd64", "arm64"
	Hostname  string
	CPUModel  string
	CPUCores  int
	RAMTotal  uint64 // в байтах
	OSVersion string // "Ubuntu 22.04", "macOS 15.3"
	Kernel    string // "6.1.0-18-amd64"
	Disks     []DiskInfo
	Location  string // "Frankfurt, Germany"
	PublicIP  string // "203.0.113.42"
}

// Detect определяет характеристики системы.
func Detect(ctx context.Context) (SystemInfo, error) {
	slog.Info("[sysinfo] detecting system information")

	hostname, err := os.Hostname()
	if err != nil {
		slog.Error("[sysinfo] failed to get hostname", "error", err)
		hostname = "unknown"
	}

	osVersion, kernel := detectOS()

	info := SystemInfo{
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
		Hostname:  hostname,
		CPUModel:  detectCPU(),
		CPUCores:  runtime.NumCPU(),
		RAMTotal:  detectRAM(),
		OSVersion: osVersion,
		Kernel:    kernel,
		Disks:     detectDisks(),
		Location:  "Unknown",
		PublicIP:  "",
	}

	// GeoIP — с таймаутом, не блокирует остальное
	ip, location := detectLocation(ctx)
	info.PublicIP = ip
	info.Location = location

	slog.Info("[sysinfo] detected",
		"os", info.OSVersion,
		"arch", info.Arch,
		"hostname", info.Hostname,
		"cpu_model", info.CPUModel,
		"cpu_cores", info.CPUCores,
		"ram_total", FormatRAM(info.RAMTotal),
		"disks", len(info.Disks),
		"location", info.Location,
	)
	slog.Debug("[sysinfo] full info",
		"kernel", info.Kernel,
		"public_ip", info.PublicIP,
	)

	return info, nil
}

// FormatRAM форматирует байты в человекочитаемый вид (например, "16 GB").
func FormatRAM(bytes uint64) string {
	return formatRAM(bytes)
}

func formatRAM(bytes uint64) string {
	const (
		gb = 1024 * 1024 * 1024
		mb = 1024 * 1024
	)

	if bytes == 0 {
		return "Unknown"
	}

	if bytes >= gb {
		val := float64(bytes) / float64(gb)
		if val == float64(int(val)) {
			return fmt.Sprintf("%d GB", int(val))
		}
		return fmt.Sprintf("%.1f GB", val)
	}

	return fmt.Sprintf("%d MB", bytes/mb)
}

// formatBytes форматирует байты в человекочитаемый вид для дисков.
func formatBytes(bytes uint64) string {
	const (
		tb = 1024 * 1024 * 1024 * 1024
		gb = 1024 * 1024 * 1024
	)

	if bytes == 0 {
		return "0"
	}

	if bytes >= tb {
		val := float64(bytes) / float64(tb)
		return fmt.Sprintf("%.1f TB", val)
	}

	val := float64(bytes) / float64(gb)
	return fmt.Sprintf("%.0f GB", val)
}

// FormatDisks форматирует список дисков в строку.
func FormatDisks(disks []DiskInfo) string {
	if len(disks) == 0 {
		return "Unknown"
	}

	parts := make([]string, 0, len(disks))
	for _, d := range disks {
		parts = append(parts, fmt.Sprintf("%s %s", d.Type, formatBytes(d.Size)))
	}

	if len(parts) == 1 {
		return parts[0]
	}

	result := ""
	for i, p := range parts {
		if i > 0 {
			result += ", "
		}
		result += p
	}
	return result
}
