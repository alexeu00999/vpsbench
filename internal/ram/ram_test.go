package ram

import (
	"context"
	"testing"
	"time"
)

func TestRAMBenchName(t *testing.T) {
	b := New()
	if b.Name() != "RAM" {
		t.Errorf("Name() = %q, want %q", b.Name(), "RAM")
	}
}

func TestRAMBenchRun(t *testing.T) {
	b := New()
	ctx := context.Background()

	result := b.Run(ctx)

	if result.Err != nil {
		t.Fatalf("Run() returned error: %v", result.Err)
	}

	if result.Module != "RAM" {
		t.Errorf("Module = %q, want %q", result.Module, "RAM")
	}

	if len(result.Results) != 2 {
		t.Fatalf("expected 2 results (Write, Read), got %d", len(result.Results))
	}

	write := result.Results[0]
	read := result.Results[1]

	if write.Name != "Write" {
		t.Errorf("result[0].Name = %q, want %q", write.Name, "Write")
	}
	if read.Name != "Read" {
		t.Errorf("result[1].Name = %q, want %q", read.Name, "Read")
	}

	if write.Value <= 0 {
		t.Errorf("Write value should be > 0, got %f", write.Value)
	}
	if read.Value <= 0 {
		t.Errorf("Read value should be > 0, got %f", read.Value)
	}

	if write.Unit != "MB/s" {
		t.Errorf("Write unit = %q, want %q", write.Unit, "MB/s")
	}
	if read.Unit != "MB/s" {
		t.Errorf("Read unit = %q, want %q", read.Unit, "MB/s")
	}

	t.Logf("Write: %.0f MB/s", write.Value)
	t.Logf("Read:  %.0f MB/s", read.Value)
}

func TestRAMBenchRunWithCancel(t *testing.T) {
	b := New()
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(500 * time.Millisecond)
		cancel()
	}()

	result := b.Run(ctx)

	if result.Module != "RAM" {
		t.Errorf("Module = %q, want %q", result.Module, "RAM")
	}
	t.Logf("Cancelled benchmark returned %d results", len(result.Results))
}
