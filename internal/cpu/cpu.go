package cpu

import (
	"context"
	"log/slog"
	"math"
	"runtime"
	"sync"
	"time"

	"github.com/user/vpsbench/internal/bench"
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

// Run выполняет CPU бенчмарк (заглушка с базовыми вычислениями).
func (c *CPUBench) Run(ctx context.Context) bench.ModuleResult {
	slog.Info("[cpu] starting benchmark")
	result := bench.ModuleResult{
		Module: "CPU",
		Info:   "CPU detection pending", // TODO: реальный детект в sysinfo
	}

	// Single-core тест
	slog.Debug("[cpu] running single-core test")
	single := runOpsTest(ctx, 1)
	result.Results = append(result.Results, bench.Result{
		Name:  "Single-core",
		Value: single,
		Unit:  "ops/s",
	})
	slog.Debug("[cpu] single-core result", "ops_per_sec", single)

	// Multi-core тест
	cores := runtime.NumCPU()
	slog.Debug("[cpu] running multi-core test", "cores", cores)
	multi := runOpsTest(ctx, cores)
	result.Results = append(result.Results, bench.Result{
		Name:  "Multi-core",
		Value: multi,
		Unit:  "ops/s",
	})
	slog.Debug("[cpu] multi-core result", "ops_per_sec", multi, "cores", cores)

	slog.Info("[cpu] benchmark completed", "results_count", len(result.Results))
	return result
}

// runOpsTest запускает вычислительный тест на указанном количестве горутин.
// Возвращает суммарное количество операций в секунду.
func runOpsTest(ctx context.Context, goroutines int) float64 {
	duration := 2 * time.Second
	var totalOps int64
	var mu sync.Mutex

	var wg sync.WaitGroup
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ops := int64(0)
			deadline := time.Now().Add(duration)
			for time.Now().Before(deadline) {
				// Вычислительная нагрузка
				for j := 0; j < 1000; j++ {
					_ = math.Sqrt(float64(j) * 2.71828)
				}
				ops += 1000
				if ctx.Err() != nil {
					break
				}
			}
			mu.Lock()
			totalOps += ops
			mu.Unlock()
		}()
	}
	wg.Wait()

	return float64(totalOps) / duration.Seconds()
}
