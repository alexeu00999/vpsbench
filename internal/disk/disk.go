package disk

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"
	"os"
	"time"

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

// Run выполняет Disk I/O бенчмарк (заглушка с базовым тестом).
func (d *DiskBench) Run(ctx context.Context) bench.ModuleResult {
	slog.Info("[disk] starting benchmark")
	result := bench.ModuleResult{
		Module: "DISK",
		Info:   "Disk detection pending", // TODO: реальный детект дисков
	}

	// Sequential Write
	slog.Debug("[disk] running sequential write test")
	writeSpeed, err := testSequentialWrite(ctx)
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
	readSpeed, err := testSequentialRead(ctx)
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

	slog.Info("[disk] benchmark completed", "results_count", len(result.Results))
	return result
}

func testSequentialWrite(ctx context.Context) (float64, error) {
	const blockSize = 1024 * 1024 // 1 MB
	const totalSize = 64           // 64 MB total

	data := make([]byte, blockSize)
	rand.Read(data)

	tmpFile, err := os.CreateTemp("", "vpsbench-disk-*")
	if err != nil {
		return 0, err
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	start := time.Now()
	for i := 0; i < totalSize; i++ {
		if ctx.Err() != nil {
			break
		}
		if _, err := tmpFile.Write(data); err != nil {
			return 0, err
		}
	}
	tmpFile.Sync()
	elapsed := time.Since(start).Seconds()

	return float64(totalSize) / elapsed, nil
}

func testSequentialRead(ctx context.Context) (float64, error) {
	const blockSize = 1024 * 1024 // 1 MB
	const totalSize = 64           // 64 MB total

	// Создаём временный файл для чтения
	data := make([]byte, blockSize)
	rand.Read(data)

	tmpFile, err := os.CreateTemp("", "vpsbench-disk-*")
	if err != nil {
		return 0, err
	}
	defer os.Remove(tmpFile.Name())

	for i := 0; i < totalSize; i++ {
		tmpFile.Write(data)
	}
	tmpFile.Sync()
	tmpFile.Seek(0, 0)

	buf := make([]byte, blockSize)
	start := time.Now()
	for i := 0; i < totalSize; i++ {
		if ctx.Err() != nil {
			break
		}
		if _, err := tmpFile.Read(buf); err != nil {
			break
		}
	}
	elapsed := time.Since(start).Seconds()
	tmpFile.Close()

	return float64(totalSize) / elapsed, nil
}
