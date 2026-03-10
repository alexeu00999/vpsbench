package bench

import (
	"context"
	"errors"
	"testing"
)

// mockBenchmark — мок для тестирования Runner.
type mockBenchmark struct {
	name   string
	result ModuleResult
}

func (m *mockBenchmark) Name() string                          { return m.name }
func (m *mockBenchmark) Run(_ context.Context) ModuleResult { return m.result }

func TestRunnerRunAll(t *testing.T) {
	benchmarks := []Benchmark{
		&mockBenchmark{
			name: "A",
			result: ModuleResult{
				Module:  "A",
				Results: []Result{{Name: "test", Value: 100, Unit: "ops/s"}},
			},
		},
		&mockBenchmark{
			name: "B",
			result: ModuleResult{
				Module:  "B",
				Results: []Result{{Name: "test", Value: 200, Unit: "ops/s"}},
			},
		},
	}

	runner := NewRunner(benchmarks)
	results := runner.RunAll(context.Background())

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	// Проверяем что оба модуля вернули результаты
	modules := map[string]bool{}
	for _, r := range results {
		modules[r.Module] = true
	}
	if !modules["A"] || !modules["B"] {
		t.Errorf("expected modules A and B, got %v", modules)
	}
}

func TestRunnerRunSelected(t *testing.T) {
	benchmarks := []Benchmark{
		&mockBenchmark{name: "CPU", result: ModuleResult{Module: "CPU"}},
		&mockBenchmark{name: "RAM", result: ModuleResult{Module: "RAM"}},
		&mockBenchmark{name: "DISK", result: ModuleResult{Module: "DISK"}},
	}

	runner := NewRunner(benchmarks)
	results := runner.RunSelected(context.Background(), []string{"CPU", "DISK"})

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	for _, r := range results {
		if r.Module != "CPU" && r.Module != "DISK" {
			t.Errorf("unexpected module %s", r.Module)
		}
	}
}

func TestRunnerGracefulDegradation(t *testing.T) {
	benchmarks := []Benchmark{
		&mockBenchmark{
			name: "OK",
			result: ModuleResult{
				Module:  "OK",
				Results: []Result{{Name: "test", Value: 100, Unit: "ops/s"}},
			},
		},
		&mockBenchmark{
			name: "FAIL",
			result: ModuleResult{
				Module: "FAIL",
				Err:    errors.New("benchmark failed"),
			},
		},
	}

	runner := NewRunner(benchmarks)
	results := runner.RunAll(context.Background())

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	// Один должен быть ок, другой с ошибкой
	var okCount, errCount int
	for _, r := range results {
		if r.Err != nil {
			errCount++
		} else {
			okCount++
		}
	}
	if okCount != 1 || errCount != 1 {
		t.Errorf("expected 1 ok and 1 error, got %d ok and %d errors", okCount, errCount)
	}
}

func TestFilterBenchmarks(t *testing.T) {
	benchmarks := []Benchmark{
		&mockBenchmark{name: "CPU"},
		&mockBenchmark{name: "RAM"},
		&mockBenchmark{name: "DISK"},
		&mockBenchmark{name: "NETWORK"},
	}

	tests := []struct {
		name     string
		opts     FilterOptions
		expected int
	}{
		{"No filter", FilterOptions{}, 4},
		{"Only CPU", FilterOptions{CPU: true}, 1},
		{"CPU and RAM", FilterOptions{CPU: true, RAM: true}, 2},
		{"DISK and Network", FilterOptions{Disk: true, Network: true}, 2},
		{"All but Disk", FilterOptions{CPU: true, RAM: true, Network: true}, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			selected := FilterBenchmarks(benchmarks, tt.opts)
			if len(selected) != tt.expected {
				t.Errorf("expected %d selected benchmarks, got %d", tt.expected, len(selected))
			}
			if tt.opts.HasFilter() {
				for _, b := range selected {
					name := b.Name()
					switch name {
					case "CPU":
						if !tt.opts.CPU {
							t.Errorf("unexpected module %s", name)
						}
					case "RAM":
						if !tt.opts.RAM {
							t.Errorf("unexpected module %s", name)
						}
					case "DISK":
						if !tt.opts.Disk {
							t.Errorf("unexpected module %s", name)
						}
					case "NETWORK":
						if !tt.opts.Network {
							t.Errorf("unexpected module %s", name)
						}
					}
				}
			}
		})
	}
}
