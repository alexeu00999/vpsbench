//go:build darwin

package sysinfo

import (
	"fmt"
	"log/slog"
	"os/exec"
	"strconv"
	"strings"
)

// detectCPU определяет модель CPU через sysctl.
func detectCPU() string {
	slog.Debug("[sysinfo] detecting CPU model via sysctl")

	out, err := exec.Command("sysctl", "-n", "machdep.cpu.brand_string").Output()
	if err != nil {
		slog.Warn("[sysinfo] sysctl machdep.cpu.brand_string failed", "error", err)
		return "Unknown"
	}

	model := strings.TrimSpace(string(out))
	slog.Debug("[sysinfo] CPU model detected", "model", model)
	return model
}

// detectRAM определяет общий объём RAM через sysctl hw.memsize.
func detectRAM() uint64 {
	slog.Debug("[sysinfo] detecting RAM via sysctl hw.memsize")

	out, err := exec.Command("sysctl", "-n", "hw.memsize").Output()
	if err != nil {
		slog.Warn("[sysinfo] sysctl hw.memsize failed", "error", err)
		return 0
	}

	bytes, err := strconv.ParseUint(strings.TrimSpace(string(out)), 10, 64)
	if err != nil {
		slog.Warn("[sysinfo] failed to parse hw.memsize", "raw", string(out), "error", err)
		return 0
	}

	slog.Debug("[sysinfo] RAM detected", "bytes", bytes, "human", formatRAM(bytes))
	return bytes
}

// detectOS определяет версию macOS через sw_vers и uname.
func detectOS() (osVersion, kernel string) {
	slog.Debug("[sysinfo] detecting OS via sw_vers")

	// macOS версия
	nameOut, err := exec.Command("sw_vers", "-productName").Output()
	if err != nil {
		slog.Warn("[sysinfo] sw_vers -productName failed", "error", err)
		osVersion = "macOS"
	} else {
		name := strings.TrimSpace(string(nameOut))
		verOut, err := exec.Command("sw_vers", "-productVersion").Output()
		if err != nil {
			osVersion = name
		} else {
			osVersion = fmt.Sprintf("%s %s", name, strings.TrimSpace(string(verOut)))
		}
	}
	slog.Debug("[sysinfo] OS detected", "os_version", osVersion)

	// Ядро через uname -r
	kernelOut, err := exec.Command("uname", "-r").Output()
	if err != nil {
		slog.Warn("[sysinfo] uname -r failed", "error", err)
		kernel = "Unknown"
	} else {
		kernel = strings.TrimSpace(string(kernelOut))
	}
	slog.Debug("[sysinfo] kernel detected", "kernel", kernel)

	return osVersion, kernel
}

// detectDisks определяет список дисков через diskutil.
func detectDisks() []DiskInfo {
	slog.Debug("[sysinfo] detecting disks via diskutil")

	return detectDisksFallback()
}

// detectDisksFallback определяет диски через diskutil list (текстовый формат).
func detectDisksFallback() []DiskInfo {
	out, err := exec.Command("diskutil", "list").Output()
	if err != nil {
		slog.Warn("[sysinfo] diskutil list failed", "error", err)
		return nil
	}

	var physicalDisks []string
	for _, line := range strings.Split(string(out), "\n") {
		// Ищем строки вида "/dev/disk0 (internal, physical):"
		if strings.HasPrefix(line, "/dev/disk") && strings.Contains(line, "physical") {
			parts := strings.Fields(line)
			if len(parts) > 0 {
				diskName := strings.TrimSuffix(parts[0], ":")
				physicalDisks = append(physicalDisks, diskName)
			}
		}
	}

	slog.Debug("[sysinfo] found physical disks", "disks", physicalDisks)

	var disks []DiskInfo
	for _, diskPath := range physicalDisks {
		info := getDiskInfo(diskPath)
		if info.Size > 0 {
			disks = append(disks, info)
		}
	}

	slog.Info("[sysinfo] disks detected", "count", len(disks))
	return disks
}

// getDiskInfo получает информацию об одном диске через diskutil info.
func getDiskInfo(diskPath string) DiskInfo {
	out, err := exec.Command("diskutil", "info", diskPath).Output()
	if err != nil {
		slog.Warn("[sysinfo] diskutil info failed", "disk", diskPath, "error", err)
		return DiskInfo{Device: diskPath}
	}

	disk := DiskInfo{Device: diskPath}
	lines := strings.Split(string(out), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "Device / Media Name:") {
			disk.Model = strings.TrimSpace(strings.TrimPrefix(line, "Device / Media Name:"))
		}
		if strings.HasPrefix(line, "Disk Size:") {
			// "Disk Size: 500.1 GB (500107862016 Bytes)..."
			if idx := strings.Index(line, "("); idx > 0 {
				if endIdx := strings.Index(line[idx:], " Bytes"); endIdx > 0 {
					bytesStr := line[idx+1 : idx+endIdx]
					if b, err := strconv.ParseUint(bytesStr, 10, 64); err == nil {
						disk.Size = b
					}
				}
			}
		}
		if strings.HasPrefix(line, "Solid State:") {
			if strings.Contains(line, "Yes") {
				disk.Type = "SSD"
			} else {
				disk.Type = "HDD"
			}
		}
		if strings.HasPrefix(line, "Protocol:") {
			if strings.Contains(line, "NVMe") {
				disk.Type = "NVMe"
			}
		}
	}

	if disk.Type == "" {
		disk.Type = "SSD" // macOS обычно SSD
	}

	slog.Debug("[sysinfo] disk info",
		"device", disk.Device,
		"model", disk.Model,
		"size_bytes", disk.Size,
		"type", disk.Type,
		"human_size", formatBytes(disk.Size),
	)

	return disk
}
