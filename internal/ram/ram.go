package ram

import (
	"context"
	"log/slog"
	"time"
	"unsafe"

	"github.com/user/vpsbench/internal/bench"
)

// RAMBench реализует bench.Benchmark для тестирования оперативной памяти.
type RAMBench struct{}

// New создаёт новый экземпляр RAMBench.
func New() *RAMBench {
	return &RAMBench{}
}

// Name возвращает имя модуля.
func (r *RAMBench) Name() string {
	return "RAM"
}

// Run выполняет RAM бенчмарк (заглушка с базовым тестом скорости).
func (r *RAMBench) Run(ctx context.Context) bench.ModuleResult {
	slog.Info("[ram] starting benchmark")
	result := bench.ModuleResult{
		Module: "RAM",
		Info:   "RAM detection pending", // TODO: реальный детект в sysinfo
	}

	speed := measureSpeed(ctx)
	result.Results = append(result.Results, bench.Result{
		Name:  "Speed (R/W)",
		Value: speed,
		Unit:  "MB/s",
	})
	slog.Debug("[ram] speed result", "mb_per_sec", speed)

	slog.Info("[ram] benchmark completed")
	return result
}

// measureSpeed измеряет скорость чтения/записи памяти в MB/s.
func measureSpeed(ctx context.Context) float64 {
	const blockSize = 64 * 1024 * 1024 // 64 MB
	buf := make([]byte, blockSize)

	// Запись
	start := time.Now()
	iterations := 0
	deadline := start.Add(2 * time.Second)

	for time.Now().Before(deadline) {
		for i := 0; i < blockSize; i += int(unsafe.Sizeof(uint64(0))) {
			buf[i] = byte(i)
		}
		iterations++
		if ctx.Err() != nil {
			break
		}
	}

	elapsed := time.Since(start).Seconds()
	totalBytes := float64(iterations) * float64(blockSize)
	mbPerSec := totalBytes / elapsed / 1024 / 1024

	slog.Debug("[ram] measurement done", "iterations", iterations, "elapsed_sec", elapsed)
	return mbPerSec
}
