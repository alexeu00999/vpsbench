package cpu

import (
	"context"
	"log/slog"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/user/vpsbench/internal/bench"
)

const (
	// testDuration — длительность каждого теста (single/multi).
	testDuration = 3 * time.Second

	// warmupDuration — длительность прогрева кэшей перед замером.
	warmupDuration = 100 * time.Millisecond
)

// CPUBench реализует bench.Benchmark для тестирования CPU.
type CPUBench struct{}

// New создаёт новый экземпляр CPUBench.
func New() *CPUBench {
	return &CPUBench{}
}

// Name возвращает имя модуля.
func (c *CPUBench) Name() string {
	return "CPU"
}

// Run выполняет CPU бенчмарк: warmup, single-core, multi-core.
func (c *CPUBench) Run(ctx context.Context) bench.ModuleResult {
	slog.Info("[cpu] starting benchmark")
	start := time.Now()

	result := bench.ModuleResult{
		Module: "CPU",
	}

	// Warmup — прогреваем кэши CPU
	slog.Debug("[cpu] running warmup", "duration", warmupDuration)
	warmupOps := runOpsTest(ctx, 1, warmupDuration)
	slog.Debug("[cpu] warmup complete", "ops", warmupOps)

	// Single-core тест
	slog.Debug("[cpu] running single-core test", "duration", testDuration)
	singleStart := time.Now()
	single := runOpsTest(ctx, 1, testDuration)
	singleElapsed := time.Since(singleStart)
	result.Results = append(result.Results, bench.Result{
		Name:  "Single-core",
		Value: single,
		Unit:  "ops/s",
	})
	slog.Debug("[cpu] single-core result", "ops_per_sec", single, "elapsed", singleElapsed)

	// Multi-core тест
	cores := runtime.NumCPU()
	slog.Debug("[cpu] running multi-core test", "cores", cores, "duration", testDuration)
	multiStart := time.Now()
	multi := runOpsTest(ctx, cores, testDuration)
	multiElapsed := time.Since(multiStart)
	result.Results = append(result.Results, bench.Result{
		Name:  "Multi-core",
		Value: multi,
		Unit:  "ops/s",
	})
	slog.Debug("[cpu] multi-core result", "ops_per_sec", multi, "cores", cores, "elapsed", multiElapsed)

	totalElapsed := time.Since(start)
	slog.Info("[cpu] benchmark completed",
		"single_core", single,
		"multi_core", multi,
		"cores", cores,
		"total_elapsed", totalElapsed,
	)
	return result
}

// runOpsTest запускает смешанную вычислительную нагрузку на goroutines горутинах
// в течение duration. Возвращает суммарное количество операций в секунду.
func runOpsTest(ctx context.Context, goroutines int, duration time.Duration) float64 {
	slog.Debug("[cpu] runOpsTest starting", "goroutines", goroutines, "duration", duration)

	var totalOps atomic.Int64
	var wg sync.WaitGroup

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			slog.Debug("[cpu] goroutine started", "id", id)
			runMixedWorkload(ctx, duration, &totalOps)
			slog.Debug("[cpu] goroutine finished", "id", id)
		}(i)
	}

	wg.Wait()

	ops := totalOps.Load()
	opsPerSec := float64(ops) / duration.Seconds()

	slog.Debug("[cpu] runOpsTest completed",
		"goroutines", goroutines,
		"total_ops", ops,
		"ops_per_sec", opsPerSec,
	)

	return opsPerSec
}
