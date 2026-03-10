package bench

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

// Runner оркестрирует запуск бенчмарков.
type Runner struct {
	benchmarks []Benchmark
}

// NewRunner создаёт Runner с заданным набором бенчмарков.
func NewRunner(benchmarks []Benchmark) *Runner {
	names := make([]string, len(benchmarks))
	for i, b := range benchmarks {
		names[i] = b.Name()
	}
	slog.Debug("[runner] created", "benchmarks", names, "count", len(benchmarks))
	return &Runner{benchmarks: benchmarks}
}

// FilterOptions — опции фильтрации бенчмарков по флагам CLI.
type FilterOptions struct {
	CPU     bool
	RAM     bool
	Disk    bool
	Network bool
}

// HasFilter возвращает true, если установлен хотя бы один фильтр.
func (o FilterOptions) HasFilter() bool {
	return o.CPU || o.RAM || o.Disk || o.Network
}

// FilterBenchmarks фильтрует список бенчмарков согласно опциям.
func FilterBenchmarks(all []Benchmark, opts FilterOptions) []Benchmark {
	if !opts.HasFilter() {
		return all
	}

	var selected []Benchmark
	for _, b := range all {
		name := b.Name()
		switch name {
		case "CPU":
			if opts.CPU {
				selected = append(selected, b)
			}
		case "RAM":
			if opts.RAM {
				selected = append(selected, b)
			}
		case "DISK":
			if opts.Disk {
				selected = append(selected, b)
			}
		case "NETWORK":
			if opts.Network {
				selected = append(selected, b)
			}
		}
	}
	slog.Debug("[bench] filtered benchmarks", "requested", opts, "count", len(selected))
	return selected
}

// RunAll запускает все зарегистрированные бенчмарки параллельно.
func (r *Runner) RunAll(ctx context.Context) []ModuleResult {
	slog.Info("[runner] starting all benchmarks", "count", len(r.benchmarks))
	return r.run(ctx, r.benchmarks)
}

// RunSelected запускает только бенчмарки с указанными именами.
func (r *Runner) RunSelected(ctx context.Context, names []string) []ModuleResult {
	nameSet := make(map[string]bool, len(names))
	for _, n := range names {
		nameSet[n] = true
	}

	var selected []Benchmark
	for _, b := range r.benchmarks {
		if nameSet[b.Name()] {
			selected = append(selected, b)
		}
	}

	slog.Info("[runner] starting selected benchmarks", "requested", names, "matched", len(selected))
	return r.run(ctx, selected)
}

func (r *Runner) run(ctx context.Context, benchmarks []Benchmark) []ModuleResult {
	results := make([]ModuleResult, len(benchmarks))
	var wg sync.WaitGroup

	start := time.Now()

	for i, b := range benchmarks {
		wg.Add(1)
		go func(idx int, bm Benchmark) {
			defer wg.Done()
			name := bm.Name()
			slog.Debug("[runner] launching benchmark", "module", name)

			bmStart := time.Now()
			results[idx] = bm.Run(ctx)
			elapsed := time.Since(bmStart)

			if results[idx].Err != nil {
				slog.Error("[runner] benchmark failed", "module", name, "error", results[idx].Err, "elapsed", elapsed)
			} else {
				slog.Debug("[runner] benchmark completed", "module", name, "elapsed", elapsed)
			}
		}(i, b)
	}

	wg.Wait()
	slog.Info("[runner] all benchmarks finished", "total_elapsed", time.Since(start))
	return results
}
