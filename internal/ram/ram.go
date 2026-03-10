package ram

import (
	"context"
	"log/slog"
	"time"

	"github.com/user/vpsbench/internal/bench"
)

const (
	// testDuration — длительность каждого теста (write/read).
	testDuration = 3 * time.Second
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

// Run выполняет RAM бенчмарк: sequential write и sequential read.
func (r *RAMBench) Run(ctx context.Context) bench.ModuleResult {
	slog.Info("[ram] starting benchmark")
	start := time.Now()

	result := bench.ModuleResult{
		Module: "RAM",
	}

	// Sequential Write
	slog.Debug("[ram] running write test", "duration", testDuration)
	writeStart := time.Now()
	writeMBs := workloadWrite(ctx, testDuration)
	writeElapsed := time.Since(writeStart)
	result.Results = append(result.Results, bench.Result{
		Name:  "Write",
		Value: writeMBs,
		Unit:  "MB/s",
	})
	slog.Debug("[ram] write result", "mb_per_sec", writeMBs, "elapsed", writeElapsed)

	// Sequential Read
	slog.Debug("[ram] running read test", "duration", testDuration)
	readStart := time.Now()
	readMBs := workloadRead(ctx, testDuration)
	readElapsed := time.Since(readStart)
	result.Results = append(result.Results, bench.Result{
		Name:  "Read",
		Value: readMBs,
		Unit:  "MB/s",
	})
	slog.Debug("[ram] read result", "mb_per_sec", readMBs, "elapsed", readElapsed)

	totalElapsed := time.Since(start)
	slog.Info("[ram] benchmark completed",
		"write_mb_s", writeMBs,
		"read_mb_s", readMBs,
		"total_elapsed", totalElapsed,
	)

	return result
}
