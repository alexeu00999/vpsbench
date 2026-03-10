package disk

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/user/vpsbench/internal/bench"
)

// DiskBench реализует bench.Benchmark для тестирования дискового I/O.
type DiskBench struct{}

// New создаёт новый экземпляр DiskBench.
func New() *DiskBench {
	return &DiskBench{}
}

// Name возвращает имя модуля.
func (d *DiskBench) Name() string {
	return "DISK"
}

// Run выполняет Disk I/O бенчмарк.
func (d *DiskBench) Run(ctx context.Context) bench.ModuleResult {
	slog.Info("[disk] starting benchmark")
	result := bench.ModuleResult{
		Module: "DISK",
		Info:   "Disk detection pending",
	}

	// Sequential Write
	slog.Debug("[disk] running sequential write test")
	writeSpeed, err := testSequentialWrite(ctx, "")
	if err != nil {
		slog.Error("[disk] sequential write failed", "error", err)
		result.Err = fmt.Errorf("disk write test: %w", err)
		return result
	}
	result.Results = append(result.Results, bench.Result{
		Name:  "Seq. Write",
		Value: writeSpeed,
		Unit:  "MB/s",
	})
	slog.Debug("[disk] seq write result", "mb_per_sec", writeSpeed)

	// Sequential Read
	slog.Debug("[disk] running sequential read test")
	readSpeed, err := testSequentialRead(ctx, "")
	if err != nil {
		slog.Error("[disk] sequential read failed", "error", err)
		result.Err = fmt.Errorf("disk read test: %w", err)
		return result
	}
	result.Results = append(result.Results, bench.Result{
		Name:  "Seq. Read",
		Value: readSpeed,
		Unit:  "MB/s",
	})
	slog.Debug("[disk] seq read result", "mb_per_sec", readSpeed)

	// Random 4K IOPS
	slog.Debug("[disk] running random 4K IOPS test")
	iops, err := testRandom4KIOPS(ctx, "")
	if err != nil {
		slog.Error("[disk] random 4K IOPS failed", "error", err)
		result.Err = fmt.Errorf("disk 4K IOPS test: %w", err)
		return result
	}
	result.Results = append(result.Results, bench.Result{
		Name:  "Rand 4K IOPS",
		Value: iops,
		Unit:  "IOPS",
	})
	slog.Debug("[disk] rand 4K IOPS result", "iops", iops)

	slog.Info("[disk] benchmark completed", "results_count", len(result.Results))
	return result
}
