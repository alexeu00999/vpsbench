package disk

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/user/vpsbench/internal/bench"
)

// DiskBench реализует bench.Benchmark для тестирования дискового I/O.
type DiskBench struct {
	targets []DiskTarget // disks to test; empty = use default temp dir
}

// New создаёт новый экземпляр DiskBench с опциональными настройками.
func New(opts ...Option) *DiskBench {
	d := &DiskBench{}
	for _, opt := range opts {
		opt(d)
	}
	return d
}

// Name возвращает имя модуля.
func (d *DiskBench) Name() string {
	return "DISK"
}

// Run выполняет Disk I/O бенчмарк.
// Если заданы targets, тестирует каждый диск отдельно.
// Иначе тестирует дефолтный temp dir.
func (d *DiskBench) Run(ctx context.Context) bench.ModuleResult {
	slog.Info("[disk] starting benchmark", "targets", len(d.targets))
	result := bench.ModuleResult{
		Module: "DISK",
	}

	if len(d.targets) == 0 {
		slog.Debug("[disk] no targets specified, using default temp dir")
		result.Info = "Default"
		return d.runForDir(ctx, "", "", result)
	}

	// Build info string from targets
	infoParts := make([]string, 0, len(d.targets))
	for _, t := range d.targets {
		infoParts = append(infoParts, fmt.Sprintf("%s %s", t.Type, t.Device))
	}
	result.Info = strings.Join(infoParts, ", ")

	for _, target := range d.targets {
		slog.Info("[disk] testing disk", "device", target.Device, "type", target.Type, "path", target.Path)
		suffix := ""
		if len(d.targets) > 1 {
			suffix = " (" + target.Device + ")"
		}
		result = d.runForDir(ctx, target.Path, suffix, result)
		if result.Err != nil {
			return result
		}
	}

	slog.Info("[disk] benchmark completed", "results_count", len(result.Results))
	return result
}

// runForDir runs all disk tests against a specific directory.
// suffix is appended to result names (e.g. " (nvme0n1)") for multi-disk.
func (d *DiskBench) runForDir(ctx context.Context, dir, suffix string, result bench.ModuleResult) bench.ModuleResult {
	slog.Debug("[disk] running tests for dir", "dir", dir, "suffix", suffix)

	// Sequential Write
	writeSpeed, err := testSequentialWrite(ctx, dir)
	if err != nil {
		slog.Error("[disk] sequential write failed", "dir", dir, "error", err)
		result.Err = fmt.Errorf("disk write test%s: %w", suffix, err)
		return result
	}
	result.Results = append(result.Results, bench.Result{
		Name:  "Seq. Write" + suffix,
		Value: writeSpeed,
		Unit:  "MB/s",
	})
	slog.Debug("[disk] seq write result", "dir", dir, "mb_per_sec", writeSpeed)

	// Sequential Read
	readSpeed, err := testSequentialRead(ctx, dir)
	if err != nil {
		slog.Error("[disk] sequential read failed", "dir", dir, "error", err)
		result.Err = fmt.Errorf("disk read test%s: %w", suffix, err)
		return result
	}
	result.Results = append(result.Results, bench.Result{
		Name:  "Seq. Read" + suffix,
		Value: readSpeed,
		Unit:  "MB/s",
	})
	slog.Debug("[disk] seq read result", "dir", dir, "mb_per_sec", readSpeed)

	// Random 4K IOPS
	iops, err := testRandom4KIOPS(ctx, dir)
	if err != nil {
		slog.Error("[disk] random 4K IOPS failed", "dir", dir, "error", err)
		result.Err = fmt.Errorf("disk 4K IOPS test%s: %w", suffix, err)
		return result
	}
	result.Results = append(result.Results, bench.Result{
		Name:  "Rand 4K IOPS" + suffix,
		Value: iops,
		Unit:  "IOPS",
	})
	slog.Debug("[disk] rand 4K IOPS result", "dir", dir, "iops", iops)

	return result
}
