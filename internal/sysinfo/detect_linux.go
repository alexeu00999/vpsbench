//go:build linux

package sysinfo

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

// detectCPU определяет модель CPU через /proc/cpuinfo.
func detectCPU() string {
	slog.Debug("[sysinfo] detecting CPU model via /proc/cpuinfo")

	data, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		slog.Warn("[sysinfo] failed to read /proc/cpuinfo", "error", err)
		return "Unknown"
	}

	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "model name") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				model := strings.TrimSpace(parts[1])
				slog.Debug("[sysinfo] CPU model detected", "model", model)
				return model
			}
		}
	}

	slog.Warn("[sysinfo] CPU model name not found in /proc/cpuinfo")
	return "Unknown"
}

// detectRAM определяет общий объём RAM через /proc/meminfo.
func detectRAM() uint64 {
	slog.Debug("[sysinfo] detecting RAM via /proc/meminfo")

	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		slog.Warn("[sysinfo] failed to read /proc/meminfo", "error", err)
		return 0
	}

	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "MemTotal:") {
			fields := strings.Fields(line)
			if len(fields) < 2 {
				continue
			}
			kb, err := strconv.ParseUint(fields[1], 10, 64)
			if err != nil {
				slog.Warn("[sysinfo] failed to parse MemTotal", "raw", fields[1], "error", err)
				return 0
			}
			bytes := kb * 1024
			slog.Debug("[sysinfo] RAM detected", "kb", kb, "bytes", bytes, "human", formatRAM(bytes))
			return bytes
		}
	}

	slog.Warn("[sysinfo] MemTotal not found in /proc/meminfo")
	return 0
}

// detectOS определяет дистрибутив и версию ОС.
func detectOS() (osVersion, kernel string) {
	slog.Debug("[sysinfo] detecting OS via /etc/os-release")

	// Дистрибутив из /etc/os-release
	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		slog.Warn("[sysinfo] failed to read /etc/os-release", "error", err)
		osVersion = "Linux"
	} else {
		var name, version string
		for _, line := range strings.Split(string(data), "\n") {
			if strings.HasPrefix(line, "NAME=") {
				name = strings.Trim(strings.TrimPrefix(line, "NAME="), "\"")
			}
			if strings.HasPrefix(line, "VERSION_ID=") {
				version = strings.Trim(strings.TrimPrefix(line, "VERSION_ID="), "\"")
			}
		}
		if name != "" {
			osVersion = name
			if version != "" {
				osVersion = fmt.Sprintf("%s %s", name, version)
			}
		} else {
			osVersion = "Linux"
		}
		slog.Debug("[sysinfo] OS detected", "os_version", osVersion)
	}

	// Ядро через uname
	kernelData, err := os.ReadFile("/proc/version")
	if err != nil {
		slog.Warn("[sysinfo] failed to read /proc/version", "error", err)
		kernel = "Unknown"
	} else {
		fields := strings.Fields(string(kernelData))
		if len(fields) >= 3 {
			kernel = fields[2]
		} else {
			kernel = strings.TrimSpace(string(kernelData))
		}
		slog.Debug("[sysinfo] kernel detected", "kernel", kernel)
	}

	return osVersion, kernel
}

// detectDisks определяет список дисков через /sys/block/.
func detectDisks() []DiskInfo {
	slog.Debug("[sysinfo] detecting disks via /sys/block/")

	entries, err := os.ReadDir("/sys/block")
	if err != nil {
		slog.Warn("[sysinfo] failed to read /sys/block", "error", err)
		return nil
	}

	var disks []DiskInfo
	for _, entry := range entries {
		name := entry.Name()

		// Пропускаем виртуальные устройства (loop, ram, dm)
		if strings.HasPrefix(name, "loop") || strings.HasPrefix(name, "ram") || strings.HasPrefix(name, "dm-") {
			continue
		}

		slog.Debug("[sysinfo] checking block device", "device", name)

		disk := DiskInfo{
			Device: "/dev/" + name,
		}

		// Размер (в секторах по 512 байт)
		sizePath := fmt.Sprintf("/sys/block/%s/size", name)
		sizeData, err := os.ReadFile(sizePath)
		if err == nil {
			sectors, err := strconv.ParseUint(strings.TrimSpace(string(sizeData)), 10, 64)
			if err == nil {
				disk.Size = sectors * 512
			}
		}

		// Пропускаем диски размером 0
		if disk.Size == 0 {
			slog.Debug("[sysinfo] skipping zero-size device", "device", name)
			continue
		}

		// Модель
		modelPath := fmt.Sprintf("/sys/block/%s/device/model", name)
		modelData, err := os.ReadFile(modelPath)
		if err == nil {
			disk.Model = strings.TrimSpace(string(modelData))
		}

		// Тип: NVMe, SSD или HDD
		if strings.HasPrefix(name, "nvme") {
			disk.Type = "NVMe"
		} else {
			rotPath := fmt.Sprintf("/sys/block/%s/queue/rotational", name)
			rotData, err := os.ReadFile(rotPath)
			if err == nil && strings.TrimSpace(string(rotData)) == "0" {
				disk.Type = "SSD"
			} else {
				disk.Type = "HDD"
			}
		}

		slog.Debug("[sysinfo] disk found",
			"device", disk.Device,
			"model", disk.Model,
			"size_bytes", disk.Size,
			"type", disk.Type,
			"human_size", formatBytes(disk.Size),
		)

		disks = append(disks, disk)
	}

	slog.Info("[sysinfo] disks detected", "count", len(disks))
	return disks
}
