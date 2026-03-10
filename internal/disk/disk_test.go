package disk

import (
	"context"
	"testing"
	"time"
)

func TestDiskBenchName(t *testing.T) {
	b := New()
	if b.Name() != "DISK" {
		t.Errorf("Name() = %q, want %q", b.Name(), "DISK")
	}
}

func TestDiskBenchRunDefaultDir(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping disk benchmark in short mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	b := New()
	result := b.Run(ctx)

	if result.Err != nil {
		t.Fatalf("Run() returned error: %v", result.Err)
	}

	if result.Module != "DISK" {
		t.Errorf("Module = %q, want %q", result.Module, "DISK")
	}

	expectedNames := map[string]bool{
		"Seq. Write":    false,
		"Seq. Read":     false,
		"Rand 4K IOPS":  false,
	}

	for _, r := range result.Results {
		if _, ok := expectedNames[r.Name]; ok {
			expectedNames[r.Name] = true
			if r.Value <= 0 {
				t.Errorf("result %q has non-positive value: %f", r.Name, r.Value)
			}
		}
	}

	for name, found := range expectedNames {
		if !found {
			t.Errorf("missing expected result %q", name)
		}
	}
}

func TestDiskBenchRunWithTargets(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping disk benchmark in short mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	tmpDir := t.TempDir()
	b := New(WithTargets([]DiskTarget{
		{Device: "test0", Type: "SSD", Path: tmpDir},
	}))

	result := b.Run(ctx)

	if result.Err != nil {
		t.Fatalf("Run() returned error: %v", result.Err)
	}

	// Single target should not have suffix
	if len(result.Results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(result.Results))
	}
}

func TestDiskBenchContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	b := New()
	result := b.Run(ctx)

	// Should either return an error or have zero/partial results
	// The important thing is it doesn't hang
	_ = result
}
