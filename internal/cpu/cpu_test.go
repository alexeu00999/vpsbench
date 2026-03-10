package cpu

import (
	"context"
	"testing"
	"time"
)

func TestCPUBenchName(t *testing.T) {
	b := New()
	if b.Name() != "CPU" {
		t.Errorf("Name() = %q, want %q", b.Name(), "CPU")
	}
}

func TestCPUBenchRun(t *testing.T) {
	b := New()
	ctx := context.Background()

	result := b.Run(ctx)

	if result.Err != nil {
		t.Fatalf("Run() returned error: %v", result.Err)
	}

	if result.Module != "CPU" {
		t.Errorf("Module = %q, want %q", result.Module, "CPU")
	}

	if len(result.Results) != 2 {
		t.Fatalf("expected 2 results (Single-core, Multi-core), got %d", len(result.Results))
	}

	single := result.Results[0]
	multi := result.Results[1]

	if single.Name != "Single-core" {
		t.Errorf("result[0].Name = %q, want %q", single.Name, "Single-core")
	}
	if multi.Name != "Multi-core" {
		t.Errorf("result[1].Name = %q, want %q", multi.Name, "Multi-core")
	}

	if single.Value <= 0 {
		t.Errorf("Single-core value should be > 0, got %f", single.Value)
	}
	if multi.Value <= 0 {
		t.Errorf("Multi-core value should be > 0, got %f", multi.Value)
	}

	if single.Unit != "ops/s" {
		t.Errorf("Single-core unit = %q, want %q", single.Unit, "ops/s")
	}

	// Multi-core должен быть >= single-core (с допуском 80% от single)
	if multi.Value < single.Value*0.8 {
		t.Errorf("Multi-core (%f) should be >= 80%% of Single-core (%f)", multi.Value, single.Value)
	}

	t.Logf("Single-core: %.0f ops/s", single.Value)
	t.Logf("Multi-core:  %.0f ops/s", multi.Value)
	t.Logf("Scaling:     %.1fx", multi.Value/single.Value)
}

func TestCPUBenchRunWithCancel(t *testing.T) {
	b := New()
	ctx, cancel := context.WithCancel(context.Background())

	// Отменяем контекст через 500ms — бенчмарк должен завершиться graceful
	go func() {
		time.Sleep(500 * time.Millisecond)
		cancel()
	}()

	result := b.Run(ctx)

	// Не должно быть паники, бенчмарк должен вернуть результат
	if result.Module != "CPU" {
		t.Errorf("Module = %q, want %q", result.Module, "CPU")
	}

	t.Logf("Cancelled benchmark returned %d results", len(result.Results))
}

func TestRunOpsTest(t *testing.T) {
	ctx := context.Background()
	duration := 500 * time.Millisecond

	single := runOpsTest(ctx, 1, duration)
	multi := runOpsTest(ctx, 4, duration)

	if single <= 0 {
		t.Errorf("single goroutine ops/s should be > 0, got %f", single)
	}
	if multi <= 0 {
		t.Errorf("multi goroutine ops/s should be > 0, got %f", multi)
	}

	// 4 горутины должны дать больше ops/s чем 1 (с допуском 1.5x)
	if multi < single*1.5 {
		t.Errorf("4 goroutines (%f) should give >= 1.5x of 1 goroutine (%f)", multi, single)
	}

	t.Logf("1 goroutine:  %.0f ops/s", single)
	t.Logf("4 goroutines: %.0f ops/s (%.1fx)", multi, multi/single)
}
