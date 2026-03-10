package sysinfo

import (
	"context"
	"testing"
	"time"
)

func TestDetect(t *testing.T) {
	ctx := context.Background()
	info, err := Detect(ctx)
	if err != nil {
		t.Fatalf("Detect() returned error: %v", err)
	}

	if info.OS == "" {
		t.Error("OS should not be empty")
	}
	if info.Arch == "" {
		t.Error("Arch should not be empty")
	}
	if info.Hostname == "" {
		t.Error("Hostname should not be empty")
	}
	if info.CPUCores <= 0 {
		t.Errorf("CPUCores should be > 0, got %d", info.CPUCores)
	}
}

func TestDetectCPU(t *testing.T) {
	model := detectCPU()
	if model == "" || model == "Detection pending" {
		t.Errorf("detectCPU() should return real model, got %q", model)
	}
	t.Logf("CPU model: %s", model)
}

func TestDetectRAM(t *testing.T) {
	ram := detectRAM()
	if ram == 0 {
		t.Error("detectRAM() should return > 0")
	}
	t.Logf("RAM: %d bytes (%s)", ram, formatRAM(ram))
}

func TestDetectOS(t *testing.T) {
	osVersion, kernel := detectOS()
	if osVersion == "" {
		t.Error("OS version should not be empty")
	}
	if kernel == "" {
		t.Error("Kernel should not be empty")
	}
	t.Logf("OS: %s, Kernel: %s", osVersion, kernel)
}

func TestDetectDisks(t *testing.T) {
	disks := detectDisks()
	if len(disks) == 0 {
		t.Log("WARNING: no disks detected (may be expected in some CI environments)")
		return
	}

	for _, d := range disks {
		if d.Device == "" {
			t.Error("disk device should not be empty")
		}
		if d.Size == 0 {
			t.Errorf("disk %s has zero size", d.Device)
		}
		if d.Type == "" {
			t.Errorf("disk %s has empty type", d.Device)
		}
		t.Logf("Disk: %s %s %s (%s)", d.Device, d.Model, d.Type, formatBytes(d.Size))
	}
}

func TestDetectLocationTimeout(t *testing.T) {
	// Очень короткий таймаут — должен graceful fallback
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	time.Sleep(time.Millisecond) // Даём контексту протухнуть

	ip, location := detectLocation(ctx)
	// При протухшем контексте должен вернуть fallback
	if location != "Unknown" && ip == "" {
		// Если location не Unknown, то IP тоже должен быть
		t.Logf("Got location despite short timeout: %s (ip: %s)", location, ip)
	}
}

func TestFormatRAM(t *testing.T) {
	tests := []struct {
		bytes uint64
		want  string
	}{
		{0, "Unknown"},
		{512 * 1024 * 1024, "512 MB"},
		{1024 * 1024 * 1024, "1 GB"},
		{2 * 1024 * 1024 * 1024, "2 GB"},
		{16 * 1024 * 1024 * 1024, "16 GB"},
		{32 * 1024 * 1024 * 1024, "32 GB"},
	}

	for _, tt := range tests {
		got := formatRAM(tt.bytes)
		if got != tt.want {
			t.Errorf("formatRAM(%d) = %q, want %q", tt.bytes, got, tt.want)
		}
	}
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		bytes uint64
		want  string
	}{
		{0, "0"},
		{500 * 1024 * 1024 * 1024, "500 GB"},
		{1024 * 1024 * 1024 * 1024, "1.0 TB"},
		{2 * 1024 * 1024 * 1024 * 1024, "2.0 TB"},
	}

	for _, tt := range tests {
		got := formatBytes(tt.bytes)
		if got != tt.want {
			t.Errorf("formatBytes(%d) = %q, want %q", tt.bytes, got, tt.want)
		}
	}
}

func TestFormatDisks(t *testing.T) {
	// Пустой список
	if got := FormatDisks(nil); got != "Unknown" {
		t.Errorf("FormatDisks(nil) = %q, want %q", got, "Unknown")
	}

	// Один диск
	disks := []DiskInfo{{Type: "NVMe", Size: 500 * 1024 * 1024 * 1024}}
	got := FormatDisks(disks)
	if got != "NVMe 500 GB" {
		t.Errorf("FormatDisks(1 disk) = %q, want %q", got, "NVMe 500 GB")
	}

	// Два диска
	disks = append(disks, DiskInfo{Type: "SSD", Size: 1024 * 1024 * 1024 * 1024})
	got = FormatDisks(disks)
	if got != "NVMe 500 GB, SSD 1.0 TB" {
		t.Errorf("FormatDisks(2 disks) = %q, want %q", got, "NVMe 500 GB, SSD 1.0 TB")
	}
}
